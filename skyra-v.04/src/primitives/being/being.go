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
}

var _ IBeing = Being{}

type Being struct {
	presentBeing
	id            string
	name          string
	Impression    string
	pathos        pathos.Pathos
	medium        medium.Medium
	operators     []string
	relationships map[string]any
}

func (b Being) ID() string            { return b.id }
func (b Being) Name() string          { return b.name }
func (b Being) Medium() medium.Medium { return b.medium }

func (b Being) Relate(r entity.Relation) entity.Entity {
	p, _ := pathos.Pathos{}.Relate(r).(pathos.Pathos)
	imp, _ := impression.Impression{}.Relate(r).(impression.Impression)
	mediumName, _ := meaning.Extract(r.Impulse, "~medium", "being")
	operatorsRaw, _ := meaning.Extract(r.Impulse, "~operators", "being")
	var operators []string
	if operatorsRaw != "" {
		for _, op := range strings.Split(operatorsRaw, ",") {
			operators = append(operators, strings.TrimSpace(op))
		}
	}
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
		operators:     operators,
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

	if len(b.operators) > 0 {
		sb.WriteString("\noperators:\n")
		for _, op := range b.operators {
			sb.WriteString("  " + op + "\n")
		}
	}

	if len(b.relationships) > 0 {
		sb.WriteString("\nrelationships:\n")
		for peer := range b.relationships {
			sb.WriteString("  " + peer + "\n")
		}
	}

	return sb.String()
}
