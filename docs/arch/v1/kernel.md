# Kernel

> Previously named "Router." All router content is canonical here.

## Overview

The kernel is Skyra's central execution boundary. Every event passes through it. It owns trust enforcement, skill invocation, job instantiation, task execution, memory provisioning, and pattern recognition. It does not own time — the Cron Service is the time-aware layer.

**The kernel is purely reactive.** It processes events. It never self-initiates. A command only fires as the result of another event, or as an event emitted by the Cron Service.

**The two-sentence mental model:**

> The kernel owns all execution. Every API call goes through the kernel. If it's not a kernel-executable command, it's reasoning happening inside a shard — until that shard produces `skyra <tool> [args]`.

Shards reason. The kernel executes. There is no other mode.

**Two axioms:**

> Everything is a Shard. (hardware layer)
> Everything is a Skill. (execution layer)

Every operation the system performs — including its own internal operations — is expressed as a skill. The three hardcoded primitives (`reply`, `fan_out`, `report`) are skills. Pattern recognition crystallizes into skills. Memory provisioning is triggered by a skill. Cron executes skills. There is no operation that is not a skill.

---

## Syntax

```
skyra <tool> [args]
```

This is the standard tool call the API Gateway emits. One prefix. All tools. The kernel resolves `tool` against Redis and dispatches.

```
skyra reply "You hit 4 workouts this week"
skyra fan_out -gym -home "cancel gym and turn off lights"
skyra report "gym session cancelled"
skyra check_nginx
skyra log_workout "chest day"
```

---

## The Stack

```
                             API Gateway          ← doc: api-gateway.md
                             (ingress / egress)
                                    ↓
              ┌─────────────┬───────┴────────┐
           Kernel        Kernel           Kernel   ← one per shard
              │              │               │
           Memory         Memory          Memory
              │              │               │
           Shard          Shard           Shard
              │
           [Cron Service]                         ← standalone shard service
```

The Cron Service is a standalone service running on a shard — provisioned by skyrad like any other service. It fires scheduled events through Ingress. The kernel has no idea they came from a cron — it just sees events.

The kernel receives `job_envelope_v1` from the API Gateway — already trusted, already dispatched to the right shard. The kernel executes. See `docs/arch/v1/api-gateway.md` for gateway design and job envelope schema.

---

## Primitives

**Skill** — learned class. A roadmap: 1-to-many tasks. Skills are learned by the kernel's pattern recognition function from observational streams — not manually defined. Registered in Redis (the trust boundary). Lives in memory.

The skill definition is the execution contract. It carries:

- **Tasks** — the ordered steps to execute (1-to-many)
- **Boundary rules** — which tools are allowed, which require approval, which are denied. Enforced by BoundaryValidator before each tool dispatch. Permission prompt: `allow_always | allow_once | deny`.
- **State contract** — whether the skill writes to working state (scratch, free, no approval) or committed state (user-gated, canonical, versioned). Working state is the scratch pad. Committed state requires `propose_commit`.
- **Severity policy** — how assumption failures are handled: trivial (log and continue), minor (adjust locally), moderate (attempt fix then replan), major (replan remaining steps), critical (halt and notify user).
- **Replan budget** — max replan attempts before escalating. Default: 3.
- **Preemption** — implicit. The kernel re-queues the job after each tool call. Higher priority work gets picked up first. The job waits. Resume is seamless — the context blob is the job state.
- **Improvement scope** — optional. Creator-defined. A bounded observational namespace where Skyra can reason freely about improving the algorithm. Scope defines what aspects she is allowed to reason about. Inside the scope: unconstrained reasoning. Outside it: nothing. Improvement proposals surface via `propose_commit`. See `docs/arch/v1/skill-improvement.md`.

The kernel runs the skill. The skill defines how.

**Job** — skill instance. Created when the kernel invokes a skill. Holds a `skill_id` reference. Same relationship as class → object in programming.

**Task** — execution unit. Atomic work inside a job. A job expands into a task tree.

**Memory** — provisioned namespace. Contains entities. Kernel provisions on user approval. The kernel does not own what lives inside memory — memory owns its own contents.

**Entity** — named things that live in memory. People, places, tools, concepts. Accumulated through observation. Memory is their home — not the kernel.

**Shard** — hardware node. Registered capabilities (voice, deep_reasoning, etc.). The kernel dispatches tasks to shards based on capability matching.

---

## Trust Boundary

**The cache (Redis) is the trust membrane.** Skills in memory are inert — roadmaps. They cannot execute. The only path to execution is through the cache. If a skill is not provisioned in Redis, the kernel has nothing to run.

When the kernel resolves a skill:

