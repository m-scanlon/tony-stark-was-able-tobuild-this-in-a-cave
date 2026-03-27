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
- command-triggered interaction
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

This keeps the runtime simple and preserves the actual temporal shape of experience.

## Relationship To Async Commands

Asynchronous command execution may produce results later.

If those results are interaction-relevant, they should surface through the unified interaction log rather than forcing a new top-level frame section by default.

The internal runtime may still track outstanding command state separately.

The frame does not need to mirror every internal mechanism directly.

## Escalation Path

If one node begins handling too many distinct interaction responsibilities, the preferred next move is:

- revise the node contract
- decompose responsibility
- birth a new node if needed

This is preferable to prematurely fragmenting the frame structure.

So the current preference is:

- decompose overloaded nodes
- do not prematurely decompose the interaction section

## Current Design Posture

The strongest current claim is:

- interaction should stay unified and chronological by default

Multiple interaction methods do not by themselves justify multiple top-level frame sections.

## Short Framing

Keep interaction unified until there is a real operational reason not to.

Use typed chronological interaction events inside one interaction log.

If a node becomes overloaded, prefer node or contract decomposition over premature frame fragmentation.
