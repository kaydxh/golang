// Package apireport 提供接口测试结果的收集、控制台输出与退出码，通用无业务。
package apireport

import "fmt"

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
}

// Collector 汇总结果。
type Collector struct {
	Results []Result
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
		case "skip":
			fmt.Printf("  %s- SKIP%s [%s/%s] %s\n", yel, nc, r.Suite, r.Name, r.Reason)
		default:
			fmt.Printf("  %s✗ FAIL%s [%s/%s] %s\n", red, nc, r.Suite, r.Name, r.Reason)
			fmt.Printf("    复现: %s\n", r.CurlRepro)
		}
	}
	pass, fail, skip := c.counts()
	fmt.Printf("\n总计 %d | 通过 %d | 失败 %d | 跳过 %d\n", len(c.Results), pass, fail, skip)
}
