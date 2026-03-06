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
- **Event Register**: durable hash table (keyed by `turn_id`) on the Voice Shard tracking sent events awaiting ACK.
- **Inbox**: durable control-plane store for accepted events.
- **ACK**: control-plane confirmation that the event for a given `turn_id` is durably stored.
- **Idempotency**: processing duplicate deliveries without duplicate side effects.
- **At-least-once delivery**: sender retries until ACK; duplicates are possible.

## 3. System Architecture

Components:

- listener/front-door event producer (Voice Shard)
- Voice Shard event register + retry sender
- transport channel (WebSocket initially, gRPC stream optional)
- control-plane ingress handler
- control-plane SQLite inbox

```text
+---------------------+       +-----------------------+       +----------------------+
| Voice Shard         |       | WS / gRPC Transport   |       | Control Plane        |
| front-door producer |-----> | send event envelope   |-----> | ingress + inbox      |
+----------+----------+       +-----------+-----------+       +----------+-----------+
           |                                  ^                           |
           v                                  |                           v
    +----------------+                  ACK(turn_id)               +-------------+
    | Voice Shard    | <------------------------------------------- | SQLite Inbox|
    | Event Register | <------------------------------------------- | event_id PK |
    +----------------+                                              +-------------+
           |
           v
   pop by turn_id only after ACK
```

## 4. Event Envelope Schema

Required fields (sent by Voice Shard):

- `schema`
- `turn_id`
- `ts`
- `device_id`
- `transcript`
- `triage_hints`
- `session_state`

Note: `event_id` is NOT sent by Voice Shard. Brain Shard generates `event_id` (ULID) on ingress. Voice Shard provides `(session_id, turn_id)` as its idempotency pair — Brain Shard uses this composite for duplicate detection. All fields except `triage_hints` are stamped by the shard transport layer during hydration — see `skyra/schemas/ingress/voice/voice_event_v1.json` for the full schema.

Example:

```json
{
  "schema": "voice_event_v1",
  "turn_id": "turn_8f4c",
  "ts": "2026-02-20T18:10:12Z",
  "device_id": "pi-livingroom-01",
  "transcript": "what did I decide about backups",
  "triage_hints": {
    "intent": {
      "summary": "user wants to know what was decided about backups",
      "confidence": 0.94
    },
    "latency_class": {
      "value": "interactive",
      "confidence": 0.88
    },
    "ack_policy": {
      "value": "spoken_if_slow",
      "confidence": 0.76
    }
  },
  "session_state": {
    "pending_job_id": null,
    "waiting_for": null
  }
}
```

## 5. ACK Protocol

ACK message:

```json
{
  "turn_id": "turn_8f4c",
  "ack_ts": "2026-02-16T23:10:22Z",
  "status": "stored"
}
```

Rules:

- ACK references `turn_id` — Voice Shard never sees `event_id`, which is internal to the control plane
- control plane sends ACK only after SQLite commit succeeds
- Voice Shard pops entry from event register by `turn_id` on matching ACK
- duplicate delivery (same `session_id` + `turn_id`) returns ACK without reinserting
- invalid envelopes return error/NACK and remain in event register for retry

## 6. Voice Shard Event Register

Storage:

- local SQLite DB on Voice Shard (`listener_event_register.db`)
- WAL mode enabled

Fields:

- `turn_id` (PRIMARY KEY — Voice Shard-generated, stable across retries)
- `session_id` (paired with `turn_id` for Brain Shard-side deduplication)
- `payload_json`
- `created_at`
- `next_attempt_at`
- `attempt_count`
- `last_error`

Retry:

- background loop sends rows where `next_attempt_at <= now`
- exponential backoff with jitter (cap max delay)
- on ACK: delete row by `turn_id`
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
  event_id        TEXT PRIMARY KEY,      -- Brain Shard-generated ULID
  session_id      TEXT NOT NULL,
  turn_id         TEXT NOT NULL,
  status          TEXT NOT NULL,
  received_at     TEXT NOT NULL,
  last_updated_at TEXT NOT NULL,
  payload         TEXT NOT NULL,
  UNIQUE(session_id, turn_id)            -- deduplication key for Voice Shard retries
);
CREATE INDEX IF NOT EXISTS idx_event_inbox_status ON event_inbox(status);
CREATE INDEX IF NOT EXISTS idx_event_inbox_received_at ON event_inbox(received_at);
```

Idempotency behavior:

- `event_id` is Brain Shard-generated (ULID) — unique per ingress attempt, not per logical event
- `(session_id, turn_id)` UNIQUE constraint prevents duplicate rows on Voice Shard retry
- duplicate delivery (same `session_id` + `turn_id`) triggers re-ACK with `turn_id`, no reinsert

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
  event_id = newULID()  # Brain Shard generates event_id — Voice Shard does not provide one

  begin tx
    try insert(event_id, session_id, turn_id, status="received", received_at=now, last_updated_at=now, payload=json)
    on duplicate (session_id, turn_id):
      event_id = lookup_existing_event_id(session_id, turn_id)  # fetch original for internal reference only
  commit tx

  sendAck(turn_id, status="stored")  # Voice Shard deletes outbox row by turn_id; event_id stays internal
```

## 9. Failure Scenarios

Network drop before ACK:

- event may already be stored
- Voice Shard retries with same `turn_id`
- control plane detects duplicate via `(session_id, turn_id)` UNIQUE constraint and re-ACKs with `turn_id`

Duplicate events:

- expected with at-least-once delivery
- safely handled by inbox PK + idempotent ACK

Control plane crash after receive:

- crash before commit: no durable row, no ACK, Voice Shard retries
- crash after commit before ACK: row exists, retry gets duplicate ACK path

Voice Shard reboot:

- event register persists locally
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
        event_register/
          store.go
          retry_loop.go
          sender.go
```

Responsibilities:

- `services/listener/internal/event_register`: local persistence + retry sender
- `internal/controlplane/ingress`: receive + validate + ACK
- `internal/controlplane/inbox`: durable write/read primitives
- `internal/event`: shared envelope and ACK protocol types

## 12. Future Extensions (Brief)

- inbox garbage collection and retention policy
- task creation worker from inbox rows
- multi-node control plane with shared durable event log
- transport/inbox metrics and telemetry
