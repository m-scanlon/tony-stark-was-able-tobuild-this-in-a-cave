package kernel

import "context"

// Heap is the kernel's work queue.
// Events are ordered by priority — higher score, higher priority.
type Heap interface {
	Push(ctx context.Context, event *Event) error
	Pop(ctx context.Context) (*Event, error)
	Len() int
}
