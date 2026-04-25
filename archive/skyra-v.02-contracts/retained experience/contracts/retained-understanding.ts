import type { RetainedArtifactBase } from "./retained-artifact-base"

// Derived meaning from traces — preserves what the experience meant
export type RetainedUnderstanding = RetainedArtifactBase & {
  kind: "understanding"
  interpretation: string        // the interpreted meaning
  source_trace_ids: string[]    // traces this understanding was derived from
}
