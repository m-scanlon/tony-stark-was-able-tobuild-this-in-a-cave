# Skyra Listener Service

The listener is the always-on Pi-side service. It handles voice capture, deterministic intent gating, triage, provisional response decisions, event delivery to the control plane, and turn reconciliation.

The Pi is non-authoritative for semantic decisions. Every request goes to the Mac. But the Pi can give a provisional answer from fresh cached context while Mac processes authoritatively in parallel.

## Responsibilities

- Always-on audio pipeline: wake word → VAD → STT
- Deterministic intent gate (no LLM, no inference)
- Pi triage: classify latency, delegation need, and provisional answer eligibility
- Front-door fast model: structure the event, generate provisional answer when eligible
- Durable local outbox for reliable event delivery to control plane
- Turn reconciliation: handle Mac response and correct/confirm provisional if needed
- TTS output of Mac-authored content

## Boundaries

Pi is allowed to:
- Emit non-semantic ACKs (earcon, LED, short wait phrase)
- Give a provisional answer from fresh LCACHE context, clearly qualified as non-authoritative
- Render Mac-authored `UPDATE | PLAN_PROGRESS | CLARIFY | PLAN_APPROVAL_REQUIRED | FINAL | ERROR`

Pi must not:
- Generate semantic answers without fresh cached context
- Claim an action completed unless confirmed by Mac
- Claim state changes occurred unless confirmed by Mac
- Write or modify system memory or project state

---

## 1. Audio Pipeline

```
wake word → VAD → STT → intent gate → triage → front-door model
```

**Wake word**: openWakeWord or Porcupine. Always-on, low power.

**VAD**: captures utterance boundaries. Avoids sending silence.

**STT**: Whisper small or base. Runs locally on Pi. Optional: stream audio to Mac for faster STT (see `docs/arch/v1/scyra.md` §7).

---

## 2. Intent Gate (Deterministic)

Runs before any LLM inference. Stateless and synchronous.

Rules (in order):
- No wake word detected → `ignore`
- Empty or filler transcript ("um", "uh", silence) → `ignore`
- Cancel intent detected ("never mind", "stop", "cancel") → `ignore`
- Otherwise → `dispatch`

No model involved. No ambiguity. Fast gate only.

---

## 3. Pi Triage Layer

Runs after intent gate passes. Fast rules or tiny classifier. Produces routing and UX hints for the rest of the pipeline.

### Triage Output

```json
{
  "latency_class": "fast | medium | slow",
  "needs_delegation": true,
  "hint_target": "control_plane | agent:<id> | gpu:<id>",
  "ack_policy": "silent | nonverbal | spoken_if_slow",
  "confidence": 0.84,
  "provisional_eligible": true,
  "cache_age_seconds": 420
}
```

`provisional_eligible` — whether the front-door model is allowed to attempt a provisional answer. Set by triage based on intent type and cache freshness check (see Section 5).

`cache_age_seconds` — age of the most relevant item in LCACHE for this request, in seconds. Used by the front-door model to decide whether to answer provisionally.

Triage output is a hint, not final authority. Mac makes all authoritative decisions.

---

## 4. Front-Door Fast Model

**Model**: Llama 3.2 3B Instruct (GGUF, Q4_K_M)
**Context target**: 4096–8192 tokens

Context budget (from Context Injector):
- 35% system instructions
- 25% live conversation
- 25% injected context package (LCACHE)
- 15% response/scratch headroom

Responsibilities:
- Structure the outgoing `voice_event_v1`
- If `provisional_eligible: true` — attempt a provisional answer from LCACHE
- Emit the appropriate user feedback ACK based on `ack_policy`

The front-door model runs event-driven — not always-on inference. It is invoked per utterance after the intent gate passes.

---

## 5. Provisional Response Model

The Pi can give a provisional answer from cached context while the request is being processed authoritatively on Mac. This reduces perceived latency for retrieval-style questions without breaking the authoritative model.

### When Pi May Answer Provisionally

All conditions must be met:

| Condition | Requirement |
|---|---|
| Intent type | Read-only / retrieval only — "what did I decide", "what's the status of", "remind me" |
| Cache freshness | Most relevant LCACHE item `retrieved_at` within **30 minutes** |
| Confidence | Triage confidence above threshold (default: `0.75`) |
| Intent clarity | Unambiguous — no multi-step, no action words |
| LCACHE hit | At least one item with meaningful relevance score to the request |

