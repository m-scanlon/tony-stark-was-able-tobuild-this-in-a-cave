# Next Steps: Actor, Recall, and Learn

This is the immediate focus for the next stretch of work.

The goal is to harden the runtime model around:

- `actor` — the durable runtime operator
- `recall` — the contract that brings retained experience into the current episode frame
- `learn` — the write path that turns completed episode activity into retained artifacts

## Current Scope

Lock only these layers:

- actor birth and actor process
- actor vs episode ownership
- recall contract
- runtime artifacts
- retained artifacts
- recall flow
- learn flow

## 5-Day Plan

### 1. Actor Foundation v0

Status:

- materially complete enough for `v1`

Covered:

- actor birth
- actor substrate interface
- actor vs episode ownership
- event intake posture
- mailbox posture
- episode reuse/closure policy
- `dependencyLedger` handling posture
- contract publication timing
- public request/response surface
- actor process framing

### 2. Recall v0

Define:

- what current episode context heavy inference reads
- when heavy inference should emit a bounded recall request stimulus
- how candidate artifacts are fetched through `anchor_set`
- how `trace`, `understanding`, `salience`, and `tension` are admitted into frame
- what recall returns inside the response envelope

### 3. Learn v0

Define:

- what episode data learning can inspect
- how runtime artifacts are included
- how `retained_trace` is extracted
- how derived artifacts are formed from traces
- when `context_artifact_ids` and `source_trace_ids` are written

### 4. Artifact Lifecycle

Define:

- duplicate vs new artifact
- merge vs separate trace
- reinforcement
- revision
- suppression
- conflict handling

### 5. Kernel Boundary

Define only the practical boundary:

- when `recall` can be called
- what the recall request payload looks like
- what learning reads after episode close
- what stays transient vs retained

## Principle

Do not overdesign the future interface before these paths are concrete.

The next useful canon is:

- `actor`
- `recall`
- `learn`
- `observe`
- `act`
- the retained artifact lifecycle between them
