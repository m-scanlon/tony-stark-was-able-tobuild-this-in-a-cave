# Ideas — Future Thinking

A place to capture ideas that aren't ready to design yet but are worth revisiting.

---

## Soul Evolution

**The idea:** Skyra's personality and preferences evolve over time through interaction rather than being statically defined once.

Two distinct layers:

**1. Learned user preferences (`preferences.md` or similar)**
As Skyra works with you across projects it gradually encodes what it learns about how you think, work, and make decisions — not as memory items buried in a vector DB, but as first-class documented traits. Things like:
- "Prefers TypeScript over JavaScript"
- "Wants plan approval before any file writes"
- "Likes concise answers, skip the preamble"

These would live in a separate document from `soul.md` — your preferences, not Skyra's identity. Periodically updated, not on every interaction. You'd review and approve changes before they're written, same as any other commit.

**2. Skyra's own soul evolution**
Skyra's `soul.md` — its own identity, values, and voice — should also be able to change over time. Not on every interaction, but gradually, as the relationship between you and the system deepens. What that cadence looks like and what triggers it is a conversation for another day.

**Why two separate documents matters:**
- `soul.md` = who Skyra is
- `preferences.md` = what it has learned about you
- Keeping them separate means Skyra's identity stays stable while your preferences grow independently

**Open questions:**
- What triggers a preference write? Threshold of repeated behavior? Explicit user signal?
- Who proposes the update — Skyra or the user?
- How does this interact with project-specific preferences vs global preferences?
- What does the approval flow look like for soul updates?

---

## Seamless UX-to-Brain Handoff (Continuous Speech)

**The idea:** When the UX model is mid-ACK and the brain's context package arrives with the real answer, the system transitions from UX model output to brain output without the user hearing a break or a non-sequitur.

The key insight is that the TTS layer is the continuity boundary — not the models. The TTS buffer stays fed regardless of which model is upstream. As long as tokens keep arriving, the user hears unbroken speech.

The hard part is semantic continuity. If the brain just starts a new thought, the user hears a seam even without silence. The fix: the brain receives the UX model's partial transcript alongside the context package and generates a *completion*, not an initiation. It picks up mid-sentence.

**The flow:**
```
wake word → UX model starts ACK → tokens feed TTS buffer (audio plays)
                                        ↓
                         brain receives: UX partial transcript + context package
                                        ↓
                         brain generates continuation of the sentence
                                        ↓
                         brain tokens replace UX tokens in TTS buffer
                                        ↓
                         user hears one unbroken voice, start to finish
```

**Why it's hard:**
- Brain must receive the partial transcript fast enough to start generating before the UX model's buffer runs dry
- The UX model should be prompted/tuned to leave sentences open-ended — natural bridge phrases the brain can complete
- The transition point needs to be at a clause boundary, not mid-word

**Open questions:**
- How does the system detect that the context package is ready and good enough to trigger the handoff?
- What's the minimum buffer size to guarantee no gap?
- Does the UX model need to be aware it's generating a handoff, or is that handled at the infrastructure layer?
- What happens if the brain's continuation doesn't match the UX model's sentence direction?

---

## UX Quality as an Emergent Property of the Shard Network

**The idea:** UX capability isn't hardcoded to any device. Shards register hardware capabilities — mic, speaker, GPU, RAM, compute class. The brain decides what runs where and pushes model packages down to the shard. The shard just executes. It doesn't own its models, it owns its hardware.

As more capable devices come online, the brain has more options to route to. UX quality improves automatically — not because the shard brought a better model, but because the brain assigned one.

**The emergent behavior:**
- New device comes online → registers capabilities → brain evaluates and may push a model package to it
- Better hardware joins the network → brain prefers it for higher-quality UX roles
- Dead shard → brain routes to next available shard with `voice` capability
- Multi-room → brain picks the shard closest to the user

UX quality becomes a function of the network's current hardware footprint, not a design decision made at build time.

**Why this matters:**
- Shards are generic — no special-casing for the Voice Shard or any specific device
- Model deployment is centrally controlled by the brain, not managed per-device
- The system gets better as hardware improves without touching the architecture

