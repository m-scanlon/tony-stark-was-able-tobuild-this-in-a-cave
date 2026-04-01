# Device Registration v0

## Purpose

This document defines the current `v0` registration shape for discovered devices and hosts.

The goal is not to make all devices look identical.

The goal is to give the system one strict registration envelope that can hold:

- what the subject is
- how the system sees it
- how it was probed
- what capabilities have actually been verified

## Core Framing

There is no single universal probe abstraction that makes:

- laptops
- phones
- speakers
- USB drives
- displays
- network appliances

all behave the same way.

That means the registration contract should be strict at the envelope level and flexible inside the device-specific details.

The current backbone is:

- `subject`
- `transport`
- `probe_strategy`
- `verified_capabilities`

Runtime should also carry:

- `registration_state`
- `last_verified_at`

## Main Principle

The system should standardize:

- registration shape
- verification boundary
- capability recording

It should not prematurely standardize:

- one fake universal device ontology
- one fake universal probe path
- one fake universal OS abstraction

The stable rule is:

- every registered subject has one typed registration envelope
- every capability inside that envelope must be verification-backed

## Registration Envelope

Conceptually:

```ts
type DeviceRegistration = {
  subject: SubjectRegistration
  transport: TransportRegistration
  probe_strategy: ProbeStrategyRegistration
  verified_capabilities: VerifiedCapability[]
  registration_state: "active" | "partial" | "offline" | "failed"
  last_verified_at: string
}
```

## 1. Subject

`subject` identifies what is being registered.

It should answer:

- what this thing is
- how the system refers to it
- enough stable identity to track it over time

Conceptually:

```ts
type SubjectRegistration = {
  subject_id: string
  subject_kind: "self_hosted" | "network_device" | "peripheral" | "mobile" | "unknown"
  display_name?: string
  identity?: Record<string, unknown>
}
```

### Subject Rules

`subject_id` should be stable enough for re-probe and re-registration.

`subject_kind` is a broad runtime class, not a full ontology.

`identity` may carry transport-specific or device-specific fields such as:

- hostname
- serial-like identifiers
- vendor/product identifiers
- service identity
- friendly name

`identity` should not be forced into one fake universal device schema.

## 2. Transport

`transport` records how the system currently sees or reaches the subject.

Conceptually:

```ts
type TransportRegistration = {
  kind: string
  attachment?: "local" | "network" | "proxied"
  details?: Record<string, unknown>
}
```

Examples of transport kinds:

- `local_os`
- `usb`
- `bluetooth`
- `wifi`
- `ethernet`
- `mdns_service`
- `adb`
- `hdmi_cec`

### Transport Rules

Transport is part of the registration truth because discovery and probing depend on it.

The same subject may later be reachable through more than one transport.

`details` may therefore include transport-specific facts such as:

- USB vendor/product identifiers
- local network address
- service type
- port
- bus path
- proxy agent identity

## 3. Probe Strategy

`probe_strategy` records how the system decided to inspect this subject.

Conceptually:

```ts
type ProbeStrategyRegistration = {
  strategy_id: string
  version?: string
  confidence?: "high" | "medium" | "low"
}
```

Examples:

- `darwin_system_profiler_v1`
- `usb_bootstrap_v0`
- `roku_ecp_probe_v1`
- `mdns_then_http_probe_v0`
- `android_adb_probe_v1`

### Why Probe Strategy Belongs Here

The registration record should preserve provenance.

That means the system should be able to answer:

- why the device was classified the way it was
- what probe path was used
- how much trust to place in the current registration

This does not replace detailed evidence storage.

It preserves the top-level probing path that produced the current registered state.

## 4. Verified Capabilities

`verified_capabilities` is the core payload of the registration.

These are the capability surfaces the system currently believes are usable because they passed some verification path.

Conceptually:

