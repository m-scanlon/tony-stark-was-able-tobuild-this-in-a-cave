import type { RetainedArtifactBase } from "./retained-artifact-base"

export type RetainedTension = RetainedArtifactBase & {
  kind: "tension"
  unresolved: string
  source_trace_ids: string[]
}
