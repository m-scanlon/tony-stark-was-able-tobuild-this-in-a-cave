# Episode Field Scoring (v0)

## Core Framing

Retrieval should not begin from a fixed anchor object.

It should begin from a dynamically scored structural field inside the current episode.

That field is composed of:

- entities resolved in the episode
- relationships resolved in the episode

The purpose of the scoring process is to determine what parts of the current episode structure are most active right now for recall.

This document defines that scoring idea at a high level.

## The Two Concrete Inputs

At recall time, the score for any entity or relationship should be produced from two concrete things:

- support from the incoming stimulus
- support from the accumulated structure of the current episode

This means:

- early in an episode, the incoming stimulus dominates
- later in an episode, the accumulated episode structure helps stabilize and refine what is active

The score should not be based on the latest turn alone.

It should also not be based on older episode residue alone.

It is produced from the interaction between the current stimulus and the episode field.

## Episode Field

The episode field is the current structural representation of what the episode is about.

It is not a separate abstract "theme" object.

It is the accumulated set of currently relevant:

- entities
- relationships

with dynamic activation scores attached.

So:

- the frame projects the current purpose, interaction, recall, and available primitives
- the episode field is the scored entity/relationship layer that sits just behind the current turn inside the episode

## Scoring Unit

The scoring unit is not a phrase or a raw text span.

The scoring unit is:

- an entity
- or a relationship

Both should be scored directly.

Relationships matter because meaning often depends on configuration rather than isolated entities.

Examples:

- `outside`
- `self -> located_in -> outside`
- `self -> doing -> construction`

`outside` alone may not be enough to drive correct recall.

The connected relationships can change what that entity means in the current episode.

## Working Intuition

For any scored item `x`, where `x` can be either an entity or a relationship:

```text
activation(x) = f(
  stimulus_support(x),
  episode_support(x),
  relational_support(x)
)
```

This is only a framing equation.

The exact combination function remains open.

The important part is:

- `stimulus_support` comes from the current incoming turn
- `episode_support` comes from prior activation within the same episode
- `relational_support` comes from connected active structure

## Stimulus Support

Stimulus support answers:

"How directly is this entity or relationship supported by the current incoming stimulus?"

Possible sources:

- direct mention
- direct extraction
- high-confidence relation resolution
- immediate semantic fit to the current turn

## Episode Support

Episode support answers:

"How active has this entity or relationship already become in the current episode?"

This is what allows the episode to accumulate shape over time.

It gives continuity to the current conversation without requiring retrieval to start from scratch on every turn.

Operationally:

- recently active items should continue to carry weight
- repeated items should reinforce
- inactive items should fade

## Relational Support

Relational support answers:

"How much does the surrounding scored structure support this item being active right now?"

This allows activation to spread through connected structure.

Examples:

- if `construction` is highly active and `self -> doing -> construction` is present, that relationship should gain support
- if `self -> doing -> construction` is highly active, `construction` should gain support
- if `self -> located_in -> outside` and `self -> doing -> construction` are both active, `outside` may become active in a more specific way than by mention alone

This support should be bounded.

The system should not allow activation to spread arbitrarily far through the episode field.

## Update Behavior Across the Episode

The scoring process should run consistently on every recall step.

The same process applies whether the episode is:

- one turn old
- ten turns old

What changes is not the process but the available episode structure.

So:

- early episode = thin field, stimulus-heavy scoring
- later episode = richer field, stronger contextual stabilization

## Retrieval Shape

Retrieval should not query against a single unweighted key.

Instead:

1. resolve entities and relationships from the current stimulus
2. update the episode field scores
3. identify the highest-activation connected slice of the field
4. match retained artifacts against that scored slice

This means the effective lookup target is:

- not the raw stimulus
- not a fixed anchor object
- but the currently dominant relational slice of the episode field

## Matching to Retained Artifacts

Retained artifacts should be matched against the scored episode field.

The retained artifact should carry an `anchor_set`.

That anchor set is composed of structural references such as:

- entity references
- relationship references

The scored episode field then determines how well the current episode matches that artifact's anchor set.

This allows:

- generic artifact matches
- more specific relational matches
- partial matches when the current structure is incomplete

More specific relational matches should generally outrank weaker entity-only matches.

## Fault Tolerance

The model needs to remain useful even when the current structure is incomplete or slightly wrong.

So:

- exact match should not be required
- partial structural overlap should still be usable
- isolated weak matches should rank lower than connected relational matches

Fault tolerance should come from scored overlap, not from hand-written special cases.

## What This Avoids

This model avoids several things that currently seem too heavy or too early:

- a separate abstract theme object
- a rigid role taxonomy
- a giant stored pattern layer
- retrieval driven only by the latest stimulus
- retrieval driven only by static long-term weights

## Current Design Posture

At this stage, the strongest claim is:

- the episode accumulates a scored entity/relationship field over time
- retrieval is driven by the currently dominant relational slice of that field

The exact math for:

- decay
- reinforcement
- propagation depth
- artifact scoring

is still open.

## Short Framing

The current episode should maintain a dynamically scored structural field of entities and relationships.

Recall should be driven by the highest-activation connected portion of that field.

That activation should be produced from the interaction of:

- the incoming stimulus
- the accumulated structure of the current episode
- the local relational support inside that structure

This makes retrieval context-native without requiring a separate abstract theme layer or a rigid stored pattern system.
