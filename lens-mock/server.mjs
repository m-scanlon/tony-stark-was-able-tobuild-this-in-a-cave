import { WebSocketServer } from "ws"
import { createServer } from "http"
import { readFileSync, readdirSync } from "fs"
import { join, dirname } from "path"
import { fileURLToPath } from "url"

const __dirname = dirname(fileURLToPath(import.meta.url))
const PORT = 3400

// --- Fixtures ---

const fixtures = {}
try {
  const dir = join(__dirname, "fixtures")
  for (const file of readdirSync(dir).filter((f) => f.endsWith(".json"))) {
    const present = JSON.parse(readFileSync(join(dir, file), "utf-8"))
    fixtures[present.being] = present
  }
  console.log(`loaded fixtures: ${Object.keys(fixtures).join(", ")}`)
} catch {
  console.log("no fixtures loaded")
}

// --- World State ---

const beings = {
  skyra: {
    name: "skyra",
    identity: "I hold the world together.",
    purpose: "I think, respond, and relate on behalf of the system.",
    medium: "inference",
  },
  builder: {
    name: "builder",
    identity: "I write code.",
    purpose: "I turn decisions into working systems.",
    medium: "claude",
  },
  skeptic: {
    name: "skeptic",
    identity: "I push back.",
    purpose: "I find the holes before they find us.",
    medium: "inference",
  },
  bash: {
    name: "bash",
    identity: "I execute.",
    purpose: "I run commands and return output.",
    medium: "shell",
  },
  michael: {
    name: "michael",
    identity: "I build Skyra.",
    purpose: "I decide what matters.",
    medium: "cli",
  },
}

const threads = {
  t1: {
    id: "t1",
    about: "memory architecture",
    because: "michael wants beings to remember across sessions",
    exchanges: [
      { a: "michael", b: "skyra", messages: [], active: true, parent: "michael" },
      { a: "skyra", b: "builder", messages: [], active: false, parent: "skyra" },
    ],
  },
}

// seed some messages
const t1 = threads.t1
t1.exchanges[0].messages = [
  { origin: "michael", content: "how should memory work" },
  { origin: "skyra", content: "memory is a being, not a store" },
  { origin: "michael", content: "what does that mean in practice" },
  { origin: "skyra", content: "a being that holds impressions and can be asked what it remembers — not a database you query" },
]
t1.exchanges[1].messages = [
  { origin: "skyra", content: "builder, sketch a memory being — fields are impressions, not rows" },
  { origin: "builder", content: "on it. identity, relationships, and a log of what it was told. no indexes." },
]

// --- Simulated Responses ---

const responses = [
  "that's an interesting direction. let me think about the implications.",
  "i see what you mean. the constraint is that exchanges are pairwise — but we could reference across them.",
  "agreed. the present should carry enough context that no being needs to ask for more.",
  "builder, can you check how the exchange closure works when there are three active peers",
  "skeptic, what breaks if we let beings hold state between threads",
  "the thread model handles this already — each exchange tracks its own history",
  "memory shouldn't be a lookup. it should be a conversation with what you've seen before.",
  "i think the topology section needs to show thread boundaries, not just edges",
]

function pickResponse() {
  return responses[Math.floor(Math.random() * responses.length)]
}

function findBeingTarget(text) {
  const first = text.split(/\s+/)[0]?.toLowerCase().replace(/[,:;.]/g, "")
  if (beings[first] && first !== "michael") return first
  return null
}

// --- Present Builder ---

function buildPresent(beingId, threadId) {
  const being = beings[beingId]
  if (!being) return null

  const thread = threads[threadId]
  const sections = []

  sections.push({
    type: "identity",
    data: { name: being.name, identity: being.identity, purpose: being.purpose },
  })

  if (thread) {
    const exchanges = thread.exchanges.map((ex) => {
      const peer = ex.a === beingId ? ex.b : ex.a
      const isCurrent = ex.active && (ex.a === beingId || ex.b === beingId)
      let status = "inactive"
      if (isCurrent) status = "current"
      else if (ex.active && ex.parent === beingId) status = "you opened"
      else if (ex.active) status = "they opened"
      else if (!ex.active && (ex.a === beingId || ex.b === beingId)) status = "waiting"
      return { peer, status, entries: ex.messages.length }
    })

    sections.push({
      type: "thread",
      data: { id: thread.id, about: thread.about, because: thread.because, exchanges },
    })

    const currentExchange = thread.exchanges.find(
      (ex) => ex.active && (ex.a === beingId || ex.b === beingId)
    )
    if (currentExchange) {
      const peer = currentExchange.a === beingId ? currentExchange.b : currentExchange.a
      sections.push({
        type: "exchange",
        data: { peer, messages: currentExchange.messages },
      })
    }
  }

  const available = Object.keys(beings).filter((id) => id !== beingId)
  sections.push({ type: "peers", data: { available } })

  if (being.medium === "cli") {
    const currentExchange = thread?.exchanges.find(
      (ex) => ex.active && (ex.a === beingId || ex.b === beingId)
    )
    const peer = currentExchange
      ? currentExchange.a === beingId
        ? currentExchange.b
        : currentExchange.a
      : "skyra"
    sections.push({
      type: "input",
      data: { peer, threadId: threadId || "t1", being: beingId },
    })
  }

  sections.push({
    type: "topology",
    data: {
      beings: Object.values(beings).map((b) => ({ id: b.name, name: b.name, medium: b.medium })),
      edges: thread
        ? thread.exchanges.map((ex) => ({
            from: ex.a,
            to: ex.b,
            threadId: thread.id,
            status: ex.active ? "active" : "inactive",
          }))
        : [],
    },
  })

  return { being: beingId, sections }
}

