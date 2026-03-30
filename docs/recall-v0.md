# Recall v0

## Core Framing

Recall is the read path from retained experience into the current episode.

That retained experience surface includes `retained_trace` as well as the more derived retained artifact types.

For `v1`, recall should stay simple.

It should be driven primarily by the current stimulus rather than by a fully accumulated and heavily scored episode field.

This keeps the first implementation streamlined while preserving the larger architectural direction.

## `v1` Posture

The `v1` recall posture is:

- start from the current stimulus
- make one light inference call over that stimulus
- extract entities and relationships
- store that result as a thin structural array in the episode field
- query retained artifacts by `anchor_set` overlap
- admit the top retained artifacts into episode recall

So in `v1`:

- `episode_field` still exists
- but it is thin
- and it is mostly the current-stimulus structural projection

Later versions may expand this into a richer accumulated and dynamically scored episode field.

## Why This Is The Right `v1` Move

This avoids overbuilding too early.

It gives the system:

- a real structural recall path
- explicit entities and relationships
- retained-artifact lookup through shared anchors
- a bounded implementation surface

Without immediately requiring:

- heavy activation propagation
- long-horizon field scoring
- complex episode-wide theme stabilization

## Inputs

The current recall inputs should be:

- the current raw stimulus
- the current episode
- the retention layer

More specifically:

- raw stimulus provides the current recall cue
- the episode provides the container where the thin field and recalled artifacts live
- the retention layer provides the candidate retained artifacts

## Thin Episode Field

For `v1`, the episode field should be treated as a thin structural cue surface derived from the current stimulus.

Conceptually:

```ts
type EpisodeField = {
  entities: EpisodeFieldEntity[]
  relationships: EpisodeFieldRelationship[]
  updated_at: string
}
```

```ts
type EpisodeFieldEntity = {
  entity_id: string
  confidence?: number
}
```

```ts
type EpisodeFieldRelationship = {
  relationship_id: string
  from_entity_id: string
  to_entity_id: string
  confidence?: number
}
```

This is intentionally thin.

It is enough to:

- hold the current structural cue set
- support anchor-based recall
- preserve the architectural slot for richer field behavior later

## Recall Flow

The `v1` recall flow should be:

1. read the current raw stimulus
2. make one light inference call over that stimulus
3. extract entities and relationships
4. write those into the current episode field
5. use the resulting entity and relationship ids to retrieve candidate retained artifacts
6. score those candidates by structural overlap
7. admit the strongest retained artifacts into episode recall
8. project the admitted artifacts into the frame's `recall` section

This makes recall stimulus-first while still keeping the episode and retention model intact.

## Working Recall Command Shape

If recall is expressed as an emitted runtime command, a good current working shape is:

```text
skyra primitive recall \
  -entity <entity_id> \
  -relationship <relationship_id> \
  -bundle <left_entity_id>:<relationship_id>:<right_entity_id> \
  -top_k <n>
```

This matches the current recall posture:

- `-entity` supports broad candidate generation
- `-relationship` supports more specific structural lookup
- `-bundle` supports the strongest bound relational cue when the full triple is known
- `-top_k` keeps recall admission bounded

Example:

```text
skyra primitive recall \
  -entity assistant \
  -entity terraform \
  -relationship help_with \
  -relationship has_property \
  -bundle assistant:help_with:terraform \
  -top_k 8
```

This should still be treated as a working command form rather than a final locked argument grammar.

Internally, that command is best normalized into a query array rather than handled as one flat string.

Conceptually:

```ts
type RecallQuery =
  | { kind: "entity"; entity_id: string }
  | { kind: "relationship"; relationship_id: string }
  | {
      kind: "bundle"
      left_entity_id: string
      relationship_id: string
      right_entity_id: string
    }
```

```ts
type RecallArgs = {
  queries: RecallQuery[]
  top_k?: number
}
```

This keeps the primitive simple:

- normalize the emitted command into `RecallQuery[]`
- loop through those queries
- build the candidate set
- rank the candidates
- return one typed recall package

## Light Inference Extraction

The extraction call should produce:

- entities
- relationships
- explicit bindings between them

The point is not to produce a full knowledge graph.

It is just to produce a bounded structural cue array for recall.

Conceptually:

```ts
type StimulusProjection = {
  entities: ProjectedEntityCue[]
  relationships: ProjectedRelationshipCue[]
}
```

```ts
type ProjectedEntityCue = {
  entity_id: string
  confidence?: number
}
```

```ts
type ProjectedRelationshipCue = {
  relationship_id: string
  from_entity_id: string
  to_entity_id: string
  confidence?: number
}
```

These cues then become the current thin episode field.

## Candidate Generation

Recall should not scan the entire retention layer.

It should use indexed candidate generation through shared anchors.

That means:

