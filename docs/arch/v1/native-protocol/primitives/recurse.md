# Recurse Primitive (v1)

Recurse is a primitive.
It opens a bounded child call and guarantees bounded return into the original parent workflow.

## Purpose

Enable reusable nested workflows while preventing unbounded context growth.

## Command Shape

Root form:

```bash
octos recurse [args...]
```

Nested form:

```bash
octos <root_skill>.recurse [args...]
```

## Required Fields

- `to`: nested target in `<root_skill>.<nested_skill>` form
- `input_refs[]`: refs passed into child call (`node:`, `edge:`, `artifact:`, `skill:`)
- `depth_left`: recursion depth budget remaining
- `context_units_left`: context-unit budget remaining
- `return_max_units`: max context units child may return to parent
- `return_max_bytes`: max bytes child may return to parent
- `on_overflow`: `collapse | summarize_once_then_collapse`
- `idempotency_key`

## Optional Fields (Recommended)

- `branch_units_used`
- `branch_bytes_used`
- `max_no_progress`
- `trace_id`
- `reason_required`: `true | false`

## Base Case

Stop recursion when any condition is true:

1. `depth_left == 0`
2. `context_units_left == 0`
3. frontier or candidate set is empty
4. `no_progress_count >= max_no_progress`

## Out-Call Stack

Runtime keeps an explicit out-call stack per job.

`call_frame`:

```
call_frame {
  frame_id
  job_id
  call_id
  parent_call_id
  root_skill
  active_skill
  input_refs[]
  return_max_units
  return_max_bytes
  depth_left_at_entry
  context_units_left_at_entry
  status            // open | returned | failed | collapsed
  error_code
  started_at
  ended_at
}
```

Rules:

1. Push frame on recurse call.
2. Pop frame on `returned | failed | collapsed`.
3. Parent workflow resumes only after frame resolves.

## Return Contract

Child returns:

- `context_delta_refs[]` (context deposit to parent)
- `returned_units`
- `returned_bytes`
- `stop_reason`
- `decision_refs[]`

Parent validates:

1. `returned_units <= return_max_units`
2. `returned_bytes <= return_max_bytes`

If valid, parent merges delta and continues.

Unwind shrink rule:

1. Child local working context is dropped at frame return.
2. Parent receives only the bounded context deposit (`context_delta_refs[]`) plus bounded metadata.
3. Final outputs are built from merged deposits, not raw child-frame context.

## Overflow Collapse

If return contract is violated:

1. mark child frame `collapsed`
2. append failure record
3. emit `ERR_CONTEXT_OVERFLOW`
4. follow `on_overflow`:
   - `collapse`: fail branch immediately
   - `summarize_once_then_collapse`: one summarize attempt, then collapse on failure or overflow

Never silently truncate child output.

## Validation Rules

1. Reject if `depth_left < 0`.
2. Reject if `context_units_left < 0`.
3. Reject if `return_max_units <= 0`.
4. Reject if `return_max_bytes <= 0`.
5. Reject if target nested skill is not provisioned by parent root skill.
6. Reject return payloads that include full child-frame context instead of bounded deposit refs.

## Output Contract

Return:

- `frame_id`
- `status`
- `returned_units`
- `returned_bytes`
- `context_delta_refs[]`
- `error_code`
- `stop_reason`

## Notes

- Recurse bounds context by contract, not model memory.
- Parent continuity is deterministic: valid bounded result or explicit failure.
