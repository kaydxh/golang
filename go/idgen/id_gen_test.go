package idgen

import (
	"regexp"
	"sync"
	"testing"
	"time"
)

// TestNewGenerator 测试生成器的创建
func TestNewGenerator(t *testing.T) {
	tests := []struct {
		name      string
		workerID  int64
		wantError bool
	}{
		{"valid worker ID 0", 0, false},
		{"valid worker ID 1", 1, false},
		{"valid worker ID 512", 512, false},
		{"valid worker ID 1023", 1023, false},
		{"invalid worker ID -1", -1, true},
		{"invalid worker ID 1024", 1024, true},
		{"invalid worker ID 2000", 2000, true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			gen, err := NewGenerator(tt.workerID)
			if tt.wantError {
				if err == nil {
					t.Errorf("NewGenerator(%d) expected error, got nil", tt.workerID)
				}
				if gen != nil {
					t.Errorf("NewGenerator(%d) expected nil generator on error", tt.workerID)
				}
			} else {
				if err != nil {
					t.Errorf("NewGenerator(%d) unexpected error: %v", tt.workerID, err)
				}
				if gen == nil {
					t.Errorf("NewGenerator(%d) returned nil generator", tt.workerID)
				}
				if gen.workerID != tt.workerID {
					t.Errorf("NewGenerator(%d) workerID = %d, want %d", tt.workerID, gen.workerID, tt.workerID)
				}
			}
		})
	}
}

// TestMustNewGenerator 测试 MustNewGenerator
func TestMustNewGenerator(t *testing.T) {
	// 正常情况
	t.Run("valid worker ID", func(t *testing.T) {
		defer func() {
			if r := recover(); r != nil {
				t.Errorf("MustNewGenerator(1) should not panic, got: %v", r)
			}
		}()
		gen := MustNewGenerator(1)
		if gen == nil {
			t.Error("MustNewGenerator(1) returned nil")
		}
	})

	// 应该 panic 的情况
	t.Run("invalid worker ID should panic", func(t *testing.T) {
		defer func() {
			if r := recover(); r == nil {
				t.Error("MustNewGenerator(1024) should panic but didn't")
			}
		}()
		MustNewGenerator(1024)
	})
}

// TestNextID 测试基本的ID生成
func TestNextID(t *testing.T) {
	gen := MustNewGenerator(1)

	// 生成多个ID，确保都能成功
	ids := make([]uint64, 100)
	for i := 0; i < 100; i++ {
		id, err := gen.NextID()
		if err != nil {
			t.Fatalf("NextID() error: %v", err)
		}
		if id == 0 {
			t.Error("NextID() returned 0")
		}
		ids[i] = id
	}

	// 检查ID是否递增
	for i := 1; i < len(ids); i++ {
		if ids[i] <= ids[i-1] {
			t.Errorf("IDs not increasing: ids[%d]=%d, ids[%d]=%d", i-1, ids[i-1], i, ids[i])
		}
	}
}

// TestMustNextID 测试 MustNextID
func TestMustNextID(t *testing.T) {
	gen := MustNewGenerator(1)

	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustNextID() should not panic, got: %v", r)
		}
	}()

	id := gen.MustNextID()
	if id == 0 {
		t.Error("MustNextID() returned 0")
	}
}

// TestIDUniqueness 测试ID的唯一性
func TestIDUniqueness(t *testing.T) {
	gen := MustNewGenerator(1)
	count := 10000
	ids := make(map[uint64]bool, count)

	for i := 0; i < count; i++ {
		id := gen.MustNextID()
		if ids[id] {
			t.Fatalf("Duplicate ID found: %d", id)
		}
		ids[id] = true
	}

	if len(ids) != count {
		t.Errorf("Expected %d unique IDs, got %d", count, len(ids))
	}
}

