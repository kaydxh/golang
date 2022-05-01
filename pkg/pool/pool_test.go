package pool_test

import (
	"fmt"
	"testing"
	"time"

	pool_ "github.com/kaydxh/golang/pkg/pool"
)

func TestProcessOk(t *testing.T) {

	p := pool_.New(5, func(task interface{}) error {
		fmt.Printf("start task[%v]\n", task)
		time.Sleep(time.Second)
		fmt.Printf("finish task[%v]\n", task)
		return nil //must return nil
	})

	//	defer p.Wait()
	for i := 0; i < 10; i++ {
		err := p.Put(fmt.Sprintf("%d", i))
		if err != nil {
			t.Errorf("put work err: %v", err)
			return
		}
		//p.Stop()
	}
	//	p.Stop()
	p.Wait()
	t.Logf("err: %v", p.Error())
}

func TestProcessFail(t *testing.T) {

	p := pool_.New(2, func(task interface{}) error {
		fmt.Printf("start task[%v]\n", task)
		time.Sleep(time.Second)
		fmt.Printf("finish task[%v]\n", task)
		return fmt.Errorf("failed to process task: %v\n", task)
	})

	//defer p.Wait()
	for i := 0; i < 4; i++ {
		err := p.Put(fmt.Sprintf("%d", i))
		if err != nil {
			t.Errorf("put work err: %v", err)
			return
		}
	}
	p.Stop()
	p.Wait()
	t.Logf("err: %v", p.Error())
}
