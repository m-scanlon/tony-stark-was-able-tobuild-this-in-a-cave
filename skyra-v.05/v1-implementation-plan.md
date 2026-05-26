# v.1 Implementation Plan — Observe and Express

## What This Is

A step-by-step plan to replace ~1,900 lines of hardcoded routing logic with one recursive traversal in two phases. Each step builds on the previous one.

## The Core Insight

The graph already exists. Every Reality already has `map[string]Reality` hashmaps pointing to other Realities. These are already edges. The topology is already there. It's just unweighted and single-phase.

The entire move from v.05 to v.1 is: split Realize into Observe and Express, add weights, and let recursion do what it already does.

Every recursive function has two phases — work before the recursive call (observation) and work after it returns (expression). v.05 mixes both into a single `Realize()` and uses ~1,900 lines of handwritten routing to compensate. v.1 makes the phases explicit.

### The Interface

```go
type Reality interface {
    ID() string
    Create(r *Relation) Reality
    Realize(r *Relation) string
    Observe(r *Relation)
    Express(r *Relation) string
}
```

`Realize` is the recursive call. `Observe` and `Express` are the phases. Every Reality implements all five. The pattern:

```go
func (x *SomeReality) Realize(rel *Relation) string {
    x.Observe(rel)
    // ... recursive traversal ...
    return x.Express(rel)
}
```

Having all three on the interface means a Reality can't skip the phases. You might have an empty `Observe` or a trivial `Express`, but you had to think about it.

Observe traverses the Relationships map. Each entity above threshold gets `Observe` called — context accumulates on the Relation. This continues recursively until weight exhausts or no more Relationships exist. That's the bottom.

At the bottom, `Express` fires. Express calls `Realize` on Expressors. Each Expressor is a full Reality — it observes through its own Relationships (which affect and configure the Expressor itself), exhausts, then expresses through its own Expressors. This repeats until a Reality has no Expressors. That's the durable thing. It executes.

`Realize` is the recursive call that contains both phases. `Observe` and `Express` are the phases themselves. The traversal calls `Observe` on entities discovered through Relationships maps. It calls `Realize` on entities discovered through Expressors maps — because a Expressor needs its own full observe/express cycle. Execution is always at the termination of the recursion.

### Reality Carries Its Own Topology

There is no separate Relationship struct. The relationship fields dissolve into Reality itself. Every Reality struct carries:

```go
Weight        float64                // global weight — Skyra's intrinsic relationship to this Reality
Usage         int
LastUsed      time.Time
Relationships map[string]Reality     // observe — context accumulation
Expressors    map[string]Reality     // express — execution
```

Every Reality is both a thing and its connections. Self has Relationships (being, memory, context, desk) and Expressors (Think, Act). Think has Relationships (operators — what it knows it can do) and Expressors (Provider — the durable thing). Provider has no Expressors. It executes.

### Two Weights — Global and Local

There are two kinds of weight in the system. They serve different purposes and live in different places.

**Global weight** lives on Base. It's Skyra's intrinsic relationship to this Reality — how much bash matters to Skyra overall, accumulated through all use across all contexts. Each being is its own god. The global weight is the being's perspective on the importance of this Reality.

**Local weight** lives on the edge — the Reality that *is* the connection between two other Realities. Bash's local weight from server-memory is different from bash's local weight from poetry. Same bash. Different connection strengths. The local weight is how strongly *this specific path* flows from one Reality to another.

The edge is a Reality. It implements Reality. It embeds Base. On Observe, it reads its local weight and contributes to the activation decision. On Express — the return path — it updates its local weight based on what just happened. Reinforcement if the traversal produced value. Decay if it didn't.

The target Reality lives in the edge's own Relationships map. The traversal passes through the edge to reach the thing.

```
Self.Relationships["bash"] → Edge Reality (local weight, usage, recency)
    Edge.Observe → reads local weight, activation check
    Edge.Relationships["target"] → Bash Reality (global weight)
        Bash.Observe → ...
        Bash.Express → ...
    Edge.Express → updates local weight on return
```

