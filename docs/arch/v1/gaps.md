# Skyra Architecture Gaps Register

This document tracks known architecture gaps that must be resolved before production-hardening.

> **Architecture revision note**: Router → Kernel. Agents removed. New primitives: Skill, Job, Task, Memory, Entity. Scheduler and Executor are kernel internals. Pattern recognition engine is a kernel function. Cron Service added as time-aware layer. Gaps below are being updated to reflect this. Gaps that reference "agent" should be read as "memory/skill" until updated.

## Priority Scale

- `P0`: blocks safe/reliable production behavior
- `P1`: high impact, should be resolved early
- `P2`: important, can follow core stabilization

## Gaps

| ID | Gap | Priority | Why it matters | Suggested owner |
| --- | --- | --- | --- | --- |
| G1 | ~~Estimation call schema not locked~~ **CLOSED** | P0 | **Resolved by architecture revision**: The estimator is a kernel internal method. The handoff contract is `job_envelope_v1` defined in `docs/arch/v1/api-gateway.md` — provisions + security + payload. The API Gateway assembles it from Redis. No separate estimation call schema needed. | API Gateway |
| G2 | End-to-end idempotency beyond ingress is incomplete | P0 | Retries can create duplicate tasks, duplicate side effects, or duplicate commits. | Ingress + Task Formation + Executor |
| G3 | ~~Stateful commit safety boundary is underdefined~~ **CLOSED** | P0 | **v1 decision**: Working state vs committed state is fully defined — executor writes freely to `working/`, committed state requires user approval via `propose_commit`. The `allow_always \| allow_once \| deny` permission model is locked. BoundaryValidator enforces at runtime before any tool dispatch. See `skyra/internal/agent/README.md`. Failure/rollback coordination between operational job registry status and semantic tasksheet status moved to G20. | Task Formation + Agent Service |
| G4 | ~~Reconciliation UX policy is not fully specified~~ **CLOSED for v1** | P1 | **v1 decision**: Provisional responses cut from v1. Voice Shard emits non-semantic ACKs only (earcon/LED/short wait phrase). Voice Shard renders only Brain Shard-authored messages (`UPDATE`, `PLAN_PROGRESS`, `CLARIFY`, `PLAN_APPROVAL_REQUIRED`, `FINAL`, `ERROR`). No split-brain risk in v1. Provisional response path (front-door model speaks before Brain Shard responds, then reconciles on `FINAL`) deferred to v2 — see note in `docs/arch/v1/scyra.md` section 7.1. | Voice Shard + Control Plane |
| G5 | AuthN/AuthZ model is not concretely implemented | P0 | Device/agent channels remain security-sensitive attack surface. | Platform/Security |
| G6 | Backpressure and overload policies are undefined | P1 | Heap growth and latency spikes can collapse responsiveness under load. Max heap depth, starvation prevention, and priority bump policies for long-waiting items are undefined. | Ingress + Heap/External Router |
| G7 | ~~Executor/Resource Manager contracts are abstract~~ **CLOSED** | P1 | **Resolved by architecture revision**: Executor is gone as a standalone service — it is a kernel internal method. Execution contract (boundary rules, severity policy, state contract, replan budget) is defined in the skill itself. Resource checking is a skill property. No separate Executor/Resource Manager interface needed. | Kernel |
| G8 | Observability and SLOs are not defined end-to-end | P1 | Hard to debug reliability and latency regressions without traceability and targets. | Platform/Operations |
| G9 | Degradation and cold-start behavior is incomplete | P1 | Unclear runtime behavior when STT/model/estimation inference/context systems fail or lag. | Voice Shard + Orchestrator |
| G10 | Data lifecycle and retention policy is missing | P1 | Inbox/checkpoints/transcripts can grow unbounded and violate privacy expectations. | Memory + Platform |
| G11 | Job Registry schema not locked | P0 | Agent Registry schema is largely defined in `skyra/internal/agent/README.md` (identity, metadata, status, last_active_at). What remains: the Job Registry schema. Needs locked fields for `(job_id, event_id, agent_id, shard_id, status, created_at, updated_at)` plus the full status transition model (`created → routed → planning → executing → done \| failed`). Everything downstream that reads job state depends on this. | Control Plane + Agent Service |
| G12 | ~~Agent state (state.json) four-section schema not locked~~ **CLOSED** | P0 | Schema locked: `metadata/knowledge/artifact/boundary`. Boundary carries `allowed_tool_categories`, `denied_tool_patterns`, `restrictions[]` (no enforcement field — all locked tool attempts prompt the user). Enforced in code via two layers: hydration (lock status attached to tools before LLM sees them) and BoundaryValidator (permission prompt at runtime). See `skyra/internal/agent/README.md`. | Agent Service |
| G13 | ~~Tool system contract not finalized~~ **CLOSED** | P1 | Global tools (always injected, small fixed set) and agent tools (files in `tools/` in the object store, discovered by the LLM walking the filesystem) are fully specified. No vector index over tools — the filesystem is the index. `requires_approval` is a display flag on tool.json. `categories[]` drives boundary enforcement via BoundaryValidator. Lock status derived at runtime from `state.json` boundary — not stored on the tool. See `skyra/internal/agent/README.md`. | Agent Service + Domain Expert |
| G14 | Voice response channel (`/v1/voice`) not implemented | P0 | `/v1/voice` returns a literal string, not a `voice_result_v1` protocol message. No Brain Shard response can reach the user over the voice path. Blocks all voice UX end-to-end. Needs proper streaming response channel (WebSocket or gRPC) emitting `message_type` in `{UPDATE, PLAN_PROGRESS, CLARIFY, PLAN_APPROVAL_REQUIRED, FINAL, ERROR}`. | Control Plane |
| G15 | Plan approval response channel undefined | P0 | Voice Shard has no mechanism to route a user's spoken approval/rejection (`APPROVE \| REVISE \| CANCEL`) back to a `PLAN_APPROVAL_REQUIRED` job. `session_state.pending_job_id` + `waiting_for` contract is defined in the event schema but not implemented anywhere. A user saying "yes" after a plan approval prompt is indistinguishable from a new request. | Voice Shard + Control Plane |
| G16 | ~~`requires_approval` vs `PLAN_APPROVAL_REQUIRED` semantic distinction is a high-risk confusion point~~ **CLOSED** | P1 | **Resolved by documentation**: the distinction is now explicitly called out in `task-formation.md`, `domain-expert/README.md`, `lifecycle.md`, and `skyra/internal/agent/README.md`. Each doc has a dedicated note. `requires_approval` = display flag, highlights tool during plan review, does not pause execution. `PLAN_APPROVAL_REQUIRED` = plan-level gate, all execution waits. Documentation confusion risk eliminated. | Domain Expert + Executor + Control Plane |
| G17 | ~~Agent routing staleness~~ **CLOSED** | P1 | **Resolved by architecture revision**: Agents are gone. The Voice Shard no longer sends `active_agent`. It emits `skyra voice_event "..." --session=... --turn=...`. The API Gateway resolves the skill from Redis dynamically on every command. No staleness problem exists. | API Gateway |
| G18 | ~~Cross-agent write protocol undefined~~ **CLOSED** | P1 | **Resolved by design decision**: Cross-memory writes are not supported. Memory namespaces are isolated. If an execution produces data that belongs in another memory, the skill emits a new `skyra <tool> [args]` command back into the kernel's queue. The kernel handles it through the full trust chain. No direct cross-memory writes ever. | Kernel |
| G19 | Context engine design incomplete | P1 | Core model is locked: proactive commit-based background loop, two partitions (session_history + retrieved), two-level weight update (real-time for routed domains, batch for non-routed), all six data sources defined. What remains: (1) turn persistence storage format not locked — SQLite, JSONL, or object store; (2) session timeout policy not defined; (3) token budget defaults not set; (4) decay model not designed — time-based, relevance-based, or hybrid; (5) real-time weight update mechanism TBD — the behavior is locked (routed domains bump immediately) but the implementation is not — may be inference, embedding similarity, or signal processing. | Context Engine + Agent Service |
| G20 | Job lifecycle failure/rollback coordination partially undefined | P0 | The Job Registry tracks operational status (`created → routed → planning → executing → done \| failed`). The tasksheet tracks semantic phases (`planning → executing → validating → replanning → done`). Agent state rollback is now defined — `git checkout` via go-git reverts `state.json` to any prior commit. What remains undefined: (1) the coordination contract between Job Registry operational status and tasksheet semantic status on failure — e.g. what does the Job Registry record when the executor hits a critical severity assumption failure and halts? (2) What happens to the tasksheet when a job is cancelled mid-execution? (3) At what point does a rollback of agent state get triggered automatically vs requiring user decision? | Task Formation + Executor + Control Plane |
| G21 | ~~Long term memory promotion process undefined~~ **CLOSED** | P2 | **Superseded by G28**: Memory provisioning is now a kernel function triggered by the pattern recognition function crossing a threshold. The kernel emits a provisioning event, the user approves, the kernel provisions the memory namespace. See G28 for the open design work. | Kernel |
| G22 | ~~Shard registration tightly coupled to domain ownership~~ **CLOSED** | P1 | Resolved by dual-registry model: shards register capabilities globally; domain tools remain agent-scoped semantic contracts with `required_capabilities[]`; runtime capability resolver binds tool invocation to the best shard using policy + live state. Mediated shard-to-shard delegation is allowed with control-plane-issued delegation tokens and full audit edges. See `docs/arch/v1/agents-services.md` section 2.8 and `docs/ideas.md` ("Dual Registry: Domain Tools vs Shard Capabilities"). | Control Plane + Agent Service + Tooling |
| G23 | External MCP tool transport contract not locked | P1 | Internal CLI contract is resolved — `skyra <tool> [args]` is the unified protocol. What remains: external capabilities via MCP. MCP dispatch is kernel work — the kernel receives the skill invocation, identifies it as MCP-backed via Redis, and handles the external call. Missing pieces: how MCP tools are provisioned into Redis, kernel MCP dispatch implementation, normalized result/error schema. | Kernel |
| G24 | Session continuity during concurrent job completion and live chat is undefined | P1 | When a long-running job is near completion and the user keeps chatting, the system can lose turn-to-job alignment or inject stale context. Missing contract: how active session turns, pending job state, and completion events are merged; whether chat turns attach to the existing job or start a new one; and how `FINAL` job output reconciles with newer user turns without context corruption or double responses. | Orchestrator + Context Engine + Voice Shard |
| G25 | ~~Skyra binary invocation contract not defined~~ **CLOSED** | P1 | **Resolved by architecture revision**: The unified command protocol `skyra <tool> [args]` is the binary contract. One syntax, all shards, both directions. No per-shard invocation formats. No shell-string mode. See `docs/arch/v1/shard-communication.md`. | Kernel + API Gateway |

