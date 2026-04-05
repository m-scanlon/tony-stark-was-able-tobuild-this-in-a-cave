# Skyra Overview

## What Skyra Is

The user does not just talk to a model. The system receives input, decides which actor should care, projects a bounded working view, performs explicit operations, writes results back into runtime state, and only later learns from that activity. That is the basic shape of Skyra: an operating substrate for cognition, not a chatbot shell.

More precisely, Skyra is a local-first cognitive runtime for continuous experience, action, and selective learning. It separates meaning, structure, execution, and memory into distinct runtime layers. The goal is a system that stays alive, routes work through bounded operators, acts through explicit audited primitives, and retains only what is worth keeping. That gives it a better shot at durability, delegation, auditability, and growth over time than a single-session agent loop ever could.

## Conceptual Rundown

The center of the system is the kernel. It owns event intake, validates commands, routes work, keeps actor registration live, and controls who is allowed to do what. The key principle is simple: a command is not trusted because a model emitted it. It is trusted only after the surrounding runtime accepts it. The kernel is the execution authority, not the model prompt.

Above that, Skyra is organized around durable actors. An actor is a long-lived runtime operator with a bounded purpose, a callable surface, and typed stimulus boundaries. Actors persist across activity, can be invoked again later, and own their local participation in runtime history. The key architectural choice is that Skyra thinks in terms of multiple bounded actors rather than one monolithic intelligence.

Those actors do their work inside episodes. An episode is a bounded unit of activity local to a given perspective or intent: the runtime container where active state lives while work is happening. From that episode, Skyra projects a frame, the smaller inference-facing page the model actually reasons over. A single episode may span multiple frames as context shifts, but the episode boundary stays fixed. Episode is scope. Frame is view.

Memory splits along the same boundary. During an episode, the system brings retained artifacts into scope through recall: selective, anchored retrieval rather than a full dump of everything ever seen. After work is done, the system may learn, writing back into retained experience. Recall and learn are not vague model behaviors. They are explicit system boundaries: one for reading longer-lived experience, one for writing it.

The action model is intentionally small. Four primitives cover the full surface of what Skyra does:

- `act` — reach outward through tools, devices, and APIs
- `observe` — take the world in; receive typed events and stimuli
- `recall` — read from retained experience
- `learn` — write back into retained experience

Rather than a sprawling command taxonomy, the system expresses everything through this compact set and specializes through typed surfaces and schemas underneath. Legibility is a feature.

External tools, devices, and APIs are not actors. They are callable capability surfaces. The thing that decides is distinct from the thing that can be used. That boundary prevents device discovery, API usage, and actor identity from collapsing into one muddy abstraction, and it keeps the runtime legible as the surface area grows.

The role split inside the system follows the same logic. `Jarvis` is user space: the attention-facing side that interprets intent, manages meaning, and owns the user side of every episode. `Stark` is system space: the structural side that shapes topology, manages publication, and keeps runtime organization coherent. They share the same runtime model but operate at different levels of the stack, keeping user concerns and system concerns from bleeding into each other as complexity scales.

At the protocol layer, authorship and execution are always explicit. The command string names the target actor. The surrounding envelope names the caller. One actor can invoke another, but the system preserves who requested the action and who executed it. Delegation is first-class. Authority is auditable. Skyra is designed so that as it becomes more compositional, it becomes more legible, not less.

That is the intention behind all of it: a system that can grow, delegate, and learn without losing the thread of what it is doing and why.
