import WebSocket from "ws"

const being = process.argv[2] || "michael"
const ws = new WebSocket(`ws://localhost:3400/lens?being=${being}`)

ws.on("open", () => console.log(`connected as ${being}\n`))

ws.on("message", (raw) => {
  const present = JSON.parse(raw)
  console.log("--- present ---")
  for (const section of present.sections) {
    console.log(`[${section.type}]`, JSON.stringify(section.data, null, 2).slice(0, 300))
  }
  console.log()
})

ws.on("close", () => console.log("disconnected"))

// send a test relation after 2 seconds
setTimeout(() => {
  const relation = {
    id: "skyra",
    origin: being,
    threadId: "t1",
    impulse: "what do you think about persistent exchanges",
  }
  console.log(`sending: ${relation.origin} → ${relation.id}: ${relation.impulse}\n`)
  ws.send(JSON.stringify(relation))
}, 2000)

// close after 10 seconds
setTimeout(() => ws.close(), 10000)
