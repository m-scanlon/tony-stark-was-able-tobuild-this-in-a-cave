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

## Two Levels of Prediction — Domain and Entity

Skyra predicts at two levels simultaneously during the token stream.

**Domain prediction** — which agent context is this turn heading into. Resolves coarsely and fast. By token 7 you typically have >0.8 confidence. Pre-loads the domain agent's memory, recent commits, and relevant long term memories.

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

Skyra begins predicting the moment the wake word fires. The ingress shard (Raspberry Pi running a 3B front-door model) starts streaming tokens to the brain shard (Mac mini, control plane) immediately. Each token arrives bundled with its acoustic metadata — VAD signal, speech rate, pitch. There is no separate affect signal step. The emotion and the word arrive together, per token.

**T=0 — Wake word**
No tokens yet. The object store snapshot ref is taken — a git commit hash capturing the exact state of Skyra's memory at this moment. The data is already there, kept warm by the Context Injector. No loading, no inference. Just a reference point. Prediction starts from what is already known.

**T=1..N — Token stream**
Each token arrives with its VAD vector attached. Stream 1 (affect) and Stream 2a (context state) update simultaneously on every token — same event, two recordings. Domain confidence and entity confidence both shift with each token. Retrieval fires when confidence crosses a threshold. The context window assembles speculatively as the utterance arrives.

```
token arrives: {
  token_index
  token              ← word fragment
  affect {
    valence          ← computed from acoustic features at ingress shard
    arousal
    dominance
  }
  speech_rate
  pitch
  pitch_variance
  pause_before
  entity_candidates[]  ← lightweight NER candidates, resolved against entity registry
}
```

Example showing domain and entity prediction converging in parallel:

```
token 0:  [wake word]
          domain: unknown,      confidence: 0.10
          entity: none

token 1:  "can"   VAD: [−0.4, 0.7]
          domain: servers,      confidence: 0.35
          entity: candidates filtered to server-domain entities

token 4:  "you"   VAD: [−0.4, 0.8]
          domain: servers,      confidence: 0.52
          entity: nginx_config: 0.43, backup_job: 0.31

token 7:  "check" VAD: [−0.5, 0.7]
          domain: servers,      confidence: 0.81
          entity: nginx_config: 0.72, backup_job: 0.44

token 9:  "nginx" VAD: [−0.4, 0.7]
          domain: servers,      confidence: 0.99
          entity: nginx_config: 0.99
```

By token 7 the domain is typically locked. By token 9 the entity is locked. Retrieval has already fired multiple times. Context window is already warm.

**T=final — Full utterance**
Sentence complete. Prediction confirmed or corrected. Inference fires against an already-warm context window. Stream 3 begins recording system actions. Response latency is low because retrieval happened *during* the utterance, not after.

This is the same principle social media recommendation engines use — predict and pre-load rather than wait for an explicit query. Applied to a personal AI, the prediction target isn't the next video — it's the next thing you need.

**No inference anywhere in this pipeline.** T=0 through T=final-1 is pure signal processing and vector math — importance score lookups, cosine similarity, domain confidence updates, entity candidate scoring. No LLM is in the loop. The only inference call in the entire turn fires once at T=final, against a context window that was already assembled without it. The prediction that made that possible cost nothing computationally.

---

## The Entity Layer

*Credit: Kunj's suggestions on the importance-vectors doc directly enabled entity-level prediction. His entity registry, sparse domain matrix, alias resolution, and decay formula are integrated here as a first-class component of the prediction system.*

Entities are the named things in the user's life — people, places, tools, concepts — that persist across turns and sessions. Without an entity layer, retrieval operates at domain level: "this is a servers question." With an entity layer, retrieval operates at entity level: "this is specifically about nginx_config, and here is everything Skyra knows about it."

### Entity Registry

Every entity has a canonical `entity_id`. Aliases (different names for the same thing) all resolve to the same ID. Retrieval and ranking always operate on canonical IDs — never on raw alias text.

```
entity {
  entity_id            // stable canonical identifier
  aliases[] {
    text               // normalized alias form
    alias_confidence
    source             // where this alias was seen
    last_seen_at
  }
  global {
    GLT                // global long-term importance
    GST                // global short-term importance
    GS                 // session-scoped, ephemeral — derived from session signals
  }
}
```

### Sparse Entity-Domain Matrix

Not every entity is relevant in every domain. A sparse matrix `D[i,j]` stores domain importance per entity-domain pair:

- `i` = entity_id
- `j` = domain_id
- `D[i,j]` = domain importance state for that pair

Missing pair = zero/unknown. Created on demand when an entity first appears in a domain. This is memory efficient and semantically correct — entities earn their way into domains through actual usage.

