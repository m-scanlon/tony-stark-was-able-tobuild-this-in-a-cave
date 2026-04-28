package being

import (
	"strings"

	"skyra-v05/src/entity"
)

type Being struct {
	id       string
	name     string
	Identity string
	Purpose  string
}

func (b Being) ID() string   { return b.id }
func (b Being) Name() string { return b.name }

func (b Being) Create(r entity.Relation) entity.Entity {
	identity, _ := Extract(r.Impulse, "~identity", "being")
	purpose, _ := Extract(r.Impulse, "~purpose", "being")

	return Being{
		id:       strings.TrimSpace(r.ID),
		name:     strings.TrimSpace(r.ID),
		Identity: identity,
		Purpose:  purpose,
	}
}

func (b Being) DerivePresent(r entity.Relation) string {
	return ""
}
