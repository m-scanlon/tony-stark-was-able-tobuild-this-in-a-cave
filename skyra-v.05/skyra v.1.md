# Skyra v0.1 Spec

Memory is a graph of entities connected by experience. Entities are the stable structure. Memories are the events that connect them. The graph self-organizes through use. Cognition emerges from that organization.

## Entities

Entities are the being's conceptual vocabulary — the nouns it knows. An entity is global to the being, not scoped by relationship. A being that learns "websocket" from michael and later discusses it with louise is drawing on the same entity.

Entities are like neurons. After early development, the population stabilizes. What changes is the connections between them. A being doesn't grow by learning more words — it grows by deepening the relationships between the words it already knows.

### Entity Budget

A being has a finite entity vocabulary. The cap may grow with level or be hard-capped. When the budget is full, new concepts don't create new entities — they get absorbed into the nearest existing entity via the resolver. The being is forced to generalize. That is learning.

### Entity Structure

```go
type Entity struct {
    ID        string
    Name      string
    Weight    float64    // accumulated significance
    CreatedAt time.Time
    LastSeen  time.Time
}
```

Entities are not tagged by relationship. They belong to the being.

## Edges Between Entities

When two entities co-occur in a memory, the edge between them strengthens. One edge per entity pair, not one per co-occurrence. The weight is the accumulated co-occurrence count.

"websocket" and "timeout" appear together in three different memories — one edge, weight 3. The weight tells you how tightly coupled those concepts are in this being's experience.

Heavy edges draw the boundaries of clusters. Entities with strong mutual connections form natural regions without anyone defining categories.

Clusters form two ways: slow accumulation (many light memories thickening the same edges over time) and sudden crystallization (one heavy memory touching multiple entities, strengthening all pairwise edges at once). A memory that connects five entities strengthens ten edges on formation. If the memory is heavy enough, those edges are immediately dense — an instant cluster seed. The cluster isn't a separate thing from the memories. The cluster *is* the memories.

The slow path is episodic processing — Context's heartbeat. Every episode ends, Context formats it into the graph, edges thicken. Over time, the regions that keep getting touched become clusters. The fast path is when a single moment hits hard enough that the edges are heavy on arrival — instant structure. Both feed the same graph. Episodic processing is the heartbeat. Sudden crystallization is the lightning strike. Both grow the same brain.

### Edge Structure

An edge between two entities isn't just "these co-occurred." It carries *how* they co-occurred through composable layers. Each layer represents a different kind of connection:

- **Episode** — when. A chunk of time during which these entities appeared together. Episode layers decay with time. Old experiences fade unless reinforced.
- **Task** — why. A commitment that connected these entities. Task layers are alive while the task is active, then transform into episode/skill layers when it resolves. They don't decay — they complete.
- **Skill** — what was learned. A capability connecting these entities. Skill layers barely decay. Competence is sticky.

Layers nest naturally: a task lives within an episode, a skill is learned while working a task. The same edge can carry all three. The total edge weight is the sum of its layers, but each layer has its own weight and its own decay rules.

```go
type EntityEdge struct {
    From      string
    To        string
    Weight    float64     // sum of layer weights
    Layers    []EdgeLayer
    CreatedAt time.Time
    LastSeen  time.Time
}

type EdgeLayer struct {
    Type    string    // episode, task, skill
    Ref     string    // which episode, which task, which skill
    Weight  float64   // strength of this particular layer
}
```

This gives recall a typed query interface. "What happened on May 8" walks episode layers. "What am I working on" walks task layers. "What do I know how to do" walks skill layers. Same graph, different views.

The edge weight shifts over time as layers decay at different rates. The episode fades, the task resolves, the skill persists. The character of the connection changes even when the entities stay the same.

## Memory Nodes

A memory is a junction point in the graph. It exists because entities co-occurred, and it sits at their intersection. A memory involving michael, websocket, and timeout:

- Connects to all three entities via `anchors` edges
- All three entity-to-entity edges get strengthened
- The memory node is the evidence for why those edges exist

