# Skyra Agents and Services Catalog (v1)

> **DEPRECATED** — The agent model has been replaced. Agents are no longer a first-class concept. New primitives: **Skill** (learned class), **Job** (skill instance), **Task** (execution unit), **Memory** (provisioned namespace), **Entity** (named thing inside memory). The kernel owns the heap and router. Redis is the trust boundary and skill registry. See `docs/arch/v1/kernel.md` for the canonical architecture. This document is preserved for historical reference only.

This document defines what counts as a `shard`, an `agent`, and a `service` in Skyra and lists the components used in v1.

## 1. Definitions

- `Service`: a long-running backend component with a clear API/contract.
- `Shard`: a lightweight daemon that Skyra deploys onto a device. It boots, fingerprints the device's hardware and software environment, registers its capabilities with the control plane, and listens for commands.
- `Agent`: a scoped domain of the user's life (e.g. home, work, health, music, servers). Agents have their own memory, tools, boundaries, and state. They are the conceptual and intelligence layer. Two types exist: system agents (predefined at init) and domain agents (user-created). See `docs/arch/v1/agents/README.md`.
- `Domain Expert`: a control-plane service module used during task formation. It is not a remote shard.

Practical rule:

- If it makes routing/task decisions in Brain Shard runtime, treat it as a service module.
- If it is a daemon deployed onto a device that executes commands from the control plane, treat it as a shard.
- If it is a scoped domain of the user's life with its own state and tools, it is an agent.

## 2. Core Control-Plane Services (Brain Shard)

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

### 2.3 Domain Agent Self-Selection

There is no central classifier. Domain agents are the doorkeepers of their own domains.

The front face transformer labels the incoming turn (in-domain or "other") using the context blob — which contains all registered agents with their relevance scores. Relevant domain agents receive the full context blob and self-select. No service chooses an agent on their behalf.

- Code (front face transformer / labeling):
  - `skyra/internal/controlplane/classifier.go` _(needs update to reflect self-selection model)_
  - `skyra/internal/taskformation/routing.go` _(needs update)_

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

Note: previously called "Project Service" and before that "Memory Service". The Agent Service supersedes those definitions with a more clearly scoped responsibility model. See `skyra/internal/agent/README.md`.

- Role: single owner of all agent state. Manages agent registry, object store commits, rollback, and audit trail. Tools live in the object store — no separate tool registry.
- Owns:
  - Agent Registry (SQLite) — fast index of all agents with status and last_active_at
  - Object Store Interface — commits, HEAD, state.json, rollback
  - Tools — stored under `tools/` in the object store, versioned via commits, discovered by the LLM walking the filesystem
- Code:
  - `skyra/internal/agent/`
  - `skyra/internal/memory/objectstore/fs/store.go`
  - `skyra/internal/memory/objectstore/s3/store.go`
  - `skyra/internal/memory/commit/`
- Design:
  - `skyra/internal/agent/README.md`

### 2.7 Tooling Services

- Role: allowlisted tool execution, shard capability registry, runtime capability resolution, and auditing.
- Owns:
  - Shard Capability Registry (global runtime index) — shard-reported capabilities, health/load, location, trust metadata
  - Capability Resolver — binds a selected domain tool to a concrete shard capability at execution time
- Code:
  - `skyra/internal/tools/registry.go`
  - `skyra/internal/tools/allowlist.go`
  - `skyra/internal/tools/exec/shell.go`
  - `skyra/internal/tools/audit/log.go`

### 2.8 Capability Binding Model (v1)

Tools live in the object store. The LLM discovers them by walking the filesystem. When the LLM selects a tool with `required_capabilities[]`, the orchestrator resolves that against the Shard Capability Registry to find the right execution target.

Two layers — one for tool definitions, one for execution targets:

- Object Store Tools (Agent Service): tool definitions the LLM reads and selects during execution. Navigable directly via filesystem.
- Shard Capability Registry (Tooling Services): concrete executable capabilities exposed by shards at runtime.

Binding flow:

1. Domain Expert selects a domain tool during planning.
2. Orchestrator requests capability resolution for that tool's `required_capabilities[]`.
3. Capability Resolver selects a shard using policy constraints (boundary/location/trust) plus runtime state (health/load/latency).
4. Tooling service dispatches execution and records the resolved binding for audit/replay.

## 3. Shards