**Open questions:**
- What does the capability registration schema look like? (RAM, compute class, mic quality, speaker present?)
- How does the brain decide when to push a model package vs use a remote call?
- How does proximity factor in — physical location, network latency, or both?
- If two shards both have `voice`, does the brain pick one or coordinate them?

---

## Skyra OS — Custom Linux Branding

**The idea:** Rather than building a custom OS from scratch, take a standard Linux distro, strip the branding, and make it feel like Skyra's own. The Shard daemon is still what does the work — the OS is just the thin layer it runs on.

**What "Skyra OS" looks like in practice:**
- Custom Plymouth boot splash with the Skyra logo
- Branded login screen (GDM/LightDM theme)
- Custom hostname (`skyra-node-1`, `skyra-voice`, etc.)
- Default wallpaper and desktop env stripped or customized
- Distro branding removed from terminal and about screens
- Non-essential packages removed at image build time

**Why not a fully custom OS:**
A custom OS would reintroduce the exact problem Shards solved — you'd be back to separate OS builds per device type. The Shard daemon already owns capability registration and hardware adaptation. The OS should be thin and generic; the Shard is what adapts to the hardware.

**Why not tie OS config to hardware capabilities:**
Same reason. Hardware-specific OS builds are fragile — hardware changes, the build breaks. Capability logic belongs in the Shard, not baked into the image.

**The right split:**
- OS = minimal, generic, Skyra-branded Linux
- Shard daemon = boots, fingerprints hardware, registers capabilities, configures itself accordingly

For the Voice Shard, this is especially clean since the image is built from scratch anyway.

**Open questions:**
- Which base distro (Raspberry Pi OS, Ubuntu, Alpine)?
- Do all Shard nodes run the same base image, or is there a Voice Shard variant and a desktop variant?
- What's the right point to automate image building?

---

## The Death of the Frontend — Synthesized UI at Runtime

**The idea:** There is no frontend. UI is synthesized at serve time by the brain, based on the shape of the data and the capabilities of the device rendering it.

The brain sends a structured data packet. A model looks at the shape of that data — time series, list of short text, tabular, hierarchical — infers the right UI pattern, and synthesizes a renderer for it. The Shard with a screen executes that renderer against its available primitives.

**Two separate problems:**
1. **Shape recognition** — what kind of data is this, what UI pattern fits it (time series → graph, list of text with metadata → feed, tabular → table)
2. **Capability adaptation** — given that inferred pattern, what can this device actually render (TV gets a rich layout, terminal gets a table, phone gets cards)

**Why this is different from templates:**
The model isn't picking from a fixed set of predefined templates. It's inferring the pattern from the data shape and generating the render description. A Twitter-like feed, a music queue, a server status dashboard — none of them exist as static artifacts. They're synthesized on demand.

**The render target problem:**
The Shard with a screen needs a consistent primitive layer — a small declarative spec of what it can draw (text, image, chart, list, input, etc.) and its device constraints (screen size, interaction model). The brain synthesizes UI that stays within those bounds. The Shard is essentially a runtime for whatever the brain emits.

**Why it matters:**
- No frontend to build, ship, or maintain
- UI gets better as the model gets better — no re-deploy
- Any new device with a screen just registers its primitive capabilities and immediately gets UIs it never had a frontend built for
- The same data renders adequately on a TV, a phone, a terminal, or anything else without a separate codebase for each

**Open questions:**
- What does the primitive capability spec look like? How does a Shard declare what it can render?
- Does the brain synthesize a declarative description (like a JSON UI tree) or actual render code?
- How do you handle interactivity — inputs, gestures, navigation — in a synthesized UI?
- What's the latency tradeoff of synthesizing UI at serve time vs caching synthesized renderers for recurring data shapes?

---

## Multi-Tenancy — Multiple Users in the Same House

**The idea:** More than one person in the house uses Skyra. Each user gets their own Brain Shard instance — their own context, agents, memory, and session. No shared brain state between users.

