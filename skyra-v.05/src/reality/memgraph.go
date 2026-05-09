package reality

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type Entity struct {
	ID        string    `json:"id"`
	Name      string    `json:"name"`
	Weight    float64   `json:"weight"`
	CreatedAt time.Time `json:"created_at"`
	LastSeen  time.Time `json:"last_seen"`
}

type EntityEdge struct {
	From      string      `json:"from"`
	To        string      `json:"to"`
	Weight    float64     `json:"weight"`
	Layers    []EdgeLayer `json:"layers"`
	CreatedAt time.Time   `json:"created_at"`
	LastSeen  time.Time   `json:"last_seen"`
}

type EdgeLayer struct {
	Type   string  `json:"type"` // episode, task, skill
	Ref    string  `json:"ref"`
	Weight float64 `json:"weight"`
}

type MemNode struct {
	ID              string    `json:"id"`
	Content         string    `json:"content"`
	Type            string    `json:"type"` // trace, salience, tension, understanding
	Weight          float64   `json:"weight"`
	ActivationCount int       `json:"activation_count"`
	Relationship    string    `json:"relationship"`
	AnchorEntities  []string  `json:"anchor_entities,omitempty"`
	CreatedAt       time.Time `json:"created_at"`
	LastActivated   time.Time `json:"last_activated"`
}

type MemEdge struct {
	From string `json:"from"`
	To   string `json:"to"`
	Type string `json:"type"` // anchors
}

type MemoryGraph struct {
	Entities    map[string]*Entity     `json:"entities"`
	EntityEdges map[string]*EntityEdge `json:"entity_edges"`
	Nodes       map[string]*MemNode    `json:"nodes"`
	Anchors     []*MemEdge             `json:"anchors"`
	anchorAdj   map[string][]string
}

func NewMemoryGraph() *MemoryGraph {
	return &MemoryGraph{
		Entities:    make(map[string]*Entity),
		EntityEdges: make(map[string]*EntityEdge),
		Nodes:       make(map[string]*MemNode),
		Anchors:     []*MemEdge{},
		anchorAdj:   make(map[string][]string),
	}
}

func entityEdgeKey(a, b string) string {
	if a < b {
		return a + "|" + b
	}
	return b + "|" + a
}

func (g *MemoryGraph) AddEntity(e *Entity) {
	g.Entities[e.ID] = e
}

func (g *MemoryGraph) GetEntity(id string) *Entity {
	return g.Entities[id]
}

func (g *MemoryGraph) EntityCount() int {
	return len(g.Entities)
}

func (g *MemoryGraph) StrengthenEdge(fromEntity, toEntity string, layer EdgeLayer) {
	key := entityEdgeKey(fromEntity, toEntity)
	now := time.Now()

	edge, ok := g.EntityEdges[key]
	if !ok {
		edge = &EntityEdge{
			From:      fromEntity,
			To:        toEntity,
			CreatedAt: now,
		}
		g.EntityEdges[key] = edge
	}

	found := false
	for i := range edge.Layers {
		if edge.Layers[i].Type == layer.Type && edge.Layers[i].Ref == layer.Ref {
			edge.Layers[i].Weight += layer.Weight
			found = true
			break
		}
	}
	if !found {
		edge.Layers = append(edge.Layers, layer)
	}

	edge.Weight = 0
	for _, l := range edge.Layers {
		edge.Weight += l.Weight
	}
	edge.LastSeen = now
}

func (g *MemoryGraph) GetEdge(a, b string) *EntityEdge {
	return g.EntityEdges[entityEdgeKey(a, b)]
}

func (g *MemoryGraph) AddNode(node *MemNode) {
	g.Nodes[node.ID] = node
}

func (g *MemoryGraph) GetNode(id string) *MemNode {
	return g.Nodes[id]
}

func (g *MemoryGraph) AddAnchor(memID, entityID string) {
	g.Anchors = append(g.Anchors, &MemEdge{
		From: memID,
		To:   entityID,
		Type: "anchors",
	})
	g.anchorAdj[entityID] = append(g.anchorAdj[entityID], memID)
	g.anchorAdj[memID] = append(g.anchorAdj[memID], entityID)
}

func (g *MemoryGraph) MemoriesForEntity(entityID string) []*MemNode {
	var result []*MemNode
	seen := map[string]bool{}
	for _, memID := range g.anchorAdj[entityID] {
		if seen[memID] {
			continue
		}
		seen[memID] = true
		if node := g.Nodes[memID]; node != nil {
			result = append(result, node)
		}
	}
	return result
}

func (g *MemoryGraph) EntitiesForMemory(memID string) []*Entity {
	var result []*Entity
	seen := map[string]bool{}
	for _, entityID := range g.anchorAdj[memID] {
		if seen[entityID] {
			continue
		}
		seen[entityID] = true
		if entity := g.Entities[entityID]; entity != nil {
			result = append(result, entity)
		}
	}
	return result
}

func (g *MemoryGraph) Neighbors(entityID string) []*Entity {
	var result []*Entity
	seen := map[string]bool{entityID: true}
	for key, edge := range g.EntityEdges {
		_ = key
		other := ""
		if edge.From == entityID {
			other = edge.To
		} else if edge.To == entityID {
			other = edge.From
		}
		if other != "" && !seen[other] {
			seen[other] = true
			if e := g.Entities[other]; e != nil {
				result = append(result, e)
			}
		}
	}
	return result
}

func (g *MemoryGraph) MemoriesByRelationship(rel string) []*MemNode {
	var result []*MemNode
	for _, node := range g.Nodes {
		if node.Relationship == rel {
			result = append(result, node)
		}
	}
	return result
}

func (g *MemoryGraph) EntitiesByWeight(limit int) []*Entity {
	all := make([]*Entity, 0, len(g.Entities))
	for _, e := range g.Entities {
		all = append(all, e)
	}
	sort.Slice(all, func(i, j int) bool {
		return all[i].Weight > all[j].Weight
	})
	if limit > 0 && len(all) > limit {
		all = all[:limit]
	}
	return all
}

func (g *MemoryGraph) NodeCount() int {
	return len(g.Nodes)
}

func (g *MemoryGraph) EdgeCount() int {
	return len(g.EntityEdges)
}

func (g *MemoryGraph) Save(dir string) error {
	if err := os.MkdirAll(dir, 0755); err != nil {
		return err
	}
	data, err := json.MarshalIndent(g, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filepath.Join(dir, "graph.json"), data, 0644)
}

func LoadMemoryGraph(dir string) *MemoryGraph {
	data, err := os.ReadFile(filepath.Join(dir, "graph.json"))
	if err != nil {
		return NewMemoryGraph()
	}
	g := &MemoryGraph{}
	if err := json.Unmarshal(data, g); err != nil {
		return NewMemoryGraph()
	}
	if g.Entities == nil {
		g.Entities = make(map[string]*Entity)
	}
	if g.EntityEdges == nil {
		g.EntityEdges = make(map[string]*EntityEdge)
	}
	if g.Nodes == nil {
		g.Nodes = make(map[string]*MemNode)
	}
	if g.Anchors == nil {
		g.Anchors = []*MemEdge{}
	}
	g.rebuildAdj()
	return g
}

func (g *MemoryGraph) rebuildAdj() {
	g.anchorAdj = make(map[string][]string)
	for _, anchor := range g.Anchors {
		g.anchorAdj[anchor.To] = append(g.anchorAdj[anchor.To], anchor.From)
		g.anchorAdj[anchor.From] = append(g.anchorAdj[anchor.From], anchor.To)
	}
}
