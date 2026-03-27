# Retained Artifact Recall (Prelim)

## Core Framing

Recall should be driven by the current episode's scored entity/relationship layer.

Retained artifacts should be retrieved by matching against that scored layer.

This creates a clean bridge between:

- the current episode
- the retention layer

without requiring:

- raw text lookup
- a separate abstract theme object
- a giant stored pattern layer

## Episode Structural Layer

Within the current episode, there should be a structural layer composed of:

- entities resolved in the episode
- relationships resolved in the episode

This layer accumulates across the episode and carries activation scores.

It is the current structural representation of what the episode is about.

It sits just behind the immediate incoming turn and provides continuity across recall steps.

## Retained Artifacts

The retention layer contains retained artifacts such as:

- retained understanding
- retained salience
- retained tension

These are different object types.

What they share is that they are structurally addressable through canonical references.

Each retained artifact should carry:

- an `anchor_set`
- artifact payload

Conceptually:

```ts
type RetainedArtifact = {
  id: string
  kind: "understanding" | "salience" | "tension"
  anchor_set: {
    entity_ids: string[]
    relationship_ids: string[]
  }
  payload: unknown
}
```

## The Recall Bridge

The tie between the episode and the retained artifacts is:

- shared entity ids
- shared relationship ids

The current episode field and retained artifacts both point into the same canonical structure through the anchor set.

That means recall can operate over structure directly rather than over text.

## Recall Flow

At a high level:

1. resolve entities and relationships from the current stimulus
2. update the scored entity/relationship layer of the episode
3. take the strongest entities and relationships from that scored layer
4. use those ids to pull candidate retained artifacts
5. score candidate artifacts by weighted structural overlap with the current episode field
6. admit the top artifacts into recall

## Candidate Generation

Recall should not scan the full retention layer.

Instead, it should use indexed candidate generation.

That means:

- entity id -> artifact ids
- relationship id -> artifact ids

The scored episode field provides the active ids.

Those ids then retrieve a bounded candidate set of artifacts.

This is the first-stage filter.

## Artifact Match

After candidate generation, each candidate artifact is scored against the current episode field.

The point of this score is:

- not exact match
- not raw semantic similarity
- but weighted structural overlap

Conceptually:

```text
artifact_match(a) =
  overlap(a.anchor_set.entity_ids, scored_episode_entities)
  + overlap(a.anchor_set.relationship_ids, scored_episode_relationships)
  + connectedness_bonus
  - mismatch_penalty
```

The important behavior is:

- entity overlap allows broad recall
- relationship overlap provides specificity
- connected relational matches should outrank weak entity-only matches

## Why Relationships Matter

An entity alone is often too weak to retrieve the right retained artifact.

Example:

- `outside`

may retrieve something generic.

But:

- `self -> located_in -> outside`
- `self -> doing -> construction`

gives a much more specific current structure.

That more specific structure should pull more specific retained artifacts.

So the recall value is not a single entity.

It is the scored relational slice currently active in the episode.

## Fault Tolerance

The model should remain useful when the current episode structure is incomplete.

So:

- partial overlap should still retrieve candidates
- exact structural identity should not be required
- connected relational matches should still dominate weak partial matches

Fault tolerance should come from overlap scoring and bounded candidate generation.

It should not depend on hand-written special cases.

## Runtime Result

The output of recall should be a mixed recalled set of retained artifacts.

Conceptually:

```ts
type RecalledArtifact = {
  artifact_id: string
  kind: "understanding" | "salience" | "tension"
  score: number
  matched_entity_ids: string[]
  matched_relationship_ids: string[]
  payload: unknown
}
```

This recalled set then enters the current episode frame for inference and runtime execution.

## Current Design Posture

The strongest current claims are:

- the current episode should maintain a scored entity/relationship layer
- retained artifacts should carry structural references into the same canonical layer
- recall should retrieve candidates through shared ids
- recall should rank those candidates by weighted structural overlap

The exact math for:

- overlap scoring
- connectedness bonus
- mismatch penalties
- admission thresholds

remains open.

## Short Framing

Recall should operate by matching the scored entity/relationship layer of the current episode against the anchor sets carried by retained artifacts.

The episode provides the active structural field.

The retention layer provides structurally addressable artifacts.

Shared entity and relationship ids form the bridge between them.
