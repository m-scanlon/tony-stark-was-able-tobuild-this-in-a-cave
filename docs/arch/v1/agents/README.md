# Agent Model

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

### Open Design: Multi-Domain Invocation

The current session model assumes one active domain agent per session. The orchestrator routes to one domain. But some requests naturally span multiple domains.

Example: Mike and a friend discuss Skyra at the gym. Two things are worth remembering:
- Gym domain: training session for the day
- Developer domain: architecture ideas that came up in conversation

The current model has no answer for this. The orchestrator picks one domain. The other is dropped or lost.

The right answer depends on how the session model is defined (see G20) and how common multi-domain requests actually are in practice.

See G19 for the open gap.
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
2. **Active domain agent** — retrieved based on classifier routing.
3. **Retrieved tools** — hydrated local tools from the active domain agent.
4. **Recent turns** — session continuity context.

This ordering ensures I always know who I am talking to before knowing what domain we are in.

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
  skyra.user/          ← system agent, always present
    HEAD.json
    state.json
    commits/
  {agent_id}/          ← domain agents, one per user-created domain
    HEAD.json
    state.json
    commits/
    jobs/
```

---

## Related Docs

- `docs/arch/v1/agents/user.md` — user agent specification
- `skyra/internal/project/README.md` — Agent Service: commits, tool registry, boundary enforcement
- `docs/arch/v1/scyra.md` — full system architecture and session injection pipeline
- `docs/arch/v1/gaps.md` — open design gaps including cross-agent write protocol
