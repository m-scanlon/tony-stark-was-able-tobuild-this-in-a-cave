package stimulus

type StimulusSource struct {
	NodeID       string `json:"node_id,omitempty"`
	CapabilityID string `json:"capability_id,omitempty"`
}

type StimulusEnvelope struct {
	StimulusType string         `json:"stimulus_type"`
	Source       StimulusSource `json:"source"`
	Payload      map[string]any `json:"payload"`
}

type StimulusType struct {
	TypeID      string         `json:"type_id"`
	Description string         `json:"description,omitempty"`
	Schema      map[string]any `json:"schema,omitempty"`
}

type StimulusRegistry struct {
	Version string         `json:"version,omitempty"`
	Types   []StimulusType `json:"types"`
}
