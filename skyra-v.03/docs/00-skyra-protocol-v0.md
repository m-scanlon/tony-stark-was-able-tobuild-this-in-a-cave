# Skyra Protocol v0

## Status

This document is the source of truth for `skyra-v.03`.

It captures the canon after the genome / factory / kernel / signed-expression
rewrite.

## The Statement

Beings relate.

To others and to themselves. Through language. Via expression. Direction is a
turn-level fact, not relationship identity.

Memory is a being.

Thinking is a being relating to its memory being.

Observing is a being relating to its present being.

There is no ontological special case.

Pre-runtime beings make runtime life possible.

Once the runtime is live, it is all relating.

## Foundation

Truth is inferred.

Correctness is learned.

Failure is signal.

Consequences are feedback with weight.

Stability is earned, not designed.

## The Atomic Unit: The Being

Everything that communicates is a being.

A human, the creator, the external launcher, host services, the being factory
being, the kernel being, the RDS being, language beings, memory beings,
boundary beings, and differentiated beings are all beings.

The distinction is phase, not ontology.

The creator, external launcher, and host services are pre-runtime beings.

They exist before the runtime is live and make it possible.

### Nature

Nature has exactly three fields:

- identity
- purpose
- verification_key

Identity is public.

Purpose is public.

Verification key is kernel-visible only.

The theory-of-mind being does not expose the verification key to other beings.

Nature is locked at the creation path.

The creation path depends on what kind of being is coming into existence:

- genome beings are locked at genome seeding
- runtime beings are locked at factory creation and kernel validation
- differentiated beings are locked at differentiation path completion

### Creation Endowment

A being is created with:

- base language
- genome relationships required by its creation path
- companion beings required by its genome template

This is not a third ontological bucket.

These are creation-time facts, not a separate kind of thing.

### Relationships

Relationships are the connections beings live inside.

Two origins exist:

- genome relationships, seeded by the genome script or instantiated from genome
  templates
- earned relationships, formed through lived relating

Both are relationships.

A relationship is still the callable unit.

Specific callable language remains relationship-owned.

Exactly one specific callable language may be current on a relationship at a
time.

Old languages move into relationship history.

## One Relationship Per Pair

Between any two beings, there is exactly one relationship.

Always.

The pair is unordered.

Direction belongs to the current expression turn, not to relationship identity.

If a relationship appears to need a second identity, that is evidence that a
being boundary is wrong.

The answer is differentiation, not parallel relationships for the same pair.

## Differentiation

Differentiation is the reorganization path triggered when a being boundary is
wrong.

The differentiator does not create a being from nothing.

It reveals, splits, and reorganizes overloaded identity into cleaner beings.

The differentiator identifies the cleaner natures, routes the request through
the being factory, and reorganizes retained experience, relationship history,
and ongoing relationships around the revealed beings.

New verification keys are generated on the differentiation path.

The registration token is never used here.

## The Creation Paths

### 1. Genome Path

The creator writes `genome.skyra`.

The external launcher, a pre-runtime being, raises the host services first.

Those host services are pre-runtime beings that host the beings that come
online when the genome executes.

The genome then executes.

The being factory being reads it, creates genome beings in dependency order,
and routes each one through the kernel being for validation and live
instantiation.

The genome script defines:

- singleton genome beings
- genome templates
- creator-baked verification keys for genome beings
- the registration token used only for unsigned entry into registration

### 2. Runtime Registration Path

Registration is the birth path for runtime beings.

An unsigned registration request reaches the kernel with the registration
token.

The kernel validates the token and routes the request to the being factory.

The being factory creates the runtime being from the birth template,
instantiates its companion beings, generates a fresh verification key, and
routes the result back to the kernel.

The kernel validates nature, records the being, establishes its genome
relationships, and makes it live.

### 3. Differentiation Path

A signed runtime being identifies an overloaded boundary.

The differentiator initiates the reorganization through its own signed
relationship with the kernel and the being factory.

