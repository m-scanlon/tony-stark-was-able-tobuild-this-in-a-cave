package entity

import "skyra-v04/src/primitives/meaning"

type PresentEntity struct{}

func (p PresentEntity) DerivePresent(r Relation) string {
	value, err := meaning.Extract(r.Impulse, "~say", "present", "|")
	if err != nil {
		return r.Impulse
	}
	return value
}
