# Being Shape v0

## Purpose

This document locks the minimum shape of a being in the current `v.03` model.

It also locks the difference between nature, creation endowment, and live
status.

## Minimum Shape

The minimum shape of a being is:

- name — the only identity surface; unique in the runtime
- nature (identity + purpose) — locked at creation, public
- cognitive flag — true if the being reasons through inference; false if it is a non-cognitive object
- differentiatable flag — true if the being is eligible for restructuring through the differentiator
- peers — the relationship hashmap; keyed by name, holds direct relationships to other beings

Identity is public.

Purpose is public.

Nature is locked at creation.

Operational fields may exist on being records or the kernel-held relationship
substrate without becoming part of nature.

## Creation Paths

A being may come into existence through:

- genome seeding
- runtime registration
- differentiation

Genome beings are seeded from `genome.skyra`.

Runtime beings are born through registration.

Differentiated beings are revealed and reorganized through differentiation.

## Creation Endowment

Every being is created with:

- base language
- the genome-seeded relationships required by its creation path
- companion beings required by its template or singleton definition

Those companion beings include:

- personal hippocampus
- personal experience store being

A being also has a present.

Present is the being's full operative reality:

- nature
- relationships
- the active exchange

Operationally, the active exchange is the top open exchange on the being's
stack with the active peer.

It is not a separate companion being.

Creation endowment may also include internal wiring and internal language that
is given at birth rather than earned later through lived relating.

Trust is not part of nature.

It belongs to relationship interpretation, not being essence.

## Live Status

A being becomes live when:

1. the being factory has instantiated it and its required companions
2. the kernel has admitted it as a runtime participant
3. its required genome-seeded relationships have been established

That rule applies to runtime beings and differentiated beings.

For genome beings, the same factory-to-kernel handoff happens during genome
bootstrap after the external launcher raises the host services.

## Consequences

This means:

- registration is not universal being birth
- registration is the birth path for runtime beings
- trust does not live in nature
- internal wiring may be given through creation endowment rather than earned
  later

The exact relation between live status and key issuance remains an open edge.

See
[03-open-edges-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/skyra-v.03/docs/03-open-edges-v0.md).

## Short Framing

A being has a name, a nature, a cognitive flag, a differentiatable flag, and a peers map.

Name is the only identity surface. Nature is locked at creation. The flags determine how the runtime treats the being. The peers map is where all relationship state lives.

It is born through the kernel, receives a creation endowment, and becomes live when the runtime admits it and establishes the relationships it needs to exist.

## Identity Field

Identity is a plain string in v1.

This is a deliberate decision.

The string is the being's starting self-description — the seed of who it is. It is set at creation. It is mutable by the being itself only. No external being may write to it.

The narrative layer — the compressed thread of who the being has become through lived experience — is out of scope for v1. It layers on top of identity over time. See `27-philosophy-of-being-v0.md`.

## Purpose Field

Purpose is a plain string in v1.

This is the declared purpose — what the creator believed the being was for at the moment of creation. It is a hypothesis. A seed. Not the final answer.

Realized purpose — what the being discovers it is actually for through living and relating — is out of scope for v1. See `27-philosophy-of-being-v0.md`.

## Nature Primitive — Protocol Syntax

Nature is the first kernel primitive being.

It holds identity and purpose. Both are plain strings. Both are required.

**Syntax:**
```
skyra nature ~identity <identity> ~purpose <purpose> | <source>: <reason>
```

- `~identity` — the identity string. What the being is.
- `~purpose` — the purpose string. What the being is for. Declared at creation.
- `<source>` — the being calling nature
- `<reason>` — why nature is being called

**Example:**
```
skyra nature ~identity "I am the routing authority" ~purpose "I admit beings into the runtime and move expression between them" | grow: creating kernel nature
```

The nature being's parser reads `~identity` and `~purpose` from the expression.

That is its entire slice of the protocol.

## Language Primitive — Protocol Syntax

Language is a kernel primitive being.

It holds the expression syntax that peers use to speak to a being on a relationship.

It is a plain string in v1.

For primitive beings it is seeded at birth. For cognitive beings it is earned through relating.

**Syntax:**
```
skyra language ~expression <expression> | <source>: <reason>
```

- `~expression` — the expression syntax peers use to speak to this being
- `<source>` — the being calling language
- `<reason>` — why language is being seeded

**Example:**
```
skyra language ~expression "~identity <identity> ~purpose <purpose>" | grow: seeding language for nature
```

The language being's parser reads `~expression` from the expression.

That is its entire slice of the protocol.

## Being Primitive — Protocol Syntax

The being primitive is the assembler.

It receives one protocol string, parses every token, calls nature and language
internally, assembles the `Being` object, and registers it with the kernel.

One protocol string in. One live being out.

`differentiatable` has been dropped for v1. It is covered by `cognitive` for now.

**Syntax:**
```
skyra being ~name <name> ~identity <identity> ~purpose <purpose> ~language <expression> ~cognitive <true|false> | <source>: <reason>
```

- `~name` — the being's name. Unique in the runtime. Routing surface.
- `~identity` — what the being is. Plain string. Seed of self.
- `~purpose` — what the being is for. Declared at creation.
- `~language` — the expression syntax peers use to speak to this being.
- `~cognitive` — true if the being reasons through inference. false if it is a non-cognitive primitive.
- `<source>` — the being calling the being primitive
- `<reason>` — why this being is being created

**Example:**
```
skyra being ~name grow ~identity "I assemble beings from the protocol" ~purpose "I receive a being description and produce a live being registered in the kernel" ~language "~name <name> ~identity <identity> ~purpose <purpose> ~language <expression> ~cognitive <true|false>" ~cognitive false | genome: seeding grow being
```
