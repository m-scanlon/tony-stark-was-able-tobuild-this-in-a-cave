package world

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

type presentWorld struct{}

func (p presentWorld) DerivePresent(r entity.Relation) string {
	value, err := meaning.Extract(r.Impulse, "~say", "present", "|")
	if err != nil {
		return r.Impulse
	}
	return value
}
