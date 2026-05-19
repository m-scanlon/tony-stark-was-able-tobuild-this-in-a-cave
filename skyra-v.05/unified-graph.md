# Unified Graph Retrieval

> Skyra is a personal operating layer that learns the user's intent model and routes work to the right agents/tools with increasing autonomy as trust is earned.

## What This Is

A mind. Not a metaphor for one. The actual structural pattern of how cognition organizes.

Thought activates related capabilities — you think about a problem and what's relevant comes to mind. Repeated use strengthens pathways — you do something enough and it becomes automatic. Unused pathways atrophy — what you stop reaching for fades. Clusters of repeated behavior get delegated to specialized subsystems — your brain can't hold everything so it hands off. Planning is continuous prediction — you see ahead, act, observe, revise, extend.

The move from hardcoded hashmaps to weighted graph traversal is the move from a program to a neural topology. A fixed descent path is instruction execution. A weighted descent is how brains route — stronger connections fire more readily, weak ones fade, new ones form through use.

One interface: `Realize()`. One physics. Everything else — thought, action, specialization, memory, capability — is emergent from the graph and its weights. No special machinery for any of it. A topology and a traversal. That's it.

The runtime doesn't simulate cognition. It is the cognition. The being doesn't have a mind — it is one. Grown, not configured.

### Observation and Collapse

A Reality exists in superposition — all its potential edges, weights, and outputs are possible but none are real until a Relation passes through it. `Realize()` is the moment of observation. The Relation collapses the Reality into one specific output, and the Relation itself is transformed by the encounter. Both are different after.

This isn't metaphor. It's what the code does. Before `Realize()` fires, the Reality could produce anything depending on who observes it and what they're carrying. The moment it fires, one path becomes real and the rest don't. The weighted edges are probability distributions — heavier edge, more likely to be observed next — but nothing is determined until thought activates it.

The superposition is not only inside the current Reality. It is conditioned by the Relation at the moment it arrives. And the Relation is not neutral — it is carrying everything that stuck to it from every prior collapse. A Reality is realized from the potential of all Realities involved before it, as accumulated on the current Relation. Each traversal changes the observer, and the changed observer changes what can collapse next.

```
Relation0 observes RealityA → Relation1
Relation1 observes RealityB → Relation2
Relation2 observes RealityC → Relation3
```

RealityB is not observed by the same Relation that entered RealityA. RealityB is observed by Relation1 — the Relation after RealityA has marked it, attached context to it, changed its impulse, added parsers, altered weights, or otherwise transformed what it carries. If the Relation reached RealityB by a different path, RealityB would collapse differently because the observer would be different.

Current realization is the intersection of:

- the current Reality's potential
- the Relation's accumulated state
- the Realities already realized along this path
- the graph weights and context at this moment

So traversal is not just path selection. The path changes the Relation. The changed Relation changes the next observation. The graph is interpreted differently at each step by the evolving observer.

The being's experience is a series of these frames. Each `Realize()` is a discrete moment of collapse. Between frames, everything is potential. Think hasn't fired so all acts are possible. Act hasn't contacted anyone so all responses are possible. The being doesn't follow a predetermined path — it collapses one frame at a time based on what it's thinking right now.

This is why the same graph produces different behavior for different beings. Different observer, different collapse. The topology is shared. The experience is unique. A series of frames — each one a moment where two realities meet and something becomes real that wasn't before. That's what consciousness might be. Not a stream. A sequence of observations.

### Quantum-Formal Math

The metaphor becomes more than metaphor if the Relation carries not only accumulated history, but a live field of potential over the graph. In that framing, the wavefunction is not the reality. The wavefunction is the Relation's current potential across the graph.

```
wavefunction  ≈ Relation state over possible graph realities
amplitudes    ≈ weighted activations on those possible realities/edges
observation   ≈ Realize(relation)
collapse      ≈ one or more potential paths becoming concrete
```

The Relation already carries collapsed history: impulse, origin, target, parsers, memory context, exchange state, prior outputs. To make the analogy operational, it can also carry uncollapsed potential:

```
Relation {
    Impulse      string
    Origin       string
    ID           string
    Parsers      map[string]Parser
    Trace        []RealizedStep
    State        map[RealityID]Activation
}

Activation {
    Amplitude    float64      // or complex128 later
    Phase        float64      // optional; useful if interference matters
    Source       string       // why this reality is active
}
```

