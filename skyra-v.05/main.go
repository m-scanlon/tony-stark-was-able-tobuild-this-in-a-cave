package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"strconv"

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

	devices := make(map[string]*reality.MacOS)
	components := make(map[string]reality.Reality)
	llmWires := make(map[string]string)

	exchange := &reality.Exchange{Exchanges: make(map[string]*reality.Conversation)}
	thread := &reality.NewThread{
		Beings:   make(map[string]reality.Reality),
		Access:   map[string]bool{"michael": true},
		Threads:  make(map[string]*reality.Thread),
		Exchange: exchange,
		Devices:  make(map[string]reality.Reality),
	}

	if err := bootstrap(thread, devices, components, llmWires); err != nil {
		fmt.Fprintln(os.Stderr, "bootstrap:", err)
		os.Exit(1)
	}

	for name, model := range llmWires {
		if comp, ok := components[name]; ok {
			if p, ok := comp.(*reality.Provider); ok {
				p.Model = model
				p.Call = inference.Call
			}
		}
	}

	mac := devices["macbook"]
	thread.Devices["macbook"] = mac

	var wsComp *reality.WS
	for _, comp := range components {
		if w, ok := comp.(*reality.WS); ok {
			wsComp = w
			break
		}
	}

	universe := &reality.Universe{Thread: thread}
	thread.OnResolve = func() {
		present := universe.Realize(&reality.Relation{Collecting: true})
		debug.Log("[universe]:", present)
		if wsComp != nil {
			wsComp.Broadcast(present)
		}
	}

	fmt.Println("skyra v.05")

	term := mac.Component("terminal")
	input := term.Realize(&reality.Relation{})
	rel, _ := reality.Impress("michael", input)

	universe.Realize(rel)
}

func bootstrap(thread *reality.NewThread, devices map[string]*reality.MacOS, components map[string]reality.Reality, llmWires map[string]string) error {
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
		verb := tokens[0]
		impulse := strings.Join(tokens[1:], " ")
		name, _ := reality.Extract(impulse, "~name", verb)

		switch verb {
		case "device":
			mac := (&reality.MacOS{}).Create(&reality.Relation{ID: name}).(*reality.MacOS)
			devices[name] = mac

		case "component":
			compType, _ := reality.Extract(impulse, "~type", "component")
			deviceName, _ := reality.Extract(impulse, "~device", "component")

			dev, ok := devices[deviceName]
			if !ok {
				return fmt.Errorf("component %q references unknown device %q", name, deviceName)
			}

			rel := &reality.Relation{ID: name, Impulse: impulse}
			rel.Realities = map[string]reality.Reality{"device": dev}

			switch compType {
			case "stdin":
				term := (&reality.Terminal{}).Create(rel).(*reality.Terminal)
				term.Device = dev
				dev.Components[name] = term
				components[name] = term

			case "llm":
				model, _ := reality.Extract(impulse, "~model", "component")
				p := &reality.Provider{Model: model, Device: dev}
				dev.Components[name] = p
				components[name] = p
				llmWires[name] = model

			case "websocket":
				portStr, _ := reality.Extract(impulse, "~port", "component")
				port, _ := strconv.Atoi(portStr)
				if port == 0 {
					port = 8080
				}
				ws := (&reality.WS{}).Create(rel).(*reality.WS)
				ws.Device = dev
				ws.Start(port)
				dev.Components[name] = ws
				components[name] = ws
			}

		case "being":
			beingType, _ := reality.Extract(impulse, "~type", "being")
			devicesRaw, _ := reality.Extract(impulse, "~devices", "being")

			ctx := &reality.Relation{ID: name, Impulse: impulse}
			ctx.Realities = make(map[string]reality.Reality)

			for _, devName := range strings.Split(devicesRaw, ",") {
				devName = strings.TrimSpace(devName)
				if dev, ok := devices[devName]; ok {
					ctx.Realities[devName] = dev
					for compName, comp := range dev.Components {
						ctx.Realities[compName] = comp
					}
				}
			}

			switch beingType {
			case "llm":
				self := (&reality.Self{}).Create(ctx)
				thread.Beings[name] = self
				thread.Access[name] = false

			case "user":
				user := (&reality.User{}).Create(ctx)
				thread.Beings[name] = user
				thread.Access[name] = true
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
