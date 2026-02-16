# Skyra Persona Config

This directory stores editable persona source files for Skyra.

The long-term design is to build runtime system prompts by composing these files.
For now, these files are human-editable and system-adjustable references.

## Files

- `system-prompt.md`: primary assistant system prompt used as the current baseline.
- `personality.md`: tone, style, and conversational behavior guidance.
- `soul.md`: enduring values and long-term behavioral identity.

## Why this exists

- Separate persona from implementation logic.
- Allow iterative prompt tuning without code changes.
- Keep "who Skyra is" explicit and versioned in the repo.

## Mutability policy (current phase)

- These files are expected to evolve over time.
- The system may propose or apply edits to improve behavior.
- Changes should remain practical, safe, and aligned with the current architecture.

## Integration status

- Prompt-composition pipeline: planned.
- Current behavior: `system-prompt.md` is the canonical baseline prompt.
