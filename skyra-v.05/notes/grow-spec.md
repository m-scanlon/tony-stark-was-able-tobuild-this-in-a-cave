# Grow — Realities Creating Realities

## Current State

Grow is a hardcoded factory on NewThread (`newthread.go:206`). It's a switch statement that knows the full internal composition of every being type. It's the one place in the system that breaks composability.

```go
switch beingType {
case "llm":
    // knows Self needs Think with recall/remember/skill
    // knows Self needs Act with plan
    // knows Think needs an LLM component
case "user":
    // knows User needs a device
}
```

This works but it means Thread knows the internals of Self, Think, Act, User, and every operator. Nothing else in the system works that way.

## The Problem

A reality needs context to create another reality:
- What devices exist?
- What components are available?
- Where does the new reality register?

That context lives above the creator in the tree. So either:
1. The creator reaches up (breaks the recursion)
2. The context comes down on the relation (stays composable)

## Target: Create Is the Factory

Every Reality already has `Create(*Relation) Reality`. If the Relation carries enough context, any Reality can create sub-realities within its own world.

- Thread grows beings — it puts available devices and components on the relation, calls Create on the being type
- Self grows operators — it puts its Think/Act layers on the relation, the operator registers itself
- MacOS grows components — it puts its component registry on the relation, the component registers itself
- A being could grow another being if it has the right context on its relation

Each level grows within its own world. No central factory. The relation is the context carrier.

## What This Means

### Create becomes richer

Currently Create is mostly a constructor — it takes a Relation and returns a new Reality with some fields set. Under this model, Create also receives the world context it needs to wire itself in.

The Relation already has:
- `Impulse` — the grow command / configuration
- `Realities` — map of available realities (could carry devices, components, peers)
- `Exports` — arbitrary key-value context

A Self.Create could:
1. Read its identity/purpose from the Impulse
2. Find available LLM components from Realities
3. Create its own Think and Act by calling their Create with the same context
4. Return a fully wired Self

Thread's Grow shrinks to: build the relation with the right context, call Create on the right type, register the result.

### The genome maps to Create calls

Bootstrap already does this — it parses the genome and calls Create. The difference is that runtime Grow would do the same thing. The genome is the language. Create is the interpreter. Bootstrap and runtime Grow are both callers.

### Beings creating beings

If a being has enough context on its relation (devices, component refs, peer list), it could call Create on a new being type. This is the reproduction path from world-physics — "creating a new being costs accumulated experience." The economics layer gates it, but the mechanism is just Create with context.

## Open Questions

1. **What context does Create need?** The minimum set for each type. Self needs an LLM component. User needs a device. Operators need nothing. What's the shape of the relation that carries this?

2. **Who builds the context relation?** Thread currently knows to put devices on the relation. If Create is the factory, does the caller always need to know what context to provide? Or does the type declare what it needs?

3. **Registration.** Create returns a new Reality, but someone has to put it in the right hashmap (Thread.Beings, MacOS.Components, Think.Operators). Does Create register itself, or does the caller register the result?

4. **The Grow operator.** Currently `grow` is a verb that Thread intercepts in the main loop. If Create is the factory, is `grow` still a system-level operator, or can any reality handle it? A being saying `<grow>~name poet ~type llm</grow>` — does that go to Thread, or to the being's own world?

5. **Alpha scope.** How much of this do we need before June? The current Grow works. The refactor is architectural purity. Worth doing if it unblocks something (like beings creating beings), not worth doing just to be clean.
