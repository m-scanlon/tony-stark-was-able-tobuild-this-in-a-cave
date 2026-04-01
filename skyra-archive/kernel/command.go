package kernel

import (
	"fmt"
	"strings"
)

// Command is a parsed skyra <skill> [args] invocation.
type Command struct {
	Skill string
	Args  []string
	Raw   string
}

// ParseCommand parses a raw "skyra <skill> [args...]" string.
func ParseCommand(raw string) (Command, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return Command{}, fmt.Errorf("empty command")
	}

	parts := strings.Fields(raw)
	if len(parts) < 2 {
		return Command{}, fmt.Errorf("invalid command: expected \"skyra <skill> [args]\", got %q", raw)
	}
	if !strings.EqualFold(parts[0], "skyra") {
		return Command{}, fmt.Errorf("invalid command: must start with \"skyra\", got %q", parts[0])
	}

	return Command{
		Skill: parts[1],
		Args:  parts[2:],
		Raw:   raw,
	}, nil
}
