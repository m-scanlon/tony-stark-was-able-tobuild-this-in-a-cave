# DerivePresent

## What This Is

`DerivePresent` is owned by the entity, not the world. Each archetype owns its own lifecycle — building its present, firing its inner entity, resolving through its base case, parsing the response, managing its exchanges. The world routes. The entity derives.

---

## The World's Job

The world shrinks to:

1. Receive a relation
2. Find the target in the hashmap
3. Pass the relation to the entity's `DerivePresent`
4. Get back outbound relations
5. Route each one (back to step 1)

The world does not know what happened inside the entity. It dispatches and loops. Everything else is physics.

---

## The Entity's Job

`DerivePresent` on an entity takes a Relation. The entity:

1. Reads its state to build its own present (exchange history, active exchanges, peers, whatever it needs)
2. Fires its inner entity if it has one
3. Resolves through its base case (the invariant at the bottom of its recursion)
4. Parses the response
5. Returns outbound relations to the world

Different archetypes do different things inside `DerivePresent`:

- **llm** — fires inner entity, builds full present with thread context and inner-thoughts, routes to inference provider, parses response, manages exchanges, returns outbound relations
- **world** — runs physics to determine which child entity gets the relation, assembles context, recurses
- **human** — pushes present to the active device, waits for input, returns it as a relation
- **shell** — executes the command, returns output. No exchange management.
- **cli** — renders present as text, reads input, returns it as a relation
- **pipe** — writes present to stdin, reads stdout, returns it as a relation

Each archetype is special. The world doesn't need to know how.

---

## The Signature

```go
DerivePresent(r entity.Relation) []entity.Relation
```

Takes a relation. Returns outbound relations. The entity manages its own state internally.

---

## What Moves Out of the World

- Present building (identity, thread context, exchange history, peers, incoming message)
- Inner entity execution
- Base case resolution
- Response parsing and retry
- Exchange management (append, open, close)

---

## Routing Rules

The world's `DerivePresent` owns the routing rules — the constraints that govern how relations resolve. These live inside the archetype's implementation.

The current rules (implemented in the system world's `DerivePresent`):

1. **Self-reference drop** — an entity cannot target itself
2. **Auto-close on return** — if the entity is returning to whoever called it, close the detour exchange
3. **~ref departure close** — if the entity is addressing a new peer with a ~ref, close the current exchange
4. **Parent block** — if the target is the entity's parent in an active exchange the entity didn't open, block the message
5. **Departure visibility** — when an entity opens a new exchange while a parent exchange is waiting, annotate the parent exchange. When the entity returns, annotate again. These are world-level annotations — not messages from the entity, just visible bookkeeping in the exchange history
6. **Relationship enforcement** — the world checks the sender's relationship list before routing. If the target is not in the sender's relationships, the message is dropped.

Different archetypes can implement different rules. The rules live inside `DerivePresent`, not in a separate primitive.

---

## What Stays on the World

- Target resolution (hashmap lookup)
- The dispatch loop (route outbound relations to their targets)
- Grow (creating entities and registering them)
- Everything else lives inside the archetype's `DerivePresent`

---

## Future: Structured Objects on Relations

The Relation is currently `{ID, Origin, ThreadID, Impulse}` — everything is packed into Impulse as text. This works for messages. It may not work for shared objects that both parties in an exchange need to see, append to, and negotiate over. The Relation may need to grow to carry structured objects beyond Impulse. Not building this now — but the entity-owned lifecycle and world-defined physics are designed to support it without a rewrite.

---

## What This Is Not

Not a refactor for its own sake. The inner entity needs to fire inside present derivation, and that logic is archetype-specific. Mixing world routing with archetype-specific lifecycle in one method means every new archetype has to work around the world's assumptions. Moving derivation to the entity makes the inner entity a natural part of the entity's own lifecycle, not a special case the world has to handle.
