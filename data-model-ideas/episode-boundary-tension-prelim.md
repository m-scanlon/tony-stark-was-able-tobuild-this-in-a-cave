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

We have a well-defined atomic unit:

Cycle = stimulus -> cognition -> interaction

This represents a single unit of execution.

- Begins with stimulus
- Ends with interaction
- Does not require a broader "completion" definition

History and cognition are naturally expressed through these cycles.

## Where the Tension Lies

The open question is not how execution works, but:

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
- Activity is recorded as discrete cycles
- Episodes are treated as semantic groupings of cycles

Boundaries are expected to emerge from:

- context shifts
- intent transitions
- salience
- time gaps

## Guiding Principle

Cycles are defined; episodes are inferred.

## Open Question

It is not yet decided whether a node episode is:

- exactly one cycle
- or a grouping of multiple cycles

This remains open.

## Goal

To eventually define a consistent way to segment activity into episodes without:

- over-fragmenting the system
- forcing artificial boundaries
- losing meaningful continuity

This will likely be driven by higher-level signals such as salience, memory relevance, and contextual coherence rather than strict mechanical rules.
