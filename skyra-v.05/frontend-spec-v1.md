# Frontend Spec — v.1 Runtime

Dante — this replaces the original spec from May 1st. The runtime architecture has changed. The core principle is the same: the frontend renders the runtime, it doesn't model it. But the runtime underneath is different now, and some of your questions surfaced things that changed how the architecture works.

This doc answers your questions, updates the wire protocol, and describes what you're rendering now.

## 0. What Changed Since the Last Spec

The v.05 runtime — hardcoded routing through NewThread → Exchange → Self → Think → Act — has been replaced. v.1 is a two-phase recursive traversal: Observe (context accumulates on the way down) and Express (execution and compression on the way back up). Every Reality implements the same interface:

```go
type Reality interface {
    ID() string
    Core() *Base
    Create(r *Relation) Reality
    Realize(r *Relation) string
    Observe(r *Relation)
    Express(r *Relation) string
}
```

Every Reality carries:

```go
type Base struct {
    Weight        float64
    Usage         int
    LastUsed      time.Time
    Relationships map[string]Reality
    Expressors    map[string]Reality
}
```

What this means for the frontend: the internal structure of a being is no longer a fixed Think/Act stack with hardcoded operator lists. It's a weighted topology that changes through use. The universe state object reflects this.

## 1. Your Questions — Answered

### Who computes the region on a Relation?

The backend. Always. The runtime knows whether something is a user or a being from the genome (`~type user` vs `~type llm`). The universe state payload tags each being with its type. The frontend does not compute relational state — it renders what the backend sends.

There are two views:

**User present** — michael's view as a being in the system. His active exchanges, his conversation entries, his peers, what's in front of him right now. Scoped and perspectival.

**Universe present** — the full state of the runtime. Every being, every exchange, every thread, the topology. No perspective. No privacy.

The frontend toggles between the two. One is a being's view. The other is the view from the top.

### WebSocket authentication

First-message auth. On connect, the first message from the client is an auth payload with the token. The server validates and responds with a handshake confirmation before any state flows. No tokens in query strings (leaks to logs). No cookies (requires public domain). No CORS issues.

```json
// Client sends first:
{ "type": "auth", "token": "..." }

// Server responds:
{ "type": "auth_ok" }
// or
{ "type": "auth_fail", "reason": "..." }
```

State flows only after `auth_ok`.

### Delta vs snapshot for topology

Snapshot on connect and reconnect. Deltas afterwards. You were right — option 2 is the correct pattern. The thread graph is append-only (edges and members only add, never remove). Deltas are safe. Full snapshot is the safety net.

### Reconnect behavior

Option 2: fresh snapshot on reconnect, drop in-flight relations.

The runtime does not know the frontend exists. The runtime does not care about connection state. It runs whether anyone is watching or not. The WebSocket is a lens — it observes the Universe's state. When the lens disconnects, nothing changes in the runtime. When it reconnects, the Universe sends what it always sends: the current state. A fresh snapshot.

No replay. No server-side event buffering. No client ack tracking. The frontend missed what it missed. The runtime didn't miss anything.

Replay of historical state is a separate system for later — a recording layer outside the runtime that logs what the Universe observed. That's a product feature, not a runtime feature.

## 2. Vocabulary Mapping (Updated)

| Frontend | Runtime Term | What It Is |
|----------|-------------|------------|
| Component name (your choice) | **Being** | An entity in the system — LLM-backed, user-backed, agent-backed, or process-backed. Has identity, purpose, relationships. |
| Component name (your choice) | **Reality** | The core interface. Everything implements it. One type, recursive composition. |
| Thread graph | **Thread** | Each thread has its own graph — members and edges. Multiple threads means multiple topology graphs. |
| Reality graph | **Topology** | The weighted graph of Relationships and Expressors inside each being. Replaces the old fixed composition tree. |
| Relation | **Relation** | Same concept. The observer that descends through the topology. Carries: Origin, ID, ThreadID, Impulse, Depth, Budget, Trace. |
| Exchange | **Exchange** | Conversation history between two beings. Entries are append-only. |

Use the runtime names on the wire. Your React components can be called whatever you want.

## 3. Runtime Architecture (What You're Rendering)

The runtime is a recursive descent with two phases.

