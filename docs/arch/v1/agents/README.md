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

### === SUGGESTION BY KUNJ ===
### System Agent Suggestion: `skyra.orchestrator`

Predefined at system.init as well if does not exist in `./skyra/agents`. Its task is to take in queries and route it to correct domain. Sits in the Mac Mini/Shard which is currently active. 

It has knowledge base from:
1) `skyra.user`
2) All Domain Agents

`skyra.orchestrator` acts as the module agent which connects user info to domain info as well as routes user query into correct domain agent/domain tool needed.

Rules:
- Created automatically on first boot if does not exist in `./skyra/agents`
- Status is always active — Agent Service rejects paused or archived for system agent IDs
- Always injected into every LLM session regardless of what is being asked
- Injection Order: Injected right after the `skyra.user` is loaded.
- Allow users and domain agents to make commits to `skyra.orchestrator`. We introduce a guardrail agent that validates the user input before committing it to `skyra.orchestrator`

Benefits:
- Provides an architectural visibility layer to identify which domain agent was invoked for a given task
- Allows invoking of multiple domain agents (for e.g. Mike and Kunj meet at gym and discuss Skyra. The entire event would require invoking of two domain agents in Mike's Skyra i.e. the Gym Domain to track fitness goals for the day, the Developer Domain to store ideas discussed during the gym session).

Note: I read later on in `domain-expert/README.md` where you mention the controlplane orchetsrator already exists. It solves the issue of architectural visibility but we can still take look into the `skyra.orchestrator`'s benefits in cross-agent or multi-agent invoking.

### === END OF SUGGESTIONS ===
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
