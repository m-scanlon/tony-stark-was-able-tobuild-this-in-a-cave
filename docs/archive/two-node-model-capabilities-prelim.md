# Two-Node Model & Capabilities (Prelim)

## Core Structure

The system is organized around two foundational nodes:

- Life — determines user-facing meaning, attentiveness, and what matters
- System — handles execution, routing, capabilities, and structure

Both share the same contract primitives, such as purpose, stimulus, and interact, but differ in their domain of concern.

## Devices as Capabilities

Devices are not modeled as nodes.

They are registered as capability surfaces within the system, such as:

- microphone
- speaker
- display
- APIs

These capabilities can be used by any node that is permitted to interact with them.

## Provisioned Nodes

The system may provision additional worker nodes to handle specific responsibilities, such as:

- email
- calendar
- homework
- retrieval

These nodes are:

- logical
- purpose-driven
- not bound to any single device

They operate through the system and can use any compatible capability surface.

## Flow

- Stimulus enters through a system capability, such as a device or API
- Life determines whether the stimulus is meaningful and requires response
- Life requests action
- System executes the response using available nodes and capabilities

## Core Principle

Life decides what matters; System decides what runs.

Nodes are functional; devices are capabilities.