// TestWorkerIDInGeneratedID 测试生成的ID中是否包含正确的机器ID
func TestWorkerIDInGeneratedID(t *testing.T) {
	tests := []int64{0, 1, 100, 512, 1023}

	for _, workerID := range tests {
		t.Run(formatWorkerID(workerID), func(t *testing.T) {
			gen := MustNewGenerator(workerID)
			id := gen.MustNextID()

			extractedWorkerID := GetWorkerID(id)
			if extractedWorkerID != workerID {
				t.Errorf("WorkerID mismatch: generated with %d, extracted %d", workerID, extractedWorkerID)
			}
		})
	}
}

// TestParseID 测试ID解析功能
func TestParseID(t *testing.T) {
	workerID := int64(123)
	gen := MustNewGenerator(workerID)

	beforeTime := time.Now()
	id := gen.MustNextID()
	afterTime := time.Now()

	timestamp, parsedWorkerID, sequence := ParseID(id)

	// 检查时间戳是否在合理范围内
	if timestamp < beforeTime.UnixMilli() || timestamp > afterTime.UnixMilli() {
		t.Errorf("Timestamp out of range: %d not between %d and %d",
			timestamp, beforeTime.UnixMilli(), afterTime.UnixMilli())
	}

	// 检查机器ID
	if parsedWorkerID != workerID {
		t.Errorf("WorkerID mismatch: expected %d, got %d", workerID, parsedWorkerID)
	}

	// 检查序列号
	if sequence < 0 || sequence > maxSequence {
		t.Errorf("Sequence out of range: %d", sequence)
	}
}

// TestGetTimestamp 测试时间戳提取
func TestGetTimestamp(t *testing.T) {
	gen := MustNewGenerator(1)

	before := time.Now()
	time.Sleep(10 * time.Millisecond)
	id := gen.MustNextID()
	time.Sleep(10 * time.Millisecond)
	after := time.Now()

	timestamp := GetTimestamp(id)

	if timestamp.Before(before) || timestamp.After(after) {
		t.Errorf("Timestamp %v not between %v and %v", timestamp, before, after)
	}
}

// TestSequenceOverflow 测试同一毫秒内的序列号溢出
func TestSequenceOverflow(t *testing.T) {
	gen := MustNewGenerator(1)

	// 快速生成大量ID，可能触发序列号溢出
	count := 5000
	ids := make([]uint64, count)

	start := time.Now()
	for i := 0; i < count; i++ {
		ids[i] = gen.MustNextID()
	}
	elapsed := time.Since(start)

	t.Logf("Generated %d IDs in %v", count, elapsed)

	// 检查唯一性
	idMap := make(map[uint64]bool)
	for _, id := range ids {
		if idMap[id] {
			t.Fatalf("Duplicate ID found: %d", id)
		}
		idMap[id] = true
	}
}

// TestConcurrentGeneration 并发测试
func TestConcurrentGeneration(t *testing.T) {
	gen := MustNewGenerator(1)

	goroutines := 100
	idsPerGoroutine := 1000
	totalIDs := goroutines * idsPerGoroutine

	var wg sync.WaitGroup
	idChan := make(chan uint64, totalIDs)

	start := time.Now()

	// 启动多个goroutine并发生成ID
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id := gen.MustNextID()
				idChan <- id
			}
		}()
	}

	// 等待所有goroutine完成
	wg.Wait()
	close(idChan)

	elapsed := time.Since(start)
	t.Logf("Generated %d IDs concurrently in %v (%.0f IDs/sec)",
		totalIDs, elapsed, float64(totalIDs)/elapsed.Seconds())

	// 检查唯一性
	idMap := make(map[uint64]bool, totalIDs)
	duplicates := 0

	for id := range idChan {
		if idMap[id] {
			duplicates++
			t.Errorf("Duplicate ID found: %d", id)
		}
		idMap[id] = true
	}

	if duplicates > 0 {
		t.Errorf("Found %d duplicate IDs", duplicates)
	}

	if len(idMap) != totalIDs {
		t.Errorf("Expected %d unique IDs, got %d", totalIDs, len(idMap))
	}
}

