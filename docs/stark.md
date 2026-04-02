# Stark

## Working Definition

`Stark` is the canonical name for the structural orchestrator node.

References to the older `Architect` and `system node` concepts should collapse into `Stark`.

Functional subtitle:

- `Stark` = structural orchestrator node

Stark is not outside the node model.

It is a node whose role is structural coordination rather than ordinary specialized task execution.

## Purpose

Stark exists to maintain and reshape system structure.

This includes:

- node topology
- node registry
- stimulus registry
- node contracts
- node contract publication
- capability attachment
- distribution of responsibility
- routing structural work to other nodes when needed

Stark is the node concerned with what should exist structurally, not with performing all ordinary structural task work itself.

## Same Model As Other Nodes

Stark uses the same broad data model as other nodes.

That means Stark still participates through:

- bounded episodes
- recall
- runtime execution
- retained experience

What differs is role and authority, not ontology.

The current important clarification is:

- Stark is not the base worker layer
- Stark is an orchestrator node

## Primary Stimulus

Stark is primarily driven by structural pressure rather than ordinary user content.

Typical stimuli include:

- contract insufficiency
- capability gaps
- topology problems
- node blockage or structural mismatch
- requests for structural revision
- typed structural results returned from other nodes

For device/bootstrap flows, Stark should also be understood as the node that creates or selects the structural type of incoming device-side stimuli.

That means Stark is not only a later classifier.

It is the structural authority that turns raw or generic device/bootstrap packages into typed structural stimuli the rest of the system can route and act on.

It should also be understood as the owner of the live stimulus registry for the runtime.

## Outputs

Stark does not directly mutate the system.

It coordinates structural work and emits valid structural commands that the kernel can apply, such as:

- `birth_node`
- revise contract
- attach capability
- restrict scope
- replace node
- deprovision node

After bootstrap, Stark is also the publisher of later node contracts.

Stark should also be understood as the structural owner of the node registry.

That registry is the live structural record of which nodes exist and where routed packages addressed by `node_id` should go.

The kernel remains the authority that applies those changes.

For `v1`, the important current `act` boundary is that Stark should remain system-facing.

That means Stark should primarily own modalities such as:

- `probe`
- `registration_write`

Stark should not emit outward human-facing `act` in `v1`.

By contrast, `birth_node` belongs to Stark's structural command surface rather than to `act`.

## Bootstrap Role

`Stark` is the bootstrap structural node.

The kernel births Stark at system startup from a hardcoded contract.

After that point, Stark becomes the structural origin for later node contract publication and structural orchestration.

## Structural History

Stark may operate over a structural slice of history, such as:

- which nodes exist
- which contracts exist
- contract revisions over time
- node lifecycle state
- capability layout

Stark does not require total system state by default.

Its view should be selected for structural decision-making.

## Design Principles

- Stark is a node, not a hidden control plane
- Stark is an orchestrator node, not a base worker node
- Stark shapes structure; it does not bypass the kernel
- Stark uses the same general runtime model as other nodes
- Stark is the canonical structural orchestrator node

## Relationship To Other Nodes

Other nodes perform bounded work within contract.

Stark governs structure around that work and coordinates structural delegation.

So:

- base worker nodes handle specialized task participation
- Stark handles structural continuity, routing, and revision

## Scope Of This Document

This document defines Stark at the role level.

It does not define:

- the full contract revision protocol
- the exact structural command schema
- the full Stark/Jarvis collaboration contract
- the exact structural routing policy Stark should use

## Short Framing

Stark is the canonical structural orchestrator node.

It decides what system structure should exist and coordinates the node graph around that decision.

It also publishes later node contracts.

The kernel makes those structural decisions executable.
