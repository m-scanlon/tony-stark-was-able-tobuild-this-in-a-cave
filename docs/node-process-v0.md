# Node Process v0

## Core Framing

The node process is the live runtime behavior of a node after birth.

It is not:

- the node contract
- the node substrate interface
- the episode
- the frame

It is the runtime behavior that:

- receives routed events
- maintains an active episode
- updates bounded runtime state
- projects frames for inference
- dispatches commands
- writes command results back
- applies published contracts at episode boundaries

So:

- contract = boundary
- substrate = public runtime surface
- process = live behavior
- episode = bounded state
- frame = inference projection

## Relationship To Node Birth

Node birth creates a live node under a contract.

The node process begins once that node exists.

A newly born node starts:

- live
- bound by its contract
- without an active episode
- without outstanding command state

The node process remains idle until a valid event arrives.

## Inputs To The Node Process

The current `v1` posture is:

- the kernel owns the global priority heap
- the heap carries typed events
- the kernel routes events to the correct node
- each node owns a lightweight mailbox for already-routed events

So the node process does not own global scheduling.

It owns local event handling once an event has already been routed to it.

At minimum, the relevant event families are:

- `stimulus`
- `command_result`
- `contract_publication`

## Core Responsibilities

The node process is responsible for:

- draining its mailbox
- checking incoming events against the active contract
- opening or reusing an episode
- writing event effects into episode state
- maintaining pending-command state
- projecting frames when the node is inference-ready
- dispatching commands allowed by contract
- writing command results back into episode state
- closing episodes after inactivity
- adopting published contracts only after episode closure

## High-Level Runtime Flow

The current high-level process is:

1. kernel routes a typed event to the node mailbox
2. node process takes the next routed event
3. node process checks whether that event is valid under the active contract
4. if valid, the node opens or reuses an episode
5. the event is written into episode-local state
6. any relevant background node machinery may update episode-local artifacts
7. when the node is inference-ready, it projects a frame
8. inference selects the next allowed command
9. the node dispatches that command
10. command results later return as routed events
11. the node writes those results back into episode state
12. after inactivity, the episode closes
13. any newly published contract becomes active only after episode closure

This is the current process skeleton.

It is intentionally more general than:

- one-shot loops
- `ReAct`
- `OODA`

Those must fit on top of this process rather than replace it.

## Contract Gating

The node process should not treat every routed event as valid.

The active contract remains the gate.

At minimum, the contract currently bounds:

- purpose
- accepted stimulus boundary
- allowed interact boundary
- allowed command namespaces and commands

So the process rule is:

- routed event arrives
- node checks contract validity
- invalid events are rejected, ignored, or deferred by runtime policy
- valid events are allowed to affect the active episode

## Episode Management

The episode is the bounded state container the node process acts on.

`v1` episode policy is time-based.

That means:

- if there is no active episode, open one
- if the active episode has gone inactive long enough to close, open a new one
- otherwise reuse the current active episode

This is intentionally simpler than topic-based or semantic boundary logic.

Episode closure should currently be understood as inactivity-driven rather than conceptually "timed out."

## Event Write Paths

Different routed event types affect different parts of node and episode state.

### `stimulus`

At minimum, a stimulus event should:

- append to `interaction_history`
- become eligible as current interaction context
- trigger any node-owned background machinery that updates episode-local derived artifacts

That may include, for example:

- structural projection updates
- episode-field updates
- recall updates

The node owns the machinery.

The episode owns the resulting bounded state.

### `command_result`

A command result event should:

- match back to a prior outstanding command
- update the pending-command registry
- write relevant result effects into the episode
- surface interaction-relevant outcomes into the interaction log when appropriate

The important split is:

- dispatch and completion are separate
- outstanding tracking belongs to node runtime mechanics
- result effects become part of episode state

### `contract_publication`

A newly published contract should not switch the node mid-episode.

Instead:

- the new contract is received by the running node
- the current episode continues under the currently active contract
- the new contract becomes active only after the current episode closes

This keeps one episode from spanning two different contract regimes.

## Pending Commands

The node process owns the pending-command registry.

This is not a contract concern.

For `v1`, the main practical requirement is:

- the node knows whether a command is currently outstanding

The minimum useful internal distinctions are:

- `outstanding`
- `completed`
- `failed`
- `timed_out`

This can remain an internal runtime mechanism.

It does not need to be surfaced as a top-level contract concern.

## Frame Projection

The frame remains a projection from episode state.

The node process owns frame assembly behavior.

The episode remains the source of truth.

The frame is what inference sees.

The timing question remains open.

The better framing is:

- not "when should the node render a frame"
- but "when has the node experienced enough to justify an inference step"

So the node process should support frame projection, but `v1` does not yet lock a final inference-readiness policy.

## Commands And Loop Flexibility

The node process should assume the namespace-based command model.

That means the process is compatible with:

- `skyra <namespace> <command> <args>`

The contract should explicitly allow:

- namespaces
- commands within those namespaces

This allows the node process to stay generic.

It does not need to hardcode:

- a single primitive-only runtime
- a fixed `ReAct` loop
- a fixed `OODA` loop

Instead:

- the substrate/process provides the generic execution pattern
- the contract bounds the allowed execution envelope
- inference may later choose how to act within that envelope

## Relationship To Other Node Docs

- [node-birth-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/node-birth-v0.md) defines how a node comes into existence
- [node-substrate-interface-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/node-substrate-interface-v0.md) defines the minimal public runtime surface
- [node-and-episode-ownership-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/node-and-episode-ownership-v0.md) defines ownership boundaries
- [node-open-questions-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/node-open-questions-v0.md) tracks remaining unresolved node questions

## Current Design Posture

The strongest current claims are:

- the node process is event-driven
- the kernel owns global intake and routing
- each node owns a lightweight mailbox
- the node process opens or reuses episodes by inactivity/time policy
- command dispatch and command completion are separate
- pending-command tracking belongs to the node process
- contracts do not switch mid-episode
- frame projection is supported, but inference-readiness is still open

## Short Framing

The node process is the live behavior that sits between routed events and bounded episode state.

It drains the node mailbox, updates or opens episodes, projects frames when appropriate, dispatches commands, writes results back, and only adopts new contracts after the current episode closes.
