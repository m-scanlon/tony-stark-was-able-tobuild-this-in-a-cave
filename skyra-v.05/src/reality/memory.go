package reality

import (
	"fmt"
	"os"
	"path/filepath"
	"skyra-v05/src/debug"
	"sort"
	"strings"
	"time"
)

type Memory struct {
	id        string
	Owner     string
	Graph     *MemoryGraph
	Extractor *Extractor
	Resolver  *Resolver
	HomeDir   string
}

func NewMemory(owner string) *Memory {
	return &Memory{
		id:        "memory",
		Owner:     owner,
		Graph:     NewMemoryGraph(),
		Extractor: NewExtractor(),
		Resolver:  NewResolver(),
	}
}

func (m *Memory) ID() string { return m.id }

func (m *Memory) Create(r *Relation) Reality {
	return NewMemory(r.ID)
}

func (m *Memory) Realize(r *Relation) string {
	return ""
}

func (m *Memory) Store(content, relationship, artifactType string, layer EdgeLayer, entities []string) string {
	log := func(args ...any) { debug.Being(m.Owner, "memory", args...) }

	var resolved []string
	if len(entities) > 0 {
		for _, e := range entities {
			resolved = append(resolved, m.Resolver.Resolve(e))
		}
	} else {
		extracted := m.Extractor.Extract(content)
		for _, e := range extracted {
			resolved = append(resolved, m.Resolver.Resolve(e))
		}
	}

	now := time.Now()

	for _, name := range resolved {
		entityID := "entity:" + name
		if existing := m.Graph.GetEntity(entityID); existing != nil {
			existing.Weight += 0.1
			existing.LastSeen = now
		} else {
			m.Graph.AddEntity(&Entity{
				ID:        entityID,
				Name:      name,
				Weight:    1.0,
				CreatedAt: now,
				LastSeen:  now,
			})
			m.Extractor.Learn(name)
			m.Resolver.AddAlias(name, name)
		}
	}

	memID := fmt.Sprintf("mem:%d", now.UnixNano())
	node := &MemNode{
		ID:             memID,
		Content:        content,
		Type:           artifactType,
		Weight:         artifactWeight(artifactType),
		Relationship:   relationship,
		AnchorEntities: resolved,
		CreatedAt:      now,
		LastActivated:  now,
	}
	m.Graph.AddNode(node)

	for _, name := range resolved {
		entityID := "entity:" + name
		m.Graph.AddAnchor(memID, entityID)
	}

	for i := 0; i < len(resolved); i++ {
		for j := i + 1; j < len(resolved); j++ {
			fromID := "entity:" + resolved[i]
			toID := "entity:" + resolved[j]
			m.Graph.StrengthenEdge(fromID, toID, layer)
		}
	}

	log("[memory]: stored", artifactType, "in", relationship, "entities:", resolved)
	m.save()
	return fmt.Sprintf("remembered [%s]: %s", artifactType, truncate(content, 60))
}

