package kernel

import "context"

// Skill is the definition retrieved from the registry.
type Skill struct {
	ID          string
	Name        string
	Description string
	Shard       string
	Tasks       []SkillTask
}

// SkillTask is one step in a skill's roadmap.
type SkillTask struct {
	Name        string
	Description string
	Args        []string
}

// SkillRegistry is the trust boundary.
// Only skills present in the registry are executable.
type SkillRegistry interface {
	Get(ctx context.Context, tool string) (*Skill, error)
}
