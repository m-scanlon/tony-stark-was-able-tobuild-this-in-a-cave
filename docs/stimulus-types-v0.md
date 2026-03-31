# Stimulus Types v0

## Purpose

This document captures the current direction for typed stimulus in the node runtime.

The main point is:

- stimulus should be typed up front

## Core Framing

If nodes are going to call other nodes, the routing surface cannot remain vague.

The system needs to know:

- what a node can receive
- what a node can emit
- when one node's output is a valid input for another node

That means stimulus typing belongs in the contract model.

## Why Typed Stimulus Matters

Typed stimulus gives the runtime a clean basis for:

- routing
- orchestration
- delegation
- validation
- composability

Without typed stimulus, node-to-node execution becomes ad hoc and hard to reason about.

## Contract Surface

Each node contract should eventually declare:

- accepted stimulus types
- emitted stimulus types

Conceptually:

```ts
type NodeContract = {
  accepted_stimulus_types: string[]
  emitted_stimulus_types: string[]
}
```

The exact schema is still open.

The current important claim is just:

- stimulus typing is part of the contract, not an afterthought

## Stimulus As Routing Surface

Typed stimulus becomes the routing layer between nodes.

That means one node may:

1. receive typed stimulus
2. produce a typed output or derived stimulus
3. send that onward to another node

This is the basis for orchestrator behavior.

An orchestrator node therefore does not need ambient access to every node's memory.

It needs:

- typed inputs
- typed outputs
- routable contracts

## Working Shape

Conceptually:

```ts
type StimulusEnvelope = {
  stimulus_type: string
  source_node_id?: string
  target_node_id?: string
  intent_id?: string
  payload: Record<string, unknown>
}
```

This should be treated as directional rather than frozen.

The important parts are:

- the type is explicit
- the payload is attached to that type
- source and target may be preserved

## Examples

Example directions might later include stimulus types such as:

- `human_request`
- `device_probe_request`
- `device_probe_result`
- `recall_request`
- `recall_result`
- `registration_write_request`
- `registration_write_result`

These are examples only.

The current point is not the exact vocabulary.

The current point is:

- the vocabulary should exist

## Relationship To Protocol

Stimulus typing and protocol commands are related but not identical.

The split is:

- typed stimulus says what runtime event or message is being passed
- protocol commands say what operation a node is emitting under its contract

This allows the system to keep:

- message routing
- command execution

as separate but compatible layers.

## Current Design Posture

The strongest current claims are:

- nodes should accept typed stimulus
- nodes should emit typed stimulus
- typed stimulus should be declared in node contracts
- orchestrator behavior should be built on typed routing rather than ambient shared memory

## Short Framing

Stimulus should be typed up front.

That typing is part of the node contract and gives the system a real routing layer for orchestration and delegation.
