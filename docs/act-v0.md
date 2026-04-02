# Act v0

## Purpose

This document defines the current working meaning of `act`.

`act` is the broad world-facing primitive.

It is not the retained read path and not the retained write path.

## Core Framing

The current primitive split is:

- `recall`
- `learn`
- `observe`
- `act`

Under that split:

- `recall` reads retained experience
- `learn` writes retained experience
- `observe` takes something in from the world
- `act` is any intentional boundary engagement that produces an effect or signal

That makes `act` the primitive for operations that leave cognition and retention and touch something external.

## What Counts As World-Facing

The current direction is that `act` includes actions such as:

- responding to a human
- probing a device
- searching the web
- calling an external service
- writing registration state
- invoking a capability surface

The common property is:

- the node is crossing an external boundary

## What `act` Is Not

`act` should not absorb everything.

It should not replace:

- `recall`
- `learn`
- `observe`

Those stay distinct because they have different invariants and different write/read boundaries.

## Working Shape

The current locked shape is:

```text
skyra <node> act \
  -target <target> \
  -content <content> \
  -modality <modality> \
  -timestamp <timestamp> \
  -reason "<why this command is being emitted>"
```

Examples:

```text
skyra jarvis act -target human -content "the user needs a response" -modality text -timestamp now -reason "the user needs a response"
```

```text
skyra stark act -target laptop -content "discover capability surface" -modality probe -timestamp now -reason "the device needs capability discovery"
```

```text
skyra stark act -target laptop -content "vendor API confirmation" -modality search -timestamp now -reason "the current device probe needs vendor API confirmation"
```

```text
skyra stark act -target laptop -content "write verified device registration" -modality registration_write -timestamp now -reason "verified capability state must be persisted"
```

## Fields

The current `act` fields are:

- `target`
- `content`
- `modality`
- `timestamp`
- `reason`

Their current meanings are:

- `target` = who or what the action is directed at
- `content` = what is being conveyed or carried out
- `modality` = how the action is delivered
- `timestamp` = when the action is going to happen
- `reason` = why it was done

`reason` is always required by the system.

## Modalities

`modality` is where world-facing specialization now lives.

Current plausible modalities include:

- `text`
- `speech`
- `probe`
- `search`
- `api_call`
- `registration_write`
- `capability_use`

The exact vocabulary is not frozen.

The important current point is:

- `act` is broad at the primitive level
- `modality` carries the operational specialization

## Jarvis vs Stark

The top-level primitive stays shared:

- `act`

What differs between major node roles is the allowed modality subset.

### Jarvis

For `v1`, `Jarvis` should primarily own human-facing `act`.

The important current boundary is:

- outward human-facing `act` belongs to `Jarvis`

### Stark

For `v1`, `Stark` should primarily own system-facing `act`.

That currently includes modalities such as:

- `probe`
- `registration_write`

This keeps `Stark` focused on structural and system-facing action rather than user-facing response.

`birth_node` should not be treated as an `act` modality.

It belongs to Stark's structural command surface rather than to the world-facing `act` primitive.

## Same Primitive, Different Allowances

This should be understood as:

- one shared world-facing primitive
- different node-specific modality allowances under contract

So the differentiation is not:

- `Jarvis` gets one primitive
- `Stark` gets another primitive

It is:

- both use `act`
- their allowed world-facing forms differ by role

## Validation Shape

The likely validation ladder is:

1. node is allowed to emit `act`
2. the selected target is allowed
3. the selected `modality` is allowed for that node or delegation edge
4. the required `act` fields are present
5. runtime validates the external boundary and permissions
6. execution result is written back into episode state

This keeps `act` broad without making it loose.

## Relationship To Registration

Registration should now be thought of as one world-facing action path rather than as a separate top-level primitive.

That means:

- registration is a valid `act` form

The same is true for:

- probing
- search
- human response

This does not trivialize them.

It groups them under one shared primitive boundary:

- intentional world-facing action

## Current Design Posture

The strongest current claims are:

- `act` should be the broad world-facing primitive
- `act` should not replace `recall`, `learn`, or `observe`
- `target`, `content`, `modality`, `timestamp`, and `reason` are the active `act` shape
- specialization should live inside `modality`

## Short Framing

`act` is the primitive for crossing out of cognition and retention into the world.

Its current shape is:

```text
skyra <node> act -target <target> -content <content> -modality <modality> -timestamp <timestamp> -reason "..."
```

The exact `modality` vocabulary is still open.
