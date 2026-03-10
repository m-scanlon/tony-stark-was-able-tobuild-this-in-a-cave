# Primitive Skill Types

All work in Skyra is a skill instantiated as a job. Skills flow through the heap as jobs, get picked up by the right shard, and execute through the same path regardless of type. This document catalogs every primitive skill type in the system.

---

## Primitive Skill Type Registry

| Type | Trigger | Live user? |
|---|---|---|
| Shard registration | Shard startup | No |
| Capability verification | Brain → shard post-registration | No |
| Conversation | User utterance | Yes |
| Cron | Timer (cron skill fires on schedule) | No |
| Retrieval | Context background loop, background loop | No |
| Batch | Nightly schedule | No |
| Sub-skill | Task ReAct loop | Inherits from parent |
| Skill acquisition | Skill gap detected | Maybe |
| Compaction | Job tree completion | No |
| Context background | State change | No |

---

## Descriptions

### Shard Registration
Fired when a shard starts up. The shard fingerprints its hardware and registers its capabilities with the brain via `skyra <registration_tool> [args]`. Flows through the heap as a system skill. See `docs/arch/v1/shard-registration.md`.

### Capability Verification
The brain's response to a shard registration skill. Brain sends a command to the newly registered shard to probe and verify its advertised capabilities. Same mechanism used on reconnection.

### Conversation
The primary skill type. Triggered by a user utterance. The API Gateway resolves the skill, assembles the job envelope, dispatches to the kernel. Has a live user waiting on the response path.

### Cron
A system primitive skill that fires on a timer. Pre-provisioned at boot. No separate service — cron is just a skill. Identical to any other skill in the execution path — `skyra <tool> [args]` through the kernel. No live user session attached. Output committed to memory.

### Retrieval
A system primitive skill. Invoked by the context background loop and any other consumer that needs ranked memory results. Pre-provisioned at boot. The context engine calls `skyra retrieve [args]` — retrieval logic lives in the skill, not embedded in the caller. See `docs/arch/v1/context-engine.md`.

### Batch
Nightly. Runs all accumulated turns against skills that were not routed to in real-time. Preserves data integrity — nothing is permanently missed. Runs at very low heap priority, picks up on idle compute. See `docs/arch/v1/scheduler.md`.

### Sub-skill
Spawned by a task during its ReAct loop when it needs to delegate or spawn replicas of itself. Carries a `parent_task_id` linking it to the parent job. Completion propagates up the job tree via the closure table. See `docs/arch/v1/kernel.md`.

### Skill Acquisition
Triggered when a skill gap is detected. Skyra uses base skills (Google Search, Code Execution) to find, build, and register a new skill. Requires user approval before the skill is provisioned in Redis. See `docs/arch/v1/skyra-skills.md`.

### Compaction
Runs after a job tree completes. Processes the raw execution trace and produces two outputs: (1) OTEL observability data — full trace, timing, errors; (2) refined session data — key decisions, outcomes, facts learned — committed to relevant memory namespaces. Design not yet started. See `docs/arch/v1/gaps.md` G26.

### Context Background
Proactive context update skill. Watches for state changes across jobs, turns, and memory state. Runs inference and commits observations to the context engine. Keeps the context package fresh without blocking the request path. See `docs/arch/v1/context-engine.md`.
