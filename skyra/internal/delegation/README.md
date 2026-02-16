# Delegation Engine

This package tree hosts delegation-time decisioning components used by the control-plane orchestrator.

## Subpackages

- `estimator/`: unsupervised task estimation interfaces and model contracts
- `telemetry/`: request/outcome telemetry contracts and rolling stats interfaces
- `api/`: HTTP handlers/contracts for estimator and telemetry endpoints

## Notes

- This is intentionally in-process with orchestration for now.
- Package boundaries are designed so the estimator can be split into a standalone service later if needed.
