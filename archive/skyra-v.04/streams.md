# Streams: Mediums as Conduits

The current `Medium` is a function — call it with a request, receive a response, return. That works for one-shot, turn-based interactions. It doesn't work for:

- **Sensors.** Cameras, microphones, file watchers — emit data continuously, unprompted.
- **Actuators.** Speakers, displays, LEDs — accept data continuously, produce nothing meaningful back.
- **Streaming LLMs.** The model generates token-by-token; our current code waits for the whole response before parsing.
- **Multi-modal beings.** A human body streams video + audio + text + proprioception all at once. One medium, many channels.

This doc proposes: **a medium is a stream, not a function.** Always open, bidirectional, multiplexable. Relations flow over it in both directions. Many logical streams share one conduit.

## What changes

### From function to conduit

Current:

```go
type Medium func(present string, r Relation) (string, error)
```

Proposed:

```go
type Medium interface {
    In() <-chan Relation   // relations coming from the medium
    Out() chan<- Relation  // relations going to the medium
    Close() error
}
```

Every medium is a long-lived stream with two channels. The runtime reads from `In` to receive, writes to `Out` to send. Either direction can be empty for mediums that don't use it (actuator has no `In`, sensor has no `Out`).

### Relations carry a stream tag

Multi-source streaming needs framing:

```go
type Relation struct {
    ID       string
    Origin   string
    ThreadID string
    Stream   string  // <— new: logical channel name within the medium
    Impulse  string
}
```

Video frames go on `Stream: "video"`. Audio on `Stream: "audio"`. Text on `Stream: "text"`. Protocol commands on `Stream: "control"`. Subscribers listen to the streams they care about.

### Turn-taking is a protocol convention, not a runtime rule

Currently the runtime blocks on `medium(present, r)` — waits for the response. That forces a turn-based cadence.

With streams, the runtime doesn't block. It dispatches whatever arrives. Turn-taking becomes something beings negotiate via protocol:

- A being emits relations while they have something to say.
- A being emits `skyra end-turn | holding` when they're done (optional marker).
- Another being starts when they're ready.
- Interrupting is just... emitting while someone else is still emitting.

The runtime doesn't enforce whose turn it is. Beings coordinate through the protocol.

### Continue-thread becomes stream-aware

Instead of:

```go
response, err := medium(present, r)
// parse response, route each line
```

Becomes:

```go
medium.Out() <- buildRequestRelation(r)
for rel := range medium.In() {
    if rel.Stream != "control" { continue }  // or handle video/audio separately
    route(rel)
}
```

Continue-thread opens the stream, attaches a listener, dispatches each arriving relation as it comes in. For short-lived mediums (a one-shot inference call), the stream closes after the response finishes. For long-lived mediums (a camera feed), the stream stays open indefinitely.

### Main.go becomes an event loop

Currently `main.go` kicks off a relation and returns when the recursion completes. With streams, it kicks off, then enters an event loop:

```go
for {
    select {
    case rel := <-worldIn:
        dispatcher.Relate(rel)
    case <-shutdown:
        return
    }
}
```

Any medium pushing relations reaches the runtime through the world's inbox. The runtime dispatches them. No single call-stack bounds the lifetime of a turn.

## How it unifies the open questions

### With "world as entity"

If worlds are entities with mediums, and mediums are streams, then:

- **Cross-world communication** is opening a stream to another world. Two worlds connect via a medium, relations flow both ways, same primitive as a being talking to a medium.
- **Remote worlds** work trivially. A world running on another machine is a medium that happens to be `tcp:host:port`. Protocol flows over the stream the same way it flows locally.
- **Nesting** becomes composition. World A has world B in its EntityMap; talking to B opens a stream into B; B's runtime picks it up.

### With one-way mediums

Sensors (input-only) are streams with `In()` that emits, `Out()` that's nil or dropped. Actuators (output-only) are streams with `Out()` that receives, `In()` that's nil. Bidirectional streams have both. One primitive, three patterns.

### With LLM streaming

Inference becomes a medium whose `In()` emits protocol lines as the LLM generates them — token by token, framed by `| reason\n`. The first protocol line arrives before the full response completes. Continue-thread routes it immediately. By the time the third line arrives, the first may have already triggered a downstream response. Pipelined, not blocked.

## Multiplexing

One physical stream, many logical streams. A human's medium carries:

```
~stream text ~say hi there | speaking
~stream video ~frame <base64 jpeg> | seeing
~stream audio ~chunk <pcm data> | hearing
~stream gesture ~vector [x, y, z] | pointing
```

Each tagged with its channel. Receivers filter by `~stream`. A text-only being reads `text`. A vision being reads `video`. A whole-body being reads all of them and fuses.

This is how bodies actually work — many senses, one nervous system, the brain multiplexing. And it's how HTTP/2, WebRTC, gRPC actually work — one connection, many logical streams.

## Implementation notes

### Framing

Relations are newline-delimited. `| reason\n` is the frame boundary. A stream reader buffers until it sees `| reason\n`, then emits a `Relation`. Same semantics as current code but applied to a continuous byte stream instead of a final string.

### Backpressure

If a medium's `Out()` channel is full, the runtime blocks trying to write. This is natural backpressure — fast producers wait for slow consumers. Like Unix pipes with their kernel buffer.

### Lifetimes

Short-lived mediums close their streams after one turn. Long-lived mediums stay open. The runtime tracks open streams per being, closes them when beings are removed from the EntityMap.

### Concurrency

