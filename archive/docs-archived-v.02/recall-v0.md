# Recall v0

## Core Framing

Recall is the read contract from retained experience into the current episode.

That retained experience family includes `retained_trace` as well as the more derived retained artifact types.

For `v1`, recall is driven by heavy inference over current episode context.

It does not depend on a separate episode-side scored field object.

## `v1` Posture

The `v1` recall posture is:

- start from current episode context
- let heavy inference decide whether recall is needed
- have heavy inference emit a bounded recall request stimulus
- query retained artifacts by `anchor_set` overlap
- admit the top retained artifacts into `episode.recall`
- project the admitted retained artifacts into frame

So in `v1`:

- recall is contract-driven
- heavy inference chooses the query
- retrieval is bounded
- `retained_trace` is part of the recallable retained surface

## Inputs

The current recall inputs should be:

- the current episode context
- the retention layer

More specifically:

- `purpose`
- current stimulus
- recent interaction history
- already active recalled artifacts
- relevant commitments or public contract context when needed
- relevant runtime artifacts when present

Heavy inference reads that bounded context and decides whether to emit recall.

## Public Contract Shape

Recall should currently be modeled as one public request stimulus plus one public response envelope.

The response envelope requires:

- `status`
- `reason`

The response payload can then carry the actual recall package.

## Working Recall Request Payload

If recall is expressed as an emitted runtime stimulus, a good current working request payload is:

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

This keeps the primitive simple:

- heavy inference emits a bounded query payload
- runtime normalizes it into a structured request
- candidate generation stays bounded
- ranking stays bounded

## Heavy Inference Role

Heavy inference does not fetch retained artifacts directly.

Its job is to:

- read current episode context
- decide whether recall is needed
- choose the bounded structural query
- choose `top_k`

That keeps retrieval explicit and inspectable.

## Candidate Generation

Recall should not scan the entire retention layer.

It should use indexed candidate generation through shared anchors.

That means:

- `entity_id -> retained_artifact_ids`
- `relationship_id -> retained_artifact_ids`

The emitted recall request provides the active ids.

Those ids retrieve a bounded candidate set.

## Candidate Ranking

Candidate retained artifacts should then be ranked by structural overlap with the bounded query selected by heavy inference.

For `v1`, the useful behavior is:

- entity overlap provides broad recall
- relationship overlap provides specificity
- connected relational overlap should outrank weak entity-only overlap

Conceptually:

```text
artifact_score(a) =
  entity_overlap(a.anchor_set.entity_ids, query.entity_ids)
  + relationship_overlap(a.anchor_set.relationship_ids, query.relationship_ids)
  + connectedness_bonus
```

The exact math can stay simple for `v1`.

## Recall Output

The primitive result of recall should be returned in the public response envelope payload.

Conceptually:

```ts
type RecalledArtifact = {
  artifact_id: string
  kind: "trace" | "understanding" | "salience" | "tension"
  score: number
  matched_entity_ids: string[]
  matched_relationship_ids: string[]
}

type RecallPackage = {
  retained_artifact_ids: string[]
  matches: RecalledArtifact[]
}
```

That bounded mixed set is then written into `episode.recall`.

## Base Case / Stop Rule

For `v1`, recall should stay one-pass and bounded.

That means:

- if heavy inference does not emit recall, recall does not run
- if the emitted query yields no usable ids, recall returns an empty payload
- recall does one candidate-generation pass through the anchor indexes
- recall does one scoring pass against the emitted query
- recall admits only the top bounded matches that clear the current minimum score threshold

## Current Design Posture

The strongest current claims are:

- recall is a contract driven by heavy inference
- recall should now be described as emitted request stimulus plus response envelope
- `retained_trace` is part of the recallable retained surface
- retrieval stays bounded through shared-anchor lookup
- admitted matches are written into `episode.recall`

## Short Framing

Recall is the retained read contract.

Heavy inference chooses a bounded recall request stimulus, retention returns a bounded recall package inside a response envelope, and the admitted artifacts are written into the current episode.
