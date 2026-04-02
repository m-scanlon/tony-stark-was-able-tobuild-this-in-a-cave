# Base Nodes v0

## Purpose

This document captures the current role of base nodes in the runtime.

The main point is:

- nodes are the extensible unit of the system
- the product should still ship with a base layer of standard nodes

## Core Framing

The current stack is becoming:

- primitives as the fixed substrate
- base nodes as the shipped standard node layer
- custom nodes as extensions built on top
- workflows as compositions of nodes over typed stimuli

Under this model, primitives stay small:

- `recall`
- `learn`
- `observe`
- `act`

Specialization should mostly happen at the node layer rather than by endlessly growing the primitive set.

## Why Base Nodes Exist

If nodes are the extensible unit, the system still needs a useful starting library.

Otherwise the product ships with:

- abstract primitives
- but no practical behavior

Base nodes solve that.

They give the runtime a standard layer of reusable operators that other nodes can call, compose, and build on top of.

## What Base Nodes Are

Base nodes are shipped runtime nodes that perform bounded specialized work using the shared primitive substrate.

They are not:

- hidden control planes
- one-off product hacks
- replacements for the primitive layer

They are:

- standard library nodes
- reusable worker roles
- building blocks for workflows and orchestration

## Relationship To Primitives

Base nodes do not replace primitives.

They use primitives.

Conceptually:

- a base node may call `act`
- a base node may call `observe`
- a base node may call `recall`
- a base node may call `learn`

The primitive set stays fixed and small.

The base node layer is where concrete operational behavior becomes reusable.

## Current Example Base Nodes

The likely early base layer includes nodes such as:

- `probe`
- `registration`
- `contract_creator`
- later possibly `search`
- later possibly `response`

These examples are still provisional.

The important current point is:

- these are specialized worker nodes
- not new primitives

## Example: `probe`

`probe` is now a strong example of a base node under `Stark`.

Its job is to:

- discover candidate capabilities on a system subject
- verify them through bounded invocation
- shape the initial capability contracts from observed behavior

Its job is not to:

- persist durable registration truth
- birth runtime nodes
- perform unconstrained exploration

Conceptually, a useful current mock contract is:

```ts
type ProbeNodeContract = {
  node_type: "probe"
  purpose: {
    summary: "Discover candidate capabilities on a system subject, verify them through bounded invocation, and shape initial capability contracts from observed behavior."
    limits: [
      "Does not persist registration truth",
      "Does not birth nodes",
      "Does not perform unconstrained exploration"
    ]
  }
  stimulus: {
    accepted_types: ["device_probe_request"]
    emitted_types: ["device_probe_result"]
  }
  cognition: {
    mode: "bounded_probe"
    max_steps: 1
  }
  commands: {
    allowed_commands: ["act", "recall"]
  }
  learning_enabled: true
}
```

And the emitted result should carry enough material for later registration and contract publication, such as:

- discovered candidate capabilities
- verified capabilities
- shaped initial capability contracts
- probe strategy id
- confidence

## Example: `registration`

`registration` is now a strong companion base node under `Stark`.

Its job is to:

- consume probe output for a system subject
- assemble the typed device registration envelope
- persist that envelope through `act` with `modality = registration_write`

Its job is not to:

- rediscover or reverify capability surfaces on its own
- shape initial capability contracts from scratch
- birth runtime nodes

Conceptually, a useful current mock contract is:

```ts
type RegistrationNodeContract = {
  node_type: "registration"
  purpose: {
    summary: "Assemble the typed registration envelope for a system subject from probe output and persist that registration through the world-facing registration write path."
    limits: [
      "Does not rediscover capabilities on its own",
      "Does not shape initial capability contracts from scratch",
      "Does not birth runtime nodes"
    ]
  }
  stimulus: {
    accepted_types: ["device_probe_result"]
    emitted_types: ["registration_write_result"]
  }
  cognition: {
    mode: "bounded_registration"
    max_steps: 1
  }
  commands: {
    allowed_commands: ["act", "recall"]
  }
  learning_enabled: true
}
```

And the registration write path should preserve enough material to recover durable truth such as:

- `subject`
- `transport`
- `probe_strategy`
- `verified_capabilities`
- `registration_state`
- `last_verified_at`

## Base Node Behavior

A base node should generally:

1. accept typed stimulus
2. perform bounded specialized work
3. emit primitive commands to runtime as needed
4. emit typed output or result stimulus back into the node graph

That means base nodes sit between:

- typed node-to-node routing
- primitive runtime execution

## Base Nodes vs Workflows

Base nodes are not the same thing as workflows.

The cleaner split is:

- base nodes are reusable operators
- workflows are compositions of nodes

So:

- `probe` may be a base node
- "onboard a new device" is a workflow that may involve `probe`, `registration`, and other nodes

This keeps the architecture modular.

## Base Nodes vs Orchestrator Nodes

Base nodes should be distinguished from orchestrator nodes.

Base nodes:

- do bounded specialized work
- generally do not own broad routing authority
- are the reusable worker layer

Orchestrator nodes:

- route typed stimuli
- delegate work to other nodes
- merge or interpret results
- coordinate multi-node execution

This distinction matters because otherwise "important nodes" become overloaded with both orchestration and execution responsibilities.

## Product Surface

The product should likely ship with some base nodes already defined.

That gives the system:

- a practical standard library
- a reusable execution layer
- something that custom nodes and orchestrators can build on immediately

Without this, the system would have a protocol and primitives but too little useful runtime structure out of the box.

## Current Design Posture

The strongest current claims are:

- nodes are the extensible unit of the system
- primitives should remain small and stable
- the product should ship with a base layer of standard nodes
- base nodes are the reusable worker layer built on the primitive substrate
- workflows should be compositions of nodes rather than new primitives

## Short Framing

Base nodes are the system's standard library of reusable worker nodes.

They perform specialized work using the shared primitives and give the runtime a practical layer that orchestrators and custom nodes can build on.
