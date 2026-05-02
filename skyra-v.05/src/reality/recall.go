package reality

import (
	"os"
	"path/filepath"
	"strings"
)

type Recall struct {
	id string
}

func (rc *Recall) ID() string { return rc.id }

func (rc *Recall) Create(r *Relation) Reality {
	return &Recall{id: "recall"}
}

func (rc *Recall) Realize(r *Relation) string {
	if r.Log != nil {
		r.Log("[recall]: searching memories")
	}

	being, ok := r.Realities["being"]
	if !ok {
		if r.Log != nil {
			r.Log("[recall]: no being on relation")
		}
		return "no memories"
	}
	b, ok := being.(Being)
	if !ok {
		return "no memories"
	}

	dir := filepath.Join(b.Home, "memories")
	entries, err := os.ReadDir(dir)
	if err != nil {
		if r.Log != nil {
			r.Log("[recall]: no memories directory")
		}
		return "no memories yet"
	}

	if len(entries) == 0 {
		return "no memories yet"
	}

	query := strings.ToLower(strings.TrimSpace(r.Impulse))

	var matches []string
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}
		data, err := os.ReadFile(filepath.Join(dir, entry.Name()))
		if err != nil {
			continue
		}
		content := string(data)
		if query == "" || strings.Contains(strings.ToLower(content), query) {
			matches = append(matches, content)
		}
	}

	if len(matches) == 0 {
		if r.Log != nil {
			r.Log("[recall]: no matches for", query)
		}
		return "no relevant memories"
	}

	var sb strings.Builder
	sb.WriteString("memories:\n")
	for _, m := range matches {
		sb.WriteString("- " + m + "\n")
	}

	if r.Log != nil {
		r.Log("[recall]: found", len(matches), "memories")
	}
	return sb.String()
}
