package scheduler

import (
	"context"
	"fmt"

	"github.com/kaydxh/golang/pkg/scheduler/task"
)

type Worker struct {
	id        int
	stopCh    chan struct{}
	working   bool
	worksChCh chan chan *task.Task
	taskCh    chan *task.Task
}

func NewWorker(id int, worksChCh chan chan *task.Task) *Worker {
	w := &Worker{
		id:        id,
		worksChCh: worksChCh,
		taskCh:    make(chan *task.Task),
		stopCh:    make(chan struct{}),
	}

	return w
}

func (w *Worker) doProcess(ctx context.Context, task *task.Task) error {
	fmt.Println(" doProcess")
	if task.PreTaskHandler != nil {
		task.PreTaskHandler(task)
	}

	//do task
	_, err := task.Run()
	if err != nil {
		return err
	}

	if task.PostTaskHandler != nil {
		task.PostTaskHandler(task)
	}

	return nil
}

func (w *Worker) Process(ctx context.Context) {
	w.worksChCh <- w.taskCh

	fmt.Println(" Process task")
	for {
		select {
		case task := <-w.taskCh:
			w.working = true
			fmt.Println(" get task")
			w.doProcess(ctx, task)
			w.worksChCh <- w.taskCh

		case <-w.stopCh:
			return
		}
	}
}

func (w *Worker) Stop() {
	if !w.working {
		return
	}

	close(w.stopCh)
	w.working = false
}
