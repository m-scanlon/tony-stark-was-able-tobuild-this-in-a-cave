# Runtime Primitives And Artifacts (Prelim)

## Core Framing

The system needs a clear distinction between:

- retained artifacts
- runtime primitives
- runtime artifacts

These are not the same layer.

Retained artifacts belong to learning and longer-lived retained experience.

Runtime primitives and runtime artifacts belong to cognition inside an active episode.

## The Split

The current working split is:

- retained artifacts are long-lived outputs of learning
- runtime primitives are callable operations the node can invoke during cognition
- runtime artifacts are transient outputs produced by those primitives inside an episode

This means the system should not treat every cognitive result as retained memory.

Most cognition should remain runtime-local unless later learning decides it should be retained.

## Runtime Primitive

A runtime primitive is a callable operation that the node may issue during an active episode.

Conceptually:

```text
skyra <primitive> -<args>
```

The primitive command is emitted by the node.

The kernel receives that command and remains the authority over what happens next.

## Kernel Role

The kernel is the organizer and authority for runtime primitive execution.

At a high level, the flow is:

1. the node sees the current frame
2. the node emits a primitive command
3. the kernel validates that command
4. the kernel applies the primitive's frame template
5. the primitive executes
6. the resulting runtime artifact or computation output is written into the frame
7. the node chooses the next primitive

This means:

- the node chooses the next operation
- the kernel controls execution and frame transition
- runtime cognition unfolds as a sequence of bounded operations

## Runtime Artifacts

A runtime artifact is a transient output produced during cognition inside an episode.

Runtime artifacts are:

- in-episode
- temporary
- usable by later runtime steps
- available to later learning
- not retained by default

The important distinction is:

- runtime artifact first
- retained artifact only if later learning selects it

## Why This Split Matters

Without this split, the system collapses:

- thinking
- memory formation
- retained experience

into one layer.

That creates confusion about what should exist only inside an active episode versus what should survive after the episode.

The intended rule is:

- runtime cognition may produce many transient artifacts
- learning later decides whether any of that becomes retained experience

## Interpret

`interpret` is a strong candidate for a runtime primitive.

Its role would be:

- help the node make sense of confusing, incomplete, or ambiguous stimulus
- produce a transient runtime artifact
- shape the current frame for later cognition

It should not create retained understanding by default.

Instead:

- `interpret` creates runtime-local output
- later learning may decide that some interpretation contributed to retained understanding, salience, tension, or trace formation

## Runtime vs Retained Artifacts

The current working distinction is:

- runtime artifact = transient cognitive output inside an episode
- retained artifact = long-lived retained experience selected by learning

Retained artifacts currently include:

- `retained_trace`
- `retained_understanding`
- `retained_salience`
- `retained_tension`

Runtime artifacts are not yet fully typed, but they are conceptually distinct from this retained family.

## Kernel Sections

This split also implies a cleaner kernel organization.

The kernel may eventually need separate sections or subsystems for:

- primitive command routing
- primitive validation
- frame template application
- bounded computation
- runtime artifact handling
- interaction handling
- learning handoff

This document does not define those subsystems in detail.

It only establishes that the split is real and architecturally important.

## Learning Handoff

Runtime artifacts should remain available to the later learning process.

That means learning may use:

- the episode record
- the episode's structural field
- runtime artifacts produced during cognition
- interaction outputs

to decide what should become retained experience.

## Current Design Posture

The strongest current claims are:

- callable runtime primitives are distinct from retained artifacts
- runtime artifacts are transient outputs of runtime primitives
- the kernel is the authority over primitive execution and frame transition
- runtime artifacts should be available to later learning
- `interpret` is a good candidate for a runtime primitive
- `interpret` should not create retained understanding by default

The exact primitive menu remains open.

## Short Framing

The system should distinguish between runtime cognition and retained experience.

Runtime primitives are callable in-episode operations.

Runtime artifacts are transient outputs of those operations.

Retained artifacts are longer-lived outputs of learning.

The kernel remains the authority that validates primitive execution, applies frame templates, and organizes runtime progression.
