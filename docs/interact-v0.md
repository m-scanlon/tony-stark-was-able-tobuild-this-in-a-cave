# Interact v0

## Purpose

This document defines the current working meaning of `interact`.

`interact` is intended to be the broad world-facing primitive.

It is not the retained read path and not the retained write path.

## Core Framing

The current primitive split is:

- `recall`
- `learn`
- `interact`

Under that split:

- `recall` reads retained experience
- `learn` writes retained experience
- `interact` crosses a world boundary

That makes `interact` the primitive for operations that leave cognition and retention and touch something external.

## What Counts As World-Facing

The current direction is that `interact` includes actions such as:

- responding to a human
- probing a device
- searching the web
- calling an external service
- writing registration state
- invoking a capability surface

The common property is:

- the node is crossing an external boundary

## What `interact` Is Not

`interact` should not absorb everything.

It should not replace:

- `recall`
- `learn`

Those stay distinct because they have different invariants and different write/read boundaries.

## Working Shape

The current likely working shape is:

```text
skyra <node> interact -method <method> ... -reason "<why this command is being emitted>"
```

Examples:

```text
skyra jarvis interact -method respond -target human -reason "the user needs a response"
```

```text
skyra stark interact -method probe -subject_id laptop -reason "the device needs capability discovery"
```

```text
skyra stark interact -method search -query "roku ecp api" -reason "the current device probe needs vendor API confirmation"
```

```text
skyra stark interact -method write_device_registration -subject_id laptop -reason "verified capability state must be persisted"
```

## Methods

The method slot is where world-facing specialization should live.

Current plausible methods include:

- `respond`
- `probe`
- `search`
- `call_api`
- `write_device_registration`
- `use_capability`

The exact taxonomy is not frozen.

The important current point is:

- `interact` is broad at the primitive level
- methods carry the operational specialization

## `channel` Is Still Open

Whether `interact` should also require:

- `-channel <channel>`

remains open.

That may become useful later, but it should not be frozen prematurely.

For now, the stable design move is:

- keep `-method`
- leave `channel` unresolved

## Validation Shape

The likely validation ladder is:

1. node is allowed to emit `interact`
2. the selected `method` is allowed for that node
3. the required arguments for that method are present
4. runtime validates the external boundary and permissions
5. execution result is written back into episode state

This keeps `interact` broad without making it loose.

## Relationship To Registration

Registration should now be thought of as one world-facing interaction path rather than as a separate top-level primitive.

That means:

- registration is likely an `interact` method

The same is true for:

- probing
- search
- human response

This does not trivialize them.

It groups them under one shared primitive boundary:

- external interaction

## Current Design Posture

The strongest current claims are:

- `interact` should be the broad world-facing primitive
- `interact` should not replace `recall` or `learn`
- method-specific specialization should live inside `interact`
- `channel` remains open

## Short Framing

`interact` is the primitive for crossing out of cognition and retention into the world.

Its likely shape is:

```text
skyra <node> interact -method <method> ... -reason "..."
```

The exact `method` taxonomy is still open.

`channel` is still an unresolved question.
