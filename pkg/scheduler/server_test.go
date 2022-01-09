package scheduler_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/google/uuid"
	scheduler_ "github.com/kaydxh/golang/pkg/scheduler"
	memory_ "github.com/kaydxh/golang/pkg/scheduler/broker/memory"
	task_ "github.com/kaydxh/golang/pkg/scheduler/task"
)

func TestServer(t *testing.T) {
	burst := 1
	broker := memory_.NewBroker()
	server := scheduler_.NewServer(burst, broker)

	server.AddTask(context.Background(), &task_.Task{
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
				Type:  "string",
				Value: "hello",
			},
		},
	})

	select {}

}
