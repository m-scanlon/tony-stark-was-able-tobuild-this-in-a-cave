# Agent Service

## What It Is

The Agent Service is the single owner of all agent state in Skyra. Nothing reads or writes agent state without going through this service. It is the foundational layer that all other services interface with.

Note: earlier architecture docs refer to this as the **Project Service** or **Memory Service**. The Agent Service supersedes those definitions with a more clearly scoped responsibility model.

## Core Principles

- **Single source of truth** — all agent data lives in exactly one canonical location. The AI never edits state directly. All changes occur through explicit commits.
- **Append-only commit history** — agent changes are immutable commit objects. Each commit contains a full snapshot, links to its parent, and can be used for rollback.
- **Vector index is derived data** — the vector DB is not the source of truth. It indexes agent state for semantic retrieval and can be deleted and rebuilt at any time.

## What It Owns

**Agent Registry (SQLite)**
A lightweight index of all agents. Used by the context engine as a fast first gate before any deeper retrieval.
- Agent identity and metadata
- Status: `active | paused | archived`
- Last active timestamp

**Object Store Interface**
The only entry point for reading and mutating agent state on disk. Filesystem in Phase 1, S3/MinIO in Phase 2 with no structural changes required.

- Current agent state (knowledge, artifact mapping, boundary definition)
- Immutable commit history
- HEAD pointer management
- Rollback and audit trail

**Local Tool Registry (Vector DB)**
Each agent has a set of local tools registered and indexed semantically. The tool description is embedded and indexed — vector search runs against these embeddings using the current request as the query. Tools are retrieved by semantic similarity, not injected wholesale.

Each tool record in the registry:

```json
{
  "id": "tool-uuid",
  "score": 0.92,
  "name": "write_file",
  "description": "Writes content to a file at the specified path.",
  "input_schema": {
    "type": "object",
    "properties": {
      "path": { "type": "string" },
      "content": { "type": "string" }
    },
    "required": ["path", "content"]
  },
  "agent_id": "jarvis",
  "categories": ["filesystem_write"],
  "requires_approval": true
}
```

`score` is the semantic similarity score from the vector search. The Agent Service applies a score threshold before hydration — results below the threshold are dropped. Lock status is not stored on the tool record. It is computed at hydration time by joining the tool record against the agent boundary in `state.json`.

## Object Store Structure

```
.skyra/agents/{agent_id}/
  HEAD.json            ← pointer to current commit
  state.json           ← materialized current state (four sections below)
  commits/             ← immutable commit history
    {commit_id}.json
  working/             ← scratch space (mutable, not versioned — see Working State below)
  jobs/
    {job_id}/
      envelope.json    ← job envelope
      tasks/
        {task_id}/
          tasksheet.json  or  workplan.json
          notes.md
```

## Working State vs Committed State

The object store has two distinct partitions:

**Working state** (`working/`) — the executor's scratch pad. The system writes freely here during job execution to test ideas, validate approaches, and reason through problems on paper. No user approval required. Working state is mutable and throwaway — it does not appear in the version history. Cleaned up after job completion.

**Committed state** (`state.json` + `commits/`) — requires user approval. When the executor produces output worth persisting canonically, it proposes a commit via `propose_commit`. The user accepts or rejects. Only accepted commits update `state.json` and enter the commit history.

This distinction gives the system room to think without making permanent decisions. The audit trail stays clean — only intentional, user-approved changes appear in the commit log.

## Domain Agent as Doorkeeper

Each domain agent is the doorkeeper of its own domain. When a turn arrives with the context blob, the domain agent self-selects — it decides whether the turn is relevant to it, checks whether the turn impacts an ongoing job, and forms an estimation call if a job is needed.

No external classifier makes routing decisions. The agent knows its domain better than any classifier can.

### state.json Structure

```json
{
  "metadata": {
    "name": "...",
    "status": "active | paused | archived",
    "created_at": "...",
    "last_active_at": "..."
  },
  "knowledge": {
    "goals": [],
    "assumptions": [],
    "decisions": [],
    "facts": []
  },
  "artifact": {
    "type": "...",
    "location": "..."
  },
  "boundary": {
    "scope": "Human-readable description of this agent and what Skyra is allowed to do in it.",
    "allowed_tool_categories": ["category_a", "category_b"],
    "denied_tool_patterns": ["pattern_*"],
    "restrictions": [
      {
        "id": "restriction-id",
        "description": "Human-readable rationale for auditing.",
        "matches": {
          "tool_categories": ["category"],
          "tool_patterns": ["pattern_*"]
        }
      }
    ]
  }
}
```

### Commit Object Format

Each commit is immutable and append-only.

```json
{
  "commit_id": "2026-02-09T21-10-33Z",
  "agent": "jarvis",
  "parent": "2026-02-09T20-55-12Z",
  "actor": {
    "type": "ai",
    "model": "qwen2.5-coder:7b",
    "user": "mike"
  },
  "message": "...",
  "timestamp": "2026-02-09T21-10-33Z",
  "changes": [
    {
      "op": "set",
      "path": "/knowledge/decisions/0",
      "value": "..."
    }
  ],
  "snapshot": {}
}
```

