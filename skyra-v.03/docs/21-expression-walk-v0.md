# Expression Walk v0

## Purpose

This document locks how expression moves through the relationship field in the
current `v.03` canon.

It also locks the distinction between expression, inference, signal,
`trace_token`, and relationship strength.

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

## Signal

The infrastructure signal is lean:

```text
Signal {
  id
  trace_token
}
```

The signal is not the meaning of the expression.

It is the minimal kernel carrier that lets the walk happen and lets trace
accumulate.

`stress` and `anger` belong on expression, not on signal.

`strain` remains internal to the being.

## Trace And Relationship Strength

The `trace_token` belongs to the kernel, not to the being.

Its purpose is not interpretation.

Its purpose is relationship strength.

Outward `stress` on expression sets the `trace_token` TTL.

Higher `stress` gives the trace longer reach before decay.

Relationship emergence is a kernel operation.

Every time a signal passes through the kernel, the kernel mechanically updates
edge weight on the relationship graph for the unordered pair.

`trace_token` is the kernel carrier that lets this update happen across the
walk without requiring the full path to be stored.

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
being's local relationship hashmap.

When edge weight between two beings crosses the relationship-emergence
threshold, the kernel adds the direct relationship to both beings' local
relationship hashmaps.

When edge weight later decays below threshold, the kernel removes that direct
relationship from both hashmaps.

## Replies And Exchange

Replies are not return paths.

A reply is another expression.

It does not mechanically retrace an earlier walk.

A later expression emerges from the being's present, specifically from the
active exchange within it.

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