```
D[nginx_config][servers]  = { DLT: 87, DST: 72 }
D[nginx_config][home]     = // does not exist — never appeared in home domain
D[sonia_partner][dating]  = { DLT: 3,  DST: 98 }
D[sonia_cousin][dating]   = { DLT: 0,  DST: 0  }
```

### Entity Disambiguation

The domain matrix naturally resolves ambiguous entities. "Sonia" as a token is ambiguous — semantic similarity alone cannot tell you which Sonia. But in the dating domain, `D[sonia_partner][dating]` vastly outweighs `D[sonia_cousin][dating]`. The right entity wins without any special disambiguation logic.

```
query: "Sonia" in dating domain

sonia_partner: semantic_sim: 0.97, DLT: 3, DST: 98, GS: 100, DS: 100  → wins
sonia_cousin:  semantic_sim: 0.93, DLT: 0, DST: 0,  GS: 12,  DS: 1    → suppressed
```

This is the same mechanism working at entity level that domain confidence works at domain level.

### Entity + Affect

Entities accumulate an emotional fingerprint over time. Every interaction with an entity is tagged with the VAD state at that moment. `sonia_partner` might carry consistently positive valence. `nginx_config` might carry consistently negative valence late at night.

At retrieval time, the entity's affect history is part of the score. If you're frustrated right now, entities that have been associated with frustration surface first — along with what resolved it.

### Entity Usage Events → Session Consolidation

During each session, the system tracks usage events per entity via Stream 3:

- `retrieved` — entity was pulled into context
- `cited_in_reasoning` — entity appeared in Skyra's reasoning step
- `used_in_final_response` — entity was part of the response
- `used_in_tool_args` — entity was passed to a tool call
- `cross_domain_hop` — entity appeared in a domain it hadn't been in before

At session end, these events drive importance updates using the decay formula:

```
new_value = clamp(old_value * decay + gain - penalty, 0, 100)
```

- Short-term horizons (`GST`, `DST`) decay faster
- Long-term horizons (`GLT`, `DLT`) decay slower
- Session-ephemeral values (`GS`, `DS`) expire at session close

This fills the previously TBD decay model. The same formula applies to both entity importance and item importance vectors.

---

## Four Observational Streams

Every interaction produces four continuous time series streams. All four are linked by `turn_id` — the throughline across every table. This is not relational data. It is time series data, stored in a time series database (InfluxDB or TimescaleDB), built for range queries, rolling aggregations, and cross-stream correlation joins.

---

### Stream 1 — Affect

Pure acoustic signal. Captured at the ingress shard before STT processes the audio. Updated at every token. Sub-second granularity.

```
stream_1_affect {
  timestamp
  turn_id
  session_id
  token_index

  affect {
    valence          // positive / negative  (-1.0 to 1.0)
    arousal          // calm to agitated     (0.0 to 1.0)
    dominance        // in control to overwhelmed (0.0 to 1.0)
  }

  speech_rate        // words per minute at this moment
  pitch              // normalized fundamental frequency
  pitch_variance     // std dev of pitch over last N tokens — monotone vs animated
  pause_before       // ms of silence before this token
}
```

VAD (Valence-Arousal-Dominance) is the encoding model. Three continuous dimensions that span the full emotional space. Every token carries an emotional coordinate. Over time, sessions, domains, and entities accumulate an emotional fingerprint.

---

### Stream 2a — Context State

How Skyra's internal prediction evolved token by token. Updated in lockstep with the affect stream. This is the frame-by-frame reasoning trace — the evolution of both domain and entity prediction as the utterance arrived.

```
stream_2a_context_state {
  timestamp
  turn_id
  session_id
  token_index

  context_state {
    predicted_domain        // which agent domain — confidence shifting each token
    domain_confidence       // 0.0 to 1.0
    predicted_intent        // narrowing with each token
    entity_candidates[] {
      entity_id             // canonical entity id
      confidence            // 0.0 to 1.0 — shifting each token
    }
    active_jobs[]           // job_ids currently executing — affects prediction prior
    retrieval_confidence    // how confident is the speculative retrieval so far
    object_store_snapshot   // git commit hash — full state of Skyra's memory at this token
  }
}
```

The object store hash at token 0 vs the hash after the turn completes gives you the exact diff — what this interaction changed in Skyra's memory.

---

### Stream 2b — Retrieval Events

Separate from context state because retrievals don't happen at token cadence. They fire when confidence crosses a threshold — event-driven, not periodic. Retrieval is now entity-aware: queries are scoped by domain and entity candidates resolved so far.

