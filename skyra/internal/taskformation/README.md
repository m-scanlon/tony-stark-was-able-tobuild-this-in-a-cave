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

Notes:

- Task formation stops at scheduler hand-off.
- Estimation is only one scheduler input and is not implemented here.
