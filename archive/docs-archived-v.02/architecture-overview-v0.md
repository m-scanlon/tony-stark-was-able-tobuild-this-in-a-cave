# Architecture Overview v0

## Purpose

This document is the high-level map of the current architecture.

It is meant to make the active model easy to grasp before dropping into the more detailed contract and runtime docs.

## Current Source Of Truth

For current design work:

- `docs/` is canonical for active design language
- `skyra-v.1` is canonical for the current implementation surface
- `stimulus-unification-reference.md` is the current reference note for the stimulus-first cleanup

## Core Framing

The system is a local-first runtime organized around:

- kernel authority
- durable actors under contract
- published stimulus contracts
- bounded episodes
- projected frames
- selective retention

It is not primarily:

- a single global graph
- a prompt transcript
- a flat collection of callable tools

## `v1` Operating Theme

For `v1`, the main practical orientation is:

- experience continuously
- act when needed
- learn selectively

Most of the architecture exists to support those three concerns.

## Main Runtime Layers

The current working model is:

1. `Structure`
2. `Retention Layer`
3. `Kernel`
4. `Actor Contract`
5. `Actor`
6. `Episode`
7. `Frame`
8. `Stimulus Contracts And Runtime Artifacts`

### `Structure`

Structure is the canonical substrate.

It contains:

- entities
- relationships

Episodes and retained artifacts refer into structure.

They do not replace it.

### `Retention Layer`

The retention layer holds selective retained experience that survives beyond an episode.

The retained artifact family is:

- `retained_trace`
- `retained_understanding`
- `retained_salience`
- `retained_tension`

All retained artifacts share an `anchor_set` into canonical structure.

### `Kernel`

The kernel is the runtime authority.

It is responsible for:

- birthing actors
- owning unified typed event intake
- validating routed stimulus against published contracts
- routing toward the correct execution surface
- keeping runtime actor registration live

The kernel births `Stark` at startup from a hardcoded contract.

Later actors are born from published contracts through the same kernel-controlled birth path.

### `Actor Contract`

Every actor exists under a durable contract.

The current contract center is:

- `purpose`
- `commitments`
- request stimuli
- response envelopes

The contract defines:

- why the actor exists
- what durable commitments it carries
- what public request stimuli it can receive
- what public response envelopes it can emit

### `Actor`

A actor is the durable runtime operator acting under a contract.

The actor owns runtime machinery such as:

- mailbox handling
- event handling behavior
- recall machinery
- frame assembly behavior
- `dependencyLedger` tracking
- the pointer to the active episode

### `Episode`

An episode is the bounded runtime state container for one span of activity.

The current core episode sections are:

- `purpose`
- `interaction_history`
- `recall`

### `Frame`

The frame is the bounded inference page projected from the episode.

The current frame layout is:

1. `purpose`
2. `interaction`
3. `recall`

The frame is not the durable owner of truth.

It is the current inference projection.

### `Stimulus Contracts And Runtime Artifacts`

Runtime execution now happens through emitted stimulus rather than older command/writeback pairs.

The outer protocol direction is:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

The registry holds the published contract.

Runtime traffic carries concrete payload instances that conform to that contract.

Runtime artifacts remain transient outputs produced while fulfilling that traffic.

## Core Boundary Rules

The strongest current boundary rules are:

- contract bounds public behavior
- actor owns machinery
- episode owns bounded runtime state
- frame is projected from the episode
- kernel remains routing authority
- recall is a read path
- learning is a write path

## Main Role Split

### `Jarvis`

`Jarvis` is the user-facing meaning and attention actor.

### `Stark`

`Stark` is the structural actor.

Its concern is:

- actor topology
- contract publication
- registry authority
- structural revision

`Stark` may also name contracts at publication time.

### `stewart`

A `stewart` actor is the execution boundary that mediates an external surface.

Its concern is:

- learning the shape of an external capability
- simplifying that request complexity into a published public stimulus contract
- handling world-facing interaction under `act`

## Execution Surface Model

`ExecutionSurface` is the parent routing concept.

The current child surface types are:

- `actor`
- `capability`

Both surface types are registered.

Actor surfaces publish:

- request stimulus schema
- response envelope schema

Capability surfaces publish `act` contracts and may also originate ingress that the kernel later normalizes into `sense`.

The public callable surface is itself an `ExecutionSurface`.

Handlers may then traverse additional execution surfaces internally by emitting more Skyra protocol.

## Runtime Flow

The high-level runtime flow is:

1. external or internal input becomes typed stimulus
2. the kernel validates that stimulus against the published contract
3. the kernel routes it to the target execution surface
4. before mailbox delivery, the kernel normalizes incoming traffic into `sense`
5. the target actor opens or reuses an episode
6. the actor admits relevant sensed input into bounded episode state
7. inference projects a frame and may emit further stimulus
8. downstream execution surfaces are traversed as needed
9. a response envelope or other returned stimulus is routed back
10. learning later decides what should survive beyond the episode

## Current Design Posture

The strongest current claims are:

- the runtime is stimulus-first
- the registry should be a stimulus registry
- actor contracts expose request/response public surfaces
- actors may hide downstream complexity behind those public surfaces
- execution surfaces remain explicit and typed

## Short Framing

The system is a kernel-routed stimulus runtime.

Actors expose public request/response surfaces under contract, episodes hold bounded state, frames support inference, and downstream execution-surface traversal remains hidden behind the public API-shaped contract.
