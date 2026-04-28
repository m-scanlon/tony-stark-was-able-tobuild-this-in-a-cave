# Inner Entity

## What This Is

Every entity that accumulates experience is split into two: an inner entity and an outer entity. The inner entity thinks. The outer entity speaks. The world only sees the outer entity's output.

---

## The Split

The outer entity does not deliberate. It receives its present — identity, purpose, thread context, exchange history, peers, the incoming message — and acts. It also receives a section called `inner-thoughts`, populated by the inner entity. The outer entity does not decide whether to use it. It is part of the frame, like identity. It is there.

The inner entity does not emit to the world. It receives the incoming message, reasons over it, and produces a read. That read becomes `inner-thoughts` on the outer entity's present. The inner entity never addresses a peer other than `remember`. It never appears in anyone else's exchange. The world enforces this through physics — the inner entity's output routes only to its outer entity, nowhere else.

One thinks. One speaks. The protocol boundary between them is internal to the entity.

---

## The Inner Entity's Present

The inner entity sees three things:

### 1. The incoming message

The message targeting the outer entity, before the outer entity sees it.

### 2. The last 10 exchange pairs

Each pair is:
- **what I thought** — the inner entity's read from that turn
- **what happened** — the outer entity's actual output

These pairs are the inner entity's working memory. They are always present — no query needed. The window slides. Older pairs fall out.

Over time, the inner entity builds a model of the gap between its reasoning and what actually gets said. It calibrates.

### 3. Remember

The inner entity has `remember` in its relationship list. It routes to `remember` with `<>` like any peer — same protocol, same interface. `remember` is an entity registered in the world that implements Entity. The inner entity calls it to pull deeper context beyond the 10-pair window.

The outer entity does not talk to `remember`. Memory access is the inner entity's job.

---

## The Outer Entity's Present

The outer entity's present is the same as it is today, plus one section:

```
inner-thoughts: <the inner entity's read for this turn>
```

This sits alongside identity, purpose, thread context, exchange history, peers, and the incoming message. The outer entity does not reason over it. It acts with it in frame.

---

## Execution Flow

1. A message targets an entity.
2. The world queues the message. The outer entity does not fire yet.
3. The inner entity fires. It sees: the incoming message, the last 10 exchange pairs, and can route to `remember` for deeper context.
4. The inner entity produces its read.
5. The queue releases. The outer entity fires. Its present includes the message, all normal context, and `inner-thoughts` from step 4.
6. The outer entity responds. The world routes the response normally.
7. The completed pair (inner thought + outer output) is stored in the sliding window.

The queue and the inner entity execution are world physics. The world fires the inner entity before the outer entity sees anything. Always.

---

## Retain

The shape of retain is not yet decided. The inner entity has recall now. How and when it stores what mattered will find its shape once the system is running.

---

## What This Is Not

- Not chain-of-thought prompting. The inner entity is a separate inference call with its own present, its own accumulating history, its own calibration loop.
- Not optional. Every entity that accumulates experience has an inner entity. An entity without one acts on the present with no weight from the past.
- Not visible. The world sees one entity. The split is internal. No peer can address the inner entity. No thread records its output except the sliding window internal to the entity.
