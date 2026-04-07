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

## 4. Live Status And Signing Dependency Are Explicitly Deferred

Resolved baseline:

- every expression carries a signed envelope
- verification keys exist as kernel-visible operational fields
- runtime and differentiated beings receive fresh verification keys from the
  being factory

Status:

- this ordering question is intentionally deferred in the current canon
- it is not treated as a contradiction in ontology closure

What remains deferred:

- whether a being may be admitted live before usable signing material is in
  hand
- how live admission and key delivery are ordered on each creation path
- what a newly admitted being may do before signing becomes available

## 5. The Registration Token Is Resolved Conceptually But Not Operationally

Resolved baseline:

- the registration token is the only unsigned entry path
- the kernel checks it before routing runtime registration to the being factory
- differentiated beings do not use it

What remains open:

- where the token lives operationally
- how the creator rotates it
- whether multiple registration tokens may coexist

## 6. Local Relationship Storage Needs Formal Schemas

Resolved baseline:

- being presence, relationship history, and exchange history are no longer
  modeled through a central language table
- each being keeps local records for the relationships and exchanges it lives
  through
- being records still include a `differentiatable` boolean that defaults to
  `true`
- before a direct relationship appears in a being's present, pre-relationship
  edge weight lives on the kernel-maintained relationship graph
- when edge weight crosses threshold, the kernel adds the direct relationship
  to both beings' local relationship hashmaps
- when edge weight decays below threshold, the kernel removes that direct
  relationship from both hashmaps

What remains open:

- exact local record shapes for relationships, exchange records, and present
  state
- how a single ontological `relationship_id` is represented across two local
  perspectives
- how exchange identity, ordering, and deduplication work when each hop is a
  fresh expression
- how and when `differentiatable` changes for a being over time

## 7. Differentiation Reorganization Still Needs Concrete Rules

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

## 8. Companion Template Coverage Needs Explicit Declaration

Resolved baseline:

- every being needs companion beings
- runtime registration and differentiation instantiate them from templates

What remains open:

- whether genome singletons all share one companion template
- which singleton genome beings, if any, can omit certain companion beings
- how purpose-bound boundary beings are declared in templates

## 9. Relationship Startup Rules Need More Detail

Resolved baseline:

- some relationships exist because the genome seeded them
- other direct relationships emerge through Hebbian wiring during expression
  walk
- every signal pass through the kernel mechanically updates edge weight on the
  relationship graph for the unordered pair
- `trace_token` is the kernel carrier used for that mechanical update
- no inference is involved in relationship emergence
- when edge weight crosses threshold, the kernel adds the direct relationship
  to both beings' local relationship hashmaps
- when edge weight decays below threshold, the kernel removes that direct
  relationship from both hashmaps
- each hop is a fresh expression fired from the receiving being's present
  after inference
- base language remains the first bootstrap expression for first contact

What remains open:

- how local relationship and exchange records are created on each side at first
  encounter
- how callable language is recognized locally without a formal language table

## 10. Retrieval And Reasoning Separation Still Needs A Hard Boundary

Resolved baseline:

- memory is a family of beings
- inference reasons over what memory beings return
- retrieval and reasoning remain separate concerns

What remains open:

- how strict that separation must be in implementation
- which memory beings are singleton, shared, personal, or differentiated

## 11. User Registration UX Is Still Open

Resolved baseline:

- runtime beings are born through kernel-guarded registration
- registration requires the registration token

What remains open:

- how the user presents the registration request
- how a runtime being receives its private key material
- how recovery or revocation is handled

## 12. Relationship As A Being Is Still Open

Resolved baseline:

- for now, a relationship is modeled as a record

What remains open:

- whether the runtime will later need relationships to become beings
- what new capabilities would justify that shift

## 13. Emotional Routing Thresholds Need Concrete Rules

Resolved baseline:

- the canonical transition ladder is locked in
  [22-conflict-and-emotional-routing-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/skyra-v.03/docs/22-conflict-and-emotional-routing-v0.md)
- `strain` is internal to the being
- when `strain` surfaces outward, it becomes `stress` or `anger` on expression
- the kernel reads emotional flags structurally on expression
- above threshold, the kernel routes an emotional copy to the prefrontal while
  the original expression continues unchanged
- the prefrontal surfaces mismatch, attempts repair, and closes the failing
  exchange if repair does not restore fit
- explicit `conflict` is post-break and follows failed repair
- the conflict being detects deviation, escalates to the prefrontal, and
  suppresses the failing signal
- the differentiator fires on repeated failure over time rather than one bad
  exchange

What remains open:

- the exact `anger` threshold
- whether `stress` has its own kernel threshold behavior beyond
  `trace_token` reach
- how the routed copy is marked so the prefrontal can distinguish it from the
  original walk
- how future emotional signal types are declared and routed

## Short Framing

The live docs are now coherent on bootstrap, signed-expression direction,
creation paths, kernel-maintained relationship emergence through Hebbian
wiring, and the conflict-emotional ladder.

The remaining work is implementation detail plus a few unresolved and deferred
edges:

- exact genome syntax
- exact crypto envelopes
- deferred live-status versus signing order
- exact storage shapes
- exact emotional-routing thresholds
- exact reorganization rules
