# Retained Artifact Recall (Prelim)

## Core Framing

Recall is the contract that reads retained experience into the current episode.

It is driven by heavy inference calls over current episode context.

It does not depend on a separate episode-side scored field object.

This creates a clean bridge between:

- the current episode
- the retention layer

without requiring:

- raw text lookup
- a separate abstract theme object
- a scored episode-side structural field

## Episode Context

The current episode context available to recall includes:

- `purpose`
- the current stimulus
- recent interaction history
- already active recalled artifacts
- available commands
- relevant runtime artifacts when present

Heavy inference reads that bounded context and decides whether recall is needed.

## Retained Artifacts

The retention layer contains retained artifacts such as:

- retained trace
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
  kind: "trace" | "understanding" | "salience" | "tension"
  anchor_set: {
    entity_ids: string[]
    relationship_ids: string[]
  }
  payload: unknown
}
```

`retained_trace` remains semantically distinct because it preserves the factual retained happening, but it is still part of the same retrievable retained experience family.

## The Recall Contract

At a high level:

1. heavy inference reads current episode context
2. if recall is needed, heavy inference emits a bounded recall command
3. the recall command carries entity, relationship, and bundle queries
4. candidate retained artifacts are fetched through `anchor_set` overlap
5. bounded ranking admits the strongest retained artifacts
6. admitted artifacts are written into `episode.recall`

## Working Command Shape

```text
skyra <node> recall \
  -entity <entity_id> \
  -relationship <relationship_id> \
  -bundle <left_entity_id>:<relationship_id>:<right_entity_id> \
  -top_k <n> \
  -reason "<why recall is needed now>"
```

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

The active ids now come from inference-chosen recall queries rather than from an episode-side scored field.

## Candidate Generation

Recall should not scan the full retention layer.

Instead, it should use indexed candidate generation.

That means:

- entity id -> artifact ids
- relationship id -> artifact ids

The recall command provides the active ids.

Those ids then retrieve a bounded candidate set of artifacts.

This is the first-stage filter.

## Artifact Match

After candidate generation, each candidate artifact is scored against the bounded query selected by heavy inference.

The point of this score is:

- not exact match
- not raw semantic similarity
- but weighted structural overlap

Conceptually:

```text
artifact_match(a) =
  overlap(a.anchor_set.entity_ids, query.entity_ids)
  + overlap(a.anchor_set.relationship_ids, query.relationship_ids)
  + connectedness_bonus
  - mismatch_penalty
```

The important behavior is:

- entity overlap allows broad recall
- relationship overlap provides specificity
- connected relational matches should outrank weak entity-only matches

## Fault Tolerance

The model should remain useful when the current episode context is incomplete.

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
  kind: "trace" | "understanding" | "salience" | "tension"
  score: number
  matched_entity_ids: string[]
  matched_relationship_ids: string[]
  payload: unknown
}
```

This recalled set then enters the current episode frame for inference and runtime execution.

## Current Design Posture

The strongest current claims are:

- recall is a contract driven by heavy inference calls
- retained artifacts carry structural references into the same canonical layer
- `retained_trace` is part of that retrievable retained artifact family
- recall retrieves candidates through shared ids
- recall ranks those candidates by bounded structural overlap

The exact math for:

- overlap scoring
- connectedness bonus
- mismatch penalties
- admission thresholds

remains open.

## Short Framing

Recall operates by issuing bounded structural queries against retained artifacts anchored into canonical structure.

Heavy inference chooses the query.

The retention layer returns a bounded recalled set.

That set is then written into the current episode.
