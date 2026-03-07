# Data Model Walk — Misc Notes

Decisions and observations captured while walking the pipeline service by service. Not formal design docs — just things worth remembering.

---

## Event Register (Voice Shard only, for now)

The Voice Shard maintains an event register — a hash table keyed by `turn_id` — to track events that have been sent but not yet ACKed.

- Send event → insert entry
- ACK received → lookup by `turn_id`, pop
- Retry loop → scan for entries past `next_attempt_at`
- Durable (SQLite) so Pi reboots don't lose in-flight events

**The Brain Shard does not need an equivalent structure for ACKing.** The ACK leg is synchronous — Brain Shard receives, writes to SQLite inbox, sends ACK back on the same connection. No tracking structure needed to facilitate that. The Brain Shard's SQLite inbox is for downstream pipeline consumption, not ACK mechanics.

**GPU Shard probably doesn't need an event register.** In v1 it's a dumb inference endpoint — request/response, no async event tracking. Revisit when we walk the GPU Shard's role. If it ever handles async or multi-step work independently, it'd need one.

**Future: Shard bootstrap package.** Brain Shard knows each Shard's capability profile from registration. It could provision each Shard with a tailored bootstrap package — Voice Shard gets event register + retry loop + voice config, GPU Shard gets inference config only. Package contents derived from registered capabilities. Capture properly when we get to Shard provisioning.

---

## ACK uses turn_id, not event_id

`event_id` is internal to the Brain Shard — generated on ingress, never crosses the wire to Voice Shard. ACK references `turn_id` only. Voice Shard clears its event register by `turn_id`.

---

## Canonical Pipeline (v1) — Updated

```
Voice Shard → Event Ingress → SQLite inbox → Internal Router → [Max-Heap] → Estimator → LLM Session (planning + execution)
                                                    ↓                              ↓
                                                   RDS                       Job Registry
```

- SQLite inbox: durability + ACK only. Not a work queue.
- Internal Router: simplified. Drops off turn data at context engine. Labels turn (in-domain | other) using context blob. Routes to relevant domain agents. Done. Does not assemble `job_envelope_v1`.
- RDS: stores all turns with routing metadata. Batch process reads this at night to run turns against agents that weren't reached in real-time.
- Domain Agent: doorkeeper. Self-selects relevance. Checks for job impact. Forms estimation call if a job is needed.
- Max-Heap: all work items ordered by importance score. Three types: estimation (very high priority), job (high priority), batch (very low priority). One heap — no separate queues.
- Estimator: reads complexity score from estimation output, matches to capable machine via shard capability profiles. Simpler placement than before — complexity score in tool calls is the primary input.
- External Router: dispatches to assigned shard.
- Job Registry: passive lifecycle tracker. Source of truth for job state.
- LLM Session: owns full lifecycle — planning AND execution. One context window. Updates Job Registry as job progresses.

**Three inference types in the heap:**
- `estimation` — very high priority. "Is this a job? How complex?" Complexity ≤ 1 → execute inline, never enters heap.
- `job` — high priority. Long-running execution.
- `batch` — very low priority. Weight updates, pattern detection. Runs on idle compute at night.

**Preemptive scheduling:** Higher priority work interrupts in-flight jobs at tool call boundaries. Interrupted job's context window serialized to FIFO stack. Resume = pop context, continue generation.

---

## What the Internal Router Needs

The Internal Router is now simple. It needs:

- `event_id` — internal reference to the inbox row
- `turn_id` + `session_id` — for tracing and dedup
- `transcript` — what the user said
- `triage_hints` — intent classification and `latency_class`
- `session_state` — new job or continuation

It drops off the turn in RDS, attaches the context blob (pushed by CIX, already available), labels the turn via the front face transformer, and routes to the relevant domain agents. It does not assemble a complex job envelope. That work has moved to the domain agent.

## What the Estimator Needs

The Estimator now reads the **estimation call output** from the domain agent — not a complex `job_envelope_v1`:

```json
{
  "is_job": true,
  "complexity": 3,
  "domain": "servers"
}
```

Complexity (in estimated tool calls) is the primary placement signal. The Estimator matches against shard capability profiles and current load, picks the best available machine, writes placement to Job Registry.

---

## Key Decisions — Router Split, Estimator-as-Scheduler, Job Registry

Locked during the pipeline design walkthrough.

**Router split (Internal + External).**
A single "Router" was too much responsibility in one component. Splitting it makes the roles explicit:
- Internal Router = turn labeling and domain agent routing. It labels the turn (in-domain | other) using the context blob, routes to relevant domain agents, and drops the turn in RDS. Job formation is NOT the Internal Router's job — that moved to the domain agent.
- External Router = dispatch. It knows which shard to call and how. No context knowledge required.
These are different concerns and should fail independently.

**Estimator is the scheduler.**
"Estimator" was previously just a complexity scorer. That framing was wrong — complexity scoring is only useful if something acts on it. The Estimator now owns the full placement decision: reads the envelope, consults the agent domain for rough complexity, checks shard capability profiles and current load, and picks a shard. It's the scheduler. The name stays "Estimator" to reflect that it's doing estimation to drive a decision, not running a planning algorithm.

**"Scheduler" as a component name is retired.**
The old "Scheduler" component conflated two things: placement decisions and lifecycle tracking. Those are now separate: "Estimator" (placement) and "Job Registry" (lifecycle). `docs/arch/v1/scheduler.md` exists but describes the heap-based scheduling *system* (unified max-heap, inference types, preemption) — not a component called "Scheduler".

**Job Registry is passive.**
The Job Registry does not decide anything. It is a state machine store — components write transitions to it, nothing more. Estimator writes on placement. LLM Session writes as the job progresses through planning, executing, done/failed. The registry is the source of truth for operational job state, queryable for monitoring and recovery, but it does not drive the pipeline.

**LLM Session owns planning and execution together.**
Splitting planning (Domain Expert) and execution (Executor) into separate LLM sessions was considered and rejected. Splitting context mid-job forces re-hydration, loses accumulated reasoning, and adds latency. One session holds the full context window for the life of the job — it plans, executes, validates, and replans without any handoff. The phases (planning → executing → validating → replanning) are semantic labels within a single session, not separate processes.

---
