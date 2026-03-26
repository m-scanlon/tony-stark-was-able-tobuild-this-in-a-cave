# Skyra High-Level Architecture Sheet (v1)

## 1. What v1 Is

Skyra v1 is a cognitive runtime that continuously receives stimuli, keeps an organized chain of thought, and separates cognition from outbound interaction.

The kernel is a standalone service.

## 2. Architecture In One View

```text
Stimulus Source
  -> Ingress
  -> Stimulus Stream
  -> Nexus
  -> Kernel Max Heap
  -> Chain of Thought
  -> understand | interact
  -> Human-to-Machine Interaction
  -> Resulting Stimulus
  -> Stimulus Stream
```

## 3. Core Runtime Rules

- The stimulus stream is unbounded.
- Attention is bounded.
- The kernel queue is a max heap.
- User messages have priority `100`.
- Internal chain-of-thought work has priority `50`.
- The top-level primitives are `understand` and `interact`.
- `understand` enters `reference`, loops with `infer`, and stops at `resolve`.
- `resolve` is the only step that creates `understanding`.
- `perception` always contains `history` and `stimulus`.
- `understanding` appears in `perception` only after `resolve`.

## 4. Runtime State

v1 keeps runtime state in memory only.

The model should start with fresh memory on each step or request. The only intentional carry-over is the active frame and its `perception`.

Restart persistence is out of scope for v1.

## 5. Frontend Split

v1 has two separate frontend surfaces over the same backend:

- a human interaction surface
- a read-only chain-of-thought surface

The chain-of-thought surface is chronological and shows raw internal output together with the active frame, primitive choice, and current perception.

## 6. Not Canonical For v1

The later skill system and registry-heavy orchestration model are future-facing and are not important for the current v1 runtime.
