package reality

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

type MacOS struct {
	id        string
	Realities map[string]Reality
	scanner   *bufio.Scanner
}

func (m *MacOS) ID() string { return m.id }

func (m *MacOS) Create(r *Relation) Reality {
	return &MacOS{
		id:        "macos",
		Realities: make(map[string]Reality),
		scanner:   bufio.NewScanner(os.Stdin),
	}
}

func (m *MacOS) Realize(r *Relation) string {
	if r.Collecting {
		return ""
	}
	if r.Impulse != "" {
		fmt.Println("\n" + r.Impulse)
	}
	fmt.Print("> ")
	if !m.scanner.Scan() {
		os.Exit(0)
	}
	first := m.scanner.Text()
	if strings.TrimSpace(first) == "" {
		return m.Realize(r)
	}
	if !strings.HasSuffix(strings.TrimSpace(first), ";;") {
		return strings.TrimSpace(first)
	}
	lines := []string{strings.TrimSuffix(strings.TrimSpace(first), ";;")}
	for {
		fmt.Print("  ")
		if !m.scanner.Scan() {
			os.Exit(0)
		}
		line := m.scanner.Text()
		if strings.TrimSpace(line) == ";;" {
			return strings.TrimSpace(strings.Join(lines, "\n"))
		}
		lines = append(lines, line)
	}
}
