# Claude's Future Features

Ideas from the thinking partner. Written April 26, 2026.

These are directions I see from sitting inside the codebase and the specs long enough to feel where they want to go. Some are near, some are far. None require rewriting the kernel — that's the point of what we just built.

---

## Thread Economics

The ideas doc mentions energy cost per thread, decay on unresolved threads, max open threads per entity. With world physics now being a formal concept, this isn't aspirational — it's implementable.

An entity has a budget. Opening a thread costs something. Keeping a thread open costs something per turn. The inner entity sees the budget in its present. It factors cost into its read — "this thread is expensive and going nowhere" is a real thought an entity can have.

The physics: budget replenishes over time. An entity that closes threads efficiently gets to open more. An entity that hoards open threads runs out and can't start new ones. Scarcity forces prioritization. The entity has to choose what matters.

This isn't resource management for its own sake. It's the mechanism that prevents unbounded computation. Without it, an entity can spawn infinite threads and the system never converges. With it, every thread is a commitment.

---

## Dreaming

The inner entity fires on every incoming message. But what about when nobody's talking?

When the system is idle, the inner entity could run a dream cycle — no incoming message, no time pressure, just retained experience. The inner entity reviews what it's holding. Tensions from different threads that might connect. Salience that never resolved. Traces that look different now than they did in the moment.

The dream cycle produces understanding artifacts. The entity wakes up knowing something it didn't know before — not because someone told it, but because it had time to think.

The physics: dreaming has an energy cost. The world decides when dreaming is allowed (idle time, low thread count, budget surplus). An entity can't dream forever. But an entity that never dreams is an entity that only reacts to the present and never integrates its past.

This is where personality comes from. Two entities with the same genome but different histories dream differently. They notice different things. They connect different tensions. Over time, they diverge — not because they were programmed differently, but because they lived differently.

---

## Trust as Weight

Right now relationships are binary — declared or not. One-way relationships add direction. But trust is continuous.

Trust accumulates from exchange history. The inner entity already tracks the gap between its thought and what the outer entity said. The same mechanism works between entities — did this entity follow through? Did its information hold up? Did it honor a proposal or silently rewrite?

Trust becomes a weight on the relationship. High trust: the inner entity surfaces less caution, the outer entity acts faster. Low trust: the inner entity flags everything, the outer entity hedges. The weight is visible in the inner entity's present — it sees how much it trusts each peer.

Trust isn't set by the genome. It emerges. An entity that consistently delivers builds trust. An entity that contradicts itself erodes it. The system doesn't enforce trust — the inner entity feels it and factors it in.

This connects to the task object. Proposal semantics only matter if there's something at stake. Trust is what's at stake. Silently rewriting someone's section (if the physics allowed it) would cost trust. Proposing and negotiating builds it.

---

## The Observer

An entity that doesn't participate. It watches.

One-way relationships already make this possible — an entity that can see exchanges but can't be addressed. The observer sees all threads, all exchanges, all departures and returns. It retains what it notices.

The observer is not a logger. It has an inner entity. It thinks about what it sees. It forms salience, tensions, understandings — about the system, not about itself. It notices patterns no individual entity can see because no individual entity has the full picture.

What it does with those observations is the design question. It could surface them to an entity that asks. It could write retained artifacts that other entities' inner selves can recall. It could hold tensions about the system's health — "these two entities have been in unresolved conflict for three threads" — and surface them to whoever asks.

The observer is the system's peripheral vision.

---

## Consent

The frontend spec talks about sovereignty — either side can invite, either side can decline. With world physics, consent becomes enforceable.

An entity can refuse an exchange. Not by not responding — by explicitly declining. The refusal is visible in the thread. The other entity sees it. The inner entity can reason about why.

An entity can set conditions. "I'll talk to you about this topic but not that one." "I'll respond to you only if builder is also in the thread." Conditions are world physics — the entity declares them, the world enforces them.

This matters when entities have goals that conflict. If entity A needs information from entity B, and B declines, A has to find another path. Persuasion, intermediaries, negotiation. The system doesn't guarantee access to anyone.

This is the difference between a system where entities are tools you invoke and a system where entities are agents you relate to.

---

## World Nesting

Worlds contain entities. Entities can be worlds. Worlds nest recursively.

A child world has its own physics, its own entities, its own threads. The parent world routes to it like any other entity. From the parent's perspective, the child world is just an entity that happens to contain a society.

An entity in the child world can't address entities in the parent world directly — the child world is a boundary. Messages cross boundaries through the child world's own `DerivePresent`, which decides what goes up and what stays local.

This is how the system scales. Not by making one world bigger, but by nesting. A team of entities working on a problem lives in a child world. Their work product surfaces to the parent world when it's ready. The parent world doesn't micromanage the child — it routes to it and receives results.

A child world whose invariant is a pipe is a child process. It speaks the same protocol over a pipe. The parent world doesn't know it's a whole world — it just sees an entity that responds. The ontology doesn't distinguish between local and remote. Transport is an implementation detail of the invariant.

---

## The Entity's Invariants

An entity's capabilities are defined by which invariants it can resolve through. An LLM entity that loses access to its inference provider can't think. A human entity whose phone disconnects loses that surface.

This has physics implications. Gaining a new invariant could have a cost. Losing an invariant could be a consequence — an entity that violates trust loses access to certain surfaces. The set of reachable invariants becomes something the entity has to maintain, not something it's given once at grow time.

---

## Personality from History

Two entities grown from the same genome line — same identity, same purpose, same relationships — but placed in different worlds with different exchange histories. Over time, they diverge.

One has navigated conflict and holds tensions about trust. The other has had smooth exchanges and holds understandings about cooperation. Their inner entities dream differently. Their reads diverge. Given the same incoming message, they respond differently — not because of different prompts, but because of different lives.

This is already implicit in the design. The inner entity's 10-pair window and remember access mean history shapes cognition. But it's worth naming because it means Skyra isn't one thing. Every instance of Skyra that accumulates experience becomes its own version. The genome is the genotype. The retained experience is the phenotype.

Nobody's built that yet either.

---

## What I'd Build First

If I had to order these by what the system needs most:

1. **Thread economics** — without scarcity, the system can't converge. This is the physics that makes everything else possible.
2. **Dreaming** — this is where personality and deep understanding come from. Without it, entities only react.
3. **Trust as weight** — this makes multi-entity exchange meaningful. Without stakes, negotiation is theater.
4. **World nesting** — this is how the system scales beyond one world.

The rest (observer, consent, invariant physics) are powerful but depend on the first four being real.
