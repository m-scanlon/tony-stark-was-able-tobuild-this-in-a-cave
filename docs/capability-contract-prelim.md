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
- what capabilities it may rely on
- what stimuli may invoke it
- what cognition envelope it may operate under
- what commands it may emit

A node contract is about bounded runtime participation.

### Capability Contract

A capability contract governs an external capability surface.

It defines:

- what invocation surface that capability exposes
- what operations on that surface may be invoked
- what inputs those operations accept
- what outputs or result types they may return
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

Stark is the publisher of capability contracts.

At a high level:

1. a device is probed or otherwise discovered
2. its usable capability surface is identified
3. Stark publishes a capability contract for that surface
4. the kernel binds that capability contract into the runtime
5. nodes may target that capability only if their own node contract allows the relevant primitive and invocation path

This keeps:

- discovery
- publication
- binding
- use

as separate concerns.

Probe should now also be understood as the first contract-formation step for a capability surface.

That means:

1. a candidate capability is discovered
2. bounded invocation is attempted
3. observed behavior shapes the initial invocation surface
4. Stark publishes the resulting capability contract

Later use and learning may refine that contract over time.

## Relationship To Runtime Commands

A capability contract should not itself become a node command surface.

Instead, it should define the router-facing invocation surface that runtime commands can target.

That means device abilities become:

- external capability-bound invocation targets

not:

- node contracts
- fake internal primitives

This fits the current runtime direction better:

- node-first
- primitive-first
- invocation-surface-based

## Short Framing

Node contracts govern nodes.

Capability contracts govern device-exposed external ability.

The cognitive system should treat those as different contract families.
