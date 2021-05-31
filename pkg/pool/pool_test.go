package pool_test

import (
	"fmt"
	"testing"
	"time"

	pool_ "github.com/kaydxh/golang/pkg/pool"
)

func TestProcessOk(t *testing.T) {

	p := pool_.New(5, func(task interface{}) error {
		fmt.Printf("proecss task[%v]\n", task)
		time.Sleep(10)
		return nil //must return nil
	})

	defer p.Wait()
	for i := 0; i < 100; i++ {
		p.Put(fmt.Sprintf("%d", i))
	}
	p.Stop()
}

/*
func TestProcessFail(t *testing.T) {

	p := pool_.New(5, func(task interface{}) error {
		fmt.Printf("proecss task[%v]\n", task)
		time.Sleep(10)
		return fmt.Errorf("failedt to process task: %v\n", task)
	})

	defer p.Wait()
	for i := 0; i < 100; i++ {
		p.Put(fmt.Sprintf("%d", i))
	}
	p.Stop()
}
*/
