package thread

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

type presentThread struct{}

func (p presentThread) DerivePresent(r entity.Relation) string {
	value, err := meaning.Extract(r.Impulse, "~say", "present", "|")
	if err != nil {
		return r.Impulse
	}
	return value
}
