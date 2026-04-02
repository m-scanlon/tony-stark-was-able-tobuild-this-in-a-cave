# Node Contract (v0)

## Core Contract Axioms

Every node exists under a contract.

For the current `v1` direction, the active node contract boundary is:

- `purpose`
- `capabilities`
- `stimulus`

A node does not directly bypass runtime validation.

It emits callable capability-bound actions that the system validates and executes.

## Locked Shape

The current locked shape is:

```text
NodeContract
- Purpose
  - Summary
- Capabilities
  - CapabilityIDs
- Stimulus
  - AcceptedTypes
  - EmittedTypes
```

This matches the current implementation surface in
`skyra-v.1/node/contracts/contracts.go`.

## 1. Purpose

`purpose` answers why the node exists.

For `v1`, the locked purpose shape is:

- `Summary`

`Limits` is not a first-class field in the current `v1` node contract.

## 2. Capabilities

`capabilities` is the node's callable action surface.

This absorbs the older `commands` framing.

In other words:

- commands are capabilities
- capability use is how the node expresses what it may emit

The current locked field is:

- `CapabilityIDs`

Those ids point at capability contracts / callable surfaces.

The current top-level primitive family remains:

- `recall`
- `learn`
- `observe`
- `act`

That means the node contract no longer needs a separate `commands` field to name them.

## 3. Stimulus

`stimulus` defines the node's typed input and typed output boundary.

The current locked fields are:

- `AcceptedTypes`
- `EmittedTypes`

If incoming stimulus does not match the contract's accepted types, the node is not eligible to act.

## Not In The Current v1 Contract

The following are intentionally not part of the current `v1` node contract shape:

- `NodeType`
- `Purpose.Limits`
- `Cognition`
- `Commands`
- `LearningEnabled`

`cognition` may still matter as a later runtime or policy layer, but it is not part of the current locked node contract shape.

## Contract Level vs Runtime Level

The contract says:

- why the node exists
- what callable capabilities it may use / emit
- what typed stimuli may wake it up
- what typed stimuli it may emit

Runtime execution then handles:

- admission and validation
- dispatch
- result writeback
- episode-local state

## Same Contract Model Across Nodes

This contract model applies across node roles.

That includes:

- user-facing or task-facing nodes
- `Jarvis`
- `Stark`
- bounded worker nodes such as `probe` and `registration`

What differs between nodes is the contract content, not the existence of a separate node ontology.

## Still Open

The following remain open:

- the exact storage schema for node contracts in the database
- the exact shape of the capability contract ids used by node contracts
- actor-to-actor delegation edges
- actor-to-capability invocation edges
- exact `observe` schema
- exact `act` content, modality, and timestamp encoding

See also:

- [capability-contract-prelim.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/capability-contract-prelim.md)
- [protocol-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/protocol-v0.md)
- [stimulus-types-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/stimulus-types-v0.md)

## Short Framing

The current `v1` node contract is:

- purpose
- capabilities
- stimulus

That is the durable node boundary.
