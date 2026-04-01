import type { RetainedSalience } from "./retained-salience"
import type { RetainedTension } from "./retained-tension"
import type { RetainedTrace } from "./retained-trace"
import type { RetainedUnderstanding } from "./retained-understanding"

export type RetainedArtifact =
  | RetainedTrace
  | RetainedUnderstanding
  | RetainedSalience
  | RetainedTension
