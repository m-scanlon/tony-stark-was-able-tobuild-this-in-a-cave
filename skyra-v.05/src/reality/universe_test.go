package reality

import (
	"testing"
)

func TestUniverseRealize(t *testing.T) {
	exchange := &Exchange{Exchanges: make(map[string]*Conversation)}
	thread := &NewThread{
		Beings:   make(map[string]Reality),
		Access:   map[string]bool{"michael": true},
		Threads:  make(map[string]*Thread),
		Exchange: exchange,
		Devices:  make(map[string]Reality),
	}

	skyra := &Self{id: "skyra", Realities: make(map[string]Reality)}
	skyra.Realities["being"] = Being{
		id:            "skyra",
		name:          "skyra",
		Identity:      "I hold the world together.",
		Purpose:       "I think, respond, and relate.",
		Relationships: []string{"michael", "louise"},
		Device:        "openrouter",
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
		Device:        "macos",
	}
	michael.Realities["device"] = &MacOS{id: "macos"}
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
