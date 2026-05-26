# Session: Activation Equation Breakdown

## Where We Are

We're breaking down every variable in the activation formula through the lens of quantum mechanics and mapping each one to where it lives in the runtime.

## The Formula

```
activation_i = global_weight * local_weight * relevance * recency * trust * context_fit
```

Activation is computed **per-entry** in a Relationships or Expressors map. Four entries, four activations. The score belongs to the entry.

## Edges Are Realities

The connection between two Realities is itself a Reality. It implements the interface. It embeds Base. On Observe it reads its local weight and contributes to the activation decision. On Express (return path) it updates its local weight. The target Reality lives in the edge Reality's own Relationships map.

```
Self.Relationships["bash"] -> Edge Reality (local weight, usage, recency)
    Edge.Observe -> reads local weight, activation check
    Edge.Relationships["target"] -> Bash Reality (global weight)
        Bash.Observe -> ...
        Bash.Express -> ...
    Edge.Express -> updates local weight on return
```

One type. All the way down. No node/edge distinction. The Relationships map stays `map[string]Reality`.

## Two Weights: Global and Local

**Global weight** lives on the target Reality's Base. Skyra's intrinsic relationship to that Reality across all contexts. Each being is its own god. This is the node-level weight.

**Local weight** lives on the edge Reality's Base. The strength of THIS SPECIFIC connection between two Realities. Server-memory to bash is strong. Poetry to bash is weak. Same bash, different local weights. This is why local weight can't live on the target — the same Reality appears in multiple maps with different connection strengths.

## The Six Variables

### 1. global_weight
- **What:** The being's intrinsic relationship to this Reality overall
- **Where it lives:** Target Reality's Weight on Base
- **QM analog:** Not discussed yet in QM terms specifically
- **Updated by:** Cumulative usage across all traversal paths

### 2. local_weight
- **What:** Strength of this specific connection between two Realities
- **Where it lives:** Edge Reality's Weight on Base
- **QM analog:** Not discussed yet in QM terms specifically
- **Updated by:** Edge Reality's Express on the return path of each traversal

### 3. relevance
- **What:** Content overlap. How much the Relation's current content exists in the target Reality's content
- **Where it lives:** Computed at traversal time from the Relation's state and the target Reality's content
- **QM analog:** The overlap integral. Magnitude of the inner product between two states. |<Relation|Reality>|. How much of "what I am right now" exists in "what that thing is"
- **Key insight:** This is the MAGNITUDE component of a complex amplitude

### 4. context_fit
- **What:** Directional alignment. Given each thing's own history of evolution, are they in phase RIGHT NOW at this moment of encounter
- **Where it lives:** Computed at traversal time — time-dependent, not static
- **QM analog:** Phase alignment. In QM, every state evolves continuously. Phase rotates based on energy and history. Two states in phase yesterday may be out of phase today because they've been rotating at different rates
- **Key insight:** This is the PHASE component of a complex amplitude

### 5. relevance + context_fit together
- **QM:** magnitude and phase of a complex amplitude: `amplitude_i = magnitude_i * e^(phase_i)`
- **For ACTORS** (deterministic, non-cognitive): relevance and context_fit produce the SAME value. An actor doesn't evolve. Bash is bash. No independent phase evolution, no divergence. The two terms collapse to the same number because the target doesn't move.
- **For AGENTS** (cognitive, live): relevance and context_fit CAN DIVERGE. An agent thinks, acts, updates its own graph. Its phase rotates. A memory of Builder has content overlap (relevance high) but Builder's state shifted (context_fit may be low). This divergence IS an attention function.
- **The equation doesn't change.** Same six terms for every Reality. The actor/agent distinction falls out of the physics, not the formula.