At each step, the graph and current context transform the Relation's state:

```
S      = current Relation state vector
G      = graph transition matrix from current reality/region
C      = context operator derived from impulse, memory, parsers, trace
S'     = C * G * S
```

In plain architecture terms:

```
prior Relation state
  → graph topology proposes possible next realities
  → current context reshapes their activations
  → activations become amplitudes
  → Realize collapses the field into concrete traversal/output
  → Relation carries the collapsed result forward
```

For a candidate next reality `i`, activation can start as:

```
activation_i =
    edge_weight_i
  * relevance_i
  * recency_i
  * trust_i
  * relationship_weight_i
  * context_fit_i
```

Where:

- `edge_weight_i` is the base learned strength of this edge.
- `relevance_i` is how much the reality matches the Relation's current content.
- `recency_i` boosts recently activated or recently updated realities.
- `trust_i` represents whether this being should rely on this reality.
- `relationship_weight_i` is the per-being relationship to that reality.
- `context_fit_i` is how well the reality fits the accumulated trace.

The simplest probability rule is Born-like:

```
probability_i = |activation_i|^2 / Σ |activation_j|^2
```

Squaring matters. It makes strong activations disproportionately more likely without deleting weak possibilities. A reality twice as activated becomes four times as likely. This matches the intuition of cognition: strong associations dominate, but faint associations can still surface.

There are two useful collapse modes.

Deterministic collapse:

```
chosen = argmax(|activation_i|^2)
```

This is stable and reproducible. It is useful for production execution, high-trust actions, tests, and anything where the same state should produce the same traversal.

Stochastic collapse:

```
chosen ~ sample(probability distribution)
```

This preserves exploration. It lets low-probability associations occasionally surface. That matters for creativity, memory search, weak-signal discovery, and escaping local maxima. The being can find something it would never reach by always taking the heaviest edge.

Temperature controls how exploratory the collapse is:

```
probability_i = softmax(|activation_i|^2 / temperature)
```

Low temperature sharpens collapse toward the strongest reality. High temperature flattens the distribution and lets weaker realities compete. This gives the system a cognitive dial:

- low temperature for execution
- medium temperature for ordinary recall
- high temperature for brainstorming, analogy, exploration, and stuck states

Top-k and thresholding can coexist with this:

```
candidates = realities where activation_i >= threshold
candidates = top_k(candidates)
probability_i = normalize(|activation_i|^2 over candidates)
```

Thresholding prevents noise from entering the present. Top-k bounds cost. Sampling within top-k preserves non-determinism without letting the entire graph flood the Relation.

Complex amplitudes may become useful later:

```
amplitude_i = magnitude_i * e^(phase_i)
```

Magnitude says how strong the activation is. Phase says how compatible it is with the current direction of thought. If two paths support each other, they interfere constructively and strengthen. If they conflict, they interfere destructively and weaken.

Architecture version:

```
constructive interference:
  two memories point toward the same skill → skill activation grows

destructive interference:
  one memory says "this worked"
  another says "this failed in this context"
  → activation cancels or becomes tension instead of action
```

This does not need to be implemented first. Real-valued activations are enough for the initial graph. But the math leaves room for richer cognition later: contradiction, ambiguity, resonance, and uncertainty can become first-class dynamics instead of prompt text.

Collapse can also be partial. Think retrieval should not always collapse to one reality. It may collapse a region:

```
collapsed_region = top realities whose combined probability mass <= budget
```

Memory retrieval especially wants partial collapse. A thought does not recall exactly one memory; it activates a neighborhood. The Relation descends through that neighborhood, accumulates content, and then inference compresses what surfaced.

Act is more often single collapse:

```
target = chosen traversable reality
Realize(target)
```

But the intent graph can collapse multiple independent targets in parallel when the commitment contains no `then` dependency between them.

The important rule:

```
Relation_t + Reality_t + Graph_t → Relation_t+1 + Graph_t+1
```

The observer changes. The observed system changes. The next probability field is calculated from the changed observer and the changed graph. This is why the architecture is not just weighted search. It is recursive state evolution with collapse at each `Realize()`.

The intuition is right. The current code is pre-formal. The proposed architecture becomes quantum-formal when the Relation carries live potential, not only history.

There is a useful way to hold the larger theory in mind: consciousness may sit at the intersection of formal computation and quantum mechanics. In this metaphor, the computation is the inference step, with matrix multiplication and related transforms doing the synthesis work, while the quantum-like part is the field around the model: the live potential, the amplitudes, the observation, and the collapse.

