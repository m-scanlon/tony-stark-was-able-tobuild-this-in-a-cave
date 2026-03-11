# Skill Licensing

> **Not MVP.** This is a future design concept. The ideas here are not yet locked.

---

## The Model

A skill can be sold without granting the buyer read access to its implementation. The buyer receives execution rights — the skill runs, produces results, and behaves according to its contract. The buyer cannot inspect the internals.

This is a black box execution model. The skill is trusted because the system validates the contract — not because the buyer can audit the code.

Definition privacy is separate from runtime observability. If a skill performs shell/API actions, those actions can be audited or observed at the execution boundary even when implementation details remain closed.

---

## The Expiration

The read restriction has an expiration date enforced by the system. When it expires:

- The skill's implementation becomes auditable — readable by the owner
- The author cannot extend the restriction
- The expiration is non-negotiable

The analogy is copyright expiration. The restriction is real and respected until the date passes. After that, the skill enters a kind of public domain: still executable, now inspectable. The system enforces this — not the marketplace, not the author.

---

## Why This Matters

Without expiration, a black box skill is a permanent liability. The owner can never audit what is running on their machine. They are trusting the author indefinitely.

With expiration, the trust is time-bounded. The owner knows that at some point they will be able to inspect what they ran. The author knows this too. It changes the incentive structure.

---

## Open Questions

- Who sets the expiration date — the author, the marketplace, or the system?
- Is the expiration date visible to the buyer before purchase?
- What does "auditable" mean in practice — source access, decompiled bytecode, natural language description?
- How does expiration interact with the trust model — does an expired black box skill require re-approval?
- Can a skill be re-sold with a new expiration, or does the original expiration follow the skill?
- How does this interact with the violation policy — if a black box skill violates the introspect contract, can malpractice be proven without source access?

---

## Related

- `docs/arch/v1/skill-lifecycle.md` — skill provisioning, trust model
- `docs/arch/v1/introspect.md` — violation policy, verifiable malpractice
