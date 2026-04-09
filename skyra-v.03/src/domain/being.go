package domain

import (
	"fmt"
	"strings"
)

type Nature struct {
	Identity string
	Purpose  string
}

func (n Nature) Validate() error {
	if strings.TrimSpace(n.Identity) == "" {
		return ErrEmptyIdentity
	}
	if strings.TrimSpace(n.Purpose) == "" {
		return ErrEmptyPurpose
	}
	return nil
}

type Being struct {
	Name             string
	Nature           Nature
	Differentiatable bool
	Cognitive        bool
	Peers            map[string]RelationshipChannel
}

func NewBeing(name string, nature Nature, differentiatable bool, cognitive bool) (*Being, error) {
	being := &Being{
		Name:             strings.TrimSpace(name),
		Nature:           nature,
		Differentiatable: differentiatable,
		Cognitive:        cognitive,
		Peers:            make(map[string]RelationshipChannel),
	}

	if err := being.Validate(); err != nil {
		return nil, err
	}

	return being, nil
}

func (b *Being) Validate() error {
	if b == nil {
		return ErrNilBeing
	}
	if b.Name == "" {
		return ErrEmptyBeingName
	}
	if err := b.Nature.Validate(); err != nil {
		return err
	}
	if b.Peers == nil {
		b.Peers = make(map[string]RelationshipChannel)
	}
	return nil
}

func (b *Being) OriginName() string {
	if b == nil {
		return ""
	}
	return b.Name
}

func (b *Being) MatchesOriginName(name string) bool {
	if b == nil {
		return false
	}
	return b.Name == strings.TrimSpace(name)
}

func (b *Being) SeedPeer(peerName string, peerNature Nature) error {
	if err := b.Validate(); err != nil {
		return err
	}
	if _, exists := b.Peers[peerName]; exists {
		return nil
	}

	peer, err := NewExchangeStack(peerName, peerNature)
	if err != nil {
		return err
	}
	b.AttachPeer(peer)
	return nil
}

func (b *Being) AttachPeer(channel RelationshipChannel) error {
	if err := b.Validate(); err != nil {
		return err
	}
	if channel == nil {
		return ErrUnknownPeer
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
		return ChannelResult{}, fmt.Errorf("%w: %s", ErrUnknownPeer, peerName)
	}

	if delivery.Raw == "" {
		return ChannelResult{}, ErrInvalidImpulse
	}
	if delivery.Parsed.Raw == "" {
		return ChannelResult{}, ErrInvalidImpulse
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
		return "", ErrNilBeing
	}

	peer, ok := b.Peers[sender.Name]
	if !ok {
		return "", fmt.Errorf("%w: %s", ErrUnknownPeer, sender.Name)
	}

	deriver, ok := peer.(presentDeriver)
	if !ok {
		return "", nil
	}

	return deriver.derivePresent(b, sender), nil
}
