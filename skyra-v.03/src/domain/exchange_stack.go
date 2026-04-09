package domain

import (
	"fmt"
	"sort"
	"strings"
	"unicode"
)

type Exchange []Impulse

func (e Exchange) IsClosed() bool {
	if len(e) == 0 {
		return false
	}
	return e[len(e)-1].IsClose()
}

func (e Exchange) Impulses() []Impulse {
	return append([]Impulse(nil), e...)
}

type ExchangeStack struct {
	peerName   string
	peerNature Nature
	stack      []Exchange
}

func NewExchangeStack(peerName string, peerNature Nature) (*ExchangeStack, error) {
	if strings.TrimSpace(peerName) == "" {
		return nil, fmt.Errorf("%w: peer name is required", ErrUnknownPeer)
	}
	if err := peerNature.Validate(); err != nil {
		return nil, err
	}
	return &ExchangeStack{
		peerName:   peerName,
		peerNature: peerNature,
		stack:      make([]Exchange, 0),
	}, nil
}

func (c *ExchangeStack) Name() string {
	if c == nil {
		return ""
	}
	return c.peerName
}

func (c *ExchangeStack) PeerNature() Nature {
	if c == nil {
		return Nature{}
	}
	return c.peerNature
}

func (c *ExchangeStack) Send(delivery DeliveredImpulse) ChannelResult {
	if c == nil {
		return ChannelResult{DropReason: ErrUnknownPeer.Error()}
	}
	if delivery.Raw == "" {
		return ChannelResult{DropReason: ErrInvalidImpulse.Error()}
	}

	stored := delivery.Raw
	if shouldSwapTargetToOrigin(delivery, c.peerName) {
		stored = rewriteImpulseTarget(delivery.Raw, delivery.OriginName)
	}

	if delivery.Parsed.IsClose() {
		if !c.HasOpenExchange() {
			return ChannelResult{DropReason: "cannot close without an open exchange"}
		}

		top := len(c.stack) - 1
		c.stack[top] = append(c.stack[top], stored)
		return ChannelResult{
			Routed:      true,
			NewExchange: false,
		}
	}

	if c.HasOpenExchange() {
		top := len(c.stack) - 1
		c.stack[top] = append(c.stack[top], stored)
		return ChannelResult{
			Routed:      true,
			NewExchange: false,
		}
	}

	c.stack = append(c.stack, Exchange{stored})
	return ChannelResult{
		Routed:      true,
		NewExchange: true,
	}
}

func (c *ExchangeStack) Exchanges() []Exchange {
	if c == nil {
		return nil
	}

	exchanges := make([]Exchange, len(c.stack))
	for i, exchange := range c.stack {
		exchanges[i] = Exchange(exchange.Impulses())
	}
	return exchanges
}

func (c *ExchangeStack) HasOpenExchange() bool {
	if c == nil || len(c.stack) == 0 {
		return false
	}
	return !c.stack[len(c.stack)-1].IsClosed()
}

func (c *ExchangeStack) CurrentOpenExchange() Exchange {
	if !c.HasOpenExchange() {
		return nil
	}
	return Exchange(c.stack[len(c.stack)-1].Impulses())
}

func (c *ExchangeStack) derivePresent(receiver *Being, sender *Being) string {
	var builder strings.Builder
	builder.WriteString("name: ")
	builder.WriteString(receiver.Name)
	builder.WriteString("\nidentity: ")
	builder.WriteString(receiver.Nature.Identity)
	builder.WriteString("\npurpose: ")
	builder.WriteString(receiver.Nature.Purpose)

	builder.WriteString("\n\nYou are in an exchange with: ")
	builder.WriteString(sender.Name)
	builder.WriteString("\nthe identity of ")
	builder.WriteString(sender.Name)
	builder.WriteString(" is: ")
	builder.WriteString(sender.Nature.Identity)
	builder.WriteString("\nthe purpose of ")
	builder.WriteString(sender.Name)
	builder.WriteString(" is: ")
	builder.WriteString(sender.Nature.Purpose)

	open := c.CurrentOpenExchange()
	for _, impulse := range open {
		builder.WriteString("\n\n")
		builder.WriteString(sender.Name)
		builder.WriteString(": ")
		builder.WriteString(formatPresentImpulse(receiver, impulse))
	}

	builder.WriteString("\n\nrelationships:")
	builder.WriteString("\nCall any of your relationships using this syntax-")
	builder.WriteString("\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>")
	builder.WriteString("\n<being> must be one of your relationships listed below")
	builder.WriteString("\n<source> is the being you are currently in exchange with")
	builder.WriteString("\n<reason> is why you are firing this expression")
	builder.WriteString("\nRespond with the protocol string only — no explanation, no markdown, no extra text")
	builder.WriteString("\n________________")
	for _, peer := range sortedPeers(receiver.Peers) {
		builder.WriteString("\n")
		builder.WriteString(peer.PeerNature().Identity)
		builder.WriteString(" - ")
		builder.WriteString(peer.PeerNature().Purpose)
	}
	return builder.String()
}

func sortedPeers(peers map[string]RelationshipChannel) []RelationshipChannel {
	channels := make([]RelationshipChannel, 0, len(peers))
	for _, peer := range peers {
		channels = append(channels, peer)
	}
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Name() < channels[j].Name()
	})
	return channels
}

func formatPresentImpulse(receiver *Being, impulse Impulse) string {
	parsed, err := impulse.Parse()
	if err != nil {
		return impulse.Raw()
	}

	if !receiver.Cognitive {
		if parsed.Expression != "" {
			return parsed.Expression
		}
		return formatFlags(parsed.Flags)
	}

	parts := make([]string, 0, 2)
	if parsed.Expression != "" {
		parts = append(parts, parsed.Expression)
	}
	if flags := formatFlags(parsed.Flags); flags != "" {
		parts = append(parts, flags)
	}
	return strings.Join(parts, " ")
}

func formatFlags(flags []Flag) string {
	parts := make([]string, 0, len(flags))
	for _, flag := range flags {
		if flag == "" {
			continue
		}
		parts = append(parts, "-"+string(flag))
	}
	return strings.Join(parts, " ")
}

func shouldSwapTargetToOrigin(delivery DeliveredImpulse, peerName string) bool {
	if delivery.Parsed.IsClose() {
		return false
	}
	if delivery.OriginName == "" {
		return false
	}
	return delivery.Parsed.TargetName != peerName
}

func rewriteImpulseTarget(raw Impulse, newTarget string) Impulse {
	s := raw.Raw()
	if s == "" || newTarget == "" {
		return raw
	}

	idx := len("skyra")
	idx = skipWhitespace(s, idx)
	targetStart := idx
	for idx < len(s) && !unicode.IsSpace(rune(s[idx])) {
		idx++
	}
	if targetStart >= len(s) || targetStart == idx {
		return raw
	}

	return Impulse(s[:targetStart] + newTarget + s[idx:])
}

func skipWhitespace(s string, idx int) int {
	for idx < len(s) && unicode.IsSpace(rune(s[idx])) {
		idx++
	}
	return idx
}
