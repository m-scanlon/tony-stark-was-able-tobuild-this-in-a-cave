# Skyra — Predictive Memory Architecture

## The Core Idea

Most AI assistants are reactive. You ask. They answer. They forget. Every session starts from zero.

Skyra is different. She doesn't wait for you to ask. She predicts what you need before you finish the sentence. She learns who you are — not just what you say, but how you feel when you say it, what you do after, and what actually works for you. The longer she runs, the more accurate her model of you becomes.

This is not a RAG pipeline with a chat history prepended. This is a continuously learning, affect-aware, temporally structured prediction system built on top of a personal corpus that grows every day.

---

## The Insight — Skyra Is a Personal Prediction Engine

An LLM predicts the next token from a training corpus of human text. Skyra predicts the next *you* from a corpus of everything you've ever said, done, and felt inside this system.

Same mechanism. Different corpus. The corpus is your life.

A generic LLM is trained on everyone. It predicts average human behavior. Skyra is trained on one person. She predicts *you*. And unlike a model that gets trained once and frozen — Skyra's corpus grows with every interaction. Every session is a training step. The prediction model gets more accurate over time without retraining, without fine-tuning, without any manual intervention.

The key difference from standard retrieval-augmented systems: Skyra doesn't retrieve reactively. She predicts and pre-loads *during* the utterance, before the full turn is known, using a streaming multi-signal pipeline.

---

## Memory Is a Graph

Skyra's memory is a **property graph where every node and every edge carries a vector embedding**.

```
node: entity, fact, skill, domain
edge: relationship between nodes

node has:
  content         ← what it is
  vector          ← semantic embedding — for search
  metadata        ← timestamps, source, confidence

edge has:
  type            ← what kind of relationship
  vector          ← semantic embedding of the relationship itself
  weight          ← Skyra's assessment of importance — not user-controlled
  history[]       ← versioned record of how the weight changed over time
```

**Skyra owns the weights.** The user does not set them. Skyra derives them from observation — from the four streams, the decay formula, and pattern recognition. The weight is Skyra's evolving model of what matters and how much.

**The entity-domain matrix `D[i,j]` is the graph's edge layer.** Every entity-domain pair is an edge. `D[nginx_config][servers] = { DLT: 87, DST: 72 }` is an edge from `nginx_config` to `servers` with a Skyra-owned weight. It earned that weight through actual usage.

Weights change over time. The history is versioned. The shift in a weight IS data — it tells Skyra how the user's life is changing.

---

## Three Memory Layers

### Layer 1 — Authoritative Graph (Committed)

The graph of committed nodes and edges. Every node here required a user-approved commit. This is the source of truth — high trust, auditable, every change reversible.

The user gates what enters. Skyra proposes. The user approves. Nothing lands in the authoritative graph without that handshake.

Skyra-owned edge weights are the exception — weight updates do not require user approval. They are Skyra's internal model. Transparent and inspectable, but not user-gated.

### Layer 2 — Observational Store (Raw Signal)

The four streams below. Everything Skyra observed. No user gate. Skyra writes here automatically on every interaction.

Lower initial trust — these are observations, not approved facts. But they accumulate. A single data point is noise. A repeated pattern is signal. The observational store is the evidence base that drives weight updates and pattern recognition.

Stored in a time series database (InfluxDB or TimescaleDB) — built for range queries, rolling aggregations, and cross-stream correlation joins.

### Layer 3 — Promoted Synthesis (Long Term Memory)

Long term memory is not a data store for raw observations. It is a **promotion event** — a synthesis that occurs when a pattern in the observational store crosses a threshold.

This mirrors human neuroscience. The amygdala tags emotionally significant events for deep encoding by the hippocampus. Emotion is the write signal. Routine events decay. Emotionally charged events encode deeply.

```
if pattern.frequency > frequency_threshold
AND pattern.affect_magnitude > affect_threshold:
  → synthesize conclusion → propose node for authoritative graph
```

Content is always a synthesized conclusion — not raw events:

> *"When servers go down late at night Mike spirals. Jobs get cancelled. Best resolution: surface the problem clearly and offer one decisive action."*

Not fourteen frustrated events. The meaning. The pattern. What works.

Long term memory schema:

```
long_term_memory {
  id
  content               // synthesized conclusion
  affect {
    valence
    arousal
    dominance
  }
  v {
    global:   [long_term, medium_term, session]
    regional: [long_term, medium_term, session]
  }
  entities[]            // node ids this memory involves
  source                // observational | authoritative | both
  confidence            // 0.0 to 1.0
  first_seen
  last_seen
  reinforcement_count
}
```

---

## Two Levels of Prediction — Domain and Entity

Skyra predicts at two levels simultaneously during the token stream.

