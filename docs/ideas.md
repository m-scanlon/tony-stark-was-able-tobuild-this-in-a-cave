# Ideas — Future Thinking

A place to capture ideas that aren't ready to design yet but are worth revisiting.

---

## Soul Evolution

**The idea:** Skyra's personality and preferences evolve over time through interaction rather than being statically defined once.

Two distinct layers:

**1. Learned user preferences (`preferences.md` or similar)**
As Skyra works with you across projects it gradually encodes what it learns about how you think, work, and make decisions — not as memory items buried in a vector DB, but as first-class documented traits. Things like:
- "Prefers TypeScript over JavaScript"
- "Wants plan approval before any file writes"
- "Likes concise answers, skip the preamble"

These would live in a separate document from `soul.md` — your preferences, not Skyra's identity. Periodically updated, not on every interaction. You'd review and approve changes before they're written, same as any other commit.

**2. Skyra's own soul evolution**
Skyra's `soul.md` — its own identity, values, and voice — should also be able to change over time. Not on every interaction, but gradually, as the relationship between you and the system deepens. What that cadence looks like and what triggers it is a conversation for another day.

**Why two separate documents matters:**
- `soul.md` = who Skyra is
- `preferences.md` = what it has learned about you
- Keeping them separate means Skyra's identity stays stable while your preferences grow independently

**Open questions:**
- What triggers a preference write? Threshold of repeated behavior? Explicit user signal?
- Who proposes the update — Skyra or the user?
- How does this interact with project-specific preferences vs global preferences?
- What does the approval flow look like for soul updates?

---

## More ideas to add here as they come up
