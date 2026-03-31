# Bootstrap And Startup (Prelim)

## Purpose

This document preserves the current startup boundary for the runtime.

The main point is:

- local substrate boot should be deterministic
- higher-level startup coordination should happen through nodes

This is not a final boot design.

It is a prelim snapshot of the current shape.

## Core Framing

The startup model should distinguish between:

- the daemon process
- the local runtime infrastructure inside that daemon
- deterministic bootstrap
- runtime nodes created after bootstrap

Without this split, the architecture gets muddy very quickly.

## Daemon

The daemon is the runtime container/process running on the device.

It is the thing that hosts the local runtime.

The daemon is not one node among other nodes.

It is the process that contains the runtime substrate and the node graph.

## Local Runtime Infrastructure

The current local infrastructure inside the daemon is:

- SQLite
- mailbox/router
- kernel
- node registry storage

These are substrate components.

They are not ordinary runtime nodes.

## Deterministic Bootstrap

Bootstrap should stand up the local substrate deterministically.

That means bootstrap is responsible for:

1. starting the daemon runtime
2. opening or initializing SQLite
3. bringing up the mailbox/router
4. bringing up the kernel
5. bringing up node registry storage
6. loading the shipped startup contract bundle
7. creating the first required runtime node

The important boundary is:

- bootstrap brings up local infrastructure
- bootstrap does not perform the higher-level reasoning of the system

## First Node

The current best first node is:

- `Stark`

That means:

- bootstrap creates `Stark`
- bootstrap sends the initial startup stimulus to `Stark`

This keeps the split clean:

- bootstrap handles local deterministic startup
- `Stark` handles structural startup from inside the node runtime

## Bootstrap Node Set

Startup may require more than one shipped bootstrap-time node.

That means the runtime should allow a small hardcoded bootstrap node set rather
than pretending `Stark` is the only hardcoded startup participant.

`Stark` is still the primary structural bootstrap node, but `Stark` may depend
on other shipped startup nodes that are also available at init time.

These nodes should be fully instantiated during bootstrap.

The important split is:

- bootstrap fully instantiates a small hardcoded startup node set
- `Stark` is the first structural node inside that set
- later non-bootstrap nodes are still born through the normal kernel birth path

## Why Stark Starts First

`Stark` is the structural orchestrator node.

Startup after substrate bring-up is a structural problem:

- what nodes should exist
- what responsibilities should be active
- what startup work should be delegated

That makes `Stark` the right first node for `v0`.

## What Stark Does Next

After bootstrap creates `Stark`, `Stark` should orchestrate the startup sequence inside the runtime.

That may include creating or coordinating nodes such as:

- `probe`
- `registration`
- `node_creator`
- later other startup or coordination nodes

The important current idea is:

- `Stark` is the first node
- `Stark` then stands up the rest of the node-level startup behavior

## Local Infra vs Node Responsibilities

The key split is:

### Bootstrap / Substrate

Bootstrap and substrate own:

- SQLite
- mailbox/router
- kernel
- node registry storage
- first-node creation

### Nodes

Nodes own higher-level runtime behavior such as:

- startup orchestration
- node creation requests
- probing
- registration
- permissions or readiness workflows
- cross-device coordination

This means a node may be responsible for making sure devices or instances are talking correctly.

But that node is not responsible for bringing up SQLite or the kernel.

## Cross-Device Coordination

There is an important distinction between:

- local runtime boot
- distributed or cross-device coordination

Local runtime boot should stay deterministic.

Distributed coordination may be handled by nodes after bootstrap.

That means:

- bootstrap brings up the local substrate
- node-level startup logic may then create or coordinate nodes responsible for network/device communication health

This keeps cross-device logic out of the deterministic boot layer.

## Working Startup Sequence

The current `v0` startup sequence is:

1. `skyra init` starts or connects to the daemon
2. daemon boot begins
3. bootstrap brings up:
   - SQLite
   - mailbox/router
   - kernel
   - node registry storage
4. bootstrap loads shipped bootstrap/startup contracts
5. bootstrap fully instantiates the hardcoded bootstrap node set, including `Stark`
6. bootstrap emits `system_boot` to `Stark`
7. `Stark` coordinates the already-shipped startup nodes and requests birth of
   any later required nodes
8. node-level startup work begins
9. probing / registration / onboarding can then start

## Still Open

The following remain open:

- exact bootstrap contract format
- exact first startup stimulus shape
- exact startup node set after `Stark`
- exact node responsible for cross-device communication readiness
- exact relation between `node_creator`, `registration`, and later infra-management nodes

## Current Design Posture

The strongest current claims are:

- the daemon is the runtime container/process
- SQLite, mailbox/router, kernel, and node registry storage are local substrate infrastructure
- local substrate boot should be deterministic
- `Stark` should be the first runtime node
- startup may include a small hardcoded bootstrap node set rather than only one
  hardcoded node
- `Stark` is the primary structural bootstrap node inside that set
- the bootstrap node set is fully born during bootstrap rather than lazily
  instantiated later
- higher-level startup and coordination work should happen through nodes after bootstrap

## Short Framing

Bootstrap should deterministically bring up the local runtime substrate and create `Stark`.

After that, startup becomes a node problem.

`Stark` orchestrates the rest of the startup sequence from inside the runtime.
