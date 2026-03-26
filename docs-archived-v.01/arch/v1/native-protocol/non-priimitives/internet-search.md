# Internet Search (Non-Primitive, v1)

`internet_search` is a system capability, not a native protocol primitive.

## Intent

Query external web sources and return bounded references/results to the active workflow.

## Why It Is Not a Primitive

- It composes existing primitives plus external adapter/runtime integration.
- It is policy-gated and environment-dependent (network, provider, trust), unlike core primitive contracts.

## Naming Note

Use `internet_search` as canonical naming in protocol docs.
