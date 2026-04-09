package domain

import (
	"fmt"
	"strings"
)

type Flag string

const FlagClose Flag = "close"

type Impulse string

func NewImpulse(raw string) (Impulse, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("%w: raw value is required", ErrInvalidImpulse)
	}
	if _, err := ParseImpulse(raw); err != nil {
		return "", err
	}
	return Impulse(raw), nil
}

func (i Impulse) Raw() string {
	return string(i)
}

func (i Impulse) Parse() (ParsedImpulse, error) {
	return ParseImpulse(i.Raw())
}

func (i Impulse) IsClose() bool {
	parsed, err := i.Parse()
	if err != nil {
		return false
	}
	return parsed.IsClose()
}

type ParsedImpulse struct {
	Raw        string
	TargetName string
	Expression string
	Flags      []Flag
	Source     string
	Reason     string
}

func ParseImpulse(raw string) (ParsedImpulse, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ParsedImpulse{}, fmt.Errorf("%w: raw value is required", ErrInvalidImpulse)
	}

	zones := strings.SplitN(raw, "|", 2)
	if len(zones) != 2 {
		return ParsedImpulse{}, fmt.Errorf("%w: missing | divider", ErrInvalidImpulse)
	}

	left := strings.TrimSpace(zones[0])
	right := strings.TrimSpace(zones[1])

	// Parse left zone: skyra <being> <expression> -<flags>
	parts := strings.Fields(left)
	if len(parts) < 2 {
		return ParsedImpulse{}, fmt.Errorf("%w: expected at least protocol and target", ErrInvalidImpulse)
	}
	if parts[0] != "skyra" {
		return ParsedImpulse{}, fmt.Errorf("%w: impulse must begin with skyra", ErrInvalidImpulse)
	}

	targetName := strings.TrimSpace(parts[1])
	if targetName == "" {
		return ParsedImpulse{}, fmt.Errorf("%w: target name is required", ErrInvalidImpulse)
	}

	flagStart := len(parts)
	for i := 2; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "-") {
			flagStart = i
			break
		}
	}

	flags := make([]Flag, 0, len(parts)-flagStart)
	for i := flagStart; i < len(parts); i++ {
		token := parts[i]
		if !strings.HasPrefix(token, "-") || len(token) == 1 {
			return ParsedImpulse{}, fmt.Errorf("%w: invalid flag %q", ErrInvalidImpulse, token)
		}
		flags = append(flags, Flag(strings.TrimPrefix(token, "-")))
	}

	expression := strings.Join(parts[2:flagStart], " ")

	// Parse right zone: <source>: <reason> ~<emotional_signals>
	colonIdx := strings.Index(right, ":")
	if colonIdx < 0 {
		return ParsedImpulse{}, fmt.Errorf("%w: missing : separator in right zone", ErrInvalidImpulse)
	}

	source := strings.TrimSpace(right[:colonIdx])
	if source == "" {
		return ParsedImpulse{}, fmt.Errorf("%w: source is required", ErrInvalidImpulse)
	}

	rest := strings.TrimSpace(right[colonIdx+1:])

	tildeIdx := strings.Index(rest, "~")
	var reason string
	if tildeIdx < 0 {
		reason = strings.TrimSpace(rest)
	} else {
		reason = strings.TrimSpace(rest[:tildeIdx])
	}

	if reason == "" {
		return ParsedImpulse{}, fmt.Errorf("%w: reason is required", ErrInvalidImpulse)
	}

	impulse := ParsedImpulse{
		Raw:        raw,
		TargetName: targetName,
		Expression: expression,
		Flags:      flags,
		Source:     source,
		Reason:     reason,
	}

	if err := impulse.Validate(); err != nil {
		return ParsedImpulse{}, err
	}

	return impulse, nil
}

func (i ParsedImpulse) Validate() error {
	if strings.TrimSpace(i.Raw) == "" {
		return fmt.Errorf("%w: raw value is required", ErrInvalidImpulse)
	}
	if strings.TrimSpace(i.TargetName) == "" {
		return fmt.Errorf("%w: target name is required", ErrInvalidImpulse)
	}
	if i.Expression == "" && !i.IsClose() {
		return fmt.Errorf("%w: expression is required unless close is present", ErrInvalidImpulse)
	}
	if strings.TrimSpace(i.Source) == "" {
		return fmt.Errorf("%w: source is required", ErrInvalidImpulse)
	}
	if strings.TrimSpace(i.Reason) == "" {
		return fmt.Errorf("%w: reason is required", ErrInvalidImpulse)
	}
	return nil
}

func (i ParsedImpulse) HasFlag(flag Flag) bool {
	for _, existing := range i.Flags {
		if existing == flag {
			return true
		}
	}
	return false
}

func (i ParsedImpulse) IsClose() bool {
	return i.HasFlag(FlagClose)
}
