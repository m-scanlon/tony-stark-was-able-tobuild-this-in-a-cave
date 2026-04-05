Skyra Cognitive Runtime
Product Requirements Document (PRD)

Version: v0.2
Status: Draft
Owner: Skyra Runtime Project

1. Overview

Skyra is a continuous cognitive runtime that processes incoming stimuli, maintains an active chain-of-thought state, and decides when to deepen understanding or interact with the outside world.

Skyra experiences an unbounded stream of stimuli and signals.

To keep cognition organized, Chain of Thought does not reason over the entire stream at once.

Instead, Chain of Thought maintains bounded attention over a stable slice of experience.

Skyra v1 now has a clean separation between Chain of Thought and Human-to-Machine Interaction.

In v1, the native primitives are protocol language for runtime control.

The later skill system belongs to v2-v3 and is not important for the current runtime contract.

Chain of Thought is responsible for cognition.

Human-to-Machine Interaction is responsible for outbound communication and machine-facing execution.

Chain of Thought does not directly own human I/O, machine I/O, transport layers, or delivery channels.

Human-to-Machine Interaction does not perform cognition.

In v1, Skyra does not use recursive reasoning trees, child chains, or multi-head cognitive structures.

The canonical top-level primitive menu for v1 is:

- understand
- interact

Inside understand, interpret must always be available.

Choosing interpret immediately enters reference.

After reference, control moves to infer.

From infer, the model may choose between:

- reference
- resolve

The model loops between reference and infer until it chooses resolve.

When the interpret cycle finishes with resolve, that result becomes understanding and is written into chain-of-thought perception.

2. Goals

The system should:

Continuously ingest stimuli from the environment.

Decide whether stimuli are relevant enough to process.

Promote relevant stimuli into deeper reasoning.

Maintain an active chain-of-thought perception for the active runtime.

Maintain an unbounded stimulus stream while keeping bounded attention for active reasoning.

Keep a strict architectural separation between cognition and Human-to-Machine Interaction.

Expose two separate frontend screens in v1:

- one screen for interacting with Skyra
- one screen for watching Skyra's internal chain of thought

These frontend surfaces must remain separate screens.

These frontend surfaces must speak to the same backend.

Support exactly two top-level primitives in v1:

- understand
- interact

Ensure interpret is always available within understand.

Allow interpret to enter reference immediately, then move between reference and infer until resolve is chosen.

Write understanding back into chain-of-thought perception only when the interpret cycle finishes with resolve.

Keep runtime state in memory for the live runtime only.

Support introspection and debugging through OpenTelemetry tracing.

3. Non-Goals (v1)

Skyra v1 will not attempt to:

Build recursive cognitive trees.

Spawn child reasoning chains.

Maintain multi-head reasoning structures.

Implement a full DAG dependency system.

Support distributed execution across machines.

Collapse Chain of Thought and Human-to-Machine Interaction into one service boundary.

Treat the later skill system or later orchestration model as part of v1.

Require runtime state to survive restarts.

These may appear in later versions, but they are explicitly out of scope for v1.

4. Core Concept

Skyra operates as a continuous loop over an unbounded stimulus stream, but cognition and interaction are separate runtime concerns.

Ingress path:

stimulus
  ->
ingress shard
  ->
stimulus stream append
  ->
Nexus
  ->
kernel max heap
  ->
Chain of Thought

Interaction path:

Chain of Thought
  ->
interaction intent
  ->
Human-to-Machine Interaction
  ->
outbound message or machine action
  ->
resulting stimulus
  ->
Nexus

Every incoming stimulus or signal is first appended to the stimulus stream.

Nexus decides whether that newly appended stimulus should also produce a heap event for active cognition.

Chain of Thought reasons over chain-of-thought perception by holding bounded attention on a stable slice of experience.

The kernel uses a max heap, so higher numeric priority executes first.

The canonical v1 priority policy is:

- user messages: priority 100
- internal chain-of-thought work: priority 50

In v1, a user message always interrupts active internal reasoning.

Skyra must stop, attend to the user message, and inspect it through the kernel before deciding what to do next.

