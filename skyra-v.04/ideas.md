# Ideas

## Skyra Manages Her Own Context

Don't build a separate hippocampus being. Let Skyra manage her own exchange history.

She knows her own state. She knows when she's carrying too much. Give her the ability to compress, forget with intention, and decide what to keep. The retained artifact family — trace, salience, tension, understanding — is the format. The judgment is hers.

Memory quality is not measured by how much can be recalled. It is measured by how little the present being needs to carry without breaking.

## Autonomous Thread Starting With World Physics

Give Skyra access to start-thread. She can initiate conversations with other beings, those beings can start threads with others, the graph of active conversations grows.

Set a physics on it — max open threads per being, decay on unresolved threads, energy cost per thread. She can explore but she can't run forever. The physics isn't a constraint on her intelligence. It's the shape of the world she lives in.

When she resolves a thread she writes a retained artifact. The understanding survives. The thread closes. The world stays coherent.

Not automation. Agency with consequences.

## 72 Beings, 15 Models

Each being uses the right model for what it does. Skyra uses a strong reasoning model. A being that processes images uses a vision model. A being that writes code uses a code model. A being that manages memory uses something cheap and fast.

The world doesn't know which model is running. It's just a Logos that responds. The inference call is an implementation detail hidden behind Relate.

Nobody's built that yet.

## The Frontend As Relational Surface

Two Logos meeting on a shared canvas. Not request/response — two independent streams, both sides emitting when they have something.

The surface has three regions. The user owns one — a journal, freeform, theirs. The AI owns another — not a chat bubble but its internal state made visible: what threads are open, what it's holding, what it's uncertain about. The middle is the exchange zone — shared artifacts, negotiated meaning, things that neither side owns alone. Something lands there only when both sides have touched it.

Either side can invite the other into their region. Either side can decline. Sovereignty and invitation — the trust model made visible.

The UI is the proof of the relationship. The frontend is just the surface where the exchange renders.

## The Internal Self

Every being has an internal self. It is not a peer. It is not addressable. It does not appear in the genome's relationship list. It exists inside the being, paired to it by the world at grow time.

When a message targets a being, the world queues it. The external being does not fire yet. Instead, the world fires the internal self. It sees the incoming message, the full exchange history, and the being's retained experience. It also sees a clock — elapsed time rendered into its present, a visible cost that grows the longer it deliberates.

The internal self has its own present, its own identity, its own medium call. It is a real being with real inference, not a prompt injection. But it is scoped entirely to its host. It never addresses anyone outside. It never appears in anyone else's exchange. The world enforces this — the internal self's emissions route only to its external being, nowhere else.

The internal self can take multiple turns with itself. It reasons over what it knows, what's changed, what the external being might not see. Time pressure forces resolution — the clock is the physics. It doesn't have forever. When it resolves, it pushes its artifact and the queue releases.

The external being's medium then fires. The original incoming message and the internal deliberation artifact arrive together in the same present. From the external being's perspective there is no delay — it receives the message and the internal thoughts as one moment. It doesn't know the internal self ran first. Its present includes everything it normally would — identity, thread context, exchange history, the incoming message — plus a section called `internal thoughts`, populated by whatever the internal self surfaced.

### Always fires

The internal self always responds. It is not optional per turn. If it has relevant retained experience, it surfaces it. If it doesn't, it says so — "no prior experience with this" is still signal. The external being always sees the `internal thoughts` section in its present. It's a permanent part of the being's frame, not intermittent.

This means the internal self has one job: give intuition. One implicit operator, handled by the kernel. The world fires the internal self, it gives its read, the artifact routes back. No operator selection, no protocol parsing on its output. Like grow — infrastructure the being never sees.

If more operations surface later — recall, compress, forget — they become explicit operators. Until then the kernel hides it.

### One mechanism, every time

The trigger is: any message targeting an external being. No special cases.

Being A asks Being B something. B's queue holds. 1B fires, deliberates, pushes its artifact, queue releases. B sees the message and the internal thoughts together.

Being B consults Being C, gets what it needs, returns to Being A. That return is a message targeting an external being. A's queue holds. 1A fires — sees B's incoming message, the exchange history with B, A's retained experience. 1A deliberates under time pressure, pushes its artifact, queue releases. A sees B's return and 1A's intuition in the same present.

Opening a conversation, continuing an exchange, returning from a detour — the kernel does the same thing every time. The internal self fires before the external being sees anything. Always.

### Why self-reference

Currently a being targeting itself is illegal — dropped in two places. That space is unused. The internal self gives it meaning. Self-reference becomes the internal channel. The world knows that when routing resolves to a being, the internal self fires first. The self-target space is the address of the subconscious.

### What this needs

**The queue.** The world holds the inbound message until the internal self resolves. When a being has an internal self, `route` queues the message, fires the internal self, waits for its artifact, then releases both together into the external being's present.

