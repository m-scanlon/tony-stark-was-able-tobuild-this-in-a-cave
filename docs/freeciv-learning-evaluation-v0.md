# Freeciv Learning Evaluation v0

## Purpose

This document defines a practical test case for proving that `Skyra` learning
changes later behavior in a measurable way.

The main goal is not:

- to build a perfect game-playing benchmark
- to maximize game score in absolute terms
- to treat Freeciv as the final product surface

The main goal is:

- to create a repeatable environment where retained experience can be shown to
  change later recall, action selection, and outcomes

## Why Freeciv

Freeciv is a better early learning testbed than a live Discord server because
it provides:

- a stable ruleset
- repeated turn structure
- explicit score progression
- less social noise
- clearer outcome measures

It is also close enough to a Civ-style strategy game to make long-horizon
learning meaningful.

## Core Framing

The experiment should answer one question:

- does retained experience improve later behavior in a way that can be
  defended rather than merely felt

The intended proof surface is:

1. baseline run without retained memory
2. later runs with retained memory enabled
3. exported score-over-time timeline
4. observer notes correlated against game outcomes

The strongest outcome is not just:

- score went up

The stronger outcome is:

- retained artifacts changed recall
- recall changed later command choice
- command choice changed the game outcome

## Main Test Shape

The current working test shape is:

1. run Skyra in Freeciv with retained memory disabled
2. record game state, score progression, primitive calls, and results
3. run the same scenario again with retained memory enabled
4. repeat multiple runs until improvement begins to flatten
5. inspect both quantitative and qualitative changes

## Scenario Requirements

The scenario should be as fixed as practical.

Prefer:

- the same ruleset
- the same map size
- the same starting faction or nation
- the same seed if Freeciv exposes it cleanly
- the same turn budget
- the same node set
- the same primitive surface

The goal is to reduce noise so that learning effects are easier to isolate.

## Node Roles

The experiment should start with at least two nodes.

### 1. Player Node

The player node is the node that actually plays the game.

Its job is to:

- receive typed game stimulus
- perform recall when useful
- choose bounded game-facing commands through `act`
- react to command results
- continue play across turns

### 2. Observer Node

The observer node watches the same run but does not directly drive gameplay.

Its job is to:

- observe turn-by-turn state
- record mistakes, missed opportunities, and strategic patterns
- record score changes and major state transitions
- write to its own retained store rather than contaminating the player store

The observer node is useful because it separates:

- action
- from reflection

That makes later learning easier to inspect.

## Memory Modes

The test should run in two explicit modes.

### Baseline Mode

- retained memory disabled for the player node
- retained memory disabled for the observer node
- no prior learned state loaded

### Learned Mode

- retained memory enabled
- prior learned state loaded from earlier runs
- player and observer stores remain separate

If desired, a third mode may be useful later:

- player memory enabled
- observer memory disabled

That helps isolate whether the observer node is materially contributing to
improvement.

## What Should Be Exported

The experiment should export at least three timelines.

### 1. Game Timeline

The game timeline should include:

- turn number
- score at each turn
- major state milestones
- game-end outcome

Examples of milestones:

- city founded
- war started
- technology unlocked
- city lost
- unit loss spike
- diplomacy shift

### 2. Runtime Timeline

The runtime timeline should include:

- stimulus received
- node that received it
- primitive emitted
- command args
- result returned
- timestamps

This is needed to connect game outcomes back to actual runtime decisions.

### 3. Observer Notes Timeline

The observer node should emit timestamped notes such as:

- over-expanded before defense was stable
- ignored nearby threat
- delayed economic setup
- repeated a previous tactical mistake
- score dipped after low-value move sequence

These notes should be stored as their own timeline and, where useful, also be
eligible for later learning.

## Score-Over-Time Layer

The score-over-time layer is the simplest quantitative readout for the test.

The intent is to:

- export score by turn from the game
- compare curves across runs
- overlay observer notes and major decision points

If Freeciv score logging is available through `scorelog` and `scorefile`, this
should be the first implementation path.

The exact server-console command syntax should be verified during
implementation, but the important architectural point is already clear:

- the game can provide a historical score timeline
- Skyra can attach runtime and observer layers on top of it

## What Improvement Should Look Like

A good learning signal would include one or more of these:

- higher score at the same turn number
- fewer repeated strategic mistakes
- better early-game stability
- more efficient primitive use to achieve similar progress
- observer notes that disappear because the player changed behavior
- recalled retained artifacts that can be shown to have affected later choices

The best-case proof is:

1. baseline run makes mistake `M`
2. observer or player retains trace / understanding / salience / tension
3. later run recalls that retained artifact
4. later command choice changes
5. score curve or outcome improves

## Diminishing Returns Arc

One run with improvement is useful.

Several runs with a flattening curve are better.

The intended experiment arc is:

- baseline run
- repeated learned runs
- stop when improvements begin to plateau

Possible plateau indicators:

- median score gain over the last 3 runs is small
- repeated mistakes have mostly disappeared
- new retained artifacts stop changing behavior materially
- score curve shape stabilizes across runs

This matters because the aim is not just to prove that learning can help once.

The stronger aim is to observe:

- how learning compounds
- where it saturates
- whether retained experience becomes redundant or noisy

## Suggested Metrics

The first metric set should stay small.

Track:

- score at fixed turn checkpoints
- final score
- survived turns or completion state
- repeated mistake count
- primitive count per turn
- recall admissions that changed later action
- observer-note recurrence across runs

If a more direct objective exists in the chosen Freeciv scenario, add it later.

## Implementation Order

The implementation order should be:

1. get Freeciv state and score timeline into exportable form
2. route game state into typed stimulus
3. run one player node end to end
4. add observer node with its own retained store
5. add baseline vs learned run harness
6. add score-curve comparison and note overlays

Do not begin with advanced multi-node strategy.

The first useful proof is:

- one player node
- one observer node
- one exported timeline
- one visible improvement after learning

## Open Questions

The following should remain open until implementation pressure resolves them:

- what exact typed stimulus package should represent one game turn
- what exact `act` modalities are needed for game moves
- whether the observer node should write notes every turn or only on thresholded
  changes
- whether observer notes should become retained traces, understandings, or both
- what checkpoint interval should be used for score comparison

## Short Framing

Freeciv should be used as a controlled learning testbed for `Skyra`.

The experiment should compare:

- play without retained memory
- play with retained memory
- score over time
- observer notes over time

The main success condition is not just better play.

It is a defensible chain from:

- retained experience
- to later recall
- to changed runtime decisions
- to improved game outcomes
