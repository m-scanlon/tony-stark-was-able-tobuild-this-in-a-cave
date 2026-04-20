# Inference Adapter

The first execution surface. A being declares this surface and the runtime can derive its present and get protocol strings back. The cognitive loop becomes real.

## What It Does

Receives the present from the runtime on stdin. Calls a model API. Writes protocol strings back to stdout. That is the entire job.

```
runtime stdin  →  present  →  model API  →  protocol strings  →  runtime stdout
```

## Wire Format

This is the contract all adapters share. Defined here first.

**Runtime → adapter stdin:**

```
your name is: prefrontal
your identity is: ...
...

---
```

The present as plain text. Terminated by `---` on its own line. The adapter reads lines until it sees `---`. Everything before it is the present. The adapter does not act until it sees the terminator.

**Adapter → runtime stdout:**

```
skyra start-exchange ~with values ~about ... ~because ... ~say ... | reason
skyra close-exchange ~with values ~learned ... | reason

---
```

One protocol string per line. Terminated by `---` on its own line. The runtime reads until it sees `---`. Everything before it is the response. Multiple protocol strings are valid — the runtime processes them in order.

**Error response:**

```
error: <description>

---
```

If the adapter cannot produce a valid response, it writes `error:` followed by a description and terminates with `---`. The runtime handles the error and does not crash.

## Configuration

The adapter receives its configuration through command line arguments declared in the genome surface expression:

```
~surface process ~command "inference-adapter --model claude-sonnet-4-6 --endpoint https://openrouter.ai/api/v1"
```

The API key comes from the environment — the runtime inherits its environment to child processes on spawn. The adapter reads `OPENROUTER_API_KEY` or `OLLAMA_API_KEY` from the environment directly. No key is passed through the wire format.

Different beings can declare different models by declaring different commands. PFC on a reasoning model. A fast reactive being on a smaller model. The genome is the control surface for model selection.

## The Loop

```
start
loop:
    read lines from stdin until "---"
    if stdin closed → exit cleanly
    call model API with accumulated present as prompt
    if API error → write error response → continue
    parse response into protocol strings
    write each protocol string to stdout
    write "---" to stdout
```

The adapter runs this loop for the lifetime of the being. It does not exit after one response. The runtime writes the next present when the next signal arrives.

## Shutdown

When the runtime wants to stop the adapter it closes the adapter's stdin. The adapter sees EOF on the read side, drains any in-flight work, and exits with code 0. The router sees the process exit cleanly and marks it stopped.

SIGTERM is the fallback. The adapter catches SIGTERM, finishes the current response if mid-flight, writes `---` to stdout, and exits.

## What The Adapter Does Not Know

The adapter does not know about beings. It does not know about exchanges or threads or the protocol. It receives text and returns text. The runtime handles everything else.

This is intentional. The adapter is the translation layer only. Keeping it ignorant of the protocol means it stays simple enough for Skyra to eventually write one herself.

## Genome Declaration

```
skyra being ~name prefrontal ~surface process ~command "inference-adapter --model claude-sonnet-4-6 --endpoint https://openrouter.ai/api/v1" ~identity <identity> ~purpose <purpose> ~relationships <a,b,c> | reason
```

`~surface process` tells the router this being has a process surface. `~command` is what the router runs to spawn the adapter. The router reads these at registration, spawns the adapter, holds the handles.

## Open Questions

- Does the adapter stream tokens back as they arrive or wait for the full response before writing to stdout? Streaming is better for latency but complicates the wire format — the runtime needs to know when a protocol string is complete mid-stream.
- What is the maximum present size the adapter should accept before truncating? Large presents with deep exchange histories could exceed model context limits.
- Does the adapter do any validation that what it writes to stdout is valid protocol syntax, or does the runtime handle that entirely?
