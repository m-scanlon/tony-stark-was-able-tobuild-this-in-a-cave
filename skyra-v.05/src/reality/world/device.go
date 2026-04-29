package world

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v05/src/reality"
)

type Device interface {
	reality.Reality
}

type CLIDevice struct {
	id      string
	scanner *bufio.Scanner
}

func NewCLIDevice(id string) *CLIDevice {
	return &CLIDevice{
		id:      id,
		scanner: bufio.NewScanner(os.Stdin),
	}
}

func (c *CLIDevice) ID() string { return c.id }

func (c *CLIDevice) Create(r reality.Relation) reality.Reality {
	return c
}

func (c *CLIDevice) Realize(r reality.Relation) string {
	if r.Impulse != "" {
		fmt.Println("\n" + r.Impulse)
	}
	fmt.Print("> ")
	if !c.scanner.Scan() {
		os.Exit(0)
	}
	input := strings.TrimSpace(c.scanner.Text())
	return input
}
