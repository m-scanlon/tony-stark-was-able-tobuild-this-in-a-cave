# Retrieve Primitive (v1)

Retrieve is a primitive.
It resolves context before skill execution.

## Purpose

Given one or more seed references, recursively collect bounded context units for a skill constructor.

This primitive is a first-class interface with multiple implementations.

## Command Shape

Root form:

```bash
octos retrieve [args...]
```

Nested form:

```bash
octos <root_skill>.retrieve [args...]
```

## Core Model

A context unit is one of:

- `node`
- `node_edge_pair`

Recursive expansion continues until a base case is reached.

## Base Case

Stop retrieval when any condition is true:

1. `context_units_left == 0`
2. `frontier_size == 0`
3. `depth >= max_depth`

This is the bounded termination rule.

## Required Fields

- `seed_refs[]`: list of `node:<id>` or `edge:<id>`
- `context_units_left`: max number of context units remaining in budget
- `unit_mode`: `node | node_edge_pair`
- `max_depth`: recursive depth limit
- `layer`: `committed | full`

## Optional Fields (Recommended)

- `impl`: `auto | ctor | frame | graph` (default `auto`)
- `direction`: `out | in | both` (default `both`)
- `edge_types_allow[]`: allowed edge types
- `node_types_allow[]`: allowed node types
- `top_k_per_hop`: max expansions per recursion step
- `min_score`: score floor for candidate acceptance
- `score_weights`: weights for `similarity`, `recency`, `edge_weight`, `confidence`
- `time_window`: optional time bound for candidate edges/nodes
- `dedupe_mode`: `node_id | path`

## Validation Rules

1. Reject if `seed_refs[]` is empty.
2. Reject if `context_units_left <= 0`.
3. Reject if `max_depth < 0`.
4. Reject if `layer=committed` and candidate source is uncommitted.
5. Reject if `impl` is not in kernel allowlist for retrieve.

## Output Contract

Return:

- `selected_units[]`
- `units_used`
- `units_left`
- `stop_reason`: `budget_exhausted | frontier_exhausted | depth_limit`
- `trace_id`

## Notes

- Retrieval must run before constructor hydration.
- Retrieval may also run during a skill workflow when policy allows.
- Retrieval enforces trust policy (`committed` vs `full`) at selection time, not post-filter time.
- All retrieve implementations (`ctor|frame|graph`) must satisfy this same output contract.
- Kernel records selected retrieve implementation in invocation metadata for replay/debug.
