package kernel

import "context"

// Heap is the kernel's work queue.
// Jobs are ordered by importance score — higher score, higher priority.
type Heap interface {
	Push(ctx context.Context, job *Job) error
	Pop(ctx context.Context) (*Job, error)
}