```ts
type VerifiedCapability = {
  capability_id: string
  name: string
  kind?: "input" | "output" | "compute" | "storage" | "sensor" | "network" | "other"
  status: "verified" | "partial" | "revoked"
  interface?: string
  constraints?: string[]
  evidence_summary?: string
}
```

Examples:

- `local_compute`
- `display_output`
- `camera_input`
- `bluetooth_scan`
- `usb_host`
- `network_access`
- `audio_output`

### Capability Rules

A capability belongs in the registration only if it is backed by:

- probe output
- verification output
- or another explicit trusted evidence path

Capabilities should not be registered only because:

- the device category makes them likely
- the model guessed them
- the OS usually has them

The registration should record what is currently verified, not what is merely plausible.

## Registration State

`registration_state` expresses the overall health of the registration record.

Suggested meanings:

- `active`: the subject is reachable and the current capability set is usable
- `partial`: the subject is known but only some capability surfaces are currently verified
- `offline`: the subject was known previously but is not currently reachable
- `failed`: registration or re-verification failed in a material way

This keeps top-level state separate from per-capability status.

## Last Verified At

`last_verified_at` records when the current registration envelope was last confirmed.

This matters because:

- capability truth can decay
- transports can disappear
- devices can change behavior
- a previously valid registration can become stale

The system should therefore treat registration as a living record rather than a one-time declaration.

## Relationship To Capability Contracts

The registration record is not the same thing as the capability contract.

The split is:

- `device_registration` says what subject is known, how it was seen, how it was probed, and what capabilities are currently verified
- capability contracts define the callable interface and constraints for each usable capability surface

In other words:

- registration is the inventory and verification envelope
- capability contracts are the callable command surface

## Relationship To Probing

The expected flow is:

1. discover or attach to a subject
2. record transport-level identity
3. select a probe strategy
4. execute bounded probing
5. shape initial capability contracts from observed probe behavior
6. write the registration envelope
7. publish capability contracts for verified surfaces
8. in `v1`, immediately follow successful registration with `birth_node`

This keeps:

- discovery
- probing
- registration
- capability publication
- node birth

as separate concerns.

## Example Shape

Conceptually:

```json
{
  "subject": {
    "subject_id": "Michaels-MacBook-Pro-10.local",
    "subject_kind": "self_hosted",
    "display_name": "MacBook Pro",
    "identity": {
      "hostname": "Michaels-MacBook-Pro-10.local",
      "platform_family": "Darwin",
      "architecture": "arm64"
    }
  },
  "transport": {
    "kind": "local_os",
    "attachment": "local",
    "details": {
      "machine_model": "MacBookPro18,1"
    }
  },
  "probe_strategy": {
    "strategy_id": "darwin_system_profiler_v1",
    "confidence": "high"
  },
  "verified_capabilities": [
    {
      "capability_id": "cap_local_compute",
      "name": "local_compute",
      "kind": "compute",
      "status": "verified",
      "interface": "compute.execute",
      "constraints": [],
      "evidence_summary": "Host introspection verified Apple M1 Pro compute resources."
    },
    {
      "capability_id": "cap_display_output",
      "name": "display_output",
      "kind": "output",
      "status": "verified",
      "interface": "display.render",
      "constraints": [],
      "evidence_summary": "Display subsystem detected through local system profiling."
    }
  ],
  "registration_state": "active",
  "last_verified_at": "2026-03-30T04:08:28Z"
}
```

## Current Design Posture

The strongest current claims are:

- there should be one strict registration envelope
- the envelope should center on `subject`, `transport`, `probe_strategy`, and `verified_capabilities`
- device-specific variation should live inside typed details and capability objects rather than exploding the top-level schema
- only verified capabilities belong in the registration record
- registration and capability contracts are related but distinct

## Short Framing

The registration contract should not pretend every device is the same.

It should provide one strict outer shape that records:

- what the subject is
- how Skyra sees it
- how it was probed
- what capabilities are currently verified

That gives the system a stable registration layer without forcing a fake universal device model.
