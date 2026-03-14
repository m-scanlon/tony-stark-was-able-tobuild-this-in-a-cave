## Problem

The kernel needs a top-level routing contract for every skill.

`compute_requirements` is part of that contract, but the obvious versions of it
are unstable:

- a scalar like `1-10`
- labels like `small | medium | large`
- labels like `fast_reasoning | deep_reasoning`

Those encode today's hardware and model landscape into the skill contract. They
will age badly.

The contract should survive:

- better local models
- cheaper hardware
- new accelerators
- larger context windows
- different runtime stacks

The kernel should route from live shard capabilities. The skill should declare
what it needs, not how "powerful" the machine should feel in a specific year.

## Core Position

`compute_requirements` should be constraint-based, not scale-based.

The skill contract should describe:

- hard requirements
- soft preferences
- stable execution properties

The router should decide placement using:

- current shard capabilities
- current model/runtime properties
- current load
- locality / affinity
- availability

## Model As Dependency

The model is a dependency, not the contract.

That means:

- a skill should not encode a frozen notion of model size or model prestige
- a skill should not say "needs a 32B model"
- a skill should declare the properties it needs from the execution environment

Examples of stable needs:

- minimum context window
- audio input support
- image input support
- tool use support
- network access
- accelerator availability
- low-latency response

The active model and runtime on a shard satisfy those needs or they do not.

This stays consistent with the existing principle that the model is a dependency
chosen underneath the system, even if trust and approval remain model-scoped in
other parts of the design.

## First Primitive: Inference Type

The first routing primitive for compute should be:

- `inference_type`

This answers the first kernel question:

- what kind of inference is this skill asking for?

Initial values:

- `text_generation`
- `embedding`
- `classification`
- `speech_to_text`
- `text_to_speech`

Later:

- `vision`

This is better than a compute scale because it describes the actual kind of
runtime the skill needs. The kernel can first route by inference type, then
apply the rest of the constraints like context window, latency, accelerator,
and locality.

## What Belongs In The Skill Contract

These are good candidates for `compute_requirements`:

- `inference_type`
  - first routing primitive for model/runtime compatibility

## Skills

- a class is not a skill
- a skill is a function
- a class emerges when a set of prompts repeatedly operates on the same state nodes
- a skill with no execution history has unbounded predicted complexity
- until history exists, the router should not pretend it has a stable cost profile hence should be routed to the highest model

## Primitive Outputs

- `primitive_call`
  - an instruction to execute a registered runtime primitive or skill
  - includes the primitive name and arguments
  - this is how the model asks the kernel to perform an action

- `artifact_ref`
  - a reference to a persistent object produced or stored by the system
  - examples: file, dataset, code artifact, media
  - points to something that exists outside the current execution

- `reply`
  - a human-facing message intended for the user
  - natural language output that terminates the interaction or communicates results

- `return_value`
  - a structured value returned to the caller skill after successful execution
  - used for programmatic composition between skills

- `error`
  - a signal that execution failed or could not proceed
  - includes error information so the caller or runtime can decide how to recover

## Input Primitves

Skills accept open inputs but must emit bounded outputs.
The runtime does not control the world’s input space; it controls the executable output space.

## Cost Derived From Inputs

1. Without history, complexity is unbounded in prediction.

If models execute one line at a time, authored complexity is probably not a
useful contract field.

Better:

- skill contract declares `inference_type`
- kernel derives cost from invocation inputs

Working shape:

```text
C = T_in + T_reason + T_out
```

- `T_in`
  - time to ingest, encode, or prepare the input
- `T_reason`
  - time spent doing the core inference work
- `T_out`
  - time to generate, decode, or emit the output

Initial input map:

- `text_generation`
  - `prompt_tokens`
  - `max_output_tokens`
  - `context_window_required`
- `embedding`
  - `item_count`
  - `tokens_per_item`
- `classification`
  - `item_count`
  - `tokens_per_item`
  - `label_count` if needed
- `speech_to_text`
  - `audio_duration_ms`
  - `sample_rate`
  - maybe `streaming | batch`
- `text_to_speech`
  - `text_length`
  - maybe `voice_id`
  - maybe `streaming | batch`

## Idea To Consider: Time Complexity

It may be useful to think about routing cost in rough time-complexity terms,
not as a contract field, but as a kernel-side heuristic.

Example:

- `embedding` may behave roughly like `O(n)` over items
- `classification` may behave like `O(n)` over items
- `speech_to_text` may behave roughly like `O(audio_duration)`
- `text_to_speech` may behave roughly like `O(text_length)`
- `text_generation` is more complicated because prompt cost and generation cost
  are different phases

This may help the router reason about growth behavior without freezing a fake
compute scale into the skill contract.