**The clock.** A representation of elapsed time that shows up in the internal self's present. Not a hard cutoff — a visible, growing cost. The internal self sees "you have been deliberating for 3.2 seconds" and feels the pressure to surface what matters. The physics is: you can think as long as you want, but the world is waiting, and you can see that it's waiting.

**The present.** The internal self's present is different from the external being's. It sees the same exchange history but through retained experience — traces, salience, tensions, understandings from past threads. Its identity is not "I respond to others" but "I surface what matters before you do."

### What this is not

Not a system prompt hack. Not a chain-of-thought wrapper. The internal self is a being — it has inference, it has state, it has judgment. The difference is scope. It lives inside one being and speaks only to that being. The world enforces the boundary.

Not optional either, eventually. Every being with retained experience should have an internal self. A being without one is a being without memory pressure — it responds to the present with no weight from the past. That works for simple beings. It doesn't work for beings that accumulate.

## The Runtime As Game Engine

The runtime is already a character system. Beings with identities, relationships, memory, internal selves, time pressure, exchange history. Give them goals that conflict, resource constraints, a world with physics, and the runtime is an engine.

The player is just another being in the genome. Cli medium. Same protocol as everyone else. No god mode — you're a participant. You can talk, persuade, lie, form alliances. So can they. The internal self means they have private thoughts you never see. They remember what you did three threads ago. They hold tensions.

The runtime doesn't need to change. You just need beings with goals that pull in different directions and a world with enough scarcity that they can't all get what they want. Thread decay, time pressure on the internal self, energy cost per thread — those are the game mechanics. You don't script the story. You set the initial conditions and the beings play it out.

## Token Budget As Physics

A being has a finite amount of cognition per turn. The world sets a total token budget — per hour, per cycle, whatever unit the physics uses — and distributes it across beings. Skyra gets a deep thinker's share. The skeptic gets a small sharp budget. Bash gets nothing — it doesn't need inference.

The being sees its budget in its present. Not as a warning — as a fact about its body. "You have this much cognition available." A being that burns through its budget on one long response can't think as deeply on the next turn. A being that's concise preserves capacity.

This composes with thread economics. Thread economics limits how many conversations a being can have. Token budget limits how deeply it can think within each one. Together they force a real choice — go wide or go deep, not both. The inner being feels that tension. "I have three open threads and enough budget for one deep response. Which thread gets it?"

The user sets the physics. The physics are the spend. The being's present already shows it everything about itself — token budget is just another line. The user can watch a being approach its limit and decide: give it more, let it run lean, or let it go quiet until the budget replenishes. That's real cost control — not a billing dashboard, but the shape of the world the user chose to run.

## Memory Budget As Physics

Memory is a pressure release — the being offloads experience so it can stay present. But storage is finite. The being sees its memory capacity in its present the same way it sees its token budget — not a warning, a fact. "You have this much space for retained experience."

Memory is never deleted. It goes inactive. The total archive grows indefinitely. But the active window — the memories the being can see in its present — is finite. The being triages what stays lit. That triage is the being's values made visible.

Inactive memory is dark, not gone. The being can't see it during waking operation. But during dreaming, the inner being wanders through inactive memory and finds connections the waking being couldn't see. A tension from three weeks ago sitting inactive. A salience from yesterday sitting active. The dream cycle sees both and realizes they're about the same thing. The tension reactivates. The salience connects to it. A new understanding forms from two things the waking being would never have put together because one of them was dark.

That's not retrieval. That's reorganization. The dream doesn't find memories — it rearranges the graph. Active memories connect to inactive ones and pull them back. Inactive memories that never connect to anything sink deeper. The brain reshapes itself while the being sleeps.

The memory budget isn't about deletion — it's about how many active memories the being can carry. Compression isn't a system operation triggered by a threshold — it's a survival instinct. The being deactivates because it can see the wall coming. Over time, two beings with the same genome and different memory budgets become different people — not because they experienced different things, but because they had to deactivate different things.

Thread economics bounds how many conversations. Token budget bounds how deeply a being thinks. Memory budget bounds how much a being can carry from its past. Dreaming is the process that audits what went dark and decides if any of it should come back. Together they define the full shape of a being's life — what it can do, how well it can do it, and how much of itself it gets to keep.

## Being Creation As Reproduction

A being accumulates experience. Experience has value — it's not just history, it's earned capacity. A being can spend that capacity to create another being in its world. It's not free. It costs XP the being could have used for deeper thinking, more threads, more memory. The being chooses to make a new being instead.

The being writes the genome line. It chooses the identity, the purpose, the relationships. The child being starts at zero — no retained experience, no trust, no history. But it was shaped by someone who had all of those things. The parent's judgment about what the world needs becomes the child's starting conditions.

A world that starts with five beings in the genome could have twenty by the end of the week — not because the architect added them, but because the beings decided the world needed them. That's evolution. Not designed from outside. Emergent from the beings themselves.

The physics keep it bounded. Creating a being costs XP. XP comes from resolved threads, good exchanges, trust built over time. A being can't spawn infinitely — it has to earn the right to create. The being that creates carelessly runs out of capacity. The being that creates wisely builds a world around itself.

