package queue

import (
	"context"
	"time"
)

const (
	TaskQueuePrefix = "taskq"
)

/*
var (
  NoExistMessage
)
*/

type QueueOptions struct {
	Name string
}

type Queue interface {
	Add(ctx context.Context, msg *Message) (string, error)
	FetchOne(ctx context.Context, waitTimeout time.Duration) (*Message, error)
	Delete(ctx context.Context, msg *Message) error
}
