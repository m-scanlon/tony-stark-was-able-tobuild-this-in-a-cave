package nature

import (
	"skyra-v04/src/primitives/identity"
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/purpose"
)

type Nature struct {
	Identity identity.Identity
	Purpose  purpose.Purpose
}

func (n Nature) ID() string   { return n.Identity.Value }
func (n Nature) Name() string { return "nature" }

func (n Nature) Relate(r logos.Relation) logos.Logos {
	n.Identity, _ = identity.Identity{}.Relate(r).(identity.Identity)
	n.Purpose, _ = purpose.Purpose{}.Relate(r).(purpose.Purpose)
	return n
}
