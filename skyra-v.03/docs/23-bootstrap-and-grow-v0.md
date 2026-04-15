# Bootstrap And The Grow Being v0

## The Bootstrap Problem

The kernel enforces a hard invariant: no being can interact with the system unless it already exists in the hashmap. This creates a chicken-and-egg problem — something has to be in the hashmap before anything can happen.

The solution is a single hardcoded bootstrap step. The grow being is constructed directly via `world.Register` — the one bootstrap exception. It cannot use itself to create itself.

From that point on, everything goes through the protocol.

---

## The Grow Being

`grow` is the only path to runtime instantiation post-bootstrap. No other being can call `NewBeing` directly. If it doesn't go through `grow`, it doesn't happen.

It is non-cognitive. It has its own syntax for being creation and relationship seeding. It is the bridge between protocol strings and the underlying Go runtime.

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

The grow being does not parse identity, purpose, or language itself. It delegates:

- `~name` extraction → `extract.Meaning`
- `~identity` + `~purpose` → `nature.CreateNature`
- `~language` → `language.CreateLanguage`
- `~cognitive` → local flag parse

### Key Rules

- Seeded at boot before any genome being
- Cannot create itself — the kernel seeds it directly
- All subsequent being creation goes through it
- If any field is missing or invalid, the expression is dropped. No partial beings are registered
- Does not relate the new being to anything — relating is a separate kernel operation

---

## Boot Order

1. `world.New()` — the being registry is created. Empty.
2. The grow being is constructed directly via `world.Register` — the bootstrap exception.
3. `metaxu.New(w)` — the router is wired to the world.
4. Genome beings are seeded. Each genome being is created by sending a creation expression through grow.
5. Relationships are seeded between beings that need to know each other at boot.

At this point the runtime is live. Signal can flow.

Genome beings created from templates — peripheral input beings, motor beings, companion beings — are not present at boot. They are created when a being's purpose requires them.

---

## What The Kernel Owns And What It Delegates

### Owns

- **Being registry** — `world.World.beings`: name-keyed map of all live beings. The kernel is the only writer.
- **Routing** — `metaxu.AcceptSignal`: resolves origin and target by name, writes to both sides of the exchange, derives present for the receiver.
- **Growth** — `world.Grow`: the external creation path.
- **Relating** — `world.seedRelationships`: seeds both sides of a peer relationship. Chooses `ExchangeStack` for cognitive beings, `ExternalDispatch` for non-cognitive beings.
- **Channel type selection** — cognitive flag on the being determines which channel type is attached.

### Delegates

- **Expression parsing** — to the primitive packages: `extract`, `nature`, `identity`, `purpose`, `language`.
- **Present derivation** — to the channel. `ExchangeStack.DerivePresent` assembles the full present string. `ExternalDispatch.DerivePresent` returns the last expression only.
- **Inference** — to `inference.Runner`. The kernel derives the present and hands it to the runner. The kernel does not know what inference does.
- **Callable language** — to the relationship. Lives on `ExchangeStack` and `ExternalDispatch`, not on the being.

### The Boundary

The kernel is the only place strings become objects. Everything arriving at the kernel is a raw string. The kernel parses, routes, and writes. It does not reason. Reasoning belongs to beings.

---

## The Genome

`genome.skyra` is not read by the kernel directly. It flows through `grow`. You write the genome, the bootstrap runs, and from there the system handles instantiation. That is the one act of creation. Everything after is the system's.