**Why a second brain, not a shared brain:**
A single brain trying to serve multiple users through a priority queue creates contention at every layer — context, agents, scheduling, model capacity. It's the wrong abstraction. Each brain owns one user. Multi-tenancy becomes "how many `control_plane` Shards are registered" rather than a scheduling problem.

**Where it still gets hard — shared hardware:**
Brains don't fight. Shards do. The hardware layer is still shared:

- **Microphone** — two users, one wake word listener. Who owns it? Whoever triggered the wake word, or the primary user, or first-come-first-served?
- **Screen** — is this display claimed by another brain's session? How does a brain know?
- **GPU** — two brains want VRAM at the same time. The GPU Shard needs to queue requests from multiple brains.
- **Speakers** — can two brains talk at once? Who gets priority?

**The coordination problem:**
If Shards can serve multiple brains, there needs to be a coordination layer between brains — or something above them — that arbitrates Shard access. Brains need to know about each other, or a neutral arbiter does. Neither is simple.

**v1 assumption:**
Shards are owned 1:1 by a brain. One user, one brain, one set of Shards. Multi-tenancy is deferred entirely. The architecture doesn't prevent it — it just doesn't solve it yet.

**Open questions:**
- Does a Shard register with one brain or can it advertise to multiple?
- What's the session ownership model for shared-space Shards (living room Pi, shared screen)?
- Is the coordination layer a new service above brains, or a peer protocol between brains?
- How does user identity get established at the Shard layer — voice recognition, proximity, device association?

---

## Soul-Personality Spectrum — A Weighted Graph Model

**The idea:** Soul and personality aren't two separate documents or two distinct systems. They're opposite ends of a single spectrum, and everything lives somewhere on it.

- **Personality** is the top layer — high influence on output shape. Tone, pacing, humor, how something is said. Low influence on what gets decided.
- **Soul** is the bottom layer — high influence on decision gates. What gets refused, what gets prioritized, what gets flagged. Low influence on surface presentation.

The structure isn't a list of rules or two flat files. It's a graph where every node is a trait or value, and every node carries two weights:

- **Decision weight** — how strongly this node influences what the system does
- **Presentation weight** — how strongly this node influences how the system says it

Soul-heavy nodes (e.g. "never fabricate", "honesty over confidence") carry high decision weight, low presentation weight. Personality-heavy nodes (e.g. "dry humor", "concise pacing") carry low decision weight, high presentation weight. Middle nodes pull both ways — "steady progress over bursts" affects what gets recommended *and* how urgency gets framed.

**Why the graph structure matters:**

When the system evaluates a decision, it traverses the graph summing decision-weighted signals. When it shapes a response, it sums presentation-weighted signals. Same graph, different traversal. No separate systems.

Conflict resolution falls out naturally. If two nodes pull in opposite directions on a decision, the weights determine which wins. No explicit priority rules needed — the structure handles it.

**The decision evaluation stack:**

This graph plugs directly into how decisions get evaluated. A complete decision check runs across multiple axes:

1. Does it satisfy the immediate goal?
2. Does it advance the long-term trajectory — or create debt?
3. If it's wrong, is that recoverable? What's the blast radius?
4. Does it conflict with a core value?

The values check (axis 4) is not a tradeoff — it's a veto. A decision can be right on every other axis and still be wrong if it violates a high-weight soul node.

**The time horizon axis:**

Related but distinct: some decisions satisfy the immediate condition while actively damaging a longer-term goal. Tech debt is the obvious case — works now, costs later. But there's also asymmetric consequence: a decision that's "good enough" in the moment but has an irreversible downside if wrong should always be caught regardless of how well it scores on other axes. The graph doesn't automatically encode this — it needs a separate reversibility signal alongside the soul/personality weights.

**This has to be trained, not prompted:**

Context can remind a model of its values. Training *is* the values. You can't just inject soul.md at runtime and expect it to actually shape decisions — it nudges, it doesn't form. The soul document becomes training data annotation criteria: here's the goal, here are the options, here's why this decision aligns and this one doesn't. The graph is what gets encoded into the model's priors through fine-tuning, not what gets passed in a system prompt.

**Node weights can shift over time:**

