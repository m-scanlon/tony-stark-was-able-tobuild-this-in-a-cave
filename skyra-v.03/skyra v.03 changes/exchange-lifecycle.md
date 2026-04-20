# Exchange Lifecycle

An exchange is not just a record of what was said. It has a lifecycle. It opens with intent and closes with a resolution artifact.

## The Rule

An exchange can only be closed by the being that opened it.

The opener set the intent. The opener judges whether the resolution is sufficient. Other beings can contribute — fire into the exchange, return work, share context — but they cannot close it. Only the opener can say: this obligation is resolved.

## Open With Intent

When a being opens an exchange it has a reason. That reason is the intent. It is not just a source field in the protocol. It is what the exchange is for. The intent is what the opener will eventually evaluate the resolution against.

## Close With Resolution Artifact

When the opener judges the resolution is sufficient, it closes the exchange by signaling `close-exchange`. Only the opener may do this — `close-exchange` verifies the caller holds the intent keying that exchange and drops the expression if not. The peer that received the work has no path to close. Only the sender can.

Current close-exchange syntax:

```
skyra close-exchange ~with <being> ~learned <synthesis> ~expression-reference <start-end> | <reason>
```

`~learned` is the synthesis the opener is taking away. It is what the opener carries back up to the next level.

**Close-exchange is not terminal. It is a return.**

When close-exchange fires:
1. The exchange is marked resolved. The learned synthesis is stored. The exchange stays in the map — it does not get deleted.
2. Dispatch re-triggers inference for the opener.
3. The opener's present shows the full picture — what resolved, what is still open, what it still owes upstream. The opener decides what to do next.

The call stack unwinds one level. The opener gets its context back. The thread is still alive.

**Eventual form — four artifacts:**

When memory beings exist (Phase 2), closing will require all four retained artifact family members:

- `trace.resolution` — what the exchange produced
- `understanding` — what it meant
- `salience` — what carried weight
- `tension` — what remained unresolved

The four artifacts will move to the experience store. The hippocampus indexes them for later retrieval. For now `~learned` carries the synthesis. The four-artifact requirement is Phase 2.

## Continue An Exchange

Between open and close there is continuation. A being explicitly signals that it is adding to an existing thread — not opening a new one, not closing the current one.

`continue-exchange` is a non-cognitive kernel-native being in the same family as start-exchange and close-exchange.

```
skyra continue-exchange ~with <being> ~thread <threadID> ~say <expression> | <reason>
```

- `~with` — the peer the exchange is with
- `~thread` — the threadID of the exchange being continued
- `~say` — what the being is adding to the thread

This makes continuation explicit. The present can clearly distinguish a new signal from a continuation of an existing thread. Without it, the being has to infer from context whether something is a continuation or something new — which is ambiguous under load.

The three operations on an exchange are now complete:

1. `start-exchange` — open with intent
2. `continue-exchange` — add to an existing thread
3. `close-exchange` — close with resolution artifact

## Relational Debt

The obligation is opened with intent. It is paid when the resolution artifact is returned. The artifact is the payment. An exchange receiving a return does not automatically close. The opener evaluates: does this resolve my intent? If yes, close. If no, the exchange stays open. More can be asked.

## Behavioral Consequence For Inference

A being that receives a return from a peer does not automatically close the exchange. It evaluates. The decision to close belongs to the opener. This is a meaningful constraint on how cognitive beings reason — closing is a judgment call, not a mechanical response to receiving something back.

## Where Closed Exchanges Go

The exchange stays in the map marked resolved with the learned synthesis stored. It does not get deleted. The opener's present shows it as resolved until the opener's own upstream obligation closes — at which point resolved threads can be cleaned up.

When memory beings exist (Phase 2), the four artifacts move to the experience store and the hippocampus indexes them for later retrieval through search-exchange or recall-exchange.
