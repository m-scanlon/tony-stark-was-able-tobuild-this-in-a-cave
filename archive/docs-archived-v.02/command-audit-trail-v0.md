# Command Audit Trail v0

## Status

This filename is historical.

The old command-string framing is no longer the active runtime model.

The active runtime boundary is emitted `stimulus`.

## Current Audit Direction

The main audit requirement still survives:

- routed runtime execution should remain inspectable

In the current stimulus-first model, that means the system should preserve an inspectable record of:

- who emitted the stimulus
- which primitive was used
- which `ExecutionSurface` was targeted
- what request payload was sent
- what response envelope came back

## Response Envelope Requirement

The currently locked public response envelope requires:

- `status`
- `reason`

with `status` currently:

- `success`
- `failed`
- `timed_out`

That makes `reason` part of the active public audit surface even though the older `-reason` command grammar is no longer canonical.

## Request-Side Audit Fields

The exact required request-side audit fields are still open.

The older rule:

- every emitted command must carry `-reason`

should now be treated as historical shorthand rather than as the final protocol shape.

What survives from it is the intent:

- emitted runtime activity should remain inspectable
- returned outcomes should carry explicit rationale

## Current Design Posture

The strongest current claims are:

- the runtime audit surface is now stimulus-first
- response envelopes must carry `status` and `reason`
- the old command-string `-reason` grammar should not be treated as locked canon
- request-side audit detail is still open

## Short Framing

The old command audit rule has been superseded by the stimulus-first protocol.

What remains canonical is the need for inspectable emitted stimulus plus a response envelope that always carries `status` and `reason`.
