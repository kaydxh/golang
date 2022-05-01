package pool

func WithBurst(burst int) PoolOptions {
	return PoolOptionsFunc(func(p *Pool) {
		p.opts.burst = burst
	})
}

func WithErrStop(errStop bool) PoolOptions {
	return PoolOptionsFunc(func(p *Pool) {
		p.opts.errStop = errStop
	})
}
