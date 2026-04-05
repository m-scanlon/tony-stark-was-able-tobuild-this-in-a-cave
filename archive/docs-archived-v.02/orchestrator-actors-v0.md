# Orchestrator Actors v0

## Purpose

This document captures the current role of orchestrator actors.

The main point is:

- some actors exist primarily to route, delegate, and merge work rather than to perform all specialized work themselves

## Core Framing

The current actor model is separating into at least two broad actor classes:

- base actors
- orchestrator actors

Base actors perform bounded specialized work.

Orchestrator actors coordinate other actors through typed stimulus and contract-bounded delegation.

## What Orchestrator Actors Do

An orchestrator actor should primarily:

- accept typed stimuli
- decide which actor should handle what
- emit typed stimuli to other actors
- request or consume bounded recall when needed
- merge outputs back into a higher-level result or routing decision

This means orchestrator actors are fundamentally coordination actors.

They should not be treated as omniscient global workers.

## What Orchestrator Actors Are Not

Orchestrator actors should not default to:

- doing all specialized work themselves
- owning ambient access to every actor's retained experience
- bypassing contracts or runtime validation

They remain ordinary actors in the runtime model.

What changes is their role, not their ontology.

## Relationship To Memory

Orchestrator actors do not imply a global memory blob.

The current direction is:

- memory remains actor-local by default
- orchestrators only receive what is intentionally exposed through typed stimuli or bounded recall

This keeps orchestration composable instead of omniscient.

## Relationship To Typed Stimulus

Typed stimulus is what makes orchestrator behavior viable.

An orchestrator actor can:

1. receive typed input
2. transform that into another typed stimulus
3. delegate that to a base actor or another orchestrator
4. receive typed results back
5. decide what happens next

This is the core of the orchestration model.

## Relationship To Primitives

Orchestrator actors still use the same primitive substrate as other actors:

- `recall`
- `learn`
- `observe`
- `act`

The difference is that much of their work consists of:

- selecting targets
- issuing delegation
- composing results

rather than being the actor that performs every concrete external interaction directly.

## Jarvis And Stark

The current important clarification is:

- `Jarvis` and `Stark` should be treated as orchestrator actors

They are not the base worker layer.

### Jarvis

`Jarvis` is the canonical user-facing orchestrator actor.

Its role is to:

- track what matters in user-facing context
- route user-facing work
- coordinate which actors should participate
- merge meaning and continuity back into higher-level action

### Stark

`Stark` is the canonical structural orchestrator actor.

Its role is to:

- coordinate structural work
- route probe, registration, and contract-related tasks
- manage structural continuity and revision through the actor graph
- merge structural outcomes back into the live system picture

## Why This Distinction Matters

If Jarvis and Stark are treated as base worker actors, they become:

- overloaded
- too magical
- too broad

If they are treated as orchestrator actors instead, the system becomes cleaner:

- worker behavior stays reusable
- orchestration stays explicit
- delegation becomes a first-class design tool

## Current Design Posture

The strongest current claims are:

- orchestrator actors are coordination actors
- orchestrator actors should route via typed stimulus rather than ambient shared memory
- Jarvis and Stark are orchestrator actors, not base worker actors
- base actors and orchestrator actors should remain distinct

## Short Framing

Orchestrator actors coordinate other actors.

They route typed stimuli, delegate bounded work, and merge results.

Jarvis and Stark should now be understood as orchestrator actors rather than base worker actors.
