package kernel

import (
	"context"
	"errors"
)

const defaultHeapCapacity = 1024

// InMemoryHeap is a simple bounded FIFO queue for local development.
// It satisfies the Heap interface and supports context-aware push/pop.
type InMemoryHeap struct {
	queue chan *Job
}

func NewInMemoryHeap(capacity int) *InMemoryHeap {
	if capacity <= 0 {
		capacity = defaultHeapCapacity
	}
	return &InMemoryHeap{
		queue: make(chan *Job, capacity),
	}
}

func (h *InMemoryHeap) Push(ctx context.Context, job *Job) error {
	if job == nil {
		return errors.New("heap: nil job")
	}

	select {
	case h.queue <- job:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

func (h *InMemoryHeap) Pop(ctx context.Context) (*Job, error) {
	select {
	case job := <-h.queue:
		return job, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}
