# Actor Birth v0

## Core Framing

Actor birth is the process by which a actor comes into existence as a live runtime operator under a contract.

The contract is the birth spec.

There is no separate birth-spec object in the current model.

The actor is born live under the contract that created it.

## Birth Authority

The kernel is the birth authority.

The kernel is responsible for actual actor instantiation.

Actors do not self-birth.

Stark does not directly instantiate other actors.

Instead:

- Stark publishes a contract
- the kernel passes that contract into the actor factory
- the actor factory instantiates the actor

## Stark Exception

`Stark` is the bootstrap exception.

At system start:

1. the kernel starts
2. the kernel uses Stark's hardcoded contract
3. the kernel sends that contract to the actor factory
4. Stark is instantiated

## Regular Actor Birth

For regular actors, the flow is:

1. Stark publishes a actor contract
2. the kernel receives that contract
3. the kernel passes the contract to the actor factory
4. the actor factory instantiates the actor
5. the kernel registers the new `actor_id`
6. the actor becomes live immediately under that contract

## Contract As Birth Spec

The actor contract itself is the birth spec.

At minimum, the major runtime-relevant shape is:

- `purpose`
- `commitments`
- request stimuli
- response envelopes

That means actor birth is not centered on an old separate command or capability-allowance object.

It is centered on the public callable request/response surface the actor is born to handle.

## Live On Instantiation

There is no separate activation gap after birth.

Once the actor is instantiated by the actor factory, it is live under the contract that created it.

That means:

- the contract is active immediately
- the actor may begin accepting valid routed stimulus immediately

## Initial Runtime State

A newly born actor should start with:

- `actor_id`
- active contract snapshot
- no active episode yet
- no `dependencyLedger` state yet

The actor is alive but idle until the first valid routed stimulus arrives.

## Relationship To Episodes

Actor birth does not automatically open an episode.

The actor lifetime and the episode lifetime are different.

The actor is durable.

The episode is bounded.

So the preferred model is:

- birth actor first
- open the first episode only when the first valid routed stimulus arrives

## Current Design Posture

The strongest current claims are:

- the contract is the birth spec
- the kernel is the actor birth authority
- Stark publishes later actor contracts, but the kernel performs instantiation
- actors are live immediately on instantiation
- actor birth does not automatically open an episode

## Short Framing

The kernel births Stark from a hardcoded contract at startup.

After that, Stark publishes contracts for later actors and the kernel births those actors through the actor factory.

The contract is the birth spec, and the actor becomes live immediately under it.
