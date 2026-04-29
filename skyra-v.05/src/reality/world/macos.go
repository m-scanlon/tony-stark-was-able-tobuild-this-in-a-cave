package world

import (
	"skyra-v05/src/reality"
)

type MacOS struct {
	World
	Terminal *CLIDevice
}

func NewMacOS() *MacOS {
	terminal := NewCLIDevice("terminal")
	m := &MacOS{
		World: World{
			id:        "macos",
			name:      "macos",
			Realities: make(map[string]reality.Reality),
		},
		Terminal: terminal,
	}
	m.Realities["terminal"] = terminal
	return m
}

func (m *MacOS) Run(system *System) {
	for {
		input := m.Terminal.Realize(reality.Relation{})
		if input == "" {
			continue
		}
		rel, err := reality.Impress("terminal", "", input)
		if err != nil {
			continue
		}
		response := system.Realize(rel)
		if response != "" {
			m.Terminal.Realize(reality.Relation{Impulse: response})
		}
	}
}
