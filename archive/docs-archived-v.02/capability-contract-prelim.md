# Capability Surface (Prelim)

## Core Framing

`capability` still exists, but it should now be understood through the stimulus-first execution-surface model.

The important point is:

- `capability` is a typed `ExecutionSurface`

It is no longer best understood as a separate public callable contract family that other actors directly use by default.

## ExecutionSurface Model

The current parent routing concept is:

- `ExecutionSurface`

The initial child surface kinds are:

- `actor`
- `capability`

Both surface kinds are registered.

Actor surfaces publish request/response contracts.

Capability surfaces publish `act` contracts and may also originate ingress signal.

## Capability Boundary Invariant

For capability surfaces, the useful current invariant is:

- ingress to Skyra is first preserved as capability-originated ingress signal
- that ingress is then normalized into `sense` at the receiving actor boundary
- outbound from Skyra is always `act`

That means capability surfaces should not be modeled as publishing:

- `observe`
- `recall`
- `learn`

Instead, a capability surface should expose:

- ingress
- `act`
- or both

## Actor Surface vs Capability Surface

### Actor Surface

An actor surface is the first/public callable surface other actors generally target.

Its contract carries:

- purpose
- commitments
- request stimuli
- response envelopes

### Capability Surface

A capability surface is a registered execution surface that usually sits behind an actor such as a `stewart`.

Its contract is still useful because the runtime needs to know:

- whether that capability exposes ingress, `act`, or both
- what ingress shape may originate from it
- what request shape applies for `act`
- what response envelope it returns for `act`
- how routing reaches that surface

The crucial difference is:

- actor surface = public callable abstraction
- capability surface = lower-level execution surface

## Why This Matters

Without this distinction, the model gets muddy:

- devices start looking like general public actor APIs
- public contracts get polluted with low-level external complexity
- execution routing and public abstraction collapse together

The cleaner split is:

- other actors generally interact with an actor surface first
- that actor may then traverse a downstream capability surface under the hood

## Relationship To `stewart`

A `stewart` actor is the usual mediation boundary for an external capability.

At a high level:

1. an external capability is discovered and verified
2. the capability surface is registered
3. a `stewart` actor is born under contract
4. that `stewart` publishes a simplified public stimulus contract for the rest of the system
5. that `stewart` later traverses the downstream capability surface under `act` and receives external ingress from that surface normalized into `sense`

This keeps:

- discovery
- registration
- public abstraction
- world-facing execution

as separate concerns.

## Current Contract Shape

The currently useful capability-surface contract shape is:

- `CapabilityID`
- `Name`
- `ExecutionSurface`
- `BoundaryMode`
- `Schema`

That code shape is still transitional.

The active design meaning is now:

- identifier
- human-readable name
- routing/execution surface
- boundary mode, currently ingress or `act`
- published surface shape

## Current Design Posture

The strongest current claims are:

- `capability` remains a typed execution surface
- capability surfaces are registered
- capability surfaces publish `act` contracts and preserve ingress shape where relevant
- capability ingress should be normalized into `sense` at the receiving actor boundary
- capability outbound from Skyra should always be modeled as `act`
- public actor-facing interaction should usually happen through an actor surface first
- `stewart` actors are the main abstraction layer over external capability complexity

## Short Framing

Capability surfaces still matter.

They are just no longer the main public contract center.

They are typed execution surfaces that actors such as `stewart` can traverse beneath a cleaner public stimulus contract, with capability ingress normalized into `sense` at the actor boundary and `act` reserved for capability-directed output.
