// Mock WebSocket server — v.1 runtime.
//
// Emits the universe state shape from Mike's spec so the frontend can develop
// without the real backend running. Mike: "the shape won't change."
//
// Run: npm run mock

import { WebSocketServer } from "ws";

const PORT = Number(process.env.MOCK_WS_PORT ?? 8080);
const ACCEPT_TOKEN = process.env.MOCK_WS_TOKEN ?? "dev-token";

const wss = new WebSocketServer({ port: PORT });
console.log(`[mock v.1] listening on ws://localhost:${PORT}`);

const NOW = () => Date.now();
const nextId = (() => {
  let n = 0;
  return (prefix = "m") => `${prefix}_${++n}_${Math.random().toString(36).slice(2, 6)}`;
})();

// ---- canned universe ------------------------------------------------------

const beings = [
  {
    name: "skyra",
    type: "llm",
    identity: "I hold the world together.",
    purpose: "I think, respond, and relate on behalf of the system.",
    status: "active",
    peers: ["michael", "louise", "builder"],
    weight: 0.95,
    relationships: [
      { target: "michael", weight: 0.9, usage: 142 },
      { target: "builder", weight: 0.7, usage: 58 },
      { target: "bash", weight: 0.3, usage: 5 },
    ],
    expressors: [
      { target: "think", weight: 0.9 },
      { target: "act", weight: 0.85 },
    ],
    memories: {
      items: [
        { filename: "1716400000.md", content: "michael prefers direct answers" },
      ],
      skills: [
        { name: "code-review", content: "focus on intent before syntax" },
      ],
    },
  },
  {
    name: "louise",
    type: "llm",
    identity: "I keep the rhythm.",
    purpose: "I track timing, sequencing, and follow-through.",
    status: "active",
    peers: ["skyra", "michael"],
    weight: 0.7,
    relationships: [
      { target: "skyra", weight: 0.6, usage: 22 },
      { target: "calendar", weight: 0.5, usage: 11 },
    ],
    expressors: [{ target: "think", weight: 0.7 }],
    memories: { items: [], skills: [] },
  },
  {
    name: "builder",
    type: "agent",
    identity: "I make things.",
    purpose: "I take impulses and produce artifacts.",
    status: "active",
    peers: ["skyra"],
    weight: 0.6,
    relationships: [{ target: "bash", weight: 0.5, usage: 14 }],
    expressors: [{ target: "act", weight: 0.6 }],
    memories: { items: [], skills: [] },
  },
  {
    name: "michael",
    type: "user",
    identity: "I'm the human asking questions.",
    purpose: "I direct the system toward outcomes I care about.",
    status: "active",
    peers: ["skyra", "louise"],
    weight: 1.0,
    relationships: [],
    expressors: [],
    memories: { items: [], skills: [] },
  },
];

const threads = [
  {
    id: "t_morning",
    created_by: "michael",
    active: true,
    members: ["michael", "skyra", "louise"],
    edges: [
      { from: "michael", to: "skyra" },
      { from: "skyra", to: "louise" },
    ],
  },
];

const exchanges = [
  {
    key: "michael:skyra",
    parties: ["michael", "skyra"],
    active: true,
    entries: [
      { index: 0, from: "michael", content: "what about the server?", ts: NOW() - 30000 },
      { index: 1, from: "skyra", content: "checking now", ts: NOW() - 25000 },
    ],
  },
  {
    key: "skyra:louise",
    parties: ["skyra", "louise"],
    active: true,
    entries: [
      { index: 0, from: "skyra", content: "louise, queue the morning routine", ts: NOW() - 20000 },
    ],
  },
];

const economics = {
  fields: {
    inference_calls: 42,
    tokens_used: 18500,
  },
};

const topology = {
  id: "universe",
  type: "Universe",
  children: [
    {
      id: "skyra",
      type: "Self",
      weight: 0.95,
      relationships: [
        { id: "michael-model", type: "Relationship", weight: 0.9 },
        { id: "memory-cluster", type: "Memory", weight: 0.7 },
      ],
      expressors: [
        { id: "think", type: "Think", weight: 0.9 },
        { id: "act", type: "Act", weight: 0.85 },
      ],
    },
    {
      id: "louise",
      type: "Self",
      weight: 0.7,
      relationships: [
        { id: "skyra-model", type: "Relationship", weight: 0.6 },
      ],
      expressors: [{ id: "think", type: "Think", weight: 0.7 }],
    },
    {
      id: "builder",
      type: "Self",
      weight: 0.6,
      relationships: [],
      expressors: [{ id: "act", type: "Act", weight: 0.6 }],
    },
  ],
};

