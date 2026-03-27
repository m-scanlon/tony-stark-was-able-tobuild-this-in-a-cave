# Node Contract (v0)

## Core Contract Axioms

Every node exists under a contract.

The contract defines:

- why the node exists
- what stimuli it may respond to
- what outward interaction forms it may emit

At the contract level, the core boundary remains:

- `purpose`
- `stimulus`
- `interact`

## 1. Purpose

A node must have a defined reason for existing.

Purpose bounds:

- its role
- its responsibilities
- its limits

Purpose belongs to the node definition, not to the episode.

## 2. Stimulus

A node may only be invoked by valid stimulus.

If incoming stimulus does not match the node's accepted form, the node is not eligible to act.

Stimulus therefore defines the node's input boundary.

## 3. Interact

A node may only emit valid interact output.

Its output must conform to forms the kernel and surrounding system can receive and interpret.

Interact therefore defines the node's outward action boundary.

## Contract Level vs Runtime Level

Contract primitives are not the same thing as runtime callable commands.

The contract says:

- what the node is for
- what can wake it up
- what kinds of outward result it may emit

Runtime commands are in-episode operations the node may issue through the kernel during runtime execution.

Those belong to runtime execution, not the core contract surface.

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

## Possible Future Additions

Potential contract-adjacent additions may later include ideas such as:

- allowed runtime command namespaces and commands
- recall policy or recall defaults
- capability constraints

Those are not part of the active v0 contract surface yet.

## Short Framing

The node contract defines why a node exists, what stimulus it may act on, and what interact output it may emit.

It is the node's durable boundary.

Runtime execution happens inside episodes under that boundary.
