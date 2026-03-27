# Structural Projection Service v0

## Core Framing

The system needs a service that breaks episode-local source objects into entities and relationships and updates the episode field.

This service exists before recall scoring.

Its job is to turn bounded episode material into a unified structural projection that the episode field can use.

## Purpose

The structural projection service:

- reads episode-local source objects
- extracts entities and relationships
- resolves them into canonical structure when possible
- updates the episode field

This is the missing bridge between:

- interaction history
- in-scope recall
- runtime artifacts when present

and:

- the episode field

## Inputs

The current `v1`-compatible source objects are:

- `interaction`
- `recall`
- `runtime_artifact`

These sources remain distinct.

The projection service turns them into one fused structural field inside the episode.

Notes:

- `interaction` is the main source in `v1`
- `recall` may also project structure back into the episode field
- `runtime_artifact` is a valid later source when transient runtime outputs become structurally useful

## Flow

The current working flow is:

1. read a source object
2. extract candidate entities and relationships
3. resolve or normalize them into structure refs
4. produce field updates
5. apply those updates to the episode field

## Output

The output is not a separate long-lived memory object.

The output is an update to the episode field.

## Projection Contract

```ts
type StructuralProjection = {
  source: "interaction" | "recall" | "runtime_artifact"
  source_id: string
  entities: ProjectedEntity[]
  relationships: ProjectedRelationship[]
  timestamp: string
}
```

```ts
type ProjectedEntity = {
  entity_id: string
  delta: number
}
```

```ts
type ProjectedRelationship = {
  relationship_id: string
  from_entity_id: string
  to_entity_id: string
  delta: number
}
```

## Episode Field Update Record

Each entity and relationship in the episode field should keep its own update record.

```ts
type FieldUpdate = {
  source: "interaction" | "recall" | "runtime_artifact"
  source_id: string
  delta: number
  timestamp: string
}
```

## Episode Field State

```ts
type EntityFieldState = {
  entity_id: string
  score: number
  updates: FieldUpdate[]
}
```

```ts
type RelationshipFieldState = {
  relationship_id: string
  from_entity_id: string
  to_entity_id: string
  score: number
  updates: FieldUpdate[]
}
```

## Design Principle

The episode field should remain unified.

Interaction, recall, and runtime artifacts should not each create separate query surfaces.

They remain separate source objects, but they project into one fused structural field.

## Current Posture

The service name and exact implementation remain open.

But the architectural role is now clear:

- source objects stay separate
- the structural projection service extracts and resolves entities and relationships
- the episode field stores the fused result with per-item update records
- `workspace` is not a canonical `v1` source object for this service
