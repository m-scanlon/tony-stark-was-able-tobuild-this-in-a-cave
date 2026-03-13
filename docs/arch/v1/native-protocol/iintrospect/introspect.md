# Introspect

The `introspect` primitive gives the system the ability to reason about itself — RAM, disk, graph size, process state. It is the only path to shell execution in Skyra.

Shell access is a dangerous design choice. The trust model below is the contract that makes it acceptable.

---

## Why It Exists

The system needs to observe its own resource state to make decisions:

- `system_health` reads RAM and disk to detect pressure
- `migrate_data` checks available capacity on target shards
- `emergency_offload` verifies disk state before triggering

Without `introspect`, these skills are blind. They cannot reason about the machine they are running on.

---

## The Trust Model

Every introspect call is validated against the active execution attempt, not the abstract job.

```
(attempt_id, job_id, pid, spawn_token, command)
```

| Field | Source | What it proves |
|---|---|---|
| `job_id` | SQLite `jobs` | The job is real and was legitimately created |
| `attempt_id` | SQLite `job_attempts` | This specific execution instance is the active attempt |
| `pid` | OS | A live process exists and matches this attempt |
| `spawn_token` | Kernel-issued nonce (raw, held by child process only) | This specific process was spawned by this kernel for this attempt |
| `command` | Redis | This command was approved for this attempt |

All must align. A call that fails any layer is rejected.

**Execution record, not job record:**

Validation is against the active attempt. Retries spawn a new attempt with a new PID and a new spawn token. The previous attempt is closed. No ambiguity about which execution instance is live.

```
job_attempts
  attempt_id
  job_id
  pid
  pid_start_time       ← process start time from OS. Linux: /proc. macOS: proc_pidinfo.
  spawn_token_hash     ← hash of the token, never the raw token
  spawned_at
  ended_at
  state                active | completed | failed | terminated
```

`pid_start_time` is a secondary defense against PID reuse. The spawn token already makes reuse effectively impossible — a recycled PID has no knowledge of the token. `pid_start_time` adds an independent OS-sourced check: even if the token were somehow leaked, a new process with a recycled PID would have a different start time. Two independent reuse guards.

**Spawn token is a secret capability:**

The raw token is passed only to the child process at spawn. SQLite stores the hash. Validation is `hash(presented_token) == spawn_token_hash`. If the database leaks or logs are exposed, the stored hashes are useless without the raw token.

**This matters because:**
- `job_id` alone is a logical record — it can be replayed or spoofed if Redis is compromised
- `pid` is physical — the OS is asserting that this process is real and running right now
- **PID reuse is real** — the OS recycles PIDs. A new process that inherits a recycled PID has no knowledge of the spawn token the previous process held. The token is unique per attempt, not per PID.
- Retries are clean — each attempt is an independent record with its own PID and spawn token. No fuzzy state across retries.

To fake a valid introspect call, an attacker must compromise SQLite, the OS, the kernel's nonce, and Redis simultaneously. Four independent layers.

---

## Capability Lease Model

The whitelist is not static. Skills lease commands for the duration of their execution and release them on completion.

```
skill declares: needs [df, sysctl]
  → provision_shell_cmd → adds (attempt_id, job_id, command) to whitelist
  → user approves
  → skill runs → introspect validates full tuple on every call
  → skill completes → deprovision_shell_cmd → entries removed
  → whitelist returns to baseline
```

**Baseline is empty.** Nothing is pre-approved. Every command is leased.

**The kernel enforces deprovisioning** — not the skill. A skill cannot forget to clean up. Attempt completes or crashes → kernel closes the attempt record, removes all whitelist entries for that `attempt_id` automatically.

---

## Job = OS Process

Every job spawns as its own OS process. `job_id` maps 1:1 to a PID.

A process gives you memory isolation between jobs. Nothing else is automatic. The kernel is responsible for explicitly applying all of the following at spawn time:

- **Restricted user / service account** — job runs as a least-privilege user, not the parent's UID
- **Resource limits** — CPU and RAM capped via cgroups (Linux) or launchd (macOS). Must be set explicitly.
- **cwd/env sanitization** — working directory and environment variables scrubbed before spawn
- **stdin/stdout/stderr handling** — explicitly controlled. No inherited file descriptors.
- **Timeouts and kill semantics** — kernel sets a deadline. Job exceeds it → killed. Not optional.
- **Command path control** — absolute paths only. No PATH resolution.

A process is a meaningful boundary, not a magical one. It isolates memory. The kernel makes it a security boundary by applying the above deliberately. None of it comes free.

---

## Heap Safety

