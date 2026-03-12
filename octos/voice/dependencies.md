# Voice Shard Dependencies

Validated against upstream docs on 2026-03-11.

## Recommended baseline (local-first)

- `python>=3.9,<3.15` (3.10+ preferred)
- `sounddevice==0.5.*` (microphone/speaker I/O via PortAudio bindings)
- `openwakeword==0.6.0` (wake-word detection, open source)
- `webrtcvad-wheels==2.0.14` (standalone VAD with up-to-date binary wheels)
- `faster-whisper==1.2.1` (STT via CTranslate2)
- `piper-tts==1.4.1` (local TTS)
- `websockets==16.0` on Python 3.10+ (if voice response channel uses WebSocket)
- `websockets==15.0.1` on Python 3.9
- `grpcio==1.78.0` and `grpcio-tools==1.78.0` (if voice response channel uses gRPC)

Example:

```bash
pip install \
  "sounddevice==0.5.*" \
  "openwakeword==0.6.0" \
  "webrtcvad-wheels==2.0.14" \
  "faster-whisper==1.2.1" \
  "piper-tts==1.4.1" \
  "websockets==15.0.1"
```

Or install from this repo file:

```bash
pip install -r octos/voice/requirements.txt
```

## Optional alternatives

- `pvporcupine==4.0.2` for wake word (actively updated, requires Picovoice `AccessKey`)
- `coqui-tts==0.27.5` for richer neural TTS/voice cloning (heavier than Piper)
- `openai-whisper` (official Whisper package) if you want compatibility with that stack

## System notes

- `faster-whisper` GPU path requires CUDA 12 + cuDNN 9 (`nvidia-cublas-cu12`, `nvidia-cudnn-cu12==9.*`).
- `openwakeword` installs `onnxruntime` (+ `tflite-runtime` on Linux) as transitive dependencies.
- `sounddevice` requires a working PortAudio backend on the host OS.
- `piper-tts` is GPL-3.0-or-later. If licensing is a concern, prefer `coqui-tts` (MPL-2.0).
- If your runtime is Python 3.9, keep `websockets==15.0.1`; for Python 3.10+, use `websockets==16.0`.
