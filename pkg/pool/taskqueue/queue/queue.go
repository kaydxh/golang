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
	Add(ctx context.Context, msg *Message) (string, error)
	FetchOne(ctx context.Context, waitTimeout time.Duration) (*Message, error)
	Delete(ctx context.Context, msg *Message) error

	AddResult(ctx context.Context, result *MessageResult, expired time.Duration) (string, error)
	FetchResult(ctx context.Context, key string) (*MessageResult, error)
}

type ResultCallbackFunc func(ctx context.Context, result *MessageResult)
