package exchange

import "skyra-v04/src/primitives/logos"

type Exchange struct {
	entries []string
}

func (e Exchange) Relate(r logos.Relation) logos.Logos {
	return Exchange{entries: append(e.entries, r.Impulse)}
}

func (e Exchange) ID() string      { return "" }
func (e Exchange) Name() string    { return "exchange" }
func (e Exchange) Entries() []string { return e.entries }
