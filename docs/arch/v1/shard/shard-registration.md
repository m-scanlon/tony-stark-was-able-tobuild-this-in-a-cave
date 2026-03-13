# Shard Registration

## Overview

Registration is not a special process. It is a sequence of heap events processed by the External Router. No separate registration service exists. The External Router gates registration — it sees the final confirmation event and writes to the agent registry.

Assumes the brain exists. Bootstrap registration (brain does not exist) is a separate process.

## Flow

### 1. Device Fingerprint

The bootstrap device walks its own hardware — CPU, GPU, RAM, storage, network interfaces, attached peripherals. It self-declares its capabilities based on what it finds.

```
device_fingerprint event → heap
```

Fields:
- `device_id` — unique hardware identifier
- `capabilities` — declared capability list with hardware evidence
- `network` — how to reach this device

### 2. Skyrad Install

The brain picks up the fingerprint event, evaluates the declared capabilities, and pushes a skyrad service package configured for this device.

Skyrad is installed on the device.

### 3. Skyrad Self-Test

Skyrad boots and self-tests each declared capability locally before reporting back.

```
capabilities_installed event → heap
```

Fields:
- `device_id`
- `capabilities` — each capability with self-test result
- `shard_agent_id` — the system agent created by skyrad on this device

Skyrad does not send this event until all self-tests pass.

### 4. Brain Verification

The External Router picks up the `capabilities_installed` event. It does not trust the self-report alone. It generates one test event per capability and puts them on the heap. The brain routes each test command back to the newly installed skyrad.

```
capability_test event (one per capability) → heap → new skyrad
```

Skyrad executes each test and responds.

### 5. Registration

Skyrad sends a `capability_test_complete` event for each passing test back to the brain.

```
capability_test_complete event → heap
```

The External Router picks up all passing `capability_test_complete` events. When all declared capabilities are confirmed, it writes the agent and its capabilities to the brain's registry.

```
agent registered → capabilities live → shard ready
```

## Event Sequence

```
device_fingerprint          → heap → brain installs skyrad
capabilities_installed      → heap → External Router generates test events
capability_test             → heap → new skyrad executes
capability_test_complete    → heap → External Router registers agent
```

## Two Rounds of Testing

Skyrad self-tests before sending `capabilities_installed`. The brain independently verifies via the External Router after. A broken install or a bad self-report does not slip through — the brain's verification round catches it before registration.

## What Gets Registered

The agent registry entry created at the end of this process:

```json
{
  "agent": "pi_living_room",
  "shard": "pi_living_room",
  "location": "living_room",
  "capabilities": ["voice", "audio_output", "lightweight_reasoning"],
  "status": "active",
  "registered_at": "2026-03-07T10:00:00Z"
}
```

Tools are registered separately by the agent after it is live — capability registration and tool registration are two distinct steps.

## Partial Registration

As long as the brain has a valid connection with skyrad, an agent is created. Whatever capabilities pass verification get registered. Capabilities that fail are omitted. Minimum viable registration is a ping — the agent exists with at least that.

Partial registration is not a degraded state. It is a valid state. The agent operates with what it has.

## Connection and Heartbeat

After registration, the brain maintains a persistent WebSocket connection with each skyrad. Skyrad is the client — it connects to the brain. The brain listens and accepts.

**Primary**: persistent WebSocket connection. Brain knows immediately when a shard drops because the connection closes.

**Fallback**: when the connection drops, the brain begins polling skyrad's HTTP `/health` endpoint with exponential backoff. Skyrad also attempts to reconnect to the brain with exponential backoff on its side.

Both sides are working the problem simultaneously.

### Shard State Machine

```
active
  → connection drops
      → brain polls /health, skyrad retries WebSocket
          → reconnection established → re-verify (see below) → active
          → polling timeout → inactive
```

Inactive is not deregistered. The agent record stays in the registry, marked inactive. The brain knows the shard exists but cannot reach it.

## Reconnection and Re-Verify

When skyrad reconnects, the brain does not restore active immediately. It puts a `shard_reconnect` event on the heap. The External Router picks it up and runs the same capability test round used during registration.

```
skyrad reconnects (WebSocket handshake)
  → brain sees known device_id
  → shard_reconnect event → heap
  → External Router generates capability_test events → skyrad executes
  → capability_test_complete events → heap
  → External Router confirms → agent marked active
```

No special reconnection flow. It is the same test mechanism as registration.

## Open Questions

- How does re-registration work when a shard's hardware changes while inactive (new GPU added, peripheral removed)? Does it run the full fingerprint flow or a delta?
