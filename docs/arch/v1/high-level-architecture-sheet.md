# OctaOS High-Level Architecture Sheet (v1)

> Project naming note: many v1 docs still use the legacy name "Skyra." This sheet uses the current project name: **OctaOS**.

## 1. What OctaOS Is

OctaOS is a personal AI operating system that runs across your own devices.

It is designed to:

- remember your decisions and context over time
- execute tasks through a trusted runtime
- improve continuously from real usage
- keep user data ownership local-first

## 2. Core Product Thesis

Most AI systems predict the next token.

OctaOS is built to **predict you**: what context you need, what action should happen next, and what response style actually helps you.

It applies social-media style behavioral modeling in reverse:

- social platforms learn behavior to maximize engagement for the platform
- OctaOS learns behavior to maximize usefulness for the user

## 3. Architecture in One View

```text
Input Shards (voice/chat/devices)
  -> API Gateway (ingress/egress)
  -> Kernel (single execution boundary)
  -> Skill execution on capable shards
  -> Memory + Observational signals
  -> Output back to user
```

All execution routes through the kernel using one command protocol:

```text
Octa <tool> [args]
```

## 4. Main Building Blocks

- **Shards**: runtime nodes on devices (voice, compute, storage, etc.).
- **API Gateway**: receives commands, validates, routes, returns results.
- **Kernel**: trust enforcement, skill invocation, job/task orchestration.
- **Skill Registry (Redis)**: execution trust gate and provisioning state.
- **Memory Graph**: long-lived committed memory + observational working layer.
- **Cron/Background Loop**: periodic reasoning, integration, pattern updates.

## 5. Data + Memory Model

OctaOS uses a two-layer memory model:

- **Committed layer**: user-approved, trusted, append-only.
- **Observational layer**: free-form system reasoning/workspace, not trusted until promoted.

The system stores enough signal to detect patterns over time and improve retrieval, planning, and execution quality.

## 6. Request Lifecycle (High Level)

1. User request arrives from a shard.
2. Gateway validates command + permissions.
3. Kernel resolves skill and creates job/tasks.
4. Work runs on the best-capable shard.
5. Actions/results are audited and returned.
6. Relevant signals update memory and prediction quality.

## 7. Trust, Privacy, and Safety

- Skill definitions can be open or closed.
- Execution is policy-bound and auditable.
- Private context can remain local to execution.
- Core invariant for market skills:

```text
output_public, derivation_private
```

Meaning: outputs/actions may be observable at execution boundaries, while internal reasoning remains sealed.

## 8. Marketplace Direction

For sellable skills, OctaOS is converging on:

- **closed definition**
- **sealed reasoning**
- **audited actions**

with model-scoped verification receipts and strike policy for unsafe behavior.

## 9. North Star Principle

**Root design decisions in reality.**

If design is rooted in reality, that is the north star for the model.

