// Package idgen 提供基于 Snowflake 算法的分布式唯一 ID 生成器
//
// ID 结构 (64位):
//   - 1位: 保留位(始终为0)
//   - 41位: 时间戳(毫秒级)
//   - 10位: 机器ID(支持0-1023台机器)
//   - 12位: 序列号(每毫秒最多4096个ID)
//
// 使用示例:
//
//	gen := idgen.NewGenerator(1) // 机器ID为1
//	id := gen.NextID()
//	fmt.Printf("生成的ID: %d\n", id)
package idgen

import (
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

const (
	// Epoch 是自定义纪元时间戳(毫秒)，默认为 2024-01-01 00:00:00 UTC
	// 可以根据项目实际启动时间调整，以获得更长的使用寿命
	Epoch int64 = 1704067200000 // 2024-01-01 00:00:00 UTC

	// 各部分的位数
	workerIDBits = 10 // 机器ID位数
	sequenceBits = 12 // 序列号位数

	// 最大值
	maxWorkerID = -1 ^ (-1 << workerIDBits) // 1023
	maxSequence = -1 ^ (-1 << sequenceBits) // 4095

	// 位移量
	workerIDShift  = sequenceBits                // 12
	timestampShift = sequenceBits + workerIDBits // 22
)

var (
	// ErrInvalidWorkerID 当机器ID超出范围时返回
	ErrInvalidWorkerID = errors.New("worker ID must be between 0 and 1023")

	// ErrClockMovedBackwards 当系统时钟回拨时返回
	ErrClockMovedBackwards = errors.New("clock moved backwards")
)

// Generator 是一个线程安全的 ID 生成器
type Generator struct {
	mu            sync.Mutex
	workerID      int64
	sequence      int64
	lastTimestamp int64
}

// NewGenerator 创建一个新的 ID 生成器
// workerID 必须在 0-1023 范围内，用于区分不同的机器或实例
func NewGenerator(workerID int64) (*Generator, error) {
	if workerID < 0 || workerID > maxWorkerID {
		return nil, ErrInvalidWorkerID
	}

	return &Generator{
		workerID:      workerID,
		sequence:      0,
		lastTimestamp: 0,
	}, nil
}

// MustNewGenerator 创建一个新的 ID 生成器，如果出错则 panic
// 适用于在初始化阶段使用
func MustNewGenerator(workerID int64) *Generator {
	gen, err := NewGenerator(workerID)
	if err != nil {
		panic(fmt.Sprintf("failed to create generator: %v", err))
	}
	return gen
}

// NextID 生成下一个唯一 ID
// 返回一个 uint64 类型的唯一标识符
func (g *Generator) NextID() (uint64, error) {
	g.mu.Lock()
	defer g.mu.Unlock()

	timestamp := g.currentMillis()

	// 检测时钟回拨
	if timestamp < g.lastTimestamp {
		return 0, fmt.Errorf("%w: refused to generate id for %d milliseconds",
			ErrClockMovedBackwards, g.lastTimestamp-timestamp)
	}

	// 同一毫秒内
	if timestamp == g.lastTimestamp {
		g.sequence = (g.sequence + 1) & maxSequence
		if g.sequence == 0 {
			// 序列号溢出，等待下一毫秒
			timestamp = g.waitNextMillis(g.lastTimestamp)
		}
	} else {
		// 新的毫秒，重置序列号
		g.sequence = 0
	}

	g.lastTimestamp = timestamp

	// 组合各部分生成最终 ID
	id := uint64((timestamp-Epoch)<<timestampShift |
		(g.workerID << workerIDShift) |
		g.sequence)

	return id, nil
}

// MustNextID 生成下一个唯一 ID，如果出错则 panic
func (g *Generator) MustNextID() uint64 {
	id, err := g.NextID()
	if err != nil {
		panic(fmt.Sprintf("failed to generate id: %v", err))
	}
	return id
}

// currentMillis 返回当前时间戳(毫秒)
func (g *Generator) currentMillis() int64 {
	return time.Now().UnixMilli()
}

// waitNextMillis 等待直到下一毫秒
func (g *Generator) waitNextMillis(lastTimestamp int64) int64 {
	timestamp := g.currentMillis()
	for timestamp <= lastTimestamp {
		timestamp = g.currentMillis()
	}
	return timestamp
}

// ParseID 解析 ID，返回时间戳、机器ID和序列号
func ParseID(id uint64) (timestamp int64, workerID int64, sequence int64) {
	sequence = int64(id) & maxSequence
	workerID = (int64(id) >> workerIDShift) & maxWorkerID
	timestamp = (int64(id) >> timestampShift) + Epoch
	return
}

// GetTimestamp 从 ID 中提取时间戳
func GetTimestamp(id uint64) time.Time {
	timestamp, _, _ := ParseID(id)
	return time.UnixMilli(timestamp)
}

// GetWorkerID 从 ID 中提取机器ID
func GetWorkerID(id uint64) int64 {
	_, workerID, _ := ParseID(id)
	return workerID
}

// GetSequence 从 ID 中提取序列号
func GetSequence(id uint64) int64 {
	_, _, sequence := ParseID(id)
	return sequence
}

// GenerateUint64FromUUID 生成uint64
func GenerateUint64FromUUID() uint64 {
	id := uuid.New()
	return UUIDToUint64XOR(id)
}

// UUIDToUint64XOR 将UUID转换为uint64
func UUIDToUint64XOR(id uuid.UUID) uint64 {
	high := binary.BigEndian.Uint64(id[:8])
	low := binary.BigEndian.Uint64(id[8:])
	return high ^ low
}

// NewUUID 生成一个 v4 UUID 字符串，形如 "xxxxxxxx-xxxx-4xxx-yxxx-xxxxxxxxxxxx"（36 字符）。
//
// 适用于需要全局唯一、不可猜的标识符（如 trace ID、session ID、请求 ID 等）。
// 内部使用 crypto/rand，熵源读取失败时返回 error；调用方若希望失败即 panic，
// 应改用 MustNewUUID。
func NewUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("idgen: new uuid: %w", err)
	}
	return id.String(), nil
}

