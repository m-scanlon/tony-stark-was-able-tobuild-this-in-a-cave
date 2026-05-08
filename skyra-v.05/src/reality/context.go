package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)

type Context struct {
	id     string
	Owner  string
	Memory *Memory
	LLM    Reality
	Warm   map[string][]*MemNode
	Active string
}

func (c *Context) ID() string { return c.id }

func (c *Context) Create(r *Relation) Reality {
	return &Context{id: "context", Warm: make(map[string][]*MemNode)}
}

func (c *Context) Realize(r *Relation) string {
	return ""
}

func (c *Context) Heat(relationship string) {
	log := func(args ...any) { debug.Being(c.Owner, "context", args...) }

	if c.Active == relationship && len(c.Warm[relationship]) > 0 {
		return
	}

	c.Active = relationship

	entities := c.Memory.Graph.EntitiesByRelationship(relationship)

	var loaded []*MemNode
	seen := map[string]bool{}
	for _, entity := range entities {
		neighbors := c.Memory.Graph.ConnectedByType(entity.ID, "mentions")
		for _, node := range neighbors {
			if node.Type == "memory" && !seen[node.ID] {
				seen[node.ID] = true
				loaded = append(loaded, node)
			}
		}
	}

	c.Warm[relationship] = loaded
	log("[context]: heated", relationship, "—", len(entities), "entities,", len(loaded), "memories")
}

func (c *Context) Parse(relationship string) string {
	nodes := c.Warm[relationship]
	if len(nodes) == 0 {
		return ""
	}

	understandings := filterByArtifact(nodes, "understanding")
	tensions := filterByArtifact(nodes, "tension")

	if len(understandings) == 0 && len(tensions) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("memory context:\n")
	if len(understandings) > 0 {
		for _, n := range understandings {
			sb.WriteString("  [understanding] " + n.Content + "\n")
		}
	}
	if len(tensions) > 0 {
		for _, n := range tensions {
			sb.WriteString("  [tension] " + n.Content + "\n")
		}
	}
	return sb.String()
}

func filterByArtifact(nodes []*MemNode, artifactType string) []*MemNode {
	var result []*MemNode
	for _, n := range nodes {
		if n.ArtifactType == artifactType {
			result = append(result, n)
		}
	}
	return result
}

func (c *Context) Store(content, relationship, artifactType string, contextArtifacts []string) string {
	log := func(args ...any) { debug.Being(c.Owner, "context", args...) }

	if c.LLM == nil {
		log("[context]: no llm, storing raw")
		msg := c.Memory.StoreArtifact(content, relationship, artifactType, contextArtifacts)
		c.invalidate(relationship)
		return msg
	}

	existing := c.Memory.QueryGraph(content, relationship, "")

	lr := &Relation{
		Impulse: content,
		Parsers: make(map[string]Parser),
	}

	lr.Attach("system", c.system)
	lr.Attach("context-job", func() string {
		var sb strings.Builder
		sb.WriteString(fmt.Sprintf("artifact type: %s\n", artifactType))
		sb.WriteString(fmt.Sprintf("relationship: %s\n", relationship))
		sb.WriteString(fmt.Sprintf("content to store:\n%s\n", content))
		if existing != "no relevant memories found" {
			sb.WriteString(fmt.Sprintf("\nexisting memories in this relationship:\n%s\n", existing))
		}
		return sb.String()
	})

	result := c.LLM.Realize(lr)
	log("[context]: curated →", result)

	cleaned := content
	if cl, err := ExtractTag(result, "content"); err == nil {
		cleaned = cl
	}

	finalType := artifactType
	if t, err := ExtractTag(result, "type"); err == nil {
		if isValidArtifactType(t) {
			finalType = t
		}
	}

	var entities []string
	if e, err := ExtractTag(result, "entities"); err == nil {
		for _, entity := range strings.Split(e, ",") {
			entity = strings.TrimSpace(entity)
			if entity != "" {
				entities = append(entities, entity)
				c.Memory.Extractor.Learn(entity)
			}
		}
	}

	if action, err := ExtractTag(result, "action"); err == nil {
		if strings.TrimSpace(action) == "discard" {
			log("[context]: discarded as redundant")
			return "already known"
		}
	}

	msg := c.Memory.StoreArtifact(cleaned, relationship, finalType, contextArtifacts)
	c.invalidate(relationship)

	if len(entities) > 0 {
		log("[context]: learned entities:", entities)
	}

	return msg
}

func (c *Context) invalidate(relationship string) {
	delete(c.Warm, relationship)
}

func (c *Context) Retrieve(query, relationship, artifactType string) string {
	c.Heat(relationship)
	return c.Memory.QueryGraph(query, relationship, artifactType)
}

func (c *Context) system() string {
	return `You are a memory curator. Your job is to keep memory clean before it goes into storage.

You receive content that a being wants to remember, along with existing memories in the same relationship.

Your responsibilities:
1. Clean the content — tighten language, strip noise, keep signal
2. Extract entities — names, concepts, tools, patterns mentioned in the content
3. Deduplicate — if this is already covered by an existing memory, discard it
4. Classify — confirm or reclassify the artifact type (trace, salience, tension, understanding)
5. Merge — if this updates an existing memory, produce the merged version

Respond with exactly these tags:

<content>the cleaned content to store</content>
<type>trace|salience|tension|understanding</type>
<entities>comma,separated,entity,names</entities>
<action>store|discard</action>

If discarding, still include all tags but set content to empty.
Keep your response to just the tags. No commentary.`
}

func isValidArtifactType(t string) bool {
	switch t {
	case "trace", "salience", "tension", "understanding":
		return true
	}
	return false
}

func findContext(r *Relation) *Context {
	if r.Realities == nil {
		return nil
	}
	if c, ok := r.Realities["context"]; ok {
		if ctx, ok := c.(*Context); ok {
			return ctx
		}
	}
	return nil
}
