# Runtime Primitives And Artifacts (Prelim)

## Core Framing

The system needs a clear distinction between:

- retained artifacts
- runtime primitives
- runtime stimulus
- runtime artifacts

These are not the same layer.

## The Split

The current working split is:

- retained artifacts are long-lived outputs of learning
- runtime stimulus is the typed request/response traffic used during execution
- runtime artifacts are transient outputs produced inside an episode while fulfilling that traffic

This means the system should not treat every runtime result as retained memory.

Most runtime execution should remain episode-local unless later learning decides it should be retained.

## Runtime Stimulus

Runtime execution should now be understood through emitted stimulus rather than older command/writeback language.

Conceptually:

```text
skyra <primitive> <actor> <surface> <stimulus_protocol>
```

The actor emits concrete stimulus payload toward a named public surface that conforms to a published contract.

The kernel routes and validates that payload against the appropriate request/response contract.

## Kernel Role

The kernel remains the authority over:

- routing
- contract validation
- execution-surface dispatch
- returned stimulus routing

At a high level, the flow is:

1. the actor sees the current frame
2. inference selects the next allowed primitive and payload
3. the actor emits Skyra protocol carrying that concrete stimulus payload
4. the kernel validates the target surface and payload against the published contract
5. the target surface handles the request
6. a response envelope or other returned stimulus is routed back
7. the actor writes the result into episode-local state

This means:

- the actor chooses the next operation
- the kernel controls routing and validation
- runtime execution unfolds as a sequence of bounded stimulus exchanges

## Runtime Artifacts

A runtime artifact is a transient output produced during runtime execution inside an episode.

Runtime artifacts are:

- in-episode
- temporary
- usable by later runtime steps
- available to later learning
- not retained by default

## Why This Split Matters

Without this split, the system collapses:

- request handling
- transient runtime work
- retained experience

into one layer.

The intended rule is:

- runtime stimulus carries the public request/response interface
- runtime artifacts are the transient internal outputs produced while fulfilling that interface
- learning later decides whether anything should become retained experience

## Learn

`learn` remains the post-episode write path.

Its role is:

- take a just-closed episode as input
- run the learning/consolidation path
- produce retained artifact and structure updates

It should not be confused with ordinary in-episode response handling.

## Current Design Posture

The strongest current claims are:

- runtime stimulus is distinct from retained artifacts
- runtime artifacts are transient outputs of runtime execution
- the kernel is the authority over routing and validation
- runtime artifacts should remain available to later learning
- `learn` remains the post-episode write path rather than ordinary request handling

## Short Framing

The runtime should distinguish:

- public request/response stimulus
- transient runtime artifacts
- retained artifacts selected later by learning
