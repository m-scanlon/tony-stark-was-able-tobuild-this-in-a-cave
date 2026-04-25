# Relational Nexus

## What It Is

The Nexus is the data structure that lives on a relationship from one being's perspective. It is what the relationship is built on. It determines the depth and nature of how a being can be related to.

Every peer channel on a being is a Nexus. The Nexus type is the being type from the relational side.

## Three Types

**CognitiveNexus** — holds the full exchange history. Thread hashmap. The being reasons over its entire relational past with this peer. Used for cognitive beings.

**FunctionalNexus** — holds only the last expression. Minimal. The being only needs the signal that was just passed to it. Used for non-cognitive transducers.

**WorldNexus** — holds an inner world. A full nested kernel with its own being registry. From outside it looks like any other peer. Inside it is a populated living space. Used for world beings.

## Two Sides Of A Being

A being has two faces. The kernel connects them.

**Left — Nexus type.** What kind of Nexus peers get when they relate to this being. The relational face. The ontological side.

**Right — Praxis / execution surface.** How the being acts when the kernel dispatches to it. Inference process, internal handler, world dispatch. The operational side.

These are separate declared properties. Neither implies the other. A world being (WorldNexus on the left) could have inference praxis on the right, or direct world dispatch, or something else.

## What This Replaces

`cognitive bool` was doing double duty — it was the ontological flag AND the channel type selector. With three Nexus types it cannot keep doing that.

The kernel branches on Nexus type to seed the right channel at `seedPeer` time. The kernel dispatches via praxis when routing to the being. Two separate kernel decisions.

`cognitive bool` goes away. Replaced by two declared fields on the being:
- Nexus type (what relational channel peers get)
- Praxis (how the kernel executes when routing to this being)

## Nexus Struct

One struct. Shared fields on the base. Type-specific payload in `type`.

```
Nexus {
    name             string  // peer name, routing surface
    callableLanguage string  // what the peer offers, how to speak to it
    praxis           Praxis  // execution surface — how to invoke this peer
    type             Nexus   // typed as Nexus — the specific kind
}
```

The `type` field holds the kind-specific Nexus. The three kinds are themselves Nexus implementations:

- **CognitiveNexus** — exchange thread hashmap, last active thread
- **FunctionalNexus** — last expression
- **WorldNexus** — inner `*World`

The kernel inspects `type` to branch. New kinds require only a new Nexus implementation — the base struct does not change.

The structure is recursive by nature. A Nexus contains a Nexus. A WorldNexus contains a world of beings, each of which has Nexuses, which may themselves be WorldNexuses. The type holds at every depth.

Note: `type` is a reserved word in Go. Field name in code will be `kind`. `type` is the conceptual name.

## Code Shape

The `RelationshipChannel` interface becomes `Nexus`.

```go
type Nexus interface {
    Name() string
    CallableLanguage() string
    Praxis() Praxis
    Kind() Nexus
    Send(delivery DeliveredImpulse) ChannelResult
    // DerivePresent(receiver *Being, sender *Being) string
    // leaning toward yes — the Nexus knows its own data, it should know how to render itself.
    // not committed yet.
}
```

The three concrete types embed `BaseNexus` and implement `Nexus`:

```go
type BaseNexus struct {
    name             string
    callableLanguage string
    praxis           Praxis
    kind             Nexus
}

type CognitiveNexus struct {
    BaseNexus
    exchanges        map[string]ExchangeThread
    lastActiveThread string
}

type FunctionalNexus struct {
    BaseNexus
    lastExpression string
}

type WorldNexus struct {
    BaseNexus
    world *World
}
```

`seedPeer` in `grow.go` branches on being type to construct the right Nexus. Not on a cognitive flag.

## Open Questions

- Exact genome syntax for declaring Nexus type and praxis on a being
- Whether Nexus type and praxis are inferred from other declarations or always explicit
- What `DerivePresent` renders for a WorldNexus peer
- Whether a being can change its Nexus type after registration or if it is fixed at birth
