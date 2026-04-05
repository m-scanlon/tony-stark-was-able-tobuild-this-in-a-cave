# Cognition Policy Ideas

## Status

This is an ideas document, not canon.

It should not be read as changing the current locked actor contract or stimulus-first runtime model.

## Core Framing

One possible later direction is that runtime may grow a cognition-policy layer near the contract boundary, and runtime may instantiate a cognition helper from that policy.

The important boundary is:

- the cognition helper is not the source of truth
- the episode remains the bounded runtime state container
- the actor and kernel remain the execution boundary

So cognition here should be thought of as a runtime helper layered on top of the existing model.

## Safer Shape

The safer shape is:

- cognition policy lives in or near the contract boundary
- runtime instantiates a cognition helper from that policy
- cognition helper emits observations, candidate stimulus, and stop/continue decisions
- the runtime still owns routing, validation, and state

So the cognition helper should not execute side effects directly.

## Why Not Make Cognition The Source Of Truth

If cognition became the source of truth, the model would start collapsing:

- runtime policy
- episode state
- stimulus execution
- loop style

into one object.

That is exactly what the current actor/episode/frame split is trying to avoid.

## Short Framing

The contract may eventually define a cognition policy.

Runtime may instantiate a cognition helper from that policy.

That helper can collect observations and candidate stimulus while the episode remains the source of truth.
