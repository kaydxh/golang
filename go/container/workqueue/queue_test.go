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
