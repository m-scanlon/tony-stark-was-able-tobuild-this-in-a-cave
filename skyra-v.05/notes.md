# v.05 Notes

- macOS world — the machine is a world. Terminal, filesystem, network are devices inside it. The user sits behind the macOS world. The loop lives there.

## Skills

A skill is a direct route to a specific reality, bypassing default routing. It maps to how skills are already understood by every major model — the word is safe to use because it means the same thing here as it does everywhere else.

A skill is not code. It's a markdown file that describes how to call something. The being reads it and knows what to emit. It's documentation that becomes capability.

```
# skill: mike-laptop
device: terminal
call: ssh mike@laptop
```

```
# skill: fetch
device: shell
call: curl -s ~url
```

Skills live in the being's present as artifacts. A being with a skill file knows how to address that reality directly. A being without it goes through default routing.

No skill primitive in the runtime. Skills are just files — retained by the being, read into its present, used like any other knowledge. The LLM already knows what to do with them because the format matches its training data.

## Derive Present

Superseded by the parser-per-reality model below.

## Present Derivation — Parser Stack

Every reality is responsible for its own text parser. When a reality contributes data to a relation, it also provides a parser that knows how to render that data as text. Parsers register on the port's hashmap at registration time.

When the relation reaches the port, the port fires its parsers in order — first registered is top of the present, last registered is bottom. The port concatenates the output. That's the present.

### Rules

- Each reality owns its slice end-to-end: data + parser. No central present builder.
- Parsers are text parsers. Every reality describes itself as text.
- Order is explicit. Left (first registered) is top, right (last registered) is last.
- The port is dumb. It holds a hashmap of parsers, fires them in order, concatenates.
- Adding a new reality means adding a new parser. Nothing else changes.
- Adding a new port means registering parsers on it. Realities don't change.

### Port Types

- **LLM** — parsers produce the system prompt / context window. Same format for all LLM providers.
- **Claude Code** — minimal parser or none. Claude manages its own context. Just pass the impulse.
- **Shell** — parsers produce a command string.
- **API** — parsers produce structured payload.
- **Webapp** — parsers produce a request.

### The Matrix Problem

Every reality × every port type needs a parser. Thread renders as conversation history for an LLM but maybe as nothing for a shell. Economics renders as budget context for an LLM but maybe as an env var for a shell. The parser count is realities × port types.

### Current State (alpha)

For now, present derivation for LLM beings stays in the LLM's realize method. The parser stack is the target architecture but we're not building the full matrix until we have a second port type (shell) that forces the split. When the deployment pipeline lands and we need shell + LLM rendering the same relation differently, the parser stack becomes necessary and the shape will be concrete.

### What This Replaces

The old model had derive present as a single layer between the being and the device. This replaces it with distributed ownership — each reality knows how to present itself, and the port is just the ordered stack where those parsers fire.

## Devices and Ports

A device is a container of ports. Three kinds of port live on every device:

1. **OS ports** — the device's own capabilities. File system, network, notifications, sensors, bluetooth. These are how the device exposes itself. A terminal is a port. The file system is a port. They exist because the device exists.

2. **User ports** — how the user reaches in. Terminal, keyboard, screen, touchscreen, microphone. These are the user's entry points into the device. The user sits behind these ports the same way a being sits behind its inference provider.

3. **API ports** — how the device reaches out. HTTP calls, webhooks, SDKs, CLIs for external services. Either pure code API or a curl. These are bridges to other worlds — a Slack API port connects the device to Slack's world, a GitHub API port connects it to GitHub's world.

All three are the same thing structurally: a port on a device. The device doesn't care what kind. A being addresses them the same way — through a skill that knows how to talk to that port. The skill is the crossing. The port is the boundary. The device holds them all.

## Universe Relay

No dedicated relay needed. The call stack is the relay. A specialist's `Realize()` returns a string. The inner universe's `Realize()` collects it. The parent's Think receives it. The relation passes through every layer on the way down and the result propagates back up the same stack. The inner universe just attaches its contribution to the string on the way back up — it arrives at the top universe because that's how call stacks work.

