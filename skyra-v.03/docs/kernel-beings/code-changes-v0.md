# Code Changes v0

## Status

In progress. Not yet locked.

## Purpose

Tracks the code changes needed to support kernel beings and the object makeover.

---

## 1. Callable Language On Relationships

### What

Add a `callableLanguage` field to both `ExchangeStack` and `ExternalDispatch`.

Expose it through the `RelationshipChannel` interface as `CallableLanguage() string`.

### Why

Present derivation currently injects peers by identity and purpose only.

The being knows who its peers are but not the specific syntax to speak to them.

Callable language is the expression syntax a being uses when addressing a specific peer.

It lives on the relationship — not on the being.

### Changes

**`RelationshipChannel` interface** — add:
```go
CallableLanguage() string
```

**`ExchangeStack`** — add field and method:
```go
callableLanguage string

func (c *ExchangeStack) CallableLanguage() string
```

**`ExternalDispatch`** — add field and method:
```go
callableLanguage string

func (d *ExternalDispatch) CallableLanguage() string
```

**`derivePresent`** — in the peer loop, after purpose, inject callable language:
```go
builder.WriteString("\n")
builder.WriteString(peer.Name())
builder.WriteString(" ")
builder.WriteString(peer.CallableLanguage())
```

### Invariants

- Callable language is mutable. It lives on the relationship and can be updated.
- Base language lives on the being. It is set at creation and is not mutable.
- These are two separate things and must never be collapsed.
- Cognitive beings earn callable language through relating. It starts empty when the relationship is seeded.
- Non-cognitive primitive beings are born with callable language. It is seeded at genome time, not earned.
- Both `ExchangeStack` and `ExternalDispatch` carry callable language. Non-cognitive beings are not exempt — they need it more than cognitive beings because they cannot reason their way to the right syntax.
- Everything into the kernel is a string. The kernel is the only place strings become objects.
- Every being is the same type: `Being`. One object. No subclasses.

---

## 2. Base Language On Being

### What

Add a `baseLanguage` field to `Being`.

Set at creation. Never mutated.

### Why

Base language is what makes first contact possible.

It is intrinsic to the being — not relationship-owned.

It needs a home in the struct.

### Changes

**`Being`** — add field:
```go
BaseLanguage string
```

**`NewBeing`** — accept base language as a parameter and set it.

### Invariants

- Base language is set at creation and locked. It does not change.
- It is not the same as callable language. Base language is the being's. Callable language is the relationship's.
- The genome defines base language for every being it seeds.

---

## 3. Non-Cognitive Being Hashmap

### What

Non-cognitive beings use `ExternalDispatch` as their peer channel — not `ExchangeStack`.

`ExternalDispatch` holds a single slot: `lastExpression`. Not a stack.

### Why

Non-cognitive beings are transducers. They receive, process, and emit.

They do not accumulate exchange history.

What matters is what just arrived — not what was said before.

### Invariants

- Cognitive beings → `ExchangeStack` → full stack of exchanges per peer.
- Non-cognitive beings → `ExternalDispatch` → single last expression per peer.
- The `cognitive` flag on `Being` determines which channel type is used when a peer relationship is seeded.
- This is already partially implemented. `ExternalDispatch` exists. The seeding logic needs to respect the cognitive flag when choosing which channel type to attach.

---

## 4. Identity Primitive

### What

A new file `src/domain/identity.go`.

Parses `~identity` from a protocol expression string and returns an `Identity` object.

### Why

Identity is a primitive. It needs its own parser. It cannot call itself — the method wraps it, parses the string, and returns the object.

### Changes

**New file: `src/domain/identity.go`**

```go
type Identity struct {
    Value string
}

func CreateIdentity(expression string) (Identity, error)
```

`CreateIdentity` reads `~identity <value>` from the expression string and returns a populated `Identity`.

Returns an error if `~identity` is missing or the value is empty.

### Invariants

