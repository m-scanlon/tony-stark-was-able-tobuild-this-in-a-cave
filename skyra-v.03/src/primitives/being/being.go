package being

import (
	"fmt"
	"strings"

	"skyra-v03/src/primitives/extract"
	"skyra-v03/src/primitives/language"
	"skyra-v03/src/primitives/nature"
)

type Being struct {
	Name      string
	Nature    nature.Nature
	Language  language.Language
	Cognitive bool
	Peers     map[string]RelationshipChannel
}

func CreateBeing(expression string) (*Being, error) {
	name, err := extract.Meaning(expression, "~name", "being")
	if err != nil {
		return nil, err
	}

	n, err := nature.CreateNature(expression)
	if err != nil {
		return nil, err
	}

	l, err := language.CreateLanguage(expression)
	if err != nil {
		return nil, err
	}

	cognitive, err := extractCognitive(expression)
	if err != nil {
		return nil, err
	}

	return &Being{
		Name:      strings.TrimSpace(name),
		Nature:    n,
		Language:  l,
		Cognitive: cognitive,
		Peers:     make(map[string]RelationshipChannel),
	}, nil
}

func NewBeing(name string, n nature.Nature, l language.Language, cognitive bool) (*Being, error) {
	b := &Being{
		Name:      strings.TrimSpace(name),
		Nature:    n,
		Language:  l,
		Cognitive: cognitive,
		Peers:     make(map[string]RelationshipChannel),
	}
	if err := b.Validate(); err != nil {
		return nil, err
	}
	return b, nil
}

func (b *Being) Validate() error {
	if b == nil {
		return fmt.Errorf("being: nil being")
	}
	if b.Name == "" {
		return fmt.Errorf("being: name is required")
	}
	if strings.TrimSpace(b.Nature.Identity.Value) == "" {
		return fmt.Errorf("being: identity is required")
	}
	if strings.TrimSpace(b.Nature.Purpose.Value) == "" {
		return fmt.Errorf("being: purpose is required")
	}
	if strings.TrimSpace(b.Language.Value) == "" {
		return fmt.Errorf("being: language is required")
	}
	if b.Peers == nil {
		b.Peers = make(map[string]RelationshipChannel)
	}
	return nil
}

func (b *Being) AttachPeer(channel RelationshipChannel) error {
	if err := b.Validate(); err != nil {
		return err
	}
	if channel == nil {
		return fmt.Errorf("being: channel is nil")
	}
	b.Peers[channel.Name()] = channel
	return nil
}

func (b *Being) PeerByName(peerName string) (RelationshipChannel, bool) {
	if b == nil {
		return nil, false
	}
	peer, ok := b.Peers[peerName]
	return peer, ok
}

func (b *Being) SendToPeer(peerName string, delivery DeliveredImpulse) (ChannelResult, error) {
	if err := b.Validate(); err != nil {
		return ChannelResult{}, err
	}
	peer, ok := b.Peers[peerName]
	if !ok {
		return ChannelResult{}, fmt.Errorf("being: unknown peer %s", peerName)
	}
	if delivery.Raw == "" {
		return ChannelResult{}, fmt.Errorf("being: invalid impulse")
	}
	return peer.Send(delivery), nil
}

func (b *Being) EmitToPeer(peerName string, impulse Impulse) (ChannelResult, error) {
	parsed, err := impulse.Parse()
	if err != nil {
		return ChannelResult{}, err
	}
	return b.SendToPeer(peerName, DeliveredImpulse{
		OriginName: b.Name,
		Raw:        impulse,
		Parsed:     parsed,
	})
}

func (b *Being) DerivePresent(sender *Being) (string, error) {
	if err := b.Validate(); err != nil {
		return "", err
	}
	if sender == nil {
		return "", fmt.Errorf("being: nil sender")
	}
	peer, ok := b.Peers[sender.Name]
	if !ok {
		return "", fmt.Errorf("being: unknown peer %s", sender.Name)
	}
	deriver, ok := peer.(PresentDeriver)
	if !ok {
		return "", nil
	}
	return deriver.DerivePresent(b, sender), nil
}

func extractCognitive(expression string) (bool, error) {
	value, err := extract.Meaning(expression, "~cognitive", "being")
	if err != nil {
		return false, err
	}
	switch strings.TrimSpace(value) {
	case "true":
		return true, nil
	case "false":
		return false, nil
	default:
		return false, fmt.Errorf("being: ~cognitive must be true or false, got %q", value)
	}
}
