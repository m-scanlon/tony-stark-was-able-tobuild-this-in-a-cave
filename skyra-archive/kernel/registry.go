package kernel

import (
	"context"
	"encoding/json"
)

// Skill is the definition retrieved from the registry.
// The registry entry is intentionally small: identity plus an opaque payload blob.
type Skill struct {
	ID          string
	Name        string
	Description string
	Payload     json.RawMessage
}

// SkillRegistry is the trust boundary.
// Only skills present in the registry are executable.
type SkillRegistry interface {
	Get(ctx context.Context, skill string) (*Skill, error)
}
