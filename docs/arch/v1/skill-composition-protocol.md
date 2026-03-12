# Skill Composition Protocol (v1)

This protocol is the object-oriented representation of system execution.
The kernel is the interpreter.

This protocol defines how one skill calls another skill and how a skill reference is passed as input to another skill.

The command language stays the same:

```
octos <skill> [args...]
```

No JSON payload is added to shard-emitted commands.

---

## Invariants

1. **Intent-bound execution is mandatory.** Nested calls are bound to the originating user intent scope. If a call cannot be tied to an approved root intent scope, it is rejected. No exceptions.
2. **Registry-gated execution is mandatory.** Root and nested skills must be provisioned in Redis.
3. **Nested skill declaration is mandatory.** Nested skill must be declared in the root skill's `provisioned_skills[]`.
4. **Boundary rules are mandatory.** Root skill must be allowed to invoke the nested skill (`allowed`, `gated`, `denied`).
5. **Skill references are first-class typed inputs.** Skill-as-input values are explicit `skill_ref` values.
6. **Caller context is gateway/kernel-owned metadata.** Shards do not self-assert parent IDs or trust metadata.
7. **No task primitive exists in this protocol.** Composition and closure operate at job/call level.
8. **Reasoning scope is provisioned by the root skill.** The model may reason over only nested skills listed in root `provisioned_skills[]`.

---

## Skill Reference Type

Canonical `skill_ref` format:

```
skill:<skill_id>
```

- `skill_id` is the committed immutable skill ID.
- Name-only references are not canonical for composition because they are ambiguous across versions/models.

---

## Skill-As-Input Encoding

A skill input is passed with a reserved flag namespace:

```
--skill.<param>=skill:<skill_id>
```

Examples:

```
octos integrate --skill.source=skill:reasoning.v1 --skill.target=skill:search.v3
octos orchestrator.integrate --skill.source=skill:reasoning.v1 --skill.target=skill:search.v3
```

Rules:

1. `--skill.<param>` may be repeated.
2. Value must parse as canonical `skill_ref`.
3. Gateway resolves each referenced skill and fails fast if any are missing/untrusted.

---

## Skill-to-Skill Call Semantics

A running call invokes another skill by emitting:

```
octos <root_skill>.<nested_skill> [args...]
```

This is a nested call if emitted from an active call context. It is a root call otherwise.

The command itself stays simple; call lineage is tracked in gateway/kernel metadata.

---

## Invocation Metadata (Gateway-Assembled)

`job_envelope_v1` must include:

```
invocation {
  mode: root | nested
  intent_scope_id      // immutable scope for this call tree
  call_id              // unique call ID for dedupe/trace
  depth                // 0 for root, +1 per nested call
  caller_skill_id      // null for root
  parent_job_id        // null for root
  parent_call_id       // null for root
}
```

Rules:

1. Root call sets `mode=root`, `depth=0`, `intent_scope_id` from ingress identity.
2. Nested call inherits `intent_scope_id` exactly from parent.
3. Nested call increments `depth`.
4. Caller identifiers are copied from current execution context, not command args.

---

## Validation Order

For every command:

1. Parse command as root form `octos <skill> [args...]` or nested form `octos <root_skill>.<nested_skill> [args...]`.
2. Resolve root skill from Redis.
3. If nested form, resolve nested skill from Redis.
4. Resolve all `--skill.<param>` references.
5. Build invocation metadata (`root` or `nested`).
6. If nested form, enforce nested skill is declared in root skill `provisioned_skills[]`.
7. Enforce caller boundary rules against nested skill.
8. Enforce intent scope continuity (`nested` must match parent scope).
9. Enforce depth/cycle/idempotency guards.
10. Enqueue job for execution.

Reject immediately on first failure.

---

## Safety Guards

Required defaults:

1. `max_call_depth = 8`
2. `max_nested_calls_per_job = 128`
3. Cycle guard: reject if the same `(caller_skill_id, nested_skill_id, normalized_args_hash)` repeats in the same branch beyond retry budget.
4. Idempotency key: `(intent_scope_id, parent_call_id, root_skill_id, nested_skill_id, normalized_args_hash)`.

---

## Decision Continuity

Models can forget. The protocol therefore treats decisions as enforceable state.

Append-only decision record:

```
decision_record {
  decision_id
  job_id
  call_id
  must_call[]      // required skill calls
  must_not_call[]  // forbidden skill calls
  close_allowed
  reason
  created_at
}
```

Enforcement:

1. Before dispatching a downstream call, reject if call matches active `must_not_call[]`.
2. A job may close only through explicit close command with a decision reference:

```
octos job.close --job=<job_id> --decision=<decision_id> --status=done|abandoned|superseded
```

3. Reject close when unresolved `must_call[]` obligations remain, unless policy/user approval explicitly waives them.

---

## Reserved Flags

The following flag namespace is reserved for protocol/meta fields and may not be user/tool-authored:

```
--_*
```

If present in a command payload, reject as protocol violation.

---

## Errors

Standard rejection reasons:

- `ERR_SKILL_NOT_FOUND`
- `ERR_SKILL_REF_INVALID`
- `ERR_SKILL_REF_NOT_TRUSTED`
- `ERR_BOUNDARY_DENIED`
- `ERR_INTENT_SCOPE_MISMATCH`
- `ERR_CALL_DEPTH_EXCEEDED`
- `ERR_CALL_CYCLE_DETECTED`
- `ERR_IDEMPOTENCY_DUPLICATE`
- `ERR_RESERVED_FLAG`
- `ERR_OBLIGATION_UNRESOLVED`
- `ERR_JOB_CLOSE_NOT_ALLOWED`
- `ERR_NESTED_SKILL_NOT_PROVISIONED`

---

## Minimal Examples

Root call:

```
octos orchestrator --targets=home,gym --intent="turn off lights and cancel gym"
```

Nested call:

```
octos orchestrator.search --query="today's gym booking"
```

Skill-as-input:

```
octos integrate --skill.source=skill:reasoning.v1 --skill.target=skill:search.v3
```
