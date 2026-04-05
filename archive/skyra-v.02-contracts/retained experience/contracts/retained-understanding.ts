import type { RetainedArtifactBase } from "./retained-artifact-base"

export type RetainedUnderstanding = RetainedArtifactBase & {
  kind: "understanding"
  interpretation: string
  source_trace_ids: string[]
}
