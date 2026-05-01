package purpose

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/invariant"
	"skyra-v04/src/primitives/meaning"
)

type Purpose struct {
	invariant.Invariant
	Value string
}

func (p Purpose) ID() string   { return p.Value }
func (p Purpose) Name() string { return "purpose" }

func (p Purpose) Relate(r entity.Relation) entity.Entity {
	value, _ := meaning.Extract(r.Impulse, "~purpose", "purpose")
	return Purpose{Value: value}
}
