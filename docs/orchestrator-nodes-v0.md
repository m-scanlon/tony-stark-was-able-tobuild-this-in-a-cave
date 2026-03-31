# Orchestrator Nodes v0

## Purpose

This document captures the current role of orchestrator nodes.

The main point is:

- some nodes exist primarily to route, delegate, and merge work rather than to perform all specialized work themselves

## Core Framing

The current node model is separating into at least two broad node classes:

- base nodes
- orchestrator nodes

Base nodes perform bounded specialized work.

Orchestrator nodes coordinate other nodes through typed stimulus and contract-bounded delegation.

## What Orchestrator Nodes Do

An orchestrator node should primarily:

- accept typed stimuli
- decide which node should handle what
- emit typed stimuli to other nodes
- request or consume bounded recall when needed
- merge outputs back into a higher-level result or routing decision

This means orchestrator nodes are fundamentally coordination nodes.

They should not be treated as omniscient global workers.

## What Orchestrator Nodes Are Not

Orchestrator nodes should not default to:

- doing all specialized work themselves
- owning ambient access to every node's retained experience
- bypassing contracts or runtime validation

They remain ordinary nodes in the runtime model.

What changes is their role, not their ontology.

## Relationship To Memory

Orchestrator nodes do not imply a global memory blob.

The current direction is:

- memory remains node-local by default
- orchestrators only receive what is intentionally exposed through typed stimuli or bounded recall

This keeps orchestration composable instead of omniscient.

## Relationship To Typed Stimulus

Typed stimulus is what makes orchestrator behavior viable.

An orchestrator node can:

1. receive typed input
2. transform that into another typed stimulus
3. delegate that to a base node or another orchestrator
4. receive typed results back
5. decide what happens next

This is the core of the orchestration model.

## Relationship To Primitives

Orchestrator nodes still use the same primitive substrate as other nodes:

- `recall`
- `learn`
- `interact`

The difference is that much of their work consists of:

- selecting targets
- issuing delegation
- composing results

rather than being the node that performs every concrete external interaction directly.

## Jarvis And Stark

The current important clarification is:

- `Jarvis` and `Stark` should be treated as orchestrator nodes

They are not the base worker layer.

### Jarvis

`Jarvis` is the canonical user-facing orchestrator node.

Its role is to:

- track what matters in user-facing context
- route user-facing work
- coordinate which nodes should participate
- merge meaning and continuity back into higher-level action

### Stark

`Stark` is the canonical structural orchestrator node.

Its role is to:

- coordinate structural work
- route probe, registration, and contract-related tasks
- manage structural continuity and revision through the node graph
- merge structural outcomes back into the live system picture

## Why This Distinction Matters

If Jarvis and Stark are treated as base worker nodes, they become:

- overloaded
- too magical
- too broad

If they are treated as orchestrator nodes instead, the system becomes cleaner:

- worker behavior stays reusable
- orchestration stays explicit
- delegation becomes a first-class design tool

## Current Design Posture

The strongest current claims are:

- orchestrator nodes are coordination nodes
- orchestrator nodes should route via typed stimulus rather than ambient shared memory
- Jarvis and Stark are orchestrator nodes, not base worker nodes
- base nodes and orchestrator nodes should remain distinct

## Short Framing

Orchestrator nodes coordinate other nodes.

They route typed stimuli, delegate bounded work, and merge results.

Jarvis and Stark should now be understood as orchestrator nodes rather than base worker nodes.
