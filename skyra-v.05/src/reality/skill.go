package reality

import (
	"os"
	"path/filepath"
	"strings"
)

type Skill struct {
	id string
}

func (s *Skill) ID() string { return s.id }

func (s *Skill) Create(r *Relation) Reality {
	return &Skill{id: "skill"}
}

func (s *Skill) Realize(r *Relation) string {
	if r.Log != nil {
		r.Log("[skill]: looking up skill")
	}

	name := strings.TrimSpace(r.Impulse)
	if name == "" {
		return "no skill name provided"
	}

	being, ok := r.Realities["being"]
	if !ok {
		return "no being context"
	}
	b, ok := being.(Being)
	if !ok {
		return "no being context"
	}

	path := filepath.Join(b.Home, "skills", name+".md")
	data, err := os.ReadFile(path)
	if err != nil {
		if r.Log != nil {
			r.Log("[skill]: not found:", path)
		}
		return "skill not found: " + name
	}

	if r.Log != nil {
		r.Log("[skill]: loaded", name)
	}
	return string(data)
}
