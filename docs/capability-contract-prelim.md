# Capability Contract (Prelim)

## Core Framing

The cognitive system's contract with a device is a capability contract, not a node contract.

This is an important distinction.

Devices and device surfaces should not be forced into the same contract shape as nodes.

## Node Contract vs Capability Contract

### Node Contract

A node contract governs a node as a runtime operator.

It defines:

- why the node exists
- what stimuli may invoke it
- what outward interaction forms it may emit
- later, what command surface it may be allowed to use

A node contract is about bounded runtime participation.

### Capability Contract

A capability contract governs an external capability surface.

It defines:

- what commands that capability exposes
- what inputs those commands accept
- what outputs or results they may return
- what verification or probe established that the capability is real
- what limits or constraints apply to using it

A capability contract is about callable external ability, not node identity.

## Why This Matters

Without this distinction, the model gets muddy:

- devices start looking like nodes
- node birth gets mixed with capability discovery
- runtime behavior and external affordances collapse together

The cleaner split is:

- nodes are runtime operators under node contracts
- devices expose capability surfaces under capability contracts

## Relationship To Stark

Stark is the most likely publisher of capability contracts.

At a high level:

1. a device is probed or otherwise discovered
2. its usable capability surface is identified
3. Stark publishes a capability contract for that surface
4. the kernel binds that capability contract into the runtime
5. nodes may use that capability only if their own node contract allows the relevant command surface

This keeps:

- discovery
- publication
- binding
- use

as separate concerns.

## Relationship To Runtime Commands

A capability contract should expose commands through the runtime command model.

That means device abilities become:

- external capability-bound commands

not:

- node contracts
- fake internal primitives

This fits the namespace-based runtime direction better.

## Short Framing

Node contracts govern nodes.

Capability contracts govern device-exposed external ability.

The cognitive system should treat those as different contract families.
