# Reply Primitive (v1)

Reply is a primitive.
It emits user-visible output from a running job.

## Purpose

Provide a bounded, explicit response primitive for user-facing updates and final answers.

## Command Shape

Root form:

```bash
octos reply [args...]
```

Nested form:

```bash
octos <root_skill>.reply [args...]
```

## Required Fields

- `text`: user-visible response text
- `status`: `update | final | error`

## Optional Fields (Recommended)

- `channel`: `voice | chat | push | auto` (default `auto`)
- `confidence`: `0.0..1.0`
- `artifact_refs[]`: refs attached to the reply (`artifact:`, `node:`, `edge:`)

## Validation Rules

1. Reject if `text` is empty.
2. Reject if `status` is outside `update|final|error`.
3. Reject if reply is emitted outside active intent scope.

## Output Contract

Return:

- `reply_id`
- `status`
- `channel`
- `delivered`: `true | false`
- `error_code`

## Notes

- `reply` is pre-provisioned as a system primitive skill.
- Job closure remains explicit via `octos job.close ...`; `reply status=final` is not implicit close.
