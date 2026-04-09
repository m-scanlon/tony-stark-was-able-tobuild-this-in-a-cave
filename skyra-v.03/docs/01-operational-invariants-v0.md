# Operational Invariants v0

## Purpose

This document translates the current `v.03` canon into runtime implications.

It locks the parts of the model that are already settled and leaves detailed
relationship-lifecycle mechanics to the dedicated lifecycle docs.

## 1. Everything That Communicates Is A Being

Humans, memory beings, genome beings, purpose-bound boundary beings,
differentiated beings, the creator, the external launcher, and host services
are all beings.

The distinction is phase and function, not ontology.

## 2. Nature Has Two Fields

The minimum shape of nature is:

- identity
- purpose

Both are public.

Nature is locked at the creation path rather than universally at registration.

## 3. Theory Of Mind Exposes Only The Public Slice

Before lived relationship experience exists, another being may inspect only:

- identity
- purpose

Private relationship interpretation does not belong in the visible slice.

## 4. The Runtime Has Three Creation Paths

The runtime must support:

1. genome seeding
2. runtime registration
3. differentiation

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
- instantiates required genome-seeded internal wiring
- routes created beings to the kernel

The kernel:

- admits beings into runtime participation
- verifies signed envelopes at the channel or participant boundary
- sheds the envelope before inserting expression into a being
- binds the incoming turn to the right source channel or participant context
- owns `id`, `origin`, and `trace_token` as system fields
- parses the incoming `raw` protocol string to extract routing target,
  emotional flags, and expression
- dispatches the signal to the target being without consulting a separately
  registered resolution method
- moves expression across existing relationships and into the receiving being's
  present
  without computing recipient sets
- updates relationship-graph edge weight mechanically on every signal pass
- adds the direct relationship to both beings' relationship hashmaps when
  edge weight crosses threshold
- removes the direct relationship from both beings' relationship hashmaps
  when edge weight decays below threshold
- reads emotional flags on expression structurally and may route copies
  according to fixed threshold rules

The kernel is not the trust boundary.

The kernel does not perform cognition.

## 7. The Genome Seeds Singletons, Templates, And Internal Wiring

`genome.skyra` must define:

- singleton genome beings
- templates for companion beings and purpose-bound boundary beings
- the internal relationships that must exist from birth

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

It makes first contact possible.

The current canonical base expression is locked in
[18-base-language-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/skyra-v.03/docs/18-base-language-v0.md).

## 10. One Ontological Relationship Exists Per Pair Of Beings

The pair is unordered.

Direction belongs to the expression turn.

If a second relationship appears necessary for the same pair, differentiation
is the answer.

## 11. Relationship Substrate Is Nested And Asymmetric

The relationship remains one real ontological thing per pair.

Operationally, the kernel holds the lived substrate as:

```text
HashMap<being_id, HashMap<peer_being_id, Stack<Exchange>>>
```

The outer key is the being.

The inner key is the peer being that the outer being is in relationship with.

The value is that being's own stack of exchanges with that peer.

The most recent exchange is on top.

The kernel peeks that stack to see whether the current exchange is open or
closed.

When a new exchange opens, the kernel pushes a new exchange onto that side's
stack.

Each side of the relationship holds its own stack independently.

There is no single shared symmetric exchange stack for both sides.

## 12. Trust Is A Private Relationship Value

Trust lives on the being's side of the relationship substrate.

It is private.

It is asymmetric.

It shapes how a being interprets what arrives from that relationship.

That is all it does.

## 13. Trust Uses A Fixed Interpretive Scale

The live trust scale is `1` to `100`.

`50` is neutral.

Below `50`, skepticism increases.

Above `50`, credibility increases.

`0` is reserved for no relationship yet.

It is not a live relationship value.

The currently locked origin values are:

- genome relationship -> `100`
- creator registers external being -> `90`
- differentiator births new being -> `50`
- unknown / pre-relationship -> `0`

Trust is static at origin in the current `v.03` model.

The movement algorithm is out of scope for now.

## 14. The Kernel Verifies And Moves Expression Locally

Expressions are verified at the channel or participant boundary through the
signed envelope.

After verification, the kernel sheds the envelope.

Only the expression enters the receiving being's present.

Whole envelope objects do not enter beings.

The kernel also owns the system fields on the transport signal:

- `id`
- `origin`
- `trace_token`

Those fields do not belong to the being.

`origin` is kernel-only and never exposed to beings.

The only outside-provided value on the signal is `raw`, the full protocol
string.

The kernel parses `raw` to extract:

- target being for routing
- emotional flags for structural routing decisions
- expression, which is passed through untouched to the being

There is no separately registered resolution method in the kernel.

Response behavior is baked into the being at birth by the being creator class.

The kernel just dispatches the signal.

The being knows what to do with the expression because of what it is.

The kernel binds each incoming turn to the right source channel or participant
context.

The kernel may move expression only across existing relationships and into the
receiving being's present.

It does not compute recipient sets.

It does not broadcast.