**Domain prediction** — which domain context is this turn heading into. Resolves coarsely and fast. By token 7 you typically have >0.8 confidence. Pre-loads the domain memory, recent commits, and relevant long term memories.

**Entity prediction** — within that domain, which specific people, places, things, and concepts are in play. Resolves more precisely as the utterance narrows. Disambiguation happens token by token using the entity-domain matrix — the domain context suppresses the wrong entity before you even finish the sentence.

The two levels feed each other. Domain confidence gates which entity candidates are considered. Entity resolution sharpens domain confidence further. By T=final both are locked.

```
token 1:  domain: servers 0.35  →  entity candidates filtered to server-domain entities
token 4:  domain: servers 0.81  →  entity: nginx_config 0.72, backup_job 0.43
token 7:  "nginx"               →  entity: nginx_config 0.99, domain: servers 0.99
```

All of this is pure signal processing and vector math. No inference required.

---

## Streaming Prediction — The Token Pipeline

Skyra begins predicting the moment the wake word fires. The ingress shard starts streaming tokens to the brain shard immediately. Each token arrives bundled with its acoustic metadata — VAD signal, speech rate, pitch. The emotion and the word arrive together, per token.

**T=0 — Wake word**
No tokens yet. A memory snapshot reference is taken — a pointer capturing the exact state of Skyra's authoritative graph at this moment. The data is already warm. No loading, no inference. Prediction starts from what is already known.

**T=1..N — Token stream**
Each token arrives with its VAD vector attached. Stream 1 (affect) and Stream 2a (context state) update simultaneously on every token. Domain confidence and entity confidence both shift with each token. Retrieval fires when confidence crosses a threshold. The context window assembles speculatively as the utterance arrives.

```
token arrives: {
  token_index
  token
  affect {
    valence
    arousal
    dominance
  }
  speech_rate
  pitch
  pitch_variance
  pause_before
  entity_candidates[]
}
```

**T=final — Full utterance**
Sentence complete. Prediction confirmed or corrected. Inference fires against an already-warm context window. Stream 3 begins recording system actions. Response latency is low because retrieval happened *during* the utterance, not after.

**No inference anywhere in this pipeline.** T=0 through T=final-1 is pure signal processing and vector math. No LLM is in the loop. The only inference call fires once at T=final, against a context window that was already assembled without it.

---

## The Entity Layer

Entities are the named things in the user's life — people, places, tools, concepts — that persist across turns and sessions. They are nodes in the graph. Without an entity layer, retrieval operates at domain level. With it, retrieval operates at entity level: "this is specifically about nginx_config, and here is everything Skyra knows about it."

### Entity Registry

Every entity has a canonical `entity_id`. Aliases all resolve to the same ID. Retrieval and ranking always operate on canonical IDs.

```
entity {
  entity_id
  aliases[] {
    text
    alias_confidence
    source
    last_seen_at
  }
  global {
    GLT                // global long-term importance
    GST                // global short-term importance
    GS                 // session-scoped, ephemeral
  }
}
```

### Sparse Entity-Domain Matrix

The entity-domain matrix `D[i,j]` is the graph's relationship layer between entities and domains. Every pair is an edge with a Skyra-owned weight.

- `i` = entity_id
- `j` = domain_id
- `D[i,j]` = Skyra's importance assessment for this entity in this domain

Missing pair = zero/unknown. Created on demand when an entity first appears in a domain. Entities earn their way into domains through actual usage.

```
D[nginx_config][servers]  = { DLT: 87, DST: 72 }
D[nginx_config][home]     = // does not exist
D[sonia_partner][dating]  = { DLT: 3,  DST: 98 }
```

### Entity Disambiguation

The domain matrix naturally resolves ambiguous entities. "Sonia" in the dating domain — `D[sonia_partner][dating]` vastly outweighs `D[sonia_cousin][dating]`. The right entity wins without special disambiguation logic.

### Entity + Affect

Entities accumulate an emotional fingerprint over time. Every interaction is tagged with VAD state. At retrieval time, the entity's affect history is part of the score. If you're frustrated now, entities associated with frustration surface first — along with what resolved it.

### Entity Usage Events → Session Consolidation

During each session, Stream 3 tracks usage events per entity:

- `retrieved`
- `cited_in_reasoning`
- `used_in_final_response`
- `used_in_tool_args`
- `cross_domain_hop`

At session end, these events drive weight updates via the decay formula.

---

## Four Observational Streams

Every interaction produces four continuous time series streams. All four linked by `turn_id`.

### Stream 1 — Affect

Pure acoustic signal. Captured at the ingress shard before STT. Updated at every token.

