# Command Audit Trail v0

## Purpose

This document defines the audit requirement on command emission.

The main rule is simple:

- every emitted command must carry an explicit reason

## Core Rule

The working command protocol is:

```text
skyra <node> <primitive> -<args> -reason "<why this command is being emitted>"
```

`-reason` is mandatory.

It is not optional commentary.

It is part of the command contract.

## Why `-reason` Is Mandatory

The command surface is the system's execution boundary.

If commands are the boundary, then the system needs an inspectable record of:

- what node emitted the command
- what command was emitted
- what arguments it carried
- why the node believed that command should happen

Without mandatory rationale, the audit trail becomes a reconstruction exercise.

That is weaker than capturing the emitted reason at command time.

## What `-reason` Means

`-reason` should be understood as:

- the node's stated rationale for emitting the command now

It is:

- required
- inspectable
- part of the audit record

It is not:

- proof that the command is correct
- proof that the command is allowed
- proof that the command will succeed

Those are still determined by runtime validation and execution results.

## Runtime Rule

Runtime should reject commands that:

- omit `-reason`
- provide an empty `-reason`
- provide a structurally malformed `-reason`

At minimum, kernel validation should treat missing rationale as an invalid command invocation.

## Relationship To Validation

The split is:

- `reason` = why the node emitted the command
- validation = whether the command is allowed
- execution result = what actually happened

This keeps the audit trail honest.

The node must state its reason, but the system does not confuse that reason with proof.

## Examples

```text
skyra jarvis interact -method respond -target human -reason "the current frame requires a user-facing response"
```

```text
skyra jarvis recall -entity terraform -top_k 8 -reason "the stimulus explicitly introduced terraform as the active structural cue"
```

```text
skyra stark interact -method write_device_registration -subject_id Michaels-MacBook-Pro-10.local -reason "verified capability state must be persisted for this subject"
```

## Current Design Posture

The strongest current claims are:

- all emitted commands must carry `-reason`
- `-reason` is part of the audit trail
- kernel validation should reject commands without it
- `-reason` is rationale, not proof

## Short Framing

Every emitted command must include `-reason`.

That rationale is mandatory because command emission is part of the system's audit trail.
