# Skyra Domain Expert (v1)

This directory is the source of truth for how the Domain Expert works.

## 1. What It Is

- Domain Expert is a control-plane service module on Mac mini.
- It is not a remote machine agent.
- Task planning and task execution are carried out by the same LLM session/context.
- It shapes an accepted event into one of:
  - `no task`
  - `WorkPlan` (ephemeral)
  - `TaskSheet + Patch` (stateful)

Code references:

- `skyra/internal/taskformation/domain_expert.go`
- `skyra/internal/taskformation/system_expert.go`
- `skyra/internal/taskformation/pipeline.go`
- `skyra/internal/taskformation/factory.go`

## 2. Inputs and Outputs

Inputs:

- user request/event text
- domain routing result (`domain_id`, confidence, candidates)
- context window + memory retrieval results
- tool registry/allowlist + policies

Output:

- deterministic task artifact with confidence:
  - `no task`
  - `WorkPlan`
  - `TaskSheet + Patch`

## 3. Runtime Loop

1. Receive routed event.
2. Retrieve domain-relevant context (semantic + recency + project scope).
3. Identify assumptions and missing evidence.
4. Optionally emit planning events (`CLARIFY`, progress) to context engine/user channel.
5. Run bounded investigation tool calls if needed (docs/history/validation).
6. Decide task type and generate plan artifact.
7. If confidence/evidence is insufficient, emit `CLARIFY`.
8. After approval (when required), continue in the same LLM session to execute the plan via tools.

## 4. Domain Expert Composition

Each domain configuration includes:

- `tool_profile`: allowed tools and argument constraints
- `retrieval_profile`: filters, top-k, rerank, recency policy
- `formation_policy`: thresholds for task type selection
- `schema_contract`: required output fields
- `investigation_tools`: read/research tools available during planning
- `interaction_tools`: user-interaction tools (for example `ask_user`)
- `attempt_budget`: max retries/time budget before fallback

## 5. Tool Reprompt Policy

When assumptions are not grounded in retrieved context, Domain Expert must reprompt itself:

- "Validate assumptions with documentation or internet sources."
- "Return citations in TaskSheet evidence."

Before emitting final `TaskSheet`:

- each critical assumption is `validated`, `unresolved`, or `conflicting`
- citations include `title`, `url`, `retrieved_at`
- unresolved critical assumptions trigger `CLARIFY` or downgrade to `WorkPlan`

Reprompt control loop (bounded, non-recursive):

- loop while confidence/evidence gates are not met
- stop if `max_attempts` or `max_time_ms` budget is reached
- stop early on stagnation (confidence gain below `delta_min` for 2 rounds)
- require confidence threshold (for example `>= 0.80`) and evidence quality threshold
- if thresholds are not met within budget, emit `CLARIFY`

`ask_user` as a planning tool:

- From Domain Expert perspective, clarification is invoked via `ask_user` tool call.
- `ask_user` emits a `CLARIFY` planning event through context engine/user channel.
- Use `ask_user` when missing data blocks safe planning or evidence validation.

## 6. TaskSheet Evidence Contract

For assumption-backed planning, include:

- `assumptions[]`
- `evidence[]`

Each evidence item:

- `claim`
- `status`: `validated | unresolved | conflicting`
- `sources[]`: `title`, `url`, `retrieved_at`
- `notes`

## 7. Boundaries

- Domain Expert may investigate and plan.
- Domain Expert may also execute approved plan steps in the same LLM session.
- Domain Expert does not directly commit memory/project state.
- State commits happen downstream through canonical execution + commit flow.
- Pi remains non-authoritative and only renders backend-authored outputs.

## 8. Plan Approval Gate

After planning is complete, execution is paused behind an explicit user approval step.

Flow:

1. Domain Expert emits plan artifact (`WorkPlan` or `TaskSheet+Patch`) with confidence/evidence.
2. Control plane sends `PLAN_APPROVAL_REQUIRED` notification to user.
3. User responds `APPROVE | REVISE | CANCEL`.
4. Only `APPROVE` allows queued job execution to continue in the same assigned LLM session.
5. `REVISE` returns artifact to Domain Expert refinement loop.
6. `CANCEL` closes the job without execution.

Planning-event contract:

- Domain Expert can issue user-facing planning events before final artifact:
  - `CLARIFY` (missing information)
  - `PLAN_PROGRESS` (optional)
  - `PLAN_APPROVAL_REQUIRED` (final pre-execution gate)
- These events are written through the context engine/event bus so they are persisted in turn history.

Tool mapping:

- `ask_user` -> `CLARIFY`
- `notify_progress` -> `PLAN_PROGRESS`
- `request_plan_approval` -> `PLAN_APPROVAL_REQUIRED`

## 9. Related Docs

- `docs/arch/v1/task-formation.md`
- `docs/arch/v1/agents-services.md`
- `docs/arch/v1/scyra.md`
