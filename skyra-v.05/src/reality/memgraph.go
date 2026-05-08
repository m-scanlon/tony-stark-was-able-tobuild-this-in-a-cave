package reality

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sort"
	"time"
)

type MemNode struct {
	ID               string    `json:"id"`
	Type             string    `json:"type"`
	Content          string    `json:"content"`
	Vector           []float64 `json:"vector,omitempty"`
	Weight           float64   `json:"weight"`
	CreatedAt        time.Time `json:"created_at"`
	LastSeen         time.Time `json:"last_seen"`
	ArtifactType     string    `json:"artifact_type,omitempty"`
	Relationship     string    `json:"relationship"`
	AnchorEntities   []string  `json:"anchor_entities,omitempty"`
	ContextArtifacts []string  `json:"context_artifacts,omitempty"`
	TrustAtFormation float64   `json:"trust_at_formation,omitempty"`
}

type MemEdge struct {
	From      string        `json:"from"`
	To        string        `json:"to"`
	Type      string        `json:"type"`
	Weight    float64       `json:"weight"`
	History   []WeightEntry `json:"history,omitempty"`
	CreatedAt time.Time     `json:"created_at"`
	LastSeen  time.Time     `json:"last_seen"`
}

type WeightEntry struct {
	Weight float64   `json:"weight"`
	At     time.Time `json:"at"`
}

type MemoryGraph struct {
	Nodes map[string]*MemNode `json:"nodes"`
	Edges []*MemEdge          `json:"edges"`
	adj   map[string][]*MemEdge
}

func NewMemoryGraph() *MemoryGraph {
	return &MemoryGraph{
		Nodes: make(map[string]*MemNode),
		Edges: []*MemEdge{},
		adj:   make(map[string][]*MemEdge),
	}
}

func (g *MemoryGraph) AddNode(node *MemNode) {
	g.Nodes[node.ID] = node
}

func (g *MemoryGraph) GetNode(id string) *MemNode {
	return g.Nodes[id]
}

func (g *MemoryGraph) AddEdge(edge *MemEdge) {
	g.Edges = append(g.Edges, edge)
	g.adj[edge.From] = append(g.adj[edge.From], edge)
	g.adj[edge.To] = append(g.adj[edge.To], edge)
}

func (g *MemoryGraph) Neighbors(nodeID string, depth int) []*MemNode {
	visited := map[string]bool{nodeID: true}
	frontier := []string{nodeID}
	var result []*MemNode

	for d := 0; d < depth && len(frontier) > 0; d++ {
		var next []string
		for _, id := range frontier {
			for _, edge := range g.adj[id] {
				other := edge.To
				if other == id {
					other = edge.From
				}
				if !visited[other] {
					visited[other] = true
					if node := g.Nodes[other]; node != nil {
						result = append(result, node)
					}
					next = append(next, other)
				}
			}
		}
		frontier = next
	}
	return result
}

func (g *MemoryGraph) ConnectedByType(nodeID, edgeType string) []*MemNode {
	var result []*MemNode
	seen := map[string]bool{}
	for _, edge := range g.adj[nodeID] {
		if edge.Type != edgeType {
			continue
		}
		other := edge.To
		if other == nodeID {
			other = edge.From
		}
		if !seen[other] {
			seen[other] = true
			if node := g.Nodes[other]; node != nil {
				result = append(result, node)
			}
		}
	}
	return result
}

func (g *MemoryGraph) EntitiesByRelationship(rel string) []*MemNode {
	var result []*MemNode
	for _, node := range g.Nodes {
		if node.Type == "entity" && node.Relationship == rel {
			result = append(result, node)
		}
	}
	sort.Slice(result, func(i, j int) bool {
		return result[i].Weight > result[j].Weight
	})
	return result
}

func (g *MemoryGraph) MemoriesByRelationship(rel string) []*MemNode {
	var result []*MemNode
	for _, node := range g.Nodes {
		if node.Type == "memory" && node.Relationship == rel {
			result = append(result, node)
		}
	}
	return result
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
	if g.Nodes == nil {
		g.Nodes = make(map[string]*MemNode)
	}
	if g.Edges == nil {
		g.Edges = []*MemEdge{}
	}
	g.rebuildAdj()
	return g
}

func (g *MemoryGraph) rebuildAdj() {
	g.adj = make(map[string][]*MemEdge)
	for _, edge := range g.Edges {
		g.adj[edge.From] = append(g.adj[edge.From], edge)
		g.adj[edge.To] = append(g.adj[edge.To], edge)
	}
}

func (g *MemoryGraph) NodeCount() int {
	return len(g.Nodes)
}

func (g *MemoryGraph) EdgeCount() int {
	return len(g.Edges)
}
