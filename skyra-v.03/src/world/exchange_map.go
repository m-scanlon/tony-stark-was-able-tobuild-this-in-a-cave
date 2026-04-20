package world

import (
	"fmt"
	"sort"
	"strings"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/nature"
)

type ExchangeThread struct {
	About          string
	Because        string
	IsOpener       bool
	ContextEntries []being.ExchangeEntry
	Entries        []being.ExchangeEntry
}

type ExchangeMap struct {
	peerName         string
	peerNature       nature.Nature
	exchanges        map[string]ExchangeThread // keyed by thread_id
	callableLanguage string
	lastActiveThread string
}

func NewExchangeMap(peerName string, peerNature nature.Nature) (*ExchangeMap, error) {
	if strings.TrimSpace(peerName) == "" {
		return nil, fmt.Errorf("%w: peer name is required", ErrUnknownPeer)
	}
	if err := peerNature.Validate(); err != nil {
		return nil, err
	}
	return &ExchangeMap{
		peerName:   peerName,
		peerNature: peerNature,
		exchanges:  make(map[string]ExchangeThread),
	}, nil
}

func (c *ExchangeMap) Name() string {
	if c == nil {
		return ""
	}
	return c.peerName
}

func (c *ExchangeMap) PeerNature() nature.Nature {
	if c == nil {
		return nature.Nature{}
	}
	return c.peerNature
}

func (c *ExchangeMap) CallableLanguage() string {
	if c == nil {
		return ""
	}
	return c.callableLanguage
}

func (c *ExchangeMap) SetCallableLanguage(language string) {
	if c == nil {
		return
	}
	c.callableLanguage = language
}

func (c *ExchangeMap) HasOpenExchanges() bool {
	if c == nil {
		return false
	}
	return len(c.exchanges) > 0
}

func (c *ExchangeMap) OpenExchange(threadID, about, because string, contextEntries []being.ExchangeEntry) error {
	if strings.TrimSpace(threadID) == "" {
		return fmt.Errorf("exchange map: thread_id is required to open an exchange")
	}
	if _, exists := c.exchanges[threadID]; exists {
		return fmt.Errorf("exchange map: exchange already open for thread %s", threadID)
	}
	c.exchanges[threadID] = ExchangeThread{About: about, Because: because, IsOpener: true, ContextEntries: contextEntries}
	return nil
}

func (c *ExchangeMap) CloseExchange(threadID string) error {
	if _, exists := c.exchanges[threadID]; !exists {
		return fmt.Errorf("exchange map: no exchange open for thread %s", threadID)
	}
	delete(c.exchanges, threadID)
	return nil
}

func (c *ExchangeMap) Send(delivery being.DeliveredImpulse) being.ChannelResult {
	if c == nil {
		return being.ChannelResult{DropReason: ErrUnknownPeer.Error()}
	}
	if delivery.Raw == "" {
		return being.ChannelResult{DropReason: ErrInvalidImpulse.Error()}
	}
	if strings.TrimSpace(delivery.ThreadID) == "" {
		return being.ChannelResult{DropReason: "exchange map: thread_id is required"}
	}
	thread, exists := c.exchanges[delivery.ThreadID]
	if !exists {
		thread = ExchangeThread{IsOpener: false, About: delivery.About, Because: delivery.Because, ContextEntries: delivery.ContextEntries}
	}
	thread.Entries = append(thread.Entries, being.ExchangeEntry{
		Author:  delivery.OriginName,
		Impulse: delivery.Raw,
	})
	c.exchanges[delivery.ThreadID] = thread
	c.lastActiveThread = delivery.ThreadID
	return being.ChannelResult{Routed: true}
}

func (c *ExchangeMap) ExchangeByThread(threadID string) (ExchangeThread, bool) {
	t, ok := c.exchanges[threadID]
	return t, ok
}

