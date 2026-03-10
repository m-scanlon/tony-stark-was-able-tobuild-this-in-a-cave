# Skill: API Compatibility

An API skill wraps an external API as an executable skill. It is a first-class skill — same schema, same trust model, same crypto guarantees. The API surface it is allowed to hit is declared in the skill contract. The kernel enforces it.

---

## The Problem

External APIs require credentials, endpoint definitions, retry logic, and rate limits. None of this can live in the skill definition. Credentials in the definition are a security leak.

The harder problem is auth. API key is the simple case. Most companies do not use simple API keys. OAuth2, CLI-based auth (AWS, gcloud), service accounts, OIDC, mTLS — each is different. The system cannot be modified every time a new auth pattern appears.

**The solution is extension, not modification.** Credentials live in Redis — the trust boundary. Auth method support is a shard capability. The kernel handles base cases. Complex or exotic auth is handled by a shard that registers the capability. The skill declares what auth it needs. The kernel routes to the right shard. The core system never changes.

Closed for modification. Open for extension.

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

### Auth Method

```
auth_method: api_key | bearer | oauth2 | shard:<capability_name>
```

**Kernel-native** (handled directly):
- `api_key` — key passed as header or query param
- `bearer` — Bearer token in Authorization header
- `oauth2` — standard OAuth2 flow, token managed in Redis

**Shard-extended** — for everything else. The skill declares `shard:<capability>`. The kernel routes the auth step to a shard registered with that capability. That shard handles the auth, returns a resolved credential to the kernel, and the kernel proceeds.

```
// AWS CLI-based auth
auth_method: shard:aws_auth

// A shard registers:
capability: aws_auth   // knows how to invoke AWS CLI, assume roles, resolve temp credentials
```

Adding support for a new auth method means registering a new shard with a new capability. The core system is unchanged.

---

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

**Credentials live in Redis only.** Not in the skill definition. Not in memory. Redis is secure — protected by mTLS at the transport layer and signed provisioning records at the authorization layer. A caller that has proven ownership through the Redis trust chain gets the credential. No other path exists.

```
Redis entry for an API skill:
  skill:send_slack_message → {
    contract:        { ...skill definition... },
    credential_ref:  "slack_api_key",
    credential:      "xoxb-..."       // lives here, protected by Redis auth
  }
```

```
provisioning flow:
  skill provisioned in Redis (via provision_skill, user-signed)
  → user provides API key
  → key written to Redis under credential_ref (via Redis write skill, user-signed)
  → skill is executable
```

**Credential rotation** — write a new key to Redis via the Redis write skill. Requires user signature. The skill definition does not change. The content hash does not change. Trust is not invalidated. Only the credential in Redis changes.

**Creator-shipped credentials** — not supported. A creator provisions a skill without a credential. The consumer writes their own key into Redis at setup time. The skill contract declares `credential_ref` — the name the kernel uses to look it up.

---

## Execution Flow

```
skill emits: skyra call_api --skill send_slack_message --payload "..."
  → kernel checks Redis: skill provisioned and trusted
  → kernel reads api_contract + credential from Redis (single lookup — same trust boundary)
  → kernel validates: target endpoint matches allowed_endpoints whitelist
  → kernel checks rate limit
  → kernel executes HTTP call with credential
  → response validated against response_schema
  → result returned to skill for reasoning
  → skill produces output → propose_commit if state_contract: committed
```

The skill never touches the credential directly. The kernel resolves it from Redis and mediates the entire call. Proving ownership of the skill proves access to the credential — same chain, same boundary.

---

## Trust Model

API skills carry the same trust guarantees as all skills.

- **Closed for modification. Open for extension.** A new API version or endpoint requires a new skill version.
- **Trust is model-scoped.** A skill committed under one model is not trusted under another.
- **Trust is proven at commit time by the owner. Trust is proven to others by history.**
- **The endpoint whitelist is signed.** It is part of the skill definition, covered by the content hash and provisioning signature. A tampered whitelist produces a different content hash — rejected by the kernel.
- **Credentials live in Redis.** Protected by the same mTLS and signing that protects all Redis entries. Proving ownership of the skill proves access to the credential. Credential rotation is a Redis write — requires user signature via the Redis write skill.

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
