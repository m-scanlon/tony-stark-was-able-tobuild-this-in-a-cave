# Agent Model

> **DEPRECATED** — The agent model has been replaced. Agents are no longer a first-class concept. The new primitives are: **Skill** (learned class), **Job** (skill instance), **Task** (execution unit), **Memory** (provisioned namespace), **Entity** (named thing inside memory). The kernel provisions memory. Redis is the skill registry and trust boundary. See `docs/arch/v1/kernel.md` for the canonical model. This document is preserved for historical reference only.

## Two Types of Agents

Skyra has two distinct categories of agent. Both use the same infrastructure — same object store, same commit model, same state.json structure, same global tools. The distinction is in initialization, injection behavior, and lifecycle rules.

### System Agents

Predefined at system init. Not created by the user. Reserved `agent_id` under the `skyra.` namespace.

Rules:
- Created automatically on first boot if they do not already exist in `.skyra/agents/`
- Status is always `active` — Agent Service rejects `paused` or `archived` for system agent IDs
- Always injected into every LLM session regardless of what is being asked
- Injection order: system agents load before the active domain agent
- Commit authority is restricted — see individual agent specs for what can be committed and by whom

System agents defined in v1:

| agent_id | Purpose |
|---|---|
| `skyra.user` | Cross-domain profile of the user — preferences, habits, life context, biographical facts |

### Domain Agents

Created by the user (or by Skyra on the user's behalf) via `create_agent`. Each one represents a scoped domain of the user's life: work, gym, servers, music, home, etc.

Rules:
- Standard lifecycle: `active | paused | archived`
- Retrieved based on relevance — not always injected
- Injected after system agents in the session context block
- Created by copying the domain agent template from `skyra/configs/agents/domain.json`

---

## Template Concept

Every agent starts from a template. Templates are default `state.json` files stored in `skyra/configs/agents/`. When an agent is created, the Agent Service copies the appropriate template into `.skyra/agents/{agent_id}/state.json` and initializes a first commit.

```
skyra/configs/agents/
  user.json        ← template for skyra.user (applied at system init)
  domain.json      ← template for user-created domain agents
```

Templates define structure and defaults — they do not carry user data. The runtime instance in `.skyra/agents/{agent_id}/` is the source of truth.

---

## Session Injection Order

When an LLM session starts, context is assembled in this order:

1. **System agents** — always present, loaded first. In v1: `skyra.user`.
2. **Active domain agent** — the domain agent that self-selected relevance from the context blob. Injected after system agents.
3. **Recent turns** — session continuity context.

Tools are **not** part of the context package. Tool retrieval is owned by the Agent Service inside the LLM session during planning — not assembled upfront. See `skyra/internal/agent/README.md`.

This ordering ensures the system always knows who it's talking to before knowing what domain is active.

---

## agent_id Namespace

| Prefix | Type | Examples |
|---|---|---|
| `skyra.` | System agents — reserved | `skyra.user` |
| _(none)_ | Domain agents — user-created | `work`, `gym`, `servers` |

The `skyra.` prefix is reserved. `create_agent` must reject any `agent_id` that starts with `skyra.`.

---

## Object Store Layout

```
.skyra/agents/
  skyra.user/          ← system agent, always present (independent git repo)
    .git/
    state.json
    tools/
  {agent_id}/          ← domain agents, one per user-created domain (independent git repo)
    .git/
    state.json
    tools/
    jobs/
```

---

## Related Docs

- `docs/arch/v1/agents/user.md` — user agent specification
- `skyra/internal/agent/README.md` — Agent Service: commits, tool registry, boundary enforcement
- `docs/arch/v1/scyra.md` — full system architecture and session injection pipeline
- `docs/arch/v1/gaps.md` — open design gaps including cross-agent write protocol
