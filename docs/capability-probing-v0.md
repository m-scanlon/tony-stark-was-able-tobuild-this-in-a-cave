# Capability Probing v0

## Purpose

This document defines the current intended flow for revealing device capabilities before daemon installation.

The core point is:

- capability discovery is not the same thing as daemon install

The system should be able to learn meaningful capability surfaces even when it cannot install a local runtime daemon on the target.

## Core Framing

Capability probing should happen in stages.

The first stage should be small, generic, and safe.

Later stages may become OS-specific or device-specific.

The output of probing should be:

- evidence-backed capability claims

not:

- guessed capability labels based only on device category

## Main Principle

The system should:

1. gather a small generic fingerprint
2. send that fingerprint back to the structural authority
3. let `Stark` classify the environment
4. select an existing probe strategy when one already exists
5. only research or synthesize a new probe strategy when the environment is novel or underspecified
6. execute the selected probe in a bounded way
7. shape initial capability contracts from observed probe behavior
8. defer capability contract publication until after registration has been written

This keeps probing:

- adaptive
- evidence-driven
- bounded
- compatible with many device types

Probe should now be understood as doing three distinct jobs:

- discover candidate capabilities
- verify them through bounded invocation
- shape the initial contract from observed behavior

That means probing is not only a labeling step.

It is also the first contract-formation step for a capability surface.

## Why This Split Matters

Without this split, the system collapses several different concerns:

- environment identification
- capability discovery
- daemon installation
- runtime binding

Those should remain separate.

The cleaner flow is:

- bootstrap fingerprint
- classification
- strategy selection
- targeted probing
- capability publication
- optional later daemon install

## Probe Stages

## 1. Bootstrap Probe

The bootstrap probe should be:

- OS-agnostic in shape
- low-risk
- lightweight
- always first

Its job is not to discover everything.

Its job is to reveal enough about the environment to choose the right next probe strategy.

### Minimum bootstrap output

At minimum, the bootstrap probe should try to reveal:

- OS family
- OS version
- architecture
- hostname or device identifier
- CPU count
- memory if cheaply available
- disk totals if cheaply available
- obvious transport surfaces such as:
  - USB presence
  - Bluetooth presence
  - network presence

The bootstrap probe should also capture:

- what methods it used
- what data was unavailable
- what confidence the bootstrap result has

### Why bootstrap stays generic

The bootstrap probe should work across:

- macOS
- Linux
- Windows
- later other device environments

The output should be normalized enough that `Stark` can classify the environment without needing a fully custom first-stage probe for every device.

## 2. Fingerprint Return To Stark

The bootstrap result should be returned to `Stark` as a small typed fingerprint package.

That package should not yet pretend the system has a full capability contract.

It is a classification input, not the final capability truth.

Conceptually:

```ts
type BootstrapFingerprint = {
  subject_id: string
  platform_family?: string
  platform_version?: string
  architecture?: string
  host_kind?: string
  cpu_count?: number
  memory_bytes?: number
  disk_bytes?: number
  observed_transports: string[]
  evidence: ProbeEvidence[]
  confidence?: number
}
```

At this stage, Stark should be understood as the authority that types this package into the next structural stimulus class used by the runtime.

## 3. Stark Classification

`Stark` should read the bootstrap fingerprint and classify the environment into a probe class.

Examples:

- `darwin_laptop`
- `linux_server`
- `windows_desktop`
- `android_device`
- `roku_endpoint`
- `sonos_endpoint`
- `unknown_network_device`

This is the point where the system decides:

- what kind of thing this probably is
- what probe strategies are worth attempting
- what capability surfaces are plausible

## 4. Probe Strategy Selection

The preferred behavior is:

- use an existing probe strategy first

Only if the environment is novel or confidence is low should `Stark` synthesize or research a new one.

This avoids regenerating probe logic for common environments every time.

### Probe strategy examples

- `darwin_system_profiler_v1`
- `linux_proc_lspci_v1`
- `windows_powershell_cim_v1`
- `roku_ecp_probe_v1`
- `sonos_local_api_probe_v1`
- `hdmi_edid_cec_probe_v1`

Conceptually:

```ts
type ProbeStrategy = {
  strategy_id: string
  target_class: string
  allowed_methods: string[]
  expected_outputs: string[]
  timeout_ms: number
}
```

