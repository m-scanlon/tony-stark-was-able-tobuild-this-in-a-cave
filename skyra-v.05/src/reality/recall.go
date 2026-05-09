package reality

import "strings"

type RetrieveContext struct {
	id string
}

func (rc *RetrieveContext) ID() string { return rc.id }

func (rc *RetrieveContext) Create(r *Relation) Reality {
	return &RetrieveContext{id: "retrieve-context"}
}

func (rc *RetrieveContext) Realize(r *Relation) string {
	impulse := strings.TrimSpace(r.Impulse)
	if impulse == "" {
		return "no query"
	}

	ctx := findContext(r)
	if ctx == nil {
		if r.Log != nil {
			r.Log("[retrieve-context]: no context on relation")
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
		r.Log("[retrieve-context]: query:", query, "relationship:", relationship, "type:", artifactType)
	}

	result := ctx.Retrieve(query, relationship, artifactType)

	if r.Log != nil {
		r.Log("[retrieve-context]: →", truncate(result, 80))
	}

	return result
}
