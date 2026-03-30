# Command Set (Prelim)

## Core Framing

The runtime should not assume every executable operation belongs to one flat primitive menu.

Earlier docs used the word `namespace`.

The current working protocol term is `command_set`.

The older shape:

```text
skyra <primitive> -<args>
```

is too narrow for the current runtime direction.

The more flexible shape is:

```text
skyra <command_set> <command> -<args>
```

## Why This Matters

The runtime may eventually need to support more than one executable command family.

Examples:

- atomic primitives
- bounded loop schemas
- later higher-order command families

If everything is forced into one flat primitive surface too early, the runtime becomes harder to extend cleanly.

## Current Direction

The command line shape should be thought of as:

```text
skyra <command_set> <command> -<args> -reason "<why this command is being emitted>"
```

Where:

- `command_set` groups a family of commands
- `command` is the executable operation inside that command set
- `-reason` is the mandatory audit-trail rationale carried on every emitted command

Nodes do not act directly on users, APIs, or the runtime.

They emit commands.

Commands without `-reason` should be treated as invalid by runtime validation.

## Example Direction

Examples might later look like:

```text
skyra primitive interact -reason "the current frame requires a user-facing response"
skyra primitive recall -reason "the current stimulus introduced structural cues worth recall lookup"
skyra loop react -reason "the contract allows bounded multi-step execution for this episode"
skyra loop ooda -reason "the contract selected an ooda-style command family for this situation"
```

These are only directional examples.

They do not yet lock the final command-set vocabulary or argument grammar.

## Working Recall Command Shape

One useful working shape for recall is:

```text
skyra primitive recall \
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
skyra primitive recall \
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

## Primitive As Command Set

`primitive` should now be thought of as one command set, not necessarily the only one.

This keeps the runtime flexible enough to support:

- atomic command execution
- bounded loop execution
- later command families without flattening them into one global surface

## Relationship To Contracts

The active node contract should define:

- which command sets are allowed
- which commands inside those command sets are allowed
- what loop or execution envelopes are permitted

This means:

- the runtime substrate stays generic
- the contract bounds what the node may actually use

## Current Design Posture

The strongest current claim is:

- command execution should be organized around command sets, not one flat primitive surface

This is a flexibility move for the node runtime.

It does not yet define:

- the final command-set vocabulary
- the final command grammar
- the exact relation between command sets and skills

It does define one important audit rule:

- command emission requires explicit rationale

## Short Framing

The runtime command surface should move from:

- `skyra <primitive> -<args>`

to:

- `skyra <command_set> <command> -<args> -reason "<why this command is being emitted>"`

`primitive` becomes one command set among others rather than the only executable family.

`-reason` is mandatory because it is part of the audit trail, not optional commentary.