| G26 | Job tree compaction engine not designed | P1 | After a job tree completes, the raw execution trace (every ReAct step, tool call, sub-job, intermediate result) needs to be compacted into two distinct outputs: (1) OTEL observability data — full raw trace, timing, errors, tool calls, goes to the observability layer, not agent memory; (2) refined session data — distilled signal: key decisions, outcomes, facts learned, state changes — committed to relevant agent object stores. Without compaction, object stores bloat with raw execution noise and signal gets buried. The compaction engine needs to run post-job, identify what matters per agent, and produce clean commits. Design not started. | Context Engine + Agent Service |

## Recommended Next Actions

1. **Lock canonical contracts (`P0`)**
   - Finalize estimation call schema (`{is_job, complexity, domain}` + any additional Estimator fields). See G1.
   - Lock Job Registry schema — fields and full status transition model. See G11.
   - Define "other" turn storage schema in RDS (what fields, what retention policy).
   - Define idempotency keys for `event_id`, `job_id`, `task_id`, and stateful commit operations. See G2.

2. **Close security and reliability baseline (`P0`)**
   - Choose transport auth strategy (`mTLS` or signed tokens + rotation).
   - Define replay protection and per-device/agent authorization scopes.
   - Define ACK + retry semantics across every boundary (outbox/inbox/task creation/executor side effects).

