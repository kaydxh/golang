package idgen

import (
	"os"
	"testing"
)

// TestGenerateWorkerID 测试自动生成 WorkerID
func TestGenerateWorkerID(t *testing.T) {
	wg := NewWorkerIDGenerator()
	workerID, err := wg.GenerateWorkerID()

	if err != nil {
		t.Fatalf("GenerateWorkerID() error: %v", err)
	}

	if workerID < 0 || workerID > maxWorkerID {
		t.Errorf("WorkerID out of range: %d, expected 0-%d", workerID, maxWorkerID)
	}

	t.Logf("Generated WorkerID: %d", workerID)
}

// TestFromEnv 测试从环境变量生成
func TestFromEnv(t *testing.T) {
	tests := []struct {
		name     string
		envVar   string
		envValue string
		want     int64
		wantOK   bool
	}{
		{"valid WORKER_ID", "WORKER_ID", "42", 42, true},
		{"valid MACHINE_ID", "MACHINE_ID", "100", 100, true},
		{"valid NODE_ID", "NODE_ID", "1023", 1023, true},
		{"invalid value", "WORKER_ID", "abc", 0, false},
		{"out of range", "WORKER_ID", "2000", 0, false},
		{"negative value", "WORKER_ID", "-1", 0, false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			os.Setenv(tt.envVar, tt.envValue)
			defer os.Unsetenv(tt.envVar)

			wg := NewWorkerIDGenerator()
			got, ok := wg.fromEnv()

			if ok != tt.wantOK {
				t.Errorf("fromEnv() ok = %v, want %v", ok, tt.wantOK)
			}

			if ok && got != tt.want {
				t.Errorf("fromEnv() = %d, want %d", got, tt.want)
			}
		})
	}
}

// TestGenerateFromIP 测试从IP生成
func TestGenerateFromIP(t *testing.T) {
	tests := []struct {
		name    string
		ip      string
		wantErr bool
	}{
		{"valid IPv4", "192.168.1.100", false},
		{"valid IPv4 2", "10.0.0.50", false},
		{"valid IPv4 3", "172.16.5.200", false},
		{"invalid IP", "256.1.1.1", true},
		{"invalid format", "not-an-ip", true},
		{"IPv6 not supported", "2001:0db8::1", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workerID, err := GenerateFromIP(tt.ip)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFromIP() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if workerID < 0 || workerID > maxWorkerID {
					t.Errorf("WorkerID out of range: %d", workerID)
				}
				t.Logf("IP %s -> WorkerID: %d", tt.ip, workerID)
			}
		})
	}
}

// TestGenerateFromMAC 测试从MAC地址生成
func TestGenerateFromMAC(t *testing.T) {
	tests := []struct {
		name    string
		mac     string
		wantErr bool
	}{
		{"valid MAC", "00:1A:2B:3C:4D:5E", false},
		{"valid MAC 2", "AA:BB:CC:DD:EE:FF", false},
		{"valid MAC hyphen", "00-1A-2B-3C-4D-5E", false},
		{"invalid MAC", "not-a-mac", true},
		{"invalid format", "00:1A:2B:3C", true},
		{"empty string", "", true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			workerID, err := GenerateFromMAC(tt.mac)

			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateFromMAC() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if err == nil {
				if workerID < 0 || workerID > maxWorkerID {
					t.Errorf("WorkerID out of range: %d", workerID)
				}
				t.Logf("MAC %s -> WorkerID: %d", tt.mac, workerID)
			}
		})
	}
}

// TestGenerateFromString 测试从字符串生成
func TestGenerateFromString(t *testing.T) {
	tests := []string{
		"web-server-1",
		"api-gateway-node-3",
		"database-replica-5",
		"my-app-pod-abc123",
		"",
		"very-long-string-with-many-characters-to-test-hashing",
	}

	for _, s := range tests {
		t.Run(s, func(t *testing.T) {
			workerID := GenerateFromString(s)

			if workerID < 0 || workerID > maxWorkerID {
				t.Errorf("WorkerID out of range: %d", workerID)
			}

			t.Logf("String '%s' -> WorkerID: %d", s, workerID)
		})
	}
}

// TestGenerateFromStringConsistency 测试字符串生成的一致性
func TestGenerateFromStringConsistency(t *testing.T) {
	testString := "test-pod-123"

	// 多次生成应该返回相同的 WorkerID
	first := GenerateFromString(testString)
	for i := 0; i < 100; i++ {
		id := GenerateFromString(testString)
		if id != first {
			t.Errorf("Inconsistent WorkerID: first=%d, iteration %d=%d", first, i, id)
		}
	}

	t.Logf("Consistent WorkerID for '%s': %d", testString, first)
}

