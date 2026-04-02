# Episode Boundary Tension (Prelim)

## Problem

The system clearly produces continuous activity, but we need a way to group that activity into meaningful units ("episodes").

However, not all activity has a natural or explicit ending.

- Nodes are long-lived and continuously capable of receiving stimulus
- Some nodes, such as user-facing nodes, are effectively always "on"
- Not all interactions are goal-driven or intent-based

This creates tension between:

- continuous operation, where there is no natural stop
- bounded understanding, where we need segments to reason, store, and recall

## What Is Clear

We do have bounded runtime activity, but the exact atomic turn shape is not fully locked.

What is currently stable is:

- nodes are event-driven
- routed events update bounded episode state
- recall, inference, command dispatch, interaction, and command-result writeback may all occur inside one episode
- episodes group that activity into bounded spans

Older shorthand such as:

```text
stimulus -> recall -> cognition -> act
```

was useful, but should not be treated as the final canonical runtime loop.

## Where the Tension Lies

The open question is not whether bounded runtime turns exist, but:

How do we group cycles into meaningful higher-level units?

These groupings are what we currently refer to as "episodes."

However:

- Episodes do not have a universally clean start/stop rule
- Intent-based boundaries work for tasks, but not for all activity
- Conversational or contextual activity lacks clear resolution points

## Current Position

We do not enforce a strict definition of episode boundaries at this time.

Instead:

- Nodes operate continuously
- Activity is recorded as bounded runtime events/turns
- Episodes are treated as bounded groupings of that activity
- A node episode is currently treated as a bounded grouping of one or more routed-event turns

Current preliminary operating rule:

- if the user has not said anything for 30 minutes, treat that inactivity boundary as the end of the current episode

This should be treated as a practical heuristic for now, not a final definition of what an episode is.

Boundaries are expected to emerge from:

- context shifts
- intent transitions
- salience
- time gaps

## Guiding Principle

Bounded runtime turns exist; episode boundaries are inferred.

## Current Resolution

For now:

- the node is event-driven
- the episode is a bounded grouping of routed-event activity
- inactivity remains the practical closure heuristic
The remaining open question is how episode boundaries should eventually be defined beyond the current inactivity heuristic.

## Goal

To eventually define a consistent way to segment activity into episodes without:

- over-fragmenting the system
- forcing artificial boundaries
- losing meaningful continuity

This will likely be driven by higher-level signals such as salience, memory relevance, and contextual coherence rather than strict mechanical rules.
