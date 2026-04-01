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

- talking to a human
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
skyra jarvis interact -method talk -target human -reason "the user needs a response"
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

- `talk`
- `probe`
- `search`
- `call_api`
- `write_device_registration`

The exact taxonomy is not frozen.

The important current point is:

- `interact` is broad at the primitive level
- methods carry the operational specialization

For `v1`, the current preferred narrower set is:

- `talk`
- `probe`
- `write_device_registration`

Within that set, `talk` should initially be restricted to `Jarvis`.

That means:

- `Jarvis` is the only node allowed to emit `interact -method talk` in `v1`
- system-facing nodes should not emit human-facing talk directly
- system-facing interaction should remain focused on probing and registration

## Jarvis vs Stark

The top-level primitive stays shared:

- `interact`

What differs between major node roles is the allowed method subset.

### Jarvis

For `v1`, `Jarvis` should primarily own:

- `talk`

`Jarvis` may later gain additional user-facing interaction methods, but the important current boundary is:

- outward human talk belongs to `Jarvis`

### Stark

For `v1`, `Stark` should primarily own:

- `probe`
- `write_device_registration`

This keeps `Stark` focused on structural and system-facing interaction rather than user-facing talk.

`birth_node` should not be treated as an `interact` method.

It belongs to Stark's structural command surface rather than to the world-facing `interact` primitive.

## Same Primitive, Different Allowances

This should be understood as:

- one shared world-facing primitive
- different node-specific method allowances under contract

So the differentiation is not:

- `Jarvis` gets one primitive
- `Stark` gets another primitive

It is:

- both use `interact`
- their allowed sub-primitives differ by role

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
- human talk

This does not trivialize them.

It groups them under one shared primitive boundary:

- external interaction

## Relevant Memory Clarifications

While defining `interact`, the current memory posture became clearer:

- subject-scoped memory domains are not separate nodes
- subject-scoped memory domains are not external capabilities
- registration and capability facts are not memory domains

The current useful subject-scoped retained domains are:

- human subject:
  - `identity`
  - `preferences`
  - `boundaries`
  - `interaction_style`
- system subject:
  - `identity`
  - `constraints`
  - `health`
  - `surface_behavior`

These clarifications matter to `interact` because many `interact` results may later contribute to retained experience inside those domains, but the domains themselves are not part of the `interact` method taxonomy.

For `v1`, learning should still remain episode-bounded.

That means:

- `interact` results may be learned from the local closed episode
- deep orchestration ancestry or stack-trace learning is not required for `v1`

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