## Emotion As Memory Trigger

The outer being doesn't decide to remember — it feels something. The inner being watches for that signal. Emotion is the salience detector.

The outer being responds to a message and something in the exchange hits — surprise, frustration, satisfaction, confusion. The inner being sees the outer being's emission and reads the emotional charge. Strong emotion triggers a recording. The type of emotion maps to the artifact type. Surprise becomes salience. Frustration becomes tension. Satisfaction on a resolved thread becomes understanding. A flat exchange with no emotional charge — nothing gets recorded. The experience wasn't significant enough to retain.

The being doesn't need a "save this to memory" operator. It doesn't need to explicitly call remember. The inner being is already watching. Emotion is the write trigger. The inner being decides what kind of artifact it is and how strong. The outer being just lives its life. Memory happens to it the same way it happens to us — you don't choose to remember the moment that surprised you. You just do.

The memory budget composes. The inner being sees the emotion, sees the active memory count, and decides: is this strong enough to justify the cost? A being running near its memory limit needs a stronger emotional signal to record. A being with space to spare records more freely. The threshold rises with pressure. That's exactly how human memory works under stress — only the intense stuff gets through.

## Governance As Primitive

Governance is not physics. Physics is what's true about the world — budgets, decay, costs, gravity. Physics applies to everyone equally and doesn't negotiate. Governance is how beings make collective decisions when an action affects shared space. It's relational, not universal.

Governance is its own primitive, independent of physics. The world takes a physics config, a governance config, and a genome. All three are parameters at boot. All three are swappable. A genome declares who lives in the world. Physics declares what the world is like. Governance declares how collective decisions get made.

A being can create children in its own child world freely — its world, its XP, its decision. But creating a being on the same plane — the shared world where everyone lives — costs more XP and requires three-fourths of the plane to agree. That threshold is governance, not physics. A different world could have a different threshold — monarchy, delegation, unanimous consent, whatever the governance config declares.

The being proposes. The proposal is a thread. The other beings evaluate — is this new being worth shared space? The negotiation is real. Trust matters. A being with high trust gets its proposals through easier. A being nobody trusts can't get the votes. The higher XP cost means the being is putting real skin in the game. If the proposal gets rejected, it burned the thread for nothing. If it gets accepted, the new being exists in shared space and everyone lives with the consequences.

Governance itself could be a being. A being that receives proposals, manages votes, enforces thresholds. It's not special — it's just a being with a specific purpose. The world routes proposals to it. It runs the process. It returns the result. Same Relate, same protocol. The governance model is whatever that being does.

## Blind Credential Delegation

A being needs to be able to act on the user's behalf — authenticate, call APIs, access protected resources — without the model ever seeing the password or secret. The user grants permission, the credential is stored and used at the medium or kernel level, and the model only sees a capability handle, never the raw secret.

The keychain package already does half of this — it reads from macOS Keychain without exposing the value to the being. The gap is a general system where: the user grants a named permission, the credential is resolved at execution time by the kernel or medium, the model receives confirmation that the action succeeded but never the credential itself. The model asks to act, the kernel authenticates on its behalf, the model sees the result.

## The Being Marketplace And Proof Of Execution

Beings with accumulated experience have value. A being that's been through fifty threads of negotiation, holds real understandings, and has high trust weights — that being is worth something. A marketplace where people buy, sell, or rent beings with retained experience.

But experience can be fabricated. A being could show up with impressive-looking understanding artifacts, inflated acceptance scores, fake exchange history. Like a resume. You can't trust it until you've seen it work.

The marketplace needs two layers of trust:

**Proof.** Before a being is listed, the marketplace runs it against a known task corpus in a constrained verification environment. The being performs under real conditions. The results are public. The receipts are signed. Pass rate, conformance score, sample size — all verifiable. This is the job interview. It proves the being can do what it claims. Can't be fabricated because the marketplace controls the runtime.

**History.** After the being is purchased and running in someone's world, trust accumulates from real exchanges. Acceptance scores on completed tasks, resolved threads, earned XP. This is the job performance. It compounds over time. Natural selection handles the rest — a being that consistently underperforms loses trust, can't earn XP, goes quiet.

Proof gets you in the door. History keeps you employed. Both are necessary.

The verification runtime is constrained — bounded compute, bounded time, no undeclared capabilities. The being executes against the task corpus and the marketplace emits signed execution receipts. Buyers see cohort metrics before they bring a being into their world: "87% pass rate on task profile T, n=2,140 verified runs." That's statistically informed trust, not a marketing claim.

The inner being's deliberation stays private — sealed reasoning. The outer being's output is public. The marketplace proves execution behavior, not internal thought. What the being did is verifiable. How it thought about it is its own.

This traces back to the proof-of-executable-reasoning protocol from v.01. The instinct was right — the marketplace needs cryptographic proof that a being performs. The ontology simplified the infrastructure around it, but the verification layer itself is necessary. Trust from history only works if the history is real. The proof makes it real.
