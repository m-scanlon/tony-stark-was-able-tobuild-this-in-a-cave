# Delegation Engine

This package tree hosts delegation-time decisioning components used by the control-plane orchestrator.

## Subpackages

- `estimator/`: placement decisioning — reads complexity score from estimation call output, matches to best available shard via capability profiles
- `telemetry/`: request/outcome telemetry contracts and rolling stats interfaces
- `api/`: HTTP handlers/contracts for estimator and telemetry endpoints

## Estimator Role (Updated)

The Estimator's responsibility is **placement**. It reads the estimation call output produced by the domain agent:

```json
{
  "is_job": true,
  "complexity": 3,
  "domain": "servers"
}
```

Complexity is measured in estimated tool calls. The Estimator matches this against registered shard capability profiles and current load, then assigns the job to the best available machine.

- Complexity ≤ 1 → execute inline (never reaches the Estimator)
- Complexity 2–5 → Mac mini class
- Complexity 6+ → GPU machine or most capable available shard

Placement ranges are illustrative. Actual routing uses capability profiles, not hardcoded rules.

The Estimator no longer reads a complex `job_envelope_v1` assembled by the Internal Router. Complexity score in tool calls is the primary placement signal.

## Related Modules

- `../taskformation/`: event-to-task formation pipeline
- `../../docs/arch/v1/task-formation.md`: architecture/design reference for task formation
- `../../docs/arch/v1/scheduler.md`: unified heap, inference types, complexity scoring

## Notes

- This is intentionally in-process with orchestration for now.
- Package boundaries are designed so the estimator can be split into a standalone service later if needed.
