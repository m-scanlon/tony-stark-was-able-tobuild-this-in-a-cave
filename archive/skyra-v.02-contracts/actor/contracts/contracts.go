package actor

import (
	"time"

	episode "../../episode/contracts"
	protocol "../../protocol/contracts"
	stimulus "../../stimulus/contracts"
)

type ActorPurpose struct {
	Summary string `json:"summary"`
}

type ActorCapabilities struct {
	CapabilityIDs []string `json:"capability_ids,omitempty"`
}

type ActorStimulus struct {
	AcceptedTypes []string `json:"accepted_types"`
	EmittedTypes  []string `json:"emitted_types"`
}

type ActorContract struct {
	Purpose      ActorPurpose      `json:"purpose"`
	Commitments  []string          `json:"commitments,omitempty"`
	Capabilities ActorCapabilities `json:"capabilities"`
	Stimulus     ActorStimulus     `json:"stimulus"`
}

type ContractPublicationEvent struct {
	ActorID     string        `json:"actor_id"`
	Contract    ActorContract `json:"contract"`
	PublishedAt time.Time     `json:"published_at"`
}

type ActorEvent struct {
	Stimulus            *stimulus.StimulusEnvelope   `json:"stimulus,omitempty"`
	CommandResult       *protocol.CommandResultEvent `json:"command_result,omitempty"`
	ContractPublication *ContractPublicationEvent    `json:"contract_publication,omitempty"`
}

type ActorUpdateResult map[string]any
type CommandDispatchResult map[string]any
type ContractPublicationResult map[string]any

type ActorSubstrate interface {
	ActorID() string
	Contract() ActorContract
	ActiveEpisodeID() string
	IngestEvent(event ActorEvent) ActorUpdateResult
	ProjectFrame() *episode.Frame
	DispatchCommand(envelope protocol.StimulusEnvelope) CommandDispatchResult
	WriteCommandResult(result protocol.CommandResultEvent) ActorUpdateResult
	ReceivePublishedContract(contract ActorContract) ContractPublicationResult
}
