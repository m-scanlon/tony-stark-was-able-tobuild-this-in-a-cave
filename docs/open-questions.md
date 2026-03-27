# Open Questions

## Purpose

This document holds the active open questions for the current data model.

It is meant to track unresolved edges without reintroducing older superseded assumptions.

## Stable Baseline

The current stable enough baseline is:

- episodes are the primary bounded unit of activity
- the frame is projected as purpose, interaction, recall, and available primitives
- the episode field is the scored entity/relationship layer active within the episode
- runtime commands and runtime artifacts belong to in-episode runtime execution
- learning is the write path from episodes into retained experience and structure
- retained experience is a family of retained artifact types rather than an understandings-only layer
- retained artifacts share an `anchor_set`
- `retained_trace` remains distinct from derived retained artifacts
- recall reads retained artifacts back through the scored episode field

## 1. Runtime Primitive Menu

The runtime-vs-retained boundary is now clearer, but the primitive menu is still open.

Questions:

- which runtime commands are first-class?
- which operations remain internal to runtime execution rather than kernel-callable?
- how fine-grained should command-specific runtime handling become?

## 2. Runtime Artifact Types

Runtime artifacts are now conceptually established, but not yet concretely typed.

Questions:

- what runtime artifact kinds exist?
- which are primitive-specific versus generic?
- what survives only for one step versus the whole episode?

## 3. Episode Field Scoring

The episode field is now the main scoring surface for recall, but the scoring math is still open.

Questions:

- how exactly should stimulus support be computed?
- how exactly should episode support decay or reinforce?
- how far should relational support propagate?
- what counts as the dominant connected slice of the field?

## 4. Recall Ranking

The recall bridge is now structurally clear, but ranking policy remains open.

Questions:

- how should anchor overlap be scored?
- how much should relationship overlap outweigh entity-only overlap?
- what are the admission thresholds into recall?
- how should partial matches compete against specific matches?

## 5. Trace Extraction

Learning now clearly preserves retained traces, but trace extraction is still underspecified.

Questions:

- what counts as one retained trace?
- how bounded should a trace be?
- how much natural language should the `happened` field contain?
- how much of the episode should be allowed into one trace?

## 6. Derived Artifact Formation

The retained artifact family is defined, but formation rules remain open.

Questions:

- when should a trace yield understanding?
- when should a trace yield salience?
- when should a trace yield tension?
- when should no derived artifact be written at all?

## 7. Provenance And Context

The current contracts allow:

- `source_trace_ids`
- `context_artifact_ids`

Questions:

- how much provenance should be stored?
- how strict should `context_artifact_ids` be?
- when is prior artifact influence strong enough to record?

## 8. Promotion Into Stable Structure

The current model distinguishes retention from canonical structure, but promotion rules remain open.

Questions:

- when does repeated retained experience harden into relied-upon structure?
- what should remain provisional indefinitely?
- how should conflicts between retained artifacts and structure be handled?

## 9. Episode Boundary Sensitivity

Episode boundaries are still operational rather than fully settled.

Questions:

- how much should episode closure affect learning?
- when should semantic shift split or segment an episode?
- should long-running episodes be internally segmented before learning?

## 10. Stark And Jarvis Role Boundaries

`Stark` and `Jarvis` are now the canonical paired names for the two major node roles.

Questions:

- what exact boundary should exist between user-facing meaning authority and structural authority?
- where should recall influence sit between Jarvis and Stark?
- what collaboration contract should be documented between them?

## Short Framing

The architecture now has a cleaner ontology.

Most remaining questions are no longer about what kinds of things exist.

They are about:

- primitive menus
- runtime artifact typing
- recall scoring
- trace extraction
- provenance
- promotion