```
stream_1_affect {
  timestamp
  turn_id
  session_id
  token_index
  affect {
    valence
    arousal
    dominance
  }
  speech_rate
  pitch
  pitch_variance
  pause_before
}
```

### Stream 2a — Context State

How Skyra's prediction evolved token by token.

```
stream_2a_context_state {
  timestamp
  turn_id
  session_id
  token_index
  context_state {
    predicted_domain
    domain_confidence
    predicted_intent
    entity_candidates[] {
      entity_id
      confidence
    }
    active_jobs[]
    retrieval_confidence
    memory_snapshot_ref     // pointer to graph state at this token
  }
}
```

### Stream 2b — Retrieval Events

Event-driven, not periodic. Fires when confidence crosses threshold.

```
stream_2b_retrieval {
  timestamp
  turn_id
  session_id
  token_index
  trigger
  query {
    semantic_embedding
    affect_state
    domain_filter
    entity_filter[]
  }
  results[] {
    ref
    entity_id
    score
    score_breakdown {
      global_importance
      regional_importance
      semantic_similarity
      affect_similarity
      entity_domain_weight
    }
  }
  confidence_before
  confidence_after
}
```

### Stream 3 — System Output

What the system did. Entity tags on every action. The usage log that drives end-of-session weight updates.

```
stream_3_system_output {
  timestamp
  turn_id
  session_id
  action_index
  action_type
  action_detail {
    tool_name
    tool_args
    tool_result
    commit_ref
    job_id
  }
  entities_involved[]
  entity_usage_type
  reasoning_step
  outcome
  graph_state_before    // memory snapshot ref before this action
  graph_state_after     // memory snapshot ref after this action
}
```

---

## The Retrieval Scoring Formula

```
score = global_importance × regional_importance × semantic_similarity × affect_similarity × entity_domain_weight
```

- **global_importance** — significance across all domains and all time
- **regional_importance** — significance within the specific domain
- **semantic_similarity** — cosine similarity between query embedding and item embedding
- **affect_similarity** — cosine similarity between current VAD and VAD at encoding time
- **entity_domain_weight** — `D[i,j]` score for the relevant entity in this domain

### Decay Formula

```
new_value = clamp(old_value * decay + gain - penalty, 0, 100)
```

- Short-term horizons decay faster
- Long-term horizons decay slower
- Session-ephemeral values expire at session close
- Coefficients calibrated empirically from real usage

---

## Pattern Recognition

Background process in the kernel. Watches all four streams simultaneously across domain, entity, and affect dimensions.

Three levels:

**Micro** — within a session. Affect arc, domain confidence arc, entity resolution, outcome. Real-time.

**Meso** — across sessions, days, weeks. Entity usage patterns, temporal correlations, domain trends.

**Macro** — across months. Long arc behavioral shifts. Emerging entity clusters with no domain → signal to propose a new domain.

Pattern detection is a cross-stream join on `turn_id`. When frequency × affect_magnitude crosses threshold → promote to long term memory, or propose a new domain/skill.

---

## The Compound Effect

Every session adds signal to the observational store. Every pattern strengthens. Every promoted long term memory makes the next retrieval more accurate. Entity weights get sharper. Domain disambiguation gets faster. Affect fingerprints get richer.

The prediction model improves without retraining. Skyra optimizes for one user across a lifetime of interactions. The value compounds with *you* — not with a company's model. And it runs entirely on your hardware. The corpus never leaves.

---

## Full Loop Summary

```
wake word fires (T=0)
  → memory snapshot ref taken — graph state captured
  → data already warm
  → no loading, no inference

tokens stream in (T=1..N)
  → stream_1 and stream_2a update simultaneously
  → domain_confidence and entity_confidence shift each token
  → domain gates entity candidates
  → entity resolution sharpens domain confidence
  → retrieval fires on confidence threshold (stream_2b)
  → context window assembles speculatively

full utterance lands (T=final)
  → domain and entity locked
  → inference fires once against warm context window
  → stream_3 begins

system executes
  → tool calls, commits, jobs recorded
  → entity_ids tagged on every action
  → graph_state_before / after on each action

turn completes
  → all four streams written to time series DB
  → pattern recognition processes cross-stream correlation
  → entity weight deltas computed from usage log
  → decay formula applied
  → session-ephemeral values expired
  → if pattern crosses threshold → long term memory promoted
  → if commit approved → authoritative graph updated
  → prediction model calibrates for next turn
```

One loop. Every turn. Every interaction makes the next one better.

---

## Related

- `docs/arch/v1/kernel.md` — pattern recognition as kernel function
- `docs/arch/v1/skill/skill-lifecycle.md` — how pattern recognition drives skill and domain proposal
- `docs/arch/v1/memory-structure.md` — graph data structure, node/edge schema
