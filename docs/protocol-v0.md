# Protocol v0

## Purpose

This document captures the current working protocol direction.

It is not intended to be the final grammar.

It is intended to preserve the current architectural shift before the shape blurs again.

## Core Framing

The protocol should not be invented in the abstract.

It should be forced into shape by what the system actually needs to express.

The current design pressure now comes from:

- node-local retained experience
- typed node-to-node stimulus passing
- explicit node authorship
- a small primitive set

That pressure suggests a different protocol shape than the earlier `command_set` direction.

## Current Working Shape

The current working shape is:

```text
skyra <node> <primitive> -<args> -reason "<why this command is being emitted>"
```

Examples:

```text
skyra jarvis interact -method talk -target human -reason "the user needs a response"
```

```text
skyra stark interact -method probe -subject_id laptop -reason "the device needs capability discovery"
```

```text
skyra stark interact -method write_device_registration -subject_id laptop -reason "verified capability state must be persisted"
```

```text
skyra jarvis recall -entity terraform -top_k 8 -reason "the current stimulus introduced terraform as an active structural cue"
```

```text
skyra stark learn -episode_id ep_123 -reason "the just-closed episode should be consolidated into retained experience"
```

## Why `node` Is Explicit

The protocol should make authorship visible.

That matters because:

- node contracts are real runtime boundaries
- different nodes have different purposes and permissions
- orchestration depends on knowing which node emitted the command
- audit trails should preserve actor identity

The protocol should therefore not pretend commands are emitted from nowhere.

## Why `primitive` Is Explicit

The current runtime is converging on a small primitive set rather than a large top-level command family taxonomy.

The current working split is:

- `recall`
- `learn`
- `interact`

These are not identical operations.

They represent different system boundaries.

## Primitive Split

### 1. `recall`

`recall` reads retained experience into current runtime work.

It is the retained read path.

### 2. `learn`

`learn` writes from completed runtime activity into retained experience.

It is the retained write path.

### 3. `interact`

`interact` crosses a world boundary.

Examples include:

- talking to a human
- probing a device
- searching the web
- calling an external API
- writing a registration record
- using a capability surface

This means the protocol is no longer best understood as:

- a generic `command_set`

It is better understood as:

- a node issuing one of a small number of primitives

## Relationship To Methods

The current direction is that `interact` may carry method-specific specialization.

Examples:

- `talk`
- `probe`
- `search`
- `write_device_registration`

So a likely working shape is:

```text
skyra <node> interact -method <method> ... -reason "..."
```

This gives the system one world-facing primitive without flattening every external action into one opaque blob.

## `channel` Remains Open

Whether `interact` should also carry a first-class `-channel` field remains open.

`channel` may turn out to be useful.

But the current design should not freeze it too early.

For now, the strongest stable claims are:

- `node` should be explicit
- `primitive` should be explicit
- `interact` should carry world-facing methods
- every emitted command must include `-reason`

## Relationship To Typed Stimuli

The protocol now sits alongside typed node-to-node stimulus passing.

Nodes should not share one ambient memory pool.

Instead:

- nodes receive typed stimuli
- nodes emit typed stimuli
- nodes issue protocol commands under their contracts

This makes the protocol part of a typed runtime rather than a flat tool-call surface.

## Current Design Posture

The strongest current claims are:

- the protocol should be node-first
- the protocol should be primitive-first
- the main primitive split is `recall`, `learn`, and `interact`
- `interact` should absorb world-facing action while `recall` and `learn` remain separate
- `-reason` remains mandatory

## Still Open

The following remain open:

- the final node vocabulary
- the final primitive argument grammar
- whether `channel` becomes canonical inside `interact`
- the final method taxonomy for `interact`
- the exact command-result and writeback grammar

## Short Framing

The current protocol direction is:

```text
skyra <node> <primitive> -<args> -reason "<why this command is being emitted>"
```

This reflects a node-based runtime with a small primitive set rather than a flat command-family protocol.
