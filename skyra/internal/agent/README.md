# Agent Service

## What It Is

The Agent Service is the single owner of all agent state in Skyra. Nothing reads or writes agent state without going through this service. It is the foundational layer that all other services interface with.

Note: earlier architecture docs refer to this as the **Project Service** or **Memory Service**. The Agent Service supersedes those definitions with a more clearly scoped responsibility model.

## Core Principles

- **Single source of truth** — all agent data lives in exactly one canonical location. The AI never edits state directly. All changes go through the commit gate.
- **Git is the version control** — each agent directory is an independent git repo. Commit history, rollback, diff, and audit trail are all git primitives. No custom commit infrastructure.
- **Domains are independently versionable** — separate git root per agent. Rolling back `jarvis` does not touch `home`. Each agent owns its own history.

## What It Owns

**Agent Registry (SQLite)**
A lightweight index of all agents. Used by the context engine as a fast first gate before any deeper retrieval.
- Agent identity and metadata
- Status: `active | paused | archived`
- Last active timestamp

**Object Store Interface**
The only entry point for reading and mutating agent state on disk. Each agent directory is an independent git repo managed via go-git — no git binary required at runtime.

- Current agent state (`state.json` — knowledge, artifact mapping, boundary definition)
- Commit history, HEAD, rollback, diff — all standard git operations via go-git
- `create_agent` initializes the directory and runs `git init`
- `propose_commit` stages changes and waits for user approval before committing
- The LLM can navigate and inspect any agent's git history directly using shell git commands

**Tools**

Tools live in the object store under `tools/`. No vector index. The LLM discovers and uses tools by walking the filesystem directly — `ls`, `cat`, execute. The object store is self-describing.

Tool definition (at `.skyra/agents/{agent_id}/tools/{tool_name}/tool.json`):

```json
{
  "name": "read_thermostat",
  "description": "Read the current thermostat temperature and mode.",
  "categories": ["sensor_read"],
  "requires_approval": false,
  "required_capabilities": ["thermostat.read"],
  "input_schema": { ... },
  "impl": {
    "kind": "shell",
    "script": ".skyra/agents/home/tools/read_thermostat/run.sh"
  }
}
```

`impl.kind` can be `shell`, `http`, or `builtin`. Shell tools point to a script in the same directory — the LLM can read and execute it directly. Tool updates are commits — same versioning and audit trail as agent state.

Lock status is computed at runtime by joining against the agent boundary in `state.json`. Not stored on the tool definition itself.

## Two Layers, Two Purposes

- Object Store Tools (this service): tool definitions live as files in `tools/` inside the agent git repo. The LLM discovers them by walking the filesystem. No registry, no vector index.
- Shard Capability Registry (tooling/orchestrator): global runtime capability inventory for shard selection and dispatch. The Agent Service does not own this.

## Object Store Structure

```
.skyra/agents/{agent_id}/    ← independent git repo (go-git managed)
  .git/                      ← git history, HEAD, rollback — all standard git
  .gitignore                 ← working/ is ignored
  state.json                 ← current committed state (knowledge, boundary, artifact)
  tools/
    {tool_name}/
      tool.json              ← tool definition (input schema, impl kind, capabilities)
      run.sh                 ← shell impl, executable directly by the LLM
  working/                   ← scratch space, gitignored, cleaned up after job
  jobs/
    {job_id}/
      tasks/
        {task_id}/
          tasksheet.json  or  workplan.json
```

Each agent is an isolated git root. The LLM navigates the filesystem with shell tools and git commands — `ls`, `cat`, `git log`, `git diff`, `git show`. No special API needed to inspect state, tools, or history.

## Working State vs Committed State

The object store has two distinct partitions:

**Working state** (`working/`) — the executor's scratch pad. The system writes freely here during job execution to test ideas, validate approaches, and reason through problems on paper. No user approval required. Working state is mutable and throwaway — it does not appear in the version history. Cleaned up after job completion.

**Committed state** (`state.json`) — requires user approval. When the executor produces output worth persisting canonically, it proposes a commit via `propose_commit`. The user accepts or rejects. Only accepted commits update `state.json` and enter the git history via go-git.

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

### Commits

Git commits via go-git. The commit message carries the actor metadata — model, user, job ID. The git log is the audit trail. Rollback is `git checkout {hash}`. Diff is `git diff`. No custom format needed.

## What It Exposes

### Global Tools
Available to every LLM session regardless of agent. Always injected directly into the session. Never retrieved — always present.

- `list_agents` — read the agent registry (SQLite)
- `create_agent` — register a new agent, `git init` its directory, write initial `state.json`
- `propose_commit` — stage changes to `state.json` and request user approval before committing via go-git
- `update_agent_status` — set active / paused / archived
- `update_last_active` — called by scheduler on job completion

History, diff, and rollback are git operations the LLM calls directly via shell: `git log`, `git diff`, `git show`, `git checkout`.

## Tool Access

Tools live in the object store filesystem. The LLM reads `tool.json` directly. Before any tool executes, the BoundaryValidator checks it against the agent boundary in `state.json` and computes lock status at that moment — not stored on the tool, derived fresh each time.

The same tool may be locked in one agent and open in another — lock status is always a function of the agent's current boundary, not a property of the tool definition itself.

## Boundary Enforcement

The `boundary` section in `state.json` is enforced in code at two layers. Plain-text restrictions are insufficient — the LLM cannot be relied upon to honor prose rules under all conditions.

### Layer 1: BoundaryValidator (runtime)

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
- Does not own shard capability inventory or shard-to-shard delegation — that is the tooling/orchestrator layer

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
