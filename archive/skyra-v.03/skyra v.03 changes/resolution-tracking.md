# Resolution Tracking

## Threads and Exchanges

A thread is not an exchange. They are different things.

A **thread** is the full arc from the moment something enters the system to the moment something exits it. One intent. One origin event. Everything that happens in between — every exchange opened, every being consulted, every downstream branch — is the thread working itself out. The thread exists before the first exchange opens and closes when a response crosses the boundary outward.

An **exchange** is a single call between two beings. Exchanges are the mechanism. The thread is what they all serve.

A thread doesn't live on any one being. It travels — implicitly — as the `threadID` born at the origin event, passing through every exchange it spawns. The runtime does not need to track the thread explicitly. The threadID is the thread.

The exchange stack is the call stack. Exchanges open downstream, work, and close. Each close returns upward toward the being that opened it. The thread winds down the same way a call stack unwinds — each return carrying what was learned back up toward the origin.

The base case is when the last exchange closes and the response crosses the boundary outward. The thread is done.

## Relationship To Execution Surface

A thread originates from an external surface — Michael typing, speaking, or sending a signal from any device. Multiple surfaces belong to the same being. All of Michael's surfaces feed into the same threadID space. The thread is Michael's intent, not the device it came from.

See execution-surface.md.

## The Problem

When a cognitive being opens a downstream exchange to serve an upstream obligation, the runtime currently loses the thread on close.

`close-exchange` deletes the exchange and dispatch exits. The being that opened the downstream exchange never gets re-triggered. It never carries what it learned back up. The call stack unwinds to nothing instead of returning.

Two things need to be true for the thread to survive:

1. When a downstream exchange closes, the opener needs to be re-triggered with what was learned
2. The opener's present needs to show it the full picture — what resolved, what is still open, what it still owes upstream

## Resolution

**Resolution** is a per-exchange field that tracks whether the opener has received and acknowledged what they needed.

- `Resolved: false` — exchange is open, downstream being owes a response
- `Resolved: true` — opener closed with `~learned`, obligation fulfilled

Resolution lives on `ExchangeThread`, not on the channel. Two exchanges with the same peer are two distinct threads with independent resolution state. The relationship (channel) is shared. The resolution is isolated.

Only opener threads resolve. A thread where `IsOpener: false` is an upstream obligation — something the being owes back. The architecture makes this visible but does not resolve it automatically.

## ExchangeThread Changes

```
Resolved bool   — set true when opener calls close-exchange
Learned  string — the synthesis the opener carries out
```

`HasOpenExchanges()` counts only threads where `!Resolved`.

## close-exchange Syntax

```
skyra close-exchange ~with <being> ~learned <synthesis> ~expression-reference <start-end> | <reason>
```

- `~learned` — the synthesis the opener is taking away from this exchange
- `~expression-reference` — optional entries to carry as context

When close-exchange fires:
1. Thread marked resolved, learned stored
2. Dispatch re-triggers inference for the opener
3. Opener sees its full resolution graph and decides what to do next

## DerivePresent Structure

```
[identity block]

current resolutions
________________
<peer>
  status: unresolved | resolved
  learned: <synthesis>         ← only if resolved
  [thread content]             ← only if unresolved
  resolve: skyra close-exchange ~with <peer> ~learned <synthesis> | <reason>
           ← only shown if thread.IsOpener and unresolved

— your cognitive processes —

your processes
________________
[peers with no unresolved exchanges]
```

## Decision Space After Resolution

When a resolved thread appears in the present alongside open obligations, the being can:

- **Return upstream** — respond to the non-opener obligation, close the loop
- **Wait** — more downstream exchanges still unresolved
- **Cancel pending** — emit multiple close-exchanges to abandon open deliberations and return upstream with what it has
- **Open new exchanges** — new information from the resolved thread changes the plan

Multi-line emission makes all of these possible in a single inference turn.

## Open Implementation Questions

**Re-trigger plumbing** — when close-exchange fires, dispatch needs to know who to re-trigger. Currently dispatch doesn't track opener identity through the close-exchange call. Who holds the opener's name at the moment close-exchange processes? How does that name get surfaced to dispatch so it can re-run inference for the right being?

**Resolved thread cleanup** — if resolved threads stay in the map rather than being deleted, they accumulate. The present will eventually show a graveyard of old resolutions. When do resolved threads get cleaned up? After the being fires its next outward signal? After the upstream obligation closes? After some TTL? This needs a policy before the code changes.

**CloseExchange code change** — currently `CloseExchange` deletes the entry from the map entirely. For resolution tracking, closed threads need to stay in the map marked `Resolved: true` with `Learned` stored. That's a real structural change to `ExchangeMap`.

## What Does Not Change

- `threadID` is the key for each thread — two exchanges with the same peer are distinct
- Only the opener can close an exchange
- `IsOpener` is still set at thread creation time
- The channel itself (relationship) is shared across all threads with a peer
- Non-opener threads (upstream obligations) stay open until the caller closes from their side
