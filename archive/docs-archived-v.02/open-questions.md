# Open Questions

## Purpose

This document holds the active open questions for the current data model.

It is meant to track unresolved edges without reintroducing superseded command-first assumptions.

## Stable Baseline

The current stable enough baseline is:

- episodes are the primary bounded unit of activity
- the frame is projected as purpose, interaction, and recall
- runtime is stimulus-first
- the kernel routes emitted stimulus by contract lookup
- `ExecutionSurface` is the typed routing concept
- current execution-surface types are `actor` and `capability`
- actors emit and receive typed stimulus
- callable public surfaces should currently be modeled as one request stimulus plus one response envelope
- the response envelope currently requires `status` and `reason`
- learning is the write path from episodes into retained experience and structure
- retained experience is a family of retained artifact types rather than an understandings-only layer
- retained artifacts share an `anchor_set`
- `retained_trace` remains distinct from derived retained artifacts and remains recallable

## 1. Primitive Capability Semantics

The top-level primitive capability set is now clear, but the exact semantics remain open.

Questions:

- how strictly should each primitive capability be typed?
- which operations remain internal to runtime execution rather than kernel-routable?
- how fine-grained should primitive-specific runtime handling become?

## 2. Runtime Artifact Types

Runtime artifacts are now conceptually established, but not yet concretely typed.

Questions:

- what runtime artifact kinds exist?
- which are primitive-specific versus generic?
- what survives only for one step versus the whole episode?

## 3. Recall Contract Tuning

The recall contract is now clear, but its triggering and ranking policy remain open.

Questions:

- when should heavy inference emit a bounded recall request stimulus?
- how specific should recall queries become before they overfit?
- how should bounded ranking and admission thresholds be tuned?
- when should inference choose another primitive instead of recall?

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

## 10. Public Surface Granularity

The public request/response contract shape is now cleaner, but the final surface grain remains open.

Questions:

- how many callable surfaces should one actor usually expose?
- when should a actor split into another actor rather than grow another public surface?
- how much public abstraction is enough before the surface becomes too broad?

## 11. Stark And Jarvis Role Boundaries

`Stark` and `Jarvis` are now the canonical paired names for the two major actor roles.

Questions:

- what exact boundary should exist between user-facing meaning authority and structural authority?
- where should recall influence sit between Jarvis and Stark?
- what collaboration contract should be documented between them?

## Short Framing

Most remaining questions are no longer about whether the runtime is command-first or stimulus-first.

They are about:

- primitive semantics
- runtime artifact typing
- recall tuning
- trace extraction
- provenance
- promotion
- public surface granularity
