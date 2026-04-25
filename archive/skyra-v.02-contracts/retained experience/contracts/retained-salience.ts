import type { RetainedArtifactBase } from "./retained-artifact-base"

// Preserves what carries weight or attention — what mattered
export type RetainedSalience = RetainedArtifactBase & {
  kind: "salience"
  signal: string                // the salient signal
  source_trace_ids: string[]    // traces this salience was derived from
}
