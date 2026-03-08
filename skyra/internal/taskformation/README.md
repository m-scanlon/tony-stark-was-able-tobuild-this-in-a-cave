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
- `skyra/internal/agent/README.md` — Agent Service, object store, tool registry

Notes:

- Task formation stops at scheduler hand-off.
- Estimation is only one scheduler input and is not implemented here.
- Task artifacts (WorkPlan/TaskSheet) are persisted in the object store under `.skyra/agents/{agent_id}/jobs/{job_id}/tasks/{task_id}/`.
- Task formation operates within a two-layer tool system: global tools (always present) and local tools (files under `tools/` in the agent's git repo, discovered by the LLM walking the filesystem during execution).
- Boundary enforcement: the BoundaryValidator (pure code check) runs before any tool dispatch, joining the tool's `categories[]` against the agent boundary in `state.json`. Locked tools trigger a permission prompt before execution.