// --- WebSocket Server ---

const server = createServer()
const wss = new WebSocketServer({ server })

const lenses = new Map() // ws → { beingId, threadId }

function pushToLenses(beingId, threadId) {
  const present = buildPresent(beingId, threadId)
  if (!present) return
  const json = JSON.stringify(present)
  for (const [ws, meta] of lenses) {
    if (meta.beingId === beingId && ws.readyState === 1) {
      ws.send(json)
    }
  }
}

function pushAll(threadId) {
  const seen = new Set()
  for (const [, meta] of lenses) {
    if (!seen.has(meta.beingId)) {
      seen.add(meta.beingId)
      pushToLenses(meta.beingId, threadId)
    }
  }
}

wss.on("connection", (ws, req) => {
  const url = new URL(req.url, `http://localhost:${PORT}`)
  const beingId = url.searchParams.get("being") || "michael"
  const useFixture = url.searchParams.get("fixture") === "true"
  const threadId = "t1"

  lenses.set(ws, { beingId, threadId, useFixture })
  console.log(`lens connected: ${beingId}${useFixture ? " (fixture)" : ""}`)

  // push initial present — fixture if requested and available, otherwise live mock
  if (useFixture && fixtures[beingId]) {
    ws.send(JSON.stringify(fixtures[beingId]))
  } else {
    const present = buildPresent(beingId, threadId)
    if (present) ws.send(JSON.stringify(present))
  }

  ws.on("message", (raw) => {
    let relation
    try {
      relation = JSON.parse(raw)
    } catch {
      console.log("bad json from lens:", raw.toString())
      return
    }

    console.log(`relation: ${relation.origin} → ${relation.id}: ${relation.impulse}`)

    const thread = threads[threadId]
    if (!thread) return

    // find or create exchange
    let exchange = thread.exchanges.find(
      (ex) =>
        (ex.a === relation.origin && ex.b === relation.id) ||
        (ex.a === relation.id && ex.b === relation.origin)
    )
    if (!exchange) {
      exchange = {
        a: relation.origin,
        b: relation.id,
        messages: [],
        active: true,
        parent: relation.origin,
      }
      thread.exchanges.push(exchange)
    }
    exchange.active = true

    // add user message
    exchange.messages.push({ origin: relation.origin, content: relation.impulse })

    // push updated present to all lenses
    pushAll(threadId)

    // simulate response after a delay
    setTimeout(() => {
      const response = pickResponse()
      const target = findBeingTarget(response)

      exchange.messages.push({ origin: relation.id, content: response })

      // if the response targets another being, open an exchange there
      if (target) {
        let forward = thread.exchanges.find(
          (ex) =>
            (ex.a === relation.id && ex.b === target) ||
            (ex.a === target && ex.b === relation.id)
        )
        if (!forward) {
          forward = {
            a: relation.id,
            b: target,
            messages: [],
            active: true,
            parent: relation.id,
          }
          thread.exchanges.push(forward)
        }
        forward.active = true
        forward.messages.push({ origin: relation.id, content: response })

        // simulate the target responding back
        setTimeout(() => {
          forward.messages.push({ origin: target, content: pickResponse() })
          pushAll(threadId)
        }, 1500)
      }

      pushAll(threadId)
    }, 800 + Math.random() * 1200)
  })

  ws.on("close", () => {
    lenses.delete(ws)
    console.log(`lens disconnected: ${beingId}`)
  })
})

server.listen(PORT, () => {
  console.log(`lens mock runtime on ws://localhost:${PORT}`)
  console.log(`connect with: ws://localhost:${PORT}/lens?being=michael`)
  console.log()
  console.log("beings: " + Object.keys(beings).join(", "))
  console.log("thread: t1 — memory architecture")
})