That inspection does not force a user-facing response.

Skyra may choose to emit `interact`, or she may choose not to respond and continue internal reasoning.

In v1, the active frame always contains perception.

Perception is the always-present object in the frame.

Perception always contains:

- history
- stimulus

understanding does not belong in perception until the interpret cycle has finished with resolve and written it.

Skyra's primitive choices do not belong in the frame.

Primitive choices are runtime control outputs that steer Chain of Thought from one step to the next.

history means past events fetched into the current reasoning cycle.

stimulus means the current stimulus being attended to in the active reasoning cycle.

Only resolve creates understanding.

At the start of a reasoning cycle, the active working perception must contain history plus stimulus.

When the interpret cycle finishes with resolve, it writes understanding inside perception.

New stimuli may continue arriving into the stimulus stream while the current attention slice remains stable.

When a user message arrives during internal reasoning, the active chain is suspended rather than discarded.

The suspended chain is resumable state, not stale state.

5. Runtime Model

The system pipeline is:

stimulus
  ->
ingress shard
  ->
stimulus stream append
  ->
Nexus
  ->
strength / ingress filtering
  ->
kernel max heap
  ->
Chain of Thought
  ->
select or resume an attention slice
  ->
load perception from history + stimulus
  ->
select top-level primitive
  ->
understand or interact
  ->
if understand:
  interpret
  ->
  reference <-> infer loop
  ->
  resolve
  ->
  update understanding in perception
  ->
  if interact:
  emit interaction intent
  ->
  Human-to-Machine Interaction executes delivery
  ->
  resolution

User interrupt path:

active internal step
  ->
user message enters kernel at priority 100
  ->
kernel interrupts active model step
  ->
active chain becomes suspended
  ->
singleton perception attends to the user message
  ->
Skyra evaluates the user message
  ->
interact or no outward response
  ->
if a suspended chain exists, resume in rebase mode

Rebase mode means:

- keep the singleton perception object
- keep the latest history
- keep the latest understanding
- restore the suspended chain's prior focus as `stimulus`
- clear transient `reference` and `infer` artifacts
- restart that suspended chain at `reference`

In this model, the stimulus stream is unbounded, but active attention is bounded.

The heap decides what asks for attention next.

6. Core Runtime Objects

6.1 Chain of Thought

Chain of Thought is the top-level cognitive service.

Responsibilities:

Receive promoted stimuli from the kernel heap.

Load and maintain chain-of-thought perception.

Hold bounded attention on a stable experience slice while reasoning.

Shift attention in an organized way when higher-priority work arrives.

Suspend and later resume internal reasoning when user messages interrupt it.

Select between the two top-level primitives:

- understand
- interact

Run the understand flow.

Ensure interpret is available inside understand.

Allow interpret to enter reference immediately, then allow infer to choose between reference and resolve until resolve is chosen.

When a suspended chain resumes after a user interrupt, restart it at `reference` so the new user message can be incorporated through the updated perception.

Write understanding into perception after resolve.

Emit interaction intents when interact is chosen.

Own the prompt template lifecycle for cognitive frames.

Keep only live runtime state for the current active session.

Chain of Thought does not directly send messages to humans or invoke machine-facing channels.

6.2 Human-to-Machine Interaction

Human-to-Machine Interaction is a separate runtime boundary from Chain of Thought.

Responsibilities:

Receive interaction intents from Chain of Thought.

Translate those intents into outbound messages, UI actions, tool calls, or machine actions.

Collect replies, results, and follow-on signals from the outside world.

Convert those external results back into stimuli that re-enter through Nexus.

Human-to-Machine Interaction does not perform reasoning and does not directly mutate chain-of-thought perception outside the normal stimulus path.

6.3 Stimulus

Stimulus represents any incoming signal.

Examples:

user text

voice input

file changes

API responses

tool results

internal runtime events

For priority purposes in v1:

- user-originated messages enter the heap at priority 100
- internal chain-of-thought events enter the heap at priority 50

Stimulus Fields

