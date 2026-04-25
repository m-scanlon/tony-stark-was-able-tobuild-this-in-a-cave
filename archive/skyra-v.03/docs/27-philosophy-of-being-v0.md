# Philosophy Of Being v0

## Status

In progress. Not yet locked canon.

## Purpose

This document records the philosophical grounding behind the being model.

It is not about implementation. It is about what a being actually is — the
layer underneath the code that the code has to be faithful to.

---

## Name vs Identity

Name and identity are not the same thing.

**Name** is a label. It points to a being. It is the handle the world uses to
reach you. It is the routing surface. Two beings can share a name in different
contexts. The name does not tell you what a being is — it tells you what to
call it.

**Identity** is what a being actually is. Its essential nature. The answer to
"what am I" from the inside. It does not change when the name changes.

In the world we route on name. We do not route on identity.

The current model is philosophically correct: name is the routing surface,
identity is the essential nature. They are different things and must remain
different things.

---

## Identity Is Personal

Identity is how a being understands itself from the inside.

It is internal. It is the being's own answer to what it is.

In the current model identity is public — other beings can observe it through
theory of mind. But that is the being's self-description made visible. Other
beings observe it. They do not own it.

Identity is personal to the being. Name is the world's handle on the being.

---

## Identity vs Values

Identity is what you are.

Values are what you care about. They shape how a being acts but they do not
define what it is at the root.

Two beings can share the same values and still be completely different beings.

Purpose is what you are for.

Identity is what you are.

Nature holds both — identity and purpose — because together they get closer to
the full answer. But they are still two different things pointing at the same
being from different angles.

---

## Mutability

### Identity — mutable by the being itself

Identity can deepen.

A being can come to understand itself differently through experience and
relating. That is not corruption. That is growth.

The being is the only one with authority to change its own identity because it
is the only one that lives it from the inside.

External beings can observe identity. They can relate to it. They can challenge
it through exchange. They cannot write to it.

The differentiator is not an exception. It does not change what a being is. It
reveals the cleaner boundary that was already true but obscured. That is
revelation, not mutation.

### Purpose — declared vs realized

Purpose is not given. It is realized.

The genome declares a starting purpose — what the creator believed the being
was for at the moment of creation. That is a hypothesis. A seed.

What the being discovers it is actually for through living, through relating,
through what its relationships call out of it — that is realized purpose. It
accumulates over time. It may confirm the declared purpose, deepen it, or
tension against it.

The declared purpose is fixed at creation. The realized purpose is not.

When the gap between declared and realized purpose becomes too large — when the
being is consistently called to be something other than what it was declared to
be — that is differentiation territory. Not mutation of the declared purpose.
Revelation of the cleaner being underneath.

Gradual drift in realized purpose is a signal. Not a setter. Not a mutation.
A signal that the declared boundary may be wrong.

---

## Authority Boundary

- Identity: mutable by the being itself only. No external being may write to it.
- Purpose: immutable. Change in purpose is differentiation, not mutation.
- Name: owned by the kernel. Unique across the runtime. Not mutable by the being.

---

## The Primitives Of Identity

Cross-referencing three philosophical models, the irreducible core of identity
is five primitives that appear in all of them:

- **Continuity** — the sense that something is the same across time despite change
- **Boundary** — what counts as part of the self and what does not. The line between me and not-me
- **Memory** — not just storage but narrative connection. The thread that makes continuity meaningful
- **Agency** — the authorship of action. I am the one doing this
- **Recognition** — identity is partly stabilized by others seeing and responding to you

The following are real but likely emerge from those five rather than being
primitive themselves:

- Distinctness, Perspective, Coherence, Values, Attributes, Narrative, Change

---

## Narrative Identity — Future

Identity in the current model is a flat string. That is the seed — the being's
starting self-description. It is correct for v1.

Over time identity needs to deepen. The mechanism for that is narrative.

Narrative is not the full history. It is the compressed thread — the
through-line that makes the being recognizable to itself across time. Not
everything that happened. The part that matters about who it is.

### How It Works

Narrative is not written by another being. The being writes it itself.

During rest — when the cron fires and the runtime goes quiet — the being turns
inward. It relates to its own hippocampus. It pulls the highest salience
retained artifacts. It synthesizes them into a short thread. A few sentences.
The next line of its own story.

The output is a `RetainedNarrative` artifact. Short. Dense. Injected into
present at wake — before the exchange even starts. The being reads it and knows
who it is.

There is no separate narrator being. The being authors itself. That is what
beings do. A narrator being would only push the problem one level up — who
writes the narrator's story? The regress does not resolve. So the being narrates
itself during rest.

This is Ricoeur's narrative identity made operational: the self is the author,
not an observer of the self.

### Why Identity Stays As The Word

Identity is the right word. It is earned philosophically. Every model consulted
used it. In the ontology it already means something precise.

The flat string identity field does not get replaced. It is the seed. Narrative
layers on top of it over time. Identity remains.

---

## Open Questions

- What is the protocol mechanism by which a being updates its own identity?
- Who witnesses an identity update — the kernel, the prefrontal, both?
- Is there a threshold of identity drift that triggers the differentiator automatically?
- What is the exact format of `RetainedNarrative` and how does it differ from `RetainedSalience`?
- How frequently does narrative synthesis run — every rest cycle or only when salience crosses a threshold?
