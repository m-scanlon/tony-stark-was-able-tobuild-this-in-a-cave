package adapt

import (
	"fmt"
	"os/exec"
	"strings"

	"skyra-v04/src/primitives/entity"
)

type IAdapter interface {
	entity.Entity
}

var _ IAdapter = (*Adapter)(nil)

type Adapter struct {
	presentAdapt
	id   string
	path string
}

func New(id, path string) *Adapter {
	return &Adapter{id: id, path: path}
}

func (a *Adapter) ID() string { return a.id }

func (a *Adapter) Relate(r entity.Relation) entity.Entity {
	line := fmt.Sprintf("skyra %s %s | %s\n", r.ID, r.Impulse, r.ThreadID)
	cmd := exec.Command(a.path)
	cmd.Stdin = strings.NewReader(line)
	out, err := cmd.Output()
	if err != nil {
		fmt.Println("adapter error:", err)
		return a
	}
	fmt.Print(string(out))
	return a
}
