# AI Beings App — Frontend

Frontend for the Skyra v.1 runtime. Renders the universe state object; doesn't model it separately.

See the backend spec for the wire contract — `Frontend Spec — v.1 Runtime` (Mike's doc).

## First-time setup

```bash
npm install
```

Then copy `.env.example` to `.env.local` and fill in any overrides:

```bash
cp .env.example .env.local
```

Defaults: WS server on `ws://localhost:8080`, auth token `dev-token`.

## Run

```bash
# Terminal 1 — mock WebSocket server (so the frontend has signs of life without the real backend)
npm run mock

# Terminal 2 — Next.js dev server
npm run dev
```

Open http://localhost:3000.

The mock server runs on `ws://localhost:8080`. The frontend reads the URL from `NEXT_PUBLIC_WS_URL` and the auth token from `NEXT_PUBLIC_WS_TOKEN` (defaults in `lib/ws-client.ts`).

## Layout

```
app/                    Next.js app router
  page.tsx              User present view (perspective being's exchanges)
  map/page.tsx          Universe present view (topology + all beings)
components/             Component implementations (names are transitional —
                        a Phase 2 rework will introduce <Being>, <Entry>,
                        <ExchangeView>, <UniversePresentView>).
lib/protocol/
  types.ts              Wire types — v.1 envelope, Being, Exchange, Thread,
                        TopologyNode, UniverseState, all deltas.
  schemas.ts            Zod schemas matching types.ts. Every inbound message
                        is validated; anything that fails to parse is dropped
                        and logged. This is how we catch backend drift.
lib/ws-client.ts        WebSocket client wrapper — first-message auth handshake,
                        reconnect with exponential backoff, outbox queue.
lib/store.ts            Zustand store. Holds the most recent universe snapshot
                        plus deltas applied since. No derived state in the store
                        body; selectors do the derivation.
mock/server.mjs         Mock WS server on :8080. Auth handshake, universe
                        snapshot on auth_ok, periodic entry + weight deltas.
```

## The two views

Per Mike's v.1 spec, the frontend toggles between two perspectives:

- **User present** (`/`) — the active being's exchanges, peers, what's in front of them. Append-only. No turns. Either party can post next.
- **Universe present** (`/map`) — full state of the runtime: every being, the topology, weighted Relationships and Expressors per being, economics, Skyra activity indicator.

## Scripts

- `npm run dev` — Next dev server
- `npm run mock` — Mock WebSocket server on :8080
- `npm run build` — Production build
- `npm run typecheck` — Type-check without emitting
- `npm run lint` — Next lint

## Status

Phase 1 (foundation) complete. v.1 protocol migration complete. Phase 2 (interaction surface depth + proper component naming) in progress.
