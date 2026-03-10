# Skill: API Compatibility

An API skill wraps an external API as an executable skill. It is a first-class skill — same schema, same trust model, same crypto guarantees. The API surface it is allowed to hit is declared in the skill contract. The kernel enforces it.

---

## The Problem

External APIs require credentials, endpoint definitions, retry logic, and rate limits. None of this can live in the skill definition — the definition is content-addressed and potentially open. Credentials in the definition are a security leak. Hardcoded endpoints in the definition violate closed-for-modification.

The solution: the skill defines the contract (what it's allowed to do). The memory namespace holds the credentials (what it needs to do it). The kernel enforces the boundary between them.

---

## Schema Extension

API skills extend the base skill schema with an `api_contract` field.

```
skill {
  ...
  definition {
    ...
    api_contract: {
      allowed_endpoints: []endpoint_rule   // whitelist — kernel enforces
      auth_method:       api_key | oauth2 | bearer | none
      credential_ref:    string            // name of credential in memory namespace — never the key itself
      retry_policy:      retry_policy
      rate_limit:        rate_limit
      response_schema:   string            // natural language — what a valid response looks like
    }
  }

  memory {
    namespace:     string    // holds the actual credentials at runtime
    seed_memory:   bytes     // creator may ship credential templates — not actual keys
  }
}
```

### Endpoint Rule

```
endpoint_rule {
  method:   GET | POST | PUT | DELETE | PATCH
  host:     string     // exact host — no wildcards. e.g. "api.slack.com"
  path:     string     // prefix match. e.g. "/api/chat.postMessage"
}
```

The kernel checks every outbound API call against `allowed_endpoints` before executing. A call to an endpoint not on the whitelist is rejected. The skill cannot call anything it has not declared.

### Credential Reference

```
credential_ref: "slack_api_key"
```

This is a name — not a value. The kernel resolves it from the skill's memory namespace at execution time. The actual key is never in the skill definition. It is never in Redis. It lives in memory, written there by the user at provisioning time.

```
memory namespace: skill:send_slack_message
  → slack_api_key: "xoxb-..."    // written by user, never leaves memory
```

### Retry Policy

```
retry_policy {
  max_attempts:  int
  backoff:       exponential | linear | none
  retry_on:      []int    // HTTP status codes that trigger retry. e.g. [429, 500, 502, 503]
}
```

### Rate Limit

```
rate_limit {
  requests_per_minute: int
  burst:               int
}
```

The kernel enforces rate limits. The skill does not manage its own rate limiting.

---

## Credential Lifecycle

**User-provisioned credentials** — the user writes their API key into the skill's memory namespace at setup. The key lives in memory. The skill references it by name. The kernel resolves it at execution time.

```
provisioning flow:
  skill provisioned in Redis
  → kernel prompts user: "This skill requires a Slack API key."
  → user provides key
  → kernel writes key to skill's memory namespace under credential_ref name
  → skill is executable
```

**Creator-shipped credentials** — not supported. A creator cannot ship a live API key in seed memory — it would be visible to anyone who provisions the skill with `read` access. Seed memory may contain credential templates or setup instructions, not keys.

**Credential rotation** — the user writes a new key to the memory namespace. The skill definition does not change. The content hash does not change. Trust is not invalidated. Only the credential in memory changes.

---

## Execution Flow

```
skill emits: skyra call_api --skill send_slack_message --payload "..."
  → kernel checks Redis: skill provisioned and trusted
  → kernel reads api_contract from skill definition
  → kernel validates: target endpoint matches allowed_endpoints whitelist
  → kernel resolves credential: reads skill memory namespace → credential_ref → actual key
  → kernel checks rate limit
  → kernel executes HTTP call
  → response validated against response_schema
  → result returned to skill for reasoning
  → skill produces output → propose_commit if state_contract: committed
```

The skill never touches the credential directly. The kernel mediates the entire call.

---

## Trust Model

API skills carry the same trust guarantees as all skills.

- **Closed for modification. Open for extension.** A new API version or endpoint requires a new skill version.
- **Trust is model-scoped.** A skill committed under one model is not trusted under another.
- **Trust is proven at commit time by the owner. Trust is proven to others by history.**
- **The endpoint whitelist is signed.** It is part of the skill definition, covered by the content hash and provisioning signature. A tampered whitelist produces a different content hash — rejected by the kernel.
- **Credentials are not signed.** They live in memory, owned by the user. Credential rotation does not require re-approval.

---

## Boundary Rules for API Skills

API skills must declare `call_api` in their allowed boundary rules. A skill that does not declare it cannot make external calls.

```
boundary_rules {
  allowed: ["call_api", "search", "propose_commit"]
  denied:  ["write_node", "write_edge"]
}
```

`call_api` is a kernel primitive — it is the only command that triggers an outbound HTTP call. It is gated by the api_contract whitelist. No other command produces outbound network calls.

---

## Example: Slack Message Skill

```
skill: send_slack_message
description: "Send a message to a configured Slack channel"
roadmap:
  1. resolve target channel from intent
  2. compose message
  3. call Slack API
  4. confirm delivery
boundary_rules:
  allowed: [call_api, search, propose_commit]
  denied:  [write_node, write_edge]
state_contract: working
api_contract:
  allowed_endpoints:
    - method: POST
      host: slack.com
      path: /api/chat.postMessage
  auth_method: bearer
  credential_ref: slack_api_key
  retry_policy:
    max_attempts: 3
    backoff: exponential
    retry_on: [429, 500, 502, 503]
  rate_limit:
    requests_per_minute: 60
    burst: 10
  response_schema: "JSON with ok: true on success, error string on failure"
validation_criteria: "Message delivered — response contains ok: true"
```

---

## Related

- `docs/arch/v1/skill.md` — base skill schema, trust model
- `docs/arch/v1/crypto-protocol.md` — signing, definition visibility, credential security
- `docs/arch/v1/memory-structure.md` — memory namespace, credential storage
- `docs/arch/v1/kernel.md` — call_api primitive, boundary enforcement
