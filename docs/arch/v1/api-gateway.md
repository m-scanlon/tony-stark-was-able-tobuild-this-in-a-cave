# API Gateway

## Overview

The API Gateway has exactly two static handlers. Everything else is dynamic.

**Ingress** — static. One entry point. Receives `octos <tool> [args]` commands from shards. Not JSON. Not envelopes. Commands. The gateway parses the command and assembles the job envelope from Redis.

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

## Decomposability

**The gateway is decomposable. The user data is not.**

The gateway holds no state. All state lives in Redis and in the committed graph — both independent of the gateway. The gateway can be taken down, replaced, and published as a new version at any time. The user data is untouched. The graph doesn't care. Redis doesn't care. Everything reconnects.

This is what makes extension possible. A new capability — API compatibility, a new auth method, a new protocol — is a plugin published as a new gateway version. The system brings down the old gateway and publishes the new one. Zero impact to user data.

```
user data (committed graph)   ← sacred, durable, independent
Redis                         ← skill registry, trust boundary, independent
API Gateway                   ← decomposable, stateless, replaceable
```

---

## Plugins

The gateway is extended through plugins — not modified. A plugin is an ingress or egress adapter that translates between an external protocol and the internal `octos <tool> [args]` command syntax.

**API compatibility** is a plugin. An external REST call, a webhook, a CLI invocation — the plugin translates it into a `octos <tool> [args]` command. The gateway resolves it against Redis. The kernel executes it. The plugin knows nothing about what happens after ingress.

```
external REST call
  → API compatibility plugin (ingress adapter)
  → translated to: octos <tool> [args]
  → skill resolution ← Redis
  → trust validation ← Redis
  → kernel executes
  → response → egress adapter → translated back to REST response
```

Publishing a new plugin means publishing a new gateway version. The old gateway comes down. The new gateway comes up. The user data is never involved.

---

## Related

- `docs/arch/v1/kernel.md` — execution boundary, receives job_envelope_v1
- `docs/arch/v1/shard-communication.md` — internal command protocol
