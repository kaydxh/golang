package queue

import (
	"context"
	"time"
)

const (
	TaskQueuePrefix = "taskq"
)

type QueueOptions struct {
	Name string
}

type Queue interface {
	Add(ctx context.Context, msg *Message) error
	FetchOne(ctx context.Context, waitTimeout time.Duration) (*Message, error)
}
