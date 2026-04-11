package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v03/src/metaxu"
	"skyra-v03/src/primitives/being"
	"skyra-v03/src/primitives/identity"
	"skyra-v03/src/primitives/language"
	"skyra-v03/src/primitives/nature"
	"skyra-v03/src/primitives/purpose"
	"skyra-v03/src/world"
)

func main() {
	w := world.New()
	m := metaxu.New(w)

	if err := seedGrow(w); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		result := m.AcceptSignal(metaxu.Signal{
			Origin:  originFrom(line),
			Impulse: line,
		})

		if result.Status == metaxu.RouteStatusDropped {
			fmt.Fprintf(os.Stderr, "dropped: %s\n", result.DropReason)
			continue
		}

		if result.ReceiverPresent != "" {
			fmt.Println(result.ReceiverPresent)
		}
	}
}

func seedGrow(w *world.World) error {
	id := identity.Identity{Value: "the instantiator"}
	p := purpose.Purpose{Value: "creates and registers beings from protocol expressions"}
	n := nature.Nature{Identity: id, Purpose: p}
	l := language.Language{Value: "skyra being ~name <name> ~identity <identity> ~purpose <purpose> ~language <expression> ~cognitive <true|false> | <source>: <reason>"}

	b, err := being.NewBeing("grow", n, l, false)
	if err != nil {
		return err
	}
	return w.Register(b)
}

// originFrom extracts the source name from the right zone of a protocol string.
// skyra <being> <expression> | <source>: <reason>
func originFrom(impulse string) string {
	parts := strings.SplitN(impulse, "|", 2)
	if len(parts) != 2 {
		return ""
	}
	right := strings.TrimSpace(parts[1])
	colonIdx := strings.Index(right, ":")
	if colonIdx < 0 {
		return ""
	}
	return strings.TrimSpace(right[:colonIdx])
}
