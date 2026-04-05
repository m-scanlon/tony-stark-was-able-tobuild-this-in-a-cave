# Operational Invariants v0

## Purpose

This document translates the current `v.03` canon into runtime implications.

## 1. Everything That Communicates Is A Being

Humans, memory beings, genome beings, purpose-bound boundary beings,
differentiated beings, the creator, the external launcher, and host services
are all beings.

The distinction is phase, not ontology.

The creator, external launcher, and host services are pre-runtime beings.

## 2. Nature Has Three Fields

The minimum shape of nature is:

- identity
- purpose
- verification_key

Identity and purpose are public.

Verification key is kernel-visible only.

Nature is locked at the creation path rather than universally at registration.

## 3. Theory Of Mind Exposes Only The Public Slice

Before callable language exists, another being may inspect only:

- identity
- purpose

The verification key does not belong in the visible slice.

The kernel may read it directly from the RDS being.

Ordinary beings may not.

## 4. The Runtime Has Three Creation Paths

The runtime must support:

1. genome seeding
2. runtime registration
3. differentiation

Genome beings are created when the genome executes inside pre-raised host
services.

Runtime beings are created through the being factory after the kernel accepts a
registration request.

Differentiated beings are created through the being factory after the
differentiator identifies an overloaded identity that must be revealed and
reorganized.

## 5. Host Services And Hosted Beings Are Different Layers

The kernel service and being factory service are pre-runtime beings acting as
hosting layers.

The kernel being and being factory being are runtime beings hosted inside those
services.

That is a layering and phase distinction, not an ontological exception.

## 6. The Being Factory And Kernel Have Different Jobs

The being factory:

- reads singleton and template definitions from the genome
- instantiates beings and companion beings
- generates fresh verification keys on runtime registration and differentiation
- routes created beings to the kernel

The kernel:

- validates nature
- validates registration-token entry into unsigned registration
- verifies signed envelopes on expression
- records beings and relationship languages
- makes beings live
- routes expression

Neither replaces the other.

## 7. The Genome Script Seeds Singletons And Templates

`genome.skyra` must define:

- singleton genome beings
- templates for companion beings and purpose-bound boundary beings
- creator-baked verification keys for genome beings
- the registration token

The being factory must read this file rather than inventing these structures at
runtime.

## 8. Bootstrap Has A Fixed Minimal Order

The external launcher raises the host services first.

The genome then executes.

The being factory being reads the genome, creates genome beings in dependency
order, and routes each one through the kernel being.

Only after that is the runtime considered live.

## 9. Base Language Is Intrinsic At Creation

Base language is intrinsic to the being when it is created.

It is not relationship-owned.

It is not registry-backed.

It makes first contact possible.

## 10. Relationship Origin And Relationship Callability Are Different Questions

Some relationships exist because the genome seeded them at creation.

Other relationships begin at first encounter.

In both cases, specific callable language remains relationship-owned.

Callable status still depends on stable language in the language record.

## 11. One Relationship Exists Per Pair Of Beings

The pair is unordered.

Direction belongs to the expression turn.

If a second relationship appears necessary for the same pair, differentiation
is the answer.

## 12. Differentiation Reorganizes Identity

Differentiation is not ex nihilo creation.

It reveals and reorganizes overloaded identity into cleaner beings.

Operationally that means:

- the differentiator identifies the cleaner natures
- the being factory instantiates the revealed beings and their companions
- the kernel validates the resulting beings
- retained experience and relationship history are reorganized around the new
  boundaries

## 13. The Registration Token Is A Single Exception Path

Unsigned runtime registration is the only path that may reach the kernel
without an existing verification key.

That path requires the registration token baked into the genome.

No other flow may bypass signature verification.

## 14. Every Expression Must Carry A Signed Envelope

The kernel must verify the envelope against the sender's verification key before
routing.

No signature means rejection.

The wrong signature means rejection.

This is structural, not a bolt-on policy.

## 15. Shared Storage Lives On The RDS Being

The RDS being stores:

- being records
- relationship records
- language records

The theory-of-mind being reads public nature from the RDS being.

The kernel reads verification keys directly from the RDS being.

The language reader reads language records from the RDS being.

## 16. Language Records And Relationship Records Are Distinct

Language records store:

- language_id
- relationship_id
- current_language
- status

Relationship records store:

- relationship_id
- unordered pair
- aggregated language history

Relationship as a being remains an open edge.

## 17. Forming Is Kernel-Implicit

When the kernel sees the language reader invoked on an unknown pair, it creates
the relationship record for that unordered pair and the initial language record
in `forming` state atomically.

No explicit action is required from either being to create those records.

`forming` means one side has reached toward the other, but mutual exchange has
not yet stabilized language.

## 18. Companion Beings Come From Templates

A being's personal hippocampus, personal experience store being, and present
being are companion beings instantiated from genome templates on the being's
creation path.

That includes runtime registration and differentiation.

Genome singletons receive their required companion beings on the genome path.

## 19. Observe And Think Are Still Relation Flow

Observe means:

- a being relating to its present being

Think means:

- a being relating to its memory beings

No separate wire primitive is introduced.

## 20. The Memory Split Still Holds

Retained experience is being-owned.

Language history lives in shared RDS-backed relationship records.

The hippocampus is an access being, not the owner.

Retained experience does not live in the shared RDS being.

## 21. Boundary Beings Follow Purpose

The main sensory being is a singleton genome being.

Peripheral input beings and motor beings are created when a being's purpose
requires them.

## 22. Registration Is The Birth Path For Runtime Beings

Registration is not universal being birth anymore.

It is the birth path for non-genome runtime beings.

Genome beings are seeded earlier through the genome path.

## 23. The Canon Now Centers On

- being
- nature
- relationship
- expression
- genome path
- RDS-backed records
- signed verification

That replaces the older center of gravity built around `primordial`,
`beinghood`, and unsigned runtime assumptions.
