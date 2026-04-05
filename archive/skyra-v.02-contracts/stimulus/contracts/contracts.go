package stimulus

type StimulusSource struct {
	ActorID      string `json:"actor_id,omitempty"`
	CapabilityID string `json:"capability_id,omitempty"`
}

type StimulusResponseStatus string

const (
	StimulusResponseStatusSuccess  StimulusResponseStatus = "success"
	StimulusResponseStatusFailed   StimulusResponseStatus = "failed"
	StimulusResponseStatusTimedOut StimulusResponseStatus = "timed_out"
)

type StimulusEnvelope struct {
	StimulusType string         `json:"stimulus_type"`
	Source       StimulusSource `json:"source"`
	Payload      map[string]any `json:"payload"`
}

// StimulusResponseEnvelope is the shared public response shape for callable
// stimulus surfaces. Actors may extend the payload, but `status` and `reason`
// remain the required common fields.
type StimulusResponseEnvelope struct {
	Status  StimulusResponseStatus `json:"status"`
	Reason  string                 `json:"reason"`
	Payload map[string]any         `json:"payload,omitempty"`
}

type StimulusType struct {
	TypeID      string         `json:"type_id"`
	Name        string         `json:"name"`
	Description string         `json:"description,omitempty"`
	Schema      map[string]any `json:"schema,omitempty"`
}

type StimulusRegistry struct {
	Version string         `json:"version,omitempty"`
	Types   []StimulusType `json:"types"`
}
