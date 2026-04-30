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
	if r.Impulse != "" {
		fmt.Println("\n" + r.Impulse)
	}
	for {
		fmt.Print("> ")
		if !m.scanner.Scan() {
			os.Exit(0)
		}
		input := strings.TrimSpace(m.scanner.Text())
		if input != "" {
			return input
		}
	}
}

func (m *MacOS) Parse() string {
	return ""
}
