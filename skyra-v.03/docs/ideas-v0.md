# Ideas v0

Speculative ideas. Not canon. Not in scope for current iteration.

---

## Relationship Worlds

Instead of an exchange stack at the top of a peer relationship — a world.

A being's peer isn't just an exchange stack holding conversation history. It's a world of beings that grew out of that relationship.

The exchange stack is the conversation history.
The world is what got built from it.

Two beings relate. Through that relating, new beings emerge — concepts, skills, shared understanding. Those live in a world that belongs to that relationship. Not a flat stack of impulses. A populated space.

This is differentiation at the relationship level. Not just at the being level.

Each being has their own world per peer. The kernel sees all worlds. But each relationship grows its own.

---

## Inference Retry On Rate Limit

When a 429 is returned from the inference provider, the response body includes a `retryDelay` field indicating exactly how long to wait. Instead of failing hard, the runner should parse that delay and retry the call once after sleeping that duration.

This keeps the cognitive chain alive through temporary quota exhaustion rather than dropping the signal.

---

## Routing To Claude Execution Surfaces

If Skyra is linked to Claude Code, a Claude Code instance should not be treated
as a being.

It should be treated as an execution surface.

That keeps the ontological layer clean.

The being remains Skyra.

The Claude worker is where one line of work for Skyra happens right now.

### The Core Distinction

Do not route to "the right thread."

Route to the right execution surface.

The kernel should own a binding like:

```text
exchange_id -> surface_id
surface_id -> {
  session_id,
  process_handle,
  cwd,
  worktree,
  writable
}
```

So if Skyra has two Claude Code instances open, there is no ambiguity.

Each instance has its own `surface_id`.

Each live line of work is bound to one surface.

### The Practical Rule

One exchange binds to one surface at a time.

If the surface is alive, the next turn goes to that surface.

If the surface dies, the kernel respawns the worker and reattaches to that
surface's `session_id`.

If Skyra wants a branch rather than a continuation, the kernel forks the Claude
session and creates a new `surface_id`.

### Why This Matters

This means Claude sessions are not Skyra's identity.

They are temporary operational placements.

Skyra's continuity lives in the kernel-owned relationship and exchange state,
not in the worker process.

### Operational Invariants For Coding Work

- one writable Claude surface per repo/worktree
- one active exchange per writable surface
- if parallel coding is needed in the same repo, use separate worktrees
- if two simultaneous coding lines become semantically different enough, that
  may be differentiation rather than one being with two surfaces

### Consequence

Skyra does not need to remember "which terminal window."

She needs the kernel to remember:

- which exchange is active
- which execution surface currently carries it
- whether that surface is continuing, branched, or dead

That is enough to route the next turn correctly.
