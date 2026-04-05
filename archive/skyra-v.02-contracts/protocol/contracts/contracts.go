package protocol

import (
	"time"

	primitives "../primitives/contracts"
)

// StimulusEnvelope is the base kernel-facing dispatch contract.
// The command string names the target actor; the envelope names the caller.
// dependencyLedger / dependency shape is intentionally left open for now.
// If that runtime dependency model is promoted into contract space, it will
// most likely live as a sub-object of this kernel envelope rather than as a
// separate top-level contract.
type StimulusEnvelope struct {
	CallingActor string `json:"calling_actor"`
	Command      string `json:"command"`
}

type CommandResultEvent struct {
	CommandID    string                   `json:"command_id"`
	CallingActor string                   `json:"calling_actor,omitempty"`
	TargetActor  string                   `json:"target_actor"`
	EpisodeID    string                   `json:"episode_id,omitempty"`
	IntentID     string                   `json:"intent_id,omitempty"`
	Primitive    primitives.PrimitiveName `json:"primitive"`
	ResultKind   string                   `json:"result_kind"`
	Result       any                      `json:"result"`
	CompletedAt  time.Time                `json:"completed_at"`
}
