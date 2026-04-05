package episode

import (
	"time"
)

type EpisodeScope string

const (
	EpisodeScopeNode   EpisodeScope = "actor"
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
	Purpose     EpisodePurpose   `json:"purpose"`
	Interaction FrameInteraction `json:"interaction"`
	Recall      EpisodeRecall    `json:"recall"`
}

type Episode struct {
	EpisodeID          string             `json:"episode_id"`
	EpisodeScope       EpisodeScope       `json:"episode_scope"`
	ActorID            string             `json:"actor_id,omitempty"`
	IntentID           string             `json:"intent_id,omitempty"`
	ActorContractID    string             `json:"actor_contract_id,omitempty"`
	Purpose            EpisodePurpose     `json:"purpose"`
	InteractionHistory InteractionHistory `json:"interaction_history"`
	Recall             EpisodeRecall      `json:"recall"`
	OpenedAt           time.Time          `json:"opened_at"`
	UpdatedAt          time.Time          `json:"updated_at"`
	ClosedAt           *time.Time         `json:"closed_at,omitempty"`
}