// TestMultipleGenerators 测试多个生成器同时工作
func TestMultipleGenerators(t *testing.T) {
	workers := 10
	idsPerWorker := 1000
	totalIDs := workers * idsPerWorker

	var wg sync.WaitGroup
	idChan := make(chan uint64, totalIDs)

	// 创建多个生成器，每个使用不同的机器ID
	for workerID := 0; workerID < workers; workerID++ {
		wg.Add(1)
		go func(wID int64) {
			defer wg.Done()
			gen := MustNewGenerator(wID)
			for i := 0; i < idsPerWorker; i++ {
				id := gen.MustNextID()
				idChan <- id
			}
		}(int64(workerID))
	}

	wg.Wait()
	close(idChan)

	// 检查所有ID的唯一性
	idMap := make(map[uint64]bool, totalIDs)
	for id := range idChan {
		if idMap[id] {
			t.Errorf("Duplicate ID found: %d", id)
		}
		idMap[id] = true
	}

	if len(idMap) != totalIDs {
		t.Errorf("Expected %d unique IDs, got %d", totalIDs, len(idMap))
	}
}

// TestConcurrentSameWorker 测试同一个生成器的高并发场景
func TestConcurrentSameWorker(t *testing.T) {
	gen := MustNewGenerator(1)

	goroutines := 1000
	idsPerGoroutine := 100
	totalIDs := goroutines * idsPerGoroutine

	var wg sync.WaitGroup
	var mu sync.Mutex
	idMap := make(map[uint64]bool, totalIDs)
	duplicates := 0

	start := time.Now()

	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < idsPerGoroutine; j++ {
				id := gen.MustNextID()

				mu.Lock()
				if idMap[id] {
					duplicates++
				}
				idMap[id] = true
				mu.Unlock()
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	t.Logf("Generated %d IDs with %d goroutines in %v (%.0f IDs/sec)",
		totalIDs, goroutines, elapsed, float64(totalIDs)/elapsed.Seconds())

	if duplicates > 0 {
		t.Errorf("Found %d duplicate IDs in concurrent test", duplicates)
	}

	if len(idMap) != totalIDs {
		t.Errorf("Expected %d unique IDs, got %d", totalIDs, len(idMap))
	}
}

// BenchmarkNextID 基准测试：单线程生成性能
func BenchmarkNextID(b *testing.B) {
	gen := MustNewGenerator(1)
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		gen.MustNextID()
	}
}

// BenchmarkNextIDParallel 基准测试：并发生成性能
func BenchmarkNextIDParallel(b *testing.B) {
	gen := MustNewGenerator(1)
	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			gen.MustNextID()
		}
	})
}

// BenchmarkParseID 基准测试：ID解析性能
func BenchmarkParseID(b *testing.B) {
	gen := MustNewGenerator(1)
	id := gen.MustNextID()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		ParseID(id)
	}
}

// 辅助函数：格式化机器ID用于测试名称
func formatWorkerID(workerID int64) string {
	return string(rune('0' + workerID%10))
}

// uuidStdRe 校验标准 v4 UUID 字符串：8-4-4-4-12，且第三段以 4 开头、第四段以 8/9/a/b 开头。
var uuidStdRe = regexp.MustCompile(`^[0-9a-f]{8}-[0-9a-f]{4}-4[0-9a-f]{3}-[89ab][0-9a-f]{3}-[0-9a-f]{12}$`)

// uuidHexRe 校验去掉连字符的 32 字符 UUID。
var uuidHexRe = regexp.MustCompile(`^[0-9a-f]{32}$`)

