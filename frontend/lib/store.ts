/**
 * Zustand store — v.1 universe state.
 *
 * Holds the most recent universe snapshot plus the deltas applied since.
 * Components derive views via selectors; the store never computes derived
 * state of its own. If you find yourself adding a derived field here, ask
 * whether the runtime should own that derivation.
 *
 * Per Mike's v.1: snapshot on connect, snapshot on reconnect, deltas in
 * between. No event replay, no client-side ack tracking.
 */

import { create } from "zustand";
import type {
  Being,
  Economics,
  Exchange,
  ExchangeEntry,
  Thread,
  TopologyNode,
  UniverseState,
} from "./protocol/types";
import type { ServerMessageParsed } from "./protocol/schemas";
import type { ConnectionStatus, WSClient, WSEvent } from "./ws-client";

interface StoreState {
  // ----- Universe state (rebuilt from snapshot, mutated by deltas) ----------
  beings: Record<string, Being>;
  threads: Record<string, Thread>;
  exchanges: Record<string, Exchange>;
  economics: Economics;
  topology: TopologyNode;

  // ----- View state (frontend-only, not from runtime) -----------------------
  /** Which top-level view is active. */
  view: "user" | "universe";
  /** The being whose perspective User-present is taken from. */
  perspectiveBeing: string;
  /** Selected being for the detail panel. */
  selectedBeing: string | null;

  // ----- Connection state ---------------------------------------------------
  status: ConnectionStatus;
  /** Most recent server-pushed error, or auth failure reason. */
  lastError: { code: string; message: string; ts: number } | null;
  events: WSEvent[];
  client: WSClient | null;

  // ----- Actions ------------------------------------------------------------
  setClient: (client: WSClient) => void;
  ingest: (msg: ServerMessageParsed) => void;
  recordEvent: (e: WSEvent) => void;
  setStatus: (s: ConnectionStatus) => void;
  setView: (v: "user" | "universe") => void;
  setPerspectiveBeing: (name: string) => void;
  selectBeing: (name: string | null) => void;
  dismissError: () => void;
}

const EVENTS_WINDOW = 200;

const EMPTY_TOPOLOGY: TopologyNode = { id: "universe", type: "Universe", children: [] };
const EMPTY_ECONOMICS: Economics = { fields: {} };

