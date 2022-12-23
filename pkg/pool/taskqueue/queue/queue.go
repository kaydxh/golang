package queue

import (
	"context"
)

const (
	TaskQueuePrefix = "taskq"
)

type QueueOptions struct {
	Name string

	BufferSize int64
}

type Queue interface {
	Add(ctx context.Context, msg *Message) error
}
