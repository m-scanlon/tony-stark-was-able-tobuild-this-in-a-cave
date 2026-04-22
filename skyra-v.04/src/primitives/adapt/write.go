package adapt

import (
	"fmt"
	"os"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = WriteLogos{}

type WriteLogos struct {
	presentAdapt
}

func (w WriteLogos) ID() string { return "write" }

func (w WriteLogos) DerivePresent(r entity.Relation) string {
	value, _ := meaning.Extract(r.Impulse, "~content", "write", "|")
	return value
}

func (w WriteLogos) Relate(rel entity.Relation) entity.Entity {
	path, err := meaning.Extract(rel.Impulse, "~path", "write")
	if err != nil {
		fmt.Println("write: missing ~path")
		return w
	}
	content := w.DerivePresent(rel)
	if content == "" {
		fmt.Println("write: missing ~content")
		return w
	}
	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		fmt.Println("write error:", err)
		return w
	}
	fmt.Println("write: ok →", path)
	return w
}
