package world

import (
	"fmt"
	"testing"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/identity"
	"skyra-v03/src/primitives/language"
	"skyra-v03/src/primitives/nature"
	"skyra-v03/src/primitives/purpose"
)

// helpers

func makeNature(id, p string) nature.Nature {
	return nature.Nature{
		Identity: identity.Identity{Value: id},
		Purpose:  purpose.Purpose{Value: p},
	}
}

func makeLang(v string) language.Language {
	return language.Language{Value: v}
}

func newBeing(t *testing.T, name, id, p string, cognitive bool) *being.Being {
	t.Helper()
	b, err := being.NewBeing(name, makeNature(id, p), makeLang("skyra being"), cognitive)
	if err != nil {
		t.Fatalf("NewBeing(%q) error = %v", name, err)
	}
	return b
}

func mustImpulse(t *testing.T, raw string) being.Impulse {
	t.Helper()
	i, err := being.NewImpulse(raw)
	if err != nil {
		t.Fatalf("NewImpulse(%q) error = %v", raw, err)
	}
	return i
}

func mustParse(t *testing.T, i being.Impulse) being.ParsedImpulse {
	t.Helper()
	p, err := i.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	return p
}

// ExchangeMap

func TestExchangeMapOpenAndSend(t *testing.T) {
	channel, err := NewExchangeMap("bob", makeNature("friend", "relate"))
	if err != nil {
		t.Fatalf("NewExchangeMap() error = %v", err)
	}

	if err := channel.OpenExchange("intent-1", "test intent", "test reason", nil); err != nil {
		t.Fatalf("OpenExchange() error = %v", err)
	}

	raw := mustImpulse(t, "skyra bob hello | opening")
	result := channel.Send(being.DeliveredImpulse{
		OriginName: "alice",
		ThreadID:   "intent-1",
		Raw:        raw,
		Parsed:     mustParse(t, raw),
	})
	if !result.Routed {
		t.Fatalf("Send() Routed = false, want true")
	}

	exchange, ok := channel.ExchangeByThread("intent-1")
	if !ok {
		t.Fatalf("ExchangeByThread() ok = false, want true")
	}
	if len(exchange.Entries) != 1 {
		t.Fatalf("len(exchange.Entries) = %d, want 1", len(exchange.Entries))
	}
}

func TestExchangeMapSendAutoOpensOnFirstReceive(t *testing.T) {
	channel, err := NewExchangeMap("bob", makeNature("friend", "relate"))
	if err != nil {
		t.Fatalf("NewExchangeMap() error = %v", err)
	}

	raw := mustImpulse(t, "skyra bob hello | opening")
	result := channel.Send(being.DeliveredImpulse{
		OriginName: "alice",
		ThreadID:   "intent-1",
		Raw:        raw,
		Parsed:     mustParse(t, raw),
	})
	if !result.Routed {
		t.Fatalf("Send() Routed = false, want true (receiver should auto-open)")
	}

	exchange, ok := channel.ExchangeByThread("intent-1")
	if !ok {
		t.Fatalf("ExchangeByThread() ok = false after auto-open")
	}
	if len(exchange.Entries) != 1 {
		t.Fatalf("len(exchange.Entries) = %d, want 1", len(exchange.Entries))
	}
}

func TestExchangeMapSendDropsWithoutThreadID(t *testing.T) {
	channel, err := NewExchangeMap("bob", makeNature("friend", "relate"))
	if err != nil {
		t.Fatalf("NewExchangeMap() error = %v", err)
	}
	if err := channel.OpenExchange("intent-1", "test intent", "test reason", nil); err != nil {
		t.Fatalf("OpenExchange() error = %v", err)
	}

	raw := mustImpulse(t, "skyra bob hello | opening")
	result := channel.Send(being.DeliveredImpulse{
		OriginName: "alice",
		Raw:        raw,
		Parsed:     mustParse(t, raw),
	})
	if result.Routed {
		t.Fatalf("Send() Routed = true, want false")
	}
}

func TestExchangeMapHasOpenExchanges(t *testing.T) {
	channel, err := NewExchangeMap("bob", makeNature("friend", "relate"))
	if err != nil {
		t.Fatalf("NewExchangeMap() error = %v", err)
	}

	if channel.HasOpenExchanges() {
		t.Fatalf("HasOpenExchanges() = true, want false before open")
	}

	if err := channel.OpenExchange("intent-1", "test intent", "test reason", nil); err != nil {
		t.Fatalf("OpenExchange() error = %v", err)
	}
	if !channel.HasOpenExchanges() {
		t.Fatalf("HasOpenExchanges() = false, want true after open")
	}

	if err := channel.CloseExchange("intent-1"); err != nil {
		t.Fatalf("CloseExchange() error = %v", err)
	}
	if channel.HasOpenExchanges() {
		t.Fatalf("HasOpenExchanges() = true, want false after close")
	}
}

