# Sense v0

## Purpose

This document defines the current working meaning of `sense`.

The main point is:

- `sense` is the normalized ingress envelope at an actor boundary

## Core Framing

`sense` is how incoming traffic looks after it has crossed into a receiving actor boundary.

It is the boundary event for things such as:

- microphone input
- camera input
- webhook ingress
- actor-to-actor requests
- sensor updates
- callback-style external API input
- other external capability-originated signal

`sense` is not the same as `observe`.

It is also not the same thing as the caller-facing wire primitive.

## `sense` vs `observe`

The clean split is:

- a caller emits `recall`, `learn`, `observe`, or `act`
- the kernel routes and normalizes that ingress into `sense`
- `observe` is the actor-side admission path the actor may invoke when it chooses to take pending sensed input into bounded episode state and make it available for frame projection

So:

- `sense` = normalized ingress at the callee boundary
- `observe` = actor-side intake/admission

Arrival does not imply immediate observation.

## Actor Surface vs Capability Surface

Actors still expose the ordinary callable handler family:

- `observe`
- `act`
- `recall`
- `learn`

`sense` should not be treated as an ordinary public handler.

It is the callee-side runtime package the kernel produces before mailbox delivery.

Capability surfaces may still:

- publish `act`
- originate ingress signal that the kernel normalizes into `sense`

## Wire vs Ingress

The caller-facing outer shape remains:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

Where:

- `<primitive>` is one of `recall`, `learn`, `observe`, or `act`
- `<actor>` names the owning actor
- `<surface>` names the callable public surface beneath that actor

After routing, the kernel produces:

```ts
type SenseEnvelope = {
  kind: "sense"
  source_surface: ExecutionSurface
  target_actor: string
  sense_schema_id: string
  source_timestamp?: string
  received_at: string
  payload: unknown
}
```

The `payload` is the concrete stimulus payload that must satisfy the published request schema for the matched contract.

`sense_schema_id` identifies the matched callable sense schema for that actor.

The primitive is derived from that schema row rather than repeated on the mailbox envelope.

`source_timestamp` is when the signal happened at the source, if known.

`received_at` is when the kernel recorded the ingress into the mailbox.

## Mailbox Semantics

`sense` is the mailbox form of ingress.

The important rule is:

- `sense` should not expire by default

That means:

- the actor is not forced to `observe` immediately
- the actor may finish current work first
- the actor may observe later, batch several pending `sense` envelopes, or decide not to observe some at all

Time should remain visible rather than hidden in scheduler policy.

## Runtime Flow

The current useful flow is:

1. an actor or capability emits `skyra <primitive> <actor> <surface> <stimulus_protocol>`
2. the kernel validates and routes that request
3. the kernel normalizes the incoming request into `sense`
4. the target actor mailbox receives `sense`
5. when ready, the actor may handle admission through `observe`
6. episode state updates
7. frame projection may then include the admitted input

That means `sense` does not bypass the actor runtime model.

It is the ingress form inside it.

## Relationship To `act`

The matching outward capability boundary is still:

- `act`

So the useful capability-facing split is:

- capability ingress = normalized into `sense` at the receiving actor boundary
- `act` = Skyra to capability

## Current Design Posture

The strongest current claims are:

- `sense` is the normalized ingress envelope at the receiving actor boundary
- `sense` is not a public caller-facing primitive
- actor and capability ingress should both be normalized into `sense`
- `sense` should not expire by default
- `observe` remains the actor-side admission path and is actor-controlled rather than automatic
- capability surfaces should distinguish ingress from outbound `act`

## Short Framing

`sense` is how inbound traffic arrives at a receiving actor after kernel normalization.

`observe` is how an actor takes that typed input into its own bounded runtime context.
