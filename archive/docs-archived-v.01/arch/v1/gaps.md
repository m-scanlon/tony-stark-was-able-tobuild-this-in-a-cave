# Skyra Architecture Gaps Register

This document tracks the current gaps for the canonical v1 runtime.

Older gaps centered on later-phase registry and orchestration abstractions are not important for the current v1 boundary.

## Priority Scale

- `P0`: blocks the current v1 runtime
- `P1`: important for the next implementation pass
- `P2`: useful, but can follow the core build

## Active Gaps

| ID | Gap | Priority | Why it matters |
| --- | --- | --- | --- |
| G1 | Kernel heap implementation is still FIFO in code | P0 | The v1 contract is a max heap with user messages above internal thought work. |
| G2 | Kernel execution loop still uses a stub executor | P0 | The runtime still needs a real step engine for `understand`, `reference`, `infer`, `resolve`, and `interact`. |
| G3 | Runtime frame schema is documented but not implemented | P0 | `perception`, `history`, `stimulus`, and post-`resolve` `understanding` need a concrete runtime type. |
| G4 | Model adapter to the Ollama WebSocket gateway is not implemented | P0 | The kernel needs a live adapter for primitive execution against the local model gateway. |
| G5 | Two frontend surfaces are specified but not implemented | P1 | v1 requires a separate human interaction surface and a separate chronological chain-of-thought surface. |
| G6 | Thought log event schema is not locked | P1 | The read-only thought surface needs a consistent per-step record for raw output, frame/template, primitive choice, and perception. |
| G7 | Observability output is not wired end-to-end | P1 | The runtime needs traceability from stimulus ingress through primitive execution and outbound interaction. |
| G8 | Backpressure and starvation policy is not yet specified beyond static priorities | P2 | A max heap alone is not enough once multiple internal events accumulate for long periods. |

## Future-Facing, Not Important For v1

The following areas belong to later phases and should not drive the current implementation:

- skill system design
- registry-gated execution design
- later orchestration instances
- transport envelopes
- restart persistence of runtime state
