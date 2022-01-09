package scheduler

import (
	"context"
	"fmt"

	"github.com/kaydxh/golang/pkg/scheduler/task"
)

type Worker struct {
	id        int
	limit     int // number of tasks which can be executed in parallel
	stopCh    chan struct{}
	working   bool
	worksChCh chan chan *task.Task

	taskCh chan *task.Task
	//	tasksWait *sync.WaitGroup // wait group for waiting all tasks finishing when we close this worker

	preTaskHandler  func(*task.Task)
	postTaskHandler func(*task.Task)
}

func NewWorker(id int, worksChCh chan chan *task.Task) *Worker {
	w := &Worker{
		id:        id,
		worksChCh: worksChCh,
		//		tasksWait: new(sync.WaitGroup),
		taskCh: make(chan *task.Task),
		stopCh: make(chan struct{}),
	}

	return w
}

func (w *Worker) doProcess(ctx context.Context, task *task.Task) error {
	//	defer w.tasksWait.Done()
	fmt.Println(" doProcess")
	/*
		if w.preTaskHandler != nil {
			w.preTaskHandler(task)
		}
	*/

	//do task
	_, err := task.Run()
	if err != nil {
		return err
	}

	if w.postTaskHandler != nil {
		w.postTaskHandler(task)
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
			//w.tasksWait.Add(1)
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

	//w.tasksWait.Wait()
	w.working = false
}
