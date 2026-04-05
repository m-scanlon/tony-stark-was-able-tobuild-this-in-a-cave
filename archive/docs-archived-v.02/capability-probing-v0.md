# Capability Probing v0

## Purpose

This document defines the current intended flow for revealing external capability surfaces before any optional daemon installation.

The core point is:

- capability discovery is not the same thing as daemon install

## Core Framing

Capability probing should happen in stages.

The first stage should be:

- small
- generic
- safe
- evidence-driven

The output of probing should be:

- verified external capability surfaces

not:

- guessed labels with no evidence

## Main Principle

The system should:

1. gather a small generic fingerprint
2. return that fingerprint to `Stark`
3. let `Stark` classify the environment
4. select an existing probe strategy when possible
5. only synthesize a new strategy when the environment is novel or underspecified
6. execute the selected probe in a bounded way
7. preserve enough observed behavior to define a registered capability surface
8. write durable registration
9. birth or revise the relevant `stewart` actor and publish its public stimulus contract

This keeps probing:

- adaptive
- evidence-driven
- bounded
- compatible with many device types

## Probe Jobs

Probe should now be understood as doing three distinct jobs:

- discover candidate capability surfaces
- verify them through bounded invocation
- preserve enough request/response detail to support registration, primitive typing, and later `stewart` abstraction

That means probing is not only a labeling step.

It is also the first contract-formation step for a downstream capability surface.

## Why This Split Matters

Without this split, the system collapses:

- environment identification
- capability discovery
- public abstraction
- daemon installation
- runtime binding

Those should remain separate.

## Probe Stages

### 1. Bootstrap Probe

The bootstrap probe should be:

- OS-agnostic in shape
- low-risk
- lightweight
- always first

Its job is not to discover everything.

Its job is to reveal enough about the environment to choose the right next probe strategy.

### 2. Fingerprint Return To Stark

The bootstrap result should be returned to `Stark` as a small typed fingerprint package.

That package is a classification input, not the final public contract.

### 3. Stark Classification

`Stark` should classify the environment into a probe class and choose the next bounded strategy.

### 4. Probe Strategy Selection

The preferred behavior is:

- use an existing strategy first

Only if the environment is novel or confidence is low should `Stark` synthesize or research a new one.

### 5. Bounded Probe Execution

The selected probe strategy should execute under explicit runtime control.

The system should prefer:

- explicit method allowlists
- explicit targets
- explicit timeout budgets
- explicit transport restrictions

### 6. Evidence-Backed Capability Surfaces

The output of probing should be:

- capability surfaces backed by concrete evidence

Examples might include:

- `camera_input`
- `gpu_acceleration`
- `roku_ecp_endpoint`
- `bluetooth`

The probe result should preserve enough detail to support later contract publication, such as:

- capability name
- verification method
- evidence
- whether the verified surface exposes ingress, `act`, or both
- what execution surface was successfully reached
- what request shape was accepted
- what response shape came back
- confidence or status

### 7. Registration Write

Once probing has verified what is real, the system should write durable registration for the subject.

Registration remains the inventory and verification envelope.

It is not the public abstraction other actors will call.

### 8. `stewart` Birth And Public Contract Publication

Once durable registration exists, `Stark` may:

- register the verified capability surface
- register whether that surface is ingress, `act`, or both
- birth or revise the relevant `stewart` actor
- publish the simplified public request/response stimulus contract that other actors will use

This is the key split:

- capability surface = downstream execution surface
- `stewart` contract = public callable abstraction

## Relationship To Daemon Install

Daemon installation is not required for first capability revelation.

The probe flow should be able to stop successfully at lower tiers such as:

- identity-only
- transport-only
- protocol-control-only
- partial capability registration

Daemon installation is an expansion path.

It is not the prerequisite for all capability knowledge.

## Example Runtime Sequence

The current intended sequence is:

1. bootstrap probe runs on or against the subject
2. bootstrap fingerprint returns to `Stark`
3. `Stark` classifies the subject
4. `Stark` selects an existing probe strategy if one exists
5. if no strategy is good enough, `Stark` synthesizes a new bounded strategy
6. the kernel executes that strategy under explicit runtime limits
7. probe returns evidence-backed capability surfaces
8. durable registration is written
9. the relevant capability surfaces are registered
10. the relevant `stewart` actor is born or revised
11. the simplified public stimulus contract is published
12. optional daemon install may happen later if supported

## Current Design Posture

The strongest current claims are:

- capability discovery should begin with a generic bootstrap probe
- the bootstrap result is a classification input, not the final public contract
- `Stark` should prefer existing probe strategies before generating new ones
- probing should preserve enough detail to define verified downstream capability surfaces
- verified capability ingress should be preserved as ingress and normalized into `sense` at the receiving actor boundary
- verified capability outbound should be typed as `act`
- durable registration should happen before public abstraction is published
- `stewart` actors are the public abstraction layer over external capability complexity
- daemon installation is separate from first capability revelation

## Short Framing

Capability probing should begin with a small generic fingerprint.

That fingerprint returns to `Stark`, which classifies the environment, selects an existing probe strategy when possible, and only synthesizes a new one when necessary.

The resulting probe should run under explicit bounds and return evidence-backed capability surfaces.

Only after registration should the system register those surfaces and publish the simplified public `stewart` contract that the rest of the runtime will call.
