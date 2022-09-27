/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package scheduler

import (
	"context"
	"fmt"
	"sync"

	"github.com/kaydxh/golang/pkg/scheduler/task"
	task_ "github.com/kaydxh/golang/pkg/scheduler/task"
)

type Dispatcher struct {
	//	tasks   *sync.Map
	workersChCh chan chan *task.Task
	taskQueCh   chan *task.Task
	workers     []*Worker
	burst       int
	stopCh      chan struct{}
}

func NewDispatcher(burst int) *Dispatcher {
	d := &Dispatcher{
		workersChCh: make(chan chan *task.Task, burst),
		taskQueCh:   make(chan *task.Task, 100),
		burst:       burst,
		stopCh:      make(chan struct{}),
	}
	go d.run()
	return d
}

func (d *Dispatcher) run() {
	for i := 0; i < d.burst; i++ {
		worker := NewWorker(i, d.workersChCh)
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
