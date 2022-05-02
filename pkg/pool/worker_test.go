package pool_test

import (
	"context"
	"fmt"
	"testing"
	"time"

	pool_ "github.com/kaydxh/golang/pkg/pool"
)

func TestWork(t *testing.T) {
	p := pool_.Worker{
		Burst: 10,
	}

	workChan := make(chan interface{})
	ctx := context.Background()
	workDoneChan := p.Work(ctx, workChan, func(work interface{}) error {
		fmt.Printf("process task[%v]\n", work)
		time.Sleep(5 * time.Second)
		//return nil
		return fmt.Errorf("error")
	}, false)
	_ = workDoneChan

	works := []int{1, 2, 3}
	for _, w := range works {
		// 这里select workDoneChan和ctx.Done可能无法感知，因为work func执
		// 行比较耗时，而for循环在burst足够多的情况下，立马就执行结束了
		select {
		case workChan <- w:
			// 这里和woker func函数执行顺序是不确定的，workChan <- w返回，
			// 代表work func已经接收数据，开始执行woker处理协程
			t.Logf("start to do work %v", w)

		case <-workDoneChan:
			t.Logf("workdone channel is closed, work goroutine is exit")
			return

		case <-ctx.Done():
			//canceled, get Error messgae
			t.Errorf("err: %v\n", p.Error())
			return

		}

		// stop后，将会在workDoneChan读取到任务协助退出
		//p.Stop()
	}

	// if need wait finish
	p.Wait()
	t.Logf("wait finish, err: %v", p.Error())
}
