# Runtime Primitives And Artifacts (Prelim)

## Core Framing

The system needs a clear distinction between:

- retained artifacts
- runtime primitives
- runtime artifacts

These are not the same layer.

Retained artifacts belong to learning and longer-lived retained experience.

Runtime commands and runtime artifacts belong to runtime execution inside an active episode.

## The Split

The current working split is:

- retained artifacts are long-lived outputs of learning
- runtime commands are callable operations the node can invoke during runtime execution
- runtime artifacts are transient outputs produced by those commands inside an episode

This means the system should not treat every runtime result as retained memory.

Most runtime execution should remain episode-local unless later learning decides it should be retained.

## Runtime Command

A runtime command is a callable operation that the node may issue during an active episode.

Conceptually:

```text
skyra <command_set> <command> -<args>
```

The command is emitted by the node.

The kernel receives that command and remains the authority over what happens next.

That includes:

- user-facing output commands
- capability or API commands
- commands that request another reasoning step

## Kernel Role

The kernel is the organizer and authority for runtime primitive execution.

At a high level, the flow is:

1. the node sees the current frame
2. inference selects an allowed runtime command
3. the node emits that command
4. the kernel validates and dispatches it
5. the command executes
6. the resulting runtime artifact or interaction-relevant output is written back into episode state
7. the node may later project another frame and choose the next command

This means:

- the node chooses the next operation
- the kernel controls execution and writeback
- runtime execution unfolds as a sequence of bounded operations

## Runtime Artifacts

A runtime artifact is a transient output produced during runtime execution inside an episode.

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

- runtime execution may produce many transient artifacts
- learning later decides whether any of that becomes retained experience

## Interpret

`interpret` is a strong candidate for a runtime command.

Its role would be:

- help the node make sense of confusing, incomplete, or ambiguous stimulus
- produce a transient runtime artifact
- shape later inference and episode-local state

One later valid direction is that the node may emit a command that requests another prompt or reasoning step rather than "thinking by itself" outside the command model.

It should not create retained understanding by default.

Instead:

- `interpret` creates runtime-local output
- later learning may decide that some interpretation contributed to retained understanding, salience, tension, or trace formation

## Learn

`learn` is a strong candidate for the learning handoff command.

Its role would be:

- take a just-closed episode as input
- run the learning / consolidation write path
- produce retained artifact and structure updates

A good current working shape is:

```text
skyra primitive learn -episode_id <episode_id>
```

This command should be understood as the kickoff into the learning path after episode closure rather than as ordinary in-episode state mutation.

## Runtime vs Retained Artifacts

The current working distinction is:

- runtime artifact = transient runtime output inside an episode
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

- command routing
- command validation
- command dispatch
- command-result writeback
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
- runtime artifacts produced during runtime execution
- interaction outputs

to decide what should become retained experience.

When learning persists retained artifacts, the same write path should also update the anchor lookup layer used later by recall.

That learning result may still return a typed write receipt containing the newly written artifact ids.

The node may use that receipt to seed a small transient recent-learned cache if useful, but the durable owner of retained experience remains the retention/index layer rather than the node itself.

## Current Design Posture

The strongest current claims are:

- callable runtime commands are distinct from retained artifacts
- runtime artifacts are transient outputs of runtime commands
- the kernel is the authority over command execution and writeback
- runtime artifacts should be available to later learning
- `interpret` is a good candidate for a runtime command
- `interpret` should not create retained understanding by default
- `learn` is a good candidate for the post-episode learning kickoff command

The exact command surface remains open.

## Short Framing

The system should distinguish between runtime execution and retained experience.

Runtime commands are callable in-episode operations.

Runtime artifacts are transient outputs of those operations.

Retained artifacts are longer-lived outputs of learning.

The kernel remains the authority that validates command execution, routes writeback, and organizes runtime progression.
