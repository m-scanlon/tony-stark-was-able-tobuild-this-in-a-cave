package reality

import (
	"testing"
)

func TestUniverseRealize(t *testing.T) {
	exchange := &Exchange{Exchanges: make(map[string]*Conversation)}

	mac := &MacOS{}
	mac = mac.Create(&Relation{ID: "macbook"}).(*MacOS)

	term := &Terminal{}
	term = term.Create(&Relation{}).(*Terminal)
	term.Device = mac
	mac.Components["terminal"] = term

	provider := &Provider{id: "openrouter", Model: "anthropic/claude-sonnet-4-5", Device: mac}
	mac.Components["openrouter"] = provider

	thread := &NewThread{
		Beings:   make(map[string]Reality),
		Access:   map[string]bool{"michael": true},
		Threads:  make(map[string]*Thread),
		Exchange: exchange,
		Devices:  map[string]Reality{"macbook": mac},
	}

	skyra := &Self{id: "skyra", Realities: make(map[string]Reality)}
	skyra.Realities["being"] = Being{
		id:            "skyra",
		name:          "skyra",
		Identity:      "I hold the world together.",
		Purpose:       "I think, respond, and relate.",
		Relationships: []string{"michael", "louise"},
		Device:        "macbook",
	}
	skyra.Realities["think"] = &Think{
		Operators: map[string]Reality{
			"recall":   &Recall{},
			"remember": &Remember{},
			"skill":    &Skill{},
		},
	}
	skyra.Realities["act"] = &Act{
		Operators: map[string]Reality{
			"plan": &Plan{},
		},
	}
	thread.Beings["skyra"] = skyra

	michael := &User{id: "michael", Realities: make(map[string]Reality)}
	michael.Realities["being"] = Being{
		id:            "michael",
		name:          "michael",
		Identity:      "I build Skyra.",
		Purpose:       "I decide what matters.",
		Relationships: []string{"skyra"},
		Device:        "macbook",
	}
	michael.Realities["device"] = mac
	thread.Beings["michael"] = michael

	th := thread.newThread("michael")
	th.Spread("michael", "skyra")

	econ := NewEconomics()
	econ.Set("inference_calls", 42)

	universe := &Universe{id: "universe", Thread: thread, Econ: econ}
	present := universe.Realize(&Relation{Collecting: true})

	checks := []string{
		`"name": "skyra"`,
		`"type": "llm"`,
		`"identity": "I hold the world together."`,
		`"name": "michael"`,
		`"type": "user"`,
		`"from": "michael"`,
		`"to": "skyra"`,
		`"inference_calls": 42`,
		`"type": "Universe"`,
		`"type": "NewThread"`,
		`"id": "macbook"`,
		`"id": "terminal"`,
		`"id": "openrouter"`,
		`"type": "Self"`,
		`"type": "Think"`,
		`"type": "Recall"`,
		`"type": "Act"`,
		`"type": "Plan"`,
		`"type": "User"`,
		`"type": "Economics"`,
	}

	for _, check := range checks {
		if !contains(present, check) {
			t.Errorf("missing in output: %s", check)
		}
	}

	t.Log(present)
}

func TestCLIRouting(t *testing.T) {
	cli := &CLI{id: "linter", Realities: make(map[string]Reality)}
	cli.Realities["being"] = Being{
		id:            "linter",
		name:          "linter",
		Entrypoints:   []string{"echo"},
		Relationships: []string{"skyra"},
	}

	r := &Relation{
		Origin:    "skyra",
		ID:        "linter",
		Impulse:   "check this",
		Parsers:   make(map[string]Parser),
		Realities: map[string]Reality{"linter": cli},
	}

	cli.Realize(r)

	if r.Origin != "linter" {
		t.Errorf("expected Origin to be 'linter', got %q", r.Origin)
	}
	if r.ID != "skyra" {
		t.Errorf("expected ID to be 'skyra' (route back to sender), got %q", r.ID)
	}
	if r.Impulse == "" {
		t.Errorf("expected Impulse to contain output, got empty")
	}
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && searchString(s, substr)
}

func searchString(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
