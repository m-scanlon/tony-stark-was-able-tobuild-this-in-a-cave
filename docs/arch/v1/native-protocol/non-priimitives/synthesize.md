# Synthesize (Non-Primitive, v1)

`synthesize` is a system capability, not a native protocol primitive.

## Intent

Convert session signal into structured graph deltas.

Typical flow:

1. Read session history + auxiliary signal (for example VAD/time markers)
2. Produce observational nodes
3. Produce observational edges
4. Emit outputs for merge/integration workflows

## Why It Is Not a Primitive

- It is composed from primitives (`retrieve`, `decide`, `recurse`) and system orchestration.
- It may run as an OS-provided capability even when not directly exposed as a primitive command.

## Naming Note

Use `synthesize` as canonical naming for this capability in protocol docs.
`compaction` is not a capability name in this system.
