# v.1 Implementation Plan — One Traversal

## What This Is

A step-by-step plan to build the v.1 runtime from the clean slate. The old v.05 routing code (~1,900 lines) is deleted. The interface and Relation are defined. Everything below is new construction.

## The Core Insight (Updated)

There is one recursive call: `Realize()`. It descends through Relationships (thought), ascends through Expressors (action formation). Providers fire on either direction when a node needs inference. One traversal does everything — storage, thought, compression, action, learning. The being visits each node once. Thoughts compound on the way down. Actions crystallize on the way up.

Think and Act are not separate systems. They are descriptions of what the traversal is doing at different depths. Deep is thought. The surface is action. There is no orchestrator, no loop, no handoff. The call stack is the control flow.

Conversation and memory are the same thing. Messages are Reality nodes attached to the entities they reference. Thread is its own Reality — the episodic binding that connects unrelated messages that co-occurred.

### The Interface (Locked — Already Implemented)

```go
type Reality interface {
    ID() string
    Core() *Base
    Create(r *Relation) Reality
    Realize(r *Relation) string
    Observe(r *Relation)
    Express(r *Relation) string
}
```

`Realize` is the recursive call. `Observe` fires on the descent — context accumulates. `Express` fires on the ascent — compression and action formation. Every Reality implements all of them.

### Three Maps on Base

```go
type Base struct {
    Weight          float64
    WeightSum       float64            // cumulative activation for EMA
    TraversalCount  int                // being's proper time — traversal count, not wall clock
    LastTraversed   int                // traversal count when last activated
    Relationships   map[string]Reality // descent — what the being knows
    Expressors      map[string]Reality // ascent — what the being can do
    Providers       map[string]Reality // either direction — inference surfaces
}
```

- **Relationships** activate during descent. Context, memory, associations. The being's experience.
- **Expressors** activate during ascent. Compression, action formation. The being's output.
- **Providers** activate on either phase. The being's ability to think. Fires when a node needs to call out — LLM, bash, API. Orthogonal to direction.

A being without Providers is an actor — purely reactive, no inference, passes content through based on weights. Add a Provider and it starts thinking. The actor/agent line is whether the Providers map has anything in it.

### The Activation Equation (v.1)

```
activation_i = global_weight * local_weight * recency * thread_alignment * thread_strength
```

Computed per-entry in any of the three maps. Five terms, all computable:

