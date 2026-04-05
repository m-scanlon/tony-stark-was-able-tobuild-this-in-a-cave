import type { RetainedArtifactBase } from "./retained-artifact-base"

export type RetainedTrace = RetainedArtifactBase & {
  kind: "trace"
  happened: string
  source_episode_ids: string[]
}
