# Modeling Interaction (Superseded)

This document is historical and is not active canon for the current runtime.

It captured an older simplification where `act` was modeled as one universal
interaction object with fields such as:

- `target`
- `content`
- `modality`
- `timestamp`
- `reason`

That is no longer the active public contract shape.

## Current Source Of Truth

For the current model, use:

- [act-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/act-v0.md)
- [protocol-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/protocol-v0.md)

## Current Direction

The active direction is now:

- `act` remains the world-facing primitive
- runtime is stimulus-first
- `act` emits stimulus toward an `ExecutionSurface`
- the public callable surface is modeled as one request stimulus plus one
  response envelope
- the response envelope requires `status` and `reason`

## Historical Value That Still Survives

The older framing still preserves one useful intuition:

- `act` is intentional boundary engagement that produces an effect or signal

Fields such as `target`, `content`, `modality`, and `timestamp` may still appear
inside specific stimulus payloads where they are useful.

They are just not the universal top-level `act` contract anymore.
