package k8sdns

func WithGroupNode(groupNode string) K8sDNSResolverOption {
	return K8sDNSResolverOptionFunc(func(r *K8sDNSResolver) {
		r.opts.groupNode = groupNode
	})
}

func WithNodeUnit(nodeUnit string) K8sDNSResolverOption {
	return K8sDNSResolverOptionFunc(func(r *K8sDNSResolver) {
		r.opts.nodeUnit = nodeUnit
	})
}

func WithkubeConfig(kubeConfig string) K8sDNSResolverOption {
	return K8sDNSResolverOptionFunc(func(r *K8sDNSResolver) {
		r.opts.kubeConfig = kubeConfig
	})
}