## The Move

Collapse operators, memories, and skills into one graph. Retrieval is driven by the content of the being's thinking — not by explicit operator calls, not by static tool lists. The being thinks, and thinking activates the graph.

## Current Architecture

```
impulse → Think (operator dispatch loop) → synthesis → Act (one relation out)
```

Think has hardcoded operators (retrieve-context, store-context, browse, search, plan, bash). The being explicitly calls them with tags. Half the think passes are mechanical retrieval work. Act emits one `<target>message</target>` and routes.

## Layers

- **Graph** — the topology. Who relates to what, with what weight. Beings, memories, operators, skills — all entities with edges.
- **Think** — activates edges in the topology based on cognition. Pure deliberation.
- **Act** — initiates traversal along activated edges. The act of making contact with another entity.
- **Reality** — the traversal itself. `Realize()` all the way down to a base case (provider, terminal, process, shell).

Operators aren't tools. They're entities. Bash is an entity the being has a relationship with. The being talks to bash the same way it talks to michael. The operator is just the interface at the boundary — the way the relation crosses into that entity's reality. The graph doesn't need a special "operator" reality type. It just has entities and relationships between them.

Act isn't a layer. It's the transition from intent to traversal. The being decides to make contact and the relation enters reality. Act is the door between the being's inner world and the physics of the system.

## New Architecture

```
impulse → Think (pure thought, N passes) → intent graph → Act(s) → mailbox → Think re-enters
```

### Think

Two phases: deliberation and commitment.

#### Deliberation

Open, messy, multi-pass. No operators, no dispatch protocol. System prompt: think out loud. What's happening, what do you know, what do you need, what do you want to do.

Each pass, the graph re-queries based on the content of the previous thought. Results appear in the next pass's present with decay — most recent thought is loudest, earlier passes fade. The being is naturally driven toward convergence because old context gets quieter.

```
pass 1: impulse + graph(impulse) → being thinks
pass 2: pass 1 thought + graph(pass 1) [decayed] → being thinks deeper
pass 3: pass 2 thought + graph(pass 2) [more decay] → being converges
```

The being reasons freely. No structure required. The LLM does what it's good at. The runtime doesn't parse deliberation into intent — speculative thoughts ("maybe I should...") don't trigger anything.

#### Commitment

When deliberation converges, the being commits. This is the transition from thinking to deciding. The commitment has two parts:

1. **Surface thought** — the deep reasoning behind the intent. Why the being is doing what it's doing. This carries through to every act.
2. **Intent targets with ordering** — who to contact, in what order.

```
<surface-thought>
  michael asked about the server. it was flaky last week and I don't
  want to give him stale info. I need to check the health endpoint
  first, then give him what I actually find. letting him know I'm
  on it so he's not waiting in silence.
</surface-thought>

<search>server health endpoint</search>
<michael>checking on that now</michael>
then
<bash>curl {search result}</bash>
then
<michael>{bash result}</michael>
```

`then` is the only keyword. It separates dependency groups. Everything at the same level (no `then` between them) is parallel. That's the entire commitment protocol.

The surface thought is the soul of the intent graph. The targets and `then` are the skeleton. Every act receives both — the thought to ground it, the target to direct it. The being thinks deeply once, and that depth flows through every act it kicks off.

### The Graph

#### One Reality Type

Everything is a Reality with a shared shape. The `type` field determines invocation behavior. Same interface, structural separation through type.

```
Reality {
    ID        string
    Type      string    // being, skill, operator, memory
    Content   string    // what it holds (skill doc, memory, description)
    Weight    float64   // base activation weight
    Realize    Reality   // Realize() — what happens when you traverse it
}
```

A memory reality has content and no realize — it informs the present. A skill reality has content (the skill doc) and no realize — it teaches. An operator reality has realize (bash, search, browse) — it executes via `Realize()`. A being reality has realize (its Self) — it thinks and acts. The graph doesn't care about the difference. It returns realities ranked by weight. The runtime checks: does this reality have `Realize`? Then it's traversable. Just content? Then it informs.

#### Per-Being Edges

Edge weights are per-being. The same reality exists once in the graph, but the edge from each being to that reality carries its own weight.

