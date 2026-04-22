package adapter

import (
	"fmt"

	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = &LearnLogos{}

type LearnLogos struct {
	LogosMap map[string]logos.Logos
}

func (s *LearnLogos) ID() string { return "learn" }

func (s *LearnLogos) Relate(r logos.Relation) logos.Logos {
	name, err := meaning.Extract(r.Impulse, "~name", "learn")
	if err != nil {
		fmt.Println("learn: missing ~name")
		return s
	}
	path, err := meaning.Extract(r.Impulse, "~path", "learn")
	if err != nil {
		fmt.Println("learn: missing ~path")
		return s
	}
	a, err := Spawn(name, path)
	if err != nil {
		fmt.Println("learn error:", err)
		return s
	}
	s.LogosMap[name] = a

	grow, ok := s.LogosMap["grow"]
	if !ok {
		fmt.Println("learn: grow not found in world")
		return s
	}
	grow.Relate(r)

	fmt.Println("learn: ok →", name)
	return s
}