This dissolves the node/edge distinction entirely. There is no graph with nodes and edges. There is one type: Reality. The edge between two neurons is itself a neuron. The Relationships map stays `map[string]Reality`. No type change. The connection between two Realities is itself a Reality containing Realities.

The activation formula uses both weights:

```
activation_i = global_weight * local_weight * relevance * recency * trust * context_fit
```

Global weight says how much this Reality matters to the being overall. Local weight says how strongly this specific connection flows. Both multiply into the activation score. A Reality can be globally important but locally irrelevant to the current traversal path. A Reality can be locally strong from this specific context but globally weak. The product captures both.

### The Activation Variables

```
activation_i = global_weight * local_weight * relevance * recency * trust * context_fit
```

Each variable and where it comes from:

- **global_weight** — the target Reality's Weight on Base. The being's intrinsic relationship to this Reality across all contexts. Skyra's relationship to bash overall. Lives on the target Reality. Updated by cumulative usage across all traversal paths.

- **local_weight** — the edge Reality's Weight on Base. The strength of this specific connection between two Realities. How often traversal flows from A to B. Lives on the edge Reality. Updated on the return path of each traversal through this edge.

- **relevance** — the overlap integral. How much of the Relation's current content exists in the target Reality's content. Pure content match — does the substance of what the Relation is carrying overlap with the substance of what this Reality holds? This is the magnitude of the inner product ⟨Relation|Reality⟩.

- **recency** — how recently this edge was traversed. Power law decay (ACT-R). Time-dependent. A recently traversed edge is warm. An untouched edge cools. The decay curve is a power law, not exponential — validated against human behavioral data.

- **trust** — the being's trust in this Reality. Position on the persuadability spectrum (Levin). How much this Reality's output reshapes the being's state. Bidirectional — the being trusts the Reality, and the Reality's track record earns or erodes that trust.

