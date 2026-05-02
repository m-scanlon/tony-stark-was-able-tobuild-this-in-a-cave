package reality

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type Remember struct {
	id string
}

func (rm *Remember) ID() string { return rm.id }

func (rm *Remember) Create(r *Relation) Reality {
	return &Remember{id: "remember"}
}

func (rm *Remember) Realize(r *Relation) string {
	if r.Log != nil {
		r.Log("[remember]: writing memory")
	}

	content := strings.TrimSpace(r.Impulse)
	if content == "" {
		if r.Log != nil {
			r.Log("[remember]: empty content")
		}
		return "nothing to remember"
	}

	being, ok := r.Realities["being"]
	if !ok {
		if r.Log != nil {
			r.Log("[remember]: no being on relation")
		}
		return "no being context"
	}
	b, ok := being.(Being)
	if !ok {
		return "no being context"
	}

	dir := filepath.Join(b.Home, "memories")
	os.MkdirAll(dir, 0755)

	filename := fmt.Sprintf("%d.md", time.Now().UnixMilli())
	path := filepath.Join(dir, filename)

	if err := os.WriteFile(path, []byte(content), 0644); err != nil {
		if r.Log != nil {
			r.Log("[remember]: write error:", err)
		}
		return "failed to remember"
	}

	if r.Log != nil {
		r.Log("[remember]: saved to", path)
	}
	return "remembered"
}
