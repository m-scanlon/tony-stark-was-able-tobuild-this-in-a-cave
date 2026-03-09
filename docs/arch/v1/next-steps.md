# Next Steps — Open Design Questions

## Architecture Revision — Locked (2026-03-08)

- **Router → Kernel** — the router is now the kernel. Same switch statement. Renamed throughout. Canonical doc: `docs/arch/v1/kernel.md`.
- **Syntax** — `skyra <tool> [args]`. One prefix. No agent middle layer. API Gateway emits this. Kernel resolves `tool` against Redis.
- **Agents removed** — replaced by primitives: Skill (class), Job (instance), Task (execution unit), Memory (provisioned namespace), Entity (named thing inside memory).
- **Skill = class, Job = instance** — same relationship as programming. Skill is the registered roadmap. Job is what you get when the kernel invokes it. `job_envelope_v1` → job now holds `skill_id` reference.
- **Skills live in memory** — the object store is memory. Skills are stored there. Redis is the trust boundary — skills in memory are inert until provisioned in cache.
- **Skills are learned** — not manually defined. Kernel's pattern recognition function watches observational streams. When a pattern crosses threshold, a skill crystallizes.
- **Memory is provisioned by the kernel** — on user approval. Not created manually. Emerges from entity observation crossing a threshold. Kernel emits provisioning event → user approves → namespace created.
- **Entities live inside memory** — entities are named things (people, places, tools, concepts) accumulated through observation. They live in memory namespaces — they don't own them.
- **Scheduler + Executor → kernel internals** — not standalone services. Internal methods of the kernel. Detail docs (`executor.md`, `scheduler.md`) preserved as kernel internals.
- **Pattern recognition engine → kernel function** — not a separate engine. Triggered on schedule by the Cron Service.
- **Cron Service = System Event Issuer** — the only time-aware component. Sits outside the kernel. Fires scheduled events onto the heap. Kernel has no clock.
- **Kernel invocation paths** — two: (1) Skyra → domain skills via `skyra <tool> [args]`. (2) Cron Service → system skills via scheduled events. Kernel never self-initiates.
- **Terminology layer** — user-configurable labels stored in `skyra.user`. Default: Skill/Job/Task/Memory. Programmer persona: Class/Instance/Task/Repo. Execution model unchanged.
- **`skyra` and `delegate` are system primitives** — hardcoded switch cases in the kernel. Redis-independent. Three hardcoded skills: `reply`, `fan_out`, `report`.

---

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

- **Skyra as delegation layer** — Skyra knows her own tools and that other agents exist, but not what tools those agents expose. She delegates via `skyra <agent> <tool> [args]`. She can fan out to multiple domains simultaneously.
- **Two-level reasoning** — Skyra gets the agent registry (agent names + domain descriptions) in her context. She reasons at domain level. Each agent gets its own tool list injected at inference time. Neither gets more than they need.
- **Unified syntax** — `skyra <agent> <tool> [args]` throughout. No skyrad prefix. Personal agent is also named `skyra` so system tools are `skyra skyra reply`.
- **Router structure** — switch on `command.agent`. `case "skyra"` handles system primitives (switch on `command.tool`, currently only `reply`). `case redis.get("agent:" + command.agent)` handles all dynamic agents — validates tool exists, builds context (Skyra's command + user message + tool list from Redis), pushes inference job to heap. `default` errors agent not found.
- **Every agent needs a delegation entry point tool** — registered in Redis at agent creation. The router validates it before building context for the inference call.
- **Router is hybrid** — `skyra` system primitives are hardcoded. All other agents are fully dynamic via Redis. Adding a new agent = registering in Redis, no router code changes.
- **Redis as live registry + SQLite as backing store** — Redis owns real-time tool/agent state. Keyspace notifications push updates. SQLite persists for brain restarts. CLI reads Redis directly on each tool call validation — no local cache.
- **Tools live in the capability registry (Redis), not the object store** — object store is purely state and memory. Tool definitions live in Redis. Router reads Redis to dispatch.
- **Agent state is distributed** — agent registry tracks location per component: `storage` (which shard holds the object store), `websocket` (which shard manages the connection). The registry is the location map.
- **Two core system agents** — `skyra` (user interaction, the face) and `delegate` (multi-agent coordinator, the engine). Both hardcoded in router. Redis-independent. Always available.
- **Three system skills** — `skyra delegate fan_out` (opens job, fans out), `skyra delegate report` (agent → delegate), `skyra skyra reply` (Skyra → user only). Nothing else talks to the user.
- **Delegate is an active coordinator** — validates exit conditions, reprompts failing agents with context (up to N retries), escalates to Skyra only when retries exhausted. Delegate is a pure state machine — no inference of its own.
- **ReAct loop per domain agent** — every agent runs Reason → Act → Observe → Repeat until task complete, then `skyra delegate report`. Multi-step skill calls within a single task. Delegate reprompt restarts the loop with failure context injected.
- **Delegator is the progressive delivery mechanism** — as each task completes and reports back via `skyra delegate report`, the delegator notifies Skyra incrementally. Skyra decides whether to reply to the user immediately or wait for all tasks. Job pops when all tasks complete or TTL expires.
- **Job/task data structure** — SQLite is source of truth (jobs table + tasks table). Redis pub/sub is the real-time signal. Router writes to both on task completion. Dispatcher subscribes to Redis, reads SQLite to confirm state, pops job when all tasks complete.
- **Progressive delivery resolved** — open questions 1 and 2 (quick reply + progressive delivery, decoupled response delivery) are closed. The delegator owns incremental updates. Skyra owns the reply decision.

---

## Open Design Questions

### 7. Job Tree Tracking — RESOLVED

The delegator owns the job tree. `skyra delegate` creates a job with N tasks (one per agent). Tasks report back via `skyra delegate report`. Delegator notifies Skyra incrementally. Job pops when all tasks complete.

SQLite schema:
```
jobs:  job_id, turn_id, session_id, status, created_at, completed_at
tasks: task_id, job_id, agent, status, result, created_at, completed_at
```

Open questions remaining:
- What happens if one task fails — does the job partial-succeed or fail entirely?
- Does Skyra always reply incrementally or does she decide per job?

---

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

## Implementation Tasks

### Skyrad Registration

Design is complete in `docs/arch/v1/shard-registration.md`. This needs to be implemented.

The algorithm to build:
1. `device_fingerprint` event → heap
2. Brain picks up fingerprint, installs skyrad service package on device
3. Skyrad boots, self-tests each capability, emits `capabilities_installed`
4. External Router picks up `capabilities_installed`, generates one `capability_test` event per capability → heap → skyrad
5. Skyrad executes each test, responds with `capability_test_complete`
6. External Router collects results, writes confirmed capabilities to agent registry → shard active

Reconnection re-runs the same capability test round. Partial registration (some capabilities fail) is a valid state.

---

## Related Docs

- `docs/arch/v1/scheduler.md` — unified heap, inference types, complexity scoring, preemption
- `docs/arch/v1/executor.md` — heap-driven execution loop, working state
- `docs/arch/v1/task-formation.md` — domain agent as doorkeeper, estimation call
- `docs/arch/v1/context-engine.md` — context blob, CIX, batch weight updates
- `docs/arch/v1/importance-vectors.md` — importance vector design, V3 background process
- `skyra/internal/agent/README.md` — object store, git model, tools
- `skyra/schemas/ingress/voice/` — voice_event schema, context_blob, location_tag
