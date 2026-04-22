package adapter

import (
	"fmt"
	"os"

	"skyra-v04/src/primitives/logos"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = MoveLogos{}

type MoveLogos struct{}

func (m MoveLogos) ID() string { return "move" }

func (m MoveLogos) Relate(rel logos.Relation) logos.Logos {
	from, err := meaning.Extract(rel.Impulse, "~from", "move")
	if err != nil {
		fmt.Println("move: missing ~from")
		return m
	}
	to, err := meaning.Extract(rel.Impulse, "~to", "move")
	if err != nil {
		fmt.Println("move: missing ~to")
		return m
	}
	if err := os.Rename(from, to); err != nil {
		fmt.Println("move error:", err)
		return m
	}
	fmt.Println("move: ok →", from, "→", to)
	return m
}
