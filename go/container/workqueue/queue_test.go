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
package workqueue_test

import (
	"sync"
	"testing"
	"time"

	workqueue_ "github.com/kaydxh/golang/go/container/workqueue"
)

func TestBasic(t *testing.T) {
	testCases := []struct {
		queue         *workqueue_.Type
		queueShutDown func(workqueue_.Interface)
	}{
		{
			queue:         workqueue_.NewQueue(),
			queueShutDown: workqueue_.Interface.ShutDown,
		},
	}

	for _, testCase := range testCases {

		// Start producers
		const producers = 10
		producerWG := sync.WaitGroup{}
		producerWG.Add(producers)
		for i := 0; i < producers; i++ {
			go func(i int) {
				defer producerWG.Done()
				for j := 0; j < 2; j++ {
					testCase.queue.Add(i)
					time.Sleep(time.Millisecond)
				}
			}(i)
		}

		// Start consumers
		const consumers = 10
		consumerWG := sync.WaitGroup{}
		consumerWG.Add(consumers)
		for i := 0; i < consumers; i++ {
			go func(i int) {
				defer consumerWG.Done()
				for {
					item, quit := testCase.queue.Get()
					if quit {
						return
					}

					t.Logf("Woker %v: begin processing %v", i, item)
					time.Sleep(3 * time.Millisecond)
					t.Logf("Woker %v: done processing %v", i, item)
					testCase.queue.Done(item)

				}
			}(i)
		}

		producerWG.Wait()
		testCase.queueShutDown(testCase.queue)
		testCase.queue.Add("add after shutdown!")
		consumerWG.Wait()

		if testCase.queue.Len() != 0 {
			t.Errorf("Expected the queue to be empty, had: %v items", testCase.queue.Len())
		}
	}

}

func TestReinsert(t *testing.T) {
	q := workqueue_.NewQueue()
	q.Add("foo")

	// Start processing
	i, _ := q.Get()
	if i != "foo" {
		t.Errorf("Expected %v, got %v", "foo", i)
	}

	// Add it back while processing
	q.Add(i)

	// Finish it up
	q.Done(i)

	// It should be back on the queue
	i, _ = q.Get()
	if i != "foo" {
		t.Errorf("Expected %v, got %v", "foo", i)
	}

	// Finish that one up
	q.Done(i)

	if a := q.Len(); a != 0 {
		t.Errorf("Expected queue to be empty. Has %v items", a)
	}
}
