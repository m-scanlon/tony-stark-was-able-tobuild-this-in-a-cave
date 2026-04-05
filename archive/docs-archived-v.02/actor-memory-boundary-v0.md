# Actor Memory Boundary v0

## Purpose

This document captures the current memory-boundary direction for actors.

The main point is:

- the system should not default to one global retained experience store that every actor reads directly

## Core Framing

Actors should remain composable.

That becomes harder if every actor shares one ambient retained memory pool.

The current direction is therefore:

- each actor owns its own retained experience boundary
- cross-actor access should be intentional and contract-bounded

## No Ambient Global Experience

The system should not assume:

- all retained experience is globally visible
- all actors can freely read all other actors' retained experience
- orchestration implies ambient memory access

That would make:

- actor boundaries weak
- contracts muddy
- delegation leaky

## Per-Actor Retained Experience

Each actor should have its own retained experience surface.

That means:

- actor-local recall is normal
- actor-local learning writes into that actor's retained layer
- another actor should not automatically inherit that retained context

This keeps:

- memory segmented
- roles composable
- delegation explicit

## Subject Memory Domains

Within a actor's retained experience boundary, it is useful to distinguish subject-scoped memory domains from actor identity or registry state.

These domains are not automatically:

- separate actors
- external capabilities
- registration records

They are retained-experience domains that help bound what kinds of traces and later derived artifacts a actor is responsible for.

### Human Subject Domains

For a human-subject-facing actor such as `Jarvis`, the current useful retained domains are:

- `identity`
- `preferences`
- `boundaries`
- `interaction_style`

These should be understood primarily as retained-experience domains.

They may begin from retained trace and later support derived understanding, salience, or tension.

### System Subject Domains

For a system-subject-facing actor, the current useful retained domains are:

- `identity`
- `constraints`
- `health`
- `surface_behavior`

These should also be treated as retained-experience domains rather than as registry truth.

For example, they are not the same thing as:

- transport metadata
- verified capability inventory
- published capability-surface contracts
- registration envelopes

Those belong to the registration and contract layers.

The retained domains above instead capture what it is like to operate with that system subject over time:

- stable identity-oriented facts
- recurring limits and incompatibilities
- degradation, failure, and recovery patterns
- practical behavioral quirks of exposed surfaces

## Actor-To-Actor Access

If one actor needs something from another actor, that access should happen through a deliberate boundary.

Examples:

- typed stimulus passing
- explicit actor-to-actor recall
- emitted runtime artifacts that another actor may consume

This means:

- access is intentional
- access is typed
- access is contract-bounded

## Orchestrator Actors

An orchestrator actor should not receive ambient access to every actor's retained experience.

Instead, an orchestrator should:

- receive typed stimuli
- route work to other actors
- receive typed outputs or recall products back

This makes orchestration composable rather than omniscient.

## Thin Global Registry

The system may still need a thin shared global layer.

But that layer should be:

- routing-oriented
- identity-oriented
- contract-oriented

not:

- a universal retained experience blob

The thin global layer may include:

- actor registry
- actor identity
- callable actor addressing
- published stimulus types
- capability registration metadata

## Recall Exposure Boundary

Each actor contract should eventually define how much of that actor's retained experience is exposed to others.

This boundary should answer questions like:

- who may ask this actor for recall
- what stimulus types may trigger that recall
- what shape of recall package may be returned

This keeps memory access from becoming implicit and uncontrolled.

## Why This Matters

This direction makes the system read more like a typed runtime:

- actors have contracts
- actors accept typed stimuli
- actors emit typed stimuli
- actors retain local experience
- actors only share through explicit boundaries

That is a stronger substrate for composition than one shared memory soup.

## Current Design Posture

The strongest current claims are:

- no default global experience store
- per-actor retained experience is the default
- actor-to-actor access should happen through typed and contract-bounded boundaries
- a thin global registry may still exist for routing and identity

## Short Framing

Actors should keep their own retained experience by default.

Other actors should only access that experience intentionally through typed, contract-bounded boundaries.

The system may still keep a thin global registry, but not a shared global memory blob.
