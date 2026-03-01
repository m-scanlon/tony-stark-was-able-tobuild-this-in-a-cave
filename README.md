# Skyra — Personal AI Assistant

## The Idea

Most people rent their tools. They use assistants built by companies, running on company servers, storing their data in company databases. When the company changes direction, raises prices, or shuts down — your history goes with it.

Skyra is the opposite of that.

It's a personal AI system built from scratch, running on the user's own hardware, designed to last a lifetime. Every conversation, every decision, every domain of the user's life lives in a database they own. Nobody can deprecate it. Nobody can monetize the data. Nobody can take it away.

The vision is simple: **one system that everything runs through**. Every product built, every task that needs doing, every light that gets turned on, every decision that needs making — it all goes through Skyra. Over time it learns how the user thinks, how they work, and what they care about. The longer it runs, the more useful it gets. And that value compounds with the user, not with someone else.

---

## Why I Built It

Two reasons.

The first is practical — I wanted to actually build something hard. Not a tutorial project. Something with real moving parts, real design decisions, and real trade-offs I have to live with. Distributed systems, AI orchestration, voice interfaces, local inference. The kind of thing that looks good because it actually is good.

The second reason is bigger. I want an assistant that grows with me over my lifetime. The hardware will get faster. The AI models will get smarter. But the memory — the context it has about my life, my decisions, every domain I care about — that only gets richer over time. Most people will never have that because they're always starting over on someone else's platform.

---

## What It Actually Does

Skyra is built on three concepts:

**A control plane — the brain**
One machine that owns everything. It receives what the user said, figures out what they're asking for, pulls in relevant context from memory, forms a plan, and coordinates execution. It knows all active domains — work, home, servers, health, music — and decides what happens next.

**Agents — the domains of the user's life**
Each area of life is an Agent: work, home, health, servers, music. Each one has its own memory, its own set of tools, and its own rules for what Skyra is allowed to do inside it. When the user asks something, Skyra figures out which domain they're in and works from there.

**Shards — Skyra's presence on every device**
A Shard is a small piece of software that runs on any device. When it starts up, it figures out what that device can do — does it have a microphone? A GPU? Can it run scripts? — and registers those capabilities with the control plane. The control plane then knows what it has available and routes work accordingly.

Every device added extends Skyra's capabilities without changing the underlying system.

---

## How It Feels to Use

You talk to it like someone who knows you and your work deeply.

> "What did I decide about the server backups last month?" — it knows.

> "Draft a plan for the next phase of this." — it pulls up everything relevant, forms a plan, and waits for approval before doing anything.

> "Turn off the lights and set a reminder for tomorrow." — it runs the tools and confirms it's done.

Skyra proposes, the user approves. For low-stakes tasks it just runs. For anything significant it surfaces the plan first. How much autonomy it has is tunable per domain.

---

## Where It's At

Skyra is actively being built. The architecture and core systems are designed and partially implemented — voice pipeline, event delivery, memory and agent model, control plane, tool execution. The current focus is the executor loop: the phase where a planned job actually runs, tools get called, state gets updated, and replanning happens when something goes sideways.

The first milestone is simple: the user says something, it thinks, it responds. From there it's iteration.

---

## In One Sentence

Skyra is a personal operating environment that executes the user's intent across machines, owns their history, and gets smarter the longer it runs.

---

## Technical Reference

For the architecture and implementation details:

- Full system architecture: `docs/arch/v1/scyra.md`
- Executor design: `docs/arch/v1/executor.md`
- Domain expert / planning phase: `docs/arch/v1/domain-expert/README.md`
- Event ingress and ACK: `docs/arch/v1/event-ingress-ack.md`
- Task formation: `docs/arch/v1/task-formation.md`
- Agent service: `skyra/internal/project/README.md`
- Scheduler: `skyra/internal/scheduler/README.md`
- Open gaps: `docs/arch/v1/gaps.md`

---

## The Bigger Picture

I also think this is where the world is going, whether the big companies want it to or not. Right now, software is built around frontends that companies control — they decide what you see, how you interact, and what you're allowed to do with your own data. I think that model breaks down. What replaces it is an OS layer that sits on top of your data, where the interface isn't something a company designed for you but something that renders based on what your data actually looks like. Skyra is built around that idea. The thing Skyra operates on is the data — your decisions, your history, your domains. The frontend is almost secondary. In a world where AI is cheap enough to run locally and powerful enough to be genuinely useful, there's no reason that intelligence has to live on someone else's server. Decentralized AI for everyday use isn't a distant idea — it's what happens when the hardware catches up with the ambition. I'm building for that world.
