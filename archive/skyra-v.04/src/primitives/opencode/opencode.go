package opencode

import (
	"strings"

	"skyra-v04/src/primitives/being"
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/medium"
)

var _ being.IBeing = OpenCode{}

type OpenCode struct {
	id            string
	name          string
	relationships map[string]any
	medium        medium.Medium
}

func (o OpenCode) ID() string            { return o.id }
func (o OpenCode) Name() string          { return o.name }
func (o OpenCode) Medium() medium.Medium { return o.medium }

func (o OpenCode) Relationships() []string {
	list := make([]string, 0, len(o.relationships))
	for peer := range o.relationships {
		list = append(list, peer)
	}
	return list
}

func (o OpenCode) Relate(r entity.Relation) entity.Entity {
	mediumName, _ := meaning.Extract(r.Impulse, "~medium", "opencode")
	relationshipsRaw, _ := meaning.Extract(r.Impulse, "~relationships", "opencode")
	relationships := make(map[string]any)
	if relationshipsRaw != "" {
		for _, peer := range strings.Split(relationshipsRaw, ",") {
			relationships[strings.TrimSpace(peer)] = nil
		}
	}
	return OpenCode{
		id:            strings.TrimSpace(r.ID),
		name:          strings.TrimSpace(r.ID),
		medium:        medium.Get(mediumName),
		relationships: relationships,
	}
}

// DerivePresent keeps the present minimal. OpenCode manages its own
// context and session state — it only needs the task.
func (o OpenCode) DerivePresent(r entity.Relation) string {
	return ""
}
