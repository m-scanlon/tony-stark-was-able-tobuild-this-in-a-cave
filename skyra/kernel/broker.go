package kernel

import "sync"

type UIEvent struct {
	Type string `json:"type"`
	Data any    `json:"data"`
}

type Broker struct {
	mu     sync.RWMutex
	nextID int
	subs   map[int]chan UIEvent
}

func NewBroker() *Broker {
	return &Broker{subs: make(map[int]chan UIEvent)}
}

func (b *Broker) Subscribe() (int, <-chan UIEvent) {
	b.mu.Lock()
	defer b.mu.Unlock()

	b.nextID++
	id := b.nextID
	ch := make(chan UIEvent, 32)
	b.subs[id] = ch
	return id, ch
}

func (b *Broker) Unsubscribe(id int) {
	b.mu.Lock()
	defer b.mu.Unlock()

	ch, ok := b.subs[id]
	if !ok {
		return
	}
	delete(b.subs, id)
	close(ch)
}

func (b *Broker) Publish(event UIEvent) {
	b.mu.RLock()
	defer b.mu.RUnlock()

	for _, ch := range b.subs {
		select {
		case ch <- event:
		default:
		}
	}
}
