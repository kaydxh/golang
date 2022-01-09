package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/kaydxh/golang/pkg/scheduler/broker"
	"github.com/kaydxh/golang/pkg/scheduler/task"
	task_ "github.com/kaydxh/golang/pkg/scheduler/task"
)

type Dispatcher struct {
	//	tasks   *sync.Map
	workersChCh chan chan *task.Task
	taskQueCh   chan *task.Task
	workers     []*Worker
	burst       int
	broker      broker.Broker
	stopCh      chan struct{}
}

func NewDispatcher(burst int, broker broker.Broker) *Dispatcher {
	d := &Dispatcher{
		workersChCh: make(chan chan *task.Task, burst),
		taskQueCh:   make(chan *task.Task, 100),
		burst:       burst,
		broker:      broker,
		stopCh:      make(chan struct{}),
	}
	go d.run()
	return d
}

func (d *Dispatcher) run() {
	for i := 0; i < d.burst; i++ {
		worker := NewWorker(i, d.workersChCh, d.broker)
		go worker.Process(context.Background())
		d.workers = append(d.workers, worker)

		go d.dispatch()
	}
}

func (d *Dispatcher) dispatch() {

	for {
		select {
		case task := <-d.taskQueCh:
			taskCh := <-d.workersChCh
			taskCh <- task

		case <-d.stopCh:
			d.stopWokrs()
			return
		}

	}
}

func (d *Dispatcher) stopWokrs() {
	wg := new(sync.WaitGroup)
	wg.Add(len(d.workers))
	for _, v := range d.workers {
		go func(w *Worker, wg *sync.WaitGroup) {
			defer wg.Done()
			if w.working {
				w.Stop()
			}
		}(v, wg)

	}

	wg.Wait()
}

func (d *Dispatcher) AddTask(task *task_.Task) error {
	err := task_.Validate(task)
	if err != nil {
		return err
	}
	d.taskQueCh <- task
	fmt.Println("add task")

	return nil
}

func (d *Dispatcher) Stop() {
	close(d.stopCh)
}
