# Actor And Episode Ownership v0

## Core Framing

The clean split is:

- the actor owns durable runtime machinery
- the episode owns bounded runtime state

This keeps process and state from collapsing into the same object.

## Actor Owns

The actor owns durable behavior, policies, and runtime machinery.

That currently includes:

- `actor_id`
- active and historical contract versions
- the actor process
- stimulus handling behavior
- primitive execution policy
- public request/response surfaces
- recall machinery
- frame assembly machinery
- the pointer to the active episode
- the `dependencyLedger`

## Episode Owns

The episode owns the bounded runtime state for one span of activity.

That currently includes:

- `episode_id`
- scope
- purpose snapshot
- interaction history
- recall
- derived episode-local artifacts
- timestamps and open/closed state

## Important Distinctions

### Contract vs Episode

The actor contract is durable.

The episode is instantiated runtime participation under that contract.

The contract bounds behavior.

The episode carries state.

### Actor vs Recall

The actor owns the recall machinery.

The episode owns the currently activated recalled artifacts.

### Actor vs Frame

The actor owns frame assembly behavior.

The frame itself is a projection from episode state.

The frame is not the durable owner of runtime truth.

## Short Framing

The actor is the durable runtime operator.

The episode is the bounded runtime state container.

The actor acts on the episode.

The frame is projected from the episode for inference.
