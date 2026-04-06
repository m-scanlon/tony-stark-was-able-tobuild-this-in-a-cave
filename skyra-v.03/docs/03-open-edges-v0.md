# Open Edges v0

## Purpose

This document records the implementation questions that remain after the genome
bootstrap and signed-expression rewrite.

These are not contradictions in the canon.

They are the first concrete edges that still need design work.

## 1. Genome Syntax Is Not Yet Formal

Resolved baseline:

- `genome.skyra` defines singleton genome beings and templates
- the external launcher executes it
- the being factory reads it

What remains open:

- the exact Skyra syntax for singleton declarations
- the exact Skyra syntax for templates
- how dependency ordering is declared explicitly

## 2. Signed Envelope Shape Is Not Yet Formal

Resolved baseline:

- every expression carries a signed envelope
- the kernel verifies it against the sender's verification key
- the verification key is kernel-visible only

What remains open:

- the precise envelope schema
- the cryptographic algorithm choice
- how replay protection is expressed

## 3. Private Key Custody Is Not Yet Formal

Resolved baseline:

- genome beings receive creator-baked verification keys
- runtime and differentiated beings receive fresh verification keys from the
  being factory

What remains open:

- how private keys are delivered to runtime beings
- whether key rotation is allowed
- what happens when a private key is lost or compromised

## 4. The Registration Token Is Resolved Conceptually But Not Operationally

Resolved baseline:

- the registration token is the only unsigned entry path
- the kernel checks it before routing runtime registration to the being factory
- differentiated beings do not use it

What remains open:

- where the token lives operationally
- how the creator rotates it
- whether multiple registration tokens may coexist

## 5. RDS Record Shapes Need Formal Schemas

Resolved baseline:

- the RDS being stores being records
- being records include a `differentiatable` boolean that defaults to `true`
- the RDS being stores relationship records
- the RDS being stores language records

What remains open:

- exact table shapes
- exact indexes for unordered pairs
- how kernel-only fields are isolated from the theory-of-mind layer
- how and when `differentiatable` changes for a being over time

## 6. Differentiation Reorganization Still Needs Concrete Rules

Resolved baseline:

- differentiation reveals and reorganizes overloaded identity
- the being factory instantiates the revealed beings
- the kernel validates them

What remains open:

- how retained experience is redistributed concretely
- how relationship history is reattached
- how provenance of the split is preserved

The `differentiatable` flag reflects current capability limits, not permanent
ontological status. As the system matures, infrastructure and external beings
may become differentiatable. Revisit when the runtime reveals the need.

## 7. Companion Template Coverage Needs Explicit Declaration

Resolved baseline:

- every being needs companion beings
- runtime registration and differentiation instantiate them from templates

What remains open:

- whether genome singletons all share one companion template
- which singleton genome beings, if any, can omit certain companion beings
- how purpose-bound boundary beings are declared in templates

## 8. Relationship Startup Rules Need More Detail

Resolved baseline:

- some relationships exist because the genome seeded them
- other relationships begin when the kernel sees the language reader invoked on
  an unknown pair
- the kernel creates the relationship record and initial language record in
  `forming` state atomically at that point

What remains open:

- which genome relationships begin immediately callable
- which still have to earn specific callable language through use
- exactly what evidence promotes `forming` to `callable`

## 9. Retrieval And Reasoning Separation Still Needs A Hard Boundary

Resolved baseline:

- memory is a family of beings
- inference reasons over what memory beings return
- retrieval and reasoning remain separate concerns

What remains open:

- how strict that separation must be in implementation
- which memory beings are singleton, shared, personal, or differentiated

## 10. User Registration UX Is Still Open

Resolved baseline:

- runtime beings are born through kernel-guarded registration
- registration requires the registration token

What remains open:

- how the user presents the registration request
- how a runtime being receives its private key material
- how recovery or revocation is handled

## 11. Relationship As A Being Is Still Open

Resolved baseline:

- for now, a relationship is modeled as a record

What remains open:

- whether the runtime will later need relationships to become beings
- what new capabilities would justify that shift

## Short Framing

The canon is now coherent on bootstrap, verification, and creation paths.

The remaining work is implementation detail:

- exact genome syntax
- exact crypto envelopes
- exact storage shapes
- exact reorganization rules
