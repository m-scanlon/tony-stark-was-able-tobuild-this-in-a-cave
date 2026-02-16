# Skyra Task Formation System

## 1. Overview

The Task Formation System converts incoming events from voice or chat into structured task objects for downstream scheduling.

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

Responsibilities:

- decide no-task vs ephemeral vs stateful
- identify `systems_affected`
- create `WorkPlan` or `TaskSheet`
- annotate assumptions and confidence

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

## 10. Hand-off to Scheduler

Task Formation hands a canonical task object to scheduler intake.

Handoff guarantees:

- task type is explicit
- artifact is structured and validated
- confidence and ambiguity flags are included

Important boundary:

- scheduler is not fully designed yet
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

Duplicate source events:

- deduplicate at task object creation boundary

## 12. Future Work (Brief)

- formal task schema versioning
- stronger ambiguity scoring
- integration with finalized scheduler policies
- richer review policies from production telemetry

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
