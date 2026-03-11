# Scheduler (Kernel Internal)

> **Architecture note**: there is no standalone Scheduler service. Scheduling is a kernel-internal heap model.

This doc exists as the canonical reference for heap behavior previously described as "scheduler."

## Core Model

All work is placed on one unified max-heap and ordered by importance score.

Work item classes:

- `estimation` — very high priority
- `job` — high priority
- `batch` — very low priority

## Heap Semantics

- One queue model for all work.
- Highest-scored item is always selected next.
- Work re-enters the heap after each tool call.
- Preemption happens naturally at tool-call boundaries.

## Placement

Estimator reads the work item's complexity and routes execution to a shard with compatible capabilities.

- Complexity `<= 1` can execute inline.
- Complexity `> 1` is heap-managed and placed to the best available shard.

## Job Registry Relationship

The Job Registry is passive state tracking, not scheduling logic.

- Scheduler/heap decides ordering and dispatch timing.
- Registry records lifecycle transitions.

## Related

- `docs/arch/v1/kernel.md`
- `docs/arch/v1/executor.md`
- `docs/arch/v1/next-steps.md`
- `docs/arch/v1/dataflow-walk-notes.md`