- **global_weight** — EMA: `α * activation_this_traversal + (1 - α) * previous`. Slow α. The being's intrinsic relationship to this Reality across all contexts. Starvation is built in — when activation is 0, the EMA decays by `(1 - α)`.
- **local_weight** — EMA: same formula, faster α. Strength of this specific connection. Absorbs relevance for v.1 (actors don't evolve, so historical strength ≈ current relevance).
- **recency** — Power law: `(traversals_since_last_use)^(-0.5)`. ACT-R validated. Long tail — old edges fade but never fully disappear. Measured in the being's proper time (traversal count), not wall clock. Different beings age at different rates. Relativistic.
- **thread_alignment** — Binary for v.1: 1 if the node is on the same thread as the Relation, 0 if not. The gate — are you on this thread at all.
- **thread_strength** — EMA of per-node hit counts within the thread. How central this node is to what the thread is doing right now. A node on the thread but not recently relevant gets `thread_alignment = 1.0` but `thread_strength → 0` — still reachable but quiet. A node the thread keeps hitting gets both — fully gated in and loud. This is what makes the cost curve invert: the thread's own weight history concentrates activation where it belongs, so traversals get cheaper as the thread gets longer.

Thread alignment and thread strength split what was one term into gate and amplitude. Alignment is binary — on or off. Strength is continuous — the thread's opinion about what matters right now.

Activation is a tensor contraction — five dimensions with independent temporal dynamics projected to a scalar. High activation = all dimensions aligned. Low = at least one suppressing.

### The Base Traversal Pattern

1. Relation enters a node. The being sees it — one frame, one look.
2. Think pass fires (if Provider activates). Thought attaches to the Relation.
3. Activation determines the next node in Relationships. Relation moves deeper.
4. At the next node, the being sees new content PLUS all thoughts from above. Thoughts compound.
5. A node can only be visited once per traversal. No loops. No revisiting.
6. Signal exhausts — no more activations fire. The Relation attaches where it stopped, with edges back to entities involved. That's storage.
7. Ascent begins. Expressors activate. Each layer compresses the thought stream. Action plans crystallize.
8. Providers fire on the ascent if compression needs inference.
9. Top frame — act sequence deploys.

One pass. Storage, thought, compression, action, learning. The call stack is the boundary.

### Tags Are Signal

Tags (`<builder>message</builder>`) are not routing. They're a frequency component on the Relation that modulates activation. Builder's nodes resonate with a `<builder>` tag — activation high. Other beings are transparent to it. No tag = pure activation-based routing.

### Thread As Reality

Thread is its own Reality — the hippocampal index. It binds messages that co-occurred regardless of entity overlap. Two retrieval paths to any message:

- **Semantic** — through entity edges (what the message is about)
- **Episodic** — through the thread Reality (what happened together)

### Message Placement

Messages find their own place through traversal. The signal propagates until it exhausts. Where it dies is where it lives. Depth encodes meaning — shallow messages attach near the surface, deep questions attach far down. Same pass as the response — the processing IS the storage.

### Conversation and Memory Are The Same Thing

A message is a Reality node with edges to entities it references. A memory is a Reality node with edges to entities it references. Same type. Same activation. Same retrieval. The only difference is recency — a message from the active thread is hot, a memory from last week is cold. Not a type distinction. A position on the decay curve.

---

## Steps

### Step 0: Update Base struct (DONE → needs update)

**Files:** `reality.go`

Add Providers map and EMA/traversal fields to Base:

```go
type Base struct {
    Weight          float64
    WeightSum       float64
    TraversalCount  int
    LastTraversed   int
    Relationships   map[string]Reality
    Expressors      map[string]Reality
    Providers       map[string]Reality
}
```

Update `Activation()` to compute the four-term equation. Add `UpdateWeight()` for EMA updates.

**Verify:** Compile.

**Risk:** None. Additive.

---

### Step 1: Update Relation for traversal

**Files:** `relation.go`

Add traversal state:

```go
Visited     map[string]bool  // visit-once constraint
ThreadID    string           // already exists — used for thread_alignment
Signal      float64          // signal strength — attenuates during propagation
Thoughts    []string         // accumulated thought stream from descent
MaxDepth    int              // safety rail only — not the mechanism
```

Remove or deprecate `Budget` — signal attenuation and visit-once are the depth limits. Keep `Depth` as a counter for diagnostics. Keep `Trace` for debugging.

Update `Impress()` to initialize `Visited`, `Signal: 1.0`, `MaxDepth` to a safe default.

**Verify:** Compile. Existing code unaffected.

**Risk:** None.

---

### Step 2: Implement activation equation

**Files:** `reality.go` (new functions on Base)

```go
func (b *Base) Activate(rel *Relation, target Reality) float64 {
    globalWeight  := target.Core().Weight
    localWeight   := b.Weight
    recency       := computeRecency(b.TraversalCount, b.LastTraversed)
    threadAlign   := computeThreadAlignment(rel, target)
    threadStr     := computeThreadStrength(rel, target)
    return globalWeight * localWeight * recency * threadAlign * threadStr
}
```

Implement:
- `computeRecency(current, last int) float64` — power law: `(current - last)^(-0.5)`. Handle `last == 0` (never traversed) as minimum recency, not zero.
- `computeThreadAlignment(rel *Relation, target Reality) float64` — binary: 1.0 if target is on same thread, 0.0 if not. Needs thread membership on the Reality somehow. Start with 1.0 for all (no thread filtering) and refine.
- `computeThreadStrength(rel *Relation, target Reality) float64` — EMA of how often the current thread has activated this node. Thread Reality maintains per-node hit counts. Returns the thread's local weight to this node. Defaults to 1.0 when no thread history exists (cold start — don't suppress anything until the thread has opinion).
- `updateEMA(current, activation, alpha float64) float64` — the EMA formula.

**Verify:** Unit tests. Known weights produce expected activation scores.

**Risk:** Low. Pure math.

---

### Step 3: Implement the traversal — Realize()

**Files:** new file `traverse.go` or on Base directly

This is the core. One function that every Reality uses:

