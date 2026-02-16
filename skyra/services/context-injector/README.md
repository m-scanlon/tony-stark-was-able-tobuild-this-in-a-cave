# Skyra Context Injector Service

Background service that proactively builds and pushes compressed context packages to the listener/front-door stack.

This service is a separate process from the control-plane API, but runs in the same control-plane trust zone.

## Purpose

- Keep listener/front-door context fresh without blocking user interaction.
- Perform heavy retrieval/ranking/compression off the listener device.
- Push ready-to-use context into listener cache with low latency.

## Why Separate Service

- Continuous background loop with different runtime profile than request/response APIs.
- Independent deployment, tuning, and restart policy.
- Failure isolation: listener and control-plane API remain responsive if injector degrades.

## Core Responsibilities

- Subscribe to events:
  - conversation turns
  - intent/project hints
  - active task updates
  - time-based events
  - high-priority memory changes
- Gather candidate context from memory/vector/task stores.
- Rank and stabilize selected context (anti-thrashing).
- Compress selected context to token budget.
- Push versioned context package to listener cache.

## Data Flow

1. Listener emits lightweight events (`intent_hint`, `project_hint`, `turn_id`).
2. Control-plane components emit task/memory/time events.
3. Context Injector consumes events and updates session snapshot.
4. Injector computes new package (rank + compress).
5. Injector pushes package to listener context cache.
6. Front-door model reads `base + live + injected` segments at inference time.

## Trigger Strategy

- Immediate refresh:
  - high-priority event
  - project switch
  - strong intent shift
- Periodic refresh:
  - every `10-20s`
- TTL refresh:
  - package TTL target `60-120s`

## Context Budget (Recommended)

For front-door model context window `T`:

- `35%` system instructions
- `25%` live conversation
- `25%` injected context package
- `15%` response/scratch headroom

## Package Schema (v0)

```json
{
  "package_id": "ctxpkg_2026-02-16T18:04:09Z_4f2a",
  "session_id": "sess_abc123",
  "version": 42,
  "created_at": "2026-02-16T18:04:09Z",
  "ttl_seconds": 90,
  "stale": false,
  "intent_hint": "work.soc2_draft",
  "project_hint": "work",
  "items": [
    {
      "id": "mem_91",
      "type": "recent_decision",
      "source": "object_store",
      "score": 0.93,
      "tokens": 120,
      "content": "Last week you chose weekly Tekkit backup snapshots at 02:00 UTC."
    }
  ],
  "budget": {
    "injected_tokens": 2048,
    "used_tokens": 560
  },
  "confidence": 0.86
}
```

## Injection Event Schema (v0)

```json
{
  "event_id": "injevt_77a1",
  "session_id": "sess_abc123",
  "listener_id": "pi-livingroom-01",
  "trigger": "intent_shift",
  "reason": "user switched from gym planning to work compliance",
  "old_version": 41,
  "new_version": 42,
  "package_id": "ctxpkg_2026-02-16T18:04:09Z_4f2a",
  "published_at": "2026-02-16T18:04:10Z"
}
```

## v0 Implementation Plan

1. Build event subscriber (Redis Streams or NATS).
2. Implement ranking + anti-thrash selection.
3. Reuse `skyra/internal/context/compress` for token-bounded compression.
4. Add push client to listener cache endpoint.
5. Emit telemetry:
   - refresh latency
   - package hit rate
   - staleness rate
   - dropped candidate count

## Failure Behavior

- Injector unavailable:
  - listener keeps last package until TTL grace expires
  - then continues with base + live context only
- Stale package:
  - mark `stale=true`
  - front-door asks short clarification when confidence is low

## Placement

- Run as `context-injector` service in control-plane zone (Mac mini initially).
- Keep listener-side cache on Raspberry Pi.
