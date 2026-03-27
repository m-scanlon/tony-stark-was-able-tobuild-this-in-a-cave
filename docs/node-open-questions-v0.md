# Node Open Questions v0

## Core Framing

These are the main remaining open questions around node design after locking:

- node birth
- node vs episode ownership
- node substrate interface

This document exists to track the next questions in a stable order and resolve them one at a time.

## Resolution Order

The current recommended order is:

1. event intake model
2. episode selection policy
3. pending command lifecycle
4. frame projection timing
5. contract publication timing
6. command allowance surface

## 1. Event Intake Model

Question:

- should the node use one unified event queue with typed events
- or separate queues for stimulus ingress and command-result writeback

Why it matters:

- this shapes `ingest_event(...)`
- this affects runtime simplicity
- this affects how writeback and stimulus timing interact

Current candidate event families:

- `stimulus`
- `command_result`
- `contract_publication`

Resolution:

- `v1` uses one unified typed event intake surface
- event ordering and global routing should sit at the kernel front
- the runtime should use the existing priority heap rather than invent a second global queue
- each node should own a lightweight mailbox for already-routed events

Current posture:

- the kernel owns the unified priority heap
- the heap may contain typed events such as:
  - `stimulus`
  - `command_result`
  - `contract_publication`
- once routed, events land in the target node mailbox
- the node remains event-driven but does not need its own priority queue
- the mailbox should stay lightweight:
  - FIFO list
  - ring buffer
  - or equivalent simple pending-event holder

## 2. Episode Selection Policy

Question:

- when a valid event arrives, when does the node:
  - reuse the current episode
  - open a new episode
  - close an old episode first

Why it matters:

- this defines the actual bounded-unit behavior of the runtime

Resolution:

- `v1` uses a time-based episode policy

Current posture:

- if the node has no active episode, open one
- if the active episode has exceeded its time window, close it and open a new one
- otherwise reuse the current active episode

Why this is the right `v1` move:

- simple
- predictable
- avoids premature semantic episode-boundary logic
- easy to replace later with richer policies

## 3. Pending Command Lifecycle

Question:

- where does the node track outstanding dispatched commands
- what statuses are needed
- how are command results matched back to prior dispatches

Why it matters:

- async command execution is already implied by the substrate interface

Resolution:

- the node should track whether a command is currently outstanding
- outstanding command tracking is a node runtime concern, not a contract concern
- if a command result never arrives, it should eventually time out

Current posture:

- the node keeps an internal pending-command registry
- the minimum useful distinction is:
  - outstanding
  - completed
  - failed
  - timed_out
- command result events should match back to a prior dispatch through runtime identifiers

Design note:

- timeout behavior may be imperfect in `v1`
- that is acceptable because it belongs to runtime mechanics, not to contract shape

## 4. Frame Projection Timing

Question:

- when should the node project a frame

Candidate timings:

- after every accepted event
- only before inference
- after every command writeback
- some combination of the above

Why it matters:

- this affects runtime cost
- this affects responsiveness
- this affects how “current” the frame really is

Better framing:

- this is really an inference-readiness question
- the deeper issue is not “when do we render a frame”
- the deeper issue is “when has the node experienced enough to justify an inference step”

Current posture:

- leave this open for now
- do not force a `v1` answer yet
- keep the question attached to node inference-readiness rather than raw frame-render timing

## 5. Contract Publication Timing

Question:

- when does a published contract become effective on a running node

Candidate timings:

- immediately
- next safe boundary
- next episode

Why it matters:

- this defines runtime stability versus agility

Resolution:

- a published contract should take effect when the current episode closes

Current posture:

- the node should not switch contracts mid-episode
- the current episode remains bounded by the contract that was active when that episode was running
- a newly published contract becomes active after the current episode resolves or closes
- in `v1`, episode closure should be driven by inactivity rather than an abstract hard timeout

Design note:

- this keeps contract transitions aligned to bounded runtime context
- it also avoids mixing one episode across two different contract regimes

## 6. Command Allowance Surface

Question:

- how should the contract express allowed namespaces and commands

Why it matters:

- this is where command namespaces meet node contracts

Resolution:

- `v1` should use a simple namespace/command allowlist

Current posture:

- the contract should name the allowed namespaces
- each allowed namespace should declare its allowed commands
- avoid wildcard or pattern matching in `v1`
- keep the command allowance surface explicit and inspectable

Why this is the right `v1` move:

- simple
- auditable
- easy to validate at runtime
- flexible enough for the namespace-based command model

## Current Design Posture

The strongest current claim is:

- these questions should be resolved one at a time in runtime dependency order

## Short Framing

Node birth and ownership are now stable enough that the next work is the remaining node process questions.

This document tracks those questions in the order they should be resolved.
