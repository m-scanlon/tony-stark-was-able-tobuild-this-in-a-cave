package being

import (
	"strings"

	"skyra-v05/src/reality"
)

type Being struct {
	id       string
	name     string
	Identity string
	Purpose  string
}

func (b Being) ID() string   { return b.id }
func (b Being) Name() string { return b.name }

func (b Being) Create(r reality.Relation) reality.Reality {
	identity, _ := Extract(r.Impulse, "~identity", "being")
	purpose, _ := Extract(r.Impulse, "~purpose", "being")

	return Being{
		id:       strings.TrimSpace(r.ID),
		name:     strings.TrimSpace(r.ID),
		Identity: identity,
		Purpose:  purpose,
	}
}

func (b Being) Realize(r reality.Relation) string {
	var sb strings.Builder
	sb.WriteString("being: " + b.name + "\n")
	if b.Identity != "" {
		sb.WriteString("identity: " + b.Identity + "\n")
	}
	if b.Purpose != "" {
		sb.WriteString("purpose: " + b.Purpose + "\n")
	}
	return sb.String()
}
