# Episode Contract v0

## Core Framing

An episode is the bounded runtime container for one span of activity.

At runtime, it acts as the main episode-local state container.

It is the container that holds the live state from which the frame is projected.

The frame is not the episode itself.

The frame is the inference page assembled from the episode.

## What Belongs In The Episode

The current core episode sections are:

- `purpose`
- `interaction_history`
- `recall`

This is the current minimum useful container.

## Important Clarification

There is no canonical separate scored field object in this version.

Recall is the bounded retained state currently admitted into the episode through prior recall requests.

`workspace` remains a possible later layer for transient episode-local cognitive artifacts.

That layer is not required in the active `v0` episode contract.

## Contract

```ts
type Episode = {
  episode_id: string
  episode_scope: "actor" | "intent"

  actor_id?: string
  intent_id?: string
  actor_contract_id?: string

  purpose: EpisodePurpose
  interaction_history: InteractionHistory
  recall: EpisodeRecall

  opened_at: string
  updated_at: string
  closed_at?: string
}
```

## Purpose

`purpose` is the governing lens active for the episode.

It comes from the actor contract.

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

In `v1`, this section is written by bounded recall results chosen through heavy inference.

## Contract Boundary

The episode contract does not carry a separate command-allowance field.

The public request/response surface lives on the active actor contract.

That contract boundary still governs what the actor may receive and emit during the episode.

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
}
```

This means:

- the episode is the container
- the frame is the inference projection

## Current Design Posture

The strongest current claims are:

- the episode is the main bounded runtime container
- it acts as the episode-local state container
- purpose, interaction history, and recall are the current core sections
- recall is a contract driven by heavy inference rather than by an episode-side scored field
- the frame is projected from the episode rather than being the episode itself

## Short Framing

An episode is the bounded container that holds the live state for one span of activity.

Its current core sections are purpose, interaction_history, and recall.

The frame is then projected from that episode for inference.