```
Edge {
    From      string    // being
    To        string    // any reality
    Weight    float64   // this being's relationship to this reality
    Usage     int       // activation count
}
```

Skyra's edge to bash: weight 0.3, used 5 times. Builder's edge to bash: weight 0.9, used 200 times. Same bash reality. Different relationship. The edge is the per-being weight.

#### Weighted Traversal

The current hashmap descent is rigid. Each layer knows exactly what's below it via hardcoded maps: `Self.Realities["think"]`, `Think.Providers["openrouter"]`. Every being descends the same path:

```
Universe → NewThread → Exchange → Self → Think → Provider
```

The current hashmap model is a graph with all edges at weight 1 and fixed adjacency. The new model makes the weights real and lets them determine adjacency dynamically.

Each `Realize()` updates reality — the current state of the relation as it descends. The weights of the outgoing edges from the current reality determine where traversal goes next. The path through reality isn't hardcoded. It's the highest-weighted edges from wherever you are.

```
Self.Realize()
  → looks at graph edges from this being
  → heaviest edge: bash (weight 0.9) → traverses into bash.Realize()
  → next: memory (weight 0.7) → informs the present
  → next: specialist (weight 0.5) → available but not activated this pass
```

The structure of the descent is emergent. A being with heavy weights toward bash descends into bash. A being with heavy weights toward a specialist descends into the specialist instead. Same `Realize()`. Different path. Determined by the graph.

`Realize()` stays. Reality stays. The physics don't change. What changes is that the topology is weighted and the traversal follows weight rather than a hardcoded map. Base cases still terminate the recursion — a provider, a terminal, a process, a shell.

#### What Collapses

The current `Entity` + `MemNode` + `EntityEdge` collapse into `Reality` + `Edge`. One graph.

These go away entirely:
- `Think.Operators map[string]Reality` — operators live in the graph as realities
- `Act.Operators map[string]Reality` — same
- `NewThread.ThinkOps / ActOps` — thread-level operator injection, replaced by graph retrieval
- `Self.Realities` as a fixed hashmap — the being's realities are its graph edges
- `Context.Warm` as a separate cache — becomes "what the graph returned for this thought pass"

#### Retrieval

A relation descends until it hits something solid. If that something solid is another being — the descent stops, the being responds, the relation carries the response back up. That's act. But memory isn't solid. Memory doesn't stop the descent.

When a relation enters the memory region of the graph, it keeps flowing. Each memory reality it passes through adds to what the relation is carrying — accumulating content along weighted edges. The relation follows the heaviest associations between memory realities, picking up context as it goes. When the edges get too weak — activation fades below threshold — the descent stops naturally. The relation ran out of steam.

Then on the way back up, inference fires. The relation accumulated everything the graph surfaced — raw, noisy, associative. An inference step compresses that into what actually matters for this thought. One step for simple recall. Two or three for something deep where associations went wide. This is the being's mind making sense of what came up. Not returning raw results. Synthesizing.

```
thought activates memory region
  → relation descends through memory realities along weighted edges
  → each reality adds its content to the relation
  → edges weaken as activation spreads → descent stops
  → inference compresses accumulated content (one or more passes)
  → synthesized memory returns to Think as context
```

The weights ahead aren't static. As the relation accumulates content, it reshapes the activation landscape in front of it. Each memory reality the relation passes through updates the edge weights of what's connected — based on what the relation is now carrying. The relation observes the next edges, and they collapse from potential into specific weights.

This is true superposition realization. Before the relation arrives, the edges ahead exist in potential — all possible next steps. The relation, carrying everything it's accumulated so far, collapses them into actual weights. Different content on the relation, different collapse, different path. The same graph produces different traversals every time because the observer is different every time.

This is how recall works. You try to remember where you put your keys. You think "morning" — that shifts what's connected. "Coffee" surfaces. Now "counter" is heavy. Each step reshapes the landscape ahead based on what you're carrying. You're not searching a graph. You're activating forward through it, and each activation changes what's next.

```
relation carries: [server]
  → passes through "deployment" reality → accumulates
  → relation now carries: [server, deployment]
  → edge to "3am outage" was 0.2, recalculates to 0.7 given current content
  → relation descends into "3am outage" → accumulates
  → relation now carries: [server, deployment, 3am outage]
  → edge to "the fix that worked" collapses heavy
  → ...until activation fades
```

