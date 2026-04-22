package adapt

import (
	"fmt"
	"os/exec"

	"skyra-v04/src/primitives/entity"
	"skyra-v04/src/primitives/meaning"
)

var _ IAdapter = BashLogos{}

type BashLogos struct {
	presentAdapt
}

func (b BashLogos) ID() string { return "bash" }

func (b BashLogos) DerivePresent(r entity.Relation) string {
	value, _ := meaning.Extract(r.Impulse, "~cmd", "bash", "|")
	return value
}

func (b BashLogos) Relate(r entity.Relation) entity.Entity {
	cmd := b.DerivePresent(r)
	if cmd == "" {
		fmt.Println("bash: missing ~cmd")
		return b
	}
	out, err := exec.Command("bash", "-c", cmd).CombinedOutput()
	if err != nil {
		fmt.Println("bash error:", err)
	}
	fmt.Print(string(out))
	return b
}
