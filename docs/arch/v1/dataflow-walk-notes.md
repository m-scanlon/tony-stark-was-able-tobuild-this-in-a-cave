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

## Canonical Pipeline (v1)

```
Voice Shard → Event Ingress → SQLite inbox → Queue → Internal Router → Estimator → External Router → LLM Session (planning + execution)
                                                                              ↓
                                                                        Job Registry
```

- SQLite inbox: durability + ACK only. Not a work queue.
- Queue: single queue for v1. Estimator owns lane assignment. No split queues.
- Internal Router: context-aware job formation. Pulls context from Context Engine, resolves the domain agent, assembles `job_envelope_v1`.
- Estimator: the scheduler. Reads `job_envelope_v1`, does a shallow consult with the agent domain for rough complexity, picks the target shard based on capability profiles and current load. Writes placement to Job Registry.
- External Router: receives the Estimator's placement decision and dispatches the job to the assigned shard.
- Job Registry: passive lifecycle tracker only. Source of truth for job state (created → routed → planning → executing → done / failed). Does not make decisions.
- LLM Session: owns the full lifecycle — planning AND execution. One session, one context window. Updates Job Registry as the job progresses. No handoff mid-job.

No split queues in v1. Revisit if lane contention becomes a real problem.

---

## What the Internal Router Needs from the Queue

To assemble `job_envelope_v1`, the Internal Router needs:

- `event_id` — internal reference to the inbox row
- `turn_id` + `session_id` — for tracing and dedup
- `transcript` — what the user said (from the raw event payload)
- `triage_hints` — intent classification and `latency_class` (from the raw event payload)
- `session_state` — new job or continuation (drives routing logic)

From there, the Internal Router pulls from the Context Engine to resolve the domain agent and build the context package. The assembled `job_envelope_v1` — containing transcript + resolved agent + context package — is what everything downstream receives.

`job_envelope_v1` schema is still to be locked. That's the first question the dataflow walk must answer.

## What the Estimator Needs from `job_envelope_v1`

To make a placement decision, the Estimator needs:

- `latency_class` — from `triage_hints`, signals urgency
- `agent_id` — to do a shallow consult with the agent domain for rough complexity
- `job_envelope_v1` reference — passed through to the assigned shard via the External Router

The Estimator picks the target shard based on capability profiles and current load, then writes the placement to the Job Registry. The Job Registry's record: (`job_id`, `event_id`, `agent_id`, `shard`, `status`, timestamps).

---

## Key Decisions — Router Split, Estimator-as-Scheduler, Job Registry

Locked during the pipeline design walkthrough.

**Router split (Internal + External).**
A single "Router" was too much responsibility in one component. Splitting it makes the roles explicit:
- Internal Router = job formation. It's a context operation — pull context, resolve agent, assemble envelope. No dispatch knowledge required.
- External Router = dispatch. It knows which shard to call and how. No context knowledge required.
These are different concerns and should fail independently.

**Estimator is the scheduler.**
"Estimator" was previously just a complexity scorer. That framing was wrong — complexity scoring is only useful if something acts on it. The Estimator now owns the full placement decision: reads the envelope, consults the agent domain for rough complexity, checks shard capability profiles and current load, and picks a shard. It's the scheduler. The name stays "Estimator" to reflect that it's doing estimation to drive a decision, not running a planning algorithm.

**"Scheduler" is retired as a term.**
The old "Scheduler" concept conflated two things: placement decisions and lifecycle tracking. Those are now separate components with separate responsibilities. Using "Scheduler" anywhere in v1 docs is wrong — replace with "Estimator" (placement) or "Job Registry" (lifecycle) as appropriate.

**Job Registry is passive.**
The Job Registry does not decide anything. It is a state machine store — components write transitions to it, nothing more. Estimator writes on placement. LLM Session writes as the job progresses through planning, executing, done/failed. The registry is the source of truth for operational job state, queryable for monitoring and recovery, but it does not drive the pipeline.

**LLM Session owns planning and execution together.**
Splitting planning (Domain Expert) and execution (Executor) into separate LLM sessions was considered and rejected. Splitting context mid-job forces re-hydration, loses accumulated reasoning, and adds latency. One session holds the full context window for the life of the job — it plans, executes, validates, and replans without any handoff. The phases (planning → executing → validating → replanning) are semantic labels within a single session, not separate processes.

---