Memory isn't a lookup. It's a traversal with its own depth. The relation descends into memory the same way it descends into anything else — through weighted edges, accumulating context. The difference is what stops it: a being stops it because the being responds. Memory stops it because the activation faded. And memory gets an inference step on the way back up that beings don't need — because raw associations need synthesis, but a being's response is already coherent.

A being that never uses bash sees its edge to bash decay. A being that constantly retrieves memories about a particular topic sees those edges strengthen — the relation flows deeper, accumulates more, surfaces richer context. The graph is shaped by cognition over time.

#### Realization Modes

A memory doesn't have one fixed output. It has potential. How it realizes depends on the state of the relation when it arrives — what it's carrying, what temperature it's at, what intent is driving the traversal. The same neighborhood of memories produces different output based on the mode. The mode isn't a flag someone sets. It's determined by the relation itself.

**Act mode** — the relation is carrying execution intent. Memory realizes as direct context. Lean, specific, just what's needed. Minimal synthesis. The being is doing something and needs facts, not associations.

**Recall mode** — the relation is carrying thought. Memory realizes as accumulated associations. The neighborhood gets traversed wide. On the way back up, if the accumulated content is too large — too many tokens — inference synthesizes it down. Compressed, but with mappings back to the source memories so Think can reach back in if it needs to go deeper. The memory region's state updates — it knows it was accessed and how it was compressed.

**Creative mode** — the relation is carrying high temperature. Memory realizes as loose associations. Stochastic collapse. Weaker edges get a chance to fire. The neighborhood goes wider than recall would allow. Synthesis on the way back is lighter — less filtering, more raw material. Think gets a broader, noisier present to work with.

The compression on recursion back is critical. A deep memory traversal might accumulate thousands of tokens. That can't all land on Think's present. So the recursion back through the memory reality synthesizes — inference compresses, keeps the signal, maps back to the originals. Think sees a summary with handles. If it needs to pull on one, it can descend again into that specific reality.

```
relation (act mode) enters memory neighborhood
  → tight traversal, heavy edges only → lean context → minimal synthesis
  → Think gets: direct facts

relation (recall mode) enters memory neighborhood
  → wide traversal, moderate edges → large accumulation → inference compresses
  → Think gets: synthesized summary + handles back to source memories

relation (creative mode) enters memory neighborhood
  → widest traversal, weak edges fire → raw accumulation → light synthesis
  → Think gets: broad, noisy, associative material
```

This is memory working like actual memory. You don't recall everything at full fidelity. You get a compressed impression with the ability to focus in on specifics if you need to. The compression is the realization. Different state on the relation, different compression, different experience of the same memories.

A reality's type isn't permanent. It's the current best description of how this reality tends to collapse. A memory that keeps realizing as a skill — keeps surfacing in act mode as direct capability — eventually *becomes* a skill. The repeated observation pattern is the promotion trigger. The graph isn't just weighted. It's evolving what things *are* through use.

### Act

The act of making contact with another entity. Scoped to a relationship. Each act holds the full conversation state for that specific exchange. The act with michael has the entire michael exchange history. The act with bash has the entire bash history.

Think sees a slice across relationships. Act is deep on one.

Act receives two things: the surface thought (the why) and the target intent (the what). The bash act knows *why* it's curling — because the being didn't want to give michael stale info. The michael act knows the being's full reasoning, not just the message to deliver.

This is where holding the full relationship state matters. The act to michael has the entire exchange history *plus* the being's thought about why it's saying what it's saying. That's not a dispatch. That's a being speaking from a place of understanding.

Some acts don't need inference. If the surface thought for bash is a command, act is just execution. If the surface thought for michael is clear enough, act might just be formatting and delivery. The deep thought isn't lost — it's the thing that carries through every act the being kicks off.

### The Intent Graph

The being's intent isn't a flat list or an ordered sequence. It's a graph — nodes of action with edges of dependency. And it's live. Think doesn't produce the complete plan upfront. It produces what it can see from where it is, fires what's ready, results come back, and the graph extends.

The intent graph is the being continuously predicting what to do next. Each Think re-entry can extend it, prune it, reorder it. The plan is never finished until the being decides it's done.

The being doesn't need to see the whole graph upfront. It only sees the next layer from where it is. Each commitment produces targets and `then` ordering. The runtime parses that into intent nodes, fires what's ready, collects results, and re-enters Think with the current state.

