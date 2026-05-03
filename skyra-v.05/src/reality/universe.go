package reality

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
)

type Universe struct {
	id     string
	Thread *NewThread
	Econ   *Economics
}

func (u *Universe) ID() string { return u.id }

func (u *Universe) Create(r *Relation) Reality {
	return &Universe{id: "universe"}
}

func (u *Universe) Realize(r *Relation) string {
	if r.Collecting {
		cr := &Relation{
			Collecting: true,
			Exports:    make(map[string]any),
		}

		u.Thread.Realize(cr)
		if u.Econ != nil {
			u.Econ.Realize(cr)
		}

		root := RealityNode{ID: "universe", Type: "Universe", Children: []RealityNode{}}
		if threadNode, ok := cr.Exports["node:root"]; ok {
			root.Children = append(root.Children, threadNode.(RealityNode))
			delete(cr.Exports, "node:root")
		}
		if u.Econ != nil {
			root.Children = append(root.Children, RealityNode{ID: "economics", Type: "Economics", Children: []RealityNode{}})
		}
		cr.Exports["node:root"] = root

		state := assembleState(cr.Exports)
		data, err := json.MarshalIndent(state, "", "  ")
		if err != nil {
			return "{}"
		}
		return string(data)
	}

	return u.Thread.Realize(r)
}

type UniverseState struct {
	Beings       []BeingSnapshot    `json:"beings"`
	Threads      []ThreadSnapshot   `json:"threads"`
	Exchanges    []ExchangeSnapshot `json:"exchanges"`
	Economics    map[string]int     `json:"economics"`
	RealityGraph RealityNode        `json:"reality_graph"`
}

type BeingSnapshot struct {
	Name     string          `json:"name"`
	Type     string          `json:"type"`
	Identity string          `json:"identity"`
	Purpose  string          `json:"purpose"`
	Peers       []string        `json:"peers"`
	Entrypoints []string        `json:"entrypoints"`
	Status      string          `json:"status"`
	Device      string          `json:"device"`
	Layers   *LayersSnapshot `json:"layers,omitempty"`
	Memories MemorySnapshot  `json:"memories"`
}

type LayersSnapshot struct {
	Think ThinkSnapshot `json:"think"`
	Act   ActSnapshot   `json:"act"`
}

type ThinkSnapshot struct {
	Budget    int               `json:"budget"`
	Operators []string          `json:"operators"`
	History   []ThoughtSnapshot `json:"history"`
}

type ThoughtSnapshot struct {
	Peer    string `json:"peer"`
	Thought string `json:"thought"`
	Ts      int64  `json:"ts"`
}

type ActSnapshot struct {
	Operators []string `json:"operators"`
}

type MemorySnapshot struct {
	Items  []MemoryItem `json:"items"`
	Skills []SkillItem  `json:"skills"`
}

type MemoryItem struct {
	Filename string `json:"filename"`
	Content  string `json:"content"`
}

type SkillItem struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type ThreadSnapshot struct {
	ID        string         `json:"id"`
	CreatedBy string         `json:"created_by"`
	Active    bool           `json:"active"`
	Members   []string       `json:"members"`
	Edges     []EdgeSnapshot `json:"edges"`
}

type EdgeSnapshot struct {
	From string `json:"from"`
	To   string `json:"to"`
}

type ExchangeSnapshot struct {
	Key      string            `json:"key"`
	Parties  [2]string         `json:"parties"`
	Active   bool              `json:"active"`
	Entries  []EntrySnapshot   `json:"entries"`
	Context  map[string]string `json:"context,omitempty"`
}

type EntrySnapshot struct {
	Index   int    `json:"index"`
	From    string `json:"from"`
	Content string `json:"content"`
	Ts      int64  `json:"ts"`
}

type RealityNode struct {
	ID       string        `json:"id"`
	Type     string        `json:"type"`
	Children []RealityNode `json:"children"`
}

func assembleState(exports map[string]any) UniverseState {
	u := UniverseState{
		Beings:    []BeingSnapshot{},
		Threads:   []ThreadSnapshot{},
		Exchanges: []ExchangeSnapshot{},
		Economics: make(map[string]int),
	}

	for key, val := range exports {
		switch {
		case strings.HasPrefix(key, "being:"):
			u.Beings = append(u.Beings, val.(BeingSnapshot))
		case strings.HasPrefix(key, "thread:"):
			u.Threads = append(u.Threads, val.(ThreadSnapshot))
		case strings.HasPrefix(key, "exchange:"):
			u.Exchanges = append(u.Exchanges, val.(ExchangeSnapshot))
		case key == "economics":
			u.Economics = val.(map[string]int)
		case key == "node:root":
			u.RealityGraph = val.(RealityNode)
		}
	}

	activeBeings := map[string]bool{}
	for _, ex := range u.Exchanges {
		if ex.Active {
			for _, party := range ex.Parties {
				activeBeings[party] = true
			}
		}
	}
	for i := range u.Beings {
		if activeBeings[u.Beings[i].Name] {
			u.Beings[i].Status = "active"
		}
	}

	return u
}

func snapshotMemories(home string) MemorySnapshot {
	snap := MemorySnapshot{
		Items:  []MemoryItem{},
		Skills: []SkillItem{},
	}

	memDir := filepath.Join(home, "memories")
	if entries, err := os.ReadDir(memDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			data, err := os.ReadFile(filepath.Join(memDir, entry.Name()))
			if err != nil {
				continue
			}
			snap.Items = append(snap.Items, MemoryItem{
				Filename: entry.Name(),
				Content:  string(data),
			})
		}
	}

	skillDir := filepath.Join(home, "skills")
	if entries, err := os.ReadDir(skillDir); err == nil {
		for _, entry := range entries {
			if entry.IsDir() {
				continue
			}
			data, err := os.ReadFile(filepath.Join(skillDir, entry.Name()))
			if err != nil {
				continue
			}
			name := entry.Name()
			name = strings.TrimSuffix(name, filepath.Ext(name))
			snap.Skills = append(snap.Skills, SkillItem{
				Name:    name,
				Content: string(data),
			})
		}
	}

	return snap
}

func capitalizeType(name string) string {
	if name == "" {
		return ""
	}
	return strings.ToUpper(name[:1]) + name[1:]
}
