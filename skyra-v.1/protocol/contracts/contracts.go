package protocol

import (
	"time"

	primitives "../primitives/contracts"
)

type CommandInvocation struct {
	CommandID string                   `json:"command_id"`
	NodeID    string                   `json:"node_id"`
	EpisodeID string                   `json:"episode_id"`
	IntentID  string                   `json:"intent_id,omitempty"`
	Primitive primitives.PrimitiveName `json:"primitive"`
	Args      map[string]any           `json:"args"`
	Reason    string                   `json:"reason"`
	EmittedAt time.Time                `json:"emitted_at"`
}

type CommandResultEvent struct {
	CommandID   string                   `json:"command_id"`
	NodeID      string                   `json:"node_id"`
	EpisodeID   string                   `json:"episode_id"`
	IntentID    string                   `json:"intent_id,omitempty"`
	Primitive   primitives.PrimitiveName `json:"primitive"`
	ResultKind  string                   `json:"result_kind"`
	Result      any                      `json:"result"`
	CompletedAt time.Time                `json:"completed_at"`
}
