# Skyra Agents and Services Catalog (v1)

This document defines what counts as a `shard`, an `agent`, and a `service` in Skyra and lists the components used in v1.

## 1. Definitions

- `Service`: a long-running backend component with a clear API/contract.
- `Shard`: a lightweight daemon that Skyra deploys onto a device. It boots, fingerprints the device's hardware and software environment, registers its capabilities with the control plane, and listens for commands.
- `Agent`: a scoped domain of the user's life (e.g. home, work, health, music, servers). Agents have their own memory, tools, boundaries, and state. They are the conceptual and intelligence layer.
- `Domain Expert`: a control-plane service module used during task formation. It is not a remote shard.

Practical rule:

- If it makes routing/task decisions in Mac runtime, treat it as a service module.
- If it is a daemon deployed onto a device that executes commands from the control plane, treat it as a shard.
- If it is a scoped domain of the user's life with its own state and tools, it is an agent.

## 2. Core Control-Plane Services (Mac mini)

### 2.1 API Gateway + Ingress

- Role: receives `/voice`, `/chat`, `/tools`, `/memory` requests and persists events.
- Code:
  - `skyra/internal/api/v1/voice.go`
  - `skyra/internal/api/v1/chat.go`
  - `skyra/internal/controlplane/server.go`

### 2.2 Orchestrator Service

- Role: authoritative runtime for routing, execution, and response emission.
- Code:
  - `skyra/internal/controlplane/orchestrator.go`
  - `skyra/internal/controlplane/router.go`

### 2.3 Classifier + Domain Routing Service

- Role: chooses domain/agent candidate set and confidence.
- Code:
  - `skyra/internal/controlplane/classifier.go`
  - `skyra/internal/taskformation/routing.go`

### 2.4 Task Formation Service

- Role: decides `no task | WorkPlan | TaskSheet+Patch`, produces canonical task object.
- Code:
  - `skyra/internal/taskformation/domain_expert.go`
  - `skyra/internal/taskformation/system_expert.go`
  - `skyra/internal/taskformation/pipeline.go`
  - `skyra/internal/taskformation/factory.go`
- Design:
  - `docs/arch/v1/task-formation.md`
  - `docs/arch/v1/domain-expert/README.md`

### 2.5 Context Engine Services

- Role: retrieve, rank, compress, and inject relevant context.
- Code:
  - `skyra/internal/context/engine.go`
  - `skyra/internal/context/assemble.go`
  - `skyra/internal/context/vector/engine.go`
  - `skyra/internal/context/compress/engine.go`
  - `skyra/services/context-injector/README.md`

### 2.6 Agent Service

Note: previously called "Project Service" and before that "Memory Service". The Agent Service supersedes those definitions with a more clearly scoped responsibility model. See `skyra/internal/project/README.md`.

- Role: single owner of all agent state. Manages agent registry, object store commits, rollback, audit trail, and the local tool registry.
- Owns:
  - Agent Registry (SQLite) — fast index of all agents with status and last_active_at
  - Object Store Interface — commits, HEAD, state.json, rollback
  - Local Tool Registry (Vector DB) — per-agent tools with `categories[]` (operation tags for boundary enforcement) and `requires_approval` flag, retrieved via vector search
- Code:
  - `skyra/internal/project/`
  - `skyra/internal/memory/objectstore/fs/store.go`
  - `skyra/internal/memory/objectstore/s3/store.go`
  - `skyra/internal/memory/vectorstore/chroma/store.go`
  - `skyra/internal/memory/vectorstore/qdrant/store.go`
  - `skyra/internal/memory/commit/`
- Design:
  - `skyra/internal/project/README.md`

### 2.7 Tooling Services

- Role: allowlisted tool execution and auditing.
- Code:
  - `skyra/internal/tools/registry.go`
  - `skyra/internal/tools/allowlist.go`
  - `skyra/internal/tools/exec/shell.go`
  - `skyra/internal/tools/audit/log.go`

## 3. Shards

A Shard is a lightweight daemon deployed by Skyra onto a device. Every device in the Skyra network runs a Shard. The Shard boots, fingerprints its hardware and software environment, registers its capabilities with the control plane, and listens for commands.

The control plane treats all Shards uniformly. A Shard is identified by its capability profile — not by its hardware type.

### 3.1 Pi Shard (Voice Capability Profile)

The Raspberry Pi runs a Shard with a voice-specific capability profile:

- `microphone`: always-on audio capture
- `wake_word`: openWakeWord or Porcupine
- `vad`: voice activity detection
- `stt`: speech-to-text (Whisper small/base)
- `tts`: text-to-speech (Piper or Coqui)
- `front_door_model`: fast local LLM (Llama 3.2 3B Instruct Q4_K_M)
- `outbox`: durable event outbox with retry sender

The Pi Shard is not architecturally special. It is a Shard that happens to have voice capabilities. The control plane treats it as a Shard whose capability profile includes audio input and output.

- Code:
  - `skyra/services/listener/app/main.py`
  - `skyra/services/listener/README.md`

Boundary:

- Pi Shard is non-authoritative for semantic responses.
- Pi emits ACK/progress and renders backend-authored `UPDATE|PLAN_PROGRESS|CLARIFY|PLAN_APPROVAL_REQUIRED|FINAL|ERROR`.

### 3.2 Machine Shards (Laptop/Desktop/Server)

General-purpose Shards deployed on workstations and servers. Capability profile includes command execution.

- Role: execute allowlisted commands sent by control plane.
- Code:
  - `skyra/cmd/skyra-agent/main.go`
  - `skyra/internal/executor/client.go`
  - `skyra/internal/executor/executor.go`
  - `skyra/internal/executor/security.go`

## 4. Single vs Multiple Vector DBs

v1 recommendation:

- one shared vector infrastructure
- strict metadata scoping per query (`agent_id`, `domain_id`, `memory_type`, `time_window`)
- optional domain namespaces only when scale/noise requires it

## 5. Ownership Summary

- Pi Shard: capture + transport + render (voice capability profile)
- Machine Shards: execute commands only
- Control-plane services: classify + form tasks + orchestrate + commit agent state
- Agent Service: own agent registry, object store, local tool registry, commit history

If you ask, "is Domain Expert a shard?":

- No. Domain Expert is a control-plane service module, not a device daemon.
- Domain Expert details live in `docs/arch/v1/domain-expert/README.md`.
