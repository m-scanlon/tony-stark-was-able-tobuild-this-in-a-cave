# Skill Improvement

Skyra's core algorithms — retrieval, reasoning, integrate, crystallization — are committed skills. They are creator-signed, trusted, and immutable. Skyra cannot change them.

But Skyra can observe how they perform and reason about how they could be better. The creator gives her a bounded space to do exactly that. Inside that space, she reasons freely. Outside it, nothing.

---

## The Contract

Every skill that carries an `improvement_scope` field opens a bounded improvement namespace in the observational layer — a private scratchpad dedicated to algorithmic reasoning.

```
skill_contract {
  ...
  improvement_scope: {
    aspects:     []string    // what Skyra is allowed to reason about improving
    namespace:   string      // observational namespace for improvement reasoning
  }
}
```

**`aspects`** — the creator-defined boundary. What Skyra may reason about. Examples:

```
retrieval skill:
  aspects: ["ranking weights", "decay formula", "session count threshold"]

reasoning skill:
  aspects: ["entity extraction confidence threshold", "edge weight derivation"]

crystallization skill:
  aspects: ["recurrence ratio threshold", "affect weight modifier"]
```

Skyra cannot reason outside the declared aspects. The kernel enforces this the same way it enforces boundary rules — before each write to the improvement namespace.

**`namespace`** — the observational scratchpad. Skyra writes here freely. No user gate. Not trusted. Her working model for improvement reasoning.

---

## The Loop

```
Skyra executes the committed algorithm
  → observes outcomes (quality, failures, edge cases)
  → reasons in bounded improvement namespace
  → improvement coheres into a proposal
  → propose_commit → creator approves
  → provision_skill → new skill created (new content hash, new model_id)
  → old version persists (append-only)
```

Skyra does the observational work. The creator approves or rejects. The algorithm improves without the creator having to think about it constantly.

---

## Trust

The improvement namespace is observational — not trusted. Skyra's reasoning about improvement is her working model. It does not change the committed algorithm.

The committed algorithm only changes through `provision_skill` of a new version (legacy term: `update_skill`) with user approval. Same trust rules as any other committed node. Same crypto guarantees. Same append-only history.

**Trust is model-scoped here too.** An improvement proposal generated under 7B is not trusted under 32B. If the model changes, pending improvement proposals are flagged. The creator re-evaluates them under the new model before approving.

---

## Why the Creator Assigns the Scope

The creator knows what the algorithm is trying to do. They know which parameters are safe for Skyra to reason about and which would break the system if changed naively. The scope is not arbitrary — it is the creator's judgment about where Skyra's reasoning adds value without introducing risk.

This is Principle 3 applied to the algorithm itself: constrain the data (the improvement scope), not the model (Skyra's reasoning inside it). Inside the declared aspects, Skyra reasons with full intelligence. The scope is the only guardrail.

---

## What This Means Over Time

The algorithms get better without the creator redesigning them. Skyra runs the retrieval algorithm thousands of times. She observes what works, what doesn't, what edge cases appear. Her improvement namespace accumulates data. Her proposals become more informed.

But trust is proven at commit time — not earned through observation. A proposal that has been accumulating in the improvement namespace for months is no more trusted than one generated yesterday. The creator's signature at commit time is the only thing that confers trust. The observational data makes the proposal better. The commit makes it trusted.

Better models produce better improvement proposals from the same observational data. The improvement loop inherits model gains. No redesign needed.

---

## Related

- `docs/arch/v1/kernel.md` — skill contract schema, improvement_scope field
- `docs/arch/v1/skill-lifecycle.md` — skill versioning, provision_skill
- `docs/arch/v1/crypto-protocol.md` — trust is model-scoped, provisioning + versioning flow
- `docs/arch/v1/memory-structure.md` — observational namespace, committed layer
- `docs/arch/v1/principles.md` — Principle 3 (constrain the data), Principle 7 (system grows with models)
