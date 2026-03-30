# Episode

## Working Clarification

An episode is a bounded unit of activity at a given scope.

There is no single global episode object.

When precision matters, the scope should be qualified.

## Scoped Episode Forms

The currently defined episode forms are:

- `node episode`
- `intent episode`

Episodes are the primary bounded unit of activity.

## Node Episode

A node episode is a bounded grouping of one or more cycles of one node's participation.

A node episode is:

- local to that node's perspective
- bounded by that node's contract
- a record of that node's participation rather than the whole system story

When intent is present, a node episode may be linked to one or more `intent_id`s.

## Intent Episode

An intent episode is a higher-level view of execution across nodes.

It is composed of related node episodes linked by shared `intent_id`.

Intent episodes are reconstructed across nodes rather than lived from one single perspective.

## Episode Frame

The active frame of an episode is projected from the episode.

The current frame layout is:

- `purpose`
- `interaction`
- `recall`
- `available_commands`

This is not the same thing as the whole episode.

The episode remains the bounded runtime container.

Interaction captures:

- incoming stimulus
- outgoing interact
- external actions
- timestamps

Interaction is factual, directional, and append-only.

Recall captures:

- retained artifacts activated into scope from retained experience
- a selected bounded set, not a full dump

Recall may include mixed retained artifact types such as:

- trace
- understanding
- salience
- tension

Available commands captures:

- the currently allowed runtime command surface for the episode
- the operations inference may choose next

The current first-class command examples are:

- `primitive recall`
- `primitive interact`

It is projected into frame because it bounds what the node may do at that moment.

## Episode Field

In addition to the frame, each active episode should maintain an episode field.

The episode field is:

- the scored entity/relationship layer of the episode
- the structural representation of what is currently active in that episode
- the main scoring surface used by recall

It sits just behind the current turn inside the episode.

It is not a separate theme object.

## Runtime Turns

The exact atomic runtime turn is not fully locked yet.

Earlier shorthand such as:

```text
stimulus -> recall -> cognition -> interact
```

was useful, but should not be treated as the final canonical runtime loop.

What is stable now is:

- nodes are event-driven
- episodes group bounded spans of runtime activity
- events may lead to recall, inference, command dispatch, interaction, and command-result writeback

A node episode contains one or more such bounded runtime turns.

## Intent And Continuity

`intent_id` provides continuity when work moves across nodes.

Intent tracks execution flow, not all activity.

Not all episodes are intent-driven.

Examples of non-intent-driven activity include:

- conversational interaction
- user sharing context

## Reconstructed History

There is no single global history object.

History is reconstructed from:

- episodes
- their included cycles
- their ordering over time
- their scope and relation
- shared `intent_id`

## Boundaries

Episode boundaries are still operationally heuristic rather than fully settled ontologically.

The current practical rule remains that inactivity may close an episode.

That rule is useful, but not final.

## Current Design Posture

The strongest current claims are:

- episodes are scoped and bounded
- node episodes are the primary local record of participation
- intent episodes are reconstructed across nodes
- the frame is projected from the episode as purpose, interaction, recall, and available commands
- the episode field is the scored structural layer active inside the episode

## Short Framing

An episode is a bounded scoped unit of activity.

Its frame is projected from episode state.

Its episode field maintains the scored structural context active during that episode.

History is reconstructed from episodes over time rather than stored as one mutable object.
