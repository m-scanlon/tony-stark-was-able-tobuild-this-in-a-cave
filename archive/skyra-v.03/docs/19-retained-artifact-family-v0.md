# Retained Artifact Family v0

## Purpose

This document locks the current retained artifact family for `v.03`.

It defines the retained artifact schema and the meaning of
`trust_at_formation`.

## Family Shape

```ts
type AnchorSet = {
  being_names: string[]              // Names of the beings this retained artifact is structurally anchored to.
  relationship_pairs: [string, string][] // Relationships as unordered name pairs [being_a, being_b].
}

type RetainedArtifactKind =
  | "trace"         // A retained factual happening.
  | "understanding" // A retained interpretation or meaning.
  | "salience"      // A retained sense of importance or weight.
  | "tension"       // A retained unresolved conflict, question, or incompletion.

type RetainedArtifactBase = {
  kind: RetainedArtifactKind      // Which member of the retained family this artifact is.
  anchor_set: AnchorSet           // The structural anchors used to ground and later retrieve it.
  context_artifacts?: AnchorSet[] // Earlier retained artifacts that shaped formation of this one, referenced by their anchor sets.
  trust_at_formation: number      // The cognitive trust judgment present at the moment this artifact was formed.
}

type RetainedTrace = RetainedArtifactBase & {
  kind: "trace"    // This artifact is a trace.
  happened: string // The bounded factual happening the being retains.
}

type RetainedUnderstanding = RetainedArtifactBase & {
  kind: "understanding"  // This artifact is an understanding.
  interpretation: string // The meaning, lesson, or interpretation drawn from experience.
}

type RetainedSalience = RetainedArtifactBase & {
  kind: "salience" // This artifact is salience.
  signal: string   // What carries weight, importance, or attention later.
}

type RetainedTension = RetainedArtifactBase & {
  kind: "tension"    // This artifact is tension.
  unresolved: string // What remains unsettled, conflicting, incomplete, or open.
}

type RetainedArtifact =
  | RetainedTrace
  | RetainedUnderstanding
  | RetainedSalience
  | RetainedTension
  // A retained artifact is always exactly one of the four family members.
```

## The Four Family Members

- `trace` preserves what happened.
- `understanding` preserves what it meant.
- `salience` preserves what carried weight.
- `tension` preserves what remained unresolved.

## Meaning Of trust_at_formation

`trust_at_formation` is the permanent snapshot of the forming being's cognitive
trust judgment at the moment the retained artifact was created.

Explicitly:

- it is written by cognition onto the artifact at formation time
- it reflects how the being interpreted the full active context at that moment
- that context may include multiple relationships, recalled experience,
  current salience, current tension, and broader situational understanding
- it is not copied from a relationship record
- it is not averaged, blended, or mechanically derived from relationship trust
  values
- it is a snapshot of interpretive posture at formation time
- once written, it stays fixed on the artifact permanently
- later trust movement does not retroactively change it

## Structural Note

`anchor_set` and `trust_at_formation` do different work.

- `anchor_set` grounds the artifact structurally for later retrieval
- `trust_at_formation` preserves the cognitive trust judgment present when the
  artifact was formed

The schema stays unchanged under this clarification.
