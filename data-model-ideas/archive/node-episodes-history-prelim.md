# Node Episodes, Intent Episodes & Reconstructed History (Prelim)

## Core Idea

Episodes are the primary unit of activity, but they exist at different scopes.

Currently defined:

- node episodes
- intent episodes

There is no single global episode object.

## Node Episodes

A node episode is a bounded grouping of one or more cycles of node participation.

Each node episode organizes the active frame of that participation into:

- interaction
- recall
- cognition

A node episode is always:

- bounded by the node's contract
- local to that node's perspective

A node does not record the full system story, only its participation across its included cycles.

When intent is present, a node episode may be associated with one or more `intent_id`s.

## Intent Episodes

An intent episode represents the lifecycle of an intent as it moves through the system.

It is composed of multiple node episodes linked by shared `intent_id`.

An intent episode provides a higher-level view of execution across nodes.

## Intent and Continuity

An `intent_id` represents a unit of intent moving through the system.

As work is delegated across nodes, the same `intent_id` is passed along.

A single `intent_id` may link multiple node episodes.

A single node episode may also participate in multiple `intent_id`s.

Intent tracks execution flow, not all activity.

Not all episodes are intent-driven.

Examples:

- user sharing context
- conversational interaction

## Episode Structure

Every episode, regardless of scope, organizes its active frame into three distinct parts.

### Interaction

Interaction captures exchange between the episode and the external world:

- stimulus
- interact
- external actions
- timestamps

Interaction is:

- append-only
- factual
- directional

### Recall

Recall captures what was activated into scope:

- anchors
- understandings

Recall is:

- selected rather than exhaustive
- bounded by scope
- activation-driven
- not a full memory dump

### Cognition

Cognition records internal processing:

- interpretation
- inference
- ambiguity handling
- decision formation
- candidate memory updates

Cognition is:

- internal
- bounded by scope
- uses recall to interpret interaction
- separate from interaction and recall

## Reconstructed History

History is not stored as a single mutable object.

History is not stored. It is reconstructed.

The system stores episodes as atomic records.

History emerges by:

- grouping episodes by scope and relation
- following cycles within each episode
- ordering them over time
- following `intent_id` when execution moves across nodes

There is no single global history object.

## Core Law

An intent flows through the system, and each node records its own bounded episode of handling that intent.

## Example

A user asks: "What's the weather today?" and this creates an `intent_id`.

A user-facing node receives the request and delegates, then records its node episode.

A weather node is called and returns data, then records its node episode.

The user-facing node formats and returns the answer, then completes its node episode.

These related node episodes form an intent episode.

Each node records only its own interaction, recall, and cognition across the cycles in that episode.

The larger history can be reconstructed by grouping and ordering the related episodes and their cycles.
