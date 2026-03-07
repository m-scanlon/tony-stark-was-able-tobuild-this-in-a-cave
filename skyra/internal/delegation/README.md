# Delegation Engine

This package tree hosts delegation-time decisioning components used by the control-plane orchestrator.

## Subpackages

- `estimator/`: placement decisioning — reads complexity score from estimation call output, matches to best available shard via capability profiles
- `telemetry/`: request/outcome telemetry contracts and rolling stats interfaces
- `api/`: HTTP handlers/contracts for estimator and telemetry endpoints

## Estimator Role

The Estimator is **an inference call, not a service**. It fires when the External Router picks up an estimation work item from the heap. The External Router owns the heap — priority ordering, preemption, and dispatch. When an estimation item is dequeued, a prompt runs. That prompt is the Estimator.

Input — estimation output from the domain agent:

```json
{
  "is_job": true,
  "complexity": 3,
  "reasoning_depth": 2,
  "cross_domain": false,
  "reversible": true,
  "output_scope": "fact",
  "domain": "servers"
}
```

Complexity is measured in estimated tool calls. The Estimator reads this, checks current shard capability profiles and load, and decides:

- `complexity ≤ 1` → execute inline immediately. The Estimator does the work itself — no job formed, no heap placement.
- `complexity > 1` → place job onto the heap targeting the best available shard.

Placement ranges are illustrative. Actual routing uses capability profiles, not hardcoded rules.

Because the Estimator is a prompt, estimation quality improves as model quality improves — no code changes required.

## Related Modules

- `../taskformation/`: event-to-task formation pipeline
- `../../docs/arch/v1/task-formation.md`: architecture/design reference for task formation
- `../../docs/arch/v1/scheduler.md`: unified heap, inference types, complexity scoring, External Router heap ownership
