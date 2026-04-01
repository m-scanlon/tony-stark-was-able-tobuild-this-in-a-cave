package episode

import (
	"time"

	primitives "../../protocol/primitives/contracts"
)

type EpisodeScope string

const (
	EpisodeScopeNode   EpisodeScope = "node"
	EpisodeScopeIntent EpisodeScope = "intent"
)

type EpisodePurpose struct {
	Text string `json:"text"`
}

type EpisodeRecall struct {
	RetainedArtifactIDs []string `json:"retained_artifact_ids"`
}

type InteractionEvent map[string]any

type InteractionHistory struct {
	Events []InteractionEvent `json:"events"`
}

type FrameInteraction struct {
	CurrentStimulus          InteractionEvent   `json:"current_stimulus,omitempty"`
	RecentInteractionHistory []InteractionEvent `json:"recent_interaction_history"`
}

type Frame struct {
	Purpose           EpisodePurpose             `json:"purpose"`
	Interaction       FrameInteraction           `json:"interaction"`
	Recall            EpisodeRecall              `json:"recall"`
	AvailableCommands []primitives.PrimitiveName `json:"available_commands"`
}

type Episode struct {
	EpisodeID          string                     `json:"episode_id"`
	EpisodeScope       EpisodeScope               `json:"episode_scope"`
	NodeID             string                     `json:"node_id,omitempty"`
	IntentID           string                     `json:"intent_id,omitempty"`
	NodeContractID     string                     `json:"node_contract_id,omitempty"`
	Purpose            EpisodePurpose             `json:"purpose"`
	InteractionHistory InteractionHistory         `json:"interaction_history"`
	Recall             EpisodeRecall              `json:"recall"`
	AvailableCommands  []primitives.PrimitiveName `json:"available_commands"`
	OpenedAt           time.Time                  `json:"opened_at"`
	UpdatedAt          time.Time                  `json:"updated_at"`
	ClosedAt           *time.Time                 `json:"closed_at,omitempty"`
}
