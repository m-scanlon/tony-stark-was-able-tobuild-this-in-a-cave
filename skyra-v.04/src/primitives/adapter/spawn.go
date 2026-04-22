package adapter

import (
	"fmt"

	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = &SpawnLogos{}

type SpawnLogos struct {
	LogosMap map[string]logos.Logos
}

func (s *SpawnLogos) ID() string { return "spawn" }

func (s *SpawnLogos) Relate(r logos.Relation) logos.Logos {
	name, err := meaning.Extract(r.Impulse, "~name", "spawn")
	if err != nil {
		fmt.Println("spawn: missing ~name")
		return s
	}
	path, err := meaning.Extract(r.Impulse, "~path", "spawn")
	if err != nil {
		fmt.Println("spawn: missing ~path")
		return s
	}
	a, err := Spawn(name, path)
	if err != nil {
		fmt.Println("spawn error:", err)
		return s
	}
	s.LogosMap[name] = a
	fmt.Println("spawn: ok →", name)
	return s
}