The being factory instantiates the revealed beings and their companion beings
from the differentiation template and generates fresh verification keys.

The kernel validates the resulting nature records and makes the reorganized
beings live.

Retained experience and relationship history are reorganized rather than
cloned blindly.

## First Encounter

First encounter happens between two already-live beings.

Before first encounter, the visible slice of another being is:

- identity
- purpose

The verification key is not visible through the registry reader.

If the other being is relevant, one side reaches toward the other by invoking
the language reader on the unknown pair.

The kernel creates the relationship record for that unordered pair and the
initial language record in `forming` state atomically.

If the other being is relevant, the first being then sends a base expression
inside a signed envelope.

The kernel verifies the signature, routes the expression, and the receiving
being takes it into its present being.

From there, language begins building through mutual exchange.

## Think

Thinking is a being relating to its memory being.

Memory is still a family of beings rather than one global surface.

Retrieval and reasoning remain separable concerns.

## Observe

Observing is a being relating to its present being.

The present being is a companion being created on the being's creation path.

It holds what has arrived and is waiting for attention.

## Language

Two language layers remain:

### Base Language

Base language is intrinsic to the being at creation.

It is not relationship-owned.

It is not registry-backed.

It makes first contact possible.

### Specific Callable Language

Specific callable language is relationship-owned.

It may emerge bilaterally or mostly through one-sided adaptation, depending on
response richness.

It becomes callable once stable and written into the current-language slot of
the language record.

## Security

Each expression carries a signed envelope.

The sending being signs the envelope with its private key.

The kernel verifies the signature against the sender's verification key in the
RDS being before routing the turn.

No signature means rejection.

The wrong signature means rejection.

The registration token is the only credential that works before a being has a
verification key.

It is used only for unsigned entry into runtime registration.

## The Genome Layer

The genome layer is what makes later relating possible.

It contains singleton genome beings such as:

- being factory
- kernel
- RDS being
- registry writer
- language writer
- language reader
- differentiator
- main sensory being
- prefrontal being
- the five specialized prefrontal beings

It also contains templates that define per-being companion beings such as:

- personal hippocampus
- personal experience store being
- present being

Purpose-bound boundary beings may also be created from genome templates when a
being's purpose requires them.

See [07-genome-beings-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/skyra-v.03/docs/07-genome-beings-v0.md) for the dedicated description.

## The Registries

## Shared Storage

The RDS being is the shared storage being for the runtime.

It stores:

- being records
- relationship records
- language records

### Being Records

Being records store:

- identity
- purpose
- verification_key

The theory-of-mind being reads from the RDS being and exposes only the public
slice:

- identity
- purpose

The kernel reads the verification key directly from the RDS being.

### Language Records

Language records store:

- language_id
- relationship_id
- current_language
- status

`current_language` is nullable until stabilization.

The allowed statuses are:

- forming
- callable
- retired

### Relationship Records

Relationship records store:

- relationship_id
- the unordered pair of being identities
- aggregated language history

For now, a relationship is modeled as a record rather than as a being.

## Two Kinds Of Memory

Retained experience is owned by the being that lived it.

Language history is stored in the shared RDS being.

Retained experience does not live there.

Those are different write paths and should not be collapsed.

## The CLI

```text
skyra relate <being> <expression>
```

The CLI is a motor channel.

It is not the caller.

The caller is the being inside the runtime that is signing and emitting the
turn through that channel.

## What Was Retired

The following language is retired in `v.03`:

- `primordial` -> replaced by `genome`
- `beinghood` as a formal third bucket
- `signature` as a stored bearer string
- `registration is being birth` as a universal rule
- differentiation as ex nihilo creation

Base language remains intrinsic to the being at creation.

Security now means:

- verification key in nature
- private key held by the being
- signed envelope on every expression

## What A Being Is In One Sentence

A being has a three-field nature, is created through the genome path, runtime
registration path, or differentiation path, inherits base language and genome
relationships at creation, and relates through signed expression into the
relationships it lives inside.
