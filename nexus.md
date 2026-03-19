# Nexus

## Overview

Nexus is the top-level ingress signal gate for Skyra.

It is the first hard boundary in front of the runtime. Every API exposed by an ingress shard sits behind Nexus. Nothing enters the system, reaches the API Gateway, or touches kernel execution without passing through Nexus first.

Its job is not deep reasoning. Its job is to guard ingress against boring signals, capture the richest possible input state, preserve it, and hand useful turns to the rest of the system without losing signal.

In practical terms, Nexus sits above the ingress-shard APIs and below the kernel execution path:

```
voice / text / device ingress
  -> ingress-shard APIs
  -> Nexus
  -> API Gateway / kernel handoff
  -> execution, memory, reply
```

## Responsibilities

- guard every ingress-shard API as the mandatory front door
- suppress low-value, repetitive, or uninteresting ingress before it reaches the runtime
- terminate ingress from user-facing channels
- decide what is signal and what is noise at the ingress edge
- preserve raw and derived ingress signals before they are flattened
- assign or attach stable ingress metadata
- normalize channel-specific input into a canonical handoff shape
- attach cached local context that is useful at request time
- emit an auditable event trail for every ingress turn

## Guard Boundary

Nexus is the choke point for ingress.

If an ingress shard exposes voice, text, webhook, device, or compatibility APIs, those APIs do not hand traffic directly into the system. They terminate at Nexus. Nexus decides whether the incoming traffic is meaningful enough to become a runtime event or whether it should be ignored, collapsed, delayed, or dropped as boring ingress.

That means Nexus is where ingress policy lives:

- what counts as meaningful signal
- what counts as boring or low-information ingress
- what should be ignored, collapsed, delayed, or dropped
- what metadata must be attached to promoted events
- what gets transformed into the canonical runtime handoff

Nexus is not the system's malicious-traffic defense layer. It is a signal-quality gate.

## Voice Capture

For voice, Nexus should preserve the richest signal set available before speech-to-text strips the audio down to transcript only.

That includes:

- transcript
- device identity
- location tag
- timestamp
- wake-word / utterance boundaries
- intent and acknowledgement hints
- session continuation state
- cached context attached at ingress time
- acoustic affect features such as `valence`, `arousal`, and `dominance`
- delivery features such as `speech_rate`, `pitch`, and `pitch_variance`
- token-level context-state evolution as the utterance arrives
- retrieval events triggered during the turn
- audit-chain records for build, emit, yield, and error

The principle is simple: if a signal may matter later, Nexus should capture it at ingress rather than trying to reconstruct it after the fact.

## Decision Principle

Nexus asks one question:

Is this signal worth waking the runtime for right now?

Every ingress event should end in one of four outcomes:

- `pass_now` — forward immediately
- `pass_collapsed` — merge with similar recent ingress, then forward
- `hold` — wait briefly for more signal
- `drop` — do not wake the runtime

This is a lightweight ingress judgment, not deep reasoning.

## nexus_decision_v1

Nexus should emit a decision artifact for every evaluated ingress event.

```text
nexus_decision_v1 {
  decision_id
  ts
  ingress_event_ref

  source {
    channel
    device_id
    location_tag
  }

  session_state {
    pending_job_id
    waiting_for
  }

  features {
    directedness
    continuity
    novelty
    urgency
    confidence
    corroboration
    cost_of_miss
    repetition_penalty
    uncertainty
  }

  voice_features {
    wake_confidence
    vad_confidence
    stt_confidence
    transcript_tokens
    interruption_marker
    continuation_marker
    valence
    arousal
    dominance
    speech_rate
    pitch
    pitch_variance
  }

  outcome {
    value              <- pass_now | pass_collapsed | hold | drop
    score
    hold_ms
    collapse_key
    reason_codes[]
  }
}
```

All scoring features should be normalized to a cheap comparable range, such as `0.0` to `1.0`. Nexus should prefer simple, inspectable features over opaque learned scoring.

## Signals That Matter

The core ingress features are:

- `directedness` — does this look aimed at the system
- `continuity` — does it continue an active job or recent turn
- `novelty` — is it meaningfully different from recent ingress
- `urgency` — does the user appear to need a fast response
- `confidence` — how complete and trustworthy is the capture
- `corroboration` — do multiple cheap signals agree
- `cost_of_miss` — how bad would it be to ignore this
- `repetition_penalty` — how duplicate is it relative to recent ingress
- `uncertainty` — how unstable or ambiguous is the event

