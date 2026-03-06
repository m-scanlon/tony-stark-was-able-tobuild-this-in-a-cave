# Job Registry

## What It Is

The Job Registry is a passive lifecycle tracker on the Brain Shard. It is the source of truth for job state from creation through completion or failure. It does not make placement, routing, or lane decisions — those are owned by the Estimator.

It is intentionally simple in v1. Single queue, state tracking only.

## Responsibilities

- Record jobs as they enter the system (`created`)
- Reflect routing decisions made by the Estimator (`routed`, `shard_id`, `lane`, `routed_at`)
- Track job status through planning, execution, and completion
- Surface active job state to other services (context injector)

## What It Does Not Do

- Does not form tasks — that is the Domain Expert's job
- Does not execute tasks — that is the assigned LLM session's job
- Does not make semantic decisions about job content
- Does not own the TaskSheet or WorkPlan — those live in the object store
- Does not assign lanes or route to shards — that is the Estimator's job

## Position in the Pipeline

```
event
  → inbox (SQLite, event_id PK)
  → queue
  → Estimator (assigns lane, shard, creates job record → Job Registry)
  → assigned LLM session (task formation + execution)
  → Job Registry (marks planning → executing → completed / failed)
```

Canonical pipeline reference: `docs/arch/v1/scyra.md` section 10.2

## Execution Lanes

| Lane | Used For |
|---|---|
| `fast_local` | Short, low-cost requests handled by local Brain Shard models |
| `deep_reasoning` | Complex requests routed to a Shard with deep_reasoning capability |

The Estimator makes the lane assignment. The Job Registry records it.

## Job Lifecycle

```
created → routed → planning → executing → completed
                                        → failed
```

- `created`: job accepted, record written to registry
- `routed`: Estimator has assigned lane and shard; `shard_id`, `lane`, and `routed_at` are set
- `planning`: assigned LLM session is forming the task
- `executing`: LLM session is actively executing
- `completed`: session finished successfully
- `failed`: unrecoverable error

## Job Envelope

Each job entering the Job Registry carries a `job_envelope_v1`:

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

The jobs table is operational state only. Owned exclusively by the Job Registry.

Access rules:
- Job Registry: read + write
- Estimator: write on routing (sets `shard_id`, `lane`, `routed_at`, `status = routed`); read otherwise
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
