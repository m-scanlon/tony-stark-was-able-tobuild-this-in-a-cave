# Data Model Prelim

## Key Decisions (Today)

### 1. Core Model Shift: From Graph to Episodes

Decision:

A knowledge graph is not the primary data model.

Instead:

The system is centered around episodes, which are bounded units of activity.

Why:

- the system is temporal, not just relational
- sequence, lifecycle, and meaning over time matter
- graphs are better as a secondary layer for patterns and relations

### 2. Episodes Are the Core Unit

Decision:

An episode is a bounded, coherent unit of activity grouped by purpose and time.

Clarifications:

- not just user interaction
- not just time grouping
- includes stimuli
- includes internal processing
- includes external actions such as API calls
- includes outcome

Key property:

"A thing that started, happened, and finished"

### 3. Clean Layered Backbone

Current structure:

- Episodes (core)
- History (derived)
- Retained Experience (memory)
- Patterns
- Standing Intents

Definitions:

- Episodes: what happened as a coherent activity
- History: what factually occurred in the episode
- Retained Experience: what mattered about the episode
- Patterns: connections across retained experiences
- Standing Intents: things that generate future episodes

### 4. Separation of Fact vs Meaning

Decision:

History and Retained Experience must be separate.

History:

- factual
- minimal
- stable
- append-only

Retained Experience:

- interpreted
- enriched with meaning, salience, and confidence
- evolves over time

Why:

- allows reinterpretation later
- prevents meaning from corrupting ground truth

### 5. Loose Coupling Between Layers

Decision:

These layers should be linked, not embedded.

Examples:

- episode -> history
- episode -> retained experience
- retained experience -> patterns

Not:

- one giant JSON object

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

Some episodes create durable objects that generate future episodes.

Example:

"check email daily at 9pm"

This produces future episodes such as daily checks.

So:

- episodes are bounded
- intents persist across episodes

### 8. Perception Definition (Critical Anchor)

Decision:

Perception is the bounded runtime frame for the active episode.

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

Purpose belongs to the node, not the episode.

### 10. Intrinsic Purpose (Strong Concept)

Decision:

Each node has a stable intrinsic purpose.

Example:

"stay attentive to Mike"

This governs:

- what it should handle
- what it should delegate
- how it maintains coherence

### 11. Frame = Scope

Decision:

The frame defines what is currently in scope for the node.

So:

- purpose = why the node exists
- frame = what it can handle right now

### 12. Delegation (Key Mechanism)

Decision:

When intent falls outside the frame, the node delegates.

Not because it cannot, but because it is outside the scope of the current frame.

This creates:

- a new bounded context
- a child node
- a more specific purpose

### 13. Node-Based Model (Instead of "Agents")

Decision:

Avoid "agent".

Use:

- node
- child node
- delegated node

Model:

- nodes have intrinsic purpose
- nodes operate within frames
- nodes can delegate

### 14. Episodes Still Matter (But Not for Purpose)

Decision:

Episodes remain the record of what happened.

But:

- they do not carry deep purpose
- they do not define system identity

They are the unit of activity, not the unit of intent.

## Final Mental Model

```text
Node (intrinsic purpose)
   ->
Frame (current scope)
   ->
Episode (bounded activity)
   ->
History (facts)
   ->
Retained Experience (meaning)
   ->
Patterns (cross-episode connections)
```

Additional dynamics:

- Standing Intents generate future episodes
- Delegation creates new nodes and frames when needed

## Important Insight

- Purpose is durable and structural
- Scope is local and dynamic
- Episodes are what happened
- Meaning is derived later

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
