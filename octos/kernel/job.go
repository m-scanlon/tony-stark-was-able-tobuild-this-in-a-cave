package kernel

import "time"

type JobStatus string

const (
	JobStatusPending  JobStatus = "pending"
	JobStatusRunning  JobStatus = "running"
	JobStatusComplete JobStatus = "complete"
	JobStatusFailed   JobStatus = "failed"
	JobStatusTimedOut JobStatus = "timed_out"
)

type TaskStatus string

const (
	TaskStatusPending  TaskStatus = "pending"
	TaskStatusRunning  TaskStatus = "running"
	TaskStatusComplete TaskStatus = "complete"
	TaskStatusFailed   TaskStatus = "failed"
)

// Job is a skill instance. Skill is the class, Job is the object.
type Job struct {
	ID           string
	SkillID      string
	ParentTaskID string
	TurnID       string
	SessionID    string
	Status       JobStatus
	CreatedAt    time.Time
	CompletedAt  *time.Time
	Tasks        []*Task
}

// Task is the atomic execution unit inside a job.
type Task struct {
	ID          string
	JobID       string
	Skill       string
	ReplicaID   string
	Status      TaskStatus
	Result      string
	CreatedAt   time.Time
	CompletedAt *time.Time
}
