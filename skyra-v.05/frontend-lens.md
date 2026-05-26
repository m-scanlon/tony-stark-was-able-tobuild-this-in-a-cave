# Frontend Lens

## The Pattern

The same principle governs all three layers. Proven blocks, composed by the being, generation as fallback.

- **Cognition**: Skyra runtime. Observe and express on two hashmaps.
- **Action**: Composio → Nango → own Go blocks. Proven blocks for doing things in the world.
- **Display**: shadcn/ui. Proven blocks for showing things to the user.

The being is the composer across all three. It doesn't generate code to act. It doesn't generate HTML to display. It composes proven blocks in both directions — outward through expressors to change the world, and outward through lenses to show the user what it's thinking.

## Why shadcn/ui

shadcn/ui is not a dependency. It's a catalog of proven React components you copy into your codebase and own. No npm package in the critical path. No version to break. No vendor to disappear. Once you copy the block, it's yours. Same ownership model as Phase 2 of the act service — start borrowed, end owned. Except the display layer is owned from day one because shadcn's model is copy-paste-own.

The components are composable and self-contained. A card, a table, a form, a sidebar — each one works independently. They combine without layout conflicts or style collisions. The being doesn't need to understand CSS. It needs to know which blocks exist and how to arrange them. That's a composition problem, not a generation problem.

The ecosystem is massive and free. 60+ official blocks, 323+ community blocks, hundreds of extensions for specialized use cases — kanban boards, charts, data tables, dashboards, landing pages. All open source. All copy-paste. All React + Tailwind.

## How It Connects to the Runtime

The lens is a blank rendering surface that receives pushed present data. The runtime pushes a JSON component tree. The lens resolves and renders.

shadcn/ui is the component registry. Each block is a proven, tested, accessible React component with a known interface. The runtime doesn't push raw HTML. It pushes a structured description — "show a card with a table inside it, three columns, this data." The lens resolves that description against its registry of shadcn components and renders natively.

```
Runtime (Go)              Pipe              Lens (React)
                           │
Being observes.            │
Being composes tree. ───►  │  ───►   Resolve block names against
                           │         shadcn registry.
                           │         Render components.
                     ◄───  │  ◄───   User interaction events.
                           │
```

Same boundary as the act service. JSON out, result back. The being doesn't know React exists. It knows it has blocks and it arranges them.

## The Composition Model

The being observes — traverses its relationships, accumulates context, collapses its understanding of what the user needs to see.

The being expresses — selects proven shadcn blocks from its registry, arranges them into a component tree, pushes the tree to the lens.

The lens renders — resolves each block name against its local registry, instantiates the React components, displays.

No step in this process generates code. Every step composes proven blocks. The being is the composer. The blocks are the material. The lens is the surface.

### What the JSON Tree Looks Like

```json
{
  "type": "card",
  "props": { "title": "Deployment Status" },
  "children": [
    {
      "type": "table",
      "props": {
        "columns": ["service", "status", "last deploy"],
        "rows": [
          ["api", "healthy", "14:32"],
          ["worker", "deploying", "14:45"],
          ["frontend", "healthy", "13:10"]
        ]
      }
    },
    {
      "type": "button",
      "props": { "label": "Rollback", "variant": "destructive" }
    }
  ]
}
```

The being composed this from proven blocks. The lens renders it. No HTML generated. No CSS written. No code produced.

## Different Lenses, Different Registries

Different lenses have different registries. A phone lens might have a simplified set of shadcn blocks optimized for small screens. A desktop lens has the full set. A TV lens has a subset optimized for display.

Same present data, different available blocks, different rendering. The being doesn't know which lens is active. It describes what to show. The lens decides how.

If the lens receives a block name it doesn't have in its registry, it falls back to a default renderer or skips it. The being's composition is a description of intent, not a command to execute. The lens interprets what it can.

## The Fallback Hierarchy

If the being needs a block that doesn't exist in the registry:

1. **Check the shadcn ecosystem** for a community extension. 323+ community blocks and growing.
2. **AI generates a component** following shadcn conventions — React, Tailwind, self-contained, composable.
3. **If the generated block proves itself through use**, it gets promoted to the registry as a proven block.

Generation as fallback. Composition as primary. Proven blocks as the substrate. Same principle as the act service. Same principle as the runtime. Same principle everywhere.

## Future Protocol

Google's A2UI (Dec 2025, Apache 2.0) is an open protocol where agents describe rich native UIs in declarative JSON. Designed for LLMs to generate incrementally. Handles security across trust boundaries — no executable code crosses the wire.

For now, pushing JSON to your own lens with your own shadcn registry is the right move. If Skyra's lens ever needs to render in contexts you don't control — third-party apps, other devices, other platforms — A2UI is the protocol layer that sits between the being's component tree and the renderer. Don't add protocol overhead until you need it.

## Cost

$0. shadcn/ui is open source and free forever. The community ecosystem is open source. The blocks are yours once you copy them. No subscription. No per-render fee. No dependency.

## The Complete Stack

| Layer | Phase 1 (borrowed) | Phase 2 (owned) |
|-------|-------------------|-----------------|
| **Cognition** | Skyra runtime | Skyra runtime |
| **Action** | Composio free tier | Own Go act service |
| **Display** | shadcn/ui blocks | shadcn/ui blocks (already owned) |

The runtime is already yours. The action layer starts borrowed, ends owned. The display layer is owned from day one.

Total infrastructure cost at alpha: $0.

Total dependencies at maturity: zero external services in the critical path.

## The Principle

Local-first. User-owned. Proven blocks. Composable. Free.

The being thinks, acts, and shows. Every layer uses proven blocks. Every layer composes rather than generates. The durable thing is the true thing. The whole way through.
