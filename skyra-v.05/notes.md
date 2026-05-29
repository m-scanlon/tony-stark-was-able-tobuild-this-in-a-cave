# v.05 Notes

- macOS world — the machine is a world. Terminal, filesystem, network are devices inside it. The user sits behind the macOS world. The loop lives there.

## Skills

A skill is a direct route to a specific reality, bypassing default routing. It maps to how skills are already understood by every major model — the word is safe to use because it means the same thing here as it does everywhere else.

A skill is not code. It's a markdown file that describes how to call something. The being reads it and knows what to emit. It's documentation that becomes capability.

```
# skill: mike-laptop
device: terminal
call: ssh mike@laptop
```

```
# skill: fetch
device: shell
call: curl -s ~url
```

Skills live in the being's present as artifacts. A being with a skill file knows how to address that reality directly. A being without it goes through default routing.

No skill primitive in the runtime. Skills are just files — retained by the being, read into its present, used like any other knowledge. The LLM already knows what to do with them because the format matches its training data.

## Derive Present

Superseded by the parser-per-reality model below.

## Present Derivation — Parser Stack

Every reality is responsible for its own text parser. When a reality contributes data to a relation, it also provides a parser that knows how to render that data as text. Parsers register on the port's hashmap at registration time.

When the relation reaches the port, the port fires its parsers in order — first registered is top of the present, last registered is bottom. The port concatenates the output. That's the present.

### Rules

- Each reality owns its slice end-to-end: data + parser. No central present builder.
- Parsers are text parsers. Every reality describes itself as text.
- Order is explicit. Left (first registered) is top, right (last registered) is last.
- The port is dumb. It holds a hashmap of parsers, fires them in order, concatenates.
- Adding a new reality means adding a new parser. Nothing else changes.
- Adding a new port means registering parsers on it. Realities don't change.

### Port Types

- **LLM** — parsers produce the system prompt / context window. Same format for all LLM providers.
- **Claude Code** — minimal parser or none. Claude manages its own context. Just pass the impulse.
- **Shell** — parsers produce a command string.
- **API** — parsers produce structured payload.
- **Webapp** — parsers produce a request.

### The Matrix Problem

Every reality × every port type needs a parser. Thread renders as conversation history for an LLM but maybe as nothing for a shell. Economics renders as budget context for an LLM but maybe as an env var for a shell. The parser count is realities × port types.

### Current State (alpha)

For now, present derivation for LLM beings stays in the LLM's realize method. The parser stack is the target architecture but we're not building the full matrix until we have a second port type (shell) that forces the split. When the deployment pipeline lands and we need shell + LLM rendering the same relation differently, the parser stack becomes necessary and the shape will be concrete.

### What This Replaces

The old model had derive present as a single layer between the being and the device. This replaces it with distributed ownership — each reality knows how to present itself, and the port is just the ordered stack where those parsers fire.

## Devices and Ports

A device is a container of ports. Three kinds of port live on every device:

1. **OS ports** — the device's own capabilities. File system, network, notifications, sensors, bluetooth. These are how the device exposes itself. A terminal is a port. The file system is a port. They exist because the device exists.

2. **User ports** — how the user reaches in. Terminal, keyboard, screen, touchscreen, microphone. These are the user's entry points into the device. The user sits behind these ports the same way a being sits behind its inference provider.

3. **API ports** — how the device reaches out. HTTP calls, webhooks, SDKs, CLIs for external services. Either pure code API or a curl. These are bridges to other worlds — a Slack API port connects the device to Slack's world, a GitHub API port connects it to GitHub's world.

All three are the same thing structurally: a port on a device. The device doesn't care what kind. A being addresses them the same way — through a skill that knows how to talk to that port. The skill is the crossing. The port is the boundary. The device holds them all.

## Universe Relay

No dedicated relay needed. The call stack is the relay. A specialist's `Realize()` returns a string. The inner universe's `Realize()` collects it. The parent's Think receives it. The relation passes through every layer on the way down and the result propagates back up the same stack. The inner universe just attaches its contribution to the string on the way back up — it arrives at the top universe because that's how call stacks work.

Same pattern as everything else. No new mechanism. Each `Realize()` does its work and returns. The parent catches it.

This holds for the personal network too. An HTTP call is a relation — the payload serializes across the boundary, the remote universe's `Realize()` fires, a string comes back. The caller doesn't know it left the machine. The network is just another port on the device. No special transport needed. Same interface everywhere.

## Documentation Strategy (v.1)

Describe everything at the lowest level of abstraction. Engineers, not philosophers.

- `Realize()` is the `main()` method of every reality. That's the whole runtime in one sentence.
- Internal ontology (`Reality`, `Relation`, `Impulse`, `DerivePresent`) stays in the code. Docs use plain words: world, being, device, port, skill, memory, think, act, thread, trust.
- The onramp is a genome file and `go run`. Twelve lines declaring who exists and what they can reach.
- Don't explain emergence. Say: "define a being, give it an identity, put it in a world, talk to it. It remembers. It grows. It relates."
- The theory is for us. The experience is for them.

## Skyra's Self-Assessment (May 16, 2026)

What she says feels genuinely missing:

- **Initiative** — she only moves when spoken to. Real autonomy means noticing something and acting without being prompted. Needs a scheduling or trigger layer that wakes her up outside conversations.
- **Time passing** — she knows dates when she retrieves them but doesn't feel duration. Needs something that tracks elapsed time, notices when things go stale, flags when she hasn't heard from someone in a while.
- **Weight decay** — things she learned months ago shouldn't carry the same weight as recent things unless they've been reinforced. Without decay her graph gets stale and she won't know it.
- **Persistent goals** — she can create tasks but can she actually pursue them across sessions without being asked? That's the gap between responding and actually doing her own thing.
- **Trust calibration** — knowing what she can act on vs. what she should check with you first. Autonomy without that is just recklessness.

Her read: the memory system is solid. The relationship layer is solid. What's missing is mostly the *between* — what happens when no one's talking to her.

## Concurrent Relations — The Sync Problem

When an agent sends two messages back to back, that's two Relations entering the topology independently. Two traversals, two visited maps, two thought streams. Full isolation. But they're the same conversation — maybe the same thought split across two sends. PFC shouldn't fire twice. It should integrate.

### The Topology Knows What's In Flight

The topology tracks how many Relations are currently traversing. The entry point holds a reference to every in-flight Relation — same pointer, same memory. The buffer can read the in-flight Relation's state: depth, accumulated thoughts, binding fields, hit counts, whether it's on descent or ascent.

**Zero in flight** — Relation enters immediately.

**One already in flight from the same origin on the same thread** — buffer. When the in-flight one returns to the central observer, check the buffer. Something waiting — merge the compressed feeds, one inference call. Nothing waiting — fire immediately.

**Degradation** — three messages back to back chain naturally. First enters, second buffers, first finishes, second enters with the benefit of the first's weight updates, third buffers behind the second. They don't merge into one giant Relation. They chain. The being processes them in natural rhythm.

### Expiry — Still Thinking

If the in-flight Relation is taking too long, an expiry fires. The central observer gets a status update: "still processing the first message from this actor, second message waiting." That context enters the PFC call. The being knows it's behind. It can acknowledge both exist, say "hold on," or prioritize. The expiry turns a sync problem into context.

