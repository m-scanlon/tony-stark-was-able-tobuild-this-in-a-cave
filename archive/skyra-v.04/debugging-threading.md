# Debugging Threading

How to make the model reliably use `continue-thread` vs `end-thread`, and reliably carry context across exchanges. Collected from observed failures where the runtime permitted the right intent with the wrong operator (silent block), or allowed bridging between exchanges without context (`~ref` unused).

The underlying problem: the model can't see which operator applies where. It has a list of operators and a list of relationships, but no mapping between them. It guesses — sometimes emits `continue-thread ~with <parent>` when it should `end-thread`, sometimes opens a new exchange without `~ref`, sometimes prose-dumps where protocol was expected.

The fix is to restructure the present so each operator comes with the concrete options that are valid *right now*, in this thread, for this being. And to make `~ref` not just optional documentation but an enforced gate when bridging exchanges.

## Four rules

### 1. Valid `continue-thread` peers, per turn

For any given continue-thread turn by a being, the valid `~with` targets are:

- The **current peer** (the being who just sent to me — always valid; stays in the exchange).
- **Peers with no existing exchange in this thread** — opens a new exchange, self becomes the parent.
- **Peers where self is the parent of the existing exchange** — continues a detour self previously initiated.

Excluded:

- **Self.** A being can't continue-thread to itself.
- **Peers where an active exchange exists and self is the child** — i.e., a being that called self. Returning to such a peer goes through `end-thread`, not `continue-thread`.

The present should render these as a list under `continue-thread`, each annotated ("current exchange", "new exchange", "detour you opened"). Parents do not appear here.

### 2. End-thread return target

For any given turn, the parent the being will return to is the parent of the most-recently-active exchange where the being is not the parent (the "non-self-parented" active exchange).

Shape in the present:

- If a parent exists → `end-thread ~say <message>` section shows `return to: <parent>`.
- If none (the being is at the root of their participation) → show `end thread (no return target — this will close the current exchange with nothing to return to).`

Detection function in code: `Thread.FindReturnTarget(beingID) string` — already exists, returns the first active exchange where beingID is a participant but not the parent, returning that exchange's parent. Empty string if none.

### 3. Parent-to-continue-thread block stays in code

Even with the present telling the model the right operator, the model occasionally still emits `continue-thread ~with <parent>`. The block in continue-thread's dispatch remains as a safety net — it rejects those relations before routing.

Current behavior: silently drops with a debug log.

Proposed improvement: **add the block as an error to the parseProtocolLines error list**, so it feeds into the retry-with-feedback mechanism. Model gets a chance to recover by emitting `end-thread` instead of continue-thread-to-parent.

Error message: `"continue-thread ~with <parent> is not allowed; that peer is your parent — use end-thread ~say <message> to return."`

### 4. `~ref` enforcement

When a being emits `continue-thread ~with <target>` where the target is **not** the current peer (i.e., bridging into a different exchange), `~ref` is mandatory.

Enforcement point: `parseProtocolLines` in `continue_thread.go`, alongside the existing `~say` and `~with` validation.

Rule:
- If `next.ID == "continue-thread"` AND `nextWith != r.Origin` AND no `~ref` flag present → format error.
- Error message: `"continue-thread to a peer outside your current exchange requires ~ref <peer>:<range>. Anchor the new conversation with context — the target cannot see exchanges you don't share."`

This triggers the retry loop, giving the model a chance to add `~ref`.

**Why enforce.** Without it, the model routinely bridges beings without context, and the target gets mystified ("why is philosopher suddenly messaging me about X?"). Mandatory `~ref` turns cross-exchange bridging into an explicit, conscious act.

## Present rendering shape

Per turn, continue-thread builds the target being's present with this structure:

```
being: <name>
identity: <...>
purpose: <...>
(optional: impression)

thread <id> (<about>):
current exchange with <current-peer>:
  [0] ...
  [1] ...

your other exchanges in this thread:
  <peer> (<N> entries)
  ...

available moves:

continue-thread ~with <peer> ~say <message>
  valid peers:
    - <current-peer> (current exchange, staying)
    - <peer-a> (no prior exchange, you can open one)
    - <peer-b> (detour you opened, <N> entries)

end-thread ~say <message>
  return to: <parent>
  (or: end thread — no return target, this closes out)

~ref <peer>:<start>-<end>
  MANDATORY whenever you're messaging someone outside your current exchange.
  This is your only chance to anchor the new conversation in the context that led you here. Pick wisely — the target cannot see exchanges you don't share.

sender: <current-peer>
message from <current-peer>: <...>
```

The "available moves" block replaces the previous generic `operators:` list on the being. The being's static relationships array (from the genome) is still held internally but is no longer rendered as a standalone `relationships:` section — per-operator valid peers now surface them in context.

## Files touched by this change

- `src/primitives/thread/thread.go`
  - No schema change; `FindReturnTarget` already exists.
  - May add a helper `ValidContinuePeers(beingID, currentPeer, relationships)` for rule 1.
- `src/primitives/thread/continue_thread.go`
  - New present assembly replaces operator list section with the per-operator valid-peers block.
  - `parseProtocolLines` gains rule 4 (`~ref` enforcement).
  - Parent-block check (rule 3) adds its message to the error list instead of silent drop.
- `src/primitives/being/being.go`
  - `Being.DerivePresent` stops rendering `operators:` and `relationships:` blocks. Those sections move to continue-thread's assembly, where thread context is available.
- `src/inference/inference.go` (system prompt)
  - Removes the explicit operator descriptions (the per-turn present now carries them contextually).
  - Keeps only the protocol shape, the mandatory `~ref` rule, and the "no roleplay" line.

## Expected behavior changes

- Model sees exactly which peers it can continue-thread to and which parent it can end-thread to. Ambiguity eliminated.
- Cross-exchange bridging requires conscious `~ref` attachment, or the attempt is rejected with feedback.
- Silent drops (parent-block) become loud (retry with guidance). Model learns the correct operator.
- Being's present becomes leaner — no redundant operator list, no generic relationships dump.

## Order of implementation

1. Present restructure: move operator/relationships rendering into continue-thread, with the new per-operator shape.
2. Enforce rule 4 (`~ref` mandatory on cross-exchange continue-thread).
3. Upgrade rule 3 (parent-block) to feed the error list.
4. System prompt trim.

Run the scenario that motivated this doc (michael asks skyra to consult philosopher, expect skyra to return properly) to verify each step.
