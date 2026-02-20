# Skyra Task Formation System

## 1. Overview

The Task Formation System converts incoming events from voice or chat into structured task objects and execution-ready plans.

Planning and execution are carried out by the same LLM session/context in the control plane.

Formation outputs exactly one of:

- no task (ephemeral reply only)
- ephemeral task (non-mutating work)
- stateful task (mutating work that must be committed)

Scope of this document:

- event-to-task pipeline
- formation artifacts and task object contract
- hand-off boundary to scheduler

Out of scope:

- full scheduler logic
- full executor logic
- full memory internals
- final estimator behavior

## 2. Terminology

- **Event**: incoming user action or model proposal.
- **Ephemeral task**: work that returns a result without changing project state.
- **Stateful task**: work that produces a patch and modifies project state via commit.
- **WorkPlan**: artifact for ephemeral execution.
- **TaskSheet**: structured stateful execution plan.
- **Patch**: concrete state diff for object-store commit.
- **Systems affected**: list of subsystems/resources touched by the task.

## 3. Task Types (Ephemeral vs Stateful)

### Ephemeral task

- used for retrieval, analysis, summarization, research, planning support
- no commit to object store
- result-only output

Artifact:

- `WorkPlan`

### Stateful task

- used when request implies project state mutation
- must produce patch and commit into project history
- must preserve provenance through commit metadata

Artifacts:

- `TaskSheet`
- `Patch`

## 4. WorkPlan vs TaskSheet

### WorkPlan (ephemeral)

Recommended fields:

- `goal`
- `inputs`
- `steps[]`
- `tools_required[]`
- `expected_output`
- `timeout_hint_ms`
- `risk_flags[]`

### TaskSheet (stateful)

Recommended fields:

- `project_id`
- `intent`
- `state_targets[]`
- `constraints[]`
- `validation_checks[]`
- `assumptions[]`
- `evidence[]` (citations for validated assumptions)
- `commit_message_hint`
- `requires_patch=true`

Patch requirements:

- stateful path must emit patch payload or patch-generation request
- patch must be attributable to source event and task id

## 5. Domain Routing

Routing maps an event to the most relevant project/domain.

Routing inputs:

- event text and metadata
- session hints
- project registry
- vector search over derived project state

Routing outputs:

- `domain_id`
- `project_id` (if resolved)
- `routing_confidence`
- `top_candidates[]`

If confidence is low, formation enters ambiguity handling (see Section 11).

## 6. Domain Expert Role

The domain expert is the first-pass task shaper.

Canonical Domain Expert spec:

- `docs/arch/v1/domain-expert/README.md`

Responsibilities:

- decide no-task vs ephemeral vs stateful
- identify `systems_affected`
- create `WorkPlan` or `TaskSheet`
- annotate assumptions and confidence
- validate critical assumptions with tools when needed
- include citations in TaskSheet evidence when external/docs lookup is used

Expected output contract:

- deterministic, machine-readable artifact
- explicit task type
- clear downstream execution intent

## 7. Cross-System Dependency Handling

If a task touches multiple systems, run a system-expert dependency pass.

When to run:

- multiple tools/services are required
- dependency ordering is non-trivial
- stateful change depends on retrieval/analysis pre-step

Output additions:

- dependency ordering
- required preconditions
- inter-system constraints
- rollback hints (stateful path)

## 8. Threshold-Based Review

Optional review by a larger model for complex/high-risk formations.

Typical triggers:

- low routing confidence
- broad stateful patch scope
- multi-system dependency count above threshold
- ambiguous intent wording

Review outcomes:

- approve
- revise
- clarification required
- reject task formation

## 9. Task Object Creation

After expert/review stages, the system creates a canonical task object.

Required fields:

- `task_id`
- `source_event_id`
- `task_type` (`ephemeral` or `stateful`)
- `project_id` (nullable for some ephemeral tasks)
- `systems_affected[]`
- `artifact_ref` (WorkPlan or TaskSheet)
- `patch_ref` (stateful only)
- `formation_confidence`
- `created_at`

Idempotency:

- task creation must be idempotent for duplicate source events
- use stable key strategy (for example `source_event_id + task_variant`)

## 10. Queue, Scheduler, and Execution Start

Task formation and execution run in the same assigned LLM context after queueing and scheduler lane assignment.

Execution model:

- task formation and execution are performed by the same LLM context
- approval gates and queueing may pause/resume work, but do not require switching to a different LLM agent

Approval gate:

- for plan-backed execution, scheduler hand-off is blocked until user approval is received
- control plane sends `PLAN_APPROVAL_REQUIRED` with plan summary and confidence/evidence
- user response contract: `APPROVE | REVISE | CANCEL`
- only `APPROVE` advances to execution
- `REVISE` returns task to formation refinement
- `CANCEL` terminates task before execution

Planner event emission:

- during formation, planner/domain expert may emit user-facing events via context engine
- planner may invoke these via interaction tools (for example `ask_user`)
- allowed planning events: `CLARIFY`, optional `PLAN_PROGRESS`, `PLAN_APPROVAL_REQUIRED`
- events must be persisted on the event/context timeline before delivery
- execution remains blocked until approval event resolves to `APPROVE`

