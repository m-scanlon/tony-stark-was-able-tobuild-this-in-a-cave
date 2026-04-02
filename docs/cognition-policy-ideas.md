# Cognition Policy Ideas

## Core Framing

This is an ideas document, not canon.

The current thought is that the contract may define a cognition policy, and runtime may instantiate a cognition object from that policy.

The important boundary is:

- the cognition object is not the source of truth
- the episode remains the bounded runtime state container
- the node substrate remains the execution surface

So cognition here should be thought of as a runtime helper or execution-policy object layered on top of the existing runtime model.

## Main Idea

The contract may eventually define cognition policy such as:

- `free_reasoning`
- `react`
- `ooda`

Conceptually:

```ts
type CognitionPolicy = {
  mode: "free_reasoning" | "react" | "ooda"
  max_steps?: number
  stop_conditions?: string[]
}
```

Then runtime might instantiate something like:

```ts
const cog = new CognitionRun(contract.cognition)
```

## Why This Fits The Current Runtime

This works because the runtime was already designed so that:

- the substrate is generic
- the contract bounds the allowed execution envelope
- the node process is event-driven
- the episode owns state
- the frame is projected from the episode

That means a cognition object can sit on top of the runtime without replacing the underlying runtime model.

## The Safer Shape

The safer shape is:

- cognition policy lives in or near the contract
- runtime instantiates a cognition helper from that policy
- cognition helper emits observations, commands, and stop/continue decisions
- the substrate/kernel still owns dispatch and writeback

So the cognition object should not execute side effects directly.

It should operate over bounded runtime objects and emit commands into the runtime.

## Example Direction

Conceptually:

```ts
class CognitionRun {
  constructor(policy: CognitionPolicy) {}

  addObservation(event: NodeEvent) {}
  addCommand(cmd: CommandEnvelope) {}
  addResult(result: CommandResultEvent) {}

  shouldStop(): boolean {}
}
```

And a command added by cognition might look like:

```ts
cog.addCommand({
  calling_actor: "jarvis",
  command: 'skyra jarvis act -target human -content "the current frame requires a user-facing response" -modality text -timestamp now -reason "the current frame requires a user-facing response"'
})
```

The important detail is:

- cognition adds or emits the command
- the runtime still dispatches it

## Why Not Make Cognition The Source Of Truth

If cognition became the source of truth, the model would start collapsing:

- runtime policy
- episode state
- command execution
- loop style

into one object.

That is exactly what the current node/episode/frame split is trying to avoid.

So the cleaner posture is:

- cognition object = helper
- episode = state
- substrate = runtime surface
- kernel = execution authority

## Current Usefulness

This idea is useful because it gives a concrete OOD shape for:

- contract-bounded `ReAct`
- contract-bounded `OODA`
- free reasoning

without forcing the whole runtime to become one fixed loop.

## Short Framing

The contract may eventually define a cognition policy.

Runtime may instantiate a cognition helper from that policy.

That helper can collect observations, commands, and results while the episode remains the source of truth and the substrate/kernel remain the execution surface.
