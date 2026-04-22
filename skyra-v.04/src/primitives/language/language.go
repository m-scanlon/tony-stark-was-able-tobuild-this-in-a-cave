package language

import (
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

type Language struct {
	Value string
}

func (l Language) ID() string   { return l.Value }
func (l Language) Name() string { return "language" }

func (l Language) Relate(r logos.Relation) logos.Logos {
	value, _ := meaning.Extract(r.Impulse, "~language", "language")
	return Language{Value: value}
}
