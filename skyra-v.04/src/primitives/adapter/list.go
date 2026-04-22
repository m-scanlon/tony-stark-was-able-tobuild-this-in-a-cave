package adapter

import (
	"fmt"
	"os"

	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = ListLogos{}

type ListLogos struct{}

func (l ListLogos) ID() string { return "list" }

func (l ListLogos) Relate(rel logos.Relation) logos.Logos {
	path, err := meaning.Extract(rel.Impulse, "~path", "list")
	if err != nil {
		fmt.Println("list: missing ~path")
		return l
	}
	entries, err := os.ReadDir(path)
	if err != nil {
		fmt.Println("list error:", err)
		return l
	}
	for _, e := range entries {
		if e.IsDir() {
			fmt.Println(e.Name() + "/")
		} else {
			fmt.Println(e.Name())
		}
	}
	return l
}
