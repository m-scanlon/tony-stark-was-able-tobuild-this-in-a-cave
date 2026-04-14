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

---

## Issue 4: Spinning without resolving

**Symptom**
Beings fire the same signal into the same exchange repeatedly — "continue", "hi", "proceed" — without the exchange going anywhere. This is distinct from the backpressure problem. The issue is not that the system fails to route outward. The issue is that it keeps moving while going nowhere. It thrashes instead of holding.

**Root cause**
There is no constraint against repetition. A being can fire the same expression at the same peer indefinitely without the system recognizing that nothing is resolving. Paralysis is acceptable — a being that cannot resolve an exchange and goes quiet is behaving honestly. Spinning is not — it is noise that looks like progress.

**Candidate solutions**

1. **Repetition detection in dispatch** — dispatch tracks the last N expressions a being produced. If the same expression fires twice in a row to the same target, the loop stops. No new protocol, no new beings. Cheap floor on the worst case.

2. **Stillness as a valid signal** — introduce a `~hold` flag. A being that cannot resolve an exchange can fire `skyra <self> ~hold | unresolved` instead of continuing to emit. Dispatch treats `~hold` as a terminal state for that exchange turn. The being went quiet honestly rather than spinning.

3. **Exchange entropy in the present** — the present surfaces how many times the current exchange has turned without a new expression appearing. If the turn count is high and the expressions are similar, the model is shown this explicitly. Forces self-awareness about the loop before the next signal fires.

---

## Issue 5: Beings have no shared account of the external world

**Symptom**
When michael speaks, the signal arrives at prefrontal as an exchange entry — something that happened in a peer channel. The external origin is lost. Prefrontal cannot deliberate over what michael said because it only sees a signal in an exchange, not a fact about the world. When it routes internally to strategy or conflict, those beings have even less — they see whatever prefrontal chose to pass, not what michael actually said. The cognitive network is not oriented toward a shared external reality. It is passing notes about notes.

**Root cause**
The present has no external account layer. Each being sees only its own exchange history with its peers. There is no persistent, shared representation of what is happening outside — who said what, what is in the world right now, what the system is responding to. The external signal arrives, gets written into an exchange stack, and the external context dissolves into the relay chain.

**Why this matters**
People deliberate over external reality, not over internal routing decisions. When someone speaks to you, the thing you think about is what they said — not the fact that your auditory cortex fired. The external account is the shared ground that makes internal deliberation meaningful. Without it, the cognitive network has nothing real to be oriented toward together.

**Candidate directions**

1. **External present layer** — the present gains a third block alongside identity and exchange: a persistent account of the external world. What external beings have said, what is happening outside, what the system is currently responding to. All cognitive beings see the same external layer. Internal routing becomes deliberation over shared external reality rather than signal passing.

2. **External exchange always in prefrontal's present** — the simplest version: prefrontal always sees the full unmodified signal from the external being, regardless of how many relay hops it passed through. The chain of custody is preserved. Prefrontal is always oriented toward michael, not toward thalamus.

3. **Shared world state being** — a non-cognitive being whose job is to hold a running account of external events. Cognitive beings can address it to read the current world state. The external account becomes queryable rather than broadcast.

---

## Issue 6: No structured record of the boundary exchange

**Symptom**
There is no clean abstraction for what crosses the external boundary. Inbound signals from michael and outbound signals from motor are written into exchange stacks the same way internal signals are. The full arc of the external conversation — what michael said, what Skyra expressed, what michael said back — is not preserved as a distinct thing. Internal beings have no shared thread to deliberate over. They are working degraded copies of the original signal, not the same ground truth.

**Root cause**
Metaxu treats all signals uniformly. Nothing records boundary crossings as a distinct event. `DerivePresent` assembles the present from peer exchange stacks, which dissolves external origin and external response into the relay chain. The architecture has an implicit boundary — sensory inbound, motor outbound — but no explicit record of what crosses it.

**Why this matters**
Internal deliberation needs to be *about* the external exchange, not about internal routing decisions. If every cognitive being's present includes the full external conversation — what arrived, what was expressed — the cognitive network is oriented toward shared ground. Strategy, values, consequence are all working the same problem. The gap between what michael said and what Skyra last expressed creates natural pressure toward resolution without requiring formal backpressure mechanisms.

**Open questions**
- Does the full boundary history go into every present, or a window? Full history gets expensive as exchanges grow.
- Does `BoundaryExchange` live in `world` or `metaxu`? It is world state, not routing logic — probably `world`.
- How does `DerivePresent` render it? Likely as the top block, before identity and peer exchange — external reality first, then who you are, then what your peers said.

**Candidate directions**

1. **`BoundaryExchange` type baked into present derivation** — a new type in `world` with `RecordInbound`, `RecordOutbound`, and `Render() string`. Metaxu writes to it when signals cross the sensory or motor boundary. `DerivePresent` on every cognitive being includes the rendered boundary exchange as its top block. No new beings. No new routing. The external conversation is just always there.

2. **Windowed boundary record** — same as above but `DerivePresent` only renders the last N turns of the boundary exchange rather than the full history. Keeps the present from growing unbounded as the conversation lengthens. N is a tunable constant.

3. **Boundary exchange on prefrontal only** — simpler version: only prefrontal's present includes the boundary record. Other cognitive beings get it indirectly through what prefrontal passes them. Less correct than option 1 but lower implementation cost.
