# Learning / Consolidation Mechanism (v1)

## Core Framing

Learning is the write path from episodes into retained experience and structure.

It is separate from recall.

Its job is to decide what from completed activity should survive beyond the episode.

## Inputs

Learning may draw from:

- completed episodes
- the episode field
- interaction records
- runtime artifacts produced during runtime execution
- previously active retained artifacts present in recall

These inputs give learning both:

- what happened during the episode
- what was active in the episode while it happened

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

## Key Constraints

### 1. Selective Retention

Not everything from an episode should survive.

Retention should remain selective and bounded.

### 2. Trace Stays Distinct

`retained_trace` should not collapse into the same semantic role as understanding, salience, or tension.

### 3. Runtime First, Retained Later

In-episode runtime execution may produce runtime artifacts freely.

Only learning decides whether any of that becomes retained experience.

### 4. No Blind Overwrite

Learning should evolve retained experience over time rather than replacing it blindly.

## Recall Separation

Recall is the read path.

Learning is the write path.

Recall brings retained artifacts into the active episode.

Learning writes new or revised retained artifacts out of completed episode activity.

## Open Edges

The following still need exact definition:

- trace extraction rules
- retention thresholds
- same-artifact resolution and merging
- structure update policy
- provenance detail level
- timing of learning passes

## Short Framing

Learning is the write path from episodes into retained experience and structure.

It should preserve factual retained traces first and derive other retained artifacts from there when appropriate.

Runtime artifacts may inform learning, but they are not retained by default.
