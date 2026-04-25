# Process Router

The router reads the execution surface a being declares and routes to it. That is the entire routing decision.

## Two Surfaces, Two Paths

**Internal** — the kernel handles it directly. No process. No spawn. The handler is baked into the runtime. `grow`, `start-exchange`, `close-exchange`. These exist before any being does.

**Process** — the kernel routes to an adapter. The adapter is a long-running process the router spawns at being registration and keeps alive. The router writes to its stdin and reads from its stdout.

Every being declares which surface it runs on. The router reads that declaration and routes accordingly. No flag on the being struct. No hardcoded switch. The surface declaration is the routing instruction.

## The Adapter Is The Stable Surface

The adapter is what the router holds a relationship with — not the API, not the model, not the HTTP connection. Those live inside the adapter. The adapter owns everything downstream.

```
runtime → adapter (stdin/stdout) → API / model / external service
```

The router's pipe is to the adapter. The adapter manages its own downstream. If the model API goes down, the adapter handles reconnection. If an HTTP session drops, the adapter handles it. The runtime never sees that instability. The only failure the router cares about is the adapter itself going down.

This is why broken pipe is an adapter problem, not a router problem. The router's stdin/stdout relationship with the adapter is stable by design. If that pipe breaks, the adapter crashed — and the router restarts it.

## What The Router Tracks

One entry per being with a process surface:

```
BeingName    string
Command      string         — what to run to start the adapter
State        running | stopped | starting
Stdin        io.WriteCloser — runtime writes signals here
Stdout       io.ReadCloser  — runtime reads responses here
PID          int
```

Internal surface beings have no entry. The kernel handles them directly.

## Lifecycle

**Registration** — when a being with a process surface is registered (via grow), the router spawns its adapter and holds the handles. The being is not routable until the adapter is running.

**Dispatch** — when a signal arrives for a being, the router checks the entry. If running, write to stdin. If stopped, apply the restart policy before routing.

**Crash** — the router monitors each adapter process. When one exits unexpectedly, the entry moves to stopped. The router restarts it. Signals that arrive during restart apply the same policy as stopped.

**Hot reload** — grow can update a being's execution surface. The router tears down the old adapter and spawns the new one. The being is unreachable during the transition.

## Restart Policy

What happens when a signal arrives and the adapter is stopped is an open question. Three options:

- **Drop** — signal is lost, error returned to sender
- **Queue** — signal held until adapter restarts, then delivered
- **Restart and retry** — router restarts the adapter immediately, retries the signal once

The policy probably belongs on the being's surface declaration, not hardcoded in the router. Different beings may have different tolerance for signal loss.

## Relationship To grow

grow is the control surface for the router the same way it is for beings. A new genome directive for a being with a process surface triggers adapter registration. An updated directive for an existing being triggers hot reload. The router does not have a separate configuration mechanism — everything flows through grow.

## Decided

- **Wire format** — present as plain text terminated by `---` on its own line. Adapter writes protocol strings back one per line, terminated by `---`. See adapter-inference.md.
- **Shutdown** — router closes adapter stdin. Adapter drains in-flight work and exits cleanly. SIGTERM as fallback. See adapter-inference.md.

## Open Questions

- Does the router handle adapter stdout in a dedicated goroutine per adapter or a shared reader?
- What is the restart backoff strategy — immediate, exponential, or capped?
- Does the queue (if adopted) have a depth limit, and what happens when it fills?
