# Skill Lifecycle

Skills are not defined. They are learned. This document describes the full lifecycle — from the first signal in a fresh system to a provisioned, executable skill.

---

## Cold Start

A fresh system has no domains. No skills. No history.

Day zero Skyra ships with:

- `skyra.user` — committed facts about the person. Always present, first-injected every session.
- **System primitive skills** — pre-provisioned in Redis at boot: `reply`, `fan_out`, `report`, `chat`, `reasoning`, `integrate`, `update_skill`, `commit`, `propose_commit`, `search`, `provision_memory`, `provision_skill`

That is a complete, functional system. Everything else grows from observation.

---

## Signal Flow

Signals accumulate in the observational layer. When a domain exists that matches the signal, a `belongs_to` edge is added. If no domain matches, the node exists without one — domain is retrieval scaffolding, not a structural requirement.

```
signal arrives
  → observational node created
  → does a domain exist for this?
      yes → belongs_to edge added
      no  → node exists without a domain — valid state
```

**RAG lookup is signal.** When the LLM needs context that isn't present, it fires a manual lookup. The act of needing to search is itself an observational data point — vectors strengthen, the pattern is recorded.

---

## The Cron Pass

When the user is offline, the Cron Service fires a background reasoning skill. Skyra reads unprocessed session history and VAD time series, married on `turn_id`, and reasons them into observational nodes and edges. This is a graph mutation event — Skyra is appending to the graph.

Nodes produced here have no domain yet. Domain edges are added as retrieval structure emerges from clustering.

---

## Domain Proposal

When observational nodes cluster into a coherent pattern around a recognizable area of the user's life:

```
cluster crosses threshold
  → system proposes domain to user
  → "I've noticed a lot of activity around your fitness routine. Want me to track that as its own domain?"
  → user approves → domain node created, belongs_to edges added to cluster
  → user declines → threshold resets, nodes continue accumulating
```

---

## Skill Crystallization

**The boundary is intent.** Skills do not emerge from arbitrary graph patterns — recurring facts, recurring behaviors, recurring entities. None of those are skills. A skill candidate is a recurring intent — something the user has expressed a desire to do, repeatedly, that the system has not yet crystallized.

The signal is already in the graph: **the intent namespace**. When the user expresses an intent and no skill exists to handle it, an intent namespace is provisioned as an observational node. That node accumulates across sessions. When Skyra observes that the same intent namespace has been provisioned repeatedly and never crystallized, that is the signal.

```
intent node: log_workout
  → appears in session 1, session 3, session 4, session 7
  → no skill node covers this intent
  → Skyra proposes crystallization
```

When the pattern is clear, Skyra proposes the skill — already formed as a natural language definition.

```
"I've noticed you log your workouts fairly consistently. Want me to set that up as a skill?"

user approves
  → skill node promoted to committed layer
  → memory namespace provisioned alongside it
  → skill provisioned in Redis
  → executable
```

---

## Summary

```
Day 0:
  skyra.user + system primitives
      ↓
  user interacts → session history + VAD accumulate
      ↓
  cron fires → session history + VAD married on turn_id
  → observational nodes + edges appended to graph
      ↓
  nodes cluster → domain identified → proposed → approved
  → domain node created, belongs_to edges added
      ↓
  patterns cross threshold → skill proposed → approved
  → skill node promoted → Redis provisioned → executable
```

---

## The Full Model

```
intent (input)
  → skill crystallizes from recurring intent node
  → skill is observational — not trusted
  → Skyra reasons freely inside the skill boundary
  → output produced
  → propose_commit → user approves → committed
```

**Intent is the input.** It bounds the problem. The skill is shaped by it.

**The skill is not trusted.** It lives in the observational layer. Skyra can reason and revise freely inside it — no guardrails on the reasoning, unconstrained intelligence inside the intent boundary. This is Principle 3: constrain the data, not the model.

**The output is the finished product.** The user's intent realized. Not the reasoning that produced it — the output. That is what gets proposed. That is what lands in the committed layer on approval. The committed layer holds finished products, not working state.

The user never sees Skyra's working model. They see the output.

---

## Bounding

- No domain for a node → node exists without one, valid state
- Domain exists, no skill covers the intent → accumulates as observational nodes, intent namespace provisioned
- Skill exists → routes to skill

The system cannot create unbounded namespaces. Domains are approved by the user. Skills are proposed by the system and approved by the user.

---

## Deprovisioning

Does not exist. Once a skill is provisioned, it stays provisioned.

---

## System Primitive Skills

Pre-provisioned in Redis at boot. The system cannot function without these.

| Skill | Purpose |
|---|---|
| `chat` | A conversation with the user. Every session is a job. Opens on first turn, closes on session end. |
| `update_skill` | The only path to modifying a skill node. Requires user approval. |
| `reasoning` | Background job triggered by cron. Decomposes session history + VAD into observational nodes, then writes edges to the graph. |
| `integrate` | Connects the mini graph from reasoning to the existing graph. Finds aliases, updates weights, adds missing edges. |
| `commit` | Write to memory (user-gated) |
| `propose_commit` | Surface a commit proposal to the user |
| `search` | Semantic search in memory — retrieval and signal |
| `provision_memory` | Create a new memory namespace |
| `provision_skill` | Add a skill to Redis |

---

## Related

- `docs/arch/v1/kernel.md` — kernel pattern recognition, memory provisioning
- `docs/arch/v1/gaps.md` — G28 (memory provisioning flow), G29 (skill learning threshold)
- `docs/arch/v1/capability-model.md` — Redis skill registry, execution gate
