package identity

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/invariant"
	"skyra-v04/src/primitives/meaning"
)

type Identity struct {
	invariant.Invariant
	Value string
}

func (i Identity) ID() string   { return i.Value }
func (i Identity) Name() string { return "identity" }

func (i Identity) Relate(r entity.Relation) entity.Entity {
	value, _ := meaning.Extract(r.Impulse, "~identity", "identity")
	return Identity{Value: value}
}