- `entity_id -> retained_artifact_ids`
- `relationship_id -> retained_artifact_ids`

The current stimulus projection supplies the active ids.

Those ids retrieve a bounded candidate set.

This is the first-stage filter.

## Candidate Ranking

Candidate retained artifacts should then be ranked by structural overlap with the thin episode field.

For `v1`, the useful behavior is:

- entity overlap provides broad recall
- relationship overlap provides specificity
- connected relational overlap should outrank weak entity-only overlap

Conceptually:

```text
artifact_score(a) =
  entity_overlap(a.anchor_set.entity_ids, field.entities)
  + relationship_overlap(a.anchor_set.relationship_ids, field.relationships)
  + connectedness_bonus
```

The exact math can stay simple for `v1`.

The important thing is the shape of the behavior, not precise scoring sophistication yet.

## Base Case / Stop Rule

For `v1`, recall should stay one-pass and bounded.

That means:

- if the current stimulus/frame projection yields no usable entity or relationship ids, recall should return an empty result
- recall should do one candidate-generation pass through the anchor indexes
- recall should do one scoring pass against the current stimulus-derived frame structure
- recall should admit only the top bounded matches that clear the current minimum score threshold
- recall should then stop

So the base case is not "walk until there are no more matching ids."

The base case is:

- one structural retrieval pass
- one scoring pass
- bounded admission

`v1` should not yet walk outward through attached retained artifacts until exhaustion.

If that kind of multi-hop recall exists later, it should arrive with explicit depth, threshold, or candidate-budget bounds.

## Recall Output

The primitive result of recall should be one typed recall package.

Conceptually:

```ts
type RecallPackage = {
  retained_artifact_ids: string[]
  matches: RecalledArtifact[]
}
```

That package should then be returned through the shared kernel result-routing/writeback path.

The node can then write the admitted retained artifact ids into episode recall and use the richer match detail for frame projection or later runtime choice.

The bounded mixed set admitted into the current episode is:

Conceptually:

```ts
type EpisodeRecall = {
  retained_artifact_ids: string[]
}
```

And for projection into frame:

```ts
type RecalledArtifact = {
  artifact_id: string
  kind: "trace" | "understanding" | "salience" | "tension"
  score: number
  matched_entity_ids: string[]
  matched_relationship_ids: string[]
}
```

The frame does not need the entire retained store.

It only needs the bounded set currently brought into scope.

## Relationship To The Larger Architecture

This `v1` choice does not discard the larger episode-field model.

It just keeps the first implementation thin.

So the progression is:

- `v1`: episode field is mostly the current-stimulus structural projection
- later: episode field becomes an accumulated, dynamically scored structural layer across the episode

This preserves the long-term direction without making it a blocker.

## Worked Example 1

Stimulus:

```text
Can you help me with Terraform? It is a new language and it is difficult.
```

Light projection:

- entities:
  - `assistant`
  - `self`
  - `terraform`
  - `language`
- relationships:
  - `assistant -> help -> self`
  - `assistant -> help_with -> terraform`
  - `terraform -> is_a -> language`
  - `terraform -> has_property -> difficult`

Recall then:

- uses `terraform`, `language`, `help_with`, `has_property`
- retrieves retained artifacts whose `anchor_set` overlaps those ids
- admits the strongest matching traces, understandings, salience, or tension into episode recall

## Worked Example 2

Stimulus:

```text
I am trying to build a runtime system and the architecture is still confusing.
```

Light projection:

- entities:
  - `self`
  - `runtime_system`
  - `architecture`
- relationships:
  - `self -> build -> runtime_system`
  - `architecture -> has_property -> confusing`

Recall then:

- uses those entity and relationship ids as the current cue surface
- retrieves retained artifacts linked to runtime systems, architecture, confusion, prior design tension, and related structure
- returns only the top bounded set into the current recall section

## What Is Deliberately Deferred

This `v1` does not yet lock:

- rich multi-turn field accumulation
- activation decay and reinforcement
- propagation across a larger connected structural field
- semantic widening beyond structural anchors
- complex recall admission policy

Those are valid later improvements.

They should not block the first real recall path.

## Current Design Posture

The strongest current claims are:

- `v1` recall should be stimulus-first
- one light inference call should produce the structural cue surface
- the current cue surface should be written into a thin episode field
- retained artifacts should be fetched through `anchor_set` overlap
- `v1` recall should stop after one bounded scored retrieval pass
- that retained artifact surface includes `retained_trace`
- the strongest retained artifacts should be admitted into episode recall and projected into frame

## Short Framing

For `v1`, recall should start from the current stimulus.

Make one light inference call, extract entities and relationships, write them into a thin episode field, and use that structural cue set to pull the top matching retained artifacts into recall.

Then stop after that one bounded scored pass.

That recalled set may include `retained_trace`, even though trace remains semantically distinct from more derived retained artifacts.
