# voice_event Schema Changelog

## v2 (planned)

New fields in `triage_hints`:
- `provisional_eligible` — whether the Pi is allowed to give a provisional spoken response before the brain replies
- `cache_age_seconds` — how stale the Pi's local context cache is, so the brain can weight it accordingly
- `needs_delegation` — whether the Pi is delegating to the brain or handling locally (v2 introduces local handling)
- `hint_target` — which shard the Pi thinks should handle this (v2 introduces multi-brain routing)

New top-level fields:
- `context_state` — Pi's view of the brain's context window token budget (total, system, live, reserve, available). Allows brain to right-size the context package it sends back.
- `pi_gave_provisional` — whether the Pi already spoke a provisional response
- `provisional_text` — what the Pi said, so the brain can reconcile or correct it
- `context_window` — Pi's local cache of session context (summary, recent turns, active agent, injected facts). Allows brain to detect staleness and reconcile what the Pi knew vs what it knows.

## v1 (current)

Initial schema. Pi always delegates. Single brain. No provisional responses.
