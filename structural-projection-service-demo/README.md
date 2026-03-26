# Structural Projection Service Demo

Small local demo for testing first-pass structural extraction from natural language.

If you want the cleanest high-level explanation first, start here:

- [HIGH_LEVEL_WALKTHROUGH.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/structural-projection-service-demo/HIGH_LEVEL_WALKTHROUGH.md)

Current backend:

- Stanford CoreNLP OpenIE

Experimental hybrid backend:

- `fastcoref` for local coreference resolution
- Stanford CoreNLP OpenIE after coref resolution

Experimental dependency projection backend:

- `fastcoref` for local coreference resolution
- spaCy dependency parsing
- spaCy `DependencyMatcher`
- small projected relation demo (`help_with`, `is_a`, `has_property`)

Experimental repair backend:

- light chat-text preprocessing
- JamSpell local spelling repair

This is not the full structural projection service.

It is a demo harness that shows:

- raw entity-relationship-entity extraction
- simple local-first entity normalization
- simple relationship normalization

before full canonical resolution and episode-field scoring.

## What It Does

Given a sentence, it prints extracted triples in this shape:

- `subject`
- `relation`
- `object`
- `confidence`

This is useful for evaluating whether OpenIE is a good first-pass binding layer for:

- entity extraction
- relationship extraction
- entity-relationship binding

## Usage

From this directory:

```bash
./run.sh "I am working on architecture with a new library and it is frustrating."
```

With theme bias:

```bash
./run.sh --theme software_work --theme-note "software architecture and tooling" "I am working on architecture with a new library and it is frustrating."
```

Hybrid coref + extraction:

```bash
./run_hybrid.sh --theme software_work --theme-note "software architecture and tooling" "Can u help me with terraform? Its a new language and its a little difficult"
```

Hybrid coref + raw entities/relationships:

```bash
./run_hybrid.sh --raw "Can u help me with terraform? Its a new language and its a little difficult"
```

Dependency graph + projected structural fragments:

```bash
./run_dependency_demo.sh --resolve-coref "Can u help me with terraform? Its a new language and its a little difficult"
```

Dependency graph with optional JamSpell repair before parse:

```bash
./run_dependency_demo.sh --repair --resolve-coref "I likee pi"
```

JamSpell repair before parsing:

```bash
./run_jamspell_demo.sh "i cnwant e smoothy"
```

Or directly through Maven:

```bash
mvn -q compile exec:java -Dexec.mainClass=StructuralProjectionDemo -Dexec.args="I am outside doing construction again today."
```

## Current Theme

The current implemented normalization theme is:

- `software_work`

It biases toward entities and relationships like:

- `self`
- `architecture`
- `library`
- `project`
- `toolset`
- `work`
- `working_on`
- `working_with`
- `using`
- `has_state`

`--theme-note` adds extra theme text to the scoring pass.

## Notes

- Output is expected to be messy.
- Coreference is still approximate.
- Theme normalization is simulated and local-first.
- This is a first-pass extractor/normalizer, not the final episode field updater.
- `run_hybrid.sh` requires a local Python virtualenv with `fastcoref` installed.
- `run_jamspell_demo.sh` requires `models/en.bin`.
