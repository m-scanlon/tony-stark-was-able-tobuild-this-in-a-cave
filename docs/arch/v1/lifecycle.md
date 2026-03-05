# Task Lifecycle — End to End

Every request that enters Skyra moves through the same pipeline, regardless of complexity. The plan artifact determines the weight of execution — not the route.

---

## Stage 1 — Ingress

`voice_event_v1` arrives at the Brain Shard.

- Brain Shard generates `event_id` (ULID) on receipt
- Event persisted to JSONL
- `(session_id, turn_id)` is the deduplication key — `turn_id` is the Voice Shard's reference, `event_id` is the Brain Shard's
- Event handed off to the ingest queue

---

## Stage 2 — Domain Routing

The event is resolved to an agent.

Inputs:
- event text and metadata
- session hints from `voice_event_v1.session_state`
- agent registry
- vector search over agent state

Outputs:
- `agent_id`
- `routing_confidence`
- `top_candidates[]`

Low confidence → ambiguity handling (clarification request or conservative ephemeral path).

---

## Stage 3 — Domain Expert (Planning)

The Domain Expert takes the routed event, retrieves context and tools, and decides what kind of task this is.

Outputs exactly one of:
- `no task` → reply-only path, done
- `WorkPlan` → ephemeral, no state commit
- `TaskSheet + Patch` → stateful, requires commit

Tool system:
- **Global tools** — always injected (agent state operations, propose commit, etc.)
- **Local tools** — retrieved via vector search, hydrated with access status by Agent Service before being handed to the LLM

Planning events the Domain Expert may emit:
- `CLARIFY` — missing information blocks safe planning
- `PLAN_PROGRESS` — optional progress signal
- `PLAN_APPROVAL_REQUIRED` — final pre-execution gate

---

## Stage 4 — Optional Threshold Review

A larger model reviews the plan if complexity or risk thresholds are crossed.

Triggers:
- low routing confidence
- broad stateful patch scope
- multi-system dependency count above threshold
- ambiguous intent

Outcomes: `approve | revise | clarify | reject`

---

## Stage 5 — Task Object Creation

A canonical task object is stamped after expert/review stages.

Required fields:
- `task_id`
- `source_event_id`
- `task_type` (`ephemeral` or `stateful`)
- `agent_id`
- `systems_affected[]`
- `artifact_ref` (WorkPlan or TaskSheet)
- `patch_ref` (stateful only)
- `formation_confidence`
- `created_at`

Task creation is idempotent — duplicate source events produce the same task object.

Artifacts persisted at: `.skyra/agents/{agent_id}/jobs/{job_id}/tasks/{task_id}/`

---

## Stage 6 — Plan Approval Gate

Execution is blocked until the user explicitly approves.

Flow:
1. Control plane sends `PLAN_APPROVAL_REQUIRED` with plan summary and confidence/evidence
2. User responds: `APPROVE | REVISE | CANCEL`
3. Only `APPROVE` advances to execution
4. `REVISE` returns the task to Domain Expert refinement
5. `CANCEL` terminates the task before execution

Note: `PLAN_APPROVAL_REQUIRED` (this gate) and `requires_approval` on a local tool are distinct.
- `PLAN_APPROVAL_REQUIRED` — plan-level gate, all execution waits
- `requires_approval` on a tool — display flag only, highlights the tool during plan review, does not add a separate approval step

---

## Stage 7 — Scheduler

The approved task enters the scheduler queue.

- Lane assigned: `fast_local | deep_reasoning`
- Jobs table tracks operational status: `queued → running → completed | failed`
- Semantic phases (planning / executing / validating / replanning / done) are tracked separately in the task artifact

---

## Stage 8 — Executor

The same LLM session that planned the task now executes it. Planning and execution share one context — approval gates and queueing may pause work but do not require a context switch.

Execution loop per step:
1. Pre-check resources (Resource Manager)
2. BoundaryValidator — check all tool calls against agent boundary. Locked tool → permission prompt (`allow_always | allow_once | deny`). `deny` triggers bounded replan.
3. Execute step with specified tools/models
4. Capture outputs + runtime metadata
5. Validate output against criteria
6. Check assumptions
7. Pass → persist checkpoint, proceed
8. Fail → classify severity, attempt local fix, replan if needed

Replan budget: 3 attempts max. Exceed budget → escalate or halt.

Progress snapshots sent to Estimator at each step boundary (and on a heartbeat for long-running steps).

---

## Flow Diagram

```
voice_event_v1
      |
      v
  [Ingress] — generates event_id, persists, queues
      |
      v
  [Domain Routing] — resolves agent_id
      |
      v
  [Domain Expert] — no_task | WorkPlan | TaskSheet+Patch
      |
      +-- no_task ---------> reply-only, done
      |
      v
  [Threshold Review] — optional, complexity/risk gate
      |
      v
  [Task Object Creation] — canonical task stamped
      |
      v
  [Plan Approval Gate] — APPROVE | REVISE | CANCEL
      |
      v
  [Scheduler] — lane assigned, queued
      |
      v
  [Executor] — step-by-step, validate, checkpoint, replan
      |
      v
     done
```

---

## Related Docs

- `docs/arch/v1/scyra.md` — full system architecture
- `docs/arch/v1/task-formation.md` — task formation detail
- `docs/arch/v1/executor.md` — executor design
- `docs/arch/v1/domain-expert/README.md` — planning phase
- `skyra/internal/agent/README.md` — agent service, tool hydration, boundary enforcement
- `skyra/internal/scheduler/README.md` — scheduler, job lifecycle, lanes
- `skyra/schemas/ingress/voice/` — voice_event schema
- `docs/arch/v1/next-steps.md` — open questions on job_envelope_v1 and executor loop
