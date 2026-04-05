# Skyra v.1 Contract Notes

This file is the top-level review note for the current `skyra-v.1` contract surface.

It records the main contract questions that still appear open after the stimulus-first cleanup in `docs/`.

These are review questions, not locked changes.

## Stable Enough Right Now

The following pieces look directionally stable enough:

- primitive names: `recall`, `learn`, `observe`, `act`
- shared response envelope shape in `stimulus/contracts`
- registration envelope direction in `registration/contracts`

## Main Contract Questions

### 1. Protocol Envelope

`protocol/contracts` still models the kernel-facing envelope as:

- `calling_actor`
- `command`

Questions:

- should `protocol.StimulusEnvelope` stop carrying a raw `command` string?
- should it instead carry explicit typed fields such as primitive, target execution surface, and emitted stimulus payload?
- if runtime is stimulus-first, what is the minimum kernel envelope shape?

### 2. Returned Runtime Events

`protocol/contracts` still defines `CommandResultEvent`.

Questions:

- should `CommandResultEvent` become a returned-stimulus event?
- should it become a response-envelope event?
- does the kernel need a separate result event at all, or should returned stimulus be the only runtime return object?

### 3. Actor Contract Shape

`actor/contracts` still centers:

- `Capabilities`
- `Stimulus.AcceptedTypes`
- `Stimulus.EmittedTypes`

Questions:

- should the active actor contract move fully to:
  - `purpose`
  - `commitments`
  - request stimuli
  - response envelopes
- if so, what is the exact `v1` object shape?
- should callable surfaces be modeled directly instead of as accepted/emitted type lists?

### 4. Actor Runtime Interface

`actor/contracts` still uses:

- `CommandResult`
- `DispatchCommand(...)`
- `WriteCommandResult(...)`

Questions:

- should `ActorEvent` stop modeling `CommandResult` directly?
- should the substrate interface rename command-centric methods to stimulus-centric ones?
- what is the correct runtime return path for response envelopes or returned stimulus?

### 5. ExecutionSurface Modeling

The docs now treat `ExecutionSurface` as a first-class typed routing concept, but `skyra-v.1` does not model it cleanly yet.

Questions:

- where should the typed `ExecutionSurface` contract live?
- should it be a shared contract under `stimulus/contracts`, `protocol/contracts`, or a separate top-level contract family?
- should the initial kinds remain exactly:
  - `actor`
  - `capability`

### 6. Stimulus Contract Shape

`stimulus/contracts` currently gives `StimulusType`:

- `type_id`
- `name`
- `description`
- `schema`

Questions:

- should `StimulusType` also carry `ExecutionSurface`?
- should it distinguish request vs response role?
- should request/response pairing be explicit in the contract model?
- should callable public surfaces be modeled as one record or as paired request/response records?

### 7. Stimulus Instance Shape

`stimulus/contracts.StimulusEnvelope` currently has:

- `stimulus_type`
- `source`
- `payload`

Questions:

- should emitted stimulus instances also carry the immediate `ExecutionSurface` explicitly?
- should source and target both live on the stimulus instance?
- should the instance carry enough information for kernel routing without depending on a command string?

### 8. Primitive Invocation Contracts

`protocol/primitives/contracts` still uses command-centric naming such as:

- `CommandID`
- `PrimitiveResultEvent`

Questions:

- should these be renamed to stimulus-centric identifiers?
- is `TargetActor` still the right field, or should that become an execution-surface reference?
- should primitive invocation contracts stay separate once the outer stimulus contract is finalized?

### 9. `act` Contract Shape

`ActArgs` still assumes:

- `target`
- `content`
- `modality`
- `timestamp`

Questions:

- should `act` stay modeled as fixed args?
- or should `act` be understood purely as emitting published stimulus toward an execution surface?
- if the latter, does `ActArgs` remain useful at all?

### 10. Episode Scope

`episode/contracts` still includes:

- `episode_scope = "actor" | "intent"`

Questions:

- should live `intent` episode scope remain part of the contract?
- or should intent become a reconstructed continuity layer instead of a live episode contract boundary?

### 11. Capability Contract Role

`capability/contracts` still exists as its own top-level family.

Questions:

- should `CapabilityContract` survive as a distinct contract family?
- or should capability surfaces be absorbed into the broader stimulus/execution-surface contract model?
- if it survives, what is the exact difference between a capability contract and a published capability execution surface?

### 12. Probe Payloads

`stimulus/contracts/payloads.go` still returns `CapabilityContract` inside probe output.

Questions:

- should probe results keep returning `CapabilityContract`?
- or should they return a richer published-surface record aligned to the stimulus registry?
- what exact contract should a `stewart` learn from when it abstracts a discovered capability?

## Suggested Order If These Change

If the contract surface is revised, the likely order is:

1. `protocol/contracts`
2. `actor/contracts`
3. `stimulus/contracts`
4. `protocol/primitives/contracts`
5. `episode/contracts`
6. `capability/contracts`

## Short Framing

The biggest unresolved contract gap is that the docs now center:

- stimulus-first runtime
- typed `ExecutionSurface`
- request stimulus + response envelope

but the implementation contracts are still partly organized around:

- command strings
- command results
- accepted/emitted stimulus lists

That is the main contract drift to resolve next.
