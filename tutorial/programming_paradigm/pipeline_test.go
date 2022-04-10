package tutorial

import (
	"fmt"
	"math"
	"sync"
	"testing"
)

//通过channel 与 gorutinue 实现Fan in/Out, 1对多/多对1
//Multiple functions can read from the same channel until that channel
//is closed; this is called fan-out
//https://go.dev/blog/pipelines

// echo 将整数数组放入到channel中，并返回 channel
func echo(nums []int) <-chan int {
	out := make(chan int)
	go func() {
		for _, n := range nums {
			out <- n
		}
		close(out)
	}()
	return out
}

// sum内部使用一个协程处理从in读取的数据, 如果多次调用sum，就可以实现将
// 同一个chan 中的数据分成多个协程进行并发处理(读取数据这块是串行，单处理是并发)
func sum(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		var sum = 0
		for n := range in {
			sum += n
		}
		out <- sum
		//close 作用的是让外部读取out chanel停止调，不会会永远阻塞
		close(out)
	}()
	return out
}

func is_prime(value int) bool {
	for i := 2; i <= int(math.Floor(float64(value)/2)); i++ {
		if value%i == 0 {
			return false
		}
	}
	return value > 1
}

func prime(in <-chan int) <-chan int {
	out := make(chan int)
	go func() {
		for n := range in {
			if is_prime(n) {
				out <- n
			}
		}
		close(out)
	}()
	return out
}

//
func merge(cs []<-chan int) <-chan int {
	var wg sync.WaitGroup
	out := make(chan int)

	wg.Add(len(cs))
	for _, c := range cs {
		go func(c <-chan int) {
			for n := range c {
				out <- n
			}
			wg.Done()
		}(c)
	}
	//协程的目的是让外部读取out的函数先工作起来，等到所有协程都完成工作
	//后，close调，中断外部读取操作
	go func() {
		wg.Wait()
		// 关闭读取前，等待所有的协程完成任务
		close(out)
	}()
	return out
}

//生成指定范围的数组
func makeRange(min, max int) []int {
	a := make([]int, max-min+1)
	for i := range a {
		a[i] = min + i
	}
	return a
}

func TestPipleline(t *testing.T) {
	nums := makeRange(1, 10000)
	// echo 将整数数组放入到channel中，并返回 channel
	in := echo(nums)

	const nProcess = 5
	var chans [nProcess]<-chan int
	for i := range chans {

		// 将chan in中的数据分成5个协程处理, 结果输出到5个chan中
		chans[i] = sum(prime(in))
	}

	/*
		n := <-sum(merge(chans[:]))
		fmt.Println(n)
	*/
	//如果有多个结果的话，使用range读取结果
	for n := range sum(merge(chans[:])) {
		fmt.Println(n)
	}

}
