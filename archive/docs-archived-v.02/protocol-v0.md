# Protocol v0

## Purpose

This document captures the current working Skyra protocol direction.

It preserves the active shape before the final grammar is fully locked.

## Core Framing

The protocol should be forced into shape by the runtime model.

The current design pressure now comes from:

- typed stimulus as the unified runtime message family
- explicit actor authorship
- typed execution surfaces
- a small public primitive family
- a normalized ingress envelope at the callee boundary

That pressure no longer points at a large command taxonomy.

It points at a primitive-first outer protocol carrying explicit actor and surface address plus published stimulus payloads.

## Current Working Shape

The current working caller-facing outer shape is:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

Where:

- `<primitive>` names the requested public boundary mode
- `<actor>` names the owning public actor
- `<surface>` names the callable public surface beneath that actor
- `<stimulus_protocol>` is the concrete payload instance that conforms to a published stimulus contract

The important point is:

- runtime emission carries concrete payload
- publication in the database provides the contract/schema the payload must conform to
- actor + surface resolution identifies the matched callable sense schema

## Why `actor` Is Explicit

The protocol should make the owning actor visible.

That matters because:

- actor contracts are real runtime boundaries
- the kernel must know which actor owns the targeted surface
- audit trails should preserve both caller and first public callee

## Why `surface` Is Explicit

The public callable surface should be named directly on the wire.

That matters because:

- one actor may expose more than one callable public surface
- the kernel should resolve a specific public surface before producing `sense`
- mailbox routing should land on a matched `sense_schema_id`, not only a broad actor id

The wire-level `surface` name is the human-facing address that resolves to that matched callable sense schema.

## Why `primitive` Is Explicit

The public primitive layer remains small.

The current working public primitive split is:

- `recall`
- `learn`
- `observe`
- `act`

The important distinction is:

- the public wire primitive family remains `recall`, `learn`, `observe`, and `act`
- `sense` is not a caller-facing wire primitive
- `sense` is the normalized ingress envelope the kernel produces at the callee boundary

These remain different system boundaries even though the runtime is now stimulus-first.

## Public Primitive Split

### 1. `recall`

`recall` is the retained read path.

It brings retained experience into current runtime work.

### 2. `learn`

`learn` is the retained write path.

It writes selected consequence from completed runtime activity back into retained experience and structure.

### 3. `observe`

`observe` is the actor-side intake and admission path.

It admits relevant sensed input into bounded actor context when the actor chooses to take pending mailbox ingress into active runtime work.

### 4. `act`

`act` is the world-facing output or engagement path.

It is the primitive used when an actor such as a `stewart` crosses from internal stimulus handling to an outside execution surface.

## Kernel Ingress: `sense`

`sense` is the normalized ingress envelope at the callee boundary.

It is not the caller-facing wire primitive.

The caller still emits:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

After routing and validation, the kernel turns that caller-facing request into a callee-facing `sense` package.

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

The clean split is:

- caller emits public protocol under `recall`, `learn`, `observe`, or `act`
- kernel routes and normalizes that ingress into `sense`
- the target actor receives `sense`
- the target actor may admit relevant input through `observe` when ready

So:

- wire primitive = what the caller is asking for
- wire surface = the named callable ingress surface the caller is targeting
- `sense` = what has crossed into the callee boundary
- `observe` = actor-side admission of that input

The `payload` inside `sense` should conform to the published request schema for the matched stimulus contract.

The matched `sense_schema_id` is the mailbox-level routing authority.

The primitive should be derived from that matched sense schema rather than duplicated as another top-level mailbox field.

`source_timestamp` is when the signal happened at the source, if known.

`received_at` is when the kernel recorded the ingress into the mailbox.

`sense` should remain pending by default until the actor chooses whether to observe it.

## Published-Before-Use Rule

Public stimulus protocols should be published before runtime use.

That means:

- the registry holds the published contract
- the runtime carries concrete payload instances
- actors do not invent ad hoc public protocols mid-flight

## Kernel Routing

The kernel should route based on:

- the emitted stimulus payload
- the original emitter
- the target actor
- the target surface
- contract lookup in the database

At minimum, the kernel should be able to check whether:

- the target actor exists
- the named target surface exists beneath that actor
- the matched sense schema allows that request stimulus
- the payload conforms to the registered contract

Before mailbox delivery, the kernel should normalize incoming traffic into `sense`.

## Responses

The public request/response posture should currently be treated more like an API:

- one public request stimulus schema
- one public response envelope schema

This currently applies to the public callable surfaces under `recall`, `learn`, `observe`, and `act`.

The current response envelope should require:

- `status`
- `reason`

With `status` currently:

- `success`
- `failed`
- `timed_out`

Everything else may remain actor-defined.

## Delegation And Traversal

Actors may still traverse additional execution surfaces while fulfilling a request.

That traversal happens by emitting more Skyra protocol.

The important split is:

- caller sees the first/public callable surface
- callee may traverse more surfaces internally
- those internal traversals remain protocol-driven rather than hidden special cases

## Current Design Posture

The strongest current claims are:

- the protocol should remain primitive-first
- the caller-facing wire shape should name primitive, actor, and surface explicitly
- the public wire primitive family currently remains `recall`, `learn`, `observe`, and `act`
- `sense` is the normalized ingress envelope at the callee boundary
- the runtime should be stimulus-first
- public protocols should be published before use
- kernel routing should be driven by emitter context, actor + surface resolution, and DB-backed contract lookup
- public request/response surfaces should currently be modeled like small APIs

## Still Open

The following remain open:

- the final outer wire grammar around the inner stimulus payload
- the exact final kernel envelope fields
- the exact final object shape of `ExecutionSurface`
- whether `contract_publication` should also be normalized into stimulus

## Short Framing

The current protocol direction is:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

The outer protocol stays small.

The runtime carries concrete payload.

The registry carries the published contract that payload must satisfy.

The kernel turns received traffic into `sense` before mailbox delivery.
