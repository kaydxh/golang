package memory

import (
	"context"
	"errors"

	"github.com/kaydxh/golang/pkg/scheduler/broker"
	"github.com/kaydxh/golang/pkg/scheduler/task"
)

type Broker struct {
	dispatcher broker.Taskdispatcher
}

func NewBroker() *Broker {
	return &Broker{}
}

func (b *Broker) SetTaskDispatcher(dispatcher broker.Taskdispatcher) {
	b.dispatcher = dispatcher

}

func (b *Broker) Publish(ctx context.Context, task *task.Task) error {

	if b.dispatcher == nil {
		return errors.New("dispatcher is not assigned in memory-mode")
	}

	return b.dispatcher.AddTask(task)
}
