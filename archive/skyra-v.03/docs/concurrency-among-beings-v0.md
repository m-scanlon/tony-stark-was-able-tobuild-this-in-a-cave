# Concurrency Among Beings v0

## Purpose

This document locks the write algorithm for impulses that touch more than one
exchange.

It exists because the current single-write flow is incomplete.

A being needs a full record of every impulse it emits and every impulse it
receives, regardless of where those impulses go.

This document defines the write rule only.

It does not define new relationship-formation behavior.

## Problem

The current data flow records only one write per impulse:

- the receiver's exchange stack

That loses part of the lived history.

If a being fires an impulse while in active exchange with one peer and directs
that impulse at another peer, multiple exchanges were touched by that one turn.

All of them need the write.

## Core Rule

Every impulse gets written to every exchange it touches, from the perspective
of the being whose exchange it is.

That is the whole rule.

Do not collapse the write history to the arrival side only.

Do not collapse the write history to the directed side only.

Do not treat the originating exchange as irrelevant just because the impulse
was fired outward toward someone else.

## The Three-Write Case

When `A` emits an impulse at `B` while `A` is in active exchange with `C`,
three writes occur:

1. `A`'s exchange with `C`
2. `A`'s exchange with `B`
3. `B`'s exchange with `A`

Why:

- `A`'s exchange with `C` gets the write because that is the exchange the turn
  originated from
- `A`'s exchange with `B` gets the write because that is where the turn was
  directed
- `B`'s exchange with `A` gets the write because that is where the turn
  arrived

This is the canonical concurrent case.

## Signal Requirement And Algorithm Constraint

The write algorithm must be uniform across every hop:

1. Signal arrives at target being
2. Write to target's exchange with origin (arrival write)
3. Target runs inference
4. Inference produces a new signal — that signal becomes the next hop with its
   own origin and target

The algorithm is the same at every hop.

This means write `1` — the origin-side source-exchange write — does not fit
the uniform algorithm.

Write `1` is a record of what the origin fired outward while it was in
exchange with a third being.

That write is a side effect of what inference produced, not of what arrived.

It cannot happen at arrival time because the kernel does not yet know what
inference will produce.

It cannot happen post-inference without breaking the uniform algorithm.

## The Source Distinction

`source` is not the sending being.

`source` is the relationship the sending being fired from.

Example:

- `B` is in exchange with `A`
- `B` fires at `C` while in that exchange

The signal carries:

- `origin` — the being who fired (`B`)
- `target` — who it was directed at (`C`, parsed from the impulse)
- `source` — the relationship `B` fired from (`B`'s exchange with `A`)

`origin` and `source` are distinct.

`origin` identifies the being.

`source` identifies which relationship the being was operating from when it
fired.

Without `source`, write `1` is impossible.

## Open Question

Whether write `1` is strictly necessary, or whether the exchange history on
the arrival side is sufficient to reconstruct what happened, is not yet
resolved.

The three-write rule is the full honest record.

The question is whether the uniform hop algorithm and the full honest record
can coexist, or whether one must yield to the other.

## Write Algorithm

For each incoming impulse turn:

1. Resolve the `origin` being.
2. Resolve the `target` being.
3. Read `source_peer`, the peer whose exchange the origin was firing from.
4. Construct the set of touched exchanges:
   - `(owner = origin, peer = source_peer)`
   - `(owner = origin, peer = target)`
   - `(owner = target, peer = origin)`
5. De-duplicate overlapping entries if two slots refer to the same exchange.
6. Write the impulse once to each unique touched exchange.
7. Write it from the perspective of the being that owns that exchange.

The rule is therefore:

- one write per unique touched exchange

The common concurrent case yields three writes.

The deeper rule is not "always exactly three."

The deeper rule is "write to every exchange touched."

## Perspective Rule

The same impulse may appear in multiple exchange histories.

That is correct.

What changes across writes is not the fact that the same turn occurred.

What changes is the exchange perspective that owns the record.

So:

- the origin-side source exchange stores the turn as part of the origin being's
  lived exchange with its active peer
- the origin-side directed exchange stores the turn as part of the origin
  being's lived exchange with the being it addressed
- the receiver-side arrival exchange stores the turn as part of the receiver
  being's lived exchange with the sender

## Overlap Rule

The three-write case is not the only case.

Sometimes two conceptual slots collapse onto the same exchange.

Example:

- `A` fires at `B`
- `A`'s active source peer is also `B`

Then these two entries:

- `(owner = A, peer = source_peer)`
- `(owner = A, peer = target)`

refer to the same exchange.

So they collapse into one unique origin-side write, plus the arrival write on
`B`'s side.

In that case there are two writes, not three.

This is not an exception.

It still follows the core rule:

- write once to every unique exchange touched

## What This Does Not Change

This document does not define:

- relationship formation
- relationship emergence thresholds
- relationship decay
- graph maintenance

It defines exchange recording only.

The question here is not whether a relationship should exist.

The question here is:

- once a turn touched these exchanges, where must it be written

## Consequence

After this rule is applied, a being can reconstruct its own lived exchange
history more honestly.

It no longer loses emitted turns just because they were routed outward to a
different peer.

Arrival history remains intact.

Directed history remains intact.

Originating-exchange history remains intact.

## Short Framing

Single-write recording is incomplete.

Every impulse must be written to every exchange it touches, from the
perspective of the being whose exchange it is.

In the canonical concurrent case that means:

- origin/source exchange
- origin/target exchange
- target/origin exchange

The signal therefore needs `source_peer` so the kernel can perform the missing
origin-side write.
