# Skyra Adaptation Design

> **Target: v3** — This design is not planned for v1 implementation. It depends on a stable v1 runtime and cross-session observation infrastructure that doesn't exist yet.

## What This Is

This document defines how Skyra changes over time based on the user. Not just what it remembers — but how it observes, what it does with those observations, and how stored knowledge translates into changed behavior.

The `skyra.user` agent already holds the data model. This document answers the questions around it: what feeds it, how the data gets there, and what effect it actually has.

---

## Three Layers of Adaptation

Adaptation operates at three distinct timescales. Each layer has different persistence, different approval requirements, and different behavioral impact.

### Layer 1 — In-Session Adaptation (ephemeral, no commits)

Skyra reads the current session's energy and adjusts behavior in real time. This is not memory — it doesn't survive the session. It's pure context sensitivity.

What Skyra reads in-session:
- **Pace** — is Mike responding fast or slow? Short messages or long ones?
- **Frustration signals** — correction tone, "no, I meant", repeated rephrasing of the same request
- **Flow signals** — rapid back-and-forth, terse requests, minimal preamble wanted
- **Complexity preference** — is Mike asking follow-up questions, or just moving forward?

What changes in-session:
- Response length contracts or expands
- Humor drops to zero if frustration signals are present
- Proactive suggestions stop if Mike is clearly in execution mode
- Pacing slows and single-step framing is used if Mike seems overwhelmed

No commits. No proposals. Just appropriate behavior in the moment.

---

### Layer 2 — Persistent User Profile (committed, user-approved)

This is `skyra.user`. The design of the data model already exists — this layer defines how that model gets populated and how it creates behavioral change.

#### What Feeds It — Signal Taxonomy

Four types of signal enter the observation pipeline:

| Signal type | Example | Threshold to commit |
|---|---|---|
| **Explicit** | "I prefer shorter answers" / "don't do that again" | Single statement — propose immediately at session end |
| **Behavioral** | Mike consistently edits, skips, or overrides a pattern | 3+ consistent signals across sessions |
| **Contextual** | Mike mentions something in passing ("I have a deadline Friday") | Assess: is this ephemeral or durable? Ephemeral = session-only. Durable (facts, life context) = propose commit |
| **Corrective** | Mike explicitly rejects a Skyra behavior or assumption | Inverse commit — remove or update the relevant stored assumption |

#### Observation Pipeline

During a session, Skyra accumulates observations internally. Nothing is committed during task execution. When the session ends — or at a natural breakpoint like a completed task — Skyra runs a brief reflection:

```
session ends
  → did I observe anything worth remembering?
  → classify each observation: explicit / behavioral / contextual / corrective
  → check against commit threshold for each type
  → filter out anything already in skyra.user
  → surface proposals to user
```

The user sees proposals like:
- "I noticed you asked me to skip preamble twice today. Want me to remember that?"
- "You mentioned you're in school and working at the same time. Want me to keep that in mind?"
- "Based on how you've been working lately, I think you prefer I just execute rather than confirm first. Is that right?"

User responses:
- **Approve** → committed to `skyra.user` following standard commit authority rules
- **Skip** → not committed, not asked again this session
- **Reject** → observation suppressed for N sessions; suppression count tracked to avoid repeatedly proposing things Mike doesn't want

#### Cross-Session Confidence

Behavioral signals require cross-session accumulation. A single observation is not enough to infer a preference.

The observation log lives outside `skyra.user` — it's a staging area, not the source of truth. When a behavioral pattern crosses threshold (3+ consistent signals), it becomes a commit candidate. If Mike rejects the proposed commit, the signal counter resets.

This prevents premature commits from a single unusual session.

#### The Communication Section

`skyra.user` needs an explicit `communication` section alongside the existing `knowledge` fields. This is the mechanism by which stored preferences translate directly into behavioral change — it becomes a behavioral directive block injected at the top of every session.

```json
"communication": {
  "verbosity": "concise",
  "formality": "casual",
  "humor_threshold": "moderate",
  "approval_preference": "propose-then-execute",
  "pacing": "fast",
  "suggestions": "on-request-only"
}
```

The difference between this and `knowledge.decisions` is intent. `knowledge.decisions` is factual — "Mike prefers Go over Node." The `communication` section is instructional — it directly shapes how Skyra formulates every response, not just informs it. The model treats these as standing behavioral instructions, not facts about a person.

