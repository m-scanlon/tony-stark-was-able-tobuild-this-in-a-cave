# Medium Abstraction

## Requirements

- Beings resolve over mediums.
- Medium is the highest abstraction. It sits at the same level as Entity.
- Beings contain mediums in their config representing what mediums they can derive their present in. This list is not representative of the mediums a being *will* derive its present in — it is the list of mediums it *can* derive its present in. The Being registry and the Medium registry are independent and not coupled to each other. The world is where coupling happens.
- Invariants can derive their present over the same mediums beings can.
- Mediums are interfaces. Always one way. In, never out.
- Affordances are scoped to the being and actualized over mediums. A being carries potential. Which affordances are available depends on what medium the being is currently operating over.
- Mediums have entry points. A laptop is a medium — its entry points are the keyboard, mouse, HTTP if a port is open. Beings have entry points. Entry points live on compatible mediums.
- Mediums have light coupling to beings. Only when instantiated together on the world.
- The output medium is not important. What matters is who the message is addressed to. "Sent from iPhone" is metadata, not the message.
- Beings address beings, not mediums.
- Beings own mediums when they live on the world. Mediums are how other beings can get ahold of them.
- The present is always created based on the being type.

## Current Direction

- Mediums always have a world initialized on them. Beings man mediums. One being is not tied to a single medium.
- The present for a being is global. There is one present per being. The present is divided into decomposable sections. A being can derive their present over multiple mediums/affordances.
- A relation always carries a string. That string may be a pointer to some other piece of data, but it is a string.
- Medium is a present concern. It is resolved at present creation time. The being declares what mediums it can operate over. Which ones are active is determined when the world derives the present. The present owns the mediums for that moment.
- Medium is not a static field on the being. It is resolved fresh each time the world derives a present.
- Entities are addressable units. Mediums are a pass through. Interfaces resolve the in-between.
- World is an entity. Medium is the runtime. Grow is the interface that the runtime gives. Grow lives on the world as an interface.
- The genome's `~medium cli` is wrong twice. It puts medium on the being instead of the world, and it names an interface instead of a medium. The medium is the runtime. CLI is an interface on the runtime. The being doesn't own the medium — it reaches the world through an interface on the medium.
- Whether something is a medium or an entity is a matter of perspective. From the beings, the laptop is an entity. From the user, the laptop is a medium. If medium is just "entity from the outside," then there are only two primitives: entities and interfaces. The medium was never a separate thing — it was an entity you're looking through instead of looking at.
- There is a thing between the entity and the interface where the present gets rendered — shaped by the constraints of whatever it's passing through. The present exists on the being. The interface determines what can come in. The lens is where the present becomes concrete. Same present through a CLI lens is text. Through a frontend lens it has space, layout, regions. The being doesn't change. The lens does. (Working name: lens.)
- The frontend is not an app. It is a protocol for pushing present data to whatever lens is available. TV, laptop, phone, watch — each one is a lens with different constraints. The runtime pushes present data. The lens shapes it. Skyra doesn't need a mobile app and a desktop app and a web app. It needs one runtime that pushes present data, and lenses that know how to render it on their surface. The being's present is the same everywhere. The lens is the only thing that changes.
- An app is just a blank lens. It holds no state, no logic, no present. It is a transparent surface that the runtime pushes to. The lens receives and renders. The state of the lens is the last time the being that manned it had its present derived.
- The entire system is push. Relations push in, presents push out to lenses. The only pull is the runtime reaching into its own storage — reading files, loading the genome, retrieving retained artifacts. That is infrastructure, not protocol. The protocol is all push. The runtime pushes your present to your surfaces the moment it changes. You don't ask for it. It arrives.
- Three primitives: entities, interfaces, and lenses. Addressable units, intake surfaces, and rendering surfaces.
- React Native as the lens framework. The lens is a thin shell — a blank React Native app with a registry of primitive components. The runtime pushes a component tree description as present data. The lens receives it, resolves the components from its registry, renders natively on whatever surface it's running on. Phone, tablet, TV, web — each lens has its own component registry tuned to its surface. The runtime pushes the same present. Each lens maps it to its own native components. Same data, different glass.
- `DerivePresent` builds JSON objects instead of flat strings. The structured present gets pushed to connected lenses. The routing, threading, and exchange tracking stay the same. The only change is the output format and an output channel (WebSocket) that isn't stdout.
- Open question: should threads be decoupled from the world? If the runtime is infrastructure, threads are an opinion. Opinions belong in the implementation layer, not the runtime. But threads currently drive present derivation, routing, and `~ref` resolution. Decoupling is the right instinct but needs more thought.

---

## What's Wrong Today

### 1. Three agent mediums that are one medium

