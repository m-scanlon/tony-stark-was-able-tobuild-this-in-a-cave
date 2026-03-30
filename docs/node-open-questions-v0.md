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
- the runtime should use the existing max heap rather than invent a second global queue
- each node should own a lightweight mailbox for already-routed events

Current posture:

- the kernel owns the unified max heap
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

- when a published contract arrives for a running node, what is the exact receipt-to-adoption flow

Why it matters:

- this defines runtime stability, auditability, and contract authority

Resolution:

- Stark may publish a contract for a running node while that node is mid-episode
- the kernel routes that contract publication through the same typed event flow
- the node receives that publication and holds the new contract in pending node state / mailbox flow
- the current episode continues under the currently active contract
- the new contract takes effect only when the current episode closes

Current posture:

- the node should not switch contracts mid-episode
- the current episode remains bounded by the contract that was active when that episode was running
- the node may receive the next contract before episode close
- that next contract remains pending until episode close
- once the episode is over, the pending contract becomes the active contract
- in `v1`, episode closure should be driven by inactivity rather than an abstract hard timeout

Design note:

- this keeps contract transitions aligned to bounded runtime context
- it also avoids mixing one episode across two different contract regimes

## 6. Command Surface

Question:

- how should the contract express allowed command sets and commands
- how much of the cognition envelope should be named alongside that command surface
- how much command argument structure should be fixed now

Why it matters:

- this is the node's real execution boundary

Resolution:

- the active node contract should include the command surface
- `v1` should use a simple command_set/command allowlist
- one working protocol shape is `skyra <command_set> <command> -<args>`

Current posture:

- the contract should name the allowed command sets
- each allowed command set should declare its allowed commands
- avoid wildcard or pattern matching in `v1`
- keep the command allowance surface explicit and inspectable
- `skyra primitive interact` is one valid example inside that surface
- exact command argument schemas remain open for now

Why this is the right `v1` move:

- simple
- auditable
- easy to validate at runtime
- flexible enough for the command-set-based model

## Current Design Posture

The strongest current claim is:

- these questions should be resolved one at a time in runtime dependency order

## Short Framing

Node birth and ownership are now stable enough that the next work is the remaining node process questions.

This document tracks those questions in the order they should be resolved.
