# Execution Surface v0

## Status

In progress. Not yet locked canon.

## The Problem

The kernel knows which being is the target and which is the source. What it does not yet answer cleanly is where — and how.

For cognitive beings: routing ends at `DerivePresent`. The present string goes to an inference runner and a protocol string comes back.

For non-cognitive beings: routing ends at `ExternalDispatch.Send()`. But `Send()` just stores `lastExpression`. Nothing executes. Non-cognitive beings need a runner — something that takes the delivered expression and does something with it.

In a distributed architecture — many machines, many inference providers, many CLI binaries — the kernel also needs to know not just which being to route to but where that being executes.

---

## The Shape

Each being has a registered execution surface binding:

```
being name → surface type → { address, handle }
```

- **surface type** — what kind of surface this is: `internal`, `cli`, or `http`
- **address** — where the surface lives: a function pointer for internal, a binary path for CLI, a URL for HTTP
- **handle** — the live instance: nothing for internal, a process handle for CLI, a connection handle for HTTP

The execution surface is not part of nature. Nature is identity and purpose — it does not change. The execution surface is operational. It answers where this being runs right now.

A cognitive being has one slot: inference — the endpoint where cognition happens.

A non-cognitive being has one slot: the surface that handles execution.

Either may be local or remote.

---

## Surface Types

### Internal

Execution stays inside the kernel process.

The address is a registered handler function. Internal surfaces are hardwired at bootstrap. They are not registered dynamically.

**Examples:**
- **Grow** — `world.Grow(expression)`. Creates and registers a new being from a protocol expression.
- **Memory** — future. Queries a memory store directly. Registered at bootstrap alongside grow.

### CLI

Execution leaves the process via a spawned subprocess.

The address is the binary path or command. The handle is the process handle returned when the subprocess starts. Multiple CLI instances can run simultaneously — each has its own handle.

CLI surfaces are registered dynamically when a non-cognitive being backed by a CLI binary is grown.

**Examples:**
- An API binary being. The binary executes over the CLI, returns a result, and the result flows back as a protocol string.
- Any external tool wrapped as a non-cognitive being.

### HTTP

Execution leaves the process via HTTP on the local network. Not public — local network only.

The address is the endpoint URL. The handle is the connection handle. Same dynamic registration pattern as CLI.

**Examples:**
- A remote being running on another machine.
- A plugin service that speaks the protocol over HTTP.

---

## The Non-Cognitive Runner

The non-cognitive runner is the dispatch layer that sits after `ExternalDispatch.Send()`. It mirrors the cognitive runner's position in the signal flow:

```
cognitive:     ExchangeStack.Send() → DerivePresent → inference.Runner → signal back
non-cognitive: ExternalDispatch.Send() → DerivePresent → non-cognitive runner → result back
```

The non-cognitive runner receives the expression from `DerivePresent`, looks up the surface binding for the target being, and dispatches to the correct surface. For internal surfaces it calls the handler directly. For CLI it writes to stdin and reads stdout back. For HTTP it posts the protocol string and reads the response.

The result in all cases is a protocol string that flows back into `AcceptSignal`.

---

## What The Kernel Owns

Distributing execution does not distribute kernel authority.

The execution surface registry is kernel-owned. Beings do not know where each other execute. The kernel resolves the surface at dispatch time.

The kernel retains:
- being registration
- relationship hashmap
- edge weight updates
- relationship emergence threshold logic
- trust values
- execution surface registry
- responsibility for spawning CLI processes and holding handles
- responsibility for binding HTTP connections and holding handles

None of that crosses the network. The only thing that crosses is the present going out and a protocol string coming back.

What would break the model is kernel authority crossing the network — a remote being registering other beings, writing to the relationship hashmap, updating edge weights. That is the line.

---

## Registration

### At Bootstrap

Internal surfaces are registered at bootstrap before the first signal. Grow is registered first — it must exist before any other being can be created. Memory surfaces are registered alongside it.

### At Runtime

When a being is grown with an external surface, the surface binding is registered at the same time. The being is registered in the world. The surface binding is registered in the execution surface registry. Both happen atomically — a being without a surface binding and a surface binding without a being are both errors.

---

## Memory As A Distributed Internal Surface

Memory is an internal surface but not a single instance. Memory is distributed — multiple memory stores exist simultaneously, each with its own handle. The surface binding for memory includes the store address and the handle for that store.

---

## Open Questions

- Whether a being's execution surface can change after registration or is fixed at birth
- Whether the kernel needs a heartbeat or health check against remote execution surfaces
- Whether a being whose execution surface goes offline becomes dormant or is removed
- The exact format of the execution surface field on the being record
- Whether local and remote beings need different trust origin values
- Whether the non-cognitive runner lives in the `metaxu` package or gets its own package
- Exact format of the surface binding struct
- Error handling when a CLI process dies mid-execution