func (c *ExchangeMap) DerivePresent(receiver *being.Being, sender *being.Being) string {
	var builder strings.Builder
	builder.WriteString("your name is: ")
	builder.WriteString(receiver.Name)
	builder.WriteString("\nyour identity is: ")
	builder.WriteString(receiver.Nature.Identity.Value)
	builder.WriteString("\nyour purpose is: ")
	builder.WriteString(receiver.Nature.Purpose.Value)

	builder.WriteString("\n\nyou received a signal from: ")
	builder.WriteString(sender.Name)
	builder.WriteString("\ntheir identity is: ")
	builder.WriteString(sender.Nature.Identity.Value)
	builder.WriteString("\ntheir purpose is: ")
	builder.WriteString(sender.Nature.Purpose.Value)

	latestInbound, hasLatestInbound := c.latestInboundFrom(sender.Name)

	// open exchange threads — context, before the network
	var openPeers []being.RelationshipChannel
	for _, peer := range sortedPeers(receiver.Peers) {
		if peer.Name() == "start-exchange" || peer.Name() == "close-exchange" {
			continue
		}
		if pm, ok := peer.(*ExchangeMap); ok && pm.HasOpenExchanges() {
			openPeers = append(openPeers, peer)
		}
	}

	if len(openPeers) > 0 {
		for _, peer := range openPeers {
			pm := peer.(*ExchangeMap)
			for _, threadID := range sortedKeys(pm.exchanges) {
				thread := pm.exchanges[threadID]
				if thread.IsOpener {
					builder.WriteString(fmt.Sprintf("\n\nyou started an exchange with %s", peer.Name()))
					builder.WriteString(fmt.Sprintf("\nabout: %s", thread.About))
					builder.WriteString(fmt.Sprintf("\nbecause: %s", thread.Because))
					if len(thread.ContextEntries) > 0 {
						builder.WriteString(fmt.Sprintf("\n\ncontext %s included from a different exchange:", peer.Name()))
						for i, entry := range thread.ContextEntries {
							builder.WriteString(fmt.Sprintf("\n\n[%d] %s: %s", i, entry.Author, formatPresentImpulse(receiver, entry.Impulse, entry.Author == receiver.Name)))
						}
					}
				} else {
					if thread.About != "" {
						builder.WriteString(fmt.Sprintf("\n\n%s wants to talk about %s because %s.", peer.Name(), thread.About, thread.Because))
					}
					if len(thread.ContextEntries) > 0 {
						builder.WriteString("\nHere is the context shared from a previous thread:")
						for _, entry := range thread.ContextEntries {
							builder.WriteString(fmt.Sprintf("\n\n%s: %s", entry.Author, formatPresentImpulse(receiver, entry.Impulse, entry.Author == receiver.Name)))
						}
					}
				}
				for i, entry := range thread.Entries {
					if hasLatestInbound && peer.Name() == sender.Name && threadID == latestInbound.ThreadID && i == latestInbound.Index {
						builder.WriteString("\n\nlatest inbound from ")
						builder.WriteString(sender.Name)
						builder.WriteString(":")
					}
					builder.WriteString(fmt.Sprintf("\n\n[%d] ", i))
					isOwn := entry.Author == receiver.Name
					if isOwn {
						builder.WriteString("you")
					} else {
						builder.WriteString(entry.Author)
					}
					builder.WriteString(": ")
					builder.WriteString(formatPresentImpulse(receiver, entry.Impulse, isOwn))
				}
			}
		}
	}

	// cognitive network
	builder.WriteString("\n\n— your cognitive processes —")
	builder.WriteString("\nThese are not separate agents — they are your own faculties. Calling on them is internal deliberation, not delegation.")
	builder.WriteString("\nRespond with a single protocol string only — no explanation, no markdown, no extra text")

	if len(openPeers) > 0 {
		builder.WriteString("\n\nactive deliberations")
		builder.WriteString("\ncontinue deliberating using: skyra <process> <expression> | <reason>")
		builder.WriteString("\n________________")
		for _, peer := range openPeers {
			pm := peer.(*ExchangeMap)
			builder.WriteString("\n")
			builder.WriteString(peer.Name())
			builder.WriteString(fmt.Sprintf("\n  continue: skyra %s <expression> [~expression-reference <start-end>] | <reason>", peer.Name()))
			if lang := pm.CallableLanguage(); lang != "" {
				builder.WriteString("\n  language: ")
				builder.WriteString(lang)
			}
			if pm.hasOpenerExchange() {
				builder.WriteString(fmt.Sprintf("\n  resolve: skyra close-exchange ~with %s ~learned <synthesis> ~expression-reference <start-end> | <reason>", peer.Name()))
			}
		}
	}

	builder.WriteString("\n\nyour processes — begin deliberation using: skyra start-exchange ~with <process> ~about <string> ~because <sentence> ~expression-reference <start-end> ~say <expression> | <reason>")
	builder.WriteString("\n________________")

	for _, peer := range sortedPeers(receiver.Peers) {
		if peer.Name() == "start-exchange" || peer.Name() == "close-exchange" {
			continue
		}
		pm, isCognitive := peer.(*ExchangeMap)
		if isCognitive && pm.HasOpenExchanges() {
			continue
		}
		builder.WriteString("\n")
		builder.WriteString(peer.Name())
		if callable := peer.PeerNature().Callable.Value; callable != "" {
			builder.WriteString("\n  deliberate when: ")
			builder.WriteString(callable)
		}
		builder.WriteString("\n  identity: ")
		builder.WriteString(peer.PeerNature().Identity.Value)
	}
	return builder.String()
}

type latestInboundEntry struct {
	ThreadID string
	Index    int
	Entry    being.ExchangeEntry
}

func (c *ExchangeMap) hasOpenerExchange() bool {
	for _, thread := range c.exchanges {
		if thread.IsOpener {
			return true
		}
	}
	return false
}

func (c *ExchangeMap) latestInboundFrom(senderName string) (latestInboundEntry, bool) {
	if c == nil || strings.TrimSpace(senderName) == "" || strings.TrimSpace(c.lastActiveThread) == "" {
		return latestInboundEntry{}, false
	}

	thread, ok := c.exchanges[c.lastActiveThread]
	if !ok {
		return latestInboundEntry{}, false
	}

	for i := len(thread.Entries) - 1; i >= 0; i-- {
		if thread.Entries[i].Author == senderName {
			return latestInboundEntry{ThreadID: c.lastActiveThread, Index: i, Entry: thread.Entries[i]}, true
		}
	}

	return latestInboundEntry{}, false
}

func sortedKeys(exchanges map[string]ExchangeThread) []string {
	keys := make([]string, 0, len(exchanges))
	for k := range exchanges {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
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
		return parsed.Expression
	}

	return parsed.Expression
}
