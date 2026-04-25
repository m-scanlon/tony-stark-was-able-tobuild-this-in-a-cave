import type { RetainedArtifactBase } from "./retained-artifact-base"

// Factual, non-interpretive record of what occurred — the grounding artifact that
// understanding, salience, and tension are later derived from
export type RetainedTrace = RetainedArtifactBase & {
  kind: "trace"
  happened: string              // bounded natural-language rendering of the occurrence
  source_episode_ids: string[]  // episodes this trace was produced from
}
