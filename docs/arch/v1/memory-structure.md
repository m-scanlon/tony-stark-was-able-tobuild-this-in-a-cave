# Memory Structure

Skyra's memory is a property graph. Nodes are things that exist. Edges are relationships between them. Both carry vector embeddings. Edges carry Skyra-owned weights that evolve over time.

---

## Node

```
node {
  id:       string                  // stable canonical identifier
  type:     entity | fact | skill | domain | long_term_memory
  content:  string                  // what this node is
  vector:   float[]                 // semantic embedding — for search
  metadata: {
    source:       observational | authoritative
    confidence:   0.0 to 1.0
    created_at:   timestamp
    last_seen_at: timestamp
  }
}
```

### Node Types

**entity** — a named thing in the user's life. Person, place, tool, concept. Accumulates through observation. Has aliases that all resolve to the same canonical id.

**fact** — a committed piece of knowledge. Requires user approval to enter the authoritative graph. High trust.

**skill** — a learned class. Roadmap, contract, validation criteria. Lives in memory as a node — discoverable via semantic search. Execution gated by Redis.

**domain** — a scoped area of the user's life. gym, work, servers, life. Contains skill partitions and entity relationships. Proposed by the system, approved by the user.

**long_term_memory** — promoted synthesis. Not raw observations. A synthesized conclusion that crossed the frequency × affect threshold. Proposed to the authoritative graph after promotion.

---

## Edge

```
edge {
  from:    node_id
  to:      node_id
  type:    string                   // what kind of relationship
  vector:  float[]                  // semantic embedding of the relationship
  weight:  float                    // Skyra's importance assessment — not user-controlled
  history: [
    { weight: float, at: timestamp },
    ...
  ]
}
```

**Skyra owns the weight.** The user does not set it. Skyra derives it from the observational streams, the decay formula, and pattern recognition. Weight updates do not require user approval.

**The history is data.** A weight shifting from 0.9 to 0.4 over six months tells Skyra something about how the user's life is changing. The shift is as meaningful as the current value.

### Edge Types

| type | meaning |
|---|---|
| `motivation` | why the user does something |
| `belongs_to` | entity or skill belongs to a domain |
| `relates_to` | general semantic relationship |
| `supports` | observation supports a fact or conclusion |
| `part_of` | task is part of a skill, skill is part of a domain |
| `alias_of` | one entity name resolves to another |
| `causes` | one thing leads to another (behavioral pattern) |

Edge types are not exhaustive — the system adds types as relationships are observed.

---

## Two Kinds of Writes

**User-gated commits** — facts, domain proposals, skill proposals, long term memory promotions. Skyra proposes. User approves. Nothing enters the authoritative graph without this handshake.

**Skyra-owned weight updates** — edge weights updated by pattern recognition and the decay formula. No user gate. Transparent and inspectable but Skyra's to manage.

```
user commit:       node enters authoritative graph
Skyra weight update: edge weight shifts, history appended
```

---

## The Entity-Domain Matrix

The entity-domain matrix `D[i,j]` from the predictive memory model is the graph's edge layer between entity nodes and domain nodes.

```
D[nginx_config][servers] = { DLT: 87, DST: 72 }

maps to:

edge {
  from:   "entity:nginx_config"
  to:     "domain:servers"
  type:   "belongs_to"
  weight: 87   // DLT — Skyra's long-term importance assessment
  history: [...]
}
```

Entities earn their way into domains through usage. A missing `D[i,j]` pair means the edge doesn't exist yet — not that the importance is zero.

---

## Domain Structure

Every domain is a subgraph — a cluster of nodes and edges scoped to that area of the user's life.

```
domain: gym
  nodes:
    skill:log_workout
    skill:cancel_session
    entity:gym_membership
    entity:personal_record_bench
    fact:"Mike trains 4x per week"
  edges:
    log_workout → belongs_to → gym
    personal_record_bench → belongs_to → gym
    log_workout → relates_to → personal_record_bench
    ...
```

Skills are nodes in the domain. Their partitioned data lives as child nodes and edges within the domain subgraph.

---

## The Life Domain

The life domain is a domain like any other — same structure, same rules, same commits. It is always provisioned. It is the fallback for nodes and edges that have no other domain yet.

Tools and skills without a domain live here. Entity clusters that haven't resolved into a domain yet live here. The scratchpad for new domain discovery lives here.

---

## Scratchpad

Every domain has a scratchpad — the system's private reasoning workspace. Not user-visible. Not part of the authoritative graph. Skyra writes freely.

The scratchpad is where intent patterns accumulate before they crystallize into skill proposals. It is not versioned in the same way as the authoritative graph — it is working memory, not committed memory.

---

## Related

- `docs/arch/v1/predictive-memory.md` — observational streams, entity weights, decay formula
- `docs/arch/v1/skill-lifecycle.md` — how nodes enter the graph (domain proposal → skill proposal)
- `docs/arch/v1/kernel.md` — pattern recognition, memory provisioning
