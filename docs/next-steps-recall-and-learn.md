# Next Steps: Node, Recall, and Learn

This is the immediate focus for the next stretch of work.

The goal is to harden the runtime model around:

- `node` — the durable runtime operator
- `recall` — the read path that brings retained experience into the current episode frame
- `learn` — the write path that turns completed episode activity into retained artifacts

This work should happen before going deeper on device capability manifests or probe interfaces.

## Current Scope

Lock only these layers:

- node birth and node process
- node vs episode ownership
- episode field
- runtime artifacts
- retained artifacts
- recall flow
- learn flow

## Out of Scope

Do not focus on these yet:

- capability manifest design
- deep Stark/Jarvis device routing
- long-term structure promotion beyond placeholders
- hardware probe interfaces

## 5-Day Plan

### 1. Node Foundation v0

Define:

- what Stark provides at node birth
- what the node owns versus what the episode owns
- how the active contract gates acceptable stimulus/events
- how the node opens or reuses an episode
- how the node updates episode-local state and projects a frame
- what generic primitive execution substrate the node uses so the runtime can support loops like `OODA`, `ReAct`, or one-shot execution without hardcoding one fixed loop
- how the active contract defines the allowed loop envelope for that node
- how inference may choose the actual next step or loop progression inside that contract-bounded envelope
- how later contract updates may be adopted safely

Outputs:

- `node-birth-v0.md`
- `node-process-v0.md`
- `contract-update-v0.md`

### 2. Recall v0

Define:

- what inputs recall reads from the current episode
- how the scored episode field becomes a cue surface
- how candidate artifacts are fetched through `anchor_set`
- how `trace`, `understanding`, `salience`, and `tension` are admitted into frame
- what recall returns to cognition

Outputs:

- `recall-v0.md`
- 2-3 worked examples

### 3. Learn v0

Define:

- what episode data learning can inspect
- how runtime artifacts are included
- how `retained_trace` is extracted
- how derived artifacts are formed from traces
- when `context_artifact_ids` and `source_trace_ids` are written

Outputs:

- `learn-v0.md`
- one end-to-end example from episode to retained artifacts

### 4. Artifact Lifecycle

Define:

- duplicate vs new artifact
- merge vs separate trace
- reinforcement
- revision
- suppression
- conflict handling

Outputs:

- `artifact-lifecycle-v0.md`

### 5. Kernel Boundary

Define only the practical boundary:

- when `recall` can be called
- what the `recall` command shape looks like
- what learning reads after episode close
- what stays transient vs retained

Outputs:

- tighten `runtime-primitives-and-artifacts-prelim.md`

### 6. Canonical Cleanup

Turn the above into a stable working spine:

- overview
- node
- recall
- learn
- lifecycle
- examples
- trimmed open questions

## Two Hard Questions

### Recall

What exactly gets admitted into the current frame, and in what shape?

### Learn

What counts as one retained trace versus multiple traces?

## Principle

Do not overdesign the future interface before these paths are concrete.

The next useful canon is:

- `node`
- `recall`
- `learn`
- the retained artifact lifecycle between them