### When Pi Must Not Answer Provisionally

- Request implies action or state mutation ("turn off", "update", "delete", "create", "set")
- Cache is stale (`cache_age_seconds > 1800`) or empty
- Triage confidence below threshold
- Intent is ambiguous or multi-part
- `stale: true` on the current context package

### What Pi Says

Pi prefixes all provisional answers with a clear non-authoritative qualifier:

> "Based on what I have — [answer]. Let me confirm that for you."

The qualifier is not optional. Pi must never state a provisional answer as fact.

### Provisional Answer Output (front-door model)

```json
{
  "provisional_answer": "Based on what I have — you decided on weekly Tekkit backups at 02:00 UTC last Tuesday. Let me confirm that for you.",
  "provisional_confidence": 0.81,
  "source_item_id": "mem_91",
  "cache_age_seconds": 420
}
```

If `provisional_eligible` is false or the front-door model cannot form a confident answer, `provisional_answer` is null and Pi emits a non-semantic ACK only.

---

## 6. Event Schema

### voice_event_v1

```json
{
  "schema": "voice_event_v1",
  "turn_id": "turn_8f4c",
  "ts": "2026-02-20T18:10:12Z",
  "transcript": "what did I decide about the Tekkit backups",
  "triage_hints": {
    "latency_class": "medium",
    "needs_delegation": true,
    "hint_target": "control_plane",
    "ack_policy": "spoken_if_slow",
    "confidence": 0.84,
    "provisional_eligible": true,
    "cache_age_seconds": 420
  },
  "pi_gave_provisional": true,
  "provisional_text": "Based on what I have — you decided on weekly Tekkit backups at 02:00 UTC last Tuesday. Let me confirm that for you.",
  "context_window": {
    "session_summary": "...",
    "recent_turns": [],
    "active_project": "server_ops",
    "injected_facts": []
  },
  "context_state": {
    "total_context_tokens": 8192,
    "system_tokens": 1420,
    "live_conversation_tokens": 980,
    "response_reserve_tokens": 512,
    "available_for_injection": 5280
  }
}
```

`pi_gave_provisional` and `provisional_text` are included so Mac knows what Pi said. Mac uses this to tune its response — it may confirm, correct, or add detail rather than restating from scratch.

`context_state` is included on every request. Pi computes `available_for_injection` as `total_context_tokens - system_tokens - live_conversation_tokens - response_reserve_tokens`. Mac fans this out internally to the Context Injector, which uses it to size the next context package. Pi does not interface with the Context Injector directly — it only speaks to the Mac API Gateway.

---

## 7. Outbox and Event Delivery

Pi uses a local SQLite outbox for durable event delivery.

```
1. Front-door produces voice_event_v1
2. Pi writes event to local outbox (before sending)
3. Pi sends event to Mac API Gateway
4. Mac writes to inbox (event_id PRIMARY KEY — idempotent)
5. Mac sends transport ACK
6. Pi deletes outbox row on ACK match
```

If no ACK arrives within timeout, Pi retries from outbox. Events are delivered at-least-once. Mac inbox deduplicates by `event_id`.

Transport ACK is machine-to-machine only. It is never spoken to the user.

Reference: `docs/arch/v1/event-ingress-ack.md`

---

## 8. Turn Loop and Reconciliation

### State Machine

```
IDLE → LISTENING → TRANSCRIBED → FORWARDED → ACKED → RUNNING → RESOLVED
```

- `RUNNING → RUNNING` on `UPDATE` or `PLAN_PROGRESS`
- `RUNNING → LISTENING` on `CLARIFY`
- `RUNNING → RESOLVED` on `FINAL | ERROR`

### Reconciliation Protocol

When Mac responds with `FINAL`, Pi reconciles against any provisional answer it gave:

| Mac result | Pi behavior |
|---|---|
| Mac agrees with provisional | Silent, or soft "confirmed" — do not repeat the answer |
| Mac has additional detail | Speak the delta: "Also worth noting — [Mac addition]" |
| Mac contradicts provisional | Correct clearly: "Update — [Mac answer]" |

Pi appends the authoritative Mac content to the local context window on `FINAL | ERROR` and closes the turn.

### Pi-Side Turn Loop (Pseudocode)

