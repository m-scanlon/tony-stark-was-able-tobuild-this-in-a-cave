# Data Model Overview

## Core Framing

The system is not primarily a global graph and it is not primarily a prompt transcript.

It is a runtime system organized around:

- durable actors under contract
- bounded episodes inside those actors
- projected frames for inference
- runtime execution inside the episode
- retained artifacts produced later through learning

Canonical structure still matters, but it is the substrate these other layers refer into rather than the sole primary model.

## V1 Emphasis

For `v1`, the highest-level operating theme is:

- experiencing
- acting
- learning

The current runtime machinery exists mainly to support those three concerns.

See also:

- [v1-operating-theme-prelim.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/v1-operating-theme-prelim.md)

## Canonical Layers

The current working model is:

- `Structure` — canonical entities and relationships
- `Actor Contract` — durable behavioral boundary
- `Actor` — durable runtime operator under contract
- `Episode` — bounded runtime container for one span of activity
- `Frame` — the bounded inference page projected from the current episode
- `Runtime Execution` — emitted stimulus and transient runtime artifacts inside the episode
- `Retention Layer` — retained artifacts that survive the episode
- `Reconstructed History` — derived views across episodes over time

This is the current backbone.

## Applies To All Actors

This model applies across actor types.

That includes:

- user-facing or task-facing actors
- `Jarvis`
- `Stark`

The same layered runtime model applies even when a actor has a different purpose.

## Actor And Contract

Every actor exists under a contract.

At the current contract level, the durable boundary remains:

- `purpose`
- `commitments`
- request stimuli
- response envelopes

These define:

- why the actor exists
- what durable commitments it carries
- what public request stimuli it can receive
- what public response envelopes it can emit

The contract is durable.

The actor is the long-lived runtime operator acting under that contract.

## Episode

An episode is the bounded runtime container for one span of activity.

It is the main episode-local state container from which the frame is projected.

The current core episode sections are:

- `purpose`
- `interaction_history`
- `recall`

The episode is not the same thing as the frame.

The episode is the source of truth for bounded runtime state.

## Frame

The frame is the full page consumed by inference.

It is projected from the current episode.

The current frame layout is:

1. `purpose`
2. `interaction`
3. `recall`

The frame should stay smaller than the episode.

It is the bounded inference page, not the durable runtime container.

The public request/response surface remains part of the active actor contract rather than a separate episode or frame field.

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

## Runtime Execution

Runtime execution occurs inside the active episode.

The current key split is:

- runtime stimulus is the callable in-episode request/response traffic
- runtime artifacts are transient outputs produced while fulfilling that traffic

Runtime execution should assume the actor-first outer protocol:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

The registry holds the published contract the payload must conform to.

The runtime carries concrete payload instances.

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

Recall is the read contract from retained experience into the current episode frame.

At a high level:

1. heavy inference reads current episode context
2. heavy inference emits bounded recall stimulus when recall is needed
3. retained artifacts with overlapping anchors are fetched and ranked
4. a bounded mixed set enters episode recall and may be projected into frame

Recall is therefore driven by:

- current episode context
- heavy inference
- bounded recall request stimulus
- structural overlap between recall queries and retained artifacts

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

- actors are the durable runtime operators
- episodes are the primary bounded unit of runtime state
- frames are projected from episodes for inference
- runtime stimulus and runtime artifacts remain episode-local by default
- retained artifacts belong to the retention layer
- recall is a contract driven by heavy inference
- learning writes from episodes into retained experience and structure
