# Being Types

Each being has a type. The type is a data structure at the top of the being that determines what it can do — what it can receive, what it can emit, how it participates in the runtime.

Right now the only distinction is the cognitive flag — true or false. That's too coarse.

A type system on beings would let the runtime enforce capabilities and constraints at the structural level rather than relying on the being's own behavior or inference to stay in bounds.

Types might include things like:
- external (outside the boundary, owns an I/O pair)
- cognitive (reasons through inference, full present)
- non-cognitive (transducer, minimal present)
- boundary (I/O, owns the external exchange)
- memory (companion being, pressure release)

This extends to the channel. Right now every cognitive relationship is an ExchangeStack — same structure regardless of what the relationship is for. But the channel between prefrontal and values is doing different work than the channel between sensory and thalamus. If beings have types, the channel between two beings of certain types probably needs a data structure matched to what that relationship actually does.

The type lives at the top of the being. The channel type follows from the being types on both ends.

Rough idea. The shape of what "type" means and what data structure it maps to is unresolved.

## Response Queue

When a being routes to someone different from the sender, and that someone is cognitive, the channel type for that relationship might include a queue. The being holds until all — or a threshold — of the expected responses come back before passing the exchange on. Flows naturally from the type: the data structure on the channel owns the queue, not the being itself.

## Communication Patterns

Being types probably determine what communication patterns a being can participate in. The current sequential hop model — signal moves one being to the next — is only one pattern.

Others:

**Fan out / collect / synthesize** — a being broadcasts to multiple peers, waits for all to return, retreats to its self-exchange to synthesize, then resolves outward. Example: prefrontal asks strategy, values, and consequence simultaneously, waits for all three, thinks by itself, then decides.

**Sequential chain** — ask A, take the answer to B, take that to C. Each hop depends on the previous.

**Single consult** — ask one peer, act immediately on the response.

The pattern a being uses probably follows from its type and the channel types on its relationships. The channel owns the coordination structure — queue, threshold, ordering — not the being itself. The being just fires and synthesizes. The infrastructure handles the rest.

What patterns exist and how they map to channel data structures is unresolved.

## Communication Pattern Beings

The pattern itself might be a being type rather than a channel property.

A fan out being receives a problem, calls a set of peers, tracks who owes what back, collects responses, synthesizes, and returns a result. The calling being doesn't manage any of that — it just knows to use the fan out being when it needs multiple opinions.

This makes patterns reusable across any being that needs them. It also makes patterns learnable — the fan out being accumulates retained experience about which peer combinations produce good synthesis, when three opinions is overkill, when one is enough.

Fits the ontology: a coordination pattern communicates, receives, manages, returns. So it's a being.

Types might include:
- fan out (broadcast, collect, synthesize)
- sequential chain (each hop depends on the previous)
- single consult (ask one, act immediately)

The channel between the calling being and the pattern being is probably its own type too — distinct from the channels between the pattern being and its targets.

Unresolved.