Trust and inference are different concerns.

Relationship interpretation happens in beings from their present.

Emotional routing is the narrow structural exception.

When an emotional threshold is crossed, the kernel may route a copy while the
original expression continues unchanged.

## 15. Shared Storage Does Not Define Discovery

There is no required central registry in the current model.

A runtime may include shared storage beings for artifacts that some beings
choose to persist centrally.

They do not own the kernel-held relationship substrate.

They do not own emergent relationship language.

They do not define who can discover whom.

## 16. The Genome Answers Who Knows Who From Birth

Genome-seeded internal relationships begin at trust `100`.

The internal language needed on those relationships is given at creation
through creation endowment.

That language does not wait to emerge.

## 17. Retained Experience Answers What Has Been Built Through Relating

Emergent external or novel relationship language lives in retained experience.

Callable language is recognized pattern retained through lived relating.

A being treats language as callable when cognition is confident enough to act.

That confidence is inferred, not stored as a protocol field.

## 18. Retained Artifacts Snapshot trust_at_formation

Every retained artifact carries `trust_at_formation`.

That field is the forming being's cognitive trust judgment at the moment the
artifact was created.

It is not copied mechanically from a relationship record.

See
[19-retained-artifact-family-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/skyra-v.03/docs/19-retained-artifact-family-v0.md).

## 19. Companion Beings Come From Templates

A being's personal hippocampus and personal experience store being are
companion beings instantiated from genome templates on the being's creation
path.

That includes runtime registration and differentiation.

Genome singletons receive their required companion beings on the genome path.

Present is not a separate companion being.

It is the being's operative reality.

## 20. Present And Thinking Stay Distinct

Observe means:

- a being operating from its present

Think means:

- a being relating to its memory beings

Present is not a passive queue.

Memory remains relational.

## 21. Expression Walk Is Fresh Firing At Each Hop

A being does not carry an expression forward mechanically.

What arrives enters inference.

The inference call may fire a new expression from the receiving being's
present.

Each hop is therefore a fresh expression turn.

What changes relationship is adjacent co-firing, not preservation of one fixed
message.

Relationship emergence is a kernel operation.

Every time a signal passes through the kernel, the kernel mechanically updates
edge weight on the relationship graph for the unordered pair.

This is Hebbian wiring in runtime form.

No inference is involved in that update.

When edge weight crosses the relationship-emergence threshold, the kernel adds
the direct relationship to both beings' relationship hashmaps.

When edge weight later decays below threshold, the kernel removes that direct
relationship from both hashmaps.

See
[21-expression-walk-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/skyra-v.03/docs/21-expression-walk-v0.md).

## 22. The Memory Split Still Holds

Retained experience is being-owned.

Retained artifacts live in the personal experience store of the being that
lived them.

Retained experience does not live in a shared central store.

## 23. Strain, Stress, Anger, And Conflict Form One Runtime Ladder

`strain` is internal pre-failure load in the being's present.

It is self-reported by the being paying the cost.

It is private by default.

When a being chooses to surface that load outward, it becomes `stress` or
`anger` on expression rather than traveling outward as `strain`.

Outward `stress` may set `trace_token` TTL for relationship-strength
accumulation.

`stress` and `anger` are structural values read by the kernel at every hop.

When `anger` crosses threshold, the kernel routes a copy of the expression to
the prefrontal while the original expression continues unchanged.

The prefrontal surfaces mismatch, attempts repair, and may propose a new path.

If repair fails and the prefrontal closes the exchange, the mismatch becomes
explicit `conflict`.

When explicit `conflict` is later carried outward, it may appear as an
optional integer field on expression.

The conflict being is not a standing monitor.

It detects deviation, escalates to the prefrontal, and suppresses the failing
signal.

Repeated failure over time may trigger differentiation.

See
[20-strain-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/skyra-v.03/docs/20-strain-v0.md)
and
[22-conflict-and-emotional-routing-v0.md](/Users/mikepersonal/tony-stark-was-able-tobuild-this-in-a-cave/skyra-v.03/docs/22-conflict-and-emotional-routing-v0.md).

## 24. Registration Is The Birth Path For Runtime Beings

Registration is not universal being birth.

It is the birth path for non-genome runtime beings.

Genome beings are seeded earlier through the genome path.

## 26. Names Are The Identity Surface

Being names are the only identity that flows through the runtime.

IDs do not appear in the protocol.

IDs do not appear in the present.

IDs do not flow through inference.

Inference produces a protocol string with a name.

The kernel routes on that name.

Name uniqueness is therefore a hard runtime invariant.

The kernel enforces uniqueness at every creation path.

No two beings may share a name at runtime.

If names were not unique, the kernel could not disambiguate routing and the
model would break.

IDs are internal kernel implementation detail at most.

They are not a first-class concept in the domain model.

## 27. The Canon Now Centers On

- being
- nature
- relationship
- present
- expression
- expression walk
- relationship strength
- trust
- genome-seeded internal wiring
- retained experience
