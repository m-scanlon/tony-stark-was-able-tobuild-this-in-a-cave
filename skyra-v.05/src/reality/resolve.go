package reality

import "strings"

type Resolver struct {
	Aliases map[string]string
}

func NewResolver() *Resolver {
	return &Resolver{Aliases: make(map[string]string)}
}

func (r *Resolver) AddAlias(alias, canonical string) {
	r.Aliases[normalizeEntity(alias)] = normalizeEntity(canonical)
}

func (r *Resolver) Resolve(name string) string {
	norm := normalizeEntity(name)
	if canonical, ok := r.Aliases[norm]; ok {
		return canonical
	}
	best := ""
	bestDist := 3
	for alias, canonical := range r.Aliases {
		d := levenshtein(norm, alias)
		if d > 0 && d < bestDist {
			bestDist = d
			best = canonical
		}
	}
	if best != "" {
		return best
	}
	return norm
}

func normalizeEntity(s string) string {
	s = strings.ToLower(strings.TrimSpace(s))
	fields := strings.Fields(s)
	return strings.Join(fields, " ")
}

func levenshtein(a, b string) int {
	if len(a) == 0 {
		return len(b)
	}
	if len(b) == 0 {
		return len(a)
	}
	prev := make([]int, len(b)+1)
	curr := make([]int, len(b)+1)
	for j := 0; j <= len(b); j++ {
		prev[j] = j
	}
	for i := 1; i <= len(a); i++ {
		curr[0] = i
		for j := 1; j <= len(b); j++ {
			cost := 1
			if a[i-1] == b[j-1] {
				cost = 0
			}
			curr[j] = min(curr[j-1]+1, min(prev[j]+1, prev[j-1]+cost))
		}
		prev, curr = curr, prev
	}
	return prev[len(b)]
}
