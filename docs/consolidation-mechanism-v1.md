# Learning / Consolidation Mechanism (v1)

## Core Framing

Learning is the write path from episodes into retained experience and structure.

It is separate from recall.

Its job is to decide what from completed activity should survive beyond the episode.

Episode closure is the current natural kickoff point for learning.

The owning node may emit:

```text
skyra <node> learn -episode_id <episode_id>
```

against the just-closed episode.

`episode_id` is the minimal primitive arg for now.

The surrounding kernel envelope should stay minimal.

For `v1`, the useful current rule is:

- `calling_actor` travels in the kernel envelope
- `episode_id` stays in the learn command args

The kernel can derive caller contract and authorization from `calling_actor` rather than requiring extra caller metadata inside the command itself.

## Inputs

Learning may draw from:

- completed episodes
- bounded recall results already written into the episode
- interaction records
- runtime artifacts produced during runtime execution
- previously active retained artifacts present in recall

These inputs give learning both:

- what happened during the episode
- what was active in the episode while it happened

## Working Learn Command Shape

If learning is kicked off as an emitted runtime command, a good current working shape is:

```text
skyra <node> learn -episode_id <episode_id>
```

Conceptually:

```ts
type LearnArgs = {
  episode_id: string
}
```

This should stay minimal for `v1`.

The closed episode already contains the main material learning needs.

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

Conceptually, the primitive-specific result may be summarized as:

```ts
type LearnPackage = {
  episode_id: string
  retained_artifact_ids: string[]
  structure_update_ids?: string[]
}
```

The exact result detail can stay open, but the important point is that learning returns typed update information through the shared kernel result-routing/writeback path.

This package should be thought of primarily as a write receipt.

The writes have already happened by the time it returns.

Its job is to acknowledge what retained artifacts and structure updates were produced, not to become the durable owner of them.

## Core Direction

The strongest current direction is:

- first preserve factual retained happenings as traces
- then derive other retained artifacts from those traces and the broader episode context

This keeps the retention layer grounded in what happened while still allowing experience to survive as more than explicit understanding.

## High-Level Shape

At a high level, learning should:

1. identify bounded factual happenings worth preserving as retained traces
2. resolve or update the structure those traces refer into
3. derive retained understanding, salience, and tension where appropriate
4. attach anchors and provenance
5. select what should actually be retained

The exact mechanics remain open.

This document only fixes the high-level direction.

## Retained Trace First

`retained_trace` should remain distinct.

Its role is to preserve a bounded factual retained happening.

That gives later retained artifacts grounding in something that occurred rather than requiring everything to collapse directly into meaning.

## Derived Retained Artifacts

Derived retained artifacts may include:

- what the happening meant
- what in it should carry weight later
- what remained unresolved

These do not need to appear for every trace.

Retention remains selective.

## Runtime Artifacts And Learning

Runtime artifacts should remain available to learning.

That means in-episode interpretation or other transient runtime execution may inform what gets retained later.

But runtime artifacts are not retained by default.

Learning is the selection boundary.

## Anchors And Provenance

All retained artifacts should be anchored into canonical structure through `anchor_set`.

Derived retained artifacts may also carry trace grounding such as:

- `source_trace_ids`

Retained artifacts may also preserve which prior retained artifacts shaped their formation through:

- `context_artifact_ids`

When learning writes retained artifacts, it should also update the anchor lookup layer in the same write path.

That means:

- persist the retained artifact record
- project its `anchor_set` into the retrieval indexes
- keep artifact storage and anchor lookup synchronized

This gives recall fast entry points without overloading the retained artifact record itself.

Additional direct lookup tables may also be maintained in that same write path.

For example:

```ts
episode_to_artifacts[episode_id] -> retained_artifact_ids[]
node_to_artifacts[node_id] -> retained_artifact_ids[]
intent_to_artifacts[intent_id] -> retained_artifact_ids[]
trace_to_derived[trace_id] -> retained_artifact_ids[]
```

These lookup tables should remain part of the retention/index layer rather than turning the node itself into the owner of retained experience.

## Key Constraints

### 1. Selective Retention

Not everything from an episode should survive.

Retention should remain selective and bounded.

### 2. Trace Stays Distinct

`retained_trace` should not collapse into the same semantic role as understanding, salience, or tension.

### 3. Runtime First, Retained Later

In-episode runtime execution may produce runtime artifacts freely.

Only learning decides whether any of that becomes retained experience.

Learning should therefore operate over a closed episode rather than mutating active in-episode runtime state directly.

### 4. No Blind Overwrite

Learning should evolve retained experience over time rather than replacing it blindly.

## Recall Separation

Recall is the read path.

Learning is the write path.

Recall brings retained artifacts into the active episode.

Learning writes new or revised retained artifacts out of completed episode activity.

The retention layer and its lookup/index surfaces remain the source of truth for retained artifacts after those writes land.

## Open Edges

The following still need exact definition:

- trace extraction rules
- retention thresholds
- same-artifact resolution and merging
- structure update policy
- provenance detail level
- learning result writeback detail beyond the minimal typed package

## Short Framing

Learning is the write path from episodes into retained experience and structure.

It should preserve factual retained traces first and derive other retained artifacts from there when appropriate.

Runtime artifacts may inform learning, but they are not retained by default.

For `v1`, the natural kickoff is `skyra <node> learn -episode_id <episode_id>` after episode closure.
