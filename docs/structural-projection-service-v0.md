# Structural Projection Service v0

The system needs a service that breaks source objects into entities and relationships and updates the episode field.

This service exists before recall scoring.

Its job is to turn incoming episode material into a unified structural projection that the episode field can use.

## Purpose

The structural projection service:

- reads episode-local source objects
- extracts entities and relationships
- resolves them into canonical structure when possible
- updates the episode field

This is the missing bridge between:

- interaction
- recall
- workspace

and:

- the episode field

## Inputs

The service should be able to read from:

- `interaction`
- `recall`
- `workspace`

Each source object remains distinct.

The projection service turns them into one fused structural field.

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
  source: "interaction" | "recall" | "workspace"
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
  source: "interaction" | "recall" | "workspace"
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

Interaction, recall, and workspace should not each create separate query surfaces.

They remain separate source objects, but they project into one fused structural field.

## Current Posture

The service name and exact implementation remain open.

But the architectural role is now clear:

- source objects stay separate
- the structural projection service extracts and resolves entities and relationships
- the episode field stores the fused result with per-item update records
