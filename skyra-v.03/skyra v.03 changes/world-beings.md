# World Beings

Two being types. Both append to the being list. Both callable via the same exchange protocol.

**Single being** — the current model. One being, one execution surface, one callable language. The genome syntax is unchanged.

**World being** — a being with a full world inside it. Its own kernel, its own being list, its own exchange state. From the outside it looks like any other being. Protocol sent to it enters its inner kernel and gets processed there.

## Genome Syntax

```
skyra world ~name <name> ~identity <identity> ~purpose <purpose> | reason
```

One new directive. Everything else stays the same. `grow` handles it — creates a world being and appends it to the outer being list.

## Calling A World Being

No new syntax. The caller opens an exchange with the world being the same way it opens an exchange with any being:

```
skyra start-exchange ~with research ~about <thread> ~because <reason> ~say <expression>
```

The signal enters the world being's kernel. The inner kernel processes it as protocol. Responses come back through the same exchange channel. The outer kernel does not know or care that there is a world on the other side.

## What Lives Inside

A world being starts empty — no beings, no exchanges. Beings are grown inside it through the same genome protocol. The world being's inner kernel accepts `grow` directives the same way the outer kernel does.

Who grows beings inside a world being is up to the being that created it. Skyra can grow a world, pass it a genome, and let its inner kernel bootstrap from there.

## Self-Extension

Skyra can write a genome directive containing `world`, emit it through the adapter-writer surface, and create a named world she owns. She grows beings inside it. She calls into it. She has a callable namespace under her control with no human in the loop.

This is differentiation at the world level. Not just new beings — new worlds with their own populations.

## Relationship To Relationship Worlds

A world being is the concrete implementation of the relationship-worlds idea. Instead of modeling a relationship as a world metaphorically, Skyra can grow a world being for a relationship. All beings that emerge from that relationship live inside that world. The relationship has a home.

## Composition Patterns

World beings compose through ordinary beings. No new syntax. No new machinery.

### Gateway Pattern

A single being sits in front of multiple world beings and presents a unified interface:

```
caller
  → gateway being
      → world A
      → world B
```

The gateway has relationships to both world beings. It opens exchanges with each of them like any other being. The caller talks to one interface. The gateway decides what goes where, aggregates responses if needed. The two worlds never touch each other — the gateway is the only thing that crosses that boundary, and only because it explicitly opened exchanges with both.

The isolation is real. World A's kernel has no visibility into world B's kernel.

### Patterns That Fall Out

**Work / personal isolation** — two world beings, one gateway. Neither world sees the other's exchanges. The user has one interface, two completely separate memory spaces.

**Experimentation** — run a new being population in an isolated world before promoting anything to the main world. If it works, grow the beings into the outer world. If it doesn't, close the world being and nothing was contaminated.

**Multi-user** — each user gets their own world being. One gateway routes by identity. Users share a single Skyra interface but their worlds never intersect.

These are not features. They are compositions of world beings and ordinary beings using the exchange protocol that already exists.

## Open Questions

- Does a world being have its own execution surface, or does it inherit the outer kernel's dispatch?
- Can a world being call back into the outer world — does it have a reference to the outer kernel?
- What happens when a world being is closed — are the beings inside it destroyed or suspended?
- Does the outer kernel's present show the contents of world beings or just their existence?
- Can world beings nest — a world inside a world?
