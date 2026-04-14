package world

import (
	"fmt"
	"sort"
	"strings"
	"unicode"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/nature"
)

type ExchangeEntry struct {
	Author  string
	Impulse being.Impulse
}

type Exchange []ExchangeEntry

func (e Exchange) IsClosed() bool {
	if len(e) == 0 {
		return false
	}
	return e[len(e)-1].Impulse.IsClose()
}

func (e Exchange) Entries() []ExchangeEntry {
	return append([]ExchangeEntry(nil), e...)
}

type ExchangeStack struct {
	peerName         string
	peerNature       nature.Nature
	stack            []Exchange
	callableLanguage string
}

func NewExchangeStack(peerName string, peerNature nature.Nature) (*ExchangeStack, error) {
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

func (c *ExchangeStack) PeerNature() nature.Nature {
	if c == nil {
		return nature.Nature{}
	}
	return c.peerNature
}

func (c *ExchangeStack) CallableLanguage() string {
	if c == nil {
		return ""
	}
	return c.callableLanguage
}

func (c *ExchangeStack) SetCallableLanguage(language string) {
	if c == nil {
		return
	}
	c.callableLanguage = language
}

func (c *ExchangeStack) Send(delivery being.DeliveredImpulse) being.ChannelResult {
	if c == nil {
		return being.ChannelResult{DropReason: ErrUnknownPeer.Error()}
	}
	if delivery.Raw == "" {
		return being.ChannelResult{DropReason: ErrInvalidImpulse.Error()}
	}

	stored := delivery.Raw
	if shouldSwapTargetToOrigin(delivery, c.peerName) {
		stored = rewriteImpulseTarget(delivery.Raw, delivery.OriginName)
	}

	entry := ExchangeEntry{Author: delivery.OriginName, Impulse: stored}

	if delivery.Parsed.IsClose() {
		if !c.HasOpenExchange() {
			return being.ChannelResult{DropReason: "cannot close without an open exchange"}
		}

		top := len(c.stack) - 1
		c.stack[top] = append(c.stack[top], entry)
		return being.ChannelResult{Routed: true}
	}

	if c.HasOpenExchange() {
		top := len(c.stack) - 1
		c.stack[top] = append(c.stack[top], entry)
		return being.ChannelResult{Routed: true}
	}

	c.stack = append(c.stack, Exchange{entry})
	return being.ChannelResult{Routed: true, NewExchange: true}
}

func (c *ExchangeStack) Exchanges() []Exchange {
	if c == nil {
		return nil
	}
	exchanges := make([]Exchange, len(c.stack))
	for i, exchange := range c.stack {
		exchanges[i] = Exchange(exchange.Entries())
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
	return Exchange(c.stack[len(c.stack)-1].Entries())
}

func (c *ExchangeStack) DerivePresent(receiver *being.Being, sender *being.Being) string {
	var builder strings.Builder
	builder.WriteString("your name is: ")
	builder.WriteString(receiver.Name)
	builder.WriteString("\nyour identity is: ")
	builder.WriteString(receiver.Nature.Identity.Value)
	builder.WriteString("\nyour purpose is: ")
	builder.WriteString(receiver.Nature.Purpose.Value)

	builder.WriteString("\n\nyou are in an exchange with: ")
	builder.WriteString(sender.Name)
	builder.WriteString("\ntheir identity is: ")
	builder.WriteString(sender.Nature.Identity.Value)
	builder.WriteString("\ntheir purpose is: ")
	builder.WriteString(sender.Nature.Purpose.Value)

	open := c.CurrentOpenExchange()
	for _, entry := range open {
		builder.WriteString("\n\n")
		author := entry.Author
		if author == "" {
			author = sender.Name
		}
		isOwn := author == receiver.Name
		if isOwn {
			builder.WriteString("you")
		} else {
			builder.WriteString(author)
		}
		builder.WriteString(": ")
		builder.WriteString(formatPresentImpulse(receiver, entry.Impulse, isOwn))
	}

	builder.WriteString("\n\nyour cognitive network — beings you can address:")
	builder.WriteString("\nTo respond, output a single protocol string:")
	builder.WriteString("\nskyra <being> <what you want to say> | <reason>")
	builder.WriteString("\n<being> is who you are sending to from the network below")
	builder.WriteString("\n<what you want to say> is the substance of your expression to that being — carry the message forward")
	builder.WriteString("\n<reason> is why you are firing this signal")
	builder.WriteString("\nRespond with the protocol string only — no explanation, no markdown, no extra text")
	builder.WriteString("\n________________")
	for _, peer := range sortedPeers(receiver.Peers) {
		builder.WriteString("\n")
		builder.WriteString(peer.Name())
		if callable := peer.PeerNature().Callable.Value; callable != "" {
			builder.WriteString("\n  call me when: ")
			builder.WriteString(callable)
		}
		builder.WriteString("\n  identity: ")
		builder.WriteString(peer.PeerNature().Identity.Value)
		if lang := peer.CallableLanguage(); lang != "" {
			builder.WriteString("\n  language: ")
			builder.WriteString(lang)
		}
	}
	return builder.String()
}

func sortedPeers(peers map[string]being.RelationshipChannel) []being.RelationshipChannel {
	channels := make([]being.RelationshipChannel, 0, len(peers))
	for _, peer := range peers {
		channels = append(channels, peer)
	}
	sort.Slice(channels, func(i, j int) bool {
		return channels[i].Name() < channels[j].Name()
	})
	return channels
}

func formatPresentImpulse(receiver *being.Being, impulse being.Impulse, isOwn bool) string {
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

	parts := make([]string, 0, 3)
	if parsed.Expression != "" {
		parts = append(parts, parsed.Expression)
	}
	if flags := formatFlags(parsed.Flags); flags != "" {
		parts = append(parts, flags)
	}
	if isOwn && parsed.Reason != "" {
		parts = append(parts, "(your reason for sending this: "+parsed.Reason+")")
	}
	return strings.Join(parts, " ")
}

func formatFlags(flags []being.Flag) string {
	parts := make([]string, 0, len(flags))
	for _, flag := range flags {
		if flag == "" {
			continue
		}
		parts = append(parts, "~"+string(flag))
	}
	return strings.Join(parts, " ")
}

func shouldSwapTargetToOrigin(delivery being.DeliveredImpulse, peerName string) bool {
	if delivery.Parsed.IsClose() {
		return false
	}
	if delivery.OriginName == "" {
		return false
	}
	return delivery.Parsed.TargetName != peerName
}

func rewriteImpulseTarget(raw being.Impulse, newTarget string) being.Impulse {
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

	return being.Impulse(s[:targetStart] + newTarget + s[idx:])
}

func skipWhitespace(s string, idx int) int {
	for idx < len(s) && unicode.IsSpace(rune(s[idx])) {
		idx++
	}
	return idx
}
