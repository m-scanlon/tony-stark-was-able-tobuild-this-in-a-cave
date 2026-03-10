# Skill

A skill is the extendable unit of Skyra. Everything the system does is a skill. Every capability the user gains is a skill. The system grows by gaining skills. The system improves by improving skills.

Skills are not defined. They are learned from intent.

---

## What a Skill Is

A skill is a natural language contract. It declares what work gets done, how it gets done, what it is allowed to touch, and what done looks like. The model interprets and executes it. The skill's capability ceiling is the model's ceiling — not an artificial one the system imposed.

A skill is not code. It is not a function. It is a class — a learned, reusable pattern of intent that the system has crystallized into a permanent, executable contract.

---

## Schema

```
skill {
  id:                   string          // SHA-256(definition_bytes || seed_memory_bytes) — content-addressed, immutable
  model_id:             string          // model under which this skill was committed — trust is model-scoped
  layer:                committed       // skills are always committed — approved by the user
  definition_visibility: open | closed  // creator-set — open: definition readable, closed: encrypted blob

  definition {
    name:               string
    description:        string          // natural language — drives semantic discovery in memory
    roadmap:            []task          // ordered steps — 1-to-many
    boundary_rules:     boundary_rules  // what the skill can and cannot touch
    state_contract:     working | committed
    severity_policy:    severity_policy
    replan_budget:      int             // max replan attempts before escalating. default: 3
    validation_criteria: string         // natural language — what done looks like
    compute_requirements: []capability  // what shard capabilities are needed to run this
    improvement_scope:  improvement_scope  // optional — creator-defined bounded space for self-improvement
  }

  memory {
    namespace:          string          // 1:1 with skill — provisioned alongside it, always
    seed_memory:        bytes           // creator-shipped pre-populated memory. protected by definition_visibility
  }

  provisioning {
    access:             read | write    // consumer-set at provisioning time
    provisioned_at:     timestamp
    signed_by:          user_public_key
    signature:          bytes           // Ed25519 over (id + model_id + access + definition_visibility + provisioned_at)
  }

  history {
    executions:         int             // total times run
    commit_rate:        float           // percentage of outputs committed vs rejected
    user_count:         int             // distinct users who have run it
    versions:           []skill_id      // append-only version chain — old versions never deleted
  }
}
```

---

## Key Properties

**Content-addressed.** The skill ID is the SHA-256 hash of the definition and seed memory combined. The content cannot change without producing a different ID. Immutable by construction.

**1:1 with memory.** Every skill has its own memory namespace, provisioned alongside it atomically. You cannot provision one without the other.

**Natural language.** The definition, roadmap, validation criteria — all prose. The model interprets and executes. The skill improves as the model improves. No redesign needed.

**Trust is model-scoped.** A skill committed under one model is not trusted under a different model. Changing the model flags the skill — visible in memory, not executable in Redis — until the user re-approves under the new model.

**Closed for modification. Open for extension.** A skill is immutable — content-addressed, cannot be changed. The system grows by adding new skills and new versions, never by modifying what's committed. Re-approval under a new model produces a new version. The old version persists. You only ever add.

**Trust is proven at commit time by the owner. Trust is proven to others by history.** The signature proves authenticity. The execution record proves quality. Both are required.

**Deprovisioning does not exist.** Once a skill is provisioned, it stays provisioned. Old versions persist. The version chain is append-only.

---

## Boundary Rules

```
boundary_rules {
  allowed:    []tool        // tools the skill can call freely
  gated:      []tool        // tools requiring runtime approval (allow_once)
  denied:     []tool        // tools the skill can never call
}
```

The kernel enforces boundary rules before every tool dispatch via BoundaryValidator. A skill cannot exceed its declared boundary.

---

## State Contract

- `working` — skill writes to the observational layer freely. No user gate. Skyra's working model.
- `committed` — skill writes require user approval via `propose_commit`. Outputs land in the committed layer.

---

## Severity Policy

How the skill handles assumption failures mid-execution.

```
trivial   → log and continue
minor     → adjust locally and continue
moderate  → attempt fix, then replan
major     → replan remaining steps
critical  → halt and notify user
```

---

## Improvement Scope

Optional. Creator-defined. A bounded observational namespace where Skyra can reason freely about improving the skill's algorithm.

```
improvement_scope {
  aspects:    []string    // what Skyra is allowed to reason about improving
  namespace:  string      // observational scratchpad — Skyra writes freely, not trusted
}
```

Inside the declared aspects: unconstrained reasoning. Outside: nothing. Improvement proposals surface via `propose_commit`. See `docs/arch/v1/skill-improvement.md`.

---

## Discovery vs. Execution

**Discovery** — skills live in memory as vector data. The model searches memory semantically to find relevant skills. No hardcoded tool list. No context injection.

**Execution** — gated by Redis. Even if the model finds a skill in memory and emits `skyra <tool> [args]`, the kernel checks Redis. If the skill is not provisioned, it does not run.

A skill can be visible in memory — discoverable, reasoned about — without being executable. Redis is the gate, not memory.

---

## Lifecycle

```
intent expressed → no skill found → intent namespace provisioned (observational)
  → intent recurs across sessions
  → recurrence_ratio = intent_session_count / total_session_count crosses threshold
  → Skyra proposes skill in natural language
  → user approves
  → skill node committed (new id, model_id stamped)
  → memory namespace provisioned alongside it
  → provisioned in Redis
  → executable
```

See `docs/arch/v1/skill-lifecycle.md` for the full lifecycle.

---

## Versioning

Skills are immutable. Modification means a new skill — new content hash, new model_id. The old version persists. `update_skill` is the only path to producing a new version.

```
skill_v1  →  skill_v2  →  skill_v3
  (all exist, append-only version chain)
```

---

## System Primitive Skills

Pre-provisioned in Redis at boot. The system cannot function without these. They are skills — same schema, same rules.

| Skill | Purpose |
|---|---|
| `reply` | Send a reply to the user's device. Only Skyra calls this. |
| `fan_out` | Open a job, fan out to N target domains. |
| `report` | Report task result back to delegate. Any task can call this. |
| `chat` | A conversation with the user. Every session is a job. |
| `reasoning` | Background cron job. Session history + VAD → observational nodes + edges. |
| `integrate` | Connect reasoning output to the existing graph. Alias resolution, weight updates, new edges. |
| `update_skill` | The only path to modifying a skill node. Requires user approval. |
| `commit` | Write to memory (user-gated). |
| `propose_commit` | Surface a commit proposal to the user. |
| `search` | Semantic search in memory — retrieval and signal. |
| `provision_memory` | Create a new memory namespace. |
| `provision_skill` | Add a skill to Redis. |

---

## Related

- `docs/arch/v1/skill-lifecycle.md` — how skills are born, crystallize, and version
- `docs/arch/v1/skill-improvement.md` — bounded self-improvement loop
- `docs/arch/v1/skill-reasoning.md` — reasoning primitive
- `docs/arch/v1/skill-integrate.md` — integrate primitive
- `docs/arch/v1/skill-update-skill.md` — update_skill primitive
- `docs/arch/v1/crypto-protocol.md` — signing, trust model, definition visibility
- `docs/arch/v1/memory-structure.md` — observational vs committed layer
- `docs/arch/v1/kernel.md` — execution, Redis trust boundary, job tree
