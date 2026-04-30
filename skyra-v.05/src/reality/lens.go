package reality

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"golang.org/x/net/websocket"
)

type Lens struct {
	id        string
	Realities map[string]Reality
	port      int
	conns     map[*websocket.Conn]string
	mu        sync.Mutex
	incoming  chan *Relation
}

func (l *Lens) ID() string { return l.id }

func (l *Lens) Create(r *Relation) Reality {
	port, _ := Extract(r.Impulse, "~port", "lens")
	p := 3400
	if port != "" {
		fmt.Sscanf(port, "%d", &p)
	}
	lens := &Lens{
		id:        "lens",
		Realities: make(map[string]Reality),
		port:      p,
		conns:     make(map[*websocket.Conn]string),
		incoming:  make(chan *Relation, 64),
	}
	go lens.serve()
	return lens
}

func (l *Lens) Realize(r *Relation) string {
	l.push(r)
	rel := <-l.incoming
	return rel.Impulse
}

func (l *Lens) Parse() string {
	return ""
}

func (l *Lens) serve() {
	mux := http.NewServeMux()
	mux.Handle("/lens", websocket.Handler(l.handle))
	addr := fmt.Sprintf(":%d", l.port)
	fmt.Printf("lens listening on ws://localhost%s/lens\n", addr)
	http.ListenAndServe(addr, mux)
}

func (l *Lens) handle(ws *websocket.Conn) {
	being := ws.Request().URL.Query().Get("being")
	l.mu.Lock()
	l.conns[ws] = being
	l.mu.Unlock()

	defer func() {
		l.mu.Lock()
		delete(l.conns, ws)
		l.mu.Unlock()
		ws.Close()
	}()

	for {
		var raw string
		if err := websocket.Message.Receive(ws, &raw); err != nil {
			return
		}
		var rel Relation
		if err := json.Unmarshal([]byte(raw), &rel); err != nil {
			continue
		}
		l.incoming <- &rel
	}
}

func (l *Lens) push(r *Relation) {
	present := l.buildPresent(r)
	data, err := json.Marshal(present)
	if err != nil {
		return
	}
	l.mu.Lock()
	defer l.mu.Unlock()
	for ws := range l.conns {
		websocket.Message.Send(ws, string(data))
	}
}

func (l *Lens) buildPresent(r *Relation) map[string]any {
	sections := []map[string]any{}
	for name, parser := range r.Parsers {
		sections = append(sections, map[string]any{
			"type": name,
			"data": parser(),
		})
	}
	return map[string]any{
		"being":    r.ID,
		"sections": sections,
	}
}