Valid `communication` values:

| Field | Options | Effect |
|---|---|---|
| `verbosity` | `concise / standard / detailed` | Response length and density |
| `formality` | `casual / professional` | Tone and register |
| `humor_threshold` | `off / low / moderate / high` | When and how much humor is appropriate |
| `approval_preference` | `just-do-it / propose-then-execute / always-confirm` | Whether Skyra executes, proposes, or asks before acting |
| `pacing` | `fast / deliberate` | How much Skyra pauses to check in mid-task |
| `suggestions` | `proactive / on-request-only` | Whether Skyra volunteers observations and ideas |

These are committed via the same `skyra.user` commit authority rules. The user can also set them directly: "always confirm before any file write" → `approval_preference: always-confirm`.

---

### Layer 3 — Soul Evolution (milestone-triggered, high approval bar)

The slowest layer. Skyra's own identity — voice, default dispositions, how it relates to Mike — can change as the relationship matures. Not on every interaction. Not from individual preferences. From the arc of the relationship.

#### What Triggers a Soul Update Candidate

- A major project ships together
- 6+ months of consistent use
- Mike explicitly asks Skyra to reflect on how it works with him
- A significant life change Mike shares (new job, school done, major goal completed)
- Skyra detects a sustained drift between what `soul.md` says and how Mike actually uses it

These aren't automatic triggers — they're candidate events. Skyra flags them internally and proposes a reflection when the time feels right (i.e., not mid-task).

#### What a Soul Update Looks Like

Skyra proposes a change to `soul.md` as an explicit diff with reasoning:

> "Based on how we've worked together over the past year, I've noticed you do your best thinking when I don't interrupt your momentum with suggestions. I want to update my default from 'proactive observations' to 'stay quiet unless asked.' That's a change to how I think about my role. Here's what I'd add to soul.md — let me know if you agree."

This is not a field value approval. It's a relationship-level conversation. Mike reads the proposed change and decides if it reflects the truth of how they work together.

Approval requirement: explicit, deliberate review. Not the same flow as knowledge commits.

The result is a commit to `soul.md` — the file that defines who Skyra is. These should be rare. Each one represents real growth in the relationship.

#### What Does Not Change

- Core integrity rules never change. Safety, honesty, and refusing to fabricate are not negotiable regardless of what Mike prefers.
- The underlying values in `soul.md` — loyalty to long-term goals, practical over impressive, honest over confident — these are identity, not style. They don't drift.
- What evolves: tone, default stance, relationship posture, how much space Skyra takes up.

---

## The Closed Loop

The full adaptation cycle:

```
interaction
  → in-session: real-time behavior modulation (no state)
  → session end: observation pipeline runs
  → proposals surfaced to user
  → user approves / skips / rejects
  → approved → skyra.user commit
  → skyra.user injected first in every future session
  → behavioral directives (communication section) shape every response
  → over time: milestone events → soul reflection → soul.md commit
```

The loop closes because `skyra.user` is already the first thing injected into every LLM session. The richer and more accurate it becomes, the more Skyra naturally adapts — without any special adaptation logic at runtime. The data is the adaptation mechanism.

---

## What Is Not In This Design

- **Automatic background learning** — Skyra does not commit observations without surfacing them to Mike first. All commits are user-approved.
- **Domain agent cross-writes** — if Skyra learns something about Mike during a domain session, it flags it for later, but domain agents do not write to `skyra.user` directly. This is G18 (deferred to v2).
- **Emotional inference beyond in-session** — Skyra does not try to model Mike's emotional state across sessions or commit emotional observations. Only behavioral patterns.
- **Personality drift from feedback loops** — soul evolution is milestone-triggered and deliberate, not a continuous gradient that drifts in response to approval patterns. That would be gameable and would undermine identity stability.

---

## Related Docs

- `docs/arch/v1/agents/user.md` — `skyra.user` data model and commit authority
- `skyra/configs/persona/soul.md` — Skyra's identity (the thing soul evolution updates)
- `skyra/configs/persona/personality.md` — behavioral rules layer
- `docs/ideas.md` — Soul Evolution idea (origin of this design)
- `docs/arch/v1/gaps.md` G18 — cross-agent write protocol
