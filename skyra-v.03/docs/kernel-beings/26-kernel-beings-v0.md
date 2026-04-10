# Kernel Beings v0

## Status

In progress. Not yet locked canon.

## Purpose

This document defines the beings the kernel needs to function.

These are not abstractions around the kernel. They are the kernel's own
population — the beings that make registration, routing, and growth possible
from within the ontology.

## Contents

- The grow being
- The registration path as a being operation
- What the kernel owns and what it delegates
- How kernel beings relate to each other at boot

---

## The Nature Being

### What It Is

The nature being is a non-cognitive primitive.

It holds exactly two fields: identity and purpose.

It is the minimum shape of every being in the system.

It does not reason. But it has explicit callable language so that any being
relating to it knows exactly how to speak to it.

### Callable Language

```
skyra nature ~identity <identity> ~purpose <purpose> | <source>: <reason>
```

- `~identity` — the identity field
- `~purpose` — the purpose field
- `<source>` — the being calling it
- `<reason>` — why it is being called

### Key Rules

- Nature is locked at creation. It is not mutable.
- The callable language on the relationship with nature is seeded at birth — not earned through use.
- Non-cognitive primitive beings are born with their callable language. Cognitive beings earn theirs through relating.
- The nature being's parser reads `~identity` and `~purpose` from the expression. That is its entire slice of the protocol.
