# Capability Model

## The Core Distinction

**Shards have capabilities** — hardware and compute. What the device can physically do or run.

**Skills have contracts** — execution requirements. What compute a skill needs to run. Declared in the skill definition. Read by the kernel router at dispatch time.

```
Shard capabilities          Skill contract (compute field)
──────────────────          ──────────────────────────────
voice                       compute: "voice"
audio_output                compute: "audio_output"
fast_reasoning              compute: "fast_reasoning"
deep_reasoning              compute: "deep_reasoning"
storage                     compute: "control_plane"
```

Two separate registries. Two separate concerns. Connected through the kernel router.

## Shards

Shards are compute nodes — hardware fingerprint, network address, runtime environment. A shard has no skills. It has capabilities that describe what it can physically support or run.

```
shard: mac_mini
capabilities: [fast_reasoning, control_plane, storage]

shard: gpu_machine
capabilities: [deep_reasoning]

shard: pi_living_room
capabilities: [voice, audio_output, lightweight_reasoning]
```

Shard capabilities are registered at boot via hardware fingerprinting and verification. See `docs/arch/v1/shard/shard-registration.md`.

## Skill Registry (Redis)

Skills live in Redis. Each skill carries its full contract — roadmap, compute requirements, boundary rules. The kernel router reads the contract to decide where the skill executes.

```
key: skill:log_workout
value: {
  status: "active",
  contract: {
    compute: "fast_reasoning",
    roadmap: [...],
    boundary_rules: {...},
    validation_criteria: {...}
  }
}

key: skill:deep_analysis
value: {
  status: "active",
  contract: {
    compute: "deep_reasoning",
    roadmap: [...],
    boundary_rules: {...},
    validation_criteria: {...}
  }
}
```

## Kernel Router

The kernel router matches skill compute requirements against registered shard capabilities. No hardcoded routing logic. The skill declares where it can run. The router honors it.

```
skill contract: { compute: "deep_reasoning" }  → GPU shard
skill contract: { compute: "voice" }            → Voice Shard
skill contract: { compute: "fast_reasoning" }   → Brain Shard
```

## Reasoning Is a Shard Capability

Reasoning is not a machine role. It is a registered shard capability. Any shard that can run an LLM registers a reasoning capability.

```
shard: gpu_machine     capability: deep_reasoning
shard: mac_mini        capability: fast_reasoning
shard: pi_living_room  capability: lightweight_reasoning
```

A reasoning job is placed on the heap. The kernel router reads the skill's compute requirement and routes to the capable shard. Reasoning and execution can happen on the same shard or different shards — always decoupled.

## Two Registries

### Shard Registry

Tracks hardware capabilities per shard. Used by the kernel router.

```
key: shard:mac_mini
value: {
  capabilities: ["fast_reasoning", "control_plane", "storage"],
  location: "office",
  status: "active"
}
```

### Skill Registry

Tracks skill contracts. Used by the API Gateway (validation) and kernel router (dispatch). Redis is the live registry. Brain is the authority.

## Registry Ownership

The brain owns both registries. The elected brain is the source of truth. Shards register capabilities with the brain at boot. Skills are provisioned into Redis by the kernel on user approval.

When the brain role shifts to a different shard, both registries transfer with it.

## Commit Proposal Flow

When a job produces a proposed change to memory state:

```
job produces commit proposal
  → proposal routed to Skyra's memory namespace
  → Skyra: propose commit to user
  → user approves → kernel commits to memory
```

Skyra gates all commits. Nothing lands in memory without her approval.

## Related

- `docs/arch/v1/kernel.md` — kernel router, execution model, trust boundary
- `docs/arch/v1/shard/shard-registration.md` — capability registration at boot
- `docs/arch/v1/api-gateway/api-gateway.md` — skill validation, job envelope assembly
