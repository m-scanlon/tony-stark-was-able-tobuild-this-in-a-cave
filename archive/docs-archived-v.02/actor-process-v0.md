# Actor Process v0

## Core Framing

The actor process is the live runtime behavior of a actor after birth.

It is not:

- the actor contract
- the actor substrate interface
- the episode
- the frame

It is the runtime behavior that:

- receives routed stimulus
- maintains an active episode
- updates bounded runtime state
- projects frames for inference
- emits stimulus
- tracks open downstream dependencies
- applies published contracts at episode boundaries

So:

- contract = boundary
- substrate = public runtime surface
- process = live behavior
- episode = bounded state
- frame = inference projection

## Relationship To Actor Birth

Actor birth creates a live actor under a contract.

The actor process begins once that actor exists.

A newly born actor starts:

- live
- bound by its contract
- without an active episode
- without `dependencyLedger` state

The actor process remains idle until a valid event arrives.

## Inputs To The Actor Process

The current `v1` posture is:

- the kernel owns the global event intake and routing surface
- the runtime carries typed events
- the kernel validates those events and routes them to the correct actor
- each actor owns a lightweight mailbox for already-routed events

At minimum, the relevant event families are:

- `stimulus`
- `contract_publication`

If the runtime still needs internal completion receipts, they should be normalized into typed response stimulus before they reach the actor boundary.

## Core Responsibilities

The actor process is responsible for:

- draining its mailbox
- checking incoming stimulus against the active contract
- opening or reusing an episode
- writing event effects into episode state
- maintaining the `dependencyLedger`
- projecting frames when the actor is inference-ready
- emitting stimulus allowed by contract
- resolving obligations through returned response envelopes
- closing episodes after inactivity
- adopting published contracts only after episode closure

## High-Level Runtime Flow

The current high-level process is:

1. a typed runtime event arrives at the kernel
2. the kernel validates and routes it to the correct actor mailbox
3. the actor process takes the next routed event
4. the actor checks whether the event is valid under the active contract
5. if valid, the actor opens or reuses an episode
6. the event is written into episode-local state
7. background actor machinery may update episode-local derived artifacts
8. when the actor is inference-ready, it projects a frame
9. inference selects the next allowed primitive and stimulus payload
10. the actor emits Skyra protocol carrying that stimulus
11. the kernel validates the target execution surface and payload against the published contract
12. the target surface handles the request
13. a response envelope or other returned stimulus is routed back
14. the actor writes that returned stimulus into episode state and resolves the matching dependency/obligation as needed
15. after inactivity, the episode closes
16. any newly published contract becomes active only after episode closure

## Contract Gating

The actor process should not treat every routed event as valid.

The active contract remains the gate.

At minimum, the contract currently bounds:

- purpose
- commitments
- request stimuli
- response envelopes

So the process rule is:

- routed event arrives
- actor checks contract validity
- unsupported or invalid payloads are rejected, ignored, or deferred by runtime policy
- valid events are allowed to affect the active episode

## Episode Management

The episode is the bounded state container the actor process acts on.

`v1` episode policy is time-based.

That means:

- if there is no active episode, open one
- if the active episode has gone inactive long enough to close, open a new one
- otherwise reuse the current active episode

Episode closure is also the natural learning handoff point.

For `v1`, the owning actor may emit learning stimulus against the just-closed episode rather than reopening the closed episode as ordinary active runtime state.

## Event Write Paths

Different routed event types affect different parts of actor and episode state.

### `stimulus`

At minimum, a stimulus event should:

- append to `interaction_history`
- become eligible as current interaction context
- trigger any actor-owned background machinery that updates episode-local derived artifacts

That may include:

- structural projection updates
- episode-field updates
- recall updates
- response-envelope handling

The actor owns the machinery.

The episode owns the resulting bounded state.

### `contract_publication`

A newly published contract should not switch the actor mid-episode.

Instead:

- Stark publishes the new contract into the runtime
- the kernel schedules and routes that publication like other typed events
- the new contract is received by the running actor
- the actor holds that contract in pending actor state / mailbox flow
- the current episode continues under the currently active contract
- the new contract becomes active only after the current episode closes

## Dependency Ledger

The actor process owns the `dependencyLedger`.

This is not a public contract concern.

For `v1`, the main practical requirement is:

- the actor knows whether a downstream dependency is currently open

The minimum useful internal distinctions are:

- `open`
- `resolved`
- `failed`
- `timed_out`

Each `dependency` entry should be keyed at minimum by a runtime dispatch identifier.

## Frame Projection

The frame remains a projection from episode state.

The actor process owns frame assembly behavior.

The episode remains the source of truth.

The frame is what inference sees.

The timing question remains open.

The better framing is:

- not "when should the actor render a frame"
- but "when has the actor experienced enough to justify an inference step"

## Skyra Protocol Traversal

The actor process should assume the stimulus-first Skyra protocol:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

That means:

- the public callable surface is itself an `ExecutionSurface`
- the handler may traverse additional execution surfaces internally
- those downstream traversals happen by emitting more Skyra protocol

The actor does not need hidden special-case invocation machinery beneath that.

## Relationship To Other Actor Docs

- [actor-birth-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-birth-v0.md) defines how a actor comes into existence
- [actor-substrate-interface-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-substrate-interface-v0.md) defines the minimal public runtime surface
- [actor-and-episode-ownership-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-and-episode-ownership-v0.md) defines ownership boundaries
- [actor-open-questions-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-open-questions-v0.md) tracks remaining unresolved actor questions

## Current Design Posture

The strongest current claims are:

- the actor process is event-driven
- the kernel owns global routing
- each actor owns a lightweight mailbox
- the actor process opens or reuses episodes by inactivity/time policy
- actors emit stimulus rather than older command/writeback pairs
- `dependencyLedger` tracking belongs to the actor process
- contracts do not switch mid-episode
- frame projection is supported, but inference-readiness is still open

## Short Framing

The actor process is the live behavior that sits between routed stimulus and bounded episode state.

It drains the actor mailbox, updates or opens episodes, projects frames when appropriate, emits stimulus, tracks dependencies, and only adopts new contracts after the current episode closes.
