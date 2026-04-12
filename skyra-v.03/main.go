package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v03/src/inference"
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
	runner := inference.New(inference.Config{
		APIKey: os.Getenv("GEMINI_API_KEY"),
		Model:  os.Getenv("GEMINI_MODEL"),
	})

	if err := seedGrow(w); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap: %v\n", err)
		os.Exit(1)
	}

	if err := bootstrap(w); err != nil {
		fmt.Fprintf(os.Stderr, "bootstrap: %v\n", err)
		os.Exit(1)
	}

	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}

		dispatch(m, w, runner, metaxu.Signal{
			Origin:  "michael",
			Impulse: fmt.Sprintf("skyra sensory %s | michael: experience", line),
		})
	}
}

func dispatch(m *metaxu.Metaxu, w *world.World, runner *inference.Runner, initial metaxu.Signal) {
	signal := initial
	for {
		result := m.AcceptSignal(signal)

		if result.Status == metaxu.RouteStatusDropped {
			fmt.Fprintf(os.Stderr, "dropped: %s\n", result.DropReason)
			return
		}

		if result.ReceiverPresent == "" {
			return
		}

		if result.ReceiverCognitive {
			next, err := runner.Run(result.ReceiverPresent, result.TargetName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "inference: %v\n", err)
				return
			}
			signal = next
			continue
		}

		// non-cognitive dispatch
		switch result.TargetName {
		case "grow":
			if _, err := w.Grow(result.ReceiverPresent); err != nil {
				fmt.Fprintf(os.Stderr, "grow: %v\n", err)
			}
		case "sensory":
			signal = metaxu.Signal{
				Origin:  "sensory",
				Impulse: fmt.Sprintf("skyra theory-of-mind %s | sensory: routing", result.ReceiverPresent),
			}
			continue
		case "motor":
			fmt.Println(result.ReceiverPresent)
		}
		return
	}
}

func bootstrap(w *world.World) error {
	data, err := os.ReadFile("genome.skyra")
	if err != nil {
		return err
	}

	var kept []string
	for _, line := range strings.Split(string(data), "\n") {
		if !strings.HasPrefix(strings.TrimSpace(line), "#") {
			kept = append(kept, line)
		}
	}

	chunks := strings.Split(strings.Join(kept, "\n"), "skyra")
	for _, chunk := range chunks {
		chunk = strings.TrimSpace(chunk)
		if chunk == "" {
			continue
		}
		if _, err := w.Grow("skyra " + chunk); err != nil {
			return fmt.Errorf("genome: %w", err)
		}
	}
	return nil
}

func seedGrow(w *world.World) error {
	id := identity.Identity{Value: "the instantiator"}
	p := purpose.Purpose{Value: "creates and registers beings from protocol expressions"}
	n := nature.Nature{Identity: id, Purpose: p}
	l := language.Language{Value: "skyra being ~name <name> ~identity <identity> ~purpose <purpose> ~language <expression> ~cognitive <true|false> ~relationships <a,b,c> | <source>: <reason>"}

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
