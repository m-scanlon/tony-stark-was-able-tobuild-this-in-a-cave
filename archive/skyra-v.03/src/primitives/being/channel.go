package being

import (
	"fmt"
	"strings"

	"skyra-v03/src/primitives/nature"
)

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

type ParsedImpulse struct {
	Raw        string
	TargetName string
	Expression string
	Reason     string
}

func ParseImpulse(raw string) (ParsedImpulse, error) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: raw value is required")
	}

	switch strings.Count(raw, "|") {
	case 0:
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: missing | divider")
	case 1:
		// expected shape
	default:
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: expected exactly one | divider")
	}

	zones := strings.SplitN(raw, "|", 2)

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

	expression := strings.Join(parts[2:], " ")

	reason := strings.TrimSpace(right)
	if reason == "" {
		return ParsedImpulse{}, fmt.Errorf("being: invalid impulse: reason is required")
	}

	impulse := ParsedImpulse{
		Raw:        raw,
		TargetName: targetName,
		Expression: expression,
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
	if strings.TrimSpace(i.Expression) == "" {
		return fmt.Errorf("being: invalid impulse: expression is required")
	}
	if strings.TrimSpace(i.Reason) == "" {
		return fmt.Errorf("being: invalid impulse: reason is required")
	}
	return nil
}

type DeliveredImpulse struct {
	OriginName     string
	ThreadID       string
	About          string
	Because        string
	ContextEntries []ExchangeEntry
	Raw            Impulse
	Parsed         ParsedImpulse
}

type ExchangeEntry struct {
	Author  string
	Impulse Impulse
}

type ChannelResult struct {
	Routed     bool
	DropReason string
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