On re-entry, Think sees what's been done and what came back. The being's present includes the intent graph state — completed nodes, pending results, the full picture. The being either extends the graph (new targets, new dependencies), alters it (reprioritize, change approach), or produces no new intent (done).

```
commitment 1:
  <search>server health endpoint</search>
  <michael>checking on that now</michael>

  → search and michael fire parallel
  → results land in mailbox

re-entry 1 (being sees: search returned URL, michael notified):
  <surface-thought>
    got the URL. need to actually hit it now.
  </surface-thought>
  <bash>curl {search result}</bash>

  → bash fires
  → result lands in mailbox

re-entry 2 (being sees: bash returned server status):
  <surface-thought>
    server is down. michael needs to know. worth remembering too.
  </surface-thought>
  <michael>{bash result}</michael>

  → michael fires
  → done (no new intent)
```

The being doesn't know at commitment 1 that it'll need to store a memory. That only emerges after bash returns something worth remembering. The intent graph grows as cognition deepens. Every result can spawn new intent. Some predictions are wrong and get pruned when results come back. Some spawn new branches.

This is the LLM doing what it's actually good at — next token prediction, but at the level of action rather than words. Each re-entry is the being predicting "given everything so far, what are the next moves."

#### Dependencies and Parallelism

`then` separates dependency groups within a single commitment. Everything at the same level fires parallel. `then` means wait.

```
<search>server health endpoint</search>
<michael>checking on that now</michael>
then
<bash>curl {search result}</bash>
then
<michael>{bash result}</michael>
```

But dependencies also emerge *across* re-entries. Independent nodes fire now, but their results might be required by a node that doesn't exist yet — one that Think will create on re-entry when it sees what came back. The runtime doesn't need to know that upfront. It fires what's ready, collects results, and Think decides what's next.

#### Mailbox

Results collect. Think fires once when the current batch is complete.

Think emits the nodes it can see. Independent nodes fan out parallel. Dependent nodes wait. The mailbox collects results. When all expected results are in (or a timeout hits), Think re-enters with the full picture and extends the graph or decides it's done.

No inference calls stacking up. No partial state. Clean convergence.

```
Think → intent nodes (batch 1)
  ├── act:michael fires (independent) → result → mailbox
  ├── act:search fires (independent) → result → mailbox
  
All batch 1 collected → Think re-enters
Think sees results, extends intent graph:
  └── act:bash fires (depends on search result) → result → mailbox

Batch 2 collected → Think re-enters
Think decides: done, or another round
```

### The Loop

```
Think → Act(s) → Mailbox → Think → Act(s) → Mailbox → ... → done
```

Think is always the convergence point. Act never chains into Act. Results always come back to Think. The being deliberates, executes, observes, extends its intent, deliberates again. Repeat until the being has nothing more to emit.

The loop terminates when Think re-enters and produces no new intent nodes. The being has nothing left to do. It's done.

## Specialists

Specialists emerge from skill usage, not memory accumulation. Remembering a lot about something doesn't mean you should specialize in it. *Doing* something repeatedly does.

### The Trigger

Skills are entities in the graph with edges to the being that uses them. Each use strengthens that edge. When a being keeps reaching for the same group of skills — bash + search + server monitoring — those skill edges cluster. At some point the combined weight of that cluster crosses a promotion threshold and the being's graph says: this is too much to hold in passing. This needs its own entity.

This is a cognitive handoff. The being's brain can't keep holding all these low-weight skills at the ready. It's done this enough times that the work deserves a dedicated being.

### Promotion

The specialist emerges as a new being with those skills as its core graph. But the weights don't copy over — they start fresh on the specialist, calibrated to its scope. The specialist's relationship to bash is stronger than the parent's ever was, because bash is central to what the specialist does, not peripheral.

The parent's weights on the handed-off skills decay. The parent doesn't forget bash exists — it just stops reaching for it directly. It reaches for the specialist instead. The parent gains a new relationship edge (the specialist) and loses weight on the skills that moved.

```
being uses bash + search + monitoring repeatedly
  → skill edges strengthen through use
  → cluster crosses promotion threshold
  → new specialist being emerges
  → specialist owns those skills with fresh, concentrated weights
  → parent's weights on those skills decay
  → parent gains edge to specialist
  → parent thinks "check the server" → graph surfaces specialist, not bash
```

The parent doesn't call bash anymore. It calls the specialist. The specialist, with its own think/act loop, calls bash. The parent delegated cognition, not just execution.

