# voice_event Schema

The `voice_event` is the only event type the Voice Shard sends to the Brain Shard. It carries the user's transcript, intent hints, session state, and the context blob pushed by CIX.

The envelope is assembled by the shard's transport layer before sending. The intent model only produces `triage_hints` — all other fields are hydrated by the shard at send time.

Implementation dependency set for the Voice Shard runtime: `dependencies.md`.

## Runtime Pathway (Lean)

Wake word to text stream is implemented as one method:

```python
from octos.voice import WakeToTextConfig, wake_word_to_text_stream

config = WakeToTextConfig()
for text in wake_word_to_text_stream(config):
    print(text)
```

Method location: `wake_to_text_stream.py`.

Wake word directly to hydrated `voice_event_v1` stream (with audit logs):

```python
from octos.voice import (
    WakeToTextConfig,
    VoiceEventStreamConfig,
    wake_word_to_voice_event_stream,
)

wake_cfg = WakeToTextConfig()
event_cfg = VoiceEventStreamConfig(
    device_id="pi-livingroom-01",
    location_tag="home-main",
    audit_log_path="octos/voice/audit.log.jsonl",
)

for event in wake_word_to_voice_event_stream(
    wake_config=wake_cfg,
    event_config=event_cfg,
):
    print(event["turn_id"], event["transcript"])
```

Method location: `voice_event_stream.py`.

Audit logging:
- append-only JSONL file
- one line per lifecycle action (`stream_started`, `voice_event_built`, `voice_event_emitted`, `voice_event_yielded`, `stream_error`)
- includes `record_hash` and `prev_record_hash` for run-local hash chaining

---

## Envelope Hydration

The intent model only produces `triage_hints`. Everything else is assembled by the shard's transport layer before the event is sent.

| Field | Source |
|---|---|
| `schema` | Transport layer — stamped at send time |
| `turn_id` | Transport layer — generated per turn |
| `ts` | Transport layer — stamped at send time |
| `device_id` | Transport layer — shard identity |
| `location_tag` | Transport layer — derived from network fingerprint at shard registration |
| `transcript` | STT model |
| `triage_hints` | Intent model — the only field the model produces |
| `session_state` | Shard local turn tracking |
| `context_blob` | Shard local context cache — kept warm by CIX, attached at hydration |

The context blob is not assembled at request time. CIX (Context Injector) proactively pushes compressed context packages to the shard's local cache whenever agent state changes. The shard attaches the cached blob at hydration. By the time the event reaches the Brain Shard, context is already embedded — no request-time context assembly needed.

---

## Fields

### `schema`
Version identifier. Always `voice_event_v1` for v1. Stamped by the shard's transport layer — not produced by the intent model.

### `turn_id`
Unique ID for this turn. Stamped by the shard's transport layer before sending — not produced by the intent model. Used for deduplication and outbox tracking. The brain generates its own `event_id` on ingress — `turn_id` is the Voice Shard's reference, not the Brain Shard's.

### `ts`
ISO 8601 timestamp. Stamped by the shard's transport layer before sending — not produced by the intent model.

### `device_id`
The identity of the shard that captured this event. Stamped by the shard's transport layer — not produced by the intent model. Allows the brain to know which shard to route responses back to.

### `location_tag`
The physical location identifier for the shard — derived from the network fingerprint (SSID, gateway MAC, subnet) recorded at shard registration. The capability resolver uses this tag to filter shard capabilities to those co-located with the ingress shard. This is how spatial awareness works: the ingress shard's location tag scopes which physical devices can be reached without an explicit location in the user's request.

Named at first registration on a new network — the system prompts once ("what should I call this location?"). Subsequent registrations on the same network fingerprint are automatic.

### `transcript`
What the user said, as plain text. Output of the Voice Shard's STT model.

### `triage_hints`
Intent classification produced by the Voice Shard's intent model. Tells the brain how to prioritize and how to instruct the UX model to respond. Each classification carries its own confidence score so the brain can independently decide how much to trust each one.

> **Note:** The fields inside `triage_hints` are not fully locked. The UX model design will determine what hints it actually needs. The examples below are directional — treat them as placeholders until the UX layer is designed.

#### `intent`
What the user wants.

- `summary` — one sentence describing the user's intent
- `confidence` — float 0.0–1.0. How confident the model is in the intent classification. Low confidence may cause the brain to re-derive intent from the transcript directly rather than trusting this field.

#### `latency_class` (? — not locked, depends on UX model design)
How urgently the user expects a response.

- `value`: `interactive | background | deferred`
  - `interactive` — user is waiting, respond as fast as possible
  - `background` — can run while the user does something else
  - `deferred` — no urgency, schedule it for later
- `confidence` — float 0.0–1.0

#### `ack_policy` (? — not locked, depends on UX model design)
How the UX model should acknowledge the request while the brain is working. The exact values here will be defined once the UX model and its acknowledgement behavior are designed.

- `value`: `spoken_if_slow | earcon_only | silent`
  - `spoken_if_slow` — speak a wait phrase if the response is taking time
  - `earcon_only` — play a sound only, no words
  - `silent` — no acknowledgement
- `confidence` — float 0.0–1.0

### `session_state`
Responsible for managing context about outgoing jobs and syncing that context with the brain. The shard tracks what jobs are in flight on its side — this field is how it communicates that state so the brain can route and continue correctly.

Shape is settled for v1 (`pending_job_id`, `waiting_for`). Routing semantics — how the Brain Shard acts on this to distinguish new jobs from continuations — are still being designed. See `docs/arch/v1/next-steps.md` §2.

#### `pending_job_id` (? — may be restructured into a job object)
- `null` — new job, the brain should start fresh
- set to a job ID — this turn continues an existing open job (waiting on user approval, clarification, etc.)

#### `waiting_for` (? — may be restructured into a job object)
What the existing job is waiting on, if anything. Examples: `user_approval`, `clarification`.

### `context_blob`
The pre-assembled context package attached at hydration from the shard's local cache. Kept warm by CIX — not assembled at request time.

#### `cache_ts`
ISO 8601 timestamp of when this context blob was last pushed by CIX. The Brain Shard can use this to detect a stale cache — if `cache_ts` is old, context may not reflect the latest agent state.

#### `agents[]`
All registered agents with their current relevance scores. Every agent is present regardless of status — the front face transformer uses these scores to label the turn (in-domain or other) and route to relevant domain agents. No active/inactive filter — a dormant agent with a low score still appears.

- `agent_id` — agent identifier
- `relevance_score` — float 0.0–1.0. How active and relevant this domain has been recently.
- `status` — `active | paused | archived`

#### `recent_turns[]`
Recent turns from the active session. Gives the Brain Shard conversational continuity without a separate lookup.

#### `active_job`
The in-flight job if `session_state.pending_job_id` is set — null otherwise. Includes enough context for the brain to resume a continuation without reading the full job artifact.

---

## Versions

See `CHANGELOG.md` for what changed between versions and what fields are deferred to v2.
