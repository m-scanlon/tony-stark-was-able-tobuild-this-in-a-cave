# Skyra Event Ingress and ACK Design

## 1. Overview

The Event Ingress and ACK system guarantees reliable transfer of listener-generated events to the control plane with durable persistence and duplicate-safe handling.

Reliability goals:

- no lost events across normal crash/restart/network failure cases
- at-least-once delivery from listener to control plane
- idempotent duplicate handling
- ACK only after durable inbox write commit

Out of scope:

- task planning/execution
- context injection design
- inbox retention/garbage collection

## 2. Terminology

- **Event**: structured message produced by listener/front-door for downstream orchestration.
- **Outbox**: durable local queue on Pi for unacked events.
- **Inbox**: durable control-plane store for accepted events.
- **ACK**: control-plane confirmation that `event_id` is durably stored.
- **Idempotency**: processing duplicate deliveries without duplicate side effects.
- **At-least-once delivery**: sender retries until ACK; duplicates are possible.

## 3. System Architecture

Components:

- listener/front-door event producer (Pi)
- Pi outbox + retry sender
- transport channel (WebSocket initially, gRPC stream optional)
- control-plane ingress handler
- control-plane SQLite inbox

```text
+---------------------+       +-----------------------+       +----------------------+
| Pi Listener Node    |       | WS / gRPC Transport   |       | Control Plane        |
| front-door producer |-----> | send event envelope   |-----> | ingress + inbox      |
+----------+----------+       +-----------+-----------+       +----------+-----------+
           |                                  ^                           |
           v                                  |                           v
    +--------------+                    ACK(event_id)               +-------------+
    | Pi Outbox    | <--------------------------------------------- | SQLite Inbox|
    | durable queue|                                                | event_id PK |
    +--------------+                                                +-------------+
           |
           v
   delete only after ACK
```

## 4. Event Envelope Schema

Required fields:

- `event_id`
- `type`
- `ts`
- `session_id`
- `device_id`
- `payload`

Example:

```json
{
  "schema_version": 1,
  "event_id": "evt_01JZ4J1NZ0A1G8R4J8X3P4H2WG",
  "type": "proposal.task",
  "ts": "2026-02-16T23:10:21Z",
  "session_id": "sess_01JZ4J1M9V52K8GSRP8YQ0N2YA",
  "device_id": "pi-livingroom-01",
  "payload": {
    "intent": "server.log_summary",
    "confidence": 0.84,
    "user_text": "summarize last night crash logs"
  }
}
```

## 5. ACK Protocol

ACK message:

```json
{
  "event_id": "evt_01JZ4J1NZ0A1G8R4J8X3P4H2WG",
  "ack_ts": "2026-02-16T23:10:22Z",
  "status": "stored"
}
```

Rules:

- control plane sends ACK only after SQLite commit succeeds
- Pi deletes outbox row only after matching ACK
- duplicate event receipt returns ACK again (idempotent behavior)
- invalid envelopes return error/NACK and remain in outbox for retry

## 6. Pi Outbox Design

Storage:

- local SQLite DB on Pi (`listener_outbox.db`)
- WAL mode enabled

Recommended outbox fields:

- `event_id` (PRIMARY KEY)
- `payload_json`
- `created_at`
- `next_attempt_at`
- `attempt_count`
- `last_error`

Retry:

- background loop sends rows where `next_attempt_at <= now`
- exponential backoff with jitter (cap max delay)
- on ACK: delete row
- on failure/no ACK: increment attempts + reschedule

## 7. Control Plane Inbox Design

SQLite is used for durable inbox storage.

Why SQLite:

- low operational overhead
- strong local durability guarantees
- good fit for single control-plane node
- straightforward backup and inspection

Durability settings:

- `PRAGMA journal_mode=WAL;`
- `PRAGMA synchronous=FULL;`
- `PRAGMA busy_timeout=5000;`

Table schema:

```sql
CREATE TABLE IF NOT EXISTS event_inbox (
  event_id TEXT PRIMARY KEY,
  status TEXT NOT NULL,
  received_at TEXT NOT NULL,
  last_updated_at TEXT NOT NULL,
  payload TEXT NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_event_inbox_status ON event_inbox(status);
CREATE INDEX IF NOT EXISTS idx_event_inbox_received_at ON event_inbox(received_at);
```

Idempotency behavior:

- `event_id` PK prevents duplicate durable rows
- duplicate delivery triggers ACK without reinserting logical duplicate

## 8. Ingress Flow

Flow:

1. event envelope received over WS/gRPC
2. envelope validated
3. DB transaction starts
4. insert into inbox by `event_id`
5. commit transaction
6. send ACK
7. downstream processors consume inbox later

Pseudocode:

```text
onEvent(envelope):
  if !valid(envelope):
    return nack("invalid envelope")

  now = utcNow()
  begin tx
    try insert(event_id, status="received", received_at=now, last_updated_at=now, payload=json)
    on duplicate key:
      no-op
  commit tx

  sendAck(event_id, status="stored")
```

## 9. Failure Scenarios

Network drop before ACK:

- event may already be stored
- Pi retries same `event_id`
- control plane detects duplicate and re-ACKs

Duplicate events:

- expected with at-least-once delivery
- safely handled by inbox PK + idempotent ACK

Control plane crash after receive:

- crash before commit: no durable row, no ACK, Pi retries
- crash after commit before ACK: row exists, retry gets duplicate ACK path

Pi reboot:

- outbox persists locally
- sender resumes unsent/unacked events after restart

## 10. Service Placement Decision

Decision:

- implement ingress + inbox + ACK as part of the control-plane stack (not separate service yet)

Why:

- tight coupling with orchestrator event lifecycle
- fewer hops and simpler failure handling
- lower operational complexity for current single-node control plane

Tradeoffs:

- ingress scaling tied to control-plane scaling
- independent rollout is harder than with separate gateway

Mitigation:

- keep ingress/inbox modules isolated so extraction into separate event gateway is possible later

## 11. Repository Layout

Placement decision for this repo:

```text
skyra/
  internal/
    event/
      envelope.go
      protocol.go
    controlplane/
      ingress/
        handler.go
        transport_ws.go
        transport_grpc.go
      inbox/
        store.go
        schema.sql
  services/
    listener/
      internal/
        outbox/
          store.go
          retry_loop.go
          sender.go
```

Responsibilities:

- `services/listener/internal/outbox`: local persistence + retry sender
- `internal/controlplane/ingress`: receive + validate + ACK
- `internal/controlplane/inbox`: durable write/read primitives
- `internal/event`: shared envelope and ACK protocol types

## 12. Future Extensions (Brief)

- inbox garbage collection and retention policy
- task creation worker from inbox rows
- multi-node control plane with shared durable event log
- transport/inbox metrics and telemetry