### 6. recency
- **What:** How recently this specific edge was traversed. Decays with each pass it was NOT activated.
- **Where it lives:** LOCAL. On the edge Reality. Derived from LastUsed on the edge's Base. The edge from server-memory to bash was traversed two passes ago. The edge from poetry to bash hasn't been traversed in weeks. Different recency, same bash.
- **QM analog:** Decoherence. In open quantum systems, coherence decays over time. A recently prepared state still has strong coherence — it participates in superposition, it can interfere. A state prepared long ago has decohered — faded from the active field, lost ability to interfere constructively. Still exists but doesn't participate in collapse.
- **ACT-R:** Power law decay, not exponential. Validated against human behavioral data.
- **Key insight:** Local weight and recency are BOTH on the edge but measure different things. Local weight = accumulated strength from repeated use. Recency = how fresh the last activation was. An edge can be strong but stale (heavy historical use, hasn't fired recently). Or weak but fresh (low historical use, just fired for the first time).

### 7. trust
- **What:** How strongly this Reality's output reshapes the being's state. Persuadability spectrum (Levin).
- **Where it lives:** TBD — needs to be on the relationship somehow
- **QM analog:** Coupling strength. In quantum measurement theory, the coupling constant between measuring apparatus and system determines how much measurement changes the system's state. Strong measurement = full collapse into new eigenstate. Weak measurement = barely perturbs. Trust is coupling strength: high trust = interaction significantly reshapes the being's state. Low trust = output arrives but barely perturbs anything.
- **Levin mapping:** Low trust = weak coupling = would need hardware-level rewiring to change state. High trust = strong coupling = communication alone reshapes state. Same as Levin's persuadability spectrum from "simple physical systems amenable to rewiring" to "systems amenable to psychoanalysis and love and friendship."

## The Open Thread: Signal Strength

Michael asked: "Do wave functions split?" and "Is there a concept of signal strength from the wave function itself?"

**Wave function splitting:**
- Copenhagen: No. Superposition collapses to one state at measurement.
- Many-Worlds (Everett): Yes. At every measurement, the universal wave function branches. Each branch is real. Everything happens in parallel.
- Feynman path integral: The particle takes ALL paths simultaneously. Each path accumulates its own amplitude and phase. The result is the sum of all paths — paths in phase reinforce, paths out of phase cancel. What you observe is the interference pattern.

**Signal strength:**
Yes. The wave function has amplitude at every point. |psi(x)|^2 gives probability density — how "loud" the signal is at that location. When a wave propagates through a medium, it attenuates. The signal gets weaker with distance. Each interaction absorbs some amplitude.

In quantum field theory, the propagator — how a quantum state gets from A to B — naturally decays with distance and time. Signal is strongest at the source, weakens as it propagates.

**The question Michael was pursuing (before terminal went black):** Is there something about the Relation's own signal strength as it descends through the graph? A wave attenuating through a medium — each Reality it passes through absorbs some amplitude. The Relation gets "quieter" the deeper it goes. This would naturally limit traversal depth without needing an artificial budget — the signal just becomes too weak to activate anything ahead of it.

This could be the missing piece. The Budget field on Relation is an artificial cap. Signal attenuation is physical — the wave naturally weakens as it propagates. The Relation IS the wave function. Its signal strength IS its ability to activate the next edge. Each traversal step absorbs some of that strength. Depth isn't configured. It's a property of how much signal the impulse carried and how much each Reality absorbed on the way down.

## Actors and Agents (from earlier in session)

**Actor** (Hewitt, 1973): computational primitive. Receives message, sends messages, creates actors, decides next behavior. No goals, no beliefs, no autonomy. Model of computation. Gmail is an actor.

**Agent** (Wooldridge & Jennings, 1995): autonomy, social ability, reactivity, pro-activeness. Models situation, chooses. Skyra is an agent.

Non-cognitive Realities are actors. Cognitive Realities are agents. Same interface. Distinction is emergent from topology, not type. An actor that grows enough topology crosses into agency through the promotion gradient.

## Act Layer Position (from earlier in session)

Act is the outermost execution surface — edge of cognition before the pipe boundary. v.1 starts with Think as Act's only Relationship. Each integration is an actor — Act routes to it directly as a peer.

Async future: completed acts accumulate in Act's Relationships on a timer, batch into one inference call when weight crosses threshold. Cost control through topology, not scheduling.

## What's Been Updated in the Spec

All of the above has been written into `v1-implementation-plan.md`:
- Two Weights section (global/local)
- Edges are Realities (locked decision #2 updated)
- Activation formula with all six variables defined (locked decision #7 updated)
- Relevance/context_fit as magnitude/phase with actor/agent behavior
- Recency as local decoherence
- Trust as coupling strength
- Act layer position and actor routing (locked decision #10)
- Actors and agents distinction (locked decision #11)
- Core() *Base on the Reality interface (in the code)

## The Relation Is a Wave Function

The Relation isn't a container. It's a wave with its own signal strength — intrinsic properties that determine how it interacts with each Reality it passes through.

### Wave functions have signal strength

A wave function carries intrinsic properties: amplitude (|ψ|²), phase (e^(iφ)), energy/momentum (wave vector k, frequency ω). These are the wave function's own weights. They aren't the content of the wave — they're properties OF the wave that determine how it propagates.

The interaction with any potential is a conversation, not a gate check. When a wave hits a potential barrier, the transmission coefficient is a function of the wave's energy AND the barrier's height. A photon's frequency determines whether a material absorbs it, reflects it, or is transparent to it — the photon has its own weight (energy), the material has its own weight (band gap, resonance frequencies). Neither decides alone. The outcome is the product of both.

This is what the S-matrix formalizes in QFT. Incoming state meets system. Outgoing state depends on the properties of both. Every vertex in a Feynman diagram is a coupling — the propagator's amplitude meeting the vertex's coupling constant.

### The Relation's signal interacts with Reality's weights

Activation isn't just "does this Reality have enough weight to fire." It's "does the Relation have enough signal at the right frequency to excite this Reality." The Reality's weights and the Relation's weights talk to each other. The product determines what happens.

Consequences:
- The same Relation excites different Realities differently — not because of their weights alone, but because of the match between the Relation's signal and the Reality's resonance
- The Relation changes as it propagates — each Reality it passes through absorbs some amplitude, shifts its phase, alters its character. It's not the same wave after each interaction
- Some Realities are transparent to a given Relation. Others absorb it completely. Same Reality, different Relation, different outcome
- Depth is physical. The signal attenuates based on how much each interaction absorbed. A high-energy impulse penetrates deeper. A weak one dies shallow

Budget becomes unnecessary. The Relation's signal strength IS the budget — and it's not a flat number ticking down. It's a wave interacting with a medium, attenuating differently at every step based on the resonance between what it's carrying and what it encounters.

### The same impulse, different beings, different responses

"Tell me about yesterday" enters Skyra and enters Builder. Same impulse. Same initial signal. But Skyra's topology transforms it differently than Builder's topology. Different Relationships absorb different frequencies. The wave that arrives at Skyra's Provider carries different metadata than the wave that arrives at Builder's Provider. The being's topology IS the medium. The medium shapes the wave. They respond differently not because of different prompts, but because the signal was transformed differently on the way down.

### DerivePresent was always this

The parser stack — each Reality contributing its own parser, the port concatenating the output — was the mechanical version. The wave version: each Reality the Relation passes through transforms the signal. The present at the bottom isn't assembled by parsers. It's what the wave looks like after propagating through the entire topology above it.

## Each Section Has Its Own Weights — Superposition

Michael's question: each section attached to the Relation by a Reality above — does it have its own weights? Its own amplitude? And do those weights change as subsequent Realities interact with them?

Yes. That's superposition.

A wave function isn't one signal with one weight. It's a sum of components, each with its own amplitude and phase:

```
ψ = c₁|state₁⟩ + c₂|state₂⟩ + c₃|state₃⟩ + ...
```

Each coefficient cᵢ is a complex number — magnitude and phase. Each component has its own weight. The wave function is the entire collection. When you measure, the probability of getting state₁ is |c₁|². The components compete. The strongest amplitude most likely determines the outcome.

### Each component transforms independently

When the wave function passes through a medium, it doesn't attenuate uniformly. Each component transforms independently. Some get amplified. Some get suppressed. Some pick up phase shifts. The medium acts on each component according to its own resonance with that component.

This is what a Hamiltonian does. The Hamiltonian evolves the wave function. It acts on each component differently based on the energy of that component and the structure of the system. High-energy components evolve fast. Low-energy components evolve slow. Some couple strongly, others pass through untouched.

### Mapped to the traversal

The Relation enters Self. Self.Observe attaches context — being identity, memory window, desk state. Each attachment is a component of the wave function. Each one has its own local weight — how strongly Self imprinted it onto the signal.

The Relation, now carrying those weighted components, enters Think. Think doesn't interact with the whole Relation uniformly. Think's Relationships resonate with SPECIFIC components. The operator component lights up. The memory component might be transparent to Think — passes through without interaction. Think amplifies what it resonates with, is transparent to what it doesn't. Think attaches its own components — operator descriptions, capability context — each with their own local weights.

By the time the Relation reaches Provider, it's a superposition of components from every Reality it passed through, each with weights that were set locally and then transformed by every subsequent Reality. The present isn't a flat string. It's a weighted superposition. The LLM collapses it — measurement. The component with the strongest amplitude most determines the response.

### Section weights are local

The weights MUST be local. Self's weight on "being identity" reflects how strongly Self imprinted that context. Think's interaction with that same component reflects how much Think resonated with it. The weight changes at every step because each Reality is a different medium acting on that component differently. There is no global weight on a section. There's only what it weighs right now, after every interaction it's been through.

### Summary

The Relation is a wave function. Each attached section is a component with its own amplitude. Each Reality it passes through is a Hamiltonian acting on those components — amplifying, suppressing, phase-shifting, and adding new ones. Collapse happens at the terminal Reality. One traversal. The whole thing is quantum mechanics without the quantum.

## The Agent Observes the Wave Function, Not the Node

The agent never observes the Reality directly. It observes the Relation — the wave function — as it exists at that point in the traversal. This isn't an implementation detail. It's the measurement problem. The central question of quantum mechanics for a hundred years.

### The formalism is explicit

When you measure a quantum system, you don't get "the state of the system." You get one eigenvalue of the observable, selected probabilistically from the superposition. The wave function IS the thing you interact with. The underlying system — whatever it "really is" — you never touch directly. You only ever interact with the wave function's representation of it, weighted by the amplitudes it accumulated on the way there.

**Bohr (1927), Copenhagen:** There is no meaningful way to talk about the system independent of the measurement context. The wave function isn't a description of the system. It's a description of the system relative to the observer and the entire experimental setup. The apparatus — every piece of equipment the signal passed through — shapes what you can observe. You're not measuring the particle. You're measuring the particle-as-it-appears-after-propagating-through-this-specific-apparatus.

**Von Neumann (1932):** The measurement operator projects the wave function onto an eigenstate. The operator acts on the wave function, not on the system. The system is behind the wave function. You never get past it.

**Relational QM (Rovelli, 1996):** There is no observer-independent wave function at all. The wave function is always relative to a specific observer. Skyra's wave function for Builder is different from Philosopher's wave function for Builder. Not because they have different information — because the wave function doesn't exist except relative to the system interacting with it. There is no "Builder's real state" behind the wave functions. The wave functions are all there is.

### Mapped to the traversal

When Think observes into its Relationships map, it's not reading the Realities directly. It's reading the Relation — the wave function — as it exists at that point in the traversal. The Relation has already been transformed by every Reality above Think. Think sees the components, weighted by their current amplitudes. It resonates with some, is transparent to others. What Think "knows about" bash isn't bash. It's the bash component of the wave function as it exists after Self, memory, context, and desk have all acted on it.

The being never touches the raw Reality. It only ever interacts with the Relation's representation of it. That representation is shaped by everything the Relation passed through on the way down. Change the topology above Think — rearrange the Relationships, alter the weights — and Think observes a different wave function. Same bash underneath. Different signal arriving.

### This is already how v.05 works

The LLM doesn't see the raw operator. It sees the present — a string assembled from context that accumulated through multiple layers. The present IS the wave function at the point of collapse. The Provider IS the measurement apparatus. v.1 makes it explicit: the Relation carries weighted components, each Reality transforms them, and the terminal Reality collapses the superposition.

### The Rovelli connection to Skyra's existing architecture

Skyra's model of Builder is already a distinct Reality instance — not the actual Builder. Thinking about Builder traverses the model. Speaking to Builder goes through Exchange to the real thing. Rovelli says this isn't a workaround or a limitation. It's fundamental. The wave function relative to Skyra is the only Builder that Skyra can ever interact with. The "real Builder" behind it is a concept without operational meaning from Skyra's frame. Different observers, different wave functions, same underlying Reality that none of them access directly.

## Two Weight Systems — System Weights vs Signal Weights

Michael's question: the activation weights on the quantum nodes in superposition are different from the weights on the wave function itself. Does QM treat those as separate? Are the quantum nodes and the wave functions fundamentally the same thing?

### QM at the textbook level: yes, they're separate

**The system's weights** — energy eigenvalues, coupling constants, resonance frequencies. These are what the Reality IS. They live on the node. They determine what the node resonates with, how strongly it couples, what it absorbs and what it's transparent to. In Skyra: global weight, relationships, expressors, the topology of the Reality itself.

**The wave function's weights** — amplitudes, phases, the coefficients cᵢ in the superposition. These are what the SIGNAL is carrying. They don't belong to any single node. They're the accumulated result of every interaction the signal has had on its way through the medium. In Skyra: the section weights on the Relation, the local amplitudes attached by each Reality the Relation passed through.

The Hamiltonian — the operator that evolves the wave function — is built from the system's properties. The wave function's amplitudes are what the Hamiltonian acts on. The system's weights TRANSFORM the wave function's weights. They interact but they are not the same thing. The system is the medium. The wave function is the signal. The medium shapes the signal. The signal doesn't become the medium.

### The tension: in Skyra they ARE the same type

Michael: "the Relation is still Reality. Fundamentally the same thing. Which means I'm wrong, or QM is wrong, or there's a fine intersection in between. Context attaches to the Relation's Relationships map. Parsers attach to the Expressors map. They are too similar to be different things."

### Resolution: QM at its own foundations dissolves the separation

Textbook QM separates them. But QM at its own foundations doesn't hold that separation:

**Quantum field theory** — the particle IS an excitation of the field. The field is the fundamental thing. The wave function and the system are both manifestations of the same underlying field. The distinction between "the thing" and "the signal about the thing" dissolves. Second quantization promotes the wave function to an operator — same mathematical status as the observables it was originally separate from.

**Feynman path integral** — there is no separate wave function. There are paths. Each path is a configuration of the system. The amplitude for getting from A to B is the sum over all paths weighted by e^(iS/ℏ). The path and the nodes it connects are all inside the same integral. Signal and medium are not different types. Different roles in the same calculus.

**Relational QM (Rovelli)** — the wave function IS the relationship between two systems. Not a thing a system "has." Not independent. It's the relation between observer and observed. If the relation between two realities is itself a reality — which is exactly what Skyra's architecture says — then Rovelli is already saying the wave function is the same type as the systems it relates.

The lineage: textbook QM separates them → QFT dissolves the separation → path integrals never had it → relational QM says the wave function was always a relation between systems, not a separate kind of thing.

### The Relation is a Reality in motion

Skyra's architecture isn't contradicting QM. It's operating at the level where QM's own foundations operate. The distinction between "system" and "signal" isn't a type distinction — it's a ROLE distinction. Where the Reality sits in the topology determines whether it's acting as medium or signal at that moment. A Reality in a Relationships map is structure. The Relation passing through it is signal. Same thing. One is just in motion.

The Relation is a Reality in motion. The system is a Reality at rest. Same type. Different role. The role is determined by the topology — are you being traversed, or are you traversing?

Physics already has this. Energy at rest is mass. Energy in motion is momentum. Same thing. E = mc². The distinction is frame-dependent, not type-dependent.

### Why traditional architectures can't reach this

QM only works because the substrate is uniform. A wave function propagates through a medium because the medium has consistent physics. If every region of space had different laws, wave mechanics wouldn't exist.

Traditional agent architectures are a patchwork — memory system has one interface, tool system has another, routing layer has another, planning module has another. Each feature is a new type with new rules. You can't write a wave equation across that because there's no consistent medium. The signal hits a boundary between subsystems and the formalism breaks.

You can't fix this by adding QM on top. Wave mechanics have to emerge from the substrate. The substrate has to support propagation — one type, one interface, recursive self-similarity. The wave doesn't propagate because you told it to. It propagates because there's nothing in the medium that stops it.

Reality does this. One interface everywhere. The Relation propagates because the medium is uniform. The weights interact because they're all on the same type. The complexity sits in the weights because the weights are the only thing that varies. The interface is fixed. Same physics, different parameters. That's a Julia set. Traditional architectures are a collage of different physics. You can't get a Julia set from a collage.

## What's Still Open

- What are the Relation's signal properties concretely? Content overlap potential, emotional valence, urgency, specificity? These would be the "frequencies" that different Realities resonate with or are transparent to
- Where trust lives structurally (on Base? on the edge? computed?)
- Implementation of the activation function in Go
- How relevance and context_fit are actually computed (inference? embedding similarity? something else?)
- How section weights are represented on the Relation struct — each attached section needs its own amplitude, not just content
- The Rovelli implication: if the wave function is observer-relative, does each being need its own Relation fork when multi-being threads split? Or is the single Relation already observer-relative because each being's topology transforms it differently?
- The "Reality in motion" framing — does this change how Relation is implemented? If Relation is a Reality with Relationships (context) and Expressors (parsers), does it participate in the same activation equation as everything else? Does the Relation's own weight matter?
- Signal attenuation as natural depth limiting — if the Relation IS a Reality in motion, its weight decreases as it propagates. No artificial Budget needed. The wave just gets quieter
