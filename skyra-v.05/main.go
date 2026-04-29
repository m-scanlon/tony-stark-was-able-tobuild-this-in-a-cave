package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v05/src/inference"
	"skyra-v05/src/reality"
	"skyra-v05/src/reality/being"
	"skyra-v05/src/reality/world"
)

func main() {
	loadEnv("../.env")

	llm := world.NewLLM()
	physics := world.NewPhysics()
	system := world.NewSystem(physics)

	if err := bootstrap(system, llm, physics); err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap:", err)
		os.Exit(1)
	}

	llm.WireCall("openrouter", inference.Call)

	fmt.Println("skyra v.05")

	mac := world.NewMacOS()
	mac.Run(system)
}

func bootstrap(system *world.System, llm *world.LLM, physics *world.Physics) error {
	data, err := os.ReadFile("genome.skyra")
	if err != nil {
		return err
	}
	for _, line := range strings.Split(string(data), "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		tokens := strings.Fields(line)
		if len(tokens) < 2 {
			continue
		}
		operator := tokens[0]
		impulse := strings.Join(tokens[1:], " ")

		switch operator {
		case "physics":
			physics.Create(reality.Relation{Impulse: impulse})

		case "llm":
			llm.Create(reality.Relation{Impulse: impulse})

		case "grow":
			name, _ := being.Extract(impulse, "~name", "grow")
			deviceName, _ := being.Extract(impulse, "~device", "grow")

			pathos := being.Being{}.Create(reality.Relation{
				ID:      name,
				Impulse: impulse,
			}).(being.Being)

			var device world.Device
			switch deviceName {
			case "cli":
				device = world.NewCLIDevice(name)
			default:
				device = llm.Device(deviceName)
			}

			if device == nil {
				return fmt.Errorf("unknown device %q for being %q", deviceName, name)
			}

			bw := world.NewBeingWorld(pathos, device)
			system.Realities[name] = bw
		}
	}
	return nil
}

func loadEnv(path string) {
	f, err := os.Open(path)
	if err != nil {
		return
	}
	defer f.Close()
	sc := bufio.NewScanner(f)
	for sc.Scan() {
		line := strings.TrimSpace(sc.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		k, v, ok := strings.Cut(line, "=")
		if !ok {
			continue
		}
		if os.Getenv(k) == "" {
			os.Setenv(k, v)
		}
	}
}