A Shard is a lightweight daemon deployed by Skyra onto a device. Every device in the Skyra network runs a Shard. The Shard boots, fingerprints its hardware and software environment, registers its capabilities with the control plane, and listens for commands.

The control plane treats all Shards uniformly. A Shard is identified by its capability profile — not by its hardware type.

### 3.1 Voice Shard (Voice Capability Profile)

The Voice Shard runs with a voice-specific capability profile:

- `microphone`: always-on audio capture
- `wake_word`: openWakeWord or Porcupine
- `vad`: voice activity detection
- `stt`: speech-to-text (Whisper small/base)
- `tts`: text-to-speech (Piper or Coqui)
- `front_door_model`: fast local LLM (Llama 3.2 3B Instruct Q4_K_M)
- `outbox`: durable event outbox with retry sender

The Voice Shard is not architecturally special. It is a Shard that happens to have voice capabilities. The control plane treats it as a Shard whose capability profile includes audio input and output.

- Code:
  - `skyra/services/listener/app/main.py`
  - `skyra/services/listener/README.md`

Boundary:

- Voice Shard is non-authoritative for semantic responses.
- Voice Shard emits ACK/progress and renders backend-authored `UPDATE|PLAN_PROGRESS|CLARIFY|PLAN_APPROVAL_REQUIRED|FINAL|ERROR`.

### 3.2 Machine Shards (Laptop/Desktop/Server)

General-purpose Shards deployed on workstations and servers. Capability profile includes command execution.

- Role: execute allowlisted commands sent by control plane.
- Code:
  - `skyra/cmd/skyra-agent/main.go`
  - `skyra/internal/executor/client.go`
  - `skyra/internal/executor/executor.go`
  - `skyra/internal/executor/security.go`

### 3.3 TV Shard (Voice Input via Remote)

Smart TVs with voice remotes register as ingress shards. The remote delivers audio to the TV OS — whether via Bluetooth, RF, or a proprietary protocol is irrelevant to Skyra. The TV OS is the recording device.

Capability profile:

- `voice_input`: TV OS receives voice from the remote and exposes it via OS voice API (Android TV voice API, Tizen SDK, webOS SDK). The Skyra TV shard hooks this API before the native assistant processes the audio.
- `speaker`: TV audio output (TTS playback)
- `display`: TV screen output (optional visual rendering)

How it works:

- User speaks into the TV remote.
- Remote transmits audio to the TV via its native protocol (Bluetooth, RF, infrared — opaque to Skyra).
- TV OS raises a voice input event.
- Skyra TV shard intercepts the event via the OS API hook — before it reaches the native assistant.
- Shard packages the audio as a `voice_event` and sends it to the brain.
- The remote is invisible to Skyra. The TV is the ingress shard.

The TV shard registers `voice_input` in its capability profile. The Shard Capability Registry sees it as an ingress-capable shard, same as the Raspberry Pi Voice Shard. Spatial awareness applies: the TV shard's network tag scopes it to the room it's in.

### 3.4 Shard-to-Shard Delegation (v1 direction)

Shard-to-shard calls are allowed, but only through control-plane mediation:

- A shard requests delegation for a capability it cannot satisfy locally.
- Control plane validates policy and emits a short-lived delegation token.
- Target shard executes and returns results through control-plane-tracked channels.
- Every delegation edge is audited (`source_shard`, `target_shard`, `capability_id`, `job_id`).

No unmanaged peer mesh is allowed in v1.

## 4. Single vs Multiple Vector DBs

v1 recommendation:

- one shared vector infrastructure
- strict metadata scoping per query (`agent_id`, `domain_id`, `memory_type`, `time_window`)
- optional domain namespaces only when scale/noise requires it

## 5. Ownership Summary

- Voice Shard: capture + transport + render (voice capability profile)
- TV Shard: intercept remote voice via OS API hook, package as voice_event, render TTS/display output
- Machine Shards: execute commands only
- Control-plane services: label turns + domain agent self-selection + form tasks + orchestrate + commit agent state
- Agent Service: own agent registry, object store, domain tools (filesystem files in object store), commit history
- Tooling Services: own shard capability registry, capability resolution, and execution dispatch/audit

If you ask, "is Domain Expert a shard?":

- No. Domain Expert is a control-plane service module, not a device daemon.
- Domain Expert details live in `docs/arch/v1/domain-expert/README.md`.
