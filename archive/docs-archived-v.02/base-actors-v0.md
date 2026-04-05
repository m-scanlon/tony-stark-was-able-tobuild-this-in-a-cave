# Base Actors v0

## Purpose

This document captures the current role of base actors in the runtime.

The main point is:

- actors are the extensible unit of the system
- the product should still ship with a base layer of standard actors

## Core Framing

The current stack is becoming:

- primitives as the fixed substrate
- base actors as the shipped standard actor layer
- custom actors as extensions built on top
- workflows as compositions of actors over typed stimuli

Under this model, primitives stay small:

- `recall`
- `learn`
- `observe`
- `act`

Specialization should mostly happen at the actor layer rather than by endlessly growing the primitive set.

## Why Base Actors Exist

If actors are the extensible unit, the system still needs a useful starting library.

Otherwise the product ships with:

- abstract primitives
- but no practical behavior

Base actors solve that.

They give the runtime a standard layer of reusable operators that other actors can call, compose, and build on top of.

## What Base Actors Are

Base actors are shipped runtime actors that perform bounded specialized work using the shared primitive substrate.

They are not:

- hidden control planes
- one-off product hacks
- replacements for the primitive layer

They are:

- standard library actors
- reusable worker roles
- building blocks for workflows and orchestration

## Relationship To Primitives

Base actors do not replace primitives.

They use primitives.

Conceptually:

- a base actor may call `act`
- a base actor may call `observe`
- a base actor may call `recall`
- a base actor may call `learn`

The primitive set stays fixed and small.

The base actor layer is where concrete operational behavior becomes reusable.

## Current Example Base Actors

The likely early base layer includes actors such as:

- `probe`
- `registration`
- `contract_creator`
- later possibly `search`
- later possibly `response`

These examples are still provisional.

The important current point is:

- these are specialized worker actors
- not new primitives

## Example: `probe`

`probe` is now a strong example of a base actor under `Stark`.

Its job is to:

- discover candidate capabilities on a system subject
- verify them through bounded invocation
- preserve enough detail to register verified capability surfaces and later support `stewart` abstraction

Its job is not to:

- persist durable registration truth
- birth runtime actors
- perform unconstrained exploration

Conceptually, a useful current mock contract is:

```ts
type ProbeActorContract = {
  purpose: {
    summary: "Discover candidate capabilities on a system subject, verify them through bounded invocation, and preserve enough detail to register verified capability surfaces and support later steward abstraction."
  }
  commitments: [
    "Does not persist registration truth",
    "Does not birth actors",
    "Does not perform unconstrained exploration"
  ]
  request_stimuli: ["device_probe_request"]
  response_envelopes: ["device_probe_result_envelope"]
}
```

This example intentionally stays within the locked `v1` actor contract shape.

And the emitted result should carry enough material for later registration and contract publication, such as:

- discovered candidate capabilities
- verified capabilities
- registered capability-surface details
- probe strategy id
- confidence

## Example: `registration`

`registration` is now a strong companion base actor under `Stark`.

Its job is to:

- consume probe output for a system subject
- assemble the typed device registration envelope
- persist that envelope through `act` with `modality = registration_write`

Its job is not to:

- rediscover or reverify capability surfaces on its own
- publish public abstraction contracts from scratch without registration input
- birth runtime actors

Conceptually, a useful current mock contract is:

```ts
type RegistrationActorContract = {
  purpose: {
    summary: "Assemble the typed registration envelope for a system subject from probe output and persist that registration through the world-facing registration write path."
  }
  commitments: [
    "Does not rediscover capabilities on its own",
    "Does not publish public abstraction contracts from scratch",
    "Does not birth runtime actors"
  ]
  request_stimuli: ["device_probe_result"]
  response_envelopes: ["registration_write_result_envelope"]
}
```

And the registration write path should preserve enough material to recover durable truth such as:

- `subject`
- `transport`
- `probe_strategy`
- `verified_capabilities`
- `registration_state`
- `last_verified_at`

## Base Actor Behavior

A base actor should generally:

1. accept typed stimulus
2. perform bounded specialized work
3. emit further stimulus to runtime as needed
4. emit typed output or result stimulus back into the actor graph

That means base actors sit between:

- typed actor-to-actor routing
- primitive runtime execution

## Base Actors vs Workflows

Base actors are not the same thing as workflows.

The cleaner split is:

- base actors are reusable operators
- workflows are compositions of actors

So:

- `probe` may be a base actor
- "onboard a new device" is a workflow that may involve `probe`, `registration`, and other actors

This keeps the architecture modular.

## Base Actors vs Orchestrator Actors

Base actors should be distinguished from orchestrator actors.

Base actors:

- do bounded specialized work
- generally do not own broad routing authority
- are the reusable worker layer

Orchestrator actors:

- route typed stimuli
- delegate work to other actors
- merge or interpret results
- coordinate multi-actor execution

This distinction matters because otherwise "important actors" become overloaded with both orchestration and execution responsibilities.

## Product Surface

The product should likely ship with some base actors already defined.

That gives the system:

- a practical standard library
- a reusable execution layer
- something that custom actors and orchestrators can build on immediately

Without this, the system would have a protocol and primitives but too little useful runtime structure out of the box.

## Current Design Posture

The strongest current claims are:

- actors are the extensible unit of the system
- primitives should remain small and stable
- the product should ship with a base layer of standard actors
- base actors are the reusable worker layer built on the primitive substrate
- workflows should be compositions of actors rather than new primitives

## Short Framing

Base actors are the system's standard library of reusable worker actors.

They perform specialized work using the shared primitives and give the runtime a practical layer that orchestrators and custom actors can build on.
