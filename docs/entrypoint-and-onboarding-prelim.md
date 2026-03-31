# Entrypoint And Onboarding (Prelim)

## Purpose

This document preserves the current working direction for:

- the user-facing entrypoint
- first-run onboarding
- local discovery
- initial probing
- early registration

This is not a stable contract.

It is a snapshot of where the design currently is so the system can keep iterating without losing the thread.

## Current Posture

The system should feel like one product with one obvious entrypoint.

The current working entrypoint is:

```text
skyra init
```

That command should own first-run onboarding.

The system should not require the user to think about:

- separate installer modes
- manual network setup
- daemon bootstrapping details
- local database setup

at the conceptual level.

Those implementation details may exist underneath, but the user-facing entrypoint should stay simple.

## Core Idea

The current design direction is:

- one binary
- one obvious entrypoint
- two broad runtime roles

The rough roles are:

- a device may become a local Skyra instance
- a device may join an existing Skyra network

This is still fluid.

The important point is that the user should begin with `skyra init`, not with a menu of install variants.

## Why This Matters

The hard problem is not packaging.

The hard problem is:

- trust establishment
- local discovery
- permission handling
- probing
- registration
- capability publication

So the product surface should center the onboarding flow, not the packaging shape.

Packaging may later include:

- Homebrew
- tarballs
- pkg installers

but those are distribution layers rather than the architectural core.

## Current First-Run Direction

The current first-run direction is:

1. user runs `skyra init`
2. Skyra checks whether it can discover an existing local Skyra core
3. if one is found, the user may join it
4. if none is found, the user may start a new local instance
5. permissions are checked and surfaced
6. local or proxied probing runs
7. the subject is registered
8. verified capabilities are published

This is a direction, not a locked final wizard.

## Local Discovery

The current same-network discovery direction is:

- use local service discovery first
- likely mDNS in `v0`

The rough idea is:

- a Skyra core advertises itself on the local network
- a new instance browses for that service during onboarding
- if found, onboarding can branch into a join flow
- if not found, onboarding can branch into a create flow

This keeps `v0` constrained to the local network rather than trying to solve internet-wide onboarding too early.

## Current Branches

The current onboarding branches are:

### 1. Join Existing Skyra

If a local Skyra core is found, the tentative flow is:

1. present the discovered instance
2. ask whether the user wants to join it
3. establish trust and authentication with that core
4. check local permissions
5. run local probing
6. report the result to `Stark`
7. write device registration
8. publish verified capabilities

### 2. Start New Local Skyra

If no local Skyra core is found, the tentative flow is:

1. ask whether this device should become a new local Skyra instance
2. initialize local state
3. start the kernel/runtime
4. start `Stark`
5. check local permissions
6. probe the local device
7. register the local subject
8. publish verified capabilities
9. advertise the instance on the local network

These two branches are still conceptual.

The stable idea is just:

- `skyra init` should discover, branch, and onboard

## Local State

The current working local state direction is:

```text
~/.skyra/
  skyra.db
  config.yml
  logs/
```

This is not yet frozen.

The important current idea is:

- local durable state should live in one predictable home

That state likely includes:

- retained experience storage
- configuration
- logs
- onboarding metadata

## Entry Boundary vs Packaging Boundary

The current design should distinguish between:

- entry boundary
- packaging boundary

The entry boundary is:

- `skyra init`

The packaging boundary is still open.

Possible packaging later includes:

- `brew install skyra`
- a tarball with the binary
- later installer packages if needed

But the product should not let packaging details decide the runtime model too early.

## Probing Direction

The current probing direction after onboarding is:

```text
transport.attach
-> transport.reveal
-> os.inference
-> probe.strategy.select
-> probe.execute
-> registration.write
-> capability publication
```

Important current posture:

- probing is staged
- probing is not the same thing as daemon install
- OS family may be inferred, not directly read in a universal way
- different device classes may require different probe strategies

## Smart Devices vs Constrained Devices

The current direction distinguishes between broad device classes.

### Self-Hosted Or Smart Devices

Examples:

- macOS laptop
- Linux box
- Raspberry Pi
- later some Android paths

These are the easiest environments because Skyra can more plausibly run local runtime code on them.

### Constrained Or Externally Probed Devices

Examples:

- speakers
- displays
- drives
- appliances
- some mobile devices

These may need:

- network probing
- transport probing
- or a nearby probe agent such as a Pi

The system should not pretend all devices can self-host a daemon.

## Phone Reality

The current design should explicitly acknowledge that phones are difficult.

Broadly:

- Android is more tractable
- iPhone is much more constrained

This means the system should avoid assuming:

- every phone can run the same local runtime path
- every device can become a symmetric Skyra node

The registration layer should therefore stay uniform even when the onboarding and probing paths differ substantially.

## Registration Direction

The current registration backbone is:

- `subject`
- `transport`
- `probe_strategy`
- `verified_capabilities`

See also:

- [device-registration-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/device-registration-v0.md)

This matters because onboarding should end in a strict registration shape even if the road to get there varies by device.

## Permissions

The current design direction is:

- permissions should be surfaced during onboarding
- missing permissions should be explained clearly
- the system may deep-link to OS settings where that is possible

This is especially important on macOS, where some permissions require explicit user action.

The system should not pretend it can grant those permissions automatically if the OS does not allow that.

## What Is Not Stable Yet

The following remain open:

- exact daemon lifecycle model
- exact trust/auth handshake
- exact mDNS advertisement format
- exact onboarding prompts and UX
- exact packaging route
- exact phone strategy
- exact multi-device network topology
- exact persistent config schema

This document should therefore be treated as preserved direction, not frozen architecture.

## Current Design Posture

The strongest current claims are:

- the main user-facing entrypoint should be `skyra init`
- onboarding should discover an existing local Skyra network before creating a new one
- probing should be staged and strategy-driven
- registration should end in a strict outer envelope even if device-specific paths differ
- packaging should be treated as a later distribution concern rather than the main design driver

## Short Framing

The current direction is to make `skyra init` the single obvious entrypoint.

That entrypoint should handle discovery, branching, permissions, probing, and initial registration.

The design is still early, but the main principle is already useful:

- keep the user-facing entry simple
- keep probing adaptive
- keep registration strict
- keep packaging secondary for now