Goroutines per medium. Each medium's Send/Receive runs in its own goroutine, feeding the world inbox. The runtime's event loop consumes the inbox sequentially (preserving dispatch order). This keeps the model simple — beings don't need locks; their state is mutated in a single goroutine (the dispatcher).

## What it costs

Big refactor.

- **`Medium` becomes an interface, not a function type.** Every existing medium (`inference`, `cli`, `shell`, `exec`) needs to be rewritten as a stream.
- **`Relation` gains `Stream` field.** Touching every construction and parse.
- **Turn semantics change.** Right now continue-thread is recursive and synchronous. It'd become async, event-driven.
- **`main.go` becomes an event loop.** No more "kick off and wait for recursion to end."
- **Protocol parser handles framing.** Impress needs to work incrementally on a stream, not just a final string.

Probably a day or two of careful refactoring. Worth doing with tests.

## What it buys

- **Sensors and actuators fit naturally.** No special case for cameras, mics, displays.
- **Streaming LLMs work.** Token-by-token pipelining becomes native.
- **Multi-modal beings work.** One medium, many logical streams.
- **World-as-entity gets its transport story.** Mediums are how worlds connect.
- **Continuous interaction becomes possible.** The runtime is no longer turn-bound.
- **Backpressure, multiplexing, pipelining** come for free from Go's channel model.
- **The Linux analogy tightens further.** Mediums become file descriptors in the full sense — streams you read/write, multiplexed via epoll-style loops.

## Sources: labeled tracks within a being's stream

A being has one medium — one conduit out to the world. But that conduit can carry many labeled sub-streams, which we call **sources**. Everything crossing the medium is tagged: this one comes from the text faculty, that one from video, that one from audio.

It matches how bodies work. A human has one body, many faculties. Voice and gesture both come from you, labeled. It matches how networked streams work — WebRTC has one peer connection with many tracks; HTTP/2 has one TCP socket with many logical streams. One conduit, many named channels.

### The shape

A relation gains a `Source` field:

```go
type Relation struct {
    ID       string
    Origin   string   // the being this comes from
    Source   string   // the labeled track within their stream
    ThreadID string
    Stream   string   // (optional) stream-wide channel if also multiplexing at that layer
    Impulse  string
}
```

`Origin` answers *who*. `Source` answers *which faculty of them*. The pair identifies the wire.

A glasses-wearing Michael's medium could carry:

```
~origin michael ~source video ~frame <base64> | streaming
~origin michael ~source audio ~chunk <pcm> | streaming
~origin michael ~source text ~say hi | speaking
~origin michael ~source gesture ~vector [x, y, z] | pointing
```

Four sources, one medium, all from Michael. A video-only listener filters `Source == "video"` and ignores the rest. A whole-body listener consumes all of them.

### Why this is the cleaner decomposition

Earlier thinking treated "sources" as a separate primitive — a list attached to a being, configurable, mountable. This collapsed once we realized:

- **Sources are labels on relations crossing a being's medium.** Not a new thing to add. Just a field.
- **Pulling from other beings is already done via relationships.** When Skyra "pulls" from bash, she's calling bash's medium. Bash's medium emits relations with `Origin: bash, Source: stdout`. That's the whole mechanism.
- **One medium per being still holds.** The medium is the conduit. What it carries can be multi-source without adding another primitive.

No separate `Source` struct, no `attach-source` operator, no new category. Sources are faculties of beings, expressed as labels on their outbound relations.

### What mediums do with sources

- `cli` emits `Source: text` only.
- `inference` emits `Source: text` only.
- `shell` emits `Source: stdout` (or `stderr` for diagnostics).
- `exec` (compiled being) emits whatever sources that binary is designed to emit.
- A hypothetical `glasses` medium emits `Source: video`, `Source: audio`, `Source: gesture`.

The medium defines what sources its being can carry. The runtime doesn't enforce — it just delivers labeled relations.

### Receivers filter

Continue-thread and other routing gain source-awareness. Filtering by source is trivial:

```
skyra continue-thread ~with michael ~source text ~say hi | chat
```

Addresses only Michael's text-faculty. A vision-being addressing Michael might route `~source video` relations; a text-being ignores those.

### Inbound and outbound are symmetric

The being's own medium is their outbound conduit — their emissions leave labeled by source. When a being receives, they're consuming *other* beings' outbound streams. Each of those beings has their own sources. Skyra attending to Michael + bash sees:

- `Origin: michael, Source: text` — his words
- `Origin: michael, Source: video` — his camera
- `Origin: bash, Source: stdout` — command output

Her inbound is the union of others' outbound. She filters by what she cares about.

### What this collapses

- No separate "sources" primitive.
- No `attach-source` / `detach-source` operators. Mount/unmount becomes adding or removing a relationship (already an operator we have).
- No forking of the ontology. Beings, operators, invariants, mediums, relations — still the five.

Sources are a label on a field. The ontology stays the same; the expressive power grows.

### What it costs

Small. Add one field to `Relation`. Update mediums to set it. Update routing operators to filter on it if they care.

Can ship in an afternoon when we need it.

## Decision point

**Not urgent.** The turn-based runtime works for what we're doing today.

**Urgent when:**

- We want to stream inference (any serious production Skyra will want this).
- We want sensors or actuators (cameras, speakers, real-time input).
- We want parallel conversations inside one being (mix-speak while listening).
- We want cross-machine transport (remote worlds).
- We want interactivity where the being can interrupt or be interrupted.

At that point, streams aren't an enhancement — they're the substrate.

**Status:** spec'd here. Deferred until first real need. Likely lands together with world-as-entity, since they solve different faces of the same problem: how does a world get addressed, and how does data flow across its boundary.
