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

## Streaming Prediction — The Token Pipeline

Skyra begins predicting the moment the wake word fires. The ingress shard (Raspberry Pi running a 3B front-door model) starts streaming tokens to the brain shard (Mac mini, control plane) immediately. Each token arrives bundled with its acoustic metadata — VAD signal, speech rate, pitch. There is no separate affect signal step. The emotion and the word arrive together, per token.

**T=0 — Wake word**
No tokens yet. The object store snapshot ref is taken — a git commit hash capturing the exact state of Skyra's memory at this moment. The data is already there, kept warm by the Context Injector. No loading, no inference. Just a reference point. Prediction starts from what is already known.

**T=1..N — Token stream**
Each token arrives with its VAD vector attached. Stream 1 (affect) and Stream 2a (context state) update simultaneously on every token — same event, two recordings. Domain confidence shifts with each token. Retrieval fires when confidence crosses a threshold. The context window assembles speculatively as the utterance arrives.

```
token arrives: {
  token_index
  token            ← word fragment
  affect {
    valence        ← computed from acoustic features at ingress shard
    arousal
    dominance
  }
  speech_rate
  pitch
  pitch_variance
  pause_before
}
```

Example confidence evolution across a single utterance:

```
token 0:  [wake word]                  domain: unknown,  confidence: 0.10
token 1:  "can"    VAD: [−0.4, 0.7]   domain: servers,  confidence: 0.35
token 4:  "you"    VAD: [−0.4, 0.8]   domain: servers,  confidence: 0.52
token 7:  "check"  VAD: [−0.5, 0.7]   domain: servers,  confidence: 0.81
token 9:  "server" VAD: [−0.4, 0.7]   domain: servers,  confidence: 0.99
```

By token 7 the domain is typically known with >0.8 confidence. Retrieval has already fired multiple times. Context window is already warm.

**T=final — Full utterance**
Sentence complete. Prediction confirmed or corrected. Inference fires against an already-warm context window. Stream 3 begins recording system actions. Response latency is low because retrieval happened *during* the utterance, not after.

This is the same principle social media recommendation engines use — predict and pre-load rather than wait for an explicit query. Applied to a personal AI, the prediction target isn't the next video — it's the next thing you need.

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

VAD (Valence-Arousal-Dominance) is the encoding model. Three continuous dimensions that span the full emotional space. Every token carries an emotional coordinate. Over time, sessions and domains accumulate an emotional fingerprint.

---

### Stream 2a — Context State

How Skyra's internal prediction evolved token by token. Updated in lockstep with the affect stream. This is the frame-by-frame reasoning trace — not the LLM inference, but the evolution of the prediction model's state as the utterance arrived.

The critical field is `object_store_snapshot` — a git commit hash referencing the exact state of Skyra's object store at that token. The object store is a git repo (managed via go-git). No data is duplicated here. The hash is a pointer into an already-versioned system. At any token you can reconstruct exactly what Skyra knew.

```
stream_2a_context_state {
  timestamp
  turn_id
  session_id
  token_index

  context_state {
    predicted_domain        // which agent domain — confidence shifting each token
    predicted_intent        // narrowing with each token
    domain_confidence       // 0.0 to 1.0
    active_jobs[]           // job_ids currently executing — affects prediction prior
    retrieval_confidence    // how confident is the speculative retrieval so far
    object_store_snapshot   // git commit hash — full state of Skyra's memory at this token
  }
}
```

Example prediction evolution across a single utterance:

```
token 0:  [wake word]      domain: unknown,  confidence: 0.10
token 1:  [affect: high arousal, negative valence]
                           domain: servers,   confidence: 0.35
token 4:  "can you"        domain: servers,   confidence: 0.52
token 7:  "check the"      domain: servers,   confidence: 0.81
token 9:  "server"         domain: servers,   confidence: 0.99
```

The object store hash at token 0 vs the hash after the turn completes gives you the exact diff — what this interaction changed in Skyra's memory.

---

### Stream 2b — Retrieval Events

