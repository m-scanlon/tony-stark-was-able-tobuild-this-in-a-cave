package purpose

import (
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

type Purpose struct {
	Value string
}

func (p Purpose) ID() string   { return p.Value }
func (p Purpose) Name() string { return "purpose" }

func (p Purpose) Relate(r logos.Relation) logos.Logos {
	value, _ := meaning.Extract(r.Impulse, "~purpose", "purpose")
	return Purpose{Value: value}
}
