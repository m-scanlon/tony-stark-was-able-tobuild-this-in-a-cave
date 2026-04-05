# Two-Actor Model & Capabilities (Prelim)

## Core Structure

The system is organized around two foundational actors:

- Life — determines user-facing meaning, attentiveness, and what matters
- System — handles execution, routing, capabilities, and structure

Both share the same contract primitives, such as purpose, stimulus, and interact, but differ in their domain of concern.

## Devices as Capabilities

Devices are not modeled as actors.

They are registered as capability surfaces within the system, such as:

- microphone
- speaker
- display
- APIs

These capabilities can be used by any actor that is permitted to interact with them.

## Provisioned Actors

The system may provision additional worker actors to handle specific responsibilities, such as:

- email
- calendar
- homework
- retrieval

These actors are:

- logical
- purpose-driven
- not bound to any single device

They operate through the system and can use any compatible capability surface.

## Flow

- Stimulus enters through a system capability, such as a device or API
- Life determines whether the stimulus is meaningful and requires response
- Life requests action
- System executes the response using available actors and capabilities

## Core Principle

Life decides what matters; System decides what runs.

Actors are functional; devices are capabilities.
