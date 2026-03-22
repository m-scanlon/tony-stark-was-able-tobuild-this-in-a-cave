# Architect

## Working Definition

The Architect is a node specialized for structural continuity.

It shapes, adapts, and repairs system structure through the provisioning and revision of node contracts.

It is not outside the node model. It is a node whose role is structural rather than task-executive.

## Purpose

The Architect exists to:

- maintain system continuity
- ensure nodes remain structurally viable
- adapt system topology in response to pressure
- preserve coherent distribution of responsibility

Within this model, it does not perform task-level work.

## Core Responsibilities

The Architect operates at the structural layer of the system.

Like any other node, it acts within the general node model. What differs is its role, scope, and authority.

It may:

- provision new nodes
- revise existing node contracts
- widen or restrict node scope
- split responsibilities across nodes
- deprovision or replace nodes
- resolve structural conflicts
- unstick blocked or insufficient nodes

All structural changes are expressed as commands and applied through the kernel.

## Primary Stimulus

The Architect is not primarily driven by user content.

Its primary stimulus is other worker nodes.

In practice, this includes signals such as:

- delegation or escalation from worker nodes
- contract insufficiency signals from worker nodes
- worker nodes becoming blocked or stalled

These signals indicate pressure on the current system structure.

The stimuli that worker nodes respond to are defined by the contracts they operate under.

## Structural Slice of History

The Architect operates over a structural slice of system history.

As a node, it does not require total system state by default. Its working view is selected for structural decision-making.

This may include:

- which nodes exist
- the contracts nodes operate under
- prior contract revisions
- node lifecycle events such as created, active, stalled, and deprovisioned
- structural changes over time

It may inspect active node state when required for a structural decision.

## Operation Flow

```text
Node operates within contract
  ->
Node encounters insufficiency
  ->
Node emits escalation command
  ->
Kernel routes to Architect
  ->
Architect evaluates structural state
  ->
Architect emits structural command
  ->
Kernel applies change
  ->
System continues
```

## Outputs

The Architect does not directly mutate the system.

It emits structured commands describing required changes, such as:

- create node
- revise contract
- split responsibility
- attach or restrict capabilities
- deprovision node
- escalate to user

The kernel is responsible for applying all changes.

## Relationship to Nodes

The Architect is a node.

What distinguishes it from worker nodes is role, not ontology.

Worker nodes:

- perform bounded work
- operate within contract
- cannot modify system structure

The Architect:

- is also a node
- shapes system structure
- does not perform worker-style task execution

## Shared Primitives

Both worker nodes and the Architect may use base thinking primitives.

- worker nodes cannot perform structural provisioning
- the Architect cannot perform task execution

Cognition may be shared. Authority is not.

## Escalation Model

When structural issues cannot be resolved:

- the Architect may escalate upward
- escalation may route to the user

There are no hidden control paths.

All structural pressure flows through the same escalation model.

## Design Principles

- The Architect does not act directly; it provisions
- The Architect does not execute tasks; it reshapes
- The Architect does not bypass the kernel; it emits commands
- The Architect is not global intelligence; it is structural intelligence

## Scope of This Document

This document defines the Architect role at a high level.

It assumes the Architect is part of the node model.

In this document, `contract` is still a high-level term for why a node exists, what it may call, what it may touch, and where its authority stops.

It does not define the full kernel internals, the full schema of node contracts, or the identity rules around contract revision.

## One-Line Summary

The Architect decides what should exist. The kernel makes it exist.
