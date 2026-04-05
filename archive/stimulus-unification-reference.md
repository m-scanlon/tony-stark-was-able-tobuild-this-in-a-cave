# Stimulus Unification Reference

## Purpose

This is a working top-level reference for the current architecture direction around `stimulus`, schema ownership, routing, and external-surface mediation.

It captures the clarifications made in the recent architecture discussion so they can be referenced directly without reconstructing the thread.

## Core Direction

The runtime is moving toward one unified message family:

- actors emit typed `stimulus`
- actors receive typed `stimulus`
- the kernel routes typed `stimulus`
- APIs, devices, actors, and other runtime/external sources are all normalized into `stimulus` at the actor boundary

The actor-facing runtime should therefore be understood as stimulus-first rather than command-first.

## What This Replaces

This direction supersedes the older framing where:

- actors emitted commands
- the kernel executed commands
- results came back as `command_result`
- capabilities were modeled as a separate callable registry family

That older command/capability split is no longer the main conceptual center.

## Registry Model

The registry should be a `stimulus registry`.

The important implication is:

- the main thing being published into the runtime is the `stimuli` contract/object
- each contract record may have a DB primary key for storage identity
- each contract record should also have a `name`
- `Stark` may assign that `name` at contract creation

This means the registry is no longer best understood as a separate capability registry of directly callable external surfaces.

All registered `stimuli` objects carry their `ExecutionSurface` in the contract itself.

The registry therefore stores contracts that already know where they execute.

## Capability Submission And `stewart` Actor Birth

When a user submits an external capability for registration:

1. a `stewart` actor is born under a contract
2. that actor learns the schema of the external capability
3. that actor emits or publishes a `stimuli` contract for the runtime
4. world-facing interaction with the external surface happens under the `act` primitive

The correction here is important:

- this is `act`, not `interact`

The primitive family remains:

- `act`
- `recall`
- `learn`
- `observe`

At actor birth, the contract assigns request stimulus to the actor.

Those request stimuli are not invented ad hoc by the actor at runtime.

They are part of the assigned callable surface the actor is born to handle.

Those request stimuli are already defined.

The actor does not invent ad hoc public request protocols at runtime.

## External Surface Mediation

The external surface is not callable directly by other actors.

Instead:

- other actors interact with the `stewart` actor
- the `stewart` actor consumes incoming `stimulus`
- the `stewart` actor makes the external call
- the `stewart` actor emits resulting `stimulus` back into the system

That means the external API/device surface sits behind the `stewart` actor rather than existing as a directly callable first-class runtime boundary for general actor use.

## Surface Model

`ExecutionSurface` is the parent type.

The initial child surface types are:

- `actor`
- `capability`

Both surface types are registered.

Both surface types publish:

- request stimulus schema
- response envelope schema

Every registered surface is therefore a contract record with a public request/response shape.

For the current simplified API posture, that public shape should be understood as:

- one request schema
- one response envelope schema

The response envelope should currently have two mandatory fields:

- `status`
- `reason`

The current `status` enum should be:

- `success`
- `failed`
- `timed_out`

Everything else in the response payload may be invented by the actor as needed.

## Actor Composition Model

A actor may be in charge of one thing or many things.

That means:

- a actor may be assigned multiple public request stimuli
- a actor may register multiple public response envelopes
- a actor may simplify and abstract over lower-level complexity
- a actor may traverse multiple downstream execution surfaces while fulfilling one public request

So the actor is not just a raw endpoint.

It can also act as:

- an adapter
- a composer
- a wrapper over lower-level surfaces

The key clarification is:

- the assigned request stimuli and the public callable surface are the same thing

That means the actor is not exposing one raw request side and then inventing a second separate public request side.

Instead:

- the actor is assigned request stimuli in its contract
- those assigned request stimuli are the actor's callable surface
- the actor absorbs the lower-level complexity behind that callable surface
- the actor registers a response envelope as the simplified return contract for that callable surface

## Contract Records

The DB stores contract records for registered stimuli/surfaces.

The important current shape is:

- each record gets a DB primary key
- each record should have a `name`
- the request/response contract is the public interface
- downstream traversal does not need to be explicitly encoded as a required relationship on the public record

Separate records may still exist for downstream surfaces.

The public record does not need to denormalize that internal chain by default.

## Actor Contract Shape

The actor contract still needs `commitments`.

The current important contract shape should be understood as:

- `purpose`
- `commitments`
- request stimuli
- response envelopes

This is not meant to say the contract is only three fields.

It is meant to identify the major runtime-relevant shape:

- what the actor is for
- what it is committed to
- what request stimuli it can receive
- what response envelopes it can emit/register

The actor contract is effectively the actor's registry entry.

A actor contract may expose multiple callable surfaces.

Each public callable surface should currently be modeled as:

- one request stimulus
- one response envelope

## ExecutionSurface

`ExecutionSurface` stays.

The clarified meaning is:

- `ExecutionSurface` is about routing
- every piece of `stimulus` has an `ExecutionSurface`

This is not just a capability-only field.

It is a universal property of runtime `stimulus`.

## Multi-Surface Execution

Execution can occur in stages.

The clarified model is:

- when caller actor `A` calls callee actor `B`, `B` is the first execution surface
- that first/public callable surface is also an `ExecutionSurface`
- the handler may then traverse additional `ExecutionSurface`s internally
- those downstream traversals happen by emitting more Skyra protocol
- those downstream surfaces may be `actor` or `capability`

