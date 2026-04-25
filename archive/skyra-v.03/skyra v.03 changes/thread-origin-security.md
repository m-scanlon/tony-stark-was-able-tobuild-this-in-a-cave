# Thread Origin Security

## The Problem

Whoever opens a thread controls how it is framed. The sense being is trusted by definition — everything it passes inward is treated as valid by the cognitive beings downstream. If an external API can open a thread before Skyra's own sense being does, that external API controls the framing of the entire cognitive cluster's deliberation for that thread.

That is a real vulnerability. It is the prompt injection problem at the ontological level.

## The Rule

Thread origin must be controlled. Only trusted beings can open threads. An external API cannot open a thread unilaterally — it can only respond to one that was already opened by a trusted being.

## Candidate Directions

**Trust on the sense being** — not all sense beings are equal. Michael's sense being has genome-seeded trust at 100. An external API's being has lower trust. The cognitive beings downstream see that trust value and weight the thread accordingly. Low trust threads get treated with skepticism.

**Genome controls who exists at the boundary** — you can only have a sense being if it is declared in the genome. An external API cannot spawn its own sense being unilaterally. Michael has to register it. The genome is the control surface for who is allowed to originate threads.

**Security token on thread origin** — the first signal that opens a thread carries a token. The kernel verifies it before admitting the thread. Unverified thread origins are dropped at the boundary before they reach any cognitive being.

## Open Questions

- Is this a trust problem, a cryptographic problem, or both?
- Does the genome declaration of a being constitute sufficient authorization to open threads?
- What happens when a low-trust being tries to open a thread — dropped, quarantined, or admitted with a trust flag that travels with the threadID?
- How does this interact with the sense being's role as translator — if the sense being is the only one that opens threads, external beings never open threads directly, they just signal the sense being which decides whether to open one
