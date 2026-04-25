# Theory of Mind Language — v0

## The Problem

theory-of-mind currently has no callable language defined on its relationship with prefrontal. As a result, when it passes a signal inward it defaults to reporting signal quality ("language stable") rather than doing its actual job: identifying who is reaching out and attributing intent to the signal.

When michael says "hi", theory-of-mind produces:

```
skyra prefrontal hi | sensory: language stable for first contact
```

prefrontal receives this and sees a signal from sensory about language quality. It has no idea michael is waiting for a response. It bounces back to theory-of-mind because nothing in the signal tells it what to do.

## What Theory of Mind Is For

In neuroscience, theory of mind is the capacity to attribute mental states to other agents — to understand that someone else has intentions, desires, and beliefs, and that those can differ from your own. It sits at the perceptual boundary and answers: **who is this from, and what do they want?**

In Skyra's architecture, theory-of-mind sits between sensory and prefrontal for this exact reason. sensory receives raw signal and does not interpret it. theory-of-mind's job is to:

1. Identify the external being sending the signal
2. Model their intent — why they are reaching out, what they want
3. Pass that model inward to prefrontal with the originating being preserved as the source

theory-of-mind's output should answer: **someone wants Skyra's attention — here is who, and here is why.**

## The Fix

Define a callable language on theory-of-mind's relationship with prefrontal. This language specifies the expression format theory-of-mind must use when calling prefrontal, encoding the requirement to carry identity and intent.

The callable language should be set in the genome on the theory-of-mind → prefrontal relationship, not hardcoded in the being's nature. It is relationship-specific knowledge.

### Proposed callable language

```
skyra prefrontal <content> | <originating being>: <their intent>
```

Where:
- `<content>` is what the external being said or sent
- `<originating being>` is who sent it — michael, or any external being — not sensory
- `<intent>` is theory-of-mind's attribution of why they are reaching out

### Example

michael says "hi" → theory-of-mind produces:

```
skyra prefrontal hi | michael: initiating contact, expects a greeting
```

prefrontal now sees michael as the source, understands there is a person waiting, and can route toward premotor with that context intact.

## What Needs to Change

- Add `SetCallableLanguage` call for theory-of-mind's peer channel with prefrontal in the genome bootstrap, or expose callable language as a field in the genome expression syntax
- The genome currently has no mechanism to set callable language on a relationship — this may require a new `~language` field on relationship seed expressions, or a separate genome line type

## Progress

### Completed

- `extract.MeaningToEnd` added — extracts token values that may contain `~` flags, stopping only at `|`
- `bootstrap` in `main.go` fixed — now splits genome on newlines instead of the word `skyra`, so expression values can safely contain protocol strings
- `seedRelationships` in `world.go` extended — looks for `~language-<peername>` on relationship lines and calls `SetCallableLanguage` on the peer channel
- Expression syntax agreed: `skyra prefrontal ~from <being> ~message <content>`

### In Progress

- Callable language ownership is being moved to the callee side. Currently the expression syntax is set on theory-of-mind's genome line (`~language-prefrontal`). The right model is for prefrontal to define its own callable interface via a `~callable-language` field on its being definition. Any being that gets prefrontal as a peer should inherit the syntax automatically via `seedRelationships`.
- This requires adding `CallableLanguage string` to the `Being` struct and updating `seedRelationships` to propagate it when seeding peer channels.

### Still Open

- The source field problem remains — theory-of-mind needs to be guided (via the callable language showing in its present) to preserve the originating being as source rather than replacing it with `sensory`

---

## Why the Source Field Matters

The source field in the protocol (`| source: reason`) is the primary carrier of identity through the system. When theory-of-mind replaces `michael` with `sensory` as the source, michael disappears from prefrontal's present entirely. By the time prefrontal fires, it has no visibility into who originated the signal.

Preserving the originating being as source is not optional — it is how intent and identity propagate inward through the cognitive layers.
