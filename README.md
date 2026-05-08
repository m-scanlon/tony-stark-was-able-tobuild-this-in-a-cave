# Skyra

An experimental runtime for autonomous beings.

## What This Is

Skyra is a recursive descent runtime where every component implements the same
interface: `Reality`. Beings, threads, exchanges, devices, the universe itself —
all Realities. A single mutable message (Relation) enters at the top and
descends through self-similar layers, each adding context or routing, until it
hits a being and comes back up.

The human is not an input source. The human is a being in the system — with
identity, purpose, relationships, and a device. The runtime doesn't get better
without you putting in the time.

5,300 lines of Go.

## What Makes Skyra Different

- **One interface for everything.** `Reality` — three methods: `ID()`, `Create()`, `Realize()`. Beings, memory, devices, the universe itself. Same shape all the way down.
- **Beings think privately.** Two layers: Think (inner, no one sees) and Act (outer, protocol-enforced). A being can recall, remember, plan, and run operators before it ever speaks.
- **The human is a being, not an input source.** Same interface, same physics. Identity, purpose, relationships, device. You teach by relating, not by configuring.
- **Memory belongs to the being that lived it.** Graph-backed, entity-anchored, relationship-scoped. No shared memory bus. No central store. What you remember is yours.
- **Context is managed, not dumped.** An LLM-curated layer between memory and thought. It warms relevant context per relationship, cleans content before storage, and keeps the graph from filling with noise.
- **Multi-party routing is recursive, not orchestrated.** No router, no dispatcher. Messages descend through the tree. Beings talk to beings. The thread tracks the graph.
- **Relationships have stakes.** Lose all your relationships, you end. XP accrues through exchange. This is the physics, not a game mechanic.
- **The genome is the creator's control surface.** A declarative file that defines beings, devices, components, and wiring. One file boots the entire universe.
- **Beings can grow other beings.** Mid-flight. No restart. The creator says `grow`, the thread creates a new being and wires it in.
- **Per-being operators.** Builder gets bash. Philosopher doesn't. Operators are capabilities, not global features.

## Architecture

```
Universe
├── NewThread
│   ├── Exchange
│   ├── Levels
│   ├── skyra (Self)
│   │   ├── Being
│   │   ├── Memory → Graph, Entities, Vectors
│   │   ├── Context → Heat, Store, Retrieve
│   │   ├── Think → Recall, Remember, Plan, Skill, Browse, Search
│   │   ├── Act
│   │   └── Desk
│   ├── builder (Self)
│   │   ├── Think → Bash, Recall, Remember, ...
│   │   └── ...
│   ├── philosopher (Self)
│   ├── louise (Self)
│   ├── michael (User)
│   │   └── MacOS → Terminal, WebSocket, LLM
│   └── claude (Agent)
└── Economics
```

Every node is a Reality. The tree is self-similar all the way down.

## Status

v.05 is the live version. Alpha targets June 1, 2026.

## Repo Map

- `skyra-v.05/` — live codebase
- `skyra-v.05/src/reality/` — the runtime
- `skyra-v.05/genome.skyra` — the universe definition
- `skyra-v.05/specs/` — design specs
- `archive/` — older generations (v.03, v.04). History, not canon.

## License

Functional Source License 1.1 with Apache 2.0 future license (`FSL-1.1-ALv2`).
Commercial competing use requires a separate license.

## The Question

Does giving intelligence somewhere to live, someone to learn from, and a history
that compounds produce something no model upgrade can?

## In 6 Words

Its a hashmap that calls itself
