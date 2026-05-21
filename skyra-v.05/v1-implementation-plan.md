# v.1 Implementation Plan — Relationships Replace Routing

## What This Is

A step-by-step plan to replace ~1,900 lines of hardcoded routing logic with weighted relationships. Each step builds on the previous one.

## The Core Insight

The graph already exists. Every Reality already has `map[string]Reality` hashmaps pointing to other Realities. These are already edges. The topology is already there. It's just unweighted.

The entire move from v.05 to v.1 is: add a weight to each hashmap entry.

```go
// v.05 — unweighted
map[string]Reality

// v.1 — weighted
map[string]*Relationship

type Relationship struct {
    Reality        Reality
    Weight         float64
    Usage          int
    LastUsed       time.Time
    Relationships  map[string]*Relationship
}
```

Relationship implements Reality. It's not just a weighted pointer — it's a full node with its own hashmap of relationships inside it. Skyra's Relationship to Builder is her model of Builder: memories of him, patterns she's noticed, trust level, interaction history all live as relationships inside it. Activation isn't a number someone set — it's emergent from the contents of the Relationship.

The `Reality` field is a reference to the actual entity — the real Builder, not a copy. It's there because the Relationship needs to know which real thing it's about, and because Builder's full state can influence Skyra's weights. But the reference is not traversable. Skyra cannot reach Builder by descending into her Relationship to him — that's thinking about him, not reaching him. To actually reach Builder, she acts. Act routes through the thread plane — the shared level where beings exist as traversable relationships. The only way to reach a being is to act. Relationships live inside the being, not between beings.

There is no separate graph data structure. No `CognitiveGraph`. No `GraphNode`. No separate memory system. Every Reality contains a `map[string]*Relationship`. Every Relationship is itself a Reality with its own `map[string]*Relationship`. The topology is recursive — each node holds its own subgraph. This is not a graph. A graph has two types (nodes and edges). Here there is one type: Reality. The closest formal analog is a sheaf on a topology — local data assigned to each region, with consistency across overlaps and no central authority assembling the global view.

Relationship's `Realize()` is the only routing logic in the system:

1. Relation comes in.
2. Check each relationship in the hashmap — weight against the Relation's current state.
3. All above threshold → fire `Realize()` on all of them concurrently. Wait for all to return.
4. Return.

Every Relationship does the same thing. The difference is what's in the hashmap. Deep inside Skyra's model of Builder, the hashmap holds memories and patterns. On the thread plane, the hashmap holds actual beings. Same `Realize()`. Different contents. Different depth. Base-level entities — Think, Provider, Terminal — are where the chain terminates. Their `Realize()` does real work (calls the LLM, runs bash, writes to disk) instead of weight routing. Everything above them is one method: route by weight, call realize, return.

The ~1,900 lines of routing code — target peeling, protocol enforcement, self-route detection, retry loops, operator dispatch, access checks — is the runtime doing by hand what this one method does automatically. Replace `map[string]Reality` with `map[string]*Relationship`, and the routing code dissolves.

## What's Being Replaced

- **Exchange.Realize** (~100 lines) — target peeling from impulse text, `isBeing()` type-switch, redirect logic, being lookup, Process special-casing
- **NewThread.Realize** (~110 lines) — the `for` loop (event loop), access checks, operator injection (`ThinkOps`/`ActOps`), error routing
- **Self.Realize** (~80 lines) — hardcoded `think → act → think-back` loop, hashmap-based reality assembly
- **Think.Realize** (~130 lines) — operator dispatch loop, tag parsing, outer-op blocking
- **Act.Realize** (~130 lines) — protocol violation retries (3 attempts), self-route detection, `ParseResponse` as router
- **Operators.Realize** (~30 lines) — verb extraction and constructor dispatch
- **main.go bootstrap** (~130 lines) — hardcoded per-being operator injection
- **Tag parsing helpers** (~200 lines) — `ParseResponse`, `parseOp`, `parseThink`, `parseThinkBack`, `isNoReply`, `isBeing`, `extractVerb`, `Peel`

The remaining ~1,000 lines are conversation management, ref handling, memory compression, context heating. That code survives — it's state, not routing.

## What Replaces It

```go
func (r *SomeReality) traverse(rel *Relation) {
    for _, edge := range r.Relationships {
        activation := edge.Weight * rel.Relevance(edge.Reality) * edge.Recency()
        if activation < rel.Threshold() {
            continue
        }
        edge.Reality.Realize(rel)
    }
}
```

One formula. Ten lines. Routing becomes physics.

---

## Steps

### Step 0: Define Relationship and swap the hashmaps

**Files:** new `src/reality/relationship.go`, `self.go`, `think.go`, `act.go`

Define the Relationship struct:

