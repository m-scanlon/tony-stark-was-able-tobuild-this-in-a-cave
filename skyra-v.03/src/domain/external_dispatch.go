package domain

import (
	"fmt"
	"strings"
)

type ExternalDispatch struct {
	peerName       string
	peerNature     Nature
	lastExpression string
}

func NewExternalDispatch(peerName string, peerNature Nature) (*ExternalDispatch, error) {
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

func (d *ExternalDispatch) PeerNature() Nature {
	if d == nil {
		return Nature{}
	}
	return d.peerNature
}

func (d *ExternalDispatch) Send(delivery DeliveredImpulse) ChannelResult {
	if d == nil {
		return ChannelResult{DropReason: ErrUnknownPeer.Error()}
	}
	if delivery.Raw == "" {
		return ChannelResult{DropReason: ErrInvalidImpulse.Error()}
	}

	d.lastExpression = delivery.Parsed.Expression
	if delivery.Parsed.IsClose() {
		d.lastExpression = ""
	}

	return ChannelResult{
		Routed:      true,
		NewExchange: false,
	}
}

func (d *ExternalDispatch) CurrentExpression() string {
	if d == nil {
		return ""
	}
	return d.lastExpression
}

func (d *ExternalDispatch) derivePresent(_ *Being, _ *Being) string {
	return d.lastExpression
}
