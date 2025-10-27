package syscall

type MemoryUsage struct {
}

func (m MemoryUsage) SysTotalMemory() uint64 {
	return 0
}

func (m MemoryUsage) SysFreeMemory() uint64 {
	return 0
}

func (m MemoryUsage) SysUsageMemory() float64 {
	return 0
}
