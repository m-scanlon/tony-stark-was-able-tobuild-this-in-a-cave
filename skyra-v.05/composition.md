# Composition

## Current Build

What's running today. Every node implements `Reality` (`ID`, `Create`, `Realize`).

### Reality Composition Tree

```
Universe                                    universe.go:22
├── NewThread                               newthread.go:10
│   ├── Exchange                            exchange.go:11
│   ├── Devices
│   │   ├── MacOS                           macos.go:10
│   │   └── Provider "openrouter"           llm.go:8
│   └── Beings
│       ├── Self "skyra"                    self.go:5
│       │   ├── Being                       being.go:24
│       │   │   ├── .Identity
│       │   │   ├── .Purpose
│       │   │   ├── .Relationships
│       │   │   └── .Home (~/.skyra/beings/skyra/)
│       │   ├── Think                       think.go:13
│       │   │   ├── .LLM → Provider
│       │   │   ├── .History []ThoughtSection
│       │   │   └── Operators
│       │   │       ├── Recall              recall.go:9
│       │   │       ├── Remember            remember.go:11
│       │   │       └── Skill               skill.go:9
│       │   └── Act                         act.go:9
│       │       ├── .LLM → Provider
│       │       └── Operators
│       │           └── Plan (stub)         plan.go:3
│       │
│       ├── Self "louise"                   (same shape as skyra)
│       │
│       └── User "michael"                  user.go:5
│           ├── Being                       being.go:24
│           └── Device → MacOS
│
├── Economics (structural, not enforced)    economics.go:8
└── OnResolve → Universe.Realize(collecting)
```

### Relation Flow

The `Relation` (`relation.go:10`) is the mutable message — it accumulates parsers and context as it descends.

```
User types at terminal
        │
        ▼
MacOS.Realize                               macos.go:26
  reads stdin, returns raw text
        │
        ▼
Impress("michael", raw)                     relation.go:36
  creates Relation{Origin:"michael", Impulse:raw}
        │
        ▼
Universe.Realize(rel)                        universe.go:22
        │
        ▼
NewThread.Realize(rel)                       newthread.go:44
  ├── finds/creates Thread
  ├── attaches thread parser
  ├── puts all beings on rel.Realities
  │       │
  │       ▼
  ├── Exchange.Realize(rel)                  exchange.go:47
  │     ├── peels target from impulse (or uses rel.ID)
  │     ├── finds/creates Conversation
  │     ├── checks <ref> for context crossing
  │     ├── blocks crossing without ref
  │     ├── records Entry
  │     ├── attaches exchange + conversation parsers
  │     │       │
  │     │       ▼
  │     └── being.Realize(rel)  ← routes to target
  │             │
  │             ▼
  │     Self.Realize(rel)                    self.go:19
  │       ├── puts Being on rel.Realities
  │       │
  │       ├── THINK PHASE                    think.go:42
  │       │   ├── attaches system prompt, being (inner parse), operators, history
  │       │   ├── loop (up to 5 passes):
  │       │   │   ├── attaches time pressure + exchange
  │       │   │   ├── Provider.Realize(rel)  → LLM call
  │       │   │   │     └── derivePresent: fires all parsers → system + present
  │       │   │   │         inference.Call(system, present)     inference.go:40
  │       │   │   ├── if <surface-thought>  → done, return thought
  │       │   │   ├── if outer op (e.g. <plan>) → blocked, system message
  │       │   │   ├── if inner op (<recall>/<remember>/<skill>) → fire operator
  │       │   │   │     operator result goes into exchange for next pass
  │       │   │   └── loop
  │       │   └── returns inner thought
  │       │
  │       ├── ACT PHASE                      act.go:24
  │       │   ├── attaches system prompt, being (outer parse), inner thought
  │       │   ├── loop (up to 3 attempts):
  │       │   │   ├── Provider.Realize(rel)  → LLM call
  │       │   │   ├── ParseResponse → extract <target>message</target>
  │       │   │   ├── if no tags → protocol violation, retry
  │       │   │   ├── if self-route → blocked, retry
  │       │   │   └── sets rel.ID = target, rel.Impulse = message
  │       │   └── returns message
  │       │
  │       └── rel.Origin = self.id
  │
  ├── Thread.Spread(from, to)                newthread.go:182
  ├── OnResolve()                            main.go:53
  │     └── Universe.Realize(collecting)     universe.go:22
  │           walks entire tree, each Reality exports → JSON snapshot
  │
  └── LOOP continues                         newthread.go:83
        rel.ID now points to next target
        parsers reset, descend again
        (if target is User → MacOS prints response, reads next input)
        (if target is another Self → Think/Act fires for that being)
```

