// Package apiload 提供接口压测的采样、延迟分位聚合与控制台输出，通用无业务。
// 与 apiassert/apireport 同级:apiassert 断字段、apireport 报功能结果、apiload 压负载。
package apiload

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Sample 单次请求采样。Name 为接口名(聚合维度)，DurationMs 端到端耗时，OK 是否成功，Code HTTP 码。
type Sample struct {
	Name       string
	DurationMs int64
	OK         bool
	Code       int
}

// Metrics 单接口(或全局)聚合指标。分位单位毫秒。
type Metrics struct {
	Name     string
	Count    int64
	OKCount  int64
	ErrCount int64
	ErrRate  float64 // ErrCount/Count
	QPS      float64 // Count/wallSeconds
	MinMs    int64
	MaxMs    int64
	MeanMs   float64
	P50Ms    int64
	P90Ms    int64
	P99Ms    int64
}

// Collector 并发安全采样收集器。多 worker 并发 Add，收尾 Compute 聚合。
type Collector struct {
	mu      sync.Mutex
	samples []Sample
}

// Add 记录一次采样(并发安全)。
func (c *Collector) Add(s Sample) {
	c.mu.Lock()
	c.samples = append(c.samples, s)
	c.mu.Unlock()
}

// Total 返回累计采样数(并发安全)。
func (c *Collector) Total() int {
	c.mu.Lock()
	defer c.mu.Unlock()
	return len(c.samples)
}

// Snap 运行期轻量快照(不排序,仅计数),供实时进度用。
type Snap struct {
	Count    int64
	OKCount  int64
	ErrCount int64
}

// Snapshot 返回当前累计计数快照(并发安全,O(n) 遍历但不排序)。
func (c *Collector) Snapshot() Snap {
	c.mu.Lock()
	defer c.mu.Unlock()
	var s Snap
	s.Count = int64(len(c.samples))
	for i := range c.samples {
		if c.samples[i].OK {
			s.OKCount++
		} else {
			s.ErrCount++
		}
	}
	return s
}

// Compute 按接口名聚合分位。wallSeconds 为压测实际墙钟秒数(算 QPS)。
// 返回 map[接口名]Metrics 外加键 "__all__" 的全局汇总。
func (c *Collector) Compute(wallSeconds float64) map[string]Metrics {
	c.mu.Lock()
	snapshot := make([]Sample, len(c.samples))
	copy(snapshot, c.samples)
	c.mu.Unlock()

	byName := map[string][]Sample{}
	for _, s := range snapshot {
		byName[s.Name] = append(byName[s.Name], s)
		byName["__all__"] = append(byName["__all__"], s)
	}

	out := make(map[string]Metrics, len(byName))
	for name, group := range byName {
		out[name] = computeGroup(name, group, wallSeconds)
	}
	return out
}

// computeGroup 聚合单组采样为 Metrics。
func computeGroup(name string, group []Sample, wallSeconds float64) Metrics {
	m := Metrics{Name: name, Count: int64(len(group))}
	if m.Count == 0 {
		return m
	}

	durs := make([]int64, 0, len(group))
	var sum int64
	for _, s := range group {
		if s.OK {
			m.OKCount++
		} else {
			m.ErrCount++
		}
		durs = append(durs, s.DurationMs)
		sum += s.DurationMs
	}
	sort.Slice(durs, func(i, j int) bool { return durs[i] < durs[j] })

	m.MinMs = durs[0]
	m.MaxMs = durs[len(durs)-1]
	m.MeanMs = float64(sum) / float64(m.Count)
	m.P50Ms = percentile(durs, 50)
	m.P90Ms = percentile(durs, 90)
	m.P99Ms = percentile(durs, 99)
	m.ErrRate = float64(m.ErrCount) / float64(m.Count)
	if wallSeconds > 0 {
		m.QPS = float64(m.Count) / wallSeconds
	}
	return m
}

// percentile 对已升序切片取 nearest-rank 分位(p 为 0~100)。空切片返回 0。
func percentile(sorted []int64, p int) int64 {
	n := len(sorted)
	if n == 0 {
		return 0
	}
	// nearest-rank: rank = ceil(p/100 * n)，索引 rank-1，收敛到 [0, n-1]。
	rank := (p*n + 99) / 100 // ceil(p*n/100) 整数实现
	if rank < 1 {
		rank = 1
	}
	if rank > n {
		rank = n
	}
	return sorted[rank-1]
}

