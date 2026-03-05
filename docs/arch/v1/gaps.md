# Skyra Architecture Gaps Register

This document tracks known architecture gaps that must be resolved before production-hardening.

## Priority Scale

- `P0`: blocks safe/reliable production behavior
- `P1`: high impact, should be resolved early
- `P2`: important, can follow core stabilization

## Gaps

| ID | Gap | Priority | Why it matters | Suggested owner |
| --- | --- | --- | --- | --- |
| G1 | `job_envelope_v1` schema not fully locked | P0 | Cross-node compatibility and routing correctness depend on one canonical contract. | Control Plane |
| G2 | End-to-end idempotency beyond ingress is incomplete | P0 | Retries can create duplicate tasks, duplicate side effects, or duplicate commits. | Ingress + Task Formation + Executor |
| G3 | Stateful commit safety boundary is underdefined | P0 | Risk of corrupt/partial state changes without strict pre-commit checks and conflict handling. Two-level status (operational in scheduler jobs table vs semantic in tasksheet) must be coordinated on failure/rollback. | Task Formation + Project Service |
| G4 | ~~Reconciliation UX policy is not fully specified~~ **CLOSED for v1** | P1 | **v1 decision**: Provisional responses cut from v1. Voice Shard emits non-semantic ACKs only (earcon/LED/short wait phrase). Voice Shard renders only Brain Shard-authored messages (`UPDATE`, `PLAN_PROGRESS`, `CLARIFY`, `PLAN_APPROVAL_REQUIRED`, `FINAL`, `ERROR`). No split-brain risk in v1. Provisional response path (front-door model speaks before Brain Shard responds, then reconciles on `FINAL`) deferred to v2 — see note in `docs/arch/v1/scyra.md` section 7.1. | Voice Shard + Control Plane |
| G5 | AuthN/AuthZ model is not concretely implemented | P0 | Device/agent channels remain security-sensitive attack surface. | Platform/Security |
| G6 | Backpressure and overload policies are undefined | P1 | Queue growth and latency spikes can collapse responsiveness under load. | Ingress + Scheduler |
| G7 | Executor/Resource Manager contracts are abstract | P1 | Stage execution decisions cannot be implemented predictably without exact API contracts. | Executor + Resource Manager |
| G8 | Observability and SLOs are not defined end-to-end | P1 | Hard to debug reliability and latency regressions without traceability and targets. | Platform/Operations |
| G9 | Degradation and cold-start behavior is incomplete | P1 | Unclear runtime behavior when STT/model/estimator/context systems fail or lag. | Voice Shard + Orchestrator |
| G10 | Data lifecycle and retention policy is missing | P1 | Inbox/checkpoints/transcripts can grow unbounded and violate privacy expectations. | Memory + Platform |
| G11 | Agent Registry and Scheduler Jobs table schemas not locked | P0 | Two separate SQLite tables. Schemas must be finalized before any service can be built against them. | Agent Service + Scheduler |
| G12 | ~~Agent state (state.json) four-section schema not locked~~ **CLOSED** | P0 | Schema locked: `metadata/knowledge/artifact/boundary`. Boundary carries `allowed_tool_categories`, `denied_tool_patterns`, `restrictions[]` (no enforcement field — all locked tool attempts prompt the user). Enforced in code via two layers: hydration (lock status attached to tools before LLM sees them) and BoundaryValidator (permission prompt at runtime). See `skyra/internal/agent/README.md`. | Agent Service |
| G13 | ~~Tool system two-layer contract not finalized~~ **CLOSED** | P1 | Global vs local tools, `requires_approval` behavior, and local tool `categories[]` field for boundary enforcement are now fully specified. See `skyra/internal/agent/README.md`. | Agent Service + Domain Expert |
| G14 | Voice response channel (`/v1/voice`) not implemented | P0 | `/v1/voice` returns a literal string, not a `voice_result_v1` protocol message. No Brain Shard response can reach the user over the voice path. Blocks all voice UX end-to-end. Needs proper streaming response channel (WebSocket or gRPC) emitting `message_type` in `{UPDATE, PLAN_PROGRESS, CLARIFY, PLAN_APPROVAL_REQUIRED, FINAL, ERROR}`. | Control Plane |
| G15 | Plan approval response channel undefined | P1 | Voice Shard has no mechanism to route a user's spoken approval/rejection (`APPROVE \| REVISE \| CANCEL`) back to a `PLAN_APPROVAL_REQUIRED` job. `session_state.pending_job_id` + `waiting_for` contract is defined in the event schema but not implemented anywhere. A user saying "yes" after a plan approval prompt is indistinguishable from a new request. | Voice Shard + Control Plane |
| G16 | `requires_approval` vs `PLAN_APPROVAL_REQUIRED` semantic distinction is a high-risk confusion point | P1 | These have similar names but completely different execution semantics: `requires_approval` on a tool is a UI hint only and does NOT pause execution; `PLAN_APPROVAL_REQUIRED` is a plan-level gate that pauses before any tool runs. High risk of implementation bugs that either silently execute restricted tools or double-gate all plans. Must be explicitly documented and enforced at code review time. | Domain Expert + Executor + Control Plane |
| G17 | Agent routing staleness — Voice Shard `active_agent` may be stale | P1 | Voice Shard sends `active_agent` from local cache in `context_window`. Brain Shard has no mechanism to verify or correct it. If the user switches agents, Brain Shard routes to the wrong context until Context Injector sync is implemented. **Design locked**: CIX is a state-aware sync daemon — not turn-triggered. It watches for state changes in the context engine and pushes compressed context packages to registered edge devices whenever their cache is stale. CIX tracks per-device last-sync timestamp. `cache_age_seconds` in `voice_event_v2` reflects time since last CIX push, not time since last turn. See `docs/arch/v1/context-engine.md`. Still not implemented. | Context Injector + Voice Shard |
| G18 | Cross-agent write protocol undefined | P1 | `skyra.user` is always injected (read direction solved). The write direction is not designed: if Skyra learns something about the user during a domain agent session, there is no defined protocol for propagating that to `skyra.user` without breaking single-ownership semantics. v1 behavior: manual commit only, Skyra flags insights at session end. v2 direction: structured learning event processed by orchestrator post-session. See `docs/arch/v1/agents/user.md`. | Agent Service + Orchestrator |
| G19 | Context engine design incomplete | P1 | Session, turn, and context package structures are designed (`docs/arch/v1/context-engine.md`) but open questions remain. No locked schema for turn persistence storage, session timeout policy, token budget defaults, or context package caching strategy. Six data sources defined: agent registry, object store, vector DB, turn history, active job context, `skyra.user`. | Context Engine + Agent Service |

## Recommended Next Actions

1. **Lock canonical contracts (`P0`)**
   - Finalize `job_envelope_v1` JSON schema with `schema_version`.
   - Define idempotency keys for `event_id`, `job_id`, `task_id`, and stateful commit operations.
   - Define commit preconditions and conflict behavior for stateful mutations.

2. **Close security and reliability baseline (`P0`)**
   - Choose transport auth strategy (`mTLS` or signed tokens + rotation).
   - Define replay protection and per-device/agent authorization scopes.
   - Define ACK + retry semantics across every boundary (outbox/inbox/task creation/executor side effects).

3. **Finalize execution interfaces (`P1`)**
   - Specify Executor <-> Resource Manager request/response schema.
   - Specify progress snapshot schema for Estimator updates.
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

## Exit Criteria

The gap register is considered materially closed for v1 when:

- All `P0` items are resolved with committed schema/interface docs.
- At least one integration test covers retry + dedupe + commit safety path.
- Reconciliation behavior is deterministic and documented.
- Observability provides traceable lifecycle from event ingress to final response.
