# Non-Cognitive Execution Surface v0

## Status

In progress. Not yet locked canon.

## Purpose

This document defines how non-cognitive beings execute.

Cognitive beings execute through inference — the kernel derives a present and
dispatches it to an inference runner. Non-cognitive beings do not reason. They
receive a signal and execute. This document defines where and how that
execution happens.

---

## The Problem

`AcceptSignal` routes to a target being by name. For cognitive beings, routing
ends at `DerivePresent` — the present string goes to an inference runner and a
protocol string comes back.

For non-cognitive beings, routing ends at `ExternalDispatch.Send()`. But
`Send()` just stores `lastExpression`. Nothing executes.

Non-cognitive beings need a runner — something that takes the delivered
expression and does something with it.

That runner is not a single function. Different non-cognitive beings execute in
different places. The router needs to answer: given this being, where does this
signal execute?

---

## The Shape

Each non-cognitive being has a registered execution surface binding:

```
being name → surface type → { address, handle }
```

- **surface type** — what kind of surface this is: `internal`, `cli`, or `http`
- **address** — where the surface lives: a function pointer for internal, a
  binary path for CLI, a URL for HTTP
- **handle** — the live instance: nothing for internal, a process handle for
  CLI, a connection handle for HTTP

The router resolves the binding at dispatch time and calls the surface.

---

## Surface Types

### Internal

Execution stays inside the kernel process.

The address is a registered handler function. The handle is not needed.

Internal surfaces are hardwired at bootstrap. They are not registered
dynamically.

**Examples:**

- **Grow** — `world.Grow(expression)`. Creates and registers a new being from a
  protocol expression. This is the primary internal surface.
- **Memory** — future. Queries a memory store directly. Registered at bootstrap
  alongside grow.

### CLI

Execution leaves the process via a spawned subprocess.

The address is the binary path or command. The handle is the process handle
returned when the subprocess is started.

CLI surfaces are dynamic. They are registered when a non-cognitive being backed
by a CLI binary is grown. The being expression and the surface binding arrive
together.

Multiple CLI instances can run simultaneously. Each has its own handle. The
router uses the handle to route the signal to the correct instance, not just
the correct binary.

**Examples:**

- An API binary being. The binary executes over the CLI, returns a result, and
  the result flows back as a protocol string.
- Any external tool wrapped as a non-cognitive being.

### HTTP

Execution leaves the process via HTTP on the local network.

The address is the endpoint URL. The handle is the connection handle.

HTTP surfaces follow the same pattern as CLI — registered dynamically when the
being is grown, multiple instances possible, each with its own handle.

**Examples:**

- A remote non-cognitive being running on another machine.
- A plugin service that speaks the protocol over HTTP.

---

## Memory as a Distributed Internal Surface

Memory is an internal surface but it is not a single instance.

Memory is distributed. Multiple memory stores exist simultaneously. Each store
is a separate instance with its own handle.

The surface binding for memory includes the store address and the handle for
that store. When a cognitive being grows a companion memory being, the surface
binding for that memory instance is registered at the same time.

The router resolves: this signal targets this memory being → find the binding
for this instance → dispatch to that store.

---

## Registration

### At Bootstrap

Internal surfaces are registered at bootstrap before the first signal.

Grow is registered first — it must exist before any other being can be created.
Memory surfaces are registered alongside it.

### At Runtime

When a cognitive being grows a non-cognitive being backed by an external
surface, the surface binding is registered at the same time as the being.

`world.Grow` for a non-cognitive being accepts the being expression and the
surface binding. The being is registered in the world. The surface binding is
registered in the execution surface registry. Both happen atomically — a being
without a surface binding and a surface binding without a being are both errors.

---

## What the Kernel Owns

The execution surface registry is kernel-owned.

Beings do not know where each other execute. The kernel resolves the surface at
dispatch time.

The kernel retains:

- the surface type registry
- the address and handle for each non-cognitive being's surface
- responsibility for spawning CLI processes and holding the returned handles
- responsibility for binding HTTP connections and holding the returned handles

---

## The Non-Cognitive Runner

The non-cognitive runner is the dispatch layer that sits after
`ExternalDispatch.Send()`.

It mirrors the cognitive runner's position in the signal flow:

```
cognitive:     ExchangeStack.Send() → DerivePresent → inference.Runner → signal back
non-cognitive: ExternalDispatch.Send() → DerivePresent → non-cognitive runner → result back
```

The non-cognitive runner receives the expression from `DerivePresent`, looks up
the surface binding for the target being, and dispatches to the correct surface.

For internal surfaces it calls the handler directly.

For CLI surfaces it writes to the process stdin and reads stdout back.

For HTTP surfaces it posts the protocol string and reads the response.

The result in all cases is a protocol string that flows back into
`AcceptSignal`.

---

## Open Questions

- Whether the non-cognitive runner lives in the `metaxu` package or gets its
  own package
- Exact format of the surface binding struct
- Error handling when a CLI process dies mid-execution — does the kernel
  respawn it or drop the signal
- Whether HTTP non-cognitive surfaces need the same heartbeat/health check
  logic as cognitive HTTP surfaces (doc 25)
- How memory store handles are acquired — whether the kernel spawns the store
  or connects to an already-running instance
