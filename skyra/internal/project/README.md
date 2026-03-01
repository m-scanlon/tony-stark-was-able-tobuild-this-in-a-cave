# Project Service

## What It Is

The Project Service is the single owner of all project state in Skyra. Nothing reads or writes project state without going through this service. It is the foundational layer that all other services interface with.

Note: earlier architecture docs refer to this as the **Memory Service**. The Project Service supersedes that definition with a more clearly scoped responsibility model.

## Core Principles

- **Single source of truth** — all project data lives in exactly one canonical location. The AI never edits state directly. All changes occur through explicit commits.
- **Append-only commit history** — project changes are immutable commit objects. Each commit contains a full snapshot, links to its parent, and can be used for rollback.
- **Vector index is derived data** — the vector DB is not the source of truth. It indexes project state for semantic retrieval and can be deleted and rebuilt at any time.

## What It Owns

**Project Registry (SQLite)**
A lightweight index of all projects. Used by the context engine as a fast first gate before any deeper retrieval.
- Project identity and metadata
- Status: `active | paused | archived`
- Last active timestamp

**Object Store Interface**
The only entry point for reading and mutating project state on disk. Filesystem in Phase 1, S3/MinIO in Phase 2 with no structural changes required.

- Current project state (knowledge, artifact mapping, boundary definition)
- Immutable commit history
- HEAD pointer management
- Rollback and audit trail

**Local Tool Registry (Vector DB)**
Each project has a set of local tools registered and indexed semantically. The tool description is embedded and indexed — vector search runs against these embeddings using the current request as the query. Tools are retrieved by semantic similarity, not injected wholesale.

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
  "project_id": "jarvis",
  "categories": ["filesystem_write"],
  "requires_approval": true
}
```

`score` is the semantic similarity score from the vector search. The Project Service applies a score threshold before hydration — results below the threshold are dropped. Lock status is not stored on the tool record. It is computed at hydration time by joining the tool record against the project boundary in `state.json`.

## Object Store Structure

```
.skyra/projects/{project_id}/
  HEAD.json            ← pointer to current commit
  state.json           ← materialized current state (four sections below)
  commits/             ← immutable commit history
    {commit_id}.json
  jobs/
    {job_id}/
      envelope.json    ← job_envelope_v1
      tasks/
        {task_id}/
          tasksheet.json  or  workplan.json
          notes.md
```

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
    "scope": "Human-readable description of this project and what Skyra is allowed to do in it.",
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
  "project": "jarvis",
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
Available to every LLM session regardless of project. Always injected directly into the session. Never retrieved — always present.

- `get_project_state` — read current state for a project
- `get_project_facts` — read facts and assumptions
- `get_commit_history` — read the commit log
- `list_projects` — read the registry
- `get_job_tasks` — read tasks under a job
- `create_project` — register a new project
- `propose_commit` — submit a patch against project state
- `apply_commit` — apply an approved patch and update HEAD
- `rollback_commit` — revert to a previous commit
- `update_project_status` — set active / paused / archived
- `update_last_active` — called by scheduler on job completion

### Local Tool Registry
Per-project tools registered at project setup. Retrieved via vector search — not injected directly. The LLM sees all retrieved tools, including locked ones, with their access status attached.

## Tool Hydration

Tool hydration is the intermediary step between raw vector DB results and what the LLM session receives. The Project Service runs hydration on every tool in the result set before returning anything to the Domain Expert.

Hydration joins the raw tool record against the project boundary in `state.json` and computes the `access` field:

```json
{
  "id": "tool-uuid",
  "score": 0.92,
  "name": "write_file",
  "description": "Writes content to a file at the specified path.",
  "input_schema": { ... },
  "project_id": "jarvis",
  "categories": ["filesystem_write"],
  "requires_approval": true,
  "access": {
    "status": "locked",
    "reason": "Filesystem writes are restricted for this project.",
    "enforcement": "prompt_user"
  }
}
```

`access.status` is either `allowed` or `locked`. Lock status is always derived fresh at hydration time — it is never stored on the tool record itself, because the same tool may be locked in one project and open in another.

The LLM receives the full hydrated list. It can see and reason over all tools — including locked ones — but the BoundaryValidator enforces what actually executes.

## Boundary Enforcement

The `boundary` section in `state.json` is enforced in code at two layers. Plain-text restrictions are insufficient — the LLM cannot be relied upon to honor prose rules under all conditions.

### Layer 1: Hydration (retrieval time)

After vector search returns results, the Project Service hydrates each tool with its access status for the current project. No tools are hidden. The LLM sees everything, with locked tools clearly marked. This gives the LLM full situational awareness — it knows what exists and what it can use.

### Layer 2: BoundaryValidator (runtime)

Before any tool call is dispatched to execution, the BoundaryValidator checks the proposed call against the project boundary. This is pure code — no LLM in the loop.

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
| `allow_always` | Tool is permanently unlocked for this project. An immediate commit is written to `state.json` before execution resumes. | Yes — boundary updated via commit |
| `allow_once` | This single invocation is allowed for the stated reason. Boundary is unchanged. | No — ephemeral |
| `deny` | Tool call is blocked. Skyra replans around the denied tool and attempts to complete the task another way. | No — ephemeral |

### What `scope` and `description` Are For

`boundary.scope` and each restriction's `description` are human-readable fields. They may be included in system prompts as a soft context hint, but enforcement never relies on them. They exist for human auditing, project setup review, and operator understanding — not for controlling LLM behavior.

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

- Context engine — reads registry and project state
- Scheduler — updates `last_active_at` on job completion
- Orchestrator — applies task outcomes via commits
- LLM session — calls global tools directly
- Domain Expert — reads project state during task formation

## Related Docs

- `docs/arch/v1/scyra.md` — system architecture and canonical pipeline
- `docs/arch/v1/agents-services.md` — full service catalog
- `docs/arch/v1/task-formation.md` — how tasks are formed against project state
- `docs/arch/v1/domain-expert/README.md` — domain expert and plan approval gate
- `skyra/internal/scheduler/README.md` — scheduler service
