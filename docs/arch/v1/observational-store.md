# Observational Store

## What It Is

The observational store is Skyra's continuous record of everything she watches. It is not the object store — those are user-approved commits. It is not long term memory — that is the promoted, synthesized output. It is the raw signal layer. The stream of everything that happened, how the user felt, how Skyra's prediction evolved, what was retrieved, and what the system did about it.

Every interaction produces records across four streams. All four are linked by `turn_id`. That is the throughline through everything.

---

## Four Streams

```
Stream 1  —  affect                per token       turn_id
Stream 2a —  context state         per token       turn_id
Stream 2b —  retrieval events      per retrieval   turn_id
Stream 3  —  system output         per action      turn_id
```

---

### Stream 1 — Affect

Pure emotional signal. Updated at every token. Continuous, sub-second granularity. Captured acoustically at the ingress shard before STT strips the audio. Nothing semantic — just how the user was feeling as they spoke.

```
stream_1_affect {
  timestamp
  turn_id
  session_id
  token_index

  affect {
    valence          ← positive / negative (-1.0 to 1.0)
    arousal          ← calm to excited (0.0 to 1.0)
    dominance        ← in control to overwhelmed (0.0 to 1.0)
  }

  speech_rate        ← words per minute at this moment
  pitch              ← normalized pitch value
  pitch_variance     ← monotone vs animated
}
```

---

### Stream 2a — Context State

How Skyra's prediction evolved token by token. Updated at every token alongside the affect stream. This is her reasoning — not the inference call, but the frame-by-frame evolution of her understanding as the utterance arrived.

Includes a snapshot ref to Skyra's object store at that exact moment. The ref is a hash — the full object store state is already versioned in git. No data duplication. The hash is the pointer into that state.

```
stream_2a_context_state {
  timestamp
  turn_id
  session_id
  token_index

  context_state {
    predicted_domain        ← confidence shifting with each token
    predicted_intent        ← narrowing with each token
    domain_confidence       ← 0.0 to 1.0
    active_jobs[]           ← what jobs are currently running
    retrieval_confidence    ← how confident is the prediction
    object_store_snapshot   ← hash ref to Skyra's object store at this token
  }
}
```

The object store snapshot ref means you have the complete state of Skyra's memory at every token. After the turn completes and commits land, the diff between the first and last snapshot ref is the measure of what this turn changed.

Example of how the prediction evolves:

```
token 0:  [wake word]      → domain: unknown,  confidence: 0.1, pre-loading temporal patterns
token 1:  [affect arrives] → domain: servers,   confidence: 0.3, pre-loading server context
token 4:  "can you"        → domain: servers,   confidence: 0.5, narrowing
token 7:  "check the"      → domain: servers,   confidence: 0.8, retrieval sharpening
token 9:  "server"         → domain: servers,   confidence: 0.99, context warm
```

---

### Stream 2b — Retrieval Events

Every retrieval that fired during the turn. Separate from context state because retrievals do not happen at token cadence — they fire when prediction confidence crosses a threshold. One token may trigger zero retrievals. Another may trigger three.

```
stream_2b_retrieval {
  timestamp
  turn_id
  session_id
  token_index           ← which token triggered this retrieval

  trigger               ← affect_shift | domain_confidence | partial_transcript
  query                 ← what was queried
  results[]             ← refs into object store / long term memory (pointers, not copies)
  confidence_before     ← retrieval_confidence before this fired
  confidence_after      ← retrieval_confidence after
}
```

Results are refs — pointers into the object store snapshot or long term memory. The data lives in its authoritative location. This stream records what was selected and why.

The confidence delta tells you how much each retrieval sharpened the prediction. Over time this calibrates which retrieval triggers are actually useful vs noise.

---

### Stream 3 — System Output

What the system did. One record per discrete action. Not per turn — per action. A single turn can produce many actions and the order matters.

Includes the complete reasoning trace — every step from retrieval through inference through tool calls through result. This is how Skyra got to the answer, preserved.

```
stream_3_system_output {
  timestamp
  turn_id
  session_id
  action_index          ← ordering within the turn

  action_type           ← tool_call | commit_proposed | commit_approved |
                          commit_denied | job_created | job_completed |
                          job_cancelled | plan_proposed | plan_approved |
                          plan_revised | response_emitted

  action_detail {
    tool_name           ← if tool_call
    tool_args           ← if tool_call
    tool_result         ← if tool_call
    commit_ref          ← if commit action
    job_id              ← if job action
  }

  reasoning_step        ← why Skyra took this action at this moment
  outcome               ← success | failure | denied | cancelled

  object_store_before   ← hash ref to object store before this action
  object_store_after    ← hash ref to object store after this action
}
```

