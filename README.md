# Skyra - Personal AI Assistant

A distributed personal AI system designed for always-on voice interaction with project-centric memory and modular hardware scaling.

## Project Vision

Skyra aims to be a private, local-first AI assistant that understands your personal context across different projects and domains. The end goal is a three-machine architecture: Raspberry Pi for voice, Mac mini as control plane, and GPU machine for heavy reasoning.

## Current State

This repository contains the architectural planning and initial project structure. The system is designed to use:

- **OpenClaw** agent runtime for orchestration
- Multiple specialized models (Llama 3.1, Qwen2.5-Coder, DeepSeek)
- Project-centric memory with semantic retrieval
- Voice interface with wake word detection

## Next Steps

1. Set up OpenClaw agent runtime
2. Implement core API endpoints
3. Configure local models (Llama 3.1 8B, Qwen2.5-Coder 7B)
4. Build memory service with vector database
5. Develop voice interface for Raspberry Pi

## Key Technologies Planned

- **Agent Runtime**: OpenClaw
- **Models**: Llama 3.1, Qwen2.5-Coder, DeepSeek-Coder
- **Databases**: PostgreSQL, Vector DB (Qdrant/Chroma)
- **Voice**: openWakeWord, Whisper, Piper/Coqui
- **API**: FastAPI/Node.js
- **Infrastructure**: Docker, private LAN