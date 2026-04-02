# Next Steps: Node, Recall, and Learn

This is the immediate focus for the next stretch of work.

The goal is to harden the runtime model around:

- `node` — the durable runtime operator
- `recall` — the contract that brings retained experience into the current episode frame
- `learn` — the write path that turns completed episode activity into retained artifacts

This work should happen before going deeper on device capability manifests or probe interfaces.

## Current Scope

Lock only these layers:

- node birth and node process
- node vs episode ownership
- recall contract
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

Status:

- materially complete enough for `v1`

Covered:

- node birth
- node substrate interface
- node vs episode ownership
- event intake posture
- mailbox posture
- episode reuse/closure policy
- pending command handling posture
- contract publication timing
- command allowance surface
- node process framing

Canonical outputs now in place:

- `node-birth-v0.md`
- `node-substrate-interface-v0.md`
- `node-and-episode-ownership-v0.md`
- `node-process-v0.md`
- `node-open-questions-v0.md`
- `command-namespace-prelim.md`
- `interaction-unification-prelim.md`

Remaining node gap:

- inference-readiness / frame projection timing remains open

This should not block moving to recall.

### 2. Recall v0

Define:

- what current episode context heavy inference reads
- when heavy inference should emit a bounded recall command
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
- keep `v1` learning episode-bounded rather than requiring full node ancestry or orchestration trace replay

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
- `observe`
- `act`
- the retained artifact lifecycle between them
