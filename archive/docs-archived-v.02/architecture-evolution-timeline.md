# Architecture Evolution Timeline

## Scope

This is a high-level but fairly detailed timeline of how the architecture changed from the beginning of the project to the current runtime model.

It is based on the git history from:

- first commit: `2026-02-08`
- latest reviewed state in this pass: `2026-03-26`

Recent pace reference:

- `159` commits on `HEAD`
- the most recent `20` commits span about `14.1` days

This is not a full changelog.

It is an architectural timeline:

- what the system was trying to be
- what documents carried that phase
- what changed
- what survived into the current model

## Executive Summary

The architecture changed in a few major waves:

1. hardware and topology first
2. service decomposition and control-plane thinking
3. job/executor and context assembly focus
4. agents, shards, and capability distribution
5. predictive/context-engine memory experiments
6. kernel + skill + graph-centric cognitive OS phase
7. runtime rewrite into episodes, retention, and primitives
8. current actor/episode/frame/process model
9. actor substrate + recall v1 simplification + canon cleanup

The biggest long-term shift was:

- from a machine-and-service architecture
- to a graph-and-skill architecture
- to the current runtime model centered on:
  - actors
  - episodes
  - retained artifacts
  - recall
  - frames
  - structural projection

## Constants Across The Whole Timeline

A few ideas survived nearly every phase:

- local-first and user-owned
- hardware as an extension of capability
- memory as the real compounding asset
- some kind of kernel/runtime boundary
- a desire for the system to grow in capability over time

What changed was not the ambition.

What changed was the shape of the runtime beneath it.

## Phase 0: Bootstrap

### Date

- `2026-02-08`
- first commit: `2f4a672`

### What existed

Almost nothing architecturally yet.

This was the seed state before the system had a real design language.

### Why it matters

This matters mainly because everything that follows happened very quickly after this point.

The repo did not spend long in an empty or purely toy phase.

## Phase 1: Personal AI Assistant / Three-Machine Architecture

### Dates

- `2026-02-09`

### Representative commits

- `15300e4` `Initial Arch Overview`
- `3c9db19` `Add comprehensive README and architecture documentation`
- `8b4c277` `WIP: Memory architecture updates`

### Representative docs

- historical `README.md` from `3c9db19`
- historical `docs/arch/v1/scyra.md` from `15300e4`

### Main model

The system started as a strongly physical topology:

- Raspberry Pi as voice actor
- Mac mini as control plane
- GPU machine as heavy inference

Memory at this point was framed more conventionally:

- relational data
- vector DB
- object store
- project-centric retrieval

OpenClaw sat near the center of orchestration.

The dominant question in this phase was:

- how do the machines fit together?

### What was strong in this phase

- clear system vision
- concrete hardware thinking
- local-first stance from the start
- immediate concern for memory and retrieval

### What changed later

This phase was still too anchored to:

- named machines
- service boxes
- a classic orchestrator pattern

Later architecture moved away from:

- fixed machine identities
- project-centric memory
- OpenClaw as the conceptual center

### What survived

- hardware still matters
- local ownership still matters
- voice/device presence still matters
- memory remained central

## Phase 2: Service Decomposition And Control Plane

### Dates

- `2026-02-11` to `2026-02-15`

### Representative commits

- `b5aa5f7` multi-device capability and future concurrency planning
- `64e32b8` voice auth, mobile interaction, tv-actor infra
- `6bf98fa` listener/control-plane scaffolding and context compression
- `3268be2` event ingress/ack
- `ec8783f` task-formation docs

### Representative docs

- event ingress and ACK docs
- task formation docs
- listener/control-plane docs

Some of the surviving historical material for this period lives under:

- [high-level-architecture-sheet.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/high-level-architecture-sheet.md)
- [api-gateway.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/api-gateway/api-gateway.md)

### Main model

This phase decomposed the system into more explicit services:

- ingress
- ACK semantics
- control-plane orchestration
- task formation
- context injection
- delegation estimation

The dominant question shifted from:

- what machines exist?

to:

- what services and boundaries exist?

### Why this phase mattered

This was the first serious move toward runtime structure.

The system stopped being just a topology diagram and became a set of named boundaries and flows.

### What changed later

Much of this service decomposition was still too infrastructure-first.

The eventual runtime model became more centered on:

- episodes
- actors
- retained artifacts
- primitive execution

rather than purely:

- gateways
- listeners
- services
- queues

### What survived

- ingress/ACK discipline
- explicit runtime boundaries
- pressure toward decomposition instead of one monolith

## Phase 3: Arc Model, Executor Loop, And Job Thinking

### Dates

- `2026-02-20` to `2026-02-28`