If nodes have weights, those weights can drift — not the values themselves, but how strongly they pull. That's a more precise model of how a system (or a person) actually matures. You don't stop valuing honesty. But "steady progress" might pull harder as the system deepens. Weight drift is the right mechanism for evolution, not rewriting the soul document.

**Open questions:**
- What's the right representation for the graph at training time — structured examples, annotations, something else?
- How do you detect when a node weight has drifted enough to warrant review vs. normal variation?
- What triggers a weight update — accumulated decisions, explicit signal, periodic review?
- How does this interact with the Soul Evolution idea — are preference updates just personality-layer weight shifts?
- Who can propose a soul-layer weight change, and what does approval look like?

---

## Decision Routing — Tiered Evaluation Depth

**The idea:** Not every decision should run through the full evaluation model. A light switch doesn't need a time horizon check. Routing the decision to the right evaluation depth before any real evaluation starts is just as important as the evaluation model itself.

**The tiers:**

| Tier | When | What runs |
|---|---|---|
| 0 | Deterministic, low-stakes, reversible | Execute directly. No evaluation. |
| 1 | Simple with constraints | Allowlist check, agent scope check. Then execute. |
| 2 | Stateful or downstream effects | Goal alignment + reversibility check. No values layer. |
| 3 | Consequential, irreversible, or sensitive | Full stack — goal, trajectory, reversibility, values. |

**The classifier has to be cheap:**

You can't run a deep model to decide what depth to evaluate at. The classifier runs before any real evaluation starts — it has to be fast and shallow. Rule-based or a tiny classifier. The signal it uses is mostly knowable upfront: tool type, affected scope, reversibility of the action. These are properties of the tool registry, not things inferred at runtime.

**Why this matters:**

Over-evaluating cheap decisions kills latency and adds noise. Under-evaluating consequential decisions is how things go wrong. The right model is not "always run the full stack" or "always run the shallow stack" — it's routing each decision to the minimum evaluation depth that's appropriate for its stakes.

**Connection to existing architecture:**

The Pi's `triage_hints` already classify requests by latency class — that's the same instinct applied to routing compute. This extends it one level deeper: not just where the job runs, but how deeply it gets evaluated once it gets there. Tier classification could live in the tool registry as a property of each tool, resolved at hydration time.

**Open questions:**
- Who owns the tier classification — the tool registry, the Internal Router, or the Estimator?
- Can a job span multiple tiers (e.g. one step is Tier 0, another is Tier 3)?
- What triggers a tier upgrade mid-execution — unexpected state, denied tool, confidence drop?
- Does the user ever see the tier? Or is it purely internal?

---

## Decision Tools — Model-Accessible Evaluation Instruments

**The idea:** Don't enforce a rigid decision pipeline. Give the model tools it can reach for when it decides it needs them. The model self-selects based on judgment — not a pre-classifier routing it.

Examples of what these tools might look like:
- `evaluate_consequence(action)` — what's the blast radius if this goes wrong?
- `check_trajectory(action)` — does this advance the long-term goal or create debt?
- `check_values_alignment(action)` — does this conflict with a core value?
- `escalate_to_user(reason)` — surface this before proceeding

A light switch calls none of them. A consequential, hard-to-reverse action might call several. The model decides.

**Why this is better than a pipeline:**
Models are getting better. A rigid classifier that pre-routes decisions becomes a bottleneck and fights model capability over time. Framing evaluation as tools means the system improves as the model improves — no architecture changes required. You can add lightweight hints at tool registration time for weaker models and remove them as capability grows.

