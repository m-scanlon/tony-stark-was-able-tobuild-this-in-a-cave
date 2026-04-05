# Jarvis

## Working Definition

`Jarvis` is the canonical name for the user-facing orchestrator actor.

References to the older `Life` concept should collapse into `Jarvis`.

Functional subtitle:

- `Jarvis` = user-facing orchestrator actor

Jarvis is not outside the actor model.

It is a actor whose role is to track what matters, shape user-facing meaning, and coordinate user-facing work across other actors.

## Purpose

Jarvis exists to:

- determine what is meaningful in the current user context
- maintain continuity around priorities and what matters
- shape interpretation from the user-facing side
- guide which concerns deserve attention, response, or follow-through
- route user-facing work to the right participating actors when necessary

Jarvis is not the structural actor.

It is the meaning-and-attention counterpart to Stark.

## Same Model As Other Actors

Jarvis uses the same broad data model as other actors.

That means Jarvis still participates through:

- bounded episodes
- recall
- runtime execution
- retained experience

What differs is role and authority, not ontology.

The current important clarification is:

- Jarvis is not the base worker layer
- Jarvis is an orchestrator actor

## Primary Stimulus

Jarvis is primarily driven by user-facing activity and what becomes meaningful within it.

Typical stimuli include:

- ambiguous user requests
- shifts in user priority
- continuity across ongoing episodes
- emotionally or practically consequential context
- questions about what matters, not just what runs
- typed outputs from other actors that need user-facing synthesis or routing

## Outputs

Jarvis does not reshape system structure.

Instead, Jarvis helps shape:

- what should be attended to
- what should be brought into recall
- what user-facing interpretation should guide the current episode
- what outward act should be prioritized
- which other actors should participate in user-facing work
- how typed outputs from other actors should be merged back into the user-facing flow

The kernel remains the authority over execution and frame transition.

For `v1`, the important current `act` boundary is:

- `Jarvis` is the only actor allowed to emit outward human-facing `act` in `v1`

This keeps outward human communication anchored to the user-facing orchestrator role rather than spreading it across system-facing actors.

## Design Principles

- Jarvis is a actor, not a hidden persona layer
- Jarvis is an orchestrator actor, not a base worker actor
- Jarvis shapes meaning and attention; it does not own structure
- Jarvis uses the same general runtime model as other actors
- Jarvis is the canonical user-facing orchestrator actor

## Relationship To Stark

Jarvis and Stark are paired major actor roles.

- Jarvis orchestrates user-facing meaning and attention
- Stark orchestrates structural arrangement and revision

They share the same overall data model, but they operate over different concerns.

## Scope Of This Document

This document defines Jarvis at the role level.

It does not define:

- the full Stark/Jarvis collaboration contract
- the exact actor-routing policy Jarvis should use
- the exact runtime primitive menu Jarvis may call
- the final policy for how Jarvis influences recall versus act

## Short Framing

Jarvis is the canonical user-facing orchestrator actor.

It tends continuity, attention, and what matters from the user-facing side of the system while coordinating other actors when needed.
