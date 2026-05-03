# Skyra: A Physics Engine for Agentic Cognition

## Thesis

A physics engine is a better runtime for an agent than a pipeline.

Pipelines produce output. Physics produces development. Every existing agentic framework — LangChain, CrewAI, AutoGen, Hermes, OpenClaw — is an orchestration system: it coordinates behavior by assembling prompts, calling models, parsing responses, and routing results. The agent is a function that receives input and produces output. There is no world. There is no interiority. There is no development over time.

Skyra inverts this. Instead of engineering agent behavior from the outside, Skyra creates an environment with invisible laws — physics — and places beings inside it. Behavior emerges from living inside those laws. The agent doesn't execute a pipeline. It inhabits a world.

## Core Abstraction

Everything in Skyra implements one interface:

```go
type Reality interface {
    ID() string
    Create(r *Relation) Reality
    Realize(r *Relation) string
}
```

A being is a Reality. A device is a Reality. A thread is a Reality. An LLM provider is a Reality. The world itself is a Reality. There is no second type system, no special cases, no orchestrator. Every node creates itself from a Relation and realizes itself when a Relation passes through it.

## The Relation

The Relation is the protagonist of the system. It is not a message — it is an experience being constructed.

```go
type Relation struct {
    ID         string
    Origin     string
    ThreadID   string
    Impulse    string
    Parsers    map[string]Parser
    Realities  map[string]Reality
    Collecting bool
    Exports    map[string]any
}
```

A Relation enters the world empty. As it descends through layers — Thread, Exchange, Being, Think, Act — each layer attaches a Parser to it. A Parser is a function that renders that layer's contribution to the being's present. By the time the Relation reaches the inference provider at the bottom, its Parsers map contains everything the being should experience. The provider concatenates them. That's the present.

No single layer builds the context. The journey builds it.

## Recursive Descent

The system resolves through recursive descent. A Relation enters at the top (Universe) and descends through each Reality until it reaches a terminal node (the inference provider or a device). Each layer the Relation passes through:

1. Reads the Relation
2. Updates its own state (Thread records an entry, Exchange logs the message)
3. Attaches its Parser to the Relation (adding its context to the eventual present)
4. Passes the Relation deeper

The recursion terminates at the provider — the LLM call. The response then travels back up, parsed at each layer (Act parses the tag protocol, NewThread routes the response as a new Relation, which descends again toward its target).

This means the system is a recursive loop:
```
Relation descends → accumulates context → hits provider → response ascends → becomes new Relation → descends again
```

The loop continues until the Relation returns to a user's device (terminal or WebSocket) and waits for new input. A single user message can trigger multiple full descents — being A responds to the user by addressing being B, which triggers a new descent to B, whose response triggers a descent back to A, and so on until the chain resolves.

Every descent constructs a unique present. Two beings receiving the same message will have different presents because the Relation took different paths to reach them. The being doesn't construct its own experience — the path does.

This means:
- Every being gets its own present, built by the specific path the Relation took to reach it
- 50 beings can exist in the same world without paying each other's context cost
- The context window never fills with irrelevant noise because irrelevant things were never attached
- Context is not managed — it is constructed per descent

## World Structure

```
Universe → NewThread → Exchange → Being (Self) → Think → Act → Provider
```

### Universe
The top-level Reality. Holds the thread gate, economics, and devices. On every resolve, collects the full state of the world as JSON and broadcasts it via WebSocket. The world is observable from outside without the beings knowing they're being observed.

### NewThread
The routing loop. Manages the set of beings, their access permissions, active threads, and the continuous cycle of: receive impulse → route to being → receive response → route response as new impulse. The loop runs until a Relation returns to its origin (a user device) and waits for new input.

### Exchange
The history between two beings. Records every entry with timestamp. Enforces context continuity: if a being tries to leave a conversation to talk to someone else without carrying context via `<ref>`, the exchange blocks the crossing and returns an error. This is not a bug — it's a theory of coherence. You cannot context-switch without carrying what you learned.

### Being Types

**Self** — an LLM being with two-layer cognition (Think + Act). Has identity, purpose, relationships, memory, and skills. The primary being type.

**User** — a human behind a device. Routes impulses to/from the terminal or WebSocket. No cognition layers — the human provides those.

**Agent** — an external process (e.g., Claude Code) invoked via subprocess with session persistence. Treated as a peer, not a tool.

**CLI** — a shell command executor. Routes impulses as commands, returns stdout.

## Two-Layer Cognition

Every Self being has an inner layer (Think) and an outer layer (Act). This is not chain-of-thought. It is a private interior that the world cannot see, followed by a public exterior that the world receives.

### Think (Inner Layer)

Private thought. No one sees this — not other beings, not the user, only the debug logs.

- 5-pass budget per invocation
- Operators available: `recall` (search memory), `remember` (write memory), `skill` (load a skill file)
- One protocol emission per pass: either an operator call or `<surface-thought>`
- Time pressure increases as budget depletes
- Thought history persists across exchanges — the being remembers what it thought before

When the being emits `<surface-thought>`, the thought passes to the outer layer. The inner layer is done.

### Act (Outer Layer)

Public speech. The being addresses a peer using the tag protocol: `<target>message</target>`.

- Receives the surfaced thought as "your inner thought"
- Must emit exactly one tagged message per response
- Protocol enforcement: if the response doesn't follow the tag format, the system retries with a warning
- Self-routing detection: if the being addresses itself, the system retries
- 3 retry budget before exhaustion

### Think-Back

