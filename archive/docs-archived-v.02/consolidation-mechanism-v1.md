# Learning / Consolidation Mechanism (v1)

## Core Framing

Learning is the write path from episodes into retained experience and structure.

It is separate from recall.

Its job is to decide what from completed activity should survive beyond the episode.

Episode closure is the current natural kickoff point for learning.

The owning actor may emit a minimal learn request stimulus against the just-closed episode.

## Learn Request Payload

For `v1`, the useful current request payload is:

```ts
type LearnRequest = {
  episode_id: string
}
```

`episode_id` is the minimal primitive payload for now.

The surrounding kernel envelope should stay minimal.

For `v1`, the useful current rule is:

- `emitter_surface` travels in the kernel envelope
- `episode_id` stays in the learn request payload

For ordinary learn traffic, that `emitter_surface` will usually be the calling actor surface.

The kernel can derive caller contract and authorization from `emitter_surface` rather than requiring extra caller metadata inside the request itself.

## Public Contract Shape

Learning should currently be modeled as:

- one request stimulus
- one response envelope

The response envelope carries:

- `status`
- `reason`
- optional actor-defined payload describing what was written

## Inputs

Learning may draw from:

- completed episodes
- bounded recall results already written into the episode
- interaction records
- runtime artifacts produced during runtime execution
- previously active retained artifacts present in recall

## Outputs

Learning may produce updates to two layers:

### 1. Structure

Structure updates include:

- entities
- relationships

### 2. Retention Layer

Retention updates include retained artifacts such as:

- `retained_trace`
- `retained_understanding`
- `retained_salience`
- `retained_tension`

Conceptually, the response payload may summarize:

```ts
type LearnPackage = {
  episode_id: string
  retained_artifact_ids: string[]
  structure_update_ids?: string[]
}
```

This package should be understood primarily as a write receipt.

## Core Direction

The strongest current direction is:

- first preserve factual retained happenings as traces
- then derive other retained artifacts from those traces and the broader episode context

This keeps the retention layer grounded in what happened while still allowing experience to survive as more than explicit understanding.

## Runtime Artifacts And Learning

Runtime artifacts should remain available to learning.

That means in-episode interpretation or other transient runtime execution may inform what gets retained later.

But runtime artifacts are not retained by default.

Learning is the selection boundary.

## Current Design Posture

The strongest current claims are:

- learning is a request/response contract, not a freeform command string
- `episode_id` is the minimal useful request payload for `v1`
- the response envelope is primarily a write receipt
- retention stays selective and grounded in retained trace

## Short Framing

Learning is the retained write contract.

A actor emits a minimal learn request stimulus for a closed episode, and the response envelope confirms what retained artifacts or structure updates were produced.
