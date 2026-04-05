package open

import actor ".."

type OpenContractStatus string

const (
	OpenContractStatusIntentionallyOpen OpenContractStatus = "intentionally_open"
)

type OpenActorContract struct {
	Purpose      actor.ActorPurpose      `json:"purpose"`
	Commitments  []string                `json:"commitments,omitempty"`
	Capabilities actor.ActorCapabilities `json:"capabilities"`
	Stimulus     actor.ActorStimulus     `json:"stimulus"`
	Status       OpenContractStatus      `json:"status"`
	OpenReason   string                  `json:"open_reason,omitempty"`
	OpenEdges    []string                `json:"open_edges"`
}

// JarvisActorContract remains intentionally open until the root-manager
// shape and jarvis-lineage clone contract are defined.
type JarvisActorContract struct {
	OpenActorContract
}

// StarkActorContract remains intentionally open until the root-manager
// shape and system-actor lineage contract are defined.
type StarkActorContract struct {
	OpenActorContract
}
