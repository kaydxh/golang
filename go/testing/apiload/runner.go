package apiload

import (
	"context"
	"fmt"
	"os"
	"sync"
	"sync/atomic"
	"time"
)

// Mode 负载模型。
type Mode int

const (
	// ModeConcurrency 固定并发:Workers 个 goroutine 各自 loop 满打,压系统最大吞吐。
	ModeConcurrency Mode = iota
	// ModeRate 固定 QPS:按 1/QPS 间隔投递,worker pool 消费,背压时记 dropped。
	ModeRate
)

// Task 单次压测动作:发一次请求并返回采样。
// workerID 标识发起 goroutine，seq 为该 worker 内递增序号(供写路径生成唯一入参)。
type Task func(ctx context.Context, workerID, seq int) Sample

// Options 压测参数。
type Options struct {
	Mode     Mode
	Workers  int           // 并发 worker 数(两种模式都需要)
	QPS      int           // ModeRate 目标 QPS(ModeConcurrency 忽略)
	Duration time.Duration // 压测时长(与 Total 谁先到谁停;Total>0 时可为 0)
	Total    int64         // 总请求数上限(0=不限,仅按 Duration)
	// Progress>0 时,每隔该间隔向 stderr 打印一行实时进度(已发/QPS/错误率)。0=不打印。
	Progress time.Duration
}

// Result 压测结果:采样收集器 + 实际墙钟 + rate 模式丢弃计数。
type Result struct {
	Collector   *Collector
	WallSeconds float64
	Dropped     int64 // ModeRate 下 worker 全忙、tick 无法投递的次数
}

// Run 执行压测。tasks 轮转选取(worker 每次迭代取 tasks[seq%len])。
// 停止条件:Duration 到点 或 已发够 Total 个请求(谁先到),cancel 后等 in-flight 收敛。
func Run(ctx context.Context, opts Options, tasks []Task) *Result {
	col := &Collector{}
	if len(tasks) == 0 || opts.Workers <= 0 || (opts.Duration <= 0 && opts.Total <= 0) {
		return &Result{Collector: col}
	}

	// Duration>0 → 带超时;仅 Total 限制时用可取消 ctx(由计数触发 cancel)。
	var runCtx context.Context
	var cancel context.CancelFunc
	if opts.Duration > 0 {
		runCtx, cancel = context.WithTimeout(ctx, opts.Duration)
	} else {
		runCtx, cancel = context.WithCancel(ctx)
	}
	defer cancel()

	start := time.Now()

	// 实时进度:每 Progress 间隔打印已发/QPS/错误率到 stderr,压测结束停。
	var progDone chan struct{}
	if opts.Progress > 0 {
		progDone = make(chan struct{})
		go runProgress(runCtx, opts, col, start, progDone)
	}

	var dropped int64
	switch opts.Mode {
	case ModeRate:
		dropped = runRate(runCtx, cancel, opts, tasks, col)
	default:
		runConcurrency(runCtx, cancel, opts, tasks, col)
	}
	wall := time.Since(start).Seconds()

	if progDone != nil {
		<-progDone // 等进度协程打完收尾行
	}

	return &Result{Collector: col, WallSeconds: wall, Dropped: dropped}
}

// runProgress 定时向 stderr 打印实时进度(\r 原地刷新),ctx 结束时打最终行 + 换行。
func runProgress(ctx context.Context, opts Options, col *Collector, start time.Time, done chan struct{}) {
	defer close(done)
	ticker := time.NewTicker(opts.Progress)
	defer ticker.Stop()

	print := func() {
		snap := col.Snapshot()
		elapsed := time.Since(start).Seconds()
		qps := 0.0
		if elapsed > 0 {
			qps = float64(snap.Count) / elapsed
		}
		errRate := 0.0
		if snap.Count > 0 {
			errRate = float64(snap.ErrCount) / float64(snap.Count) * 100
		}
		suffix := ""
		if opts.Total > 0 {
			suffix = fmt.Sprintf("/%d", opts.Total)
		}
		fmt.Fprintf(os.Stderr, "\r  进度: 已发 %d%s | %.1fs | QPS %.1f | 错误率 %.2f%%   ",
			snap.Count, suffix, elapsed, qps, errRate)
	}

	for {
		select {
		case <-ctx.Done():
			print()
			fmt.Fprintln(os.Stderr) // 收尾换行,避免覆盖后续汇总
			return
		case <-ticker.C:
			print()
		}
	}
}

// runConcurrency 起 Workers 个 goroutine，各自 loop 轮流跑 tasks 直到 ctx 到期或发够 Total。
// Total>0 时用共享原子计数器抢占配额:抢到序号 >Total 即停,保证恰好执行 Total 个请求。
func runConcurrency(ctx context.Context, cancel context.CancelFunc, opts Options, tasks []Task, col *Collector) {
	var counter int64
	var wg sync.WaitGroup
	for w := 0; w < opts.Workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			seq := 0
			for {
				if ctx.Err() != nil {
					return
				}
				if opts.Total > 0 {
					// 抢配额:第 n 个(从 1 计)。超出 Total 则触发收尾并退出。
					if n := atomic.AddInt64(&counter, 1); n > opts.Total {
						cancel()
						return
					}
				}
				task := tasks[seq%len(tasks)]
				col.Add(task(ctx, workerID, seq))
				seq++
			}
		}(w)
	}
	wg.Wait()
}

// runRate 固定 QPS:ticker 按 1/QPS 投递到 job 通道，Workers 个消费者执行。
// 通道满(worker 全忙)时丢弃该 tick 并计数，避免阻塞节流器。返回 dropped 计数。
// Total>0 时:成功投递够 Total 个即停止投递(丢弃的不计入 Total)。
func runRate(ctx context.Context, cancel context.CancelFunc, opts Options, tasks []Task, col *Collector) int64 {
	qps := opts.QPS
	if qps <= 0 {
		qps = 1
	}
	interval := time.Second / time.Duration(qps)

	type job struct{ seq int }
	jobs := make(chan job, opts.Workers) // 缓冲=worker 数,背压即丢
	var dropped int64
	var seq int64
	var sent int64

	var wg sync.WaitGroup
	for w := 0; w < opts.Workers; w++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			for j := range jobs {
				col.Add(tasks[j.seq%len(tasks)](ctx, workerID, j.seq))
			}
		}(w)
	}

	ticker := time.NewTicker(interval)
	defer ticker.Stop()
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		case <-ticker.C:
			s := int(atomic.AddInt64(&seq, 1) - 1)
			select {
			case jobs <- job{seq: s}:
				if opts.Total > 0 && atomic.AddInt64(&sent, 1) >= opts.Total {
					cancel() // 发够 Total,停止投递
					break loop
				}
			default:
				atomic.AddInt64(&dropped, 1) // worker 全忙,丢弃本 tick
			}
		}
	}
	close(jobs)
	wg.Wait()
	return dropped
}
