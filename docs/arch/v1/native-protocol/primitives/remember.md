# Remember Primitive (v1)

Remember is a primitive.
It writes bounded memory artifacts from active workflow state.

## Purpose

Provide one memory-write primitive interface with multiple implementations while preserving trust, bounds, and intent scope.

## Command Shape

Root form:

```bash
octos remember [args...]
```

Nested form:

```bash
octos <root_skill>.remember [args...]
```

## Required Fields

- `memory_refs[]`: refs to store (`node:`, `edge:`, `artifact:`)
- `target_namespace`: destination memory namespace
- `layer`: `working | committed`
- `idempotency_key`

## Optional Fields (Recommended)

- `impl`: `auto | episodic | semantic | procedural` (default `auto`)
- `reason`
- `ttl`
- `dedupe_mode`: `hash | semantic`

## Validation Rules

1. Reject if `memory_refs[]` is empty.
2. Reject if `target_namespace` is missing.
3. Reject if `impl` is not in kernel allowlist for remember.
4. Reject committed writes without required approval path/policy.
5. Reject writes outside active intent scope.

## Output Contract

Return:

- `remember_id`
- `stored_refs[]`
- `rejected_refs[]`
- `layer`
- `status`: `stored | partial | rejected`
- `error_code`

## Notes

- `remember` is an interface primitive; implementations are extensible.
- All remember implementations must satisfy this output contract.
- Use `layer=working` for scratch memory and `layer=committed` for user-gated canonical memory.
