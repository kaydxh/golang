package instance

import "time"

func WithName(name string) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.name = name
	})
}

func WithResevePoolSizePerGpu(resevePoolSizePerGpu int64) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.resevePoolSizePerGpu = resevePoolSizePerGpu
	})
}

func WithCapacityPoolSizePerGpu(capacityPoolSizePerGpu int64) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.capacityPoolSizePerGpu = capacityPoolSizePerGpu
	})
}

func WithWaitTimeoutOnce(waitTimeout time.Duration) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.waitTimeoutOnce = waitTimeout
	})
}

func WithWaitTimeoutTotal(waitTimeout time.Duration) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.waitTimeoutTotal = waitTimeout
	})
}

func WithLoadBalanceMode(loadBalanceMode LoadBalanceMode) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.loadBalanceMode = loadBalanceMode
	})
}

func WithGpuIDs(gpuIDs []int64) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.gpuIDs = gpuIDs
	})
}

func WithModelPaths(modelPaths []string) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.modelPaths = modelPaths
	})
}

func WithBatchSize(batchSize int64) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.batchSize = batchSize
	})
}

func WithGlobalInitFunc(f GlobalInitFunc) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.globalInitFunc = f
	})
}

func WithGlobalReleaseFunc(f GlobalReleaseFunc) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.globalReleaseFunc = f
	})
}

func WithLocalInitFunc(f LocalInitFunc) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.localInitFunc = f
	})
}

func WithLocalReleaseFunc(f LocalReleaseFunc) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.localReleaseFunc = f
	})
}

func WithDeleteFunc(f DeleteFunc) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.deleteFunc = f
	})
}

func WithEnabledPrintCostTime(enabled bool) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.enabledPrintCostTime = enabled
	})
}