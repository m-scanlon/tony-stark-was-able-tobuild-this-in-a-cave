package impression

import (
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

type Impression struct {
	Value string
}

func (i Impression) ID() string   { return i.Value }
func (i Impression) Name() string { return "impression" }

func (i Impression) Relate(r logos.Relation) logos.Logos {
	value, _ := meaning.Extract(r.Impulse, "~impression", "impression")
	return Impression{Value: value}
}
