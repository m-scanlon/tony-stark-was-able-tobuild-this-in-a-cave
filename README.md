# Skyra

A runtime for autonomous beings.

Not an agent framework. Not a chatbot wrapper. Not a prompt pipeline with memory
bolted on. A runtime in which minds live.

## What This Is

Skyra is a recursive descent runtime where every component implements the same
interface: `Reality`. Beings, threads, exchanges, devices, the universe itself —
all Realities. A single mutable message (Relation) enters at the top and
descends through self-similar layers, each adding context or routing, until it
hits a being and comes back up.

Beings have two layers: Think (private inner thought) and Act (outer speech and
routing). Think has operators — recall, remember, skill — and a budget. Act
enforces protocol and routes to peers. No one sees what a being thinks. Everyone
sees what it says.

The human is not an input source. The human is a being in the system — with
identity, purpose, relationships, and a device. Teaching, correction, memory,
and relationship all become part of the universe the beings inhabit. The runtime
doesn't get better without you putting in the time.

2,618 lines of Go.

## What Is Settled

- Reality is the interface: `ID()`, `Create()`, `Realize()`
- Being is the atomic unit: identity, purpose, relationships, device
- Self contains Think and Act — private thought and public speech
- The relation bus carries a single mutable message through recursive layers
- Exchange is append-only, context crossing requires explicit `<ref>` tags
- Memory is local to the being that lived it (filesystem-backed)
- NewThread owns the loop — multi-party routing is recursive, not orchestrated
- The genome is the creator's control surface
- Universe is the outermost Reality — its present is the full state of everything

## Architecture

```
Universe
├── NewThread
│   ├── Exchange
│   ├── skyra (Self)
│   │   ├── Being
│   │   ├── Think → Recall, Remember, Skill
│   │   └── Act → Plan
│   ├── louise (Self)
│   │   └── ...
│   └── michael (User)
│       ├── Being
│       └── MacOS
└── Economics
```

Every node is a Reality. Every node can contain other Realities. The tree is
self-similar all the way down.

## Status

v.05 is the live version. Alpha targets June 1, 2026.

Done: recursive descent engine, multi-party threads, Think/Act planes, memory,
context crossing, mid-flight grow, universe serialization, frontend contract.

Next: WebSocket device, Inference (energy per being), Economics (task economy).

## Repo Map

- `skyra-v.05/` — live codebase
- `skyra-v.05/notes/` — specs and design notes
- `skyra-v.05/specs/` — future features and roadmap
- `architecture-evolution-timeline.md` — the full arc from February to now
- `archive/` — older generations (v.03, v.04). History, not canon.

## The Question

Does giving intelligence somewhere to live, someone to learn from, and a history
that compounds produce something no model upgrade can?