## What It Exposes

### Global Tools
Available to every LLM session regardless of agent. Always injected directly into the session. Never retrieved — always present.

- `get_agent_state` — read current state for an agent
- `get_agent_facts` — read facts and assumptions
- `get_commit_history` — read the commit log
- `list_agents` — read the registry
- `get_job_tasks` — read tasks under a job
- `create_agent` — register a new agent
- `propose_commit` — submit a patch against agent state
- `apply_commit` — apply an approved patch and update HEAD
- `rollback_commit` — revert to a previous commit
- `update_agent_status` — set active / paused / archived
- `update_last_active` — called by scheduler on job completion

### Local Tool Registry
Per-agent tools registered at agent setup. Retrieved via vector search — not injected directly. The LLM sees all retrieved tools, including locked ones, with their access status attached.

## Tool Hydration

Tool hydration is the intermediary step between raw vector DB results and what the LLM session receives. The Agent Service runs hydration on every tool in the result set before returning anything to the Domain Expert.

Hydration joins the raw tool record against the agent boundary in `state.json` and computes the `access` field:

```json
{
  "id": "tool-uuid",
  "score": 0.92,
  "name": "write_file",
  "description": "Writes content to a file at the specified path.",
  "input_schema": { ... },
  "agent_id": "jarvis",
  "categories": ["filesystem_write"],
  "requires_approval": true,
  "access": {
    "status": "locked",
    "reason": "Filesystem writes are restricted for this agent."
  }
}
```

`access.status` is either `allowed` or `locked`. Lock status is always derived fresh at hydration time — it is never stored on the tool record itself, because the same tool may be locked in one agent and open in another.

The LLM receives the full hydrated list. It can see and reason over all tools — including locked ones — but the BoundaryValidator enforces what actually executes.

## Boundary Enforcement

The `boundary` section in `state.json` is enforced in code at two layers. Plain-text restrictions are insufficient — the LLM cannot be relied upon to honor prose rules under all conditions.

### Layer 1: Hydration (retrieval time)

After vector search returns results, the Agent Service hydrates each tool with its access status for the current agent. No tools are hidden. The LLM sees everything, with locked tools clearly marked. This gives the LLM full situational awareness — it knows what exists and what it can use.

### Layer 2: BoundaryValidator (runtime)

Before any tool call is dispatched to execution, the BoundaryValidator checks the proposed call against the agent boundary. This is pure code — no LLM in the loop.

If the tool is locked, execution pauses and a permission prompt is sent to the user. The prompt clearly states what Skyra wants to do and why:

**Prompt payload (Skyra → user):**
```json
{
  "tool": "write_file",
  "why": "To save the updated nginx config after modifying the upstream block.",
  "how": "write_file({ path: \"/etc/nginx/nginx.conf\", content: \"...\" })"
}
```

**Response (user → Skyra):**
```json
{
  "decision": "allow_always | allow_once | deny"
}
```

| Decision | Behavior | Persisted? |
|---|---|---|
| `allow_always` | Tool is permanently unlocked for this agent. An immediate commit is written to `state.json` before execution resumes. | Yes — boundary updated via commit |
| `allow_once` | This single invocation is allowed for the stated reason. Boundary is unchanged. | No — ephemeral |
| `deny` | Tool call is blocked. Skyra replans around the denied tool and attempts to complete the task another way. | No — ephemeral |

### What `scope` and `description` Are For

`boundary.scope` and each restriction's `description` are human-readable fields. They may be included in system prompts as a soft context hint, but enforcement never relies on them. They exist for human auditing, agent setup review, and operator understanding — not for controlling LLM behavior.

## What It Does Not Do

- Does not propose changes — that is the Domain Expert's job
- Does not execute tasks — that is the Orchestrator's job
- Does not decide which tools are relevant — that is the context engine's job
- Does not enforce approval during execution — approval is surfaced at plan review only
- Does not manage system-level tool execution (shell, HTTP) — that is the tooling service's job

## Tool Approval vs Plan Approval Gate

These are two distinct concepts:

**Plan Approval Gate** — after the Domain Expert forms a full plan, the user reviews and approves the entire plan before execution begins. Defined in `docs/arch/v1/domain-expert/README.md`.

**Tool `requires_approval` flag** — determines which tools are surfaced and highlighted to the user during plan review. Does not pause execution mid-run. Shapes what the user sees when reviewing the plan.

## Who Calls It

- Context engine — reads registry and agent state
- Scheduler — updates `last_active_at` on job completion
- Orchestrator — applies task outcomes via commits
- LLM session — calls global tools directly
- Domain Expert — reads agent state during task formation

## Related Docs

- `docs/arch/v1/scyra.md` — system architecture and canonical pipeline
- `docs/arch/v1/agents-services.md` — full service and shard catalog
- `docs/arch/v1/task-formation.md` — how tasks are formed against agent state
- `docs/arch/v1/domain-expert/README.md` — domain expert and plan approval gate
- `skyra/internal/scheduler/README.md` — scheduler service
