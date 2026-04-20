# Public vs Private Nature

## The Problem

Identity and purpose are currently treated as public. They flow directly into other beings' presents — a being can see its sender's identity and purpose, and each peer's identity in its network section.

That is wrong. Identity is how a being understands itself from the inside. Purpose is what it is for. Neither belongs in another being's present.

## What Is Private

- Identity — the being's own self-description. Visible to itself only.
- Purpose — what the being is for. Visible to itself only.

A being's own identity and purpose appear in its present as self-knowledge. They do not cross the relationship boundary.

## What Is Public

- Name — the routing surface. Always public. The kernel enforces uniqueness.
- Callable — what the being advertises about when to reach for it. This is the public interface. The only thing a peer needs to know about another being before signaling it.

`Callable` is already in the code on `Nature`. It is already surfaced in `DerivePresent` as "deliberate when: ...". It is already the right field. It just needs to be the only thing shown for peers — not identity, not purpose.

## What Changes In DerivePresent

**What a being sees about itself** — identity and purpose stay. That is private self-knowledge. Correct as-is.

**What a being sees about its sender** — drop identity and purpose. Show name and callable only.

**What a being sees about its peers in the network** — drop identity. Show name and callable only.

## Relationship To Theory Of Mind

Doc 01 operational invariants states: "Before lived relationship experience exists, another being may inspect only identity and purpose."

This needs revisiting. If identity and purpose are private, theory-of-mind cannot expose them either. What theory-of-mind exposes is the public slice — name and callable. Perceived nature over time may deepen that picture through retained experience, but the starting point is name and callable, not identity and purpose.

## Open Questions

- Does callable need to be richer to carry enough signal for a being to know when to reach for a peer?
- How does perceived nature — a being's interpretation of a peer built through lived experience — surface in the present if not through identity and purpose?
- What does theory-of-mind actually return if not identity and purpose?
