# Stimulus Types v0

## Purpose

This document captures the current direction for typed stimulus in the runtime.

The main point is:

- stimulus should be typed up front

## Core Framing

Typed stimulus is a runtime boundary of its own.

It should not be owned by the node boundary.

The system needs to know:

- what kind of stimulus is arriving
- where that stimulus came from
- what payload it carries
- which nodes accept or emit that type

That means stimulus typing belongs in a top-level stimulus contract family.

## Why Typed Stimulus Matters

Typed stimulus gives the runtime a clean basis for:

- routing
- orchestration
- validation
- composability

Without typed stimulus, execution becomes ad hoc and hard to reason about.

## Node Contract Surface

Each node contract should declare:

- accepted stimulus types
- emitted stimulus types

Conceptually:

```ts
type NodeContract = {
  accepted_stimulus_types: string[]
  emitted_stimulus_types: string[]
}
```

The node contract should name stimulus types.

It should not own the top-level stimulus envelope.

## Top-Level Stimulus Shape

Conceptually:

```ts
type StimulusSource = {
  node_id?: string
  capability_id?: string
}
```

```ts
type StimulusEnvelope = {
  stimulus_type: string
  source: StimulusSource
  payload: Record<string, unknown>
}
```

Rule:

- exactly one source path should be set

That means typed stimulus may come from:

- a capability
- a node

This avoids baking node-to-node routing assumptions into every stimulus.

## Stimulus Registry

The runtime should also have a structural registry of known stimulus types.

Conceptually:

```ts
type StimulusType = {
  type_id: string
  description?: string
  schema?: Record<string, unknown>
}
```

```ts
type StimulusRegistry = {
  version?: string
  types: StimulusType[]
}
```

The registry is the structural record of:

- what stimulus types exist
- what schema shape they carry

## Stark As Type Authority

`Stark` is the structural authority over stimulus typing.

That means `Stark` owns:

- stimulus type creation
- stimulus type revision
- the live stimulus registry
- classification of raw incoming capability output into registered stimulus types

Nodes do not invent ambient unregistered stimulus kinds on their own.

They consume and emit registered types under contract.

## Examples

Example stimulus types might include:

- `human_request`
- `bootstrap_fingerprint`
- `device_probe_request`
- `device_probe_result`
- `registration_write_request`
- `registration_write_result`

These are examples only.

The important point is:

- the vocabulary should exist in a registry

## Relationship To Protocol

Stimulus typing and protocol commands are related but not identical.

The split is:

- typed stimulus says what runtime package is being passed
- protocol commands say what operation a node is emitting under its contract

This allows the system to keep:

- message typing
- command execution

as separate but compatible layers.

## Current Design Posture

The strongest current claims are:

- nodes should accept typed stimulus
- nodes should emit typed stimulus
- node contracts should declare stimulus types, not own the envelope
- the stimulus envelope is a top-level boundary
- `Stark` owns the live stimulus registry and stimulus typing authority

## Short Framing

Stimulus should be typed up front.

The top-level stimulus contract defines the envelope, source, and registry.

`Stark` owns the live registry.

Nodes only declare which types they accept and emit.
