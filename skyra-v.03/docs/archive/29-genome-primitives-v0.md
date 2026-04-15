# Genome Primitives v0

## Status

In progress. Not yet locked canon.

## Purpose

This document defines the primitive beings seeded at boot before the genome
runs.

These are not genome beings in the ordinary sense. They are the beings that
make the genome possible. They exist before any genome expression fires.

---

## The Grow Being

### What It Is

The grow being is a non-cognitive primitive.

It is the only path for creating beings in the runtime from outside direct
bootstrap registration.

When it receives an expression, it parses the full being creation protocol,
constructs the being, and registers it. One protocol string in. One live being
registered in the world out.

In the code this is `world.Grow(expression string)` — which calls
`being.CreateBeing(expression)` then `world.Register(b)`.

### Callable Language

```
skyra being ~name <name> ~identity <identity> ~purpose <purpose> ~language <expression> ~cognitive <true|false> | <source>: <reason>
```

- `~name` — the new being's routing name, unique in the runtime
- `~identity` — the new being's identity
- `~purpose` — the new being's purpose
- `~language` — the new being's base language
- `~cognitive` — true if the being reasons through inference; false if it is a transducer
- `<source>` — the being requesting creation
- `<reason>` — why the being is being created

### What It Delegates

The grow being does not parse identity, purpose, or language itself.

It delegates:

- `~name` extraction → `extract.Meaning`
- `~identity` + `~purpose` → `nature.CreateNature` (which delegates further to
  `identity.CreateIdentity` and `purpose.CreatePurpose`)
- `~language` → `language.CreateLanguage`
- `~cognitive` → local flag parse

### Key Rules

- The grow being is seeded at boot before any genome being.
- It cannot create itself. The kernel seeds it directly at init — this is the
  one bootstrap exception.
- All subsequent being creation goes through it.
- If any field is missing or invalid, the expression is dropped. No partial
  beings are registered.
- The grow being does not relate the new being to anything. Relating is a
  separate kernel operation.

---

## The Nature Being

### What It Is

The nature being is a non-cognitive primitive.

It holds exactly two fields: identity and purpose.

It is the minimum shape of every being in the system.

It does not reason. But it has explicit callable language so that any being
relating to it knows exactly how to speak to it.

### Callable Language

```
skyra nature ~identity <identity> ~purpose <purpose> | <source>: <reason>
```

- `~identity` — the identity field
- `~purpose` — the purpose field
- `<source>` — the being calling it
- `<reason>` — why it is being called

### Key Rules

- Nature is locked at creation. It is not mutable.
- The callable language on the relationship with nature is seeded at birth —
  not earned through use.
- Non-cognitive primitive beings are born with their callable language.
  Cognitive beings earn theirs through relating.
- The nature being's parser reads `~identity` and `~purpose` from the
  expression. That is its entire slice of the protocol.

---

## Boot Order

1. `world.New()` — the being registry is created. Empty.
2. The grow being is constructed directly via `world.Register` — the bootstrap
   exception. It cannot use itself to create itself.
3. `metaxu.New(w)` — the router is wired to the world.
4. Genome beings are seeded. Each genome being is created by sending a creation
   expression to the grow being, or by calling `world.Grow` directly at init.
5. `world.Relate` is called to seed the peer relationships between beings that
   need to know each other at boot.

At this point the runtime is live. Signal can flow.

### What Does Not Exist at Boot

Genome beings created later from templates — peripheral input beings, motor
beings, companion beings — are not present at boot. They are created when a
being's purpose requires them.

---

## What the Kernel Owns and What It Delegates

### The Kernel Owns

- **The being registry** — `world.World.beings`: a name-keyed map of all live
  beings in the runtime. The kernel is the only writer.
- **Routing** — `metaxu.AcceptSignal`: resolves origin and target by name,
  writes to both sides of the exchange, derives present for the receiver.
- **Growth** — `world.Grow`: the external creation path.
- **Relating** — `world.Relate`: seeds both sides of a peer relationship.
  Chooses `ExchangeStack` for cognitive beings, `ExternalDispatch` for
  non-cognitive beings.
- **Channel type selection** — `world.seedPeer`: cognitive flag on the being
  determines which channel type is attached.

### The Kernel Delegates

- **Expression parsing** — to the primitive packages: `extract`, `nature`,
  `identity`, `purpose`, `language`.
- **Present derivation** — to the channel. `ExchangeStack.DerivePresent`
  assembles the full present string from identity, purpose, open exchange, and
  peer list with callable languages. `ExternalDispatch.DerivePresent` returns
  the last expression only.
- **Inference** — to `inference.Runner`. The kernel derives the present and
  hands it to the runner. The runner calls the model and returns a `Signal`.
  The kernel does not know what inference does.
- **Callable language** — to the relationship. Lives on `ExchangeStack` and
  `ExternalDispatch`, not on the being. The kernel exposes it through the
  `RelationshipChannel` interface via `CallableLanguage() string`.

### The Boundary

The kernel is the only place strings become objects.

Everything arriving at the kernel is a raw string. The kernel parses, routes,
and writes. It does not reason. Reasoning belongs to beings.
