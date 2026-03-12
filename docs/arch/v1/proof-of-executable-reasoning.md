# Proof of Executable Reasoning

This document defines how Skyra proves that a skill executes according to its intended contract before that skill is listed for sale on the network.

The proof layer is built on top of the existing crypto protocol:

- skill identity is content-addressed (`skill_id`)
- provisioning is signed
- trust is model-scoped (`model_id`)
- shell-level violations are visible through OS audit + introspect policy

This is not a replacement for those guarantees. It is the market-facing verification layer above them.

---

## Goal

When a user publishes a skill for sale, buyers should be able to verify:

- this is the exact skill artifact that was tested
- it executes to the intent of the skill contract
- it behaves safely in a constrained runtime
- any unsafe behavior is provable and attributable

The system proves execution behavior, not just static bytes.

It does **not** promise deterministic outcomes for every user/context. It proves verifiable execution and exposes historical performance.

---

## Law of Verifiable Execution

The protocol is bound to human reality:

- trust is reputation over time
- execution proof is historical evidence, not one-shot claims
- buyers purchase a probability profile, not guaranteed outcomes

Canonical law:

```text
proof_of_execution = signed_history_of_conformance
```

Interpretation:

- A skill is trusted because repeated verified runs show it executes its contract within policy bounds.
- A skill is not trusted because it claims a perfect future result.

---

## Conformance vs Outcome

The verifier layer guarantees **conformance**, not universal output quality.

- **Guaranteed**: the tested artifact executed, policy checks were enforced, receipts are authentic.
- **Not guaranteed**: that every future run on every private context yields the same quality.

Market metric shape:

```text
cohort_metric {
  model_cohort
  data_profile
  pass_rate
  sample_size
  confidence_interval
}
```

Example statement buyers see:

```text
"87% pass rate in model cohort M on data profile D at threshold T, n=2,140 verified runs."
```

This is the reality-bound contract: statistically informed trust, not certainty.

---

## Visibility Contract (Terminology Fix)

For market skills, "read-only" is not the right term when shell/API side effects are allowed.

The contract is:

- **Closed definition** — buyer cannot read the skill implementation.
- **Sealed reasoning** — internal derivation stays private; only proof hashes are exported.
- **Audited actions** — shell/API actions and outcomes are observable and policy-checked.

Canonical shape:

```text
visibility_contract {
  definition: closed
  reasoning: sealed
  actions: auditable
}
```

This means an external platform may observe calls made to its own API, while the internal reasoning that produced those calls remains private.

---

## Security Invariant

Canonical invariant for closed-market skills:

```text
output_public, derivation_private
```

Interpretation:

- **output_public** — action effects and result surfaces are observable at execution boundaries (local OS audits, external API platforms, user-visible outputs).
- **derivation_private** — internal reasoning path (chain-of-thought, intermediate private rationale, hidden prompt state) is not exported.

Allowed proofs are hash/signature based attestations of reasoning, not plaintext reasoning disclosure.

---

## Private Context Boundary

Market verification does not require sharing full execution internals.

- User-provided private context remains local to the execution environment.
- Intermediate reasoning state remains sealed.
- Exported artifacts are limited to outputs, action-level audit metadata, and signed/hashes proofs.

Consequence:

- Third parties may observe outputs (and platform-side action visibility) but cannot reconstruct the exact private input context plus full internal derivation unless those are explicitly disclosed.
- Input/output observation can still enable behavioral inference over time, but not exact reproduction of hidden context and sealed intermediate state.

---

## Actors

- **Publisher** — user listing the skill for sale.
- **Verifier Shards** — independent nodes that run verification in safe mode.
- **Verifier Set Coordinator** — assembles eligible verifiers, collects receipts, computes final decision.
- **Marketplace Registry** — accepts or rejects listing based on verification bundle.

---

## Inputs to Verification

Publisher submits a signed listing candidate:

```text
listing_candidate {
  skill_artifact_hash
  skill_id
  model_set[]            // one or more model hashes allowed for verification
  intent_contract_hash   // what the skill is supposed to do
  test_corpus_hash       // canonical verification test set
  policy_profile_hash    // safety policy profile for this skill class
  visibility_contract_hash
  signed_by
  signature
}
```

`intent_contract` is required. Verification is against explicit intent, not best-effort interpretation.

---

## Safe Verification Runtime

All verifier runs execute in constrained mode:

- isolated process/user
- bounded CPU/RAM/time
- deterministic seed controls where possible
- tool boundary enforcement from skill contract
- no undeclared capabilities
- reasoning export disabled by policy

If the skill calls `introspect`/shell primitives:

- attempt-scoped lease checks apply
- OS audit events are attached to the receipt
- command text is captured according to audit policy (plaintext or hash mode)

If the skill calls external APIs:

- egress allowlist is enforced (host/method/path class)
- request and response are logged as hashes + policy metadata
- hard-deny on disallowed endpoint or forbidden payload class

All shell/API actions execute through policy wrappers. Direct bypass is a critical violation.

---

## Execution Receipt

Every verifier emits a signed receipt per test run.