| State | Action |
|---|---|
| Not in cache | Deprovisioned or never trusted → permission request or hard denial → back to heap |
| In cache, gated | Needs runtime approval → BoundaryValidator → user approval → back to heap |
| In cache, trusted | Kernel instantiates job, expands into tasks, dispatches |

---

## Kernel Structure

```
function kernel(event):
    command = parse(event.payload)
    // command = { tool, args }

    switch command.tool:

        case "reply":
            device = get_user_device(event.session_id)
            inference = build_reply_inference(command.args)
            ws.send(device, inference)
            emit_completed(event)

        case "fan_out":
            job = create_job(event.turn_id, command.targets)
            for target in command.targets:
                task = create_task(job.id, target)
                heap.push(build_task_event(target, command.message, task.id))
            // Skyra is free. delegate state machine owns it from here.

        case "report":
            complete_task(event.task_id, command.args)
            // delegate notifies Skyra incrementally
            // all tasks complete → job pops → Skyra composes reply

        case redis.get("skill:" + command.tool):
            skill = redis.get("skill:" + command.tool)

            if not skill:
                emit_error(event, "skill not found or not trusted")
                return

            job = instantiate_job(skill, event)
            tasks = expand_to_tasks(skill, job)
            heap.push(tasks)

        default:
            emit_error(event, "skill not found")
```

---

## The Full Loop

```
user: "cancel my gym session and turn off the lights"

API Gateway resolves: two domains in play
  → skyra fan_out -gym -home "cancel gym, turn off lights"

kernel: case "fan_out"
  → creates job with two tasks
  → heap.push(gym task)
  → heap.push(home task)
  → Skyra is free immediately

gym task completes
  → skyra report "gym session cancelled"
  → delegate marks task complete
  → notifies Skyra incrementally
  → Skyra: skyra reply "Gym cancelled, still working on lights..."

home task completes
  → skyra report "lights off"
  → delegate marks task complete → all tasks done → job pops
  → Skyra: skyra reply "All done."

kernel: case "reply"
  → sends to user's device
```

---

## System Primitives

Every command resolves against Redis. No skill bypasses the registry — not even system primitives. If a skill is not in Redis, it does not run.

All system primitive skills are pre-provisioned in Redis at boot.

| Skill | Purpose |
|---|---|
| `reply` | Sends reply to user's device. Only Skyra calls this. |
| `fan_out` | Opens a job, fans out to N target domains. |
| `report` | Reports task result back to delegate. Any task can call this. |
| `chat` | A conversation with the user. Every session is a job. Opens on first turn, closes on session end. |
| `reasoning` | Background job triggered by cron. Decomposes session history + VAD into observational nodes, then writes edges to the graph. |
| `integrate` | Connects the mini graph from reasoning to the existing graph. Finds aliases, updates weights, adds missing edges. |
| `update_skill` | The only path to modifying a skill node. Requires user approval. |
| `commit` | Write to memory (user-gated) |
| `propose_commit` | Surface a commit proposal to the user |
| `search` | Semantic search in memory — retrieval and signal |
| `provision_memory` | Create a new memory namespace |
| `provision_skill` | Add a skill to Redis |

---

## Delegate — Pure State Machine

Delegate coordinates multi-skill execution. No inference. No reasoning. Tasks reason about their own results in their ReAct loops before reporting back. Delegate just tracks state and routes.

```
task reports success  → mark complete → check if job done
task reports failure  → retry N times with failure context injected
retries exhausted     → escalate to Skyra
all tasks complete    → job pops → Skyra notified
```

### Exit Conditions

```
1. All tasks complete successfully  → job pops → Skyra replies
2. Task fails → retry succeeds      → job continues
3. Retries exhausted                → escalate to Skyra → she decides
```

### Reprompt Flow

```
task fails → skyra report "failed: couldn't find lights API"

delegate:
  → validates exit condition → not met
  → reprompts failing task with context:
      "previous attempt failed because X. Retry with this context."
  → task retries
      → succeeds → job continues
      → fails again → after N retries → escalate to Skyra
```

---

## ReAct Loop — Per Task

Every task runs a ReAct loop. Not a single pick-and-call — a full reasoning loop that can make multiple skill calls before reporting back.

```
Reason  → what do I need to do next?
Act     → call a skill → heap
Observe → result comes back
Repeat  → reason about result, decide next step
Exit    → skyra report "result"
```

Example:

```
fan_out: "cancel the gym session"

gym task ReAct loop:
  Reason:  need to check if a session exists first
  Act:     skyra check_schedule → result: "session at 6pm"
  Observe: session exists
  Reason:  now cancel it
  Act:     skyra cancel_session → result: "cancelled"
  Observe: success
  Exit:    skyra report "gym session cancelled"
```

If the loop exits with failure, delegate reprompts with failure context. The task restarts its loop with adjusted reasoning.