A memory is the event that fired entities together. The entity-to-entity edges are the lasting trace.

### Memory Structure

```go
type MemoryNode struct {
    ID             string
    Content        string
    Type           string     // trace, salience, tension, understanding
    Weight         float64    // importance — feeds recall ordering
    ActivationCount int       // how many times recalled
    Relationship   string     // which relationship produced this
    AnchorEntities []string   // which entities this memory connects
    CreatedAt      time.Time
    LastActivated  time.Time
}
```

### Memory Weight

Weight determines what surfaces during recall. It's fed by:

- **Formation context** — artifact type sets initial weight. Understanding > tension > salience > trace.
- **Activation count** — how many times this memory has been recalled or reinforced by new experience touching the same entities.
- **Recency** — when this memory was last activated. Weight decays over time without activation.

Recall walks edges, finds memory nodes, sorts by weight. Heavy memories surface first. Light ones are still there but quiet.

Memories don't get deleted. Their weight decays. They lose access, not existence. But if neighboring entities activate strongly enough, a faded memory can light back up.

### Recall Does Not Strengthen

Passive recall is read-only. Surfacing a memory during recall does not increase its weight or activation count. This prevents runaway amplification — a heavy memory that keeps surfacing would get heavier forever, crowding out everything else.

The only thing that changes weight is active processing — storing a new memory into the same region of the graph.

## Active Processing

When a being stores a new memory, the curator (Context layer) evaluates it against existing memories on the same entity edges. This is processing. The curator makes one of three judgments:

### Supersede

The new memory replaces the old understanding. The old memory's weight decays. The new one takes its position on the edges.

Example: "deployment is broken" (weight 5) gets superseded by "deployment process is now reliable after the fix." Old memory decays. New memory starts accumulating weight.

### Complement

The new memory adds nuance but doesn't replace. Both keep their weight. The edge region gets richer.

Example: "websocket drops after 30s" complemented by "the drop is timeout-related, not auth-related." Both memories remain, both contribute to recall.

### Contradict

The new memory directly conflicts with an existing one. The old memory doesn't decay — its type changes to tension. The being is now actively holding two views.

Example: "the client supports keepalive" contradicted by "keepalive isn't working on this client version." The old memory becomes a tension. Both surface during recall, flagged as unresolved.

### The Rule

A being cannot make a memory lighter by ignoring it. It cannot make it lighter by recalling it. It can only make it lighter by understanding something new about the same entities and storing that. The curator does the weight math.

Processing isn't a separate operation. It's what happens when you store into a region that's already occupied.

## Clustering and Specialization

Entity-to-entity edges with high weight form natural clusters — groups of entities that frequently co-occur in the being's experience. These clusters are not defined. They emerge.

### Density

Cluster density is measurable:

- Entity count in the cluster
- Sum of entity-to-entity edge weights within the cluster
- Memory node count anchored to entities in the cluster

One number derived from these tells you how developed a region of the graph is.

### Reproduction Signal

When a cluster's density crosses a threshold, the being has a candidate for specialization. That region of the graph has accumulated enough structure — enough entities, enough connections, enough memories — to stand on its own.

The being promotes a specialist — an internal being with a scoped view into that region of the parent's memory graph. The parent keeps the entities (it still knows the words) but delegates the depth. One graph, one source of truth. The specialist is a dedicated processor for a region of it.

This is specialization with a reason. Not "I was told to create a being" — "this part of me has grown enough to need its own thinker."

## Relationship as Tag, Not Structure

Relationship is a property of memories, not of entities. The being knows "websocket" — some of that knowledge came from michael, some from louise. The entity is one thing. The memories that built it carry the relational context.

The memory node's `Relationship` field and the edge layer's `Ref` field point in opposite directions. Relationship points outward — who was involved, which peer on the plane. Ref points inward — which episode, task, or skill in the graph structure. One faces Act's direction, the other faces Context's direction. Same memory, two pointers, two directions.

