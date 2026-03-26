# Data Model Overview

## Core Framing

The system is not primarily a global graph.

It is a temporal system organized around bounded episodes, runtime cognition inside those episodes, and a retained experience layer produced later through learning.

Canonical structure still matters, but it is the substrate the other layers refer into rather than the sole primary model.

## Canonical Layers

The current working model is:

- `Structure` — canonical entities and relationships
- `Episode` — bounded activity at node or intent scope
- `Episode Field` — the scored entity/relationship layer active within the current episode
- `Runtime Cognition` — callable runtime primitives and transient runtime artifacts inside the episode
- `Retention Layer` — retained artifacts that survive the episode
- `Reconstructed History` — derived views across episodes over time

This is the main layered backbone.

## Applies To All Nodes

This model applies to any node in the system.

That includes:

- user-facing or task-facing nodes
- `Jarvis` (user-facing meaning node)
- `Stark` (structural node)

It is not outside the node model.

## Cycles And Episodes

A cycle is the atomic unit of execution.

The working cycle shape remains:

```text
stimulus -> recall -> cognition -> interact
```

A node episode is a bounded grouping of one or more cycles of one node's participation.

An intent episode is a higher-level grouping of related node episodes linked by shared `intent_id`.

There is no single global episode object.

## Episode Frame

Every episode organizes its active frame into:

- `interaction`
- `recall`
- `cognition`

Interaction:

- incoming stimulus
- outgoing interact
- external actions
- timestamps

Recall:

- retained artifacts activated into scope from retained experience
- selected rather than exhaustive
- bounded by episode context

Cognition:

- in-episode reasoning and decision formation
- runtime primitive execution
- transient runtime artifact production

## Episode Field

Each active episode should also maintain an episode field.

The episode field is:

- the scored entity/relationship layer of the current episode
- the structural representation of what the episode is about right now
- the main scoring surface that recall uses

It is not a separate abstract theme object.

It is the accumulated scored structure of the episode itself.

## Runtime Cognition

Runtime cognition occurs inside the active episode.

The current key distinction is:

- runtime primitives are callable in-episode operations
- runtime artifacts are transient outputs of those operations

Runtime artifacts are not retained by default.

They remain episode-local unless later learning selects them into retained experience.

## Retained Experience

The retention layer is the layer of retained experience.

It is not the episode itself.

It is composed of retained artifacts, currently:

- `retained_trace`
- `retained_understanding`
- `retained_salience`
- `retained_tension`

All retained artifacts share an `anchor_set` into canonical structure.

`retained_trace` remains distinct from the derived retained artifact types.

This preserves the boundary between:

- what happened
- what it meant
- what mattered
- what remained unresolved

## Recall

Recall is the read path from retained experience into the active episode frame.

At a high level:

1. the current stimulus updates the episode field
2. the episode field scores entities and relationships
3. the dominant connected slice of that field becomes the active recall surface
4. retained artifacts with overlapping anchors are fetched and ranked
5. a bounded mixed set enters recall

Recall is therefore driven by:

- the incoming stimulus
- the accumulated structure of the current episode
- structural overlap between the episode field and retained artifacts

## Learning

Learning is the write path from episodes into retained experience and structure.

Learning is not ordinary runtime cognition.

It is the later process that decides what from an episode should survive as:

- retained traces
- derived retained artifacts
- structure updates

Runtime artifacts may inform learning, but they do not become retained experience automatically.

## Structure

Structure remains canonical.

It contains:

- entities
- relationships

Retained artifacts do not replace structure.

They refer into it through `anchor_set`.

The episode field also scores over that same structural substrate.

## History

History is not stored as one mutable object.

It is reconstructed from:

- episodes
- their cycles
- their ordering over time
- shared `intent_id`

Different scopes may reconstruct different histories from the same underlying episodic records.

## Node Contract

Every node exists under a contract.

At the contract level, the core primitives remain:

- `purpose`
- `stimulus`
- `interact`

These define the node's boundary and eligibility to act.

They are not the same thing as runtime callable primitives inside an episode.

## Current Design Posture

The strongest current claims are:

- episodes are the primary bounded unit of activity
- the episode field is the scored structural layer active within an episode
- runtime primitives and runtime artifacts belong to in-episode cognition
- retained artifacts belong to retained experience
- recall reads from retained experience through the scored episode field
- learning writes from episodes into retained experience and structure
- the same model applies across node types, including `Jarvis` and `Stark`

## Short Framing

The current data model is a layered episode-based architecture.

Episodes hold active work.

Episode fields score the current structural context.

Runtime cognition operates inside episodes.

Learning turns selected episode outcomes into retained artifacts.

Recall brings retained artifacts back into later episodes through shared structural anchors.
