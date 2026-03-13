# Shard Communication

## The Protocol

`octos <tool> [args]` is the protocol. No JSON. No custom message formats. No envelopes. One unified command language across every shard in the network.

Every shard speaks this syntax. Voice Shard, Cron Service, GPU Shard, Brain Shard — all emit the same commands. There is no special wire format per shard type.

---

## Brain vs Non-Brain Shards

**Brain Shard** — full command authority. Elected role. Owns the kernel, the API Gateway, Redis, and the shard registry. Can issue any command.

**Non-brain shards** — limited vocabulary. Can only communicate back to the brain using their registered primitives. Cannot send arbitrary data. Cannot go off-script. Everything outside their registered primitives is rejected.

This extends the trust model to communication itself. Just as skills must be provisioned in Redis to execute — the commands a shard can emit are provisioned at registration time. The brain defines the vocabulary when it sends the service package to the shard on boot.

---

## Primitives Are Registered at Boot

When a shard registers with the brain:

1. Shard fingerprints hardware, reports capabilities
2. Brain evaluates capabilities, sends service package
3. Service package includes: services to run + **registered command primitives**
4. Shard can only emit those primitives — nothing else

Adding a new primitive requires a new service package from the brain. Shards cannot self-authorize new commands.

---

## Responsibility Split — Brain vs Ingress Shards

**Brain Shard** owns session and turn identity. It generates `session_id` and `turn_id`. Ingress shards never produce or track these — they don't know about sessions. That is the brain's concern.

**Ingress shards** are responsible for one thing: issuing commands. They stream data to the brain via the command protocol. The brain handles everything else.

## Session Data Is Shardless

Sessions live in the brain. Not on any shard. No shard owns a session.

This means:

- **Walk between rooms** — different ingress shard picks up your voice, same session continues. The brain doesn't care which shard is streaming.
- **Multiple shards, one session** — your phone and the Pi can both stream commands into the same session simultaneously.
- **Shard failure doesn't kill a session** — the brain holds all session state. The shard was just a command emitter. Another shard can pick up immediately.
- **Shards are stateless** — they stream and they forget. State lives in the brain.

Sessions are a brain primitive. Shards are command emitters. The session has no physical location.

## Example — Voice Shard Primitives

The Voice Shard issues commands. It does not track session or turn IDs:

```
octos stream --token="what" --valence=-0.4 --arousal=0.7 --dominance=0.5
octos stream --token="did" --valence=-0.4 --arousal=0.7 --dominance=0.5
octos stream --token="nginx" --valence=-0.5 --arousal=0.8 --dominance=0.4
```

That is its entire vocabulary. Raw data. Raw vectors. No session. No turn. No context. The brain receives the stream, assigns identity, tracks the session.

Everything else is rejected at Ingress.

---

## Example — Cron Service Primitives

The Cron Service emits whatever skills it is scheduled to execute:

```
octos pattern_scan
octos memory_snapshot
```

Registered at boot. Fired on schedule. Nothing else.

---

## ACK Is a Command Too

ACK flows back as a command, not a JSON response:

```
octos ack --turn=turn_8f4c --status=stored
```

One language. Both directions.

---

## What This Means for the API Gateway

The API Gateway's Ingress receives commands — not JSON blobs. It does not parse message envelopes. It parses `octos <tool> [args]`.

The provisions, security metadata, and job envelope are assembled BY the gateway after parsing the command. The shard is responsible for forming a valid command. The gateway is responsible for what happens next.

---

## Why This Works

- **LLMs can't reliably emit JSON contracts** — JSON is rigid, schema-dependent, and error-prone for LLMs to produce correctly at scale. A CLI command is natural. LLMs already think in terms of tool calls. `octos <tool> [args]` maps directly to how an LLM reasons about what to do next. The protocol is designed around what LLMs can reliably emit.
- **Dynamic by nature** — JSON schemas break when you add fields. CLI commands are additive. New args, new tools — the protocol absorbs change without a schema migration.
- **The LLM's output IS the protocol** — the LLM inside Skyra issues `octos <tool> [args]`. The same syntax the shards use. No translation layer. No serialization step. The reasoning and the wire format are the same thing.
- **Security** — non-brain shards have a fixed, auditable vocabulary. No shard can issue commands it wasn't provisioned for.
- **Simplicity** — one protocol everywhere. No per-shard message schemas to maintain.
- **Trust** — the command language IS the trust boundary. If it's not a registered primitive, it doesn't execute.
- **Debuggability** — every inter-shard communication is a readable command. No opaque JSON blobs.

---

## Related

- `docs/arch/v1/api-gateway/api-gateway.md` — Ingress receives commands, assembles job envelope
- `docs/arch/v1/kernel.md` — kernel executes commands
- `docs/arch/v1/shard/shard-registration.md` — how primitives are registered at boot
- `docs/arch/v1/gaps.md` — G31: syntax propagation, G32: API Gateway domain resolution
