package scheduler_test

import (
	"fmt"
	"testing"

	"github.com/google/uuid"
	scheduler_ "github.com/kaydxh/golang/pkg/scheduler"
	memory_ "github.com/kaydxh/golang/pkg/scheduler/broker/memory"
	task_ "github.com/kaydxh/golang/pkg/scheduler/task"
)

func TestDispatcher(t *testing.T) {
	burst := 1
	broker := memory_.NewBroker()
	dispatcher := scheduler_.NewDispatcher(burst, broker)
	dispatcher.AddTask(&task_.Task{
		TaskId: uuid.New().String(),
		TaskFunc: func(a int, b string) error {
			fmt.Printf("do task: %v:%v\n", a, b)
			return nil
		},
		Args: []task_.TaskArgument{
			task_.TaskArgument{
				Type:  "int",
				Value: 8,
			},
			task_.TaskArgument{
				Type:  "int",
				Value: "hello",
			},
		},
	})
	dispatcher.Stop()
	fmt.Println("out")

}