stimulus_id
stimulus_type
stimulus_source
raw_payload
timestamp
normalized_summary

stimulus_source identifies the origin of the incoming stimulus or signal.

6.4 Experience and Chain-of-Thought Perception

In v1, the stimulus stream is an append-only, unbounded stream of stimuli and signals.

Experience is the live cognitive record formed from the stimulus stream and retained history during the active runtime.

Experience is not a recursive actor tree.

The stimulus stream continuously grows as Skyra receives new input from the world or from internal runtime activity.

Experience is the source material from which history is retrieved and attention is formed.

Chain-of-thought perception is the active working view used by Chain of Thought during one reasoning cycle.

Chain-of-thought perception is not the whole stimulus stream.

It is the bounded attention frame over that stream.

Chain-of-thought perception is composed of:

- history
- stimulus

In v1, there is exactly one mutable chain-of-thought perception object.

That perception may be updated over time, but the runtime should not create multiple competing perceptions.

understanding is optional and is absent until the interpret cycle has finished with resolve and written it.

history

Retrieved past events relevant to the current attention slice.

stimulus

The current stimulus being actively attended to in the frame.

understanding

The understanding held in perception.

understanding is absent before the interpret cycle has finished with resolve.

When the interpret cycle finishes with resolve, understanding must be written into perception alongside history and stimulus.

Only the latest understanding is retained in perception in v1.

Attention rules:

The stimulus stream may continue growing while cognition is active.

The active attention slice must remain stable for the duration of a reasoning step.

Attention may shift only at defined execution boundaries.

Suggested perception fields:

perception_id
attention_id
attention_cursor
stimulus_ids[]
history_events[]
stimulus_summary
understanding?

perception_id should be stable in v1 because the runtime keeps one mutable perception object.
updated_at

6.5 Kernel Event

Kernel events are work items processed by the max heap.

In v1, priority is currently fixed by event class:

- user messages: 100
- internal chain-of-thought work: 50

Because the heap is a max heap, user messages are processed before internal chain-of-thought work when both are pending.

Suggested event fields:

event_id
priority
priority_class
event_kind
event_source
caused_by_stimulus_id
payload

event_source records the origin of the stimulus or signal that produced the event.

For user-originated work, event_source should preserve the original stimulus source.

For internal chain-of-thought work, event_source may be an internal runtime source such as chain_of_thought.

6.6 Interaction Intent

Interaction Intent is the handoff object from Chain of Thought to Human-to-Machine Interaction.

Suggested intent fields:

intent_id
source_stimulus_id
source_understanding_id
channel
payload
created_at

6.7 Frontend Contract

Skyra v1 exposes two separate frontend surfaces that speak to the same backend.

Interaction Surface

The interaction surface is only for human-to-Skyra interaction.

It should show only the human-facing exchange.

It should not show internal chain-of-thought frames or primitive choices.

Chain-of-Thought Surface

The chain-of-thought surface is a separate read-only frontend.

It should show Skyra's internal chain of thought in chronological order.

It should show the full raw internal output.

It should also show a structured view of:

- the current perception
- the currently running frame or template
- the primitive choice made at each step

Primitive choices may appear on this frontend, but they are runtime outputs and are not part of the frame itself.

The chain-of-thought surface should show only the latest understanding carried in perception.

Per-Step Presentation

The chain-of-thought surface should render one chronological sequence of internal steps.

Each running step should expose a stable step identifier so the frontend can show what Skyra is currently doing while output is still being generated.

Suggested active step identifiers:

- step_id
- step_index
- step_status

step_id is the unique identifier for the current step.

step_index is the chronological position of that step within the active reasoning cycle.

step_status should indicate whether the step is:

- queued
- streaming
- completed

Live Step Streaming

While a step is actively running, the chain-of-thought surface may show streamed raw deltas associated with that active step identifier.

This allows the observer to see Skyra typing or generating output live.

Completed Step Record

When the active step finishes, the streamed output should collapse into one completed chronological step record.

That completed step record is the durable thought event shown in the step stream.

Each displayed step should show together:

