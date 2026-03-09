# Memory Structure

Skyra's memory is a property graph. Nodes are things that exist. Edges are relationships between them. Both carry vector embeddings. Edges carry Skyra-owned weights that evolve over time.

---

## Two-Tier Graph

The graph has two trust levels. This is the foundation of data integrity.

```
committed layer       ← user-approved, high trust, authoritative, stable
        ↑
    promotion
        ↑
observational layer   ← Skyra writes freely, working model, low trust
```

**Observational layer** — Skyra's working model. She writes here freely. Fragments, signals, partial facts, inferred entities. Can be incomplete. Can be wrong. No user gate.

**Committed layer** — the source of truth. Nothing lands here without user approval. Stable, auditable, trusted. This is where data integrity lives.

Skyra reasons in the observational layer. Truth lives in the committed layer. The user only ever touches the committed layer.

---

## Node Promotion

When observational nodes cohere into something real — a cluster of fragments that together represent a fact, an entity, a pattern — Skyra proposes promotion.

```
Skyra accumulates observational nodes around an entity
  → fragments, signals, partial facts build up over time
  → nodes cluster and cohere
  → Skyra recognizes: this cluster represents something real
  → proposes the cluster for promotion
  → user approves
  → nodes become committed facts in the authoritative layer
```

A group of nodes representing the same entity can be promoted together. Promotion is atomic — the cluster lands as a unit.

Committed nodes are not re-evaluated. Once promoted and approved, they are stable. Corrections require an explicit new commit.

---

## Node

```
node {
  id:       string
  type:     entity | fact | skill | domain | artifact | long_term_memory
  layer:    observational | committed
  content:  string
  vector:   float[]
  ref:      { type, url/path }    // optional — for artifact nodes pointing to external data
  metadata: {
    confidence:   0.0 to 1.0
    created_at:   timestamp
    last_seen_at: timestamp
    promoted_at:  timestamp       // set when node moves from observational → committed
  }
}
```

### Node Types

**entity** — a named thing in the user's life. Person, place, tool, concept. Accumulates through observation. Has aliases resolving to one canonical id. Starts observational. Promoted when the cluster is coherent.

**fact** — a committed piece of knowledge. Always in the committed layer. Requires user approval.

**skill** — a learned class. Roadmap, contract, validation criteria. Discoverable via semantic search. Execution gated by Redis. Starts observational (intent namespace). Promoted when proposed and approved.

**domain** — a scoped area of the user's life. A node like any other. Proposed by Skyra, approved by user. Domains are nodes — containment is expressed through edges, not structural hierarchy.

**artifact** — a pointer to real digital data. A file, a git repo, a database. The node is the semantic representation. The `ref` field points to where the actual data lives. Multiple nodes can share the same ref. Artifact nodes can belong to multiple domains via edges.

**long_term_memory** — promoted synthesis. A synthesized conclusion from the pattern recognition engine. Not raw observations — the meaning, the pattern, what works.

---

## Edge

```
edge {
  from:    node_id
  to:      node_id
  type:    string
  vector:  float[]     // semantic embedding of the relationship itself
  weight:  float       // Skyra's importance assessment — not user-controlled
  history: [
    { weight: float, at: timestamp },
    ...
  ]
}
```

**Two nodes can have multiple edges between them.** Each edge has a different type. Your relationship to something is rarely just one thing.

```
mike → motivation  → skyra_project
mike → works_on    → skyra_project
mike → proud_of    → skyra_project
```

**Skyra owns the weight.** Derived from observational streams, decay formula, pattern recognition. Not user-controlled. Transparent and inspectable.

**The history is data.** A weight shifting from 0.9 to 0.4 over six months is as meaningful as the current value. It tells Skyra how the user's life is changing.

**Edges are not nodes.** They are a separate primitive. If you need to reason about a relationship as a thing, reify it — create a node that represents the relationship and draw edges to it.

**Edges cross domain boundaries freely.** Domains are nodes. There are no walls. The graph is one connected structure.

### Edge Types

| type | meaning |
|---|---|
| `motivation` | why the user does something |
| `belongs_to` | node belongs to a domain |
| `relates_to` | general semantic relationship |
| `supports` | observation supports a fact or conclusion |
| `part_of` | task part of skill, skill part of domain |
| `alias_of` | entity name resolves to another |
| `causes` | one thing leads to another |
| `same_ref` | two artifact nodes point to the same external resource |

Edge types are not exhaustive — new types emerge as relationships are observed.

---

## Artifact Nodes and Cross-Domain Refs

A single artifact can be referenced from multiple domains. The artifact node is one node. Multiple domain edges point to it.

```
node: skyra_repo (artifact)
  ref: { type: git, url: "github.com/m-scanlon/skyra" }

edge: skyra_repo → belongs_to → home     weight: 0.8
edge: skyra_repo → belongs_to → school   weight: 0.6
```

`home` and `school` are domain nodes. `skyra_repo` is an artifact node. Same ref. Two domain contexts. Different Skyra-owned weights because she's observed where the work actually happens.

Images, files, and other data live on the hard drive. Nodes point to them. The graph is the semantic layer — not the storage layer.

---

## Write Ownership

```
System provisions:   life domain, primitive skills
User approves:       new domains, new skills, node promotion (committed layer)
Skyra writes freely: observational nodes, edges, weights, scratchpad
```

Data integrity lives in the committed layer. Skyra's working model lives in the observational layer. The user owns the truth.

---

## The Entity-Domain Matrix

The entity-domain matrix `D[i,j]` from the predictive memory model maps directly to graph edges between entity nodes and domain nodes.

```
D[nginx_config][servers] = { DLT: 87, DST: 72 }

→ edge {
    from:    "entity:nginx_config"
    to:      "domain:servers"
    type:    "belongs_to"
    weight:  87
    history: [...]
  }
```

Entities earn their way into domains through usage. A missing pair means the edge doesn't exist yet.

---

## Domain Structure

Every domain is a subgraph — a cluster of nodes and edges. Domains are nodes themselves. Containment is expressed through `belongs_to` edges, not structural hierarchy.

Every domain has a **scratchpad** — Skyra's private reasoning workspace. Not user-visible. Not committed. Skyra writes freely. Where intent patterns accumulate before crystallizing into skill proposals.

---

## The Life Domain

A domain like any other. Always provisioned from day zero. The fallback for nodes that have no other domain yet. Tools without a domain live here. Entity clusters that haven't resolved into a domain yet live here. New domains are born from the life domain scratchpad.

---

## Related

- `docs/arch/v1/predictive-memory.md` — observational streams, entity weights, decay formula
- `docs/arch/v1/skill-lifecycle.md` — node promotion path for skills and domains
- `docs/arch/v1/kernel.md` — pattern recognition, memory provisioning