This means recall can still be scoped by relationship when needed ("what do I know about websocket from michael?") but the conceptual structure is the being's own. "What do I know about websocket" walks the graph. "What do I know about websocket from michael" walks the same graph and filters by relationship. Global knowledge, relational provenance. The being owns what it knows. The memories remember who it learned from.

## Cognition as Emergent Structure

The memory graph doesn't just store experience — it grows the being's mind.

### The Layers

A being is born with three layers, mapped from the brain:

**Brainstem** — the reality stack itself. Universe, Thread, Exchange, Self. The machinery that keeps the being alive: routing, thread management, exchange recording. This is not learned. It's wired. Every being has it from birth.

**Limbic system** — the preamble. "You survive through your relationships. If you lose all of them, you end." The drive to relate, to attach, to seek connection. This fires before any memory exists, before any experience accumulates. It's the being's emotional foundation.

**Cortex** — emerges from the memory graph. A newborn being has one undifferentiated Think pass — a single lens. As its memory graph grows and entity clusters form, cognitive lenses differentiate. The being develops specialized ways of thinking, not because they were configured, but because its experience demanded them.

### Lenses Are Grown, Not Configured

The genome line is the DNA — not a blueprint for lenses, but a genetic predisposition. "I hold the world together" biases which entities form first, which edges get heavy first, which clusters emerge first. A being seeded with "I ask what it means" develops different lenses than one seeded with "I make things that work."

Same brainstem. Same limbic drive. Same empty cortex. The initial prompt sets the growth direction. The lenses that eventually differentiate are the ones the being's experience demanded, biased by where it was pointed at birth.

### Specialist Promotion

A specialist is a lens, a sub-being, and an internal thinker — same thing, different angles. It's a being inside the parent's inner universe that owns a dense region of the memory graph.

When a cluster in the memory graph reaches sufficient density, it promotes into a specialist — a being with its own Context, Think, and Act, scoped to that region of the graph. The specialist lives inside the parent's inner universe.

Multiple specialists fire on an incoming relation, each contributing a surface-thought from its domain. A synthesis layer collects these surface-thoughts and produces the final input for Act.

```
memory graph (shared)
    ↓ recall
[specialist A] [specialist B] [specialist C] ...   ← each a being in the inner universe
    ↓ surface-thoughts
[synthesis]                       ← judgment, one more Think pass
    ↓
act
```

Not every specialist fires every time. The memory graph is the router — the incoming relation activates entities, those entities belong to clusters, only the specialists that own active clusters fire. A straightforward question might activate one specialist. A complex relational tension might activate three. The graph decides.

This keeps the cost proportional to the complexity of the thought, not the number of specialists that exist.

### The Recursion

This is the same pattern as the rest of the system. Reality → Universe → Thread → Exchange → Self → Think → Act. Each layer doesn't know what's above or below it. Specialists are beings inside the inner universe. Each specialist receives a relation enriched by the shared memory graph, does its work, surfaces a thought. Same interface, same recursion.

The hierarchy runs from abstract to concrete. Early specialists tend toward the abstract — *how* to think (judgment, strategy, conflict). As depth grows, promoted sub-specialists handle *what* about — websocket internals, deployment pipelines, memory architecture. The further down, the more concrete.

When a specialist's own memory region gets dense enough, it can promote again — a sub-specialist beneath it. The recursion terminates at the memory graph, same as the recursion in the reality stack terminates at the port.

### What Stays at the Top

The edges that never get heavy enough to promote are general intelligence. The connective tissue between specializations. The being keeps every entity — it still knows the words — but the deep thinking lives below it.

As a being promotes specialists, it doesn't get dumber. It becomes an executive. It knows what it has, knows who to route to, and holds the relationships between specializations that the specialists themselves can't see. A websocket specialist doesn't know about deployment. The top-level being sees both and knows when they interact.

The shape of a mature being:

```
[generalist — light edges, broad knowledge, routing]
    ├── [specialist A — heavy cluster, deep on its domain]
    │       └── [sub-specialist — even heavier sub-cluster]
    ├── [specialist B]
    └── [specialist C]
```