func (m *Memory) Query(query, relationship, artifactType string) string {
	log := func(args ...any) { debug.Being(m.Owner, "memory", args...) }

	entities := m.Extractor.Extract(query)
	if len(entities) == 0 {
		words := strings.Fields(strings.ToLower(query))
		for _, w := range words {
			if !isCommonWord(w) {
				entities = append(entities, w)
			}
		}
	}

	var results []*MemNode
	seen := map[string]bool{}

	for _, e := range entities {
		resolved := m.Resolver.Resolve(e)
		entityID := "entity:" + resolved
		memories := m.Graph.MemoriesForEntity(entityID)
		for _, node := range memories {
			if seen[node.ID] {
				continue
			}
			if relationship != "" && node.Relationship != relationship {
				continue
			}
			if artifactType != "" && node.Type != artifactType {
				continue
			}
			seen[node.ID] = true
			results = append(results, node)
		}
	}

	if len(results) == 0 {
		all := m.Graph.MemoriesByRelationship(relationship)
		if artifactType != "" {
			var filtered []*MemNode
			for _, n := range all {
				if n.Type == artifactType {
					filtered = append(filtered, n)
				}
			}
			all = filtered
		}
		queryLower := strings.ToLower(query)
		for _, node := range all {
			if strings.Contains(strings.ToLower(node.Content), queryLower) {
				results = append(results, node)
			}
		}
	}

	if len(results) == 0 {
		all := m.Graph.MemoriesByRelationship(relationship)
		if artifactType != "" {
			var filtered []*MemNode
			for _, n := range all {
				if n.Type == artifactType {
					filtered = append(filtered, n)
				}
			}
			all = filtered
		}
		sort.Slice(all, func(i, j int) bool {
			return all[i].CreatedAt.After(all[j].CreatedAt)
		})
		if len(all) > 5 {
			all = all[:5]
		}
		results = all
		if len(results) > 0 {
			log("[memory]: recall", query, "in", relationship, "→ fallback to", len(results), "recent")
		}
	}

	sort.Slice(results, func(i, j int) bool {
		if results[i].Weight != results[j].Weight {
			return results[i].Weight > results[j].Weight
		}
		return results[i].CreatedAt.After(results[j].CreatedAt)
	})

	if len(results) > 10 {
		results = results[:10]
	}

	log("[memory]: recall", query, "in", relationship, "→", len(results), "results")

	if len(results) == 0 {
		return "no relevant memories found"
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("recalled %d memories:\n", len(results)))
	for _, node := range results {
		sb.WriteString(fmt.Sprintf("[%s] %s\n", node.Type, node.Content))
	}
	return sb.String()
}

func (m *Memory) Compress(entries []Entry, relationship string) {
	log := func(args ...any) { debug.Being(m.Owner, "memory", args...) }

	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%s: %s\n", e.From, e.Content))
	}

	m.Store(sb.String(), relationship, "trace", EdgeLayer{Type: "episode", Weight: 1.0}, nil)
	log("[memory]: compressed", len(entries), "entries into trace for", relationship)
}

func (m *Memory) Load() {
	if m.HomeDir == "" {
		return
	}
	dir := filepath.Join(m.HomeDir, "memory")
	m.Graph = LoadMemoryGraph(dir)

	for _, entity := range m.Graph.Entities {
		m.Extractor.Learn(entity.Name)
		m.Resolver.AddAlias(entity.Name, entity.Name)
	}
}

func (m *Memory) SeedSkills(skillsDir string) {
	log := func(args ...any) { debug.Being(m.Owner, "memory", args...) }

	entries, err := os.ReadDir(skillsDir)
	if err != nil {
		log("[memory]: no skills dir:", skillsDir)
		return
	}

	seeded := 0
	for _, entry := range entries {
		if entry.IsDir() || !strings.HasSuffix(entry.Name(), ".md") {
			continue
		}

		name := strings.TrimSuffix(entry.Name(), ".md")
		ref := "skill:" + name

		already := false
		for _, node := range m.Graph.Nodes {
			if node.SourceRef == ref {
				already = true
				break
			}
		}
		if already {
			continue
		}

		data, err := os.ReadFile(filepath.Join(skillsDir, entry.Name()))
		if err != nil {
			continue
		}

		m.Store(string(data), "self", "understanding", EdgeLayer{
			Type:   "skill",
			Ref:    name,
			Weight: 1.0,
		}, nil)

		var newest *MemNode
		for _, node := range m.Graph.Nodes {
			if node.SourceRef == "" && node.Type == "understanding" && node.Relationship == "self" {
				if newest == nil || node.CreatedAt.After(newest.CreatedAt) {
					newest = node
				}
			}
		}
		if newest != nil {
			newest.SourceRef = ref
			m.save()
		}
		seeded++
	}

	if seeded > 0 {
		log("[memory]: seeded", seeded, "skills from", skillsDir)
	}
}

func (m *Memory) save() {
	if m.HomeDir == "" {
		return
	}
	dir := filepath.Join(m.HomeDir, "memory")
	m.Graph.Save(dir)
}

func (m *Memory) GraphStats() (int, int, int) {
	return m.Graph.EntityCount(), m.Graph.NodeCount(), m.Graph.EdgeCount()
}

func artifactWeight(artifactType string) float64 {
	switch artifactType {
	case "understanding":
		return 1.0
	case "tension":
		return 0.8
	case "salience":
		return 0.5
	case "trace":
		return 0.2
	default:
		return 0.3
	}
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