The object store before/after refs on each action give you a causal chain:

```
user said X
→ Skyra knew Y  (context state snapshot)
→ she did Z     (action)
→ world changed from A to B  (object store diff)
```

Full accountability. Full reproducibility.

---

## The Join — Where Pattern Recognition Lives

`turn_id` links all four streams. The pattern recognition engine joins across them.

Neither stream is sufficient alone. The intelligence is in the correlation:

```
stream_1:   frustrated + high arousal
stream_2a:  servers domain, high confidence by token 7
stream_2b:  long term memory retrieved — "server outage pattern"
stream_3:   one decisive tool call + commit approved + response emitted
→ outcome: user active after, job completed
→ pattern: direct action resolves frustration. this works.
```

vs

```
stream_1:   frustrated + high arousal
stream_2a:  servers domain, low confidence — prediction never sharpened
stream_2b:  multiple retrievals, confidence_after < confidence_before
stream_3:   job cancelled + commit denied
→ outcome: user silent after
→ pattern: when prediction fails to sharpen, system fails to help.
           retrieval strategy needs adjustment for this affect state.
```

Skyra learns not just what you felt — but how well she predicted what you needed — and whether what she did about it worked.

---

## Long Term Memory — The Promotion Event

Long term memory is not a fourth stream. It is a promotion event — a synthesis.

A long term memory is created when a pattern crosses an emotional threshold. This mirrors human long term memory: emotion is the write signal. The amygdala fires during emotionally significant events and signals deep encoding.

```
if pattern.frequency > frequency_threshold
AND pattern.affect_magnitude > affect_threshold:
  → promote to long term memory
```

Two sources feed promotion:

**From the observational store** — a recurring pattern across all four streams. Skyra synthesizes a conclusion, not raw data.

**From the object store** — a user-approved commit significant enough to warrant deep encoding. High initial importance vector.

Once promoted, both sources carry equal authority in long term memory. The distinction disappears.

```
long_term_memory {
  id
  content               ← synthesized conclusion, not raw data
  affect {
    valence
    arousal
    dominance
  }
  v {
    global:   [long_term, medium_term, session]
    regional: [long_term, medium_term, session]
  }
  source                ← observational | authoritative | both
  confidence            ← 0.0 to 1.0
  first_seen
  last_seen
  reinforcement_count
}
```

Content is always a synthesis:

```
"When servers go down late at night Mike spirals. Jobs get cancelled.
Best resolution: surface the problem clearly and offer one decisive action."
```

Not fourteen data points. The meaning. The pattern. The resolution.

---

## Predictive Retrieval

The four streams enable prediction — not reactive retrieval. Skyra assembles context before the turn completes.

```
T=0  wake word fires
     → temporal patterns load (time of day, day of week, last session)
     → stream_2a begins, object_store_snapshot ref taken

T=1  affect signal arrives
     → stream_1 begins
     → emotional patterns query long term memory
     → stream_2b: first retrieval fires

T=2  partial transcript arrives
     → domain confidence starts rising
     → stream_2b: retrieval sharpens to predicted domain
     → context window warming

T=3  full utterance lands
     → prediction confirmed or corrected
     → inference fires against warm context
     → stream_3 begins recording actions
```

By T=3 the context window is already assembled. Response latency is low because the work happened during the utterance, not after it.

---

## Storage

All four streams are time series data. Not a fit for SQLite or the object store.

Storage: **time series database** (exact engine TBD — InfluxDB, TimescaleDB, or equivalent).

Queries the pattern recognition engine needs:
- Range queries across any stream within a time window
- Rolling aggregates: average arousal over last 7 sessions
- Correlation joins across all four streams on `turn_id`
- Pattern detection: recurring affect + prediction + outcome combinations
- Confidence delta analysis: which retrieval triggers sharpen predictions

---

## What This Is Not

- **Not the object store** — no user approval, no commits, no source of truth for agent state
- **Not long term memory** — raw signal, not synthesized conclusions
- **Not session history** — session history is turn text for conversational continuity. This is signal for pattern recognition and prediction calibration.
- **Not emotional data in the object store** — the user's emotional history does not pollute authoritative commits. It lives here.

---

## Related Docs

- `docs/arch/v1/importance-vectors.md` — vector model, affect dimension, retrieval scoring
- `docs/arch/v1/context-engine.md` — context package assembly, retrieval pipeline
- `docs/arch/v1/job-types.md` — context background job watches all four streams
- `skyra/internal/agent/README.md` — object store, authoritative memory, commit model
