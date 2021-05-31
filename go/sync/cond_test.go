package sync_test

import (
	"fmt"
	"testing"
	"time"

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
			cond.SignalDo(func() {
				a--
				fmt.Printf("a: %v\n", a)
			})
		}
	}()

	cond.WaitCond(func() bool {
		return a == 3
	}, func() {
		a += 100
		fmt.Printf("a: %v\n", a)
		assert.Equal(103, a)
	})

}
