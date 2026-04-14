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
		BaseURL: os.Getenv("OLLAMA_BASE_URL"),
		Model:   os.Getenv("OLLAMA_MODEL"),
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
			Impulse: fmt.Sprintf("skyra sensory %s | experience", line),
		})
	}
}

func dispatch(m *metaxu.Metaxu, w *world.World, runner *inference.Runner, initial metaxu.Signal) {
	signal := initial
	var lastCognitiveName string
	var lastSignal string

	for {
		result := m.AcceptSignal(signal)

		if result.Status == metaxu.RouteStatusDropped {
			fmt.Fprintf(os.Stderr, "dropped: %s\n", result.DropReason)
			if result.OriginName != "" && result.OriginName == lastCognitiveName && lastSignal != "" {
				correction := "you wrote this last time:\n" + lastSignal + "\nyour last signal was dropped: " + result.DropReason + " — try again"
				lastSignal = ""
				lastCognitiveName = ""
				next, err := runner.Run(correction, result.OriginName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "inference: %v\n", err)
					return
				}
				signal = next
				continue
			}
			return
		}

		if result.ReceiverPresent == "" {
			return
		}

		if result.ReceiverCognitive {
			lastCognitiveName = result.TargetName
			lastSignal = signal.Impulse
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
				Impulse: fmt.Sprintf("skyra thalamus %s | routing", result.ReceiverPresent),
			}
			continue
		case "thalamus":
			signal = metaxu.Signal{
				Origin:  "thalamus",
				Impulse: fmt.Sprintf("skyra prefrontal %s | relaying", result.ReceiverPresent),
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

	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if _, err := w.Grow(line); err != nil {
			return fmt.Errorf("genome: %w", err)
		}
	}
	return nil
}

func seedGrow(w *world.World) error {
	id := identity.Identity{Value: "the instantiator"}
	p := purpose.Purpose{Value: "creates and registers beings from protocol expressions"}
	n := nature.Nature{Identity: id, Purpose: p}
	l := language.Language{Value: "skyra being ~name <name> ~identity <identity> ~purpose <purpose> ~language <expression> ~cognitive <true|false> ~relationships <a,b,c> | <reason>"}

	b, err := being.NewBeing("grow", n, l, false)
	if err != nil {
		return err
	}
	return w.Register(b)
}

