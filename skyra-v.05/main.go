package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"skyra-v05/src/debug"
	"skyra-v05/src/inference"
	"skyra-v05/src/reality"
)

func main() {
	loadEnv("../.env")

	if err := debug.Init("logs"); err != nil {
		fmt.Fprintln(os.Stderr, "debug:", err)
		os.Exit(1)
	}
	defer debug.Close()

	if err := reality.InitHome(); err != nil {
		fmt.Fprintln(os.Stderr, "home:", err)
		os.Exit(1)
	}

	llm := reality.NewLLM()
	exchange := &reality.Exchange{Exchanges: make(map[string]*reality.Conversation)}
	thread := &reality.NewThread{
		Beings:   make(map[string]reality.Reality),
		Access:   map[string]bool{"michael": true},
		Threads:  make(map[string]*reality.Thread),
		Exchange: exchange,
		Devices:  make(map[string]reality.Reality),
	}
	mac := &reality.MacOS{}
	mac = mac.Create(&reality.Relation{}).(*reality.MacOS)

	if err := bootstrap(thread, llm, mac); err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap:", err)
		os.Exit(1)
	}

	llm.WireCall("openrouter", inference.Call)

	thread.Devices["macos"] = mac
	if p := llm.Provider("openrouter"); p != nil {
		thread.Devices["openrouter"] = p
	}

	fmt.Println("skyra v.05")

	input := mac.Realize(&reality.Relation{})
	rel, _ := reality.Impress("michael", input)

	thread.Realize(rel)
}

func bootstrap(thread *reality.NewThread, llm *reality.LLM, mac *reality.MacOS) error {
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
		case "llm":
			llm.Create(&reality.Relation{Impulse: impulse})

		case "grow":
			name, _ := reality.Extract(impulse, "~name", "grow")
			beingType, _ := reality.Extract(impulse, "~type", "grow")
			deviceName, _ := reality.Extract(impulse, "~device", "grow")

			being := reality.Being{}.Create(&reality.Relation{
				ID:      name,
				Impulse: impulse,
			}).(reality.Being)

			var device reality.Reality
			switch deviceName {
			case "macos":
				device = mac
			default:
				device = llm.Provider(deviceName)
				if device == nil {
					return fmt.Errorf("unknown device %q for being %q", deviceName, name)
				}
			}

			switch beingType {
			case "llm":
				self := &reality.Self{}
				self = self.Create(&reality.Relation{ID: name}).(*reality.Self)
				self.Realities["being"] = being

				think := &reality.Think{
					Operators: map[string]reality.Reality{
						"recall":   &reality.Recall{},
						"remember": &reality.Remember{},
						"skill":    &reality.Skill{},
					},
					LLM: device,
				}

				act := &reality.Act{
					Operators: map[string]reality.Reality{
						"plan": &reality.Plan{},
					},
					LLM: device,
				}

				self.Realities["think"] = think
				self.Realities["act"] = act

				thread.Beings[name] = self

			case "user":
				user := &reality.User{}
				user = user.Create(&reality.Relation{ID: name}).(*reality.User)
				user.Realities["being"] = being
				user.Realities["device"] = device

				thread.Beings[name] = user
			}
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