export const useAppStore = create<StoreState>((set) => ({
  beings: {},
  threads: {},
  exchanges: {},
  economics: EMPTY_ECONOMICS,
  topology: EMPTY_TOPOLOGY,

  view: "user",
  perspectiveBeing: "michael",
  selectedBeing: null,

  status: "idle",
  lastError: null,
  events: [],
  client: null,

  setClient: (client) => set({ client }),

  ingest: (msg) =>
    set((state) => {
      switch (msg.type) {
        case "universe": {
          const u: UniverseState = msg.payload as UniverseState;
          return {
            beings: indexBy(u.beings, (b) => b.name),
            threads: indexBy(u.threads, (t) => t.id),
            exchanges: indexBy(u.exchanges, (e) => e.key),
            economics: u.economics,
            topology: u.topology,
          };
        }

        case "entry": {
          const d = msg.payload;
          const ex = state.exchanges[d.exchange];
          if (!ex) return {}; // unknown exchange — drop until next snapshot
          const entry: ExchangeEntry = {
            index: d.index,
            from: d.from,
            content: d.content,
            ts: msg.ts,
          };
          // Idempotent append: ignore if we've already seen this index.
          if (ex.entries.some((e) => e.index === entry.index)) return {};
          return {
            exchanges: {
              ...state.exchanges,
              [d.exchange]: { ...ex, entries: [...ex.entries, entry] },
            },
          };
        }

        case "edge": {
          const d = msg.payload;
          const t = state.threads[d.thread_id];
          if (!t) return {};
          if (t.edges.some((e) => e.from === d.from && e.to === d.to)) return {};
          return {
            threads: {
              ...state.threads,
              [d.thread_id]: { ...t, edges: [...t.edges, { from: d.from, to: d.to }] },
            },
          };
        }

        case "being": {
          const b = msg.payload as Being;
          return { beings: { ...state.beings, [b.name]: b } };
        }

        case "weight": {
          const d = msg.payload;
          const b = state.beings[d.being];
          if (!b) return {};
          const list = d.kind === "relationship" ? b.relationships : b.expressors;
          const next = list.map((edge) =>
            edge.target === d.target ? { ...edge, weight: d.weight } : edge,
          );
          const updated: Being =
            d.kind === "relationship"
              ? { ...b, relationships: next }
              : { ...b, expressors: next };
          return { beings: { ...state.beings, [b.name]: updated } };
        }

        case "memory": {
          const d = msg.payload;
          const b = state.beings[d.being];
          if (!b) return {};
          const memories = { ...b.memories };
          if (d.item) memories.items = [...(memories.items ?? []), d.item];
          if (d.skill) memories.skills = [...(memories.skills ?? []), d.skill];
          return {
            beings: { ...state.beings, [b.name]: { ...b, memories } },
          };
        }

        case "topology": {
          const d = msg.payload;
          // Append the subtree under the named parent. If parent is "universe",
          // it goes on the root's children. Otherwise we search the tree.
          const nextTopology = appendSubtree(state.topology, d.parent, d.subtree);
          return { topology: nextTopology };
        }

        case "error": {
          return {
            lastError: { code: msg.payload.origin, message: msg.payload.message, ts: msg.ts },
          };
        }

        case "auth_fail": {
          return {
            lastError: { code: "auth_fail", message: msg.payload.reason, ts: msg.ts },
          };
        }

        case "auth_ok":
        default:
          return {};
      }
    }),

  recordEvent: (e) =>
    set((state) => ({ events: [...state.events, e].slice(-EVENTS_WINDOW) })),

  setStatus: (s) => set({ status: s }),

  setView: (v) => set({ view: v }),

  setPerspectiveBeing: (name) => set({ perspectiveBeing: name }),

  selectBeing: (name) => set({ selectedBeing: name }),

  dismissError: () => set({ lastError: null }),
}));

// ---------------------------------------------------------------------------
// Selectors
// ---------------------------------------------------------------------------

export const selectBeings = (s: StoreState) => Object.values(s.beings);

export const selectBeing = (name: string) => (s: StoreState) => s.beings[name];

export const selectThreads = (s: StoreState) => Object.values(s.threads);

export const selectExchanges = (s: StoreState) => Object.values(s.exchanges);

/** Exchanges the given being participates in. */
export const selectExchangesFor = (beingName: string) => (s: StoreState) =>
  Object.values(s.exchanges).filter((e) => e.parties.includes(beingName));

/** Peers of the given being (exchange partners). */
export const selectPeersOf = (beingName: string) => (s: StoreState) => {
  const out = new Set<string>();
  for (const e of Object.values(s.exchanges)) {
    if (e.parties[0] === beingName) out.add(e.parties[1]);
    else if (e.parties[1] === beingName) out.add(e.parties[0]);
  }
  return Array.from(out);
};

/**
 * Skyra activity indicator — high recent weight movement means the system
 * is shaping. Low means stable. This selector returns the latest known
 * being weight for skyra; richer activity tracking belongs in a derived
 * timeseries the store does not maintain (per design pressure).
 */
export const selectSkyraWeight = (s: StoreState) =>
  s.beings["skyra"]?.weight ?? null;

// ---------------------------------------------------------------------------
// helpers
// ---------------------------------------------------------------------------

function indexBy<T>(items: T[], key: (t: T) => string): Record<string, T> {
  const out: Record<string, T> = {};
  for (const item of items) out[key(item)] = item;
  return out;
}

function appendSubtree(
  node: TopologyNode,
  parentId: string,
  subtree: TopologyNode,
): TopologyNode {
  if (node.id === parentId) {
    return { ...node, children: [...(node.children ?? []), subtree] };
  }
  const recurse = (arr?: TopologyNode[]) =>
    arr?.map((n) => appendSubtree(n, parentId, subtree));
  return {
    ...node,
    children: recurse(node.children),
    relationships: recurse(node.relationships),
    expressors: recurse(node.expressors),
  };
}