```python
def on_user_utterance(audio_chunk_stream):
    turn_id = new_turn_id()
    transcript = stt(audio_chunk_stream)
    triage = pi_fast_triage(transcript)
    context_window = context_manager.snapshot_for_turn(turn_id)

    provisional = None
    if triage["provisional_eligible"]:
        provisional = front_door_model.attempt_provisional(transcript, lcache)
        if provisional:
            tts_speak(provisional["provisional_answer"])

    event = build_voice_event_v1(
        turn_id=turn_id,
        transcript=transcript,
        triage=triage,
        provisional=provisional,
        context_window=context_window,
    )

    outbox.persist(event)
    if not provisional:
        emit_user_ack(triage["ack_policy"])
    transport.send(event)

    while True:
        msg = transport.recv_for_turn(turn_id, timeout=TURN_TIMEOUT_S)
        if msg is None:
            transport.retry_from_outbox(event["event_id"])
            continue

        if msg["message_type"] in ("UPDATE", "PLAN_PROGRESS"):
            maybe_speak_progress(msg["text"], triage["ack_policy"])
            continue

        if msg["message_type"] == "CLARIFY":
            tts_speak(msg["text"])
            context_manager.append_assistant(turn_id, msg["text"], authoritative=True)
            return "needs_user_input"

        if msg["message_type"] == "PLAN_APPROVAL_REQUIRED":
            tts_speak(msg["text"])
            context_manager.append_assistant(turn_id, msg["text"], authoritative=True)
            return "awaiting_plan_approval"

        if msg["message_type"] in ("FINAL", "ERROR"):
            reconciled = reconcile(provisional, msg["text"])
            tts_speak(reconciled)
            context_manager.append_assistant(turn_id, msg["text"], authoritative=True)
            outbox.delete_if_acked(event["event_id"])
            return "resolved"
```

---

## 9. LCACHE (Listener Context Cache)

LCACHE is the Pi-side context store fed by the Context Injector service on Mac. It holds compressed, ranked context packages ready for front-door model use.

### Freshness Check for Provisional Answers

Each item in a context package carries a `retrieved_at` timestamp. Before giving a provisional answer, the front-door model checks the most relevant item's `retrieved_at` against the 30-minute threshold:

```python
PROVISIONAL_MAX_AGE_SECONDS = 1800  # 30 minutes, configurable

def is_fresh_enough_for_provisional(item):
    age = now() - parse_iso(item["retrieved_at"])
    return age.total_seconds() < PROVISIONAL_MAX_AGE_SECONDS
```

If no relevant item meets the freshness threshold, `provisional_eligible` is set to false by triage and no provisional answer is attempted.

### Cache Failure Behavior

- Context Injector unavailable → keep last package until TTL grace expires, then run on base + live context only
- Stale package (`stale: true`) → do not attempt provisional answers, emit non-semantic ACK only

---

## 10. Recommended Front-Door Model

- `Llama-3.2-3B-Instruct` (GGUF, `Q4_K_M` preferred)
- Context target: `4096` to `8192` tokens
- Keep long memory retrieval and heavy context assembly in the control plane

Context compression reference: `skyra/internal/context/compress`

---

## 11. v2 Planned Upgrade — Reactive Layer

A two-tier model stack is planned for v2 to make interactions feel more natural. The gap between utterance end and first response is where the robotic feeling comes from.

- **Tier 1 (reactive)** — rule-based phrase pool, fires ~100ms after utterance ends. Maps triage output to natural acknowledgements ("on it", "let me check", "mhm"). No model, no latency.
- **Tier 2 (front-door)** — existing 3B model, unchanged. Provisional answers and event structuring follow after.

This is a drop-in addition before the front-door model invocation. Nothing in the core architecture changes. Deferred until v1 is stable.

---

## 12. Related Docs

- `skyra/services/context-injector/README.md` — context package format, push strategy, trigger model
- `docs/arch/v1/scyra.md` — full system architecture and voice request flow
- `docs/arch/v1/event-ingress-ack.md` — outbox/inbox reliability contract
- `skyra/internal/project/README.md` — project service, boundary enforcement

---

## Run Locally

```bash
cd skyra/services/listener
python -m venv .venv
source .venv/bin/activate
pip install -r requirements.txt
uvicorn app.main:app --host 0.0.0.0 --port 8090
```

## Run with Docker

```bash
cd skyra/services/listener
docker build -t skyra-listener:dev .
docker run --rm -p 8090:8090 skyra-listener:dev
```

## Endpoints

- `GET /health`
- `POST /listener/event`
