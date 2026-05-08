package reality

import "strings"

type Remember struct {
	id string
}

func (rm *Remember) ID() string { return rm.id }

func (rm *Remember) Create(r *Relation) Reality {
	return &Remember{id: "remember"}
}

func (rm *Remember) Realize(r *Relation) string {
	impulse := strings.TrimSpace(r.Impulse)
	if impulse == "" {
		return "nothing to remember"
	}

	ctx := findContext(r)
	if ctx == nil {
		if r.Log != nil {
			r.Log("[remember]: no context on relation")
		}
		return "no memory available"
	}

	artifactType := "trace"
	if t, err := ExtractTag(impulse, "type"); err == nil {
		artifactType = t
	}

	content := impulse
	if c, err := ExtractTag(impulse, "content"); err == nil {
		content = c
	}

	var contextArtifacts []string
	if ca, err := ExtractTag(impulse, "context"); err == nil {
		for _, a := range strings.Split(ca, ",") {
			a = strings.TrimSpace(a)
			if a != "" {
				contextArtifacts = append(contextArtifacts, a)
			}
		}
	}

	relationship := r.Origin

	if r.Log != nil {
		r.Log("[remember]:", artifactType, "for", relationship, "→", truncate(content, 60))
	}

	result := ctx.Store(content, relationship, artifactType, contextArtifacts)

	if r.Log != nil {
		r.Log("[remember]: →", result)
	}

	return result
}
