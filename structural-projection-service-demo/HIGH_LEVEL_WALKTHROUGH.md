# Structural Projection Demo Walkthrough

## What This Demo Is

This directory is a sandbox for one specific problem:

- take natural language
- preserve entities and relationships
- output bound structural fragments the runtime can later use for recall

It is not the final structural projection service.

It is a local harness for testing the pipeline shape.

## The Current Recommended Path

If you only look at one path in this demo, look at the dependency projection path:

```text
raw text
-> light preprocess
-> optional spelling repair
-> optional coreference
-> dependency parse
-> bounded structural projection
-> projected fragments
```

That path lives in:

- [dependency_projection_demo.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/dependency_projection_demo.py)
- [run_dependency_demo.sh](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/run_dependency_demo.sh)

This is the path that currently maps best to the retained-experience runtime model.

## Why This Exists

OpenIE was fast, but it produced bad relation candidates.

The dependency path is trying to fix that by doing something more controlled:

- use parsing to recover sentence structure
- use coref only to clarify references
- project a bounded set of relation shapes from that structure

So the goal is not:

- let a generic extractor invent the whole graph

The goal is:

- recover structure
- then project the structure we actually want

## The Main Pieces

### 1. Preprocess

File:

- [coref_resolve.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/coref_resolve.py)

This does light cleanup for casual user text:

- expands a few shorthand forms
- normalizes spacing

This is intentionally small.

## 2. Optional Repair

Files:

- [jamspell_repair_demo.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/jamspell_repair_demo.py)
- [run_jamspell_demo.sh](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/run_jamspell_demo.sh)

This uses JamSpell to repair mildly broken text before parsing.

It helps on inputs like:

- `I likee pi`

It is not reliable enough to treat as perfect repair.

In the dependency demo, repair is optional via `--repair`.

## 3. Optional Coreference

File:

- [coref_resolve.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/coref_resolve.py)

This uses `fastcoref` locally.

Its job is not to produce final structure.

Its job is to help the next stage by resolving vague references like:

- `it`
- `they`
- `that`

Important design choice:

- in the dependency path, coref is used as annotation
- we do not rewrite the sentence before parsing anymore

That preserves grammar while still giving the projector better referent information.

## 4. Dependency Parse

File:

- [dependency_projection_demo.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/dependency_projection_demo.py)

This uses spaCy to recover the sentence structure:

- tokens
- dependency links
- noun chunks
- clause structure

This is the structural substrate the projector reads from.

## 5. Projection

File:

- [dependency_projection_demo.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/dependency_projection_demo.py)

This is the custom layer that maps parse shapes into fragments like:

- `subject -> relation -> object`

Examples:

- `you -> help -> me`
- `you -> help_with -> terraform`
- `terraform -> is_a -> a new language`
- `terraform -> has_property -> difficult`

This is the most important part of the demo.

The parser gives us structure.

The projector turns that structure into the form the runtime actually wants.

## What The Output Means

The dependency demo prints several views:

- `INPUT`
- `PREPROCESSED`
- `REPAIRED`
- `COREFERENCE CLUSTERS`
- `TOKENS`
- `NOUN CHUNKS`
- `DEPENDENCY MATCHES`
- `PROJECTED FRAGMENTS`

The part that matters most for the runtime is:

- `PROJECTED FRAGMENTS`

Everything above that is there to help explain where those fragments came from.

## Current Recommended Commands

### Basic dependency path

```bash
./run_dependency_demo.sh --resolve-coref "Can u help me with terraform? Its a new language and its a little difficult"
```

### Dependency path with repair

```bash
./run_dependency_demo.sh --repair --resolve-coref "I likee pi"
```

### Repair only

```bash
./run_jamspell_demo.sh "i cnwant e smoothy"
```

## What To Ignore If You Are New

These files are still useful, but they are not the main path anymore:

- [StructuralProjectionDemo.java](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/src/main/java/StructuralProjectionDemo.java)
- [run.sh](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/run.sh)
- [run_hybrid.sh](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/run_hybrid.sh)
- [score_triples.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/score_triples.py)
- [run_hybrid_scored.sh](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/run_hybrid_scored.sh)

Those were part of the earlier OpenIE-first exploration.

They are still useful as historical comparison or baseline tooling, but they are not the recommended architecture.

## What This Demo Proved

At a high level, this demo proved:

- we can get bound entity/relationship structure locally
- a dependency-first pipeline is more controllable than OpenIE
- repair and coref help, but they are support layers
- the projector is the real bridge into the runtime memory model

## What This Demo Does Not Do Yet

It does not yet do:

- final canonical entity resolution
- final canonical relationship ontology
- episode field scoring
- retained-artifact lookup
- full runtime integration

So this demo stops at:

- local structural fragment generation

That is enough for now, because that was the missing layer recall needed.

## Suggested Reading Order

If someone wants to understand this quickly, read in this order:

1. [HIGH_LEVEL_WALKTHROUGH.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/HIGH_LEVEL_WALKTHROUGH.md)
2. [README.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/README.md)
3. [dependency_projection_demo.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/dependency_projection_demo.py)
4. [coref_resolve.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/coref_resolve.py)
5. [jamspell_repair_demo.py](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/jamspell_repair_demo.py)

## Short Framing

This demo is a local testbed for turning natural language into bound structural fragments.

The current recommended path is:

- preprocess
- optional repair
- optional coref
- dependency parse
- structural projection

The projector is the important part.
