package capability

type CapabilityContract struct {
	CapabilityID      string         `json:"capability_id"`
	Name              string         `json:"name"`
	InvocationSurface string         `json:"invocation_surface,omitempty"`
	Schema            map[string]any `json:"schema,omitempty"`
	// Constraints []string `json:"constraints,omitempty"`
}
