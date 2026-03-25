# Memory Open Questions (Prelim)

## Purpose

This document preserves the main open questions in the memory model after the recent convergence on episodes, understandings, and consolidation.

It exists so unresolved edges remain explicit and can be answered later without losing the current direction.

This is not the canonical model.

It is a holding document for questions that still need formal decisions.

## Current Baseline

The following points are currently stable enough to treat as the working model:

- a cycle is the execution atom
- a node episode is a bounded grouping of multiple cycles
- the active frame of an episode is organized as interaction, recall, and cognition
- the current cycle shape is `stimulus -> activation/recall -> cognition -> interact`
- the node contract currently contains `purpose`, `stimulus`, and `interact`
- an understanding is the smallest unit of retained experience
- an understanding currently contains `core`, `interpretation`, `strength`, and `activation`
- `core` is stimulus-grounded
- `core.entities` is required
- `core.relationships` is optional
- consolidation is the write path from episodes into structure and understandings

## Why This Doc Exists

Several newer ideas have started to sharpen the model, especially around primitives, frames, and artifacts.

Those ideas are promising, but they are not fully settled yet.

This document preserves them in a structured way so they remain available for future refinement.

## Open Area 1 - Primitive, Frame, and Artifact

One emerging idea is that cognition may be described in terms of:

- a frame, which defines what is in scope
- a primitive, which is the action applied over that frame
- an artifact, which is the output produced by that action

Working framing:

`primitive(frame) -> artifact`

This appears to be a strong abstraction, but several questions remain:

- should this become the general model for cognition, or only a local description of some cognition steps?
- which operations count as primitives?
- which outputs count as artifacts?
- are artifacts always intermediate, or can some be retained directly?

## Open Area 2 - Interpret as a Primitive

One current candidate is:

`interpret(core, interaction, recall, cognitive_artifacts) -> understanding_artifact`

This framing currently means:

- `core`: the stimulus-grounded structural substrate of the current target
- `interaction`: the current episode interaction trace in scope
- `recall`: already-retrieved understanding artifacts now present in the episode frame
- `cognitive_artifacts`: prior artifacts produced earlier in the same episode
- `understanding_artifact`: the meaning output produced by the interpret step

This direction appears promising, but the following remain open:

- should `interpret` formally consume `cognitive_artifacts` rather than `cognition`?
- is `understanding_artifact` only an intermediate runtime object, or is it already the same thing as a retained understanding?
- if they are different, what is the exact boundary between runtime artifact and retained record?
- should `interpret` be treated as one primitive among many, or as the central primitive of cognition?

## Open Area 3 - Cognitive Artifacts

The model now points toward the existence of prior artifacts produced during the same episode.

This is useful, but still vague.

The main open questions are:

- what exactly counts as a cognitive artifact?
- are cognitive artifacts only outputs of primitives, or can they include other episode-local structures?
- how long do cognitive artifacts live?
- are they scoped to a cycle, an episode, or both?
- can cognitive artifacts themselves be recalled later, or only consolidated through understanding?

## Open Area 4 - Understanding Artifact vs Retained Understanding

Another unresolved edge is the relationship between:

- an `understanding_artifact` produced during an episode
- a retained `understanding` stored in longer-term experience

Possible models include:

- they are the same object at different stages
- the runtime artifact is transformed before retention
- only some runtime artifacts ever become retained understandings

This distinction matters because it affects:

- consolidation
- retention thresholds
- comparison and matching
- provenance

## Open Area 5 - Same Understanding Resolution

Consolidation still needs a clear identity rule for understandings.

The key question is:

- when does a new candidate count as the same understanding as an existing one?

Possible signals may include:

- same core
- same interpretation
- similar interpretation
- same intent or situational context
- compatible prior activation and strength

Related open questions:

- can multiple conflicting understandings coexist on the same core?
- when should an understanding be refined versus split into a new understanding?
- how much natural-language variation should count as the same interpretation?

## Open Area 6 - Strength

`Strength` is required in the current understanding model, but its stored shape remains open.

The main questions are:

- is strength a bundled profile or a set of explicit fields?
- which fields are primary, such as confidence, salience, stability, reinforcement, or recency?
- how should conflicting evidence weaken strength?
- how should repeated use or repeated confirmation reinforce strength?

## Open Area 7 - Activation

`Activation` is also required, but its representation remains unsettled.

The main questions are:

- what exactly is stored in activation?
- is activation primarily semantic, structural, contextual, or all three?
- how are entities, relationships, and situations encoded in activation?
- how is activation queried during recall?
- how compressed can activation become before retrieval quality breaks down?

## Open Area 8 - Interaction Scope for Interpretation

The model currently treats `interaction` as an input to interpretation.

That is useful, but not fully specified.

Open questions:

- does `interaction` mean the full interaction trace of the episode so far?
- does it mean only the current local interaction window?
- what parts of interaction should be available to `interpret` directly?
- when does older interaction stop mattering for interpretation inside the same episode?

## Open Area 9 - Consolidation Policy

The overall consolidation flow is now clear, but several policy questions remain.

These include:

- what exact retention threshold should determine whether an understanding is written at all?
- how much provenance from the episode should remain attached to a retained understanding?
- should consolidation happen only on episode close, or can some partial consolidation happen earlier?
- how should multi-intent episodes affect activation and later recall?

## Open Area 10 - Episode Heuristic Sensitivity

Episodes are currently closed using a preliminary heuristic of 30 minutes of user inactivity.

This is operationally useful, but it still affects memory formation.

Open questions:

- how much should episode boundaries shape consolidation?
- when should semantic shift override time-based closure?
- should long-running episodes be internally segmented before consolidation?

## Short Framing

The architecture is close enough that most remaining questions are no longer about the overall shape of memory.

They are about the exact behavior of primitives, artifacts, matching, weighting, and recall.

That is progress, but these unresolved edges still matter and should remain explicit until they are formally answered.
