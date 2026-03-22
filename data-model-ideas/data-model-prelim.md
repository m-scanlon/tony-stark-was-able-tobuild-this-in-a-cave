# Data Model Prelim

## Key Decisions (Today)

### 1. Core Model Shift: From Graph to Scoped Episodes

Decision:

A knowledge graph is not the primary data model.

Instead:

The system is composed of scoped episodes rather than a single global episode object.

The currently defined scopes are:

- node episode
- intent episode

History is reconstructed from episodes over time.

Why:

- the system is temporal, not just relational
- scoped participation matters
- intent can coordinate execution across nodes
- sequence, lifecycle, and reconstruction over time matter
- graphs are better as a secondary layer for patterns and relations

### 2. Node Episodes Are the Atomic Unit

Decision:

A node episode is the atomic unit of execution history.

Clarifications:

- local to one node
- bounded by the node's contract
- captures stimulus
- captures permitted experience
- captures produced interaction
- does not attempt to store the full system story
- may be linked through `intent_id` when execution is intent-driven

Key property:

"What this node received, used, and produced"

### 3. Clean Layered Backbone

Current structure:

- Episodes (scoped units)
- History (factual component and reconstructed views)
- Cognition (internal processing component)
- Retained Experience (memory)
- Patterns
- Standing Intents

Definitions:

- Episodes: bounded units of activity at node or intent scope
- Node Episodes: atomic local records of node participation
- Intent Episodes: higher-level episodes composed of related node episodes linked by `intent_id`
- History: the factual record within episodes; larger history views are reconstructed from episodes over time
- Cognition: the internal processing within episodes, separate from history
- Retained Experience: the longer-lived experience layer that nodes draw from through contract-bounded experience access
- Patterns: connections across retained experiences
- Standing Intents: things that generate future episodes

### 4. Episode Structure: History + Cognition

Decision:

Every episode contains two distinct components: history and cognition.

History:

- stimulus
- interaction
- external actions
- state changes
- timestamps
- append-only
- objective
- non-interpretive

Cognition:

- interpretation
- inference
- ambiguity handling
- decision formation
- emerging intent
- internal
- bounded by scope
- separate from history

Why:

- keeps facts separate from internal processing
- prevents internal reasoning from corrupting the factual record

### 5. History Is Derived, Not Stored

Decision:

There is no single mutable history object.

The system stores episodes as atomic records.

History is produced by:

- grouping episodes by scope and relation
- ordering them over time
- following `intent_id` when execution moves across nodes

This supports node-scoped and intent-scoped views without requiring a global history object.

Why:

- different lifecycles
- future reinterpretation
- cleaner evolution

### 6. Patterns Are a First-Class Layer

Decision:

Patterns exist between experiences, not inside them.

They represent:

- habits
- repetition
- similarity
- behavioral arcs

Important realization:

You cannot encode cross-time meaning inside a single experience object.

### 7. Standing Intents (Continuity)

Decision:

Some executions create durable objects that generate future episodes.

Example:

"check email daily at 9pm"

This produces future episodes such as daily checks.

So:

- episodes are bounded
- standing intents persist across episodes

### 8. Perception Definition (Critical Anchor)

Decision:

Perception is the bounded runtime frame through which a node handles the current work within its contract.

Composed of:

- stimuli
- selected history
- selected retained experience
- optionally patterns

Not:

- global state
- full history
- storage object

### 9. Purpose Reframed (Important Correction)

Rejected:

- "episodic purpose" as a field

Decision:

Purpose does not belong to an episode.

It belongs with the node definition.

### 10. Contract as Node Definition (Current Direction)

Decision:

Each node exists under a contract.

For now, `contract` is the high-level term for:

- why the node exists
- what stimuli it responds to
- what it can call
- what it can touch
- where its authority stops

The exact contract schema is not yet defined.

### 11. Frame = Scope

Decision:

The frame defines what is currently in scope for the node.

So:

- contract = why the node exists and its operating boundary
- frame = what it can handle right now

### 12. Delegation (Key Mechanism)

Decision:

When intent falls outside the frame, the node delegates.

Not because it cannot, but because it is outside the scope of the current frame.

This can create:

- a new bounded context
- delegation into more specific handling
- structural pressure when the current arrangement is insufficient

Important:

- this does not, by itself, imply a new `intent_id`
- structural change should not be treated as part of ordinary node execution

### 13. Node-Based Model (Instead of "Agents")

Decision:

Avoid "agent".

Use:

- node
- child node
- delegated node

Model:

- nodes exist under contracts
- nodes operate within frames
- nodes can delegate

### 14. Episode Needs Qualification

Decision:

Episodes are the primary unit of activity, but they exist at different scopes.

There is no single global episode object.

Use a qualified form when precision matters.

Right now, the defined forms are:

- node episode
- intent episode

Not all episodes are intent-driven.

History remains derived from episodes rather than stored as one mutable object.

The exact relationship between episode cognition and retained experience remains open.

## Final Mental Model

```text
Node (contract)
  operates within a
Frame (current scope)
  and records a
Node Episode (local scoped episode)

Related Node Episodes
  linked by intent_id
  can form an
Intent Episode (cross-node scoped episode)

Every Episode contains
  History (facts)
  and
  Cognition (internal processing)

History views are reconstructed from episodes over time.

Retained Experience and Patterns remain separate longer-lived layers.
```

Additional dynamics:

- Standing Intents generate future episodes
- Delegation can move work into more specific handling contexts
- Delegation does not automatically create a new `intent_id`

## Important Insight

- Episodes are scoped, not global
- Contract is durable and structural
- History records what happened
- Cognition captures how it was understood
- Intent coordinates execution across nodes when needed
- History is reconstructed when needed

This is a strong foundation.

## Idea for Later (Not Locked In)

### Chain of Thought as a Shared Tool

Idea:

Chain of Thought is not the front-facing mind.

It is a reusable reasoning and orchestration tool available to any node.

Used when:

- execution is unclear
- ambiguity exists
- structured reasoning is needed

Not used when:

- direct execution is sufficient
