package adapter

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"sync"

	"skyra-v04/src/primitives/logos"
)

type IAdapter interface {
	logos.Logos
}

var _ IAdapter = (*Adapter)(nil)

type Adapter struct {
	id     string
	stdin  io.Writer
	stdout *bufio.Reader
	mu     sync.Mutex
}

func Spawn(id, path string, args ...string) (*Adapter, error) {
	cmd := exec.Command(path, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return nil, err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}
	if err := cmd.Start(); err != nil {
		return nil, err
	}
	return &Adapter{
		id:     id,
		stdin:  stdin,
		stdout: bufio.NewReader(stdout),
	}, nil
}

func (a *Adapter) ID() string { return a.id }

func (a *Adapter) Relate(r logos.Relation) logos.Logos {
	a.mu.Lock()
	defer a.mu.Unlock()

	line := fmt.Sprintf("skyra %s %s | %s\n", r.ID, r.Impulse, r.ThreadID)
	if _, err := fmt.Fprint(a.stdin, line); err != nil {
		return a
	}

	response, err := a.stdout.ReadString('\n')
	if err != nil {
		return a
	}

	fmt.Println(strings.TrimSpace(response))
	return a
}