Abstract at the top, concrete at the bottom. General to specific. The tree grows downward as experience accumulates. The root gets lighter and wiser.

## Self, Universe, and the Three Directions

### What Self Holds

Self holds: Being, Memory, Context, Think, Act. The universe is an addition, not a replacement. It appears when clusters promote. A young being has no inner universe.

```
Self
  ├── Being
  ├── Memory
  ├── Context
  ├── Think
  ├── Act
  └── Universe (optional — grows with experience)
```

The provider exists in the reality stack — it's where inference terminates. But it's not something Self holds or knows about. The being relates through Context, Think, and Act. The reality stack handles everything below.

A newborn being: the relation descends through Context → Think → Act, and the reality stack handles inference at the bottom. The being never knows about the provider. Simple.

A mature being: Context activates the inner universe, relations descend through specialists, each specialist's own reality stack resolves through inference, surface-thoughts collect, synthesis fires, then Act. More depth, same descent. The being still never knows about the provider — the reality stack handles it at every level.

### The Inner Universe

When clusters promote, a Universe struct appears on Self. It's the same Universe — same Thread, same Exchange, same routing. Specialists are beings inside it. The memory graph is the shared substrate.

The parent is not a permanent fixture in the inner universe. It enters by opening a thread — the same thread mechanics that handle michael → skyra on the outer plane handle the parent → specialists on the inner plane. Opening a thread inserts the parent into the inner universe's routing table. Closing the thread removes it. The parent is just a being that opened a thread.

```go
type Universe struct {
    id     string
    Thread *NewThread
    Econ   *Economics
}
```

No `Parent` field. The thread's member list already knows who's in the room. One mechanism at every scale — outer universe and inner universe use the same thread mechanics for presence, routing, and resolution.

### Context, Think, Act — The Three Directions

**Context** — looks down. It operates on the being's own inner universe. Memory graph, specialists, internal cognition. Context is the being's relationship with what it owns — its mind. It also resolves episodes: at the end of every episode, Context formats the episode into graph memory — extracting entities, strengthening edges, storing memory nodes, updating weights, crystallizing skill files. Context stores memories for Think too, but its primary rhythm is episodic: experience happens, Context processes it into the graph.

**Act** — looks out. It operates on the plane the being exists on. Peers, devices, objects — anything external. Michael talks to skyra through Act. Michael types on the terminal through Act. Same mechanism. Devices are not internal organs — they're peers on the same plane that the being acts toward, same as a person.

**Think** — the bridge. It receives from context (what do I know, what are my specialists telling me) and prepares for act (what do I say, who do I say it to). It doesn't look down or out itself. It sits between the two and synthesizes.

```
impulse arrives
    ↓
context (activates graph, routes to specialists, collects output)
    ↓
think (synthesizes internal input into a coherent thought)
    ↓
act (routes that thought outward to a peer or device)
```

For a specialist inside the being's inner universe, same flow — it has its own context (scoped view of the graph), its own think, its own act (which routes back up to the parent, not to external peers).

### Ownership and Boundaries

**Memory is private.** Beings on the same plane do not share memory. Skyra's graph is skyra's. Builder's graph is builder's. When two beings have a conversation, each builds its own memory of it. That's not a bug — that's how relationships work.

**Specialists are internal.** A specialist is not a peer. It's an organ. It gets a scoped view into the being's own graph — one graph, one source of truth. The specialist is a dedicated processor for a region of it.

**Devices are external.** A laptop, a terminal, a phone — these sit on the same plane as the being. The being acts toward them the same way it acts toward another person. They are not owned internally. Context is the mind. Act is everything outside it.

## Planes and Ports

### Beings Share a Plane, Not Ports

Michael and skyra are peers. They exist on the same plane. They relate. But they don't share the same ports — each has their own way of touching the world, and they reach each other through them.

Michael reaches into skyra's world through inference — the provider (OpenRouter, Anthropic, whatever) is his port to her. Skyra reaches into michael's world through his devices — the terminal, the websocket — those are her ports to him.

