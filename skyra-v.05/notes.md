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

There needs to be a present derivation layer between the being and the LLM call. It lives on the being layer but is invisible to the being — same pattern as physics on the surface world. It assembles the present from all the accumulated context — physics output, pathos, exchange history, skills, retained experience — into the thing the LLM actually sees. The being doesn't know it's there. It just gets a present that's already assembled by the time the device fires.

Devices have their own derive present too. Each device renders the present for its medium — the CLI renders for a human, the LLM device renders for a model (system prompt, context window, pathos), a screen renders for a display, an API renders for a webhook. Same relation goes in, different present comes out. The device knows how to render for its medium. That's the lens — the present is the same everywhere, only the glass changes. Each device has its own custom implementation.
