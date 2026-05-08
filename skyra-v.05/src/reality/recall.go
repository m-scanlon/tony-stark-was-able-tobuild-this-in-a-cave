package reality

import "strings"

type Recall struct {
	id string
}

func (rc *Recall) ID() string { return rc.id }

func (rc *Recall) Create(r *Relation) Reality {
	return &Recall{id: "recall"}
}

func (rc *Recall) Realize(r *Relation) string {
	impulse := strings.TrimSpace(r.Impulse)
	if impulse == "" {
		return "no recall query"
	}

	ctx := findContext(r)
	if ctx == nil {
		if r.Log != nil {
			r.Log("[recall]: no context on relation")
		}
		return "no memories"
	}

	query := impulse
	if q, err := ExtractTag(impulse, "about"); err == nil {
		query = q
	}

	relationship := r.Origin
	if rel, err := ExtractTag(impulse, "relationship"); err == nil {
		relationship = rel
	}

	artifactType := ""
	if t, err := ExtractTag(impulse, "type"); err == nil {
		artifactType = t
	}

	if r.Log != nil {
		r.Log("[recall]: query:", query, "relationship:", relationship, "type:", artifactType)
	}

	result := ctx.Retrieve(query, relationship, artifactType)

	if r.Log != nil {
		r.Log("[recall]: →", truncate(result, 80))
	}

	return result
}

func findMemory(r *Relation) *Memory {
	if r.Realities == nil {
		return nil
	}
	if m, ok := r.Realities["memory"]; ok {
		if mem, ok := m.(*Memory); ok {
			return mem
		}
	}
	return nil
}