Separate from context state because retrievals don't happen at token cadence. They fire when prediction confidence crosses a threshold — event-driven, not periodic. One token may trigger zero retrievals. Another may trigger three.

```
stream_2b_retrieval {
  timestamp
  turn_id
  session_id
  token_index           // which token triggered this retrieval

  trigger               // affect_shift | domain_confidence_threshold | partial_transcript_match
  query {
    semantic_embedding  // vector used for similarity search
    affect_state        // VAD at time of retrieval — used for affect similarity matching
    domain_filter       // which agent's memory space was searched
  }
  results[] {
    ref                 // pointer into object store or long term memory — no data copy
    score               // final retrieval score (see scoring formula below)
    score_breakdown {
      global_importance
      regional_importance
      semantic_similarity
      affect_similarity
    }
  }
  confidence_before     // retrieval_confidence before this fired
  confidence_after      // retrieval_confidence after — delta measures retrieval value
}
```

The confidence delta is the feedback signal for the retrieval strategy. Over time: which triggers actually sharpen the prediction? Which retrievals surface useful context vs noise? This is the data that eventually lets Skyra self-tune her own retrieval thresholds.

---

### Stream 3 — System Output

What the system did. One record per discrete action — not per turn. Order matters. A single turn produces many actions and the causal chain between them is the reasoning trace.

```
stream_3_system_output {
  timestamp
  turn_id
  session_id
  action_index          // ordering within the turn

  action_type           // tool_call | commit_proposed | commit_approved |
                        // commit_denied | job_created | job_completed |
                        // job_cancelled | plan_proposed | plan_approved |
                        // plan_revised | response_emitted

  action_detail {
    tool_name
    tool_args
    tool_result
    commit_ref          // git hash of proposed commit
    job_id
  }

  reasoning_step        // why Skyra took this action at this moment
  outcome               // success | failure | denied | cancelled

  object_store_before   // git hash — state before this action
  object_store_after    // git hash — state after this action
}
```

The before/after git hashes on every action give you a complete causal chain:

```
user said X
→ Skyra knew Y         (context state snapshot — git hash)
→ she did Z            (action)
→ world changed A → B  (git diff between before and after hashes)
```

Full reproducibility. Full accountability.

---

## Three Memory Layers

### Layer 1 — Object Store (Authoritative)

Git repository per agent, managed via go-git. Every state change requires explicit user approval via `propose_commit`. No agent writes to its own object store without a human seeing it first.

This is the source of truth. High trust. Every commit is auditable. Every mistake is reversible via `git checkout`.

Structure:
```
.skyra/agents/{agent_id}/
  state.json            // knowledge, decisions, boundary — committed state
  tools/                // agent's skill implementations
  working/              // gitignored scratch space — free writes during execution
  jobs/{job_id}/        // task artifacts per job
```

### Layer 2 — Observational Store (Raw Signal)

The four streams above. Everything Skyra observed. No user gate. Skyra writes here automatically on every interaction.

Lower initial trust — these are observations, not approved facts. But they accumulate. A single data point is noise. A repeated pattern is signal.

This is where the corpus lives. This is what gets mined for patterns.

### Layer 3 — Long Term Memory (Promoted Synthesis)

Long term memory is not a data store for raw observations. It is a **promotion event** — a synthesis that occurs when a pattern crosses an emotional threshold.

This mirrors human neuroscience directly. The amygdala tags emotionally significant events for deep encoding by the hippocampus. Emotion is the write signal. Routine events decay. Emotionally charged events encode deeply.

Skyra's model:

```
// from observational store
if pattern.frequency > frequency_threshold
AND pattern.affect_magnitude > affect_threshold:
  → synthesize conclusion → write to long term memory

// from object store
if commit.significance > significance_threshold:
  → write to long term memory
```