### Representative commits

- `bbe8ae3` `New Arc Model`
- `b0d94e5` tool hydration, listener lifecycle, job phase model
- `8490dc2` add `next-steps.md` for executor loop design

### Representative docs

Historical docs from this phase still survive in the archive:

- [agents-services.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/agents-services.md)
- [context-engine.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/context-engine.md)
- [lifecycle.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/lifecycle.md)

### Main model

The architecture became more pipeline-oriented:

- jobs
- stages/phases
- hydration
- listeners
- executors
- routing

This was a strong attempt to define the actual operational loop of the system.

### Why this phase mattered

This is where the project started trying to become a runtime instead of just an assistant architecture.

The executor loop became a central concern.

### What changed later

The job/stage framing did not become the final backbone.

It eventually gave way to a more cognitive/runtime model built around:

- episodes
- frames
- retained experience
- actor process

### What survived

- obsession with execution flow
- awareness that the runtime loop matters more than high-level vision
- concern for hydration and staging before action

## Phase 4: Rebrand To Personal AI OS, Agents, And Shards

### Dates

- `2026-03-01` to `2026-03-04`

### Representative commits

- `1783049` `Rebrand: Personal AI Assistant → Personal AI OS`
- `090a9be` rename projects to agents, device daemons to shards
- `28fcf3c` reframe around control plane + variable GPUs + variable shards
- `9603d1c` add agent model

### Representative docs

- historical `README.md` at `1783049`
- [shard-model.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/shard/shard-model.md)
- [distributed-brain.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/native-protocol/shard/distributed-brain.md)

### Main model

This phase changed the nouns dramatically:

- projects became agents
- device daemons became shards
- GPU machines stopped being a separate architectural class and became shards too

Capabilities became more important than named machines.

This was a real conceptual improvement.

### Why this phase mattered

This was the first time the architecture moved from:

- fixed boxes on a network

to:

- distributed capability-bearing runtime surfaces

### What changed later

The idea of shards and capabilities survived.

The older agent model did not survive intact.

Later, the project shifted away from:

- many domain agents as the main conceptual unit

toward:

- actors
- contracts
- episodes
- runtime primitives

### What survived

- capability thinking
- distributed hardware as an extension of the system
- the idea that the runtime should not be tied to one machine form factor

## Phase 5: Context Engine, Predictive Memory, And Unified Heap Experiments

### Dates

- `2026-03-05` to `2026-03-07`

### Representative commits

- `93a7bb7` context engine design
- `aba797d` architecture overhaul with routers and registry
- `373372e` retrieval model with weights and decay
- `22d5fe2` domain agents as doorkeepers, unified heap
- `1adee5e` universal daemon model + execution loop + object store redesign
- `e367163` predictive memory architecture
- `47b137e` entity-level prediction in predictive memory

### Representative docs

- [context-engine.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/context-engine.md)
- [importance-vectors.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/memory/importance-vectors.md)
- [predictive-memory.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/native-protocol/retrieve/predictive-memory.md)

### Main model

This was the first really serious memory/runtime experimentation phase.

The system explored:

- a context engine as short-term memory
- warm state instead of request-time assembly only
- importance vectors
- predictive memory
- unified heap / observational store ideas
- domain agents acting as doorkeepers

### Why this phase mattered

A lot of current thinking can be traced back here.

This phase asked some of the right questions even if the answers changed:

- should context be kept warm?
- can memory be predictive?
- should retrieval happen before full inference?
- how do weights, recency, and structure act?

### What changed later

A lot of the specifics did not survive:

- session/turn/job-heavy context packaging
- unified heap framing
- context engine as the central object
- some of the blob/hydration assumptions

But this phase clearly prefigured later ideas like:

- background recall
- anticipatory retrieval
- always-on experience
- structural cue surfaces

### What survived

- predictive recall instinct
- warm context instinct
- memory as an active runtime participant, not just storage

## Phase 6: Cognitive OS, Kernel, Graph, And Skills

### Dates

- `2026-03-08` to `2026-03-14`

### Representative commits

- `ca27d7e` `cognitive OS — memory graph, kernel v1, soul/personality removed`
- `8fca363` `kernel v2 — unified protocol, skill/job/task primitives`
- `29b1e1c` `skill as first-class object`
- `56a961a` add OctaOS architecture sheet
- `03ea955` migrate to octos runtime
- `8f265dd` add skill contract notes

### Representative docs

- [kernel.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/kernel/kernel.md)
- [skill.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/skill/skill.md)
- [memory-structure.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/memory/memory-structure.md)
- [native-protocol.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01/arch/v1/native-protocol/native-protocol.md)