### Parser Accumulation

Each Reality attaches its own parser to the Relation as it passes through. The Provider fires them all.

```
Provider.derivePresent(rel)                  llm.go:41
  ├── r.Parsers["system"]     → Think.System() or Act.System()
  ├── r.Parsers["being"]      → Being.ParseInner() (think) or Being.Parse() (act)
  ├── r.Parsers["think-operators"] → Think.Parse() — available ops
  ├── r.Parsers["think-time"]     → timePressure() — budget remaining
  ├── r.Parsers["think-exchange"] → prior passes in this think session
  ├── r.Parsers["thought-history"] → recent thoughts across exchanges
  ├── r.Parsers["thread"]     → Thread.Parse() — id, members, status
  ├── r.Parsers["exchange"]   → Conversation.Parse() — full exchange log
  ├── r.Parsers["conversation"] → Conversation.ParseRecent(10)
  ├── r.Parsers["ref-context"]  → carried context from another exchange
  ├── r.Parsers["inner"]      → (act only) the surfaced thought
  └── r.Impulse               → appended as "message: ..."
```

---

## Target Architecture

The device layer separated from beings. Devices are shared infrastructure. Beings determine which slice lights up.

### Key Shifts

- **Beings don't own devices.** A being has references to devices, not ownership. Many beings can share a device. One being can exist on many devices.
- **User and Self are the same shape.** Both are beings with inner realities and a device layer at the bottom. The difference is composition — Self has Think/Act, User doesn't — not type.
- **Devices are a registry.** A shared hashmap on the device world (MacOS). Beings reference into it. The being shapes the relation so the device layer knows what to resolve.
- **MacOS is a world, not a device.** The machine is the world. Beings live inside it. Everything running on the machine is inside MacOS.
- **Devices own components.** A device is hardware — the machine, the phone, the server. Components are the capabilities that run on a device: terminal, websocket, API, inference provider. A device has a list of components. Beings reference devices, and the device routes to the right component.
- **The user is a being whose invariant is a human.** Michael is a being like any other. His descent terminates at a device component (terminal, WS) that waits for human input instead of computing a response. Thread routes to him by name, same as skyra.
- **Providers are not devices.** OpenRouter is a component that runs on a device, not a device itself.

### Reality Composition Tree (Target)

```
Universe
├── NewThread
│   ├── Exchange
│   └── MacOS (the machine world)
│       ├── Devices (registry)
│       │   └── macbook
│       │       ├── Terminal (stdin/stdout)
│       │       ├── WS (websocket server)
│       │       └── OpenRouter (inference provider)
│       └── Beings
│           ├── Self "skyra"
│           │   ├── Being (identity, purpose, relationships)
│           │   ├── Think (inner layer)
│           │   │   └── Operators: recall, remember, skill
│           │   └── Act (outer layer)
│           │       └── Operators: plan
│           │   devices: [macbook]
│           │
│           ├── Self "louise"
│           │   devices: [macbook]
│           │
│           └── User "michael"
│               ├── Being (identity, purpose, relationships)
│               devices: [macbook, phone]
│
├── Economics
└── OnResolve → Universe.Realize(collecting)
```

