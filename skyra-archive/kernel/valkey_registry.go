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
// - String key: skill:<skill> -> JSON object
// - Hash key:   skill:<skill> -> id/name/description/payload_json
type ValkeyRegistry struct {
	client    *ValkeyClient
	keyPrefix string
}

type skillRecord struct {
	ID          string          `json:"id"`
	Name        string          `json:"name"`
	Description string          `json:"description"`
	Payload     json.RawMessage `json:"payload"`
}

func NewValkeyRegistry(client *ValkeyClient) *ValkeyRegistry {
	return &ValkeyRegistry{
		client:    client,
		keyPrefix: defaultSkillKeyPrefix,
	}
}

func (r *ValkeyRegistry) Get(ctx context.Context, skill string) (*Skill, error) {
	if r == nil || r.client == nil {
		return nil, fmt.Errorf("valkey registry is not configured")
	}

	skill = strings.TrimSpace(skill)
	if skill == "" {
		return nil, nil
	}

	key := r.keyPrefix + skill

	raw, err := r.client.Get(ctx, key)
	switch {
	case err == nil:
		var record skillRecord
		if err := json.Unmarshal([]byte(raw), &record); err != nil {
			return nil, fmt.Errorf("skill %q: invalid JSON value: %w", skill, err)
		}
		return normalizeSkillRecord(skill, &record), nil
	case err == ErrValkeyNil:
		// not a string key or missing key; fall through to hash lookup
	case strings.Contains(err.Error(), "WRONGTYPE"):
		// hash key, continue
	default:
		return nil, fmt.Errorf("skill %q: GET failed: %w", skill, err)
	}

	fields, err := r.client.HGetAll(ctx, key)
	if err != nil {
		return nil, fmt.Errorf("skill %q: HGETALL failed: %w", skill, err)
	}
	if len(fields) == 0 {
		return nil, nil
	}

	record := &skillRecord{
		ID:          strings.TrimSpace(fields["id"]),
		Name:        strings.TrimSpace(fields["name"]),
		Description: strings.TrimSpace(fields["description"]),
	}

	payloadRaw := strings.TrimSpace(fields["payload_json"])
	if payloadRaw == "" {
		payloadRaw = strings.TrimSpace(fields["payload"])
	}
	if payloadRaw != "" {
		record.Payload = normalizePayload([]byte(payloadRaw))
	}

	return normalizeSkillRecord(skill, record), nil
}

func normalizeSkillRecord(skill string, record *skillRecord) *Skill {
	skillID := strings.TrimSpace(record.ID)
	if skillID == "" {
		skillID = skill
	}

	name := strings.TrimSpace(record.Name)
	if name == "" {
		name = skill
	}

	return &Skill{
		ID:          skillID,
		Name:        name,
		Description: strings.TrimSpace(record.Description),
		Payload:     defaultPayload(skill, record.Payload),
	}
}

func defaultPayload(skill string, payload json.RawMessage) json.RawMessage {
	if normalized := normalizePayload(payload); len(normalized) != 0 {
		return normalized
	}

	return normalizePayload([]byte(fmt.Sprintf(`{"skill":%q}`, skill)))
}

func normalizePayload(payload []byte) json.RawMessage {
	trimmed := strings.TrimSpace(string(payload))
	if trimmed == "" {
		return nil
	}
	if json.Valid([]byte(trimmed)) {
		return append(json.RawMessage(nil), []byte(trimmed)...)
	}

	encoded, err := json.Marshal(trimmed)
	if err != nil {
		return nil
	}
	return json.RawMessage(encoded)
}
