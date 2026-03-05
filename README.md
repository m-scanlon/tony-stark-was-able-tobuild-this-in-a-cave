# Skyra — Personal AI OS

## The Idea

Most people rent their tools. They use assistants built by companies, running on company servers, storing their data in company databases. When the company changes direction, raises prices, or shuts down — your history goes with it.

Skyra is the opposite of that.

It's a personal AI system built from scratch, running on the user's own hardware, designed to last a lifetime. Every conversation, every decision, every domain of the user's life lives in a database they own. Nobody can deprecate it. Nobody can monetize the data. Nobody can take it away.

The vision is simple: **one system that everything runs through**. Every product built, every task that needs doing, every light that gets turned on, every decision that needs making — it all goes through Skyra. Over time it learns how the user thinks, how they work, and what they care about. The longer it runs, the more useful it gets. And that value compounds with the user, not with someone else.

---

## The Bigger Picture

Right now, software is built around frontends that companies control. They decide what you see, how you interact, and what you're allowed to do with your own data. Every conversation you have with their assistant makes their model better. You get the answer. They get the data point.

That's a bad deal. And I think it breaks down.

What replaces it is an OS layer that sits on top of your data — where the interface isn't something a company designed for you, but something that renders based on what your data actually looks like. The thing the system operates on is your decisions, your history, your domains. The frontend is almost secondary.

Most people think about AI assistants as a better search box. This is something closer to a personal runtime. The app is incidental. What matters is the layer underneath that knows how you think.

In a world where AI is cheap enough to run locally and powerful enough to be genuinely useful, there's no reason that intelligence has to live on someone else's server. Local-first AI for everyday use isn't a distant idea — it's what happens when the hardware catches up with the ambition.

I'm building for that world.

---

## Why I'm Building It

Two reasons.

The first is bigger. I want an assistant that grows with me over my lifetime. The hardware will get faster. The AI models will get smarter. But the memory — the context it has about my life, my decisions, every domain I care about — that only gets richer over time. Most people will never have that because they're always starting over on someone else's platform.

The second is practical — I want to actually build something hard. Not a tutorial project. Something with real moving parts, real design decisions, and real trade-offs I have to live with. Distributed systems, AI orchestration, voice interfaces, local inference. The kind of thing that looks good because it actually is good.

---

## What It Actually Does

Skyra is built on two concepts:

**Shards — Skyra's presence on every device**
A Shard is a small piece of software that runs on any device. When it starts up, it fingerprints what that device can do — does it have a microphone? A GPU? Can it run scripts? — and registers those capabilities with the rest of the network. The network routes work based on what each Shard advertises, not what kind of machine it is. One Shard runs the control plane because it's currently the most capable node. Another handles voice because it has a microphone and a speaker. A third handles deep reasoning because it has a GPU. These are capability designations, not permanent hardware roles. Every device added extends Skyra's capabilities without changing the underlying system.

**Agents — the domains of the user's life**
Each area of life is an Agent: work, home, health, servers, music. Each one has its own memory, its own set of tools, and its own rules for what Skyra is allowed to do inside it. When the user asks something, Skyra figures out which domain they're in and works from there.

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

Skyra executes your intent across every area of your life, owns your history, and gets smarter the longer it runs.

---

## Technical Reference

For the architecture and implementation details:

- Full system architecture: `docs/arch/v1/scyra.md`
- Executor design: `docs/arch/v1/executor.md`
- Domain expert / planning phase: `docs/arch/v1/domain-expert/README.md`
- Event ingress and ACK: `docs/arch/v1/event-ingress-ack.md`
- Task formation: `docs/arch/v1/task-formation.md`
- Agent service: `skyra/internal/agent/README.md`
- Scheduler: `skyra/internal/scheduler/README.md`
- Open gaps: `docs/arch/v1/gaps.md`