func TestExchangeMapParallelExchangesSamePeer(t *testing.T) {
	channel, err := NewExchangeMap("bob", makeNature("friend", "relate"))
	if err != nil {
		t.Fatalf("NewExchangeMap() error = %v", err)
	}

	if err := channel.OpenExchange("intent-1", "test intent", "test reason", nil); err != nil {
		t.Fatalf("OpenExchange(intent-1) error = %v", err)
	}
	if err := channel.OpenExchange("intent-2", "test intent 2", "test reason 2", nil); err != nil {
		t.Fatalf("OpenExchange(intent-2) error = %v", err)
	}

	r1 := mustImpulse(t, "skyra bob first thread | thread one")
	r2 := mustImpulse(t, "skyra bob second thread | thread two")

	channel.Send(being.DeliveredImpulse{OriginName: "alice", ThreadID: "intent-1", Raw: r1, Parsed: mustParse(t, r1)})
	channel.Send(being.DeliveredImpulse{OriginName: "alice", ThreadID: "intent-2", Raw: r2, Parsed: mustParse(t, r2)})

	e1, ok := channel.ExchangeByThread("intent-1")
	if !ok || len(e1.Entries) != 1 {
		t.Fatalf("intent-1 exchange len = %d, want 1", len(e1.Entries))
	}
	e2, ok := channel.ExchangeByThread("intent-2")
	if !ok || len(e2.Entries) != 1 {
		t.Fatalf("intent-2 exchange len = %d, want 1", len(e2.Entries))
	}
}

func TestExchangeMapOpenDuplicateIntentErrors(t *testing.T) {
	channel, err := NewExchangeMap("bob", makeNature("friend", "relate"))
	if err != nil {
		t.Fatalf("NewExchangeMap() error = %v", err)
	}
	if err := channel.OpenExchange("intent-1", "test intent", "test reason", nil); err != nil {
		t.Fatalf("first OpenExchange() error = %v", err)
	}
	if err := channel.OpenExchange("intent-1", "test intent", "test reason", nil); err == nil {
		t.Fatalf("second OpenExchange() error = nil, want error")
	}
}

// DerivePresent

