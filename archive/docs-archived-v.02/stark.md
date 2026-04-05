# Stark

## Working Definition

`Stark` is the canonical name for the structural orchestrator actor.

References to the older `Architect` and `system actor` concepts should collapse into `Stark`.

Functional subtitle:

- `Stark` = structural orchestrator actor

## Purpose

Stark exists to maintain and reshape system structure.

This includes:

- actor topology
- actor registry
- stimulus registry
- actor contracts
- actor contract publication
- capability attachment and steward mediation
- distribution of responsibility
- routing structural work to other actors when needed

## Same Model As Other Actors

Stark is not outside the actor model.

It still participates through:

- bounded episodes
- recall
- retained experience
- emitted and received stimulus

What differs is role and authority, not ontology.

## Primary Stimulus

Stark is primarily driven by structural pressure rather than ordinary user content.

Typical stimuli include:

- contract insufficiency
- capability gaps
- topology problems
- actor blockage or structural mismatch
- requests for structural revision
- typed structural results returned from other actors

For device/bootstrap flows, Stark should also be understood as the actor that creates or selects the structural type of incoming device-side stimuli.

## Outputs

Stark does not directly mutate the system.

It coordinates structural work and emits valid structural stimulus that the kernel can route and apply, such as:

- `birth_actor`
- revise contract
- attach capability
- restrict scope
- replace actor
- deprovision actor

After bootstrap, Stark is also the publisher of later actor contracts.

Stark should also be understood as the structural owner of the live actor registry and stimulus registry.

The kernel remains the authority that applies those structural changes.

For `v1`, the important current `act` boundary is that Stark should remain system-facing.

That means Stark should primarily own work such as:

- probing
- registration writes
- structural mediation over capability surfaces

Stark should not emit outward human-facing `act` in `v1`.

## Bootstrap Role

`Stark` is the bootstrap structural actor.

The kernel births Stark at system startup from a hardcoded contract.

After that point, Stark becomes the structural origin for later actor contract publication and structural orchestration.

## Structural History

Stark may operate over a structural slice of history, such as:

- which actors exist
- which contracts exist
- contract revisions over time
- actor lifecycle state
- capability layout

It does not require total system state by default.

## Design Principles

- Stark is a actor, not a hidden control plane
- Stark is an orchestrator actor, not a base worker actor
- Stark shapes structure; it does not bypass the kernel
- Stark uses the same general runtime model as other actors
- Stark is the canonical structural orchestrator actor

## Relationship To Other Actors

Other actors perform bounded work within contract.

Stark governs structure around that work and coordinates structural delegation.

So:

- base worker actors handle specialized task participation
- Stark handles structural continuity, routing, and revision

## Scope Of This Document

This document defines Stark at the role level.

It does not define:

- the full contract revision protocol
- the exact final structural stimulus schema
- the full Stark/Jarvis collaboration contract
- the exact structural routing policy Stark should use

## Short Framing

Stark is the canonical structural orchestrator actor.

It decides what system structure should exist and coordinates the actor graph around that decision.

It also publishes later actor contracts and registry entries, while the kernel makes those structural decisions executable.
