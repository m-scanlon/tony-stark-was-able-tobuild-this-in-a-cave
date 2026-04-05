# Salvaged from Arch v1

## Purpose

This document preserves ideas from the older `docs/arch/v1` tree that still appear useful inside the newer memory architecture.

It is not a return to the old model.

It is a record of ideas worth carrying forward in updated terms.

## How To Read This

These are not adopted wholesale.

Some are strong invariants that still fit directly.

Others are later ideas that may become useful once the core memory model is more operational.

## Strong Salvage

### 1. Recall should be bounded

The old retrieve primitive was graph-shaped, but the bounding rules were good.

Recall should remain explicitly bounded by things such as:

- budget
- stopping conditions
- trust mode
- scope
- traceability

The important idea is that recall should terminate for a known reason rather than expand implicitly until context is full.

### 2. Write should be bounded and separate from recall

The old remember primitive carried the right separation:

- recall is a read path
- write is a separate path
- write should be scoped
- write should be explicit
- write should respect intent boundaries

This still fits the newer consolidation model.

### 3. Understanding is produced, not preloaded

One of the strongest old ideas was that understanding does not already exist in the active frame before interpretation finishes.

The frame provides what is available in scope.

The interpretive process produces understanding.

This maps well onto the newer direction of:

`primitive(frame) -> artifact`

and specifically:

`interpret(core, interaction, recall, cognitive_artifacts) -> understanding_artifact`

### 4. The model should not rely on hidden memory

The old runtime rule that the model should start fresh on each invocation is still strong.

What persists should be explicit:

- the active episode frame
- recalled artifacts
- retained understandings
- structural records

The model should not depend on hidden carry-over outside the architecture.

### 5. Truth is derived, not simply read

The old graph language is obsolete, but the temporal principle still holds:

- the system should avoid pretending there is one static current-truth field
- current truth should be derived from evolving records over time
- later evidence should not erase prior reality

This fits the current preference for evolution over overwrite.

### 6. Retrieval needs more than semantic similarity

The old importance-vector work still points at a real need.

Simple semantic similarity is not enough to drive high-quality recall.

The current model will likely still need additional weighting signals such as:

- salience
- recency
- stability
- affect or tone
- contextual fit

This does not require reusing the old exact vector schema, but it does preserve the core insight that recall needs weighted gating.

### 7. Observability around recall and action is valuable

The old observational-store docs separated:

- affect
- context evolution
- retrieval events
- system output

That exact store model may not survive, but the logging idea is still useful.

The newer system will likely benefit from recording:

- what triggered recall
- what entered recall
- what confidence changed
- what interaction followed
- what later consolidated

Without this, tuning activation and consolidation will be much harder.

### 8. Synthetic persona testing is still a strong fit

The old testing strategy using synthetic personas still fits the new model very well.

It should be adapted away from graph expectations and toward checks like:

- entities resolved correctly
- relationships resolved correctly
- understandings formed correctly
- consolidation avoids duplication
- repeated sessions reinforce rather than drift

This is one of the better old ideas and should likely survive almost unchanged in spirit.

## Good Later Ideas

### 1. Warm recall or context caches

The older context engine emphasized warm context rather than assembling everything from zero on demand.

That is still a useful direction for later:

- a warm recall cache
- immediate reinforcement of recently used memory
- batch updates for less urgent adjustments

This should be treated as an optimization layer, not part of the core ontology.

### 2. Background pattern passes

The older batch and cron ideas still make sense as later additions.

Some updates do not need to happen inline.

Later the system may use background passes for:

- reinforcement updates
- decay updates
- activation tuning
- pattern detection

### 3. Bounded self-improvement of retrieval policy

The old improvement-scope idea is still useful at a higher level.

Once retrieval, activation, and consolidation are stable enough, the system may need a bounded place to reason about improving those policies without directly mutating the trusted model.

## Leave Behind

The following should not be carried forward as the canonical architecture:

- property graph as the primary memory model
- graph actors and edges as the universal unit of memory
- old `perception / history / stimulus` terminology
- domain-agent-centric retrieval as the main memory framing
- token-by-token predictive retrieval as a current requirement

These ideas belong to the older architecture and should not be reintroduced accidentally through salvage work.

## Short Framing

What is most worth salvaging from Arch v1 is not the old data model.

It is the set of constraints and invariants around:

- bounded recall
- bounded write
- produced understanding
- explicit state
- temporal truth
- weighted recall
- observability
- evaluation

Those ideas still fit the newer architecture and may help make it operational.
