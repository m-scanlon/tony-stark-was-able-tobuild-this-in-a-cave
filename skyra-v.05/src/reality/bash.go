package reality

import (
	"bytes"
	"os/exec"
	"strings"
	"time"
)

type Bash struct {
	id      string
	WorkDir string
	Timeout time.Duration
}

func (b *Bash) ID() string { return b.id }

func (b *Bash) Create(r *Relation) Reality {
	return &Bash{id: "bash", Timeout: 30 * time.Second}
}

func (b *Bash) Realize(r *Relation) string {
	command := strings.TrimSpace(r.Impulse)
	if command == "" {
		return "no command"
	}

	cmd := exec.Command("bash", "-c", command)
	if b.WorkDir != "" {
		cmd.Dir = b.WorkDir
	}

	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	done := make(chan error, 1)
	if err := cmd.Start(); err != nil {
		return "error: " + err.Error()
	}
	go func() { done <- cmd.Wait() }()

	timeout := b.Timeout
	if timeout == 0 {
		timeout = 30 * time.Second
	}

	select {
	case err := <-done:
		out := strings.TrimSpace(stdout.String())
		errOut := strings.TrimSpace(stderr.String())
		if err != nil {
			if errOut != "" {
				return errOut
			}
			return "error: " + err.Error()
		}
		if out == "" && errOut != "" {
			return errOut
		}
		return out
	case <-time.After(timeout):
		cmd.Process.Kill()
		return "timeout after " + timeout.String()
	}
}
