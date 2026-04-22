# Next Steps

## Where We Are

The base runtime is running. One loop. One world. One process. Skyra responds to Michael.

The Logos interface is proven. Everything composes. 674 lines.

## What's Next: The Child Process

The runtime is not extendable yet. All beings live in the same process, in the same map, sharing the same memory. That works for a demo. It doesn't work for a real system.

A being needs to be able to run as its own process.

That means:
- A being can be spawned as a child process by the world
- The child process communicates via stdin/stdout using the same protocol
- The world holds a `Logos` adapter that wraps the child process — it looks like any other Logos from the inside
- `Relate` on the adapter sends the relation to the child's stdin and reads the response from stdout

The protocol is already the wire format. `skyra <target> <expression> | <reason>` goes in. A response comes back. The adapter is just an IO wrapper around that.

This is the execution surface. Once it exists, a being can be:
- A Go function in the same process
- A child process on the same machine
- Anything that speaks the protocol over a pipe

The world doesn't know the difference. The LogosMap doesn't know the difference. That's the point.

## Why This Before Everything Else

Memory compression, relationship emergence, retained artifacts — all of that requires beings that can run independently and accumulate state outside the parent process. Without the child process, every new capability has to live in the same binary. The runtime can't grow.

The child process is the extensibility primitive. Everything else builds on top of it.

## Shape Of The Work

1. `src/adapter/process.go` — `ProcessAdapter` implements `Logos`. Wraps a `*exec.Cmd`. `Relate` writes the relation to stdin, reads the response from stdout.
2. `world.go` — add a `Spawn` operator that starts a child process and registers the adapter in the LogosMap.
3. `genome.skyra` — a being can declare itself as a process path instead of an inline identity.
4. The child process binary implements `main.go` with the same stdin loop — it IS a Logos runtime, just smaller.

The child process is a Skyra runtime running inside a Skyra runtime.
