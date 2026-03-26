# Structural Extraction Pipeline v0

This document defines the handoff contracts for:

`tokens stream in -> bulk processing -> entity and relationship representation`

It is only about new natural-language stimulus.

It does not define recall projection or frame assembly.

## Pipeline

The current working pipeline is:

1. token stream enters a stimulus buffer
2. the buffer emits a stable chunk
3. the chunk is processed in bulk
4. the processor emits structural fragments
5. the fragments are normalized into entity and relationship representation

## 1. Stimulus Buffer

The stimulus buffer accumulates streaming input until a stable chunk boundary is reached.

```ts
type StimulusBuffer = {
  id: string
  source_id: string
  tokens: StimulusToken[]
  opened_at: string
  updated_at: string
}
```

```ts
type StimulusToken = {
  value: string
  timestamp: string
}
```

## 2. Chunk Boundary

The buffer should emit a chunk when the input is stable enough for bulk parsing.

The exact policy remains open, but likely boundaries are:

- pause
- punctuation
- clause boundary
- utterance end

## 3. Stimulus Chunk

```ts
type StimulusChunk = {
  chunk_id: string
  source_id: string
  text: string
  started_at: string
  ended_at: string
}
```

This is the unit sent to the extraction processor.

## 4. Structural Extraction

Bulk processing reads one chunk and emits one or more structural fragments.

```ts
type StructuralExtraction = {
  chunk_id: string
  fragments: StructuralFragment[]
  created_at: string
}
```

## 5. Structural Fragment

A chunk may produce multiple fragments.

Each fragment should preserve entity-relationship binding.

```ts
type StructuralFragment = {
  fragment_id: string
  chunk_id: string
  text: string
  entities: EntityCandidate[]
  relationships: RelationshipCandidate[]
}
```

## 6. Entity Candidate

```ts
type EntityCandidate = {
  candidate_id: string
  surface_form: string
  span: TextSpan
  confidence: number
}
```

## 7. Relationship Candidate

Each relationship candidate should point to the entity candidates it binds.

```ts
type TextSpan = {
  start_char: number
  end_char: number
}
```

```ts
type RelationshipCandidate = {
  candidate_id: string
  surface_form: string
  span: TextSpan
  from_candidate_id: string
  to_candidate_id: string
  confidence: number
}
```

## 8. Representation Output

The extraction layer should output a representation that preserves:

- entity candidates
- relationship candidates
- the binding between them

This layer should preserve surface evidence, not finalize canonical structure.

That means:

- extraction produces candidates
- normalization happens later
- canonical resolution happens later

The handoff out of extraction is therefore a bound candidate representation, not final entities and relationships.

## Design Posture

- token streaming is continuous
- parsing is chunked, not per-token
- extraction is bulk
- one chunk may produce multiple fragments
- the minimum useful output is bound entity-relationship structure