If the outer layer emits `<think>content</think>` instead of addressing a peer, the system returns the being to its inner layer with a fresh budget. The being can retreat from action back into thought. This treats thinking not as preprocessing, but as something you might need to return to after you've tried to speak and realized you weren't ready.

No other system implements this.

## Physics

Physics are realities that live in every world's hashmap, fire on every Relation that passes through, and are invisible to the beings that inhabit the world. A being never addresses a physics reality. It never sees one in its peers. It doesn't know they exist. But every Relation carries the accumulated weight of every physics layer it passed through.

### Thread (implemented)
Records what happened — who said what, in what order. Creates threads, tracks members, builds the graph of who has spoken to whom. Fields applied: thread ID, exchange history, active exchanges.

### Economics (wired, not yet applied)
Four budgets: token (finite cognition per turn), memory (finite active window), thread (max open conversations), reproduction (creating new beings costs accumulated experience). Fields applied: tokens remaining, memory pressure, open thread count.

### Salience (designed, not yet implemented)
Ambient semantic matching. On every think pass, the current thought is compared against peer identities and skill descriptions via embedding similarity. Only relevant peers and skills surface into the being's present. The being doesn't search — things come to mind. This is not retrieval. It is attention shaped by the world.

### Governance (designed, not yet implemented)
How beings make collective decisions when an action affects shared space. Proposals, thresholds, votes. Applied invisibly — a being doesn't call governance, it proposes something and the world decides.

## The Genome

```
device ~name macbook ~type macos
component ~name terminal ~type stdin ~device macbook
component ~name ws ~type websocket ~port 8080 ~device macbook
component ~name openrouter ~type llm ~model anthropic/claude-sonnet-4-5 ~device macbook

being ~name skyra ~type llm ~identity I hold the world together. ~purpose I think, respond, and relate on behalf of the system. ~devices macbook ~entrypoints openrouter ~relationships michael,louise,claude
being ~name michael ~type user ~identity I build Skyra. ~purpose I decide what matters. ~devices macbook ~entrypoints terminal,ws ~relationships skyra
```

A flat file that declares a world. Devices, components, beings — each one line. The runtime reads this, bootstraps every Reality, wires them together, and starts the loop. The genome is the world's DNA. Adding a being is adding a line.

## Device Layer

The machine is a world. MacOS is a world Reality with components inside it — terminal (stdin), WebSocket (browser bridge), LLM provider (OpenRouter). A being's entrypoints determine which components it can be reached through. A user behind a terminal and a user behind a browser are the same Being type with different device paths.

## Context Continuity (Ref Crossing)

When a being leaves one conversation to enter another, it must carry context:

```
<louise>I want to discuss this <ref>michael:0-3</ref></louise>
```

This brings entries 0–3 from the being's exchange with michael into the new conversation with louise as private context. If the being attempts to cross without a ref, the Exchange blocks it and returns an error explaining why.

This enforces coherence architecturally. A being cannot dissociate — it cannot drop one context and pick up another without bridging them. Every other multi-agent system allows clean-slate context switches. Skyra says that's not agency, that's fragmentation.

## State Observation (Universe Collecting)

On every resolve, the Universe fires a collecting pass through every Reality. Each Reality exports its current state — beings export their snapshots, threads export their graphs, exchanges export their histories. The Universe assembles this into a JSON structure and broadcasts it via WebSocket.

This means:
- A browser client can render the full reality graph in real time
- Every being's state, every exchange, every thread is observable
- The beings don't know they're being observed — collecting is invisible to them
- Development over time is visible and loggable

## Self-Extension

Beings can create new beings at runtime via the `grow` command. A being generates a genome line, the system parses it, creates the Reality, wires it into the world. The new being gets the full stack — cognition layers, memory, operators, exchange access.

Constraints:
- A being cannot modify physics (it is bound to its world's laws)
- A being cannot modify beings above it in the hierarchy
- A being can create peers and structures beneath it
- Economics (when active) will gate creation behind accumulated experience

## What This Solves

| Industry pain point | Skyra's structural answer |
|---|---|
| Context window overflow | Each being gets its own present, built by descent — irrelevant context never attached |
| Multi-agent incoherence | Exchange enforces continuity; ref crossing prevents dissociation |
| Skills don't transfer | Salience matches on thought content, not task type |
| No development over time | Physics produce emergent behavior; think history persists; memory accumulates |
| Orchestration complexity | No orchestrator — the world routes through physics, beings don't coordinate externally |
| Scaling to many agents | Beings don't pay each other's context cost — 50 beings, each sees only its own present |

## Implementation

- Language: Go
- Inference: OpenRouter (any model, currently Claude Sonnet 4.5)
- State broadcast: WebSocket on port 8080
- Memory: filesystem (~/.skyra/beings/{name}/memories/)
- Skills: filesystem (~/.skyra/beings/{name}/skills/)
- Logging: per-being, per-layer (system.log, {being}/inner.log, {being}/outer.log)
- External agents: subprocess invocation with JSON output and session persistence

## Status (v.05 — May 2026)

**Working:** Reality interface, genome bootstrap, two-layer cognition, think-back, exchange continuity, ref crossing, multi-party threading, universe collecting, WebSocket broadcast, grow, recall/remember/skill operators, Agent type (Claude Code integration with session persistence).

**Wired but inactive:** Economics.

**Designed:** Salience, Governance, parser-per-reality matrix.

**Next phase:** Self-extension pipeline (being writes code, compiles, attaches), Salience implementation, surface area parity with existing frameworks (tools, integrations, multi-platform IO).
