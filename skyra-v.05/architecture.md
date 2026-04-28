# Architecture

## World Nesting

```
system world → being worlds → llm world → invariants
```

The system world contains being worlds. Each being world contains an inner and outer entity that resolve through an LLM world. The LLM world contains inference provider invariants. The recursion terminates at the provider.

Each level is a world with its own `DerivePresent`. Each level doesn't know what's above or below it. It routes in and gets a response back.

## Entity Types

**Being** — the pathos object. Identity, purpose. Data only. Not a world.

**World** — a container with a hashmap of entities and a `DerivePresent` that determines how relations resolve inside it. World types are where the specialization lives.

### World Types

**System world** — contains being worlds. Routes messages between them. Manages threads, exchanges, routing rules.

**Being world** — contains an inner entity, an outer entity, and an LLM world. Its `DerivePresent` fires the inner entity first, assembles the present with inner-thoughts, fires the outer entity, parses the response, returns outbound relations.

**LLM world** — contains inference provider invariants. Its `DerivePresent` selects a provider and routes the present to it.

### Invariants

The base case. An API endpoint, a pipe, a CLI, a shell, a screen. Terminates the recursion.
