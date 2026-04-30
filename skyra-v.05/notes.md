# v.05 Notes

- macOS world — the machine is a world. Terminal, filesystem, network are devices inside it. The user sits behind the macOS world. The loop lives there.

## Skills

A skill is a direct route to a specific reality, bypassing default routing. It maps to how skills are already understood by every major model — the word is safe to use because it means the same thing here as it does everywhere else.

A skill is not code. It's a markdown file that describes how to call something. The being reads it and knows what to emit. It's documentation that becomes capability.

```
# skill: mike-laptop
device: terminal
call: ssh mike@laptop
```

```
# skill: fetch
device: shell
call: curl -s ~url
```

Skills live in the being's present as artifacts. A being with a skill file knows how to address that reality directly. A being without it goes through default routing.

No skill primitive in the runtime. Skills are just files — retained by the being, read into its present, used like any other knowledge. The LLM already knows what to do with them because the format matches its training data.

## Derive Present

Superseded by the parser-per-reality model below.

## Present Derivation — Parser Stack

Every reality is responsible for its own text parser. When a reality contributes data to a relation, it also provides a parser that knows how to render that data as text. Parsers register on the invariant's hashmap at registration time.

When the relation reaches the invariant, the invariant fires its parsers in order — first registered is top of the present, last registered is bottom. The invariant concatenates the output. That's the present.

### Rules

- Each reality owns its slice end-to-end: data + parser. No central present builder.
- Parsers are text parsers. Every reality describes itself as text.
- Order is explicit. Left (first registered) is top, right (last registered) is last.
- The invariant is dumb. It holds a hashmap of parsers, fires them in order, concatenates.
- Adding a new reality means adding a new parser. Nothing else changes.
- Adding a new invariant means registering parsers on it. Realities don't change.

### Invariant Types

- **LLM** — parsers produce the system prompt / context window. Same format for all LLM providers.
- **Claude Code** — minimal parser or none. Claude manages its own context. Just pass the impulse.
- **Shell** — parsers produce a command string.
- **API** — parsers produce structured payload.
- **Webapp** — parsers produce a request.

### The Matrix Problem

Every reality × every invariant type needs a parser. Thread renders as conversation history for an LLM but maybe as nothing for a shell. Economics renders as budget context for an LLM but maybe as an env var for a shell. The parser count is realities × invariant types.

### Current State (alpha)

For now, present derivation for LLM beings stays in the LLM's realize method. The parser stack is the target architecture but we're not building the full matrix until we have a second invariant type (shell) that forces the split. When the deployment pipeline lands and we need shell + LLM rendering the same relation differently, the parser stack becomes necessary and the shape will be concrete.

### What This Replaces

The old model had derive present as a single layer between the being and the device. This replaces it with distributed ownership — each reality knows how to present itself, and the invariant is just the ordered stack where those parsers fire.
