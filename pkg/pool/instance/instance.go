package instance

// CoreInstanceHolder is a instance which is used by gpu or cpu.
type CoreInstanceHolder struct {
	Instance   interface{}
	Name       string
	GpuID      int64
	ModelPaths []string
	BatchSize  int64
	Thread
}
