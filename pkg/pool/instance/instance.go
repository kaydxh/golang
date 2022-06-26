package instance

// CoreInstanceHolder is a instance which is used by core (gpu or cpu or npu).
type CoreInstanceHolder struct {
	Instance   interface{}
	Name       string
	CoreID     int64
	ModelPaths []string
	BatchSize  int64
	Thread
}
