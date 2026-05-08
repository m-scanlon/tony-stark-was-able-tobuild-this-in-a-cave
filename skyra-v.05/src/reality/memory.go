package reality

import (
	"fmt"
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
	Vectors   *VecIndex
	Embed     func(string) ([]float64, error)
	HomeDir   string
}

func NewMemory(owner string) *Memory {
	return &Memory{
		id:        "memory",
		Owner:     owner,
		Graph:     NewMemoryGraph(),
		Extractor: NewExtractor(),
		Resolver:  NewResolver(),
		Vectors:   NewVecIndex(),
	}
}

func (m *Memory) ID() string { return m.id }

func (m *Memory) Create(r *Relation) Reality {
	return NewMemory(r.ID)
}

func (m *Memory) Realize(r *Relation) string {
	if r.Collecting {
		return ""
	}
	return ""
}

func (m *Memory) StoreArtifact(content, relationship, artifactType string, contextArtifacts []string) string {
	log := func(args ...any) { debug.Being(m.Owner, "memory", args...) }

	entities := m.Extractor.Extract(content)

	resolved := make([]string, 0, len(entities))
	for _, e := range entities {
		resolved = append(resolved, m.Resolver.Resolve(e))
	}

	now := time.Now()

	for _, e := range resolved {
		entityID := "entity:" + relationship + ":" + e
		if existing := m.Graph.GetNode(entityID); existing != nil {
			existing.Weight += 0.1
			existing.LastSeen = now
		} else {
			m.Graph.AddNode(&MemNode{
				ID:           entityID,
				Type:         "entity",
				Content:      e,
				Weight:       1.0,
				Relationship: relationship,
				CreatedAt:    now,
				LastSeen:     now,
			})
			m.Extractor.Learn(e)
			m.Resolver.AddAlias(e, e)
		}
	}

	memID := fmt.Sprintf("mem:%d", now.UnixNano())
	memNode := &MemNode{
		ID:               memID,
		Type:             "memory",
		Content:          content,
		ArtifactType:     artifactType,
		Relationship:     relationship,
		AnchorEntities:   resolved,
		ContextArtifacts: contextArtifacts,
		Weight:           artifactWeight(artifactType),
		CreatedAt:        now,
		LastSeen:         now,
	}

	if m.Embed != nil {
		if vec, err := m.Embed(content); err == nil {
			memNode.Vector = vec
			m.Vectors.Add(memID, vec)
		}
	}

	m.Graph.AddNode(memNode)

	for _, e := range resolved {
		entityID := "entity:" + relationship + ":" + e
		m.Graph.AddEdge(&MemEdge{
			From:      memID,
			To:        entityID,
			Type:      "mentions",
			Weight:    1.0,
			CreatedAt: now,
			LastSeen:  now,
		})
	}

	log("[memory]: stored", artifactType, "in", relationship, "entities:", resolved)
	m.save()
	return fmt.Sprintf("remembered [%s]: %s", artifactType, truncate(content, 60))
}

func (m *Memory) QueryGraph(query, relationship, artifactType string) string {
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
		entityID := "entity:" + relationship + ":" + resolved
		connected := m.Graph.ConnectedByType(entityID, "mentions")
		for _, node := range connected {
			if node.Type != "memory" || seen[node.ID] {
				continue
			}
			if node.Relationship != relationship {
				continue
			}
			if artifactType != "" && node.ArtifactType != artifactType {
				continue
			}
			seen[node.ID] = true
			results = append(results, node)
		}
	}

	if m.Embed != nil && len(m.Vectors.Vectors) > 0 {
		if queryVec, err := m.Embed(query); err == nil {
			vecResults := m.Vectors.Search(queryVec, 5)
			for _, vr := range vecResults {
				if seen[vr.ID] || vr.Score < 0.7 {
					continue
				}
				if node := m.Graph.GetNode(vr.ID); node != nil {
					if node.Relationship == relationship && node.Type == "memory" {
						if artifactType == "" || node.ArtifactType == artifactType {
							seen[vr.ID] = true
							results = append(results, node)
						}
					}
				}
			}
		}
	}

	if len(results) == 0 {
		all := m.Graph.MemoriesByRelationship(relationship)
		if artifactType != "" {
			var filtered []*MemNode
			for _, n := range all {
				if n.ArtifactType == artifactType {
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
				if n.ArtifactType == artifactType {
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
		sb.WriteString(fmt.Sprintf("[%s] %s\n", node.ArtifactType, node.Content))
	}
	return sb.String()
}

func (m *Memory) Compress(entries []Entry, relationship string) {
	log := func(args ...any) { debug.Being(m.Owner, "memory", args...) }

	var sb strings.Builder
	for _, e := range entries {
		sb.WriteString(fmt.Sprintf("%s: %s\n", e.From, e.Content))
	}
	content := sb.String()

	m.StoreArtifact(content, relationship, "trace", nil)
	log("[memory]: compressed", len(entries), "entries into trace for", relationship)
}

func (m *Memory) Load() {
	if m.HomeDir == "" {
		return
	}
	dir := filepath.Join(m.HomeDir, "memory")
	m.Graph = LoadMemoryGraph(dir)

	for _, node := range m.Graph.Nodes {
		if node.Type == "entity" {
			m.Extractor.Learn(node.Content)
			m.Resolver.AddAlias(node.Content, node.Content)
		}
		if node.Vector != nil && len(node.Vector) > 0 {
			m.Vectors.Add(node.ID, node.Vector)
		}
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
	entities := 0
	memories := 0
	for _, node := range m.Graph.Nodes {
		switch node.Type {
		case "entity":
			entities++
		case "memory":
			memories++
		}
	}
	return entities, memories, m.Graph.EdgeCount()
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

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n] + "..."
}