## 5. Dynamic Strategy Synthesis

If no satisfactory probe strategy already exists, `Stark` may:

- research the environment
- determine what probing methods are valid for that OS or device class
- produce a new bounded probe script or probe plan

This should be treated as:

- fallback behavior

not:

- the default path for every common system

The output of this step should still be:

- a bounded probe strategy
- not unconstrained autonomous exploration

## 6. Bounded Probe Execution

The selected probe strategy should then execute under runtime control.

The execution boundary matters.

The system should prefer:

- explicit command allowlists
- explicit target hosts
- explicit timeout budgets
- explicit transport restrictions

The system should avoid:

- unrestricted subnet scans by default
- uncontrolled shell exploration
- silent capability escalation

### Safe defaults

Good defaults include:

- probe only the local host
- probe only explicitly named remote targets
- no network-wide discovery unless the network is trusted and the user approved it
- collect evidence for each capability claim

## 7. Evidence-Backed Capability Claims

The output of probing should be:

- capability claims backed by concrete evidence

Examples:

- `camera_input` because the local OS reported camera devices
- `gpu_acceleration` because the graphics stack reported a usable compute surface
- `roku_ecp_endpoint` because the target answered on port `8060`
- `bluetooth` because the controller and device surfaces were observed

This means the probe result should preserve:

- capability name
- verification method
- evidence
- constraints

Where possible, probe should also preserve enough observed invocation behavior to support the initial capability contract, such as:

- what invocation surface was successfully reached
- what operation or operation family was exercised
- what argument shape was accepted
- what result shape came back
- confidence or status

Conceptually:

```ts
type ProbedCapability = {
  name: string
  status: "verified" | "partial" | "failed"
  verification: string
  evidence: ProbeEvidence[]
  constraints?: string[]
}
```

## 8. Capability Contract Publication

Once registration has been written from the probe result, `Stark` may publish capability contracts for the verified surfaces.

That published contract should only include:

- capability surfaces that actually passed verification

Failed or weakly supported capabilities should not be promoted as if they were fully real.

Partial registration remains valid.

The system should prefer:

- smaller true capability surfaces

over:

- inflated guessed capability surfaces

## Relationship To Daemon Install

Daemon installation is not required for first capability revelation.

The probe flow should be able to stop successfully at lower tiers such as:

- identity-only
- transport-only
- protocol-control-only
- partial capability registration

Daemon installation is an expansion path.

It is not the prerequisite for all capability knowledge.

So the correct ordering is:

1. probe
2. verify
3. register what is real
4. optionally install a daemon later if the environment supports it

## Example Runtime Sequence

The current intended sequence is:

1. bootstrap probe runs locally on the subject device
2. bootstrap fingerprint returns to `Stark`
3. `Stark` classifies the subject
4. `Stark` selects an existing probe strategy if one exists
5. if no strategy is good enough, `Stark` synthesizes a new bounded strategy
6. kernel executes that strategy under explicit runtime limits
7. probe returns evidence-backed capability results
8. `Stark` publishes a capability contract from verified surfaces
9. kernel binds that contract into runtime
10. optional daemon install may happen later if supported

## Current Design Posture

The strongest current claims are:

- capability discovery should begin with a generic bootstrap probe
- the bootstrap result is a classification input, not the final capability truth
- `Stark` should prefer existing probe strategies before generating new ones
- device- or OS-specific probing should be bounded and evidence-driven
- capability contracts should be published from verified probe results
- publication should happen after registration write, not before it
- daemon installation is separate from first capability revelation

## Open Questions

The following still need fuller design:

- exact bootstrap fingerprint schema
- exact probe strategy storage and versioning model
- exact boundary for dynamic strategy synthesis
- how signed or trusted probe strategies should be handled
- how capability confidence and expiry should be represented
- how re-probing and capability drift should work over time

## Short Framing

Capability probing should begin with a small generic fingerprint.

That fingerprint returns to `Stark`, which classifies the environment, selects an existing probe strategy when possible, and only synthesizes a new one when necessary.

The resulting probe should run under explicit bounds and return evidence-backed capability claims.

Only then should `Stark` publish the capability contract.

Daemon install is an optional later expansion path, not a prerequisite for revealing capabilities.
