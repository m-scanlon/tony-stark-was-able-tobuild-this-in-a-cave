# Episode

## Working Clarification

An episode is a bounded unit of activity at a given scope.

There is no single global episode object.

When precision matters, the scope should be qualified.

## Scoped Episode Forms

The currently defined episode forms are:

- `actor episode`
- `intent episode`

Episodes are the primary bounded unit of activity.

## Actor Episode

A actor episode is a bounded grouping of one or more cycles of one actor's participation.

A actor episode is:

- local to that actor's perspective
- bounded by that actor's contract
- a record of that actor's participation rather than the whole system story

## Intent Episode

An intent episode is a higher-level continuity view across actors.

It is reconstructed across actor episodes rather than lived from one single perspective.

It should not currently be treated as the thing that enforces dependency or response routing.

## Episode Frame

The active frame of an episode is projected from the episode.

The current frame layout is:

- `purpose`
- `interaction`
- `recall`

Interaction captures:

- incoming stimulus
- outgoing participation
- external actions
- timestamps

Recall captures:

- retained artifacts activated into scope from retained experience
- a selected bounded set, not a full dump
- the admitted result of prior bounded recall requests

## Runtime Turns

The exact atomic runtime turn is not fully locked yet.

Earlier shorthand such as:

```text
stimulus -> recall -> inference -> act
```

was useful, but should not be treated as the final canonical runtime loop.

What is stable now is:

- actors are event-driven
- episodes group bounded spans of runtime activity
- heavy inference may emit bounded recall request stimulus during the episode
- events may lead to recall, inference, stimulus emission, interaction, and returned response handling

A actor episode contains one or more such bounded runtime turns.

## Reconstructed History

There is no single global history object.

History is reconstructed from:

- episodes
- their included cycles
- their ordering over time
- their scope and relation
- shared continuity where available

## Boundaries

Episode boundaries are still operationally heuristic rather than fully settled ontologically.

The current practical rule remains that inactivity may close an episode.

That rule is useful, but not final.

## Current Design Posture

The strongest current claims are:

- episodes are scoped and bounded
- actor episodes are the primary local record of participation
- intent episodes are reconstructed continuity views across actors
- the frame is projected from the episode as purpose, interaction, and recall
- recall is a contract driven by heavy inference rather than an episode-side scored field

## Short Framing

An episode is a bounded scoped unit of activity.

Its frame is projected from episode state.

Its recall section holds the retained artifacts currently brought into scope.
