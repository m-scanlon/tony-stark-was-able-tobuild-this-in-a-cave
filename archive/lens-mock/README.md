# Lens Mock Server

Mock runtime that pushes presents over WebSocket so you can build lens components without the Go runtime.

## Setup

```
cd lens-mock
npm install ws
```

## Run

```
node server.mjs
```

Starts on `ws://localhost:3400`.

## Connect

### Live mode — synthetic data, simulated responses

```
ws://localhost:3400/lens?being=michael
```

Beings available: `michael`, `skyra`, `builder`, `skeptic`, `bash`.

Send a relation as JSON and the mock will simulate a response after a short delay. If the response targets another being, the mock opens a new exchange and simulates that too.

### Fixture mode — real session data, static

```
ws://localhost:3400/lens?being=builder&fixture=true
```

Serves a present captured from a real runtime session. Use this for static component development — the data is richer and more realistic than live mode.

Available fixtures are loaded from `fixtures/`. Any `.json` file with a `being` field gets picked up on startup.

## Relation format

This is what the lens sends back through the socket when the user acts:

```json
{
  "id": "skyra",
  "origin": "michael",
  "threadId": "t1",
  "impulse": "what the user typed"
}
```

- `id` — target being
- `origin` — the being this lens renders for
- `threadId` — current thread from the present
- `impulse` — user input

## Present format

This is what the server pushes. The lens replaces its entire state with each push.

```json
{
  "being": "michael",
  "sections": [
    { "type": "identity", "data": { ... } },
    { "type": "thread", "data": { ... } },
    { "type": "exchange", "data": { ... } },
    { "type": "peers", "data": { ... } },
    { "type": "input", "data": { ... } },
    { "type": "topology", "data": { ... } }
  ]
}
```

Each section type maps to a component. See `lens-implementation-spec.md` in the project root for full data shapes.

## Test client

```
node test-client.mjs michael
```

Connects, prints presents, sends a test relation after 2 seconds, disconnects after 10.

## Adding fixtures

Run a real session, copy the present log, structure it as JSON, drop it in `fixtures/`. The server picks it up on next start. The `being` field in the JSON determines which `?being=` param serves it.
