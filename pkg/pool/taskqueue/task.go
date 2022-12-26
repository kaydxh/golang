package taskqueue

import (
	"sync"

	queue_ "github.com/kaydxh/golang/pkg/pool/taskqueue/queue"
)

type TaskOptions struct {
	Name string
}

type Task struct {
	opts TaskOptions
}

type TaskerMap sync.Map

type Tasker interface {
	TaskHandler(message *queue_.Message) (*queue_.MessageResult, error)
	Scheme() string
}
