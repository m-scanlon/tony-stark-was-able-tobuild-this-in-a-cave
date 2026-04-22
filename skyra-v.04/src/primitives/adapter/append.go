package adapter

import (
	"fmt"
	"os"

	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = AppendLogos{}

type AppendLogos struct{}

func (a AppendLogos) ID() string { return "append" }

func (a AppendLogos) Relate(rel logos.Relation) logos.Logos {
	path, err := meaning.Extract(rel.Impulse, "~path", "append")
	if err != nil {
		fmt.Println("append: missing ~path")
		return a
	}
	content, err := meaning.ExtractToEnd(rel.Impulse, "~content", "append")
	if err != nil {
		fmt.Println("append: missing ~content")
		return a
	}
	f, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("append error:", err)
		return a
	}
	defer f.Close()
	if _, err := fmt.Fprintln(f, content); err != nil {
		fmt.Println("append error:", err)
		return a
	}
	fmt.Println("append: ok →", path)
	return a
}