Two sources, same destination. Once promoted, both carry equal authority.

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
    global:   [long_term, medium_term, session]  // importance scores
    regional: [long_term, medium_term, session]  // domain-scoped importance
  }
  source                // observational | authoritative | both
  confidence            // 0.0 to 1.0 — based on frequency and affect magnitude
  first_seen
  last_seen
  reinforcement_count
}
```

Content is always a synthesized conclusion — not raw data:

> *"When servers go down late at night Mike spirals. Jobs get cancelled. Best resolution: surface the problem clearly and offer one decisive action."*

Not fourteen frustrated events. The meaning. The pattern. What works.

---

## The Retrieval Scoring Formula

Every candidate memory item — whether from the object store, long term memory, or session history — is scored at retrieval time:

```
score = global_importance × regional_importance × semantic_similarity × affect_similarity
```

**global_importance** — scalar from the importance vector `v.global`. Significance across all domains and all time. Major life decisions, defining patterns, cross-domain signals score high.

**regional_importance** — scalar from `v.regional`. Significance within the specific domain being queried. Naturally scopes data without hardcoded domain filters.

**semantic_similarity** — cosine similarity between the query embedding and the item embedding. Standard dense retrieval.

**affect_similarity** — cosine similarity between the current VAD vector and the VAD vector attached to the memory item at write time. This is the key differentiator. Memories encoded during the same emotional state surface preferentially.

Items below a minimum vector threshold are excluded before semantic similarity is even computed. The vector score gates retrieval. Semantic similarity confirms relevance within the gate.

The importance vector has temporal structure:

```
v {
  global:   [long_term, medium_term, session]  // [months/years, days/weeks, right now]
  regional: [long_term, medium_term, session]
  affect:   [valence, arousal, dominance]      // VAD at time of encoding
}
```

Vectors are not static. They decay when items go unaccessed and get bumped when referenced. The decay function is TBD — time-based, relevance-based, or hybrid. Eventually Skyra self-tunes these thresholds based on what surfaces useful context vs noise.

---

## Pattern Recognition Engine

Background process. Always running. Watches all four streams simultaneously. Looks for correlations that repeat across a time window.

Three levels:

**Micro** — within a session. Affect arc, prediction confidence arc, outcome. Real-time signal.

**Meso** — across sessions, days, weeks. Behavioral patterns, temporal correlations, domain usage trends.

**Macro** — across months. Long arc changes. Stress trends. Behavioral shifts. Emerging domains (data clustering around a topic with no existing agent — signal to propose a new agent).

Pattern detection is a cross-stream join on `turn_id`:

```sql
SELECT
  s1.affect,
  s2a.predicted_domain,
  s2a.domain_confidence,
  s2b.retrieval_confidence_delta,
  s3.action_type,
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

Every session adds signal to the observational store. Every pattern strengthens. Every promoted long term memory makes the next retrieval more accurate. The prediction model improves without retraining.

This is the property that makes Skyra fundamentally different from cloud assistants. They optimize for the average user across millions of interactions. Skyra optimizes for one user across a lifetime of interactions. The value compounds with *you* — not with a company's model.

And it runs entirely on your hardware. The corpus never leaves.

---

## Full Loop Summary

```
wake word fires (T=0)
  → object_store_snapshot hash taken — Skyra's memory state captured
  → data already warm, kept current by Context Injector
  → no loading, no inference — just a reference point

tokens stream in, each carrying VAD attached (T=1..N)
  → stream_1 and stream_2a update simultaneously on every token
  → domain_confidence shifts with each token
  → retrieval fires when confidence crosses threshold (stream_2b)
  → context window assembles speculatively as utterance arrives

full utterance lands (T=final)
  → prediction confirmed or corrected
  → inference fires against warm context window
  → stream_3 begins — action_index 0

system executes
  → tool calls, commits, job creation recorded per action
  → object_store_before / object_store_after hashes on each action
  → reasoning_step recorded for each action

turn completes
  → stream_3 closes
  → all four streams written to time series DB
  → pattern recognition engine processes cross-stream correlation
  → if pattern crosses threshold → long term memory written
  → if commit approved → object store updated, importance vector initialized
  → prediction model calibrates for next turn
```

One loop. Every turn. Every interaction makes the next one better.
