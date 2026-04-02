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

This command string names the target actor.

It does not by itself carry command authorship.

Command authorship should instead travel in the surrounding kernel envelope as a separate `calling_actor` field.

The minimal kernel envelope is therefore:

```text
calling_actor: <actor>
command: skyra <node> <primitive> -<args> -reason "<why this command is being emitted>"
```

Examples:

```text
skyra jarvis act -target human -content "the user needs a response" -modality text -timestamp now -reason "the user needs a response"
```

```text
skyra stark act -target laptop -content "discover capability surface" -modality probe -timestamp now -reason "the device needs capability discovery"
```

```text
skyra stark act -target laptop -content "write verified device registration" -modality registration_write -timestamp now -reason "verified capability state must be persisted"
```

```text
skyra jarvis recall -entity terraform -top_k 8 -reason "the current stimulus introduced terraform as an active structural cue"
```

```text
skyra stark learn -episode_id ep_123 -reason "the just-closed episode should be consolidated into retained experience"
```

## Why `node` Is Explicit

The protocol should make the target actor visible.

That matters because:

- node contracts are real runtime boundaries
- different nodes have different purposes and permissions
- the kernel must know which actor should execute the command
- audit trails should preserve both caller identity and target identity

The protocol should therefore not pretend commands are targetless.

Authorship is still visible, but it belongs in the kernel envelope rather than in the command string itself.

## Why `primitive` Is Explicit

The current runtime is converging on a small primitive set rather than a large top-level command family taxonomy.

The current working split is:

- `recall`
- `learn`
- `observe`
- `act`

These are not identical operations.

They represent different system boundaries.

## Primitive Split

### 1. `recall`

`recall` reads retained experience into current runtime work.

It is the retained read path.

### 2. `learn`

`learn` writes from completed runtime activity into retained experience.

It is the retained write path.

### 3. `observe`

`observe` takes something in from the world.

It is the world-facing intake path.

Examples include:

- reading a screen
- listening to a human
- inspecting device state
- checking an external result surface

### 4. `act`

`act` is any intentional boundary engagement that produces an effect or signal.

Examples include:

- responding to a human
- probing a device
- searching the web
- calling an external API
- writing a registration record
- using a capability surface

This means the protocol is no longer best understood as:

- a generic `command_set`

It is better understood as:

- a node issuing one of a small number of primitives

## Relationship To The `act` Shape

The current locked `act` shape is:

```text
skyra <node> act \
  -target <target> \
  -content <content> \
  -modality <modality> \
  -timestamp <timestamp> \
  -reason "<why this command is being emitted>"
```

This gives the system one shared world-facing primitive without flattening every external action into one opaque blob.

For now, the strongest stable claims are:

- `node` should be explicit
- `primitive` should be explicit
- `act` should carry `target`, `content`, `modality`, `timestamp`, and `reason`
- every emitted command must include `-reason`

## Relationship To Typed Stimuli

The protocol now sits alongside typed node-to-node stimulus passing.

Nodes should not share one ambient memory pool.

Instead:

- nodes receive typed stimuli
- nodes emit typed stimuli
- nodes issue protocol commands under their contracts

This makes the protocol part of a typed runtime rather than a flat tool-call surface.

## Kernel Envelope

The command string is only one half of command dispatch.

The other half is the minimal kernel envelope surrounding it.

For `v1`, the useful current shape is:

- `calling_actor`
- `command`

The kernel should use `calling_actor` to load the caller contract from the database.

The kernel should then parse the command to determine:

- the target actor
- the primitive
- the primitive arguments

This keeps the transport envelope small while preserving an explicit authorization boundary.

## Delegation

Actors may invoke other actors.

That is valid behavior when:

- the target actor is explicitly allowed by the caller contract
- the requested primitive is allowed on that delegation edge
- the selected `modality`, when present, is allowed on that delegation edge

The kernel should not reject actor-to-actor invocation merely because the caller and target differ.

It should reject only unauthorized invocation.

The important split is:

- `calling_actor` = who is asking for the operation
- `<node>` in the command = which actor should execute the operation

This keeps delegation auditable without forcing the command grammar itself to grow a separate caller slot.

## Current Design Posture

The strongest current claims are:

- the protocol should be node-first
- the protocol should be primitive-first
- the main primitive split is `recall`, `learn`, `observe`, and `act`
- `observe` should absorb world-facing intake while `recall` and `learn` remain separate
- `act` should absorb world-facing action while `recall`, `learn`, and `observe` remain separate
- `act` should use the shape `target`, `content`, `modality`, `timestamp`, `reason`
- `calling_actor` should live in the surrounding kernel envelope rather than in the command string
- the node slot in the command should name the target actor
- `-reason` remains mandatory

## Still Open

The following remain open:

- the final node vocabulary
- the final primitive argument grammar
- the final `modality` vocabulary for `act`
- the final encoding shape for `content`
- the final timestamp conventions for `act`
- the exact command-result and writeback grammar

## Short Framing

The current protocol direction is:

```text
skyra <node> <primitive> -<args> -reason "<why this command is being emitted>"
```

This reflects a node-based runtime with a small primitive set rather than a flat command-family protocol.
