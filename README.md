# OctOS - Personal Runtime

OctOS is a local-first personal runtime that executes intent across your own devices, learns from real usage over time, and keeps trust anchored to user-controlled commits.

---

## Core Model

- **Execution boundary**: the kernel is the canonical runtime boundary for all work.
- **Command protocol**: all work is expressed as `octos <tool> [args]`.
- **Shards**: devices register capabilities; routing is capability-based, not hardcoded.
- **Skills**: executable contracts discovered in memory and gated by Redis provisioning.
- **Memory**: property graph with two trust layers:
  - `observational`: writable working model, untrusted
  - `committed`: user-approved, append-only, trusted

Mental model:

> Nodes are identity. Edges are history. Truth is derived, not stored.

---

## Trust Model

- User keypair (Ed25519) is the root of trust.
- Committed writes require user-approved signed commits.
- Skills are content-addressed, immutable, and model-scoped for trust.
- Redis is the live execution gate; memory discovery alone does not grant execution.
- Non-brain shards are restricted to registered command primitives.

---

## Current v1 Status

- Kernel-first architecture is canonical.
- API Gateway + command protocol + shard capability model are documented.
- Graph memory model, skill model, and cryptographic trust model are defined.
- Voice event schema and ingress ACK semantics are defined.
- Open gaps remain around idempotency, job schema, plan approval return path, and Redis auth/write controls.

---

## Naming

- Canonical project name in this repo: **OctOS**
- Canonical command prefix: `octos`
- Legacy names in v1 docs: **Skyra** and **OctaOS**

When names conflict, treat `OctOS`/`octos` as canonical.

---

## Documentation Map

- Architecture index: `docs/arch/v1/README.md`
- Kernel (canonical execution model): `docs/arch/v1/kernel.md`
- API Gateway: `docs/arch/v1/api-gateway.md`
- Command protocol: `docs/arch/v1/command-parser.md`
- Shard communication: `docs/arch/v1/shard-communication.md`
- Shard registration: `docs/arch/v1/shard-registration.md`
- Capability model: `docs/arch/v1/capability-model.md`
- Memory model: `docs/arch/v1/memory-structure.md`
- Skill model: `docs/arch/v1/skill.md`
- Skill lifecycle: `docs/arch/v1/skill-lifecycle.md`
- Crypto protocol: `docs/arch/v1/crypto-protocol.md`
- Gaps register: `docs/arch/v1/gaps.md`
- Next steps: `docs/arch/v1/next-steps.md`
- Voice schema: `octos/voice/README.md`
- White paper (draft): `docs/whitepaper.md`
