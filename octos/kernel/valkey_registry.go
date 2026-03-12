package kernel

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
)

const defaultSkillKeyPrefix = "skill:"

// ValkeyRegistry resolves skills from Valkey/Redis.
// Supported key format:
// - String key: skill:<tool> -> JSON object
// - Hash key:   skill:<tool> -> id/name/description/shard/tasks_json
type ValkeyRegistry struct {
	client    *ValkeyClient
	keyPrefix string
}

type skillRecord struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Shard       string      `json:"shard"`
	Tasks       []SkillTask `json:"tasks"`
}

func NewValkeyRegistry(client *ValkeyClient) *ValkeyRegistry {
	return &ValkeyRegistry{
		client:    client,
		keyPrefix: defaultSkillKeyPrefix,
	}
}

func (r *ValkeyRegistry) Get(ctx context.Context, tool string) (*Skill, error) {
	if r == nil || r.client == nil {
		return nil, fmt.Errorf("valkey registry is not configured")
	}

	tool = strings.TrimSpace(tool)
	if tool == "" {
		return nil, nil
	}

	key := r.keyPrefix + tool

	raw, err := r.client.Get(ctx, key)
	switch {
	case err == nil:
		var record skillRecord
		if err := json.Unmarshal([]byte(raw), &record); err != nil {
			return nil, fmt.Errorf("skill %q: invalid JSON value: %w", tool, err)
		}
		return normalizeSkillRecord(tool, &record), nil
	case err == ErrValkeyNil:
		// not a string key or missing key; fall through to hash lookup
	case strings.Contains(err.Error(), "WRONGTYPE"):
		// hash key, continue
	default:
		return nil, fmt.Errorf("skill %q: GET failed: %w", tool, err)
	}

	fields, err := r.client.HGetAll(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("skill %q: HGETALL failed: %w", tool, err)
	}
	if len(fields) == 0 {
		return nil, nil
	}

	record := &skillRecord{
		ID:          strings.TrimSpace(fields["id"]),
		Name:        strings.TrimSpace(fields["name"]),
		Description: strings.TrimSpace(fields["description"]),
		Shard:       strings.TrimSpace(fields["shard"]),
	}

	tasksRaw := strings.TrimSpace(fields["tasks_json"])
	if tasksRaw == "" {
		tasksRaw = strings.TrimSpace(fields["tasks"])
	}
	if tasksRaw != "" {
		if err := json.Unmarshal([]byte(tasksRaw), &record.Tasks); err != nil {
			return nil, fmt.Errorf("skill %q: invalid tasks JSON: %w", tool, err)
		}
	}

	return normalizeSkillRecord(tool, record), nil
}

func normalizeSkillRecord(tool string, record *skillRecord) *Skill {
	skillID := strings.TrimSpace(record.ID)
	if skillID == "" {
		skillID = tool
	}

	name := strings.TrimSpace(record.Name)
	if name == "" {
		name = tool
	}

	tasks := record.Tasks
	if len(tasks) == 0 {
		tasks = []SkillTask{
			{Name: tool, Description: "default task"},
		}
	}

	return &Skill{
		ID:          skillID,
		Name:        name,
		Description: strings.TrimSpace(record.Description),
		Shard:       strings.TrimSpace(record.Shard),
		Tasks:       tasks,
	}
}
