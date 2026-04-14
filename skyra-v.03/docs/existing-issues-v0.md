# Existing Issues

Active problems observed in the running runtime. Each issue includes the symptom, root cause, and candidate solutions.

---

## Issue 1: Relationships are bidirectional by default

**Symptom**
When thalamus seeds a relationship with prefrontal, prefrontal gets thalamus in its addressable network. Prefrontal then routes back to thalamus — treating the upstream relay as a valid target. This creates a ping-pong loop where signal bounces between prefrontal and thalamus indefinitely.

**Root cause**
`seedRelationships` in `world.go` always seeds both sides of a relationship. Every peer gets a channel on both ends. There is no concept of directionality at the relationship layer. Thalamus is a one-way relay inward — prefrontal should never address it — but the architecture has no way to express that.

**Candidate solutions**

1. **Filter non-cognitive peers from DerivePresent** — the "your cognitive network" section only shows cognitive beings. Thalamus, sensory, and motor are non-cognitive and should never appear as addressable targets for a cognitive being. Cheap fix, doesn't require architectural change.

2. **Directional relationship seeding** — extend the genome protocol with a direction flag (`~direction inbound|outbound|bidirectional`). `seedRelationships` seeds only the specified side. Correct at the architectural level but requires protocol change and more complex seeding logic.

3. **Separate inbound and outbound peer maps on Being** — Being holds two peer maps: one for beings it can address, one for beings that can address it. Seeding writes to the appropriate map. Present only shows the addressable map. More invasive but cleanly separates the two concerns.

---

## Issue 2: The kernel drops malformed signals with no correction path — RESOLVED 2026-04-13

**Symptom**
When a cognitive being emits a malformed protocol string — missing expression, hallucinated source name, empty field — the kernel drops it, logs to stderr, and returns. The being that produced the bad signal never finds out. No correction happens. The exchange dies silently.

**Root cause**
The kernel has no feedback loop to inference. `AcceptSignal` returns a drop result. `dispatch` in `main.go` logs the drop reason to stderr and exits the loop. The model that produced the bad output is never told what it did wrong or given a chance to retry with correction.

**Candidate solutions**

1. **Retry with error appended to present** — when the kernel drops a signal, the next inference call includes the drop reason appended to the bottom of the same present: `your last signal was dropped: <reason> — try again`. Model gets one more shot with the correction in context. Simple, no new protocol, no new beings.

2. **Error as a being** — introduce a corrector being whose job is to surface malformed signals back to the sender as a new exchange. Keeps the protocol intact and makes error correction a first-class citizen of the system. Heavier to implement and feels like the wrong abstraction for a parse failure.

3. **Kernel validates before dispatching, retries inline in the runner** — extend the inference runner's existing retry logic. If the model's response fails to parse as a valid impulse, treat it as a retriable error. Retry carries the original present plus the parse failure reason. No new beings, no protocol changes — failure handling stays at the inference layer where HTTP retries already live.

4. **Append correction to system prompt on retry** — similar to option 3 but the error goes in the system message on the retry call rather than the user message. Keeps the present clean. Model sees: `your previous response was invalid: <reason> — respond again`. Separates correction from context.

**Resolution**

Implemented in `dispatch` (`main.go`). Tracks the last signal a cognitive being produced and its name. On drop, if the dropping being matches, re-runs inference with a short correction:

```
you wrote this last time:
<the bad signal>
your last signal was dropped: <reason> — try again
```

No new beings. No protocol changes. One correction shot — state clears after retry so a second consecutive drop stops the loop.

**Observation from first live run (2026-04-13)**

When prefrontal got stuck ping-ponging with thalamus, it eventually fired `skyra conflict | thalamus: escalation requested due to perceived redundancy in signals`. Malformed (no expression), but the intent was right — the model detected the loop and tried to escalate. The correction path lets it fix the syntax and fire again instead of dying silently.

---

## Issue 3: No backpressure to respond to external beings

**Symptom**
Prefrontal receives an inbound signal from michael (via thalamus) and loops internally — routing between strategy, conflict, values — indefinitely. It never routes to premotor. Nothing in the system tells prefrontal that michael is waiting, that a response is expected, or that the exchange needs to terminate outward. The loop is not a logic error — it is the correct behavior for a system with no backpressure. Prefrontal has no reason to stop.

**Root cause**
The protocol has no concept of an external being waiting for a response. Michael sends a signal in. That signal passes through sensory and thalamus and arrives at prefrontal as a present. But the chain of custody is lost — prefrontal sees a signal, not a request from an external being. Nothing in the present says "someone outside is waiting." Nothing in the system enforces that an inbound signal from an external being must eventually produce an outbound signal through motor. The cognitive layer can run forever without violating any constraint.

**Candidate solutions**

1. **Mandatory `~about` flag** — signal carries what it is responding to. `skyra premotor hello back ~about hi | ready to express`. Kernel validates `~about` references content from the current open exchange. Model cannot fire freely — it must name what it is processing. Lightweight protocol change, enforces grounding at the syntax level.

2. **Typed input/output on beings** — each being declares what types it can receive and emit. Theory-of-mind receives `raw_signal`, emits `attributed_contact`. Prefrontal receives `attributed_contact`, emits `judgment`. Premotor receives `judgment`, emits `shaped_expression`. Kernel rejects signals whose emitted type does not match the target being's declared input type. Misrouting becomes a type error. Architecturally rich but requires significant new infrastructure.

3. **Grounded expression requirement in the present** — the present explicitly separates "what you received" from "your network" and requires the model to carry the received content forward. If the exchange has `thalamus: hi`, the emitted signal must carry `hi` or a derivative. No protocol change — enforced through how the present is constructed and the instruction given to inference. Lightest to implement.