```
stream_2b_retrieval {
  timestamp
  turn_id
  session_id
  token_index           // which token triggered this retrieval

  trigger               // affect_shift | domain_confidence_threshold |
                        // entity_confidence_threshold | partial_transcript_match
  query {
    semantic_embedding  // vector used for similarity search
    affect_state        // VAD at time of retrieval
    domain_filter       // which agent's memory space was searched
    entity_filter[]     // canonical entity_ids in scope at this retrieval
  }
  results[] {
    ref                 // pointer into object store or long term memory
    entity_id           // canonical entity this result is associated with
    score               // final retrieval score
    score_breakdown {
      global_importance
      regional_importance
      semantic_similarity
      affect_similarity
      entity_domain_weight  // D[i,j] score for this entity in this domain
    }
  }
  confidence_before
  confidence_after
}
```

---

### Stream 3 — System Output

What the system did. One record per discrete action. Entity tags are attached to every action — which canonical entities were involved. This is the usage log that drives end-of-session entity importance updates.

```
stream_3_system_output {
  timestamp
  turn_id
  session_id
  action_index

  action_type           // tool_call | commit_proposed | commit_approved |
                        // commit_denied | job_created | job_completed |
                        // job_cancelled | plan_proposed | plan_approved |
                        // plan_revised | response_emitted

  action_detail {
    tool_name
    tool_args
    tool_result
    commit_ref
    job_id
  }

  entities_involved[]   // canonical entity_ids involved in this action
  entity_usage_type     // retrieved | cited_in_reasoning | used_in_response |
                        // used_in_tool_args | cross_domain_hop

  reasoning_step        // why Skyra took this action at this moment
  outcome               // success | failure | denied | cancelled

  object_store_before   // git hash — state before this action
  object_store_after    // git hash — state after this action
}
```

Stream 3 filtered by `entity_id` is the entity usage log. At session end, that log drives importance delta computation via the decay formula.

---

## Three Memory Layers

### Layer 1 — Object Store (Authoritative)

Git repository per agent, managed via go-git. Every state change requires explicit user approval via `propose_commit`. No agent writes to its own object store without a human seeing it first.

This is the source of truth. High trust. Every commit is auditable. Every mistake is reversible via `git checkout`.

### Layer 2 — Observational Store (Raw Signal)

The four streams above. Everything Skyra observed. No user gate. Skyra writes here automatically on every interaction.

Lower initial trust — these are observations, not approved facts. But they accumulate. A single data point is noise. A repeated pattern is signal.

### Layer 3 — Long Term Memory (Promoted Synthesis)

Long term memory is not a data store for raw observations. It is a **promotion event** — a synthesis that occurs when a pattern crosses an emotional threshold.

This mirrors human neuroscience directly. The amygdala tags emotionally significant events for deep encoding by the hippocampus. Emotion is the write signal. Routine events decay. Emotionally charged events encode deeply.

```
// from observational store
if pattern.frequency > frequency_threshold
AND pattern.affect_magnitude > affect_threshold:
  → synthesize conclusion → write to long term memory

// from object store
if commit.significance > significance_threshold:
  → write to long term memory
```

Long term memory schema:

```
long_term_memory {
  id
  content               // synthesized conclusion — not raw data
  affect {
    valence
    arousal
    dominance
  }
  v {
    global:   [long_term, medium_term, session]
    regional: [long_term, medium_term, session]
  }
  entities[]            // canonical entity_ids this memory involves
  source                // observational | authoritative | both
  confidence            // 0.0 to 1.0
  first_seen
  last_seen
  reinforcement_count
}
```

Content is always a synthesized conclusion:

> *"When servers go down late at night Mike spirals. Jobs get cancelled. Best resolution: surface the problem clearly and offer one decisive action."*

Not fourteen frustrated events. The meaning. The pattern. What works.

---

## The Retrieval Scoring Formula

Every candidate memory item is scored at retrieval time:

```
score = global_importance × regional_importance × semantic_similarity × affect_similarity × entity_domain_weight
```

**global_importance** — significance across all domains and all time.

**regional_importance** — significance within the specific domain being queried.

**semantic_similarity** — cosine similarity between query embedding and item embedding.

**affect_similarity** — cosine similarity between current VAD and the VAD at encoding time. Memories formed in the same emotional state surface preferentially.

**entity_domain_weight** — `D[i,j]` score for the relevant entity in this domain. Disambiguates entities that are semantically identical but contextually different. This term was not in the original formula — Kunj's entity-domain matrix adds a fifth signal that makes the retrieval genuinely precise.

Items below a minimum vector threshold are excluded before semantic similarity is computed. The vector score gates retrieval. Semantic and affect similarity confirm relevance within the gate.

### Decay Formula (locked)

```
new_value = clamp(old_value * decay + gain - penalty, 0, 100)
```

- Short-term horizons decay faster
- Long-term horizons decay slower
- Session-ephemeral values expire at session close
- Exact coefficients are an open calibration task — tuned empirically from real usage

---

