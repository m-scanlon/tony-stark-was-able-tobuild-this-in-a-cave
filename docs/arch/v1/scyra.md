# Personal AI "Jarvis" AKA Skyra – End-Goal Architecture

## 1. Goals

- Always-on personal assistant with voice interaction.
- Project-centric memory (gym, work, music, servers, etc.).
- Private, local-first inference.
- Modular hardware that can scale over time.
- Fast local responses with automatic escalation to a high-reasoning GPU model.

## 2. High-Level Architecture

The system is composed of three main machines:

- **Raspberry Pi** → Voice interface (wake word, STT, TTS)
- **Mac mini** → Control plane (API, OpenClaw agent, memory, tools, fast local models)
- **GPU Machine** → Heavy reasoning model (DeepSeek LLM server)

Each machine has a clear responsibility.

## 3. Distributed Agent Architecture

### 3.1 Agent Model Overview

Each computer (laptop, desktop, server) runs a lightweight "Jarvis Agent" that extends the system's reach beyond the central control plane. The Mac mini remains the central orchestrator, sending high-level commands to agents for execution.

**Jarvis Agents are execution-only components and do not perform reasoning, memory access, or model inference.**

**Key Concepts:**

- **Control Plane**: Mac mini maintains intelligence, memory, and decision-making
- **Jarvis Agents**: Lightweight services on target machines
- **Command Distribution**: High-level intents sent to agents for execution
- **Secure Execution**: Allowlisted actions only, with authenticated connections

### 3.2 Distributed System Diagram

```mermaid
flowchart TB
  %% Voice Node
  subgraph PI[Raspberry Pi • Voice Node]
    WW[Wake Word]
    STT[Speech-to-Text]
    VCLIENT[Voice Client]
    WW --> STT --> VCLIENT
  end

  %% Control Plane
  subgraph CTRL[Mac mini • Control Plane]
    APIGW[API Gateway]
    ORCH[Orchestrator]
    CORE[OpenClaw Core]

    CHAT[Conversational Model]
    CODER[Coding Model]

    MEM[Memory Service]
    OBJ[Object Store]
    VDB[Vector DB]

    APIGW --> ORCH
    ORCH --> CORE
    CORE --> CHAT
    CORE --> CODER
    CORE --> MEM
    MEM --> OBJ
    MEM --> VDB
  end

  %% GPU Machine
  subgraph GPU[GPU Machine • Compute]
    DEEP[DeepSeek Model]
  end

  %% Agent Machines
  subgraph LAPTOP[Laptop • Agent]
    LAGENT[Jarvis Agent\nWebSocket Client]
    LEXEC[Command Executor]
    LAGENT --> LEXEC
  end

  subgraph DESKTOP[Desktop • Agent]
    DAGENT[Jarvis Agent\nWebSocket Client]
    DEXEC[Command Executor]
    DAGENT --> DEXEC
  end

  subgraph SERVER[Server • Agent]
    SAGENT[Jarvis Agent\nWebSocket Client]
    SEXEC[Command Executor]
    SAGENT --> SEXEC
  end

  %% Connections
  VCLIENT -->|voice| APIGW
  AGENT -->|complex tasks| DEEP
  ORCH -->|high-level commands| LAGENT
  ORCH -->|high-level commands| DAGENT
  ORCH -->|high-level commands| SAGENT

  %% Secure outbound connections
  LAGENT -.->|outbound WSS| CTRL
  DAGENT -.->|outbound WSS| CTRL
  SAGENT -.->|outbound WSS| CTRL
```

### 3.3 Agent Security Model

#### Authentication & Authorization

- **Token-based authentication** using mTLS or JWT tokens
- **Allowlisted commands only** - agents reject unknown actions
- **Non-root execution** - agents run as unprivileged users
- **Audit logging** - all commands logged with timestamps and results

#### Command Allowlist

