# Phase 1 Model Endpoint Example

This project currently uses an OpenAI-compatible `/v1/chat/completions` endpoint exposed over HTTPS.

## Quick Test

```bash
curl -s https://league-nasa-sculpture-dos.trycloudflare.com/v1/chat/completions \
  -H "Content-Type: application/json" \
  -d '{
    "model": "/content/models/DeepSeek-R1-Distill-Qwen-32B-Q4_K_M.gguf",
    "messages":[{"role":"user","content":"Say hi in one sentence."}],
    "max_tokens":64,
    "temperature":0.2
  }'
```

## Notes

- Endpoint is OpenAI-compatible, so it can be integrated with standard chat-completions clients.
- Current model path:
  - `/content/models/DeepSeek-R1-Distill-Qwen-32B-Q4_K_M.gguf`
- Example responses may include reasoning artifacts such as `</think>`.

## Output Sanitization (Required)

Before sending model output to voice or chat clients:

1. Remove `<think>...</think>` blocks.
2. Remove stray `</think>` tokens.
3. Trim whitespace and avoid duplicate assistant text.

## Suggested Runtime Config (Phase 1)

- Timeout: 30s
- Retries: 2 (exponential backoff)
- Temperature: `0.2` for deterministic assistant behavior
- Log per request: latency, prompt tokens, completion tokens, total tokens, errors
