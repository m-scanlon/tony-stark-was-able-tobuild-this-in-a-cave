# Skyra Listener Service

Always-on listener service for wake/audio events.

Scope for this service:
- Always-on listener-facing API surface
- Deterministic intent gate
- Event-driven handoff to front-door LLM (not always running inference)

## Recommended Front-Door Model (Small Device)

- `Llama-3.2-3B-Instruct` (GGUF, `Q4_K_M` preferred)
- Context target on Raspberry Pi: `4096` to `8192`
- Keep long memory retrieval and heavy context assembly in the control plane

## Context Compression Integration

Use control-plane compression before passing retrieved memory back to listener/front-door prompts:

- Package: `skyra/internal/context/compress`
- Purpose: rank + trim + budget context for low-latency prompts

## Run locally

```bash
cd skyra/services/listener
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8090
```

## Run with Docker

```bash
cd skyra/services/listener
docker build -t skyra-listener:dev .
docker run --rm -p 8090:8090 skyra-listener:dev
```

## Endpoints

- `GET /health`
- `POST /listener/event`
