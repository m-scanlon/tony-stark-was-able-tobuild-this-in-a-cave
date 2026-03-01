# Skyra — Personal AI Assistant

## The Idea

Most people rent their tools. They use assistants built by companies, running on company servers, storing their data in company databases. When the company changes direction, raises prices, or shuts down — your history goes with it.

Skyra is the opposite of that.

It's a personal AI system I'm building from scratch, running on my own hardware, that I'll use for the rest of my life. Every conversation, every decision, every project lives in a database I own. Nobody can deprecate it. Nobody can monetize my data. Nobody can take it away.

The vision is simple: **one system that everything runs through**. Every product I build, every task I need done, every light I turn on, every decision I need to make — it all goes through Skyra. Over time it learns how I think, how I work, and what I care about. The longer it runs, the more useful it gets. And that value compounds with me, not with someone else.

---

## Why I Built It

Two reasons.

The first is practical — I wanted to actually build something hard. Not a tutorial project. Something with real moving parts, real design decisions, and real trade-offs I have to live with. Distributed systems, AI orchestration, voice interfaces, local inference. The kind of thing that looks good because it actually is good.

The second reason is bigger. I want an assistant that grows with me over my lifetime. The hardware will get faster. The AI models will get smarter. But the memory — the context it has about my life, my projects, my decisions — that only gets richer over time. Most people will never have that because they're always starting over on someone else's platform. I won't be.

---

## What It Actually Does

Skyra sits across three machines that work together:

**Raspberry Pi — always on, always listening**
This is the voice layer. It detects when I'm talking, converts speech to text, and can give a quick answer from recent context while the heavier thinking happens in the background. It's always on, low power, and sits on my desk.

**Mac mini — the brain**
This is where decisions get made. It receives what I said, figures out what I'm asking for, pulls in relevant context from memory, forms a plan, and coordinates execution. It knows all my active projects and can take actions on my behalf — running scripts, managing files, calling APIs, whatever the task needs.

**GPU machine — heavy thinking**
For complex reasoning, deep coding problems, or anything that needs serious horsepower, the Mac delegates to a dedicated GPU machine running a large language model locally. No cloud, no API keys, no usage limits.

---

## How It Feels to Use

I talk to it like I'd talk to someone who knows me and my work deeply.

"What did I decide about the server backups last month?" — it knows.

"Draft a plan for the next phase of this project" — it pulls up everything relevant, forms a plan, and asks me to approve before doing anything.

"Turn off the lights and set a reminder for tomorrow" — it runs the tools, confirms it's done.

The key thing is that I stay in control. Skyra proposes, I approve. For low-stakes tasks it just runs. For anything significant it surfaces the plan first. I've designed exactly how much autonomy it has, and I can tune that per project.

---

## Where It's At

Skyra is actively being built. The architecture and core systems are designed and partially implemented — voice pipeline, event delivery, memory and project model, control plane, tool execution. The current focus is the executor loop: the phase where a planned job actually runs, tools get called, state gets updated, and replanning happens when something goes sideways.

The first milestone is simple: I say something, it thinks, it responds. From there it's iteration.

---

## In One Sentence

Skyra is a personal operating environment that executes my intent across machines, owns my history, and gets smarter the longer I run it.

---

## Technical Reference

For the architecture and implementation details:

- Full system architecture: `docs/arch/v1/scyra.md`
- Executor design: `docs/arch/v1/executor.md`
- Domain expert / planning phase: `docs/arch/v1/domain-expert/README.md`
- Event ingress and ACK: `docs/arch/v1/event-ingress-ack.md`
- Task formation: `docs/arch/v1/task-formation.md`
- Project service: `skyra/internal/project/README.md`
- Scheduler: `skyra/internal/scheduler/README.md`
- Open gaps: `docs/arch/v1/gaps.md`
