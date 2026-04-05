# Entrypoint And Onboarding (Prelim)

## Purpose

This document preserves the current working direction for:

- the user-facing entrypoint
- first-run onboarding
- local discovery
- initial probing
- early registration

This is not a stable contract.

## Current Posture

The system should feel like one product with one obvious entrypoint.

The current working entrypoint is:

```text
skyra init
```

That entrypoint should own first-run onboarding.

The system should not require the user to think about:

- separate installer modes
- manual network setup
- daemon bootstrapping details
- local database setup

## Core Idea

The current design direction is:

- one binary
- one obvious entrypoint
- one onboarding flow that discovers whether the device should join or host

The important point is that the user should begin with `skyra init`, not with a menu of install variants.

## Current First-Run Direction

The current first-run direction is:

1. user runs `skyra init`
2. Skyra checks whether it can discover an existing local Skyra core
3. if one is found, the user may join it
4. if none is found, the user may start a new local instance
5. permissions are checked and surfaced
6. local or proxied probing runs
7. the subject is registered
8. verified capability surfaces are published
9. relevant steward or worker actors are born under contract

## Short Framing

`skyra init` should discover, branch, and onboard.

Packaging details may vary later, but they should not define the architecture.
