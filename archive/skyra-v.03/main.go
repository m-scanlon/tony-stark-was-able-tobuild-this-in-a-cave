package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v03/src/inference"
	"skyra-v03/src/keychain"
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
	apiKey := keychain.Get("OpenRouter_API_KEY")
	if apiKey == "" {
		apiKey = os.Getenv("OpenRouter_API_KEY")
	}
	runner := inference.New(inference.Config{
		BaseURL: os.Getenv("OLLAMA_BASE_URL"),
		Model:   os.Getenv("OLLAMA_MODEL"),
		APIKey:  apiKey,
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
			Origin:   "michael",
			ThreadID: "experience",
			Impulse:  fmt.Sprintf("skyra skyra %s | experience", line),
		})
	}
}

func dispatch(m *metaxu.Metaxu, w *world.World, runner *inference.Runner, initial metaxu.Signal) {
	pending := []metaxu.Signal{initial}
	var lastCognitiveName string
	var lastPresent string
	var lastSignal string

	for len(pending) > 0 {
		signal := pending[0]
		pending = pending[1:]

		result := m.AcceptSignal(signal)

		if result.Status == metaxu.RouteStatusDropped {
			fmt.Fprintf(os.Stderr, "dropped: %s\n", result.DropReason)
			if result.OriginName != "" && result.OriginName == lastCognitiveName && lastSignal != "" && lastPresent != "" {
				correction := lastPresent + "\n\nIMPORTANT:\nyou wrote this last time:\n" + lastSignal + "\nyour last signal was dropped: " + result.DropReason + "\nrespond again with one corrected protocol string only"
				lastSignal = ""
				lastCognitiveName = ""
				lastPresent = ""
				signals, err := runner.Run(correction, result.OriginName)
				if err != nil {
					fmt.Fprintf(os.Stderr, "inference: %v\n", err)
					return
				}
				for i := range signals {
					signals[i].ThreadID = signal.ThreadID
				}
				pending = append(signals, pending...)
				continue
			}
			return
		}

		if result.ReceiverPresent == "" {
			return
		}

		if result.ReceiverCognitive {
			lastCognitiveName = result.TargetName
			lastPresent = result.ReceiverPresent
			lastSignal = signal.Impulse
			signals, err := runner.Run(result.ReceiverPresent, result.TargetName)
			if err != nil {
				fmt.Fprintf(os.Stderr, "inference: %v\n", err)
				return
			}
			for i := range signals {
				if signals[i].ThreadID == "" {
					signals[i].ThreadID = result.ThreadID
				}
			}
			pending = append(signals, pending...)
			continue
		}

		// non-cognitive dispatch
		switch result.TargetName {
		case "grow":
			if _, err := w.Grow(result.ReceiverPresent); err != nil {
				fmt.Fprintf(os.Stderr, "grow: %v\n", err)
			}
		case "close-exchange":
			res, err := w.ParseCloseExchange(result.ReceiverPresent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "close-exchange: %v\n", err)
				return
			}
			if err := w.CloseExchange(result.OriginName, res.PeerName, result.ThreadID); err != nil {
				fmt.Fprintf(os.Stderr, "close-exchange: %v\n", err)
			}
		case "start-exchange":
			res, err := w.StartExchange(result.ReceiverPresent)
			if err != nil {
				fmt.Fprintf(os.Stderr, "start-exchange: %v\n", err)
				return
			}
			contextEntries := w.ResolveExpressionRef(result.OriginName, result.ThreadID, res.ExpressionRef)
			if err := w.OpenExchange(result.OriginName, res.PeerName, res.ThreadID, res.About, res.Because, contextEntries); err != nil {
				fmt.Fprintf(os.Stderr, "start-exchange: %v\n", err)
				return
			}
			openSignal := metaxu.Signal{
				Origin:         result.OriginName,
				ThreadID:       res.ThreadID,
				About:          res.About,
				Because:        res.Because,
				ContextEntries: contextEntries,
				Impulse:        fmt.Sprintf("skyra %s %s | opening exchange", res.PeerName, res.Said),
			}
			pending = append([]metaxu.Signal{openSignal}, pending...)
		case "motor":
			fmt.Println(result.ReceiverPresent)
		}
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
