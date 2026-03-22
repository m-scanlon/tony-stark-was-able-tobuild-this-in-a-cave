# Node Episodes, Intent Episodes & History (Prelim)

## Core Idea

Episodes are the primary unit of activity, but they exist at different scopes.

Currently defined:

- node episodes
- intent episodes

There is no single global episode object.

## Node Episodes

A node episode is the atomic unit of execution history.

Each node episode represents:

- the stimulus the node received
- the experience it was allowed to use
- the interaction it produced

A node episode is always:

- bounded by the node's contract
- local to that node's perspective

A node does not record the full system story, only its participation.

## Intent Episodes

An intent episode represents the lifecycle of an intent as it moves through the system.

It is composed of multiple node episodes linked by shared `intent_id`.

An intent episode provides a higher-level view of execution across nodes.

## Intent and Continuity

An `intent_id` represents a unit of intent moving through the system.

As work is delegated across nodes, the same `intent_id` is passed along.

A node episode may be associated with multiple `intent_id`s.

Intent tracks execution flow, not all activity.

Not all episodes are intent-driven.

Examples:

- user sharing context
- conversational interaction

## Episode Structure

Every episode, regardless of scope, contains two distinct components.

### History

History records what occurred:

- stimulus
- interaction
- external actions
- state changes
- timestamps

History is:

- append-only
- objective
- non-interpretive

### Cognition

Cognition records internal processing:

- interpretation
- inference
- ambiguity handling
- decision formation
- emerging intent

Cognition is:

- internal
- bounded by scope
- separate from history

## History

History is not stored as a single mutable object.

History is not stored. It is reconstructed.

The system stores episodes as atomic records.

History emerges by:

- grouping episodes by scope and relation
- ordering them over time

There is no single global history object.

## Core Law

An intent flows through the system, and each node records its own bounded episode of handling that intent.

## Example

A user asks: "What's the weather today?" and this creates an `intent_id`.

A user-facing node receives the request and delegates, then records its node episode.

A weather node is called and returns data, then records its node episode.

The user-facing node formats and returns the answer, then completes its node episode.

These related node episodes form an intent episode.

Each node records only its own view.

The larger history can be reconstructed by grouping and ordering the related episodes.
