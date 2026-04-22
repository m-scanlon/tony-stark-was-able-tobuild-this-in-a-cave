package being

import (
	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

type presentBeing struct{}

func (p presentBeing) DerivePresent(r entity.Relation) string {
	value, err := meaning.Extract(r.Impulse, "~say", "present", "|")
	if err != nil {
		return r.Impulse
	}
	return value
}
