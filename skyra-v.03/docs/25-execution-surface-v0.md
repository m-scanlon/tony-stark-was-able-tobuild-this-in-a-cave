# Execution Surface v0

## Status

This document records emerging design direction from `2026-04-09`.

It is not yet locked canon.

It is open for pressure and revision.

## The Problem

The current model satisfies the question of who.

The kernel knows which being is the target.

The kernel knows which being is the source.

What the model does not yet answer cleanly is where.

In a distributed architecture — many machines, many inference providers, many
CLI binaries — the kernel needs to know not just which being to route to but
where that being executes.

## The Proposed Shape

Each being has a registered execution surface.

The execution surface is not part of nature.

Nature is identity and purpose. It does not change.

The execution surface is operational. It answers where this being runs right
now.

A being has two possible execution surface slots:

- one for inference — the endpoint where cognition happens for cognitive beings
- one for CLI — the binary or process that handles execution for non-cognitive beings

Either slot may be local or remote.

## HTTP As The Transport

Everything resolves over HTTP on the local network.

Not public. Local network only.

The kernel holds the registered execution surfaces for all beings.

Beings do not know where each other live.

The kernel resolves the surface at dispatch time.

When the target being is remote, the kernel posts the protocol string to that
being's execution surface over HTTP and receives a protocol string back.

The transport is invisible to the beings on either side.

They still speak the protocol.

The network hop is an implementation detail of the kernel's routing layer.

## What This Means For The API Layer

If everything flows through the protocol over HTTP, a separate API layer
becomes unnecessary.

The protocol string is the API.

Any machine that can receive a protocol string and emit one back is a valid
execution surface.

## What The Kernel Retains

Distributing execution does not distribute kernel authority.

The kernel retains sole ownership of:

- being registration
- relationship hashmap
- edge weight updates
- relationship emergence threshold logic
- trust values
- execution surface registry

None of those cross the network.

The only thing that crosses the network is the present going out and a protocol
string coming back.

Cognition happens in a being that lives on a different machine.

The kernel still routed it.

The kernel still registered it.

The kernel still owns the result.

## The Key Distinction

Cognition being executed outside the kernel's process is not the same as
cognition being outside the kernel's authority.

The kernel has never done cognition.

It routes.

Beings think.

Some beings happen to be far away.

That does not change what the kernel owns.

What would break the model is kernel authority crossing the network — a remote
being registering other beings, writing to the relationship hashmap, updating
edge weights.

That is the line.

## Open Questions

- Whether a being's execution surface can change after registration or is fixed
  at birth
- Whether the kernel needs a heartbeat or health check against remote execution
  surfaces
- Whether a being whose execution surface goes offline becomes dormant or is
  removed from the runtime
- The exact format of the execution surface field on the being record
- Whether local and remote beings need different trust origin values