- Identity is a plain string in v1. One field. One token.
- The parser owns only its slice of the expression — `~identity`. Nothing else.
- Returns an object. The caller assembles.

---

## 5. Purpose Primitive

### What

A new file `src/domain/purpose.go`.

Parses `~purpose` from a protocol expression string and returns a `Purpose` object.

### Why

Purpose is a primitive. It needs its own parser. It cannot call itself — the method wraps it, parses the string, and returns the object.

### Changes

**New file: `src/domain/purpose.go`**

```go
type Purpose struct {
    Value string
}

func CreatePurpose(expression string) (Purpose, error)
```

`CreatePurpose` reads `~purpose <value>` from the expression string and returns a populated `Purpose`.

Returns an error if `~purpose` is missing or the value is empty.

### Invariants

- Purpose is a plain string in v1. One field. One token.
- This is the declared purpose — what the being is for at creation. Not the realized purpose.
- The parser owns only its slice of the expression — `~purpose`. Nothing else.
- Returns an object. The caller assembles.

---

## 6. Nature Primitive

### What

A new file `src/domain/nature.go`.

Parses `~identity` and `~purpose` from a protocol expression string, calls `CreateIdentity` and `CreatePurpose` internally, and returns a `Nature` object.

### Why

Nature is the assembler of identity and purpose. It does not parse identity or purpose itself — it delegates to their own primitives and assembles the result.

### Changes

**New file: `src/domain/nature.go`**

```go
func CreateNature(expression string) (Nature, error)
```

`CreateNature` reads the full expression string, calls `CreateIdentity` and `CreatePurpose`, and returns a populated `Nature`.

Returns an error if either identity or purpose fails.

### Invariants

- Nature owns no parsing logic of its own. It delegates to `CreateIdentity` and `CreatePurpose`.
- Nature is locked at creation. The object returned is the final shape.
- Both identity and purpose are required. Missing either is an error.
- Returns an object. The caller assembles.

---

## 7. Language Primitive

### What

A new file `src/domain/language.go`.

Parses `~expression` from a protocol expression string and returns a `Language` object.

### Why

Language is a primitive. It holds the expression syntax peers use to speak to a being on a relationship. It needs its own parser.

### Changes

**New file: `src/domain/language.go`**

```go
type Language struct {
    Value string
}

func CreateLanguage(expression string) (Language, error)
```

`CreateLanguage` reads `~expression <value>` from the expression string and returns a populated `Language`.

Returns an error if `~expression` is missing or the value is empty.

### Invariants

- Language is a plain string in v1. One field. One token.
- The parser owns only its slice of the expression — `~expression`. Nothing else.
- For primitive beings language is seeded at birth. For cognitive beings it is earned through relating.
- Returns an object. The caller assembles.

---

## 8. Being Primitive

### What

A new file `src/domain/create_being.go`.

Parses the full being protocol string, calls `CreateNature` and `CreateLanguage` internally, assembles a `Being`, and registers it with the kernel.

### Why

The being primitive is the assembler of everything. One protocol string in. One live being registered in the kernel out.

### Changes

**New file: `src/domain/create_being.go`**

```go
func CreateBeing(expression string, state *KernelState) (*Being, error)
```

`CreateBeing` reads the full expression string and:

1. Parses `~name` — the being's name
2. Calls `CreateNature` — delegates identity and purpose parsing
3. Calls `CreateLanguage` — delegates language parsing
4. Parses `~cognitive` — true or false
5. Calls `NewBeing` with the assembled parts
6. Calls `state.InsertBeing` to register with the kernel

Returns the live `Being` and an error if any step fails.

### Invariants

- `CreateBeing` owns no parsing logic for identity, purpose, or language. It delegates entirely.
- All fields are required. Missing any token is an error.
- The kernel is the only place a being is registered. `CreateBeing` is the one path.
- Everything into the kernel is a string. `CreateBeing` is the bridge — string in, object registered.
- One being type. `CreateBeing` always produces a `Being`. No subclasses.
