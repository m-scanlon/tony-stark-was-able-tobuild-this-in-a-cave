package pathos

import (
	"skyra-v04/src/primitives/identity"
	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/purpose"
)

type Pathos struct {
	Identity identity.Identity
	Purpose  purpose.Purpose
}

func (p Pathos) ID() string { return p.Identity.Value }

func (p Pathos) Relate(r logos.Relation) logos.Logos {
	p.Identity, _ = identity.Identity{}.Relate(r).(identity.Identity)
	p.Purpose, _ = purpose.Purpose{}.Relate(r).(purpose.Purpose)
	return p
}
