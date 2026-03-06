# Next Steps — Post Architecture Revision

## What Changed

A major design session revised the core pipeline. The central classifier is gone. Domain agents are the doorkeepers. All work flows through a unified max-heap. Three inference types govern priority. The Internal Router is simplified. The Estimator reads complexity scores, not complex job envelopes.

See `docs/arch/v1/dataflow-walk-notes.md` for the updated canonical pipeline.
See `docs/arch/v1/scheduler.md` for the full heap, inference type, and preemptive scheduling design.

---

## What's Resolved

The `job_envelope_v1` design problem that previously blocked this doc is largely resolved. The domain agent owns job formation — the envelope it passes to the heap is simple:

```json
{
  "is_job": true,
  "complexity": 3,
  "domain": "servers"
}
```

No complex context assembly at the routing layer. The context came with the event.

---

## Open Design Questions

### 1. Lock the Estimation Call Schema

The estimation call output needs a locked schema before anything downstream can be implemented. Current draft:

```json
{
  "is_job": true,
  "complexity": 3,       // estimated tool calls
  "domain": "servers",
  "importance": 85,      // composite importance score for heap ordering
  "latency_class": "interactive | background"
}
```

Questions:
- What other fields does the Estimator need for placement?
- Is `importance` derived here or by the front face transformer upstream?
- How does a continuation (existing job) flow through vs a new job?

### 2. Define "Other" Turn Storage in RDS

Turns labeled "other" by the front face transformer get stored in RDS for batch pickup. The schema needs to capture enough for the batch process to route correctly at night.

Minimum fields:
- `turn_id`, `session_id`, `event_id`
- `transcript`
- `context_blob_ref` — reference to the context blob at ingress time
- `routed_agents[]` — which agents were routed to in real-time (batch skips these)
- `created_at`

Questions:
- Does the context blob need to be snapshotted at ingress, or can the batch process reconstruct it?
- What is the retention policy for "other" turns in RDS?

### 3. Define the Batch Job Contract

The nightly batch process runs all agents against accumulated session context. The contract:

**Input:** all turns since the last batch run (from RDS), per agent
**What it does:**
- Updates importance vectors (weight updates deferred from real-time)
- Runs each agent against turns it didn't receive in real-time
- Detects cross-domain patterns (V3 behavior — new domain proposals)

**Output:**
- Updated importance vectors written back to agent object stores
- Pattern observations committed to context engine state
- New domain proposals surfaced to user (if applicable)

Questions:
- Does the batch process run as one heap item per agent, or one item per turn-agent pair?
- What model runs the batch inference — lightweight or full?
- How does the batch process handle a domain agent that's been archived?

### 4. Working State Schema in the Object Store

The executor can write freely to working state without user approval. The object store needs a defined partition for this.

Current thinking:
- `state.json` = committed state (canonical, user-approved)
- `working/` = scratch space (mutable, throwaway, not versioned)

Questions:
- Does working state get cleaned up automatically after a job completes?
- Can the system read from working state during planning, or only during execution?
- Does working state have a size limit?

### 5. Preemptive Scheduling — FIFO Stack Contract

The FIFO stack for interrupted jobs needs a defined schema:

```json
{
  "job_id": "...",
  "interrupted_at": "...",
  "context_snapshot": "...",   // serialized context window
  "priority": 72,
  "domain": "servers"
}
```

Questions:
- Where does the FIFO stack live — in-memory or persisted?
- What's the maximum depth?
- Does a repeatedly interrupted job get a priority bump to prevent starvation?

---

## Next Design Session

Walk the estimation call end to end:

1. Turn arrives on the heap as an estimation work item
2. Machine picks it up
3. Domain agent receives turn + context blob
4. Estimation call runs — outputs `{is_job, complexity, domain, importance}`
5. Complexity ≤ 1 → execute inline. What does inline execution look like exactly?
6. Complexity > 1 → job formed, pushed to heap. What does the job item in the heap look like?

Getting this right unblocks everything downstream — execution, scheduling, batch integration.

---

## Related Docs

- `docs/arch/v1/scheduler.md` — unified heap, inference types, complexity scoring, preemption
- `docs/arch/v1/dataflow-walk-notes.md` — updated canonical pipeline
- `docs/arch/v1/task-formation.md` — domain agent as doorkeeper, estimation call
- `docs/arch/v1/executor.md` — preemptive scheduling, working state
- `docs/arch/v1/context-engine.md` — context blob with all agents, batch weight updates
- `docs/arch/v1/importance-vectors.md` — importance vector design, V3 background process