```json
{
  "allowed_actions": {
    "open_vscode": {
      "cmd": "code",
      "args": ["."]
    },
    "open_browser": {
      "cmd": "google-chrome",
      "args": []
    },
    "open_terminal": {
      "cmd": "gnome-terminal",
      "args": []
    },
    "start_docker": {
      "cmd": "docker-compose",
      "args": ["up", "-d"]
    },
    "stop_docker": {
      "cmd": "docker-compose",
      "args": ["down"]
    },
    "start_minecraft_server": {
      "cmd": "systemctl",
      "args": ["start", "minecraft"]
    }
  }
}
```

#### Network Security

- **Outbound connections only** - agents initiate contact with control plane
- **WebSocket or HTTPS** for secure command channels
- **No inbound ports** - reduces attack surface on agent machines
- **Command validation** - parameters validated against schemas

### 3.4 Agent Communication Protocol

#### Command Format

```json
{
  "command_id": "cmd_12345",
  "intent": "start_development_environment",
  "action": "start_docker",
  "parameters": {
    "compose_file": "/path/to/docker-compose.yml",
    "services": ["database", "redis"]
  },
  "timeout": 30
}
```

#### Response Format

```json
{
  "command_id": "cmd_12345",
  "status": "success|error|timeout",
  "result": "Docker containers started successfully",
  "exit_code": 0,
  "timestamp": "2026-02-11T22:15:30Z"
}
```

### 3.5 Example Command Flow

**User**: "Jarvis, open VS Code on my laptop."

1. **Voice node** captures audio → sends text to control plane
2. **Orchestrator** selects target: laptop, action: open_vscode
3. **Control plane** sends command to laptop Jarvis Agent
4. **Agent** executes local command (`code .`)
5. **Agent** returns result to control plane
6. **Control plane** responds to user: "VS Code opened on your laptop"

---

## 4. Concurrency and Job Model (Planned)

The system will eventually support multiple concurrent user requests across voice, chat, and agent interactions. A formal job queue, request lifecycle management, and concurrency limits will be designed in a future iteration to ensure proper resource allocation and response ordering. This portion of the architecture is intentionally deferred for a later version to focus on core functionality first.

## 5. Model Roles

The system uses multiple specialized models instead of a single monolithic model.

### Mac mini (fast, always-on models)

**Conversational Model**

- Example: Llama 3.1 8B Instruct
- Handles:
  - General conversation
  - Intent detection
  - Clarifications
  - Task routing

**Coding / Tool Model**

- Example: Qwen2.5-Coder 7B
- Handles:
  - Script generation
  - Code editing
  - Tool execution
  - CLI-style tasks

### GPU Machine (heavy reasoning model)

**Primary Reasoning Model**

- Example: DeepSeek-Coder 33B+
- Handles:
  - Complex coding
  - Architecture decisions
  - Multi-file reasoning
  - Deep debugging
  - Long-context tasks

## 6. System Topology Diagram

```mermaid
flowchart LR
  %% ===== Nodes =====
  subgraph PI[Raspberry Pi • Voice Node]
    WW[Wake Word\n(openWakeWord/Porcupine)]
    STT[Speech-to-Text\n(Whisper small/base)]
    TTS[Text-to-Speech\n(Piper/Coqui)]
    VCLIENT[Voice Client\nHTTP/gRPC to API]
    WW --> STT --> VCLIENT
  end

  subgraph MAC[Mac mini (M4, 24GB) • Control Plane]
    APIGW[API Gateway\n(FastAPI/Node)\n/voice /chat /tools /memory]
    AGENT[OpenClaw Agent Runtime\nOrchestrator + Router]
    CLASS[Project + Intent Classifier]

    CHAT[Conversational Model\nLlama 3.1 8B]
    CODER[Coding Model\nQwen2.5-Coder 7B]

    MEMSVC[Memory Service\n(Read/Write, Summaries)]
    TOOLS[Tool/Skills Runner\n(SSH, scripts, Slack, etc.)]
    OBJ[(Object Store\n.skyra/projects\nversioned state)]
    VDB[(Vector DB\nChroma\nsemantic index)]

    APIGW --> AGENT
    AGENT --> CLASS
    AGENT --> CHAT
    AGENT --> CODER
    AGENT --> MEMSVC
    MEMSVC --> OBJ
    MEMSVC --> VDB
    AGENT --> TOOLS
  end

  subgraph GPU[GPU Machine • Compute Plane]
    LLM[DeepSeek Reasoning Model\n(33B+)\nLLM Server]
  end

  %% ===== Links between machines =====
  VCLIENT -->|text transcript| APIGW
  AGENT -->|complex task| LLM
  LLM -->|completion| AGENT
  AGENT -->|final response| APIGW
  APIGW -->|text for speech| TTS
```

