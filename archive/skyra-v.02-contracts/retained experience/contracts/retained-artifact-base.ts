import type { AnchorSet } from "./anchor-set"

export type RetainedArtifactKind =
  | "trace"
  | "understanding"
  | "salience"
  | "tension"

export type RetainedArtifactBase = {
  id: string
  kind: RetainedArtifactKind
  anchor_set: AnchorSet
  context_artifact_ids?: string[]
}