### Main model

This phase was the most ambitious and the most ontology-heavy.

The system became:

- a cognitive OS
- graph-first
- kernel-centric
- skill-centric
- protocol-centric

Important ideas in this phase:

- observational vs committed memory layers
- graph as the truth substrate
- skills as first-class executable objects
- native protocol primitives
- stronger trust and verification language

### Why this phase mattered

This phase overreached in some ways, but it also produced a lot of valuable instincts:

- contract thinking
- primitive thinking
- kernel boundary thinking
- memory/experience separation pressure

### What changed later

Several major ideas were later demoted or rewritten:

- the graph stopped being the primary model
- “everything is a skill” was rejected
- a lot of the trust/crypto/provisioning language became too heavy for the current runtime layer
- the memory model shifted from graph-first to episode-first

### What survived

- primitives matter
- kernel boundaries matter
- some form of contract is necessary
- the runtime needs a real execution substrate

This phase was not wasted.

It was a necessary overbuild that exposed what actually mattered.

## Phase 7: Runtime Rewrite And Episode-Based Data Model

### Dates

- `2026-03-18` to `2026-03-24`

### Representative commits

- `b9c4e91` `Checkpoint Skyra runtime rewrite and docs`
- `56dfa49` `docs: expand data model notes and simplify readme`
- `38e0095` `Refactor data model docs around retention and runtime cognition`

### Representative docs

Current canonical docs that emerged from this phase:

- [data-model-prelim.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/data-model-prelim.md)
- [retention-layer-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/retention-layer-v0.md)
- [runtime-primitives-and-artifacts-prelim.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/runtime-primitives-and-artifacts-prelim.md)
- [episode.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/episode.md)

### Main model

This was the big pivot away from the older graph/skill-centered worldview.

The architecture became:

- episode-based
- temporal
- runtime-first
- retention-aware

Major breakthroughs in this phase:

- episodes as the bounded unit of activity
- the later-superseded scored episode-side field as the structural layer
- runtime primitives versus runtime artifacts
- retained artifact family:
  - `retained_trace`
  - `retained_understanding`
  - `retained_salience`
  - `retained_tension`

This is where the architecture started feeling like a real runtime instead of a set of interesting beliefs.

### Why this phase mattered

This phase separated several things that had been collapsing into one blob:

- runtime cognition vs retained memory
- structure vs retention
- current activity vs history
- primitive execution vs longer-lived learning

### What survived

Almost everything central to the current model starts here.

## Phase 8: Actor / Episode / Frame / Process Clarification

### Dates

- `2026-03-25` to `2026-03-26`

### Representative docs

- [episode-contract-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/episode-contract-v0.md)
- [frame-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/frame-v0.md)
- [actor-contract-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-contract-v0.md)
- [actor-and-episode-ownership-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-and-episode-ownership-v0.md)
- [next-steps-recall-and-learn.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/next-steps-recall-and-learn.md)

### Main model

The architecture is now clarifying:

- `actor contract` as durable boundary
- `actor` as durable runtime operator
- `episode` as bounded runtime container
- `frame` as the projected inference page
- the then-current scored episode-side field as the structural layer recall used

This phase also sharpened the question of the actor process:

- how actors are born
- how events are accepted
- how episodes open and close
- how frames are projected
- how primitive execution stays flexible enough to support different loop styles

This is also where the current execution-flexibility requirement became explicit:

- the runtime needs a generic execution substrate
- the contract defines the allowed loop envelope
- inference can choose the next step within that envelope

That is a much more mature runtime idea than anything early in the repo.

### What survived

- actor as durable operator
- episode as bounded runtime container
- frame as inference projection
- the older scored episode-side field model as the recall driver
- the idea that execution flexibility should be bounded by contract rather than hardcoded globally

## Parallel Track: Structural Projection Demo

### Dates

- heavily explored during `2026-03-25` and `2026-03-26`

### Representative docs and files

- [structural-projection-service-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/structural-projection-service-v0.md)
- [archive/structural-extraction-pipeline-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/archive/structural-extraction-pipeline-v0.md)
- [dependency-projection-mapping-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/dependency-projection-mapping-v0.md)
- [dependency-pattern-coverage-map-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/dependency-pattern-coverage-map-v0.md)
- [HIGH_LEVEL_WALKTHROUGH.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/HIGH_LEVEL_WALKTHROUGH.md)

### Why it matters

This demo is important because it solved a hidden blocker:

- recall needed a bounded structural query input
- the earlier scored episode-side field model needed a real input path
- raw language needs to become bound entity/relationship fragments

The OpenIE path proved too noisy.

The current dependency-first projection path is the first one that feels aligned with the retained-experience runtime model.

