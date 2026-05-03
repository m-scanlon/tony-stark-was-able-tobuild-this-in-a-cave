package reality

import (
	"encoding/json"
	"fmt"
	"net/http"
	"skyra-v05/src/debug"
	"sync"
	"time"

	"golang.org/x/net/websocket"
)

type WS struct {
	id      string
	Device  Reality
	Port    int
	clients map[*websocket.Conn]bool
	mu      sync.Mutex
	inbox   chan string
}

func (w *WS) ID() string { return w.id }

func (w *WS) Create(r *Relation) Reality {
	return &WS{
		id:      "ws",
		clients: make(map[*websocket.Conn]bool),
		inbox:   make(chan string),
	}
}

func (w *WS) Connected() bool {
	w.mu.Lock()
	defer w.mu.Unlock()
	return len(w.clients) > 0
}

func (w *WS) Start(port int) {
	w.Port = port
	mux := http.NewServeMux()
	mux.Handle("/ws", websocket.Handler(w.handle))
	addr := fmt.Sprintf(":%d", port)
	debug.Log("[ws]: starting on", addr)
	go http.ListenAndServe(addr, mux)
}

func (w *WS) handle(conn *websocket.Conn) {
	debug.Log("[ws]: client connected")
	w.mu.Lock()
	w.clients[conn] = true
	w.mu.Unlock()

	defer func() {
		w.mu.Lock()
		delete(w.clients, conn)
		w.mu.Unlock()
		conn.Close()
		debug.Log("[ws]: client disconnected")
	}()

	for {
		var raw string
		err := websocket.Message.Receive(conn, &raw)
		if err != nil {
			debug.Log("[ws]: read error:", err)
			return
		}
		debug.Log("[ws]: received:", raw)

		var msg wsInbound
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			debug.Log("[ws]: parse error:", err)
			continue
		}

		if msg.Type == "input" && msg.Payload.Content != "" {
			w.inbox <- msg.Payload.Content
		}
	}
}

func (w *WS) Broadcast(state string) {
	msg := wsOutbound{
		Type:    "universe",
		Ts:      time.Now().UnixMilli(),
		Payload: json.RawMessage(state),
	}
	data, err := json.Marshal(msg)
	if err != nil {
		debug.Log("[ws]: marshal error:", err)
		return
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	for conn := range w.clients {
		if err := websocket.Message.Send(conn, string(data)); err != nil {
			debug.Log("[ws]: send error:", err)
			conn.Close()
			delete(w.clients, conn)
		}
	}
}

func (w *WS) Realize(r *Relation) string {
	if r.Collecting {
		return ""
	}

	if r.Impulse != "" {
		msg := wsOutbound{
			Type: "impulse",
			Ts:   time.Now().UnixMilli(),
			Payload: marshalJSON(wsImpulse{
				From:    r.Origin,
				Content: r.Impulse,
			}),
		}
		data, _ := json.Marshal(msg)

		w.mu.Lock()
		for conn := range w.clients {
			websocket.Message.Send(conn, string(data))
		}
		w.mu.Unlock()
	}

	return <-w.inbox
}

func marshalJSON(v any) json.RawMessage {
	data, _ := json.Marshal(v)
	return data
}

type wsInbound struct {
	Type    string `json:"type"`
	Payload struct {
		Content string `json:"content"`
	} `json:"payload"`
}

type wsOutbound struct {
	Type    string          `json:"type"`
	Ts      int64           `json:"ts"`
	Payload json.RawMessage `json:"payload"`
}

type wsImpulse struct {
	From    string `json:"from"`
	Content string `json:"content"`
}
