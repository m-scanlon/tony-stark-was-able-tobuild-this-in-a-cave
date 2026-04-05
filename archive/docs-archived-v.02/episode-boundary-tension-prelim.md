# Episode Boundary Tension (Prelim)

## Problem

The system clearly produces continuous activity, but we need a way to group that activity into meaningful units ("episodes").

However, not all activity has a natural or explicit ending.

- actors are long-lived and continuously capable of receiving stimulus
- some actors, such as user-facing actors, are effectively always on
- not all interactions are goal-driven or intent-based

This creates tension between:

- continuous operation, where there is no natural stop
- bounded understanding, where we need segments to reason, store, and recall

## What Is Clear

We do have bounded runtime activity, but the exact atomic turn shape is not fully locked.

What is currently stable is:

- actors are event-driven
- routed stimulus updates bounded episode state
- recall, inference, stimulus emission, interaction, and returned response handling may all occur inside one episode
- episodes group that activity into bounded spans

Older shorthand such as:

```text
stimulus -> recall -> inference -> act
```

was useful, but should not be treated as the final canonical runtime loop.

## Current Position

We do not enforce a strict definition of episode boundaries at this time.

Instead:

- actors operate continuously
- activity is recorded as bounded runtime turns
- episodes are treated as bounded groupings of that activity
- a actor episode is currently treated as a bounded grouping of one or more routed-stimulus turns

Current preliminary operating rule:

- if the user has not said anything for 30 minutes, treat that inactivity boundary as the end of the current episode

This should be treated as a practical heuristic for now, not a final definition of what an episode is.

## Guiding Principle

Bounded runtime turns exist; episode boundaries are inferred.

## Short Framing

The actor is event-driven.

The episode is a bounded grouping of routed-stimulus activity.

The remaining open question is how episode boundaries should eventually be defined beyond the current inactivity heuristic.