## Phase 9: Actor Substrate, Early Recall Simplification, And Canon Cleanup

### Dates

- `2026-03-26`

### Representative commit

- `37d480e` `docs: refine actor runtime and recall model`

### Representative docs

- [actor-birth-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-birth-v0.md)
- [actor-substrate-interface-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-substrate-interface-v0.md)
- [actor-process-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-process-v0.md)
- [actor-open-questions-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/actor-open-questions-v0.md)
- [command-namespace-prelim.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/command-namespace-prelim.md)
- [interaction-unification-prelim.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/interaction-unification-prelim.md)
- [recall-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs/recall-v0.md)

### Main model

This phase turned the runtime model from a strong sketch into something much closer to an implementable base system.

The architecture clarified:

- kernel births Stark
- Stark publishes later actor contracts
- the contract itself acts as the birth spec
- actors are live immediately on instantiation
- the kernel owns the global priority heap
- each actor owns a lightweight mailbox for already-routed events
- the actor substrate exposes a small explicit runtime surface
- command execution is actor-first and primitive-based rather than command-set-based

This phase also simplified recall in an important way.

Instead of requiring the full long-horizon episode-field scoring model immediately, the docs briefly moved to:

- current stimulus
- one bounded inference call
- entity/relationship extraction
- a thin scored episode-side field
- retained-artifact lookup by `anchor_set` overlap

That was a major simplification at the time, though later canon replaced it with the recall contract driven by heavy inference calls.

### Why this phase mattered

This phase answered several questions that would otherwise have blocked implementation:

- who births actors
- what a actor's public runtime surface is
- how routed events reach a actor
- when contracts change
- how `v1` recall could first work without overbuilding the scored episode-side field

It also cleaned the canon enough that the docs now mostly agree on:

- actor
- episode
- frame
- the older scored episode-side field model
- runtime execution
- retained artifacts

### What survived

- actor/episode/frame separation
- kernel authority
- contract-bounded runtime behavior
- anchor-based recall through retained anchors
- the need for a stronger recall contract and ranking policy

## Structural Changes In The Repo

The repository structure also tells the story.

### Early layout

- root `README.md`
- heavy use of `docs/arch/v1/*`
- architecture centered in archived `v1` trees

### Mid-period layout

- `next-steps.md` moved into the architecture tree
- more and more subdomains under `docs/arch/v1`
- increasing specialization of docs

### Current layout

- [docs](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs) is canonical
- [docs-archived-v.01](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/docs-archived-v.01) is legacy
- current docs are flatter and more contract-oriented
- current docs are also more implementation-facing than the earlier architecture trees
- archive docs preserve the older `arch/v1` worldview

That structural change mirrors the architectural change:

- from sprawling future-system design
- to a smaller set of runtime contracts and working model documents

## What Actually Changed The Most

The most important transformations were:

### 1. Topology -> Runtime

Early work asked:

- what machines exist?
- what services exist?

Current work asks:

- what is the runtime boundary?
- what is an episode?
- what does inference actually consume?

### 2. Graph Truth -> Episode Truth

Earlier work often treated:

- graph structure as the primary truth model

Current work treats:

- the episode as the bounded runtime truth
- retained artifacts as selective long-lived consequence
- structure as substrate rather than sole primary object

### 3. Skills Everywhere -> Runtime Surface + Namespaced Commands

Earlier work made skills do too much.

Current work is more disciplined:

- a small runtime surface first
- explicit command sets first
- runtime state first
- higher-order compositions later

### 4. System Decomposition -> Ownership Boundaries

Earlier work named a lot of services.

Current work is getting better at naming:

- what the actor owns
- what the episode owns
- what the frame is
- what recall reads from

That is a more powerful kind of clarity.

## My Read On The Overall Arc

The project did not progress linearly.

It progressed by overbuilding several candidate architectures and then extracting the durable truths from them.

The phases that look messiest in git were often the ones that produced the strongest later boundaries.

That is especially true for:

- context-engine / predictive-memory experiments
- graph/kernel/skill overbuild
- the runtime rewrite into episodes and retained artifacts

The current architecture is stronger because it has already survived several bad or overgeneralized shapes.

## Short Framing

The project started as a local-first distributed personal assistant built around three machines and OpenClaw orchestration.

It evolved through service decomposition, agents and shards, predictive memory, graph-and-skill runtime ideas, and finally into the current actor/episode/runtime model.

The biggest maturation was the shift from broad future-system diagrams to explicit runtime contracts:

- actor
- episode
- frame
- actor substrate
- retained artifacts
- recall
- structural projection

That is the clearest sign of real architectural progress.