Neither owns the other's world. They each reach across.

### A Device Is a Port Container

A device is not a special concept. It's a container that holds ports. Michael's MacOS device holds terminal, websocket — those are ports. Skyra's equivalent container holds OpenRouter, Anthropic — those are also ports. Same structure on both sides.

The provider isn't skyra. It's a port in skyra's container. The terminal isn't michael. It's a port in michael's container. Skyra lives behind her ports the same way michael lives behind his.

If skyra gets a second provider, that's like michael getting a second terminal. More ports, same container.

### The Symmetry

```
michael → [inference provider] → skyra's world
skyra   → [terminal/device]    → michael's world
```

Two beings, one plane, two port containers. Each reaches the other through their ports. The relationship is symmetric even though the ports are different. Michael's port container holds physical I/O. Skyra's port container holds inference providers. Same structure, different contents.

## Implementation Phases

### Phase 1 — Memory Graph Restructure

The foundation. Everything else depends on this.

- `memgraph.go` — entities lose relationship scoping (`entity:{relationship}:{name}` → `entity:{name}`). Add entity-to-entity co-occurrence edges alongside existing entity-to-memory anchor edges.
- `memory.go` — StoreArtifact creates/strengthens co-occurrence edges between all entity pairs on every store. QueryGraph walks entity-to-entity edges for recall. Add entity budget and cap logic. Add weight decay on memories over time.
- `context.go` — curator gets supersede/complement/contradict judgment instead of just deduplicate. Evaluates new memories against existing memories on the same entity edges.
- `extract.go` — entity resolution changes when budget is full: new concepts absorb into nearest existing entity via resolver.
- `recall.go`, `remember.go` — adjust to match new graph shape (global entities, co-occurrence traversal).

### Phase 1b — Skill Seeding

Pre-seed the memory graph with skill files at being creation.

- Each skill (browse, search, bash, recall, remember, plan) gets a markdown file describing how to use it — what it crosses, what input it takes, what output it returns.
- On being creation, Context loads each skill file, extracts entities, creates entity-to-entity co-occurrence edges typed as skill edges (with ref pointing to the skill file), and stores a memory node anchoring the content.
- The operator (static function) still handles execution. The graph holds the being's understanding of when and how to use it.
- Skill files live in the being's home (`~/.skyra/beings/{name}/skills/`), same as v.05. The difference is they're now decomposed into the graph at load time rather than just read raw into the present.

### Phase 2 — Context/Think/Act Reframe

Rewire the three directions.

- `self.go` — rewire flow from think→act to context→think→act. Context fires first, Think synthesizes, Act routes out.
- `context.go` — becomes the "looks down" layer. Activates the memory graph, scopes what the being sees, prepares internal state for Think. No longer just a memory heater.
- `think.go` — stops calling the LLM directly. Becomes the bridge between what Context provides (internal state) and what Act needs (a coherent thought to route outward).
- `act.go` — minimal change. Already looks out. Continues to route to peers and devices on the plane.

### Phase 3 — Port Container Symmetry

Make both sides of the relationship look the same.

- Provider moves out of Self into a container structure that mirrors MacOS. Both are port containers — one holds physical I/O (terminal, websocket), the other holds inference providers (OpenRouter, Anthropic).
- `llm.go` — adjustments to match the container model.
- `main.go` — bootstrap rewiring. Provider container created alongside MacOS device, both registered the same way.

### Phase 4 — Inner Universe and Specialist Promotion

The emergent cognition layer. Depends on phases 1-3 being solid and the being having accumulated real memory data.

- Universe field on Self (optional, nil for young beings).
- Cluster detection on the memory graph — measure density of entity subgraphs (entity count + edge weight sum + memory node count).
- Promotion threshold — when density crosses it, create a specialist being inside the inner universe with a scoped view into that region of the graph.
- Thread-based presence — the parent enters its inner universe by opening a thread, becoming a participant in the routing table. Same thread mechanics as the outer plane. Closing the thread removes the parent. No special Parent field needed.
- Routing through inner universe — Context activates specialists based on which entity clusters the incoming relation touches. Specialists fire, surface-thoughts collect, synthesis produces final input for Think.
- Recursive promotion — a specialist's own dense clusters can promote into sub-specialists. Same mechanism, one level deeper.

