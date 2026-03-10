package kernel

import (
	"context"
	"errors"
	"fmt"
	"log"
)

var ErrSkillNotFound = errors.New("skill not found or not provisioned")

// Kernel is the central execution boundary.
// Shards reason. The kernel executes. There is no other mode.
//
// Every command resolves against the skill registry.
// No command bypasses Redis — not even system primitives.
// If a skill is not in Redis, it does not run.
type Kernel struct {
	registry SkillRegistry
	heap     Heap
}

func New(registry SkillRegistry, heap Heap) *Kernel {
	return &Kernel{registry: registry, heap: heap}
}

// Run starts the execution loop. Blocks until ctx is cancelled.
// Must be called in its own goroutine.
func (k *Kernel) Run(ctx context.Context) {
	log.Println("kernel: execution loop started")
	for {
		job, err := k.heap.Pop(ctx)
		if err != nil {
			if ctx.Err() != nil {
				log.Println("kernel: execution loop stopped")
				return
			}
			log.Printf("kernel: heap pop error: %v", err)
			continue
		}
		if job == nil {
			continue
		}
		go k.execute(ctx, job)
	}
}

// execute runs a single job. Each job runs in its own goroutine.
func (k *Kernel) execute(ctx context.Context, job *Job) {
	log.Printf("kernel: executing job=%s skill=%s", job.ID, job.SkillID)
	// TODO: route job to shard, run task loop, handle ReAct cycle
}

// Dispatch parses a command, checks the registry, and pushes a job to the heap.
func (k *Kernel) Dispatch(ctx context.Context, raw string) error {
	cmd, err := ParseCommand(raw)
	if err != nil {
		return fmt.Errorf("dispatch: %w", err)
	}

	skill, err := k.registry.Get(ctx, cmd.Tool)
	if err != nil {
		return fmt.Errorf("kernel: registry lookup failed for %q: %w", cmd.Tool, err)
	}
	if skill == nil {
		return fmt.Errorf("kernel: %w: %q", ErrSkillNotFound, cmd.Tool)
	}

	job := k.instantiateJob(skill, cmd)

	if err := k.heap.Push(ctx, job); err != nil {
		return fmt.Errorf("kernel: failed to push job to heap: %w", err)
	}

	log.Printf("kernel: skill=%q job=%s pushed to heap", skill.Name, job.ID)
	return nil
}

func (k *Kernel) instantiateJob(skill *Skill, cmd Command) *Job {
	job := &Job{
		ID:      newID(),
		SkillID: skill.ID,
		Status:  JobStatusPending,
	}
	for _, st := range skill.Tasks {
		job.Tasks = append(job.Tasks, &Task{
			ID:     newID(),
			JobID:  job.ID,
			Skill:  st.Name,
			Status: TaskStatusPending,
		})
	}
	return job
}
