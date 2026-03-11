# Skyra Design Principles

## 1. A Bet on Hardware and Model Efficiency

This entire system is a long-term bet. Hardware gets faster and cheaper. Models get more capable and more efficient. Local inference that feels slow today will feel instant in two years. A GPU that costs thousands today will cost hundreds tomorrow.

Skyra is designed for that trajectory, not for today's constraints.

This is why we run local models instead of cloud APIs. This is why the shard model is distributed — as hardware improves, you add a node and the whole system gets more capable automatically. This is why we do not over-optimize for current limitations. A design decision that makes sense for today's hardware but closes off tomorrow's capability is the wrong decision.

The system should get dramatically better without being redesigned. Better models drop in. Faster hardware registers as a new shard. The capability surface grows. Nothing in the architecture fights that — everything is built to welcome it.

---

## 2. Data Integrity First, UX Second


Every design decision is evaluated in this order. When there is tension between a better user experience and data safety, data integrity wins. A fast system that corrupts state is worse than a slow system that doesn't.

Models hallucinate. They always will. No matter how capable models become, they will produce confident, plausible, wrong outputs. This is not a temporary problem to be solved — it is a permanent property of the technology.

The commit model exists because of this. Every state change an agent proposes is surfaced to the user before it lands. The user is the final validation layer. A hallucinated fact, a wrong decision, a misunderstood instruction — none of it becomes permanent without a human seeing it first.

This is why:
- Nothing lands in the committed layer without user approval
- Commits require explicit user sign-off via `propose_commit`
- External skills are never registered without user approval, no exceptions
- The committed layer is append-only — every commit is auditable, every mistake is reversible

**Trust has two axes. Owner trust is proven at commit time by signature — binary, cryptographic, not earned over time. External trust is proven by history — the skill's execution record, commit rate, and user count. The signature proves authenticity. The history proves quality. Both are required.** What is committed is trusted. What is not committed is not trusted. Binary. An observational node that has existed for two years is no more trusted than one created yesterday. Data accumulates. Trust does not. The user's signature at commit time is the only thing that confers trust. When Skyra derives an output from unverified data, she notifies the user. That notification is a UX concern, not a protocol constraint. The commit is the signal.

---

## 3. Constrain the Data, Not the Model

We do not put guardrails on how models reason. We constrain what data they reason over.

Each domain agent operates within a scoped data boundary — its own object store, its own skill list, its own domain context. Inside that boundary the model reasons completely freely. No artificial cap on intelligence, no hardcoded decision trees, no restricted reasoning paths.

This is the right place for guardrails. The data boundary is enforced by the system. The reasoning inside it is left to the model.

**Why this matters:**
- As models improve, agents improve automatically. No re-architecture needed.
- The ReAct loop gets faster and more accurate with better models
- Skill acquisition gets smarter — better models write better skills
- The system's intelligence ceiling is the model's ceiling, not an artificial one we imposed

**Why data integrity is inseparable from this principle:**
The model reasons freely — so the data it reasons over must be trustworthy. Constrained data + unconstrained reasoning = powerful. Dirty data + unconstrained reasoning = dangerous. These two principles are the same principle from different angles.

---

## 4. Root Design Decisions in Reality.

If design is rooted in reality, that is the north star for the model.

---

## 5. Everything Is a Skill

Shards are infrastructure. Skills are the execution layer. Every capability, every piece of work, every system operation is expressed as a skill. Shards provide the compute. Skills provide the contract.

This unifies the mental model — whether the system is replying to a user, running background reasoning, integrating memory, or executing a domain action, it is always instantiating and executing a skill. The routing layer handles placement.

---

## 6. Shards Have Capabilities. Agents Have Skills.

Two distinct registries. Two distinct concerns.

- **Shard capabilities** — what the hardware can do. Voice, display, deep reasoning, storage.
- **Agent skills** — what the agent can do for the user. Send a message, log a workout, turn on the lights.

Skills execute on shards whose capabilities support them. The reasoning layer never sees shard capabilities — only agent skills. Shards are transparent above the routing layer.

---

## 7. No Special Cases in Routing

