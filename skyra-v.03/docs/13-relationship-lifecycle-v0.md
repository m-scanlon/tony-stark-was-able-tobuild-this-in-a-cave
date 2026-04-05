# Relationship Lifecycle v0

## Purpose

This document locks startup, callability, and record semantics for earned peer
relationships in the current `v.03` model.

Genome-seeded relationships already exist on the relevant creation path.

This document is about relationships that begin through first encounter.

## Startup Boundary

By the time first encounter happens, both beings are already live.

Both already carry base language from their creation path.

Before the kernel sees a reach toward an unknown pair, no earned peer
relationship record exists between them.

Only bounded visibility exists:

- identity
- purpose

## The Base Expression

The base expression is the bootstrap every being can use at first contact.

One side reaches toward the other by invoking the language reader on an unknown
pair.

The kernel creates the relationship record for that unordered pair and the
initial language record in `forming` state atomically.

The base expression is then sent inside a signed envelope.

The kernel verifies the envelope before routing the turn.

Base language is intrinsic to the being at creation.

It is not relationship-owned.

It is not registry-backed.

## Forming

`forming` means the pair is known to the kernel, but mutual exchange has not
yet stabilized language.

Specific callable language begins building only after both beings exchange
expressions.

The transition into `forming` is implicit in the kernel.

## Callable

Once stabilized specific callable language fills the current-language slot of
the language record, the relationship becomes callable.

The language record is the proof that stabilization happened.

## Live States

The protocol needs:

- forming
- callable
- retired

Retired is the only terminal state.

Breakage returns the language record to `forming` while repair is brokered.

## Records

The RDS being stores two related record shapes for earned peer relationships.

### Language Record

- language_id
- relationship_id
- current_language
- status

`current_language` is nullable until stabilization.

### Relationship Record

- relationship_id
- unordered pair of beings
- aggregated language history
- created atomically with the initial language record

Language records reference relationship records by `relationship_id`.

## Short Framing

Genome relationships are seeded earlier.

Earned peer relationships enter `forming` when the kernel sees a reach toward
an unknown pair, creates both records atomically, and become callable when
stable language fills the current language slot of the language record.
