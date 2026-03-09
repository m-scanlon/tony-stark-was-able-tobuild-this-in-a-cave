# API Gateway

## Overview

The API Gateway has exactly two static handlers. Everything else is dynamic.

**Ingress** — static. One entry point. Receives `skyra <tool> [args]` commands from shards. Not JSON. Not envelopes. Commands. The gateway parses the command and assembles the job envelope from Redis.

**Egress** — static. One exit point. All responses leave here as commands — same syntax, other direction.

**Everything between** — dynamic. Skill resolution, trust validation, dispatch, shard routing — all read from Redis at runtime. No hardcoded routes. No hardcoded shards. A new shard registers in Redis → the gateway starts routing to it immediately. A new skill is provisioned → the gateway resolves it immediately. Zero gateway code changes required.

See `docs/arch/v1/shard-communication.md` for the full communication protocol.

---

## Flow

```
Ingress (static)
    ↓
Skill resolution   ← Redis
    ↓
Trust validation   ← Redis (provisions + security metadata)
    ↓
Dispatch           ← Redis (which shards can run this)
    ↓
Route to shard kernel
    ↓
Egress (static)    ← response returns here
```

The job arrives at the kernel already on the correct shard. The kernel does not make routing decisions — it executes.

---

## Job Envelope

The job the API Gateway assembles and hands to the kernel:

```
job_envelope_v1 {
  job_id
  skill_id          ← the skill being instantiated
  provisions        ← what this skill is allowed to do (from Redis)
  security {
    trust_level
    boundary_rules  ← allowed tools, denied tools, approval-required tools
    permissions[]
  }
  payload {
    roadmap         ← the skill's task definitions (1-to-many tasks)
    args            ← arguments from the incoming event
    turn_id
    session_id
  }
}
```

Dispatch is resolved by the API Gateway before the envelope is sent. It does not travel with the job — the kernel receives work, not routing instructions.

---

## The Full Stack

```
                            API Gateway
                            Ingress (static)      ← all events enter here
                                   ↓
                          [skill resolution]
                          [trust validation]  ← Redis
                          [dispatch]
                                   ↓
              ┌────────────┬───────┴────────┐
           Kernel       Kernel           Kernel      ← one per shard
              │             │                │
           Memory        Memory           Memory
              │             │                │
           Shard          Shard            Shard
              │
           [Cron Service]                           ← standalone shard service
                                   ↓
                            API Gateway
                            Egress (static)       ← all responses leave here
```

The Cron Service is a standalone service running on a shard — provisioned by skyrad. It fires scheduled events through Ingress like any other event source. Ingress does not distinguish.

---

## Related

- `docs/arch/v1/kernel.md` — execution boundary, receives job_envelope_v1
- `docs/arch/v1/gaps.md` — G32: API Gateway domain resolution not designed