// TestNewUUID 测试标准 UUID 字符串的格式与唯一性。
func TestNewUUID(t *testing.T) {
	const n = 1000
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		id, err := NewUUID()
		if err != nil {
			t.Fatalf("NewUUID() error: %v", err)
		}
		if !uuidStdRe.MatchString(id) {
			t.Fatalf("NewUUID() = %q, not a valid v4 UUID string", id)
		}
		if _, ok := seen[id]; ok {
			t.Fatalf("NewUUID() collision at i=%d: %s", i, id)
		}
		seen[id] = struct{}{}
	}
}

// TestMustNewUUID 在正常环境下不应 panic。
func TestMustNewUUID(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustNewUUID() unexpected panic: %v", r)
		}
	}()
	id := MustNewUUID()
	if !uuidStdRe.MatchString(id) {
		t.Errorf("MustNewUUID() = %q, not a valid v4 UUID string", id)
	}
}

// TestNewUUIDHex 测试紧凑 UUID 的格式与唯一性。
func TestNewUUIDHex(t *testing.T) {
	const n = 1000
	seen := make(map[string]struct{}, n)
	for i := 0; i < n; i++ {
		id, err := NewUUIDHex()
		if err != nil {
			t.Fatalf("NewUUIDHex() error: %v", err)
		}
		if !uuidHexRe.MatchString(id) {
			t.Fatalf("NewUUIDHex() = %q, not 32-char hex", id)
		}
		if _, ok := seen[id]; ok {
			t.Fatalf("NewUUIDHex() collision at i=%d: %s", i, id)
		}
		seen[id] = struct{}{}
	}
}

// TestMustNewUUIDHex 在正常环境下不应 panic。
func TestMustNewUUIDHex(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustNewUUIDHex() unexpected panic: %v", r)
		}
	}()
	id := MustNewUUIDHex()
	if !uuidHexRe.MatchString(id) {
		t.Errorf("MustNewUUIDHex() = %q, not 32-char hex", id)
	}
}

// TestNewUUIDBytes 测试 16 字节 UUID。
func TestNewUUIDBytes(t *testing.T) {
	a, err := NewUUIDBytes()
	if err != nil {
		t.Fatalf("NewUUIDBytes() error: %v", err)
	}
	b, err := NewUUIDBytes()
	if err != nil {
		t.Fatalf("NewUUIDBytes() error: %v", err)
	}
	if a == b {
		t.Fatalf("NewUUIDBytes() collision: %x", a)
	}
	// v4 版本号在第 7 字节高 4 位为 0100 (=4)
	if a[6]>>4 != 4 {
		t.Errorf("NewUUIDBytes() version nibble = %x, want 4", a[6]>>4)
	}
	// variant 在第 9 字节高 2 位为 10
	if a[8]>>6 != 0b10 {
		t.Errorf("NewUUIDBytes() variant bits = %b, want 10", a[8]>>6)
	}
}

// TestUUIDConcurrent 并发生成不应碰撞。
func TestUUIDConcurrent(t *testing.T) {
	const goroutines = 50
	const perG = 200
	total := goroutines * perG

	idCh := make(chan string, total)
	var wg sync.WaitGroup
	for i := 0; i < goroutines; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perG; j++ {
				idCh <- MustNewUUID()
			}
		}()
	}
	wg.Wait()
	close(idCh)

	seen := make(map[string]struct{}, total)
	for id := range idCh {
		if _, ok := seen[id]; ok {
			t.Fatalf("concurrent UUID collision: %s", id)
		}
		seen[id] = struct{}{}
	}
	if len(seen) != total {
		t.Errorf("expected %d unique UUIDs, got %d", total, len(seen))
	}

	_ = time.Now() // 仅占位，保持 time 包被使用（其它测试需要）
}

// BenchmarkNewUUID 基准测试：标准 UUID 字符串生成。
func BenchmarkNewUUID(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MustNewUUID()
	}
}

// BenchmarkNewUUIDHex 基准测试：紧凑 UUID 字符串生成。
func BenchmarkNewUUIDHex(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = MustNewUUIDHex()
	}
}
