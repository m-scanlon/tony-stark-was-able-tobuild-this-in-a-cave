package metaxu

import (
	"testing"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/identity"
	"skyra-v03/src/primitives/language"
	"skyra-v03/src/primitives/nature"
	"skyra-v03/src/primitives/purpose"
	"skyra-v03/src/world"
)

func TestAcceptSignalRoutesToReceiverAndStartsNewExchange(t *testing.T) {
	w, origin, target := seededWorld(t)
	m := New(w)

	result := m.AcceptSignal(Signal{
		ID:         "sig-1",
		Origin:     origin.Name,
		TraceToken: "trace-1",
		Impulse:    "skyra skyra hello there | skyra: greeting",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if !result.NewExchange {
		t.Fatalf("NewExchange = false, want true")
	}
	if result.WrittenBeingName != target.Name {
		t.Fatalf("WrittenBeingName = %q, want %q", result.WrittenBeingName, target.Name)
	}
	if result.WrittenPeerName != origin.Name {
		t.Fatalf("WrittenPeerName = %q, want %q", result.WrittenPeerName, origin.Name)
	}
	want := "name: skyra\nidentity: system\npurpose: relate\n\nYou are in an exchange with: michael\nthe identity of michael is: builder\nthe purpose of michael is: hold the line\n\nmichael: hello there\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nbuilder - hold the line"
	if result.ReceiverPresent != want {
		t.Fatalf("ReceiverPresent = %q, want %q", result.ReceiverPresent, want)
	}
}

func TestAcceptSignalAppendsToExistingOpenExchange(t *testing.T) {
	w, origin, target := seededWorld(t)
	m := New(w)

	first, err := being.NewImpulse("skyra skyra first | skyra: opening")
	if err != nil {
		t.Fatalf("NewImpulse(first) error = %v", err)
	}
	if _, err := target.SendToPeer(origin.Name, deliveredFrom(t, origin.Name, first)); err != nil {
		t.Fatalf("SendToPeer(first) error = %v", err)
	}

	result := m.AcceptSignal(Signal{
		ID:         "sig-2",
		Origin:     origin.Name,
		TraceToken: "trace-2",
		Impulse:    "skyra skyra second | skyra: follow-up",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if result.NewExchange {
		t.Fatalf("NewExchange = true, want false")
	}
	want := "name: skyra\nidentity: system\npurpose: relate\n\nYou are in an exchange with: michael\nthe identity of michael is: builder\nthe purpose of michael is: hold the line\n\nmichael: first\n\nmichael: second\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nbuilder - hold the line"
	if result.ReceiverPresent != want {
		t.Fatalf("ReceiverPresent = %q, want %q", result.ReceiverPresent, want)
	}
}

func TestAcceptSignalStartsNewExchangeAfterClosedTop(t *testing.T) {
	w, origin, target := seededWorld(t)
	m := New(w)

	first, err := being.NewImpulse("skyra skyra first | skyra: opening")
	if err != nil {
		t.Fatalf("NewImpulse(first) error = %v", err)
	}
	closeImpulse, err := being.NewImpulse("skyra skyra -close | skyra: closing")
	if err != nil {
		t.Fatalf("NewImpulse(close) error = %v", err)
	}
	if _, err := target.SendToPeer(origin.Name, deliveredFrom(t, origin.Name, first)); err != nil {
		t.Fatalf("SendToPeer(first) error = %v", err)
	}
	if _, err := target.SendToPeer(origin.Name, deliveredFrom(t, origin.Name, closeImpulse)); err != nil {
		t.Fatalf("SendToPeer(close) error = %v", err)
	}

	result := m.AcceptSignal(Signal{
		ID:         "sig-3",
		Origin:     origin.Name,
		TraceToken: "trace-3",
		Impulse:    "skyra skyra after close | skyra: resuming",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if !result.NewExchange {
		t.Fatalf("NewExchange = false, want true")
	}

	targetPeer, ok := target.PeerByName(origin.Name)
	if !ok {
		t.Fatalf("PeerByName() ok = false, want true")
	}
	stack, ok := targetPeer.(*world.ExchangeStack)
	if !ok {
		t.Fatalf("target peer type = %T, want *world.ExchangeStack", targetPeer)
	}
	if len(stack.Exchanges()) != 2 {
		t.Fatalf("len(targetPeer.Exchanges()) = %d, want 2", len(stack.Exchanges()))
	}
	if len(stack.Exchanges()[0]) != 2 {
		t.Fatalf("len(targetPeer.Exchanges()[0]) = %d, want 2", len(stack.Exchanges()[0]))
	}
	if len(stack.Exchanges()[1]) != 1 {
		t.Fatalf("len(targetPeer.Exchanges()[1]) = %d, want 1", len(stack.Exchanges()[1]))
	}
	want := "name: skyra\nidentity: system\npurpose: relate\n\nYou are in an exchange with: michael\nthe identity of michael is: builder\nthe purpose of michael is: hold the line\n\nmichael: after close\n\nrelationships:\nCall any of your relationships using this syntax-\nskyra <being> <expression> | <source>: <reason> ~<emotional_signals>\n<being> must be one of your relationships listed below\n<source> is the being you are currently in exchange with\n<reason> is why you are firing this expression\nRespond with the protocol string only — no explanation, no markdown, no extra text\n________________\nbuilder - hold the line"
	if result.ReceiverPresent != want {
		t.Fatalf("ReceiverPresent = %q, want %q", result.ReceiverPresent, want)
	}
}

func TestAcceptSignalCloseWritesOnlyToOriginStack(t *testing.T) {
	w, origin, target := seededWorld(t)
	m := New(w)

	open, err := being.NewImpulse("skyra skyra open | skyra: opening")
	if err != nil {
		t.Fatalf("NewImpulse(open) error = %v", err)
	}
	if _, err := origin.EmitToPeer(target.Name, open); err != nil {
		t.Fatalf("EmitToPeer(open) error = %v", err)
	}

	result := m.AcceptSignal(Signal{
		ID:         "sig-4",
		Origin:     origin.Name,
		TraceToken: "trace-4",
		Impulse:    "skyra skyra -close | skyra: closing",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if result.WrittenBeingName != origin.Name {
		t.Fatalf("WrittenBeingName = %q, want %q", result.WrittenBeingName, origin.Name)
	}
	if result.WrittenPeerName != target.Name {
		t.Fatalf("WrittenPeerName = %q, want %q", result.WrittenPeerName, target.Name)
	}

	originPeer, ok := origin.PeerByName(target.Name)
	if !ok {
		t.Fatalf("origin PeerByName() ok = false, want true")
	}
	targetPeer, ok := target.PeerByName(origin.Name)
	if !ok {
		t.Fatalf("target PeerByName() ok = false, want true")
	}
	originStack, ok := originPeer.(*world.ExchangeStack)
	if !ok {
		t.Fatalf("origin peer type = %T, want *world.ExchangeStack", originPeer)
	}
	targetStack, ok := targetPeer.(*world.ExchangeStack)
	if !ok {
		t.Fatalf("target peer type = %T, want *world.ExchangeStack", targetPeer)
	}
	if len(originStack.Exchanges()) != 1 {
		t.Fatalf("len(originStack.Exchanges()) = %d, want 1", len(originStack.Exchanges()))
	}
	if len(originStack.Exchanges()[0]) != 2 {
		t.Fatalf("len(originStack.Exchanges()[0]) = %d, want 2", len(originStack.Exchanges()[0]))
	}
	if len(targetStack.Exchanges()) != 0 {
		t.Fatalf("len(targetStack.Exchanges()) = %d, want 0", len(targetStack.Exchanges()))
	}
	if string(originStack.Exchanges()[0][1]) != "skyra skyra -close | skyra: closing" {
		t.Fatalf("close stored raw impulse = %q, want raw value", string(originStack.Exchanges()[0][1]))
	}
}

func TestAcceptSignalDropsWhenTargetCannotBeResolved(t *testing.T) {
	w, origin, _ := seededWorld(t)
	m := New(w)

	result := m.AcceptSignal(Signal{
		ID:         "sig-5",
		Origin:     origin.Name,
		TraceToken: "trace-5",
		Impulse:    "skyra unknown hello | skyra: testing",
	})

	if result.Status != RouteStatusDropped {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusDropped)
	}
	if result.DropReason == "" {
		t.Fatalf("DropReason = empty, want value")
	}
}

func TestAcceptSignalRoutesToExternalDispatchAndDerivesExpressionOnlyPresent(t *testing.T) {
	w := world.New()

	origin, err := being.NewBeing(
		"michael",
		nature.Nature{Identity: identity.Identity{Value: "builder"}, Purpose: purpose.Purpose{Value: "hold the line"}},
		language.Language{Value: "skyra being"},
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(origin) error = %v", err)
	}
	sensor, err := being.NewBeing(
		"sensor",
		nature.Nature{Identity: identity.Identity{Value: "sensor"}, Purpose: purpose.Purpose{Value: "receive"}},
		language.Language{Value: "skyra being"},
		false,
	)
	if err != nil {
		t.Fatalf("NewBeing(sensor) error = %v", err)
	}
	if err := w.Register(origin); err != nil {
		t.Fatalf("Register(origin) error = %v", err)
	}
	if err := w.Register(sensor); err != nil {
		t.Fatalf("Register(sensor) error = %v", err)
	}
	if err := w.Relate(origin.Name, sensor.Name); err != nil {
		t.Fatalf("Relate() error = %v", err)
	}

	// Replace sensor's peer channel for origin with an ExternalDispatch
	dispatch, err := world.NewExternalDispatch(origin.Name, origin.Nature)
	if err != nil {
		t.Fatalf("NewExternalDispatch() error = %v", err)
	}
	if err := sensor.AttachPeer(dispatch); err != nil {
		t.Fatalf("AttachPeer() error = %v", err)
	}

	m := New(w)
	result := m.AcceptSignal(Signal{
		ID:         "sig-6",
		Origin:     origin.Name,
		TraceToken: "trace-6",
		Impulse:    "skyra sensor heat spike | sensor: reporting",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if result.ReceiverPresent != "heat spike" {
		t.Fatalf("ReceiverPresent = %q, want %q", result.ReceiverPresent, "heat spike")
	}
}

func seededWorld(t *testing.T) (*world.World, *being.Being, *being.Being) {
	t.Helper()

	w := world.New()

	origin, err := being.NewBeing(
		"michael",
		nature.Nature{Identity: identity.Identity{Value: "builder"}, Purpose: purpose.Purpose{Value: "hold the line"}},
		language.Language{Value: "skyra being"},
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(origin) error = %v", err)
	}

	target, err := being.NewBeing(
		"skyra",
		nature.Nature{Identity: identity.Identity{Value: "system"}, Purpose: purpose.Purpose{Value: "relate"}},
		language.Language{Value: "skyra being"},
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(target) error = %v", err)
	}

	if err := w.Register(origin); err != nil {
		t.Fatalf("Register(origin) error = %v", err)
	}
	if err := w.Register(target); err != nil {
		t.Fatalf("Register(target) error = %v", err)
	}
	if err := w.Relate(origin.Name, target.Name); err != nil {
		t.Fatalf("Relate() error = %v", err)
	}

	return w, origin, target
}

func deliveredFrom(t *testing.T, originName string, impulse being.Impulse) being.DeliveredImpulse {
	t.Helper()

	parsed, err := impulse.Parse()
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	return being.DeliveredImpulse{
		OriginName: originName,
		Raw:        impulse,
		Parsed:     parsed,
	}
}
