# Dependency Projection Mapping v0

This note maps spaCy dependency parsing and `DependencyMatcher` onto the structural projection layer we need after coreference.

## Core Principle

The parser is not the retained-experience substrate.

The parser gives us a local syntactic tree:

- tokens
- heads
- dependency labels
- sentence boundaries

The structural projector reads that tree and emits bound fragments that fit our model:

- `entity -> relationship -> entity`
- `entity -> property`
- other small structural fragments if needed later

## What spaCy Gives Us

From the dependency parse:

- each token has exactly one head
- each token has a dependency label
- each token exposes children, subtree, and ancestors
- noun chunks give flat noun phrases rooted at a noun head

That gives us a stable local graph to project from.

## What `DependencyMatcher` Gives Us

`DependencyMatcher` does not extract memory structure directly.

It matches reusable dependency shapes inside the parse tree:

- verb with subject and object
- verb with subject, preposition, and object
- copula with attribute
- copula with adjectival complement
- later: clausal complements, adverbial clauses, conjunctions, control structures

So the matcher is a pattern engine over the parse tree, not the final projector.

## Mapping To Our Pipeline

The current target pipeline should be:

1. raw text arrives
2. light preprocessing
3. dependency parse over the original text
4. coreference over the original text
5. coreference stored as annotation, not rewrite
6. dependency patterns matched over the parse
7. structural projector maps matched shapes into bound fragments

That means:

- grammar comes from the parse
- referent identity comes from coreference
- final structural fragments come from the projector

## Why Coreference Must Be Annotation

Rewriting text before parsing can damage syntax.

Example:

- original: `assistant is helping me design`
- rewritten badly: `assistant is helping I design`

That makes the parse worse.

So coreference should annotate spans or tokens with referents, while the parser still sees the original surface text.

## Current Pattern Families

The current dependency demo supports:

- `verb + direct object`
  - `I -> build -> system`
- `verb + preposition + object`
  - `you -> help_with -> terraform`
- `copula + attr`
  - `terraform -> is_a -> language`
- `copula + acomp`
  - `terraform -> has_property -> difficult`

These are generic structural mappings, not example-specific hacks.

## Next Pattern Families

The next useful projector upgrades are:

- `ccomp` / embedded verb patterns
  - `assistant is helping me design the architecture`
- `advcl` patterns
  - `because it is getting complicated`
- conjunction-aware predicate projection
  - multiple linked predicates in one sentence

Those shapes should let the projector emit:

- `assistant -> help -> me`
- `me -> design -> architecture`
- `architecture -> get_complicated`

instead of flattening everything into only the easiest direct-object edges.

## Design Posture

- keep parser output structural
- keep coreference referential
- keep projection explicit
- do not let generic extractors invent the final relation space
- let the projector be the bridge from syntax into retained-experience structure
