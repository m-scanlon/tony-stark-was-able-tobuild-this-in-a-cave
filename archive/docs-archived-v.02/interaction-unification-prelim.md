# Interaction Unification (Prelim)

## Core Framing

Interaction should remain unified by default.

The system may support multiple interaction methods and targets, but that does not mean the frame should immediately split them into separate top-level structures.

The default runtime posture should be:

- one interaction section
- one chronological interaction log
- typed interaction events inside that log

## Why This Matters

As the runtime grows, there will be pressure to split interaction into many separate channels:

- user-facing interaction
- device-facing interaction
- stimulus-triggered interaction
- capability-specific interaction

That may eventually be useful.

But doing it too early risks:

- fragmenting the frame
- losing temporal coherence
- overfitting the ontology before behavior is stable

## Current Direction

The current preferred model is:

- `interaction` remains one frame section
- the underlying interaction history remains one chronological log
- interaction events may be typed, but they remain part of one ordered history

## Relationship To Async Runtime Work

Asynchronous stimulus traversal may produce returned stimulus or response envelopes later.

If those results are interaction-relevant, they should surface through the unified interaction log rather than forcing a new top-level frame section by default.

The internal runtime may still track the `dependencyLedger` separately.

The frame does not need to mirror every internal mechanism directly.

## Escalation Path

If one actor begins handling too many distinct interaction responsibilities, the preferred next move is:

- revise the actor contract
- decompose responsibility
- birth a new actor if needed

This is preferable to prematurely fragmenting the frame structure.

## Current Design Posture

The strongest current claim is:

- interaction should stay unified and chronological by default

Multiple interaction methods do not by themselves justify multiple top-level frame sections.

## Short Framing

Keep interaction unified until there is a real operational reason not to.

Use typed chronological interaction events inside one interaction log.