func TestExchangeMapDerivePresentShowsEngageWhenNoOpenExchange(t *testing.T) {
	receiver := newBeing(t, "michael", "builder", "hold the line", true)
	sender := newBeing(t, "skyra", "system", "relate", true)

	if err := SeedPeer(receiver, sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	present, err := receiver.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	if present == "" {
		t.Fatalf("DerivePresent() = empty")
	}
	if !contains(present, "skyra start-exchange ~with <process>") {
		t.Fatalf("present missing start-exchange syntax, got:\n%s", present)
	}
}

func TestExchangeMapDerivePresentShowsCallableLanguageWhenExchangeOpen(t *testing.T) {
	receiver := newBeing(t, "michael", "builder", "hold the line", true)
	sender := newBeing(t, "skyra", "system", "relate", true)

	if err := SeedPeer(receiver, sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	peer, ok := receiver.PeerByName(sender.Name)
	if !ok {
		t.Fatalf("PeerByName() ok = false")
	}
	em, ok := peer.(*ExchangeMap)
	if !ok {
		t.Fatalf("peer type = %T, want *ExchangeMap", peer)
	}
	em.SetCallableLanguage("skyra skyra <expression> | <reason>")
	if err := em.OpenExchange("intent-1", "test intent", "test reason", nil); err != nil {
		t.Fatalf("OpenExchange() error = %v", err)
	}

	present, err := receiver.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	if !contains(present, "language: skyra skyra") {
		t.Fatalf("present missing callable language, got:\n%s", present)
	}
	if !contains(present, "continue: skyra skyra <expression>") {
		t.Fatalf("present missing continue syntax, got:\n%s", present)
	}
	if !contains(present, "resolve: skyra close-exchange ~with skyra") {
		t.Fatalf("present missing close-exchange operator syntax, got:\n%s", present)
	}
	if contains(present, "— your network —") && contains(present, "\nskyra\n") {
		t.Fatalf("present should not show skyra in network section when exchange is open, got:\n%s", present)
	}
}

func TestExchangeMapDerivePresentShowsLatestInboundFromSender(t *testing.T) {
	receiver := newBeing(t, "prefrontal", "system", "relate", true)
	sender := newBeing(t, "michael", "builder", "hold the line", true)

	if err := SeedPeer(receiver, sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	peer, ok := receiver.PeerByName(sender.Name)
	if !ok {
		t.Fatalf("PeerByName() ok = false")
	}
	em, ok := peer.(*ExchangeMap)
	if !ok {
		t.Fatalf("peer type = %T, want *ExchangeMap", peer)
	}

	first := mustImpulse(t, "skyra michael hi | opening")
	em.Send(being.DeliveredImpulse{
		OriginName: sender.Name,
		ThreadID:   "intent-1",
		Raw:        first,
		Parsed:     mustParse(t, first),
	})

	reply := mustImpulse(t, "skyra michael hello back | replying")
	em.Send(being.DeliveredImpulse{
		OriginName: receiver.Name,
		ThreadID:   "intent-1",
		Raw:        reply,
		Parsed:     mustParse(t, reply),
	})

	latest := mustImpulse(t, "skyra michael What is your name? | asking")
	em.Send(being.DeliveredImpulse{
		OriginName: sender.Name,
		ThreadID:   "intent-1",
		Raw:        latest,
		Parsed:     mustParse(t, latest),
	})

	present, err := receiver.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	if !contains(present, "[0] michael: hi") {
		t.Fatalf("present missing first transcript entry, got:\n%s", present)
	}
	if !contains(present, "[1] you: hello back") {
		t.Fatalf("present missing reply transcript entry, got:\n%s", present)
	}
	if !contains(present, "latest inbound from michael:\n\n[2] michael: What is your name?") {
		t.Fatalf("present missing inline latest inbound marker, got:\n%s", present)
	}
}

func TestExchangeMapDerivePresentSortsPeersByName(t *testing.T) {
	receiver := newBeing(t, "michael", "builder", "hold the line", true)
	sender := newBeing(t, "zoe", "late", "second", true)
	early := newBeing(t, "adam", "early", "first", true)

	if err := SeedPeer(receiver, sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer(sender) error = %v", err)
	}
	if err := SeedPeer(receiver, early.Name, early.Nature); err != nil {
		t.Fatalf("SeedPeer(early) error = %v", err)
	}

	present, err := receiver.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	// search only in the network section, after the separator
	networkSection := present
	if sepIdx := indexOf(present, "________________"); sepIdx >= 0 {
		networkSection = present[sepIdx:]
	}
	adamIdx := indexOf(networkSection, "\nadam\n")
	zoeIdx := indexOf(networkSection, "\nzoe\n")
	if adamIdx == -1 || zoeIdx == -1 {
		t.Fatalf("present missing peers in network section, got:\n%s", present)
	}
	if adamIdx > zoeIdx {
		t.Fatalf("adam should appear before zoe in network section")
	}
}

// ExternalDispatch

func TestExternalDispatchDerivesExpressionOnly(t *testing.T) {
	b := newBeing(t, "sensor", "sensor", "receive", false)

	channel, err := NewExternalDispatch("world", makeNature("world", "emit"))
	if err != nil {
		t.Fatalf("NewExternalDispatch() error = %v", err)
	}
	if err := b.AttachPeer(channel); err != nil {
		t.Fatalf("AttachPeer() error = %v", err)
	}

	impulse := mustImpulse(t, "skyra world incoming signal | observing")
	if _, err := b.EmitToPeer("world", impulse); err != nil {
		t.Fatalf("EmitToPeer() error = %v", err)
	}

	sender := newBeing(t, "world", "world", "emit", false)
	present, err := b.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}
	if present != "incoming signal" {
		t.Fatalf("Present = %q, want %q", present, "incoming signal")
	}
}

// World

func TestWorldRelatesBeingsAndResolvesByName(t *testing.T) {
	w := New()

	left := newBeing(t, "michael", "builder", "hold the line", true)
	right := newBeing(t, "skyra", "system", "relate", true)

	if err := w.Register(left); err != nil {
		t.Fatalf("Register(left) error = %v", err)
	}
	if err := w.Register(right); err != nil {
		t.Fatalf("Register(right) error = %v", err)
	}
	if _, err := w.Grow(fmt.Sprintf("skyra being ~name %s ~relationships %s", left.Name, right.Name)); err != nil {
		t.Fatalf("Grow(relate) error = %v", err)
	}

	resolved, ok := w.BeingByName("michael")
	if !ok {
		t.Fatalf("BeingByName() ok = false, want true")
	}
	if resolved.Name != left.Name {
		t.Fatalf("BeingByName() = %q, want %q", resolved.Name, left.Name)
	}

	leftPeer, ok := left.PeerByName(right.Name)
	if !ok {
		t.Fatalf("left PeerByName() ok = false, want true")
	}
	if leftPeer.Name() != "skyra" {
		t.Fatalf("left peer name = %q, want %q", leftPeer.Name(), "skyra")
	}

	rightPeer, ok := right.PeerByName(left.Name)
	if !ok {
		t.Fatalf("right PeerByName() ok = false, want true")
	}
	if rightPeer.Name() != "michael" {
		t.Fatalf("right peer name = %q, want %q", rightPeer.Name(), "michael")
	}
}

// helpers

func contains(s, substr string) bool {
	return len(s) >= len(substr) && (s == substr || len(s) > 0 && indexOf(s, substr) >= 0)
}

func indexOf(s, substr string) int {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return i
		}
	}
	return -1
}