---

## Job Tree

Tasks can spawn sub-jobs and replicas during their ReAct loops. The result is a tree, not a flat list.

### Three Scaling Dimensions

**Depth — tasks spawn sub-jobs**
```
fan_out → gym task
  gym ReAct: "I need calendar data"
    → skyra fan_out -calendar "get schedule"
    → gym task pauses, waiting on child job
    → calendar completes → reports to delegate
    → delegate propagates up → gym resumes
```

**Width — arbitrary fan-out at any level**
Any task can call `skyra fan_out`. Not just Skyra.

**Replicas — task spawns copies of itself**
```
gym task: "10 workout logs to process"
  → skyra fan_out -log_workout -log_workout -log_workout "process batch"
  → 3 replicas run in parallel
  → all complete → parent task unblocks
```

### Schema

```
jobs
  job_id
  skill_id          ← reference to parent skill
  parent_task_id    ← null if root, task_id if spawned by a task
  turn_id
  session_id
  status            pending | complete | failed | timed_out
  created_at
  completed_at

tasks
  task_id
  job_id
  skill
  replica_id        ← for parallel instances of same skill
  status            pending | complete | failed
  result
  created_at
  completed_at

job_tree_closure
  ancestor_id
  descendant_id
  depth             ← 0 = self, 1 = direct parent, etc.
```

### Data Structures

**Closure table (SQLite)** — all ancestor-descendant relationships. "Find all tasks under this root job" is one query regardless of depth. No recursive joins.

**Redis atomic counter** — hot path completion check.
```
key: job:{job_id}:pending
value: N  ← DECR on each task completion, atomic, O(1)
```
When counter hits 0 → job complete. No SQL needed for the check.

**SQLite jobs + tasks** — source of truth, durability.

### Completion Propagation

```
task completes:
  1. Redis DECR job:{job_id}:pending → hits 0 → job complete
  2. Closure table lookup → find parent task instantly
  3. Redis DECR parent job counter → propagates up the tree
  4. SQLite updated → durable record
  → repeat until root job pops
```

Propagation up the tree is a chain of O(1) Redis operations. No recursive queries.

---

## Skill Discovery vs. Skill Execution

**Two separate layers. Two separate concerns.**

**Discovery** — Skills live in memory as vector data. Indistinguishable from any other piece of data. The LLM searches memory semantically to find relevant skills. No hardcoded tool list. No context injection. The model reasons about what tools exist the same way it reasons about any other fact — by searching.

**Execution** — gated by Redis. Even if the model finds a skill in memory and emits `skyra <tool> [args]`, that command hits the kernel. The kernel checks Redis. If the skill is not provisioned, it doesn't run. End of story.

```
Memory (discovery)          Redis (execution gate)
──────────────────          ─────────────────────
skill: log_workout          skill:log_workout → { status: active, ... }
skill: check_nginx          skill:check_nginx → null (deprovisioned)
skill: deep_analysis        skill:deep_analysis → { status: active, ... }

LLM finds log_workout in memory via semantic search
  → emits: skyra log_workout --type=run --duration=30
  → kernel checks Redis: skill:log_workout → present, trusted
  → executes

LLM finds check_nginx in memory via semantic search
  → emits: skyra check_nginx
  → kernel checks Redis: skill:check_nginx → null
  → rejected (visible in memory, not executable)
```

A skill can exist in memory — be known, discoverable, semantically queryable — without being executable. Deprovisioned skills are still in memory. The LLM can reason about them ("you used to have check_nginx"), but cannot invoke them. Redis is the gate, not memory.

This eliminates the tool injection problem entirely. The context window doesn't carry a hardcoded tool list. The model discovers what's available the same way it discovers anything else — by searching its own memory.

---

## Two-Level Reasoning

Skyra discovers relevant skills by searching memory semantically — skill descriptions, domains, past usage. She reasons at domain level, selecting which domain to fan out to.

Each task discovers its own relevant tools by searching memory — skills matching its specific workload. It reasons at skill level.

```
Skyra:  memory search → semantically relevant skills → domain selection → fan_out
Task:   memory search → tools matching current step → select → emit skyra <tool> [args]
```

Execution is always gated by Redis. Discovery is always via memory.

---

## Redis Validation

The kernel reads Redis directly on every dispatch. No local cache. The registry is always live.

```
redis.get("skill:" + command.tool)
  → returns: { status, shard, location, tasks: [{ name, description, args }] }
  → or null → default case → skill not found error
```

---

## Memory Provisioning

Memory is not created manually. It emerges from observation. The kernel's pattern recognition function watches the observational streams. When an entity accumulates enough signal, the kernel emits a provisioning event — the user decides.

