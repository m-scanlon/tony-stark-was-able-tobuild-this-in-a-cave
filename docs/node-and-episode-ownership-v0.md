# Node And Episode Ownership v0

## Core Framing

The clean split is:

- the node owns durable runtime machinery
- the episode owns bounded runtime state

This keeps process and state from collapsing into the same object.

## Node Owns

The node owns durable behavior, policies, and runtime machinery.

That currently includes:

- `node_id`
- active and historical `node_contract` versions
- the node process
- event handling behavior
- primitive definitions
- primitive execution policy
- capability bindings
- recall machinery
- frame assembly machinery
- the pointer to the active episode

The node owns the mechanisms that operate over runtime state.

## Episode Owns

The episode owns the bounded runtime state for one span of activity.

That currently includes:

- `episode_id`
- episode scope
- purpose snapshot for that episode
- `interaction_history`
- `recall`
- derived episode-local artifacts
- the available command surface for that episode
- episode timestamps and open/closed state

The episode owns what is true inside that bounded span of runtime participation.

## Important Distinctions

### Contract vs Episode

The node contract is durable.

The episode is instantiated runtime participation under that contract.

The contract bounds behavior.

The episode carries state.

### Node vs Recall

The node owns the recall machinery.

The episode owns the currently activated recalled artifacts.

### Node vs Frame

The node owns frame assembly behavior.

The frame itself is a projection from episode state.

The frame is not the durable owner of runtime truth.

## Programming Analogy

A useful approximation is:

- `node contract` = durable interface / class boundary
- `node` = long-lived runtime object or service
- `episode` = scoped execution context
- `frame` = projected input page for one inference step

This analogy is only approximate, but it helps preserve the main separation:

- node = operator
- episode = bounded working state

## Current Design Posture

The strongest current claims are:

- node and episode should remain separate objects
- the node owns process and mechanisms
- the episode owns bounded runtime state
- the frame is projected from episode state
- the contract bounds behavior but does not store episode state

## Short Framing

The node is the durable runtime operator.

The episode is the bounded runtime state container.

The node acts on the episode.

The frame is projected from the episode for inference.