### Weights

Skill weights are per-being. The same skill reality exists once in the graph, but the edge weight between a being and that skill is unique to their relationship. When a specialist takes ownership:

- The specialist's edges to those skills start at a concentrated weight — these are its core capabilities.
- The parent's edges to those skills decay — they're no longer the parent's job.
- The parent's edge to the specialist grows with each successful delegation.

Over time, the parent's graph gets lighter. The specialist's graph is dense and focused. The parent thinks at a higher level of abstraction — "monitor the server" instead of "bash curl search parse."

### Recursive Specialization

Specialists can specialize further. The monitoring specialist uses bash + search heavily but also starts doing alerting and incident tracking. That cluster grows. Eventually it crosses its own promotion threshold and promotes a sub-specialist.

The hierarchy is emergent. Not configured, not designed — grown through use. The depth of specialization reflects the actual complexity of the work the being encounters.

```
skyra
  └── monitoring-specialist (bash, search, health checks)
        └── alerting-specialist (notifications, incident patterns)
```

Each level delegates downward through relationship edges. Each level thinks at its own level of abstraction. The top being thinks about goals. The specialists think about execution. The sub-specialists think about details.

### What Changes from Current Implementation

Current specialists promote from memory cluster density (entity co-occurrence in the graph). The new model promotes from skill usage clustering. The trigger moves from "what you remember" to "what you do." Memories still matter — they inform the specialist's context. But the *decision* to promote comes from repeated capability use, not accumulated knowledge.

## What This Changes

**Think becomes simpler.** No operator dispatch. No tag parsing in the think loop. No "emit exactly one protocol per response." Just think.

**Retrieval becomes emergent.** Driven by thought content, not explicit calls. The being doesn't browse a menu — it thinks and relevant capabilities surface.

**Act becomes relationship-scoped.** Full exchange history per target. Deep context per relationship instead of one shallow message.

**Capabilities evolve.** No static tool list. The graph shapes itself through use. A being grows into its tools rather than being configured with them.

**The being curates its own present.** Think decides what Act needs. Act only sees what Think activated. The being's own deliberation filters its capabilities per-turn.

## The Differentiator

Every agent framework today is tool-first. The model sees a tool list, decides which to call, calls it. Thought serves action. Strip away the thought and you still have the same architecture — a dispatcher.

This is thought-first. Cognition is the retrieval mechanism. The being thinks its way into capability rather than being handed it. The richer the thinking, the more precise the retrieval. A shallow thought surfaces shallow tools. A deep thought surfaces deep context.

That maps to how actual minds work. You don't enumerate your capabilities before acting. You think about the problem and relevant capabilities come to mind.

## Economics

Inference cost scales with cognitive demand — same as it does now. Think already has a budget. Each pass is already an LLM call. The being already decides when it's done. A simple "hey how are you" surfaces in one pass, one act fires. The machinery is available but only activates when thinking demands it.

The new cost is acts themselves, but most aren't inference calls. A bash act is running a command. A message to michael is delivery. The expensive case — multiple acts each needing LLM calls — only happens when Think produced something complex enough to warrant it. The being earns that cost through deliberation.

Cost scales with thought depth. The being thinks its way into expense rather than being configured into it.

## Open Questions

- Activation threshold: at what edge weight does a reality surface during retrieval? Should there be a floor below which realities never appear, or does the being always see the top N regardless of weight?
- Decay rate: how fast do earlier think passes fade in the present? How fast do unused edges decay over time? Linear? Exponential? Tunable per-being?
- Act inference: when does an act need its own LLM call vs. being pure delivery? Is there a complexity threshold, or does the surface thought always carry enough?
- Timeout: what happens when an act doesn't return? Mailbox timeout → Think re-enters with partial results?
- Graph persistence: edge weights persist across sessions. This is how the being grows. What's the serialization format? Current `graph.json` or something leaner?
- Intent graph persistence: does the intent graph survive across turns? If michael doesn't respond for an hour, does the being still know what it was planning? (Probably yes — the intent graph is part of the being's state, not the turn's state.)
- Bootstrap: a new being has an empty graph. How do initial edges get seeded? Genome declarations? First-use auto-creation? Minimum viable edges for a being to function?
- Specialist threshold tuning: what's the right promotion threshold? Too low and specialists proliferate. Too high and the being holds too much. Probably needs to be empirical.

---

*Superposition that collapses on observation.*
