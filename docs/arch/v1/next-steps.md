# Next Steps — Open Design Questions

## Architecture — Locked (2026-03-08)

- **Router → Kernel** — the router is now the kernel. Canonical doc: `docs/arch/v1/kernel.md`.
- **Syntax** — `octos <tool> [args]`. One prefix. All tools. API Gateway emits this. Kernel resolves `tool` against Redis.
- **Agents removed** — replaced by: Skill (class), Job (instance), Task (execution unit), Memory (provisioned namespace), Entity (named thing inside memory).
- **Skills live in memory** — Redis is the trust boundary. Skills in memory are inert until provisioned in Redis.
- **Skills are learned** — not manually defined. Kernel pattern recognition watches observational streams. Pattern crosses threshold → skill crystallizes.
- **Memory emerges from observation** — kernel emits provisioning event when entity signal crosses threshold → user approves → namespace created.
- **Scheduler + Executor → kernel internals** — not standalone services.
- **Pattern recognition → kernel function** — triggered on schedule by the Cron Service.
- **Cron Service** — the only time-aware component. Fires scheduled events. Kernel has no clock. Primary scheduled skill: `reasoning`.
- **Kernel invocation paths** — two: (1) Skyra → domain skills via `octos <tool> [args]`. (2) Cron Service → system skills via scheduled events.
- **Terminology layer** — user-configurable labels in `skyra.user`. Default: Skill/Job/Task/Memory.

---

## What's Locked

- **Unified max-heap** — all work ordered by importance score. Three inference types: estimation (very high), job (high), batch (very low).
- **Estimation call schema** — `{is_job, complexity, reasoning_depth, cross_domain, reversible, output_scope, domain}`. Complexity ≤ 1 → inline. Complexity > 1 → heap.
- **Estimator is an inference call** — not a service. Fires when estimation work item is picked up. Kernel owns the heap.
- **Heap-driven execution loop** — every tool call re-queues. Preemption is free.
- **Skyrad universal daemon** — one binary, all devices. Brain sends capability-based service package at registration. Brain is an elected role.
- **Spatial awareness** — ingress shard network fingerprint is the location anchor. Capability resolver filters to co-located shards.
- **Preemptive scheduling** — natural property of heap re-entry. No FIFO stack needed.

---

## Open Design Questions

### Mic Auto-Switching + Duplicate Tiebreaker

Active ingress shard is whichever shard most recently received user audio. If two shards pick up the same utterance, duplicate detection via `(session_id, turn_id)` fires — tiebreaker is amplitude. Louder = closer = right shard.

Questions:
- Where does amplitude get captured and attached to the event — at STT time or as a separate field on voice_event?
- How does the brain know which shard to treat as active ingress for the response path?
- What's the session handoff model when the active shard switches mid-session?

### Estimation Call Schema — Remaining Fields

Schema is expanded but two fields remain unresolved:
- `importance` — composite heap ordering score. Derived here or by the front face transformer upstream?
- `latency_class` — `interactive | background`. Already on `triage_hints` from the ingress shard. Does it flow through or get re-derived?

---

## Implementation Tasks

### Skyrad Registration

Design is complete in `docs/arch/v1/shard-registration.md`. This needs to be implemented.

The algorithm to build:
1. `device_fingerprint` event → heap
2. Brain picks up fingerprint, installs skyrad service package on device
3. Skyrad boots, self-tests each capability, emits `capabilities_installed`
4. Kernel picks up `capabilities_installed`, generates one `capability_test` event per capability → heap → skyrad
5. Skyrad executes each test, responds with `capability_test_complete`
6. Kernel collects results, writes confirmed capabilities to shard registry → shard active

Reconnection re-runs the same capability test round. Partial registration (some capabilities fail) is a valid state.

---

## Related Docs

- `docs/arch/v1/kernel.md` — canonical kernel architecture
- `docs/arch/v1/scheduler.md` — unified heap, inference types, complexity scoring, preemption
- `docs/arch/v1/executor.md` — heap-driven execution loop
- `docs/arch/v1/context-engine.md` — context assembly, retrieval pipeline
- `docs/arch/v1/importance-vectors.md` — importance vector design
- `docs/arch/v1/shard-registration.md` — shard registration algorithm