So a caller does not need to know the full downstream chain.

It targets the immediate callable surface and its first execution surface.

The callee owns the downstream translation.

## How `capability` Still Fits

`capability` does not disappear entirely.

What changes is its role.

It is no longer best understood as:

- a separate public callable registry that other actors directly use

It is better understood as:

- one typed `ExecutionSurface`
- typically downstream of an owning actor such as a `stewart`

This preserves `ExecutionSurface` while still keeping general actor-to-actor interaction stimulus-first.

## Actor Boundary Knowledge

The callee-side knowledge model remains:

- each actor knows what it received in its stimulus boundary
- the callee does not know the full lineage of the request
- the callee does know who it owes dependency resolution to

So request lineage and return responsibility remain separate concepts.

## Request And Response Clarification

The request/response abstraction is now clearer:

- request stimuli are assigned to the actor in its contract
- those request stimuli are already the callable public abstraction
- other actors send those request stimuli directly
- the actor handles the lower-level complexity behind them
- the actor registers one public response envelope as the return-side contract

So the public actor-facing pair is:

- request stimulus
- response stimulus envelope

Those are the abstraction layer other actors interact with.

A actor may expose multiple such public pairs, but each individual callable surface should still be modeled as:

- one request stimulus
- one response envelope

This should be thought of more like an API:

- request schema = what the caller sends
- response envelope schema = what the caller gets back
- internal operations are hidden in the handler

The runtime should not treat response stimuli as if they were just exposed handler internals.

The current recommended simplification is:

- one public request
- one public response envelope
- internal success/failure/result branching stays inside the response payload rather than becoming multiple public response contracts

The current response-envelope rule is:

- mandatory field: `status`
- mandatory field: `reason`
- actor-defined fields may carry the rest of the return payload

## Dependencies And Obligations

The current split remains:

- `dependencyLedger` is one record per downstream dispatch
- obligations are identified when the actor observes stimulus
- obligations have their own registry
- active obligations remain projected in the frame while the obligation exists

This means:

- downstream waiting is not the same thing as upstream responsibility
- `dependencyLedger` is not the obligation model

Current runtime naming for dependencies is:

- collection: `dependencyLedger`
- entry: `dependency`
- states: `open`, `resolved`, `failed`, `timed_out`

The exact `dependency` shape remains open.

The exact `obligation` shape remains open.

## Transitional Mismatch In Current Contracts

Some current contract names are transitional and still carry older command-centric assumptions.

Most notably:

- the top-level dispatch contract is now named `StimulusEnvelope`
- but its current fields still include:
  - `calling_actor`
  - `command`

That means the name has moved toward the new model faster than the field shape.

The current contract surface should therefore be treated as partially transitional rather than fully reconciled.

## What Now Looks Stale In The Docs

If this newer model becomes the active canon, the following doc families are now the main stale areas.

### Command-Centric Runtime Docs

- `docs/protocol-v0.md`
- `docs/command-namespace-prelim.md`
- `docs/runtime-primitives-and-artifacts-prelim.md`
- `docs/data-model-prelim.md`
- `docs/actor-substrate-interface-v0.md`
- `docs/actor-process-v0.md`
- `docs/architecture-overview-v0.md`

These still assume some version of:

- actor emits command
- kernel executes command
- result comes back as `command_result`

### Capability-Contract And Capability-Registry Docs

- `docs/capability-contract-prelim.md`
- `docs/database-schema-v0.md`
- `docs/capability-probing-v0.md`
- `docs/device-registration-v0.md`
- parts of `docs/architecture-overview-v0.md`

These still assume:

- a distinct capability-contract family
- a distinct capability-registry/storage center
- published capabilities as standalone callable surfaces
- `ExecutionSurface` as a capability-surface concern rather than a universal stimulus-routing concern

### Docs Closest To The New Model

These are already closer to the newer direction:

- `docs/stimulus-types-v0.md`
- `docs/actor-contract-v0.md`
- `docs/frame-v0.md`
- `docs/episode.md`

Even there, some language still reflects the older command/stimulus split.

## Open Questions Still Left Open

The following are still intentionally unresolved:

- what exact object is emitted and queued when one actor sends work to another actor
- whether `command_result` survives as a separate event family or is normalized into typed `stimulus`
- whether `contract_publication` is also `stimulus` or remains a control-plane event
- the exact schema shape of the `stimulus registry`
- the exact typed shape of `ExecutionSurface`
- the exact `dependency` object shape
- the exact `obligation` object shape
- whether obligation metadata should eventually live as a sub-object under the base `StimulusEnvelope`

## Short Framing

The clarified model is:

- the runtime should think in terms of `stimulus`
- the registry should be a `stimulus registry`
- the registry stores `stimuli` contract records
- actors emit and receive typed `stimulus`
- `ExecutionSurface` is a universal routing property on every registered surface
- `ExecutionSurface` stays and becomes a universal stimulus-routing concept
- `actor` and `capability` are the initial typed execution surfaces
- both `actor` and `capability` surfaces publish request/response schemas
- each public callable surface should currently be modeled as one request schema plus one response envelope schema
- the public callable surface is itself an `ExecutionSurface`
- handlers may traverse additional execution surfaces by emitting more Skyra protocol
- external surfaces are mediated through actors such as `stewart` actors rather than being directly callable by other actors
