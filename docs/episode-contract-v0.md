# Episode Contract v0

## Core Framing

An episode is the bounded runtime container for one span of activity.

At runtime, it acts as the main episode-local blackboard.

It is the container that holds the live state from which the frame is projected.

The frame is not the episode itself.

The frame is the inference page assembled from the episode.

## What Belongs In The Episode

The current core episode sections are:

- `purpose`
- `interaction_history`
- `recall`
- `episode_field`
- `available_primitives`

This is the current minimum useful container.

## Important Correction

The entity and relationship projection layer is the `episode_field`.

It is not the `workspace`.

`workspace` remains a possible later layer for transient episode-local cognitive artifacts.

That layer is not required in the active v0 episode contract.

## Contract

```ts
type Episode = {
  episode_id: string
  episode_scope: "node" | "intent"

  node_id?: string
  intent_id?: string
  node_contract_id?: string

  purpose: EpisodePurpose
  interaction_history: InteractionHistory
  recall: EpisodeRecall
  episode_field: EpisodeField
  available_primitives: AvailablePrimitive[]

  opened_at: string
  updated_at: string
  closed_at?: string
}
```

## Purpose

`purpose` is the governing lens active for the episode.

It comes from the node contract.

The episode does not invent purpose for itself.

It carries the active purpose under which the episode is operating.

## Interaction History

`interaction_history` is the factual record of episode-local exchange and external action.

It is append-only.

It is broader than the current frame's interaction slice.

The frame later projects:

- `current_stimulus`
- `recent_interaction_history`

from this larger interaction history.

## Recall

`recall` contains the retained artifacts currently activated into scope for the episode.

It is the retained-experience state currently available to the episode, not the full retained store.

It may contain mixed retained artifact families such as:

- `trace`
- `understanding`
- `salience`
- `tension`

## Episode Field

`episode_field` is the scored entity and relationship layer active inside the episode.

It is:

- the structural projection of what is active right now
- the main query surface for recall
- updated from interaction, recall, and later other episode-local sources

It is not the same thing as interaction history.

It is not the same thing as retained recall.

It is the structural layer that sits behind both.

## Available Primitives

`available_primitives` is the current runtime action menu available to the node during the episode.

The current first-class primitives are:

- `recall`
- `interact`

This menu belongs in the episode because it helps define what the frame may expose to inference at that moment.

## Minimal Supporting Types

```ts
type EpisodePurpose = {
  text: string
}
```

```ts
type InteractionHistory = {
  events: InteractionEvent[]
}
```

```ts
type EpisodeRecall = {
  retained_artifact_ids: string[]
}
```

```ts
type AvailablePrimitive = {
  primitive_id: string
}
```

## Frame Projection

The current frame can be projected from the episode as:

```ts
type Frame = {
  purpose: EpisodePurpose
  interaction: {
    current_stimulus: InteractionEvent | null
    recent_interaction_history: InteractionEvent[]
  }
  recall: EpisodeRecall
  available_primitives: AvailablePrimitive[]
}
```

This means:

- the episode is the container
- the frame is the inference projection

## Current Design Posture

The strongest current claims are:

- the episode is the main bounded runtime container
- it acts as the episode-local blackboard
- purpose, interaction history, recall, episode field, and available primitives are the current core sections
- the entity/relationship projection belongs to `episode_field`, not `workspace`
- the frame is projected from the episode rather than being the episode itself

## Short Framing

An episode is the bounded container that holds the live state for one span of activity.

Its current core sections are purpose, interaction history, recall, episode field, and available primitives.

The frame is then projected from that episode for inference.
