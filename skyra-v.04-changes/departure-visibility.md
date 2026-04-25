# Departure Visibility

## What This Is

When a being leaves its current exchange to talk to someone else, the parent should see it. When the being returns, the parent should see that too. This is world-level bookkeeping — the world annotates what happened while a being was away.

---

## The Problem

Currently: skyra goes to talk to builder, and michael's exchange just sits there. Michael's present shows the exchange frozen at the last message, with no context about what happened. If skyra talked to three beings before coming back, michael has no idea.

---

## The Fix

When a being opens a new exchange while a parent exchange is waiting, the world appends a system annotation to the waiting exchange:

```
[skyra left to talk to builder]
```

When the being returns:

```
[skyra returned from builder]
```

These are world-level annotations, not messages from the being. The being doesn't emit them — the world writes them when it opens or closes a detour exchange. They show up in the exchange history as read-only entries. They don't trigger inference. They're just visible to whoever reads the exchange.

---

## What the parent sees

If skyra has been gone for three exchanges before coming back, michael's exchange shows the trail:

```
  [0] michael: can you investigate this?
  [1] skyra: let me check with a few people.
  [skyra left to talk to philosopher]
  [skyra left to talk to builder]
  [skyra returned from philosopher]
  [skyra returned from builder]
  [2] skyra: here's what I found...
```

The annotations are interleaved chronologically. Michael can see where skyra went, in what order, and when she came back. The exchange history is honest about what happened.

---

## What this changes in the runtime

- `src/primitives/world/world.go` — when opening a new exchange while another is active for the same being, append a departure annotation to the waiting exchange. On close/return, append a return annotation.
- `src/primitives/exchange/exchange.go` — annotations are stored as relations with a system origin (e.g. `origin: "world"`) so they render in history but don't trigger routing.

---

## What this is not

Not related to the internal self, memory, or re-orientation. This is pure world bookkeeping — making the exchange history complete. Nothing else in the system depends on it and it depends on nothing else.