## 7. Voice Request Flow

```mermaid
sequenceDiagram
  participant User
  participant Pi as Raspberry Pi (Voice)
  participant Mac as Mac mini (OpenClaw + Models)
  participant DB as Memory (SQL+Vector+Docs)
  participant GPU as GPU Box (DeepSeek)

  User->>Pi: "Hey Skyra..." (audio)
  Pi->>Pi: Wake word detect
  Pi->>Pi: STT (audio → text)
  Pi->>Mac: POST /voice {text, timestamp}

  Mac->>Mac: Project/Intent classify
  Mac->>DB: Retrieve relevant context
  DB-->>Mac: Context snippets

  Mac->>Mac: Route task
  alt Simple or coding task
    Mac->>Mac: Use local model (Llama or Qwen)
  else Complex reasoning task
    Mac->>GPU: Send to DeepSeek
    GPU-->>Mac: Completion
  end

  Mac->>DB: Write memory update
  Mac-->>Pi: Response text
  Pi->>Pi: TTS
  Pi-->>User: Spoken response
```

## 8. Component Responsibilities

### 8.1 Raspberry Pi – Voice Node

**Purpose**: Always-on audio interface.

**Services**:

- Wake word detection
- Speech-to-text (STT)
- Text-to-speech (TTS)
- Voice client that sends text to Mac mini

**Characteristics**:

- Lightweight compute
- Always powered on
- Local network only

### 8.2 Mac mini – Control Plane

**Purpose**: Orchestration, memory, APIs, tools, and fast local models.

**Services**:

- API gateway (/chat, /voice, /tools, /memory)
- OpenClaw agent runtime (orchestrator + router)
- Project classifier
- Model router
- Memory service
- Tool execution engine
- Databases
- Local conversational model
- Local coding/tool model

**Local Models**:

| Model                 | Role                                 |
| --------------------- | ------------------------------------ |
| Llama 3.1 8B Instruct | Conversational interface and routing |
| Qwen2.5-Coder 7B      | Coding and tool execution            |

**Datastores**:

- Relational DB (projects, events, preferences)
- Vector DB (embeddings)
- Object storage (documents)

### 8.3 GPU Machine – Compute Plane

**Purpose**: Heavy reasoning and large-model inference.

**Services**:

- DeepSeek-Coder (33B+)
- LLM server (vLLM, TGI, or Ollama)

**Characteristics**:

- Dedicated GPU
- High VRAM
- Private network access only

**Model Role**:

| Model               | Purpose              |
| ------------------- | -------------------- |
| DeepSeek-Coder 33B+ | Main reasoning brain |

## 9. Memory Architecture

### 9.1 Object Store (System of Record)

**Structure**:

- `.skyra/projects/{project}/`
- HEAD.json (current commit pointer)
- state.json (materialized current state)
- commits/ (immutable commit objects)
- attachments/ (files, documents)

**Usage**:

- Versioned project state via commit objects
- AI modifications through explicit commits only
- File-based storage (local) or S3/MinIO (distributed)

### 9.3 Vector Store (Derived Data)

**Stores**:

- Embedded project state snapshots
- Embedded documents and attachments
- Semantic index for fast retrieval

**Characteristics**:

- Can be rebuilt from object store
- Not source of truth
- Used for semantic search only

## 9.4. Retrieval Strategy (Commit + Semantic + Temporal)

