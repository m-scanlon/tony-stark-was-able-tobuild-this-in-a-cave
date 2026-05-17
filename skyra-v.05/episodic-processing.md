# Episodic Processing

## The Two-Stream Model

Every reality that holds experience retains two views of the same data:

- **Full trace** — what actually happened. Retained until Context has processed it. Then cleared.
- **Present view** — what the being sees going forward. Trimmed for the context window.

This pattern applies to both Exchange and Think.

### Exchange

- **Full history** — every entry from both parties. What was said, in order. Retained intact until Context processes it.
- **Present view** — last N entries rendered into the being's present via the parser. This is what the LLM sees during conversation.

Compaction is a view concern, not a data concern. The present trims for the context window. The full history stays until processed.

Current code conflates these — compaction trims the entries AND destroys the data by dumping it into a raw trace. These need to split.

### Think

- **Full trace** — every pass of a thinking session. Every operator call, every result returned, every intermediate thought. The complete cognitive cycle.
- **Present view** — the surface thought. What escapes to Act, what gets shown in thought history going forward.

Currently the intermediate passes exist only as local variables in `Think.Realize()` (the `exchange []thinkEntry` slice). They're thrown away when the function returns. Think needs to retain the full trace per session until Context processes it.

```
Think session:
  pass 0: thought X
  pass 1: <retrieve-context>michael</retrieve-context> → got Y
  pass 2: <store-context>...</store-context> → stored
  pass 3: <surface-thought>final synthesis</surface-thought>

Present view: final synthesis only
Full trace: all 4 passes (retained for Context)
```

## The Trigger

Episodic processing fires at compaction time. The exchange accumulates entries. When it hits the threshold (currently 20), the oldest batch is ready for processing. Same trigger that exists today — the work is different.

No timers. No silence detection. No topic shift parsing. The exchange overflows, Context processes the overflow. Rolling window.

Later refinement: silence-based trigger (consolidation during rest). Not now.

## What Context Processes

Context receives both streams and processes them together:

- **Exchange stream** — what happened externally. Who said what to whom.
- **Think stream** — what happened internally. What the being deliberated, what it recalled, what it considered.

The gap between Think and Act is information. A being that thinks deeply about websocket timeouts but says something simple is still learning about websocket timeouts. The weight update reflects the full cognitive cycle, not just the output.

## The Processing Loop

When the trigger fires, Context receives the batch and processes it:

1. **Chunking** — the batch is split into digestible pieces. Fixed-size chunks (N turn-pairs per chunk). Semantic chunking (topic-shift detection) is a future refinement.

2. **Per-chunk extraction** — for each chunk, Context calls its provider:
   - Extract entities mentioned
   - Detect skill usage (which operators were called, how)
   - Identify what was learned (new understanding, resolved tensions, patterns)
   - Assign weight updates (what got reinforced, what got superseded)

3. **Graph updates** — each chunk's output feeds the graph:
   - Entity weights updated based on activation
   - Entity-to-entity edge weights strengthened where co-occurrence happened
   - New memory nodes stored with appropriate types (trace, understanding, etc.)
   - Skill edges updated when skill usage detected

4. **Skill file rewrites** — after all chunks processed, check skill edge weights. If any skill region crosses the maturation threshold, Context rewrites the skill file from lived experience.

5. **Cleanup** — full history and full trace cleared. The graph holds what was learned. The raw data is gone.

## What This Replaces

Currently `Memory.Compress()` takes old exchange entries and dumps them as a single raw trace node:

```go
func (m *Memory) Compress(entries []Entry, relationship string) {
    var sb strings.Builder
    for _, e := range entries {
        sb.WriteString(fmt.Sprintf("%s: %s\n", e.From, e.Content))
    }
    m.Store(sb.String(), relationship, "trace", EdgeLayer{Type: "episode", Weight: 1.0})
}
```

This becomes real episodic processing. Same trigger point (compaction threshold in Exchange.Realize), but instead of raw dump → proper extraction, entity strengthening, typed memory storage.

## The Extraction Framework

Episodic processing is not `Store()` in a loop. It's a batch of experience that needs to become graph structure. The batch is seen as a whole — patterns emerge across entries that aren't visible in any single one.

### Open Questions

1. **What gets extracted?** Not every entry is an entity. Not every exchange is significant. The batch needs to be reduced to signal. What's noise vs what's structure?

2. **What type of memory?** A batch might produce traces (what happened), understandings (what was learned), tensions (what conflicted), salience (what mattered). One batch might produce all four types.

3. **How much weight?** Not all memories are equal. Something the being struggled with across 5 passes of Think is heavier than something it handled in one pass. Depth of processing = weight.

4. **What's the grain?** Does the framework produce one memory node per batch? Multiple? One per detected topic? One per entity cluster activated?

5. **Edges** — co-occurrence isn't "they appeared in the same entry." It's "they were part of the same thread of thought across multiple turns." Co-occurrence is at the episode level, not the entry level.

### The Difference From Store()

Current `Store()`: content in → extract entities → create node → strengthen edges. Mechanical. One piece of content at a time.

Episodic processing: batch in → what happened here? → what matters? → what type of structure does this become? The middle step needs the LLM. The framework is the prompt structure and output format for that extraction call.

## Implementation Shape

### Exchange changes
- Split `Entries` into `History []Entry` (full, retained) and keep the present view as a sliding window over it
- Compaction trims the present view but does NOT destroy the history
- After Context processes a batch, those entries are cleared from history

### Think changes
- Retain the full `[]thinkEntry` trace per session on the Think instance (new field: `Traces []ThinkTrace`)
- After the session completes, the trace is available for Context
- After Context processes it, the trace is cleared
- `History []ThoughtSection` continues to hold only surface thoughts (the present view)

### Context changes
- New method: episodic processing entry point
- Receives both streams (exchange batch + think traces for the same period)
- Calls provider for extraction
- Updates graph (entities, edges, memory nodes)
- Clears consumed data from Exchange and Think

### Trigger
- Same location as current compaction (`Exchange.Realize`, threshold check)
- Instead of calling `mem.Compress()`, calls Context's episodic processing
- Context needs to be reachable from the exchange (already on the relation's Realities map)
