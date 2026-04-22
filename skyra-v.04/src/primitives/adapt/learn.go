package adapt

import (
	"fmt"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = &LearnLogos{}

type LearnLogos struct {
	presentAdapt
	EntityMap map[string]entity.Entity
}

func (s *LearnLogos) ID() string { return "learn" }

func (s *LearnLogos) Relate(r entity.Relation) entity.Entity {
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
	s.EntityMap[name] = New(name, path)

	grow, ok := s.EntityMap["grow"]
	if !ok {
		fmt.Println("learn: grow not found in world")
		return s
	}
	grow.Relate(r)

	fmt.Println("learn: ok →", name)
	return s
}
