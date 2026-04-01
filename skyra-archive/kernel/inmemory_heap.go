package kernel

import (
	"container/heap"
	"context"
	"errors"
	"sync"
)

const defaultHeapCapacity = 1024

var ErrHeapFull = errors.New("heap: capacity reached")

type eventPriorityQueue []*Event

func (q eventPriorityQueue) Len() int { return len(q) }

func (q eventPriorityQueue) Less(i, j int) bool {
	if q[i].Priority == q[j].Priority {
		return q[i].sequence < q[j].sequence
	}
	return q[i].Priority > q[j].Priority
}

func (q eventPriorityQueue) Swap(i, j int) {
	q[i], q[j] = q[j], q[i]
}

func (q *eventPriorityQueue) Push(x any) {
	*q = append(*q, x.(*Event))
}

func (q *eventPriorityQueue) Pop() any {
	old := *q
	n := len(old)
	item := old[n-1]
	*q = old[:n-1]
	return item
}

// InMemoryHeap is an in-process max heap used by the v1 kernel runtime.
type InMemoryHeap struct {
	mu       sync.Mutex
	notify   chan struct{}
	items    eventPriorityQueue
	capacity int
	sequence uint64
}

func NewInMemoryHeap(capacity int) *InMemoryHeap {
	if capacity <= 0 {
		capacity = defaultHeapCapacity
	}
	h := &InMemoryHeap{
		notify:   make(chan struct{}, 1),
		capacity: capacity,
		items:    make(eventPriorityQueue, 0, capacity),
	}
	heap.Init(&h.items)
	return h
}

func (h *InMemoryHeap) Push(ctx context.Context, event *Event) error {
	if event == nil {
		return errors.New("heap: nil event")
	}

	h.mu.Lock()
	defer h.mu.Unlock()

	if h.capacity > 0 && len(h.items) >= h.capacity {
		return ErrHeapFull
	}

	h.sequence++
	event.sequence = h.sequence
	heap.Push(&h.items, event)

	select {
	case h.notify <- struct{}{}:
	default:
	}

	return nil
}

func (h *InMemoryHeap) Pop(ctx context.Context) (*Event, error) {
	for {
		h.mu.Lock()
		if len(h.items) > 0 {
			item := heap.Pop(&h.items).(*Event)
			h.mu.Unlock()
			return item, nil
		}
		wait := h.notify
		h.mu.Unlock()

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-wait:
		}
	}
}

func (h *InMemoryHeap) Len() int {
	h.mu.Lock()
	defer h.mu.Unlock()
	return len(h.items)
}
