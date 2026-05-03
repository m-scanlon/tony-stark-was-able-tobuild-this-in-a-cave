package reality

import (
	"bytes"
	"os/exec"
	"skyra-v05/src/debug"
	"strings"
)

type CLI struct {
	id        string
	Realities map[string]Reality
}

func (c *CLI) ID() string { return c.id }

func (c *CLI) Create(r *Relation) Reality {
	cli := &CLI{
		id:        r.ID,
		Realities: make(map[string]Reality),
	}

	if r.Impulse != "" {
		being := Being{}.Create(r).(Being)
		cli.Realities["being"] = being
	}

	return cli
}

func (c *CLI) Realize(r *Relation) string {
	if r.Collecting {
		node := RealityNode{ID: c.id, Type: "CLI", Children: []RealityNode{}}
		snap := BeingSnapshot{
			Name: c.id, Type: "cli", Status: "idle",
			Peers:       []string{},
			Entrypoints: []string{},
			Memories:    MemorySnapshot{Items: []MemoryItem{}, Skills: []SkillItem{}},
		}

		if being, ok := c.Realities["being"].(Being); ok {
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
			node.Children = append(node.Children, RealityNode{ID: c.id + "-being", Type: "Being", Children: []RealityNode{}})
		}

		r.Export("being:"+c.id, snap)
		r.Export("node:"+c.id, node)
		return ""
	}

	debug.Log("[cli]: realizing", c.id)

	being, ok := c.Realities["being"].(Being)
	if !ok {
		debug.Log("[cli]: no being")
		return ""
	}

	if len(being.Entrypoints) == 0 {
		debug.Log("[cli]: no entrypoints")
		return ""
	}

	entrypoint := being.Entrypoints[0]
	debug.Log("[cli]: calling", entrypoint)

	cmd := exec.Command(entrypoint, r.Impulse)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		debug.Log("[cli]: error:", err.Error(), stderr.String())
		return ""
	}

	result := strings.TrimSpace(stdout.String())
	debug.Log("[cli]: response length:", len(result))

	r.Impulse = result
	r.ID = r.Origin
	r.Origin = c.id
	return result
}