- the raw internal output for that step
- the active frame or template for that step
- the primitive choice produced at that step
- the current perception at that step

This allows the observer to see both the raw chain of thought and the structured runtime state side by side.

v1 should therefore support this UI behavior:

- live raw typing during an active step
- one finalized chronological event after the step completes

Visual Direction

The interaction surface and the chain-of-thought surface should be visually distinct.

Session Scope

v1 supports one active Skyra session.

7. Primitive Model

7.1 Top-Level Primitive Menu

The top-level primitive menu for v1 is exactly:

- understand
- interact

No additional top-level primitives are part of the v1 canonical model.

7.2 Understand

understand is the top-level cognitive primitive used to form meaning from perception.

Inside understand, interpret must always be selectable.

The understand frame menu must therefore always include:

- interpret

7.3 Interpret

interpret is the active reasoning frame inside understand.

There is no standalone interpret selector template in v1.

When understand selects interpret, execution enters reference immediately.

reference

Pull relevant past events or supporting context into the active history view.

reference always hands off to infer when it finishes.

infer

Generate or refine meaning from the combination of history and stimulus inside perception.

From infer, the model may choose between:

- reference
- resolve

resolve

End the current interpret cycle and commit the result as understanding.

The model may move from infer back to reference repeatedly until it chooses resolve.

resolve is the stop condition for interpret in v1.

7.4 Interact

interact is the top-level primitive used when the system needs to communicate or act outside cognition.

interact does not merge Human-to-Machine Interaction into Chain of Thought.

Instead, interact produces an interaction intent that Human-to-Machine Interaction executes through a separate boundary.

8. Prompt Frame Templates

Each primitive or reasoning frame uses a predefined prompt template.

Templates provide the scaffold for model behavior.

Templates are stored in a registry keyed by frame or primitive.

Chain of Thought owns the prompt template lifecycle for cognition.

The minimum template set for v1 is:

frame_templates["understand"]
frame_templates["reference"]
frame_templates["infer"]
frame_templates["resolve"]
frame_templates["interact"]

Template rules:

The understand template must always expose interpret as an available selection.

There is no standalone interpret selector template in v1.

Choosing interpret must enter reference immediately.

The reference template must hand off to infer.

The infer template must expose reference and resolve as available selections.

The resolve template must produce an understanding payload that updates perception.

The interact template must produce an interaction intent rather than directly performing the interaction.

9. Execution Flow

Execution follows this pattern:

stimulus enters stimulus stream
  ->
Nexus evaluates appended stimulus
  ->
if promoted, enqueue heap event
  ->
stimulus promoted from heap
  ->
Chain of Thought picks it up
  ->
create or resume bounded attention slice
  ->
retrieve relevant history
  ->
load working perception with history + stimulus
  ->
select top-level primitive
  ->
understand
  ->
interpret
  ->
reference
  ->
infer
  ->
reference or resolve
  ->
resolve
  ->
update understanding in perception
  ->
optional interact
  ->
emit interaction intent
  ->
Human-to-Machine Interaction performs outbound action
  ->
resulting stimulus returns through Nexus

Important execution rules:

There is no recursive child-chain spawn in v1.

There is no parent-child cognitive tree in v1.

There is no multi-head execution model in v1.

If resolve is not chosen yet, interpret remains active.

Entering interpret begins with reference.

reference always advances to infer.

infer may return to reference or continue to resolve.

When the interpret cycle completes with resolve, understanding must be immediately available in perception.

New stimuli may append to the stimulus stream at any time.

New stimuli do not mutate the active attention slice in the middle of a reasoning step.

Attention shifts may occur only at safe boundaries:

- before interpret
- after reference
- after infer
- after resolve

If higher-priority work arrives while lower-priority work is active, the runtime should save the current attention state and return to it later.

10. Telemetry and Observability

Skyra uses OpenTelemetry for runtime observability.

Telemetry must be detailed enough to reconstruct:

what stimulus arrived

why it was dropped or promoted

what was appended to the stimulus stream

what attention slice was selected

what history was retrieved

how interpret moved between reference and infer

