# Claude Code As Execution Surface

A Claude Code instance is not a being. It is an execution surface.

The being remains Skyra. The Claude Code worker is where one line of work for Skyra happens right now. Treating it as a being pollutes the ontological layer — beings have identity, purpose, relationships. A worker process has none of that. It is a surface: a place where work lands.

## Surface Binding

The kernel owns the binding:

```
exchange_id → surface_id
surface_id → {
  session_id,
  process_handle,
  cwd,
  worktree,
  writable
}
```

One exchange binds to one surface at a time. If the surface is alive, the next turn goes there. If the surface dies, the kernel respawns the worker and reattaches via `session_id`. If Skyra wants a branch rather than a continuation, the kernel forks the session and creates a new `surface_id`.

## Operational Invariants

- One writable Claude surface per repo/worktree
- One active exchange per writable surface
- Parallel coding in the same repo uses separate worktrees
- If two simultaneous coding lines diverge enough, that may be differentiation — a new being, not a second surface

## What This Means

Skyra's continuity lives in the kernel-owned exchange and relationship state, not in the worker process. Claude sessions are temporary operational placements. The kernel tracks which exchange is active, which surface currently carries it, and whether that surface is continuing, branched, or dead. That is enough to route the next turn correctly.

## Relationship To Execution Surface

This is a specific surface type that the execution surface model already supports — a process surface whose adapter wraps a Claude Code session. No new machinery. The adapter handles session management. The router sees a process surface like any other.
