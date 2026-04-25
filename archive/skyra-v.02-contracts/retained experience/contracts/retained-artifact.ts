import type { RetainedSalience } from "./retained-salience"
import type { RetainedTension } from "./retained-tension"
import type { RetainedTrace } from "./retained-trace"
import type { RetainedUnderstanding } from "./retained-understanding"

// Union of all retained artifact types — trace is factual, the rest are derived consequences
export type RetainedArtifact =
  | RetainedTrace
  | RetainedUnderstanding
  | RetainedSalience
  | RetainedTension