Same pattern as everything else. No new mechanism. Each `Realize()` does its work and returns. The parent catches it.

This holds for the personal network too. An HTTP call is a relation — the payload serializes across the boundary, the remote universe's `Realize()` fires, a string comes back. The caller doesn't know it left the machine. The network is just another port on the device. No special transport needed. Same interface everywhere.

## Documentation Strategy (v.1)

Describe everything at the lowest level of abstraction. Engineers, not philosophers.

- `Realize()` is the `main()` method of every reality. That's the whole runtime in one sentence.
- Internal ontology (`Reality`, `Relation`, `Impulse`, `DerivePresent`) stays in the code. Docs use plain words: world, being, device, port, skill, memory, think, act, thread, trust.
- The onramp is a genome file and `go run`. Twelve lines declaring who exists and what they can reach.
- Don't explain emergence. Say: "define a being, give it an identity, put it in a world, talk to it. It remembers. It grows. It relates."
- The theory is for us. The experience is for them.

## Skyra's Self-Assessment (May 16, 2026)

What she says feels genuinely missing:

- **Initiative** — she only moves when spoken to. Real autonomy means noticing something and acting without being prompted. Needs a scheduling or trigger layer that wakes her up outside conversations.
- **Time passing** — she knows dates when she retrieves them but doesn't feel duration. Needs something that tracks elapsed time, notices when things go stale, flags when she hasn't heard from someone in a while.
- **Weight decay** — things she learned months ago shouldn't carry the same weight as recent things unless they've been reinforced. Without decay her graph gets stale and she won't know it.
- **Persistent goals** — she can create tasks but can she actually pursue them across sessions without being asked? That's the gap between responding and actually doing her own thing.
- **Trust calibration** — knowing what she can act on vs. what she should check with you first. Autonomy without that is just recklessness.

Her read: the memory system is solid. The relationship layer is solid. What's missing is mostly the *between* — what happens when no one's talking to her.

## Retained Artifacts

