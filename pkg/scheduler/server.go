package scheduler

import (
	"github.com/kaydxh/golang/pkg/scheduler/broker"
	task_ "github.com/kaydxh/golang/pkg/scheduler/task"
	"golang.org/x/net/context"
)

type Server struct {
	broker broker.Broker
	burst  int
}

func NewServer(burst int, broker broker.Broker) *Server {
	s := &Server{
		broker: broker,
	}
	s.broker.SetTaskDispatcher(NewDispatcher(burst))

	return s
}

func (s *Server) AddTask(ctx context.Context, task *task_.Task) error {
	err := task_.Validate(task)
	if err != nil {
		return err
	}

	return s.broker.Publish(ctx, task)

}
