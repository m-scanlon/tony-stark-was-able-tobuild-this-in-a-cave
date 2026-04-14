package being

import (
	"fmt"
	"strings"

	"skyra-v03/src/primitives/nature"
)

type Flag string

const FlagClose Flag = "close"

type Impulse string

func NewImpulse(raw string) (Impulse, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return "", fmt.Errorf("being: invalid impulse: raw value is required")
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
	Reason     string
}

func ParseImpulse(raw string) (ParsedImpulse, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: raw value is required")
	}

	zones := strings.SplitN(raw, "|", 2)
	if len(zones) != 2 {
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: missing | divider")
	}

	left := strings.TrimSpace(zones[0])
	right := strings.TrimSpace(zones[1])

	parts := strings.Fields(left)
	if len(parts) < 2 {
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: expected at least protocol and target")
	}
	if parts[0] != "skyra" {
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: impulse must begin with skyra")
	}

	targetName := strings.TrimSpace(parts[1])
	if targetName == "" {
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: target name is required")
	}

	flagStart := len(parts)
	for i := 2; i < len(parts); i++ {
		if strings.HasPrefix(parts[i], "~") {
			flagStart = i
			break
		}
	}

	flags := make([]Flag, 0, len(parts)-flagStart)
	for i := flagStart; i < len(parts); i++ {
		token := parts[i]
		if !strings.HasPrefix(token, "~") || len(token) == 1 {
			return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: invalid flag %q", token)
		}
		flags = append(flags, Flag(strings.TrimPrefix(token, "~")))
	}

	expression := strings.Join(parts[2:flagStart], " ")

	reason := strings.TrimSpace(right)
	if reason == "" {
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: reason is required")
	}

	impulse := ParsedImpulse{
		Raw:        raw,
		TargetName: targetName,
		Expression: expression,
		Flags:      flags,
		Reason:     reason,
	}

	if err := impulse.Validate(); err != nil {
		return ParsedImpulse{}, err
	}

	return impulse, nil
}

func (i ParsedImpulse) Validate() error {
	if strings.TrimSpace(i.Raw) == "" {
		return fmt.Errorf("being: invalid impulse: raw value is required")
	}
	if strings.TrimSpace(i.TargetName) == "" {
		return fmt.Errorf("being: invalid impulse: target name is required")
	}
	if i.Expression == "" && !i.IsClose() {
		return fmt.Errorf("being: invalid impulse: expression is required unless close is present")
	}
	if strings.TrimSpace(i.Reason) == "" {
		return fmt.Errorf("being: invalid impulse: reason is required")
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

type DeliveredImpulse struct {
	OriginName string
	Raw        Impulse
	Parsed     ParsedImpulse
}

type ChannelResult struct {
	Routed      bool
	NewExchange bool
	DropReason  string
}

type RelationshipChannel interface {
	Send(delivery DeliveredImpulse) ChannelResult
	Name() string
	PeerNature() nature.Nature
	CallableLanguage() string
}

type PresentDeriver interface {
	DerivePresent(receiver *Being, sender *Being) string
}