```
Cron Service fires snapshot event → heap
  → kernel evaluates observational state
  → threshold crossed → provisioning event emitted → heap
  → user approval requested
  → approved → kernel provisions memory namespace
  → denied → threshold resets, observation continues
```

---

## Skill Learning

Skills are learned, not defined. The kernel's pattern recognition function watches behavioral patterns across the observational streams. When a pattern crosses a frequency × affect threshold:

```
pattern crosses threshold
  → skill crystallized (roadmap: 1-to-many tasks)
  → provisioned in Redis (trust established)
  → skill lives in memory
  → available for invocation
```

---

## Job Execution Model

### From Command to Heap

```
skyra <tool> [args] + credentials arrive at API Gateway
    ↓
Redis check: skill exists AND shard is authorized
    ↓
No  → rejected
    ↓
Yes → Redis returns the full skill definition
    ↓
command args + full skill enter the heap as a job
    ↓
kernel router reads the skill's contract
    ↓
routes to the right shard based on compute/capability requirements
    ↓
shard executes
```

Redis does not just approve — it returns the entire skill on approval. The job on the heap carries both the command (args) and the full skill definition (roadmap + contract). The skill is the routing manifest.

### The Skill Contract

Every skill shares the same schema. The contract declares:

- **Roadmap** — the ordered tasks (1-to-many)
- **Compute requirements** — what kind of shard can run this (voice, deep_reasoning, control_plane, etc.)
- **Execution contract** — boundary rules, severity policy, state contract, replan budget
- **Validation criteria** — what "done" looks like

The kernel router reads the compute requirements and routes accordingly. No hardcoded routing logic — the skill itself is the manifest.

### The Kernel Is the Heap and the Router

The kernel owns both. The heap is the kernel's queue. The router is the kernel's dispatch logic. They are not separate services.

### Execution Loop

**Security is front-loaded.** Once past Redis, trust is established. No further security checks inside the kernel.

**State travels with the job.** The job carries its own context on the heap — routing instructions and a growing state blob. Every task execution appends to the context. The heap is a queue of stateful jobs, not raw work items.

```
job on heap (command args + full skill)
    ↓
kernel router reads skill contract → routes to shard
    ↓
task executes
    ↓
context grows
    ↓
job re-queues on heap
    ↓
repeat until self-validation passes
    ↓
skyra reply "work done"
    ↓
routed to Skyra's memory namespace
```

**Self-validation.** The skill defines what "done" means. The job validates against its own contract. No external validator.

**`skyra reply` routes to Skyra's memory.** Completed work is committed to Skyra's memory namespace. Skyra decides what to surface to the user.

---

## Kernel Internal Methods

These are not external interfaces — they are kernel functions.

- **Estimation** — reads job complexity, matches against shard capability profiles, selects target shard. Fires as an inference call when an estimation heap item is picked up.
- **Execution** — heap-driven loop. Reads the skill's execution contract (boundary rules, severity policy, state contract, replan budget). State travels with the job. Preemption is free between tasks — higher priority jobs get picked up first, the current job waits on the heap.
- **Pattern recognition** — watches the four observational streams. Crystallizes skills. Emits memory provisioning events. Triggered on schedule by the Cron Service.

Detail on heap and scheduling: `docs/arch/v1/scheduler.md` (kernel internal).

---

## Cron Service

A standalone service running on a shard — provisioned by skyrad like any other shard service. The only component that knows about time.

The Cron Service executes skills on a schedule. It fires `skyra <tool> [args]` through Ingress at configured intervals. The kernel receives them as ordinary events — has no idea they came from a cron. The Cron Service is the invoker for system skills. Skyra is the invoker for domain skills.

The primary scheduled skill is `reasoning` — fires when the user is offline, reads unprocessed session history + VAD, and produces observational nodes and edges. Exact schedule TBD. See `docs/arch/v1/skill-reasoning.md`.

---

## Terminology Layer

Canonical primitive names are the default. Users configure alternate labels in `skyra.user`. Every label-emitting surface reads from this config.

| Default | Programmer |
|---|---|
| Skill | Class |
| Job | Instance |
| Memory | Repo |
| Task | Task |

Execution model is unchanged. Only the labels differ.

---

## Related

- `docs/arch/v1/executor.md` — execution loop detail (kernel internal)
- `docs/arch/v1/scheduler.md` — heap design (kernel internal)
- `docs/arch/v1/predictive-memory.md` — observational streams, pattern recognition (kernel function)
- `docs/arch/v1/capability-model.md` — shard capability registration
- `docs/arch/v1/skill-lifecycle.md` — full skill lifecycle: observation → intent namespace → validation → skill building → provisioned
- `docs/arch/v1/gaps.md` — open gaps: Cron Service design, memory provisioning flow, skill learning thresholds, terminology layer implementation