```go
type Relationship struct {
    Reality  Reality
    Weight   float64
    Usage    int
    LastUsed time.Time
}
```

Replace `map[string]Reality` with `map[string]*Relationship` on every Reality that holds references to others:

- `Self.Realities` → `Self.Relationships`
- `Think.Operators` → `Think.Relationships`
- `Act.Operators` → `Act.Relationships`
- `Act.Providers` → `Act.Relationships` (merged)

All initial weights set to 1.0. The runtime behaves identically — every lookup that previously did `ops[name]` now does `rels[name].Reality`. Pure mechanical replacement. No behavior change.

Add `Activation(rel *Relation) float64` method on Relationship. For now returns `Weight` directly — stubs for relevance/recency/context_fit.

**Verify:** Compile. Run. Everything works exactly as before.

**Risk:** None. Mechanical rename.

---

### Step 1: Extend Relation with traversal state

**Files:** `relation.go`

Add:

```go
Trace     []RealizedStep
Depth     int
Budget    float64
Mode      string    // "act", "recall", "creative"
```

Additive. Existing fields stay. Zero-valued new fields change nothing. Sets up Steps 4-5.

**Verify:** Compile. Existing tests pass.

**Risk:** None.

---

### Step 2: Think operator dispatch → relationship traversal

**Files:** `think.go`

**Remove:** `parseOp()`, `collectOps()`, outer-op blocking (`isOuterOp`), `renderOps`/`renderOpsWithOuter`, the tag-dispatch loop.

**Replace:** Think traverses its relationships by weight. Each think pass:

1. Relationships sorted by activation. Above threshold → available to the being.
2. Operator-type relationships listed in present as capabilities.
3. Memory-type relationships attached as context.
4. LLM sees operators because the weights surfaced them, not because they're in a hardcoded map.
5. Tag parsing for invocation stays temporarily — the being still says `<bash>command</bash>`.

Inner/outer operator distinction goes away. One relationship map. Weights decide reachability.

**Fallback:** If no relationship activates above threshold, show all (identical to v.05).

**Verify:** Send impulse. Think's prompt includes operators from weighted relationships. Invocation still works.

**Risk:** Medium. Fallback makes it safe.

---

### Step 3: Act protocol enforcement → doesNotUnderstand

**Files:** `act.go`, `exchange.go`

**Remove from Act:**
- 3-retry loop for protocol violations
- Self-route detection and warning

**Replace:**
- Act calls provider once.
- Valid `<target>message</target>` tags → route.
- No valid tags → doesNotUnderstand: seed new relationship at minimum weight, inform the being.
- Self-route is structurally impossible — no relationship from self to self.

**Remove from Exchange:**
- `Peel()` for target extraction. Act sets `r.ID` directly.
- `isBeing()` type-switch. Traversability is whether the relationship has a Reality, not a type assertion.

**Verify:** Protocol violation → doesNotUnderstand. Explicit tags still route. Self-address doesn't retry.

**Risk:** Behavior change. Correct per spec. Log doesNotUnderstand events clearly.

---

### Step 4: Operator injection → relationship seeding

**Files:** `newthread.go`, `self.go`, `main.go`

**Remove:**
- `NewThread.ThinkOps` and injection loops
- `NewThread.ActOps` and injection loops
- Per-being operator wiring in `main.go` bootstrap
- `Self.Create`'s hardcoded Think/Act assembly from `r.Realities`

**Replace:**
- Genome declares operators as relationships. Bootstrap seeds them on each being with initial weights.
- `Self.Create` reads from its own relationships for topology.
- NewThread injects beings into `r.Realities` (Exchange needs them). Operators are already on the being.

**Verify:** Bootstrap. Each being has correct operator relationships. Builder's bash has higher weight.

**Risk:** Low-medium.

---

### Step 5: Self's Think-Act loop → weighted traversal

**Files:** `self.go`

**Remove:** The `for { think; act; think-back }` loop. The heart of v.05.

**Replace:** `Self.Realize`:

1. Attach identity, memory, context, desk (same as now).
2. Sort relationships by activation against current Relation.
3. Traverse in weight order:
   - Informational → attach content to Relation.
   - Traversable → call `rel.Reality.Realize(r)`.
   - Decrement `r.Budget`.
   - Budget exhausted → stop. Weight exhaustion.
4. Think is a high-weight relationship (skeleton). Act is a high-weight relationship. Think-back: Act sets `r.ID = "_think"`, traversal re-enters Think. Budget limits depth.

Thread → Exchange → Self stays hardcoded (skeleton). Weighted traversal governs what happens inside Self.

**Feature flag:** If relationships are empty, fall back to v.05 loop. Migrate one being at a time.

**Verify:** Think fires (high weight). Act fires after. Think-back works (bounded by budget).