3. **Finalize execution interfaces (`P1`)**
   - Specify Executor <-> Resource Manager request/response schema.
   - Specify progress snapshot schema for Job Registry updates during execution.
   - Lock tool transport split contract: internal via CLI adapter (`argv`), external via MCP, with one normalized invocation/result schema. See G23.
   - Define the internal Skyra binary contract (`skyra <shard> <tool> <args...>`) and result/error mapping. See G25.
   - Define session/job continuity contract for concurrent completion + live chat turns (`pending_job_id`, turn attachment, FINAL reconciliation). See G24.
   - Define timeout, retry, and error classification semantics.

4. **Define user-facing consistency behavior (`P1`)**
   - ~~Formalize provisional phrasing policy on Voice Shard.~~ **v1 decision**: provisional responses cut. Voice Shard speaks non-semantic ACKs only. See G4.
   - Implement voice response channel (`voice_result_v1` protocol over WebSocket/gRPC). See G14.
   - Define and implement plan approval response channel so Voice Shard can route user decisions back to pending jobs. See G15.
   - Add explicit voice UX tests for `CLARIFY` and `PLAN_APPROVAL_REQUIRED` flows.

5. **Operationalize (`P1`)**
   - Define SLOs: ingress ACK latency, task start latency, completion latency classes.
   - Add trace propagation: `event_id -> job_id -> task_id -> execution_id`.
   - Set retention and compaction strategy for inbox, checkpoints, and transcript artifacts.