`claude.go`, `opencode.go`, and `codex.go` are identical except for the binary path. Each one:
- takes `r.Impulse`
- shells out to a CLI binary
- sanitizes the output
- returns `<origin> <response>`

This is one medium — `agent` — parameterized by path. The `exec` medium already does this for generic binaries. The agent mediums are `exec` with extra steps.

### 2. Three being types that are one being

`src/primitives/claude/`, `src/primitives/opencode/`, `src/primitives/codex/` — each is a copy of `Being` with `DerivePresent` returning `""`. They exist because agent mediums don't want the full present — they manage their own context and only need the raw task.

The being type shouldn't change because the medium changed. Whether a being gets a full present or a minimal one is a property of the medium, not the being.

### 3. The system prompt lives in the wrong place

`inference.go` hardcodes the protocol instructions — `<>` delimiter, `~ref` syntax, peer addressing rules. This is the world's protocol, not the medium's. The medium should receive the system prompt from the caller. The world knows the protocol. The medium knows how to call an API.

### 4. The medium signature hides a decision

```go
type Medium func(present string, r entity.Relation) (string, error)
```

Every medium receives `present` — the full derived context. But agent mediums ignore it. Shell ignores it. Only `inference` and `cli` use it. The signature doesn't distinguish between "I need context" and "I need a task."

---

## The Abstraction

### Medium becomes an interface

```go
type Medium interface {
    Call(input MediumInput) (string, error)
}
```

### MediumInput carries what the medium needs

```go
type MediumInput struct {
    Present  string           // full derived present (identity, thread, exchange, peers)
    Relation entity.Relation  // the raw relation
    System   string           // protocol instructions from the world
}
```

The medium decides what to use. Inference uses `Present` and `System`. Shell uses `Relation.Impulse`. An agent medium uses `Relation.Impulse`. CLI uses `Present`. The input is the same — the medium picks what matters.

### Agent is one medium, parameterized

```go
medium.Agent("/Users/mikepersonal/.local/bin/claude")
medium.Agent("/opt/homebrew/bin/opencode")
medium.Agent("/opt/homebrew/bin/codex")
```

One implementation. The binary path comes from the genome:

```
grow ~name claude ~medium agent:/Users/mikepersonal/.local/bin/claude ~relationships skyra,michael,builder
```

Same pattern as `exec:` today, but the agent medium handles the input/output protocol that all CLI agents share — pass the impulse as an argument, sanitize the response, format the return.

### The three being types collapse back into Being

With agent mediums ignoring `Present` on their own, there's no reason for separate being types. Claude, OpenCode, and Codex become regular beings with an `agent:` medium. `DerivePresent` always runs — the medium decides whether to read it.

The `claude/`, `opencode/`, and `codex/` packages are deleted. `grow.go` no longer needs special cases.

### The system prompt moves to the world

The world owns the protocol. When it fires a medium, it passes the system prompt as `MediumInput.System`. The inference medium puts it in the API call's system field. Other mediums ignore it.

```go
// world.go — when firing the medium
input := medium.MediumInput{
    Present:  present,
    Relation: r,
    System:   w.systemPrompt(),
}
response, err := m.Call(input)
```

The protocol prompt is defined once, on the world. If a child world runs a different protocol, it passes a different system prompt. The medium never changes.

---

## What Changes

| Before | After |
|---|---|
| `Medium` is `func(string, Relation) (string, error)` | `Medium` is an interface with `Call(MediumInput)` |
| 3 agent mediums (claude, opencode, codex) | 1 agent medium, parameterized by path |
| 3 being types (Claude, OpenCode, Codex) | all collapse into Being |
| system prompt hardcoded in inference.go | system prompt passed in from the world |
| `grow.go` switches on medium name for being type | `grow.go` always creates a Being |
| `exec` medium is separate from agent mediums | `agent` subsumes `exec` or they share the pattern |

---

## What Doesn't Change

- The genome format. `~medium agent:/path/to/binary` already works with the existing `exec:` pattern.
- The protocol. `<>` delimiter, `~ref`, peer addressing — all the same.
- The world's routing. `DerivePresent` → fire medium → parse response → route. Same loop.
- The `cli`, `shell`, and `inference` mediums still exist. They just implement the interface instead of being bare functions.

---

## What This Enables

- **New agent mediums without new code.** Any CLI tool that takes a message and returns text is an `agent:` medium. No new package, no new being type, no registration.
- **Per-world protocol.** Child worlds can define their own system prompt. A game world and a work world don't have to speak the same protocol.
- **Medium-level decisions.** A medium that wants to preprocess the present, add its own context, or transform the impulse can do so inside `Call`. The world doesn't need to know.
- **Testability.** A medium is an interface. Mock it.
