# Dependency Pattern Coverage Map v0

This note maps the bounded dependency-pattern families that matter for the structural projection layer.

This is not a catalog of every possible parser label.

It is the coverage map for the pattern families that are relevant to:

- episode structural projection
- recall-oriented entity / relationship representation
- retained experience grounding

## Core Claim

The current design can cover these failures.

Why:

- the parser already gives bound local structure
- coreference can stay annotation-only
- the projector can grow pattern-family by pattern-family
- the missing behavior is mostly coverage, not architectural failure

So the projector problem is bounded.

## Coverage Status

`implemented`
- already in the dependency demo

`next`
- natural next additions that fit the current design directly

`later`
- still fits the design, but not urgent for `v1`

`defer`
- probably not first-pass projector work

## 1. Core Clause Patterns

### 1.1 Active verb with direct object

Parse shape:

- `verb`
- `nsubj`
- `dobj` / `obj`

Example:

- `I am building a software system`

Projection:

- `I -> build -> software system`

Status:

- `next`

### 1.2 Verb with prepositional object

Parse shape:

- `verb`
- `nsubj`
- `prep`
- `pobj`

Example:

- `you help me with terraform`

Projection:

- `you -> help_with -> terraform`

Status:

- `implemented`

### 1.3 Copula with nominal attribute

Parse shape:

- `be`
- `nsubj`
- `attr`

Example:

- `terraform is a language`

Projection:

- `terraform -> is_a -> language`

Status:

- `implemented`

### 1.4 Copula with adjectival complement

Parse shape:

- `be`
- `nsubj`
- `acomp`

Example:

- `terraform is difficult`

Projection:

- `terraform -> has_property -> difficult`

Status:

- `implemented`

### 1.5 Verb with clausal complement (`ccomp`)

Parse shape:

- outer `verb`
- outer `nsubj`
- child `ccomp`
- inner `nsubj`

Example:

- `the assistant is helping me design the architecture`

Projection:

- `assistant -> help -> me`
- `me -> design -> architecture`

Status:

- `implemented` in basic form

Notes:

- this is one of the biggest upgrades over the earlier shallow projector

### 1.6 Verb with open complement (`xcomp`)

Parse shape:

- outer `verb`
- complement `xcomp`
- no explicit inner subject

Example:

- `I want to build a system`

Projection:

- `I -> want -> build`
- `I -> build -> system`

Status:

- `implemented`

Why it matters:

- many intention / planning / action sentences use this pattern

### 1.7 Adverbial clause (`advcl`)

Parse shape:

- main predicate
- child `advcl`
- subordinate `nsubj`
- often subordinate `acomp` / object / prep structure

Example:

- `... because it is getting complicated`

Projection:

- `architecture -> get_complicated -> complicated`

Status:

- `next`

Notes:

- current output is usable but awkward
- likely needs refinement toward:
  - `architecture -> becoming -> complicated`
  - or `architecture -> has_property -> complicated`

### 1.8 Passive voice with agent

Parse shape:

- `nsubjpass`
- passive verb
- `agent` / `by`
- `pobj`

Example:

- `MiroFish was built by Guo Hangjiang`

Projection:

- `Guo Hangjiang -> build -> MiroFish`

Status:

- `implemented`

Why it matters:

- this was one of the biggest misses in the stress test

### 1.9 Ditransitive / indirect object

Parse shape:

- `verb`
- `nsubj`
- `iobj` or object-like second argument

Example:

- `she gave him a tool`

Projection:

- `she -> give -> tool`
- `she -> give_to -> him`

Status:

- `implemented`

### 1.10 Particle verbs / phrasal verbs

Parse shape:

- `verb`
- particle / `prt`

Example:

- `set up`
- `turn on`

Projection:

- `set_up`
- `turn_on`

Status:

- `implemented` in first-pass form

## 2. Identity And Nominal Structure

### 2.1 Appositive identity (`appos`)

Parse shape:

- head noun
- `appos`

Example:

- `the project's predecessor, BettaFish`
- `OASIS, an open-source project`

Projection:

- `predecessor -> identity -> BettaFish`
- `OASIS -> is_a -> project`

Status:

- `implemented` in first-pass form

Why it matters:

- this was another major miss in the stress test

### 2.2 Relative clause anchoring (`relcl`)

Parse shape:

- nominal antecedent
- child `relcl`
- relative pronoun subject or object

Example:

- `community that supports one million interactions`

Projection:

- `community -> support -> interactions`

Status:

- `implemented` in first-pass form

Why it matters:

- current demo now resolves simple `that`-subject relative clauses, but broader relative-clause handling still needs more coverage

### 2.3 Possessive / ownership / predecessor structure

Parse shape:

- possessive nominal
- `poss`
- head noun

Example:

- `the project's predecessor`

Projection:

- `project -> has_predecessor -> predecessor`
- or better once appos is added:
  - `project -> has_predecessor -> BettaFish`

Status:

