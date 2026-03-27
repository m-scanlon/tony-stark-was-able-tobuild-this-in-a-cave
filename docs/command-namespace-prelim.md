# Command Namespace (Prelim)

## Core Framing

The runtime should not assume every executable operation belongs to one flat primitive menu.

The older shape:

```text
skyra <primitive> <args>
```

is too narrow for the current runtime direction.

The more flexible shape is:

```text
skyra <namespace> <command> <args>
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
skyra <namespace> <command> <args>
```

Where:

- `namespace` groups a family of commands
- `command` is the executable operation inside that namespace

## Example Direction

Examples might later look like:

```text
skyra primitive interact ...
skyra primitive recall ...
skyra loop react ...
skyra loop ooda ...
```

These are only directional examples.

They do not yet lock the final namespace set.

## Primitive As Namespace

`primitive` should now be thought of as one command namespace, not necessarily the only one.

This keeps the runtime flexible enough to support:

- atomic command execution
- bounded loop execution
- later command families without flattening them into one global surface

## Relationship To Contracts

The active node contract should later define:

- which namespaces are allowed
- which commands inside those namespaces are allowed
- what loop or execution envelopes are permitted

This means:

- the runtime substrate stays generic
- the contract bounds what the node may actually use

## Current Design Posture

The strongest current claim is:

- command execution should be organized around namespaces, not one flat primitive surface

This is a flexibility move for the node runtime.

It does not yet define:

- the final namespace list
- the final command grammar
- the exact relation between namespaces and skills

## Short Framing

The runtime command surface should move from:

- `skyra <primitive> <args>`

to:

- `skyra <namespace> <command> <args>`

`primitive` becomes one namespace among others rather than the only executable family.
