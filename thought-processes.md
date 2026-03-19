# Thought Processes

Before the system defines a protocol, it defines the outer shape of system
flow.

The top-level execution primitive set is:

- `stimuli`
- `process`
- `interact`

A primitive must participate in execution.
These three do:

- `stimuli` starts execution by bringing something into the system
- `process` advances execution by transforming that input
- `interact` completes execution by emitting a bounded result

## `stimuli`

What enters the system.

Examples:

- user text
- user voice
- artifacts
- shard events
- scheduled triggers

## `process`

What the system does with incoming stimuli.
`process` is the parent transformation layer.
`thought` is a child form of `process`.

### `thought`

Thought is one internal process.
It is the shaping step that turns stimuli into bounded next actions or outputs.

Examples:

- infer intent
- retrieve context
- validate boundaries
- decide next action

#### `repair`

Repair is a subprocess of `thought`.

Other process work can include:

- transform state

### Candidate Processes

- `stimuli reasoning selection`

## `interact`

What leaves the system in bounded form.
`interact` is the parent outward layer.
`agree` and `reply` are child forms of `interact`.

### `agree`

Agreement is an interaction state.
It may be carried by a `reply`, but it is not the same as `reply`.

### `reply`

Reply is an interaction form.
It carries bounded outward response.

Examples:

- `agree`
- `reply`
- `primitive_call`
- `artifact_ref`
- `return_value`
- `error`

## Why This Sits Above The Protocol

The protocol is not the process itself.
The protocol is the executable encoding of these top-level forms.

That means:

- the top-level primitives participate directly in execution
- the top-level set is stable even when implementations change
- `process` is the parent layer
- `thought` is one child process within it
- native protocol primitives are secondary forms under this layer
- the runtime consumes contracts, not introspection
- the world sees bounded interactions, not internal reasoning

## Relationship To Native Protocol

`docs/arch/v1/native-protocol/native-protocol.md` defines how these forms
become executable protocol surfaces.
