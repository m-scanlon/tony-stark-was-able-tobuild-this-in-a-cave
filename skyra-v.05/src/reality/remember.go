package reality

import "strings"

type StoreContext struct {
	id string
}

func (sc *StoreContext) ID() string { return sc.id }

func (sc *StoreContext) Create(r *Relation) Reality {
	return &StoreContext{id: "store-context"}
}

func (sc *StoreContext) Realize(r *Relation) string {
	impulse := strings.TrimSpace(r.Impulse)
	if impulse == "" {
		return "nothing to store"
	}

	ctx := findContext(r)
	if ctx == nil {
		if r.Log != nil {
			r.Log("[store-context]: no context on relation")
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

	relationship := r.Origin

	if r.Log != nil {
		r.Log("[store-context]:", artifactType, "for", relationship, "→", truncate(content, 60))
	}

	result := ctx.Store(content, relationship, artifactType)

	if r.Log != nil {
		r.Log("[store-context]: →", result)
	}

	return result
}
