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

- Receive `context_state` from Event Ingress on every Voice Shard request (fan-out from `voice_event_v1`). Use `available_for_injection` as the live package budget. Fall back to static budget percentages on cold start (no state received yet).
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

The Context Injector does not interface with Voice Shard directly. Voice Shard sends `context_state` to the Brain Shard API Gateway as part of every `voice_event_v1`. Event Ingress fans this out internally to the Context Injector on Brain Shard.

## Data Flow

1. Voice Shard sends `voice_event_v1` to Brain Shard API Gateway on every request. Event includes `context_state` (`available_for_injection` and token breakdown).
2. Brain Shard Event Ingress fans `context_state` out to Context Injector internally. Context Injector updates its live budget for the next package.
3. Control-plane components emit task/memory/time events to Context Injector.
4. Context Injector consumes events and updates session snapshot.
5. Injector computes new package: rank + compress to fit `available_for_injection` tokens exactly.
6. Injector pushes package to Voice Shard listener context cache (LCACHE).
7. Front-door model reads `base + live + injected` segments at inference time.

## Trigger Strategy

- Immediate refresh:
  - high-priority event
  - project switch
  - strong intent shift
- Periodic refresh:
  - every `10-20s`
- TTL refresh:
  - package TTL target `60-120s`

## Context Budget

The injected package budget is set dynamically from `context_state.available_for_injection` received via Event Ingress fan-out on every Voice Shard request.

`available_for_injection = total_context_tokens - system_tokens - live_conversation_tokens - response_reserve_tokens`

Voice Shard computes this. Context Injector uses it directly as the token budget for the next package. This means the package grows when conversation is sparse and shrinks as the conversation accumulates — always filling the available headroom exactly.

**Cold start fallback** (no `context_state` received yet):

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
  "agent_hint": "work",
  "items": [
    {
      "id": "mem_91",
      "type": "recent_decision",
      "source": "object_store",
      "score": 0.93,
      "tokens": 120,
      "retrieved_at": "2026-02-20T17:41:00Z",
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

`retrieved_at` is required on every item. The Voice Shard listener uses it to determine whether an item is fresh enough to support a provisional answer (threshold: 30 minutes). Items missing `retrieved_at` are treated as stale.

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

- Run as `context-injector` service in control-plane zone (Brain Shard initially).
- Keep listener-side cache on Voice Shard.