```go
func (b *Base) Realize(self Reality, rel *Relation) string {
    // Visit-once check
    if rel.Visited[self.ID()] {
        return ""
    }
    rel.Visited[self.ID()] = true

    // === DESCENT — Observe through Relationships ===
    self.Observe(rel)

    // Think pass — fire Providers if activated
    for _, p := range activatedEntries(b.Providers, rel) {
        thought := p.Realize(rel)
        rel.Thoughts = append(rel.Thoughts, thought)
    }

    // Descend into activated Relationships
    for _, r := range activatedEntries(b.Relationships, rel) {
        r.Realize(rel)
    }

    // === ASCENT — Express through Expressors ===
    // Expressors see accumulated thoughts, produce action
    for _, e := range activatedEntries(b.Expressors, rel) {
        // Provider may fire again during compression
        e.Realize(rel)
    }

    result := self.Express(rel)

    // === WEIGHT UPDATES ===
    b.TraversalCount++
    b.LastTraversed = b.TraversalCount
    // EMA update on return path
    b.Weight = updateEMA(b.Weight, 1.0, b.alpha())

    return result
}
```

`activatedEntries` sorts entries by activation score, returns those above threshold.

This is the DNA. Every Reality runs this. The topology determines the behavior.

**Verify:** Simple topology — three Realities in a chain. Relation enters, descends, ascends. Visit-once holds. Thoughts accumulate. Express fires on return.

**Risk:** High. This is the core mechanism. Get this right and everything else follows.

---

### Step 4: Implement concrete Realities — Being, Thread, Provider

**Files:** `being.go`, `thread.go`, `provider.go`

**Being** — the Self. Implements Reality. Its Relationships hold: models of other beings, memory nodes, message nodes. Its Expressors hold: action surfaces (response formation, tag emission). Its Providers hold: LLM providers (DeepSeek, Claude, etc).

```go
type Being struct {
    Base
    Name     string
    Identity string
    Purpose  string
    Type     string // "llm", "user", "agent"
}
```

Being's `Observe` attaches identity and purpose to the Relation as context.
Being's `Express` returns the compressed result from the ascent.

