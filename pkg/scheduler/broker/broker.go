package broker

import (
	"context"

	"github.com/kaydxh/golang/pkg/scheduler/task"
)

// Broker - a common interface for all brokers
// local node
// remote(ip:port) node
// the mq
type Broker interface {
	Publish(ctx context.Context, task *task.Task) error
}

// TaskProcessor - can process a delivered task
// This will probably always be a worker instance
type TaskProcessor interface {
	Process(task *task.Task) error
}