// PrintConsole 打印每接口 + 全局汇总(QPS/错误率/延迟分位)。
// 使用 displayWidth 手动对齐:%-Ns 是按 UTF-8 字节计算填充,CJK 字符(每个 3 字节)会造成列错位。
func (c *Collector) PrintConsole(metrics map[string]Metrics) {
	const (
		red   = "\033[0;31m"
		green = "\033[0;32m"
		yel   = "\033[1;33m"
		nc    = "\033[0m"
	)

	// 稳定输出:接口名排序,__all__ 置底。
	names := make([]string, 0, len(metrics))
	for n := range metrics {
		if n != "__all__" {
			names = append(names, n)
		}
	}
	sort.Strings(names)

	// 列宽:接口名 32(左对齐),其余 8(右对齐)。
	const nameW, numW = 32, 8
	headers := []string{"请求数", "QPS", "错误率", "p50ms", "p90ms", "p99ms", "maxms"}
	fmt.Print(padRight("接口", nameW))
	for _, h := range headers {
		fmt.Print(" " + padLeft(h, numW))
	}
	fmt.Println()

	printRow := func(m Metrics) {
		color := green
		if m.ErrRate > 0 {
			color = yel
		}
		if m.ErrRate >= 0.5 {
			color = red
		}
		errStr := fmt.Sprintf("%.2f%%", m.ErrRate*100)
		fmt.Print(padRight(truncName(m.Name, nameW), nameW))
		fmt.Print(" " + padLeft(fmt.Sprintf("%d", m.Count), numW))
		fmt.Print(" " + padLeft(fmt.Sprintf("%.1f", m.QPS), numW))
		fmt.Print(" " + color + padLeft(errStr, numW) + nc)
		fmt.Print(" " + padLeft(fmt.Sprintf("%d", m.P50Ms), numW))
		fmt.Print(" " + padLeft(fmt.Sprintf("%d", m.P90Ms), numW))
		fmt.Print(" " + padLeft(fmt.Sprintf("%d", m.P99Ms), numW))
		fmt.Print(" " + padLeft(fmt.Sprintf("%d", m.MaxMs), numW))
		fmt.Println()
	}
	for _, n := range names {
		printRow(metrics[n])
	}
	if all, ok := metrics["__all__"]; ok {
		fmt.Println("────────────────────────────────────────")
		all.Name = "全局汇总"
		printRow(all)
	}
}

// displayWidth 返回字符串在等宽终端下的显示列数。ASCII 1 列,CJK/全角 2 列。
// 不处理零宽 combining marks(测试报告场景足够)。
func displayWidth(s string) int {
	w := 0
	for _, r := range s {
		w += runeWidth(r)
	}
	return w
}

func runeWidth(r rune) int {
	switch {
	case r < 0x7F:
		return 1
	case r >= 0x1100 && r <= 0x115F,
		r >= 0x2E80 && r <= 0x303E,
		r >= 0x3041 && r <= 0x33FF,
		r >= 0x3400 && r <= 0x4DBF,
		r >= 0x4E00 && r <= 0x9FFF,
		r >= 0xA000 && r <= 0xA4CF,
		r >= 0xAC00 && r <= 0xD7A3,
		r >= 0xF900 && r <= 0xFAFF,
		r >= 0xFE30 && r <= 0xFE4F,
		r >= 0xFF00 && r <= 0xFF60,
		r >= 0xFFE0 && r <= 0xFFE6,
		r >= 0x20000 && r <= 0x2FFFD,
		r >= 0x30000 && r <= 0x3FFFD:
		return 2
	}
	return 1
}

func padRight(s string, width int) string {
	if w := displayWidth(s); w < width {
		return s + strings.Repeat(" ", width-w)
	}
	return s
}

func padLeft(s string, width int) string {
	if w := displayWidth(s); w < width {
		return strings.Repeat(" ", width-w) + s
	}
	return s
}

// truncName 按显示列宽截断名字,尾部加 …(1 列)。ASCII/CJK 混合安全。
func truncName(s string, w int) string {
	if displayWidth(s) <= w {
		return s
	}
	used := 0
	var b strings.Builder
	for _, r := range s {
		rw := runeWidth(r)
		if used+rw > w-1 {
			break
		}
		b.WriteRune(r)
		used += rw
	}
	b.WriteRune('…')
	return b.String()
}
