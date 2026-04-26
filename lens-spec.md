# Lens Spec

## What This Is

The lens is a blank surface that renders present data pushed from the runtime. It holds no state, no logic, no routing. It receives JSON and renders components. The runtime is the source of truth. The lens is glass.

This replaces the frontend-spec's framing. That spec described a "surface where two Logos meet." The lens model is simpler: the runtime pushes, the lens renders. The interaction model (user space, being space, exchange zone) still holds — but those are sections of the present, not sections of the app.

---

## Architecture

### The runtime pushes

When `DerivePresent` runs for a being, the result is a structured JSON object. The runtime pushes that object over a WebSocket to every lens registered for that being. The lens does not request data. It does not poll. It receives.

The state of the lens is the last present that was pushed to it.

### The lens renders

The lens is a React Native app with a component registry. When a present arrives, the lens walks the JSON, resolves each section to a registered component, and renders. No client-side state management. No routing library. No Redux, no context providers.

If the lens is managing state, something is wrong.

### The interface takes input

The lens exposes interfaces — text input, touch, gestures. These are intake surfaces. When the user acts through an interface, the lens sends a relation back to the runtime. The runtime processes it, derives a new present, pushes it back. The loop is: interface in, present out.

---

## Present Structure

The present is a JSON object with decomposable sections. Each section maps to a component in the lens registry.

```json
{
  "being": "skyra",
  "sections": [
    {
      "type": "identity",
      "data": {
        "name": "skyra",
        "identity": "I hold the world together.",
        "purpose": "I think, respond, and relate on behalf of the system."
      }
    },
    {
      "type": "thread",
      "data": {
        "id": "a1b2c3",
        "about": "memory architecture",
        "exchanges": [
          {
            "peer": "michael",
            "status": "current",
            "entries": 4
          },
          {
            "peer": "builder",
            "status": "you opened",
            "entries": 2
          }
        ]
      }
    },
    {
      "type": "exchange",
      "data": {
        "peer": "michael",
        "messages": [
          {"origin": "michael", "content": "how should memory work"},
          {"origin": "skyra", "content": "memory is a being, not a store"}
        ]
      }
    },
    {
      "type": "peers",
      "data": {
        "available": ["michael", "builder", "skeptic", "bash"]
      }
    }
  ]
}
```

The lens doesn't know what these sections mean. It knows how to render `identity`, `thread`, `exchange`, `peers` because those components are in its registry. New section types require a new component in the registry, not new logic in the lens.

---

## Component Registry

The lens ships with a small set of primitive components:

| Component | Renders |
|---|---|
| `identity` | Being name, identity, purpose |
| `thread` | Thread metadata, active exchanges |
| `exchange` | Message history between two beings |
| `peers` | Available peers to address |
| `input` | Text input interface — sends relations back to runtime |
| `topology` | Graph of beings and relationships in the world |

Each component is a pure function: JSON in, rendered output. No side effects, no local state, no API calls.

New components are added to the registry as the present grows new section types. The lens binary updates. The runtime doesn't change.

---

## Connection

```
runtime <-- WebSocket --> lens
```

- Lens connects to the runtime on startup
- Lens registers which being(s) it is rendering for
- Runtime pushes present JSON whenever `DerivePresent` fires for that being
- Lens sends relations back through the same socket when the user acts through an interface
- Multiple lenses can connect for the same being — phone and laptop showing the same present, each through their own component registry

---

## Multi-Device

The same being's present renders on every connected lens. Phone, laptop, TV, watch. Each lens has its own component registry tuned to its surface constraints:

- **Phone** — compact layout, touch interfaces, single section visible at a time
- **Laptop** — full layout, keyboard interface, multiple sections visible
- **TV** — read-only, large type, no input interfaces
- **Watch** — minimal, last message only, notification-style

The runtime doesn't know which device it's pushing to. It pushes present. The lens decides how to render based on its own constraints.

---

## What the Lens Is Not

- Not an app with its own state. The runtime owns all state.
- Not a client that fetches data. The runtime pushes.
- Not a single-platform build. The lens is React Native — one component registry, native rendering per surface.
- Not the interaction model. User space, being space, exchange zone — those are present sections, not app architecture. The lens renders them. It doesn't define them.