1. Classifier determines project domain.
2. Vector store retrieves semantically similar content with temporal metadata.
3. Object store provides recent commit context.
4. Results are re-ranked by:
   - Commit recency
   - Semantic similarity
   - Project relevance
   - Temporal importance (recent events weighted higher)
5. Top results injected into LLM prompt with timestamp context.

**Rule**: Vector store finds relevant content; object store provides authoritative state; temporal metadata adds context.

## 10. OpenClaw Integration

OpenClaw runs on the Mac mini as the central orchestrator and model router.

**Agent pipeline**:

1. Receive message
2. Run project and intent classifier
3. Retrieve memory context
4. Route task:
   - Llama (conversation)
   - Qwen (coding/tools)
   - DeepSeek (GPU reasoning)
5. Execute tools if needed
6. Write memory updates
7. Return final response

## 10. Network Layout

### Logical network

- All machines on same LAN/VLAN
- Token or mTLS between services
- No public exposure of GPU machine

### Trust zones

| Zone    | Machine      | Role                                 |
| ------- | ------------ | ------------------------------------ |
| Edge    | Raspberry Pi | Voice input/output                   |
| Control | Mac mini     | Orchestration + memory + fast models |
| Compute | GPU box      | Deep reasoning model                 |

## 11. Deployment Strategy

### Phase 1 (single machine)

Mac mini runs:

- Llama conversational model
- Qwen coding model
- API
- Memory
- OpenClaw agent

### Phase 2 (two machines)

- GPU box hosts DeepSeek
- Mac mini runs control plane + local models

### Phase 3 (three machines)

- Raspberry Pi handles voice
- Full modular architecture

## 12. Security Baseline

- Token-based service authentication
- Separate service accounts
- Tool allow-list
- Audit logs
- Encrypted backups

## 13. End-State Role Assignment

| Machine      | Role                                             |
| ------------ | ------------------------------------------------ |
| GPU Box      | DeepSeek reasoning model                         |
| Mac mini     | OpenClaw agent, APIs, memory, tools, fast models |
| Raspberry Pi | Voice interface                                  |

## 14. Example Capabilities

- "What did I decide about the Tekkit backups last week?"
- "Switch to work mode—draft a SOC2 response."
- "Gym mode—suggest next week's lifts."
- "Music mode—ideas for SKANZ Vol. 3."
- "Server mode—summarize crash logs."

For every request:

**Input**:

- User message
- Recent chat history
- Project registry

**Output**:

- project_id
- intent
- confidence score

If confidence is low:

- Ask clarification
- Or search across multiple projects

## 15. Network Layout

### Logical network

- All machines on same LAN/VLAN
- Token or mTLS between services
- No public exposure of GPU machine

### Trust zones

| Zone    | Machine      | Role                   |
| ------- | ------------ | ---------------------- |
| Edge    | Raspberry Pi | Voice input/output     |
| Control | Mac mini     | Orchestration + memory |
| Compute | GPU box      | Model inference        |

## 16. Deployment Strategy

### Phase 1 (single machine)

Mac mini runs:

- Local model
- API
- Memory
- Agent

### Phase 2 (two machines)

- GPU box hosts model
- Mac mini runs control plane

### Phase 3 (three machines)

- Raspberry Pi handles voice
- Full modular architecture

## 17. Security Baseline

- Token-based service authentication
- Separate service accounts
- Tool allow-list
- Audit logs
- Encrypted backups

## 18. Telemetry and Monitoring

### 18.1 Model-Level Metrics

- Request count per model (Llama/Qwen vs DeepSeek)
- Response times and token counts
- GPU utilization during DeepSeek calls
- Memory usage per inference request
- Cache hit rates for local models

### 18.2 User Interaction Metrics

- Voice-to-text latency
- End-to-end request/response time
- Routing decision accuracy (was DeepSeek escalation needed?)
- User satisfaction signals (voice follow-ups, task completion)

### 18.3 System Health

- Network latency between machines
- Database query performance
- Model warm-up times
- Error rates and fallback triggers

### 18.4 Cost Tracking