```
Relation enters (carrying: impulse from michael)
  ↓ Thread     attaches thread context
    ↓ Exchange  attaches conversation history
      ↓ Self    attaches identity, stress fields
        ↓ weighted traversal along Relationships
          ↓ memory (weight 0.8) — attaches content
            ↓ connected memory (weight 0.5) — attaches content
              ↓ weight exhausted — bottom
            ↑ synthesizes accumulated content
          ↑ compresses, updates weights
        ↑ Expressors fire (Think → Provider → LLM call)
      ↑ Self integrates, updates weights
    ↑ Exchange records entry
  ↑ Thread routes result
Relation returns
```

Key differences from the old spec:

- **No hardcoded Think/Act stack.** Beings have Relationships (context accumulation) and Expressors (execution). What used to be Think and Act are now Expressors on Self — they fire by weight, not by hardcoded sequence.
- **No operator lists.** Operators (bash, search, browse) are Realities in the being's Relationships map. They surface into the being's context by weight, not by explicit listing. The being thinks its way into capabilities rather than being handed a tool list.
- **Weight-driven depth.** How deep the traversal goes depends on the Relation's content and the edge weights. A deep thought goes deep because the content activates heavy edges. A shallow greeting returns shallow.
- **Thread still owns the loop.** When skyra responds to michael by addressing builder, Thread routes the Relation to builder. The frontend doesn't manage this — it receives events as they happen.

### Skyra's Role

Skyra is not a regular being. She maintains the universal principles — the curvature of the space all other beings operate inside. During boot she is maximally active, shaping the topology. As the principles take hold, she fades to resting state. If the topology tears, she reactivates.

The frontend can reflect this: Skyra's activity level is an indicator of system stability. High activity = system is shaping. Low activity = system is stable.

## 4. What the Runtime Exposes — The Universe State

The Universe observes itself. It sends a Relation to itself, and the return carries the full state. The frontend receives this observation — not a database dump, but what the Universe sees when it observes its own topology.

### Full snapshot (on connect and reconnect)

```json
{
  "type": "universe",
  "ts": 1716400000,
  "payload": {
    "beings": [],
    "threads": [],
    "exchanges": [],
    "economics": {},
    "topology": {}
  }
}
```

### beings[]

```json
{
  "name": "skyra",
  "type": "llm",
  "identity": "I hold the world together.",
  "purpose": "I think, respond, and relate on behalf of the system.",
  "status": "active",
  "peers": ["michael", "louise", "claude", "builder", "philosopher"],
  "weight": 0.95,
  "relationships": [
    { "target": "michael", "weight": 0.9, "usage": 142 },
    { "target": "builder", "weight": 0.7, "usage": 58 },
    { "target": "bash", "weight": 0.3, "usage": 5 }
  ],
  "expressors": [
    { "target": "think", "weight": 0.9 },
    { "target": "act", "weight": 0.85 }
  ],
  "memories": {
    "items": [
      { "filename": "1716400000.md", "content": "michael prefers direct answers" }
    ],
    "skills": [
      { "name": "code-review", "content": "when reviewing code, focus on..." }
    ]
  }
}
```

The being's topology is now visible — Relationships and Expressors with weights. The frontend can render the internal graph of any being. Weights change over time as the being uses things.

### threads[]

```json
{
  "id": "a1b2c3d4",
  "created_by": "michael",
  "active": true,
  "members": ["michael", "skyra", "louise"],
  "edges": [
    { "from": "michael", "to": "skyra" },
    { "from": "skyra", "to": "louise" }
  ]
}
```

Unchanged from v.05. Edges are append-only.

### exchanges[]

```json
{
  "key": "michael:skyra",
  "parties": ["michael", "skyra"],
  "active": true,
  "entries": [
    { "index": 0, "from": "michael", "content": "what about the server?", "ts": 1716400000 },
    { "index": 1, "from": "skyra", "content": "checking now", "ts": 1716400005 }
  ]
}
```

Unchanged from v.05. Every exchange between every pair of beings.

### economics{}

```json
{
  "fields": {
    "inference_calls": 42,
    "tokens_used": 18500
  }
}
```

Not yet enforced in the descent. When active, tracks token spend, inference call count, budget remaining per being.

### topology{}

Replaces the old `reality_graph`. The weighted composition structure — what contains what, with what weight.

```json
{
  "id": "universe",
  "type": "Universe",
  "children": [
    {
      "id": "skyra",
      "type": "Self",
      "weight": 0.95,
      "relationships": [
        { "id": "michael-model", "type": "Relationship", "weight": 0.9 },
        { "id": "memory-cluster", "type": "Memory", "weight": 0.7 }
      ],
      "expressors": [
        { "id": "think", "type": "Think", "weight": 0.9 },
        { "id": "act", "type": "Act", "weight": 0.85 }
      ]
    }
  ]
}
```

