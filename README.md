# Skyra - Personal Runtime

Skyra v1 is a local-first cognitive runtime built around a standalone kernel service, an unbounded stimulus stream, and a bounded chain-of-thought attention loop.

## Current v1 Model

- The kernel is the canonical runtime boundary.
- Incoming stimuli append to the stimulus stream first.
- Nexus decides what becomes active kernel work.
- The kernel queue is a max heap.
- User messages run at priority `100`.
- Internal chain-of-thought work runs at priority `50`.
- The only top-level primitives are `understand` and `interact`.
- Native primitives are the v1 protocol language. The later skill system is not important for v1.
- Chain of Thought and Human-to-Machine Interaction are separate runtime boundaries.
- `perception` always contains `history` and `stimulus`.
- `understanding` is absent until `resolve` finishes an interpret cycle.
- v1 runtime state is in-memory only. It does not need to survive restarts.
- v1 has two separate frontend surfaces:
  - a human interaction surface
  - a read-only internal chain-of-thought surface

## Not Canonical For v1

Older docs in this repo still describe a future-facing architecture built around later-phase registry and orchestration abstractions.

That material belongs to later phases and should not be treated as important for the current v1 runtime.

## Documentation Map

- Runtime PRD: `skyra/chain-of-thought/experience.md`
- Kernel contract: `docs/arch/v1/kernel/kernel.md`
- Gateway/transport contract: `docs/arch/v1/api-gateway/api-gateway.md`
- Architecture sheet: `docs/arch/v1/high-level-architecture-sheet.md`
- Active gaps: `docs/arch/v1/gaps.md`

## Run v1

1. Copy `.env.example` to `.env` and set `OLLAMA_GATEWAY_WS_URL` and `OLLAMA_GATEWAY_TOKEN` for your local gateway.
2. Start the stack with `docker compose up --build`.
3. Open `http://127.0.0.1:9090/interact` for the human surface.
4. Open `http://127.0.0.1:9090/thoughts` for the read-only chain-of-thought surface.
