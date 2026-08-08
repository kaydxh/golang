// Package apireport 提供接口测试结果的收集、控制台输出与退出码，通用无业务。
package apireport

import (
	"fmt"
	"sort"
)

// Result 单条用例结果。
type Result struct {
	Suite      string
	Name       string
	Path       string
	Code       int
	DurationMs int64
	Status     string // pass/fail/skip
	Reason     string
	CurlRepro  string
	Body       string // 响应体(verbose 模式打印)
	TraceID    string // X-TraceId,便于服务端日志按串反查(打印在 verbose/failure 输出中)
}

// Collector 汇总结果。
type Collector struct {
	Results []Result
	// BodyLimit 响应体打印上限(字节;0=不限)。超过则截断显示,避免长响应/二进制刷屏。
	BodyLimit int
}

func (c *Collector) Add(r Result) { c.Results = append(c.Results, r) }

func (c *Collector) counts() (pass, fail, skip int) {
	for _, r := range c.Results {
		switch r.Status {
		case "pass":
			pass++
		case "fail":
			fail++
		case "skip":
			skip++
		}
	}
	return
}

// ExitCode：有 fail → 1，否则 0。（配置/连通性错误由调用方直接 os.Exit(2)。）
func (c *Collector) ExitCode() int {
	_, fail, _ := c.counts()
	if fail > 0 {
		return 1
	}
	return 0
}

// PrintConsole 打印彩色汇总；失败项打印响应/curl。slowMs>0 时标注慢接口。
// 末尾额外打印接口耗时 Top 汇总(按 suite/name 去重,降序;skip 不计)。
func (c *Collector) PrintConsole(verbose bool, slowMs int64) {
	const (
		red   = "\033[0;31m"
		green = "\033[0;32m"
		yel   = "\033[1;33m"
		nc    = "\033[0m"
	)
	for _, r := range c.Results {
		switch r.Status {
		case "pass":
			fmt.Printf("  %s✓ PASS%s [%s/%s] %dms\n", green, nc, r.Suite, r.Name, r.DurationMs)
			if slowMs > 0 && r.DurationMs > slowMs {
				fmt.Printf("    %s⚠ SLOW%s %dms (>%dms)\n", yel, nc, r.DurationMs, slowMs)
			}
			if verbose {
				if r.TraceID != "" {
					fmt.Printf("    TraceId: %s\n", r.TraceID)
				}
				if r.CurlRepro != "" {
					fmt.Printf("    请求: %s\n", r.CurlRepro)
				}
				if r.Body != "" {
					fmt.Printf("    响应: %s\n", truncateBody(r.Body, c.BodyLimit))
				}
			}
		case "skip":
			fmt.Printf("  %s- SKIP%s [%s/%s] %s\n", yel, nc, r.Suite, r.Name, r.Reason)
		default:
			fmt.Printf("  %s✗ FAIL%s [%s/%s] %dms %s\n", red, nc, r.Suite, r.Name, r.DurationMs, r.Reason)
			if r.TraceID != "" {
				fmt.Printf("    TraceId: %s\n", r.TraceID)
			}
			if r.Body != "" {
				fmt.Printf("    响应: %s\n", truncateBody(r.Body, c.BodyLimit))
			}
			fmt.Printf("    复现: %s\n", r.CurlRepro)
		}
	}
	pass, fail, skip := c.counts()
	fmt.Printf("\n总计 %d | 通过 %d | 失败 %d | 跳过 %d\n", len(c.Results), pass, fail, skip)

	c.printDurationTop(slowMs)
}

// truncateBody 按字节截断,附总长提示。limit<=0 或未超时原样返。
// 注意:字节级截断可能切断多字节 UTF-8 字符尾部,显示上略难看但不 panic,可接受。
func truncateBody(s string, limit int) string {
	if limit <= 0 || len(s) <= limit {
		return s
	}
	return s[:limit] + fmt.Sprintf("...(总 %d 字节, 已截断)", len(s))
}
// 慢接口(> slowMs)标黄。用例数 ≤ 1 或全 skip 时不打印。
func (c *Collector) printDurationTop(slowMs int64) {
	const (
		yel = "\033[1;33m"
		nc  = "\033[0m"
	)
	type entry struct {
		suite    string
		name     string
		ms       int64
		status   string // pass/fail
	}
	dedup := map[string]*entry{} // key=suite/name
	for _, r := range c.Results {
		if r.Status == "skip" {
			continue
		}
		key := r.Suite + "/" + r.Name
		if e, ok := dedup[key]; ok {
			if r.DurationMs > e.ms {
				e.ms = r.DurationMs
			}
			continue
		}
		dedup[key] = &entry{suite: r.Suite, name: r.Name, ms: r.DurationMs, status: r.Status}
	}
	if len(dedup) <= 1 {
		return
	}
	list := make([]*entry, 0, len(dedup))
	for _, e := range dedup {
		list = append(list, e)
	}
	sort.Slice(list, func(i, j int) bool { return list[i].ms > list[j].ms })
	limit := min(len(list), 10)
	fmt.Println("\n接口耗时 Top:")
	for i := range limit {
		e := list[i]
		mark := ""
		if slowMs > 0 && e.ms > slowMs {
			mark = fmt.Sprintf(" %s⚠ SLOW(>%dms)%s", yel, slowMs, nc)
		}
		fmt.Printf("  %4dms  [%s/%s]%s\n", e.ms, e.suite, e.name, mark)
	}
}
