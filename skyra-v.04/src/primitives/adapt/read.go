package adapt

import (
	"fmt"
	"os"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = ReadLogos{}

type ReadLogos struct {
	presentAdapt
}

func (r ReadLogos) ID() string { return "read" }

func (r ReadLogos) Relate(rel entity.Relation) entity.Entity {
	path, err := meaning.Extract(rel.Impulse, "~path", "read")
	if err != nil {
		fmt.Println("read: missing ~path")
		return r
	}
	content, err := os.ReadFile(path)
	if err != nil {
		fmt.Println("read error:", err)
		return r
	}
	fmt.Print(string(content))
	return r
}
