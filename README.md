# Skyra - Personal AI Assistant

Distributed personal AI system for always-on voice interaction, project-centric memory, and modular hardware scaling.

## Architecture (Current Direction)

- Raspberry Pi: listener service (wake word, VAD, STT, deterministic intent gate)
- Mac mini: control plane (LangGraph orchestration, APIs, memory, routing)
- GPU machine: heavy reasoning model (DeepSeek)

The listener pipeline is always-on, while front-door LLM inference is event-driven (invoked after wake/audio capture).

## Model Strategy

- Listener/front-door on small device:
  - `Llama-3.2-3B-Instruct` (recommended) using GGUF quant such as `Q4_K_M`
  - Target context window: `4k-8k` for responsiveness
- Control plane / heavy reasoning:
  - Llama/Qwen for local control-plane tasks
  - DeepSeek on GPU for complex reasoning

## Implemented So Far

- Secured control-plane endpoints (`/v1/chat`, `/v1/voice`) with API key auth and rate limiting
- Working `/v1/chat` path to OpenAI-compatible model endpoint
- Response sanitization for reasoning artifacts (`<think>...</think>`)
- Listener service scaffold as separate service (`skyra/services/listener`)
- Context compression engine for prompt budgeting (`skyra/internal/context/compress`)
- Event ingress and ACK design with durable outbox/inbox contract (`docs/arch/v1/event-ingress-ack.md`)
- Task Formation architecture and module scaffolding (`docs/arch/v1/task-formation.md`, `skyra/internal/taskformation`)

## Context Compression

Skyra now includes a deterministic context compression engine to keep prompts small and fast.

- Package: `skyra/internal/context/compress`
- Behavior:
  - ranks chunks by score + recency
  - trims chunk length
  - enforces token budget
  - emits prompt-ready context block

Defaults:
- `MaxTokens: 700`
- `MaxChunks: 8`
- `MaxWordsPerHit: 60`

## Docs

- Main architecture: `docs/arch/v1/scyra.md`
- Model endpoint example: `docs/examples/model-endpoint-phase1.md`
- Event ingress and ACK: `docs/arch/v1/event-ingress-ack.md`
- Task Formation: `docs/arch/v1/task-formation.md`
- Listener service: `skyra/services/listener/README.md`

## Next Steps

Implement the Context Injector service (`skyra/services/context-injector`) as a first-class control-plane background service.

The Context Injector continuously watches conversation state, active tasks, and time-relevant memory signals, then builds compressed context packages and pushes them to the listener/front-door cache. This keeps front-door responses relevant and personal without blocking the always-on listener pipeline on heavy retrieval.
