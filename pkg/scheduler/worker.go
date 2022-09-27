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
