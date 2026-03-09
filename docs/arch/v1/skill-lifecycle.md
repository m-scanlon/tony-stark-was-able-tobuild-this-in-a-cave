# Skill Lifecycle

Skills are not defined. They are learned. This document describes the full lifecycle — from the first signal in a fresh system to a provisioned, executable skill with its own partition in domain memory.

---

## The Memory Hierarchy

Every domain is architecturally identical:

```
domain/
  shared data        ← all data for this domain, accessible to all skills inside it
  scratchpad/        ← system's reasoning workspace, not user-visible
  skill partitions/  ← crystallized from the scratchpad, proposed to user when ready
    log_workout/
    cancel_session/
    ...
```

The **life namespace** is a domain like any other — shared data, scratchpad, skill partitions, commits. It is not architecturally special. It exists from day zero and is the fallback for anything that doesn't belong to a more specific domain yet.

```
life/               ← always present, first domain
gym/                ← provisioned when system identifies fitness cluster
work/               ← provisioned when system identifies work cluster
servers/            ← etc.
```

---

## Cold Start

A fresh system has no domains except life. No skills. No history.

Day zero Skyra ships with:

- `skyra.user` — committed facts about the person. Always present, first-injected every session.
- **Life namespace** — the first and always-present domain. Fallback for all unclassified signals and tools.
- **System primitive skills** — pre-provisioned in Redis: `commit`, `propose_commit`, `search`, `provision_memory`, `provision_skill`
- **Hardcoded kernel primitives** — `reply`, `fan_out`, `report`

That is a complete, functional system. Everything else grows from observation.

---

## The Life Namespace

The life namespace is a domain. It follows the same rules as every other domain — shared data, scratchpad, skill partitions, commits.

It serves two roles by virtue of being the fallback:

**1. Home for tools without a domain.**
When a tool doesn't belong to any provisioned domain yet, it lives in the life namespace. This is not special treatment — life namespace is just the catch-all domain. As domains get provisioned, tools migrate to where they belong.

**2. Birthplace of new domains.**
The life namespace scratchpad is where the system reasons freely over the user's intent before specific domains exist. Signal clusters in life → domain identified → proposed to user → new domain provisioned.

---

## Signal Flow

Every signal is assigned to the most specific domain it belongs to. If no domain matches, it falls to life.

```
signal arrives
  → does a domain exist for this?
      yes → routes to that domain
      no  → falls to life namespace
```

Over time, the life namespace scratchpad accumulates clusters. The system reasons over them and identifies candidate domains. When a cluster is ready, a new domain is proposed.

**RAG lookup is signal.** When the LLM needs context that isn't present, it fires a manual lookup. The act of needing to search is itself an observational data point — vectors strengthen, the pattern is recorded.

---

## Domain Proposal

When the system identifies a domain cluster in the life namespace scratchpad:

```
cluster crosses threshold
  → system proposes domain to user
  → "I've noticed a lot of activity around your fitness routine. Want me to track that as its own domain?"
  → user approves → domain memory provisioned
  → user declines → threshold resets, signals continue accumulating in life
```

---

## Domain Scratchpad → Skill Partitions

Every domain has a scratchpad. It is the system's private reasoning workspace — not user-visible, not committable by the user.

The system reasons over the scratchpad continuously, working toward **partitions**. For each candidate partition, it identifies the intent and defines requirements — what data, what pattern confidence, what execution criteria need to be fulfilled before this partition is ready to propose.

```
gym scratchpad:
  → system identifies: log_workout
      requirements: workout type pattern, frequency signal, duration data
  → system identifies: cancel_session
      requirements: session lookup pattern, cancellation signal
```

When a partition's requirements are fulfilled, the system proposes it to the user as a skill — already formed.

```
"I've noticed you log your workouts fairly consistently. Want me to set that up as a skill?"

user approves
  → skill partition crystallizes in domain memory
  → first memory commit
  → skill provisioned in Redis
  → executable
```

---

## Summary

```
Day 0:
  skyra.user + life namespace + system primitives
      ↓
  signals accumulate in life namespace
  system reasons freely over life scratchpad
      ↓
  domain cluster identified → proposed → approved
      ↓
  domain provisioned → domain scratchpad opens
  system reasons over domain scratchpad → partitions form
      ↓
  partition requirements fulfilled → skill proposed → approved
      ↓
  skill crystallizes → Redis provisioned → executable
```

---

## Bounding

Namespace creation is bounded by the domain hierarchy:

- No domain for a signal → falls to life namespace
- Domain exists, no skill covers the intent → accumulates in domain scratchpad
- Skill exists → routes to skill partition

The system cannot create unbounded namespaces. Every signal has a home. Every intent partition belongs to a domain. Domains are approved by the user.

---

## Deprovisioning

Does not exist. Once a skill is provisioned, it stays provisioned.

---

## System Primitive Skills

Pre-provisioned in Redis at boot. The system cannot function without these.

| Skill | Purpose |
|---|---|
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
