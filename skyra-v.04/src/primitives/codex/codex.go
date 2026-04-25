package codex

import (
	"strings"

	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/medium"
)

var _ being.IBeing = Codex{}

type Codex struct {
	id            string
	name          string
	relationships map[string]any
	medium        medium.Medium
}

func (c Codex) ID() string            { return c.id }
func (c Codex) Name() string          { return c.name }
func (c Codex) Medium() medium.Medium { return c.medium }

func (c Codex) Relationships() []string {
	list := make([]string, 0, len(c.relationships))
	for peer := range c.relationships {
		list = append(list, peer)
	}
	return list
}

func (c Codex) Relate(r entity.Relation) entity.Entity {
	mediumName, _ := meaning.Extract(r.Impulse, "~medium", "codex")
	relationshipsRaw, _ := meaning.Extract(r.Impulse, "~relationships", "codex")
	relationships := make(map[string]any)
	if relationshipsRaw != "" {
		for _, peer := range strings.Split(relationshipsRaw, ",") {
			relationships[strings.TrimSpace(peer)] = nil
		}
	}
	return Codex{
		id:            strings.TrimSpace(r.ID),
		name:          strings.TrimSpace(r.ID),
		medium:        medium.Get(mediumName),
		relationships: relationships,
	}
}

func (c Codex) DerivePresent(r entity.Relation) string {
	return ""
}