**Risk:** High. Feature flag is mitigation.

---

### Step 6: NewThread for-loop → relationship-driven re-entry

**Files:** `newthread.go`, `exchange.go`

**Remove:** NewThread's `for { ... }` infinite loop.

**Replace:**

1. Validate access, create thread (same).
2. Inject beings (same).
3. `Exchange.Realize(r)` once. Records entry.
4. Exchange calls `being.Realize(r)` → weighted traversal from Step 5.
5. Relation returns up through Exchange → Thread.
6. `r.ID` set to new target → Thread calls `Exchange.Realize(r)` again. Explicit re-entry.
7. Stops when `r.ID` empty or `r.Budget` exhausted.

Exchange narrows to: conversation state, entries, compression, parsers. No target guessing.

**Verify:** Multi-hop (michael → skyra → builder → skyra → michael). Entries recorded. Terminates.

**Risk:** High. Feature flag — old loop alongside new.

---

### Step 7: doesNotUnderstand as growth

**Files:** `act.go`

When Act targets something not in the being's relationships:

1. New Relationship created: `Weight: 0.01`, Reality is a stub.
2. Being receives doesNotUnderstand response.
3. Subsequent encounters increment weight.
4. Threshold crossing → participates in traversal.

The being grows by reaching into the unknown.

**Verify:** Nonexistent target → new relationship. 5 attempts → weight increases. Crosses threshold → appears in context.

**Risk:** Low. Additive.

---

### Step 8: Weight updates from usage

**Files:** `relationship.go`, `self.go`

After each traversal:

- Traversed: `Usage++`, `Weight += reinforcement`
- Untraversed: `Weight *= decay`
- Power law decay (ACT-R): `decay = (time_since_last_use)^(-d)`

The graph learns. Frequent use → stronger. Neglect → fades.

**Verify:** 20 interactions using bash → weight up. Browse untouched → weight down. Traversal ordering changes.

**Risk:** Low-medium. Start conservative (0.99 decay).

---

### Step 9: Remove MemoryGraph — memory is already relationships

**Files:** `memgraph.go`, `memory.go`, `context.go`

**Remove:** `MemoryGraph`, `Entity`, `EntityEdge`, `MemNode`, `MemEdge`. `Context.Warm` cache. The entire separate memory data structure.

Memory was never a separate system. It was always relationships — the MemoryGraph was a parallel structure doing what Relationships already do. Skyra's memories of Builder don't live in a separate graph and get queried — they live inside her Relationship to Builder. When she needs context about him, she descends into the Relationship. When she needs to reach him, she follows the Reality pointer inside it. Same traversal, same mechanism, different depth.

**Replace:**
- `Memory.Store` → creates a Relationship inside the relevant parent Relationship (e.g., a memory about Builder becomes a Relationship inside the Builder Relationship)
- `Memory.Query` → descent into a Relationship's subgraph, activation-weighted
- `Context.Heat` → same descent, scoped

Migration for existing `graph.json` — entities become Relationships, entity edges become nested Relationships.

**Verify:** Existing memories load into Relationship subgraphs. Retrieval returns same results. Storage works.

**Risk:** High. Persistent state. Migrate on copy first. Keep old code behind build tag.

---

### Step 10: Clean up dead code

**Files:** All touched files.

**Remove:**
- `Think.Operators` / `Act.Operators` (now `Relationships`)
- `NewThread.ThinkOps` / `NewThread.ActOps`
- `Operators` struct and `operators.go`
- `ParseResponse` as router
- `isBeing()` type-switch
- `Peel()` for target extraction
- Feature flags
- Old `MemoryGraph` code

**Verify:** Test suite passes. `go vet` clean. No dead code.

**Risk:** Low. Cleanup.

---

## Decisions Before Starting

### 1. Activation formula scope

Start with `Weight * recency` only? Or all six factors (`edge_weight * relevance * recency * trust * relationship_weight * context_fit`) from day one?

### 2. Deterministic first

`argmax` for collapse initially, add stochastic/temperature later?

### 3. Skeleton stays

Thread → Exchange → Self stays hardcoded for now?

### 4. Intent graph is separate

Weighted relationships first, intent graph/mailbox as follow-up?

### 5. Relation signature

Keep `Realize(r *Relation) string` for now, change to `Realize(r *Relation)` in a later pass?

---

## Sequencing

Steps 0-1 are mechanical. Rename hashmaps, add fields. Zero risk.

Steps 2-4 replace dispatch logic with weight-based surfacing. Medium risk, fallbacks available.

Steps 5-6 replace the core loops. High risk, feature flags.

Steps 7-8 are the payoff. The graph grows and learns.

Step 9 is unification. One topology.

Step 10 is compression. Throw away what's left.
