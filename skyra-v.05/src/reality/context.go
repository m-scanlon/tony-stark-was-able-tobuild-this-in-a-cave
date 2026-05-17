package reality

import (
	"fmt"
	"skyra-v05/src/debug"
	"strings"
)

type Context struct {
	id          string
	Owner       string
	Memory      *Memory
	Providers   map[string]Reality
	Warm        map[string][]*MemNode
	Active      string
	Scope       []string
	Claimed     map[string]string
	Specialists map[string]*Self
	OnPromote   func(*Cluster)
	promoting   bool
}

func (c *Context) ID() string { return c.id }

func (c *Context) Create(r *Relation) Reality {
	return &Context{id: "context", Warm: make(map[string][]*MemNode)}
}

func (c *Context) Realize(r *Relation) string {
	return ""
}

func (c *Context) provider() Reality {
	for _, p := range c.Providers {
		return p
	}
	return nil
}

func (c *Context) Heat(relationship string) {
	log := func(args ...any) { debug.Being(c.Owner, "context", args...) }

	if c.Active == relationship && len(c.Warm[relationship]) > 0 {
		return
	}

	c.Active = relationship

	var loaded []*MemNode
	seen := map[string]bool{}

	entities := c.scopedEntities()
	for _, entity := range entities {
		memories := c.Memory.Graph.MemoriesForEntity(entity.ID)
		for _, node := range memories {
			if node.Relationship != relationship || seen[node.ID] {
				continue
			}
			seen[node.ID] = true
			loaded = append(loaded, node)
		}
	}

	c.Warm[relationship] = loaded
	log("[context]: heated", relationship, "—", len(loaded), "memories")
}

func (c *Context) scopedEntities() []*Entity {
	if c.Scope == nil {
		all := make([]*Entity, 0, len(c.Memory.Graph.Entities))
		for _, e := range c.Memory.Graph.Entities {
			all = append(all, e)
		}
		return all
	}
	var result []*Entity
	for _, id := range c.Scope {
		if e := c.Memory.Graph.GetEntity(id); e != nil {
			result = append(result, e)
		}
	}
	return result
}

func (c *Context) Parse(relationship string) string {
	nodes := c.Warm[relationship]
	if len(nodes) == 0 {
		return ""
	}

	understandings := filterByType(nodes, "understanding")
	tensions := filterByType(nodes, "tension")

	if len(understandings) == 0 && len(tensions) == 0 {
		return ""
	}

	var sb strings.Builder
	sb.WriteString("memory context:\n")
	for _, n := range understandings {
		sb.WriteString("  [understanding] " + n.Content + "\n")
	}
	for _, n := range tensions {
		sb.WriteString("  [tension] " + n.Content + "\n")
	}
	return sb.String()
}

func filterByType(nodes []*MemNode, artifactType string) []*MemNode {
	var result []*MemNode
	for _, n := range nodes {
		if n.Type == artifactType {
			result = append(result, n)
		}
	}
	return result
}

func (c *Context) Store(content, relationship, artifactType string) string {
	log := func(args ...any) { debug.Being(c.Owner, "context", args...) }

	provider := c.provider()
	if provider == nil {
		log("[context]: no provider, storing raw")
		msg := c.Memory.Store(content, relationship, artifactType, EdgeLayer{Type: "episode", Weight: 1.0}, nil)
		c.invalidate(relationship)
		return msg
	}

	candidates := c.Memory.Extractor.Extract(content)
	existing := c.Memory.Query(content, relationship, "")

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
		if len(candidates) > 0 {
			sb.WriteString(fmt.Sprintf("\nentity candidates (from extractor — filter these, add missing ones, drop noise):\n%s\n", strings.Join(candidates, ", ")))
		}
		if existing != "no relevant memories found" {
			sb.WriteString(fmt.Sprintf("\nexisting memories in this relationship:\n%s\n", existing))
		}
		return sb.String()
	})

	result := provider.Realize(lr)
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
			}
		}
	}

	action := "store"
	if a, err := ExtractTag(result, "action"); err == nil {
		action = strings.TrimSpace(a)
	}

	judgment := "complement"
	if j, err := ExtractTag(result, "judgment"); err == nil {
		judgment = strings.TrimSpace(j)
	}

	switch action {
	case "discard":
		log("[context]: discarded as redundant")
		return "already known"
	}

	switch judgment {
	case "supersede":
		if target, err := ExtractTag(result, "supersedes"); err == nil {
			c.decayByContent(target, relationship)
			log("[context]: superseded existing memory")
		}
	case "contradict":
		if target, err := ExtractTag(result, "contradicts"); err == nil {
			c.markTension(target, relationship)
			log("[context]: marked contradiction as tension")
		}
		finalType = "tension"
	}

	msg := c.Memory.Store(cleaned, relationship, finalType, EdgeLayer{Type: "episode", Weight: 1.0}, entities)
	c.invalidate(relationship)

	if len(entities) > 0 {
		log("[context]: curator entities:", entities)
	}

	c.checkPromotion(log)

	return msg
}

func (c *Context) checkPromotion(log func(...any)) {
	if c.OnPromote == nil || c.promoting {
		return
	}
	c.promoting = true
	defer func() { c.promoting = false }()

	exclude := map[string]bool{}
	if c.Claimed != nil {
		for id := range c.Claimed {
			exclude[id] = true
		}
	}

	clusters := c.Memory.Graph.DetectClusters(PromotionThreshold, exclude)
	for _, cluster := range clusters {
		log("[context]: cluster detected, density:", cluster.Density, "entities:", cluster.Entities)
		c.OnPromote(&cluster)
	}
}

