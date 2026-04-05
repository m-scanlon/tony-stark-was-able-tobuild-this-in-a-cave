# Bootstrap And Startup (Prelim)

## Purpose

This document preserves the current startup boundary for the runtime.

The main point is:

- local substrate boot should be deterministic
- higher-level startup coordination should happen through actors

This is not a final boot design.

It is a prelim snapshot of the current shape.

## Core Framing

The startup model should distinguish between:

- the daemon process
- the local runtime infrastructure inside that daemon
- deterministic bootstrap
- runtime actors created after bootstrap

Without this split, the architecture gets muddy very quickly.

## Daemon

The daemon is the runtime container/process running on the device.

It is the thing that hosts the local runtime.

The daemon is not one actor among other actors.

It is the process that contains the runtime substrate and the actor graph.

## Local Runtime Infrastructure

The current local infrastructure inside the daemon is:

- SQLite
- mailbox/router
- kernel
- actor registry storage

These are substrate components.

They are not ordinary runtime actors.

## Deterministic Bootstrap

Bootstrap should stand up the local substrate deterministically.

That means bootstrap is responsible for:

1. starting the daemon runtime
2. opening or initializing SQLite
3. bringing up the mailbox/router
4. bringing up the kernel
5. bringing up actor registry storage
6. loading the shipped startup contract bundle
7. creating the first required runtime actor

The important boundary is:

- bootstrap brings up local infrastructure
- bootstrap does not perform the higher-level reasoning of the system

## First Actor

The current best first actor is:

- `Stark`

That means:

- bootstrap creates `Stark`
- bootstrap sends the initial startup stimulus to `Stark`

This keeps the split clean:

- bootstrap handles local deterministic startup
- `Stark` handles structural startup from inside the actor runtime

## Bootstrap Actor Set

Startup may require more than one shipped bootstrap-time actor.

That means the runtime should allow a small hardcoded bootstrap actor set rather
than pretending `Stark` is the only hardcoded startup participant.

`Stark` is still the primary structural bootstrap actor, but `Stark` may depend
on other shipped startup actors that are also available at init time.

These actors should be fully instantiated during bootstrap.

The important split is:

- bootstrap fully instantiates a small hardcoded startup actor set
- `Stark` is the first structural actor inside that set
- later non-bootstrap actors are still born through the normal kernel birth path

## Why Stark Starts First

`Stark` is the structural orchestrator actor.

Startup after substrate bring-up is a structural problem:

- what actors should exist
- what responsibilities should be active
- what startup work should be delegated

That makes `Stark` the right first actor for `v0`.

## What Stark Does Next

After bootstrap creates `Stark`, `Stark` should orchestrate the startup sequence inside the runtime.

That may include creating or coordinating actors such as:

- `probe`
- `registration`
- later other startup or coordination actors

The important current idea is:

- `Stark` is the first actor
- `Stark` then stands up the rest of the actor-level startup behavior

## Local Infra vs Actor Responsibilities

The key split is:

### Bootstrap / Substrate

Bootstrap and substrate own:

- SQLite
- mailbox/router
- kernel
- actor registry storage
- bootstrap-time actor birth

### Actors

Actors own higher-level runtime behavior such as:

- startup orchestration
- `birth_actor` requests
- probing
- registration
- permissions or readiness workflows
- cross-device coordination

This means a actor may be responsible for making sure devices or instances are talking correctly.

But that actor is not responsible for bringing up SQLite or the kernel.

## Cross-Device Coordination

There is an important distinction between:

- local runtime boot
- distributed or cross-device coordination

Local runtime boot should stay deterministic.

Distributed coordination may be handled by actors after bootstrap.

That means:

- bootstrap brings up the local substrate
- actor-level startup logic may then create or coordinate actors responsible for network/device communication health

This keeps cross-device logic out of the deterministic boot layer.

## Working Startup Sequence

The current `v0` startup sequence is:

1. `skyra init` starts or connects to the daemon
2. daemon boot begins
3. bootstrap brings up:
   - SQLite
   - mailbox/router
   - kernel
   - actor registry storage
4. bootstrap loads shipped bootstrap/startup contracts
5. bootstrap fully instantiates the hardcoded bootstrap actor set, including `Stark`
6. bootstrap emits `system_boot` to `Stark`
7. `Stark` coordinates the already-shipped startup actors and requests birth of
   any later required actors
8. actor-level startup work begins
9. probing / registration / onboarding can then start
10. in `v1`, successful registration should be followed immediately by `birth_actor`

## Still Open

The following remain open:

- exact bootstrap contract format
- exact first startup stimulus shape
- exact startup actor set after `Stark`
- exact actor responsible for cross-device communication readiness
- exact relation between `registration`, later infra-management actors, and actor-birth follow-on policy beyond the default `v1` path

## Current Design Posture

The strongest current claims are:

- the daemon is the runtime container/process
- SQLite, mailbox/router, kernel, and actor registry storage are local substrate infrastructure
- local substrate boot should be deterministic
- `Stark` should be the first runtime actor
- startup may include a small hardcoded bootstrap actor set rather than only one
  hardcoded actor
- `Stark` is the primary structural bootstrap actor inside that set
- the bootstrap actor set is fully born during bootstrap rather than lazily
  instantiated later
- higher-level startup and coordination work should happen through actors after bootstrap

## Short Framing

Bootstrap should deterministically bring up the local runtime substrate and create `Stark`.

After that, startup becomes a actor problem.

`Stark` orchestrates the rest of the startup sequence from inside the runtime.
