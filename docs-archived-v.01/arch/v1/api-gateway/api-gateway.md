# API Gateway (v1)

## Overview

The API Gateway is a transport adapter in front of the standalone kernel service.

Its v1 responsibility is simple:

- receive inbound commands or stimuli
- normalize transport details
- forward canonical runtime input to the kernel
- return outbound responses from the kernel side

There is no canonical transport envelope in v1.

## Flow

```text
external input
  -> ingress adapter
  -> kernel standalone service
  -> kernel max heap
  -> Chain of Thought / Human-to-Machine Interaction
  -> egress adapter
  -> external output
```

## Transport Rule

The gateway should not add hidden planning structures around the request.

It forwards canonical input into the kernel runtime. The kernel owns queueing, attention, perception, and primitive progression.

## Relationship To The Kernel

The kernel is the canonical runtime boundary.

The gateway is optional transport and compatibility infrastructure around that boundary.

## Not Canonical For v1

The following ideas are not part of the current v1 contract:

- transport envelopes
- routing manifests
- orchestration lifecycle contracts
- skill-resolution-heavy gateway behavior
