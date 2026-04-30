# Lens Implementation Spec

Handoff spec for the React Native lens. Everything your partner needs to build the client. Nothing more.

---

## The Rule

The lens holds one piece of state: the last present the runtime pushed. If you're adding a second piece of state, stop. Something is wrong.

---

## WebSocket Protocol

### Connection

```
ws://{host}:{port}/lens?being={beingId}
```

The lens connects on startup with the ID of the being it renders for. The runtime accepts the connection and begins pushing presents for that being.

Multiple lenses can connect for the same being. Each receives the same present.

### Inbound (runtime → lens)

Every message is a complete present. Never a patch, never a delta. The lens replaces its current present entirely on each message.

```typescript
type Present = {
  being: string
  sections: Section[]
}

type Section = {
  type: string
  data: Record<string, any>
}
```

### Outbound (lens → runtime)

When the user acts through an input interface, the lens sends a relation:

```typescript
type Relation = {
  id: string        // target being
  origin: string    // the being this lens renders for
  threadId: string   // current thread, from the present
  impulse: string   // what the user typed/did
}
```

That's the only thing the lens ever sends. One shape. No variants.

### Lifecycle

| Event | Behavior |
|---|---|
| Socket opens | Lens sends nothing. Waits for first present. |
| Present arrives | Lens replaces current state, re-renders. |
| User submits input | Lens sends a relation. Does not optimistically update. Waits for next present. |
| Socket drops | Lens shows disconnected state. Attempts reconnect with backoff. |
| Reconnect succeeds | Runtime pushes current present. Lens is caught up. |

No heartbeat. No ping/pong beyond what WebSocket provides natively. If the socket is open, the lens is live.

---

## Section Types and Data Shapes

These are the section types the runtime will push at launch. Each one maps to exactly one component.

### `identity`

The being's self-description.

```typescript
type IdentityData = {
  name: string
  identity: string
  purpose: string
}
```

Render: name prominent, identity and purpose as secondary text. Static display — no interaction.

### `thread`

The current thread and its active exchanges.

```typescript
type ThreadData = {
  id: string
  about: string
  because: string
  exchanges: Exchange[]
}

type Exchange = {
  peer: string
  status: "current" | "waiting" | "you opened" | "they opened"
  entries: number
}
```

Render: thread topic at top, exchanges as a list. Status determines visual treatment:
- `current` — highlighted, this is the active conversation
- `waiting` — dimmed, the being opened this but is talking to someone else
- `you opened` / `they opened` — directional indicator for who initiated

### `exchange`

The message history for the current exchange.

```typescript
type ExchangeData = {
  peer: string
  messages: Message[]
}

type Message = {
  origin: string
  content: string
}
```

Render: message list. Align/style based on whether `origin` matches the being or the peer. New messages appear at the bottom. Auto-scroll to latest.

### `peers`

Who the being can address.

```typescript
type PeersData = {
  available: string[]
}
```

Render: list of names. Tapping a peer name could pre-fill the input with that name as the target — that's the only interaction, and it's optional.

### `input`

The input surface.

```typescript
type InputData = {
  peer: string       // current exchange peer, default target
  threadId: string   // current thread
  being: string      // who this lens renders for (the origin)
}
```

Render: text input field. On submit, send a relation:

```json
{
  "id": "{peer}",
  "origin": "{being}",
  "threadId": "{threadId}",
  "impulse": "{user text}"
}
```

Clear the input after sending. Do not echo the message locally — wait for the next present to show it in the exchange.

If the user types a different being name at the start of the impulse (e.g., `builder check this`), the lens sends `id: "builder"` instead of the default peer. The runtime handles validation.

### `topology`

The world graph.

```typescript
type TopologyData = {
  beings: TopologyBeing[]
  edges: TopologyEdge[]
}

type TopologyBeing = {
  id: string
  name: string
  medium: string
}

type TopologyEdge = {
  from: string
  to: string
  threadId: string
  status: "active" | "inactive"
}
```

Render: node-edge graph. Beings are nodes, edges are active exchanges. This is read-only. Layout is up to the implementation — force-directed, radial, whatever works on the surface.

---

## Component Registry

```typescript
type SectionComponent = React.FC<{
  data: any
  send: (relation: Relation) => void
}>

const registry: Record<string, SectionComponent> = {
  identity: IdentitySection,
  thread: ThreadSection,
  exchange: ExchangeSection,
  peers: PeersSection,
  input: InputSection,
  topology: TopologySection,
}
```

Every component receives:
- `data` — the section's data object, typed per section
- `send` — function to send a relation back through the socket

Most components ignore `send`. Only `input` (and optionally `peers`) uses it.

If the present contains a section type not in the registry, skip it. Don't crash, don't warn. The runtime may push new section types before the lens has components for them.

---

## The App

```typescript
function Lens() {
  const [present, setPresent] = useState<Present | null>(null)
  const ws = useRef<WebSocket | null>(null)

  useEffect(() => {
    const socket = new WebSocket(`ws://${HOST}/lens?being=${BEING}`)
    socket.onmessage = (e) => setPresent(JSON.parse(e.data))
    ws.current = socket
    return () => socket.close()
  }, [])

  const send = useCallback((relation: Relation) => {
    ws.current?.send(JSON.stringify(relation))
  }, [])

  if (!present) return <Loading />

  return (
    <ScrollView>
      {present.sections.map((section, i) => {
        const Component = registry[section.type]
        if (!Component) return null
        return <Component key={`${section.type}-${i}`} data={section.data} send={send} />
      })}
    </ScrollView>
  )
}
```

That is the app. If it grows significantly beyond this, revisit the architecture.

---

## What Not to Build

| Don't | Why |
|---|---|
| Client-side routing | There's one screen. The present determines what's on it. |
| State management (Redux, Zustand, etc.) | One `useState` for the present. That's it. |
| API layer / fetch calls | The runtime pushes. The lens never requests. |
| Optimistic updates | The runtime is the source of truth. Wait for the next present. |
| Message persistence / local cache | The present is ephemeral. If the socket drops, reconnect and receive current state. |
| Auth (for now) | Single-user, local runtime. Auth is a later problem. |
| Animations between presents | The present replaces wholesale. Animating transitions between states adds client logic that fights the model. Revisit only if it's jarring. |

---

## Multi-Device

Same app, different component registries. The present JSON is identical regardless of surface.

| Surface | Registry adjustments |
|---|---|
| Phone | Single section visible at a time. Swipe between sections. Input always anchored at bottom. |
| Laptop | All sections visible. Input at bottom of exchange section. Topology gets more space. |
| TV | Drop `input`. Read-only. Large type. |
| Watch | Only render `exchange` (last message) and `input` (voice or quick reply). |

The runtime doesn't know which surface is connected. It pushes the same present. The lens decides what to render.

Start with one surface. Phone or laptop. Multi-device is the same app with a different layout — not a different architecture.

---

## Contract

**Runtime guarantees to the lens:**
- Every WebSocket message is a complete, valid `Present` JSON object
- Sections array order is intentional — render top to bottom
- A new present is pushed after every relation the runtime processes
- On reconnect, the runtime pushes the current present immediately

**Lens guarantees to the runtime:**
- Every WebSocket message is a valid `Relation` JSON object
- The lens never sends unprompted — only in response to user action
- The lens never caches or replays old presents
- The lens renders whatever sections it has components for and silently skips the rest
