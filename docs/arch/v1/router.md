# Router

## Overview

The router is the central dispatch mechanism. Every heap event passes through it. It owns tool validation, context assembly, and execution dispatch. It is a hybrid — a small set of system primitives are hardcoded, everything else is fully dynamic via Redis.

## Syntax

```
skyra <agent> <tool> [args]
```

The personal agent is also named `skyra`. System tools follow the same syntax:

```
skyra skyra reply "You hit 4 workouts this week"
skyra delegate -gym_agent -home_agent "cancel gym and turn off lights"
skyra delegate report "gym session cancelled"
```

## System Tools

Three hardcoded system tools. Nothing else is hardcoded.

| Tool | Caller | Description |
|---|---|---|
| `skyra delegate` | Skyra | Opens a job, fans out to N agents |
| `skyra delegate report` | Any agent | Reports task result back to delegator |
| `skyra skyra reply` | Skyra only | Sends reply to the user's device |

Only Skyra calls `reply`. Only agents call `delegate report`. Only Skyra calls `delegate`.

## Structure

```
function router(event):
    command = parse(event.payload)
    // command = { agent, tool, args }

    switch command.agent:

        case "skyra":
            switch command.tool:

                case "reply":
                    device = get_user_device(event.session_id)
                    inference = build_reply_inference(command.args)
                    ws.send(device, inference)
                    emit_completed(event)

                default:
                    emit_error(event, "unknown skyra command")

        case "delegate":
            switch command.tool:

                case "delegate":
                    job = dispatcher.create_job(event.turn_id, command.agents)
                    for agent in command.agents:
                        task = dispatcher.create_task(job.id, agent)
                        heap.push(build_agent_job(agent, command.message, task.id))
                    // Skyra is free. dispatcher owns it from here.

                case "report":
                    dispatcher.complete_task(event.task_id, command.args)
                    // dispatcher notifies Skyra incrementally
                    // if all tasks complete → job pops → Skyra composes reply

                default:
                    emit_error(event, "unknown delegate command")

        case redis.get("agent:" + command.agent):
            agent = redis.get("agent:" + command.agent)

            skill = agent.skills.find(command.tool)

            if not skill:
                emit_error(event, "skill not found")
                return

            context = {
                skyra_command: command,
                user_message: event.user_message,
                skills: agent.skills
            }

            inference_job = build_inference(context)
            heap.push(inference_job)

        default:
            emit_error(event, "agent not found in registry")
```

## The Full Loop

```
user: "cancel my gym session and turn off the lights"

Skyra reasons: two domains in play
  → skyra delegate -gym_agent -home_agent "cancel gym, turn off lights" → heap

router: case "delegate"
  → creates job with two tasks
  → heap.push(gym_agent task)
  → heap.push(home_agent task)
  → Skyra is free immediately

gym_agent completes
  → skyra delegate report "gym session cancelled" → heap
  → dispatcher marks task complete
  → notifies Skyra incrementally
  → Skyra: skyra skyra reply "Gym cancelled, still working on lights..." → heap

home_agent completes
  → skyra delegate report "lights off" → heap
  → dispatcher marks task complete
  → all tasks done → job pops
  → Skyra: skyra skyra reply "All done." → heap

router: case "skyra" → reply
  → sends to user's device
```

## Delegate — Pure State Machine

Delegate is a pure state machine. No inference. No reasoning. The agent already reasoned about its own result inside the ReAct loop before reporting back. Delegate just tracks state and routes.

```
task reports success  → mark complete → check if job done
task reports failure  → retry N times with same prompt
retries exhausted     → escalate to Skyra
all tasks complete    → job pops → Skyra notified
```

## Job Tree — Scalable State Machine

The delegate state machine supports arbitrary depth and width. Agents can spawn sub-jobs and replicas of themselves during their ReAct loop. The result is a job tree, not a flat list.

### Three Scaling Dimensions

**Depth — agents spawn sub-jobs**
```
Skyra → delegate → gym_agent
  gym_agent ReAct: "I need calendar data"
    → skyra delegate fan_out -calendar_agent "get schedule"
    → gym_agent pauses, waiting on child job
    → calendar_agent completes → reports to delegate
    → delegate propagates up → gym_agent resumes
```

**Width — arbitrary fan-out at any level**
Any agent can call `skyra delegate fan_out`. Not just Skyra.

**Replicas — agent spawns copies of itself**
```
gym_agent: "10 workout logs to process"
  → skyra delegate fan_out -gym_agent -gym_agent -gym_agent "process batch"
  → 3 replicas run in parallel
  → all complete → parent task unblocks
```

### Schema

