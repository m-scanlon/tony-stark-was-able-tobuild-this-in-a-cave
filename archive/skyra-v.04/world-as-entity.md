# World as Entity: Three Open Tensions

Three design frictions in v.04. Held together, they point at the same unresolved shape: **the world isn't a first-class entity yet.** Once it is, all three dissolve at once.

## 1. The world-operator cycle

At bootstrap, the world constructs its operators and shares its EntityMap with them:

```go
l := make(map[string]entity.Entity)
newWorld := World{EntityMap: l}
l["grow"] = &Grow{EntityMap: l}
l["continue-thread"] = &thread.ContinueThread{EntityMap: l}
l["parent"] = w
```

- World holds the map.
- Map holds the operators.
- Each operator holds the same map.

The reference graph is a loop. Go handles it, but the shape is tangled. You can't serialize a world without breaking cycles, and you can't cleanly move an operator between worlds.

## 2. Reference flow

Operators get their map reference at construction time and carry it forever. This forces:

- One operator instance per world.
- Cross-world operations aren't possible — operators are locked to their home map.
- The "context" (which world is active) is implicit in which operator instance you called, rather than something the relation carries.

If operators are reusable across worlds, the reference needs to flow *with the relation*, not be stored on the operator.

## 3. The world has no surface

A being has a `Medium` — the function that reaches outside the runtime. A world has nothing analogous. No way to answer: how does another world address this world? How does this world expose itself for external interaction?

## These are one question

Each tension is a symptom of the world not being a first-class entity. If the world becomes an entity with the same machinery as beings — identity, medium, DerivePresent, Relate — the three dissolve:

- A world with a **medium** has a surface.
- A world that exposes itself via its medium doesn't need operators to hold its map — **the relation carries the world**.
- The cycle resolves because operators are stateless. They ask the relation where they are.

## Operators go away

Operators are not entities in the map. The world *is* the routing. When a being emits `philosopher you're wrong about that`, the world's Relate handles it — looks up the target, builds the present, calls the medium, dispatches the response. That's not an operator doing work. That's the world deriving its present.

Grow, threading, exchange tracking, close-exchange — all internal to how the world works. The being never sees any of it. The map holds beings. Only beings. The world is the thing that holds the map.

## The fix

### Relation carries the world

```go
type Relation struct {
    ID       string
    Origin   string
    ThreadID string
    Impulse  string
    World    *World
}
```

Operators become stateless — they look up targets in `r.World.EntityMap`, not their own stored field. The cycle is gone. The world owns the map.

### World.Medium

- **Being.Medium** — how a being reaches outside the runtime. Outbound. (`inference`, `cli`, `shell`, `exec`.)
- **World.Medium** — how a world is reached from outside itself. Inbound. (`stdio`, `tcp`, `pipe`.)

Same primitive shape, opposite direction. In Linux: a process makes syscalls (outbound), and is reachable via sockets/pipes/signals (inbound).

### Cross-world communication

Multiple flat worlds. Each world is a peer. Communication between worlds goes through each world's medium — IPC, not hierarchy. A being lives in one world. If two worlds need to talk, there's a bridge.

## Status

Spec'd here. Deferred until first real need.
