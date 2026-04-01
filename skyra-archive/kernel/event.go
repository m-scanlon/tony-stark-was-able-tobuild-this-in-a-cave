package kernel

import "time"

type EventKind string

const (
	EventKindStimulus   EventKind = "stimulus"
	EventKindExperience EventKind = "experience"
	EventKindUnderstand EventKind = "understand"
	EventKindReference  EventKind = "reference"
	EventKindInfer      EventKind = "infer"
	EventKindResolve    EventKind = "resolve"
	EventKindInteract   EventKind = "interact"
)

type Event struct {
	ID           string    `json:"id"`
	StimulusID   string    `json:"stimulus_id,omitempty"`
	Source       string    `json:"source"`
	Kind         EventKind `json:"kind"`
	Priority     int       `json:"priority"`
	ChainID      string    `json:"chain_id,omitempty"`
	ChainVersion int       `json:"chain_version,omitempty"`
	CreatedAt    time.Time `json:"created_at"`
	sequence     uint64
}

type Stimulus struct {
	ID                string    `json:"id"`
	Type              string    `json:"type"`
	Source            string    `json:"source"`
	RawPayload        string    `json:"raw_payload"`
	NormalizedSummary string    `json:"normalized_summary"`
	Timestamp         time.Time `json:"timestamp"`
}

type HistoryEntry struct {
	ID        string    `json:"id"`
	Source    string    `json:"source"`
	Role      string    `json:"role"`
	Content   string    `json:"content"`
	Timestamp time.Time `json:"timestamp"`
}

type Understanding struct {
	ID        string    `json:"id"`
	Text      string    `json:"text"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Perception struct {
	ID            string         `json:"id"`
	History       []HistoryEntry `json:"history"`
	Stimulus      Stimulus       `json:"stimulus"`
	Understanding *Understanding `json:"understanding,omitempty"`
	UpdatedAt     time.Time      `json:"updated_at"`
}

type ChainState struct {
	ID               string     `json:"id"`
	Version          int        `json:"version"`
	Status           string     `json:"status"`
	RootStimulusID   string     `json:"root_stimulus_id"`
	StimulusSource   string     `json:"stimulus_source"`
	ReviewingUser    bool       `json:"reviewing_user,omitempty"`
	CurrentPrimitive string     `json:"current_primitive,omitempty"`
	Perception       Perception `json:"perception"`
	ReferenceOutput  string     `json:"reference_output,omitempty"`
	InferOutput      string     `json:"infer_output,omitempty"`
	StepCount        int        `json:"step_count"`
	ResolveCount     int        `json:"resolve_count"`
	InferLoopCount   int        `json:"infer_loop_count"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type InteractionMessage struct {
	ID         string    `json:"id"`
	Role       string    `json:"role"`
	Content    string    `json:"content"`
	StimulusID string    `json:"stimulus_id,omitempty"`
	CreatedAt  time.Time `json:"created_at"`
}

type ThoughtStepState struct {
	ChainID         string     `json:"chain_id"`
	StepID          string     `json:"step_id"`
	StepIndex       int        `json:"step_index"`
	Frame           string     `json:"frame"`
	Status          string     `json:"status"`
	PrimitiveChoice string     `json:"primitive_choice,omitempty"`
	RawOutput       string     `json:"raw_output"`
	StartedAt       time.Time  `json:"started_at"`
	CompletedAt     *time.Time `json:"completed_at,omitempty"`
}

type ThoughtStepRecord struct {
	ChainID              string     `json:"chain_id"`
	StepID               string     `json:"step_id"`
	StepIndex            int        `json:"step_index"`
	Frame                string     `json:"frame"`
	Status               string     `json:"status"`
	PrimitiveChoice      string     `json:"primitive_choice,omitempty"`
	RawOutput            string     `json:"raw_output"`
	PerceptionSnapshot   Perception `json:"perception_snapshot"`
	StartedAt            time.Time  `json:"started_at"`
	CompletedAt          time.Time  `json:"completed_at"`
	InteractionMessage   string     `json:"interaction_message,omitempty"`
	InteractionChannel   string     `json:"interaction_channel,omitempty"`
	CurrentUnderstanding string     `json:"current_understanding,omitempty"`
}
