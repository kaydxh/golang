package idgen

import (
	"crypto/md5"
	"encoding/binary"
	"fmt"
	"net"
	"os"
	"strings"
)

// WorkerIDGenerator 用于根据机器信息自动生成 WorkerID
type WorkerIDGenerator struct{}

// NewWorkerIDGenerator 创建一个 WorkerID 生成器
func NewWorkerIDGenerator() *WorkerIDGenerator {
	return &WorkerIDGenerator{}
}

// GenerateWorkerID 根据机器信息自动生成 WorkerID (0-1023)
// 优先级: 环境变量 > IP地址 > MAC地址 > 主机名
func (w *WorkerIDGenerator) GenerateWorkerID() (int64, error) {
	// 方法1: 从环境变量读取
	if id, ok := w.fromEnv(); ok {
		return id, nil
	}

	// 方法2: 基于IP地址生成
	if id, ok := w.fromIP(); ok {
		return id, nil
	}

	// 方法3: 基于MAC地址生成
	if id, ok := w.fromMAC(); ok {
		return id, nil
	}

	// 方法4: 基于主机名生成
	if id, ok := w.fromHostname(); ok {
		return id, nil
	}

	return 0, fmt.Errorf("failed to generate worker ID from machine info")
}

// fromEnv 从环境变量读取 WorkerID
func (w *WorkerIDGenerator) fromEnv() (int64, bool) {
	envVars := []string{"WORKER_ID", "MACHINE_ID", "NODE_ID", "POD_ID"}

	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			var id int64
			_, err := fmt.Sscanf(value, "%d", &id)
			if err == nil && id >= 0 && id <= maxWorkerID {
				return id, true
			}
		}
	}

	return 0, false
}

// fromIP 基于本机IP地址生成 WorkerID
// 优先使用非回环的IPv4地址
func (w *WorkerIDGenerator) fromIP() (int64, bool) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return 0, false
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipv4 := ipnet.IP.To4(); ipv4 != nil {
				// 使用IP地址的后两个字节
				// 例如: 192.168.1.100 -> (1 << 8) + 100 = 356
				id := int64(ipv4[2])<<8 | int64(ipv4[3])
				return id & maxWorkerID, true // 确保在0-1023范围内
			}
		}
	}

	return 0, false
}

// fromMAC 基于MAC地址生成 WorkerID
func (w *WorkerIDGenerator) fromMAC() (int64, bool) {
	interfaces, err := net.Interfaces()
	if err != nil {
		return 0, false
	}

	for _, iface := range interfaces {
		// 跳过回环接口和没有MAC地址的接口
		if iface.Flags&net.FlagLoopback != 0 || len(iface.HardwareAddr) == 0 {
			continue
		}

		// 使用MAC地址的后两个字节
		mac := iface.HardwareAddr
		if len(mac) >= 6 {
			id := int64(mac[4])<<8 | int64(mac[5])
			return id & maxWorkerID, true
		}
	}

	return 0, false
}

// fromHostname 基于主机名生成 WorkerID
func (w *WorkerIDGenerator) fromHostname() (int64, bool) {
	hostname, err := os.Hostname()
	if err != nil || hostname == "" {
		return 0, false
	}

	// 使用MD5哈希主机名，取前8字节转换为int64
	hash := md5.Sum([]byte(hostname))
	id := int64(binary.BigEndian.Uint64(hash[:8]))
	return (id & maxWorkerID), true
}

// GetMachineInfo 获取机器信息用于调试
func (w *WorkerIDGenerator) GetMachineInfo() map[string]string {
	info := make(map[string]string)

	// 主机名
	if hostname, err := os.Hostname(); err == nil {
		info["hostname"] = hostname
	}

	// IP地址
	if ips := w.getIPAddresses(); len(ips) > 0 {
		info["ips"] = strings.Join(ips, ", ")
	}

	// MAC地址
	if macs := w.getMACAddresses(); len(macs) > 0 {
		info["macs"] = strings.Join(macs, ", ")
	}

	// 环境变量
	envVars := []string{"WORKER_ID", "MACHINE_ID", "NODE_ID", "POD_ID"}
	for _, envVar := range envVars {
		if value := os.Getenv(envVar); value != "" {
			info[envVar] = value
		}
	}

	return info
}

// getIPAddresses 获取所有非回环的IPv4地址
func (w *WorkerIDGenerator) getIPAddresses() []string {
	var ips []string
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ips
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipv4 := ipnet.IP.To4(); ipv4 != nil {
				ips = append(ips, ipv4.String())
			}
		}
	}

	return ips
}

// getMACAddresses 获取所有网络接口的MAC地址
func (w *WorkerIDGenerator) getMACAddresses() []string {
	var macs []string
	interfaces, err := net.Interfaces()
	if err != nil {
		return macs
	}

	for _, iface := range interfaces {
		if iface.Flags&net.FlagLoopback == 0 && len(iface.HardwareAddr) > 0 {
			macs = append(macs, iface.HardwareAddr.String())
		}
	}

	return macs
}

// GenerateFromIP 根据指定IP地址生成 WorkerID
func GenerateFromIP(ip string) (int64, error) {
	parsedIP := net.ParseIP(ip)
	if parsedIP == nil {
		return 0, fmt.Errorf("invalid IP address: %s", ip)
	}

	ipv4 := parsedIP.To4()
	if ipv4 == nil {
		return 0, fmt.Errorf("only IPv4 addresses are supported")
	}

	// 使用IP的后两个字节
	id := int64(ipv4[2])<<8 | int64(ipv4[3])
	return id & maxWorkerID, nil
}

// GenerateFromMAC 根据指定MAC地址生成 WorkerID
func GenerateFromMAC(mac string) (int64, error) {
	hwAddr, err := net.ParseMAC(mac)
	if err != nil {
		return 0, fmt.Errorf("invalid MAC address: %s", mac)
	}

	if len(hwAddr) < 6 {
		return 0, fmt.Errorf("MAC address too short")
	}

	// 使用MAC地址的后两个字节
	id := int64(hwAddr[4])<<8 | int64(hwAddr[5])
	return id & maxWorkerID, nil
}

// GenerateFromString 根据任意字符串生成 WorkerID
func GenerateFromString(s string) int64 {
	hash := md5.Sum([]byte(s))
	id := int64(binary.BigEndian.Uint64(hash[:8]))
	return id & maxWorkerID
}

// AutoGenerator 自动生成并创建 Generator
// 这是一个便捷函数，结合了 WorkerID 生成和 Generator 创建
func AutoGenerator() (*Generator, error) {
	wg := NewWorkerIDGenerator()
	workerID, err := wg.GenerateWorkerID()
	if err != nil {
		return nil, fmt.Errorf("failed to generate worker ID: %w", err)
	}

	gen, err := NewGenerator(workerID)
	if err != nil {
		return nil, fmt.Errorf("failed to create generator with worker ID %d: %w", workerID, err)
	}

	return gen, nil
}

// MustAutoGenerator 自动生成并创建 Generator，失败则 panic
func MustAutoGenerator() *Generator {
	gen, err := AutoGenerator()
	if err != nil {
		panic(fmt.Sprintf("failed to auto-generate: %v", err))
	}
	return gen
}
