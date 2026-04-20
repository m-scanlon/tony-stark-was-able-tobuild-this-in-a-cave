# Execution Surface

Every being has an execution surface. The surface is how a being acts in the world. The being is the constant. The surface is the implementation detail underneath it.

Every being declares its execution surface and callable language like everything else. No exceptions. Uniform all the way down.

## Two Execution Surfaces

**Internal** — the kernel handles it directly. Deterministic. No external call. This is the floor — the baseline the protocol needs to function. grow, start-exchange, close-exchange. Hardcoded by necessity. The protocol cannot bootstrap without it.

**Process** — every other being. The kernel invokes a process over stdin/stdout. One hardcoded primitive: invoke a process, send it data, get something back. Everything above that primitive is declared on the being and hot-reloadable through grow.

Process is uniform. Inference beings, external beings, Michael — all processes. The difference is what the process does:

- **Inference process** — receives the present, calls a model, returns protocol strings
- **External process** — receives a signal, acts in the world, optionally returns a response
- **Adapter process** — multiplexes multiple input sources into stdin, fans output from stdout back out to surfaces

The distinction belongs to the process. Not the kernel. The kernel just invokes.

## The Adapter Pattern

Every being's process is its adapter. The adapter handles both directions — stdin for input into the runtime, stdout for output from the runtime. One process. Unified.

The adapter template handles:
- The input/output channel
- Reading from sources and writing to stdin
- Reading from stdout and routing to the right output surface

Each source is plug-in. Adding a new surface to a being is registering a new source with its adapter. No new process from scratch.

```
adapter := NewAdapter()
adapter.Register(HTTPSource(":8080"))
adapter.Register(VoiceSource("whisper"))
adapter.Register(StdinSource())
adapter.Run()
```

An inference being's adapter is simpler — one source (the present from the kernel), one output (protocol strings back to the kernel).

## Callable Language Is The Interface

A being declares its execution surface and what it accepts through its callable language. The caller reads the callable language and decides how to invoke it. The being owns what it offers. The caller decides how to use it.

A being with one inference provider declares that provider in its callable language. A being with multiple providers declares all of them. The caller chooses. No exceptions. No hidden configuration.

## What This Enables

PFC can spawn four strategy replicas each on a different provider if strategy declares multiple providers in its callable language. That is not a special feature. That is beings in relationship doing what they already do.

Differentiation falls out naturally — spawn a new being with a different surface declared in the genome. No new machinery.

## Hot Reload

Execution surface is hot-reloadable through the protocol. grow already handles this — if a being exists, grow updates it rather than creating a new one. Firing a new genome directive for an existing being with a different surface is already hot reload. No new mechanism needed.

## The Cognitive Flag Goes Away

`Cognitive bool` on `Being` is a proxy for execution surface. All it currently does is decide whether a signal routes to inference or executes directly. That is an execution surface question, not an ontological one.

If every being declares its execution surface, the cognitive flag is redundant. The surface IS the distinction. A being on an inference process reasons. A being on an internal handler executes directly. The binary flag disappears. And the constraint disappears with it.

## Present Is Shaped By The Being

`DerivePresent` does not split on cognitive vs non-cognitive. It renders what exists and skips what doesn't.

If identity is blank — don't render the identity block. If there are no open exchanges — don't render the exchange section. If there are no peers — don't render the network. The present becomes naturally minimal for simple beings and naturally full for complex ones.

Same function. The being's own state determines what comes out. No separate code path. No binary split.

## Decided

- **Genome syntax** — `~surface process ~command "<path>"`. See adapter-inference.md.
- **Spawn model** — long-running process. The router spawns the adapter at being registration and keeps it alive. Not spawn-per-signal.
- **Wire format** — present as plain text terminated by `---`. Adapter writes protocol strings back one per line terminated by `---`. See adapter-inference.md.
- **Process failure** — router monitors adapter health, restarts on crash. Restart policy is per-being. See process-router.md.
- **Michael's surfaces** — all merge at the boundary adapter. Michael is one being. The adapter multiplexes all surfaces into one stdout stream. See adapter-user.md.

## Open Questions

- HTTP adapters receive two kinds of responses — transport acknowledgments (200 OK) and content responses (actual signal). The adapter needs to distinguish between them and only pass content responses through to the runtime. Transport responses are plumbing, not signal. How the adapter makes that distinction is unresolved.
