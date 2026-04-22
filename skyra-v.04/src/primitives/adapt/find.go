package adapt

import (
	"fmt"
	"io/fs"
	"path/filepath"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = FindLogos{}

type FindLogos struct{ presentAdapt }

func (f FindLogos) ID() string { return "find" }

func (f FindLogos) Relate(rel entity.Relation) entity.Entity {
	root, err := meaning.Extract(rel.Impulse, "~path", "find")
	if err != nil {
		fmt.Println("find: missing ~path")
		return f
	}
	pattern, _ := meaning.Extract(rel.Impulse, "~name", "find")

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return nil
		}
		if pattern != "" {
			matched, _ := filepath.Match(pattern, d.Name())
			if !matched {
				return nil
			}
		}
		fmt.Println(path)
		return nil
	})
	return f
}
