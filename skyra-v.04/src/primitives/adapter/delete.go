package adapter

import (
	"fmt"
	"os"

	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = DeleteLogos{}

type DeleteLogos struct{}

func (d DeleteLogos) ID() string { return "delete" }

func (d DeleteLogos) Relate(rel logos.Relation) logos.Logos {
	path, err := meaning.Extract(rel.Impulse, "~path", "delete")
	if err != nil {
		fmt.Println("delete: missing ~path")
		return d
	}
	if err := os.Remove(path); err != nil {
		fmt.Println("delete error:", err)
		return d
	}
	fmt.Println("delete: ok →", path)
	return d
}