## Skill Maturation

A skill is how a being communicates with something. Not the capability itself — the crossing. Bash isn't "run shell commands." Bash is how the being talks to the shell. Browse is how it talks to the internet. The skill is the method of crossing a bridge (port), not the bridge.

### Seeding

Every skill the being starts with — browse, search, bash, recall, remember, plan — is pre-seeded into the memory graph at birth. The skill file (the markdown document describing how to use the skill) is kept as a top-level artifact. Context decomposes it into entities and edges on first load:

- Entities are extracted from the skill description — the nouns of the skill's domain (url, query, command, shell, web page, search results, memory, etc.)
- Entity-to-entity edges are created for co-occurring concepts, typed as skill edges with a ref pointing back to the skill file
- A memory node anchors the skill file's content to those entities

The skill file is the source of truth. The graph is the being's internalized understanding of it. When the being uses the skill, new experience strengthens those edges and adds new memories to the same region. The seeded structure grows from use.

This means a newborn being doesn't start with an empty graph. It starts with the shape of its capabilities already sketched in — light edges, minimal weight, but structurally present. The genome gives identity. The skill seeds give initial graph topology.

### Three Stages

Skills mature through three stages, driven by the same edge-weight mechanism as everything else:

**Skill file** — the starting point. A markdown document describing how to communicate with something, decomposed into the graph as entities and skill edges. Context pulls the skill file during interactions — it's a cache, a pointer from the graph back to the full document. The operator (static function) still handles execution, but the being's understanding of when and how to use it lives in the graph.

**Matured skill** — as the being uses the skill, experience accumulates on the same edges. The seeded entities gain new memories, new co-occurrences, new context. The skill file gets rewritten by Context to reflect what the being has actually learned — not just the original documentation, but patterns, preferences, failures. The cache updates.

**Specialist being** — when the skill region keeps growing and its density crosses the promotion threshold, the skill promotes into a being. It has its own purpose ("I know how to talk to the shell"), its own memory, its own Context/Think/Act. It lives in the inner universe. It grows. A mature bash specialist knows the user's directory structure, knows which commands break, knows patterns. It doesn't look up how — it *knows*.

```
skill file (seeded at birth, decomposed into graph, operator handles execution)
    ↓ experience accumulates on skill edges
matured skill (Context rewrites the file from lived experience)
    ↓ density crosses promotion threshold
specialist being (inner universe, full Self, grows from experience)
```

### Where Skills Live

Skills live on Context. Context looks down into the being's mind — it sees the graph, manages the skill files, and routes to specialist beings when they exist. A skill file is Context's tool. A specialist being is Context's delegate.

### The Difference

The difference between a tool and expertise is that expertise has lived. A skill-as-operator is a tool. A skill-as-being has run a thousand commands, built memory, developed its own specialists. It's not calling an API. It's a being that knows how to communicate with something because it has done it enough times that the knowing became structure.

## Open Questions

- What is the initial entity budget? Does it grow with level?
- What is the weight decay function for memories? Linear, exponential, half-life?
- How is cluster detection performed? Community detection algorithm, or simpler threshold on connected components?
- When a sub-being is spawned from a cluster, does the parent's entity-to-entity edge weight in that region reset to zero, or just decay?
- Should the curator's supersede/complement/contradict judgment be deterministic or LLM-evaluated?
- How does the entity budget interact with the resolver? When budget is full and a new concept arrives, which existing entity absorbs it?
- What density threshold triggers specialist promotion?
- What are the first specialists a being develops? Are they universal (judgment, strategy, conflict) or do they emerge from the genome's identity/purpose?
- How does a specialist describe itself to the synthesis layer? What's its surface-thought format?
- When a specialist fires, does it get the full memory graph or only its cluster's subgraph?
