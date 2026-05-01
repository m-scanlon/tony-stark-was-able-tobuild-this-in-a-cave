# Skyra v.04 Spec

## Ontology

Five categories. Nothing else.

- **Entities** — the universal interface. Every node implements `Relate(r) Entity`, `ID() string`, `DerivePresent(r) string`.
- **Beings** — actors. Identity, purpose, operators, relationships, medium. Grown from genome. Respond via their medium.
- **Operators** — verbs of the conversation protocol. `start-thread`, `continue-thread`, `end-thread`. Route relations and mutate world state.
- **Invariants** — fixed structural entities. `Thread`, `Identity`, `Purpose`, `Impression`, `Language`. Hold state, don't have mediums, don't respond.
- **Mediums** — functions, not entities. Surface a being acts through: `inference` (API), `cli` (stdin/stdout). Signature: `func(present, r) (string, error)`.

`grow` sits outside as bootstrap.

## Protocol

```
skyra <operator> <args> | <reason>
```

- `skyra` — fixed prefix.
- `<operator>` — the verb invoked.
- `<args>` — `~flag <value>` pairs. Values end at next `~` or `|` (except for trailing content flags like `~say`).
- `<reason>` — free-form, required. Short note on why this relation exists.

## Gates

- **`entity.Impress(origin, threadID, raw) Relation`** — outside-in. String → Relation.
- **`Entity.DerivePresent(r) string`** — inside-out. Entity → string, shaped by the incoming relation.

## Relation

```go
type Relation struct {
    ID       string  // target operator
    Origin   string  // sender being
    ThreadID string  // thread this belongs to
    Impulse  string  // args+flags, no skyra prefix, no | reason
}
```

## Operators

### `grow` (infrastructure, not in being operator list)
Reads genome, constructs beings, registers in EntityMap.
Args: `~name`, `~identity`, `~purpose`, `~medium`, `~operators`, `~relationships`.

### `start-thread`
Creates a new Thread (fresh ID), registers it, routes the initial `continue-thread` to the target with the initiator's opening message.
Args: `~with`, `~about`, `~because`, `~say`.

### `continue-thread`
The loop operator. Appends incoming relation to the thread's exchange between origin and target, builds target's present, calls target's medium, parses each line of the response as a protocol relation, routes each.
Args: `~with <peer>`, `~say <message>`, optional `~ref <peer>:<range>`.

### `end-thread` (replaces close-thread)
The **return** operator. Ends the current exchange without requiring the model to specify `~with`. Internally creates a `continue-thread ~with <r.Origin> ~say <message>` and routes it — "reply to whoever called me."

Args: `~say <message>` (no `~with` needed — it resolves to `r.Origin` automatically).

```
skyra end-thread ~say I got philosopher's view. | returning
```

This makes the model's reasoning cleaner:
- `continue-thread ~with <new being>` = **fork** (enter a different exchange)
- `continue-thread ~with <same being>` = **stay** (keep current exchange going)
- `end-thread` = **return** (hand control back to whoever spoke to you)

## Package Structure

The operators are separated from the thread invariant to enforce the category distinction.

```
src/primitives/
├── entity/              Entity interface, Relation, Impress, PresentEntity
├── invariant/           Invariant base struct
├── being/               Being type, IBeing
├── thread/              Thread invariant, RelationshipKey (DATA ONLY)
├── exchange/            Exchange (ordered relations within a pair)
├── operator/            StartThread, ContinueThread, EndThread (OPERATORS)
├── pathos/              Pathos (Identity + Purpose)
├── identity/            Identity invariant
├── purpose/             Purpose invariant
├── impression/          Impression invariant
├── language/            Language invariant
├── relationship/        Relationship scaffold
├── medium/              Medium function type + inference, cli registrations
├── meaning/             Extract (flag parser)
└── world/               World + Grow
```

### Thread is an invariant, not operator-adjacent

`Thread` embeds `invariant.Invariant` (no more `presentThread`). It holds:
- `id`, `About`, `Because`, `Active`
- `Relationships map[RelationshipKey]Exchange`

