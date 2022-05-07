package k8sdns

import (
	"context"
	"errors"
	"fmt"
	"time"

	context_ "github.com/kaydxh/golang/go/context"
	net_ "github.com/kaydxh/golang/go/net"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/apimachinery/pkg/selection"
	"k8s.io/client-go/informers"
	informerv1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	defaultInClusterConfig = true
	defaultUseInformer     = true
	defaultNamespace       = corev1.NamespaceAll
	defaultTimeout         = 10 * time.Second
)

type ResolverConfig struct {
	nodeGroup          string
	nodeUnit           string
	namespace          string
	useInClusterConfig bool
	useInformer        bool
	kubeConfig         string
	timeout            time.Duration
}

type K8sDNSResolver struct {
	opts        ResolverConfig
	stopCh      chan struct{}
	kubeClient  kubernetes.Interface
	factory     informers.SharedInformerFactory
	podInformer informerv1.PodInformer
}

func NewK8sDNSResolver(opts ...K8sDNSResolverOption) (*K8sDNSResolver, error) {

	r := &K8sDNSResolver{
		stopCh: make(chan struct{}),
	}
	r.opts.useInClusterConfig = defaultInClusterConfig
	r.opts.useInformer = defaultUseInformer
	r.opts.namespace = defaultNamespace
	r.opts.timeout = defaultTimeout
	r.ApplyOptions(opts...)

	var (
		restConfig *rest.Config
		err        error
	)

	if r.opts.useInClusterConfig {
		restConfig, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}

	} else if r.opts.kubeConfig != "" {
		restConfig, err = clientcmd.BuildConfigFromFlags("", r.opts.kubeConfig)
		if err != nil {
			return nil, err
		}

	} else {
		return nil, fmt.Errorf("useInClusterConfig false and kubeConfig is empty")
	}

	r.kubeClient, err = kubernetes.NewForConfig(restConfig)
	if err != nil {
		return nil, err
	}

	//r.factory = informers.NewSharedInformerFactory(r.kubeClient, 0, informers.WithNamespace(r.opts.namespace))
	r.factory = informers.NewSharedInformerFactory(r.kubeClient, 0)
	r.podInformer = r.factory.Core().V1().Pods()

	//only run once
	go r.podInformer.Informer().Run(r.stopCh)

	return r, nil
}

func (r *K8sDNSResolver) Pods(ctx context.Context, svcs ...string) ([]corev1.Pod, error) {

	var pods []corev1.Pod
	serviceLabel, _ := labels.NewRequirement("k8s-app", selection.Equals, svcs)
	selector := labels.NewSelector()
	selector = selector.Add(*serviceLabel)

	ctx, cancel := context_.WithTimeout(ctx, r.opts.timeout)
	defer cancel()

	waitCh := make(chan struct{})
	//var selector string
	if r.opts.useInformer {

		go func() {
			select {
			case <-ctx.Done():
				r.stopCh <- struct{}{}
			case <-waitCh:
			}
		}()

		if synced := cache.WaitForCacheSync(r.stopCh, r.podInformer.Informer().HasSynced); !synced {
			return nil, errors.New("pod cache sync failed")

		}
		waitCh <- struct{}{}

		/*
			s, err := labels.Parse(selector)
			if err != nil {
				return nil, err
			}
		*/

		podLister := r.podInformer.Lister()
		es, err := podLister.(v1.PodLister).Pods(r.opts.namespace).List(selector)
		if err != nil {
			return nil, err
		}

		for _, e := range es {
			pods = append(pods, *e)
		}

	} else {
		options := metav1.ListOptions{LabelSelector: selector.String()}
		podList, err := r.kubeClient.CoreV1().Pods(r.opts.namespace).List(ctx, options)
		if err != nil {
			return nil, err
		}

		pods = podList.Items
	}

	return pods, nil
}

func (r *K8sDNSResolver) LookupHostIPv4(ctx context.Context, svc string) ([]string, error) {
	var addrs []string

	pods, err := r.Pods(ctx, svc)
	if err != nil {
		return nil, err
	}

	for _, pod := range pods {
		if pod.Status.Phase != corev1.PodRunning {
			continue
		}

		nodeUnit, ok := pod.Spec.NodeSelector[r.opts.nodeGroup]
		if !ok {
			continue
		}
		if nodeUnit != r.opts.nodeUnit {
			continue
		}

		if net_.IsIPv4String(pod.Status.PodIP) {
			addrs = append(addrs, pod.Status.PodIP)
		}
	}

	return addrs, nil
}