Execution start guarantees:

- task type is explicit
- artifact is structured and validated
- confidence and ambiguity flags are included

Important boundary:

- scheduler v1 is intentionally simple (single queue + lane assignment)
- estimator is only one scheduler component, not the scheduler itself
- formation does not decide final scheduling policy

## 11. Failure and Ambiguity Handling

No domain match:

- emit clarification request or conservative ephemeral path

Conflicting system dependencies:

- run system-expert pass
- if unresolved, hold and request clarification

Insufficient context for stateful mutation:

- do not emit stateful task
- fall back to ephemeral investigation or clarification

Patch invalid/unavailable:

- retain TaskSheet
- mark patch generation pending/failed
- return to review path

## 12. Assumption Validation and Citation Policy

Canonical policy reference:

- `docs/arch/v1/domain-expert/README.md`

For stateful or high-impact plans, Domain Expert must validate critical assumptions before final TaskSheet emission.

Reprompt directive (internal):

- "Validate assumptions with documentation or internet sources."
- "Cite sources in TaskSheet evidence."

TaskSheet evidence minimums:

- `claim`: assumption being validated
- `status`: `validated | unresolved | conflicting`
- `sources[]`: each with `title`, `url`, `retrieved_at`
- `notes`: short interpretation of relevance

If evidence is insufficient within attempt budget:

- emit `CLARIFY`, or
- emit `WorkPlan` with unresolved assumptions explicitly marked

Duplicate source events:

- deduplicate at task object creation boundary

## 13. Future Work (Brief)

- formal task schema versioning
- stronger ambiguity scoring
- integration with finalized scheduler policies
- richer review policies from production telemetry

## 14. Related Docs

- Executor runtime design (draft): `docs/arch/v1/executor.md`

## 15. Appendix A: Estimator Documentation Agent Prompt

Use this prompt when generating the Estimator design documentation for Skyra Task Formation.

```text
Write a comprehensive design document for the Estimator component of Skyra's Task Formation System, based on the following specifications. The Estimator is responsible for predicting task duration and resource needs, and for dynamically updating those predictions during execution. The job execution layer is not yet fully defined, so define interfaces abstractly.

Context:
Skyra processes events into tasks via a Task Formation pipeline (see attached task-formation.md). After the Domain Expert creates a WorkPlan or TaskSheet, the Estimator produces initial estimates. During task execution, the Estimator receives progress updates and refines estimates, which are used for user communication and scheduler hints.

Requirements:

Purpose

Predict how long a task will take (initial estimate).

Classify tasks into duration classes (instant, short, long, unknown) to guide scheduling and user interaction.

Suggest checkpoint intervals for long tasks.

Provide resource hints (GPU, network, etc.) to the scheduler.

Dynamically re-estimate remaining time during execution based on progress.

Inputs – Initial Estimation

Hydrated job (after Domain Expert) including WorkPlan/TaskSheet, project context, user ID, etc.

Snapshot of current system resources (GPU load, memory, network latency).

Historical data: embeddings of similar past tasks with their actual durations and resource usage.

Inputs – Dynamic Re‑estimation (during execution)

Progress snapshots from the executor (abstractly defined): e.g., elapsed time, steps completed, current step, partial results, resource usage, errors.

Original task features and initial estimate.

Outputs

Initial: duration class, estimated seconds (with confidence), checkpoint interval, resource hints, complexity score.

Re‑estimation: updated remaining seconds, new confidence, optionally a reason for change.

All outputs may be used by scheduler, threshold review, and user notification system.

Learning & Adaptation

Store features and actual outcomes (duration, resource consumption) for every completed task.

Periodically retrain a model (e.g., gradient boosting, small neural net) to improve accuracy.

Include progress snapshots in training to improve re‑estimation.

Cold-start fallback rules until enough data exists.

Architecture & Integration

The Estimator is a service on the Mac mini, exposed via an internal API.

Initial estimation occurs after Domain Expert, before threshold review (if any).

During execution, the executor (to be defined) sends progress updates to the Estimator; the Estimator returns updated estimates.

Estimates are attached to the task object and can be queried by the scheduler or notification system.

Open Points

The exact executor interface is TBD; define the expected progress snapshot schema.

How frequently re‑estimation occurs (e.g., every 30 seconds, after each step) is configurable and may depend on duration class.

Notification logic (when to inform the user) is outside this doc but should reference the estimator's outputs.
```

## Task Formation Flow Diagram

```text
Event (voice/chat)
   |
   v
[Domain Routing]
   |
   v
[Domain Expert]
  |---- no_task --------> reply-only path
  |
  +---- ephemeral ------> WorkPlan
  |
  +---- stateful -------> TaskSheet + Patch
               |
               v
      [Optional Threshold Review]
               |
               v
        [Task Object Creation]
               |
               v
         [Scheduler Hand-off]
  (estimator contributes, but is not the scheduler)
```
