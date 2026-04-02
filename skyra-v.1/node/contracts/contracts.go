package node

import (
	"time"

	episode "../../episode/contracts"
	protocol "../../protocol/contracts"
	stimulus "../../stimulus/contracts"
)

type NodePurpose struct {
	Summary string `json:"summary"`
}

type NodeCapabilities struct {
	CapabilityIDs []string `json:"capability_ids,omitempty"`
}

type NodeStimulus struct {
	AcceptedTypes []string `json:"accepted_types"`
	EmittedTypes  []string `json:"emitted_types"`
}

type NodeContract struct {
	Purpose      NodePurpose      `json:"purpose"`
	Capabilities NodeCapabilities `json:"capabilities"`
	Stimulus     NodeStimulus     `json:"stimulus"`
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
	DispatchCommand(envelope protocol.CommandEnvelope) CommandDispatchResult
	WriteCommandResult(result protocol.CommandResultEvent) NodeUpdateResult
	ReceivePublishedContract(contract NodeContract) ContractPublicationResult
}
