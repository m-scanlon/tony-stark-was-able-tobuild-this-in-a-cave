# Task Formation

This document defines the planning artifact boundary before execution.

## Overview

Task formation shapes a routed request into one of:

- `no task`
- `WorkPlan` (ephemeral)
- `TaskSheet + Patch` (stateful)

Planning and execution run in the same LLM session. Formation is the plan boundary, not a separate runtime.

## Inputs

- routed user request/event
- context blob
- domain self-selection signal
- tool policies and boundary constraints

## Outputs

- deterministic planning artifact
- confidence/evidence metadata
- optional planning events: `CLARIFY`, `PLAN_PROGRESS`, `PLAN_APPROVAL_REQUIRED`

## Plan Artifacts

### `no task`

Reply-only path. No stateful execution artifact.

### `WorkPlan`

Ephemeral plan. Executes without stateful patch semantics.

### `TaskSheet + Patch`

Stateful plan with explicit patch/commit implications and stronger validation gates.

## Approval Gate

`PLAN_APPROVAL_REQUIRED` is a plan-level execution gate:

- `APPROVE` -> execute
- `REVISE` -> return to formation
- `CANCEL` -> terminate before execution

`requires_approval` on a tool is a plan review display hint only. It is not a separate execution gate.

## Related

- `docs/arch/v1/domain-expert/README.md`
- `docs/arch/v1/lifecycle.md`
- `docs/arch/v1/kernel.md`

