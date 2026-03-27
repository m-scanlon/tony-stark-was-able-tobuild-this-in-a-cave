# Node Birth v0

## Core Framing

Node birth is the process by which a node comes into existence as a live runtime operator under a contract.

The contract is the birth spec.

There is no separate birth-spec object in the current model.

The node is born live under the contract that created it.

## Birth Authority

The kernel is the birth authority.

The kernel is responsible for actual node instantiation.

Nodes do not self-birth.

Stark does not directly instantiate other nodes.

Instead:

- Stark publishes a contract
- the kernel passes that contract into the node factory
- the node factory instantiates the node

## Stark Exception

`Stark` is the bootstrap exception.

At system start:

1. the kernel starts
2. the kernel uses Stark's hardcoded contract
3. the kernel sends that contract to the node factory
4. Stark is instantiated

This is the one intentional special case in the current model.

After Stark exists, later node creation becomes regular.

## Regular Node Birth

For regular nodes, the flow is:

1. Stark publishes a node contract
2. the kernel receives that contract
3. the kernel passes the contract to the node factory
4. the node factory instantiates the node
5. the node becomes live immediately under that contract

## Contract As Birth Spec

The node contract itself is the birth spec.

At minimum, it currently provides:

- purpose
- accepted stimulus boundary
- allowed interact boundary

This means node birth does not currently require separate birth-only metadata for:

- capability state
- loop envelope
- node role

Those are not part of the active birth contract in this model.

## Live On Instantiation

There is no separate activation gap after birth.

Once the node is instantiated by the node factory, it is live under the contract that created it.

That means:

- the contract is active immediately
- the node may begin accepting valid events immediately

## Initial Runtime State

A newly born node should start with:

- `node_id`
- active contract snapshot
- no active episode yet
- no pending command state yet

The node is alive but idle until the first valid event arrives.

## Relationship To Episodes

Node birth does not automatically open an episode.

The node lifetime and the episode lifetime are different.

The node is durable.

The episode is bounded.

So the preferred model is:

- birth node first
- open the first episode only when the first valid event arrives

## Relationship To Kernel And Stark

The split is:

- kernel = instantiation authority
- Stark = publisher of later node contracts
- node factory = constructor path used by the kernel

This keeps:

- publication
- instantiation
- runtime execution

as separate responsibilities.

## Current Design Posture

The strongest current claims are:

- the contract is the birth spec
- the kernel is the node birth authority
- Stark is the bootstrap exception at system start
- Stark publishes contracts for later nodes, but the kernel still performs instantiation
- nodes are live immediately on instantiation
- node birth does not automatically open an episode

## Short Framing

The kernel births Stark from a hardcoded contract at startup.

After that, Stark publishes contracts for later nodes, and the kernel births those nodes through the node factory.

The contract is the birth spec.

Nodes become live immediately on instantiation and wait idle until the first valid event opens an episode.
