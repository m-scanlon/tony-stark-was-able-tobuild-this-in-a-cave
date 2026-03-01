# Next Steps — Executor Loop Design

## Why This Document Exists

Before writing any executor code, we need to fully understand the data that flows into it. The executor is the phase where a planned job actually runs — tools get called, state gets mutated, replanning happens. Getting the data model wrong here means the executor loop will be brittle or impossible to reason about.

This document captures the current state of our thinking on the event model and job lifecycle, and identifies what still needs to be locked before we can design the executor loop.

---

## What We've Established

### 1. The Canonical Pi → Mac Event

Every user request starts on the Pi as a `voice_event_v1`. This is the only event type the Pi sends to Mac. It carries:

- `transcript` — what the user said
- `triage_hints` — latency class, provisional eligibility, confidence
- `context_window` — session summary, recent turns, active project
- `context_state` — how much token headroom the Pi's model has (Mac uses this to size the context package it sends back)
- `session_state` — whether this is a new job or a continuation of an existing one
- `pi_gave_provisional` / `provisional_text` — whether Pi already said something, so Mac can reconcile

The `event_id` is **not** in the payload — it's stamped by the Pi's outbox layer before sending.

### 2. The Pi Decides: New Job or Continuation

The Pi owns the turn loop. It knows if there's already an open job. It encodes that in `session_state`:

```json
"session_state": {
  "pending_job_id": "job_abc123",
  "waiting_for": "user_approval"
}
```

- `pending_job_id: null` → new job
- `pending_job_id` set → continuation of an existing job (waiting on approval, clarification, etc.)

Mac's Event Ingress reads `session_state` and routes accordingly. No separate classification service.

### 3. The Mac Decides: WorkPlan vs TaskSheet

The Pi does **not** decide the task type. That's the Domain Expert's job (the planning phase on Mac) after it has:

- Retrieved project context
- Retrieved relevant tools
- Understood the intent

| Task type | When | Artifact |
|---|---|---|
| `WorkPlan` | Ephemeral — no state commit needed | `workplan.json` |
| `TaskSheet + Patch` | Stateful — requires project state commit | `tasksheet.json` + patch |

### 4. Job Lifecycle Phases

One job = one LLM session. A job moves through phases:

```
planning → executing → validating → replanning → done
                                  ↑___________↓
```

- **planning** — Domain Expert retrieves context, forms the plan artifact
- **executing** — Executor runs tool calls from the plan
- **validating** — checks outcomes against expected state
- **replanning** — if validation fails or a tool is denied, revises the plan and loops back
- **done / cancelled / failed** — terminal states

These are **semantic phases** tracked in the tasksheet. The scheduler's jobs table tracks operational status separately (queued / running / completed / failed).

### 5. Open Question: Project Routing

How does Mac know which project a new request belongs to? The Pi carries `active_project` in `context_window`, but this comes from the Pi's local cache — it could be stale.

This is unresolved. Options discussed but not decided:
- Trust Pi's `active_project` as authoritative (simple, but stale risk)
- Mac re-derives project from the transcript + context on every request (slower, more reliable)
- Pi and Mac stay in sync via context package pushes from the Context Injector (preferred direction — already in the design, needs formalization)

---

## What We Need to Define Next

To design the executor loop, we need to answer:

1. **What exactly does the executor receive?** — the full `job_envelope_v1` schema. What fields does it contain? What does it look like when handed off from the scheduler?

2. **How does the executor know what tools to call?** — global tools are always injected. Local tools are retrieved by the Domain Expert during planning and embedded in the task artifact. Does the executor re-retrieve, or trust the plan?

3. **What does the replanning trigger look like?** — tool denied (BoundaryValidator), validation failure, or confidence threshold not met. How does the executor signal replan vs halt vs clarify?

4. **How does the executor emit progress to the user?** — `UPDATE` and `PLAN_PROGRESS` events flow back to Pi during execution. When are they emitted? Who controls the cadence?

5. **How does a state commit happen mid-execution?** — the executor calls `propose_commit` → `apply_commit` through the Project Service global tools. Is this per-step or at the end of the plan?

---

## Related Docs

- `docs/arch/v1/scyra.md` — full system architecture
- `docs/arch/v1/executor.md` — executor design (current state, incomplete)
- `docs/arch/v1/domain-expert/README.md` — planning phase, plan approval gate
- `skyra/internal/project/README.md` — project service, tool hydration, boundary enforcement
- `skyra/internal/scheduler/README.md` — scheduler, job lifecycle, lanes
