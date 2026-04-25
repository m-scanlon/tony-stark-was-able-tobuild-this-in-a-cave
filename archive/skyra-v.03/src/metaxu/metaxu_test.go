package metaxu

import (
	"fmt"
	"testing"

	being "skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/identity"
	"skyra-v03/src/primitives/language"
	"skyra-v03/src/primitives/nature"
	"skyra-v03/src/primitives/purpose"
	"skyra-v03/src/world"
)

func TestAcceptSignalRoutesAndWritesToExchange(t *testing.T) {
	w, origin, target := seededWorld(t)
	m := New(w)

	// open exchange on both sides before sending
	openExchangeOnPeer(t, origin, target.Name, "intent-1")
	openExchangeOnPeer(t, target, origin.Name, "intent-1")

	result := m.AcceptSignal(Signal{
		Origin:   origin.Name,
		ThreadID: "intent-1",
		Impulse:  "skyra skyra hello there | greeting",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if result.WrittenBeingName != target.Name {
		t.Fatalf("WrittenBeingName = %q, want %q", result.WrittenBeingName, target.Name)
	}
	if result.WrittenPeerName != origin.Name {
		t.Fatalf("WrittenPeerName = %q, want %q", result.WrittenPeerName, origin.Name)
	}
	if result.ReceiverPresent == "" {
		t.Fatalf("ReceiverPresent = empty, want value")
	}
}

func TestAcceptSignalDropsWhenTargetCannotBeResolved(t *testing.T) {
	w, origin, _ := seededWorld(t)
	m := New(w)

	result := m.AcceptSignal(Signal{
		Origin:   origin.Name,
		ThreadID: "intent-1",
		Impulse:  "skyra unknown hello | testing",
	})

	if result.Status != RouteStatusDropped {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusDropped)
	}
	if result.DropReason == "" {
		t.Fatalf("DropReason = empty, want value")
	}
}

func TestAcceptSignalDropsWhenNoThreadID(t *testing.T) {
	w, origin, target := seededWorld(t)
	m := New(w)

	openExchangeOnPeer(t, target, origin.Name, "intent-1")

	result := m.AcceptSignal(Signal{
		Origin:  origin.Name,
		Impulse: "skyra skyra hello | greeting",
	})

	if result.Status != RouteStatusDropped {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusDropped)
	}
}

func TestAcceptSignalDropsWhenCognitiveOriginHasNoOpenExchange(t *testing.T) {
	w := world.New()

	prefrontal, err := being.NewBeing(
		"prefrontal",
		nature.Nature{Identity: identity.Identity{Value: "system"}, Purpose: purpose.Purpose{Value: "relate"}},
		language.Language{Value: "skyra being"},
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(prefrontal) error = %v", err)
	}
	values, err := being.NewBeing(
		"values",
		nature.Nature{Identity: identity.Identity{Value: "judge"}, Purpose: purpose.Purpose{Value: "align"}},
		language.Language{Value: "skyra being"},
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(values) error = %v", err)
	}

	if err := w.Register(prefrontal); err != nil {
		t.Fatalf("Register(prefrontal) error = %v", err)
	}
	if err := w.Register(values); err != nil {
		t.Fatalf("Register(values) error = %v", err)
	}
	if _, err := w.Grow(fmt.Sprintf("skyra being ~name %s ~relationships %s", prefrontal.Name, values.Name)); err != nil {
		t.Fatalf("Grow(relate) error = %v", err)
	}

	m := New(w)
	result := m.AcceptSignal(Signal{
		Origin:   prefrontal.Name,
		ThreadID: "experience",
		Impulse:  "skyra values ask what michael should be asked | consulting values",
	})

	if result.Status != RouteStatusDropped {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusDropped)
	}
	want := "no open exchange with values for thread experience; open one with start-exchange first"
	if result.DropReason != want {
		t.Fatalf("DropReason = %q, want %q", result.DropReason, want)
	}
}

