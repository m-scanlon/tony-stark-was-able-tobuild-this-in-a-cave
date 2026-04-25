# Memory

## What This Is

Memory is not stored on a being. Memory is a being.

`remember` is a single being registered in the world at startup. All beings route to it. It knows who is asking by the origin on the Relation — `r.Origin` is the being whose memory is being operated on. Because `remember` sees all beings, it can observe patterns across the whole world — not just one being's retained experience.

`remember` is not exempt from its own operations. It retains its own compression decisions, recalls its own history when asked, and compresses itself when it grows too large.

---

## Two Access Patterns

### Implicit — via the internal self

The internal self routes to `remember` on every turn to retrieve retained experience relevant to the incoming message. The being never chooses this. It happens automatically as part of the internal self's deliberation. This is involuntary recall — the subconscious surfacing what it knows.

### Explicit — via the being directly

A being can route to `remember` on its own. Retain an artifact. Recall something specific. Trigger compression. Forget. These are conscious operations — the being decides to interact with its own memory.

Both paths hit the same store. One is automatic, one is intentional.

---

## Operations

### Retain

A being routes an artifact to `remember` to make it permanent.

```
<> remember ~retain <artifact>
```

Retained artifacts have a type:

- **trace** — a record of what happened in an exchange. Lowest weight. First to be compressed away.
- **salience** — something noticed as significant. Mid-weight. Survives one compression cycle before it either becomes understanding or is discarded.
- **tension** — an unresolved question or contradiction. High weight. Persists through compression until resolved.
- **understanding** — a conclusion reached, something now known. Highest weight. Survives indefinitely.

### Recall

A being asks `remember` for relevant context.

```
<> remember ~recall ~about <topic>
<> remember ~recall ~from <being>
<> remember ~recall ~type <trace|salience|tension|understanding>
<> remember ~recall ~after <timestamp>
<> remember ~recall ~before <timestamp>
<> remember ~recall ~strength <min>
```

Filters compose. A being can recall all tension artifacts from a specific relationship after a given date. Without filters, relevance is weighted by type and recency — understanding surfaces first, traces surface last.

### Compress

A being triggers compression when its exchange history is too large to carry.

```
<> remember ~compress
```

`remember` reads the origin being's exchange history, derives understanding artifacts from it, retains those, and discards what no longer needs to be carried. The being doesn't decide what to keep — `remember` does.

Memory quality is not measured by how much can be recalled. It is measured by how little the present being needs to carry without breaking.

### Forget

A being deletes its retained artifacts before a given timestamp.

```
<> remember ~forget ~before <timestamp>
```

All artifacts retained by the origin being older than the given value are permanently removed. This is the only delete operation in the memory system.

---

## Voluntary Orientation

A being can ground itself in any of its exchange perspectives without sending a message.

```
<> remember <peer>
```

The being lands on its own perspective of the exchange with that peer. No message is sent. The target is the being itself — it is just grounding itself in that exchange context. The being then emits normally from there.

This lets a being:
- Check in on a dormant exchange before deciding to act
- Recall context before opening a new exchange with someone else
- Decide not to act after orienting — silence is a valid outcome

"Remember" is the read. Every other emission is a write. This is the conscious counterpart to the internal self's involuntary recall.

---

## Retained Artifact Family

Every artifact retained by `remember` has a structure:

**Trace**
```
type: trace
core:           the exchange or event
beings:         who was involved
thread:         thread ID if applicable
strength:       how significant the exchange was
```

**Salience**
```
type: salience
core:           what was noticed
context:        what was happening when it was noticed
strength:       how strongly it registered
```

**Tension**
```
type: tension
core:           the unresolved question or contradiction
origin:         where it came from
strength:       how much it is still active
```

**Understanding**
```
type: understanding
core:           the conclusion
derived_from:   trace or salience artifacts that produced it
strength:       confidence in the understanding
```

---

## What This Is Not

- `remember` does not store raw conversation history. That lives in the exchange map on the thread. `remember` stores derived understanding.
- `remember` does not decide what a being believes. It surfaces what the being has retained. The being reasons over it.
- Memory is not a database. It is a being with judgment.