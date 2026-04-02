package protocol

import (
	"time"

	primitives "../primitives/contracts"
)

// CommandEnvelope is the minimal kernel-facing dispatch wrapper.
// The command string names the target actor; the envelope names the caller.
type CommandEnvelope struct {
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