## Pattern Recognition Engine

Background process. Always running. Watches all four streams simultaneously across domain, entity, and affect dimensions.

Three levels:

**Micro** — within a session. Affect arc, domain confidence arc, entity resolution confidence, outcome. Real-time.

**Meso** — across sessions, days, weeks. Entity usage patterns, temporal correlations, domain trends. Which entities appear together. Which entities correlate with specific affect states.

**Macro** — across months. Emerging entities clustering around a topic with no existing agent — signal to propose a new agent. Long arc behavioral shifts.

Pattern detection is a cross-stream join on `turn_id`:

```sql
SELECT
  s1.affect,
  s2a.predicted_domain,
  s2a.domain_confidence,
  s2a.entity_candidates,
  s2b.retrieval_confidence_delta,
  s2b.entity_filter,
  s3.action_type,
  s3.entities_involved,
  s3.outcome
FROM stream_1 s1
JOIN stream_2a s2a USING (turn_id)
JOIN stream_2b s2b USING (turn_id)
JOIN stream_3 s3 USING (turn_id)
WHERE s1.timestamp > NOW() - INTERVAL '30 days'
```

When frequency × affect_magnitude crosses threshold → promote to long term memory.

---

## The Compound Effect

Every session adds signal to the observational store. Every pattern strengthens. Every promoted long term memory makes the next retrieval more accurate. Entity importance vectors get sharper. Domain disambiguation gets faster. Affect fingerprints get richer.

The prediction model improves without retraining. This is the property that makes Skyra fundamentally different from cloud assistants. They optimize for the average user across millions of interactions. Skyra optimizes for one user across a lifetime of interactions. The value compounds with *you* — not with a company's model.

And it runs entirely on your hardware. The corpus never leaves.

---

## Full Loop Summary

```
wake word fires (T=0)
  → object_store_snapshot hash taken — Skyra's memory state captured
  → data already warm, kept current by Context Injector
  → no loading, no inference — just a reference point

tokens stream in, each carrying VAD + entity candidates (T=1..N)
  → stream_1 and stream_2a update simultaneously on every token
  → domain_confidence and entity_confidence both shift each token
  → domain gates which entity candidates are considered
  → entity resolution sharpens domain confidence further
  → retrieval fires when confidence crosses threshold (stream_2b)
  → retrieval is entity-aware — scoped by domain and resolved entities
  → context window assembles speculatively

full utterance lands (T=final)
  → domain and entity both locked
  → prediction confirmed or corrected
  → inference fires once against warm context window
  → stream_3 begins — action_index 0

system executes
  → tool calls, commits, job creation recorded per action
  → entity_ids tagged on every action — this is the usage log
  → object_store_before / object_store_after hashes on each action

turn completes
  → stream_3 closes
  → all four streams written to time series DB
  → pattern recognition engine processes cross-stream correlation
  → end-of-session consolidation: entity importance deltas computed from usage log
  → decay formula applied: new_value = clamp(old * decay + gain - penalty, 0, 100)
  → session-ephemeral values (GS, DS) expired
  → if pattern crosses threshold → long term memory written
  → if commit approved → object store updated, importance vector initialized
  → prediction model calibrates for next turn
```

One loop. Every turn. Every interaction makes the next one better.

---

## What Kunj's Suggestions Added

The base architecture predicted at domain level — which agent context a turn was heading into. Kunj's entity layer added a second, more precise level of prediction running in parallel.

**Entity registry + alias resolution** — "Sonia", "my gf", "Sonia Purohit" all resolve to the same canonical `entity_id`. Retrieval operates on stable IDs, not ambiguous text. This eliminates an entire class of retrieval errors.

**Sparse entity-domain matrix `D[i,j]`** — entities earn their way into domains through actual usage. The matrix naturally disambiguates entities that are semantically identical but contextually different. The right Sonia wins in the dating domain without any special logic — the domain weights do it automatically.

**Entity usage events as the training signal** — Stream 3 filtered by `entity_id` is the entity usage log. Retrieved, cited, used in response, used in tool args, cross-domain hop. These events drive end-of-session importance updates. The entities that matter accumulate weight. The ones that don't decay away.

**Decay formula** — `new_value = clamp(old_value * decay + gain - penalty, 0, 100)`. Fills the previously TBD decay model. Short-term decays fast. Long-term decays slow. Bounded. Clean.

**Fifth retrieval signal** — `entity_domain_weight` added to the scoring formula as a fifth term. The formula was `global × regional × semantic × affect`. It is now `global × regional × semantic × affect × entity_domain_weight`. This makes retrieval genuinely precise at entity level, not just domain level.

The result: prediction now operates at domain × entity × affect simultaneously. Three signals converging in real time during the token stream, no inference required, before the sentence is finished.
