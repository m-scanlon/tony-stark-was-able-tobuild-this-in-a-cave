# Retained Artifact Family

The retained artifact family is Skyra's long-term experience layer. These are the outputs of learning — what the system keeps after an episode closes.

There are four artifact types. They share a common base shape but preserve different aspects of experience.

---

## AnchorSet

**File:** `anchor-set.ts`

The structural link into canonical structure. Every retained artifact carries one. It holds entity and relationship ids that ground the artifact into the world the system already knows about.

This is what makes recall work — when the system needs to retrieve relevant experience, it matches against anchor sets rather than scanning raw text.

```ts
{
  entity_ids: string[]
  relationship_ids: string[]
}
```

---

## RetainedArtifactBase

**File:** `retained-artifact-base.ts`

The shared base shape across all four artifact types. Carries:

- `id` — unique identifier
- `kind` — discriminator: `"trace" | "understanding" | "salience" | "tension"`
- `anchor_set` — structural grounding into canonical entities and relationships
- `context_artifact_ids` (optional) — references to prior retained artifacts that influenced this artifact's formation. This preserves the layered nature of experience: a happening can shape an understanding, and a later happening can occur through that understanding.

---

## RetainedTrace

**File:** `retained-trace.ts`

**What happened.** The only artifact type that is factual and non-interpretive. A trace is a bounded natural-language record of an occurrence — not what it meant, not whether it mattered, just what took place.

This is the grounding artifact. Understanding, salience, and tension are all derived from traces.

- `happened` — bounded natural-language rendering of the occurrence
- `source_episode_ids` — which episodes produced this trace

---

## RetainedUnderstanding

**File:** `retained-understanding.ts`

**What it meant.** A derived artifact that preserves interpreted meaning from one or more traces. This is the system's record of having made sense of something.

- `interpretation` — the interpreted meaning
- `source_trace_ids` — which traces this understanding was derived from

---

## RetainedSalience

**File:** `retained-salience.ts`

**What mattered.** A derived artifact that preserves what carried weight or attention. Not every trace becomes salient — salience marks the subset of experience that the system registered as significant.

- `signal` — the salient signal
- `source_trace_ids` — which traces this salience was derived from

---

## RetainedTension

**File:** `retained-tension.ts`

**What remained unresolved.** A derived artifact that preserves open edges — conflicting signals, unresolved questions, things that don't yet cohere. Tension is not a failure state; it is the system's record that something is still in play.

- `unresolved` — the unresolved or conflicting significance
- `source_trace_ids` — which traces this tension was derived from

---

## The Core Distinction

Trace preserves occurrence. The other three preserve consequence.

All four share anchors and can reference each other through `context_artifact_ids`, but trace stays on the factual side of the line. Understanding, salience, and tension are interpretation — they are shaped by traces but are not traces themselves.
