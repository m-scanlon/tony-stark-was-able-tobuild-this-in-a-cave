# Crypto Protocol

Skyra's trust model is cryptographic, not positional. Nothing is trusted because of where it came from. Everything is verified.

---

## Keys

The user holds one keypair. Ed25519.

```
user_private_key  — signs commits, approvals, skill provisioning, access profiles
user_public_key   — stored on-device, used by kernel for verification
```

The private key never leaves the device. The kernel holds the public key. No cloud. No escrow. No recovery. This is consistent with Principle 11: your keys, your data, your consequences.

---

## Access Profiles — The Foundation of Multi-Party Trust

An access profile is how one party grants another party access to their system. It lives in Redis. It is signed by the user. Either party can revoke it.

```
access_profile {
  profile_id:     string
  granted_to:     bytes            // grantee's public key — their identity
  scope:          []permission     // exactly what they can access — no more
  provisioned_at: timestamp
  signed_by:      user_public_key  // user grants it
  signature:      bytes            // Ed25519 over (profile_id + granted_to + scope + provisioned_at)
  revocable_by:   [user, grantee]  // either party can terminate
}
```

**The user provisions it.** They define the scope — which skills, which data, which operations the grantee can access. The signature locks it. The kernel enforces it.

**Either party can deprovision it.** The user can revoke access at any time. The grantee can remove themselves at any time. The moment either side pulls it, the profile is gone from Redis. The relationship terminates.

**Deprovisioning is not deletion.** The committed layer records that the profile existed. The history is preserved. Only the live access is terminated.

**This is two-way sovereignty.** The user controls who enters their system. The grantee controls whether they participate. Neither is locked in by the other.

```
use cases:
  skill creator ships a skill → user provisions access profile → creator can revoke support
  developer given read access → developer can remove themselves
  third-party service → user can terminate at any time, no permission needed from the service
```

Authentication proves identity — mTLS certificate, public key. Authorization is resolved from the access profile in Redis. These are separate steps. The auth layer never sees the user profile. It proves the certificate is valid and passes the verified identity to the kernel. The kernel resolves permissions from the access profile. The user profile is never in the auth layer's hands.

Compromise of auth means a stolen certificate. Real but bounded. It does not mean the attacker has the user's committed graph, skill permissions, or access profiles. Those are gated by the kernel after auth.

---

## Skills and Memory Are Provisioned as a Unit

A skill and its memory namespace are inseparable. When a skill is provisioned, its memory namespace is provisioned alongside it — atomically. You cannot provision one without the other.

A creator who distributes a skill ships two things as a single artifact:

```
skill_artifact {
  skill_definition:  bytes   // the roadmap, contract, boundary rules
  seed_memory:       bytes   // pre-populated memory namespace — the creator's knowledge, patterns, context
}
```

Both are hashed together. Both are subject to the same `definition_visibility` rules. A closed skill means the consumer cannot read the definition or the seed memory. An open skill means both are visible.

**Two kinds of memory in a skill's namespace:**

- **Seed memory** — shipped by the creator. Pre-populated before the skill ever runs. Protected by `definition_visibility`. This is the creator's IP.
- **Runtime memory** — written by the skill during execution on the consumer's device. Belongs to the consumer. Their data, on their hardware, subject to their sovereignty.

The kernel maintains this distinction. Seed memory is read-only to the consumer (protected by creator signature). Runtime memory is the consumer's and follows standard write ownership rules.

---

## Skills Are Immutable

A skill is content-addressed. When a skill crystallizes, its definition is hashed (SHA-256). That hash is the skill's identity.

```
skill_id = SHA-256(skill_definition_bytes || seed_memory_bytes)
```

The skill artifact — definition + seed memory — is hashed together as a unit. The definition is the roadmap, contract, validation criteria, and boundary rules. The seed memory is the pre-populated namespace the creator ships with it. Both are locked together at the moment of approval. Either changing produces a different ID.

**There is no edit operation for skills.** Modification means a new skill with a new ID. The old skill continues to exist at its original hash. `update_skill` is the only path — it produces a new node with a new ID, not a mutation of the old one.

This is enforced structurally: the kernel stores skills by content hash. A write to an existing hash is rejected.

---

## Skill Provisioning

When a skill is approved and provisioned into Redis, it carries a signed provisioning record.

```
provisioning_record {
  skill_id:               string         // SHA-256 hash of the skill definition
  model_id:               string         // model under which the skill was committed — trust is model-scoped
  provisioned_at:         timestamp
  access:                 read | write   // consumer-set at provisioning time — memory access
  definition_visibility:  open | closed  // creator-set at distribution time — definition readability
  creator_public_key:     bytes          // creator's public key (may differ from consumer)
  signed_by:              user_public_key
  signature:              bytes          // Ed25519 signature over (skill_id + model_id + access + definition_visibility + provisioned_at)
}
```

## Two Axes of Trust

**Owner trust** — cryptographic. Proven at commit time by the owner's signature. Binary. A skill is trusted by its owner because they signed it. This does not change.

**External trust** — historical. Proven by the skill's execution record. When a consumer evaluates a skill before provisioning, the signature proves authenticity — the skill is what the creator claims, unmodified. The history proves quality — the skill produces outputs that users actually commit.

```
skill execution record {
  executions:     int        // total times run across all users
  commit_rate:    float      // percentage of outputs committed vs rejected
  user_count:     int        // distinct users who have run it
}
```

A skill with 10,000 executions and a 95% commit rate tells a consumer something a signature alone cannot. The signature is necessary. The history is the reputation. Both together are trust.

This is the same model as open source software. The GPG signature proves the release is authentic. The years of commits and community use prove it is worth trusting.

