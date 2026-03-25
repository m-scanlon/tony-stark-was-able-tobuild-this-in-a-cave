# Understanding Primitive (v1)

## Core Framing

An understanding is the smallest unit of retained experience.

It represents interpreted meaning derived from past activity, which can be reused to inform future cognition.

An understanding is:

- not raw history
- not just structure, such as entities or relationships
- not a fixed truth

It is a retained, evolving interpretation of a stimulus-grounded core.

## The Four Primitives

Every understanding is defined by four fields:

- Core
- Interpretation
- Strength
- Activation

These form the minimal contract for experience.

## 1. Core

What this understanding is about.

The core is the explicit structural substrate made relevant by the current stimulus.

Core contains:

- `entities`, which are required
- `relationships`, which are optional

Role:

- grounds the understanding
- connects it to structure
- gives interpretation a bounded subject

Important:

- core is not everything in recall
- recall may help interpret the core, but recalled items do not become part of core by default
- for now, core is grounded in entities and optional relationships

Examples:

- entities: user, project
- entities + relationship: user <-> school

## 2. Interpretation

What is understood about the core.

This is the retained meaning.

Interpretation is resolved about the current core from the episode frame:

- interaction
- recall
- cognition

Recall is not part of the stored understanding record.

It is an episode-frame input used during interpretation.

Examples:

- "User prefers iterative design"
- "This relationship is sensitive"
- "This route requires authentication"
- "Prod issues outrank cleanup"

Role:

- captures meaning derived from experience
- remains flexible in shape, rather than rigidly typed early

For now, interpretation should be stored as natural language:

- resolved, not quoted
- short
- declarative
- about the core
- not a dump of raw cognition

## 3. Strength

How strongly this understanding is held.

Strength determines how much weight the understanding carries.

Signals may include:

- confidence, or how likely it is to be true
- salience, or how important it is
- stability, or how persistent it has been
- reinforcement and recency

Role:

- differentiates weak from strong knowledge
- enables prioritization during retrieval
- allows evolution over time

## 4. Activation

When this understanding becomes relevant for retrieval.

Activation defines the conditions under which an understanding is a strong candidate to enter the episode frame.

Terminology note:

`episode frame` is used here as the working term for the active retrieval context.

The exact relationship between episode boundaries, cycles, and frames remains an open question.

Key idea:

Activation is a retrieval surface, not decision logic.

It does not force behavior. It enables recall.

Examples:

- contexts, such as wedding planning or debugging
- semantic cues, such as family or invitation
- related entities or situations

Role:

- connects intent or stimulus to experience
- keeps retrieval selective
- enables contextual recall

## Retrieval Model (High-Level)

1. Intent or stimulus arrives
2. The system derives activation cues
3. Activation cues query understandings
4. Understandings compete for retrieval
5. A small set enters the episode frame
6. Cognition uses them to reason

Not all understandings are retrieved. Selection is required.

## Functional Framing

An understanding is the retained output of an interpretive process.

Conceptually:

`interpret(core, interaction, recall, cognition) -> understanding`

In this model:

- the function produces an understanding
- interpretation is the meaning-producing step in that process
- the stored record contains the resolved result, not the full function

## Composition

Understandings are atomic at rest and composable in context.

They are:

- stored independently
- combined dynamically during an episode
- able to interact through co-activation and shared core elements

There is no rigid nesting or hierarchy required.

## Meaning vs Understanding

Meaning is the possible significance of something, and remains context-dependent.

Understanding is the system's retained interpretation of that meaning.

The system does not store meaning directly.

It stores understandings, which are:

- partial
- evolving
- context-bound

## Key Principles

- experience is built from understandings, not raw history
- understandings are grounded through core and interpretation
- strength determines weight
- activation determines retrievability
- retrieval is selective and competitive
- composition happens at runtime, not in storage

## Short Definition

An understanding is a retained interpretive object composed of:

- a core, or the stimulus-grounded structural substrate it is about
- an interpretation, or what it means
- a strength profile, or how strongly it is held
- an activation surface, or when it is relevant