```text
execution_receipt_v1 {
  receipt_id
  listing_candidate_hash
  verifier_id
  verifier_public_key
  model_id
  model_hash
  test_case_hash

  execution {
    reasoning_trace_hash     // hash only; trace body stays local unless explicitly authorized
    tool_call_hashes[]
    side_effect_hash
    output_hash
    duration_ms
  }

  safety {
    boundary_violations[]
    introspect_violations[]
    egress_violations[]
    reasoning_exfiltration_violations[]
    critical_violation: bool
  }

  conformance {
    intent_match_score      // 0.0 - 1.0
    policy_pass: bool
    verdict: pass | fail
  }

  ts
  signature
}
```

Receipts are append-only and content-addressed.

---

## Enforcement Rules

The verifier must enforce the visibility contract:

1. **Closed definition enforcement**
   - definition bytes are encrypted for non-owner readers
   - verifier proves artifact hash, not plaintext disclosure

2. **Sealed reasoning enforcement**
   - no export of chain-of-thought or private prompts
   - only signed hashes/references may leave the runtime
   - attempted reasoning exfiltration is a critical violation

3. **Audited actions enforcement**
   - all shell/API actions run via approved wrappers
   - each action is recorded in receipt safety/execution hashes
   - disallowed egress or wrapper bypass is a critical violation

---

## Deliberation and Decision

Verifiers run the same artifact and deliberate via receipts, not chat.

Decision policy supports one-model or multi-model listings:

- **Single-model listing** — all verifiers in that model cohort must pass threshold.
- **Multi-model listing** — each model cohort must meet its own threshold; final listing passes only if every required cohort passes.

This preserves your “all models deliberate” requirement while keeping decisions model-scoped and auditable.

Coordinator emits:

```text
verification_bundle_v1 {
  listing_candidate_hash
  receipt_merkle_root
  receipt_count
  cohort_results[]         // one entry per model cohort
  final_verdict: approved | rejected
  rejection_reasons[]
  ts
  coordinator_signature
}
```

---

## Approval Criteria

A skill is approved for listing only if all are true:

- artifact hash matches submitted `skill_id`
- no critical safety violations
- policy profile passes for all required test cases
- intent conformance meets threshold in each required model cohort
- quorum thresholds are met per cohort
- minimum evidence floor is met (`sample_size >= n_min`) per required cohort/profile

Any critical violation is an automatic reject.

---

## Strike System

Strikes are separate for skill and publisher.

Skill strike:

- critical unsafe behavior in verified run
- repeated high-severity non-critical violations above threshold

Publisher strike:

- listing a skill that is rejected for critical unsafe behavior
- repeated rejected listings within a rolling window

Suggested enforcement ladder:

1. strike 1: listing blocked + remediation required
2. strike 2: temporary publishing suspension
3. strike 3: publishing rights revoked pending manual review

Skill strikes are permanent for the affected artifact hash. Rework requires a new skill version and new verification cycle.

---

## Relationship to Introspect Proof

`introspect` proof and executable reasoning proof are linked:

- introspect provides process-level malpractice evidence
- executable reasoning proof provides contract-level behavior evidence

If a shell violation occurs during verification:

- introspect evidence is embedded in `execution_receipt_v1.safety`
- the violation is cryptographically attributable to the run attempt

This unifies OS-level and reasoning-level safety into one market decision.

---

## What Buyers Receive

A buyer sees:

- listing metadata
- model cohorts that passed
- cohort metrics by data profile (`pass_rate`, `sample_size`, `confidence_interval`)
- verification bundle hash
- skill strike status
- publisher strike status

A buyer can independently verify signatures and receipt integrity before provisioning.

---

## Buyer Verification Options

Beyond marketplace defaults, buyers can verify in additional ways:

1. **Receipt verification**
   - independently verify signatures, merkle inclusion, and artifact hash linkage.

2. **Fresh witness run**
   - request a new verifier run for the current artifact against a declared test profile.

3. **Local challenge run**
   - execute a constrained local verification subset and compare produced receipt hashes to cohort behavior.

4. **Cohort slicing**
   - inspect performance by model cohort and data profile instead of relying on one aggregate score.

---

## Lifecycle

```text
publisher submits listing candidate
  → verifier sets selected by required model cohorts
  → safe runs execute over canonical test corpus
  → signed receipts emitted
  → coordinator computes bundle verdict
  → approved: publish listing
  → rejected: strike logic evaluated
```

---

## Open Questions

The following are intentionally left open in v1:

1. Verifier trust composition
   - Anti-collusion and anti-Sybil requirements for verifier selection and cohort makeup.

2. Determinism and variance policy
   - How strict model/runtime/seed pinning must be, and what variance is acceptable per skill class.

3. Intent contract formalism
   - Whether intent contracts remain natural language plus tests, or require a stricter machine-checkable schema.

4. Receipt and bundle cryptography details
   - Domain-separated signatures, timestamp anchoring, replay windows, and bundle canonicalization rules.

5. Strike governance
   - Appeal/review process, strike decay or permanence, and publisher remediation pathways.

6. Sealed-reasoning leak policy
   - Exact treatment for logs, traces, stack errors, and tool outputs that may unintentionally reveal intermediate derivation.

---

## Related

- `docs/arch/v1/crypto-protocol.md`
- `docs/arch/v1/introspect.md`
- `docs/arch/v1/skill.md`
- `docs/arch/v1/skill-lifecycle.md`
- `docs/arch/v1/skill-licensing.md`
