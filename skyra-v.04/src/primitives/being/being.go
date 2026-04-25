package being

import (
	"strings"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/impression"
	"skyra-v04/src/primitives/meaning"
	"skyra-v04/src/primitives/medium"
	"skyra-v04/src/primitives/pathos"
)

type IBeing interface {
	entity.Entity
	Name() string
	Medium() medium.Medium
	Relationships() []string
}

var _ IBeing = Being{}

type Being struct {
	id            string
	name          string
	Impression    string
	pathos        pathos.Pathos
	medium        medium.Medium
	relationships map[string]any
}

func (b Being) ID() string            { return b.id }
func (b Being) Name() string          { return b.name }
func (b Being) Medium() medium.Medium { return b.medium }

func (b Being) Relationships() []string {
	list := make([]string, 0, len(b.relationships))
	for peer := range b.relationships {
		list = append(list, peer)
	}
	return list
}

func (b Being) Relate(r entity.Relation) entity.Entity {
	p, _ := pathos.Pathos{}.Relate(r).(pathos.Pathos)
	imp, _ := impression.Impression{}.Relate(r).(impression.Impression)
	mediumName, _ := meaning.Extract(r.Impulse, "~medium", "being")
	relationshipsRaw, _ := meaning.Extract(r.Impulse, "~relationships", "being")
	relationships := make(map[string]any)
	if relationshipsRaw != "" {
		for _, peer := range strings.Split(relationshipsRaw, ",") {
			relationships[strings.TrimSpace(peer)] = nil
		}
	}
	return Being{
		id:            strings.TrimSpace(r.ID),
		name:          strings.TrimSpace(r.ID),
		Impression:    imp.Value,
		pathos:        p,
		medium:        medium.Get(mediumName),
		relationships: relationships,
	}
}

func (b Being) DerivePresent(r entity.Relation) string {
	var sb strings.Builder

	sb.WriteString("being: " + b.name + "\n")
	if b.pathos.Identity.Value != "" {
		sb.WriteString("identity: " + b.pathos.Identity.Value + "\n")
	}
	if b.pathos.Purpose.Value != "" {
		sb.WriteString("purpose: " + b.pathos.Purpose.Value + "\n")
	}
	if b.Impression != "" {
		sb.WriteString("impression: " + b.Impression + "\n")
	}

	return sb.String()
}
