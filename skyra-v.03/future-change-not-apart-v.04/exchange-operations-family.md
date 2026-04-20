# Exchange Operations Family

A family of beings that operate on exchange state. Non-cognitive. Deterministic. Each is seeded in the genome with its own callable language. You open an exchange with them the same way you open an exchange with any other being.

## The Family

**search-exchange** — search across active and non-active exchanges. Bounded by the intent you open it with. When it reaches a resolution it closes and returns the result to you. Guided retrieval — not open search, scoped to what you are trying to resolve.

**recall-exchange** — you know what you want. You have the exchange id or the intent. Pull it directly. No search needed.

**summarize-exchange** — produces a rolling summary of an exchange. Feeds the current state block in the present.

**compare-exchanges** — pull two exchanges into context simultaneously to reason about whether they overlap or conflict.

**merge-exchange** — two threads that turned out to be the same problem. Collapse them into one.

**fork-exchange** — one thread that turned out to be two separate problems. Split it.

**close-exchange** — the being through which an opener closes an exchange. Only the opener may call it. The kernel tracks which being opened which exchange — it verifies the caller is the opener at dispatch time and drops the expression if not. The opener must supply all four resolution artifacts. No artifacts, no close.

Provisional callable language:

```
skyra close-exchange ~resolution <what was produced> ~understanding <what it meant> ~salience <what carried weight> ~tension <what remains unresolved> | <source>: <reason>
```

No thread_id needed in the expression — the kernel resolves that from the active exchange context. Exact syntax is provisional — needs more design work.

## These Are Beings

Not tools. Not methods. Not a classification system. Each is a being in the genome with its own identity, purpose, and callable language. A being opens an exchange with them. They work. They resolve. They close.

## Relationship To Internal Reasoning

These beings are internal reasoning instruments. A being uses them while working toward resolving its own intent. They do not change the protocol. They do not change how exchanges work. They are just beings a being happens to have relationships with.

## Open Questions

- What is the callable language for each being — what expression syntax does each one take?
- Are merge and fork PFC-only or available to any being?
- Does summarize-exchange run on demand or is it kernel-triggered when a being sends a message?