when resolve was chosen

what understanding was written into perception

when an interaction intent was emitted

what Human-to-Machine Interaction did with that intent

what frame or template was active at each internal step

what primitive choice was made at each internal step

11. Trace Structure

A trace represents one promoted stimulus reasoning cycle or one dropped ingress cycle at Nexus.

Trace
 - nexus.ingress
 - nexus.threshold
 - nexus.drop or kernel.heap.enqueue
 - kernel.heap.pop
 - chain_of_thought.pickup
 - attention.select
 - perception.load
 - interpret.reference
 - interpret.infer
 - interpret.resolve
 - perception.update
 - interact.emit
 - hmi.dispatch
 - hmi.result_ingress

12. Core Telemetry Attributes

Telemetry should track:

Identity

trace_id
span_id

Stimulus

stimulus_id
stimulus_type
stimulus_source
threshold_score
interest_score
promotion_decision

Cognition

top_level_primitive
step_id
step_index
step_status
frame
prompt_template_id
prompt_template_version
model_id
event_priority
priority_class
event_source

Perception

perception_id
attention_id
attention_cursor
history_event_count
stimulus_count
has_understanding
understanding_id

Interaction

interaction_intent_id
interaction_channel
interaction_result_kind

Performance

latency_ms
token_input
token_output

13. Runtime Data Structures

Primary structures:

Stimulus Stream

events: AppendOnlyLog<StimulusOrSignal>

ChainOfThoughtState

state: idle | understanding | interacting
current_attention_id
perception_id
last_understanding_id
active_step_id
active_step_index
active_step_status

AttentionState

attention_id
attention_cursor
stimulus_ids[]
status: active | suspended | resolved

ThoughtStepState

step_id
step_index
frame
status: queued | streaming | completed
started_at
completed_at?
live_delta_buffer

Thought Step Log

steps: AppendOnlyLog<ThoughtStepRecord>

ThoughtStepRecord

step_id
step_index
frame
primitive_choice
raw_output
perception_snapshot
started_at
completed_at

Perception Store

perception: ChainOfThoughtPerception

History Store (tentative)

events: EventStore

Current working model:

- retrieved past events populate history in perception
- the attended stimulus populates stimulus in perception
- the backing storage shape for history is still open

Frame Templates

frame_templates: HashMap<frame, FrameTemplate>

Event Heap

events: MaxHeap<Event>

Interaction Intent Queue

interaction_intents: Queue<InteractionIntent>

Runtime State Boundary

Stimulus stream state, chain-of-thought state, attention state, perception state, history state, frame templates, heap state, and outstanding interaction intents are live runtime state in v1.

They do not need to survive restarts.

14. System Invariants

The system must maintain the following rules:

Stimuli first pass through ingress shards, append to the stimulus stream, and then reach Nexus.

Every incoming stimulus or signal must also append to the stimulus stream.

Low-strength or boring ingress is dropped at Nexus.

Dropped ingress still emits full OpenTelemetry data.

Only Nexus-promoted stimuli enter the kernel heap.

The kernel heap is a max heap ordered by event priority.

In v1, user messages are enqueued at priority 100.

In v1, internal chain-of-thought work is enqueued at priority 50.

User messages therefore outrank internal chain-of-thought work in the heap.

Every kernel event must carry event_source.

event_source must identify the origin of the stimulus or signal that produced the event.

Chain of Thought and Human-to-Machine Interaction are separate runtime boundaries.

The frontend must expose two separate screens connected to the same backend:

- an interaction screen
- an internal chain-of-thought screen

The interaction screen shows only human-facing interaction.

The internal chain-of-thought screen is read-only.

The internal chain-of-thought screen shows chronological raw internal output.

The internal chain-of-thought screen also shows structured perception and the active frame or template.

Primitive choices may be shown on the internal chain-of-thought screen, but they are not part of the frame.

The internal chain-of-thought screen presents one chronological step stream.

The internal chain-of-thought screen exposes the active step identifier while output is streaming.

The active step identifier includes step_id, step_index, and step_status.