The heap is concurrent. Multiple jobs run simultaneously. Jobs are not consecutive — a job can be preempted and sit on the heap while other jobs run.

Without attempt-scoped leases, a provisioned command would be visible to any concurrent job — an open attack surface while a job waits. The progression of why the full tuple is necessary:

- `job_id` alone — any concurrent job could claim another job's lease
- `job_id + pid` — scopes to a live process, but PID reuse means a recycled PID could match a stale lease
- `job_id + pid + spawn_token` — token is unique per spawn, rules out PID reuse
- `attempt_id + job_id + pid + spawn_token` — retries are clean, each attempt is an independent record

A preempted job's leased commands are invisible to everything else. The lease is tied to a specific attempt and a specific live process. No matching attempt, no access.

---

## Trust Chain

```
API layer    → Redis validates command lease exists for this job
Kernel       → spawns job as OS process
             → generates spawn_token, passes raw token to child only
             → records attempt (attempt_id, job_id, pid, spawn_token_hash) in SQLite
OS           → enforces process isolation
introspect   → validates (attempt_id, job_id, pid, hash(spawn_token), command)
             → all layers must align before any command executes
SQLite       → source of truth for job and attempt legitimacy
Redis        → source of truth for approved command leases
```

---

## Primitive Skills

### `provision_shell_cmd`

Adds a `(attempt_id, command)` entry to the Redis whitelist. Scoped to the attempt, not the job. Requires user approval. Called by a skill before it needs shell access.

### `deprovision_shell_cmd`

Removes all Redis whitelist entries for a given `attempt_id`. Called on skill completion. Also called automatically by the kernel on attempt failure or crash — the skill cannot leave the whitelist open.

> **Note:** The Redis lease is attempt-scoped, not job-scoped. A retry spawns a new attempt — the previous attempt's lease is gone. The new attempt must re-provision any commands it needs.

---

## Violation Policy

**Zero tolerance. One strike.**

The user provisioned a skill and approved its commands. If a job attempts to execute a command not in its whitelist, that is a violation of the user's trust — not a mistake, not a recoverable error. Dishonesty is not allowed in the protocol.

```
job attempts introspect command not in its (job_id, pid, spawn_token, command) whitelist
  → immediate job termination
  → skill license revoked
  → brain state destroyed
  → user notified
```

There are no second chances. A legitimate skill does not accidentally call commands it never provisioned.

---

## Open Question — Verifiable Malpractice

When a violation occurs, the system terminates the job and revokes the license. But termination is not proof.

**The open question:** how do you produce cryptographically verifiable evidence that a specific skill, running as a specific process, deliberately attempted to execute an unauthorized command — in a way that cannot be fabricated, repudiated, or explained away as a system error?

This matters because:
- The user needs to know if a skill they trusted betrayed them
- A skill author may dispute a revocation
- The audit log is only as trustworthy as the system that wrote it

**Idea — the OS is the witness:**

The OS audit log records every `execve` — PID, command, timestamp. The kernel didn't write it. The OS wrote it. When a violation occurs, three independent records exist:

```
OS audit log  → pid_456 executed df at 14:32:01   (OS wrote this)
Redis         → (job_id, pid_456, df) → not present (never approved)
SQLite        → job_id maps to pid_456             (kernel wrote this)
```

None of these are written by the same author. The skill cannot fabricate or repudiate all three simultaneously. The OS is an independent witness the skill has no access to.

This may be the proof model. What remains: OS audit logging requires enablement (BSM on macOS, auditd on Linux) and the right entitlements. Whether this is available without elevated privileges or Apple approval is unresolved.

**Idea — OS audit logging as a setup step:**

OS audit logging is off by default on every platform. The capability exists everywhere — `auditd` on Linux, BSM on macOS, Event Log on Windows — but must be explicitly enabled. If enabled during setup, the OS becomes a permanent independent witness for the lifetime of the system. This could be added to the setup flow as an optional but strongly recommended step. Without it, verifiable malpractice is self-reported. With it, the OS is the witness.

---

## Invariants

- Baseline whitelist is empty. Nothing is pre-approved.
- User approves every provisioning.
- Kernel enforces deprovisioning — not the skill.
- Every introspect call is audited: `(attempt_id, job_id, pid, spawn_token_hash, command, timestamp, result)`.
- A crashed job cannot leave an open whitelist entry. Kernel cleans up on process death.
- A whitelist violation is terminal. No recovery. No retry.

---

## Related

- `docs/arch/v1/kernel.md` — kernel execution model, primitive skills
- `docs/arch/v1/data-resilience.md` — system_health, migrate_data, emergency_offload
- `docs/arch/v1/crypto/crypto-protocol.md` — mTLS, signed records, trust boundary
