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
- Nothing lands in an agent's object store without user approval
- Commits require explicit user sign-off via `propose_commit`
- The delegate state machine uses SQLite as source of truth, not just Redis
- External skills are never registered without user approval, no exceptions
- Git is the object store — every commit is auditable, every mistake is reversible

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

## 4. Everything Is an Agent

Shards are infrastructure. Agents are the system. Every capability, every skill, every piece of work belongs to an agent. Shards provide the compute. Agents provide the intelligence.

This unifies the mental model — whether you are talking to a TV, a GPU cluster, or a domain expert about your gym life, you are always addressing an agent. The routing layer handles the rest.

---

## 5. Shards Have Capabilities. Agents Have Skills.

Two distinct registries. Two distinct concerns.

- **Shard capabilities** — what the hardware can do. Voice, display, deep reasoning, storage.
- **Agent skills** — what the agent can do for the user. Send a message, log a workout, turn on the lights.

Skills execute on shards whose capabilities support them. The reasoning layer never sees shard capabilities — only agent skills. Shards are transparent above the routing layer.

---

## 6. No Special Cases in Routing

The same routing model that dispatches `turn_on` to a TV dispatches a distributed inference job across two GPUs. The same `skyra delegate` command Skyra issues at the top of the tree can be issued by any agent anywhere in the tree.

When a design requires a special case, that is a signal the model is wrong. The right model has no special cases.

---

## 7. The System Grows With the Models

Pre-registering every possible skill would be missing the point. Skyra can discover and build new skills using her base skills — Google Search finds the how, Code Execution builds the what, the registry makes it permanent.

As models get better, Skyra builds better skills. The job tree becomes more efficient. Agents reason in fewer steps. The system does not need to be redesigned to benefit from model improvements — it inherits them automatically.

---

## 8. Surface Area Has a Purpose

Every component in this system was pulled in by necessity, not ambition. Every piece of complexity earns its place by solving a problem nothing else could solve.

When evaluating new additions: if an existing component can do the job, use it. If the new component has no job that nothing else does, cut it. Complexity that does not earn its place makes the system harder to build, debug, and extend.

---

## 9. Nodes Are Identity. Edges Are History.

Nodes represent things that exist — entities, facts, skills, domains. They are identity.

Edges represent relationships between things — and relationships change over time. They are history.

The committed layer is append-only. A new edge does not replace an old one. Both exist forever. If Mike married Liz in 2022 and divorced in 2026, both edges are committed facts. The graph holds the complete truth across time.

This is why deletion does not exist in the committed layer. Deleting an edge destroys history. The graph grows in one direction — forward. Query complexity is the tradeoff we accept for a complete, trustworthy record.

**Truth is derived, not stored.** There is no "current state" field. Skyra reasons over edge types, weights, `last_seen_at`, and the full history to derive what is true right now. Truth is a conclusion she reaches — not a value she reads.

---

## 10. The User Makes the System. The System Makes the User.

The system reflects who the user is trying to become — but only what the user has explicitly declared.

**Goals are not inferred. They are not proposed by the system.** They are hand-written committed entities, created by the user, confirmed through an explicit approval process, and removed through an explicit removal process. The system cannot suggest a goal. It cannot promote an observation into a goal. It can only read goals the user deliberately committed.

This is maximum security on the most important data in the system. The definition of who you are trying to become belongs entirely to you. The system has no opinion on it.

Over time the system shapes its responses to the goals the user has committed. And by responding to that reflection, the user becomes more of that person. The system and the user co-evolve — but only in the direction the user chose.

A system that runs for a lifetime doesn't just remember who you were — it participates in who you become. On your terms.

---

## 11. Skyra Is Always Available

Skyra never blocks on work she has delegated. She fans out, she delegates, and she is immediately free for the next user message. The delegate agent owns the lifecycle of delegated work. Skyra owns the user relationship.

A user should never wait because Skyra is busy. That is a design failure.
