# Episode

## Working Clarification

An episode is a bounded unit of activity at a given scope.

There is no single global episode object.

When precision matters, the episode scope should be qualified.

## Scoped Episode Forms

Currently defined:

- node episode
- intent episode

Episodes are the primary unit of activity.

Intent is optional and exists when execution requires coordination.

Each episode captures activity relative to its scope.

## Node Episode

A node episode represents a node's bounded participation in system activity.

It records:

- the stimulus the node received
- the experience it was allowed to use
- the interaction it produced

A node episode is:

- bounded by the node's contract
- local to that node's perspective

A node episode records only that node's participation.

When intent is present, a node episode may be linked to one or more `intent_id`s.

## Intent Episode

An intent episode represents the lifecycle of an intent as it moves through the system.

It is composed of multiple node episodes linked by shared `intent_id`.

An intent episode provides a higher-level view of execution across nodes.

## Episode Structure

Every episode, regardless of scope, contains two distinct components: history and cognition.

History:

- records what occurred
- is append-only
- is objective
- is non-interpretive

Cognition:

- records internal processing
- is bounded by scope
- is separate from history

## Intent and Continuity

Continuity across the system is maintained through shared `intent_id`.

Intent links related node episodes when execution flows across nodes.

Intent tracks execution flow, not all activity.

Not all episodes are intent-driven.

Examples:

- user sharing context
- conversational interaction

## History

There is no single global history object.

History is reconstructed from episodes by grouping them and ordering them over time.

## Current Design Posture

The currently defined forms are `node episode` and `intent episode`.

The exact relationship between episode cognition and retained experience remains open.
