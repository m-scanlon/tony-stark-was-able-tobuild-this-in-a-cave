# Present Derivation — Exchange-Aware Network

The cognitive network section of the present must reflect the live exchange state of the being at signal time. It is not a flat list of peers. It is a dynamic view shaped by which exchanges are currently open.

## Two Tiers

### Open Exchange

A peer the being is currently in an open exchange with. The relationship is live. The callable language on that channel is available and shown.

```
strategy  [open]
  language: <callable language on the relationship>
```

`close-exchange` also appears here — seeded when the exchange opened, removed when it closes.

```
  close-exchange
    language: skyra close-exchange ~with <being> ~learned <synthesis> | <reason>
```

### No Open Exchange

A peer the being knows but is not currently in exchange with. The only path in is `start-exchange`. The start-exchange syntax is shown instead of callable language.

```
values
  start-exchange: skyra start-exchange ~with values ~about <string> ~because <sentence> ~say <expression> | <reason>
```

No callable language is shown here. The relationship language is earned inside the exchange, not before it.

## Why This Matters For Inference

The present tells the being exactly what it can do right now. Open exchanges have live callable language — the being knows how to speak to that peer. Peers with no open exchange show only the door in via start-exchange. The being cannot hallucinate callable language on a relationship it has not opened.

`close-exchange` appearing only under open exchanges is also load-bearing — it makes the resolution path visible only when there is something to resolve.

## Present Structure

```
[who you are]
[who you are in exchange with]
[the active exchange entries]

[your cognitive network]
________________
open exchanges:
  <peer>  [open]
    language: <callable language>
    close-exchange
      language: <resolution callable language>

not in exchange:
  <peer>
    start-exchange: skyra start-exchange ~with <peer> ~about <string> ~because <sentence> ~say <expression> | <reason>
```

## Implementation Note

`DerivePresent` checks exchange state per peer by inspecting the thread hashmap. If the peer's `HashMap<thread_id, Exchange>` has any entries, the exchange is open — render tier one. If the hashmap is empty, render tier two. No separate flag. No additional state. The hashmap is the source of truth.
