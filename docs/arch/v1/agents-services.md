# Skyra Agents and Services Catalog (v1)

This document defines what counts as an `agent` vs a `service` in Skyra and lists the components used in v1.

## 1. Definitions

- `Service`: a long-running backend component with a clear API/contract.
- `Agent`: a constrained executor that performs actions from control-plane commands.
- `Domain Expert`: a control-plane service module used during task formation. It is not a remote machine agent.

Practical rule:

- If it makes routing/task decisions in Mac runtime, treat it as a service module.
- If it executes allowlisted commands on a target machine, treat it as an agent.

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

- Role: chooses domain/project candidate set and confidence.
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

### 2.6 Memory Services

- Role: authoritative state/object commits + semantic retrieval store.
- Code:
  - `skyra/internal/memory/store.go`
  - `skyra/internal/memory/objectstore/fs/store.go`
  - `skyra/internal/memory/project/registry.go`
  - `skyra/internal/memory/vectorstore/chroma/store.go`
  - `skyra/internal/memory/vectorstore/qdrant/store.go`

### 2.7 Tooling Services

- Role: allowlisted tool execution and auditing.
- Code:
  - `skyra/internal/tools/registry.go`
  - `skyra/internal/tools/allowlist.go`
  - `skyra/internal/tools/exec/shell.go`
  - `skyra/internal/tools/audit/log.go`

## 3. Edge Service (Raspberry Pi)

### 3.1 Listener Service

- Role: wake word, VAD, STT/TTS, triage, outbox transport.
- Code:
  - `skyra/services/listener/app/main.py`
  - `skyra/services/listener/README.md`

Boundary:

- Pi listener is non-authoritative for semantic responses.
- Pi emits ACK/progress and renders backend-authored `UPDATE|PLAN_PROGRESS|CLARIFY|PLAN_APPROVAL_REQUIRED|FINAL|ERROR`.

## 4. Remote Execution Agents

### 4.1 Machine Agents (Laptop/Desktop/Server)

- Role: execute allowlisted commands sent by control plane.
- Code:
  - `skyra/cmd/skyra-agent/main.go`
  - `skyra/internal/agent/client.go`
  - `skyra/internal/agent/executor.go`
  - `skyra/internal/agent/security.go`

These are true `agents` in the architecture.

## 5. Single vs Multiple Vector DBs

v1 recommendation:

- one shared vector infrastructure
- strict metadata scoping per query (`project_id`, `domain_id`, `memory_type`, `time_window`)
- optional domain namespaces only when scale/noise requires it

## 6. Ownership Summary

- Pi listener: capture + transport + render
- Control-plane services: classify + form tasks + orchestrate + commit memory
- Remote agents: execute commands only

If you ask, "is Domain Expert an agent?":

- In v1 terminology, it is a control-plane service module, not a remote executor agent.
- Domain Expert details live in `docs/arch/v1/domain-expert/README.md`.