- **context_fit** — the phase alignment. Not "does this match what I'm carrying" (that's relevance) but "given each thing's own history of evolution, are they in phase right now at this moment of encounter." In QM, every state evolves continuously — phase rotates based on energy and history. Two states that were in phase yesterday may be out of phase today because they've been rotating at different rates. Context_fit captures whether the Relation's trajectory and the target Reality's current state are aligned at the moment of encounter.

### Relevance and Context_fit — One Equation, Two Behaviors

Relevance and context_fit are the magnitude and phase of a complex amplitude:

```
amplitude_i = magnitude_i * e^(phase_i)
```

Magnitude is relevance — content overlap. Phase is context_fit — directional alignment. They're the two components of the same number. QM keeps them separate because they do different work during collapse. Two edges with equal relevance but different phase alignment produce different outcomes — constructive interference when aligned, destructive interference when misaligned.

For **actors** (deterministic, non-cognitive Realities), relevance and context_fit produce the same value. An actor doesn't evolve between traversals. Bash is bash. Its state doesn't rotate. There's no independent phase evolution, so there's no divergence between content overlap and directional alignment. The two terms collapse to the same number — not because the formula simplifies, but because the target doesn't move.

For **agents** (cognitive Realities), relevance and context_fit can diverge. An agent is live — it thinks, acts, updates its own graph between encounters. Its phase has rotated. A memory of Builder from yesterday has content overlap (relevance is high) but Builder's state has shifted since then (context_fit may be low). The divergence between relevance and context_fit for agents is an **attention function** — the mechanism by which an evolving system selects what to attend to from its field of potential, given both content overlap and directional alignment.

The equation doesn't change. The same six terms compute for every Reality. The actor/agent distinction is emergent from the physics of the target, not from the formula. One system.

### Relationships

**Relationships** activate during observation. The traversal enters, checks activations on each entry (the edge Reality reads its local weight, the target Reality contributes its global weight), and calls `Observe` on those above threshold. Context accumulates on the Relation as it passes through — parsers attach, state enriches. This is association, thinking, following weighted paths deeper.

**Expressors** activate during expression. When Relationships exhaust and Express fires, it calls `Realize` on Expressors by weight. Each Expressor runs a full cycle — observing through its own Relationships (which shape and configure the Expressor), then expressing through its own Expressors. An Expressor's Relationships affect the Expressor itself: Think's operators aren't things Think routes to, they're what Think knows it can do when it reaches the LLM. The recursion continues until a Reality has no Expressors — the durable thing. It executes.

Skyra's model of Builder is a distinct Reality instance living in her Relationships map — not the actual Builder. The real Builder is on the thread plane, reachable only through action. Skyra's model carries its own Weight, Relationships (memories, patterns, trust signals), and Expressors (how actions toward Builder surface). Thinking about Builder traverses the model. Speaking to Builder goes through Exchange to the real thing.

Both maps hold `Reality`. The recursion is natural — every entry can have its own Relationships and Expressors, all the way down.

### How Memories Work

A memory of a conversation with Builder and Philosopher is a Reality. It lives in the Relationships maps of both Skyra's Builder model and Skyra's Philosopher model — same Reality instance, different weights in different maps. The memory itself has Relationships pointing back to Builder and Philosopher as entities.

When Skyra thinks about Philosopher, she observes into that Reality's Relationships. She hits a memory. That memory's Relationships map has an entry pointing to Builder. If the weight is high enough, traversal follows it. She arrived at Builder through Philosopher without anyone routing her there. The path emerged from the topology.

### One Traversal

```
Observe: Relation enters. Activate Relationships by weight.
         Each fires Observe. Context accumulates.
         Sub-relationships activate. Deeper.
         Weight exhausts or no more Relationships.

         — bottom —

Express: Fires. Calls Realize on Expressors by weight.
         Each Expressor is a Reality. It observes through
         its own Relationships (they affect the Expressor).
         Those exhaust. Expressor expresses through its own Expressors.
         Recurse until no more Expressors.

         — durable thing — EXECUTE.

         Result propagates back up the call stack.
```

One recursive pass. Two phases. Observe accumulates context until weight exhausts. Express calls Realize on Expressors, each of which runs its own full observe/express cycle. Execution happens at the base case — a Reality with no Expressors. Everything above it is weight resolution. This is what recursion already does. The plan just makes it formal.

### Concrete Example

```
Self.Realize:
  Self.Observe → Relationships: being, memory, context, desk
    each accumulates context on the Relation
    sub-relationships observe until weight exhausts
  — exhausted —
  Self.Express → Realize on Expressors:

    Think.Realize:
      Think.Observe → Relationships: operators (bash, retrieve-context...)
        operators affect Think — they're what Think knows it can do
        observe until exhausted
      — exhausted —
      Think.Express → Realize on Expressors:

        Provider.Realize:
          Provider.Observe → Relationships: (model config)
          — exhausted —
          Provider.Express → no Expressors → EXECUTE (call LLM)
          ← result

      ← result back to Think
    ← result back to Self

    Act.Realize:
      Act.Observe → Relationships: (inner thought, peer context)
      — exhausted —
      Act.Express → Realize on Expressors:

        Provider.Realize:
          Provider.Express → EXECUTE (call LLM)
          ← result

      ← result back to Act
    ← result back to Self

  ← result propagates up
```

### Act's Position in the Architecture

Act is the outermost execution surface — the edge of cognition before results cross the pipe to the world. In v.1, Act's only Relationship is Think, so it can only go one direction: inward to thought. That's correct for synchronous single-turn operation.

Each integration is an actor — a Reality that participates in the topology but doesn't model, choose, or grow. Gmail, Slack, the calendar — same interface, same traversal, same pipe contract. Act routes to them directly, same way it routes to any other Reality on the thread. The pipe boundary and the act service sit behind the actor's Express — invisible to Act. Act just sees a peer in its Relationships map.

When Act moves to async, completed acts accumulate in Act's Relationships on a timer. Instead of firing a new inference call per completed action, the being batches — waits for results to arrive, lets them accumulate as context in Act's Relationships map, then fires one Think cycle with all of them present. The weight system governs when the batch is worth processing: enough accumulated weight from completed acts crosses threshold, Act observes them all at once, expresses once. Cost control through topology, not scheduling logic.

```
```

Relationships on a Expressor are not separate from the Expressor — they shape the Expressor before it fires. Think's operators accumulate on the Relation during Think's observation. The Provider at the bottom consumes all of it.

### What This Replaces

The ~1,900 lines of routing code — target peeling, protocol enforcement, self-route detection, retry loops, operator dispatch, access checks — is the runtime doing by hand what this one traversal does automatically. Specifically:

- **Exchange.Realize** (~100 lines) — target peeling, `isBeing()` type-switch, redirect logic, being lookup, Process special-casing
- **NewThread.Realize** (~110 lines) — the `for` loop (event loop), access checks, operator injection (`ThinkOps`/`ActOps`), error routing
- **Self.Realize** (~80 lines) — hardcoded `think → act → think-back` loop, hashmap-based reality assembly
- **Think.Realize** (~130 lines) — operator dispatch loop, tag parsing, outer-op blocking
- **Act.Realize** (~130 lines) — protocol violation retries (3 attempts), self-route detection, `ParseResponse` as router
- **Operators.Realize** (~30 lines) — verb extraction and constructor dispatch
- **main.go bootstrap** (~130 lines) — hardcoded per-being operator injection
- **Tag parsing helpers** (~200 lines) — `ParseResponse`, `parseOp`, `parseThink`, `parseThinkBack`, `isNoReply`, `isBeing`, `extractVerb`, `Peel`

The remaining ~1,000 lines are conversation management, ref handling, memory compression, context heating. That code survives — it's state, not routing.

### The Sheaf Structure

There is no separate graph data structure. No `CognitiveGraph`. No `GraphNode`. No separate memory system. Every Reality holds two maps of Reality. Every Reality found in those maps may itself have its own maps. The topology is recursive.

This is not a graph. A graph has two types (nodes and edges). Here there is one type: Reality. The closest formal analog is a sheaf on a topology — local data assigned to each region, with consistency across overlaps and no central authority assembling the global view. The Mandelbrot-Julia duality holds: `Realize` is the generating function, each being's internal topology is a different Julia set, and the thread plane indexes all of them.

---

## Steps

### Step 0: Split the interface, add relationship fields to every Reality

**Files:** `reality.go`, all Reality implementations

Split the Reality interface from 3 methods to 5:

```go
type Reality interface {
    ID() string
    Create(r *Relation) Reality
    Realize(r *Relation) string
    Observe(r *Relation)
    Express(r *Relation) string
}
```

Every existing Reality gets:
1. Stub `Observe` (empty) and `Express` (returns `""`)
2. New fields on the struct: `Weight float64`, `Usage int`, `LastUsed time.Time`, `Relationships map[string]Reality`, `Expressors map[string]Reality`

Existing `Realize` methods stay unchanged. The new methods and fields exist but nothing uses them yet.

No `relationship.go`. No Relationship struct. The fields live on every Reality directly.

**Verify:** Compile. Run. Everything works exactly as before. No behavior change.

**Risk:** None. Additive.

---

### Step 1: Extend Relation with traversal state

**Files:** `relation.go`

Add fields for the two-phase traversal:

```go
Depth     int
Budget    float64
Trace     []string
```

Additive. Existing fields stay. Zero-valued new fields change nothing. Sets up later steps.

**Verify:** Compile. Existing tests pass.

**Risk:** None.

---

### Step 2: Swap Self's hashmaps to Relationships

**Files:** `self.go`

Replace `Self.Realities map[string]Reality` with `Self.Relationships map[string]Reality`.

All initial lookups that did `s.Realities[name]` now do `s.Relationships[name]`. Pure mechanical rename. No Relationship structs injected yet — the map still holds concrete Realities directly.

This is the first file to move. Self is the center — being, memory, context, desk, think, act all hang off it. Moving Self's map first means the central node uses the new field name while leaves still use the old one.

**Verify:** Compile. Run. Identical behavior.

**Risk:** None. Mechanical rename.

---

### Step 3: Think operator dispatch → relationship traversal

**Files:** `think.go`

**Remove:** `parseOp()`, `collectOps()`, outer-op blocking (`isOuterOp`), `renderOps`/`renderOpsWithOuter`, the tag-dispatch loop.

**Replace:** Think's operators move into its Relationships map. Rename `Think.Operators` → `Think.Relationships`. Think is a Expressor on Self — when Self.Ascend calls Think.Realize, Think runs its own cycle:

`Think.Observe`: Sort Relationships by activation. Operators above threshold accumulate on the Relation as capabilities the being can use. Memory-type relationships attach as context. The LLM sees operators because the weights surfaced them, not because they're in a hardcoded map. Relationships that affect Think shape what Think knows it can do.

`Think.Express`: Call Realize on Think's Expressors (Provider). Provider is the durable thing — it calls the LLM with everything that accumulated during observation.

Inner/outer operator distinction goes away. One map. Weights decide reachability. Tag parsing for invocation stays temporarily — the being still says `<bash>command</bash>`.

**Fallback:** If no relationship activates above threshold, show all (identical to v.05).

**Verify:** Send impulse. Think's prompt includes operators from weighted relationships. Invocation still works.

**Risk:** Medium. Fallback makes it safe.

---

### Step 4: Act protocol enforcement → doesNotUnderstand

**Files:** `act.go`, `exchange.go`

**Remove from Act:**
- 3-retry loop for protocol violations
- Self-route detection and warning

**Replace:** Act is a Expressor on Self. When Self.Express calls Act.Realize, Act runs its own cycle:
- Act's `Observe` traverses its Relationships — inner thought, peer context, conversation state accumulate on the Relation. These affect Act, shaping what it knows before it speaks.
- Act's `Express` calls Realize on Act's Expressors (Provider). Provider is the durable thing — calls the LLM once, parses the response.
- Valid `<target>message</target>` tags → route.
- No valid tags → doesNotUnderstand: seed new relationship at minimum weight, inform the being.
- Self-route is structurally impossible — no relationship from self to self.

**Remove from Exchange:**
- `Peel()` for target extraction. Act sets `r.ID` directly.
- `isBeing()` type-switch. Traversability is whether the Relationship has a Reality, not a type assertion.

**Verify:** Protocol violation → doesNotUnderstand. Explicit tags still route. Self-address doesn't retry.

**Risk:** Behavior change. Correct per spec. Log doesNotUnderstand events clearly.

---

### Step 5: Operator injection → relationship seeding

**Files:** `newthread.go`, `self.go`, `main.go`

**Remove:**
- `NewThread.ThinkOps` and injection loops
- `NewThread.ActOps` and injection loops
- Per-being operator wiring in `main.go` bootstrap
- `Self.Create`'s hardcoded Think/Act assembly from Realities

**Replace:**
- Genome declares operators as relationships. Bootstrap seeds them on each being with initial weights.
- `Self.Create` reads from its own Relationships map for topology.
- NewThread injects beings as Relationships. Operators are already on the being.

**Verify:** Bootstrap. Each being has correct operator relationships. Builder's bash has higher weight.

**Risk:** Low-medium.

---

### Step 6: Self's Think-Act loop → Observe/Express traversal

**Files:** `self.go`

**Remove:** The `for { think; act; think-back }` loop. The heart of v.05.

**Replace:**

`Self.Observe`: Sort Relationships by activation. Call `Observe` on those above threshold. Being, memory, context, desk accumulate on the Relation as it passes through. Sub-relationships observe until weight exhausts.

`Self.Express`: Call `Realize` on Expressors by activation. Think and Act are Expressors on Self. Think.Realize runs a full cycle — observes through its own Relationships (operators affect Think, shaping what it knows it can do), exhausts, then expresses through its Expressors (Provider — the durable thing that calls the LLM). Act.Realize does the same — observes through its Relationships (inner thought, peer context), expresses through its Expressors (Provider — calls the LLM to produce speech).

Think-back: Act's result signals re-entry. Self calls Think.Realize again. Budget limits depth — not a hardcoded loop count, but weight/budget exhaustion on the Relation.

`Self.Realize` orchestrates: `Observe(rel)` then `Express(rel)`. The two-phase pattern.

Thread → Exchange → Self stays hardcoded (skeleton). Weighted traversal governs what happens inside Self.

**Feature flag:** If Relationships map is empty, fall back to v.05 loop. Migrate one being at a time.

**Verify:** Think fires (high weight). Act fires after. Think-back works (bounded by budget).

**Risk:** High. Feature flag is mitigation.

---

### Step 7: NewThread for-loop → relationship-driven re-entry

**Files:** `newthread.go`, `exchange.go`

**Remove:** NewThread's `for { ... }` infinite loop.

**Replace:**

1. Validate access, create thread (same).
2. Inject beings as Relationships (same).
3. `Exchange.Realize(r)` once. Records entry.
4. Exchange calls `being.Realize(r)` → Observe/Express from Step 6.
5. Relation returns up through Exchange → Thread.
6. `r.ID` set to new target → Thread calls `Exchange.Realize(r)` again. Explicit re-entry.
7. Stops when `r.ID` empty or `r.Budget` exhausted.

Exchange narrows to: conversation state, entries, compression, parsers. No target guessing. Exchange's `Observe` records the entry and attaches conversation context. Exchange's `Express` handles compression and cleanup.

**Verify:** Multi-hop (michael → skyra → builder → skyra → michael). Entries recorded. Terminates.

**Risk:** High. Feature flag — old loop alongside new.

---

### Step 8: doesNotUnderstand as growth

**Files:** `act.go`

When Act targets something not in the being's Relationships:

1. New Reality created in Relationships map: `Weight: 0.01`, stub.
2. Being receives doesNotUnderstand response.
3. Subsequent encounters increment weight.
4. Threshold crossing → participates in traversal.

The being grows by reaching into the unknown.

**Verify:** Nonexistent target → new Reality in Relationships. 5 attempts → weight increases. Crosses threshold → appears in context.

**Risk:** Low. Additive.

---

### Step 9: Weight updates from usage

**Files:** All Reality implementations with Relationships/Expressors

After each traversal:

- Traversed: `Usage++`, `Weight += reinforcement`
- Untraversed: `Weight *= decay`
- Power law decay (ACT-R): `decay = (time_since_last_use)^(-d)`

The graph learns. Frequent use → stronger. Neglect → fades. This applies to both the Relationships and Expressors maps — frequently used expressors strengthen, neglected ones fade.

**Verify:** 20 interactions using bash → weight up. Browse untouched → weight down. Traversal ordering changes.

**Risk:** Low-medium. Start conservative (0.99 decay).

---

### Step 10: Remove MemoryGraph — memory is already relationships

**Files:** `memgraph.go`, `memory.go`, `context.go`

**Remove:** `MemoryGraph`, `Entity`, `EntityEdge`, `MemNode`, `MemEdge`. `Context.Warm` cache. The entire separate memory data structure.

Memory was never a separate system. It was always relationships. The MemoryGraph was a parallel structure doing what Relationships maps already do. Skyra's memories of Builder live inside her Builder model's Relationships map — as Realities with their own Relationships maps pointing back to entities involved. When she needs context about Builder, she observes into that Reality. The memory's cross-pointers to Philosopher let her arrive at Builder through Philosopher without explicit routing.

**Replace:**
- `Memory.Store` → creates a Reality inside the relevant parent's Relationships map. A memory involving Builder and Philosopher creates entries in both with cross-pointers.
- `Memory.Query` → observation into a Reality's Relationships subgraph, activation-weighted.
- `Context.Heat` → same observation, scoped.

Migration for existing `graph.json` — entities become Realities with Weight/Relationships/Expressors fields, entity edges become nested Realities with cross-pointers.

**Verify:** Existing memories load into Relationships maps. Retrieval returns same results. Storage works.

**Risk:** High. Persistent state. Migrate on copy first. Keep old code behind build tag.

---

### Step 11: Clean up dead code

**Files:** All touched files.

**Remove:**
- `Think.Operators` / `Act.Operators` (now in Relationships map)
- `NewThread.ThinkOps` / `NewThread.ActOps`
- `Operators` struct and `operators.go`
- `ParseResponse` as router
- `isBeing()` type-switch
- `Peel()` for target extraction
- Feature flags
- Old `MemoryGraph` code
- Stub `Observe`/`Express` implementations that were never filled in

**Verify:** Test suite passes. `go vet` clean. No dead code.

**Risk:** Low. Cleanup.

---

## Decisions

### Locked

1. **Interface:** `ID`, `Create`, `Realize`, `Observe`, `Express`. Realize contains both phases. Traversal calls `Observe` on entities found through Relationships maps. Traversal calls `Realize` on entities found through Expressors maps (Expressors need their own full observe/express cycle).

2. **No separate Relationship struct. Edges are Realities.** Weight, Usage, LastUsed, Relationships, Expressors dissolve into every Reality. Every Reality is both a thing and its connections. The connection between two Realities is itself a Reality — it implements the interface, embeds Base, observes (reads local weight), expresses (updates local weight on return). The target lives in the edge Reality's own Relationships map. One type all the way down. No node/edge distinction.

3. **Two-phase traversal:** Observe activates Relationships maps, accumulates context until weight exhausts. Express calls Realize on Expressors. Each Expressor observes through its own Relationships (which affect/configure the Expressor), then expresses through its own Expressors. Recursion terminates at a Reality with no Expressors — the durable thing. Execution happens there. One recursive pass.

4. **Memories as cross-pointed Realities:** A memory involving multiple entities exists in multiple Relationships maps. Cross-pointers enable emergent traversal paths.

5. **Deterministic first:** `argmax` for collapse. Stochastic/temperature later.

6. **Skeleton stays:** Thread → Exchange → Self stays hardcoded for now.

7. **Activation formula:** Activation is computed **per-entry** in a Relationships or Expressors map. If a Reality has four entries in its Relationships map, four separate activations are computed — one for each entry. The traversal visits entries above threshold, skips entries below. The activation score belongs to the entry, not to the map.

   ```
   activation_i = global_weight_i * local_weight_i * relevance_i * recency_i * trust_i * context_fit_i
   ```

   Where `i` is one entry in the map. **Global weight** lives on the target Reality's Base — the being's intrinsic relationship to that Reality across all contexts. **Local weight** lives on the edge Reality — the strength of this specific connection between two Realities. Both are per-entry because each entry is its own edge Reality pointing to its own target Reality. Start with `global_weight * local_weight * recency`. Add factors as they prove necessary.

8. **Relation signature:** Keep `Realize(r *Relation) string` and `Express(r *Relation) string`. Change in a later pass if needed.

9. **Intent graph:** Separate follow-up. Weighted relationships first.

10. **Act as the outer layer:** Act is the outermost execution surface — the edge of cognition. v.1 starts with Think as Act's only Relationship (synchronous, one direction). Each integration (Gmail, Slack, calendar) is an actor — Act routes to it directly as a peer in its Relationships map. Async: completed acts accumulate in Act's Relationships on a timer, batch into one inference call when weight crosses threshold. Cost control through topology, not scheduling.

11. **Actors and agents:** The distinction between cognitive and non-cognitive Realities follows Hewitt (1973) and Wooldridge & Jennings (1995). An **agent** is a Reality that observes, models, and chooses — it has a Self, Relationships it models, Expressors it selects between. Skyra is an agent. Michael is an agent. An **actor** is a Reality that receives and computes — same interface, same traversal, deterministic output, no interiority. Gmail is an actor. This is not a type split. The interface is identical. The distinction is emergent — what makes a Reality cognitive is the topology behind it, not a different type. An actor that accumulates enough complexity (gains Relationships, starts modeling, gets wired with choice) crosses into agency through the promotion gradient. Not by relabeling. By growing the topology that constitutes it.

---

## Sequencing

Steps 0-1 are mechanical. Add Observe/Express stubs, add fields to every Reality. Zero risk.

Step 2 renames Self's hashmap. Zero risk.

Steps 3-5 replace dispatch logic with weight-based surfacing. Medium risk, fallbacks available.

Steps 6-7 replace the core loops with Observe/Express. High risk, feature flags.

Steps 8-9 are the payoff. The graph grows and learns.

Step 10 is unification. One topology. Memory dissolves into Relationships.

Step 11 is compression. Throw away what's left.
