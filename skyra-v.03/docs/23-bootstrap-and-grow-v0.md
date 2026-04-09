# 23 — Bootstrap and the Grow Being

## The Bootstrap Problem

The kernel enforces a hard invariant: no being can interact with the system unless it already exists in the hashmap. This creates a chicken-and-egg problem — something has to be in the hashmap before anything can happen.

The solution is a single hardcoded bootstrap command.

## skyra born

`skyra born` is the one command that bypasses routing. The kernel handles it as a special case before normal routing rules apply. It directly instantiates a single being: `grow`.

From that point on, everything goes through the protocol.

## The Grow Being

`grow` is the bridge between protocol strings and the underlying Go runtime. It is non-cognitive. It has its own syntax for being creation, relationship seeding, and any other operation that needs to touch the runtime directly.

Skyra talks to `grow`. `grow` does the instantiation.

## The Hard Rule

`grow` is the only path to runtime instantiation post-bootstrap. No other being can call `NewBeing` directly. If it doesn't go through `grow`, it doesn't happen.

- `skyra born` creates `grow`
- `grow` creates everything else
- Nothing else touches the constructors

## What This Means

The genome (`genome.skyra`) is not read by the kernel directly. It flows through `grow`. You write the genome, you run `skyra born`, and from there Skyra and `grow` handle instantiation. That is the one act of creation. Everything after is the system's.