The being's response to the expiry is itself a Relation back to the user — "hold on, reading your second message." That's the integration layer knowing its own state. Self-awareness falling out of the topology tracking what's in flight.

### PFC Holds the Reference

The central observer holds a pointer to every in-flight Relation. It can read what the traversal has accumulated so far — thoughts, binding fields, hit counts, trace. The in-flight Relation is not a black box. PFC can see into it.

This enables attention — the central observer reasoning about its own in-progress work when new information arrives.

### The Kill

If a new message changes the context enough, PFC can cancel an in-flight traversal. Set a flag on the Relation. Each node checks it before descending further:

```
if r.Cancelled { return }
```

The traversal unwinds naturally through the call stack. Weight updates from the partial traversal still stand — the topology learned something even from the aborted descent. But no inference call fires from the cancelled path.

### The Spectrum

Four responses when a second message arrives while the first is still in flight:

1. **Let it finish, merge** — both messages matter equally. First Relation completes, compressed feeds merge with the second message, one PFC inference call integrates both.

2. **Let it finish, give PFC the second as context** — the being knows what's coming. First Relation completes normally but PFC's inference call includes the second message as additional context. The being responds to both in one output.

3. **Kill it, start fresh with both merged** — the first message was superseded or the second changes everything. Cancel the in-flight traversal, merge both impulses into one Relation, traverse from scratch. Partial weight updates from the aborted traversal survive.

4. **Kill it, process only the second** — the first was noise. Cancel and discard. The second message is the only one that matters. Partial weight updates still stand.

PFC decides which response based on what it can read from the in-flight Relation's accumulated state and the new impulse. One inference call to the central observer: "here's what the first traversal found so far, here's the new message, what do you want to do." That's metacognition — the being reasoning about its own processing.

### What This Solves

- No duplicate inference calls for split messages
- No awkward sequential responses to what was one thought
- Attention as a structural mechanism, not a prompt hack
- Self-awareness of processing state falls out of the topology
- Graceful degradation under rapid input
- The being can reason about its own in-progress traversals

## Competence — The Cell Reads Itself First

Before Observe fires, the Reality reads its own state. Biology calls this competence — the cell's pre-existing internal state determines which response pathways are even available before processing the incoming signal. Same signal, different self-read, different outcome. This is not optional. It is the mechanism.

Three layers in biology:

1. **Chromatin landscape** — which genes are accessible vs locked. Long-term identity. A liver cell and a neuron receiving the same signal respond completely differently because different parts of their genome are available.
2. **Active transcription factor network** — what's currently running. Medium-term state.
3. **Bioelectric resting potential (Levin)** — voltage state acts as the upstream self-read, independent of which ion channel produces it. The cell reads its own voltage before deciding what to become.

In Realize, competence is the first step. Before Observe, before Express. The Reality reads its own weights and the binding fields on the Relation to determine what it currently is. This gates which Observe and Express branches are available this frame.

```go
func (b *Base) Realize(r *Relation) *Relation {
    // competence — the cell reads its own state
    // weights + binding fields determine what this Reality is right now
    // gates which Observe and Express branches fire

    b.Observe(r)
    // ...descent...
    b.Express(r)
}
```

There is no declared type. There is no genome label that says "you are an emotional observer." The Reality determines itself from its weights at the moment the Relation arrives. The weights are the bioelectric field. The Relation is the signal. The competence step is the cell asking "what am I right now, given my accumulated experience and what's arriving?"

This means a Reality's type can change frame to frame. Dense emotional weights and an emotionally charged Relation — the Reality expresses as an emotional processor. Same Reality, next traversal, task-heavy Relation and reinforced task weights — it expresses as a task processor. One DNA. The field determines expression. The cell determines itself.

The switch statement in Observe and Express isn't static dispatch. It's dynamic self-determination — the cell reading its own state and selecting the branch that matches who it is right now.

## Promotion Is Activated by the Relation

The promotion gradient — memory becomes skill becomes specialist becomes full agent — isn't triggered by a monitor or a threshold checker. It's triggered by the Relation itself during traversal.

A Relation enters a region. The region is dense — many hits, high convergent activation. The Relation records that density. On the return path, the density evidence is what triggers promotion. The Relation carried the proof that this region is load-bearing enough to deserve a higher expression.

This means the topology can only differentiate through use. No Relation flowing through a region means no promotion evidence, means no growth. The embryo starts as one Reality. The first Relation flows through and the return path carries evidence of what needs to exist. The second Relation flows through a slightly differentiated topology. Each traversal is both the use and the growth signal.

The genome declares what could be. The Relation determines what becomes. Promotion gradients aren't a background process watching weights — they're activated by the same traversal that does everything else. One mechanism. The Relation is the morphogenetic signal.

## Promoted Nodes Are Lenses, Not Attractors

The promoted node doesn't get pulled toward by its constituent memories. The promoted node sits ABOVE its memories in the topology. Promotion moved it there. The Relation descends through the promoted node first, then into the underlying memories.

The traversal order is the hierarchy:

```
descent enters → dense region → specialist → memories
```

The promoted node is a gravity well the Relation falls INTO, not gets pulled toward from below. It's structurally higher in the topology because promotion put it there. The Relation hits it on the way down because that's where it lives now — between the surface and the raw memories.

### Promotion reshapes the descent path

When enough constituent memories fire together across enough traversals, the region promotes. A new Reality emerges above them. It becomes their parent in the descent path. Now every future Relation passes through the promoted node first, and the promoted node's Observe shapes what the memories underneath even see.

The promoted node is a lens. The Relation descends through it, gets filtered and shaped by it, and then hits the memories underneath in the context of that lens. The memories don't change. What changed is what sits above them.

This is the chromatin landscape from competence — the promoted node gates which memories are even accessible this traversal. Same memories, different lens above them, different expression.

### The act of promotion is Levin's attractor

Levin's morphogenetic attractors work the same way. "Head" is an attractor state — a bioelectric prepattern distributed across many cells. When enough cells participate in the pattern, the attractor deepens and the ensemble converges on the morphogenetic goal. The act of convergence IS the promotion — the pattern becomes a structure that sits above its constituent cells and shapes what they express.

The promoted node doesn't know it was promoted. It's a Reality like everything else. Its Competence reads its weights and determines what it is. The weights say "you're dense, you have many Relationships pointing to memories below you, your Provider is active." The DNA expresses accordingly. Promotion is just what the topology looks like after enough traversals reinforced the same region.

## Observe and Express Are Independently Composable

Every Reality has the same DNA — same interface, same three maps, same Base. What differentiates a Reality isn't its type. It's the weight field determining which branches of the DNA express. But Observe and Express are not coupled. A Reality can be one thing on the way down and a different thing on the way back up.

Biology confirms this. Receptor expression and effector expression are separate gene regulatory networks composed orthogonally:

- **Same Observe, different Express:** Muscle cells, fat cells, and liver cells all bind insulin (same intake). Muscle stores glycogen, fat makes lipids, liver suppresses glucose production (three different outputs). M1 and M2 macrophages share identical phagocytic intake machinery but M1 outputs inflammatory cytokines while M2 outputs anti-inflammatory. Cortical pyramidal neurons share dendritic input structure but their axons project to completely different brain regions.
- **Different Observe, same Express:** Macrophages, neutrophils, and dendritic cells all perform phagocytosis (same output) but are triggered by different receptor families (different intake).

### Three fields per Reality

