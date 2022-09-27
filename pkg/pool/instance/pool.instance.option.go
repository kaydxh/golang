/*
 *Copyright (c) 2022, kaydxh
 *
 *Permission is hereby granted, free of charge, to any person obtaining a copy
 *of this software and associated documentation files (the "Software"), to deal
 *in the Software without restriction, including without limitation the rights
 *to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 *copies of the Software, and to permit persons to whom the Software is
 *furnished to do so, subject to the following conditions:
 *
 *The above copyright notice and this permission notice shall be included in all
 *copies or substantial portions of the Software.
 *
 *THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 *IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 *FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 *AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 *LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 *OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 *SOFTWARE.
 */
package instance

import "time"

func WithName(name string) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.name = name
	})
}

func WithResevePoolSizePerCore(resevePoolSizePerCore int64) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.resevePoolSizePerCore = resevePoolSizePerCore
	})
}

func WithCapacityPoolSizePerCore(capacityPoolSizePerCore int64) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.capacityPoolSizePerCore = capacityPoolSizePerCore
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

func WithCoreIDs(coreIDs []int64) PoolOption {
	return PoolOptionFunc(func(p *Pool) {
		p.opts.coreIDs = coreIDs
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
