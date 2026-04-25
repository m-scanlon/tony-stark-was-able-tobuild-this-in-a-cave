package claude

import (
	"strings"

	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/medium"
)

var _ being.IBeing = Claude{}

type Claude struct {
	id            string
	name          string
	relationships map[string]any
	medium        medium.Medium
}

func (c Claude) ID() string            { return c.id }
func (c Claude) Name() string          { return c.name }
func (c Claude) Medium() medium.Medium { return c.medium }

func (c Claude) Relationships() []string {
	list := make([]string, 0, len(c.relationships))
	for peer := range c.relationships {
		list = append(list, peer)
	}
	return list
}

func (c Claude) Relate(r entity.Relation) entity.Entity {
	mediumName, _ := meaning.Extract(r.Impulse, "~medium", "claude")
	relationshipsRaw, _ := meaning.Extract(r.Impulse, "~relationships", "claude")
	relationships := make(map[string]any)
	if relationshipsRaw != "" {
		for _, peer := range strings.Split(relationshipsRaw, ",") {
			relationships[strings.TrimSpace(peer)] = nil
		}
	}
	return Claude{
		id:            strings.TrimSpace(r.ID),
		name:          strings.TrimSpace(r.ID),
		medium:        medium.Get(mediumName),
		relationships: relationships,
	}
}

// DerivePresent keeps the present minimal. Claude Code manages its own
// context, memory, and session history — it only needs the task.
func (c Claude) DerivePresent(r entity.Relation) string {
	return ""
}