func (c *Context) decayByContent(content, relationship string) {
	for _, node := range c.Memory.Graph.Nodes {
		if node.Relationship == relationship && strings.Contains(node.Content, content) {
			node.Weight *= 0.3
		}
	}
}

func (c *Context) markTension(content, relationship string) {
	for _, node := range c.Memory.Graph.Nodes {
		if node.Relationship == relationship && strings.Contains(node.Content, content) {
			node.Type = "tension"
		}
	}
}

func (c *Context) invalidate(relationship string) {
	delete(c.Warm, relationship)
}

func (c *Context) Retrieve(query, relationship, artifactType string) string {
	log := func(args ...any) { debug.Being(c.Owner, "context", args...) }

	c.Heat(relationship)
	memories := c.Memory.Query(query, relationship, artifactType)

	activated := c.activateSpecialists(query)
	if len(activated) == 0 {
		return memories
	}

	var thoughts []specialistThought
	for name, specialist := range activated {
		log("[context]: consulting specialist", name)
		thought := c.consultSpecialist(specialist, query, relationship)
		if thought != "" {
			thoughts = append(thoughts, specialistThought{name: name, thought: thought})
			log("[context]: specialist", name, "→", truncate(thought, 80))
		}
	}

	if len(thoughts) == 0 {
		return memories
	}

	return c.synthesize(query, memories, thoughts, log)
}

type specialistThought struct {
	name    string
	thought string
}

func (c *Context) activateSpecialists(query string) map[string]*Self {
	if c.Specialists == nil || len(c.Specialists) == 0 || c.Memory == nil {
		return nil
	}

	entities := c.Memory.Extractor.Extract(query)
	resolved := make([]string, 0, len(entities))
	for _, e := range entities {
		r := c.Memory.Resolver.Resolve(e)
		resolved = append(resolved, "entity:"+r)
	}

	activated := map[string]*Self{}
	for _, eid := range resolved {
		if specName, ok := c.Claimed[eid]; ok {
			if spec, ok := c.Specialists[specName]; ok {
				activated[specName] = spec
			}
		}
	}
	return activated
}

func (c *Context) consultSpecialist(specialist *Self, query, relationship string) string {
	think, ok := specialist.Realities["think"]
	if !ok {
		return ""
	}

	lr := &Relation{
		Impulse:   query,
		Origin:    relationship,
		Parsers:   make(map[string]Parser),
		Realities: specialist.Realities,
	}

	return think.Realize(lr)
}

func (c *Context) synthesize(query, memories string, thoughts []specialistThought, log func(...any)) string {
	provider := c.provider()
	if provider == nil {
		var sb strings.Builder
		sb.WriteString(memories)
		sb.WriteString("\n\nspecialist perspectives:\n")
		for _, t := range thoughts {
			sb.WriteString(fmt.Sprintf("  [%s]: %s\n", t.name, t.thought))
		}
		return sb.String()
	}

	lr := &Relation{
		Impulse: query,
		Parsers: make(map[string]Parser),
	}

	lr.Attach("system", func() string {
		return "You are a synthesis layer. You receive a query, memory results, and specialist perspectives. Synthesize the information into one coherent response. Return only the synthesis — no tags, no commentary."
	})

	lr.Attach("synthesis-input", func() string {
		var sb strings.Builder
		sb.WriteString("query: " + query + "\n\n")
		sb.WriteString("memories:\n" + memories + "\n\n")
		sb.WriteString("specialist perspectives:\n")
		for _, t := range thoughts {
			sb.WriteString(fmt.Sprintf("  [%s]: %s\n", t.name, t.thought))
		}
		return sb.String()
	})

	result := provider.Realize(lr)
	log("[context]: synthesis →", truncate(result, 80))
	return result
}

func (c *Context) system() string {
	return `You are a memory curator. Your job is to evaluate new memories against existing ones before storage.

You receive content that a being wants to remember, along with existing memories in the same relationship.

Your responsibilities:
1. Clean the content — tighten language, strip noise, keep signal
2. Curate entities — you may receive entity candidates from an extractor. Filter out noise (common words, pronouns, vague terms). Keep real names, concepts, tools, and patterns. Add any the extractor missed. Your entity list is authoritative — only what you return enters the graph.
3. Judge the relationship to existing memories:
   - supersede: the new memory replaces an old understanding. Name what it supersedes.
   - complement: the new memory adds nuance without replacing. Both keep their weight.
   - contradict: the new memory conflicts with an existing one. The old becomes a tension.
   - discard: this is already fully covered by existing memories.
4. Classify — confirm or reclassify the artifact type (trace, salience, tension, understanding)

Respond with exactly these tags:

<content>the cleaned content to store</content>
<type>trace|salience|tension|understanding</type>
<entities>comma,separated,entity,names</entities>
<judgment>supersede|complement|contradict</judgment>
<action>store|discard</action>

If judgment is supersede, also include:
<supersedes>the content of the memory being replaced</supersedes>

If judgment is contradict, also include:
<contradicts>the content of the conflicting memory</contradicts>

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
