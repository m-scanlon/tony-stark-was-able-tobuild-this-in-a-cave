package world

import (
	"fmt"
	"strings"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/nature"
)

type ExternalDispatch struct {
	peerName         string
	peerNature       nature.Nature
	lastExpression   string
	callableLanguage string
}

func NewExternalDispatch(peerName string, peerNature nature.Nature) (*ExternalDispatch, error) {
	if strings.TrimSpace(peerName) == "" {
		return nil, fmt.Errorf("%w: peer name is required", ErrUnknownPeer)
	}
	if err := peerNature.Validate(); err != nil {
		return nil, err
	}
	return &ExternalDispatch{
		peerName:   peerName,
		peerNature: peerNature,
	}, nil
}

func (d *ExternalDispatch) Name() string {
	if d == nil {
		return ""
	}
	return d.peerName
}

func (d *ExternalDispatch) PeerNature() nature.Nature {
	if d == nil {
		return nature.Nature{}
	}
	return d.peerNature
}

func (d *ExternalDispatch) CallableLanguage() string {
	if d == nil {
		return ""
	}
	return d.callableLanguage
}

func (d *ExternalDispatch) SetCallableLanguage(language string) {
	if d == nil {
		return
	}
	d.callableLanguage = language
}

func (d *ExternalDispatch) Send(delivery being.DeliveredImpulse) being.ChannelResult {
	if d == nil {
		return being.ChannelResult{DropReason: ErrUnknownPeer.Error()}
	}
	if delivery.Raw == "" {
		return being.ChannelResult{DropReason: ErrInvalidImpulse.Error()}
	}

	d.lastExpression = delivery.Parsed.Expression

	return being.ChannelResult{Routed: true}
}

func (d *ExternalDispatch) CurrentExpression() string {
	if d == nil {
		return ""
	}
	return d.lastExpression
}

func (d *ExternalDispatch) DerivePresent(_ *being.Being, _ *being.Being) string {
	return d.lastExpression
}
