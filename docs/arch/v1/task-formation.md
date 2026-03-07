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

- `WorkPlan` — persisted at `.skyra/agents/{agent_id}/jobs/{job_id}/tasks/{task_id}/workplan.json`

### Stateful task

- used when request implies project state mutation
- must produce patch and commit into project history
- must preserve provenance through commit metadata

Artifacts:

- `TaskSheet` — persisted at `.skyra/agents/{agent_id}/jobs/{job_id}/tasks/{task_id}/tasksheet.json`
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

- `agent_id`
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

## 5. Domain Routing — Domain Agent as Doorkeeper

There is no central classifier. Domain agents self-select relevance.

The front face transformer reads the context blob — which contains **all registered agents with their relevance scores** — and labels the turn as in-domain or "other." For in-domain turns, it routes to the relevant domain agents. Each agent receives the full context blob and makes its own decision about whether the turn belongs to it.

The domain agent is the doorkeeper of its own domain. No external system knows a domain better than the agent that owns it.

**"Other" turns** — turns that don't clearly fit any current agent — are stored in RDS and picked up by the nightly batch process. Every agent runs against accumulated session context at night, so nothing is permanently missed. A turn deposited into a domain that doesn't quite fit yet is acceptable; the V3 background process will detect accumulating patterns and propose new agents over time.

**Routing inputs (front face transformer):**
- turn transcript
- context blob (all agents + relevance scores)
- importance score assigned at ingress

**Domain agent self-selection:**
- receives full context blob
- decides if turn is relevant to its domain
- checks whether turn impacts an ongoing job
- forms a job if warranted (see Section 5a)

If a turn spans multiple domains, multiple agents can each self-select independently. No special-casing needed.

## 5a. Estimation Call

After the domain agent receives the turn and determines a job is needed, it produces an **estimation call** — the first inference call for any actionable turn.

Output:

```json
{
  "is_job": true,
  "complexity": 3,
  "domain": "servers"
}
```

Complexity is measured in **estimated tool calls**.

- `complexity ≤ 1` → execute inline immediately. Never enters the heap.
- `complexity > 1` → form job, push to heap. Estimator routes to best available machine.

The complexity threshold is currently **1** and will be tuned from real usage data. See `docs/arch/v1/scheduler.md` for full heap and inference type design.

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

Note: before local tools are returned to the Domain Expert, the Agent Service runs a hydration step — each tool is enriched with an `access` field derived from the agent boundary in `state.json`. The Domain Expert receives all retrieved tools, including locked ones, with their access status attached. Locked tools that the LLM proposes calling are caught by the BoundaryValidator at runtime before execution. See `skyra/internal/agent/README.md` for the full hydration and enforcement model.

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
- `agent_id` (nullable for some ephemeral tasks)
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

Note on `requires_approval` vs `PLAN_APPROVAL_REQUIRED`:

These are two distinct concepts and must not be confused.

- `PLAN_APPROVAL_REQUIRED` — a plan-level gate. The entire plan waits for user approval before any execution begins. This is what is described above.
- `requires_approval` on a local tool — a tool-level flag in the tool registry. It means the tool is surfaced and highlighted to the user during plan review so they can see it clearly. It does NOT pause execution mid-run. The user approves the full plan once and execution continues uninterrupted.

Domain Expert event emission:

- during formation, Domain Expert may emit user-facing events via context engine
- Domain Expert may invoke these via interaction tools (for example `ask_user`)
- allowed planning events: `CLARIFY`, optional `PLAN_PROGRESS`, `PLAN_APPROVAL_REQUIRED`
- events must be persisted on the event/context timeline before delivery
- execution remains blocked until approval event resolves to `APPROVE`

Execution start guarantees:

- task type is explicit
- artifact is structured and validated
- confidence and ambiguity flags are included

Important boundary:

- scheduling is handled by the unified max-heap — see `docs/arch/v1/scheduler.md`
- the Estimator reads the estimation call output and makes placement decisions
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
- Scheduler — unified heap, inference types, complexity scoring: `docs/arch/v1/scheduler.md`
- Agent Service (object store, commits, tool registry): `skyra/internal/agent/README.md`
- Delegation Engine (Estimator): `skyra/internal/delegation/README.md`

## 15. Estimator Role (Updated)

> Note: The Estimator's role changed significantly with the architecture revision. The old prompt-based design spec below this section is superseded. See `skyra/internal/delegation/README.md` for the current Estimator design.

The Estimator now has a single, clear responsibility: **placement**. It reads the estimation call output from the domain agent and routes the job to the best available machine.

Input:
```json
{
  "is_job": true,
  "complexity": 3,
  "domain": "servers"
}
```

Complexity in estimated tool calls is the primary placement signal. The Estimator matches against registered shard capability profiles and current load. No duration prediction, no checkpoint intervals, no duration classes — those were part of the old design that has been retired.

## Task Formation Flow Diagram

```text
Event (voice/chat)
   |
   v
[Front Face Transformer]
  (labels turn: in-domain | other, reads all agents + relevance scores)
   |
   +---- other ---------> RDS (batch picks up at night)
   |
   v
[Domain Agent] (self-selects — doorkeeper of its own domain)
   |
   v
[Estimation Call]
  (is_job? complexity in tool calls?)
   |
   +---- complexity ≤ 1 --> execute inline immediately
   |
   +---- complexity > 1 --> [Heap] → Estimator routes to capable machine
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
                                            [Execution]
```
