# 6. Executor (Draft)

> Draft status: this design is based on conceptual notes and **requires product owner validation** before implementation lock-in.

## 6.1 Overview

The Executor is the heap-driven execution loop for Skyra task execution. It consumes `WorkPlan` (ephemeral) or `TaskSheet` (stateful) and executes via repeated heap re-entry — one tool call per heap cycle.

**The heap is the execution loop.** The LLM session picks up a job from the heap, runs until it hits a tool call, the tool executes, and the job re-enters the heap with the updated context blob. This repeats until one of two exit conditions is reached:

- `finished` tool — job is done, context blob committed
- `contact_user` tool — job parks on the heap, waiting for user response

Every other tool call re-queues. The context window is the job state — no separate checkpoint or serialization needed. Preemption is free: higher priority work simply gets picked up first while the job waits on the heap.

The Executor handles:

- tool call validation (BoundaryValidator before each dispatch)
- assumption drift detection and bounded replanning
- working state writes (scratch pad, no approval needed)
- commit proposals for canonical state changes (user-gated)

## 6.2 Role in Pipeline

Placement in end-to-end flow:

```text
Event → Internal Router → Domain Agent (estimation call) → [Max-Heap] → Estimator → External Router → LLM Session (Domain Expert + Executor)
                                                                               |
                                                                         Job Registry
```

Authority boundary:

- Voice Shard can provide provisional responses but does not execute authoritative task pipelines.
- Brain Shard control plane owns orchestration authority and coordinates LLM Session runs.
- Shards with appropriate capabilities are execution targets selected by the Estimator based on capability profiles and current load.

The Estimator is responsible for placement decisions. The External Router dispatches the job to the selected shard. The LLM Session owns both planning and execution — one context window, no handoff mid-job. The Job Registry tracks lifecycle state passively throughout.

## 6.3 Core Responsibilities

The Executor must:

1. Execute step-by-step with validation.
2. Detect and classify assumption drift severity.
3. Attempt local fixes before replanning.
4. Trigger bounded, incremental replanning via Domain Expert.
5. Send progress snapshots to Estimator.
6. Consult Resource Manager before constrained steps.
7. Persist checkpoints for crash/restart recovery.
8. Emit user-facing lifecycle notifications through Notifier.
9. Run BoundaryValidator before each step — check proposed tool calls against the project boundary. If a tool is locked, pause execution and send a permission prompt to the user. User responds `allow_always | allow_once | deny`. On `deny`, trigger bounded replan rather than halting the task.

## 6.4 Execution Loop

The execution loop is heap-driven. Each cycle:

1. Job is picked up from the heap.
2. LLM session resumes from the current context blob.
3. LLM runs and emits a tool call.
4. BoundaryValidator checks the tool against the agent boundary. If locked: permission prompt to user (`allow_always | allow_once | deny`). `allow_always` commits a boundary update to `state.json` before resuming. `deny` triggers bounded replan.
5. Tool executes.
6. Result appended to context blob.
7. Check exit conditions:
   - `finished` tool → job complete, commit if stateful, done
   - `contact_user` tool → park job on heap, waiting for user response
   - any other tool → re-queue job on heap with updated context blob, repeat

The importance score on the heap item can be updated between cycles — if the job discovers something that raises the stakes, its score bumps up and it gets prioritized accordingly.

Steps from the `WorkPlan` or `TaskSheet` remain valid as planning artifacts — they describe intent and expected sequence. But execution doesn't follow them mechanically. The LLM drives tool selection turn by turn, with the plan as context.

## 6.5 Severity Levels for Assumption Changes

| Level    | Example                                         | Action                                                          |
| -------- | ----------------------------------------------- | --------------------------------------------------------------- |
| Trivial  | API slightly slower than expected               | Log and continue                                                |
| Minor    | Missing optional field with safe default        | Adjust locally and continue                                     |
| Moderate | Required file missing but can be generated      | Attempt local fix; if fix fails, replan affected remaining path |
| Major    | Core assumption false (intent/goal mismatch)    | Pause and replan remaining steps                               |
| Critical | Security risk, impossible/safety-unsafe step, or user denied tool and no replan path exists | Halt task, notify user, abort                            |

## 6.6 Controlled Adaptation (Loop Prevention)

### Local fixes first

Try low-cost recovery before calling Domain Expert:

- retry with exponential backoff
- fallback tool/model path
- parameter adjustment (batch size/model variant/timeouts)

### Incremental replanning only

On replan request, send:

- original TaskSheet/WorkPlan
- completed step outputs/checkpoints
- explicit failure/assumption delta
- request to modify remaining steps only

### Replan budget

- configurable max replan attempts per task (example: 3)
- exceed budget -> escalate to safe fallback / human-required path

