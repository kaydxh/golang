package taskqueue

import (
	"context"

	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
)

type PoolOptions struct {
	burst      uint32
	bufferSize uint32
}

type Pool struct {
	buffer chan *queue_.Message
	taskq  queue_.Queue
}

func NewPool(taskq queue_.Queue) *Pool {
	p := &Pool{
		taskq: taskq,
	}

	return p
}

func (p *Pool) Publish(ctx context.Context, msg *queue_.Message) error {
	return p.taskq.Add(ctx, msg)
}
