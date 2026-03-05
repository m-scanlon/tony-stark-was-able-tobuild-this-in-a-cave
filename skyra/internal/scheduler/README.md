# Scheduler Service

## What It Is

The Scheduler is a control-plane service on the Brain Shard. It sits between the Estimator and the assigned LLM session. Its job is to receive annotated jobs, assign them an execution lane, and track their operational lifecycle.

It is intentionally simple in v1. Single queue, lane assignment only.

## Responsibilities

- Receive jobs from the queue after Estimator annotation
- Assign an execution lane (`fast_local` or `deep_reasoning`)
- Track operational job status from queued through completion
- Surface active and queued jobs to other services (context injector)

## What It Does Not Do

- Does not form tasks — that is the Domain Expert's job
- Does not execute tasks — that is the assigned LLM session's job
- Does not make semantic decisions about job content
- Does not own the TaskSheet or WorkPlan — those live in the object store
- Does not decide final scheduling policy — the Estimator informs, the Scheduler decides

## Position in the Pipeline

```
event
  → inbox (SQLite, event_id PK)
  → queue
  → Estimator (annotates lane hints)
  → Scheduler (assigns lane, creates job record)
  → assigned LLM session (task formation + execution)
  → Scheduler (marks completed or failed)
```

Canonical pipeline reference: `docs/arch/v1/scyra.md` section 10.2

## Execution Lanes

| Lane | Used For |
|---|---|
| `fast_local` | Short, low-cost requests handled by local Brain Shard models |
| `deep_reasoning` | Complex requests routed to a Shard with deep_reasoning capability |

The Estimator provides lane hints. The Scheduler makes the final assignment.

## Job Lifecycle

```
queued → running → completed
                 → failed
```

- `queued`: job accepted, waiting for lane assignment
- `running`: lane assigned, LLM session is active
- `completed`: session finished successfully
- `failed`: unrecoverable error

## Job Envelope

Each job entering the scheduler carries a `job_envelope_v1`:

- `job_id`
- `parent_job_id`
- `agent_id`
- `intent`
- `priority`
- `required_tools`
- `target` (`none | control_plane | shard:<id>`)
- `risk_level` (`low | med | high`)
- `expect_response_by`
- `schema_version`

Note: `job_envelope_v1` schema is not yet locked. See `docs/arch/v1/gaps.md` G1.

## Data Model

See `schema.sql`.

The jobs table is operational state only. Owned exclusively by the Scheduler.

Access rules:
- Scheduler: read + write
- Estimator: read only
- Context Injector: read only (`status`, `agent_id`)

## v1 Constraints

- Single queue, no priority tiers
- Lane assignment is the only routing decision
- No backpressure or overload handling
- No job cancellation or pause/resume
- transport ACK confirms durable ingest only — execution may occur later from queue

## Related Docs

- `docs/arch/v1/scyra.md` — canonical pipeline and job envelope
- `docs/arch/v1/task-formation.md` — what happens inside the assigned LLM session
- `docs/arch/v1/gaps.md` — known open issues including job_envelope_v1 schema lock
- `skyra/internal/delegation/estimator/DESIGN.md` — estimator design
