# Frame v0

The frame is the full page consumed by inference.

It is the bounded in-scope runtime context for the current inference step.

## Layout

The current working frame layout is:

1. `purpose`
2. `interaction`
3. `recall`
4. `time`

## Purpose

`purpose` comes from the actor contract.

It provides the governing lens for how the actor should interpret the current situation.

## Interaction

`interaction` captures the current external exchange.

It should contain:

- `current_stimulus`
- `recent_interaction_history`

This keeps the active thing the model may need to respond to clearly distinguished from recent context.

It may also include a pending-sense view when the actor is deciding what to observe next.

## Recall

`recall` contains retained artifacts currently brought into scope.

The current retained artifact families are:

- `trace`
- `understanding`
- `salience`
- `tension`

Recall is supportive and provisional.

It informs interpretation but does not by itself finalize meaning.

## Time

`time` makes current temporal context explicit to the actor.

It should currently expose:

- `now`

This allows the actor to reason about:

- how old a pending `sense` is
- whether an external signal is stale
- whether it is still worth observing or acting on something

## Contract Boundary

The frame does not carry a separate command-allowance field in the current contract.

The actor's public request/response surface remains part of the active actor contract rather than the frame.

## Current Framing

The frame should be thought of as a structured inference page, not a visual screen layout.

Order and section labels matter more than any left/right or top/bottom visual metaphor.
