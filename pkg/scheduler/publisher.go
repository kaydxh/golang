package scheduler

import (
	"encoding/json"
	"fmt"

	"github.com/kaydxh/golang/pkg/scheduler/task"
)

type publisher struct {
}

func (s *publisher) Publish(task *task.Task) error {

	msg, err := json.Marshal(task)
	if err != nil {
		return fmt.Errorf("JSON marshal error: %s", err)
	}

	//send msg to broker
	_ = msg
	return nil
}
