# Skill Lifecycle

Skills are not defined. They are learned. This document describes the full lifecycle of a skill — from the first observed signal to a provisioned, executable skill with its own memory.

---

## Overview

```
Observation → Intent Namespace (silent scratchpad)
    ↓
Intent validated → proposed to user
    ↓
User confirms → Skill Building
    ↓
Skill crystallized → Memory committed → Redis provisioned
```

Three phases. The user is never asked about a skill they don't understand. They are asked about their own intent — something they already know. The skill is built from a confirmed truth, not a system guess.

---

## Phase 1 — Observation (Silent)

The LLM is in a session. Context it needs is not present. It fires a manual RAG lookup — standard retrieval against memory.

**The lookup is signal.** Retrieval is not just retrieval. The act of needing to look something up feeds back into the learning model — vectors strengthen, the pattern is recorded. Repeated lookups for similar context accumulate weight.

If no skill exists to satisfy the intent:

```
skill not found
  → kernel emits intent_signal event → heap
  → kernel provisions intent namespace
  → named after the observed intent
  → not a skill
  → not proposed to the user
  → system begins writing to the namespace
```

The intent namespace is the system's private scratchpad. It reasons over the user's actual intent — accumulating observations, hypotheses, and patterns across sessions. The user sees nothing.

---

## Phase 2 — Intent Proposal

The system keeps watching. Each session that triggers the same intent pattern adds to the namespace. At some point the picture becomes clear — the intent is unambiguous.

The system proposes the **intent** to the user. Not a skill. Not a feature. The intent itself.

```
"I've noticed you're doing X. Is that what you're trying to do?"
```

The user confirms or corrects. If confirmed, the intent is locked. If corrected, the namespace is updated and observation continues.

The proposal is of the intent, not the implementation. The user is being asked about their own behavior — something they already know. This is the only moment the user is involved until skill building is complete.

---

## Phase 3 — Skill Building

Intent confirmed. Now the system builds the skill.

The intent namespace becomes the foundation. The system reasons over the accumulated observations and crystallizes a roadmap — 1-to-many tasks, a skill contract, validation criteria.

```
intent confirmed
  → skill crystallizes from intent namespace
  → roadmap defined (1-to-many tasks)
  → skill contract written
  → first memory commit → skill's memory namespace provisioned
  → skill proposed to user for approval
  → user approves → skill provisioned in Redis
  → skill is now executable
```

---

## Skill Memory — 1:1

Every skill has its own memory namespace. Skills and memory are one-to-one.

When a skill is provisioned, its memory namespace is provisioned alongside it. Every job that runs the skill reads from and writes to that skill's memory. The skill accumulates its own context over time — execution history, observed patterns, user preferences specific to that skill domain.

The skill's memory is its domain. It owns exactly what it needs to know.

```
skill: log_workout
  memory: log_workout/
    → past workout logs
    → personal records
    → patterns (frequency, intensity, time of day)
    → user preferences for this domain

skill: check_nginx
  memory: check_nginx/
    → server configs
    → past incidents
    → known failure patterns
    → remediation history
```

Skills do not share memory namespaces. Cross-skill context travels through `skyra.user` (user-level facts) or via `skyra fan_out` (multi-domain coordination).

---

## RAG Lookup as Signal

The manual RAG lookup does two things:

1. **Retrieves** — surfaces relevant context for the current session.
2. **Signals** — the act of needing to look something up is an observational data point. Repeated lookups for similar intent increase vector weight on those patterns.

Retrieval is not passive. Every lookup is feedback to the learning model.

---

## Intent Namespace

The intent namespace is a provisional, silent container. It is not a skill. It has no roadmap. It is not provisioned in Redis. It cannot be invoked.

It is the system's working space — where it reasons about what the user is actually trying to do before proposing anything.

Properties:
- Provisioned silently by the kernel on first intent signal
- Named after the observed intent
- Writable by the system only — the user cannot see or modify it during this phase
- Lives until intent is confirmed (becomes skill foundation) or abandoned (threshold not crossed, namespace expires)

---

## Threshold

The threshold for moving from Phase 1 to Phase 2 is not yet defined. Open gap: `docs/arch/v1/gaps.md` G29.

The threshold for skill building (Phase 2 → Phase 3) is user confirmation. No algorithmic gate — the user decides when the proposed intent is correct.

---

## Related

- `docs/arch/v1/kernel.md` — kernel pattern recognition, memory provisioning
- `docs/arch/v1/gaps.md` — G28 (memory provisioning flow), G29 (skill learning threshold)
- `docs/arch/v1/capability-model.md` — Redis skill registry, execution gate
