package tasks

import (
	"context"
	"time"

	"github.com/google/uuid"
)

type Task struct {
	ID        uuid.UUID
	ExecuteAt time.Time     // для однократных
	Interval  time.Duration // для периодических (0 - однократная)
	Command   func(ctx context.Context) error
	Recurring bool
}

func NewTask(command func(ctx context.Context) error) Task {
	return Task{
		ID:      uuid.New(),
		Command: command,
	}
}