// MustNewUUID 生成一个 v4 UUID 字符串；熵源失败时 panic。
//
// 适用于初始化阶段或对失败概率不敏感的场景（如本地短期 ID）。
// 在长期运行的请求路径上，倾向使用 NewUUID 显式处理错误。
func MustNewUUID() string {
	s, err := NewUUID()
	if err != nil {
		panic(fmt.Sprintf("failed to generate uuid: %v", err))
	}
	return s
}

// NewUUIDHex 生成一个去掉连字符的 v4 UUID 字符串（32 字符，等价 16 字节随机的 hex 表示）。
//
// 相比 NewUUID 更紧凑，常用作 URL 路径段、Redis key 后缀、不希望出现 '-' 的场景。
// 安全强度与 NewUUID 一致（同样 122 bit 熵）。
func NewUUIDHex() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("idgen: new uuid hex: %w", err)
	}
	return hex.EncodeToString(id[:]), nil
}

// MustNewUUIDHex 生成 32 字符的紧凑 UUID 字符串；熵源失败时 panic。
func MustNewUUIDHex() string {
	s, err := NewUUIDHex()
	if err != nil {
		panic(fmt.Sprintf("failed to generate uuid hex: %v", err))
	}
	return s
}

// NewUUIDBytes 生成一个 16 字节的 v4 UUID。
//
// 适用于需要把 UUID 作为二进制写入存储（如数据库 BINARY(16) 列、二进制协议）的场景，
// 避免字符串形态的 36 字节开销。
func NewUUIDBytes() ([16]byte, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return [16]byte{}, fmt.Errorf("idgen: new uuid bytes: %w", err)
	}
	return id, nil
}
