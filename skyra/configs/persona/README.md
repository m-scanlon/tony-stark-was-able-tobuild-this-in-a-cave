# Skyra Persona Config

## Files

- `soul.md` — enduring values and long-horizon behavioral identity. **Design document only — not loaded at runtime.** Informs fine-tuning data selection and model behavior decisions over time.
- `personality.md` — actionable behavioral rules. The thin runtime layer: self-reference style, humor guardrails, frustration handling, response style.
- `tool-contract.md` — the structured tool call contract Skyra uses when delegating work to Pi.

## Design Intent

The persona is split intentionally:

**`soul.md`** defines what Skyra should be over time. It is not a prompt — it is a specification. As the system accumulates real interaction data, soul.md becomes the basis for fine-tuning examples and behavioral alignment decisions. It is never loaded into the context window.

**`personality.md`** is the runtime layer. It contains only rules that actively change output behavior. Kept as thin as possible to preserve context budget on smaller models (Pi front-door model has 4k-8k context). No descriptive or observer-facing language.

## What Is Not Here

Runtime system prompts have been removed. The persona is not composed at runtime from these files yet — that pipeline is planned. For now, personality.md is the source of truth for what gets loaded.

## Mutability

These files are expected to evolve. soul.md changes rarely (enduring values). personality.md changes when specific behavioral rules need adjustment based on real usage.
