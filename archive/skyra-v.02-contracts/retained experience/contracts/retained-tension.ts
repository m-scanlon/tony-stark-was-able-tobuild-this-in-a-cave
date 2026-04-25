import type { RetainedArtifactBase } from "./retained-artifact-base"

// Preserves what remains unresolved or conflicting — open edges in experience
export type RetainedTension = RetainedArtifactBase & {
  kind: "tension"
  unresolved: string            // the unresolved or conflicting significance
  source_trace_ids: string[]    // traces this tension was derived from
}
