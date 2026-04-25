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

## Blind Credential Delegation

A being needs to be able to act on the user's behalf — authenticate, call APIs, access protected resources — without the model ever seeing the password or secret. The user grants permission, the credential is stored and used at the medium or kernel level, and the model only sees a capability handle, never the raw secret.

The keychain package already does half of this — it reads from macOS Keychain without exposing the value to the being. The gap is a general system where: the user grants a named permission, the credential is resolved at execution time by the kernel or medium, the model receives confirmation that the action succeeded but never the credential itself. The model asks to act, the kernel authenticates on its behalf, the model sees the result.
