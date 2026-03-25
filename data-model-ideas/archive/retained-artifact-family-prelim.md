# Retained Artifact Family (Prelim)

## Core Framing

The retention layer should contain multiple retained artifact types.

These artifact types should share a common family resemblance, but they should not be flattened into one identical object.

The most important distinction to preserve is:

- `retained_trace` is conceptually different from the other retained artifacts

This matters because `retained_trace` preserves a factual retained occurrence, while the other retained artifact types preserve derived consequences of experience.

## The Retention Layer

The current retained artifact family is:

- `retained_trace`
- `retained_understanding`
- `retained_salience`
- `retained_tension`

These are all retained artifacts.

They belong to the same retention layer.

But they do not all mean the same kind of thing.

## Shared Family Shape

Across the retained artifact family, the strongest shared elements are:

- `anchor_set`
- contextual linkage to previously active artifacts
- artifact-specific payload

This means the family can share a loose common shape without forcing identical semantics.

Conceptually:

```ts
type RetainedArtifactFamilyBase = {
  id: string
  kind: "trace" | "understanding" | "salience" | "tension"
  anchor_set: AnchorSet
  context_artifact_ids?: string[]
}
```

At this stage, this is a family resemblance contract, not a final strict schema.

## Anchor Set

Every retained artifact should carry an `anchor_set`.

The `anchor_set` is the shared structural link into canonical structure.

It is composed of references such as:

- entity ids
- relationship ids

This is what makes recall possible across different retained artifact types.

The anchor is not a perfect pointer.

It is the common structural surface through which retained artifacts can be linked, recalled, and compared.

## Context Artifact References

Retained artifacts may also carry references to previously active retained artifacts that influenced their formation.

This preserves an important property of experience:

- happenings create understandings
- later happenings can occur through previously formed understandings
- future retained artifacts can therefore be shaped by prior retained artifacts

The exact name and schema of this field remain open.

For now, the current working idea is:

- `context_artifact_ids`

This is intentionally broad enough to apply across retained artifact types.

## Why Trace Must Stay Distinct

`retained_trace` should not be treated as just another interpretation-shaped artifact.

It is different because it preserves:

- a factual retained happening
- a bounded occurrence over structure
- the thing that later understanding, salience, and tension may be derived from

So while `retained_trace` belongs to the retained artifact family, it should remain conceptually distinct.

This preserves the boundary between:

- what happened
- what it meant
- what mattered
- what remained unresolved

## Retained Trace

The current working view of `retained_trace` is:

- it carries an `anchor_set`
- it may carry `context_artifact_ids`
- it preserves a bounded factual rendering of what happened

Conceptually:

```ts
type RetainedTrace = {
  id: string
  kind: "trace"
  anchor_set: AnchorSet
  context_artifact_ids?: string[]
  happened: string
}
```

The `happened` field is intended to remain:

- factual
- bounded
- natural language
- non-interpretive

## Derived Retained Artifacts

The other retained artifact types are derived forms.

Examples:

- `retained_understanding` preserves interpreted meaning
- `retained_salience` preserves what carries weight or attention
- `retained_tension` preserves unresolved or conflicting significance

These may share anchors with one another and with retained traces.

But they should not be treated as identical to retained traces.

## Current Design Posture

The strongest current claims are:

- the retention layer contains multiple retained artifact types
- those artifact types share a common structural grounding through `anchor_set`
- those artifact types may preserve contextual influence through prior artifact refs
- `retained_trace` must remain conceptually distinct from the derived retained artifact types

The exact final schema for:

- `context_artifact_ids`
- cross-artifact provenance
- trace-to-derived-artifact grounding

remains open.

## Short Framing

Retained artifacts should share a common family shape, but `retained_trace` should remain a distinct retained artifact type.

All retained artifacts are structurally grounded through `anchor_set`.

But `retained_trace` preserves the factual retained happening, while the other retained artifact types preserve derived consequences of experience.
