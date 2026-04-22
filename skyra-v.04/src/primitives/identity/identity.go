package identity

import (
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

type Identity struct {
	Value string
}

func (i Identity) ID() string   { return i.Value }
func (i Identity) Name() string { return "identity" }

func (i Identity) Relate(r logos.Relation) logos.Logos {
	value, _ := meaning.Extract(r.Impulse, "~identity", "identity")
	return Identity{Value: value}
}
