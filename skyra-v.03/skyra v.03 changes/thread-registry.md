# Thread Registry

## What A Thread Is

A thread and a threadID are the same thing. One concept. The threadID born at the origin event is the thread. It travels through every exchange it spawns. The runtime does not need to track the thread separately from the threadID — the threadID is the thread.

## The Problem

When a being closes its direct exchange, the thread may still be alive downstream. The being that originated the thread has no visibility into whether its intent is still in flight. The system has no global picture of what threads are active, who owns them, and who is working in them.

## Per-Being Thread Registry

Every being carries its own local thread registry. Three states:

**Root owner** — threads this being originated. Only the root owner can declare a thread fully resolved.

**Active** — the thread this being is currently working in.

**Background** — threads this being is not currently working in but is participating in.

Membership is implicit — a being is a member of any thread that is root, active, or background for it.

## No Global Registry

Thread state is local to each being. There is no global registry. Each being tracks its own root, active, and background threads independently.

Thread traversal is how membership propagates — the threadID travels with the exchange. When a being opens a downstream exchange it passes the threadID along. The receiving being is now in that thread. Membership is implicit in having an exchange carrying that threadID.

## Thread Visibility In The Present

A being's present shows its local thread registry — not the full exchange history of each thread, just awareness:

```
active thread
________________
<threadID> — currently working in this

background threads
________________
<threadID> — not currently working in

root owner of
________________
<threadID> — you originated this, still in flight
```

Context load from thread state is an open problem. Deferred for now.

## Thread Lifetime

A thread is alive as long as any being is active in it. It does not die when the root owner's direct exchange closes. It dies when the root owner explicitly resolves it — declaring the original intent fulfilled.

Only the root owner can resolve a thread. Other beings can close their exchanges and move to background. The thread survives until the root says it is done.

## Open Questions

- How does the root owner know when all downstream exchanges have resolved — does the global registry surface this, or does it rely on the present?
- What happens to background threads that never resolve — do they decay, get explicitly cancelled, or persist indefinitely?
- How does the per-being registry get updated — at the kernel level on every exchange open and close, or declared by the being through the protocol?
- Does a being's present show the content of background threads or just their existence?
- How does thread visibility interact with resolution tracking — are they the same mechanism or layered on top of each other?
