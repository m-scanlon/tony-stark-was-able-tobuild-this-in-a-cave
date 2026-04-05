# Skyra: Autonomous Programmable Cognition

Skyra is an experimental runtime for autonomous programmable cognition.

The project is moving toward a typed cognitive substrate built around:

- small fixed primitives
- extensible actors
- typed stimulus
- contract-bounded orchestration
- actor-local retained experience
- world interaction as a first-class runtime concern

This repository combines runtime code, architecture documentation, protocol work, product notes, and exploratory system-model design.

## Project Direction

The current direction is to treat Skyra less like a single agent and more like a runtime that can support many agent patterns.

The emerging model is:

- primitives provide the fixed substrate
- actors are the extensible unit
- base actors form the shipped standard library
- orchestrator actors coordinate other actors
- workflows emerge from typed actor composition rather than one hardcoded loop

This is still early-stage architecture work.

Some parts of the repo are active implementation, while others are preliminary design notes that are still evolving.

## Repository Map

- `skyra/`: runtime and application code
- `docs/`: architecture, protocol, memory, and design documentation
- `data-model-ideas/`: exploratory system-model notes
- `landing/`: landing page and related frontend work
- `prds/`: product requirement drafts
- `ideas/`: loose concept notes and future-facing explorations

## Current Focus

The current work is centered on stabilizing the core runtime shape:

- protocol and primitive boundaries
- typed stimulus and actor contracts
- actor-local memory and bounded recall
- world-facing interaction
- base actors and orchestrator actors
- device registration and onboarding

## Notes

- `docs/` contains both active canon and prelim documents; not every file should be treated as fully settled.
- `data-model-ideas/` should be read as exploratory unless a decision is clearly stated.
