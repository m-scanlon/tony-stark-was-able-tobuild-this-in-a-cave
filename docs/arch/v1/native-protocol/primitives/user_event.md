# User Event Primitive (v1)

User event is a primitive.
It normalizes inbound user signal into typed input for skill execution.

## Purpose

Provide one input primitive interface for user-originated events, with multiple implementations behind the same contract.

## Command Shape

Root form:

```bash
octos user_event [args...]
```

Nested form:

```bash
octos <root_skill>.user_event [args...]
```

## Required Fields

- `impl`: `voice | text`
- `payload_ref`: input payload reference
- `session_ref`: session reference

## Optional Fields (Recommended)

- `locale`
- `channel`: `voice | chat | auto`
- `confidence`
- `source_shard`

## Validation Rules

1. Reject if `impl` is not in kernel allowlist for user_event.
2. Reject if required refs are missing.
3. Reject if payload exceeds configured ingress bounds.

## Output Contract

Return:

- `event_id`
- `impl`
- `normalized_refs[]`
- `status`: `accepted | rejected`
- `error_code`

## Notes

- `user_event` is an interface primitive; implementations are extensible.
- `voice` and `text` are v1 implementations.
- Additional implementations must preserve this output contract.
