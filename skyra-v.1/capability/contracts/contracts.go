package capability

type CapabilityContract struct {
	CapabilityID     string         `json:"capability_id"`
	Name             string         `json:"name"`
	ExecutionSurface string         `json:"execution_surface,omitempty"`
	Schema           map[string]any `json:"schema,omitempty"`
}