- Energy consumption per machine
- GPU compute time cost estimates
- Storage usage trends
- Peak vs average resource utilization

### 18.5 Implementation

- Prometheus metrics on each machine
- Simple dashboard in the Mac mini control plane
- Alerts when DeepSeek usage spikes unusually
- Weekly summaries of model usage patterns

### 18.6 Key Focus: Routing Efficiency

Track the system's ability to correctly choose local vs GPU models and analyze the tradeoffs between response time and accuracy.

## 18. End-State Role Assignment

| Machine      | Role                                             |
| ------------ | ------------------------------------------------ |
| GPU Box      | DeepSeek reasoning model                         |
| Mac mini     | OpenClaw agent, APIs, memory, tools, fast models |
| Raspberry Pi | Voice interface                                  |

## 18. Example Capabilities

- "What did I decide about the Tekkit backups last week?"
- "Switch to work mode—draft a SOC2 response."
- "Gym mode—suggest next week's lifts."
- "Music mode—ideas for SKANZ Vol. 3."
- "Server mode—summarize crash logs."

## 19. Planned Extensions

The following sections describe high-level architectural concepts that will be designed in more detail in future iterations. These are included to guide the long-term evolution of the system without constraining the initial implementation.

### 19.1 TV Node Architecture (Planned)

The system may include a dedicated TV node consisting of a small computer (mini PC or Raspberry Pi) connected to a television via HDMI. This node will run a Jarvis/Skyra agent and act as the primary media execution environment.

In this model:

The TV becomes a display only.

All media playback occurs on the local TV node through a browser or local apps.

Jarvis controls playback by interacting directly with the local execution environment.

High-level goals:

Avoid limitations of smart TV and streaming device APIs.

Allow full control of streaming services through web interfaces.

Enable the system to read what is currently playing.

Provide reliable, scriptable media control.

Hardware selection, browser automation, and media control logic will be designed in a later iteration.

### 19.2 Mobile Interaction via Progressive Web App (Planned)

Jarvis will expose a secure web interface that can be installed on a phone as a Progressive Web App (PWA). This will provide a native-like experience without requiring a dedicated mobile application.

High-level goals:

Chat-style interface for natural language commands.

Quick action buttons for common tasks.

Real-time system and device status.

Optional push notifications.

The exact UI design, authentication model, and notification system will be defined in a future version.

### 19.3 Voice Authorization Model (Planned)

The voice subsystem will include a speaker-aware authorization layer so that only approved users can control the system.

High-level components:

Wake word detection.

Speaker identification.

Command authorization based on recognized voice profiles.

This ensures that unauthorized users cannot issue commands.
High-risk or destructive actions may require additional spoken confirmation.

Specific models, thresholds, and security policies will be defined in a later design phase.

### 19.4 Streaming Device Integration Strategy (Planned)

The system will support multiple types of media devices, using the most appropriate control method for each platform.

High-level strategy:

Use direct APIs where available (e.g., Roku power or input control).

Use casting or device-level control for platforms like Android TV.

Prefer a local TV node for full automation when APIs are limited.

Media endpoints will be treated as device nodes with defined capabilities.
Detailed control logic for each platform will be designed later.

### 19.5 External Device Control Model (Planned)

Not all devices will run native Jarvis agents. Some will be controlled through network protocols or automation bridges.

High-level concepts:

Certain devices (TVs, smart devices, streaming boxes) are controlled via:

Local network APIs

Casting protocols

Automation bridges

The control plane maintains a centralized device registry.

High-level intents are mapped to device-specific actions.

Multiple device types are supported under a unified abstraction.

The device capability schema and integration patterns will be defined in a future design phase.

### 19.6 Remote Access Model (Planned)

Remote clients (such as phones or laptops outside the home network) will access Jarvis through a secure overlay network.

High-level concept:

The control plane is not exposed directly to the public internet.

Remote devices connect through a secure mesh or VPN-style network.

All interactions occur over authenticated, encrypted channels.

Authentication, session management, and final network topology will be specified in a later iteration.
