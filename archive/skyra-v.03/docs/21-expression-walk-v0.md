# Expression Walk v0

## Purpose

This document locks how expression moves through the relationship field in the
current `v.03` canon.

It also locks the distinction between expression, inference, signal,
`trace_token`, and relationship strength.

## Directionality

The relationship pair is unordered. Two beings. One relationship.

Direction is turn-level flow, not relationship identity. If A expresses to B, that uses the same relationship as B expressing to A. The relationship does not split by direction.

The relationship-graph key is the unordered pair `(being_a, being_b)`. Direction is expressed at the turn level through the expression itself — routing context, not relationship identity.

## The Medium

The medium is the relationship field.

It is not a router.

It is not a registry.

It is not a broadcast surface.

No component computes a recipient set.

No component asks who should get this.

What happens instead is local:

- a being fires from its present
- expression moves only across existing relationships
- what happens next is decided by the being that receives it

One narrow kernel exception exists:

- emotional flags on expression may cause the kernel to route a copy according
  to fixed threshold rules while the original expression continues unchanged

## The Present

The present is the complete context a being operates from at any moment.

It always carries:

- nature
- relationships
- the active exchange

These are not optional layers that turn on and off.

They are the being's living present.

## The Expression Walk

A being does not carry an expression forward mechanically.

It fires an expression from its present across existing relationships.

What arrives at another being enters inference.

The interpretation is the inference call itself.

That inference call may fire a new expression from the receiving being's
present.

That new expression may preserve what came before, bend it, replace it, or end
it.

Each hop is therefore a fresh expression.

What matters for relationship change is not message fidelity.

What matters is that adjacent beings fired in relation to one another.

Each newly fired expression is its own turn.

The existing signed-envelope rule still applies to that turn.

The signed envelope belongs to channel-boundary verification, not to the
inside of a being.

The kernel sheds the envelope before expression enters a being's present.

A being never receives the whole envelope object as part of its operative
reality.

## Signal

The transport signal is fixed:

```text
Signal {
  id           // system generated
  origin       // system assigned; kernel only, never exposed to beings
  trace_token  // system generated; kernel only, Hebbian learning
  raw          // the only outside-provided and user-facing value; the full protocol string
}
```

`raw` is the full protocol string:

```text
skyra <being> <expression> | <source>: <reason> ~<emotional_signal> <value>
```

## Protocol Zones

The `|` character is the hard divider between two zones.

**Left of `|`:**

```text
skyra <being> <expression>
```

The kernel parses only the protocol prefix (`skyra`) and the being name.

Everything after the being name and before `|` is expression.

Expression is passed through untouched to the target being.

The kernel never scans the left zone for tokens beyond the being name.

A `~` appearing in expression is not interpreted. It is part of expression.

**Right of `|`:**

```text
<source>: <reason> ~<emotional_signal> <value>
```

The right zone is fully parsed by the kernel.

`<source>` is the name of the peer the origin was in exchange with when it
fired. It is delimited by `:`.

`<reason>` is the origin being's internal record of why it fired. It runs
until the first `~` or end of string.

`~<emotional_signal> <value>` pairs follow. Each `~` token begins a new
emotional signal. The value is the number immediately after it.

## Mandatory Fields

The `|` divider is mandatory on every impulse.

`<source>` is mandatory. A being must name which relationship it fired from.

`<reason>` is mandatory. A being must have a reason for every expression it
fires. There are no silent mechanical hops.

## Kernel Parse Targets

The kernel extracts from the right zone:

- source peer name — for the three-write algorithm
- reason — stored in the origin being's exchange record
- emotional signals and values — for structural routing decisions

## Stripping

The entire right zone is stripped before the expression enters the target
being's present.

The target being receives only expression.

The target being interprets from its own present without the origin's framing,
reasoning, or emotional signal values.

There is no separately registered resolution method in the kernel.

Response behavior is baked into the being at birth by the being creator class.

The kernel just dispatches the signal.

The signal is not the meaning lived by the being.

It is the minimal kernel carrier that lets the walk happen and lets trace
accumulate.

`stress` and `anger` belong on expression, not on signal.

`strain` remains internal to the being.

`id`, `origin`, and `trace_token` are system-handled.

The being does not receive them.

The being receives only expression after kernel parsing and envelope shedding.

## Trace And Relationship Strength

The `trace_token` belongs to the kernel, not to the being.

Its purpose is not interpretation.

Its purpose is relationship strength.

Outward `stress` on expression sets the `trace_token` TTL.

Higher `stress` gives the trace longer reach before decay.

Relationship emergence is a kernel operation.

Every time a signal passes through the kernel, the kernel mechanically updates
edge weight on the relationship graph for the unordered pair.

`trace_token` is kernel-internal. It is never visible to beings. It lets the kernel register adjacent co-firing across the walk without requiring the full path to be stored.

This is Hebbian wiring in the current canon.

Beings that fire together wire together.

The trace does not need to preserve the full path.

It needs only enough local trace for the kernel to register adjacent
co-firing.

When the TTL reaches zero, the trace dies.

Expression may continue, but that walk stops contributing to edge weight.

No inference is involved in this relationship-graph update.

Before threshold, the pair has only pre-relationship edge weight on the
kernel-maintained graph rather than a live direct relationship in either
being's relationship hashmap.

When edge weight between two beings crosses the relationship-emergence
threshold, the kernel adds the direct relationship to both beings'
relationship hashmaps.

When edge weight later decays below threshold, the kernel removes that direct
relationship from both hashmaps.

## Replies And Exchange

Replies are not return paths.

A reply is another expression.

It does not mechanically retrace an earlier walk.

A later expression emerges from the being's present, specifically from the top
open exchange on that being's stack with the active peer.

The continuity is cognitive, not transport-level.

## Emergence

Global structure emerges from repeated local events:

- expression
- inference
- adjacent co-firing
- relationship strengthening
- decay
- pruning
- differentiation

No central planner decides the structure in advance.

Beings fire from their present across existing relationships.

The field changes because repeated local relating changes relationship
strength.

## Short Framing

Skyra does not model propagation as message forwarding to recipients.

A being fires an expression from its present.

What arrives enters inference.

Inference may fire a new expression.

The kernel tracks `trace_token` so adjacent co-firing can strengthen
relationship.

No component computes a recipient set.

The global structure emerges from local expression and relationship change.