---

**Trust is model-scoped.** A skill committed under one model is not trusted under a different model. The user approved the skill in the context of what that specific model produced. A different model reasons differently — the user has not seen its output. When the active model changes, skills carrying a different `model_id` are flagged: still visible in memory, not executable in Redis. The user must re-approve each skill under the new model before it is trusted again.

**Two independent permission axes:**

### Memory Access — set by the consumer at provisioning time

- `read` — the skill can read from memory (observational + committed). It cannot write nodes or edges. It cannot issue `write_node` or `write_edge`. Kernel enforces at dispatch time.
- `write` — the skill can write to the observational layer. Standard permissions. Cannot write to the committed layer without a separate `propose_commit` flow.

The consumer sets memory access when they provision the skill. A creator who provisions their own skill sets their own access. A user who purchases or receives a skill sets theirs independently.

### Definition Visibility — set by the creator at distribution time

- `open` — the skill definition is readable. The consumer can inspect the roadmap, contract, boundary rules, and validation criteria. Source-visible.
- `closed` — the skill definition is an encrypted blob. The consumer cannot read what is inside. The kernel can execute it. The user cannot inspect it.

A closed skill is encrypted by the creator with a symmetric key. The creator distributes the encrypted blob + the `skill_id` (the hash of the plaintext definition). The consumer provisions and executes it without ever seeing the definition. This is skill IP protection — the creator ships a black box.

**Definition visibility cannot be changed by the consumer.** It is set by the creator and is part of the creator's signature on the skill artifact. A consumer attempting to change `closed` to `open` in the provisioning record fails signature verification against `creator_public_key`.

### Provisioning Is Append-Only

Access level cannot be changed after provisioning. To change access, the skill must be re-provisioned with a new user signature. A new provisioning record is written. The old record is archived — same append-only principle as the committed graph layer.

---

## Commits

Every write to the committed layer requires a user-signed commit.

```
commit {
  node_ids:      []string        // nodes being committed
  edge_ids:      []string        // edges being committed
  committed_at:  timestamp
  signed_by:     user_public_key
  signature:     bytes           // Ed25519 signature over (node_ids + edge_ids + committed_at)
}
```

The kernel verifies the signature before writing anything to the committed layer. An unsigned commit is rejected. A commit signed by anyone other than the registered user key is rejected.

Commits are append-only. A commit is never mutated after landing.

---

## Shard Authentication

Shards communicate over mTLS. Each shard holds its own keypair, generated at first boot. The brain shard (control plane) holds a CA cert. Shards present their cert to the brain at registration. The brain validates and issues a registration token.

```
shard registration:
  shard generates Ed25519 keypair at boot
  shard presents cert to brain
  brain validates, issues registration token
  registration token is scoped to shard_id + capabilities
  token is signed by brain's private key
```

A shard without a valid registration token cannot receive dispatched commands. A shard that presents a mismatched cert is rejected.

This means: a rogue process on the network cannot impersonate a shard. A compromised shard cannot escalate beyond its registered capabilities.

---

## Redis Trust Boundary

Redis is the live skill registry and trust membrane. Skills in memory are inert. Skills provisioned in Redis are executable.

Redis entries for skills are signed. The kernel verifies the provisioning signature before treating a Redis entry as authoritative. A Redis entry without a valid user signature is treated as untrusted and rejected.

This means: even if Redis is compromised, an attacker cannot inject skills without the user's private key.

---

## Verification Summary

| Operation | Verified by | Key used |
|---|---|---|
| Skill execution | Kernel checks Redis provisioning record signature | User public key |
| Committed layer write | Kernel checks commit signature | User public key |
| Shard dispatch | mTLS cert + registration token | Brain CA + shard keypair |
| Skill content integrity | Kernel re-hashes skill definition, compares to skill_id | SHA-256 |
| Memory access enforcement | Kernel reads signed access field in provisioning record | User public key |
| Definition visibility enforcement | Kernel checks definition_visibility in creator-signed artifact | Creator public key |
| Closed skill execution | Kernel decrypts blob, verifies hash matches skill_id, executes | Creator symmetric key |

---

## What This Prevents

- **Skill tampering** — content-addressing makes silent modification impossible
- **Privilege escalation** — memory access levels are signed and immutable post-provisioning
- **Definition exposure** — closed skill definitions are encrypted; the consumer never sees plaintext
- **Creator IP theft** — definition visibility is creator-signed; consumer cannot override it
- **Rogue shards** — mTLS + registration tokens gate shard participation
- **Redis poisoning** — signed provisioning records mean unsigned entries are rejected
- **Committed layer corruption** — every committed write requires user signature
- **Replay attacks** — timestamps are part of signed payloads; the kernel rejects stale signatures outside a configurable window

---

## What This Does Not Prevent

- **Biased model output** — the model that produces skill output is a dependency, not a component. The crypto layer protects the committed layer from unauthorized writes. It does not audit the reasoning that produced what gets proposed. See Principle 12.
- **User self-sabotage** — if the user signs a bad commit, it lands. Sovereignty means owning the downside. The system does not protect the user from themselves.
- **Compromised private key** — if the private key is stolen, the attacker becomes the user. No escrow, no recovery. Hardware key storage (TPM, YubiKey) is the mitigation.

---

## Related

- `docs/arch/v1/kernel.md` — Redis trust boundary, dispatch flow
- `docs/arch/v1/skill-lifecycle.md` — skill crystallization, provisioning flow
- `docs/arch/v1/memory-structure.md` — committed layer, append-only model
- `docs/arch/v1/principles.md` — Principle 2 (data integrity), Principle 11 (sovereignty), Principle 12 (model as dependency)
