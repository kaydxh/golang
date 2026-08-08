package apireport

import (
	"bytes"
	"io"
	"os"
	"strings"
	"testing"
)

// captureStdout 捕获 PrintConsole 输出便于断言。
func captureStdout(f func()) string {
	r, w, _ := os.Pipe()
	orig := os.Stdout
	os.Stdout = w
	f()
	w.Close()
	os.Stdout = orig
	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

// topSection 抽取 "接口耗时 Top:" 之后的段落用于断言,避免 PASS/SKIP 行干扰。
func topSection(out string) string {
	idx := strings.Index(out, "接口耗时 Top:")
	if idx < 0 {
		return ""
	}
	return out[idx:]
}

func TestPrintConsole_DurationTop_Sorted(t *testing.T) {
	c := &Collector{}
	c.Add(Result{Suite: "sdk-device", Name: "Fast", Status: "pass", DurationMs: 50})
	c.Add(Result{Suite: "sdk-device", Name: "Slow", Status: "pass", DurationMs: 500})
	c.Add(Result{Suite: "sdk-device", Name: "Mid", Status: "pass", DurationMs: 150})
	out := captureStdout(func() { c.PrintConsole(false, 0) })

	top := topSection(out)
	if top == "" {
		t.Fatalf("missing top section\n%s", out)
	}
	slow := strings.Index(top, "[sdk-device/Slow]")
	mid := strings.Index(top, "[sdk-device/Mid]")
	fast := strings.Index(top, "[sdk-device/Fast]")
	if slow < 0 || mid < 0 || fast < 0 {
		t.Fatalf("missing entries in top\n%s", top)
	}
	if !(slow < mid && mid < fast) {
		t.Errorf("ordering wrong: slow=%d mid=%d fast=%d\n%s", slow, mid, fast, top)
	}
}

func TestPrintConsole_DurationTop_SkipsSkipStatus(t *testing.T) {
	c := &Collector{}
	c.Add(Result{Suite: "s", Name: "A", Status: "pass", DurationMs: 100})
	c.Add(Result{Suite: "s", Name: "B", Status: "skip"})
	out := captureStdout(func() { c.PrintConsole(false, 0) })

	// 只有 1 个非 skip → 不打印 Top
	if strings.Contains(out, "接口耗时 Top:") {
		t.Errorf("should skip Top for single case\n%s", out)
	}
	// 加一个 pass 让 Top 出现
	c.Add(Result{Suite: "s", Name: "C", Status: "pass", DurationMs: 200})
	out = captureStdout(func() { c.PrintConsole(false, 0) })
	top := topSection(out)
	if top == "" {
		t.Fatalf("should show Top for 2+ cases\n%s", out)
	}
	if strings.Contains(top, "[s/B]") {
		t.Errorf("skip status should be excluded from Top section\n%s", top)
	}
}

func TestPrintConsole_DurationTop_DedupSameCase(t *testing.T) {
	c := &Collector{}
	c.Add(Result{Suite: "s", Name: "A", Status: "pass", DurationMs: 100})
	c.Add(Result{Suite: "s", Name: "A", Status: "pass", DurationMs: 100})
	c.Add(Result{Suite: "s", Name: "B", Status: "pass", DurationMs: 200})
	out := captureStdout(func() { c.PrintConsole(false, 0) })

	top := topSection(out)
	if strings.Count(top, "[s/A]") != 1 {
		t.Errorf("[s/A] should appear once in Top section, got %d\n%s", strings.Count(top, "[s/A]"), top)
	}
}

func TestPrintConsole_SlowMark(t *testing.T) {
	c := &Collector{}
	c.Add(Result{Suite: "s", Name: "Slow", Status: "pass", DurationMs: 2000})
	c.Add(Result{Suite: "s", Name: "Fast", Status: "pass", DurationMs: 50})
	out := captureStdout(func() { c.PrintConsole(false, 1000) })

	if !strings.Contains(out, "SLOW(>1000ms)") {
		t.Errorf("missing SLOW mark on Slow entry\n%s", out)
	}
	// Fast(50ms) 不该标 SLOW
	// SLOW 出现次数:内联 PASS 段 1 次 + Top 段 1 次 = 2 次(仅 Slow)
	if strings.Count(out, "SLOW") != 2 {
		t.Errorf("SLOW should appear exactly twice (inline + Top), got %d\n%s", strings.Count(out, "SLOW"), out)
	}
}

func TestPrintConsole_FailShowsDuration(t *testing.T) {
	c := &Collector{}
	c.Add(Result{Suite: "s", Name: "X", Status: "fail", DurationMs: 300, Reason: "boom", CurlRepro: "curl ..."})
	c.Add(Result{Suite: "s", Name: "Y", Status: "pass", DurationMs: 100})
	out := captureStdout(func() { c.PrintConsole(false, 0) })

	// FAIL 行含 "300ms"
	if !strings.Contains(out, "✗ FAIL") || !strings.Contains(out, "300ms") {
		t.Errorf("FAIL line missing duration\n%s", out)
	}
}

func TestPrintConsole_BodyLimit_Truncates(t *testing.T) {
	longBody := strings.Repeat("A", 2000)
	c := &Collector{BodyLimit: 100}
	c.Add(Result{Suite: "s", Name: "X", Status: "fail", DurationMs: 10, Reason: "r", CurlRepro: "cc", Body: longBody})
	out := captureStdout(func() { c.PrintConsole(false, 0) })

	if !strings.Contains(out, "已截断") {
		t.Errorf("expect truncation marker\n%s", out)
	}
	if !strings.Contains(out, "总 2000 字节") {
		t.Errorf("expect total length marker\n%s", out)
	}
	// 不应包含完整 2000 A(会截到 100)
	if strings.Contains(out, strings.Repeat("A", 200)) {
		t.Errorf("body not truncated\n%s", out)
	}
}

func TestPrintConsole_BodyLimit_NoLimitPassThru(t *testing.T) {
	longBody := strings.Repeat("A", 2000)
	c := &Collector{BodyLimit: -1} // 不限
	c.Add(Result{Suite: "s", Name: "X", Status: "fail", DurationMs: 10, Reason: "r", CurlRepro: "cc", Body: longBody})
	out := captureStdout(func() { c.PrintConsole(false, 0) })

	if strings.Contains(out, "已截断") {
		t.Errorf("BodyLimit=-1 should not truncate\n%s", out)
	}
	if !strings.Contains(out, strings.Repeat("A", 2000)) {
		t.Errorf("full body missing\n%s", out)
	}
}

func TestPrintConsole_BodyLimit_ShortBodyUnaffected(t *testing.T) {
	c := &Collector{BodyLimit: 500}
	c.Add(Result{Suite: "s", Name: "X", Status: "fail", DurationMs: 10, Reason: "r", CurlRepro: "cc", Body: "short"})
	out := captureStdout(func() { c.PrintConsole(false, 0) })

	if strings.Contains(out, "已截断") {
		t.Errorf("short body should not be truncated\n%s", out)
	}
	if !strings.Contains(out, "响应: short") {
		t.Errorf("body missing\n%s", out)
	}
}
