# Job Types

All work in Skyra is a job. Jobs flow through the heap, get picked up by the right shard, and execute through the same path regardless of type. This document catalogs every job type in the system.

---

## Job Type Registry

| Type | Trigger | Live user? |
|---|---|---|
| Shard registration | Shard startup | No |
| Capability verification | Brain → shard post-registration | No |
| Conversation | User utterance | Yes |
| Cron | Timer | No |
| Batch | Nightly schedule | No |
| Sub-job | Agent ReAct loop | Inherits from parent |
| Skill acquisition | Agent skill gap | Maybe |
| Compaction | Job tree completion | No |
| Context background | State change | No |

---

## Descriptions

### Shard Registration
Fired when a shard starts up. The shard fingerprints its hardware and registers its capabilities with the brain. Flows through the heap as a system event. See `docs/arch/v1/shard-registration.md`.

### Capability Verification
The brain's response to a shard registration event. Brain sends an owner event to the newly registered shard to probe and verify its advertised capabilities. Same mechanism used on reconnection.

### Conversation
The primary job type. Triggered by a user utterance. The domain agent receives the turn and context, produces a complexity score via the estimation prompt, and either executes inline (complexity ≤ 1) or forms a job on the heap (complexity > 1). Has a live user waiting on the response path.

### Cron
A recurring skill invocation scheduled by a domain agent. Fires a heap event at the scheduled time — identical to a conversation job in the execution path. No live user session attached. Output is committed to the agent object store. See `docs/arch/v1/skyra-skills.md` and `skyra/internal/agent/README.md`.

### Batch
Nightly. Runs all accumulated turns against domain agents that were not routed to in real-time. Preserves data integrity — nothing is permanently missed. Runs at very low heap priority, picks up on idle compute. See `docs/arch/v1/scheduler.md`.

### Sub-job
Spawned by a domain agent during its ReAct loop when it needs to delegate to another agent or spawn replicas of itself. Carries a `parent_task_id` linking it to the parent job. Completion propagates up the job tree via the closure table. See `docs/arch/v1/router.md`.

### Skill Acquisition
Triggered when an agent encounters a gap in its skill set. Skyra uses base skills (Google Search, Code Execution) to find, build, and register a new skill. Requires user approval before the skill is committed to the registry. See `docs/arch/v1/skyra-skills.md`.

### Compaction
Runs after a job tree completes. Processes the raw execution trace and produces two outputs: (1) OTEL observability data — full trace, timing, errors; (2) refined session data — key decisions, outcomes, facts learned — committed to relevant agent object stores. Design not yet started. See `docs/arch/v1/gaps.md` G26.

### Context Background
Proactive context update loop. Watches for state changes across jobs, turns, and agent state. Runs inference and commits observations to the context engine. Keeps the context package fresh without blocking the request path. See `docs/arch/v1/context-engine.md`.
