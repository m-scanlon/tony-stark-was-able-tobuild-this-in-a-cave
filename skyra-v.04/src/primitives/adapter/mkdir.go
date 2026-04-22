package adapter

import (
	"fmt"
	"os"

	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = MkdirLogos{}

type MkdirLogos struct{}

func (m MkdirLogos) ID() string { return "mkdir" }

func (m MkdirLogos) Relate(rel logos.Relation) logos.Logos {
	path, err := meaning.Extract(rel.Impulse, "~path", "mkdir")
	if err != nil {
		fmt.Println("mkdir: missing ~path")
		return m
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		fmt.Println("mkdir error:", err)
		return m
	}
	fmt.Println("mkdir: ok →", path)
	return m
}
