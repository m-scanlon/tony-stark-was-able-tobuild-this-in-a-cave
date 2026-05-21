# Skyra: One interface, three methods, a hashmap, and a graph that grows a brain

I'm building a (free) cognitive runtime (just a hobby, won't be big and professional like langchain) for beings that think.

I spent three months on it — most of that figuring out the right primitive.

It's ~6,000 lines of Go. It's pre-alpha — the core idea is live, the runtime works, but it's early.

## The problem

Every AI agent framework I used felt heavy. Fixed tool registries. Orchestration layers. Thousands of lines of glue to make an LLM do something, remember it, and do something else. There has to be a simpler primitive underneath all of this.

## The insight

What if a being, a world, a device, and an operator are all the same thing?

```go
type Reality interface {
    ID() string
    Create(r *Relation) Reality
    Realize(r *Relation) string
}
```

A being implements it. A world implements it. A terminal implements it. A memory graph implements it. An LLM provider implements it. They're all realities. They all realize relations. The runtime doesn't know the difference.

A world boots from a genome file:

```
device ~name macbook ~type macos
component ~name terminal ~type stdin ~device macbook
component ~name openrouter ~type llm ~model anthropic/claude-sonnet-4-5 ~device macbook

being ~name skyra ~type llm ~identity I hold the world together. ~purpose I think, respond, and relate on behalf of the system.
being ~name michael ~type user ~identity I build Skyra. ~purpose I decide what matters.
```

`go run .` and you have a world with beings that think, talk to each other, remember, and grow.

## How it works

A single mutable Relation descends recursively through nested Reality layers, accumulating context at each level:

```
Universe → NewThread → Exchange → Self → Think/Act → Provider
```

Each layer attaches parsers to the relation as it descends. The LLM provider at the bottom evaluates all parsers into a present — a system prompt and context window that are an emergent property of the descent path, not something anyone assembled manually.

Every being has two layers. An inner layer (Think) where it reflects privately — no one sees this. An outer layer (Act) where it speaks and routes messages to peers. The being always returns to reflection between actions. The system enforces deliberation.

## Memory that grows a brain

Memory is an entity graph. Entities are like neurons — weighted by frequency and recency. Edges are like synapses — strengthened when entities co-occur.

When a cluster becomes dense enough, it automatically promotes into a specialist — an internal being with a scoped view of the graph, its own Think layer, and its own identity. Not configured. Emergent.

Specialists' heavy clusters promote into sub-specialists. Abstract at the top, concrete at the bottom. The being grows a brain from what it experiences.

## Try it

```
git clone https://github.com/skyraOS/tony-stark-was-able-tobuild-this-in-a-cave
cd skyra-v.05
export ANTHROPIC_API_KEY=your-key  # or DEEPSEEK_API_KEY or OPENROUTER_API_KEY
go run .              # terminal mode
go run ./cmd/tui      # TUI with sidebar (beings, threads, memory)
```

Requires Go 1.25+. Bring your own key — Anthropic, DeepSeek, or OpenRouter.
