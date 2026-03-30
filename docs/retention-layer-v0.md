# Retention Layer (v0)

## Core Framing

The retention layer is the layer of retained experience.

It is not the episode itself.

It is also not the same as canonical structure.

The retention layer contains retained artifacts that survive past an episode and can later influence recall, cognition, and future experience.

Retained experience is retrievable.

## Structural Position

The current structural split is:

- `Structure` — canonical entities and relationships
- `Episode` — bounded activity record
- `Retention Layer` — retained artifacts derived from experience

This means:

- episodes are historical and bounded
- retained artifacts are selective retained forms
- structure is the canonical substrate they refer into

## Retained Artifact Family

The current retained artifact family is:

- `retained_trace`
- `retained_understanding`
- `retained_salience`
- `retained_tension`

These are different retained artifact types.

They belong to the same retained experience layer.

They are all part of the retrievable retained experience surface.

## Anchor Set

All retained artifacts should carry an `anchor_set`.

The anchor set is the shared structural surface that links retained artifacts into canonical structure.

Conceptually:

```ts
type AnchorSet = {
  entity_ids: string[]
  relationship_ids: string[]
}
```

The anchor set is:

- structural
- shared across retained artifact types
- the main bridge between retention and recall

It is not a perfect pointer.

It is a fault-tolerant structural link surface.

## Shared Family Contract

All retained artifact types share a common family contract.

Conceptually:

```ts
type RetainedArtifactBase = {
  id: string
  kind: "trace" | "understanding" | "salience" | "tension"
  anchor_set: AnchorSet
  context_artifact_ids?: string[]
}
```

## Context Artifact Ids

`context_artifact_ids` represents which previously active retained artifacts shaped the formation of the current artifact.

This field is intended to stay inside the retained artifact family.

It preserves the fact that:

- prior retained artifacts can shape future experience
- later retained artifacts can therefore be formed through the lens of earlier ones

This field should mean:

- influencing retained artifacts

It should not mean:

- arbitrary related artifacts
- all nearby artifacts

## Retained Trace

`retained_trace` is a distinct retained artifact type.

It preserves a bounded factual retained happening.

Conceptually:

```ts
type RetainedTrace = RetainedArtifactBase & {
  kind: "trace"
  happened: string
  source_episode_ids: string[]
}
```

The `happened` field is intended to remain:

- factual
- bounded
- natural language
- non-interpretive

## Retained Understanding

`retained_understanding` preserves interpreted meaning derived from experience.

Conceptually:

```ts
type RetainedUnderstanding = RetainedArtifactBase & {
  kind: "understanding"
  interpretation: string
  source_trace_ids: string[]
}
```

## Retained Salience

`retained_salience` preserves what carries weight, attention, or importance in later cognition.

Conceptually:

```ts
type RetainedSalience = RetainedArtifactBase & {
  kind: "salience"
  signal: string
  source_trace_ids: string[]
}
```

## Retained Tension

`retained_tension` preserves what remains unresolved, conflicting, or incomplete.

Conceptually:

```ts
type RetainedTension = RetainedArtifactBase & {
  kind: "tension"
  unresolved: string
  source_trace_ids: string[]
}
```

## Why Trace Stays Distinct

`retained_trace` should remain distinct from the other retained artifact types.

This preserves the boundary between:

- what happened
- what it meant
- what mattered
- what remained unresolved

The other retained artifact types may be derived from traces, but they should not collapse into the same semantic role.

## Current Design Posture

The strongest current claims are:

- the retention layer contains multiple retained artifact types
- all retained artifacts share a structural grounding through `anchor_set`
- retained artifacts may preserve shaping influence through `context_artifact_ids`
- retained experience is retrievable
- all retained artifact types belong to that retrievable retained experience layer
- `retained_trace` remains a distinct retained artifact type
- derived retained artifacts may carry `source_trace_ids` for factual grounding

This document defines the current contract surface only.

It does not define:

- extraction logic
- consolidation steps
- scoring rules
- retrieval ranking

## Short Framing

The retention layer is composed of multiple retained artifact types.

All retained artifacts share an `anchor_set` into canonical structure and may record which prior retained artifacts shaped them.

They all belong to the retrievable retained experience layer.

`retained_trace` preserves the factual retained happening.

`retained_understanding`, `retained_salience`, and `retained_tension` preserve derived consequences of experience.
