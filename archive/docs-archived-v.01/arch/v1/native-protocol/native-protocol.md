# Native Protocol (v1)

This protocol is the object-oriented representation of system execution.
The kernel is the interpreter.

This file defines the executable contract for native command syntax, primitive interfaces, and kernel enforcement rules.

## Object Model

1. `skill` = class
2. `job` = object instance of a skill
3. nested skill = callable composition target declared inside a parent skill
4. tasks do not exist as a protocol primitive

Everything composes through skills calling skills.

## Command Grammar

Canonical shape:

```bash
skyra [global_flags] <skill>[.<method>] [local_flags]
```

Root and nested forms:

```bash
skyra [global_flags] <skill> [local_flags]
skyra [global_flags] <root_skill>.<nested_skill> [local_flags]
```

Flag scope:

1. Global flags apply to the whole call tree unless policy blocks inheritance.
2. Local flags apply to the current call only.
3. Primitive implementation choice is local (`--impl=...`), not implicit global inheritance.

## Top-Level Execution Primitives

The top-level process model is defined in `thought-processes.md`.

The canonical top-level execution primitive set is:

- `stimuli`
- `process`
- `interact`

A primitive must participate in execution.
At this layer:

- `stimuli` enters execution
- `process` transforms execution
- `interact` emits execution results

`process` can contain child forms such as `thought`.
`thought` can contain child forms such as `repair`.
`interact` can contain child forms such as `agree` and `reply`.
`agree` may be carried by `reply`.

The native runtime interfaces in this file are secondary primitives that sit
under that top-level model.

## Secondary Primitive Interface Policy

Secondary primitives are interfaces, not single hardcoded implementations.

Policy:

1. One primitive command contract, many internal implementations.
2. `--impl=auto` is default unless policy requires explicit implementation.
3. Kernel validates selected implementation against per-primitive allowlist.
4. All implementations of a primitive must satisfy the same output contract.
5. Depth/cycle/context/byte/time limits apply identically across implementations.
6. Selected implementation is recorded in invocation metadata for replay/debug.

## Non-Primitive System Capabilities

- `synthesize` — session-to-graph synthesis capability (not primitive).
  Spec: `docs/arch/v1/native-protocol/non-priimitives/synthesize.md`
- `conversation` — user interaction capability (not primitive).
- `internet_search` — external web lookup capability (not primitive).
  Spec: `docs/arch/v1/native-protocol/non-priimitives/internet-search.md`
- `shard_registration` — system encapsulation capability (classification non-primitive).

## Constructor Contract

A skill constructor is the context actor set used at skill start.

Constructor input encoding:

```bash
--ctor.<param>=actor:<actor_id>
```

Example:

```bash
skyra planner --ctor.user=actor:user.mike --ctor.goal=actor:goal.fitness --query="best next action"
```

## Composition Contract

Nested call:

```bash
skyra <root_skill>.<nested_skill> [args...]
```

Skill-as-input:

```bash
--skill.<param>=skill:<skill_id>
```

Provisioning declaration in root contract:

```text
provisioned_skills: [<nested_skill_1>, <nested_skill_2>, ...]
```

Validity requirements:

1. Root skill is defined and committed.
2. Nested skill is defined and committed.
3. Both skills are provisioned in Redis.
4. Nested skill is declared in root `provisioned_skills[]`.
5. Boundary policy allows root -> nested call.
6. Call stays inside originating user intent scope.

If any check fails, gateway/kernel rejects dispatch.

Detailed lineage/idempotency/depth/error rules:
`docs/arch/v1/skill/skill-composition-protocol.md`.

## Registration and Routing Contract

Execution gate:

1. A command must resolve to a skill contract in Redis to execute.
2. If skill is not defined/provisioned in Redis, kernel rejects execution.

Routing:

1. Routing is dynamic at runtime.
2. Kernel routes by skill contract + live shard capability/policy state.
3. Shards are defined by the kernel; unrecognized shard identities are invalid.

This is static contract + dynamic dispatch.

## Secondary Primitives (Runtime Interfaces)

- `user_event` — user input primitive interface.
  Spec: `docs/arch/v1/native-protocol/primitives/user_event.md`
- `reply` — user-visible output primitive.
  Spec: `docs/arch/v1/native-protocol/primitives/reply.md`
- `remember` — bounded memory-write primitive interface.
  Spec: `docs/arch/v1/native-protocol/primitives/remember.md`
- `retrieve` — bounded context retrieval primitive interface.
  Spec: `docs/arch/v1/native-protocol/primitives/retrieve/retrieve.md`
- `decide` — bounded decision primitive interface.
  Spec: `docs/arch/v1/native-protocol/primitives/decide.md`
- `recurse` — bounded child-call primitive interface.
  Spec: `docs/arch/v1/native-protocol/primitives/recurse.md`

## Context Shrink and Deposit

Skills are context reducers.

Rules:

1. A skill executes with local working context in its call frame.
2. On return, local working context is shed.
3. Child returns only bounded deposit (`context_delta_refs[]` + bounded metadata).
4. Parent continuation and final output are composed from bounded deposits.
5. Overflow follows explicit failure policy; no silent carry-through.

## Decision Continuity and Job Closure

The protocol assumes model memory is unreliable; decisions are enforceable state.

Decision record:

```text
decision_record {
  decision_id
  job_id
  call_id
  must_call[]
  must_not_call[]
  close_allowed
  reason
  created_at
}
```

Rules:

1. Material decisions append a `decision_record`.
2. Downstream calls are checked against active `must_not_call[]`.
3. Job closure is explicit:

```bash
skyra job.close --job=<job_id> --decision=<decision_id> --status=done|abandoned|superseded
```

4. Close is rejected when unresolved `must_call[]` remain unless explicitly waived by policy/user approval.

## Versioned Contracts

Protocol and schema versions are part of execution contract.

Rules:

1. Existing version behavior is immutable.
2. New behavior is introduced by new version.
3. Kernel dispatches by versioned handlers and rejects unsupported versions.

## Notes

- `skill:<skill_id>` is the canonical skill reference format.
- `actor:<actor_id>` is the canonical constructor reference format.
- "method" is informal terminology; it is not a protocol primitive.
- Command syntax stays native CLI; shards do not emit JSON envelopes.