**Open questions:**
- What's the right granularity? One tool per evaluation axis, or a single `evaluate(action, axes[])` call?
- How do hints at tool registration work — metadata on the tool schema, or a separate policy layer?
- Does the model need to explain why it called (or didn't call) an evaluation tool? That matters for auditability.

---

## Decision Telemetry — Spans, Ending Conditions, Audit Trail

**The idea:** Let the model operate with relatively open judgment, but record everything. Every decision is a span — a start, an end, and the signals in between. That data is the foundation for auditing decision quality and eventually training on it.

**The span model:**

Each decision (or job) emits a telemetry span:
- **Start** — what was the intent, what tier/tools were invoked, what context was present
- **End** — how did it resolve, what was the outcome state

**Ending conditions:**

| Signal | Type | Strength |
|---|---|---|
| User correction or undo | Explicit | Strong negative |
| User follow-up to fix something | Implicit | Weak negative |
| User frustration signal | Implicit | Negative |
| Job completed, never revisited | Implicit | Neutral / positive |
| Subsequent job reverses this one | Outcome | Negative |
| User explicit approval | Explicit | Strong positive |
| Result committed to agent state | Outcome | Positive |
| Commit followed by revert | Outcome | Negative |
| No commit produced (expected one) | Outcome | Weak negative |

**The long-horizon problem:**

Some decisions don't surface outcomes for days or weeks — tech debt, trajectory drift. Those probably require manual annotation early on. Skyra flags it, the user reviews it, labels it. Slow but honest. The data needs to exist before you can train on it.

**Why capture now even without a training plan:**

Telemetry data can't be retroactively captured. The annotation criteria and training approach can be figured out later. The spans can't. Start logging everything — intent, evaluation tools called, outcome signals — and the dataset exists when the training strategy is ready.

**Storage options:**
- **ClickHouse** — columnar, fast analytical queries, handles high write volume. Best fit for querying decision patterns across large span history (e.g. "all jobs where X tool was called and result was reverted"). Runs on a single node.
- **Grafana stack (Tempo + Loki + Prometheus)** — full observability picture. Tempo for traces, Loki for logs, Prometheus for metrics. OTel-native. Good dashboards out of the box. More infrastructure but standard open source setup. Can sit on top of ClickHouse.

**Open questions:**
- What's the minimum span schema? What fields are required vs. optional?
- How do you connect a future corrective job back to the original decision that caused it?
- Who owns span storage — the Job Registry, a separate telemetry sink, or both?
- What's the review UX for manual annotation? Does Skyra surface candidates, or does the user browse?
- How do you handle spans that never get a clean ending condition?

---

## Long-Term Memory Store — PostgreSQL + pgvector

**The idea:** Distill committed observations and facts into a durable long-term memory store separate from the operational vector DBs. Long-term memories are structured, queryable, and persistent — not just embeddings floating in a vector index.

**Why a separate store:**

The existing vector DBs (Chroma, Qdrant) are operational — tool retrieval, short-term context, session relevance. They're optimized for fast semantic search during a session. Long-term memory has different requirements: durability, structured filtering, and the ability to ask precise questions across the full history of what the system knows.

**Why PostgreSQL + pgvector:**

Long-term memories aren't schema-flexible documents — they have consistent structure (content, agent, importance score, timestamp, source, access history). SQL is the right query model for the questions you'll actually ask:

- "All memories for agent X with importance above 0.7"
- "Memories not accessed in 30 days that are candidates for decay"
- "Everything the system knows about this topic, ranked by importance"

pgvector handles semantic retrieval in the same DB. No need for a separate vector store for this layer.

NoSQL isn't bad but doesn't have an advantage here — the flexibility it offers doesn't help when memories are structured, and SQL is more expressive for the queries that matter.

**The split:**

| Store | Role |
|---|---|
| Chroma / Qdrant | Operational vector search — tool retrieval, short-term context, session relevance |
| PostgreSQL + pgvector | Long-term memory — distilled, committed, durable facts |

**Connection to importance vectors:**

Importance scores map naturally onto PostgreSQL columns — filterable, sortable, queryable. The importance vector work already defines the signal; this is where it gets persisted and queried against.

**Open questions:**
- What triggers promotion from operational vector store to long-term memory? Importance threshold? Time? Explicit commit?
- What's the schema for a long-term memory record?
- How does memory decay work — soft delete, importance decay, archival?
- How does the context engine query long-term memory vs. the operational stores?
- Who owns writes — the Context Engine background loop, the Agent Service, or both?

---

## Context Engine Makes the Front-Door Model Cheap

**The idea:** Because the context engine does its heavy lifting in the background — continuously watching, inferring, and committing — the model sitting at the front door doesn't need to be powerful for most conversations. It receives a pre-assembled, pre-reasoned context package. The hard cognitive work already happened.

This inverts the usual assumption about voice AI: that you need a capable model on the hot path. You don't. You need a capable background loop. The front-door model's job is routing and synthesis, not reasoning from scratch.

The implication: a 3B parameter model on a Pi can carry on a competent conversation because the context package handed to it already contains the domain knowledge, the relevant memory, the active job state. The model is pattern-completing over good data, not doing raw inference.

**Why this matters:**
- Lower power consumption on always-on edge devices
- Lower latency on the hot path — smaller model, faster tokens
- Gradual shard escalation becomes natural: start cheap, promote to a heavier model only when complexity actually warrants it

**Open questions:**
- What's the right signal for "this conversation has outgrown the front-door model"?
- How does the handoff feel to the user when the conversation promotes to a heavier shard?

---

## Dynamic Shard Escalation — Conversations Self-Promote

**The idea:** The context window is job state. Planning and execution happen in one LLM session. That session can live on any shard. Which means when a conversation gets more complex than the current shard can handle well — or simpler than what it's running on — the system can migrate the session to a better-matched shard mid-conversation, seamlessly.

The user starts a casual exchange on the Pi's 3B model. The topic shifts toward heavy reasoning. The Estimator sees the complexity signal, puts the job back on the heap with updated complexity, and the External Router places it on the GPU shard instead. The user experiences unbroken continuity. The system silently promoted.

This also runs in reverse: a session that started heavy but resolved can demote to a cheaper shard for follow-up. The network self-regulates.

**The key insight:** The architecture already supports this. Because the context window is the state, "moving a session" is just moving a context window. The infrastructure — heap, External Router, capability profiles — already handles it. This isn't a new feature; it's an emergent property of the existing design.

**The dynamic escalation flow:**
```
session on Pi 3B (low complexity)
  → conversation shifts, complexity signal rises
  → Estimator re-evaluates, promotes job on heap
  → External Router places to GPU shard
  → context window transferred
  → session continues — user hears no seam
```

**Open questions:**
- What's the trigger threshold for promotion/demotion — complexity estimate, confidence drop, explicit shard hint?
- How does context window transfer work at the wire level — serialized state, re-hydrated session, or something else?
- What happens to in-flight tool calls during a mid-session shard transfer?
- Can the user override — "stay on the fast model, this doesn't need heavy reasoning"?

---

## Multi-Location Spatial Awareness via Network Co-Location

**The insight:** The ingress shard is the location anchor. It's physically in the space with the user, on the same network as every other shard at that location. When a request comes in, the capability resolver doesn't need to infer location from the request text — it looks at which network the ingress shard is on and matches capabilities from shards sharing that network fingerprint. Spatial awareness falls out of network co-location.

**How it works:**
- At registration, every shard fingerprints its network: SSID, gateway MAC, local subnet. This becomes its location tag.
- First time a shard comes online on an unknown network, the system prompts once to name the location ("what should I call this?"). After that, automatic.
- The capability resolver uses the ingress shard's network tag as the default location filter. All shards sharing that tag are "in the room."

**The flow:**
```
"turn the TV off"
→ came in through living room Pi (network: home-main)
→ resolver filters capabilities to shards on home-main
→ finds TV shard
→ dispatches
```

Travel to vacation home — the Pi there is on `vacation-home` fingerprint. Same request resolves to the vacation home TV automatically. No location specified, no disambiguation. The ingress shard is the location.

**Remote control (explicit cross-location):**
"Turn off the vacation home TV" while at home — the user names the location. Resolver filters by that location tag instead of the ingress shard's network. Same resolver, different filter input.

**Open questions:**
- What happens when a shard is reachable but on a different network (VPN, tunnel)? Does tunnel origin or physical network take precedence?
- How does the system handle a shard that moves networks — laptop taken to vacation home?
- What's the location naming UX — voice prompt, app, something else?

---

## More ideas to add here as they come up
