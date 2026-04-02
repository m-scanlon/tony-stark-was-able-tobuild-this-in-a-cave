# Database Schema v0

## Scope

This document currently locks:

- the `actors` table
- the `capabilities` table
- the capability contract shape referenced by `capabilities.current_contract_id`

It does **not** lock the rest of the database schema yet.

At this stage, the following remain intentionally open:

- actor contract table shape
- capability contract storage / versioning table shape
- delegation / invocation edge tables
- validated invocation tables
- schema registry tables
- result / event persistence tables

The goal here is to preserve stable identity records before expanding outward into contracts, validation, and persistence.

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

Examples might later include names like `jarvis`, `stark`, or a more specific live actor name.

Rules:

- should be unique
- should be the name used when the actor appears in the command string

### `status`

The actor's current lifecycle state.

The field is locked.

The exact status enum may still be refined later.

### `current_contract_id`

Identifier pointing at the actor's currently active contract.

This field is locked even though the contract table itself is not yet locked.

For now, treat it as a reference by identifier only.

### `created_by`

Records how the actor came into existence.

Rules:

- `created_by = bootstrap` for actors present at runtime startup
- otherwise `created_by` must be a valid `actor_id`

This gives the system immediate lineage without requiring a separate lineage table.

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

## Actor Non-Goals

The `actors` table does **not** currently lock:

- actor type
- actor class
- actor purpose
- actor origin as a separate field
- contract versioning

At this stage:

- purpose belongs in the contract, not in the actor row
- bootstrap lineage is represented through `created_by = bootstrap`

## Conceptual SQL Shape

This is a conceptual shape only.

It is not yet the final DDL for the whole database.

```sql
CREATE TABLE actor (
  actor_id TEXT PRIMARY KEY,
  actor_name TEXT NOT NULL UNIQUE,
  status TEXT NOT NULL,
  current_contract_id TEXT NOT NULL,
  created_by TEXT NOT NULL,
  birth_timestamp TIMESTAMPTZ NOT NULL
);
```

## Actor Short Framing

`actors` is one of the first locked database tables.

It stores live actor identity, current contract pointer, lineage through `created_by`, and birth time.

The rest of the database schema remains open for now.

## `capabilities`

The current locked capability table is:

```text
capabilities
- capability_id
- capability_name
- status
- current_contract_id
- published_by
- published_timestamp
```

## Capability Field Meanings

### `capability_id`

Stable internal identifier for the capability surface.

Rules:

- must be unique
- must be immutable
- must not be reused for a different capability surface later

### `capability_name`

Semantic name for the capability surface.

Examples might later include names like `local_compute`, `display_output`, or `roku_ecp_endpoint`.

Rules:

- names the callable surface in human / system terms
- does not need to be globally unique across all runtime contexts
- uniqueness should live in `capability_id`, not in `capability_name`

### `status`

The capability's current lifecycle or availability state.

The field is locked.

The exact status enum may still be refined later.

The important current rule is:

- this status should describe a published callable surface
- not a merely guessed or plausible capability

### `current_contract_id`

Identifier pointing at the capability's currently active contract.

This field is locked.

It points to a capability contract with the canonical shape from
`skyra-v.1/capability/contracts/contracts.go`.

That shape is exactly:

```text
CapabilityContract
- CapabilityID
- Name
- ExecutionSurface
- Schema
```

The storage table for capability contracts is still open.

### `published_by`

Records which actor published this capability surface into the runtime.

Rules:

- must be a valid `actor_id`
- should identify the publishing actor, not a generic sentinel like `bootstrap`

In the current design direction, this will usually be `Stark`.

### `published_timestamp`

Timestamp recording when the capability surface was published into the runtime.

## Capability Locked Rules

The following rules are now locked for `capabilities`:

- every capability row has one stable `capability_id`
- every capability row has one `capability_name`
- every capability row has one `status`
- every capability row points to one `current_contract_id`
- every capability row records `published_by`
- every capability row records `published_timestamp`
- `published_by` must be a valid `actor_id`
- a capability row represents a published callable surface, not a raw guessed capability claim
- every capability row points to a capability contract with exactly this shape:
  - `CapabilityID`
  - `Name`
  - `ExecutionSurface`
  - `Schema`

## Capability Boundary

The capability row is **not** the same thing as the registration envelope.

The current split is:

- registration says what subject is known, how it was seen, how it was probed, and what capabilities are currently verified
- the capability row identifies one published callable surface in the runtime
- the capability contract defines the callable surface for that published capability through:
  - `CapabilityID`
  - `Name`
  - `ExecutionSurface`
  - `Schema`

This means capability publication should happen only after verification-backed registration has been written.

## Capability Non-Goals

The `capabilities` table does **not** currently lock:

- capability versioning
- probe evidence storage
- probe strategy storage
- registration envelope storage
- transport modeling
- subject / device binding fields

At this stage:

- this table is for published callable surfaces
- capability contract storage remains a separate layer
- evidence and registration remain separate layers

## Conceptual SQL Shape

This is a conceptual shape only.

It is not yet the final DDL for the whole database.

```sql
CREATE TABLE capabilities (
  capability_id TEXT PRIMARY KEY,
  capability_name TEXT NOT NULL,
  status TEXT NOT NULL,
  current_contract_id TEXT NOT NULL,
  published_by TEXT NOT NULL,
  published_timestamp TIMESTAMPTZ NOT NULL
);
```

## Short Framing

`actors` stores live runtime operators.

`capabilities` stores published callable surfaces.

`capabilities.current_contract_id` points to the canonical four-field capability contract shape.

The storage tables behind contracts remain open.
