# 6. Executor (Draft)

> Draft status: this design is based on conceptual notes and **requires product owner validation** before implementation lock-in.

## 6.1 Overview

The Executor is the stage runner for Skyra task execution. It consumes `WorkPlan` (ephemeral) or `TaskSheet` (stateful), executes stage-by-stage, validates outcomes, handles assumption drift, and coordinates controlled replanning when needed.

The Executor is not a blind script runner. It is an adaptive runtime loop with:

- stage validation
- assumption checking
- local corrective actions
- bounded replanning
- progress reporting to Estimator
- checkpoint persistence for resume/recovery

## 6.2 Role in Pipeline

Placement in end-to-end flow:

```text
Event -> JobEnvelope v1 -> Task Formation -> Task Object -> Estimator (initial) -> Scheduler -> Executor
                                                                                  ^
                                                                                  |
                                                        progress snapshots --------+
```

Authority boundary:

- Pi front-door can provide provisional responses but does not execute authoritative task pipelines.
- Mac mini control plane owns orchestration authority and coordinates Executor runs.
- GPU and agent nodes are execution targets selected by scheduler/runtime policy.

The Scheduler remains responsible for placement policy. The Executor runs the selected task and reports progress/outcomes.

## 6.3 Core Responsibilities

The Executor must:

1. Execute stage-by-stage with validation.
2. Detect and classify assumption drift severity.
3. Attempt local fixes before replanning.
4. Trigger bounded, incremental replanning via Domain Expert.
5. Send progress snapshots to Estimator.
6. Consult Resource Manager before constrained stages.
7. Persist checkpoints for crash/restart recovery.
8. Emit user-facing lifecycle notifications through Notifier.

## 6.4 Stage Execution Model

Each stage contains:

- `stage_id`
- `goal`
- `inputs`
- `tools_required`
- `expected_output`
- `validation_criteria`
- `resource_hints`
- `timeout_seconds`

Execution loop per stage:

1. Pre-check resources with Resource Manager.
2. Execute stage with specified tools/models.
3. Capture outputs + runtime metadata.
4. Validate output against criteria.
5. Check assumptions and environmental expectations.
6. If pass: persist checkpoint and proceed.
7. If fail: classify severity, attempt local fixes, then replan if needed.

## 6.5 Severity Levels for Assumption Changes

| Level    | Example                                         | Action                                                          |
| -------- | ----------------------------------------------- | --------------------------------------------------------------- |
| Trivial  | API slightly slower than expected               | Log and continue                                                |
| Minor    | Missing optional field with safe default        | Adjust locally and continue                                     |
| Moderate | Required file missing but can be generated      | Attempt local fix; if fix fails, replan affected remaining path |
| Major    | Core assumption false (intent/goal mismatch)    | Pause and replan remaining stages                               |
| Critical | Security risk or impossible/safety-unsafe stage | Halt task, notify user/system, abort                            |

## 6.6 Controlled Adaptation (Loop Prevention)

### Local fixes first

Try low-cost recovery before calling Domain Expert:

- retry with exponential backoff
- fallback tool/model path
- parameter adjustment (batch size/model variant/timeouts)

### Incremental replanning only

On replan request, send:

- original TaskSheet/WorkPlan
- completed stage outputs/checkpoints
- explicit failure/assumption delta
- request to modify remaining stages only

### Replan budget

- configurable max replan attempts per task (example: 3)
- exceed budget -> escalate to safe fallback / human-required path

### Checkpointed state

Persist after each stage:

- completed stages
- artifacts/results
- current assumptions
- replan count
- next stage pointer

## 6.7 Estimator Feedback Interface

During execution, send progress snapshots:

- `task_id`
- `elapsed_seconds`
- `stages_completed`
- `total_stages`
- `current_stage`
- `partial_results` (optional)
- `resource_usage`
- `errors` (optional)

Estimator returns updated remaining-time estimate and confidence for user communication and scheduler hints.

## 6.8 Resource Manager Interaction

Before resource-sensitive stages, Executor checks:

- GPU utilization and free VRAM
- system memory pressure
- network latency / rate limit health
- agent machine availability

If constrained:

- wait/retry (temporary contention)
- fallback to lighter execution mode
- pause and hand control back to scheduler if no viable path

Placement note:

- Executor control loop runs on the Mac mini control plane.
- Individual stages may execute locally or be delegated to GPU/agent targets.

## 6.9 Persistence and Fault Tolerance

Long-running tasks must survive process/node interruptions.

Requirements:

- periodic and stage-boundary checkpoints to Memory Service
- restart recovery from last good checkpoint
- idempotent stage replay behavior when checkpoint boundary is ambiguous

On restart:

1. load latest checkpoint
2. verify current environment assumptions
3. resume next eligible stage or request replan

## 6.10 Notification Responsibilities

Executor triggers Notifier events for:

- task started (with initial estimate)
- significant progress (stage completion / estimate shift)
- task completed (final result summary)
- errors requiring user attention

Notifier owns channel fanout (Pi voice, mobile, etc.).

## 6.11 Abstract Interfaces (TBD Contracts)

### Executor <-> Estimator

- `UpdateProgress(snapshot) -> UpdatedEstimate`

### Executor <-> Resource Manager

- `CheckStageResources(task_id, stage_id, hints) -> AvailabilityDecision`

### Executor <-> Domain Expert (Replan)

- `RevisePlan(replan_request) -> RevisedRemainingPlan`

### Executor <-> Memory Service

- `SaveCheckpoint(checkpoint)`
- `LoadCheckpoint(task_id) -> checkpoint`

### Executor <-> Notifier

- `Publish(event)`

## 6.12 Open Questions (Needs Vetting)

1. Validation criteria schema representation:
   - JSON Schema vs declarative assertions vs executable validators.
2. Machine-usable severity encoding:
   - static policy table vs learned mapping.
3. Parallel stage execution model:
   - dependency graph semantics, validation ordering, conflict handling.
4. Sub-task spawning semantics:
   - parent-child lifecycle and rollback model.
5. Resource Manager API mode:
   - polling vs push/subscription.
6. Estimator update mode:
   - push vs poll vs hybrid.
7. Replan failure fallback:
   - timeout policy, degraded execution, escalation path.
8. Local fix policy granularity:
   - per-stage profiles vs per-tool profiles vs global defaults.

## 6.13 Recommended Initial Defaults (Draft)

- sequential execution by default
- optional parallelism only when stage graph explicitly declares independence
- replan budget default: 3
- estimator updates:
  - every stage completion
  - plus periodic heartbeat for long stages (e.g., 30s)
- checkpoint:
  - at stage boundaries
  - plus timed checkpoint for long-running single stages