The internal chain-of-thought screen may show live raw deltas while the active step is streaming.

When the step completes, those live deltas collapse into one finalized chronological step event.

Each displayed step should show together the raw internal output, the active frame or template, the primitive choice, and the current perception for that step.

Only the latest understanding is retained in perception in v1.

v1 supports one active Skyra session.

The two frontend surfaces should be visually distinct.

Chain of Thought never directly performs human-facing or machine-facing delivery.

Human-to-Machine Interaction never bypasses Nexus when returning results to cognition.

The top-level primitive menu is exactly:

- understand
- interact

The understand frame always has interpret available.

Entering interpret begins at reference.

reference always advances to infer.

infer may move back to reference repeatedly.

interpret stops only when resolve is selected.

Only resolve creates understanding.

Understanding must be written into chain-of-thought perception.

Chain-of-thought perception contains history and stimulus for the active reasoning cycle.

When the interpret cycle finishes with resolve, that same perception also contains understanding.

There is exactly one mutable perception object in v1.

The stimulus stream is unbounded.

Attention is bounded.

The active attention slice must remain stable during a reasoning step.

New stimuli may arrive continuously without directly mutating an in-flight reasoning step.

Recursive cognitive trees are out of scope for v1.

Child reasoning chains are out of scope for v1.

Multi-head reasoning is out of scope for v1.

15. Example Runtime Flow

Example interaction:

Stimulus: "Why did the child cry?"

Runtime:

stimulus enters ingress shard
  ->
stimulus stream append
  ->
Nexus threshold check
  ->
kernel max heap enqueue
priority: 100
  ->
Chain of Thought picks stimulus up at top of heap
  ->
create bounded attention slice
  ->
retrieve past events
history:
  child touched stove
  child withdrew hand
  child started crying
  ->
load perception with history + stimulus
  ->
select understand
  ->
select interpret
  ->
reference
  ->
infer
  ->
reference
  ->
infer
  ->
resolve
understanding:
  the child cried because touching the stove caused pain
  ->
update understanding in perception
  ->
if a response is needed, select interact
  ->
emit interaction intent to Human-to-Machine Interaction
  ->
Human-to-Machine Interaction sends the reply

16. Future Work

Future iterations may include:

recursive cognitive models after v1 if explicitly reintroduced

richer interaction frame behavior

structured history graphs

learning-based priority tuning

distributed runtime components

17. Summary

Skyra v1 is a continuous cognitive runtime built around an unbounded stimulus stream, bounded attention, retrieved history, and an active chain-of-thought perception.

The canonical v1 primitive model is now:

- understand
- interact

Inside understand, interpret is always available.

Choosing interpret enters reference immediately.

From infer, the model may loop back to reference until it selects resolve.

When the interpret cycle finishes with resolve, the resulting understanding is written into perception alongside history and stimulus.

New stimuli may continue arriving into the stimulus stream without breaking the active reasoning step because Chain of Thought reasons over bounded attention rather than the full stream.

Chain of Thought and Human-to-Machine Interaction remain separate runtime boundaries throughout the system.

The frontend is split into two separate surfaces over the same backend:

- a human interaction surface
- a read-only internal chain-of-thought surface

The chain-of-thought surface is chronological, raw, and structurally inspectable, while the interaction surface remains human-facing only.

18. Open Questions

1. History storage shape

History retrieval is required, but the exact storage model is still open. Should history be stored as an event log, a graph, or a hybrid structure?

2. Stimulus scope policy

How large can the active stimulus in perception be in v1, and what exact rules govern when new stimuli replace it versus waiting for a later attention shift?

3. Interaction intent schema

What exact fields are mandatory across all interaction channels in v1?

4. Interact selection policy

Does Chain of Thought always choose interact explicitly, or are there cases where the runtime auto-selects interact after resolve based on policy?

5. Additional priority classes

Beyond user messages at 100 and internal chain-of-thought work at 50, are any other fixed priority classes needed in v1?

6. Template lifecycle details

What exact lifecycle operations should prompt template management include: create, update, version, retire, cache, rollback?
