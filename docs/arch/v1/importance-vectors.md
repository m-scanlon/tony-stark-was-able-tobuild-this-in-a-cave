# Importance Vectors

Every piece of data in the Skyra system — turns, sessions, commits, long term memory entries — carries an importance vector attached directly to that piece of data. The vector is how the Context Retriever decides what to surface and what to ignore.

This is a big feature. The full scoring algorithm is a V2/V3 problem — the design here establishes the foundation.

---

## Problems This Solves

**Semantic similarity alone is not enough.** Two pieces of text can be semantically similar without one being relevant to the current moment. Importance vectors add a second independent signal — if something scores high on both relevance and importance, you can trust it. Low importance items don't even make it to the semantic search step.

**Context retrieval has no memory of what matters.** Without vectors, every retrieval pass starts from scratch — recency and text similarity are the only signals. The system has no way to know that a decision made six months ago is still load-bearing, or that a rough patch Mike went through is worth remembering. Vectors encode that knowledge and keep it current.

**Everything buried in the git log.** The object store's commit history is an audit trail, not a memory system. Finding a meaningful past decision means grepping through hundreds of commits. Long term memory with importance vectors makes meaningful data purpose-built for retrieval — structured, scored, instantly accessible.

**No way to scope data to where it matters.** Without regional importance, a piece of data from the servers domain could surface inside a gym session query. Regional vectors naturally scope data to where it's relevant without hardcoded rules or domain filters.

**The system can't detect when your life changes.** A pure retrieval system has no awareness of emerging patterns. Importance vectors accumulating around a topic with no existing agent is the signal that a new domain is forming. The system can surface that observation and propose a new agent — the AI adapts to the user's life rather than waiting to be reconfigured.

**Tuning requires real usage data.** Static retrieval rules can't be calibrated without knowing what actually surfaces useful context. Because vectors are numeric and configurable, thresholds and decay rates can be tuned empirically — and eventually by Skyra herself.

---

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

### ===== SUGGESTIONS BY KUNJ =====
## 8. Proposed Entity-Memory Refinements (Draft)

These refinements preserve the current importance-vector model and add a stricter entity layer for v2+.

### 8.1 Entity Creation Policy

- Not every noun becomes memory automatically.
- Entity extraction should be hybrid:
  - deterministic candidate extraction (rules/NER)
  - LLM confirmation/classification
- Create a memory entity only when the candidate is likely to be reused beyond a one-off turn.
- Track all entities by stable canonical `entity_id` (never by alias text).

### 8.2 Entity-Importance Dimensions

For each memory entity:

- Global:
  - `GLT` (global long-term)
  - `GST` (global short-term)
  - `GS` (global session, ephemeral)
- Domain:
  - `DLT` (domain long-term)
  - `DST` (domain short-term)
  - `DS` (domain session, ephemeral)

`GS` and `DS` are session-scoped and derived from deterministic session signals (for example mention frequency, retrieval usage, final-answer usage). Avoid random initialization.

### 8.3 Sparse Entity-Domain Matrix

Use a sparse matrix `D[i,j]` where:

- `i` = entity id
- `j` = domain id
- `D[i,j]` stores domain importance state for that entity-domain pair

Do not pre-materialize all entity-domain pairs. Missing pair implies zero/unknown and is created only when an entity appears in that domain.

At session start:

- load entity global values (`GLT`, `GST`)
- load active-domain values from `D[i,j]`
- derive ephemeral `GS`, `DS` from session semantic similarity + deterministic session evidence

### 8.4 Session Updates and End-of-Session Consolidation

During session, track measurable usage events per entity:

- retrieved
- cited in reasoning
- used in final response
- used in tool arguments
- cross-domain dependency hop

At session end:

- expire session-only values (`GS`, `DS`) via session table/TTL
- compute update deltas for `GI` (global) and `DI` (domain) from the usage log
- apply deltas to persistent horizons (`GLT`, `GST`, `DLT`, `DST`) with bounded update rules

### 8.5 Aliases and Disambiguation

Each canonical entity may hold multiple aliases (for example two different "Sonia" entities):

- `aliases[]` with normalized form
- `alias_confidence`
- `source` (where alias came from)
- `last_seen_at`

Alias text maps to canonical `entity_id` through resolution logic; retrieval and ranking always operate on canonical ids.

### 8.6 LangMem Integration Boundary

LangMem can be used for:

- candidate entity extraction
- candidate alias extraction
- candidate merge/split suggestions

But persistent writes should pass deterministic validation:

- dedupe checks
- merge/split guardrails
- canonical id enforcement
- bounds checks on importance updates

### 8.7 Equation Placeholder (To Be Locked Later)

For each persistent horizon use a bounded decay-and-gain form:

`new_value = clamp(old_value * decay + gain - penalty, 0, 100)`

- short-term (`GST`, `DST`) should decay faster
- long-term (`GLT`, `DLT`) should decay slower
- exact coefficients remain an open calibration task

### 8.8 Worked Example — "Sonia" Disambiguation in Dating Domain

User query:

`Find a good date spot for Sonia based on her preferences.`

Entity candidates extracted from the turn:

- `date spot`
- `Sonia`
- `preferences`

Creation decisions:

- `date spot` -> no new person/entity record required (handled by domain + place entities)
- `Sonia` -> resolve to existing canonical entity or create one if not found
- `preferences` -> attach to resolved person entity or pull from related memories

Assume two canonical entities already exist:

```json
{
  "entity_id": "sonia_cousin",
  "aliases": ["Sonia", "cousin sonia", "Sonia Joshi", "Sonia from Ahmedabad"],
  "global": { "GLT": 56, "GST": 34 },
  "domain_matrix": {
    "family": { "DLT": 92, "DST": 61 },
    "dating": { "DLT": 0, "DST": 0 }
  }
}
```

```json
{
  "entity_id": "sonia_partner",
  "aliases": ["Sonia", "my gf Sonia", "Sonia Purohit", "Sonia from Worcester"],
  "global": { "GLT": 23, "GST": 98 },
  "domain_matrix": {
    "dating": { "DLT": 3, "DST": 98 },
    "social": { "DLT": 15, "DST": 73 }
  }
}
```

At session start in `dating` domain:

- both entities may have high semantic similarity to the token "Sonia"
- domain-weighted scoring should strongly favor `sonia_partner`
- ephemeral `GS`/`DS` values are derived from deterministic session signals (not random)

Illustrative session-scoped outcomes:

```json
{
  "session_entity_scores": [
    { "entity_id": "sonia_partner", "SSS": 0.97, "GS": 100, "DS": 100 },
    { "entity_id": "sonia_cousin", "SSS": 0.93, "GS": 12, "DS": 1 }
  ]
}
```

Interpretation:

- semantic signal alone is not enough (both look similar)
- `sonia_partner` wins because dating `DLT/DST` and session `DS` are high
- `sonia_cousin` is suppressed in the dating context due to near-zero domain weights

If planning then evaluates candidate date spots (for example `tridents_booksellers_cafe`, `kelly_ice_rink`, `jp_licks`) and selects one:

- selected spot -> domain importance bump (`DI` up)
- repeatedly used winning entity (`sonia_partner`) -> `DI` up, possible `GI` up
- retrieved but unused competitor entities (`sonia_cousin`, rejected spots) -> domain penalty/decay path

This is the intended disambiguation behavior of the entity + importance-vector model.

---

### ==== END OF SUGGESTIONS ====
## 9. Related Docs

- `docs/arch/v1/context-engine.md` — Context Retriever reads importance vectors at query time
- `docs/arch/v1/agents/README.md` — agent state lives in the object store, carries vectors
