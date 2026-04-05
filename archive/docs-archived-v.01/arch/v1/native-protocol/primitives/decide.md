# Decide Primitive (v1)

Decide is a primitive.
It performs bounded recursive decision-making over candidate outputs/context.

## Purpose

Given candidate references, choose one action:

- `keep`
- `drop`
- `split`
- `delegate`
- `close`

`extract` is not a primitive. It is derived behavior from `retrieve + decide`.

For failure-handling inside a job loop, the same decision surface may also be
used to choose one of:

- `retry_same`
- `retry_revised`
- `escalate`
- `close`

This distinguishes mechanical retry from cognitive retry. Mechanical retry
belongs to kernel execution policy. These outcomes are for reasoning over a
failed attempt and choosing what to do next.

## Command Shape

Root form:

```bash
skyra decide [args...]
```

Nested form:

```bash
skyra <root_skill>.decide [args...]
```

## Required Fields

- `candidate_refs[]`: list of candidate refs (`actor:<id>`, `edge:<id>`, `artifact:<id>`, or output refs)
- `decision_budget_left`: remaining decision budget
- `max_depth`: recursion depth cap
- `close_threshold`: minimum confidence required to close

## Optional Fields (Recommended)

- `policy`: `strict | balanced | aggressive`
- `allow_partial`: `true | false`
- `max_no_progress`: max consecutive non-progress recursion steps
- `reason_required`: `true | false`
- `layer`: `committed | full`

## Base Case

Stop when any condition is true:

1. `decision_budget_left == 0`
2. `candidate_refs[]` is empty
3. `depth >= max_depth`
4. `action == close` and `confidence >= close_threshold`
5. `no_progress_count >= max_no_progress`

## Recursive Step

1. Evaluate candidate set.
2. Choose action (`keep|drop|split|delegate|close|retry_same|retry_revised|escalate`).
3. Append a `decision_record`.
4. If action is `split`, `delegate`, `retry_same`, or `retry_revised`, recurse on child candidate sets with decremented `decision_budget_left`.

## Output Contract

Return:

- `decision_id`
- `action`
- `selected_refs[]`
- `discarded_refs[]`
- `confidence`
- `reason`
- `decision_budget_used`
- `decision_budget_left`
- `stop_reason`: `budget_exhausted | empty_candidates | depth_limit | close_threshold | no_progress`

## Validation Rules

1. Reject if `decision_budget_left <= 0` at start.
2. Reject if `max_depth < 0`.
3. Reject if `action=close` and unresolved required obligations exist.
4. Reject if `layer=committed` and candidate source is uncommitted.

## Notes

- Decision continuity is append-only via `decision_record`.
- Use this primitive to prevent reliance on model memory in deep recursive workflows.
