# Kernel Cleanup

## What This Is

`DerivePresent` is the world kernel. It derives the present moment for a being — ingress, context, medium, response, routing. That's one flow and it stays one flow. The problem is not that the method is too big. The problem is that it does structured work inline where it should call into organized things.

Two areas to extract: the response cycle and the routing rules.

---

## Response Cycle

The block from medium response through to valid outbound relations is currently inline in `DerivePresent` (~60 lines, `world.go:106-166`). It parses the response, checks for format errors, builds retry feedback, re-fires the medium, re-parses, and returns valid relations plus anything it had to drop.

This is a self-contained cycle. The kernel fires the medium and gets a string back. Then it needs valid relations. Everything between those two moments — parsing, validation, retry feedback construction, the second medium call, re-parsing — is the response cycle.

Extract it so the kernel calls something like:

```go
valid := w.resolveResponse(response, name, threadID, currentPeer, present, medium)
```

The kernel doesn't see the retry loop. It gets back the relations that survived. Drops are still logged inside the cycle.

### What moves

- `parseResponse` stays as-is but becomes internal to the cycle
- The retry feedback construction (building the feedback string, re-firing the medium, re-parsing) moves into the cycle
- The kernel's dispatch loop receives only the final valid set

---

## Routing Rules

After the response cycle produces valid relations, the kernel loops through them and applies four checks before dispatching each one:

1. **Self-reference drop** — a being cannot target itself
2. **Auto-close on return** — if the being is returning to whoever called it, close the detour exchange
3. **~ref departure close** — if the being is addressing a new peer with a ~ref, close the current exchange
4. **Parent block** — if the target is the being's parent in an active exchange the being didn't open, block the message

These are routing rules. They share a shape: take an outbound relation, look at thread state, decide whether to route, drop, or modify exchange state. Right now they're four if-blocks with interleaved concerns (some mutate the thread, some log drops, some set flags that later blocks check).

Extract them into a named structure. Each rule takes the outbound relation and the thread, and returns a decision: route, drop (with reason), or route-and-close (with which exchange to close). The kernel applies the rules in order, acts on the decision, and dispatches or drops.

```go
type RouteDecision struct {
    Action  string // "route", "drop", "route-and-close"
    Reason  string // why, for logging
    CloseA  string // exchange to close (if action is route-and-close)
    CloseB  string
}
```

The rules become a slice the kernel walks:

```go
rules := []RouteRule{
    selfReferenceRule,
    autoCloseReturnRule,
    refDepartureRule,
    parentBlockRule,
}
```

Adding a new routing rule means adding a function to the slice. Reading the routing logic means reading each rule in isolation. The kernel's dispatch loop becomes: for each outbound relation, apply rules, act on decision, dispatch or drop.

### What moves

- The four inline checks move into named rule functions
- Thread mutation (close exchange) moves into the kernel's action on the decision, not into the rule itself — rules decide, the kernel acts
- Drop logging stays in the kernel, driven by the decision's reason

---

## What stays

`DerivePresent` stays as the single kernel method. The flow doesn't change:

1. Ingress — resolve target, thread, validate
2. Append to exchange
3. Derive present (lowercase helper, unchanged)
4. Fire medium
5. Resolve response (new — encapsulates parse/retry)
6. Apply routing rules and dispatch (new — structured instead of inline)

Same flow. Same method. Organized internals.

---

## What this is not

Not a refactor of the threading model. The kernel still reaches into thread state for routing decisions — that coupling stays for now. This is about making the kernel's own logic maintainable, not about moving responsibilities across packages.