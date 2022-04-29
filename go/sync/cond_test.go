package sync_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	sync_ "github.com/kaydxh/golang/go/sync"
	"github.com/stretchr/testify/assert"
)

func TestWaitForDo(t *testing.T) {
	assert := assert.New(t)

	l := new(sync.Mutex)
	cond := sync_.NewCond(l)
	a := 5
	timout := 4 * time.Second

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		err := cond.WaitForDo(timout, func() bool {
			fmt.Printf("==>a = %v\n", a)
			return a == 3
		}, func() error {
			a += 100
			fmt.Printf("a: %v\n", a)
			assert.Equal(103, a)
			return nil
		})
		assert.Equal(nil, err)
	}()

	go func() {
		for {
			cond.SignalDo(func() error {
				a--
				//this sleep can increase the probability to waitforDo
				//miss condition = > a==3
				time.Sleep(1 * time.Second)
				fmt.Printf("a: %v\n", a)
				return nil
			})
			//time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

// condion only be fit for range condition
func TestWaitUntilDo(t *testing.T) {
	assert := assert.New(t)

	l := new(sync.Mutex)
	cond := sync_.NewCond(l)
	a := 5

	var wg sync.WaitGroup
	wg.Add(1)
	go func() {
		defer wg.Done()
		cond.WaitUntilDo(func() bool {
			fmt.Printf("==>a = %v\n", a)
			return a < 0
		}, func() error {
			fmt.Printf("a: %v\n", a)
			assert.LessOrEqual(a, 0)
			return nil
		})
	}()

	go func() {
		for {
			cond.SignalDo(func() error {
				a--
				//this sleep can increase the probability to waitforDo
				//miss condition = > a==3
				time.Sleep(1 * time.Second)
				fmt.Printf("a: %v\n", a)
				return nil
			})
			//time.Sleep(1 * time.Second)
		}
	}()

	wg.Wait()
}

func TestBroadCast(t *testing.T) {
	assert := assert.New(t)

	l := new(sync.Mutex)
	cond := sync_.NewCond(l)

	initValue := 10
	expected := 3
	threads := 5

	var wg sync.WaitGroup

	for i := 0; i < threads; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			cond.WaitForDo(10*time.Second, func() bool {
				return expected == initValue
			}, func() error {
				assert.Equal(initValue, expected)
				t.Logf("wait done")
				return nil
			})
		}()
	}

	go func() {
		for {
			time.Sleep(time.Second)
			cond.BroadcastDo(func() error {
				initValue--
				t.Logf("init: %v\n", initValue)
				return nil
			})
		}
	}()

	wg.Wait()
}
