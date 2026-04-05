# Actor Substrate Interface v0

## Core Framing

The actor substrate is the underlying runtime execution surface of an actor.

It is not:

- the actor contract
- the episode
- the frame

It is the runtime layer that lets an actor:

- accept routed ingress
- update episode state
- project frames
- emit stimulus
- receive published contract updates from Stark

This is the minimal interface the actor process should sit on top of.

## Design Goal

The substrate should stay generic enough to support:

- one-shot execution
- bounded multi-step execution
- contract-bounded loops
- different actor implementations over time

It should not hardcode one reasoning loop such as:

- `ReAct`
- `OODA`

Those should be supported by the substrate, not baked into it.

## Minimal Interface

```ts
type ActorSubstrate = {
  actor_id: string
  contract: ActorContractSnapshot
  active_episode_id?: string

  ingest_event(event: ActorEvent): ActorUpdateResult
  project_frame(): Frame | null
  emit_stimulus(envelope: StimulusEnvelope): StimulusEmitResult
  receive_published_contract(contract: ActorContractSnapshot): ContractPublicationResult
}
```

This is the current minimal public surface.

## Why These Methods Exist

### `ingest_event(...)`

Accepts incoming runtime events such as:

- normalized `sense` ingress
- contract publication

This is the main event intake boundary.

The main split is:

- `sense` is the callee-side ingress envelope
- `observe` remains the actor-side admission handler the actor may invoke after that ingress has reached the actor runtime boundary

### `project_frame()`

Builds the current inference page from active episode state.

The frame is not the source of truth.

It is a projection from the episode.

### `emit_stimulus(...)`

Emits caller-facing stimulus into the runtime.

This does not assume synchronous completion.

The actor may emit stimulus toward:

- another actor execution surface
- a capability execution surface

The substrate should not care which loop style produced that decision.

It should only care about routing the emitted stimulus through the kernel-facing runtime boundary.

### `receive_published_contract(...)`

Receives a newly published contract snapshot for the actor.

In practice, this is the actor-side receipt path for a contract publication originated by Stark.

The actor should then:

- record that contract in pending actor state
- keep the current active contract in force while the current episode is still open
- let the pending contract take effect only after the current episode closes in `v1`

## Public Stimulus Envelope

The active runtime model is stimulus-first.

Conceptually, the substrate should assume a caller-facing stimulus envelope that carries:

- the original emitter
- the target execution surface
- the requested primitive
- the concrete stimulus payload

Conceptually:

```ts
type StimulusEnvelope = {
  emitter_surface: ExecutionSurface
  target_surface: ExecutionSurface
  primitive: "recall" | "learn" | "observe" | "act"
  payload: unknown
}
```

## Sense Envelope

Before a routed request is delivered to an actor mailbox, the kernel should normalize it into `sense`.

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

The exact final envelope shape is still open.

The important point is:

- actors emit public stimulus, not older command/writeback pairs
- actors receive normalized `sense` at the mailbox boundary

The current useful distinction is:

- actors emit `recall`, `learn`, `observe`, or `act`
- actor and capability ingress are both normalized into `sense`
- `sense` remains pending in the mailbox until the actor chooses whether to `observe` it
- `sense_schema_id` identifies which callable sense schema was matched for this ingress
- the `payload` inside `sense` must conform to the matched published request schema
- `source_timestamp` and `received_at` let the actor reason about freshness and staleness explicitly

## Event Shape

The substrate should assume a typed event model.

Conceptually:

```ts
type ActorEvent =
  | SenseEnvelope
  | ContractPublicationEvent
```

This does not force a final event schema yet.

It only establishes that the substrate is event-driven rather than step-machine-driven.

## Event Heap / Mailbox Model

Because stimulus emission and returned response handling are separate, the runtime still needs a place for deferred work and contract publication to land.

`v1` should use one unified typed event intake surface at the kernel front rather than separate global queues.

At minimum, that heap should be able to carry:

- incoming `sense` envelopes
- published contract events

Routing should remain thin.

`v1` should not force expiry of `sense` by default.

Age should instead remain available to the actor through the timestamp fields carried on each pending `sense`.

At minimum, the router should:

- accept a typed package or event
- resolve the target actor through the live actor registry
- place the routed package into that actor's mailbox

Once an event has been routed to an actor, it should land in a lightweight actor-local mailbox.

That mailbox does not need to be another scheduler.

It is just the pending holding area for already-routed work waiting on that actor.

## Relationship To Other Objects

- the actor contract bounds what the actor may do publicly
- the actor substrate provides the runtime execution surface
- the episode holds bounded runtime state
- the frame is projected from episode state

So:

- contract = boundary
- substrate = runtime surface
- episode = state container
- frame = inference projection

## Current Design Posture

The strongest current claims are:

- the actor substrate should expose a small explicit interface
- runtime intake should be typed and event-driven
- actors should emit stimulus rather than older command/writeback pairs
- incoming traffic should be normalized into `sense` before mailbox delivery
- the substrate should support multiple loop styles without hardcoding one

## Short Framing

The actor substrate is the minimal runtime interface beneath the actor process.

Its current core surface is:

- `ingest_event(...)`
- `project_frame()`
- `emit_stimulus(...)`
- `receive_published_contract(...)`

That surface should stay generic while the higher-level stimulus contract model continues to settle.
