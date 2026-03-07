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
  → Internal Router (labels turn, routes to domain agents)
  → Domain Agent (estimation call → {is_job, complexity, domain})
  → Max-Heap (all work ordered by importance score)
  → Estimator (reads complexity, matches to capable shard → Job Registry)
  → assigned LLM session (task formation + execution)
  → Job Registry (marks planning → executing → completed / failed)
```

Canonical pipeline reference: `docs/arch/v1/dataflow-walk-notes.md`, `docs/arch/v1/scheduler.md`

## Shard Placement

The Estimator matches complexity score (in estimated tool calls) against registered shard capability profiles. There are no hardcoded lanes — routing is capability-profile-based.

| Complexity | Likely target |
|---|---|
| ≤ 1 | Inline execution — never reaches heap or Estimator |
| 2–5 | Mac mini class |
| 6+ | GPU machine or most capable available shard |

The Estimator makes the placement decision. The Job Registry records it.

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

## Job Entry Contract

Each job entering the heap originates from an estimation call produced by the domain agent:

```json
{
  "is_job": true,
  "complexity": 3,
  "domain": "servers"
}
```

Complexity is measured in estimated tool calls. This is the primary placement signal — the Estimator reads it and matches against shard capability profiles. The estimation call schema is not yet locked. See `docs/arch/v1/gaps.md` G1.

## Data Model

See `schema.sql`.

The jobs table is operational state only. Owned exclusively by the Job Registry.

Access rules:
- Job Registry: read + write
- Estimator: write on routing (sets `shard_id`, `lane`, `routed_at`, `status = routed`); read otherwise
- Context Injector: read only (`status`, `agent_id`)

## v1 Constraints

- Backpressure and overload policies are undefined — see `docs/arch/v1/gaps.md` G6
- Preemptive scheduling is supported — higher priority work can interrupt in-flight jobs. Interrupted job's context window is serialized to a FIFO stack and resumed when the machine is free. See `docs/arch/v1/scheduler.md`.
- Transport ACK confirms durable ingest only — execution may occur later from the heap

## Related Docs

- `docs/arch/v1/scheduler.md` — unified heap, three inference types, complexity scoring, preemptive scheduling
- `docs/arch/v1/dataflow-walk-notes.md` — updated canonical pipeline
- `docs/arch/v1/task-formation.md` — domain agent as doorkeeper, estimation call
- `docs/arch/v1/gaps.md` — known open issues including estimation call schema lock (G1)
- `skyra/internal/delegation/README.md` — estimator placement role
