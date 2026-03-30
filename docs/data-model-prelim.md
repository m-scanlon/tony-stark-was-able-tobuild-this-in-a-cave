# Data Model Overview

## Core Framing

The system is not primarily a global graph and it is not primarily a prompt transcript.

It is a runtime system organized around:

- durable nodes under contract
- bounded episodes inside those nodes
- a scored structural layer active within each episode
- projected frames for inference
- retained artifacts produced later through learning

Canonical structure still matters, but it is the substrate these other layers refer into rather than the sole primary model.

## V1 Emphasis

For `v1`, the highest-level operating theme is:

- experiencing
- interacting
- learning

The current runtime machinery exists mainly to support those three concerns.

See also:

- [v1-operating-theme-prelim.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/v1-operating-theme-prelim.md)

## Canonical Layers

The current working model is:

- `Structure` — canonical entities and relationships
- `Node` — durable runtime operator under contract
- `Episode` — bounded runtime container for one span of activity
- `Episode Field` — the scored entity/relationship layer active within the current episode
- `Frame` — the bounded inference page projected from the current episode
- `Runtime Execution` — emitted commands and transient runtime artifacts inside the episode
- `Retention Layer` — retained artifacts that survive the episode
- `Reconstructed History` — derived views across episodes over time

This is the current backbone.

## Applies To All Nodes

This model applies across node types.

That includes:

- user-facing or task-facing nodes
- `Jarvis`
- `Stark`

The same layered runtime model applies even when a node has a different purpose.

## Node And Contract

Every node exists under a contract.

At the current contract level, the durable boundary remains:

- `purpose`
- `capabilities`
- `stimulus`
- `cognition`
- `commands`

These define:

- why the node exists
- what capability surfaces it may rely on
- what may wake it up
- how cognition is bounded
- what commands it may emit

The contract is durable.

The node is the long-lived runtime operator acting under that contract.

## Episode

An episode is the bounded runtime container for one span of activity.

It is the main episode-local state container from which the frame is projected.

The current core episode sections are:

- `purpose`
- `interaction_history`
- `recall`
- `episode_field`
- `available_commands`

The episode is not the same thing as the frame.

The episode is the source of truth for bounded runtime state.

## Frame

The frame is the full page consumed by inference.

It is projected from the current episode.

The current frame layout is:

1. `purpose`
2. `interaction`
3. `recall`
4. `available_commands`

The frame should stay smaller than the episode.

It is the bounded inference page, not the durable runtime container.

## Interaction

Interaction should remain unified and chronological by default.

That means:

- one chronological interaction log inside the episode
- typed events inside that log
- no premature splitting into many separate frame channels

The frame later projects:

- `current_stimulus`
- `recent_interaction_history`

from the larger interaction history.

## Episode Field

Each episode maintains an `episode_field`.

The episode field is:

- the scored entity/relationship layer active within the current episode
- the structural representation of what the episode is about right now
- the main cue surface that recall uses

It is not:

- interaction history
- retained recall itself
- a separate abstract theme object

It is the dynamically updated structural layer sitting behind the current frame.

## Runtime Execution

Runtime execution occurs inside the active episode.

The current key split is:

- runtime commands are callable in-episode operations
- runtime artifacts are transient outputs of those operations

Runtime execution should assume the command-set-based command surface:

```text
skyra <command_set> <command> -<args>
```

This keeps runtime execution more flexible than a flat primitive-only model.

Runtime artifacts remain episode-local unless later learning selects something from them into retained experience.

## Retention Layer

The retention layer is the layer of retained experience.

It is not the episode and it is not canonical structure.

The current retained artifact family is:

- `retained_trace`
- `retained_understanding`
- `retained_salience`
- `retained_tension`

All retained artifacts share an `anchor_set` into canonical structure.

That anchor set is the main bridge between retained experience and recall.

## Recall

Recall is the read path from retained experience into the current episode frame.

At a high level:

1. current episode activity updates the `episode_field`
2. the episode field scores entities and relationships
3. the dominant connected slice of that field becomes the active cue surface
4. retained artifacts with overlapping anchors are fetched and ranked
5. a bounded mixed set enters episode recall and may be projected into frame

Recall is therefore driven by:

- current interaction/stimulus
- accumulated episode structure
- structural overlap between the episode field and retained artifacts

## Learning

Learning is the write path from episodes into retained experience and structure.

Learning is not ordinary runtime execution.

It is the later process that decides what from an episode should survive as:

- retained traces
- derived retained artifacts
- structure updates

Runtime artifacts may inform learning, but they are not retained by default.

## Structure

Structure remains canonical.

It contains:

- entities
- relationships

Retained artifacts do not replace structure.

Episodes do not replace structure.

Both refer into it.

## History

History is not one mutable object.

It is reconstructed from:

- episodes
- their ordering over time
- shared identifiers such as `intent_id`
- retained artifacts that survive across episodes

Different scopes may reconstruct different historical views from the same underlying records.

## Current Design Posture

The strongest current claims are:

- nodes are the durable runtime operators
- episodes are the primary bounded unit of runtime state
- frames are projected from episodes for inference
- the episode field is the scored structural layer active within an episode
- runtime commands and runtime artifacts remain episode-local by default
- retained artifacts belong to the retention layer
- recall reads from retained experience through the episode field
- learning writes from episodes into retained experience and structure

## Short Framing

The current data model is a layered runtime architecture.

Nodes operate under durable contracts.

Episodes hold bounded runtime state.

Frames are projected from episodes for inference.

Episode fields score the current structural context.

Runtime execution stays local to the episode unless later learning turns selected outcomes into retained artifacts.
