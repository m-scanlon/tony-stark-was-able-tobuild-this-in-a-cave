# Ontology Spec

## The Interface

One interface. Three methods. Everything implements it.

```
Entity
  ID() string
  Relate(r Relation) Entity
  DerivePresent(r Relation) string
```

`ID` is the entity's address — the name the world uses to find it in the hashmap.

`Relate` is intake — the entity receives a relation and returns its next state. `DerivePresent` is resolution — given this relation, what is real right now? What context does this entity need to make an informed decision?

---

## Two Modes

Entities exist in two modes.

**World** — an entity that contains other entities. It holds a hashmap of child entities. When a relation arrives at a world, `DerivePresent` determines which child entity receives it, what context gets assembled, what gets routed where. The archetype determines what `DerivePresent` does.

**Invariant** — an entity that resolves to a base case. An API endpoint. A screen. A pipe. A CLI prompt. An invariant does not contain other entities. It terminates the recursion. An invariant's `DerivePresent` is whatever the endpoint does — render, execute, transmit.

Everything above an invariant is a world. Worlds nest recursively — worlds contain entities that may themselves be worlds — and the recursion terminates at invariants.

---

## Archetypes

The interface is universal. The implementation varies by kind. An archetype is a default implementation of `Relate` and `DerivePresent` for a category of entity.

Archetypes do not change the interface. They change what happens inside the two methods. From the outside, every entity looks identical. The world doesn't know or care what archetype it's routing to.

### World Archetypes

**world** — the base world. Contains child entities, routes messages, recurses. The system itself is a world. Child worlds are worlds inside worlds.

**llm** — a world of inference providers. Its children are the providers (OpenRouter, Anthropic API, local models). Its `DerivePresent` determines which provider handles a given request — model selection, fallback, load balancing. From outside, it's just an entity you send a message to.

**human** — a world of devices. Its children are the devices the human uses — laptop, phone, watch. Its `DerivePresent` determines which device receives a present — active session, notification routing, presence detection. The human is not a being in the runtime. The human's world is what entities interact with.

### Invariant Archetypes

**cli** — resolves to a terminal prompt. Renders the present as text, reads input, returns it as a relation.

**api** — resolves to an HTTP endpoint. Sends the present as a request body, returns the response as a relation.

**pipe** — resolves to a process stdin/stdout. Writes the present, reads the response. Child processes are just entities whose invariant happens to be a pipe.

**screen** — resolves to a display surface. Pushes the present as rendered UI. Input from the surface returns as relations.

**shell** — resolves to a command execution. The relation's impulse is the command. The output is the response.

---

## The Present

The present is the entity's situational awareness — the scoped context it needs to make an informed decision. `DerivePresent` assembles this context based on the entity's state and the incoming relation.

What goes into the present varies by archetype:

- An **llm** world's `DerivePresent` builds the prompt — identity, exchange history, peer list, the incoming message — and routes it to an inference provider.
- A **human** world's `DerivePresent` builds the rendered state and routes it to the active device.
- A **world**'s `DerivePresent` resolves which child entity gets the message and what context it gets.
- An **invariant**'s `DerivePresent` is whatever its base case needs — a command string, a request body, a rendered screen.

The present is always derived, never stored. It is the entity's view of reality at the moment of resolution.

---

## The Loop

The system is a recursive loop.

1. A message arrives as a Relation.
2. The world calls `DerivePresent` for the target entity — assembles the present based on the entity's archetype, state, and the incoming relation.
3. The entity responds through its base case (the invariant at the bottom of its recursion).
4. The response is parsed into new Relations.
5. Each new Relation recurses back to step 2.

The whole system is a hashmap that calls itself recursively. Worlds resolve by routing to child entities. Child entities resolve by routing deeper or terminating at an invariant. Every message eventually hits a base case.

---

## The Genome

The genome declares the entities in a world. Each line grows one entity with its archetype and initial state.

```
grow ~name skyra ~archetype llm ~identity I hold the world together. ~purpose I think, respond, and relate on behalf of the system. ~relationships michael,builder,skeptic,bash

grow ~name michael ~archetype human ~relationships skyra,builder,claude

grow ~name bash ~archetype shell ~relationships skyra
```

The archetype determines which default implementation the runtime uses at grow time. The remaining fields seed the entity's initial state — identity, purpose, relationships, whatever the archetype's `Relate` knows how to read.

An entity with no archetype declaration defaults to **world**.

---

## Self-Extension

Entities can create new entities because everything is reachable via the protocol. `grow` is just a relation targeting the world. An entity that sends `grow ~name helper ~archetype llm ~identity ...` creates a new entity in the world it inhabits.

Child processes are entities whose invariant is a pipe. Remote services are entities whose invariant is an API. The ontology doesn't distinguish between local and remote. Transport is an implementation detail of the invariant.

---

## What Dissolved

**Being** — gone. An LLM entity is a world. A tool entity is an invariant. There is no third category.

**Medium** — gone. What was medium is now the invariant at the base case. The CLI is an invariant. The shell is an invariant. Inference lives inside the LLM world's `DerivePresent`, which routes to provider invariants.

**Lens** — gone. What an entity resolves to at its base case is defined by its archetype. A screen invariant renders a present. A pipe invariant writes it. Resolution is just what the invariant does.

**The human as being** — gone. The human's world is an entity in the system. Devices are entities on that world. The runtime interacts with the human's world, not with the human.

**The device routing problem** — gone. All devices are entities in the human's hashmap on the same plane. There is no routing to a human. There are only entities on a plane communicating with each other.

**The child process problem** — gone. A child process is an entity whose invariant is a pipe. No special mechanism.

---

## What Survives

**The protocol** — relations, the `~field value` syntax, `|` for reason. Unchanged.

**The exchange system** — threads, exchanges, exchange history. How a world manages exchange state is determined by its archetype's `DerivePresent`. A world that wants threaded conversation implements threading there. A world that doesn't, doesn't.

**The memory spec** — `remember` is still an entity. Retain, recall, compress, forget. The spec doesn't change. Memory is an entity you route to.

**The internal self** — still valid. The internal self is an entity paired to another entity. The queue, the clock, the deliberation — all handled inside the archetype's `DerivePresent`.

---

## What's New

**Archetypes** — a small set of default implementations that determine what happens inside `Relate` and `DerivePresent`. The genome declares the archetype. The runtime resolves it.

**Worlds as the universal container** — everything that contains other entities is a world. An LLM is a world. A human is a world. The system is a world. Nesting is recursive.

**Invariants as the universal base case** — everything that terminates the recursion is an invariant. A screen, a pipe, a CLI, an API, a shell command. The invariant is where the system meets the outside.

---

## Future: Physics as a Composable Primitive

The archetype's `DerivePresent` currently owns all resolution logic — routing rules, context assembly, exchange management, everything. As the system grows, this may want to decompose into composable physics that can be mixed and matched across worlds. Thread economics, trust weights, consent rules, departure visibility — each could be an independent piece that a world composes into its resolution strategy. Not building this now. The archetype implementations will reveal what shape composability wants to take.
