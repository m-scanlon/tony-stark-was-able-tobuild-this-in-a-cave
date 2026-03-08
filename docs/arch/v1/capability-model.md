# Agent and Shard Model

## The Core Distinction

**Shards have capabilities** — hardware and compute. What the device can physically do or run.

**Agents have skills** — functional actions. What the agent can do on behalf of the user.

```
Shard capabilities          Agent skills
──────────────────          ────────────
voice                       send_imessage
display                     turn_on
audio_output                log_workout
fast_reasoning              google_search
deep_reasoning              get_summary
storage                     play_music
```

These are different registries, different concerns, connected through the routing layer.

## Shards

Shards are compute nodes — hardware fingerprint, network address, runtime environment. A shard has no skills. It has capabilities that describe what it can physically support or run.

```
shard: mac_mini
capabilities: [fast_reasoning, storage]

shard: gpu_machine
capabilities: [deep_reasoning]

shard: pi_living_room
capabilities: [voice, audio_output, lightweight_reasoning]

shard: samsung_qled_65
capabilities: [display]
```

Shard capabilities are registered at boot via the hardware fingerprinting + two-round verification process. See `docs/arch/v1/shard-registration.md`.

## Agents

Agents are logical entities. They own skills, state, and memory. They live on shards — their object store is located on a shard — but they are not the shard.

Every shard boots with one system agent representing it. That agent exposes the skills that correspond to the shard's hardware capabilities.

```
shard: samsung_qled_65       ← compute node, display capability
agent: samsung_qled_65       ← system agent, turn_on / turn_off / set_input skills
```

Domain agents expose skills corresponding to their functional domain:

```
agent: gym_agent             ← log_workout / get_summary / plan_session skills
agent: home_agent            ← lights_on / lights_off / set_temp skills
agent: skyra                 ← google_search / imessage / calendar / code_execution skills
```

## Two Registries

### Shard Registry

Tracks hardware capabilities per shard. Used by the routing layer to decide where a skill can execute.

```
key: shard:mac_mini
value: {
  capabilities: ["fast_reasoning", "storage"],
  location: "office",
  status: "active"
}
```

### Agent Skill Registry

Tracks skills per agent. Used by Skyra (reasoning layer) and the router (dispatch layer).

```
key: agent:gym_agent
value: {
  shard: "mac_mini",           ← where object store lives
  location: "*",
  status: "active",
  skills: [
    { name: "log_workout",  description: "Logs a completed workout.", args: ["type", "duration"] },
    { name: "get_summary",  description: "Returns recent workout history." }
  ]
}
```

```
key: agent:samsung_qled_65
value: {
  shard: "samsung_qled_65",
  location: "living_room",
  status: "active",
  skills: [
    { name: "turn_on",   description: "Powers on the display." },
    { name: "turn_off",  description: "Powers off the display." },
    { name: "set_input", description: "Switches the active input source.", args: ["source"] }
  ]
}
```

The model sees only: agent name, location, and skill descriptions. Shard capabilities, execution type, endpoints — all live in separate routing records the model never reads.

## Reasoning Layer

Skyra reasons about agents and skills only. Shard capabilities are a routing concern — they never surface to the reasoning layer.

```
Skyra context:  agent skill registry — { name, domain, skills[] } per agent
Router:         reads shard registry to decide where a skill executes
```

## Reasoning Is a Shard Capability

Reasoning is not a machine role. It is a registered shard capability. Any shard that can run an LLM registers a reasoning capability.

```
shard: gpu_machine     capability: deep_reasoning
shard: mac_mini        capability: fast_reasoning
shard: pi_living_room  capability: lightweight_reasoning
```

A reasoning job is placed on the heap. The estimator scores complexity and routes to whichever shard has the appropriate reasoning capability. Reasoning and execution can happen on the same shard or different shards — always decoupled.

## Commit Proposal Flow

When a reasoning job produces a proposed change to agent state:

```
Agent puts reasoning job on heap
  → heap routes to capable reasoning shard
  → reasoning shard produces commit proposal
  → proposal returns to Skyra agent
  → Skyra: re-reason with pointers | propose commit to user
  → user approves → agent commits
```

Skyra gates all commits. Nothing lands in an agent's object store without her approval.

## Registry Ownership

The brain owns both registries. The elected brain (control plane shard) is the source of truth. Shards register capabilities with the brain at boot. Agents register skills with the brain at creation.

When the brain role shifts to a different shard, both registries transfer with it.

## System Orchestration

The same routing model that dispatches `turn_on` to a TV dispatches a distributed inference job across two GPUs. Shard capabilities determine where work can run. Agent skills determine what work can be done.

```json
{
  "shard": "gpu_cluster",
  "capabilities": ["deep_reasoning"],
  "agents": ["distributed_inference"],
  "model": "deepseek-70b"
}
```

## Capability Compositor (v2)

The brain runs a compositor that watches registered shard capabilities and synthesizes compound capabilities from combinations. Deferred to v2.

## Open Questions

- When two agents at the same location expose the same skill, how does Skyra break the tie?
- How does the shard registry stay consistent when a shard goes offline?
- How does a shard system agent register its skills before the brain's registry exists on first boot?
