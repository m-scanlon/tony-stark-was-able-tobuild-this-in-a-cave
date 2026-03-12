# Skyra Base Skills — v1

## Overview

Skyra's skills are her own internal capabilities — tools she calls herself to reason, gather information, and take action. They are not hardware-derived. They are registered under `agent:skyra` in Redis at brain boot by a `skyra_bootstrap` process.

Same registry, same format, different source.

## Registration

At brain boot, `skyra_bootstrap` writes Skyra's skill set to Redis:

```
key: agent:skyra
value: {
  status: "active",
  type: "system",
  capabilities: [ ...skills ]
}
```

Skills are available immediately. No shard registration required.

## v1 Base Skills

### iMessage

Send and read messages on behalf of the user.

```
tools:
  send_message   — send an iMessage to a contact. args: [contact, message]
  read_messages  — read recent messages from a contact. args: [contact, limit]
  search_messages — search message history. args: [query]
```

### Google Calendar

Create, read, update, and delete calendar events via Google Calendar API.

```
tools:
  create_event  — create a new calendar event. args: [title, start, end, description]
  read_events   — read upcoming events. args: [from, to]
  update_event  — update an existing event. args: [event_id, fields]
  delete_event  — delete an event. args: [event_id]
  search_events — search calendar events. args: [query]
```

### Apple Calendar

Sync Apple Calendar locally. Skyra manages the sync herself without user intervention.

```
tools:
  sync          — bidirectional sync with Apple Calendar
  create_event  — create a local calendar event. args: [title, start, end]
  read_events   — read upcoming local events. args: [from, to]
  update_event  — update a local event. args: [event_id, fields]
  delete_event  — delete a local event. args: [event_id]
```

### Google Search

Search the web and fetch content from URLs.

```
tools:
  search    — perform a web search, return ranked results. args: [query, limit]
  read_url  — fetch and parse content from a URL. args: [url]
```

### Code Execution

Execute code and shell commands in a sandboxed environment.

```
tools:
  run        — execute code. args: [language, code]
  run_shell  — run a shell command. args: [command]
```

## Router Handling

Skyra's skills are NOT hardcoded in the router's `case "skyra"` switch. Only `reply` and `delegate` are hardcoded system primitives. All other Skyra tools fall through to Redis:

```
case "skyra":
    switch command.tool:

        case "reply": ...
        case "delegate": ...

        default:
            // read from Redis — Skyra's dynamic skills
            agent = redis.get("agent:skyra")
            tool = agent.capabilities.find(command.tool)
            // dispatch as normal
```

This means new skills can be added to Skyra at any time by updating her Redis entry. No router code changes needed.

## External Skill Trust Policy

**No external skill is ever registered without explicit user approval. No exceptions.**

This applies to:
- OpenClaw / ClawHub community skills
- Any skill Skyra discovers via web search
- Any skill Skyra builds herself via code execution

The registry is trusted space. Everything trying to enter it is untrusted until the user approves it.

### Approval Flow

```
Skyra identifies a skill that fills a capability gap
  → NEVER registers automatically
  → surfaces to user via reply:
      "I found a skill that can do X.
       Here's exactly what it does: [full skill content]
       Approve to register?"
  → user reviews full content
  → user approves → registered in Redis → active
  → user denies → discarded, never registered
```

### Audit Trail

Every approved skill is commit-logged to the user's object store — what was approved, when, and by whom. If a skill ever behaves unexpectedly, the trail is there. Skills can be revoked at any time by removing them from Redis and logging the revocation.

### OpenClaw Compatibility

Skyra is compatible with the OpenClaw `SKILL.md` format. Skills from ClawHub can be pulled and parsed, but they follow the same approval flow as any other external skill. The 13,000+ community skills are a resource, not a trusted library.

## Skill Acquisition Engine

The five base skills are not just five capabilities. They are a bootstrapping surface. **Google Search + Code Execution together form a skill acquisition engine** — Skyra can discover and build virtually any capability she encounters a need for.

### Example — Skill Not Yet Registered

```
user: "send a message to my friend Coney"

Skyra: I need to send a message
  → checks Redis agent:skyra → no messaging skill found
  → I don't have this capability yet
  → but I have Google Search and Code Execution

Skyra: octos search "how to send iMessage programmatically on macOS"
  → returns: osascript / AppleScript can send iMessages natively

Skyra: octos run_shell "osascript -e 'send message...'"
  → works

Skyra: I just built a messaging capability
  → writes iMessage skill to Redis under agent:skyra
  → registered permanently for future use
  → sends the message
  → octos reply "Message sent to Coney"
```

Next time a user asks to send a message, `iMessage.send_message` is already registered. Skyra never has to rediscover it.

### What This Means

- v1 does not need to pre-register every possible tool
- Skyra builds her own skill set organically as she encounters needs
- The registry grows with her use
- The base skills are minimal by design — the real capability is the ability to grow beyond them
- Any skill Skyra can find instructions for online and execute via code, she can acquire and register permanently

## v2 Skills

### System UI Control

Skyra needs to interact with macOS GUI dialogs and apps that don't have APIs — permission prompts, app windows, OS-level interactions. Requires screen watching, dialog reading, and mouse/keyboard control via the macOS Accessibility API.

```
tools:
  watch_screen   — detect and read UI elements on screen
  click          — click at screen coordinates or named UI element
  type           — type text into focused UI element
```

Deferred to v2. iMessage Automation permission is already granted on the host Mac for v1 so the immediate need is covered.

## Open Questions

- What sandbox environment does code execution run in? Docker, WASM, or a restricted shell?
- Does iMessage require a Mac-local bridge process or can it be called via API?
- Apple Calendar sync frequency — on-demand only or background periodic sync?
- Auth token storage for Google APIs — where do credentials live and how are they rotated?
