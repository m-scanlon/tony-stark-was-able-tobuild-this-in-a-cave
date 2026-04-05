# Database Schema v0

## Scope

This document currently locks:

- the `actors` table
- the high-level direction for the stimulus registry

It does **not** lock the rest of the database schema yet.

At this stage, the following remain intentionally open:

- actor contract table shape
- final normalization of stimulus contract storage
- delegation / invocation edge tables
- validated dispatch tables
- result / event persistence tables

The goal here is to preserve stable identity records while the stimulus-first contract model settles.

## `actors`

The current locked actor table is:

```text
actors
- actor_id
- actor_name
- status
- current_contract_id
- created_by
- birth_timestamp
```

## Actor Field Meanings

### `actor_id`

Stable internal identifier for the actor.

Rules:

- must be unique
- must be immutable
- must not be reused for a different actor later

### `actor_name`

The actor's runtime / protocol-facing name.

Rules:

- should be unique
- should be the name used in the outer Skyra protocol

### `status`

The actor's current lifecycle state.

The exact enum may still be refined later.

### `current_contract_id`

Identifier pointing at the actor's currently active contract.

This field is locked even though the contract table itself is not yet fully locked.

### `created_by`

Records how the actor came into existence.

Rules:

- `created_by = bootstrap` for actors present at runtime startup
- otherwise `created_by` must be a valid `actor_id`

### `birth_timestamp`

Timestamp recording when the actor was created / born into the runtime.

## Actor Locked Rules

The following rules are now locked:

- every actor row has one stable `actor_id`
- every actor row has one runtime-facing `actor_name`
- every actor row has one `status`
- every actor row points to one `current_contract_id`
- every actor row records `created_by`
- every actor row records `birth_timestamp`
- `created_by` is either `bootstrap` or a valid `actor_id`

## Stimulus Registry Direction

The old capability-table framing is superseded by the active stimulus-registry direction.

The registry should store published contract records for routable public or downstream surfaces.

At minimum, each record should preserve:

- stable storage identity
- human-readable name
- typed `ExecutionSurface`
- published public primitive or boundary mode
- published request stimulus schema or ingress schema
- published response envelope schema when applicable
- publishing actor
- publication timestamp

The exact final normalization of those fields across tables is still open.

One useful current split is:

- actor-level contract rows remain separate from callable sense-schema rows
- multiple callable sense schemas may belong to one actor through `actor_id`

## Minimum Stable Registry Facts

Even though the final table layout is still open, the current stable direction is:

- every registered surface should have a stable storage identifier
- every registered surface should have a `name`
- every registered surface should carry an `ExecutionSurface`
- the initial execution-surface kinds are `actor` and `capability`
- actor surfaces should publish one request stimulus schema and one response envelope schema
- capability surfaces should publish `act` contracts and preserve ingress schema where relevant
- `published_by` should point to the actor that published the contract
- publication time should be recorded

For actor-facing ingress, a useful normalized table is:

```text
actor_sense_schemas
- sense_schema_id
- actor_id
- sense_schema_name
- primitive
- request_schema
- response_envelope_schema
- published_by
- published_timestamp
```

This lets one actor expose multiple callable surfaces without collapsing them into the top-level actor contract row.

## Why The Old Capability Table Was Dropped

The older `capabilities` table centered the wrong abstraction.

The active model is now:

- registry centers on published stimulus contracts and execution surfaces
- public actor-facing interaction goes through request/response contracts
- capability surfaces remain real, but as one typed execution-surface kind inside that broader registry
- capability ingress is preserved as ingress shape and normalized into `sense` at the receiving actor boundary

That is a better fit for the current runtime than keeping a separate capability-registry center.

## Conceptual Shape

One plausible conceptual shape is:

```text
registered_surfaces
- surface_id
- surface_name
- execution_surface_kind
- execution_surface_id
- current_contract_id
- published_by
- published_timestamp
```

This is not yet locked as final DDL.

It only captures the current storage direction at a high level.

## Short Framing

`actors` stores live runtime operators.

The registry layer should store published contracts and boundary schemas keyed to typed execution surfaces.

The exact final storage normalization remains open, but the older capability-table framing is no longer the active center.
