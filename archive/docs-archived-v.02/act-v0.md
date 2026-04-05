# Act v0

## Purpose

This document defines the current working meaning of `act`.

`act` is the world-facing primitive.

It is not the retained read path and not the retained write path.

## Core Framing

The current actor-side primitive split is:

- `recall`
- `learn`
- `observe`
- `act`

Under that split:

- `recall` reads retained experience
- `learn` writes retained experience
- `observe` admits typed input into bounded actor runtime context
- `act` emits world-facing stimulus toward an execution surface

There is now also a separate ingress-envelope layer:

- `sense` is the normalized ingress envelope produced by the kernel at the receiving actor boundary

That means:

- `sense` is not the same thing as actor-side `observe`
- `sense` is not a caller-facing wire primitive

That makes `act` the primitive for intentional participation that crosses out of the current actor and touches some external or downstream boundary.

## What Counts As `act`

The current direction is that `act` includes forms such as:

- responding to a human
- probing a device
- searching the web
- calling an external service
- writing registration state
- emitting stimulus to a capability surface

The common property is:

- the actor is intentionally crossing an execution boundary

## What `act` Is Not

`act` should not absorb everything.

It should not replace:

- `recall`
- `learn`
- `observe`

Those stay distinct because they have different invariants and different read/write boundaries.

## Public Contract Shape

In the current stimulus-first model, `act` should be understood through published stimulus contracts rather than through a freeform command grammar.

That means a callable `act` surface should currently be modeled as:

- one request stimulus
- one response envelope

The response envelope currently requires:

- `status`
- `reason`

with `status` currently:

- `success`
- `failed`
- `timed_out`

Any other return fields may be actor-defined inside the response payload.

## ExecutionSurface

Every emitted stimulus carries an `ExecutionSurface`.

For `act`, that means:

- the first/public callable surface is itself an `ExecutionSurface`
- the handler may then traverse additional `ExecutionSurface`s internally by emitting more Skyra protocol

The initial typed execution surfaces remain:

- `actor`
- `capability`

## Jarvis vs Stark

The top-level primitive stays shared:

- `act`

What differs between major actor roles is the world-facing responsibility they are expected to own.

### Jarvis

For `v1`, `Jarvis` should primarily own human-facing `act`.

### Stark

For `v1`, `Stark` should primarily own system-facing `act`.

That currently includes work such as:

- probing
- registration writes
- structural mediation around capability surfaces

`birth_actor` should not be treated as ordinary world-facing `act`.

It belongs to Stark's structural surface rather than to the ordinary world-facing `act` boundary.

## Validation Shape

The likely validation ladder is:

1. the actor is allowed to emit `act`
2. the request stimulus matches a published contract
3. the target `ExecutionSurface` is valid for that published surface
4. the request payload is structurally valid
5. runtime validates the external boundary and permissions
6. a response envelope or returned stimulus is routed back

## Relationship To Registration

Registration should still be understood as valid world-facing work under `act`.

That same broad `act` boundary also covers:

- probing
- search
- external service use
- human response

The primitive stays broad.

The specific stimulus contract and execution surface carry the operational detail.

## Current Design Posture

The strongest current claims are:

- `act` is the world-facing primitive
- `act` should be modeled through published request/response stimulus contracts
- `ExecutionSurface` remains part of every emitted stimulus
- downstream traversal still happens by emitting more Skyra protocol

## Short Framing

`act` is the primitive for intentional world-facing participation.

In the current model, it should be understood as emitting published `stimulus` toward an `ExecutionSurface`, not as emitting an old freeform command string.
