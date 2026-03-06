# Importance Vectors

Every piece of data in the Skyra system — turns, sessions, commits, long term memory entries — carries an importance vector attached directly to that piece of data. The vector is how the Context Retriever decides what to surface and what to ignore.

This is a big feature. The full scoring algorithm is a V2/V3 problem — the design here establishes the foundation.

---

## 1. The Vector

```
v: {
  global:   [long_term, medium_term, session],
  regional: [long_term, medium_term, session]
}
```

- **global** — importance across the entire system, all domains, all time. Life significance. Major decisions, defining moments, patterns that hold across all domains.
- **regional** — importance within a specific domain or context. High in the relevant domain, near zero in others. Scopes data naturally without hardcoded rules.

Each dimension has its own time horizon profile:
- **long_term** — importance over months and years
- **medium_term** — importance over days and weeks
- **session** — importance right now, in this conversation

Scale: 0–100. Arbitrary — designed to be tuned empirically.

---

## 2. Examples

```
turn {
  data: "what did I decide about backups"
  v: {
    global:   [20, 30, 85],  // not historically significant, somewhat recent, very live right now
    regional: [10, 40, 90]   // low cross-domain, medium recent, high in current domain
  }
}

commit {
  data: "decided to move everything to S3"
  v: {
    global:   [90, 60, 30],  // major decision, fading medium term, not session relevant
    regional: [80, 70, 40]   // high in servers domain across time
  }
}

long_term_memory {
  data: "Mike had a rough February and pushed through it"
  v: {
    global:   [75, 30, 5],   // life significant, fading, not session relevant
    regional: [20, 10, 5]    // not domain specific
  }
}

session {
  data: "Feb 2026 debugging session"
  v: {
    global:   [70, 40, 5],
    regional: [60, 50, 10]
  }
}
```

---

## 3. Retrieval Score

At query time the Context Retriever scores every candidate item:

```
score = global * regional * semantic_similarity
```

- Items must meet a minimum vector score threshold before semantic similarity is even computed. Below threshold — not considered.
- Highest scoring items that fit the token budget surface.
- The retrieval engine doesn't care if an item is a turn, a commit, or a long term memory entry — everything is a scored item in the same index.

Semantic similarity alone is not enough. The vector score gates retrieval. Semantic similarity confirms relevance within the gate.

Threshold values are configurable and will be tuned empirically.

---

## 4. How Vectors Are Updated

Vectors are not static. They evolve based on:

- **Mention** — when the context blob arriving with an event references a piece of data, that item's relevant dimension scores go up
- **Access** — items that keep getting surfaced retain their scores
- **Time passing** — unmentioned, unaccessed items decay down
- **Decay** — strategy TBD (time-based, relevance-based, or both)

The context blob attached to each incoming event (pushed by CIX, attached at hydration) is the primary relevance signal. When the domain agent receives an event, it reads the context blob and applies it to the object store — bumping vectors on related items, letting everything else decay.

---

## 5. Long Term Memory Store

A dedicated store for significant moments, patterns, and emotionally weighted data. Second class citizens — below the object store in authority, above session history in permanence.

Purpose-built for retrieval — not buried in git logs:

```
// instead of:
git log → 847 commits → grep → maybe find it

// it's:
long_term_memory.query("rough patch February") → v:[75, 20] → surfaces instantly
```

- **Object store** — authoritative state, source of truth
- **Git log** — audit trail
- **Long term memory** — the meaning layer

---

## 6. Versioned Rollout

**V1** — static vectors. Manually assigned or simple rules. Gets the plumbing in place.

**V2** — vectors updated based on access frequency, recency, and context blob signals. Decay runs automatically.

**V3** — full background importance process. Continuously measures importance across agents and across time. Uses real access patterns to keep vectors current.

Because V3 runs across agents it can spot cross-domain patterns. Something that scores low regionally in every individual domain but keeps appearing everywhere accumulates global importance. This is also how the system detects new domains emerging in the user's life — data clustering around a topic with no existing agent, growing global relevance, low regional fit anywhere:

> "Looks like photography is becoming a real part of your life — want me to create an agent for it?"

**V3+ — self-tuning** — Skyra adjusts her own thresholds and decay rates based on what surfaces useful context vs noise. The system becomes self-calibrating.

---

## 7. Open Questions

- What is the exact decay function — linear, exponential, step?
- What signals drive the initial vector assignment at creation time?
- How does the scoring algorithm assign vectors to new data that has no history yet?
- What is the minimum threshold score before semantic similarity is computed?
- How does the V3 background process handle conflicting cross-domain signals?
- What model (if any) drives inference in the background loop — lightweight LLM or pure signal/rule based?

---

## 8. Related Docs

- `docs/arch/v1/context-engine.md` — Context Retriever reads importance vectors at query time
- `docs/arch/v1/agents/README.md` — agent state lives in the object store, carries vectors
