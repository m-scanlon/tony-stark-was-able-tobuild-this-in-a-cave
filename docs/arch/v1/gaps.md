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
| G3 | Stateful commit safety boundary is underdefined | P0 | Risk of corrupt/partial state changes without strict pre-commit checks and conflict handling. | Task Formation + Memory/Object Store |
| G4 | Reconciliation UX policy is not fully specified | P1 | Split-brain user experience if Pi provisional speech conflicts with Mac final answer. | Voice Node + Control Plane |
| G5 | AuthN/AuthZ model is not concretely implemented | P0 | Device/agent channels remain security-sensitive attack surface. | Platform/Security |
| G6 | Backpressure and overload policies are undefined | P1 | Queue growth and latency spikes can collapse responsiveness under load. | Ingress + Scheduler |
| G7 | Executor/Resource Manager contracts are abstract | P1 | Stage execution decisions cannot be implemented predictably without exact API contracts. | Executor + Resource Manager |
| G8 | Observability and SLOs are not defined end-to-end | P1 | Hard to debug reliability and latency regressions without traceability and targets. | Platform/Operations |
| G9 | Degradation and cold-start behavior is incomplete | P1 | Unclear runtime behavior when STT/model/estimator/context systems fail or lag. | Voice Node + Orchestrator |
| G10 | Data lifecycle and retention policy is missing | P1 | Inbox/checkpoints/transcripts can grow unbounded and violate privacy expectations. | Memory + Platform |

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
   - Formalize provisional phrasing policy on Pi.
   - Define reconciliation rules for `UPDATE`, `FINAL`, `CLARIFY`, `ERROR`.
   - Add explicit voice UX tests for contradiction scenarios.

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
