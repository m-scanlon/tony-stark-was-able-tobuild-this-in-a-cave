# Command Surface (Prelim)

## Core Framing

This document preserves an intermediate protocol idea while clarifying the newer active direction.

Earlier docs moved from:

- `namespace`
- to `command_set`

That was useful, but the current architecture now points somewhere else.

The current active direction is:

```text
skyra <node> <primitive> -<args> -reason "<why this command is being emitted>"
```

This makes the node explicit and keeps the primitive layer small.

## Why This Shift Happened

The newer direction comes from three pressures:

- nodes are becoming explicit runtime operators
- typed stimulus is becoming the routing surface
- the primitive split is becoming smaller and cleaner than a large top-level command taxonomy

The current primitive split is:

- `recall`
- `learn`
- `observe`
- `act`

That makes a node-first protocol read more naturally than the earlier `command_set` framing.

## Example Direction

Examples might later look like:

```text
skyra jarvis act -target human -content "the current frame requires a user-facing response" -modality text -timestamp now -reason "the current frame requires a user-facing response"
skyra jarvis observe -target screen -reason "the current frame requires fresh world intake before responding"
skyra jarvis recall -reason "the current stimulus introduced structural cues worth recall lookup"
skyra stark act -target laptop -content "discover capability surface" -modality probe -timestamp now -reason "the device needs capability discovery"
skyra stark learn -episode_id ep_123 -reason "the episode should be consolidated into retained experience"
```

These are only directional examples.

They do not yet lock the final primitive grammar.

## Working Recall Command Shape

One useful working shape for recall is:

```text
skyra <node> recall \
  -entity <entity_id> \
  -relationship <relationship_id> \
  -bundle <left_entity_id>:<relationship_id>:<right_entity_id> \
  -top_k <n> \
  -reason "<why recall is being invoked now>"
```

This is a good fit for the current recall direction because recall is structural rather than freeform-text retrieval.

The retrieval surface may therefore accept:

- entity-only signals for broad retrieval
- relationship-only signals for more specific retrieval
- fully bound relational bundles for the strongest structural match

Example:

```text
skyra jarvis recall \
  -entity assistant \
  -entity terraform \
  -relationship help_with \
  -relationship has_property \
  -bundle assistant:help_with:terraform \
  -top_k 8 \
  -reason "the current stimulus explicitly mentions assistant help and terraform difficulty"
```

This should still be treated as a working command shape rather than a frozen final grammar.

The important `v0` rule is already stable, though:

- every emitted command must carry `-reason`

## Relationship To Contracts

The active node contract should define:

- which primitives are allowed
- what stimulus types are accepted and emitted
- what `act` modalities are allowed
- how `target`, `content`, `modality`, and `timestamp` are constrained inside `act`
- what loop or execution envelopes are permitted

This means:

- the runtime substrate stays generic
- the contract bounds what the node may actually use

## Current Design Posture

The strongest current claim is:

- command execution should now be understood through explicit nodes and a small primitive set

This is a cleaner fit for the current runtime direction.

It does not yet define:

- the final node vocabulary
- the final primitive grammar
- the final `act` modality vocabulary
- the final `content` encoding rules for `act`
- the final timestamp conventions for `act`

It does define one important audit rule:

- command emission requires explicit rationale

## Short Framing

The runtime command surface should now be thought of as:

- `skyra <node> <primitive> -<args> -reason "<why this command is being emitted>"`

`-reason` is mandatory because it is part of the audit trail, not optional commentary.
