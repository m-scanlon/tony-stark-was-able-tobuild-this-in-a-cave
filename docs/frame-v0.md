# Frame v0

The frame is the full page consumed by inference.

It is the bounded in-scope runtime context for the current inference step.

## Layout

The current working frame layout is:

1. `purpose`
2. `interaction`
3. `recall`
4. `available_commands`

## Purpose

`purpose` comes from the node contract.

It provides the governing lens for how the node should interpret the current situation.

## Interaction

`interaction` captures the current external exchange.

It should contain:

- `current_stimulus`
- `recent_interaction_history`

This keeps the active thing the model may need to respond to clearly distinguished from recent context.

## Recall

`recall` contains retained artifacts currently brought into scope.

The current retained artifact families are:

- `trace`
- `understanding`
- `salience`
- `tension`

Recall is supportive and provisional.

It informs interpretation but does not by itself finalize meaning.

## Available Commands

`available_commands` tells the model what command surface it is allowed to use next.

The current first-class command examples are:

- `recall`
- `interact`

## Current Framing

The frame should be thought of as a structured inference page, not a visual screen layout.

Order and section labels matter more than any left/right or top/bottom visual metaphor.