**Thread** — episodic binding. Implements Reality. Its Relationships hold: message nodes that co-occurred. Its Expressors: none (threads don't act). Its Providers: none.

```go
type Thread struct {
    Base
    ThreadID  string
    Members   map[string]Reality // beings participating
    Active    bool
}
```

Thread's `Observe` attaches thread context — who's here, what's been said.
Thread's `Express` handles re-entry — if the result tags a new target, the thread routes.

**Provider** — inference surface. Implements Reality. Terminal node — has no Relationships, no Expressors, no Providers of its own. It executes.

```go
type Provider struct {
    Base
    Model    string
    Endpoint string
    Call     func(prompt string) string // wraps inference layer
}
```

Provider's `Observe` assembles the prompt from accumulated context on the Relation.
Provider's `Express` calls the LLM and returns the result.

**Verify:** Create a Being with one Provider. Send a Relation. Being observes (identity context), Provider fires (LLM call), result returns. End to end.

**Risk:** Medium. First real integration with inference layer.

---

### Step 5: Message and memory as Reality nodes

**Files:** `message.go`

A message is a Reality. Same struct for both messages and memories — they're the same thing at different recency.

```go
type Message struct {
    Base
    Content   string
    From      string
    ThreadRef string   // which thread this belongs to
    Entities  []string // which entities this references
}
```

Message's `Observe` attaches its content to the Relation.
Message's `Express` returns empty — messages are context, not action.

**Placement:** When a traversal creates a new message, it attaches as a Reality in the Relationships maps of the entities it references AND in the Thread's Relationships map.

**Verify:** Send a message about Builder. Message node appears in Builder model's Relationships AND in the thread's Relationships. Traversal through Builder finds the message. Traversal through the thread finds the message.

**Risk:** Low. Additive.

---

### Step 6: Weight updates on the return path

**Files:** `reality.go`, `traverse.go`

On the return path of every traversal:

- **Traversed edges:** EMA update upward. `Weight = α * 1.0 + (1 - α) * Weight`. `LastTraversed = current traversal count`.
- **Untraversed edges:** EMA decays naturally. `Weight = α * 0.0 + (1 - α) * Weight` = `(1 - α) * Weight`. No separate decay pass — the being's traversal count increments, recency drops via power law.

Global weight (on target) and local weight (on edge) both update via EMA with different α values.

The being's `TraversalCount` increments once per traversal through it. This is the being's proper time.

**Verify:** 20 traversals using bash → bash weight up. Browse untouched → browse weight drifting down. Traversal ordering changes over time.

**Risk:** Low-medium. Start with conservative α values (0.1 global, 0.3 local).

---

### Step 7: doesNotUnderstand as growth

**Files:** `being.go`

When a being's Expressor produces a tag targeting something not in its Relationships:

1. New Reality seeded in Relationships: minimum weight, stub.
2. Being receives doesNotUnderstand — "I don't know this yet."
3. Subsequent encounters reinforce the edge via EMA.
4. Weight crosses threshold → participates in traversal.
5. The being grew a new connection by reaching into the unknown.

**Verify:** Being targets nonexistent entity → stub created. Five traversals → weight increases. Crosses threshold → appears in context.

**Risk:** Low. Additive.

---

### Step 8: Tags as signal modulation

**Files:** `relation.go`, `traverse.go`

Parse tags from the Relation's impulse. Tags become frequency components on the Relation that modulate activation:

- Tag matches a being's name → activation boost for that being's topology
- Tag doesn't match → activation dampened (not zero — still reachable if other factors are strong enough)
- No tag → pure activation-based routing. Whoever resonates most responds.

Tags are parsed once at entry. The frequencies live on the Relation for the duration of the traversal.

**Verify:** `<builder>check the server</builder>` → Builder's topology activates strongly. Skyra's dampened but not silent. No tag → whoever has the strongest activation responds.

**Risk:** Low. Modulates existing activation, doesn't replace it.

---

### Step 9: Genome bootstrap

**Files:** `main.go`, `genome.go`

The genome declares what could be. Bootstrap reads it and seeds topology:

1. Parse `genome.skyra` — beings, devices, components, relationships, initial weights.
2. Create Being Realities for each declared being with identity, purpose, type.
3. Seed Relationships between beings at declared initial weights.
4. Seed Providers on each being (which LLM, which endpoint).
5. Seed Expressors on each being (action surfaces — response, tag emission).
6. Create Thread Reality for the first thread.
7. Wire devices and components as actor Realities in the topology.

Skyra boots first. She shapes the topology — active judgment about what needs to be real for this world. Other beings boot into a space that already has curvature.

**Verify:** `genome.skyra` → running world. Each being has correct topology. Providers wired to inference layer. Send an impulse → full traversal → response.

**Risk:** Medium. Integration of everything above.

---

### Step 10: WebSocket device

**Files:** `ws.go`

The Universe observes itself. The frontend receives that observation.

1. WebSocket server on port 8080.
2. First-message auth.
3. On connect: snapshot — serialize the universe state (beings with weights, threads, messages, topology).
4. On traversal complete: delta events — entry, weight, being, topology changes.
5. Client → server: impulse messages that become Relations.

The serialization walks the topology and produces the JSON defined in `frontend-spec-v1.md`. The universe state shape is locked.

**Verify:** Connect via WebSocket. Receive snapshot. Send impulse. Receive deltas as the traversal runs. Reconnect → fresh snapshot.

**Risk:** Medium. Serialization of recursive topology needs cycle detection.

---

### Step 11: End-to-end integration

**Files:** All.

Full loop:

1. Genome boots the world.
2. Skyra activates, shapes initial topology.
3. Michael sends impulse via terminal or WebSocket.
4. Relation enters, descends through Skyra's topology.
5. Think passes fire at nodes (Providers activate).
6. Thoughts compound. Message attaches at exhaustion point.
7. Ascent compresses. Expressors form action.
8. Provider fires for final response generation.
9. Result returns to top frame. Tags parsed. Thread routes if needed.
10. Weights update on return path. Topology shifts.
11. Frontend receives deltas. Topology visualization updates.
12. Multi-hop: Skyra → Builder → Skyra → Michael. Thread manages re-entry.

**Verify:** Multi-party conversation works. Weights shift visibly over 20+ exchanges. Topology evolves. Frontend shows it happening.

**Risk:** High. Full integration. But every piece was verified independently in prior steps.

---

## Decisions

### Locked

1. **Interface:** `ID`, `Core`, `Create`, `Realize`, `Observe`, `Express`. Realize contains both phases. Locked and implemented.

2. **Three maps on Base.** Relationships (descent), Expressors (ascent), Providers (either direction). Same type: `map[string]Reality`. Same activation equation. Three roles.

3. **One traversal.** Descent through Relationships accumulates context with think passes. Ascent through Expressors compresses and forms action. Providers fire on either phase. No orchestrator. The call stack is the control flow.

4. **Think and Act dissolved.** They are not Expressors. They are descriptions of depth. Deep is thought. Surface is action. The traversal doesn't hand off between them. One pass.

5. **Conversation and memory are the same thing.** Messages are Reality nodes on entities. Recency is the only difference. No separate Exchange, no separate Memory system.

6. **Thread is a Reality.** Episodic binding. The hippocampal index. Binds co-occurring messages regardless of entity overlap. Two retrieval paths: semantic (entities) and episodic (thread).

7. **Message placement through traversal.** Where the signal exhausts is where the message attaches. Depth encodes meaning. The processing IS the storage.

8. **Visit-once per traversal.** A node can only be seen once per traversal. No loops. No revisiting. This is the finiteness constraint — no budget needed.

9. **Tags are signal.** Not routing. Frequency components on the Relation that modulate activation. No tag = pure weight-based routing.

10. **Activation equation (v.1):**
    ```
    activation_i = global_weight * local_weight * recency * thread_alignment * thread_strength
    ```
    Global and local weight: EMA with different α. Recency: power law. Thread alignment: binary gate. Thread strength: EMA of per-node hit counts within the thread — the thread's opinion about what matters right now. Tensor contraction — five dimensions to scalar.

11. **Global weight as EMA.** `α * activation + (1 - α) * previous`. α is per-being — controls temporal sensitivity. Slow α = stable. Fast α = reactive. Starvation built into the formula.

12. **Recency is traversal count.** Not wall clock. The system doesn't experience time. Each being has its own proper time. Different beings age at different rates. Relativistic.

13. **Deterministic first.** `argmax` for collapse. Stochastic/temperature later.

14. **Actors and agents.** Same interface. Distinction is emergent. Actor = no Providers (purely reactive). Agent = has Providers (can think). The line is whether the Providers map has entries.

15. **doesNotUnderstand as growth.** Reaching for something unknown seeds a new Reality at minimum weight. The being grows at the edges.

16. **Max depth as safety rail.** Not the mechanism. Signal attenuation and visit-once are the mechanism. Max depth catches runaway cases the mechanism misses.

### Deferred

- **Relevance** — content overlap. Absorbed by local_weight EMA for v.1. Own term when agents start evolving.
- **Trust** — coupling strength. Not computable yet. Levin persuadability spectrum. Where it lives structurally TBD.
- **Context_fit** — phase alignment beyond thread. Thread is the first component. Continuous phase comes when we understand what "rotation rate" means for a Reality.
- **Signal attenuation** — replacing Budget with physical depth. The Relation's signal strength decreasing per-step. v.1 uses visit-once + max depth. Signal attenuation is the v.2 mechanism.
- **Section weights on the Relation** — each component should have its own amplitude. Not in v.1. The Relation carries flat content for now.
- **Streaming** — descent and ascent overlapping temporally. v.1 is synchronous. Streaming is the next frontier.
- **Cognitive nervous system** — provider swap as circuit breaker for recursive patterns. Falls out of the Providers map naturally. Later.

---

## Sequencing

Step 0-1 are struct updates. Zero risk.

Step 2 is math. Low risk. Unit testable.

Step 3 is the core — one recursive traversal through three maps. High risk. Get this right and everything follows.

Step 4 builds concrete types on the traversal. Medium risk. First inference integration.

Step 5 makes messages and memory work. Low risk. Additive.

Step 6 makes the topology learn. Low-medium risk. Conservative α values.

Steps 7-8 are growth and signal modulation. Low risk. Additive.

Step 9 is bootstrap — genome becomes running world. Medium risk. Integration.

Step 10 is the frontend connection. Medium risk. Serialization.

Step 11 is the proof — full loop, multi-party, visible topology evolution.

---

## What This Replaces

The previous implementation plan (11 steps organized around replacing v.05 routing code) is superseded. Key differences:

- **Think and Act as Expressors on Self** — dead. Dissolved into depth.
- **Self's think-act loop replaced with weighted version** — dead. One traversal, no loop.
- **Exchange as a concept** — dead. Messages are Reality nodes.
- **Operator dispatch → relationship traversal** — absorbed. Operators are Realities in the topology. No dispatch needed.
- **Budget field** — deprecated. Visit-once + max depth + signal attenuation.
- **Two maps** — now three. Providers are orthogonal to direction.
- **Separate memory system removal** — unnecessary. Memory was never built as separate in v.1. It was always the topology.