```
jobs
  job_id
  parent_task_id    ← null if root (from Skyra), task_id if spawned by agent
  turn_id
  session_id
  status            pending | complete | failed | timed_out
  created_at
  completed_at

tasks
  task_id
  job_id
  agent
  replica_id        ← for parallel instances of same agent
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

Three layers working together:

**Closure table (SQLite)** — stores ALL ancestor-descendant relationships, not just direct parent-child. "Find all tasks under this root job" is one query regardless of depth. No recursive joins.

**Redis atomic counter** — hot path completion check.
```
key: job:{job_id}:pending
value: N    ← DECR on each task completion, atomic, O(1)
```
When counter hits 0 → job complete. No SQL needed for the check.

**SQLite jobs + tasks** — source of truth, durability.

### Completion Propagation

```
task completes:
  1. Redis DECR job:{job_id}:pending
      → hits 0 → job complete
  2. Closure table lookup → find parent task instantly
  3. Redis DECR parent job counter → propagates up the tree
  4. SQLite updated → durable record
  → repeat until root job pops
```

Propagation up the tree is a chain of O(1) Redis operations. SQLite is the safety net underneath. No recursive queries, no round trips per level.

## Skyra Fan-Out

Skyra emits one `delegate` command regardless of how many agents are involved. She identifies which domains are in play, lists them, and fires. The dispatcher handles the rest.

```
skyra delegate -gym_agent -home_agent -home_agent "message"
```

Skyra does not coordinate, wait, or track. She is immediately free for the next user message.

New agent registers in Redis → Skyra sees it in her next context assembly → she can start delegating to it immediately. No router changes needed.

## Two-Level Reasoning

Skyra gets the agent registry in her context — agent names and domain descriptions only. She reasons at domain level and delegates.

Each agent gets its own tool list injected at inference time. It reasons at tool level and issues specific commands.

```
Skyra context:  agent skill registry — { name, domain } per agent
Agent context:  skill list from Redis + user message + Skyra's command
```

## Redis Validation

The router reads Redis directly on every dispatch. No local cache. The registry is always live.

```
redis.get("agent:" + command.agent)
  → returns: { status, shard, location, skills: [{ name, description, args }] }
  → or null → default case → agent not found error
```

## Delegate — Active Coordinator

The delegate agent is not a passive tracker. It actively coordinates task lifecycle.

### Exit Conditions

```
1. All tasks complete successfully  → job pops → Skyra replies
2. Task fails → retry succeeds      → job continues
3. Retries exhausted                → escalate to Skyra → she decides
```

### Reprompt Flow

```
task fails → skyra delegate report "failed: couldn't find lights API"

delegate:
  → validates exit condition → not met
  → reprompts failing agent with context:
      "home_agent: previous attempt failed because X.
       Retry with this context: [failure reason + original request]"
  → agent retries
      → succeeds → job continues
      → fails again → after N retries → escalate to Skyra
```

### Skills

```
fan_out    → dispatch tasks to agents
report     → receive task results
validate   → check exit condition
reprompt   → retry failing agent with adjusted context
escalate   → surface to Skyra when retries exhausted
```

Delegate needs lightweight inference — enough to understand why a task failed and form a useful reprompt. Not a full LLM session.

## ReAct Loop — Per Domain Agent

Every domain agent runs a ReAct (Reason, Act, Observe, Repeat) loop for its assigned task. It is not a single pick-and-call — it is a full reasoning loop that can make multiple skill calls and adapt before reporting back to delegate.

```
Reason  → what do I need to do next?
Act     → call a skill → heap
Observe → result comes back
Repeat  → reason about result, decide next step
Exit    → skyra delegate report "result"
```

### Example — gym_agent cancels a session

```
delegate: "cancel the gym session"

gym_agent ReAct loop:
  Reason:  need to check if a session exists first
  Act:     check_schedule → result: "session at 6pm"
  Observe: session exists
  Reason:  now cancel it
  Act:     cancel_session → result: "cancelled"
  Observe: success
  Reason:  task complete
  Exit:    skyra delegate report "gym session cancelled"
```

### Reprompt Restarts the Loop

If the ReAct loop exits with failure, delegate catches it and restarts the loop with failure context injected. The agent adapts its reasoning on the next pass.

```
ReAct exits: failure
  → delegate reprompts with context
  → new ReAct loop begins
  → agent reasons differently this time
```

## Related

- `docs/arch/v1/capability-model.md` — agent model, registry structure, distributed state
- `docs/arch/v1/command-parser.md` — command syntax and resolution loop
- `docs/arch/v1/scheduler.md` — heap, inference types, job ordering
