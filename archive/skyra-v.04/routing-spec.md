# Routing Spec

## 1. Return routing

When a being closes an exchange and routes to someone outside it, the world should not deliver that message directly. Instead it should route the being back to its own perspective first — showing the parent exchange context — and let the being compose its response from there.

Currently: skyra finishes with philosopher, emits `michael <message> ~ref philosopher:0-5`, message goes straight to michael.

Should be: skyra finishes with philosopher → world closes philosopher exchange → world routes skyra to its own perspective with the michael exchange visible → skyra sees the full michael history, remembers where she left off, composes a response → that response goes to michael.

This is a function return. You don't jump to the caller mid-thought. You return to your own frame, orient, then speak.

The being's present at re-entry shows:
- The parent exchange history in full
- The ~ref context it carried back (injected under active exchanges as "context you returned with")
- No sender/message line — the being is just orienting, not responding to a prompt

The being then emits normally. Whatever it says goes to whoever it addresses.


## 2. Remember

A being should be able to ground itself in any of its exchange perspectives without sending a message. This is orientation — the being routes to its own perspective of an exchange, sees the state, and decides what to do next.

Protocol: `<> remember <peer>` — the being lands on its own perspective of the exchange with that peer. No message is sent to the peer. The target is always the being itself — it is just grounding itself in that exchange context. The being then emits normally from there.

This lets a being check in on a dormant exchange, recall context before opening a new one, or decide not to act after orienting.

"Remember" is the read. Every other emission is a write.


## 3. Departure visibility

When a being leaves its current exchange to go talk to someone else, the parent should see it in their present.

Currently: skyra goes to builder and michael's exchange just sits there with no indication skyra left.

Should be: when skyra opens an exchange with builder while michael is waiting, the world appends a system note to the michael exchange:

```
[skyra left to talk to builder]
```

This shows up in the exchange history. Michael can see where skyra went and roughly when she left. It's not a message from skyra — it's a world-level annotation. The being doesn't emit it, the world writes it when it opens the new exchange.

When skyra returns, a matching note:

```
[skyra returned from builder]
```

These are read-only annotations in the exchange — they don't trigger inference, they're just visible to the human reading the exchange.