This is recursive — any node can contain relationships and expressors, which themselves contain relationships and expressors. The frontend can render this as a tree, a graph, or a nested view.

## 5. Delta Events (After Initial Snapshot)

```json
{
  "type": "entry",
  "ts": 1716400010,
  "payload": {
    "exchange": "michael:skyra",
    "index": 2,
    "from": "skyra",
    "content": "server is healthy"
  }
}
```

Delta types:

| Type | What Happened |
|------|--------------|
| `entry` | New exchange entry appended |
| `edge` | New edge in thread graph |
| `being` | New being grown mid-flight (full being snapshot) |
| `weight` | Weight changed on a relationship or expressor |
| `memory` | New memory written |
| `topology` | Topology changed (new being adds a subtree) |
| `error` | System error (origin, message) — recoverable |

New in v.1: `weight` events. The frontend can animate weight changes in real time — edges getting stronger or weaker as the being uses them.

## 6. Message Envelope

Every message over the WebSocket follows one format:

```json
{
  "id": "uuid",
  "type": "string",
  "ts": 0,
  "payload": {}
}
```

Types: `auth`, `auth_ok`, `auth_fail`, `universe`, `entry`, `edge`, `being`, `weight`, `memory`, `topology`, `error`, `impulse`.

The `impulse` type is client → server: the user sending input.

```json
{
  "id": "uuid",
  "type": "impulse",
  "ts": 1716400020,
  "payload": {
    "origin": "michael",
    "content": "check the server",
    "target": "skyra"
  }
}
```

Target is optional. If omitted, the runtime routes by weight (who the user has been talking to most recently).

## 7. Two Views

**User present** — michael's perspective. His active exchanges, his peers, what's in front of him. The frontend filters the universe state to show only what's relevant to michael. Entries are append-only. No turns, no request/response. Michael can send three messages in a row. Skyra can talk to louise for ten exchanges before coming back.

**Universe present** — everything. Every being with full internal state, every thread, every exchange, the full topology with weights. No perspective. The view from the top.

The frontend toggles between these two. One is "I'm michael." The other is "I'm watching the universe."

## 8. What Doesn't Exist Yet (Backend)

1. **v.1 Reality implementations** — The interface is defined. The old implementations are deleted. New implementations using Observe/Express and weighted traversal are being written. The universe state object shape is stable — the internals that produce it are in progress.

2. **WebSocket device** — Needs to be rewritten for the v.1 architecture. Same port (8080), same pattern (push state on every resolve).

3. **Weight event emission** — New for v.1. Hooks needed at weight update points so the frontend can show topology changes in real time.

4. **Serialization** — JSON serialization for the new structs (Base, Relationships, Expressors).

## 9. Frontend Obligations

Build against the universe state object defined in this doc. The shape is stable. What you need:

1. **Connection layer** — WebSocket client. First-message auth. Snapshot on connect. Delta application. Fresh snapshot on reconnect.

2. **User present view** — Render michael's active exchanges. Filter beings to peers. Show conversation entries. Handle the user sending impulses.

3. **Universe present view** — Render the full topology. Every being, every exchange, every thread. Weighted graph visualization for being internals (Relationships and Expressors with weights). This is the new thing — the old spec didn't have weighted internal topology.

4. **Being detail** — When a being is selected, show its Relationships and Expressors with weights, its memories, its skills. Weights change over time — animate or indicate changes.

5. **Thread graph** — Nodes are members, edges are connections. Append-only. Multiple threads means multiple graphs.

6. **Skyra status indicator** — Skyra's activity level indicates system stability. High weight activity = system is shaping. Low = stable.

## 10. Stack

- **Framework**: React + shadcn/ui (copy-paste-own, no vendor dependency)
- **Transport**: WebSocket on port 8080
- **State**: Client-side, rebuilt from snapshot on reconnect
- **Rendering**: Component composition from proven blocks, not generation

## 11. Timeline

| Week | Backend | Frontend |
|------|---------|----------|
| May 22–28 | v.1 implementations in progress. Wire format locked (this doc). | Build against mock data using this spec. Connection layer. |
| May 29–Jun 1 | WS device live. Real beings producing real state. | Cut over from mock to real WS. Integration testing. |
| Jun 1 | Demo. | Demo. |

Mock data: generate JSON matching the schemas in this doc. The shape won't change.
