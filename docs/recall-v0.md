# Recall v0

## Core Framing

Recall is the read contract from retained experience into the current episode.

That retained experience family includes `retained_trace` as well as the more derived retained artifact types.

For `v1`, recall is driven by heavy inference calls.

It does not depend on a separate episode-side scored field object.

## `v1` Posture

The `v1` recall posture is:

- start from current episode context
- let heavy inference decide whether recall is needed
- have heavy inference emit a bounded recall command
- query retained artifacts by `anchor_set` overlap
- admit the top retained artifacts into `episode.recall`
- project the admitted retained artifacts into frame

So in `v1`:

- recall is contract-driven
- heavy inference chooses the query
- retrieval is bounded
- `retained_trace` is part of the recallable retained surface

## Why This Is The Right `v1` Move

This avoids overbuilding too early.

It gives the system:

- an explicit retained read path
- explicit entities and relationships in recall queries
- retained-artifact lookup through shared anchors
- a bounded implementation surface

Without requiring:

- a scored episode-side structural field
- activation propagation inside the episode
- long-horizon field stabilization rules

## Inputs

The current recall inputs should be:

- the current episode context
- the retention layer

More specifically:

- `purpose`
- current stimulus
- recent interaction history
- already active recalled artifacts
- available commands
- relevant runtime artifacts when present

Heavy inference reads that bounded context and decides whether to emit recall.

## Working Recall Command Shape

If recall is expressed as an emitted runtime command, a good current working shape is:

```text
skyra <node> recall \
  -entity <entity_id> \
  -relationship <relationship_id> \
  -bundle <left_entity_id>:<relationship_id>:<right_entity_id> \
  -top_k <n> \
  -reason "<why recall is needed now>"
```

This matches the current recall posture:

- `-entity` supports broad candidate generation
- `-relationship` supports more specific structural lookup
- `-bundle` supports the strongest bound relational cue when the full triple is known
- `-top_k` keeps recall admission bounded

Example:

```text
skyra jarvis recall \
  -entity assistant \
  -entity terraform \
  -relationship help_with \
  -relationship has_property \
  -reason "assistant and terraform were activated in the current stimulus"
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

The recall command supplies the active ids.

Those ids retrieve a bounded candidate set.

This is the first-stage filter.

## Candidate Ranking

Candidate retained artifacts should then be ranked by structural overlap with the bounded query emitted by heavy inference.

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

The important thing is the shape of the behavior, not precise scoring sophistication yet.

## Base Case / Stop Rule

For `v1`, recall should stay one-pass and bounded.

That means:

- if heavy inference does not emit recall, recall does not run
- if the emitted query yields no usable ids, recall returns an empty result
- recall does one candidate-generation pass through the anchor indexes
- recall does one scoring pass against the emitted query
- recall admits only the top bounded matches that clear the current minimum score threshold
- recall then stops

So the base case is not "walk until there are no more matching ids."

The base case is:

- one query
- one bounded retrieval pass
- one bounded ranking pass

If multi-hop recall exists later, it should arrive with explicit depth, threshold, or candidate-budget bounds.

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

## Worked Example 1

Episode context:

```text
Can you help me with Terraform? It is a new language and it is difficult.
```

Heavy inference may emit:

- `entity assistant`
- `entity terraform`
- `relationship help_with`
- `relationship has_property`
- `bundle assistant:help_with:terraform`

Recall then:

- retrieves retained artifacts whose `anchor_set` overlaps those ids
- admits the strongest matching traces, understandings, salience, or tension into episode recall

## Worked Example 2

Episode context:

```text
I am trying to build a runtime system and the architecture is still confusing.
```

Heavy inference may emit:

- `entity self`
- `entity runtime_system`
- `entity architecture`
- `relationship build`
- `relationship has_property`

Recall then:

- retrieves retained artifacts linked to runtime systems, architecture, confusion, prior design tension, and related structure
- returns only the top bounded set into the current recall section

## What Is Deliberately Deferred

This `v1` does not yet lock:

- multi-hop recall expansion
- semantic widening beyond structural anchors
- complex recall admission policy
- automatic recall retries without another explicit inference decision

Those are valid later improvements.

They should not block the first real recall path.

## Current Design Posture

The strongest current claims are:

- recall is the read contract from retained experience into the current episode
- heavy inference decides when recall should happen
- heavy inference emits the bounded structural query
- retained artifacts are fetched through `anchor_set` overlap
- `v1` recall stops after one bounded scored retrieval pass
- that retained artifact surface includes `retained_trace`
- the strongest retained artifacts are written into episode recall and projected into frame

## Short Framing

For `v1`, recall starts from current episode context.

Heavy inference decides whether recall is needed, emits a bounded recall command, and the retention layer returns the top matching retained artifacts into the episode.
