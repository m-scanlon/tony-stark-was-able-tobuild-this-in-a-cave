import type { RetainedArtifactBase } from "./retained-artifact-base"

export type RetainedSalience = RetainedArtifactBase & {
  kind: "salience"
  signal: string
  source_trace_ids: string[]
}