- `next`

### 2.4 Compound noun structure

Parse shape:

- `compound`

Example:

- `software system`
- `simulation engine`
- `research community`

Projection:

- usually stays inside entity phrase construction

Status:

- `implemented` implicitly in noun phrase expansion

Notes:

- this is mostly entity phrase quality, not relation generation

### 2.5 Adjectival nominal modifier (`amod`)

Parse shape:

- noun head
- adjectival modifier

Example:

- `open-source project`
- `peer-reviewed research`

Projection:

- either keep inside phrase
- or emit optional attribute:
  - `project -> has_property -> open-source`

Status:

- `later`

### 2.6 Numeric / quantity structure (`nummod`, `quantmod`)

Parse shape:

- number modifying noun

Example:

- `one million agent interactions`

Projection:

- keep in phrase now
- optionally later:
  - `interactions -> quantity -> one million`

Status:

- `later`

## 3. Coordination Patterns

### 3.1 Conjoined predicates (`conj`)

Parse shape:

- first predicate as head
- later predicates attached via `conj`

Example:

- `topped ... and has attracted ...`

Projection:

- preserve each predicate separately with shared subject:
  - `MiroFish -> top -> list`
  - `MiroFish -> attract -> investment`

Status:

- `implemented` in simple inherited-subject form

### 3.2 Conjoined nominal predicates / states

Parse shape:

- copula head
- conjoined `be` / complement structures

Example:

- `It is a language and it is difficult`

Projection:

- `terraform -> is_a -> language`
- `terraform -> has_property -> difficult`

Status:

- `implemented`

### 3.3 Conjoined entities / objects

Parse shape:

- object plus `conj`

Example:

- `uses Python and Rust`

Projection:

- `entity -> use -> Python`
- `entity -> use -> Rust`

Status:

- `next`

## 4. Metadata And Event Modifiers

### 4.1 Temporal attachment

Parse shape:

- `prep in/on/by`
- temporal `pobj`
- sometimes `npadvmod` / `obl:tmod`-like behavior

Example:

- `in March 2026`
- `in late 2024`

Projection:

- either:
  - keep in phrase
  - or attach as event metadata:
    - `top -> at_time -> March 2026`

Status:

- `later`

Notes:

- important, but not necessary for first structural spine

### 4.2 Source / origin / provenance

Parse shape:

- verb + `from`

Example:

- `engine comes from OASIS`

Projection:

- `engine -> come_from -> OASIS`

Status:

- `implemented`

### 4.3 Location

Parse shape:

- `in` / `at` / `on`

Example:

- `student in China`

Projection:

- either entity phrase only
- or optional:
  - `student -> located_in -> China`

Status:

- `later`

### 4.4 Publication / venue / index placement

Parse shape:

- event verb + prepositional attachment

Example:

- `published in peer-reviewed research`
- `hit #1 on GitHub Trending`

Projection:

- event relation plus optional metadata:
  - `project -> publish_in -> research`
  - `BettaFish -> hit_on -> GitHub Trending`

Status:

- partly `implemented`
- phrase cleanup still `next`

## 5. Reference And Resolution Patterns

### 5.1 Pronoun coreference

Pattern:

- pronouns resolve to earlier anchor spans

Example:

- `it -> the architecture`
- `It -> MiroFish`

Projection role:

- referent identity used during projection
- original text preserved for parse

Status:

- `implemented`

### 5.2 Relative pronoun resolution

Pattern:

- `that`, `which`, `who`

Example:

- `community that supports ...`

Projection role:

- relative pronoun should inherit antecedent during relation projection

Status:

- `next`

### 5.3 Possessive pronoun / possessive nominal resolution

Pattern:

- `its`, `the project's`

Projection role:

- identity and relation anchoring

Status:

- `next`

## 6. Harder But Still Bounded Patterns

These still fit the current design, but are probably not first projector steps:

- reported speech / attribution
- conditional clauses
- negation-sensitive relation projection
- modality (`may`, `might`, `should`)
- comparison (`better than`, `more complex than`)
- ellipsis / fragment recovery

Status:

- `later`

## 7. Patterns To Avoid Overcommitting Early

These should not be first-pass projector goals:

- discourse-wide narrative summarization
- deep causal inference from syntax alone
- large latent relation ontologies
- global semantic rewriting before parse

Status:

- `defer`

## 8. Practical Coverage Summary

The most important `v1` projector set is:

- active transitive
- prepositional object
- copula attribute
- copula property
- `ccomp`
- `xcomp`
- passive + agent
- appositive identity
- relative clause anchoring
- conjoined objects

If those are covered well, the projector should handle a large share of user/runtime text in a structurally sane way.

## 9. Main Conclusion

The failures we saw are bounded pattern misses:

- passive voice
- appositives
- relative clauses
- phrase cleanup

Those are all compatible with the current design.

So the right move is not to replace the architecture.

The right move is to keep expanding the projector pattern family by family.
