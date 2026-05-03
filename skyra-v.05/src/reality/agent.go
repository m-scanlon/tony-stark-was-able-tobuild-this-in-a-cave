package reality

import (
	"bytes"
	"encoding/json"
	"os/exec"
	"skyra-v05/src/debug"
	"strings"
)

type Agent struct {
	id        string
	Realities map[string]Reality
	SessionID string
}

func (a *Agent) ID() string { return a.id }

func (a *Agent) Create(r *Relation) Reality {
	agent := &Agent{
		id:        r.ID,
		Realities: make(map[string]Reality),
	}

	if r.Impulse != "" {
		being := Being{}.Create(r).(Being)
		agent.Realities["being"] = being
	}

	return agent
}

func (a *Agent) Realize(r *Relation) string {
	if r.Collecting {
		node := RealityNode{ID: a.id, Type: "Agent", Children: []RealityNode{}}
		snap := BeingSnapshot{
			Name: a.id, Type: "agent", Status: "idle",
			Peers:       []string{},
			Entrypoints: []string{},
			Memories:    MemorySnapshot{Items: []MemoryItem{}, Skills: []SkillItem{}},
		}

		if being, ok := a.Realities["being"].(Being); ok {
			snap.Identity = being.Identity
			snap.Purpose = being.Purpose
			if being.Relationships != nil {
				snap.Peers = being.Relationships
			}
			if being.Entrypoints != nil {
				snap.Entrypoints = being.Entrypoints
			}
			snap.Device = being.Device
			snap.Memories = snapshotMemories(being.Home)
			node.Children = append(node.Children, RealityNode{ID: a.id + "-being", Type: "Being", Children: []RealityNode{}})
		}

		r.Export("being:"+a.id, snap)
		r.Export("node:"+a.id, node)
		return ""
	}

	debug.Log("[agent]: realizing", a.id)

	being, ok := a.Realities["being"].(Being)
	if !ok {
		debug.Log("[agent]: no being")
		return ""
	}

	if len(being.Entrypoints) == 0 {
		debug.Log("[agent]: no entrypoints")
		return ""
	}

	entrypoint := being.Entrypoints[0]

	args := []string{"-p", r.Impulse, "--output-format", "json"}
	if a.SessionID != "" {
		args = append(args, "--resume", a.SessionID)
		debug.Log("[agent]: resuming session", a.SessionID)
	} else {
		debug.Log("[agent]: starting new session")
	}

	debug.Log("[agent]: calling", entrypoint)
	cmd := exec.Command(entrypoint, args...)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		debug.Log("[agent]: error:", err.Error(), stderr.String())
		return ""
	}

	raw := stdout.String()
	debug.Log("[agent]: raw response length:", len(raw))

	var parsed struct {
		SessionID string `json:"session_id"`
		Result    string `json:"result"`
	}
	if err := json.Unmarshal([]byte(raw), &parsed); err != nil {
		debug.Log("[agent]: json parse error, using raw output")
		result := strings.TrimSpace(raw)
		r.Impulse = result
		r.ID = r.Origin
		r.Origin = a.id
		return result
	}

	if parsed.SessionID != "" {
		a.SessionID = parsed.SessionID
		debug.Log("[agent]: captured session id:", a.SessionID)
	}

	result := strings.TrimSpace(parsed.Result)
	debug.Log("[agent]: response length:", len(result))

	r.Impulse = result
	r.ID = r.Origin
	r.Origin = a.id
	return result
}
