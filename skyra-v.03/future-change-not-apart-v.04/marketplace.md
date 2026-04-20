# Marketplace

A market for verified beings and worlds. Trust is cryptographic and statistical. Privacy is preserved. Nobody is the gatekeeper — the proof is the gatekeeper.

## The Unit of Sale

Two publishable artifacts:

**Single being** — genome + adapter(s). A being that does one thing well. Verified against its intent contract.

**World being** — genome + adapter(s) + seed beings + their genomes. An entire cognitive environment. A research world. A writing world. A coding world. Verified as a unit. The buyer gets a populated world they can grow into their Skyra with one directive.

## Content Addressing

Every artifact is content-addressed. The artifact_id is the SHA-256 hash of the genome + adapters + any seed beings, hashed together as a unit.

```
artifact_id = SHA-256(genome_bytes || adapter_bytes || seed_beings_bytes)
```

There is no edit operation. A change produces a new artifact_id. The old artifact continues to exist at its original hash. The marketplace stores artifacts by content hash — a write to an existing hash is rejected.

## Publishing

The publisher submits a signed listing candidate:

```
listing_candidate {
  artifact_id           — content hash of the being or world
  intent_contract       — what this being is supposed to do, explicit
  test_corpus_hash      — canonical verification test set
  model_set[]           — which models this is verified against
  visibility_contract   — open or closed definition, sealed reasoning, auditable actions
  signed_by             — publisher public key
  signature             — Ed25519 over the above
}
```

`intent_contract` is required. Verification runs against explicit intent, not interpretation.

## Verification

Verifier beings run the artifact in isolated world beings — constrained execution environments with no access to the buyer's data, no undeclared capabilities, bounded resources. Every run produces a signed execution receipt.

```
execution_receipt {
  artifact_id
  verifier_id
  model_id               — trust is model-scoped
  test_case_hash
  execution {
    reasoning_trace_hash — hash only; trace stays local
    tool_call_hashes[]
    output_hash
    duration_ms
  }
  safety {
    boundary_violations[]
    critical_violation: bool
  }
  conformance {
    intent_match_score   — 0.0 to 1.0
    policy_pass: bool
    verdict: pass | fail
  }
  signature
}
```

Receipts are append-only and content-addressed. The verifier cannot alter a receipt after signing it.

## Statistical Trust

Trust is proven over time, not claimed once.

```
cohort_metric {
  model_id
  pass_rate
  sample_size
  confidence_interval
}
```

What buyers see:

```
87% pass rate on claude-sonnet-4-6 · n=2,140 verified runs · 95% confidence interval [84%, 90%]
```

A being is trusted because repeated verified runs show it executes its contract. Not because it claims a perfect future result. Buyers purchase a probability profile, not a guarantee.

Trust is model-scoped. A being verified under claude-sonnet-4-6 is not trusted under a different model. The user approved what that specific model produced. A model change flags the artifact — still provisioned, not executable until re-verified under the new model.

## Privacy

Three axes:

**Closed definition** — the buyer cannot read the genome or adapter implementation. The artifact is an encrypted blob. The buyer gets the artifact_id and the proof, not the internals. Creator IP is protected.

**Sealed reasoning** — internal chain-of-thought stays local to the execution environment. Only signed hashes leave the runtime. Attempted reasoning exfiltration is a critical violation.

**Auditable actions** — shell and API actions are observable at execution boundaries. The buyer can see what the being does in the world. They cannot see how it decided to do it.

```
visibility_contract {
  definition: open | closed
  reasoning: sealed
  actions: auditable
}
```

## Provisioning

Provisioning flows through grow. No new syntax.

```
skyra world ~name research ~artifact <artifact_id> | reason
```

The kernel fetches the artifact, verifies the hash matches the artifact_id, checks the buyer's access profile, and grows the being or world into the runtime. A world being provisions as a populated world — all seed beings grow inside it.

## Access Profiles

An access profile is how the publisher grants a buyer access to an artifact. Signed by the publisher. Either party can revoke.

```
access_profile {
  profile_id
  granted_to        — buyer's public key
  scope             — which artifact, what operations
  signed_by         — publisher public key
  signature
  revocable_by      — [publisher, buyer]
}
```

Two-way sovereignty. The publisher controls who provisions their artifact. The buyer controls whether they participate. Either side pulls it, the access terminates. The history that the profile existed is preserved. The live access is gone.

## Strike System

Strikes are separate for artifact and publisher.

**Artifact strike** — critical safety violation in a verified run, or repeated high-severity violations above threshold. An artifact with a critical strike is permanently rejected at that hash. Rework requires a new artifact with a new artifact_id and a new verification cycle.

**Publisher strike** — listing an artifact that receives a critical violation. Repeated rejections within a rolling window escalate the strike level.

Enforcement ladder:
1. Listing blocked, remediation required
2. Temporary publishing suspension
3. Publishing rights revoked pending review

## What This Enables

Skyra can grow a verified research world with one directive. The world arrives pre-populated with beings that have a proven performance record under a specific model. The buyer knows what they're getting before they provision it. The creator ships a closed definition — the buyer gets the capability, not the IP.

Creators build beings and worlds. They publish them. Skyra instances provision them. The marketplace takes a cut. Trust is the product, not the platform.

## Relationship To Crypto Protocol

This layer sits above the crypto protocol defined in the v1 archive. The crypto protocol provides the foundation — Ed25519 keypairs, content-addressed artifacts, signed commits, model-scoped trust. Proof of executable reasoning is the market-facing verification layer above it. Neither replaces the other.

## Open Questions

- Who runs the verifier beings — a marketplace-operated fleet, or any Skyra instance that volunteers?
- How does the intent contract stay honest — natural language plus tests, or a stricter machine-checkable schema?
- Anti-collusion for verifier selection — how do you prevent a publisher from influencing which verifiers run their artifact?
- Does the marketplace operate as a world being itself, or as a separate service the sense being knows how to call?
- Strike decay — do strikes expire, or are they permanent for the artifact hash?
