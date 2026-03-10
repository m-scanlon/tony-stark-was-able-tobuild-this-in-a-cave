# Testing Strategies

---

## Synthetic User Personas

The primary testing strategy for the memory pipeline. A synthetic persona is an LLM configured to behave like a real user — consistent habits, vocabulary, domains, emotional patterns. It chats with Skyra, generating realistic session data. The output of the reasoning and integrate passes can then be audited against what the graph *should* look like.

### Why This Works

- You control the persona — you know the ground truth
- You can verify the graph output against expectations
- You can run it repeatedly with the same persona to test consistency
- You can run it with different personas to test variety

### Persona Design

Each persona needs:

- **Identity** — name, age, rough life context
- **Domains** — what areas of life they talk about (work, fitness, family, etc.)
- **Habits** — recurring patterns (always works late, logs workouts on Mondays, etc.)
- **Vocabulary** — how they refer to things. Do they say "Mike" or "Michael". "gym" or "the gym". This tests alias resolution in integrate.
- **Affect profile** — generally calm, anxious, frustrated easily, high energy. Tests VAD signal quality.
- **Quirks** — things that make them hard to parse. Vague references, topic switching, shorthand. Tests reasoning skill robustness.

### Example Personas

**Persona A — The Engineer**
- Works late, high arousal at night
- Talks about servers, code, deployments
- Uses technical shorthand — "pushed to prod", "the box", "nginx"
- Calm baseline, spikes when something breaks
- Expected graph: server entities, late-night work pattern, deployment skill candidate

**Persona B — The Athlete**
- Consistent morning routine, logs workouts
- Talks about gym, nutrition, sleep
- Refers to exercises by shorthand — "chest day", "leg day"
- High arousal in the morning, low at night
- Expected graph: fitness entities, workout pattern, log_workout skill candidate

**Persona C — The Ambiguous One**
- Switches topics mid-conversation
- Uses different names for the same thing across sessions
- Low confidence signal — hard to predict
- Tests alias resolution, structural similarity, graceful degradation
- Expected graph: messier, more alias_of edges, lower confidence scores

---

## What to Audit

### After the Reasoning Pass

- Every entity mentioned in the session has a node
- No node `content` is a sentence — atomic labels only
- Every node has a `reasoning` field that makes sense
- Confidence scores reflect frequency + affect signal
- VAD signal is visible in the reasoning fields

### After the Integrate Pass

- Known aliases resolved correctly (Mike → Michael, "the gym" → gym entity)
- No committed nodes mutated
- Weights updated on existing edges where new signal corroborates them
- New edges added where they didn't exist but should
- Every edge has a `reasoning` field
- Low confidence pairs correctly skipped

---

## Test Runs

A test run is one full cycle:

```
synthetic persona chats with Skyra (N turns)
  → session history + VAD produced
  → cron fires → reasoning job runs
  → reasoning job completes → integrate job runs
  → graph audited against expected output
```

Run the same persona multiple times. The graph should strengthen and stabilize — not drift. Repeated sessions should increase weights, not produce duplicate nodes.

---

## Failure Modes to Watch For

- **Sentences in content** — model drifted, not writing atomic labels
- **Duplicate nodes** — integrate not resolving aliases
- **Missing reasoning fields** — model skipping the audit trail
- **Wrong aliases** — Mike and a different Mike merged
- **Weight drift** — weights not reflecting actual signal frequency
- **Empty graph** — reasoning produced nothing, silent failure

---

## Related

- `docs/arch/v1/skill-reasoning.md` — what the reasoning pass produces
- `docs/arch/v1/skill-integrate.md` — what integrate does with the output
- `docs/arch/v1/memory-structure.md` — node + edge schema
- `docs/arch/v1/observational-store.md` — session history + VAD streams
