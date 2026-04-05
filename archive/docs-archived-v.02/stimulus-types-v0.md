# Stimulus Types v0

## Purpose

This document captures the current direction for typed `stimulus` in the runtime.

The main point is:

- `stimulus` is now the unified runtime message family

## Core Framing

The system should think in terms of:

- actors emit typed `stimulus`
- actors receive typed `stimulus`
- the kernel routes typed `stimulus`

This applies across:

- user-facing traffic
- actor-to-actor traffic
- capability-facing traffic
- external ingress traffic
- response traffic

The runtime should therefore be understood as stimulus-first rather than command-first.

## Execution Surfaces

Every registered stimulus contract carries an `ExecutionSurface`.

`ExecutionSurface` is the routing surface for that stimulus.

The initial typed surface kinds are:

- `actor`
- `capability`

Conceptually:

```ts
type ExecutionSurface = {
  kind: "actor" | "capability"
  id: string
}
```

## Registered Stimulus Contract

The registry should store published stimulus contracts.

Conceptually:

```ts
type StimulusContract = {
  type_id: string
  name: string
  description?: string
  primitive: "observe" | "act" | "recall" | "learn"
  execution_surface: ExecutionSurface
  request_schema: Record<string, unknown>
  response_envelope_schema: Record<string, unknown>
}
```

```ts
type StimulusRegistry = {
  version?: string
  contracts: StimulusContract[]
}
```

The important point is:

- the registry stores contract records that already know where they execute
- the registry should also preserve which public boundary mode they represent

## Primitive Tag

Registered stimulus contracts should make the public boundary mode explicit.

The useful current set is:

- `observe`
- `act`
- `recall`
- `learn`

The clean split is:

- actor surfaces publish `observe`, `act`, `recall`, or `learn`
- capability surfaces publish `act` contracts and may also originate ingress that the kernel normalizes into `sense`

This matters because ingress normalization should not be confused with a public caller-facing contract primitive.

## Ingress Normalization

`sense` is still a real runtime object, but it should now be understood as the kernel-defined ingress envelope rather than a published public primitive.

Conceptually:

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

The `payload` is the concrete request payload that must satisfy the published request schema for the matched stimulus contract.

`sense_schema_id` identifies the matched callable sense schema for the target actor.

The primitive should be derived from that schema row rather than duplicated on the normalized mailbox envelope.

`source_timestamp` is when the signal happened at the source, if known.

`received_at` is when the kernel recorded the ingress into the mailbox.

## Request And Response

For the current simplified API posture, each public callable surface should be modeled as:

- one request schema
- one response envelope schema

This currently applies to the public callable surfaces under `observe`, `act`, `recall`, and `learn`.

The response envelope should currently require:

- `status`
- `reason`

The current `status` enum is:

- `success`
- `failed`
- `timed_out`

Everything else in the response payload may be invented by the actor as needed.

## Publication Rule

Stimulus protocols should be published before runtime use.

That means:

- actors do not invent ad hoc public protocols mid-flight
- an actor is born under contract
- the callable request/response contracts it handles are already registered

## Runtime Emission

The outer Skyra protocol should currently be understood as:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

The runtime message carries a concrete payload that conforms to a published stimulus contract.

The schema itself does not need to travel inline every time.

The important split is:

- the registry holds the published contract
- the runtime carries concrete payload instances that conform to it
- actor + surface resolution identify the matched callable sense schema
- the kernel normalizes received traffic into `sense` before mailbox delivery

## Kernel Routing

The kernel should route based on:

- the emitted stimulus
- the original emitter
- the target `ExecutionSurface`
- contract lookup in the database

At minimum, the kernel should be able to check whether:

- the target actor or capability exists
- the emitted stimulus is registered for that surface
- the request payload conforms to the published contract

Routing should stay thin.

## Stark As Registry Authority

`Stark` is the structural authority over stimulus publication.

That means `Stark` owns:

- stimulus contract creation
- contract naming at publication time
- revision/publication of the live registry
- the structural side of `stewart` actor birth when new external surfaces are revealed

## Capability Submission

When an external capability is submitted for registration:

1. a `stewart` actor is born under contract
2. that actor learns the external request complexity
3. the relevant capability-facing `act` contracts and ingress schemas are registered
4. that actor publishes a simplified actor-facing stimulus contract into the registry
5. that actor later mediates world-facing interaction under `act` and receives external ingress normalized into `sense`

The external surface is therefore not the general public actor-facing interface.

The public interface is the published stimulus contract the system routes through.

## Current Design Posture

The strongest current claims are:

- the registry should be a `stimulus registry`
- every registered stimulus contract carries an `ExecutionSurface`
- every registered stimulus contract should also carry its boundary mode / primitive tag
- the public primitive family currently includes `observe`, `act`, `recall`, and `learn`
- the initial execution-surface kinds are `actor` and `capability`
- actor surfaces publish request/response contracts
- capability surfaces publish `act` contracts and may also preserve ingress schemas
- actors emit and receive typed `stimulus`
- capability ingress should be normalized into `sense`
- capability outbound should always be modeled explicitly as `act`
- public callable surfaces should currently be modeled as one request plus one response envelope

## Short Framing

The runtime should think in terms of published stimulus contracts.

Those contracts live in the stimulus registry, carry their own `ExecutionSurface`, and define the request/response shape that routed runtime payloads must follow.

The kernel turns incoming traffic into `sense` before the target actor sees it.
