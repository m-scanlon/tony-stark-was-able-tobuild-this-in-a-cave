# User Adapter

Michael's execution surface. The boundary between a human and the runtime.

This adapter is different from the inference adapter in one fundamental way — it is not request-response. Michael can send at any time. The runtime can send at any time. Both directions are live simultaneously.

## What It Does

**Inbound** — reads from Michael's surfaces (terminal, HTTP, voice). Wraps his raw input into a protocol string. Writes it to the runtime via stdout.

**Outbound** — reads from the runtime via stdin. Renders the response back to Michael on the right surface.

Both directions run concurrently. Two goroutines. One pipe.

```
Michael's surfaces  →  adapter stdout  →  runtime
runtime             →  adapter stdin   →  Michael's surfaces
```

## Wire Format

Same contract as all adapters.

**Adapter → runtime stdout (inbound from Michael):**

```
skyra skyra <what michael said> | experience

---
```

The adapter wraps Michael's raw input into a valid protocol string. The runtime receives a fully formed signal — not raw text. Michael just talks. The adapter handles the wrapping.

**Runtime → adapter stdin (outbound to Michael):**

```
<what skyra said>

---
```

Plain text response terminated by `---`. The adapter reads it and renders it to whichever surface Michael is currently on.

## Multiple Surfaces

Michael is one being. It does not matter which surface his input arrives from — terminal, HTTP, voice. The adapter multiplexes them all into one stdout stream going to the runtime. The runtime sees one being named michael regardless of the device.

Each surface is registered with the adapter at startup:

```go
adapter := NewAdapter()
adapter.Register(CLISource(os.Stdin))
adapter.Register(HTTPSource(":8080"))
adapter.Register(VoiceSource("whisper"))
adapter.Run()
```

Adding a new surface is registering a new source. No new process. No runtime change.

## The Two Goroutines

**Inbound goroutine** — listens across all registered sources. When Michael sends something on any surface, wraps it into the wire format and writes to stdout. Runs continuously.

**Outbound goroutine** — reads from stdin. When the runtime sends a response, renders it to the right surface. Runs continuously.

They run independently. Michael can send while the runtime is mid-response. The runtime can send while Michael is mid-sentence. The adapter handles both without blocking either direction.

## Surface Routing On Outbound

When the runtime sends a response, the adapter needs to know where to render it. The simplest policy: render to whichever surface the most recent inbound signal came from. Michael asked from the terminal — respond to the terminal. Michael asked via HTTP — respond via HTTP.

The adapter tracks the last active surface internally. No routing information needs to travel through the runtime.

## Genome Declaration

```
skyra being ~name michael ~surface process ~command "user-adapter --http :8080" ~identity <identity> ~purpose <purpose> | reason
```

Michael's adapter starts with whatever surfaces are declared in the command. Additional surfaces can be added through grow — a new genome directive for michael with an updated command triggers a hot reload of the adapter.

## Shutdown

Same as inference adapter. Runtime closes stdin → adapter drains → exits cleanly. SIGTERM as fallback.

## What This Enables

Once this adapter exists, Michael does not interact with the runtime through main.go's stdin loop. He interacts through his own adapter. The hardcoded scanner in main.go goes away. Michael is a being like any other — declared in the genome, reachable through the router, with a process surface the router spawned.

## Open Questions

- How does the outbound goroutine match a runtime response to the right HTTP request? HTTP is request-response — the adapter needs to hold the response writer open until the runtime sends back. Needs a correlation mechanism internal to the adapter.
- What happens when Michael sends a second message before the runtime has responded to the first? Queue them, or send immediately and let the runtime handle concurrency?
- Does voice input get transcribed by the adapter before wrapping into the protocol string, or does the adapter send raw audio and something else transcribes?
- When a new surface comes online mid-session (Michael opens a browser tab while already in a terminal session), does the adapter treat it as the same being or a new signal source?
