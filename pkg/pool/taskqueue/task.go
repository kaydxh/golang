package taskqueue

import "sync"

type TaskOptions struct {
	Name string
}

type Task struct {
	opts TaskOptions
}

type TaskerMap sync.Map

type Tasker interface {
	TaskHandler(task interface{}) error
	Scheme() string
}
