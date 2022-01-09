package memory

import (
	"context"

	"github.com/kaydxh/golang/pkg/scheduler/broker"
	"github.com/kaydxh/golang/pkg/scheduler/task"
)

type Broker struct {
	worker broker.TaskProcessor
}

func NewBroker() *Broker {
	return &Broker{
		//	worker: w,
	}
}

func (b *Broker) Publish(ctx context.Context, task *task.Task) error {

	if b.worker == nil {
		//	return errors.New("worker is not assigned in memory-mode")
	}

	// faking the behavior to marshal input into json
	// and unmarshal it back
	/*
		message, err := json.Marshal(task)
		if err != nil {
			return fmt.Errorf("JSON marshal error: %s", err)
		}
	*/

	/*
		signature := new(task.Signature)
		var
		decoder := json.NewDecoder(bytes.NewReader(message))
		decoder.UseNumber()
		if err := decoder.Decode(signature); err != nil {
			return fmt.Errorf("JSON unmarshal error: %s", err)
		}
	*/
	_, err := task.Run()
	return err

	// blocking call to the task directly
	//return b.worker.Process(task)

}
