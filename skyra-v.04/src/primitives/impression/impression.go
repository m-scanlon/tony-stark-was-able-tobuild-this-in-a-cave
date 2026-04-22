package impression

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/invariant"
	"skyra-v04/src/primitives/meaning"
)

type Impression struct {
	invariant.Invariant
	Value string
}

func (i Impression) ID() string   { return i.Value }
func (i Impression) Name() string { return "impression" }

func (i Impression) Relate(r entity.Relation) entity.Entity {
	value, _ := meaning.Extract(r.Impulse, "~impression", "impression")
	return Impression{Value: value}
}
