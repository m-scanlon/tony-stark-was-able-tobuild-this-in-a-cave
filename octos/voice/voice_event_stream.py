"""Wake-word audio ingress to hydrated voice_event stream with audit logging."""

from __future__ import annotations

from dataclasses import dataclass
from datetime import datetime, timezone
from hashlib import sha256
import copy
import json
from pathlib import Path
from typing import Any, Callable, Iterator, Mapping
import uuid

from .wake_to_text_stream import WakeToTextConfig, wake_word_to_text_stream


def _utc_now_iso() -> str:
    return datetime.now(timezone.utc).replace(microsecond=0).isoformat().replace("+00:00", "Z")


@dataclass(frozen=True)
class VoiceEventStreamConfig:
    device_id: str
    location_tag: str
    audit_log_path: str = "octos/voice/audit.log.jsonl"
    pending_job_id: str | None = None
    waiting_for: str | None = None
    context_blob: Mapping[str, Any] | None = None


def wake_word_to_voice_event_stream(
    *,
    wake_config: WakeToTextConfig = WakeToTextConfig(),
    event_config: VoiceEventStreamConfig,
    emit: Callable[[dict[str, Any]], None] | None = None,
) -> Iterator[dict[str, Any]]:
    """Yield hydrated `voice_event_v1` envelopes from live wake-word audio input.

    Pipeline:
        microphone -> wake word -> utterance transcript -> voice_event_v1

    Audit:
        Appends structured JSON records to `event_config.audit_log_path` for every
        event build and emit result, with chained hashes for run-local integrity.
    """
    audit_path = Path(event_config.audit_log_path)
    audit_path.parent.mkdir(parents=True, exist_ok=True)

    run_id = f"voice_run_{uuid.uuid4().hex[:12]}"
    previous_record_hash = "GENESIS"

    def audit(action: str, **fields: Any) -> None:
        nonlocal previous_record_hash
        record = {
            "ts": _utc_now_iso(),
            "run_id": run_id,
            "action": action,
            "prev_record_hash": previous_record_hash,
            **fields,
        }
        canonical = json.dumps(record, sort_keys=True, separators=(",", ":"), ensure_ascii=True)
        record_hash = sha256(canonical.encode("utf-8")).hexdigest()
        record["record_hash"] = record_hash
        with audit_path.open("a", encoding="utf-8") as f:
            f.write(json.dumps(record, sort_keys=True, ensure_ascii=True) + "\n")
        previous_record_hash = record_hash

    def default_context_blob(now_iso: str) -> dict[str, Any]:
        return {
            "cache_ts": now_iso,
            "agents": [],
            "recent_turns": [],
            "active_job": None,
        }

    def default_triage_hints(transcript: str) -> dict[str, Any]:
        return {
            "intent": {
                "summary": transcript[:240],
                "confidence": 0.5,
            },
            "latency_class": {
                "value": "interactive",
                "confidence": 0.5,
            },
            "ack_policy": {
                "value": "spoken_if_slow",
                "confidence": 0.5,
            },
        }

    audit(
        "stream_started",
        device_id=event_config.device_id,
        location_tag=event_config.location_tag,
    )

    try:
        for transcript in wake_word_to_text_stream(wake_config):
            turn_id = f"turn_{uuid.uuid4().hex[:8]}"
            now_iso = _utc_now_iso()

            context_blob: dict[str, Any]
            if event_config.context_blob is None:
                context_blob = default_context_blob(now_iso)
            else:
                context_blob = copy.deepcopy(dict(event_config.context_blob))
                context_blob.setdefault("cache_ts", now_iso)
                context_blob.setdefault("agents", [])
                context_blob.setdefault("recent_turns", [])
                context_blob.setdefault("active_job", None)

            voice_event = {
                "schema": "voice_event_v1",
                "turn_id": turn_id,
                "ts": now_iso,
                "device_id": event_config.device_id,
                "location_tag": event_config.location_tag,
                "transcript": transcript,
                "triage_hints": default_triage_hints(transcript),
                "session_state": {
                    "pending_job_id": event_config.pending_job_id,
                    "waiting_for": event_config.waiting_for,
                },
                "context_blob": context_blob,
            }

            event_hash = sha256(
                json.dumps(voice_event, sort_keys=True, separators=(",", ":"), ensure_ascii=True).encode("utf-8")
            ).hexdigest()

            audit(
                "voice_event_built",
                turn_id=turn_id,
                event_hash=event_hash,
                transcript=transcript,
                transcript_chars=len(transcript),
            )

            if emit is not None:
                try:
                    emit(voice_event)
                    audit("voice_event_emitted", turn_id=turn_id, event_hash=event_hash)
                except Exception as exc:  # pragma: no cover - external transport dependent
                    audit(
                        "voice_event_emit_error",
                        turn_id=turn_id,
                        event_hash=event_hash,
                        error=f"{type(exc).__name__}: {exc}",
                    )
                    raise

            audit("voice_event_yielded", turn_id=turn_id, event_hash=event_hash)
            yield voice_event
    except Exception as exc:
        audit("stream_error", error=f"{type(exc).__name__}: {exc}")
        raise
