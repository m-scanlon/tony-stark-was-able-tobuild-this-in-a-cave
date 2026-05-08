package reality

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"skyra-v05/src/debug"
	"strings"
	"time"
)

type Process struct {
	id      string
	Cmd     *exec.Cmd
	Stdin   io.WriteCloser
	Output  chan string
	Command string
	Args    []string
}

func (p *Process) ID() string { return p.id }

func (p *Process) Create(r *Relation) Reality {
	return &Process{id: r.ID}
}

func (p *Process) Start(command string, args ...string) error {
	cmd := exec.Command(command, args...)
	stdin, err := cmd.StdinPipe()
	if err != nil {
		return err
	}
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}
	cmd.Stderr = cmd.Stdout

	if err := cmd.Start(); err != nil {
		return err
	}

	p.Cmd = cmd
	p.Stdin = stdin
	p.Output = make(chan string, 100)

	go func() {
		scanner := bufio.NewScanner(stdout)
		for scanner.Scan() {
			p.Output <- scanner.Text()
		}
		close(p.Output)
	}()

	startup := p.drain()
	if startup != "" {
		debug.Log("[process]:", p.id, "startup →", startup)
		fmt.Printf("[%s] %s\n", p.id, startup)
	}

	return nil
}

func (p *Process) Realize(r *Relation) string {
	if r.Collecting {
		r.Export("node:"+p.id, RealityNode{ID: p.id, Type: "Process", Children: []RealityNode{}})
		return ""
	}

	debug.Log("[process]:", p.id, "← impulse:", r.Impulse)

	if p.Cmd == nil {
		if p.Command == "" {
			return "process not configured"
		}
		debug.Log("[process]:", p.id, "lazy start →", p.Command)
		if err := p.Start(p.Command, p.Args...); err != nil {
			debug.Log("[process]:", p.id, "start error:", err)
			return "process failed to start: " + err.Error()
		}
	}

	command := strings.TrimSpace(r.Impulse)
	if command != "" {
		fmt.Printf("\n[%s] > %s\n", p.id, command)

		_, err := io.WriteString(p.Stdin, command+"\n")
		if err != nil {
			debug.Log("[process]:", p.id, "write error:", err)
			return "process error: " + err.Error()
		}
	}

	output := p.drain()
	fmt.Printf("[%s] %s\n", p.id, output)
	debug.Log("[process]:", p.id, "→", output)

	return output
}

func (p *Process) drain() string {
	var lines []string
	timeout := time.After(2 * time.Second)
	settling := time.After(500 * time.Millisecond)

	for {
		select {
		case line, ok := <-p.Output:
			if !ok {
				return strings.Join(lines, "\n")
			}
			lines = append(lines, line)
			settling = time.After(500 * time.Millisecond)
		case <-settling:
			if len(lines) > 0 {
				return strings.Join(lines, "\n")
			}
		case <-timeout:
			return strings.Join(lines, "\n")
		}
	}
}
