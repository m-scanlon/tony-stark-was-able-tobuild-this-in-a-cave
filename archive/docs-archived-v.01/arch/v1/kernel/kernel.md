# Kernel (v1)

## Overview

The kernel is Skyra's standalone execution service.

It receives promoted work from the ingress path, orders work in a max heap, and advances the runtime one reasoning or interaction step at a time.

In v1, the kernel is not organized around the later orchestration envelope model.

In v1, the native primitives are a protocol language for runtime control. The later skill system belongs to v2-v3 and is not important here.

## Core Flow

```text
stimulus
  -> ingress
  -> stimulus stream append
  -> Nexus
  -> max heap enqueue
  -> kernel pop
  -> Chain of Thought
  -> understand | interact
```

## Queue Contract

The kernel queue is a max heap ordered by numeric priority.

Current v1 policy:

- user messages: `100`
- internal chain-of-thought work: `50`

Higher numeric priority runs first.

When a user message arrives, the kernel must interrupt active internal reasoning and inspect that message first.

That interrupt does not force an outward reply.

## Runtime Boundaries

The kernel coordinates two distinct runtime boundaries:

- Chain of Thought
- Human-to-Machine Interaction

Chain of Thought handles cognition.

Human-to-Machine Interaction handles outbound communication and machine-facing action.

The kernel does not collapse those boundaries together.

## Primitive Language

The canonical top-level primitive menu for v1 is:

- `understand`
- `interact`

The canonical interpret flow is:

```text
understand
  -> interpret
  -> reference
  -> infer
  -> reference | resolve
```

`resolve` is the only step that creates `understanding`.

## Frame Contract

The active frame always contains `perception`.

`perception` always contains:

- `history`
- `stimulus`

In v1 there is exactly one mutable perception object for the runtime.

The fields inside that perception may change over time, but the runtime does not create multiple competing perceptions.

`understanding` is absent until the interpret cycle finishes with `resolve`.

Primitive choices are runtime outputs. They do not belong inside the frame.

## Interrupt And Resume

If internal reasoning is active and a user message arrives:

- the kernel interrupts the active model step
- the active chain becomes `suspended`
- the singleton perception shifts attention to the user message
- the kernel evaluates that user message through the normal primitive flow

After that user review:

- Skyra may emit `interact`
- Skyra may choose not to respond outwardly

If there is a suspended chain, the kernel resumes it in `rebase` mode.

In `rebase` mode:

- the same singleton perception object is retained
- latest `history` and latest `understanding` stay in perception
- the suspended chain's prior focus is restored as `stimulus`
- transient `reference` and `infer` outputs are cleared
- the resumed chain restarts at `reference`

`suspended` and `stale` are not the same thing.

`suspended` means live reasoning state that can resume.

`stale` means an obsolete queued event that the state machine should ignore.

## Model Memory Rule

The model should start with fresh memory every time it is invoked.

The only intentional carry-over inside the live runtime is the active frame and its `perception`.

No hidden long-lived model memory is part of the v1 contract.

## Persistence Boundary

v1 runtime state is in-memory only.

The kernel does not need to preserve queue state, frame state, thought logs, or interaction state across restarts.

## Not Canonical For v1

The following ideas are future-facing and not important for the current v1 kernel:

- skills as the main execution abstraction
- registry-gated skill execution
- later orchestration instances
- transport envelopes
- orchestration trees
- restart persistence of runtime state
