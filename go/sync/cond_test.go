package sync_test

import (
	"fmt"
	"testing"
	"time"

	"sync"

	sync_ "github.com/kaydxh/golang/go/sync"
	"github.com/stretchr/testify/assert"
)

func TestSignal(t *testing.T) {
	assert := assert.New(t)

	cond := sync_.NewCond()
	a := 5

	go func() {
		for {
			time.Sleep(time.Second)
			cond.SignalDo(func() error {
				a--
				fmt.Printf("a: %v\n", a)
				return nil
			})
		}
	}()

	cond.WaitCondDo(func() bool {
		return a == 3
	}, func() error {
		a += 100
		fmt.Printf("a: %v\n", a)
		assert.Equal(103, a)
		return nil
	})

}

func TestBroadCast(t *testing.T) {
	assert := assert.New(t)

	cond := sync_.NewCond()

	testCases := []struct {
		threads  int
		init     int
		expected int
	}{
		{
			threads:  5,
			init:     10,
			expected: 3,
		},
	}

	var wg sync.WaitGroup

	for i, testCase := range testCases {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cond.WaitCondDo(func() bool {
				return testCase.expected == testCases[i].init
			}, func() error {
				assert.Equal(testCases[i].init, testCase.expected)
				return nil
			})
		}()
	}

	go func() {
		for {
			time.Sleep(time.Second)
			cond.BroadcastDo(func() error {
				testCases[0].init--
				fmt.Printf("init: %v\n", testCases[0].init)
				return nil
			})
		}
	}()

	wg.Wait()
}
