# Setup

Setup is the only time the system requires terminal access. Once the system is running, terminal access is closed. The system operates headlessly from that point.

---

## Prerequisites

Four components. All can run on the brain shard for v1.

| Component | Purpose |
|---|---|
| Kernel | Central execution boundary. Dispatch, trust enforcement, job tree. |
| Graph DB | The memory. Committed + observational layers. Nodes, edges, vectors. |
| Redis | Trust membrane. Skill registry, shard capabilities, credentials. |
| Model | Local inference. Interprets and executes skills. |

The brain does not need to run the model permanently. For v1 it does. When a GPU shard comes online and registers `deep_reasoning`, the brain routes there. The brain never changes — the routing does.

---

## Setup Flow

```
1. terminal in
2. install dependencies (kernel, graph DB, Redis, model runtime)
3. boot Redis
4. seed Redis with system primitive skills (pre-provisioned at boot)
5. initialize graph DB (empty committed + observational layers)
6. generate user keypair (Ed25519) — private key stays on device
7. register user public key with kernel
8. load model
9. boot kernel
10. boot API gateway (minimal, on brain for v1)
11. boot cron service (on brain for v1)
12. verify: /health checks pass on all services
13. close terminal
```

After step 13 the system is sealed. Terminal access is a security surface — closing it reduces the attack surface to what's exposed through the kernel and Redis, both behind mTLS and signed records.

---

## What Runs Headlessly

After setup, the system runs as managed services with no terminal required:

- Kernel daemon
- Redis
- Graph DB
- Model runtime (Ollama or equivalent)
- API Gateway
- Cron Service

Normal operation — conversations, skill execution, background reasoning, graph mutation — all happen without terminal access.

---

## Shard Registration

When a new shard comes online (GPU shard, Raspberry Pi, etc.), it does not require terminal access on the brain. The shard boots, generates its own keypair, and registers capabilities through the kernel's standard shard registration flow. The brain issues a registration token over mTLS. See `docs/arch/v1/shard-registration.md`.

Terminal access is only required on the new shard itself during its own initial setup.

---

## Security Properties After Setup

- No SSH access to brain
- No open terminal
- Attack surface: kernel ingress (mTLS) + Redis (mTLS + signed records)
- User data (committed graph) — protected by user keypair, append-only, never exposed through the gateway
- Credentials (API keys) — live in Redis, protected by mTLS + Redis auth

---

## Related

- `docs/arch/v1/shard-registration.md` — adding new shards after initial setup
- `docs/arch/v1/crypto-protocol.md` — keypair generation, mTLS, signed records
- `docs/arch/v1/kernel.md` — kernel boot, primitive skill provisioning
- `docs/arch/v1/api-gateway.md` — gateway decomposability, plugin model
