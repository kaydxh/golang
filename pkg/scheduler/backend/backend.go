package backend

import "github.com/kaydxh/golang/pkg/scheduler/task"

type Backend interface {
	SetState(state types.TaskState, task *task.Task) error
	GetState(taskId string) (types.TaskState, error)
}
