# Memory Spec

## What This Is

Memory is not stored on a being. Memory is a being.

Every being in the world has a paired `remember` node. When a being needs to retain, compress, or recall — it routes to its own `remember` node the same way it routes to any other Logos. `remember` is recursive — it manages its own memory using the same operations it provides to others.

---

## The Remember Being

`remember` is a single Logos node in the world. All beings route to it. It knows who is asking by the origin on the Relation — `r.Origin` is the being whose memory is being operated on.

Because `remember` sees all beings, it can observe patterns across the whole world — not just one being's retained experience.

```
skyra remember ~retain <artifact> | <reason>
skyra remember ~recall ~about <topic> | <reason>
skyra remember ~compress | <reason>
skyra remember ~forget ~before <timestamp> | <reason>
```

---

## Operations

### Retain

A being routes an artifact to its remember node to make it permanent.

```
skyra remember.<name> ~retain <artifact> | <reason>
```

Retained artifacts have a type:

- **trace** — a record of what happened in an exchange
- **salience** — something noticed as significant
- **tension** — an unresolved question or contradiction
- **understanding** — a conclusion reached, something now known

The type determines how the artifact is weighted during recall.

---

### Recall

A being asks `remember` for relevant context. Recall can be filtered by any combination of fields.

```
skyra remember ~recall ~about <topic> | <reason>
skyra remember ~recall ~from <being> | <reason>
skyra remember ~recall ~type <trace|salience|tension|understanding> | <reason>
skyra remember ~recall ~after <timestamp> | <reason>
skyra remember ~recall ~before <timestamp> | <reason>
skyra remember ~recall ~strength <min> | <reason>
```

Filters compose. A being can recall all tension artifacts from a specific relationship after a given date.

`remember` returns the most relevant retained artifacts for the query. Without filters, relevance is weighted by type and recency — understanding artifacts surface first, traces surface last.

**Recall fields:**

| Field | Description |
|---|---|
| `~about` | semantic topic match across artifact cores |
| `~from` | filter by the being the artifact involves (relationship) |
| `~type` | filter by artifact type: trace, salience, tension, understanding |
| `~after` | artifacts retained after this timestamp |
| `~before` | artifacts retained before this timestamp |
| `~strength` | minimum strength threshold (0.0 – 1.0) |

---

### Compress

A being triggers compression when its exchange history is too large to carry.

```
skyra remember ~compress | <reason>
```

`remember` reads the origin being's exchange history, derives understanding artifacts from it, retains those, and discards what no longer needs to be carried. The being doesn't decide what to keep — `remember` does. Memory quality is not measured by how much can be recalled. It is measured by how little the present being needs to carry without breaking.

### Forget

A being deletes its retained artifacts before a given timestamp.

```
skyra remember ~forget ~before <timestamp> | <reason>
```

All artifacts retained by the origin being with a timestamp older than the given value are permanently removed. This is the only delete operation in the memory system.

---

## Remember Routes to Itself

`remember` is not exempt from its own operations. It retains its own compression decisions, recalls its own history when asked, and compresses itself when it grows too large.

---

## Topology

`remember` is a single being registered in the world at startup:

```
skyra grow ~name remember ~identity I retain and compress experience for all beings in this world | genome
```

All beings route to it by name. The origin on the Relation identifies whose memory is being operated on.

---

## Retained Experience Artifact Family

Every artifact retained by `remember` has a type. The type determines how it is weighted during recall and whether it survives compression.

### Trace

A record of what happened. The raw shape of an exchange — who said what, what was routed, what responded. Traces are the lowest-weight artifact. They are the first to be compressed away.

```
type: trace
core:           the exchange or event
beings:         who was involved
thread:         thread ID if applicable
strength:       how significant the exchange was
```

### Salience

Something noticed as significant in the moment. A detail that stood out, a pattern that emerged, something that felt important but is not yet understood. Salience artifacts are mid-weight — they survive one compression cycle before they either become understanding or are discarded.

```
type: salience
core:           what was noticed
context:        what was happening when it was noticed
strength:       how strongly it registered
```

### Tension

An unresolved question or contradiction. Something the being is holding that has not been settled. Tension artifacts have high weight — they surface early in recall and persist through compression until resolved.

```
type: tension
core:           the unresolved question or contradiction
origin:         where it came from
strength:       how much it is still active
```

### Understanding

A conclusion reached. Something now known. The highest-weight artifact — it survives indefinitely and surfaces first in recall. Understanding is the output of compression: traces and salience collapse into understanding over time.

```
type: understanding
core:           the conclusion
derived_from:   trace or salience artifacts that produced it
strength:       confidence in the understanding
```

---

## What This Is Not

- `remember` does not store raw conversation history. That lives in the exchange map on the being. `remember` stores derived understanding.
- `remember` does not decide what a being believes. It surfaces what the being has retained. The being reasons over it.
- Memory is not a database. It is a being with judgment.
