# Retained Artifact Recall (Prelim)

## Core Framing

Recall is the contract that reads retained experience into the current episode.

It is driven by heavy inference over current episode context.

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
- current stimulus
- recent interaction history
- already active recalled artifacts
- relevant commitments or public contract context when needed
- relevant runtime artifacts when present

Heavy inference reads that bounded context and decides whether recall is needed.

## Retained Artifacts

The retention layer contains retained artifacts such as:

- retained trace
- retained understanding
- retained salience
- retained tension

Each retained artifact should carry:

- an `anchor_set`
- artifact payload

## The Recall Contract

At a high level:

1. heavy inference reads current episode context
2. if recall is needed, heavy inference emits a bounded recall request stimulus
3. that request carries entity, relationship, and bundle queries
4. candidate retained artifacts are fetched through `anchor_set` overlap
5. bounded ranking admits the strongest retained artifacts
6. admitted artifacts are written into `episode.recall`

## Working Request Payload

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

type RecallRequest = {
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

The recall request provides the active ids.

Those ids then retrieve a bounded candidate set of artifacts.

## Runtime Result

The output of recall should be a mixed recalled set of retained artifacts returned inside the public response envelope payload.

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

- recall is a contract driven by heavy inference
- retained artifacts remain structurally addressable through shared anchors
- `retained_trace` is part of that retrievable retained family
- recall should be described as request stimulus plus response envelope
- ranking remains bounded and overlap-based

## Short Framing

Recall operates by issuing bounded structural requests against retained artifacts anchored into canonical structure.

Heavy inference chooses the request.

The retention layer returns a bounded recalled set inside the response envelope.
