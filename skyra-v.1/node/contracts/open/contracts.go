package open

import node ".."

type OpenContractStatus string

const (
	OpenContractStatusIntentionallyOpen OpenContractStatus = "intentionally_open"
)

type OpenNodeContract struct {
	Purpose      node.NodePurpose      `json:"purpose"`
	Capabilities node.NodeCapabilities `json:"capabilities"`
	Stimulus     node.NodeStimulus     `json:"stimulus"`
	Status       OpenContractStatus    `json:"status"`
	OpenReason   string                `json:"open_reason,omitempty"`
	OpenEdges    []string              `json:"open_edges"`
}

// JarvisNodeContract remains intentionally open until the root-manager
// shape and jarvis-lineage clone contract are defined.
type JarvisNodeContract struct {
	OpenNodeContract
}

// StarkNodeContract remains intentionally open until the root-manager
// shape and system-node lineage contract are defined.
type StarkNodeContract struct {
	OpenNodeContract
}
