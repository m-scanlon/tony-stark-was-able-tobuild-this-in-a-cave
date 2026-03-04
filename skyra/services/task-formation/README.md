# Task Formation Module

This package tree defines the event-to-task formation boundary inside the control plane.

Purpose:

- convert inbound events into one of:
  - no task
  - ephemeral task (`WorkPlan`)
  - stateful task (`TaskSheet` + `Patch`)
- produce canonical task objects for scheduler intake

Design reference:

- `docs/arch/v1/task-formation.md`
- `skyra/internal/project/README.md` — project service, object store, tool registry

Notes:

- Task formation stops at scheduler hand-off.
- Estimation is only one scheduler input and is not implemented here.
- Task artifacts (WorkPlan/TaskSheet) are persisted in the object store under `.skyra/agents/{agent_id}/jobs/{job_id}/tasks/{task_id}/`.
- Task formation operates within a two-layer tool system: global tools (always present) and local tools (retrieved per request via vector search from the agent's tool registry).
- Boundary enforcement runs at two code layers: (1) the Project Service hydrates retrieved tools with access status — locked tools are visible to the LLM but marked; (2) the Executor's BoundaryValidator prompts the user before any locked tool executes.