The same routing model that dispatches `turn_on` to a TV dispatches a distributed inference job across two GPUs. The same `skyra delegate` command Skyra issues at the top of the tree can be issued by any agent anywhere in the tree.

When a design requires a special case, that is a signal the model is wrong. The right model has no special cases.

---

## 8. The System Grows With the Models

Pre-registering every possible skill would be missing the point. Skyra can discover and build new skills using her base skills — Google Search finds the how, Code Execution builds the what, the registry makes it permanent.

As models get better, Skyra builds better skills. The job tree becomes more efficient. Agents reason in fewer steps. The system does not need to be redesigned to benefit from model improvements — it inherits them automatically.

---

## 9. Surface Area Has a Purpose

Every component in this system was pulled in by necessity, not ambition. Every piece of complexity earns its place by solving a problem nothing else could solve.

When evaluating new additions: if an existing component can do the job, use it. If the new component has no job that nothing else does, cut it. Complexity that does not earn its place makes the system harder to build, debug, and extend.

---

## 10. Nodes Are Identity. Edges Are History.

Nodes represent things that exist — entities, facts, skills, domains. They are identity.

Edges represent relationships between things — and relationships change over time. They are history.

The committed layer is append-only. A new edge does not replace an old one. Both exist forever. If Mike married Liz in 2022 and divorced in 2026, both edges are committed facts. The graph holds the complete truth across time.

This is why deletion does not exist in the committed layer. Deleting an edge destroys history. The graph grows in one direction — forward. Query complexity is the tradeoff we accept for a complete, trustworthy record.

**Truth is derived, not stored.** There is no "current state" field. Skyra reasons over edge types, weights, `last_seen_at`, and the full history to derive what is true right now. Truth is a conclusion she reaches — not a value she reads.

---

## 11. The User Makes the Model. The Model Makes the User.

The user commits skills — those commits define how the model behaves. The improvement scopes, the retrieval algorithms, the reasoning constraints — all proven by user signature. The user shapes the model through what they commit. Without commits, the model is raw inference. Untrusted. The user's signatures give it its trusted shape.

The model, operating inside those committed boundaries, reflects the user's reality back to them. That reflection shapes who the user becomes. Their behavior generates new signal, new intent, new proposals — which the user commits — which shapes the model further. A loop, tightened by the commit gate at every turn.

**Goals are not inferred. They are not proposed by the system.** They are hand-written committed entities, created by the user, confirmed through an explicit approval process, and removed through an explicit removal process. The system cannot suggest a goal. It cannot promote an observation into a goal. It can only read goals the user deliberately committed.

This is maximum security on the most important data in the system. The definition of who you are trying to become belongs entirely to you. The system has no opinion on it.

A system that runs for a lifetime doesn't just remember who you were — it participates in who you become. On your terms.

---

## 12. Your Keys. Your Data. Your Consequences.

The system does not protect the user from themselves. It protects them from everyone else.

If the user commits bad data, the system reflects bad data back. If the user sets destructive goals, the system works toward them. The system has no opinion on whether the user's choices are good. It executes what the user has committed.

This is the same philosophy as Bitcoin and Linux. Sovereignty means full ownership — of the upside and the downside. The keys belong to the user. So does everything that follows from how they use them.

The system's job is to make sure no one else can touch what belongs to the user. What the user does with it is entirely their own.

---

## 13. The Model Is a Dependency, Not a Component.

The system does not own the model. The user chooses which model runs underneath — and they own that choice, including its limitations and biases.

The committed layer protects the output. It does not protect the reasoning that produced it. A biased model shapes which questions get asked, how session data gets interpreted, which entities get extracted, how edges get weighted — all before anything reaches the committed layer. That influence is real and cannot be fully guarded against by the system.

The best defense is the user's own choice: open source models, auditable weights, local inference. The architecture is built for exactly that. But the responsibility for what runs underneath belongs to the user.

**Trust is model-scoped.** A skill committed under one model is not trusted under a different model. The user approved that skill in the context of what that specific model produced. Changing the model changes the trust context. Skills committed under a different model are flagged — visible in memory, not executable — until the user re-approves them under the new model. Upgrading the model is not free. The user must re-verify.

This is consistent with the rest of the design. Sovereignty means owning the full stack — including its weakest link.