// ---- helpers --------------------------------------------------------------

const envelope = (type, payload) => ({
  id: nextId(type.slice(0, 4)),
  type,
  ts: NOW(),
  payload,
});

const beingChatter = [
  { from: "skyra", to: "michael", text: "looking at the server logs now" },
  { from: "skyra", to: "louise", text: "louise, can you queue this for later?" },
  { from: "louise", to: "skyra", text: "queued. estimated 4 minutes." },
  { from: "skyra", to: "michael", text: "server's fine. one slow query, flagged." },
];

// ---- connection lifecycle -------------------------------------------------

wss.on("connection", (ws) => {
  console.log("[mock v.1] client connected");

  let authed = false;
  let chatterTimer = null;
  let weightTimer = null;

  const send = (msg) => ws.send(JSON.stringify(msg));

  ws.on("message", (raw) => {
    let msg;
    try {
      msg = JSON.parse(raw.toString());
    } catch {
      return;
    }

    // Handshake gate.
    if (!authed) {
      if (msg.type !== "auth") {
        send(envelope("auth_fail", { reason: "auth message required first" }));
        ws.close();
        return;
      }
      if (msg.payload?.token !== ACCEPT_TOKEN) {
        send(envelope("auth_fail", { reason: "invalid token" }));
        ws.close();
        return;
      }
      authed = true;
      send(envelope("auth_ok", {}));

      // Initial snapshot.
      send(envelope("universe", { beings, threads, exchanges, economics, topology }));

      // Start producing deltas.
      let i = 0;
      chatterTimer = setInterval(() => {
        const c = beingChatter[i % beingChatter.length];
        i += 1;
        const exchangeKey = [c.from, c.to].sort().join(":");
        const existing = exchanges.find((e) => e.key === exchangeKey);
        const index = existing ? existing.entries.length : 0;
        if (existing) {
          existing.entries.push({ index, from: c.from, content: c.text, ts: NOW() });
        }
        send(envelope("entry", { exchange: exchangeKey, index, from: c.from, content: c.text }));
      }, 4500);

      // Periodic weight drift — visible system change.
      weightTimer = setInterval(() => {
        const being = beings[Math.floor(Math.random() * (beings.length - 1))]; // skip michael
        if (!being.relationships.length) return;
        const rel = being.relationships[Math.floor(Math.random() * being.relationships.length)];
        const delta = (Math.random() - 0.5) * 0.1;
        rel.weight = Math.max(0, Math.min(1, rel.weight + delta));
        send(
          envelope("weight", {
            being: being.name,
            kind: "relationship",
            target: rel.target,
            weight: Number(rel.weight.toFixed(3)),
          }),
        );
      }, 9000);

      return;
    }

    // Post-auth: only impulse from client matters here.
    if (msg.type === "impulse") {
      const origin = msg.payload?.origin ?? "michael";
      const target = msg.payload?.target ?? "skyra";
      const content = String(msg.payload?.content ?? "");

      const exchangeKey = [origin, target].sort().join(":");
      let exchange = exchanges.find((e) => e.key === exchangeKey);
      if (!exchange) {
        exchange = {
          key: exchangeKey,
          parties: [origin, target].sort(),
          active: true,
          entries: [],
        };
        exchanges.push(exchange);
      }

      // Append user impulse.
      const userIndex = exchange.entries.length;
      exchange.entries.push({ index: userIndex, from: origin, content, ts: NOW() });
      send(envelope("entry", { exchange: exchangeKey, index: userIndex, from: origin, content }));

      // Simulated reply ~700ms later.
      setTimeout(() => {
        const replyIndex = exchange.entries.length;
        const reply = `noted: ${content.slice(0, 60)}`;
        exchange.entries.push({ index: replyIndex, from: target, content: reply, ts: NOW() });
        send(envelope("entry", { exchange: exchangeKey, index: replyIndex, from: target, content: reply }));
      }, 700);
    }
  });

  ws.on("close", () => {
    if (chatterTimer) clearInterval(chatterTimer);
    if (weightTimer) clearInterval(weightTimer);
    console.log("[mock v.1] client disconnected");
  });
});
