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

// impulse parsing

func TestParseImpulseWithExpressionAndFlags(t *testing.T) {
	impulse, err := being.ParseImpulse("skyra bob hello there ~close ~custom | testing expression and flags")
	if err != nil {
		t.Fatalf("ParseImpulse() error = %v", err)
	}
	if impulse.TargetName != "bob" {
		t.Fatalf("TargetName = %q, want %q", impulse.TargetName, "bob")
	}
	if impulse.Expression != "hello there" {
		t.Fatalf("Expression = %q, want %q", impulse.Expression, "hello there")
	}
	if len(impulse.Flags) != 2 {
		t.Fatalf("len(Flags) = %d, want 2", len(impulse.Flags))
	}
	if !impulse.IsClose() {
		t.Fatalf("IsClose() = false, want true")
	}
	if impulse.Reason != "testing expression and flags" {
		t.Fatalf("Reason = %q, want %q", impulse.Reason, "testing expression and flags")
	}
}

func TestParseImpulseAllowsCloseWithoutExpression(t *testing.T) {
	impulse, err := being.ParseImpulse("skyra bob ~close | closing exchange")
	if err != nil {
		t.Fatalf("ParseImpulse() error = %v", err)
	}
	if impulse.Expression != "" {
		t.Fatalf("Expression = %q, want empty", impulse.Expression)
	}
	if !impulse.IsClose() {
		t.Fatalf("IsClose() = false, want true")
	}
}

// ExchangeStack

func TestExchangeStackStoresExplicitExchanges(t *testing.T) {
	channel, err := NewExchangeStack("bob", makeNature("friend", "relate"))
	if err != nil {
		t.Fatalf("NewExchangeStack() error = %v", err)
	}

	first := mustImpulse(t, "skyra bob hello | opening")
	close := mustImpulse(t, "skyra bob ~close | done")
	second := mustImpulse(t, "skyra bob new run | resuming")

	deliver := func(i being.Impulse) being.DeliveredImpulse {
		return being.DeliveredImpulse{Raw: i, Parsed: mustParse(t, i)}
	}

	if result := channel.Send(deliver(first)); !result.Routed {
		t.Fatalf("Send(first) Routed = false, want true")
	}
	if result := channel.Send(deliver(close)); !result.Routed {
		t.Fatalf("Send(close) Routed = false, want true")
	}
	if result := channel.Send(deliver(second)); !result.Routed {
		t.Fatalf("Send(second) Routed = false, want true")
	}

	exchanges := channel.Exchanges()
	if len(exchanges) != 2 {
		t.Fatalf("len(Exchanges()) = %d, want 2", len(exchanges))
	}
	if len(exchanges[0]) != 2 {
		t.Fatalf("len(Exchanges()[0]) = %d, want 2", len(exchanges[0]))
	}
	if len(exchanges[1]) != 1 {
		t.Fatalf("len(Exchanges()[1]) = %d, want 1", len(exchanges[1]))
	}

	open := channel.CurrentOpenExchange()
	if len(open) != 1 {
		t.Fatalf("len(CurrentOpenExchange()) = %d, want 1", len(open))
	}
	parsed, err := open[0].Impulse.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if parsed.Expression != "new run" {
		t.Fatalf("CurrentOpenExchange()[0].Expression = %q, want %q", parsed.Expression, "new run")
	}
}

func TestExchangeStackSwapsTargetToOriginForReceiverView(t *testing.T) {
	channel, err := NewExchangeStack("michael", makeNature("builder", "hold the line"))
	if err != nil {
		t.Fatalf("NewExchangeStack() error = %v", err)
	}

	raw := mustImpulse(t, "skyra skyra hello there | hello")
	result := channel.Send(being.DeliveredImpulse{
		OriginName: "michael",
		Raw:        raw,
		Parsed:     mustParse(t, raw),
	})
	if !result.Routed {
		t.Fatalf("Send() Routed = false, want true")
	}

	exchanges := channel.Exchanges()
	if len(exchanges) != 1 || len(exchanges[0]) != 1 {
		t.Fatalf("Exchanges() shape = %#v, want one exchange with one impulse", exchanges)
	}
	if string(exchanges[0][0].Impulse) != "skyra michael hello there | hello" {
		t.Fatalf("stored impulse = %q, want %q", string(exchanges[0][0].Impulse), "skyra michael hello there | hello")
	}
}

func TestExchangeStackRejectsCloseWithoutOpenExchange(t *testing.T) {
	channel, err := NewExchangeStack("bob", makeNature("friend", "relate"))
	if err != nil {
		t.Fatalf("NewExchangeStack() error = %v", err)
	}

	closeImpulse := mustImpulse(t, "skyra bob ~close | closing")
	result := channel.Send(being.DeliveredImpulse{Raw: closeImpulse, Parsed: mustParse(t, closeImpulse)})
	if result.Routed {
		t.Fatalf("Send(close) Routed = true, want false")
	}
	if result.DropReason == "" {
		t.Fatalf("DropReason = empty, want value")
	}
	if len(channel.Exchanges()) != 0 {
		t.Fatalf("len(Exchanges()) = %d, want 0", len(channel.Exchanges()))
	}
}