A device is hardware. Components are what run on it.

```
macbook (device)
├── terminal    (component)
├── ws          (component)
└── openrouter  (component)

phone (device)
├── push        (component)
└── ws          (component)
```

### Relation Flow (Target)

Two levels of routing. Thread routes by being name. The device routes to the right component.

```
Thread: "who is this for?" → michael
  Michael.Realize: attaches device routing context to relation
    macbook: "what component?" → terminal
      Terminal: prints to screen, waits for input

Thread: "who is this for?" → skyra
  Skyra.Realize: Think → Act → attaches device routing context
    macbook: "what component?" → openrouter (inference), ws (output)
      OpenRouter: LLM call, returns response
      WS: sends to browser
```

### Multi-Device Beings

A being can exist on multiple devices simultaneously. Same being, same Think/Act, same memory. Different surfaces. Each device has its own components.

```
michael:
  macbook → terminal, ws
  phone   → push, ws

skyra:
  macbook → openrouter, terminal, ws
```

The being is the constant. The devices are where it shows up. The components are how it shows up. Identity doesn't change because you picked up your phone instead of your laptop.

### Genome (Target)

Devices, components, and beings declared separately.

```
# devices (hardware)
device ~name macbook ~type macos

# components (run on devices)
component ~name terminal ~type stdin ~device macbook
component ~name ws ~type websocket ~port 8080 ~device macbook
component ~name openrouter ~type llm ~model anthropic/claude-sonnet-4-5 ~device macbook

# beings (reference devices, not components)
grow ~name skyra ~type llm ~devices macbook
grow ~name michael ~type user ~devices macbook
grow ~name louise ~type llm ~devices macbook
```

---

## File Map

```
skyra-v.05/
├── genome.skyra              ← being/provider declarations
├── main.go                   ← bootstrap, wiring, loop entry
├── go.mod
├── architecture.md
├── world-physics.md
├── notes.md
├── README.md
├── notes/
│   ├── data-spec.md          ← frontend JSON contract
│   ├── economics-spec.md     ← task-based ledger design
│   ├── frontend-runtime-mapping.md
│   └── inference-spec.md     ← energy pool design
├── specs/
│   ├── claudes-future-features.md
│   ├── memory.md
│   ├── memory-implementation.md
│   ├── one-way-relationships.md
│   └── routing-rules.md
└── src/
    ├── debug/debug.go        ← per-being, per-layer log files
    ├── inference/inference.go ← OpenRouter HTTP call
    ├── keychain/keychain.go  ← macOS Keychain lookup
    └── reality/
        ├── reality.go        ← the interface (7 lines)
        ├── relation.go       ← Relation, Impress, ParseResponse, Extract
        ├── meaning.go        ← Extract, ExtractTag, StripTag
        ├── universe.go       ← Universe, all snapshot types, assembleState
        ├── universe_test.go
        ├── newthread.go      ← NewThread (system world), Thread, Grow
        ├── exchange.go       ← Exchange, Conversation, ref parsing
        ├── self.go           ← Self (being world for LLM beings)
        ├── being.go          ← Being (pathos object, Parse/ParseInner)
        ├── think.go          ← Think (inner layer, budget loop)
        ├── act.go            ← Act (outer layer, protocol enforcement)
        ├── user.go           ← User (being world for human beings)
        ├── llm.go            ← LLM, Provider (inference invariant)
        ├── macos.go          ← MacOS (terminal device → target: machine world)
        ├── economics.go      ← Economics (fields map, no enforcement)
        ├── operators.go      ← Operators registry (not currently wired)
        ├── recall.go         ← Recall operator (search memories on disk)
        ├── remember.go       ← Remember operator (write memory to disk)
        ├── skill.go          ← Skill operator (load skill file)
        └── plan.go           ← Plan operator (stub)
```

24 source files, ~2,600 lines.