For voice, useful low-cost features include:

- wake-word confidence
- VAD end-of-utterance confidence
- STT confidence
- transcript length and completeness
- interruption markers such as `stop`, `wait`, `no`, `actually`
- continuation markers such as `yes`, `do it`, `that one`
- affect features such as `valence`, `arousal`, and `dominance`
- delivery features such as `speech_rate`, `pitch`, and `pitch_variance`

Affect should be a modifier, not the main trigger. High arousal alone should not wake the runtime if the event is otherwise undirected and low-information.

## Hard Rules

Some ingress should bypass scoring:

- `pass_now` for explicit directed commands or questions
- `pass_now` for short replies to an active `pending_job_id`
- `pass_now` for interruption or control turns such as `stop`, `cancel`, `yes`, `no`
- `drop` for empty transcript, failed utterance, or obvious non-speech
- `pass_collapsed` for near-duplicate ingress inside a short collapse window

These rules keep Nexus responsive without forcing every event through a score.

## Scored Decision

If no hard rule fires, Nexus can use a simple additive score:

```text
pass_score =
  directedness * 3 +
  continuity * 3 +
  cost_of_miss * 4 +
  novelty * 2 +
  urgency * 1 +
  corroboration * 1 -
  repetition_penalty * 2 -
  uncertainty * 2
```

Suggested mapping:

- high score -> `pass_now`
- medium score -> `hold` or `pass_collapsed`
- low score -> `drop`

The point is not mathematical purity. The point is to make Nexus predictable, debuggable, and cheap to run on the ingress edge.

## Decision Flow

1. Receive ingress event from the shard-facing API.
2. Attach cheap local context such as recent turns, active job state, and collapse-window history.
3. Extract lightweight features from the event.
4. Apply hard pass/drop/collapse rules.
5. If no hard rule fires, compute `pass_score`.
6. Choose `pass_now`, `pass_collapsed`, `hold`, or `drop`.
7. Emit `nexus_decision_v1` with reason codes.
8. If the event passes, hand off the canonical ingress package downstream.

## Hold And Collapse Windows

Nexus should keep a short-lived local buffer for suppressed or unresolved ingress.

- `hold` exists for ambiguous but potentially meaningful events
- `pass_collapsed` exists for repeated ingress that should become one runtime event instead of many
- `drop` should still leave behind a minimal audit record

This keeps Nexus from being trigger-happy without discarding useful evidence too early.

## Reason Codes

Reason codes should be simple and inspectable. Good examples:

- `wake_word_present`
- `active_job_reply`
- `interrupt_request`
- `high_novelty`
- `duplicate_recent`
- `low_transcript_confidence`
- `empty_utterance`
- `affect_only_without_directedness`

## Canonical Role

Nexus is a boundary service, not the intelligence itself.

It does not own:

- deep reasoning
- long-running planning
- skill execution
- final truth in memory

It does own:

- ingress signal filtering
- ingress fidelity
- channel normalization
- turn-level observability
- reliable handoff into the runtime

## Outputs

Nexus should hand the rest of the system a canonical ingress package that can include:

- top-level event metadata
- decision outcome, score, and reason codes
- normalized user content
- rich voice-side affect and delivery signals
- local context snapshot
- session/job continuation state
- routing and UX hints
- audit references

The exact envelope can evolve, but the service contract should remain stable: Nexus preserves ingress richness and hands off a trustworthy, replayable event.

## Design Constraints

- nothing bypasses Nexus
- no loss of high-value ingress signal at capture time
- transport-specific adapters should stay thin
- Nexus should filter boring traffic, not become a general security perimeter
- reasoning and routing policy should remain downstream
- auditability must be first-class
- the service must support both lean v1 envelopes and richer future schemas

## Non-Goals

- replacing the kernel
- replacing the API Gateway
- becoming a second brain
- embedding business logic that belongs in skills or memory promotion

## Summary

Nexus is the top-level ingress signal gate that turns messy real-world input into a preserved, normalized, auditable ingress event for Skyra.
