package adapt

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = GrepLogos{}

type GrepLogos struct {
	presentAdapt
}

func (g GrepLogos) ID() string { return "grep" }

func (g GrepLogos) Relate(rel entity.Relation) entity.Entity {
	pattern, err := meaning.Extract(rel.Impulse, "~pattern", "grep")
	if err != nil {
		fmt.Println("grep: missing ~pattern")
		return g
	}
	root, err := meaning.Extract(rel.Impulse, "~path", "grep")
	if err != nil {
		fmt.Println("grep: missing ~path")
		return g
	}

	filepath.WalkDir(root, func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return nil
		}
		f, err := os.Open(path)
		if err != nil {
			return nil
		}
		defer f.Close()

		scanner := bufio.NewScanner(f)
		lineNum := 0
		for scanner.Scan() {
			lineNum++
			if strings.Contains(scanner.Text(), pattern) {
				fmt.Printf("%s:%d: %s\n", path, lineNum, scanner.Text())
			}
		}
		return nil
	})
	return g
}