## New Gaps — Architecture Revision

| ID | Gap | Priority | Why it matters | Suggested owner |
| --- | --- | --- | --- | --- |
| G27 | Cron Service design not started | P1 | The Cron Service is the time-aware layer — the only component allowed to know about time. It fires system events onto the heap on schedule. No design exists: what schedules are defined, how they're registered, lifecycle (start/stop/pause), failure behavior, how system skills are associated with schedules. Blocks: memory provisioning loop, pattern recognition trigger, any other time-driven kernel work. | Kernel + Platform |
| G28 | Memory provisioning flow not locked | P1 | The kernel provisions memory namespaces on user approval, triggered by the pattern recognition function crossing a threshold. The event format, the user-facing approval channel, the approval contract (`approve | deny`), and the threshold signal definition are all undefined. Without this, memory namespaces can't be created. | Kernel + UX |
| G29 | Skill learning threshold not defined | P1 | Skills are crystallized from observational streams when a pattern crosses a frequency × affect threshold. The exact threshold values, what constitutes a "pattern," the schema for a learned skill, and how the kernel decides a roadmap (1-to-many tasks) from a behavioral pattern are all undefined. | Kernel + Observational Store |
| G30 | Terminology layer not implemented | P2 | Users can configure primitive labels (skill→class, job→instance, memory→repo) stored in `skyra.user`. No surface currently reads from this config. Needs: config schema in `skyra.user`, label resolution at every output surface (CLI, voice responses, logs). | Platform + UX |
| G31 | Tool call syntax change not propagated | P1 | The canonical syntax changed from `skyra <agent> <skill> [args]` to `skyra <tool> [args]`. The API Gateway now owns domain resolution before emitting the command. This change needs to propagate through: command-parser.md, all docs referencing the old syntax, and the CLI/API Gateway implementation. | Kernel + API Gateway |
| G32 | ~~API Gateway domain resolution not designed~~ **CLOSED** | P1 | **Resolved**: Skill reasoning happens at the shard level — the front-door model on the shard reasons about the input and decides which skill to call. The API Gateway is a validator, not a resolver. Command arrives → Redis check → pass or fail. | Shard + API Gateway |
| G33 | Redis-level auth credentials not designed | P0 | Possessing a valid skill name is not sufficient for authorization. The calling shard must prove identity before any command is parsed. **Direction**: two layers — (1) mTLS at the transport layer: every shard gets a certificate issued by the brain at registration. Brain is the CA. Redis only accepts connections from brain-issued certs. No cert, no connection — the command never arrives. (2) Redis authorization layer: after connection is established, Redis checks whether this shard is allowed to invoke this specific skill. Identity at transport. Authorization at registry. Two separate concerns. What remains: certificate issuance flow at registration, cert rotation and revocation policy, per-shard skill authorization schema in Redis. | Platform/Security + API Gateway |
| G34 | Redis write skill not designed | P0 | There must be exactly one skill that can write to Redis — the skill that provisions other skills, updates the registry, and modifies shard authorization. This skill is the most sensitive in the entire system. If compromised, the entire trust boundary collapses — an attacker could provision any skill, authorize any shard. Needs: strict definition of what the skill can write, who can invoke it (brain only? user approval required?), audit trail for every Redis write, and whether it requires a separate credential tier above standard shard auth. | Platform/Security + Kernel |

## Exit Criteria

The gap register is considered materially closed for v1 when:

- All `P0` items are resolved with committed schema/interface docs.
- At least one integration test covers retry + dedupe + commit safety path.
- Reconciliation behavior is deterministic and documented.
- Observability provides traceable lifecycle from event ingress to final response.