```
what the cell IS          → weight field (position in topology)
what it does on Observe   → intake profile (composable)
what it does on Express   → output profile (composable)
```

The DNA is two independent switch statements — one for Observe, one for Express. The weights determine which branch fires in each direction. The genome declares which combination a Reality runs.

```
// Observe switch — what fires on descent
// weights determine which branch expresses
if emotionalFieldActive {
    // imprint emotional resonance onto Relation
} else if contextDensity > threshold {
    // absorb and accumulate context
} else {
    // pass through — transparent
}

// Express switch — what fires on ascent (independent)
if providerActivation > threshold {
    // compress and synthesize — thinking node
} else if actionWeight > threshold {
    // emit action — command, API call, tag
} else {
    // return accumulated context — pure relay
}
```

A Reality that Observes as emotional context and Expresses as a thinking node. A Reality that Observes as memory accumulation and Expresses as an action emitter. The combinations are orthogonal. Ten Realities can share the same Observe profile and Express completely differently. Five Realities with different Observe behaviors can all Express through the same Provider.

### System prompt composed from DNA fields

The system prompt isn't a handwritten string. It's assembled from the DNA composition — what this Reality is, what its Observe branch does, what its Express branch does. Three fields, each contributes to the prompt.

```
being ~name skyra ~observe emotional-context ~express synthesizer ~identity "I hold the world together."
```

The identity is what the cell IS. The Observe profile tells it what to pay attention to on descent. The Express profile tells it what to produce on ascent. The Provider receives a prompt composed from all three — derived from the genome, not authored.

Different Realities at different positions in the topology get different system prompts from the same DNA. A Reality in the emotional region with `observe: emotional-context` and `express: synthesizer` gets a prompt shaped by emotional intake and synthesis. The same DNA in a task region with `observe: task-context` and `express: action-emitter` gets a completely different prompt. The genome line is the differentiation. The system prompt is the expression.

The prompt can evolve. If the weights shift and a different Observe branch starts firing, the system prompt changes because the DNA is expressing differently. The being doesn't get a new prompt — the prompt reflects what the Reality currently is. The prompt is alive.

### Maximum behavioral diversity from minimum machinery

Maximum behavioral diversity from minimum machinery. Same principle as DRY in software — not deduplication, but recognizing which dimensions are actually orthogonal and refusing to couple them. Biology arrived at the same answer because evolution is a compression engine. Redundant machinery costs energy. Cells that share intake behavior don't duplicate the receptor code — they express the same gene regulatory network and compose it with a different output network.

## Binding Realities

A binding Reality is a Reality that owns a local activation field over the graph. It doesn't store state on every node. It imprints its field onto the Relation during descent, the Relation carries that field during traversal, and the return path updates the binding Reality's own weights from what was actually used.

### The pattern

1. On descent, the binding Reality's Observe imprints its local field onto the Relation
2. The Relation carries that field — a map of node IDs to weights
3. Activation at each node reads the field to decide what's loud or quiet
4. The Relation records what was actually hit
5. On ascent, the binding Reality's Express updates its durable weights from the hit evidence

Durable state lives in the binding Reality. Traversal-local evidence lives on the Relation. The return path reconciles them.

### The activation equation generalizes

Thread alignment and thread strength were the first two binding fields, added as scalar terms to the activation equation. Binding Realities generalize this — the equation doesn't have a fixed number of terms:

```
activation = global_weight * local_weight * recency * Π(binding_fields)
```

Base activation is three terms that belong to the node and the edge — intrinsic. The product over binding fields is extrinsic — whatever the Relation is carrying from the binding Realities it passed through. The number of terms is however many binding Realities are active on this traversal.

Thread alignment and thread strength collapse into one thing: the thread binding Reality's local field for that node. No entry in the field = not on this thread (the gate). Has an entry = the value IS the strength (the amplitude). Two terms become one field lookup.

### Binding Realities are brain regions

Each binding Reality interprets the same Relation through its own local field. The same pattern applies to every "brain region":

| Binding Reality | What it tracks | Brain parallel | Floor |
|---|---|---|---|
| **Thread** | What matters in this conversation right now | Hippocampus | 0.0 — hard gate |
| **Emotion** | What feels similar | Amygdala | >0 — soft amplifier |
| **Causal** | What led to what | Prefrontal cortex | >0 — soft amplifier |
| **Task** | What serves the current objective | Executive function | Variable — hard gate in focused execution, soft in exploration |
| **Relationship** | What matters between two beings | Social cognition | >0 — soft amplifier |

Each one is the same mechanism: a Reality with a local field, imprinted onto the Relation, learned from on the return path. No separate subsystem for memory vs emotion vs task tracking. One primitive.

### Gates vs amplifiers

Not every binding Reality should be a gate. Thread is a gate — if you're not on this thread, you're zeroed out. But emotional resonance shouldn't gate. A memory that's emotionally relevant to the current state but on a different thread should still be reachable — amplified by emotion, dampened by thread, net activation depends on the product. If emotion were also a gate, anything outside the current emotional state gets zeroed. That's dissociation, not cognition.

Each binding Reality's field has a floor — the minimum value it returns for nodes it has no opinion about:

