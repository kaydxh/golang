package k8sdns

func WithNodeGroup(nodeGroup string) K8sDNSResolverOption {
	return K8sDNSResolverOptionFunc(func(r *K8sDNSResolver) {
		r.opts.nodeGroup = nodeGroup
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

func WithNamespace(namespace string) K8sDNSResolverOption {
	return K8sDNSResolverOptionFunc(func(r *K8sDNSResolver) {
		r.opts.namespace = namespace
	})
}

func WithUseInClusterConfig(useInClusterConfig bool) K8sDNSResolverOption {
	return K8sDNSResolverOptionFunc(func(r *K8sDNSResolver) {
		r.opts.useInClusterConfig = useInClusterConfig
	})
}

func WithUseInformer(useInformer bool) K8sDNSResolverOption {
	return K8sDNSResolverOptionFunc(func(r *K8sDNSResolver) {
		r.opts.useInformer = useInformer
	})
}
