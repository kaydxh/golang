package resolver

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/kaydxh/golang/go/errors"
	errors_ "github.com/kaydxh/golang/go/errors"
	time_ "github.com/kaydxh/golang/go/time"
	dns_ "github.com/kaydxh/golang/pkg/resolver/dns"
	k8sdns_ "github.com/kaydxh/golang/pkg/resolver/dns/k8s-resolver"
	netdns_ "github.com/kaydxh/golang/pkg/resolver/dns/net-resolver"
	"github.com/serialx/hashring"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

type ResolverQueryMap sync.Map

type ResolverOptions struct {
	resolverType    Resolver_ResolverType
	loadBalanceMode Resolver_LoadBalanceMode
	//k8s
	nodeGroup string
	nodeUnit  string
}

type ResolverQuery struct {
	// domain: domain for net resolver, Ex: cube-svc.ns.svc
	// domain: svc name for k8s resolver, Ex: cube
	Domain   string
	nodes    []string
	hashring *hashring.HashRing
	opts     ResolverOptions
	resolver dns_.DNSResolver
}

//NewDefaultResolverQuery, dns reslover, bls consist hash
func NewDefaultResolverQuery(domain string) ResolverQuery {
	rq := ResolverQuery{
		Domain: domain,
	}
	rq.SetDefault()
	return rq
}

func defaultResolverOptions() ResolverOptions {
	return ResolverOptions{
		resolverType:    Resolver_resolver_type_dns,
		loadBalanceMode: Resolver_load_balance_mode_consist,
	}
}

func (r *ResolverQuery) SetDefault() {
	r.opts = defaultResolverOptions()
}

func NewResolverQuery(domain string, opts ...ResolverQueryOption) (ResolverQuery, error) {
	rq := NewDefaultResolverQuery(domain)
	rq.ApplyOptions(opts...)
	err := rq.SetResolver()
	if err != nil {
		return rq, err
	}

	return rq, nil
}

func (r *ResolverQuery) SetResolver() error {
	var err error
	switch r.opts.resolverType {
	case Resolver_resolver_type_k8s:
		r.resolver, err = k8sdns_.NewK8sDNSResolver(
			k8sdns_.WithNodeGroup(r.opts.nodeGroup),
			k8sdns_.WithNodeUnit(r.opts.nodeUnit),
		)
		if err != nil {
			logrus.Errorf("new k8s resolver, err: %v", err)
			return err
		}
	default:
		r.resolver = netdns_.DefaultResolver
	}

	return nil
}

type ResolverService struct {
	resolverInterval time.Duration
	serviceByName    ResolverQueryMap
	inShutdown       atomic.Bool
	mu               sync.Mutex
	cancel           func()
}

func NewDefaultResolverServices(resolverInterval time.Duration, domains ...string) *ResolverService {
	rs := NewResolverService(resolverInterval)
	for _, domain := range domains {
		rq := NewDefaultResolverQuery(domain)
		rs.AddService(rq)
	}
	return rs
}

func NewResolverService(resolverInterval time.Duration, services ...ResolverQuery) *ResolverService {
	rs := &ResolverService{
		resolverInterval: resolverInterval,
	}
	if rs.resolverInterval == 0 {
		rs.resolverInterval = 10 * time.Second
	}

	rs.AddServices(services...)
	return rs
}

func (srv *ResolverService) logger() logrus.FieldLogger {
	return logrus.WithField("module", "ResolverService")
}

func (srv *ResolverService) Run(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("ResolverService Run")
	if srv.inShutdown.Load() {
		logger.Infoln("ResolverService Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(srv.Serve(ctx))
	}()
	return nil

}

func (srv *ResolverService) Serve(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("ServiceResolver Serve")

	if srv.inShutdown.Load() {
		err := fmt.Errorf("server closed")
		logger.WithError(err).Errorf("ServiceResolver Serve canceled")
		return err
	}

	defer srv.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	srv.mu.Lock()
	srv.cancel = cancel
	srv.mu.Unlock()

	time_.UntilWithContxt(ctx, func(ctx context.Context) {
		logger.Debugf("querying services")
		err := srv.QueryServices()
		if err != nil {
			logger.WithError(err).Errorf("query services failed")
			return
		}
		logger.Debugf("query services")
	}, srv.resolverInterval)
	logger.Info("stopped query services")
	return nil
}

func (srv *ResolverService) Shutdown() {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.cancel != nil {
		srv.cancel()
	}
}

func (srv *ResolverService) QueryServices() (err error) {
	var (
		errs []error
	)

	logger := srv.logger()
	srv.serviceByName.Range(func(name string, service ResolverQuery) bool {

		service.nodes, err = service.resolver.LookupHostIPv4(context.Background(), service.Domain)
		if err != nil {
			errs = append(errs, fmt.Errorf("failed to query service: %v, err: %v", service.Domain, err))
			return true
		}
		logger.Debugf("the results of lookup domain[%s] are nodes[%v]", service.Domain, service.nodes)

		if service.opts.loadBalanceMode == Resolver_load_balance_mode_consist {
			service.hashring = hashring.New(service.nodes)
		}
		srv.serviceByName.Store(name, service)
		return true

	})

	return errors_.NewAggregate(errs)
}

func (srv *ResolverService) PickNode(name string, consistKey string) (node string, has bool) {
	service, has := srv.serviceByName.Load(name)
	if !has {
		return "", false
	}

	if service.opts.loadBalanceMode == Resolver_load_balance_mode_consist {
		if service.hashring == nil {
			err := srv.QueryServices()
			if err != nil {
				return "", false
			}
			service, has = srv.serviceByName.Load(name)
			if !has {
				return "", false
			}
		}
		return service.hashring.GetNode(consistKey)
	}

	if len(service.nodes) == 0 {
		return "", false
	}
	s := service.nodes[rand.Intn(len(service.nodes))]
	return s, true
}

// AddService, if exist, it will update
func (srv *ResolverService) AddService(service ResolverQuery) {
	srv.serviceByName.Store(service.Domain, service)
}

// AddServices, if exist, it will update
func (srv *ResolverService) AddServices(services ...ResolverQuery) {
	for _, service := range services {
		srv.AddService(service)
	}
}
