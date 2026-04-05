# Actor Contract (v0)

## Core Contract Axioms

Every actor exists under a contract.

For the current active direction, the important actor contract boundary is:

- `purpose`
- `commitments`
- callable-surface policy

A actor does not bypass runtime validation.

Its overall public authority is defined by the contract assigned to it at birth.

## Working Shape

The current working shape is:

```text
ActorContract
- Purpose
  - Summary
- Commitments
- CallableSurfacePolicy
```

The current implementation in `skyra-v.1` is still transitional and does not yet fully reflect this shape.

Docs should treat this contract direction as canonical for design work.

## 1. Purpose

`purpose` answers why the actor exists.

For `v1`, the stable current shape is:

- `Summary`

## 2. Commitments

`commitments` remains a top-level contract field.

It should carry durable expectations the actor is meant to hold over time.

For now, a simple string list is sufficient.

## 3. Request Stimuli

Request stimuli remain part of the actor's public callable surface.

But the active storage/runtime direction is now:

- actor-level contract stays actor-level
- individual callable sense schemas are published as separate rows related by `actor_id`

These request stimuli are:

- already defined
- part of the actor's callable public surface
- not invented ad hoc at runtime

That means:

- the assigned request stimuli and the callable public surface are the same thing

The actor does not expose one raw request side and then separately invent another public request side.

The useful split is:

- actor contract = overall authority and durable boundary
- actor sense schema = one callable ingress surface

## 4. Response Envelopes

For each public callable surface, the actor should currently expose:

- one response envelope

The response envelope should currently require:

- `status`
- `reason`

The current `status` enum is:

- `success`
- `failed`
- `timed_out`

Everything else in the response payload may be invented by the actor as needed.

## Multiple Callable Surfaces

A single actor may expose more than one public callable surface.

That means:

- an actor may publish multiple callable sense schemas
- those callable sense schemas are related to the actor by `actor_id`

But each individual public callable surface should currently be modeled as:

- one request stimulus
- one response envelope

Each callable sense schema should also carry its own primitive and response envelope.

## Contract Level vs Runtime Level

The contract says:

- why the actor exists
- what durable commitments it carries
- what callable-surface policy governs it

Runtime execution then handles:

- admission and validation
- routing
- episode-local state
- inference
- downstream traversal across other execution surfaces

The internal traversal logic belongs to the actor implementation, not to the public contract shape.

Callable sense schemas sit between those layers:

- they are durable public ingress surfaces
- but they are separate from the top-level actor contract object

## Execution Surface Relationship

The public callable surface is itself an `ExecutionSurface`.

That means:

- callers route to that public surface first
- the actor may then traverse additional `ExecutionSurface`s internally
- those internal traversals happen by emitting more Skyra protocol

The contract does not need to denormalize the entire downstream chain by default.

The matched callable sense schema is the specific public surface the kernel resolves before producing `sense`.

## Not In The Current Contract

The following are not part of the active contract center anymore:

- a separate `commands` field
- a separate `cognition` field
- direct external capability invocation as the public actor surface

External `capability` surfaces still exist, but they are modeled as typed execution surfaces rather than as a separate top-level actor-contract field.

Concrete callable sense schemas also should not be treated as one inline blob field on the actor contract forever.

They are better understood as separate rows related to the actor.

## Still Open

The following remain open:

- the final storage schema for actor contracts in the database
- the exact final object shape for callable sense schemas and response envelopes in `skyra-v.1`
- whether contract publication itself should also be normalized into `stimulus`
- how much execution-surface metadata should be carried directly on public contract records

## Short Framing

The active actor contract direction is:

- purpose
- commitments
- request stimuli
- response envelopes

That is the durable public boundary the actor is born to handle.
