# Node Substrate Interface v0

## Core Framing

The node substrate is the underlying runtime execution surface of a node.

It is not:

- the node contract
- the episode
- the frame

It is the runtime layer that lets a node:

- accept events
- update episode state
- project frames
- dispatch commands
- write command results back
- receive published contract updates from Stark

This is the minimal interface the node process should sit on top of.

## Design Goal

The substrate should stay generic enough to support:

- one-shot execution
- bounded multi-step execution
- contract-bounded loop schemas
- later command sets beyond primitives

It should not hardcode one reasoning loop such as:

- `ReAct`
- `OODA`

Those should be supported by the substrate, not baked into it.

## Minimal Interface

```ts
type NodeSubstrate = {
  node_id: string
  contract: NodeContractSnapshot
  active_episode_id?: string

  ingest_event(event: NodeEvent): NodeUpdateResult
  project_frame(): Frame | null
  dispatch_command(invocation: CommandInvocation): CommandDispatchResult
  write_command_result(result: CommandResultEvent): NodeUpdateResult
  receive_published_contract(contract: NodeContractSnapshot): ContractPublicationResult
}
```

This is the current minimal public surface.

## Why These Methods Exist

### `ingest_event(...)`

Accepts incoming runtime events such as:

- stimulus
- command result writeback
- contract publication

This is the main event intake boundary.

## `project_frame()`

Builds the current inference page from active episode state.

The frame is not the source of truth.

It is a projection from the episode.

## `dispatch_command(...)`

Starts or emits a command invocation into the runtime.

This does not assume that command execution is synchronous.

When inference selects a command, the emitted invocation should already carry the runtime identifiers needed to correlate completion and route the result back to the correct node and episode.

## `write_command_result(...)`

Writes the result of a previously dispatched command back into episode state.

This is separate from dispatch because command completion may be:

- asynchronous
- delayed
- external
- capability-bound

The node should not pretend dispatch and completion are the same event.

Primitive-specific execution should return typed result data.

The kernel should then use one shared result-routing/writeback path to format that result into a `CommandResultEvent` and route it back by the stamped command identifiers.

## `receive_published_contract(...)`

Receives a newly published contract snapshot for the node.

In practice, this is the node-side receipt path for a contract publication originated by Stark.

The node should then:

- record that contract in pending node state
- keep the current active contract in force while the current episode is still open
- let the pending contract take effect only after the current episode closes in `v1`

The point is:

- publication may happen now
- receipt may happen now
- adoption may not happen mid-episode in `v1`

## Command Shape

The substrate should assume the runtime command surface is:

```ts
type CommandInvocation = {
  command_id: string
  node_id: string
  episode_id: string
  intent_id?: string
  command_set: string
  command: string
  args: Record<string, unknown>
  emitted_at: string
}
```

The important point is:

- the node chooses the command
- the emitted invocation carries the correlation and routing identifiers needed for completion
- the primitive does its own work
- the shared kernel result-routing path handles completion formatting and return

Conceptually, command completion should come back as:

```ts
type CommandResultEvent = {
  command_id: string
  node_id: string
  episode_id: string
  intent_id?: string
  command_set: string
  command: string
  result_kind: string
  result: unknown
  completed_at: string
}
```

This keeps the system flexible enough to support:

- `primitive` as one command set
- later loop or other command sets

without flattening everything into one global primitive menu.

## Event Shape

The substrate should assume a typed event model.

Conceptually:

```ts
type NodeEvent =
  | StimulusEvent
  | CommandResultEvent
  | ContractPublicationEvent
```

This does not force a final event schema yet.

It only establishes that the substrate is event-driven rather than step-machine-driven.

## Event Heap / Mailbox Model

Because command dispatch and command completion are separate, the runtime still needs a place for deferred writeback and contract publication to land.

`v1` should use one unified typed event intake surface at the kernel front rather than separate global queues.

The existing max heap is the current best fit for that role.

At minimum, that heap should be able to carry:

- incoming stimulus/events
- deferred command result events
- published contract events

Routing should remain thin.

At minimum, the router should:

- accept a typed package or event
- resolve the target node by `node_id` through the live node registry
- place the routed package into that node's mailbox

It should not need primitive-specific routing logic.

Once an event has been routed to a node, it should land in a lightweight node-local mailbox.

That mailbox does not need to be another scheduler.

It is just the pending holding area for already-routed events waiting on that node.

Inference may know that commands are outstanding.

Inference should not need to manage completion timing directly.

That timing belongs to the runtime beneath the node.

## Relationship To Other Objects

- the node contract bounds what the node may do
- the node substrate provides the runtime execution surface
- the episode holds bounded runtime state
- the frame is projected from episode state

So:

- contract = boundary
- substrate = runtime surface
- episode = state container
- frame = inference projection

## Current Design Posture

The strongest current claims are:

- the node substrate should expose a small explicit interface
- event intake and command writeback should remain separate
- the runtime should be event-driven
- the substrate should support multiple loop styles without hardcoding one
- command execution should already assume the command-set-based command surface

## Short Framing

The node substrate is the minimal runtime interface beneath the node process.

Its current core surface is:

- `ingest_event(...)`
- `project_frame()`
- `dispatch_command(...)`
- `write_command_result(...)`
- `receive_published_contract(...)`

That surface should be generic enough to support different contract-bounded execution loops later.
