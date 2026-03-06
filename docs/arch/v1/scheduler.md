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
  "domain": "servers"
}
```

Complexity is measured in **estimated tool calls**. This is the unit of complexity — concrete, measurable, directly tied to compute cost.

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

Weight updates, pattern detection, cross-domain analysis. Runs on idle compute — typically at night when estimation and job work is quiet.

On ingress, every turn is stored in RDS with a record of which domain agents were routed to in real-time. The batch process reads this at night and runs the turn against every agent that wasn't reached in real-time. This is how data integrity is preserved — nothing is permanently missed.

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

**The context window is the job state.** To interrupt a job:

1. Wait for the current tool call boundary — jobs are never interrupted mid-tool-call
2. Serialize the full context window at that boundary
3. Push serialized context onto a **FIFO stack**
4. Machine handles the high priority work
5. When machine is free, pop context from FIFO and resume generation

The LLM does not know it was interrupted. The context contains everything — completed steps, tool outputs, remaining intent. Resume is seamless.

FIFO ordering within the interrupted job stack: first interrupted, first resumed. Fair and simple.

---

## Estimator

The Estimator's job is placement. When a job enters the heap with `complexity > 1`, the Estimator:

1. Reads the complexity score and domain from the estimation output
2. Checks available shard capability profiles and current load
3. Assigns the job to the best available machine

| Complexity range | Likely target |
|---|---|
| ≤ 1 | Inline (never reaches Estimator) |
| 2–5 | Mac mini (fast, capable) |
| 6–15 | Mac mini or GPU machine depending on load |
| 16+ | GPU machine (deep reasoning) |

These ranges are illustrative. Actual routing is based on registered capability profiles, not hardcoded rules.

---

## Open Questions

- What is the exact composite formula for importance score?
- How does wait time factor in — do low priority items get a priority bump over time to prevent starvation?
- What triggers a complexity score revision — can the domain agent revise the estimate mid-execution?
- What is the maximum FIFO depth for interrupted jobs?
- How does the system handle a job that gets interrupted repeatedly?