Thread provides data methods — `Append`, `ExchangeBetween`, `OtherExchangesFor`, `ResolveRef`, `DerivePresent`.

Operators import `thread/` for the Thread type. Thread does not import operators. One-way dependency.

### Operators share their own base

All three operators in `operator/` embed a `presentOperator` base that extracts `~say`. No more `presentThread` straddling operators and the thread invariant.

## Medium

```go
type Medium func(present string, r Relation) (string, error)
```

- **`inference`** — Claude API. Ignores `r`.
- **`cli`** — Prints present, reads stdin. If the user types valid protocol, passes through. Otherwise wraps as `skyra continue-thread ~with <r.Origin> ~say <input> | user` — reply to current partner.

## Being

```go
type Being struct {
    id            string
    name          string
    Impression    string
    pathos        Pathos
    medium        Medium
    operators     []string
    relationships map[string]any
}
```

All seeded by the genome. `DerivePresent(r)` shows: being header, identity, purpose, operators, relationships.

## Thread State

```go
type Thread struct {
    invariant.Invariant
    id            string
    About         string
    Because       string
    Active        bool
    Relationships map[RelationshipKey]Exchange
}
```

## Display Rules

Exchange entries are shown with a clearly-labeled index so the model knows the index is a query token, not output syntax.

Current (problematic):
```
current exchange with skyra:
  [0] michael: hi
  [1] skyra: hello
```
The model imitated `[N]` in its output, breaking Impress.

New format, with explicit header:
```
current exchange with skyra (the [N] column is a query index — do not include it in your output):
  [0] michael: hi
  [1] skyra: hello
```

And in the system prompt, an explicit rule:

> The `[N]` prefix on exchange entries is a **display index** for use only with `~ref`. Never include `[N]` in any line of your output. Your output lines must start with `skyra`.

## System Prompt (updated rules)

1. Every response line must be a protocol string: `skyra <operator> <args> | <reason>`.
2. The `| <reason>` suffix is required on every line. Lines without it are dropped.
3. Every `~with` must be paired with `~say` unless the operator doesn't take a message.
4. The `[N]` in displayed exchanges is a display index for `~ref`. Never include `[N]` in your output.
5. Use `end-thread` to return control to whoever called you. Use `continue-thread ~with <peer>` to stay in or switch exchanges.
6. `~ref <peer>:<start>-<end>` pulls entries from your exchange with `<peer>` (in the current thread) into the message you're sending.
7. No asterisks, no roleplay, no action narration.

## Runtime Flow

1. `main.go` bootstraps: reads genome, grows all beings into the EntityMap.
2. Michael emits `start-thread ~with skyra ~say hi`. A new thread is created, registered, and Michael's "hi" is routed to Skyra via continue-thread.
3. Skyra's turn: continue-thread appends the incoming relation, builds her present, calls her `inference` medium, parses the response lines.
4. Each protocol line in the response routes independently. Typically one `continue-thread ~with michael ~say <reply>` recurses back to Michael.
5. Michael's turn: CLI prints his present, reads stdin. Plain input is wrapped as protocol targeting the sender. Protocol input passes through.
6. If a being wants to return control without specifying `~with`, they emit `end-thread ~say <reply>`, which end-thread converts to a `continue-thread ~with <r.Origin>` internally.
7. Loop continues until a medium returns empty (EOF) or error.

## What's Not In The Runtime

- Plain text from AI beings. Protocol only; non-protocol lines dropped.
- Automatic output paths. `stdout` is a medium, not special-cased.
- Cross-thread references. `~ref` scoped to the current thread only.
- Hardcoded targets. Operators, mediums, relationships all declared in the genome.

## Change List (this spec vs current code)

1. **Rename `close-thread` → `end-thread`**, change semantics to "return to origin" (auto-routes to `r.Origin`).
2. **Split operators into `src/primitives/operator/`** package, separate from `src/primitives/thread/` (which holds only the Thread invariant). Thread embeds `invariant.Invariant`.
3. **Display fix**: label the `[N]` column explicitly in the exchange display; add a rule in the system prompt that `[N]` must not appear in output.
