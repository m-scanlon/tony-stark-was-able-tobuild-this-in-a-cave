package open

import node ".."

type OpenContractStatus string

const (
	OpenContractStatusIntentionallyOpen OpenContractStatus = "intentionally_open"
)

type OpenNodeContract struct {
	NodeType        string             `json:"node_type,omitempty"`
	Purpose         node.NodePurpose   `json:"purpose"`
	Stimulus        node.NodeStimulus  `json:"stimulus"`
	Cognition       node.NodeCognition `json:"cognition"`
	Commands        node.NodeCommands  `json:"commands"`
	LearningEnabled bool               `json:"learning_enabled,omitempty"`
	Status          OpenContractStatus `json:"status"`
	OpenReason      string             `json:"open_reason,omitempty"`
	OpenEdges       []string           `json:"open_edges"`
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
