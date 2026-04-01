# Node Contract (v0)

## Core Contract Axioms

Every node exists under a contract.

The contract defines:

- why the node exists
- what capabilities it may rely on
- what typed stimuli it may respond to
- what typed stimuli it may emit
- what cognition envelope it may operate under
- what commands it may emit

A node does not act directly on the user, an API, or the runtime.

It emits commands that the system validates and executes.

At the contract level, the active boundary is:

- `purpose`
- `capabilities`
- `stimulus`
- `cognition`
- `commands`

## 1. Purpose

A node must have a defined reason for existing.

Purpose bounds:

- its role
- its responsibilities
- its limits

Purpose belongs to the node definition, not to the episode.

## 2. Capabilities

A node contract should name what capability surfaces the node is allowed to rely on.

This does not replace the capability contract itself.

It is the node-side allowance boundary over those capability surfaces.

## 3. Stimulus

A node may only be invoked by valid stimulus.

If incoming stimulus does not match the node's accepted form, the node is not eligible to act.

Stimulus therefore defines the node's input boundary.

The important current direction is:

- stimulus should be typed up front

That means a node contract should eventually define:

- accepted stimulus types
- emitted stimulus types

This is part of what makes node-to-node routing composable.

## 4. Cognition

Cognition is not free-standing inner autonomy.

It is the contract-bounded envelope within which the node may continue reasoning and choose the next command.

A node may emit a command that causes another inference or prompt step.

The exact budgeting and stop rules are not fully locked yet, but that envelope belongs to the contract.

## 5. Commands

A node does not interact directly.

A node emits commands.

That includes:

- user-facing output
- capability or API use
- commands that request another reasoning step

The current working protocol shape is:

```text
skyra <node> <primitive> -<args> -reason "<why this command is being emitted>"
```

The explicit node slot matters because command authorship is part of the runtime boundary.

The current primitive split is:

- `recall`
- `learn`
- `interact`

`interact` is therefore not a separate direct action path outside the command surface.

It is one primitive inside the contract-allowed command surface.

For example:

```text
skyra jarvis interact -method talk -target human -reason "the current frame requires an outward response"
```

`reason` should be treated as mandatory.

The node's emitted command surface is part of the system's audit trail.

That means:

- every emitted command must include an explicit rationale
- runtime should reject commands that omit `reason`
- `reason` explains why the node emitted the command
- `reason` does not replace later execution validation or evidence

## Contract Level vs Runtime Level

The contract says:

- what the node is for
- what capabilities it may rely on
- what can wake it up
- what typed stimuli it may emit
- how cognition is bounded
- what commands it may emit

Runtime execution then handles:

- actual command dispatch
- pending-command state
- command-result writeback
- episode-local state updates

So the command surface is part of the contract.

Execution mechanics remain part of runtime.

## Same Contract Model Across Nodes

This contract model applies across node roles.

That includes:

- user-facing or task-facing nodes
- `Jarvis` as the user-facing meaning node
- `Stark` as the structural node

What differs between nodes is role and allowed behavior, not the existence of a separate ontology.

## Episode Relation

A node contract bounds what may happen inside a node episode.

The contract does not store the episode itself.

Instead:

- the node contract is durable
- the node episode is bounded runtime participation under that contract

## Still Open

The following remain open even with the command surface inside the active contract:

- exact stimulus type vocabulary
- exact primitive argument grammar
- exact cognition budgeting and stop rules
- recall policy or recall defaults
- the final `interact` method taxonomy
- whether `channel` becomes canonical inside `interact`

Device-facing capability surfaces should be treated as capability contracts rather than as node contracts.

See also:

- [capability-contract-prelim.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/capability-contract-prelim.md)
- [protocol-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/protocol-v0.md)
- [stimulus-types-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/stimulus-types-v0.md)
- [node-memory-boundary-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/node-memory-boundary-v0.md)

## Short Framing

The node contract defines why a node exists, what capabilities and typed stimuli it may operate on, how cognition is bounded, and what commands it may emit.

It is the node's durable boundary.

Runtime execution happens inside episodes under that boundary.
