package apiload

import (
	"context"
	"sync/atomic"
	"testing"
	"time"
)

// TestPercentile_NearestRank 验 nearest-rank 分位实现。
// 样本 1..100(升序)：p50→50, p90→90, p99→99, p100→100。
func TestPercentile_NearestRank(t *testing.T) {
	sorted := make([]int64, 100)
	for i := range sorted {
		sorted[i] = int64(i + 1)
	}
	cases := []struct {
		p    int
		want int64
	}{{50, 50}, {90, 90}, {99, 99}, {100, 100}, {1, 1}}
	for _, c := range cases {
		if got := percentile(sorted, c.p); got != c.want {
			t.Errorf("percentile(p%d) = %d, want %d", c.p, got, c.want)
		}
	}
}

func TestPercentile_Empty(t *testing.T) {
	if got := percentile(nil, 90); got != 0 {
		t.Errorf("percentile(nil) = %d, want 0", got)
	}
}

// TestComputeGroup_Aggregation 验聚合:计数/错误率/min/max/mean。
func TestCollector_Compute(t *testing.T) {
	c := &Collector{}
	// 10 条:8 成功 2 失败,耗时 10,20,...,100。
	for i := 1; i <= 10; i++ {
		c.Add(Sample{Name: "X", DurationMs: int64(i * 10), OK: i > 2, Code: 200})
	}
	m := c.Compute(1.0) // wall=1s → QPS=count
	x := m["X"]
	if x.Count != 10 {
		t.Fatalf("Count = %d, want 10", x.Count)
	}
	if x.ErrCount != 2 || x.OKCount != 8 {
		t.Fatalf("OK/Err = %d/%d, want 8/2", x.OKCount, x.ErrCount)
	}
	if x.ErrRate != 0.2 {
		t.Errorf("ErrRate = %v, want 0.2", x.ErrRate)
	}
	if x.MinMs != 10 || x.MaxMs != 100 {
		t.Errorf("Min/Max = %d/%d, want 10/100", x.MinMs, x.MaxMs)
	}
	if x.MeanMs != 55 {
		t.Errorf("Mean = %v, want 55", x.MeanMs)
	}
	if x.QPS != 10 {
		t.Errorf("QPS = %v, want 10", x.QPS)
	}
	// __all__ 汇总应与单接口一致(仅一个接口)。
	if all := m["__all__"]; all.Count != 10 {
		t.Errorf("__all__ Count = %d, want 10", all.Count)
	}
}

// TestRun_Concurrency 验并发模式:短 duration 内有样本且 duration 边界收敛。
func TestRun_Concurrency(t *testing.T) {
	var calls int64
	task := func(ctx context.Context, workerID, seq int) Sample {
		atomic.AddInt64(&calls, 1)
		time.Sleep(2 * time.Millisecond)
		return Sample{Name: "T", DurationMs: 2, OK: true, Code: 200}
	}
	start := time.Now()
	res := Run(context.Background(), Options{
		Mode: ModeConcurrency, Workers: 4, Duration: 100 * time.Millisecond,
	}, []Task{task})
	elapsed := time.Since(start)

	if res.Collector.Total() == 0 {
		t.Fatal("no samples collected")
	}
	if int64(res.Collector.Total()) != atomic.LoadInt64(&calls) {
		t.Errorf("Total %d != calls %d", res.Collector.Total(), calls)
	}
	// duration 到点应收敛(留宽松上限,含 in-flight 收尾)。
	if elapsed > 500*time.Millisecond {
		t.Errorf("did not converge near duration: %v", elapsed)
	}
}

// TestRun_Rate 验限速模式:样本数量级 ~ QPS*Duration(宽松,含调度抖动)。
func TestRun_Rate(t *testing.T) {
	task := func(ctx context.Context, workerID, seq int) Sample {
		return Sample{Name: "R", DurationMs: 1, OK: true, Code: 200}
	}
	res := Run(context.Background(), Options{
		Mode: ModeRate, Workers: 4, QPS: 100, Duration: 200 * time.Millisecond,
	}, []Task{task})

	total := res.Collector.Total()
	// 理论 ~20 (100qps*0.2s)。宽松区间 [5,50] 容调度抖动。
	if total < 5 || total > 50 {
		t.Errorf("rate samples = %d, want ~20 (range 5-50)", total)
	}
}

// TestRun_Guard 验空 tasks / 零 worker / 零 duration 直接返回空结果不 panic。
func TestRun_Guard(t *testing.T) {
	res := Run(context.Background(), Options{Mode: ModeConcurrency, Workers: 0, Duration: time.Second}, nil)
	if res.Collector.Total() != 0 {
		t.Errorf("expected empty, got %d", res.Collector.Total())
	}
}

// TestRun_TotalConcurrency 验并发模式 Total 限制:恰好执行 Total 个请求(无 duration)。
func TestRun_TotalConcurrency(t *testing.T) {
	task := func(ctx context.Context, workerID, seq int) Sample {
		return Sample{Name: "T", DurationMs: 1, OK: true, Code: 200}
	}
	res := Run(context.Background(), Options{
		Mode: ModeConcurrency, Workers: 8, Total: 500, // 无 Duration,仅按总数停
	}, []Task{task})
	if got := res.Collector.Total(); got != 500 {
		t.Errorf("Total-limited concurrency = %d, want exactly 500", got)
	}
}

// TestRun_TotalRate 验限速模式 Total 限制:发够 Total 即停(允许 in-flight 略溢,不少于 Total)。
func TestRun_TotalRate(t *testing.T) {
	task := func(ctx context.Context, workerID, seq int) Sample {
		return Sample{Name: "R", DurationMs: 1, OK: true, Code: 200}
	}
	res := Run(context.Background(), Options{
		Mode: ModeRate, Workers: 4, QPS: 1000, Total: 50,
	}, []Task{task})
	got := res.Collector.Total()
	// rate 下成功投递够 50 即停;宽松容 in-flight 边界 [50,60]。
	if got < 50 || got > 60 {
		t.Errorf("Total-limited rate = %d, want ~50 (range 50-60)", got)
	}
}
