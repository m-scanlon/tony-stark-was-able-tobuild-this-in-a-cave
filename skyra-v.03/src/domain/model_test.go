package domain

import "testing"

func TestParseImpulseWithExpressionAndFlags(t *testing.T) {
	impulse, err := ParseImpulse("skyra bob hello there -close -custom | bob: testing expression and flags")
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
	if impulse.Source != "bob" {
		t.Fatalf("Source = %q, want %q", impulse.Source, "bob")
	}
	if impulse.Reason != "testing expression and flags" {
		t.Fatalf("Reason = %q, want %q", impulse.Reason, "testing expression and flags")
	}
}

func TestParseImpulseAllowsCloseWithoutExpression(t *testing.T) {
	impulse, err := ParseImpulse("skyra bob -close | bob: closing exchange")
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

func TestExchangeStackStoresExplicitExchanges(t *testing.T) {
	channel, err := NewExchangeStack("bob", Nature{Identity: "friend", Purpose: "relate"})
	if err != nil {
		t.Fatalf("NewExchangeStack() error = %v", err)
	}

	first, _ := NewImpulse("skyra bob hello | bob: opening")
	close, _ := NewImpulse("skyra bob -close | bob: done")
	second, _ := NewImpulse("skyra bob new run | bob: resuming")

	if result := channel.Send(DeliveredImpulse{Raw: first, Parsed: mustParseImpulse(t, first)}); !result.Routed {
		t.Fatalf("Send(first) Routed = false, want true")
	}
	if result := channel.Send(DeliveredImpulse{Raw: close, Parsed: mustParseImpulse(t, close)}); !result.Routed {
		t.Fatalf("Send(close) Routed = false, want true")
	}
	if result := channel.Send(DeliveredImpulse{Raw: second, Parsed: mustParseImpulse(t, second)}); !result.Routed {
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
	parsed, err := open[0].Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if parsed.Expression != "new run" {
		t.Fatalf("CurrentOpenExchange()[0].Expression = %q, want %q", parsed.Expression, "new run")
	}
}

func TestExchangeStackSwapsTargetToOriginForReceiverView(t *testing.T) {
	channel, err := NewExchangeStack("michael", Nature{Identity: "builder", Purpose: "hold the line"})
	if err != nil {
		t.Fatalf("NewExchangeStack() error = %v", err)
	}

	raw, _ := NewImpulse("skyra skyra hello there | michael: hello")
	result := channel.Send(DeliveredImpulse{
		OriginName: "michael",
		Raw:        raw,
		Parsed:     mustParseImpulse(t, raw),
	})
	if !result.Routed {
		t.Fatalf("Send() Routed = false, want true")
	}

	exchanges := channel.Exchanges()
	if len(exchanges) != 1 || len(exchanges[0]) != 1 {
		t.Fatalf("Exchanges() shape = %#v, want one exchange with one impulse", exchanges)
	}
	if string(exchanges[0][0]) != "skyra michael hello there | michael: hello" {
		t.Fatalf("stored impulse = %q, want %q", string(exchanges[0][0]), "skyra michael hello there | michael: hello")
	}
}

func TestExchangeStackRejectsCloseWithoutOpenExchange(t *testing.T) {
	channel, err := NewExchangeStack("bob", Nature{Identity: "friend", Purpose: "relate"})
	if err != nil {
		t.Fatalf("NewExchangeStack() error = %v", err)
	}

	closeImpulse, _ := NewImpulse("skyra bob -close | bob: closing")
	result := channel.Send(DeliveredImpulse{Raw: closeImpulse, Parsed: mustParseImpulse(t, closeImpulse)})
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
	nature := Nature{
		Identity: "builder",
		Purpose:  "hold the line",
	}

	being, err := NewBeing("michael", nature, true, true)
	if err != nil {
		t.Fatalf("NewBeing() error = %v", err)
	}
	sender, err := NewBeing("skyra", Nature{Identity: "system", Purpose: "relate"}, true, true)
	if err != nil {
		t.Fatalf("NewBeing(sender) error = %v", err)
	}
	if err := being.SeedPeer(sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	open, _ := NewImpulse("skyra skyra relate | michael: initiating")
	if _, err := being.EmitToPeer(sender.Name, open); err != nil {
		t.Fatalf("EmitToPeer() error = %v", err)
	}

	present, err := being.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	want := "name: michael\nidentity: builder\npurpose: hold the line\n\nYou are in an exchange with: skyra\nthe identity of skyra is: system\nthe purpose of skyra is: relate\n\nskyra: relate\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nsystem - relate"
	if present != want {
		t.Fatalf("Present = %q, want %q", present, want)
	}
}

func TestClosedTopProducesNoOpenExchangeInPresent(t *testing.T) {
	being, err := NewBeing(
		"michael",
		Nature{Identity: "builder", Purpose: "hold the line"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing() error = %v", err)
	}
	sender, err := NewBeing("skyra", Nature{Identity: "system", Purpose: "relate"}, true, true)
	if err != nil {
		t.Fatalf("NewBeing(sender) error = %v", err)
	}
	if err := being.SeedPeer(sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	impulse, _ := NewImpulse("skyra skyra hello | michael: hello")
	close, _ := NewImpulse("skyra skyra -close | michael: closing")
	_, _ = being.EmitToPeer(sender.Name, impulse)
	_, _ = being.EmitToPeer(sender.Name, close)

	present, err := being.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}
	want := "name: michael\nidentity: builder\npurpose: hold the line\n\nYou are in an exchange with: skyra\nthe identity of skyra is: system\nthe purpose of skyra is: relate\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nsystem - relate"
	if present != want {
		t.Fatalf("Present = %q, want %q", present, want)
	}
}

func TestExternalDispatchDerivesExpressionOnly(t *testing.T) {
	being, err := NewBeing(
		"sensor",
		Nature{Identity: "sensor", Purpose: "receive"},
		true,
		false,
	)
	if err != nil {
		t.Fatalf("NewBeing() error = %v", err)
	}

	channel, err := NewExternalDispatch("world", Nature{Identity: "world", Purpose: "emit"})
	if err != nil {
		t.Fatalf("NewExternalDispatch() error = %v", err)
	}
	if err := being.AttachPeer(channel); err != nil {
		t.Fatalf("AttachPeer() error = %v", err)
	}

	impulse, _ := NewImpulse("skyra world incoming signal | sensor: observing")
	if _, err := being.EmitToPeer("world", impulse); err != nil {
		t.Fatalf("EmitToPeer() error = %v", err)
	}

	sender, err := NewBeing("world", Nature{Identity: "world", Purpose: "emit"}, false, false)
	if err != nil {
		t.Fatalf("NewBeing(sender) error = %v", err)
	}
	present, err := being.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	if present != "incoming signal" {
		t.Fatalf("Present = %q, want %q", present, "incoming signal")
	}
}

func TestKernelStateSeedsRelationshipAndResolvesByName(t *testing.T) {
	kernelState := NewKernelState()

	left, err := NewBeing(
		"michael",
		Nature{Identity: "builder", Purpose: "hold the line"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(left) error = %v", err)
	}
	right, err := NewBeing(
		"skyra",
		Nature{Identity: "system", Purpose: "relate"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(right) error = %v", err)
	}

	if err := kernelState.InsertBeing(left); err != nil {
		t.Fatalf("InsertBeing(left) error = %v", err)
	}
	if err := kernelState.InsertBeing(right); err != nil {
		t.Fatalf("InsertBeing(right) error = %v", err)
	}
	if err := kernelState.SeedRelationship(left.Name, right.Name); err != nil {
		t.Fatalf("SeedRelationship() error = %v", err)
	}

	resolved, ok := kernelState.BeingByName("michael")
	if !ok {
		t.Fatalf("BeingByName() ok = false, want true")
	}
	if resolved.Name != left.Name {
		t.Fatalf("BeingByName() = %q, want %q", resolved.Name, left.Name)
	}

	leftPeer, ok := left.PeerByName(right.Name)
	if !ok {
		t.Fatalf("PeerByName() ok = false, want true")
	}
	if leftPeer.Name() != "skyra" {
		t.Fatalf("left peer name = %q, want %q", leftPeer.Name(), "skyra")
	}

	rightPeer, ok := right.PeerByName(left.Name)
	if !ok {
		t.Fatalf("PeerByName() ok = false, want true")
	}
	if rightPeer.Name() != "michael" {
		t.Fatalf("right peer name = %q, want %q", rightPeer.Name(), "michael")
	}
}

func TestExchangeStackDerivePresentSortsRelationshipsByPeerName(t *testing.T) {
	receiver, err := NewBeing(
		"michael",
		Nature{Identity: "builder", Purpose: "hold the line"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(receiver) error = %v", err)
	}
	sender, err := NewBeing(
		"zoe",
		Nature{Identity: "late", Purpose: "second"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(sender) error = %v", err)
	}
	early, err := NewBeing(
		"adam",
		Nature{Identity: "early", Purpose: "first"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(early) error = %v", err)
	}

	if err := receiver.SeedPeer(sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer(sender) error = %v", err)
	}
	if err := receiver.SeedPeer(early.Name, early.Nature); err != nil {
		t.Fatalf("SeedPeer(early) error = %v", err)
	}

	impulse, _ := NewImpulse("skyra zoe hello | michael: greeting")
	if _, err := receiver.EmitToPeer(sender.Name, impulse); err != nil {
		t.Fatalf("EmitToPeer() error = %v", err)
	}

	present, err := receiver.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	want := "name: michael\nidentity: builder\npurpose: hold the line\n\nYou are in an exchange with: zoe\nthe identity of zoe is: late\nthe purpose of zoe is: second\n\nzoe: hello\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nearly - first\nlate - second"
	if present != want {
		t.Fatalf("Present = %q, want %q", present, want)
	}
}

func TestExchangeStackDerivePresentUsesExpressionOnlyForNonCognitiveReceiver(t *testing.T) {
	receiver, err := NewBeing(
		"sensor",
		Nature{Identity: "sensor", Purpose: "receive"},
		true,
		false,
	)
	if err != nil {
		t.Fatalf("NewBeing(receiver) error = %v", err)
	}
	sender, err := NewBeing(
		"michael",
		Nature{Identity: "builder", Purpose: "hold the line"},
		true,
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(sender) error = %v", err)
	}
	if err := receiver.SeedPeer(sender.Name, sender.Nature); err != nil {
		t.Fatalf("SeedPeer() error = %v", err)
	}

	impulse, _ := NewImpulse("skyra michael heat spike -custom | sensor: spike detected")
	if _, err := receiver.EmitToPeer(sender.Name, impulse); err != nil {
		t.Fatalf("EmitToPeer() error = %v", err)
	}

	present, err := receiver.DerivePresent(sender)
	if err != nil {
		t.Fatalf("DerivePresent() error = %v", err)
	}

	want := "name: sensor\nidentity: sensor\npurpose: receive\n\nYou are in an exchange with: michael\nthe identity of michael is: builder\nthe purpose of michael is: hold the line\n\nmichael: heat spike\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nbuilder - hold the line"
	if present != want {
		t.Fatalf("Present = %q, want %q", present, want)
	}
}

func TestSignalCarriesRawImpulseUntilKernelParsing(t *testing.T) {
	impulse, err := NewImpulse("skyra skyra relate to me -close | michael: closing")
	if err != nil {
		t.Fatalf("NewImpulse() error = %v", err)
	}

	signal := Signal{
		ID:         "sig-1",
		Origin:     "michael",
		TraceToken: "trace-1",
		Impulse:    impulse.Raw(),
	}

	parsed, err := ParseImpulse(signal.Impulse)
	if err != nil {
		t.Fatalf("ParseImpulse(signal.Impulse) error = %v", err)
	}

	if parsed.TargetName != "skyra" {
		t.Fatalf("TargetName = %q, want %q", parsed.TargetName, "skyra")
	}
	if !parsed.IsClose() {
		t.Fatalf("IsClose() = false, want true")
	}
	if parsed.Source != "michael" {
		t.Fatalf("Source = %q, want %q", parsed.Source, "michael")
	}
	if parsed.Reason != "closing" {
		t.Fatalf("Reason = %q, want %q", parsed.Reason, "closing")
	}
}

func mustParseImpulse(t *testing.T, impulse Impulse) ParsedImpulse {
	t.Helper()

	parsed, err := impulse.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	return parsed
}
