package node

import (
	"time"

	episode "../../episode/contracts"
	protocol "../../protocol/contracts"
	primitives "../../protocol/primitives/contracts"
	stimulus "../../stimulus/contracts"
)

type NodePurpose struct {
	Summary string   `json:"summary"`
	Limits  []string `json:"limits,omitempty"`
}

type NodeCapabilities struct {
	AllowedCapabilitySurfaces []string `json:"allowed_capability_surfaces"`
}

type NodeStimulus struct {
	AcceptedTypes []string `json:"accepted_types"`
	EmittedTypes  []string `json:"emitted_types"`
}

type NodeCognition struct {
	Mode           string   `json:"mode"`
	MaxSteps       int      `json:"max_steps,omitempty"`
	StopConditions []string `json:"stop_conditions,omitempty"`
}

type NodeCommands struct {
	AllowedCommands []primitives.PrimitiveName `json:"allowed_commands,omitempty"`
}

type NodeContract struct {
	NodeType        string            `json:"node_type,omitempty"`
	Purpose         NodePurpose       `json:"purpose"`
	Capabilities    *NodeCapabilities `json:"capabilities,omitempty"`
	Stimulus        NodeStimulus      `json:"stimulus"`
	Cognition       NodeCognition     `json:"cognition"`
	Commands        NodeCommands      `json:"commands"`
	LearningEnabled bool              `json:"learning_enabled,omitempty"`
}

type ContractPublicationEvent struct {
	NodeID      string       `json:"node_id"`
	Contract    NodeContract `json:"contract"`
	PublishedAt time.Time    `json:"published_at"`
}

type NodeEvent struct {
	Stimulus            *stimulus.StimulusEnvelope   `json:"stimulus,omitempty"`
	CommandResult       *protocol.CommandResultEvent `json:"command_result,omitempty"`
	ContractPublication *ContractPublicationEvent    `json:"contract_publication,omitempty"`
}

type NodeUpdateResult map[string]any
type CommandDispatchResult map[string]any
type ContractPublicationResult map[string]any

type NodeSubstrate interface {
	NodeID() string
	Contract() NodeContract
	ActiveEpisodeID() string
	IngestEvent(event NodeEvent) NodeUpdateResult
	ProjectFrame() *episode.Frame
	DispatchCommand(invocation protocol.CommandInvocation) CommandDispatchResult
	WriteCommandResult(result protocol.CommandResultEvent) NodeUpdateResult
	ReceivePublishedContract(contract NodeContract) ContractPublicationResult
}