// TestGenerateFromStringUniqueness 测试不同字符串生成不同的 WorkerID
func TestGenerateFromStringUniqueness(t *testing.T) {
	strings := []string{
		"pod-1", "pod-2", "pod-3", "pod-4", "pod-5",
		"node-a", "node-b", "node-c", "node-d", "node-e",
	}

	seen := make(map[int64]string)
	collisions := 0

	for _, s := range strings {
		id := GenerateFromString(s)
		if existing, exists := seen[id]; exists {
			collisions++
			t.Logf("Collision: '%s' and '%s' both generate WorkerID %d", s, existing, id)
		} else {
			seen[id] = s
		}
	}

	t.Logf("Generated %d WorkerIDs from %d strings, %d collisions",
		len(seen), len(strings), collisions)

	// 碰撞率应该很低（对于10个字符串，期望0-1个碰撞）
	if collisions > 2 {
		t.Errorf("Too many collisions: %d", collisions)
	}
}

// TestAutoGenerator 测试自动生成器
func TestAutoGenerator(t *testing.T) {
	gen, err := AutoGenerator()
	if err != nil {
		t.Fatalf("AutoGenerator() error: %v", err)
	}

	if gen == nil {
		t.Fatal("AutoGenerator() returned nil")
	}

	// 测试生成ID
	id, err := gen.NextID()
	if err != nil {
		t.Errorf("NextID() error: %v", err)
	}

	if id == 0 {
		t.Error("NextID() returned 0")
	}

	t.Logf("AutoGenerator created successfully, generated ID: %d", id)
}

// TestMustAutoGenerator 测试 MustAutoGenerator
func TestMustAutoGenerator(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			t.Errorf("MustAutoGenerator() should not panic, got: %v", r)
		}
	}()

	gen := MustAutoGenerator()
	if gen == nil {
		t.Fatal("MustAutoGenerator() returned nil")
	}

	id := gen.MustNextID()
	if id == 0 {
		t.Error("MustNextID() returned 0")
	}

	t.Logf("MustAutoGenerator created successfully, generated ID: %d", id)
}

// TestGetMachineInfo 测试获取机器信息
func TestGetMachineInfo(t *testing.T) {
	wg := NewWorkerIDGenerator()
	info := wg.GetMachineInfo()

	if len(info) == 0 {
		t.Error("GetMachineInfo() returned empty map")
	}

	t.Log("Machine Info:")
	for key, value := range info {
		t.Logf("  %s: %s", key, value)
	}
}

// TestWorkerIDRange 测试所有生成方法都返回有效范围的 WorkerID
func TestWorkerIDRange(t *testing.T) {
	wg := NewWorkerIDGenerator()

	// 测试自动生成
	if id, err := wg.GenerateWorkerID(); err == nil {
		if id < 0 || id > maxWorkerID {
			t.Errorf("GenerateWorkerID() = %d, out of range [0, %d]", id, maxWorkerID)
		}
	}

	// 测试从字符串生成
	for i := 0; i < 100; i++ {
		id := GenerateFromString(string(rune(i)))
		if id < 0 || id > maxWorkerID {
			t.Errorf("GenerateFromString() = %d, out of range [0, %d]", id, maxWorkerID)
		}
	}
}

// TestEnvironmentPriority 测试环境变量优先级
func TestEnvironmentPriority(t *testing.T) {
	// 设置环境变量
	os.Setenv("WORKER_ID", "999")
	defer os.Unsetenv("WORKER_ID")

	wg := NewWorkerIDGenerator()
	workerID, err := wg.GenerateWorkerID()

	if err != nil {
		t.Fatalf("GenerateWorkerID() error: %v", err)
	}

	// 应该使用环境变量的值
	if workerID != 999 {
		t.Errorf("Expected WorkerID from env (999), got %d", workerID)
	}

	t.Logf("Environment variable correctly prioritized: %d", workerID)
}

// BenchmarkGenerateWorkerID 基准测试：自动生成
func BenchmarkGenerateWorkerID(b *testing.B) {
	wg := NewWorkerIDGenerator()
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		wg.GenerateWorkerID()
	}
}

// BenchmarkGenerateFromString 基准测试：从字符串生成
func BenchmarkGenerateFromString(b *testing.B) {
	testString := "my-pod-name-123"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GenerateFromString(testString)
	}
}

// BenchmarkGenerateFromIP 基准测试：从IP生成
func BenchmarkGenerateFromIP(b *testing.B) {
	testIP := "192.168.1.100"
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		GenerateFromIP(testIP)
	}
}

// BenchmarkAutoGenerator 基准测试：自动生成器
func BenchmarkAutoGenerator(b *testing.B) {
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		AutoGenerator()
	}
}
