package medium

import "skyra-v04/src/primitives/entity"

type Medium func(present string, r entity.Relation) (string, error)

var registry = map[string]Medium{}

func Register(name string, m Medium) {
	registry[name] = m
}

func Get(name string) Medium {
	return registry[name]
}
