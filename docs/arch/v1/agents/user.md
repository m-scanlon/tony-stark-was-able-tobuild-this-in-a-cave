# User Memory (`skyra.user`)

> **DEPRECATED** — `skyra.user` is no longer a "system agent." It is Skyra's memory namespace — a provisioned memory that holds cross-domain facts about the user and Skyra's terminology config. First-injected into every LLM session. The agent model is replaced by Memory/Skill/Entity primitives. See `docs/arch/v1/kernel.md`. This document is preserved for historical reference.

## What It Is

`skyra.user` is Skyra's memory namespace that holds cross-domain facts about the user. It exists because some things are true about a person regardless of which domain they are working in — preferences, habits, communication style, recurring patterns, biographical context, life priorities.

Without it, every domain agent session starts with no knowledge of who the user is. With it, that context is always present before any domain-specific reasoning begins.

## What Problem It Solves

Agent-scoped memory works well for domain knowledge — what decisions were made in the work agent, what the gym agent knows about current training goals. What it cannot hold is knowledge that belongs to the user, not the domain.

Examples of what lives here and nowhere else:

- "Mike prefers concise plans over detailed explanations."
- "Mike is in school and working at the same time — context switching is expensive."
- "Mike's default stack is Go and Python."
- "Mike works better in the mornings."
- "Mike values privacy and self-hosted infra."

None of these belong in `work` or `gym` or `servers`. They belong to the person.

---

## Fixed Properties

| Property | Value |
|---|---|
| `agent_id` | `skyra.user` |
| Type | System agent |
| Created | Automatically at system init |
| Status | Always `active` — cannot be paused or archived |
| Injection | Every LLM session, before active domain agent |
| Injection priority | Highest — loaded first in context assembly |

---

## state.json Structure

`skyra.user` uses a modified `state.json`. The `artifact` section (which describes a project or codebase) does not apply to a person and is replaced with `identity`.

```json
{
  "metadata": {
    "name": "User Profile",
    "status": "active",
    "created_at": "",
    "last_active_at": ""
  },
  "identity": {
    "name": "",
    "timezone": "",
    "language": "en"
  },
  "knowledge": {
    "goals": [],
    "assumptions": [],
    "decisions": [],
    "facts": []
  },
  "boundary": {
    "scope": "Cross-domain user profile. Skyra may update knowledge fields based on explicit user statements. Identity fields require explicit user confirmation before commit.",
    "allowed_tool_categories": [],
    "denied_tool_patterns": [],
    "restrictions": [
      {
        "id": "identity-write-requires-confirmation",
        "description": "Identity fields (name, timezone, language) must not be updated without explicit user confirmation.",
        "matches": {
          "tool_categories": ["identity_write"],
          "tool_patterns": []
        }
      }
    ]
  }
}
```

### Section Semantics

**`identity`** — stable biographical anchors. Name, timezone, language. Slow-changing. Requires explicit confirmation to update.

**`knowledge`** — the primary working memory for user-level context:
- `goals` — long-horizon life goals and priorities (e.g. "ship Skyra", "finish degree")
- `assumptions` — things Skyra believes to be true about the user based on observed patterns
- `decisions` — explicitly stated preferences (e.g. "prefers Go over Node", "wants approval before any destructive action")
- `facts` — biographical and contextual facts (e.g. "works from home", "has a 4090 Reasoning Shard on the local network")

**`boundary`** — what Skyra is allowed to update autonomously vs what requires confirmation. Knowledge fields can be updated from strong inference. Identity fields require explicit user confirmation.

---

## Commit Authority

`skyra.user` has restricted commit authority compared to domain agents.

| Field group | Who can commit | When |
|---|---|---|
| `identity.*` | Orchestrator only | Explicit user statement + user confirmation |
| `knowledge.goals` | Orchestrator only | Explicit user statement |
| `knowledge.decisions` | Orchestrator only | Explicit user preference statement |
| `knowledge.assumptions` | Orchestrator only | Based on consistent observed patterns — not single-session inference |
| `knowledge.facts` | Orchestrator only | Explicit user statement or confirmed external fact |

Rules:
- Skyra must never commit inferences to `skyra.user` mid-task as a side effect of domain work
- Commits to `skyra.user` are separate deliberate acts, not incidental to task execution
- When Skyra learns something that should update the user profile during a domain session, it flags it for review at the end of the session — it does not commit inline

---

## Cross-Agent Write Protocol (Deferred — v2)

An open design question: if Skyra learns something about the user during a domain agent session, how does that propagate to `skyra.user`?

The read direction is already solved — `skyra.user` is always in context. The write direction is not yet designed. A domain agent session cannot write directly to `skyra.user` without breaking single-ownership semantics.

**v1 behavior**: Skyra flags insights for user review at end of session. Manual commit only. No automatic cross-agent propagation.

**v2 direction**: A structured learning event that the orchestrator processes after session close — validates the insight, checks confidence threshold, proposes a commit to `skyra.user`. Domain agents never write to `skyra.user` directly.

See `docs/arch/v1/gaps.md` G18 for the open gap.

---

## What It Is Not

- Not a settings file — preferences are expressed through commits, not config keys
- Not a scratchpad — knowledge commits require deliberate intent and reasonable confidence
- Not a domain agent — `skyra.user` has no tools, no tasks, and no jobs. It is read context, not an execution domain.
- Not a session log — recent turns live in the voice event context window, not here

---

## Template

Default state at system init: `skyra/configs/agents/user.json`

---

## Related Docs

- `docs/arch/v1/agents/README.md` — agent model overview: system vs domain agents
- `skyra/internal/agent/README.md` — Agent Service: commits, boundary enforcement
- `docs/arch/v1/gaps.md` G18 — cross-agent write protocol
