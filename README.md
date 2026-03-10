# Skyra — Cognitive OS

## The Idea

Most people rent their tools. They use assistants built by companies, running on company servers, storing their data in company databases. When the company changes direction, raises prices, or shuts down — your history goes with it.

Skyra is the opposite of that.

It's a cognitive OS built from scratch, running on the user's own hardware, designed to last a lifetime. Every conversation, every decision, every domain of the user's life lives in a graph they own. Nobody can deprecate it. Nobody can monetize the data. Nobody can take it away.

The vision is simple: **one system that everything runs through**. Every product built, every task that needs doing, every light that gets turned on, every decision that needs making — it all goes through Skyra. Over time it learns how the user thinks, how they work, and what they care about. The longer it runs, the more useful it gets. And that value compounds with the user, not with someone else.

---

## The Mental Model

> Nodes are identity. Edges are history. Truth is derived, not stored.

Skyra's memory is a property graph. Nodes are the things in your life — people, places, tools, concepts, skills. Edges are the relationships between them, timestamped and weighted, accumulating over time.

There is no "current state" field. When Skyra needs to know what's true right now, she reasons over the graph — edge types, weights, recency, history — and derives a conclusion. The graph never loses data. It only grows.

---

## The Bigger Picture

Right now, software is built around frontends that companies control. They decide what you see, how you interact, and what you're allowed to do with your own data. Every conversation you have with their assistant makes their model better. You get the answer. They get the data point.

That's a bad deal. And I think it breaks down.

What replaces it is an OS layer that sits on top of your data — where the interface isn't something a company designed for you, but something that renders based on what your data actually looks like. The thing the system operates on is your decisions, your history, your domains. The frontend is almost secondary.

Most people think about AI assistants as a better search box. This is something closer to a cognitive runtime. The app is incidental. What matters is the layer underneath that knows how you think.

In a world where AI is cheap enough to run locally and powerful enough to be genuinely useful, there's no reason that intelligence has to live on someone else's server. Local-first AI for everyday use isn't a distant idea — it's what happens when the hardware catches up with the ambition.

I'm building for that world.

---

## Why I'm Building It

Two reasons.

The first is bigger. I want a system that grows with me over my lifetime. The hardware will get faster. The AI models will get smarter. But the memory — the context it has about my life, my decisions, every domain I care about — that only gets richer over time. Most people will never have that because they're always starting over on someone else's platform.

The second is practical — I want to actually build something hard. Not a tutorial project. Something with real moving parts, real design decisions, and real trade-offs I have to live with. Distributed systems, AI orchestration, voice interfaces, local inference. The kind of thing that looks good because it actually is good.

---

## How It Works

Skyra is built on two concepts:

**Shards — Skyra's presence on every device**
A Shard is a small piece of software that runs on any device. When it starts up, it fingerprints what that device can do — does it have a microphone? A GPU? Can it run scripts? — and registers those capabilities with the rest of the network. The network routes work based on what each Shard advertises. One Shard runs the control plane. Another handles voice. A third handles deep reasoning. Every device added extends Skyra's capabilities without changing the underlying system.

**Skills — learned capabilities**
Skills are not defined. They are learned. Skyra watches how the user works, identifies patterns, and crystallizes them into skills — executable capabilities with their own memory namespace. The longer the system runs, the more capable it becomes.

---

## How It Feels to Use

You talk to it like someone who knows you and your work deeply.

> "What did I decide about the server backups last month?" — it knows.

> "Draft a plan for the next phase of this." — it pulls up everything relevant, forms a plan, and waits for approval before doing anything.

> "Turn off the lights and set a reminder for tomorrow." — it runs the tools and confirms it's done.

Skyra proposes, the user approves. For low-stakes tasks it just runs. For anything significant it surfaces the plan first.

---

## Where It's At

The kernel is built — the central execution boundary that every command passes through. The memory model is designed — a property graph with a two-tier trust system (observational and committed), append-only committed layer, and a background reasoning process that turns session history into graph nodes and edges.

The current focus is the shard model — how Skyra's execution layer runs across distributed hardware.

---

## In One Sentence

Skyra is a cognitive OS that reasons over the history of your life to derive what's true right now, gets smarter the longer it runs, and belongs entirely to you.

---

## Technical Reference

- Architecture: `docs/arch/v1/kernel.md`
- Memory model: `docs/arch/v1/memory-structure.md`
- Skill lifecycle: `docs/arch/v1/skill-lifecycle.md`
- Design principles: `docs/arch/v1/principles.md`
- Open gaps: `docs/arch/v1/gaps.md`
