# voice_event Schema Changelog

## v2 (planned)

New fields in `triage_hints`:
- `provisional_eligible` — whether the Voice Shard is allowed to give a provisional spoken response before the brain replies
- `cache_age_seconds` — how stale the Voice Shard's local context cache is, so the brain can weight it accordingly
- `needs_delegation` — whether the Voice Shard is delegating to the brain or handling locally (v2 introduces local handling)
- `hint_target` — which shard the Voice Shard thinks should handle this (v2 introduces multi-brain routing)

New top-level fields:
- `context_state` — Voice Shard's view of the Brain Shard's context window token budget (total, system, live, reserve, available). Allows brain to right-size the context package it sends back.
- `pi_gave_provisional` — whether the Voice Shard already spoke a provisional response
- `provisional_text` — what the Voice Shard said, so the Brain Shard can reconcile or correct it
- `context_window` — Voice Shard's local cache of session context (summary, recent turns, active agent, injected facts). Allows brain to detect staleness and reconcile what the Voice Shard knew vs what the Brain Shard knows.

## v1 (current)

Initial schema. Voice Shard always delegates. Single Brain Shard. No provisional responses.