### Checkpointed state

Persist after each step:

- completed steps
- artifacts/results
- current assumptions
- replan count
- next step pointer

## 6.7 Estimator Feedback Interface

The Estimator made the original placement decision — it read the estimation call output (`{is_job, complexity, domain}`), matched the complexity score against shard capability profiles and current load, then selected the target shard. It wrote that placement to the Job Registry.

During execution, the LLM Session sends progress snapshots back to the Estimator:

- `task_id`
- `elapsed_seconds`
- `steps_completed`
- `total_steps`
- `current_step`
- `partial_results` (optional)
- `resource_usage`
- `errors` (optional)

Estimator returns updated remaining-time estimate and confidence for user communication. These updates are also written to the Job Registry to keep lifecycle state current.

## 6.8 Resource Manager Interaction

Before resource-sensitive steps, Executor checks:

- GPU utilization and free VRAM
- system memory pressure
- network latency / rate limit health
- Shard availability

If constrained:

- wait/retry (temporary contention)
- fallback to lighter execution mode
- pause and surface the constraint — Estimator or Job Registry can flag for rescheduling if no viable path

Placement note:

- Executor control loop runs on the Brain Shard control plane.
- Individual steps may execute locally or be delegated to Shards with matching capability profiles.

## 6.9 Preemptive Scheduling

Preemption is a natural property of the heap-driven execution model — no special mechanism needed.

Between tool calls, the job is on the heap. Higher priority work just gets picked up first. The job waits. When the machine is free and the job is the highest priority item, it gets picked up again and resumes from its context blob.

**The context window is the job state.** Nothing is serialized specially for preemption — the context blob that re-entered the heap already contains everything: completed tool calls, outputs, remaining intent. Resume is seamless. The LLM does not know it waited.

The only scheduling constraint: a running tool call is never interrupted mid-execution. The preemption point is always the re-queue after tool completion.

## 6.9a Persistence and Fault Tolerance

Long-running tasks must survive process/node interruptions (crashes, power loss) — distinct from intentional preemption.

Requirements:

- periodic and step-boundary checkpoints to Agent Service
- restart recovery from last good checkpoint
- idempotent step replay behavior when checkpoint boundary is ambiguous

On restart:

1. load latest checkpoint
2. verify current environment assumptions
3. resume next eligible step or request replan

## 6.9b Working State vs Committed State

The executor can write freely to the **working state** partition of the object store. This is a scratch pad — the system uses it to test ideas, validate approaches, run intermediate computations, and reason through problems on paper. No user approval required.

**Committed state** requires user approval. When the executor produces output worth persisting canonically, it proposes a commit. The user accepts or rejects. Only accepted commits enter the version history.

The distinction:
- Working state: mutable, free, throwaway. The system's thinking space.
- Committed state: user-gated, canonical, versioned, permanent.

Version history tracks commits only. Working state does not pollute the audit trail.

## 6.10 Notification Responsibilities

Executor triggers Notifier events for:

- task started (with initial estimate)
- significant progress (step completion / estimate shift)
- task completed (final result summary)
- errors requiring user attention

Notifier owns channel fanout (Voice Shard, mobile, etc.).

## 6.11 Abstract Interfaces (TBD Contracts)

### Executor <-> Estimator

- `UpdateProgress(snapshot) -> UpdatedEstimate`

Note: the Estimator also owns the upstream placement decision (`PlaceJob(estimation_output) -> PlacementDecision`). That contract belongs to the Estimator's own interface spec, not here. The Executor only interacts with the Estimator via progress updates during execution.

### Executor <-> Resource Manager

- `CheckStepResources(task_id, step_id, hints) -> AvailabilityDecision`

### Executor <-> Domain Expert (Replan)

- `RevisePlan(replan_request) -> RevisedRemainingPlan`

### Executor <-> Agent Service

- `SaveCheckpoint(checkpoint)`
- `LoadCheckpoint(task_id) -> checkpoint`

### Executor <-> Notifier

- `Publish(event)`

## 6.12 Open Questions (Needs Vetting)

1. Validation criteria schema representation:
   - JSON Schema vs declarative assertions vs executable validators.
2. Machine-usable severity encoding:
   - static policy table vs learned mapping.
3. Parallel step execution model:
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
   - per-step profiles vs per-tool profiles vs global defaults.

## 6.13 Recommended Initial Defaults (Draft)

- sequential execution by default
- optional parallelism only when step graph explicitly declares independence
- replan budget default: 3
- estimator updates:
  - every step completion
  - plus periodic heartbeat for long steps (e.g., 30s)
- checkpoint:
  - at step boundaries
  - plus timed checkpoint for long-running single steps
