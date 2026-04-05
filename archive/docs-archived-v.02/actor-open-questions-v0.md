# Actor Open Questions v0

## Core Framing

These are the main remaining open questions around actor design after the current stimulus-first cleanup.

This document exists to track the next questions in a stable order and resolve them one at a time.

## Resolution Order

The current recommended order is:

1. event intake model
2. episode selection policy
3. `dependencyLedger` lifecycle
4. frame projection timing
5. contract publication timing
6. public surface modeling

## 1. Event Intake Model

Question:

- should the actor use one unified event queue with typed events
- or separate queues for stimulus ingress and non-stimulus control flow

Why it matters:

- this shapes `ingest_event(...)`
- this affects runtime simplicity
- this affects how returned response stimulus is normalized

Current candidate event families:

- `stimulus`
- `contract_publication`

Resolution:

- `v1` uses one unified typed event intake surface
- event ordering and global routing sit at the kernel front
- each actor owns a lightweight mailbox for already-routed events

## 2. Episode Selection Policy

Question:

- when a valid event arrives, when does the actor:
  - reuse the current episode
  - open a new episode
  - close an old episode first

Resolution:

- `v1` uses a time-based episode policy

Current posture:

- if the actor has no active episode, open one
- if the active episode has exceeded its time window, close it and open a new one
- otherwise reuse the current active episode

## 3. Dependency Ledger Lifecycle

Question:

- where does the actor track open downstream dependencies
- what statuses are needed
- how are returned response envelopes matched back to prior dispatches

Resolution:

- the actor should track whether a downstream dependency is currently open
- dependency tracking is a actor runtime concern, not a public contract concern
- if a response never arrives, that dependency should eventually time out

Current posture:

- the actor keeps an internal `dependencyLedger`
- the minimum useful distinction is:
  - `open`
  - `resolved`
  - `failed`
  - `timed_out`
- returned response stimulus should match back to a prior `dependency` through runtime identifiers

## 4. Frame Projection Timing

Question:

- when should the actor project a frame

Better framing:

- this is really an inference-readiness question

Current posture:

- leave this open for now
- keep the question attached to actor inference-readiness rather than raw frame-render timing

## 5. Contract Publication Timing

Question:

- when a published contract arrives for a running actor, what is the exact receipt-to-adoption flow

Resolution:

- the kernel routes contract publication through the same typed event flow
- the actor receives that publication and holds the new contract in pending actor state
- the current episode continues under the currently active contract
- the new contract takes effect only when the current episode closes

## 6. Public Surface Modeling

Question:

- how much of a actor's public request/response surface should be visible in the contract itself
- how much downstream execution-surface traversal should remain implementation detail

Current posture:

- the contract should expose the public request stimuli and response envelopes
- the public callable surface should look more like an API than a handler graph
- downstream traversal across other `ExecutionSurface`s does not need to be denormalized into the public contract by default

## Short Framing

The next actor work is no longer about command sets.

It is about event intake, dependency tracking, contract adoption timing, and how much of the API-shaped public stimulus surface should be explicit in the contract.
