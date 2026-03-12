"""Minimal wake-word to text stream pathway for the Voice Shard.

Public API:
    wake_word_to_text_stream(...)

This method owns the full ingress loop:
    microphone -> wake word -> VAD end-of-utterance -> STT text yield
"""

from __future__ import annotations

from dataclasses import dataclass
from typing import Iterator, Sequence


@dataclass(frozen=True)
class WakeToTextConfig:
    sample_rate: int = 16000
    frame_ms: int = 30

    wake_threshold: float = 0.35
    wakeword_models: Sequence[str] | None = None

    vad_aggressiveness: int = 2
    end_silence_ms: int = 700
    min_speech_ms: int = 250
    max_utterance_ms: int = 10000

    stt_model: str = "small"
    stt_device: str = "cpu"
    stt_compute_type: str = "int8"
    stt_language: str | None = "en"
    stt_beam_size: int = 1


def wake_word_to_text_stream(config: WakeToTextConfig = WakeToTextConfig()) -> Iterator[str]:
    """Yield transcribed text after each wake-word-triggered utterance.

    The method runs continuously and yields one final text string per utterance.
    It is intentionally one clean pathway with a single public method.
    """
    try:
        import numpy as np
        import openwakeword
        import sounddevice as sd
        import webrtcvad
        from faster_whisper import WhisperModel
        from openwakeword.model import Model as WakeModel
    except ImportError as exc:  # pragma: no cover - environment dependent
        raise RuntimeError(
            "Voice dependencies are missing. Install with: "
            "pip install -r octos/voice/requirements.txt"
        ) from exc

    if config.frame_ms not in (10, 20, 30):
        raise ValueError("frame_ms must be one of 10, 20, or 30 for webrtcvad")
    if not 0 <= config.vad_aggressiveness <= 3:
        raise ValueError("vad_aggressiveness must be in [0, 3]")
    if not 0.0 <= config.wake_threshold <= 1.0:
        raise ValueError("wake_threshold must be in [0.0, 1.0]")

    frame_samples = int(config.sample_rate * config.frame_ms / 1000)
    end_silence_frames = max(1, config.end_silence_ms // config.frame_ms)
    min_speech_frames = max(1, config.min_speech_ms // config.frame_ms)
    max_utterance_frames = max(1, config.max_utterance_ms // config.frame_ms)

    # Ensure built-in wake-word models exist locally when no custom path is set.
    if not config.wakeword_models:
        openwakeword.utils.download_models()

    wake_model_kwargs = {}
    if config.wakeword_models:
        wake_model_kwargs["wakeword_models"] = list(config.wakeword_models)
    wake_model = WakeModel(**wake_model_kwargs)

    vad = webrtcvad.Vad(config.vad_aggressiveness)
    stt = WhisperModel(
        config.stt_model,
        device=config.stt_device,
        compute_type=config.stt_compute_type,
    )

    armed = False
    speech_started = False
    speech_frames = 0
    silence_frames = 0
    utterance_frames: list[bytes] = []

    with sd.RawInputStream(
        samplerate=config.sample_rate,
        channels=1,
        dtype="int16",
        blocksize=frame_samples,
    ) as mic:
        while True:
            chunk, overflowed = mic.read(frame_samples)
            if overflowed:
                continue

            pcm = bytes(chunk)
            pcm_i16 = np.frombuffer(pcm, dtype=np.int16)

            if not armed:
                wake_scores = wake_model.predict(pcm_i16)
                max_score = max(
                    (
                        float(score)
                        for score in wake_scores.values()
                        if isinstance(score, (int, float))
                    ),
                    default=0.0,
                )
                if max_score < config.wake_threshold:
                    continue

                armed = True
                speech_started = False
                speech_frames = 0
                silence_frames = 0
                utterance_frames = [pcm]
                continue

            utterance_frames.append(pcm)
            if vad.is_speech(pcm, config.sample_rate):
                speech_started = True
                speech_frames += 1
                silence_frames = 0
            elif speech_started:
                silence_frames += 1

            done_for_silence = speech_started and silence_frames >= end_silence_frames
            done_for_length = len(utterance_frames) >= max_utterance_frames
            if not (done_for_silence or done_for_length):
                continue

            text = ""
            if speech_frames >= min_speech_frames:
                utterance_pcm = b"".join(utterance_frames)
                audio = (
                    np.frombuffer(utterance_pcm, dtype=np.int16).astype(np.float32)
                    / 32768.0
                )
                segments, _ = stt.transcribe(
                    audio,
                    language=config.stt_language,
                    beam_size=config.stt_beam_size,
                    vad_filter=False,
                    condition_on_previous_text=False,
                )
                text = " ".join(seg.text.strip() for seg in segments).strip()

            armed = False
            speech_started = False
            speech_frames = 0
            silence_frames = 0
            utterance_frames = []

            if text:
                yield text
