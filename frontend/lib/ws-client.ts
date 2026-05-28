/**
 * WebSocket client wrapper — v.1 runtime.
 *
 * Responsibilities:
 *   - Connect, reconnect with exponential backoff
 *   - Send first-message auth on every open; gate everything else until auth_ok
 *   - Queue outbound messages while disconnected or unauthenticated
 *   - Validate every inbound message against ServerMessageSchema
 *   - Push validated messages to a sink (the Zustand store)
 *   - Maintain a structured event log
 *
 * Per Mike's v.1 spec: no replay, no server-side event buffering, no ack
 * tracking. Reconnect drops in-flight relations; the universe snapshot is
 * the safety net.
 */

import { ServerMessageSchema, type ServerMessageParsed } from "./protocol/schemas";
import type { ClientMessage } from "./protocol/types";

export type ConnectionStatus =
  | "idle"
  | "connecting"
  | "authenticating"
  | "open"
  | "closed"
  | "error"
  | "auth_failed";

export interface WSEvent {
  ts: number;
  kind: "status" | "in" | "out" | "drop" | "error";
  detail: unknown;
}

interface WSClientOptions {
  url: string;
  /** Auth token sent in the first message after open. */
  token: string;
  onMessage: (msg: ServerMessageParsed) => void;
  onStatus?: (status: ConnectionStatus) => void;
  onEvent?: (event: WSEvent) => void;
  /** Initial backoff in ms; doubles each retry up to maxBackoffMs. */
  baseBackoffMs?: number;
  maxBackoffMs?: number;
}

const newId = () =>
  `c_${Date.now().toString(36)}_${Math.random().toString(36).slice(2, 8)}`;

export class WSClient {
  private socket: WebSocket | null = null;
  private status: ConnectionStatus = "idle";
  /** Outbox holds non-auth messages queued while not yet authenticated. */
  private outbox: ClientMessage[] = [];
  private retry = 0;
  private reconnectTimer: ReturnType<typeof setTimeout> | null = null;
  private destroyed = false;
  private readonly opts: Required<Omit<WSClientOptions, "onStatus" | "onEvent">> &
    Pick<WSClientOptions, "onStatus" | "onEvent">;

  constructor(opts: WSClientOptions) {
    this.opts = {
      baseBackoffMs: 500,
      maxBackoffMs: 15_000,
      onStatus: undefined,
      onEvent: undefined,
      ...opts,
    };
  }

  connect() {
    if (this.destroyed) return;
    if (this.socket && (this.status === "open" || this.status === "connecting" || this.status === "authenticating")) {
      return;
    }

    this.setStatus("connecting");
    try {
      this.socket = new WebSocket(this.opts.url);
    } catch (err) {
      this.emit({ ts: Date.now(), kind: "error", detail: err });
      this.scheduleReconnect();
      return;
    }

    this.socket.onopen = () => {
      this.setStatus("authenticating");
      this.sendAuth();
    };

    this.socket.onmessage = (e) => {
      let raw: unknown;
      try {
        raw = JSON.parse(typeof e.data === "string" ? e.data : "");
      } catch (err) {
        this.emit({ ts: Date.now(), kind: "drop", detail: { reason: "invalid_json", err } });
        return;
      }

      const parsed = ServerMessageSchema.safeParse(raw);
      if (!parsed.success) {
        this.emit({
          ts: Date.now(),
          kind: "drop",
          detail: { reason: "schema_invalid", issues: parsed.error.issues, raw },
        });
        return;
      }

      this.emit({ ts: Date.now(), kind: "in", detail: parsed.data });

      // Handshake handling — promote to "open" on auth_ok, terminate on auth_fail.
      if (parsed.data.type === "auth_ok") {
        this.retry = 0;
        this.setStatus("open");
        this.flushOutbox();
        return;
      }
      if (parsed.data.type === "auth_fail") {
        this.setStatus("auth_failed");
        this.opts.onMessage(parsed.data); // surface to the store so UI can show the reason
        this.socket?.close();
        return;
      }

      // All other messages only reach the store after auth_ok.
      if (this.status !== "open") {
        this.emit({
          ts: Date.now(),
          kind: "drop",
          detail: { reason: "pre_auth", type: parsed.data.type },
        });
        return;
      }

      this.opts.onMessage(parsed.data);
    };

    this.socket.onclose = () => {
      // Don't reconnect on a deliberate auth failure.
      if (this.status === "auth_failed") return;
      this.setStatus("closed");
      this.scheduleReconnect();
    };

    this.socket.onerror = (err) => {
      this.emit({ ts: Date.now(), kind: "error", detail: err });
      // Don't override auth_failed; let onclose handle reconnect for other errors.
      if (this.status !== "auth_failed") this.setStatus("error");
    };
  }

  /**
   * Queue an outbound message. Sent immediately if authenticated, otherwise
   * held and flushed after auth_ok.
   */
  send(msg: ClientMessage) {
    if (this.status === "open" && this.socket) {
      this.socket.send(JSON.stringify(msg));
      this.emit({ ts: Date.now(), kind: "out", detail: msg });
    } else {
      this.outbox.push(msg);
      this.emit({
        ts: Date.now(),
        kind: "out",
        detail: { queued: true, msg, status: this.status },
      });
    }
  }

  destroy() {
    this.destroyed = true;
    if (this.reconnectTimer) clearTimeout(this.reconnectTimer);
    this.socket?.close();
    this.socket = null;
  }

  getStatus() {
    return this.status;
  }

  private sendAuth() {
    if (!this.socket) return;
    const msg: ClientMessage = {
      id: newId(),
      ts: Date.now(),
      type: "auth",
      payload: { token: this.opts.token },
    };
    this.socket.send(JSON.stringify(msg));
    this.emit({ ts: Date.now(), kind: "out", detail: msg });
  }

  private flushOutbox() {
    while (this.outbox.length && this.socket && this.status === "open") {
      const msg = this.outbox.shift()!;
      this.socket.send(JSON.stringify(msg));
      this.emit({ ts: Date.now(), kind: "out", detail: { flushed: true, msg } });
    }
  }

  private scheduleReconnect() {
    if (this.destroyed) return;
    if (this.reconnectTimer) return;
    const backoff = Math.min(
      this.opts.baseBackoffMs * 2 ** this.retry,
      this.opts.maxBackoffMs,
    );
    this.retry += 1;
    this.reconnectTimer = setTimeout(() => {
      this.reconnectTimer = null;
      this.connect();
    }, backoff);
  }

  private setStatus(next: ConnectionStatus) {
    if (this.status === next) return;
    this.status = next;
    this.emit({ ts: Date.now(), kind: "status", detail: next });
    this.opts.onStatus?.(next);
  }

  private emit(event: WSEvent) {
    this.opts.onEvent?.(event);
  }
}

export const DEFAULT_WS_URL =
  (typeof process !== "undefined" && process.env.NEXT_PUBLIC_WS_URL) ||
  "ws://localhost:8080";

export const DEFAULT_AUTH_TOKEN =
  (typeof process !== "undefined" && process.env.NEXT_PUBLIC_WS_TOKEN) ||
  "dev-token";
