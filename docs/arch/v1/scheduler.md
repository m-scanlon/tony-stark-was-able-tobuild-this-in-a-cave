# Scheduler — Unified Heap and Inference Types

All work in Skyra flows through a single unified max-heap. There are no separate batch schedulers, job queues, or priority lanes. Everything is a work item with an importance score. The heap handles ordering. Machines handle execution.

---

## The Heap

A **max-heap** ordered by importance score. The highest importance item always sits at the root. Insert and extraction are O(log n).

Every inference call — regardless of type — enters the heap. The heap does not know or care about inference type. It only sees importance scores. Low priority batch work sits at the bottom and gets picked up when machines are idle. High priority estimation calls jump to the top immediately.

Within the same priority tier, ordering is FIFO — first in, first out.

---

## Three Inference Types

### 1. Estimation

**Priority: Very high.**

The first inference call for any turn. Answers one question: *is this a job, and if so how complex?*

Output:
```json
{
  "is_job": true,
  "complexity": 3,
  "reasoning_depth": 2,
  "cross_domain": false,
  "reversible": true,
  "output_scope": "fact",
  "domain": "servers"
}
```

- `complexity` — estimated tool calls. Primary placement signal.
- `reasoning_depth` — inferential steps required (1 = direct lookup, 2 = moderate synthesis, 3 = deep multi-step reasoning).
- `cross_domain` — true if the request spans multiple domain agents.
- `reversible` — false if the action cannot be undone (send message, delete data, external API write).
- `output_scope` — `fact | plan | commit`. Commits carry higher evaluation cost than facts.

`complexity` (tool call count) drives the inline vs heap threshold. The other fields inform shard selection and evaluation depth — a low tool-call job that is irreversible and cross-domain should not be treated as lightweight.

**Inline execution threshold: complexity ≤ 1.**

If the estimation output has `complexity ≤ 1`, the work is executed inline immediately — it never enters the heap. A light switch, a timer, a one-shot lookup. No queue, no placement decision, no latency. Just execute and done.

If `complexity > 1`, a job is formed and pushed into the heap. The Estimator matches the complexity score against available shard capability profiles and assigns placement.

The complexity threshold is currently set at **1** and will be tuned empirically from real usage data.

### 2. Job

**Priority: High.**

Long-running work. The domain agent owns job formation — it receives the turn and context, decides what needs to be done, and produces the job. The job enters the heap and gets picked up by the best available machine.

Jobs can be interrupted by higher priority work. See: Preemptive Scheduling below.

### 3. Batch

**Priority: Very low.**

Weight updates for non-routed domains, pattern detection, cross-domain analysis. Runs on idle compute — typically at night when estimation and job work is quiet.

**Important distinction:** domains that were routed to in real-time get their weight updates immediately via lightweight signal processing — not deferred to batch. Deferring the weight update for a selected domain would cause it to decay even though it was just used. The batch only handles domains that were NOT reached in real-time.

On ingress, every turn is stored in RDS with a record of which domain agents were routed to. The batch process reads this at night and runs the turn against every agent that wasn't in that set. This is how data integrity is preserved — nothing is permanently missed.

Batch work items enter the heap at very low importance. They sit at the bottom and get picked up when no higher priority work is waiting. No separate scheduler needed.

---

## Importance Score

The importance score on each heap item is a composite signal:

- **Inference type** — estimation is always high, batch is always low, jobs inherit from the estimation that created them
- **Latency class** — voice requests waiting on a response score higher than background work
- **Explicit user signal** — "this is important", "remind me", "don't forget" — caught by intent classification, bumps score
- **Domain history** — domains with recent high-stakes activity score higher

The front face transformer's importance assessment is the primary input. The system asks the user only when signals genuinely conflict or stakes are high enough to warrant it. Over time, as vectors build up and batch jobs run, the system calibrates without needing to ask.

---

## Preemptive Scheduling

Higher priority work can interrupt a running job. A new estimation call arriving while all machines are busy preempts the lowest priority in-flight job.

**The context window is the job state.** Preemption is a natural property of the heap-driven execution model — no special mechanism needed.

Between tool calls, the job is on the heap. Higher priority work gets picked up first. The job waits. When the machine is free and the job is the highest priority item, it resumes from its context blob. The LLM does not know it waited — the context contains everything: completed tool calls, outputs, remaining intent. Resume is seamless.

The only constraint: a running tool call is never interrupted mid-execution. The preemption point is always the re-queue after tool completion.

---

## Estimator

The Estimator is **an inference call, not a service**. It fires when the External Router picks up an estimation work item from the heap. There is no separate Estimator process — the External Router owns the heap, handles priority ordering, manages preemption, and dispatches work items. When an estimation item is dequeued, a prompt runs against the domain agent's estimation output and available shard state. That prompt is the Estimator.

The Estimator prompt:

1. Reads the estimation output (`is_job`, `complexity`, `domain`) from the domain agent
2. Checks current shard capability profiles and load
3. Decides: execute inline or place on heap as a full job

If `complexity ≤ 1`: execute inline immediately — the Estimator does the work itself, never forming a job.

If `complexity > 1`: place the job back onto the heap targeting the best available shard.

| Complexity range | Likely target |
|---|---|
| ≤ 1 | Inline (Estimator executes directly) |
| 2–5 | Mac mini (fast, capable) |
| 6–15 | Mac mini or GPU machine depending on load |
| 16+ | GPU machine (deep reasoning) |

These ranges are illustrative. Actual routing is based on registered capability profiles, not hardcoded rules.

Because the Estimator is a prompt, estimation quality improves as model quality improves — no code changes required.

---

## Open Questions

- What is the exact composite formula for importance score?
- How does wait time factor in — do low priority items get a priority bump over time to prevent starvation?
- What triggers a complexity score revision — can the domain agent revise the estimate mid-execution?
- What is the maximum FIFO depth for interrupted jobs?
- How does the system handle a job that gets interrupted repeatedly?
