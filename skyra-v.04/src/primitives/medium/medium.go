package medium

import (
	"strings"

	"skyra-v04/src/primitives/entity"
)

type Medium func(present string, r entity.Relation) (string, error)

var registry = map[string]Medium{}

func Register(name string, m Medium) {
	registry[name] = m
}

// Get resolves a medium by name. Supports parametric mediums in the form
// "<name>:<arg>", e.g. "exec:./bin/fetch". Currently only "exec" takes a parameter.
func Get(name string) Medium {
	if strings.HasPrefix(name, "exec:") {
		return execMedium(strings.TrimPrefix(name, "exec:"))
	}
	return registry[name]
}