func TestExchangeStackDerivePresentUsesNameNatureAndCurrentOpenExchange(t *testing.T) {
	b := newBeing(t, "michael", "builder", "hold the line", true)
	sender := newBeing(t, "skyra", "system", "relate", true)

	if err := SeedPeer(b, sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	open := mustImpulse(t, "skyra skyra relate | initiating")
	if _, err := b.EmitToPeer(sender.Name, open); err != nil {
		t.Fatalf("EmitToPeer() error = %v", err)
	}

	present, err := b.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	want := "your name is: michael\nyour identity is: builder\nyour purpose is: hold the line\n\nyou are in an exchange with: skyra\ntheir identity is: system\ntheir purpose is: relate\n\nyou: relate (your reason for sending this: initiating)\n\nyour cognitive network — beings you can address:\nTo respond, output a single protocol string:\nskyra <being> <what you want to say> | <reason>\n<being> is who you are sending to from the network below\n<what you want to say> is the substance of your expression to that being — carry the message forward\n<reason> is why you are firing this signal\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nskyra\n  identity: system"
	if present != want {
		t.Fatalf("Present = %q, want %q", present, want)
	}
}

func TestClosedTopProducesNoOpenExchangeInPresent(t *testing.T) {
	b := newBeing(t, "michael", "builder", "hold the line", true)
	sender := newBeing(t, "skyra", "system", "relate", true)

	if err := SeedPeer(b, sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	impulse := mustImpulse(t, "skyra skyra hello | hello")
	close := mustImpulse(t, "skyra skyra ~close | closing")
	_, _ = b.EmitToPeer(sender.Name, impulse)
	_, _ = b.EmitToPeer(sender.Name, close)

	present, err := b.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}
	want := "your name is: michael\nyour identity is: builder\nyour purpose is: hold the line\n\nyou are in an exchange with: skyra\ntheir identity is: system\ntheir purpose is: relate\n\nyour cognitive network — beings you can address:\nTo respond, output a single protocol string:\nskyra <being> <what you want to say> | <reason>\n<being> is who you are sending to from the network below\n<what you want to say> is the substance of your expression to that being — carry the message forward\n<reason> is why you are firing this signal\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nskyra\n  identity: system"
	if present != want {
		t.Fatalf("Present = %q, want %q", present, want)
	}
}

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

func TestExchangeStackDerivePresentSortsRelationshipsByPeerName(t *testing.T) {
	receiver := newBeing(t, "michael", "builder", "hold the line", true)
	sender := newBeing(t, "zoe", "late", "second", true)
	early := newBeing(t, "adam", "early", "first", true)

	if err := SeedPeer(receiver, sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer(sender) error = %v", err)
	}
	if err := SeedPeer(receiver, early.Name, early.Nature); err != nil {
		t.Fatalf("SeedPeer(early) error = %v", err)
	}

	impulse := mustImpulse(t, "skyra zoe hello | greeting")
	if _, err := receiver.EmitToPeer(sender.Name, impulse); err != nil {
		t.Fatalf("EmitToPeer() error = %v", err)
	}

	present, err := receiver.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	want := "your name is: michael\nyour identity is: builder\nyour purpose is: hold the line\n\nyou are in an exchange with: zoe\ntheir identity is: late\ntheir purpose is: second\n\nyou: hello (your reason for sending this: greeting)\n\nyour cognitive network — beings you can address:\nTo respond, output a single protocol string:\nskyra <being> <what you want to say> | <reason>\n<being> is who you are sending to from the network below\n<what you want to say> is the substance of your expression to that being — carry the message forward\n<reason> is why you are firing this signal\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nadam\n  identity: early\nzoe\n  identity: late"
	if present != want {
		t.Fatalf("Present = %q, want %q", present, want)
	}
}

func TestExchangeStackDerivePresentUsesExpressionOnlyForNonCognitiveReceiver(t *testing.T) {
	receiver := newBeing(t, "sensor", "sensor", "receive", false)
	sender := newBeing(t, "michael", "builder", "hold the line", true)

	if err := SeedPeer(receiver, sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	impulse := mustImpulse(t, "skyra michael heat spike ~custom | spike detected")
	if _, err := receiver.EmitToPeer(sender.Name, impulse); err != nil {
		t.Fatalf("EmitToPeer() error = %v", err)
	}

	present, err := receiver.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}
	if present != "heat spike" {
		t.Fatalf("Present = %q, want %q", present, "heat spike")
	}
}

func TestSignalImpulseCanBeRoundTripped(t *testing.T) {
	raw := "skyra skyra relate to me ~close | closing"
	impulse := mustImpulse(t, raw)

	parsed, err := being.ParseImpulse(impulse.Raw())
	if err != nil {
		t.Fatalf("ParseImpulse() error = %v", err)
	}
	if parsed.TargetName != "skyra" {
		t.Fatalf("TargetName = %q, want %q", parsed.TargetName, "skyra")
	}
	if !parsed.IsClose() {
		t.Fatalf("IsClose() = false, want true")
	}
	if parsed.Reason != "closing" {
		t.Fatalf("Reason = %q, want %q", parsed.Reason, "closing")
	}
}