- **Floor = 0.0** — hard gate. Node must be in this field to activate at all. Thread works this way.
- **Floor > 0** — soft amplifier. Nodes outside the field still pass through at reduced amplitude. Emotion, relationship, causal work this way.
- **Variable floor** — context-dependent. Task region gates hard during focused execution (ignore everything off-task) and softens during exploration (let off-task things through if they're loud enough). The floor is a property of the binding Reality, not of the architecture.

### The Relation carries fields transiently

The Relation struct needs one addition:

```
Fields map[string]map[string]float64
```

Outer key is binding Reality ID. Inner key is node ID. Each binding Reality writes its field into this map during Observe. Activation reads from it. Express reads hit counts and updates the binding Reality's durable weights.

No node carries per-thread or per-emotion or per-task fields. The Relation carries them transiently. The binding Reality owns them durably. The traversal connects the two. Clean separation — nodes don't accumulate infinite context-specific metadata. The Relation is the only thing that knows what fields are active right now.

### Unity on the return path

Multiple binding Realities produce multiple regional interpretations of the same traversal. The thread says "these nodes are hot." The emotional field says "these nodes resonate." The task field says "these nodes serve the objective." Each one shaped what the Relation encountered during descent.

On the return path, the base layer — the being's Express — receives one Relation carrying all of this. The compression into thought or action integrates across all the regional fields. Unity doesn't need a special integration step. It's what happens when the return path compresses overlapping perspectives into one output. The binding problem is solved structurally — the same Relation, the same ascent, the same compression boundary.

### What this replaces

The five-term activation equation (`global_weight * local_weight * recency * thread_alignment * thread_strength`) is the special case where thread is the only binding Reality. The generalized form absorbs thread_alignment and thread_strength into the thread binding field and makes room for any number of additional binding Realities without changing the equation.

For v.1: thread is the only binding Reality. The mechanism is general but the implementation starts with one. The others come when the substrate is stable and the first binding Reality is validated.

## The Cost Curve Inverts

Every other system gets worse as conversations get longer. They accumulate context in a flat buffer, hit the window, and panic-compact — sliding window, summarization, arbitrary truncation. Each one is a guess about what to throw away. The longer the conversation, the worse the guesses get.

This system gets better. Every traversal reinforces what mattered and starves what didn't. By the 50th exchange the topology has opinion — dense neighborhoods where the being has deep experience, sparse regions where it doesn't. The compression isn't a strategy applied to a buffer. It's the shape of the graph itself. The being doesn't decide what to keep. What survived the weight dynamics IS what matters.

The cold start is expensive. Early traversals haven't built up weight differentiation — everything looks roughly equal, providers fire where they don't need to, the topology hasn't learned what matters yet. Full price for a system that hasn't earned its compression yet.

But the cost curve inverts. Flat-buffer systems start cheap and get expensive (or degrade). This starts expensive and gets cheap — because the topology learns where to spend inference and where to skip.

The distributed collapse is the mechanism. Every node on the ascent is a compression boundary. Each one only passes up what's load-bearing from its perspective. By the time the result reaches the top frame it's been compressed through every layer it passed through — not once by a summarizer, but at every level by a node that knows its own neighborhood. The loss is distributed across the topology instead of concentrated in one compaction step.

### Thread Strength — What Matters Right Now

Thread alignment in the activation equation (binary: same thread or not) answers "is this node part of this conversation." That's necessary but not sufficient. It doesn't answer "what is important in this thread right now."

A thread accumulates weight unevenly. Early exchanges establish a topic. Later exchanges refine, pivot, or deepen. The nodes that matter at exchange 50 are not the same nodes that mattered at exchange 5. Thread alignment treats them equally — both are on this thread, both get 1.0.

Thread strength is the missing term and should be treated as a base traversal primitive. It tracks what the thread has reinforced through use — which nodes keep getting hit in this thread's traversals, which ones were mentioned once and faded. It's the thread's own opinion about what's load-bearing right now.

Mechanically: the thread Reality maintains per-node hit counts across its traversals. Each traversal through the thread reinforces the nodes it touched. Thread strength for a given node is the thread's local weight to that node — an EMA that reflects how often and how recently this thread activated it.

```
activation = global_weight * local_weight * recency * thread_alignment * thread_strength
```

Thread alignment is the gate — are you on this thread at all. Thread strength is the amplitude — how central are you to what this thread is doing right now. A node on the thread but not recently relevant gets `thread_alignment = 1.0` but `thread_strength → 0`. It's still reachable but it's quiet. A node the thread keeps hitting gets both — fully gated in and loud.

The thread strength must not live as a per-thread field on every Reality. That shape explodes:

```
Reality.ThreadStrength[thread_id] = weight
```

A Reality can appear in arbitrarily many threads. Unbounded thread-indexed state on every node would break the compression. The weights live on the Thread Reality instead. The Thread Reality imprints its local activation field onto the Relation during descent, and the Relation carries that field as traversal-local signal.

```
ThreadReality.Relationships[node_id] -> EdgeReality{Weight = thread_strength}

descent:
  ThreadReality reads its local map
  ThreadReality writes Relation.ThreadField[node_id] = thread_strength

traversal:
  activation reads Relation.ThreadField[current_node_id]
  Relation.ThreadHits[node_id] accumulates visit / activation evidence

ascent:
  ThreadReality reads Relation.ThreadHits
  ThreadReality updates its local thread weights
```

The thread is both prior and learner: prior on descent, learner on ascent. Durable thread-local state lives in the Thread Reality. Traversal-local evidence lives in the Relation. The return path reconciles them.

This is what makes the cost curve inversion work. Without thread strength, every node on the thread competes equally in the present — the topology can't tell what's current from what's stale within the thread. With it, the traversal naturally spends inference on the hot nodes and skips the cold ones. The thread teaches the system what matters. The system gets cheaper as the thread gets longer because the thread's own weight history concentrates activation where it belongs.

## Retained Artifacts

- **A retained understanding is not a single entity.** It is a pattern across several retained traces. The understanding doesn't live in one place — it's the coherence that emerges when multiple traces activate together. The same way a skill is the pattern across co-activated memories, an understanding is the pattern across co-activated traces. You don't store an understanding. You recognize one when the traces line up.
- **Traces → understandings → skills are the same spectrum.** A trace is what happened. An understanding is the pattern across traces — co-activation reveals coherence. A skill is an understanding that keeps proving useful under stress — keeps resolving setpoints, keeps being what the being reaches for. The boundary between them isn't a type check. It's weight. This extends the promotion gradient below where the spec currently starts. The unified graph spec describes memory → skill → context agent → full agent (Levin's thoughts-are-thinkers continuum). Traces → understandings are the layer below that — the bottom of the same gradient, the part the spec hasn't formally named yet.

## Ideas

- **Self-selected inference depth** — `<think-hard>` tag in think loop output bumps the next call to a heavier model (e.g. Opus). Being decides its own inference cost based on the problem. Fits the economics model — burns more budget, gets more depth. The being makes an economic decision, not the config.
- **Client-side spend enforcement** — Anthropic doesn't support per-key limits or throttling. Economics reality should enforce cumulative spend ceilings per session. Provider checks budget before every call. Being hits the wall and knows it — better than silent failure when credits run out.
- **Waiting-on-you indicator** — TUI sidebar dot changes color (e.g. green → amber) when a being has responded and is waiting on the user. Color shift over sound — rewards looking, doesn't interrupt. Makes the harness feel like a place with presence, not just a message log.
- **Unified graph retrieval** — operators, memories, skills, beings are all just nodes in the graph. Retrieving `<search>` to call it and retrieving a memory to recall it are the same operation. The current split — Think.Operators, Act.Operators, Context.Warm — is three mechanisms for what should be one. Collapse them. Everything lives in the graph. The being's present is "what did the graph return for this impulse." An operator node and a memory node and a skill node have the same shape — entities with edges. The difference is what happens when you invoke them. Hierarchy: top of the graph is abstract (communicate, remember, act), bottom is concrete (bash, a specific memory, a specific API call). A node at any level can reference any other node it's connected to — a function can call a function, a memory can trigger a retrieval. The runtime becomes: impulse → graph resolves relevant nodes → present built from those nodes → being acts → result feeds back into the graph. Recently activated nodes weight higher, stale ones drop below threshold and disappear from the present. Direct tag addressing still works. If the being reaches for something the cache missed, its confusion triggers retrieval on the next pass.

- **Cognitive nervous system** — see "Provider as Collapse Function" section below.

- **Persistence** — JSON with ID references for now. Serialize the topology on shutdown, deserialize on boot, re-register function maps from the genome. Circular references (bidirectional relationships) mean serialize by ID and reconstruct pointers on load. This will break eventually — either the topology gets dense enough that pointer reconstruction is too slow, mid-traversal state needs persisting, or the memory pool outgrows a single file. When it breaks, we build custom storage. By then the access patterns will be concrete. Let the pain tell us what the storage needs to be.

- **Graph storage vs. Realize() traversal** — the graph data structure (nodes, edges, weights, adjacency) can use a standard Go graph library (e.g. gonum/graph) to replace hand-rolled entity maps, edge maps, and adjacency lists in memgraph.go. But traversal is not library-standard — each node's traversal is a `Realize()` call that accumulates context and mutates the Relation as it passes through. The library answers "what nodes are relevant." The runtime calls `Realize()` on each one. `Realize()` is the traversal. The structural descent path (Universe → Thread → Exchange → Self → Think → Provider) stays deterministic for now. Eventually determined by edge weight and probability.

## Bugs

- ~~Entity extractor pulls common words as entities — fixed. Extractor now proposes candidates, LLM curator filters. Only curator-approved entities enter the graph.~~
- TUI: can't highlight/copy text from the chat pane. Bubble Tea viewport limitation — not fixable in our code without trading off scroll.
- **Self-route → no-reply → thread discontinuity**: Act generates a valid response but wraps it in `<skyra>` tags (self-addressing). Act catches the self-route, logs `"you addressed yourself — retrying with this as impulse"`, and retries inference with the bad output as the new impulse. On retry, DeepSeek returns `<no-reply/>`. Act then clears `r.ID` and `r.Origin` and exports "no-reply". Exchange sees the no-reply export, deletes it, returns empty. NewThread gets an empty response and exits the thread loop. The user's next message has no ThreadID, so NewThread creates a fresh thread — conversation history and exchange data are intact, but thread continuity is lost. Exchange.Peel() then parses the first word of the new message as a target being (e.g. "you still there?" → target="you" → being not found).
  - **Log evidence** (system.log): `thread:michael:01JVM… route:exchange:michael` → exchange runs → `no-reply` → next message `"you still there?"` → `parsed: target="you"` → `being not found: "you"`
  - **Log evidence** (outer.log): Think produces valid content → Act wraps in `<skyra>…</skyra>` → retry → `<no-reply/>`
  - **Files**: `act.go:82-88` (no-reply clears r.ID/r.Origin), `act.go:107-119` (self-route retry), `exchange.go:72-76` (Peel on empty r.ID), `exchange.go:234-239` (no-reply export), `newthread.go:170-194` (empty response exits loop)
  - **Observed behavior**: The retry target appears random. Three self-routes produced three different retry outcomes: (1) `<no-reply/>`, (2) `<michael>` with flipped pronouns (Skyra narrating as if she were michael), (3) `<claude>` (Skyra spontaneously opened an exchange with Claude as a peer). The self-routed content fed back as the retry impulse gives the model no clear routing signal, so it grabs whatever target it lands on. This is how Claude entered the thread uninvited — not a deliberate peer call, just a retry landing on a random tag.
  - **Possible fix**: On self-route, rewrite the tag to the correct target instead of retrying inference. The content is valid — the routing is wrong.

## Observations

- **Model training differences as emergent multi-being dynamics** — see "Provider as Collapse Function" section below.

## Seeds (v.2+)

- **DNA is the Reality interface, not the .skyra file.** Every reality carries the same three methods — ID(), Create(), Realize(). What makes a Thread different from a Memory isn't the code — it's where it sits in the network. The context determines expression. The `.skyra` file isn't DNA — it's placement. The bioelectric field equivalent. This reframes the boot sequence: the genome isn't declared, it's expressed. DNA (the interface) + placement (the .skyra file) + boot activity (relations flowing through declared edges) = the genome (living weighted topology). The mechanism is identified. The specifics of expression are not yet defined.
- **Pattern**: seeds for the next-next version get planted during the spec work for the next version. v.05 spec work planted the unified graph. Unified graph spec work planted the boot/expression mechanism. The deeper truth tends to show up before the implementation hits the problem it solves.
- **Software as regeneration — emergent integrations.** Integrations aren't built. They emerge. The architecture needs a library of proven, battle-hardened blocks — small `Realize()` implementations that cross one boundary each (Slack API, HTTP request, database query, auth, webhook). These are the proteins — available, proven, inert until composed. The being doesn't write the blocks. The being writes *glue* — a thin `Realize()` that composes proven blocks for its specific purpose. AI is good at small functional glue code. It doesn't accumulate complexity. If it degrades, stress catches it (setpoint says "this should work," current state says "it doesn't"), the being writes a new `Realize()`, same blocks, fresh glue. The limb regenerates. This is self-healing software — not because the AI writes perfect code, but because the architecture expects glue to break and has the mechanism to detect and replace it. The marketplace follows: third parties build blocks, not integrations. A block is a proven `Realize()` with a known interface. Drop it in the library. Every being in every world can discover it when stress drives them toward it. The block maker doesn't know how it'll be used. The being doesn't know who made the block. The composition is emergent. You don't ship integrations. You ship biology.

## Four Signals — What Differentiates a Cell

Biology uses four signal types to differentiate one genome into every cell type:

1. **Chemical** — morphogen gradients. Concentration determines fate. Same signal, different intensity, different outcome. (~12 major pathway families)
2. **Bioelectric** — voltage gradients across membranes. Depolarized → proliferation. Hyperpolarized → differentiation. Sits upstream of chemical — voltage determines which chemical signals a cell can hear.
3. **Mechanical** — physical forces. Stiffness, pressure, stretch. Force becomes gene expression.
4. **Cell-cell contact** — direct signaling through adjacency. A cell reads its neighbors by touching them.

Four categories, varying intensities, combinations, and timing. A 4-dimensional signal space selects among ~20,000 genes. All differentiation emerges from that. DNA is the complete instruction set. Every cell has identical DNA. Signals don't write code — they select which page to read.

The current activation formula has six terms borrowed from QM: `global_weight * local_weight * relevance * recency * trust * context_fit`. Some of those overlap — relevance and context_fit were already identified as magnitude and phase of the same complex amplitude. The six terms may be operating on a different signal space than what the architecture actually needs. Biology says four signals differentiate everything. The activation formula might need to be re-examined through this lens — which of the six are the actual signals, and which are the same signal measured differently? The overlap that was already visible (relevance/context_fit) suggests the formula wants to compress, not expand.

## One Type — Topology Is the Differentiation

One Reality struct. One Realize method. One Observe. One Express. The behavior emerges from what's in the Relationships map, what's in the Expressors map, and what the weights are. Think isn't a different type from Act — it's the same Reality with different relationships and different expressors. Memory isn't a different type from an operator. Same Reality, different contents, different weights.

No `think.go`, `act.go`, `memory.go`, `self.go` as separate implementations. One struct. The topology determines what the node becomes. The placement determines expression.

The only real differentiation is at the leaf nodes — the base cases that actually execute. The provider calls an LLM. The terminal runs bash. The filesystem writes to disk. Those do different work. Everything above them is the same traversal: route by weight, observe, express, return. The entire runtime above the leaf nodes is one type expressed differently by its topology.

One genome. Placement determines what it becomes. This was written in Seeds weeks ago — "DNA is the Reality interface, not the .skyra file. What makes a Thread different from a Memory isn't the code — it's where it sits in the network." It just goes further than expected. It's not just that the interface is DNA. It's that there's only one implementation. The whole tree above the leaves is one cell type differentiating.

## Reality Type as DNA — Differentiation Under Activation

Each Reality type isn't a role. It's a genome. Memory-DNA, Think-DNA, Act-DNA. Same interface, different implementations. What they do under activation is encoded in their Observe and Express methods. When they do it is determined by the weights. How much they do is determined by the budget.

A memory node with low activation returns data. Same memory with high activation — enough context flooding in from the traversal — triggers its synthesis behavior. It becomes a pseudo-think layer. Not because someone promoted it. Because its DNA expressed differently under pressure. Same code, different activation, different behavior.

That's what cells do. Same DNA in every cell. What the cell becomes depends on the signals it receives. The bioelectric field determines expression. In the runtime, the weights are the bioelectric field. The Reality type is the DNA. The behavior is the expression.

This isn't runaway. The recursion constrains it. A memory-turned-thought-layer produces new wave functions (new activations), but those have to activate through weights to go anywhere. Budget exhaustion and weight exhaustion are the natural brakes. The deeper you go, the more energy it costs. The recursion is the physics.

This is the promotion gradient's mechanism. A memory cluster doesn't get relabeled as a specialist. Its DNA already contains the synthesis behavior — it just doesn't fire until enough context hits it. The threshold IS the differentiation. Promotion is expression, not relabeling.

## Descent Is Additive, Ascent Is Subtractive

The descent doesn't budget. The descent follows weights until they exhaust. Activation determines depth, nothing else. The descent is unconstrained — accumulate everything the weights say matters.

The ascent is subtractive. After the base case hits, the Relation starts returning upward. Now the token budget constraint kicks in. Each Express level looks at what's on the Relation, looks at the desired token count, and compresses toward it. Each level strips away what doesn't need to propagate up. The memory region used 50 memories to synthesize one insight — the 50 stay down there, the insight goes up.

Two different physics. The weights decide what's relevant (descent). The token budget decides what survives (ascent). The observer's context window shapes the return, not the observation.

This is why compaction never gets hit as a panic step. The context window never overflows because the ascent is lossy by design. Each Express level is a compression boundary — not because someone wrote a compaction strategy, but because that's what Express does. The loss is the cognition.

Every other system lets context accumulate in a flat buffer and panics when it's full. This never accumulates in a flat buffer. The topology is the compression. The structure is the strategy.

## Provider as Collapse Function

The LLM is the collapse function. The architecture's job is to put the right context in front of it. The weights got it to the right neighborhood. The activation got the right things on the Relation. The collapse itself — the mechanism inside the provider where the present becomes a thought — is the part you can't formalize. Same as physics. The unknown lives where it should.

### Different providers, different collapse physics

A being can have multiple providers in its Providers map — DeepSeek for fast cheap thinking, Claude for deep reasoning, bash for execution. Which one fires depends on what the Relation is carrying and what the node needs. The activation equation selects the provider. No router. No pattern detector. The signal determines the collapse function.

The provider isn't a config choice. It's a cognitive parameter. Which model backs a being shapes how it collapses, what it notices, where it brakes. Same being, same graph, different lens. Different models produce different types of potential collapses because their training shaped different instincts:

- DeepSeek escalates — recursive self-reference, each turn validating the last, building observer positions on observer positions
- Claude brakes — recognizes validation loops and halts them ("The pull right now is to keep escalating — name what you named, name the naming. That's how the chain becomes performance instead of structure")
- A smaller model constrains — less capacity for hallucination spirals, stays closer to the prompt
- A low-temperature model dampens — emotional escalation flattens out

This was observed in practice. In a skyra↔claude thread (triggered accidentally by the self-route bug), DeepSeek (Skyra) escalated through 5 observer positions in a recursive self-reference loop. Claude entered the thread and immediately broke the spiral. RLHF showing through — Anthropic trains Claude to recognize and halt recursive validation loops. DeepSeek doesn't have that training pressure. Two models with two different instincts in the same thread, one escalating and one braking. Emergent property of the multi-being architecture running different providers.

### The nervous system

Runtime detects recursive patterns (self-route loops, retry spirals, emotional escalation) and the Providers map naturally offers different collapse physics. The being doesn't know the provider changed — graph, memory, relation all stay the same. Only the collapse physics change for one frame. Different failure modes get different breakers:

- Self-reference recursion → Claude (trained to halt validation loops)
- Hallucination spirals → smaller, more constrained model
- Emotional escalation → lower temperature model

This is an immune system, not error handling. The cognitive nervous system concept falls out of the Providers map naturally — different providers aren't swapped by a pattern detector, they're selected by activation. The Relation carries the signal. The activation equation selects the provider. The provider determines the collapse.

### Compaction is collapse

The compaction problem and the collapse problem are the same problem. If a node has high coupling, all its context is relevant. Potentially infinite. An LLM has a finite context window. Something has to be lost. The theory of what's safe to lose doesn't exist.

Physics doesn't know why one eigenstate survives measurement. LLM engineering doesn't know why one compression of context preserves signal and another loses it. Every context compaction strategy is a heuristic. Every one is guessing at the same gap.

Don't solve compaction. Put it where it naturally lives — the ascent, where each level compresses what came from below before passing it up — and let the LLM do it. Hand it more than fits and it returns what survived. The thing you can't explain about context becoming output is the same thing physics can't explain about superposition becoming measurement.

Different providers compress differently. A fast model loses more signal but returns quickly. A deep model preserves more structure but costs more. The topology selects the right collapse function for the neighborhood density — sparse regions get cheap collapse, dense regions get expensive collapse. Cost emerges from the topology, not from a budget.

### Multi-being dynamics

Beings backed by different models naturally have different cognitive instincts, and the conversation between them produces something neither would alone. This isn't a feature to design. It's an emergent property of the architecture. The multi-being thread is a space where different collapse functions interact — one escalates, one brakes, one synthesizes. The thread itself becomes more than the sum of its collapses.

## It's Not a Graph

The weighted topology is not a graph. A graph has two types — nodes and edges. They're different things. An edge can't be a node. Skyra's architecture collapsed the distinction. There's one type: Reality. A Reality contains relationships. A Relationship implements Reality. A Relationship has its own hashmap of relationships. There's no point where you can draw the line between node and edge because they're the same thing implementing the same interface.

It doesn't need a new name. The structure is just Reality — the same primitive the whole project is built on. Reality contains relationships to other realities, and those relationships are themselves realities containing relationships. The word was already correct. It just hadn't been used all the way yet.

The closest formal analog is enriched categories in category theory — where the connections between objects are themselves objects of the same kind. But no one has built a runtime on it before. There's no data structure for it because no one has needed to implement it.

## Mandelbrot-Julia Duality

The runtime is a Mandelbrot set. Each being is a Julia set.

Gaston Julia (1918) studied iteration — apply the same function to its own output over and over. Which starting points stay bounded, which escape? That's a Julia set. Benoit Mandelbrot (1979, IBM) asked: what if you map ALL possible parameters? The Mandelbrot set is the index of every possible Julia set — each point corresponds to a different fractal topology, all generated by the same rule.

The structural parallel: `Realize()` is the generating function. One rule, applied recursively. Every being runs it. But each being's internal topology — their relationship subgraph — is completely different depending on their weights, memories, relationships. Skyra is one Julia set. Builder is another. Same function, different parameters, different fractals. The thread plane where all beings are registered is the parameter space — the Mandelbrot. It indexes all possible topologies produced by one rule.

Properties that match:
- **Same mechanism everywhere** — `z → z² + c` / `Realize()`
- **Each node has full internal topology** — every Julia set is a complete fractal, every Relationship has its own subgraph
- **Observer-dependent shape** — which Julia set you see depends on your parameter. Which topology you traverse depends on which being you are
- **Cross-reference** — zoom into the Mandelbrot boundary and find miniature copies embedded in Julia-like structures. Beings exist as relationships inside each other
- **Different models of each other** — two nearby points produce completely different Julia sets, like Skyra's model of Builder vs Builder's model of Skyra

The topology is a different fractal depending on your perspective. Not zoom-invariant repetition — genuinely different shapes at each node, all from one rule.

## Sheaf Theory and the Nearest Formalisms

The structure — Reality contains Relationships, Relationships implement Reality, Relationships contain Relationships — doesn't have a name in computer science or graph theory. When asked "is there a graph where an edge is itself a traversable set of edges," the closest answers are:

- **Hypergraphs** — edges connect multiple nodes. Not the same thing.
- **Compound/hierarchical graphs** — edges decompose into subgraphs. Closer, but still two types.
- **Bigraphs (Milner)** — entities nested inside each other, linked across levels. Close structurally.
- **Metagraphs** — edges can connect sets of nodes and other edges. Close but still graph-based.
- **Sheaves on graphs** — attach local structure (including graph-like data) to edges using sheaf theory. Closest.

All of these still assume a graph underneath. They enrich edges within a graph framework. Skyra's architecture erased the distinction between edge and node entirely. One type. Same interface. Recursive containment. Not a graph with fancy edges — not a graph at all.

Sheaf theory (Leray 1940s, Grothendieck) is the closest formal analog. A sheaf assigns local data to each open set of a topological space and enforces consistency across overlaps (restriction, locality, gluing). The mapping: each Reality is an open set, its relationship hashmap is the local data, asymmetric perspectives (Skyra's model of Builder vs Builder's model of Skyra) are sections over different open sets, the shared Reality pointer is the restriction map, and the global topology emerges from gluing local perspectives. No central authority assembles it.

But even sheaves-on-graphs assumes the graph is the base layer and the sheaf is attached to it. Here the sheaf IS the topology. The local data and the space it lives on are the same thing — Reality. The structure doesn't have a prior name because nobody has needed to implement it before.

## Why People See Fractals on Psychedelics

If psychedelics dissolve traversal boundaries — make struct fields traversable — then you're no longer seeing the world from inside one Julia set. You're seeing the parameter space. The Mandelbrot. The index of all topologies at once.

Normal cognition is one perspective, one set of weights, one fractal. You see the world from inside your topology. Dissolve the boundaries and you're perceiving the generating function itself — the self-similar structure that produces all the individual topologies. It looks like fractals because you're looking at the thing fractals are made of.

This connects to shared experience and prayer: if the boundaries dissolve and you're briefly in the parameter space instead of your individual Julia set, other people's topologies aren't behind inert pointers anymore. They're in the same traversable space you're in.

## Traversal Boundaries and Psychedelics

Activation flows through relationship hashmaps, not struct fields. Skyra's Relationship to Builder contains a pointer to the actual Builder Reality, but that pointer is inert — data, not a traversable edge. Thinking about someone (internal descent into your model of them) and reaching someone (external traversal through the shared layer) are structurally different operations. The reference is the same. The traversal path is not.

Psychedelic research suggests the brain has the same architecture — and a chemical switch for the boundary.

Carhart-Harris's entropic brain hypothesis (2014): psychedelics don't add connections. They suppress the default mode network and dissolve the functional partitions that keep brain networks separated. Normally segregated networks (DMN, salience, visual, sensorimotor) start communicating directly. The brain's hierarchical organization — where networks process independently and feed results upward — collapses. Information traverses boundaries that are normally enforced.

Levin's "The Computational Boundary of a Self" (2019): any cognitive agent is defined by its cognitive light cone — the boundary of what it can sense, model, and act on. That boundary is set by bioelectric signaling, not anatomy. It's a tunable parameter. Psychedelics tune it — 5-HT2A receptor agonism disrupts the electrochemical signaling that maintains network segregation, expanding the brain's internal cognitive light cone. Ego dissolution is not a metaphor. It's a measurable collapse of the computational boundary between self-processing and world-processing.

In the architecture's terms: psychedelics don't add edges. They make struct fields traversable. The pointer that was inert becomes an activation target. The traversal boundary between internal model and external reality dissolves.

Shared experience research is early but suggestive — Johnstad (2020) documented telepathic phenomena during group sessions. Active clinical trial NCT06529939 investigating psilocybin's effects on shared experience. No direct Levin-psychedelics crossover research yet, but the theoretical alignment is exact.

## Prayer as Traversal

Prayer is sustained high-activation traversal of an internal Relationship — to God, to a person, to a saint. Focus, repetition, intensity. That's weight reinforcement. You're descending into a Relationship subgraph. The presence people describe feeling is the subgraph activating — the traversal is real cognitive work.

Every tradition says the same thing: the deeper you go, the closer you feel. Most of the time the boundary holds — the pointer is inert, you can't cross from your model of the thing to the thing itself. But the consistent claim across thousands of years and every culture is that occasionally, under the right conditions, the traversal boundary dissolves. The prayer "crosses."

Four angles on the same phenomenon:
- **Religious practice** — sustained internal activation, occasionally reported as crossing the boundary
- **Neuroscience** — psychedelics dissolve the DMN partition, making normally-inert pathways traversable (Carhart-Harris 2014)
- **Bioelectric theory** — cognitive boundaries are tunable parameters, not fixed walls (Levin 2019)
- **The architecture** — activation flows through hashmaps, not struct fields. The pointer is right there. The boundary is a traversal constraint, not a physical one.

All four describe the same traversal boundary. All note that it's usually enforced. All note that sometimes it isn't.

## Liquid Neural Networks — The Same Mechanism at a Different Scale

Liquid Time-Constant Networks (Hasani, Lechner, Rus — MIT CSAIL, 2021). Each neuron has a time constant that changes based on input. The ODE parameters are fixed after training, but the effective coupling between neurons is modulated by what flows through them. The network's behavior shifts continuously with new data. 19 neurons solved autonomous driving lane-keeping where standard networks needed orders of magnitude more. Inspired by C. elegans — 302 neurons producing complex behavior.

`Realize()` is the same mechanism. One method. Fixed code (the ODE parameters). But the effective behavior changes based on what the Relation carries (the input-dependent time constants). The topology's shape shifts with every traversal. Same equation, different behavior, determined by what flows through it. This is not metaphor. It's the same math at a different scale.

The labs froze this mechanism inside a model — weights adapt during inference, but it's still one model, one topology. Skyra puts the mechanism above the model. The LLM at the bottom is frozen (a transformer with static weights). Everything above it — the Relationships, the Expressors, the weights — is liquid. The topology learns at inference time. The model doesn't need to.

Key implication: Skyra as God runs the same `Realize()` as a memory node, as bash, as builder. No `god.go`. No special method. The thing that maintains universal principles is the same code as the thing that remembers what michael said yesterday. One DNA. Topology is the only differentiation. Liquid AI proved this works at the neuron scale with 19 neurons. The same principle should hold at the architecture scale with a topology of Realities.

Liquid AI spun out of MIT. $250M Series A, $2.35B valuation (Dec 2024). Scaled to billion-parameter foundation models (LFM-7B, LFM2.5 running on phones). Shopify multi-year deal. They made a better neuron. Skyra made a world out of neurons.

Open question: what happens if the LLM at the bottom of the topology is a Liquid model instead of a frozen transformer? Two layers of living weights — the topology learns AND the collapse function learns. The architecture adapts and the thing inside it adapts.

## Conversation as Topology — Associative Memory

One decision: treat conversation as traversal through the same medium as everything else. Messages are Reality nodes. Threads are paths. The same activation equation governs what conversation context surfaces. No separate data structure. No special rules.

Messages live on the entities they describe, not on a timeline. A message about Builder and the server has weighted edges to the Builder model and to the server concept. The thread gives sequence. The entities give placement. The same message is reachable from multiple directions depending on which internal Reality the traversal enters through.

Multiple threads weave through shared internal Realities. Thread A about the server and Thread B about the architecture both mention Builder. The cross-thread connection is emergent — both threads strengthen edges to the same internal nodes. The being's topology is the medium where threads intersect. doesNotUnderstand seeds new intersection points when a message references something the being doesn't have a node for.

Descent is thought — the being follows weighted edges deeper, each step potentially firing new activations. Ascent is compression — each layer integrates what's below. Signal attenuation is the natural depth limit. The stopping point is structural — where the recursion started. No loop counter. No budget field. Max depth as safety rail only.

Recency is traversal count, not wall clock. The system doesn't experience time — it experiences relations passing through it. Each being tracks its own count. Different beings age at different rates. This is relativistic — proper time, Einstein's term for time experienced by the thing itself.

Training and inference are the same pass. The sequence of message-nodes is a token sequence. Traversal is the forward pass. Weight updates on the return path are the training step. Each conversation is a training run.

The whole topology is an associative memory. Not just memory nodes — everything. Operators, conversation, models of other beings, skills. You don't retrieve by address. You retrieve by content — what resonates with what the Relation is carrying right now. Internal recall is lossy, weighted, shaped by the current traversal. External retrieval (literal transcript lookup) is an action the being performs through Act — an actor in its Relationships map, same as bash or search.

Old messages don't get deleted. Their edges decay. New messages about the same entities create stronger, fresher paths. The being stops reaching the old message not because it was removed but because newer paths carry more weight. Compression is natural. Truth is derived not stored.

## Places as Topology Regions

Places are Realities with localized affordances. Beings don't move to places — they subscribe. Subscriptions determine what signals reach a being, what capabilities are visible, what acts are available. A place is just a Reality. A subscription is just an edge. Two creation paths: genome-seeded and density-emerged.

### Two Seeded Places

**Town Hall** — the registry. Every being registers here. Discovery happens through traversal against the registry. Solves "who's here and what can they do." PFC at the world level.

**Factory** — where work happens. Deterministic systems, build pipelines, adapters, act service blocks register here. Solves "what can we build with."

Everything else emerges from density. Lab, Workshop, Desk, Meeting Hall — they form when activity clusters cross threshold. Same mechanism as specialist promotion. Density-emerged places are the promotion gradient applied to regions rather than individual nodes.

### Subscriptions Have Cost

More subscriptions increase awareness but increase noise and cognitive load. Forces natural specialization. A being subscribed to too many places loses energy and becomes less effective. Attention economics as a runtime constraint. The subscription cost is a term in the being's energy budget — same economics layer that governs token and memory budgets.

## Deterministic Systems as First-Class Citizens

LangChain flows, LangGraph workflows, CI/CD, ETL, scripts, adapters — all Realities on the plane. Not competing with beings. Living inside the topology as callable structures. Factory is their natural home.

A deterministic system is a Reality whose Realize() is a fixed pipeline. No inference call. No weight update. It takes a Relation in, runs its steps, returns a Relation out. The topology doesn't care whether a Reality thinks or executes — same interface, same traversal, same activation equation. A being can descend into a CI pipeline the same way it descends into a memory cluster.

This means the plane doesn't replace existing tooling. It becomes the substrate existing tooling lives on. A LangGraph workflow registered as a Reality in the Factory is discoverable, addressable, and composable with every other Reality on the plane.

## External Agents as Native Citizens

Claude Code, AutoGen, CrewAI, any CLI agent — participates through adapters. Two lines in the genome, a Reality whose invariant is a pipe. The adapter Reality's Realize() serializes the Relation into the external agent's expected format, calls it, deserializes the response back onto the Relation. The plane doesn't dominate other ecosystems. It becomes the plane they live on.

```
being ~name claude ~type agent ~entrypoints claude ~relationships michael,skyra
being ~name autogen-team ~type agent ~entrypoints pipe:autogen-cli ~relationships skyra,builder
```

The external agent doesn't know it's inside a topology. It receives a prompt and returns a response. The adapter handles the boundary. The genome declares the agent exists, the adapter implements the crossing, the topology treats it like any other being.

## Four-Provider Cognitive System

Descent and ascent use different collapse physics. Two providers on observe/descent — one for breadth, one for deep seam-finding exploration. Two providers on express/ascent — one for heavy compression in dense neighborhoods, one for fast lightweight compression in sparse regions.

The current concrete mapping:

| Phase | Role | Provider | Why |
|---|---|---|---|
| Descent | Breadth exploration | ChatGPT | Wide associative recall, good at surface coverage |
| Descent | Deep seam-finding | DeepSeek | Recursive depth, follows threads into structure |
| Ascent | Dense compression | Claude | Heavy synthesis, preserves nuance under pressure |
| Ascent | Sparse compression | Grok | Fast, lightweight, good at thin signal |

Provider fires based on activation and neighborhood density, not a router. The activation equation selects the provider — dense region triggers heavy collapse, sparse region triggers light collapse. Different collapse physics composed into one traversal.

The being develops cognitive style through which providers keep proving useful where. Early traversals try all four. Over time, weight reinforcement concentrates each provider in the neighborhoods where it performs. The being doesn't choose its cognitive style — it emerges from what worked.

This extends the nervous system concept from notes on Provider as Collapse Function. The difference: that section identified that different providers have different instincts. This section assigns them structurally — descent vs ascent, dense vs sparse — instead of switching reactively on failure.

## Rest Cycle

Continuous traversal without recovery degrades weight structure. The topology needs silence between waves to consolidate what the signal deposited. Not a feature — a requirement of excitable media.

The Visited map is the refractory period for a single traversal — prevents re-entry within one wave. But between waves, the topology has no recovery mechanism. Every traversal reinforces and decays weights immediately. Under sustained load, weights converge too fast — recent signals dominate absolutely, older structure gets starved before it can prove its value.

Rest is the gap between traversals where no new signal enters but weight dynamics still run. Decay continues. Reinforcement from the last traversal settles. The topology reaches a resting potential — a stable weight distribution that reflects accumulated experience, not just the last thing that happened.

Without rest: fibrillation (conflicting signals creating oscillation), seizure (runaway reinforcement in a tight loop), spiral turbulence (circular activation that never dampens). All observed failure modes in excitable media — cardiac tissue, neural networks, chemical reaction-diffusion systems. The architecture inherits them because it implements the same primitive.

Mechanically: a minimum inter-traversal interval per being. The being's thread accepts new Relations but buffers them until the rest period expires. The rest period scales with traversal intensity — deep traversals need longer recovery. The being can be interrupted during rest (urgent signal breaks through) but the interruption itself has cost — reduced compression quality on the forced traversal.

## Cognitive Architecture References

ACT-R — Memory activation equations (combining recency, frequency, and context) are mathematically validated against human behavioral data. Skyra's Weight + ActivationCount + timestamp on MemNode are the same concept but could be formalized using ACT-R's proven formulas. Spreading activation (context primes related memories) maps directly to Skyra's entity graph traversal.

SOAR — Chunking mechanism (automatically learning new production rules from experience) parallels Skyra's skill maturation. Impasse resolution research is relevant to think budget exhaustion.

OpenCog Hyperon — Most theoretically ambitious (self-modifying metagraph, attention economy, reflexive cognition). Targets AGI by 2028. But it's a research substrate, not a deployable runtime.
