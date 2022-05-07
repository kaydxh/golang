package k8sdns

import (
	"context"
	"errors"
	"fmt"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/client-go/informers"
	Informerv1 "k8s.io/client-go/informers/core/v1"
	"k8s.io/client-go/kubernetes"
	v1 "k8s.io/client-go/listers/core/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	defaultInClusterConfig = true
	defaultUseInformer     = true
)

type ResolverConfig struct {
	groupNode          string
	nodeUnit           string
	namespace          string
	useInClusterConfig bool
	useInformer        bool
	kubeConfig         string
}

type K8sDNSResolver struct {
	opts       ResolverConfig
	stopCh     chan struct{}
	kubeClient kubernetes.Interface
	//	client      *kubernetes.Clientset
	factory     informers.SharedInformerFactory
	podInformer Informerv1.PodInformer
}

func NewK8sDNSResolver(opts ...K8sDNSResolverOption) (*K8sDNSResolver, error) {

	r := &K8sDNSResolver{
		stopCh: make(chan struct{}),
	}
	r.opts.useInClusterConfig = defaultInClusterConfig
	r.opts.useInformer = defaultUseInformer
	r.opts.namespace = corev1.NamespaceAll
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

	return r, nil
}

func (r *K8sDNSResolver) Pods() ([]corev1.Pod, error) {

	var pods []corev1.Pod
	var selector string
	if r.opts.useInformer {
		go r.podInformer.Informer().Run(r.stopCh)
		if synced := cache.WaitForCacheSync(r.stopCh, r.podInformer.Informer().HasSynced); !synced {
			return nil, errors.New("pod cache sync failed")
		}

		s, err := labels.Parse(selector)
		if err != nil {
			return nil, err
		}

		podLister := r.podInformer.Lister()
		es, err := podLister.(v1.PodLister).Pods(r.opts.namespace).List(s)
		if err != nil {
			return nil, err
		}

		for _, e := range es {
			pods = append(pods, *e)
		}

	} else {
		options := metav1.ListOptions{LabelSelector: selector}

		podList, err := r.kubeClient.CoreV1().Pods(r.opts.namespace).List(context.Background(), options)
		if err != nil {
			return nil, err
		}

		pods = podList.Items
	}

	return pods, nil
}
