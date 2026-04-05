# Protocol Surface (Prelim)

## Core Framing

This document preserves the shift away from older `namespace` and `command_set` thinking.

The active direction is no longer:

- large command families
- namespace-heavy top-level grammar
- direct capability-call framing

The active direction is:

- a small primitive set
- a primitive-first outer protocol with explicit actor and surface address
- published stimulus contracts
- concrete runtime payloads routed by the kernel

## Current Active Direction

The current working outer protocol is:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

This means:

- the primitive slot names the current boundary mode
- the actor slot names the owning public actor
- the surface slot names the targeted callable public surface
- the payload slot carries the concrete stimulus payload that conforms to a published contract

## Why This Shift Happened

The newer direction comes from three pressures:

- typed stimulus became the real routing surface
- execution surfaces became explicit
- request/response contracts started looking more like APIs than freeform commands

That makes a primitive-first stimulus protocol with explicit actor and surface address a better fit than the older command-surface framing.

## Primitive Split

The current working primitive split remains:

- `recall`
- `learn`
- `observe`
- `act`

These primitives still matter.

What changed is the runtime surface they carry.

## Relationship To Contracts

Public runtime use should now be understood through published request/response stimulus contracts.

That means:

- the registry holds the contract
- runtime emission carries concrete payload
- the kernel validates payloads against the published contract

This is a cleaner fit than treating the public runtime surface as an ever-growing command taxonomy.

## Current Design Posture

The strongest current claims are:

- the old command-surface framing is no longer primary
- public callable surfaces should be modeled as request/response stimulus contracts
- the outer Skyra protocol should stay small
- actor-first routing and contract lookup should do most of the real work

## Short Framing

The active runtime surface is best understood as:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

The older namespace/command-surface idea is now historical.
