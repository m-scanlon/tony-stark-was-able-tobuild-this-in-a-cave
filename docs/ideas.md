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

## More ideas to add here as they come up
