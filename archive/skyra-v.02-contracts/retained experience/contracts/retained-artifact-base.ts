import type { AnchorSet } from "./anchor-set"

export type RetainedArtifactKind =
  | "trace"
  | "understanding"
  | "salience"
  | "tension"

// Shared base shape for all retained artifacts — carries the anchor into canonical structure
// and optional references to prior artifacts that influenced this artifact's formation
export type RetainedArtifactBase = {
  id: string
  kind: RetainedArtifactKind
  anchor_set: AnchorSet
  context_artifact_ids?: string[]
}