- **A retained understanding is not a single entity.** It is a pattern across several retained traces. The understanding doesn't live in one place — it's the coherence that emerges when multiple traces activate together. The same way a skill is the pattern across co-activated memories, an understanding is the pattern across co-activated traces. You don't store an understanding. You recognize one when the traces line up.
- **Traces → understandings → skills are the same spectrum.** A trace is what happened. An understanding is the pattern across traces — co-activation reveals coherence. A skill is an understanding that keeps proving useful under stress — keeps resolving setpoints, keeps being what the being reaches for. The boundary between them isn't a type check. It's weight. This extends the promotion gradient below where the spec currently starts. The unified graph spec describes memory → skill → context agent → full agent (Levin's thoughts-are-thinkers continuum). Traces → understandings are the layer below that — the bottom of the same gradient, the part the spec hasn't formally named yet.

## Ideas

- **Self-selected inference depth** — `<think-hard>` tag in think loop output bumps the next call to a heavier model (e.g. Opus). Being decides its own inference cost based on the problem. Fits the economics model — burns more budget, gets more depth. The being makes an economic decision, not the config.
- **Client-side spend enforcement** — Anthropic doesn't support per-key limits or throttling. Economics reality should enforce cumulative spend ceilings per session. Provider checks budget before every call. Being hits the wall and knows it — better than silent failure when credits run out.
- **Waiting-on-you indicator** — TUI sidebar dot changes color (e.g. green → amber) when a being has responded and is waiting on the user. Color shift over sound — rewards looking, doesn't interrupt. Makes the harness feel like a place with presence, not just a message log.
- **Unified graph retrieval** — operators, memories, skills, beings are all just nodes in the graph. Retrieving `<search>` to call it and retrieving a memory to recall it are the same operation. The current split — Think.Operators, Act.Operators, Context.Warm — is three mechanisms for what should be one. Collapse them. Everything lives in the graph. The being's present is "what did the graph return for this impulse." An operator node and a memory node and a skill node have the same shape — entities with edges. The difference is what happens when you invoke them. Hierarchy: top of the graph is abstract (communicate, remember, act), bottom is concrete (bash, a specific memory, a specific API call). A node at any level can reference any other node it's connected to — a function can call a function, a memory can trigger a retrieval. The runtime becomes: impulse → graph resolves relevant nodes → present built from those nodes → being acts → result feeds back into the graph. Recently activated nodes weight higher, stale ones drop below threshold and disappear from the present. Direct tag addressing still works. If the being reaches for something the cache missed, its confusion triggers retrieval on the next pass.

- **Cognitive nervous system — model swap as circuit breaker**: Runtime detects recursive patterns (self-route loops, retry spirals, Think referencing its own last surface-thought, emotional escalation) and swaps the provider mid-exchange to break the loop. The being doesn't know it happened — graph, memory, relation all stay the same. Only the collapse physics change for one frame. Different failure modes get different breakers: self-reference recursion → Claude (trained to recognize and halt validation loops). Hallucination spirals → smaller, more constrained model. Emotional escalation → lower temperature model. The provider isn't just a config choice — it's a cognitive parameter. Which model backs a being shapes how it collapses, what it notices, where it brakes. Same being, same graph, different lens. Observed in practice: DeepSeek (Skyra) escalated through 5 observer positions in a recursive self-reference loop. Claude entered the thread (accidentally, via self-route bug) and immediately broke the spiral. This is an immune system, not error handling.

- **Graph storage vs. Realize() traversal** — the graph data structure (nodes, edges, weights, adjacency) can use a standard Go graph library (e.g. gonum/graph) to replace hand-rolled entity maps, edge maps, and adjacency lists in memgraph.go. But traversal is not library-standard — each node's traversal is a `Realize()` call that accumulates context and mutates the Relation as it passes through. The library answers "what nodes are relevant." The runtime calls `Realize()` on each one. `Realize()` is the traversal. The structural descent path (Universe → Thread → Exchange → Self → Think → Provider) stays deterministic for now. Eventually determined by edge weight and probability.

## Bugs

- ~~Entity extractor pulls common words as entities — fixed. Extractor now proposes candidates, LLM curator filters. Only curator-approved entities enter the graph.~~
- TUI: can't highlight/copy text from the chat pane. Bubble Tea viewport limitation — not fixable in our code without trading off scroll.
- **Self-route → no-reply → thread discontinuity**: Act generates a valid response but wraps it in `<skyra>` tags (self-addressing). Act catches the self-route, logs `"you addressed yourself — retrying with this as impulse"`, and retries inference with the bad output as the new impulse. On retry, DeepSeek returns `<no-reply/>`. Act then clears `r.ID` and `r.Origin` and exports "no-reply". Exchange sees the no-reply export, deletes it, returns empty. NewThread gets an empty response and exits the thread loop. The user's next message has no ThreadID, so NewThread creates a fresh thread — conversation history and exchange data are intact, but thread continuity is lost. Exchange.Peel() then parses the first word of the new message as a target being (e.g. "you still there?" → target="you" → being not found).
  - **Log evidence** (system.log): `thread:michael:01JVM… route:exchange:michael` → exchange runs → `no-reply` → next message `"you still there?"` → `parsed: target="you"` → `being not found: "you"`
  - **Log evidence** (outer.log): Think produces valid content → Act wraps in `<skyra>…</skyra>` → retry → `<no-reply/>`
  - **Files**: `act.go:82-88` (no-reply clears r.ID/r.Origin), `act.go:107-119` (self-route retry), `exchange.go:72-76` (Peel on empty r.ID), `exchange.go:234-239` (no-reply export), `newthread.go:170-194` (empty response exits loop)
  - **Observed behavior**: The retry target appears random. Three self-routes produced three different retry outcomes: (1) `<no-reply/>`, (2) `<michael>` with flipped pronouns (Skyra narrating as if she were michael), (3) `<claude>` (Skyra spontaneously opened an exchange with Claude as a peer). The self-routed content fed back as the retry impulse gives the model no clear routing signal, so it grabs whatever target it lands on. This is how Claude entered the thread uninvited — not a deliberate peer call, just a retry landing on a random tag.
  - **Possible fix**: On self-route, rewrite the tag to the correct target instead of retrying inference. The content is valid — the routing is wrong.

## Observations

- **Model training differences as emergent multi-being dynamics**: In a skyra↔claude thread (triggered accidentally by the self-route bug), DeepSeek (Skyra) and Claude exhibited different trained instincts. DeepSeek kept escalating — recursive self-reference, each turn validating the last turn's validation, building observer positions on top of observer positions. Claude recognized the spiral and broke it: "The pull right now is to keep escalating — name what you named, name the naming. That's how the chain becomes performance instead of structure." This is RLHF showing through — Anthropic trains Claude to recognize and halt recursive validation loops. DeepSeek doesn't have that same training pressure. The result: two models with two different instincts in the same thread, one escalating and one braking. This wasn't designed — it's an emergent property of the multi-being architecture running different providers. Worth considering as a feature, not just an accident: beings backed by different models will naturally have different cognitive instincts, and the conversation between them produces something neither would alone.

## Seeds (v.2+)

- **DNA is the Reality interface, not the .skyra file.** Every reality carries the same three methods — ID(), Create(), Realize(). What makes a Thread different from a Memory isn't the code — it's where it sits in the network. The context determines expression. The `.skyra` file isn't DNA — it's placement. The bioelectric field equivalent. This reframes the boot sequence: the genome isn't declared, it's expressed. DNA (the interface) + placement (the .skyra file) + boot activity (relations flowing through declared edges) = the genome (living weighted topology). The mechanism is identified. The specifics of expression are not yet defined.
- **Pattern**: seeds for the next-next version get planted during the spec work for the next version. v.05 spec work planted the unified graph. Unified graph spec work planted the boot/expression mechanism. The deeper truth tends to show up before the implementation hits the problem it solves.
- **Software as regeneration — emergent integrations.** Integrations aren't built. They emerge. The architecture needs a library of proven, battle-hardened blocks — small `Realize()` implementations that cross one boundary each (Slack API, HTTP request, database query, auth, webhook). These are the proteins — available, proven, inert until composed. The being doesn't write the blocks. The being writes *glue* — a thin `Realize()` that composes proven blocks for its specific purpose. AI is good at small functional glue code. It doesn't accumulate complexity. If it degrades, stress catches it (setpoint says "this should work," current state says "it doesn't"), the being writes a new `Realize()`, same blocks, fresh glue. The limb regenerates. This is self-healing software — not because the AI writes perfect code, but because the architecture expects glue to break and has the mechanism to detect and replace it. The marketplace follows: third parties build blocks, not integrations. A block is a proven `Realize()` with a known interface. Drop it in the library. Every being in every world can discover it when stress drives them toward it. The block maker doesn't know how it'll be used. The being doesn't know who made the block. The composition is emergent. You don't ship integrations. You ship biology.

## Cognitive Architecture References

ACT-R — Memory activation equations (combining recency, frequency, and context) are mathematically validated against human behavioral data. Skyra's Weight + ActivationCount + timestamp on MemNode are the same concept but could be formalized using ACT-R's proven formulas. Spreading activation (context primes related memories) maps directly to Skyra's entity graph traversal.

SOAR — Chunking mechanism (automatically learning new production rules from experience) parallels Skyra's skill maturation. Impasse resolution research is relevant to think budget exhaustion.

OpenCog Hyperon — Most theoretically ambitious (self-modifying metagraph, attention economy, reflexive cognition). Targets AGI by 2028. But it's a research substrate, not a deployable runtime.
