package reality

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type Terminal struct {
	id      string
	Device  Reality
	scanner *bufio.Scanner
}

func (t *Terminal) ID() string { return t.id }

func (t *Terminal) Create(r *Relation) Reality {
	return &Terminal{
		id:      "terminal",
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (t *Terminal) Realize(r *Relation) string {
	if r.Collecting {
		return ""
	}
	if r.Impulse != "" {
		fmt.Println("\n" + r.Impulse)
	}
	fmt.Print("> ")
	if !t.scanner.Scan() {
		os.Exit(0)
	}
	first := t.scanner.Text()
	if strings.TrimSpace(first) == "" {
		return t.Realize(r)
	}
	if !strings.HasSuffix(strings.TrimSpace(first), ";;") {
		return strings.TrimSpace(first)
	}
	lines := []string{strings.TrimSuffix(strings.TrimSpace(first), ";;")}
	for {
		fmt.Print("  ")
		if !t.scanner.Scan() {
			os.Exit(0)
		}
		line := t.scanner.Text()
		if strings.TrimSpace(line) == ";;" {
			return strings.TrimSpace(strings.Join(lines, "\n"))
		}
		lines = append(lines, line)
	}
}
