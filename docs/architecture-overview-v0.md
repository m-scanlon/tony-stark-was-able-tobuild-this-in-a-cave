# Architecture Overview v0

## Purpose

This document is the high-level map of the current architecture.

It is meant to make the active model easy to grasp before dropping into the more detailed contract and runtime docs.

It is not the place for final schema detail.

## Core Framing

The system is a local-first runtime for:

- experiencing
- interacting
- learning

It is not primarily:

- a single global graph
- a prompt transcript
- a collection of fixed service boxes

It is organized around:

- kernel authority
- durable nodes under contract
- bounded episodes
- projected frames
- inference-driven recall
- selective retention

## `v1` Operating Theme

For `v1`, the main practical orientation is:

- experience continuously
- interact when needed
- learn selectively

Most of the architecture exists to support those three concerns.

## Main Runtime Layers

The current working model is:

1. `Structure`
2. `Retention Layer`
3. `Kernel`
4. `Node Contract`
5. `Node`
6. `Episode`
7. `Frame`
8. `Runtime Commands And Runtime Artifacts`

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

That anchor set is the bridge recall uses later.

### `Kernel`

The kernel is the runtime authority.

It is responsible for:

- birthing nodes
- owning the unified typed event intake surface
- owning the global priority heap
- validating and dispatching commands
- routing command results back as typed events
- keeping runtime node registration live

The kernel births `Stark` at startup from a hardcoded contract.

Later nodes are born from published contracts through the same kernel-controlled birth path.

### `Node Contract`

Every node exists under a durable contract.

The active contract boundary is:

- `purpose`
- `capabilities`
- `stimulus`
- `cognition`
- `commands`

The contract defines:

- why the node exists
- what it may rely on
- what can wake it up
- how cognition is bounded
- what commands it may emit

### `Node`

A node is the durable runtime operator acting under a contract.

The node owns runtime machinery such as:

- mailbox handling
- event handling behavior
- recall machinery
- frame assembly behavior
- pending-command tracking
- the pointer to the active episode

The node is durable.

Its episodes are not.

### `Episode`

An episode is the bounded runtime state container for one span of activity.

The current core episode sections are:

- `purpose`
- `interaction_history`
- `recall`
- `available_commands`

The episode is the bounded source of truth for runtime state inside one node's participation.

### `Frame`

The frame is the bounded inference page projected from the episode.

The current frame layout is:

1. `purpose`
2. `interaction`
3. `recall`
4. `available_commands`

The frame is not the durable owner of truth.

It is the current inference projection.

### Runtime Commands And Runtime Artifacts

Runtime execution happens through emitted commands.

The current working command shape is:

```text
skyra <node> <command> -<args> -reason "<why this command is being emitted>"
```

Runtime commands are callable operations inside an active episode.

Runtime artifacts are transient outputs of those commands.

They stay episode-local unless later learning chooses to retain something derived from them.

## Core Boundary Rules

The strongest current boundary rules are:

- contract bounds behavior
- node owns machinery
- episode owns bounded runtime state
- frame is projected from the episode
- kernel remains execution authority
- recall is a read contract
- learning is a write path

This is the current spine of the system.

## Main Role Split

### `Jarvis`

`Jarvis` is the user-facing meaning and attention node.

Its concern is:

- what matters in the current user context
- what deserves attention
- what user-facing interpretation should guide the current episode

### `Stark`

`Stark` is the structural node.

Its concern is:

- node topology
- node contracts
- capability attachment
- structural revision

`Stark` publishes later node contracts, but the kernel still performs node instantiation and execution authority.

### Shared Model

`Jarvis` and `Stark` use the same broad runtime model:

- durable node under contract
- bounded episodes
- projected frames
- retained experience

What differs is role and authority, not ontology.

## Capability Model

Devices and external surfaces are not modeled as nodes by default.

They are modeled through capability contracts.

The split is:

- node contract = governs a runtime operator
- capability contract = governs an external callable surface

This keeps node identity separate from device or API affordance.

## Runtime Flow

The high-level runtime flow is:

1. external input is normalized into a typed event
2. the kernel validates that event and routes it to the target node mailbox
3. the node process checks the event against the active contract
4. the node opens or reuses an episode under the current inactivity/time policy
5. the event is written into episode-local state
6. the node projects a frame from the current episode state
7. heavy inference may emit a bounded recall command
8. the kernel validates and dispatches that recall command
9. recall may bring retained artifacts into scope
10. the node writes those results back into episode state
11. inference selects the next allowed command
12. the node emits that command with the routing identifiers needed for completion
13. the kernel validates and dispatches the command
14. command execution returns typed result data
15. the kernel routes that result back as a typed `command_result` event
16. the node writes the result back into episode state
17. after inactivity, the episode closes
18. learning may be kicked off for that closed episode

This keeps dispatch and completion separate.

It also keeps the node event-driven rather than forcing one hardcoded loop style.

## Episode And History Model

The current scoped episode forms are:

- `node episode`
- `intent episode`

A node episode is the primary bounded local record of one node's participation.

An intent episode is reconstructed across related node episodes linked by shared `intent_id`.

There is no single global mutable history object.

History is reconstructed from:

- episodes
- their order over time
- their scope
- shared identifiers such as `intent_id`

For `v1`, episode closure is still operationally simple:

- inactivity/time closes the episode

This is a practical heuristic, not the final theory of episode boundaries.

## Recall

Recall is the read contract from retained experience into the current episode.

The active `v1` posture is:

1. heavy inference reads current episode context
2. heavy inference decides whether recall is needed
3. heavy inference emits a bounded recall command
4. candidate retained artifacts are fetched through `anchor_set` overlap
5. bounded ranking and admission select the strongest matches
6. the node writes those retained artifacts into `episode.recall`

So for `v1`:

- recall is inference-driven
- retrieval is bounded
- `retained_trace` is part of the recallable retained surface
- no separate episode-side scored field object is required in the active contracts

## Learning

Learning is the write path from completed episodes into retained experience and structure.

The current high-level direction is:

1. preserve bounded factual retained happenings as `retained_trace`
2. derive `retained_understanding`, `retained_salience`, and `retained_tension` where appropriate
3. attach anchors and provenance
4. update the retrieval indexes used later by recall

The current working kickoff shape is:

```text
skyra primitive learn -episode_id <episode_id>
```

Learning is not ordinary in-episode runtime mutation.

It is the later retention step that operates over a closed episode.

## Contract Publication And Runtime Stability

Published contracts do not become active mid-episode in `v1`.

The practical rule is:

- a node may receive a published contract while running
- that contract is held in pending state
- the current episode continues under the currently active contract
- the new contract becomes active only after the current episode closes

This keeps one episode from spanning two contract regimes.

## Interaction Model

Interaction remains unified by default.

That means:

- one chronological interaction history
- typed events inside that history
- no premature split into multiple top-level interaction channels

If a node becomes overloaded, the preferred move is:

- revise the contract
- decompose responsibility
- birth another node if needed

The first move is not to fragment the frame.

## What The Architecture Is Optimizing For

The current architecture is trying to preserve:

- explicit runtime boundaries
- bounded recall and bounded write
- structurally grounded memory
- local-first ownership
- flexible execution without hidden autonomy
- continuity over time without one giant mutable history object

## What Is Still Open

The main unresolved edges are now more behavioral than ontological.

They include:

- exact cognition budgeting and stop rules
- exact command vocabulary and argument grammar
- final inference-readiness / frame projection timing
- exact artifact lifecycle and merge policy
- richer episode-boundary logic beyond inactivity
- richer recall-command generation and ranking policy

## Short Framing

The system is a local-first runtime organized around durable nodes under contract, bounded episodes, projected frames, inference-driven recall, and selective learning.

The kernel is the execution authority.

Nodes own runtime machinery.

Episodes own bounded state.

Recall reads retained experience back into the current episode.

Learning writes selected consequences of completed episodes back into retention and structure.

## See Also

- [data-model-prelim.md](./data-model-prelim.md)
- [node-contract-v0.md](./node-contract-v0.md)
- [node-process-v0.md](./node-process-v0.md)
- [recall-v0.md](./recall-v0.md)
- [consolidation-mechanism-v1.md](./consolidation-mechanism-v1.md)
- [structural-projection-service-v0.md](./structural-projection-service-v0.md)
