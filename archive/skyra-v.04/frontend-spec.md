# Frontend Spec

## What This Is

The frontend is not a wrapper around the runtime. It is the surface where two Logos meet. The human and the being relate through the same medium — the UI is the exchange.

---

## Two Layers

### 1. Interaction Layer — Talking to a Being

The surface where a human and a being exchange. Not a chat interface. Three owned regions:

- **User space** — freeform, belongs to the human. A journal. Theirs.
- **Being space** — the being's internal state made visible. What threads are open, what it's holding, what it's uncertain about. Not a chat bubble — actual state.
- **Exchange zone** — the middle. Shared artifacts, negotiated meaning. Something lands here only when both sides have touched it.

Either side can invite the other into their region. Either side can decline. Sovereignty and invitation — the trust model made visible.

Both sides post independently. Not request/response. Two streams meeting on a canvas.

---

### 2. Management Layer — The Cockpit

The view above individual exchanges. Navigate the full topology of worlds and beings.

---

## Feature Set

### World Management

- **Name your world** — a world has an identity, not just an ID
- **Create worlds** — spawn child worlds from the current world
- **Navigate worlds** — click into a child world, go up to the parent
- **Default topology** — show the default implementation on first load (the genome bootstrap)
- **Set a topology** — define the shape of a world before deploying it

### Being Management

- **Create beings** — name, identity, purpose
- **View beings** — see all beings in the current world
- **Deploy beings and worlds** — push a topology live

### Talking

- **Talk to a being** — enter the interaction layer for any being in the current world
- **Talk to a world** — worlds are Logos too; they can be addressed directly

### Threads

- **Track threads** — see all open threads across the world
- **Thread statistics** — open, closed, duration, hop count, beings involved
- **View thread history** — read the full exchange of any closed thread

### Exchange Data

- **View exchange history** — all exchanges for a being, all exchanges in a world
- **Browse old exchanges** — searchable, filterable by being, thread, date

### Topology View

- **Visual topology** — graph of worlds and beings, edges showing relationships and active threads
- **Metadata overlay** — being created, world created, exchange count, thread count
- **World nesting** — parent/child world structure visible in the graph

---

## Navigation Model

```
root world
  └── child world A
        └── being 1
        └── being 2
  └── child world B
        └── being 3
```

Click a world → enter it. Breadcrumb shows depth. Up arrow → parent world. Click a being → enter interaction layer.

---

## What the Frontend Is Not

Not a dashboard bolted onto a backend. The UI state IS the exchange state. The topology view IS the LogosMap rendered. Thread stats ARE the exchange map rendered. Nothing is duplicated — the runtime is the source of truth and the frontend makes it visible.
