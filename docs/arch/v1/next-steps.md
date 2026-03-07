# Next Steps — Open Design Questions

## What's Locked

- **Unified max-heap** — all work ordered by importance score. Three inference types: estimation (very high), job (high), batch (very low).
- **Estimation call schema** — `{is_job, complexity, reasoning_depth, cross_domain, reversible, output_scope, domain}`. Complexity ≤ 1 → inline. Complexity > 1 → heap.
- **Estimator is an inference call** — not a service. Fires when estimation work item is picked up. External Router owns the heap.
- **Heap-driven execution loop** — every tool call re-queues. Two exits: `finished` tool and `contact_user` tool. Preemption is free.
- **Object store is git** — each agent is an independent git repo (go-git). Rollback is `git checkout`. Audit trail is `git log`. No custom commit infrastructure.
- **Tools are filesystem files** — live under `tools/` in the agent git repo. LLM discovers via grep/cat/shell. No vector index over tools.
- **Skyrad universal daemon** — one binary, all devices. Brain sends capability-based service package at registration. Brain is an elected role.
- **Spatial awareness** — ingress shard network fingerprint is the location anchor. Capability resolver filters to co-located shards.
- **Working state** — `working/` is gitignored scratch space. Committed state requires user approval via `propose_commit`.
- **Preemptive scheduling** — natural property of heap re-entry. No FIFO stack needed.

---

## Open Design Questions

### 1. Quick Reply + Progressive Delivery Model

Every user request generates two things:

1. **Quick reply job** — high priority, hits the heap immediately. LLM's only job is to respond to the user right now. Direct answer, acknowledgment, or "working on it" — never keeps the user waiting. Response committed to context engine, CIX pushes it, ingress shard renders it.
2. **Deeper work job** — if the quick reply determines more work is needed, it re-queues a new job onto the heap. Tool calls, deep reasoning, multi-step execution. Updates flow to context engine → CIX → ingress shard as they arrive.

Questions:
- Is the quick reply a separate inference type on the heap, or does the estimation call double as the quick reply?
- What's the schema for a quick reply commitment to the context engine?
- How does the ingress shard know a new response is ready — does it watch the context cache for new completed turns, or does CIX signal it explicitly?
- How does the quick reply LLM decide whether to answer directly vs acknowledge and defer?

### 2. Decoupled Response Delivery

Responses don't come back through a direct reply channel. They get committed to the context engine. CIX pushes an updated context package to the ingress shard. The ingress shard renders new completed turns as they arrive.

Questions:
- What does the context package look like when it carries a pending response — is it a completed turn, a partial turn, or a separate field?
- How does the ingress shard distinguish a "render this now" update from a background context refresh?
- What happens if the ingress shard changes between request and response (user moved rooms)?
- Does the response target a specific ingress shard or broadcast to all active shards?

### 3. Mic Auto-Switching + Duplicate Tiebreaker

Active ingress shard is whichever shard most recently received user audio. If two shards pick up the same utterance, duplicate detection via `(session_id, turn_id)` fires — tiebreaker is amplitude. Louder = closer = right shard.

Questions:
- Where does amplitude get captured and attached to the event — at STT time or as a separate field on voice_event?
- How does the brain know which shard to treat as active ingress for the response path?
- What's the session handoff model when the active shard switches mid-session?

### 4. Estimation Call Schema — Remaining Fields

Schema is expanded but two fields remain unresolved:
- `importance` — composite heap ordering score. Derived here or by the front face transformer upstream?
- `latency_class` — `interactive | background`. Already on `triage_hints` from the ingress shard. Does it flow through or get re-derived?

### 5. "Other" Turn Storage

Turns labeled "other" by the front face transformer get stored for batch pickup.

Minimum fields: `turn_id`, `session_id`, `event_id`, `transcript`, `context_blob_ref`, `routed_agents[]`, `created_at`

Questions:
- Does the context blob need to be snapshotted at ingress, or can batch reconstruct it?
- What is the retention policy?

### 6. Batch Job Contract

Nightly batch runs all agents against accumulated turns they didn't receive in real-time.

Questions:
- One heap item per agent or per turn-agent pair?
- What model runs batch inference — lightweight or full?
- How does batch handle an archived agent?

---

## Related Docs

- `docs/arch/v1/scheduler.md` — unified heap, inference types, complexity scoring, preemption
- `docs/arch/v1/executor.md` — heap-driven execution loop, working state
- `docs/arch/v1/task-formation.md` — domain agent as doorkeeper, estimation call
- `docs/arch/v1/context-engine.md` — context blob, CIX, batch weight updates
- `docs/arch/v1/importance-vectors.md` — importance vector design, V3 background process
- `skyra/internal/agent/README.md` — object store, git model, tools
- `skyra/schemas/ingress/voice/` — voice_event schema, context_blob, location_tag