func TestAcceptSignalAllowsNonCognitiveOriginWithoutOpenExchange(t *testing.T) {
	w := world.New()

	michael, err := being.NewBeing(
		"michael",
		nature.Nature{Identity: identity.Identity{Value: "builder"}, Purpose: purpose.Purpose{Value: "hold the line"}},
		language.Language{Value: "skyra being"},
		false,
	)
	if err != nil {
		t.Fatalf("NewBeing(michael) error = %v", err)
	}
	prefrontal, err := being.NewBeing(
		"prefrontal",
		nature.Nature{Identity: identity.Identity{Value: "system"}, Purpose: purpose.Purpose{Value: "relate"}},
		language.Language{Value: "skyra being"},
		true,
	)
	if err != nil {
		t.Fatalf("NewBeing(prefrontal) error = %v", err)
	}

	if err := w.Register(michael); err != nil {
		t.Fatalf("Register(michael) error = %v", err)
	}
	if err := w.Register(prefrontal); err != nil {
		t.Fatalf("Register(prefrontal) error = %v", err)
	}
	if _, err := w.Grow(fmt.Sprintf("skyra being ~name %s ~relationships %s", prefrontal.Name, michael.Name)); err != nil {
		t.Fatalf("Grow(relate) error = %v", err)
	}

	m := New(w)
	result := m.AcceptSignal(Signal{
		Origin:   michael.Name,
		ThreadID: "experience",
		Impulse:  "skyra prefrontal hi | greeting",
	})

	if result.Status != RouteStatusRouted {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusRouted)
	}
	if result.ReceiverPresent == "" {
		t.Fatalf("ReceiverPresent = empty, want value")
	}
}

func TestAcceptSignalDropsWhenExchangeOperatorIsNotTargetedDirectly(t *testing.T) {
	w, origin, target := seededWorld(t)
	m := New(w)

	result := m.AcceptSignal(Signal{
		Origin:   origin.Name,
		ThreadID: "intent-1",
		Impulse:  fmt.Sprintf("skyra %s close-exchange ~with %s | testing", target.Name, origin.Name),
	})

	if result.Status != RouteStatusDropped {
		t.Fatalf("Status = %q, want %q", result.Status, RouteStatusDropped)
	}
	want := fmt.Sprintf("target close-exchange directly; do not wrap close-exchange inside a signal to %s", target.Name)
	if result.DropReason != want {
		t.Fatalf("DropReason = %q, want %q", result.DropReason, want)
	}
}

func TestAcceptSignalRoutesToExternalDispatch(t *testing.T) {
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
	if _, err := w.Grow(fmt.Sprintf("skyra being ~name %s ~relationships %s", origin.Name, sensor.Name)); err != nil {
		t.Fatalf("Grow(relate) error = %v", err)
	}
	openExchangeOnPeer(t, origin, sensor.Name, "intent-1")

	dispatch, err := world.NewExternalDispatch(origin.Name, origin.Nature)
	if err != nil {
		t.Fatalf("NewExternalDispatch() error = %v", err)
	}
	if err := sensor.AttachPeer(dispatch); err != nil {
		t.Fatalf("AttachPeer() error = %v", err)
	}

	m := New(w)
	result := m.AcceptSignal(Signal{
		Origin:   origin.Name,
		ThreadID: "intent-1",
		Impulse:  "skyra sensor heat spike | reporting",
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
	if _, err := w.Grow(fmt.Sprintf("skyra being ~name %s ~relationships %s", origin.Name, target.Name)); err != nil {
		t.Fatalf("Grow(relate) error = %v", err)
	}

	return w, origin, target
}

func openExchangeOnPeer(t *testing.T, b *being.Being, peerName string, threadID string) {
	t.Helper()
	peer, ok := b.PeerByName(peerName)
	if !ok {
		t.Fatalf("openExchangeOnPeer: PeerByName(%q) ok = false", peerName)
	}
	em, ok := peer.(interface {
		OpenExchange(string, string, string, []being.ExchangeEntry) error
	})
	if !ok {
		t.Fatalf("openExchangeOnPeer: peer %q does not implement OpenExchange", peerName)
	}
	if err := em.OpenExchange(threadID, "test intent", "test reason", nil); err != nil {
		t.Fatalf("openExchangeOnPeer: OpenExchange(%q) error = %v", threadID, err)
	}
}
