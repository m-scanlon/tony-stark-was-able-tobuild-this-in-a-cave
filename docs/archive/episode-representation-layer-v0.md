# Episode Representation Layer v0

Archived as superseded.

This document described a separate episode-side representation object for recall.

That role is now handled more directly by:

- the current episode contract
- the current recall contract
- bounded recall commands driven by heavy inference calls

The remaining useful ideas here are:

- optional semantic sidecars
- optional runtime-artifact references inside a recall-oriented query layer

Those ideas may return later, but this document should not be treated as active canon for `v1`.

This layer is separate from the episode itself.

It is a derived, mutable representation of what is active in the current episode.

Its job is to give recall a usable query layer without turning the episode into the retrieval object directly.

## Purpose

The episode representation layer is:

- derived from the current episode
- mutable while the episode is open
- rebuildable from episode data
- optimized for recall

It is not:

- the authoritative episode record
- retained memory
- long-term structure

## Core Idea

The current episode may contain a lot of historical activity.

Recall does not need the full historical episode object directly.

It needs a live representation of:

- what entities are active
- what relationships are active
- how strongly they are active
- what runtime artifacts are currently in frame

This representation can also include a semantic sidecar to help with fuzzy retrieval, but structure remains primary.

## Contract

```ts
type EpisodeRepresentationLayer = {
  episode_id: string

  structural: {
    entities: EpisodeEntityState[]
    relationships: EpisodeRelationshipState[]
  }

  semantic?: {
    frame_embedding?: number[]
    entity_embeddings?: Record<string, number[]>
    relationship_embeddings?: Record<string, number[]>
  }

  runtime_context: {
    runtime_artifact_ids: string[]
  }

  updated_at: string
}
```

```ts
type EpisodeEntityState = {
  entity_id: string
  score: number
  mention_count: number
  first_seen_turn: string
  last_seen_turn: string
}
```

```ts
type EpisodeRelationshipState = {
  relationship_id: string
  from_entity_id: string
  to_entity_id: string
  score: number
  mention_count: number
  first_seen_turn: string
  last_seen_turn: string
}
```

## Structural View

The structural view is the primary recall driver.

It keeps explicit entity and relationship references plus their current activation state within the episode.

This preserves:

- exact structural grounding
- inspectability
- compatibility with `anchor_set` lookup
- bounded scoring during recall

## Semantic View

The semantic view is optional.

If present, it acts as a helper surface for:

- fuzzy candidate expansion
- paraphrase support
- semantic widening when exact structural overlap is weak

The semantic view should not replace the structural view.

The structural view remains the main recall address.

## Runtime Context

The representation layer can reference runtime artifacts that are currently active inside the episode.

These are transient artifacts produced during runtime cognition.

They are not retained memory.

The purpose of including them here is to let recall and cognition see what transient outputs are already shaping the frame.

## Usage

The expected recall flow is:

1. Read the episode representation layer.
2. Use the structural view as the primary cue source.
3. Fetch candidate retained artifacts through shared `anchor_set` overlap.
4. Optionally widen candidates through the semantic view.
5. Return the strongest retained artifacts into the current frame.

## Design Principle

This layer is a query layer, not a source of truth.

It exists to represent the current episode in a form that recall can use efficiently.
