package resolver

import (
	"context"
	"fmt"
	"math/rand"
	"sync"
	"time"

	"github.com/kaydxh/golang/go/errors"
	errors_ "github.com/kaydxh/golang/go/errors"
	net_ "github.com/kaydxh/golang/go/net"
	time_ "github.com/kaydxh/golang/go/time"
	dns_ "github.com/kaydxh/golang/pkg/resolver/dns"
	"github.com/serialx/hashring"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

type ResolverQueryMap sync.Map

type ResolverOptions struct {
	ResolverType    Resolver_ResolverType
	LoadBalanceMode Resolver_LoadBalanceMode
}

type ResolverQuery struct {
	Domain   string
	nodes    []string
	hashring *hashring.HashRing
	Opts     ResolverOptions
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
		ResolverType:    Resolver_resolver_type_dns,
		LoadBalanceMode: Resolver_load_balance_mode_consist,
	}
}

func (r *ResolverQuery) SetDefault() {
	r.Opts = defaultResolverOptions()
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

	for _, s := range services {
		_ = rs.AddService(s)
	}
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
		if service.Opts.ResolverType == Resolver_resolver_type_dns {
			service.nodes, err = net_.LookupHostIPv4(service.Domain)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to query service: %v, err: %v", service.Domain, err))
				return true
			}
			logger.Debugf("the results of lookup domain[%s] are nodes[%v]", service.Domain, service.nodes)

			if service.Opts.LoadBalanceMode == Resolver_load_balance_mode_consist {
				service.hashring = hashring.New(service.nodes)
			}
			srv.serviceByName.Store(name, service)
			return true
		}

		return true
	})

	return errors_.NewAggregate(errs)
}

func (srv *ResolverService) PickNode(name string, consistKey string) (node string, has bool) {
	service, has := srv.serviceByName.Load(name)
	if !has {
		return "", false
	}

	if service.Opts.LoadBalanceMode == Resolver_load_balance_mode_consist {
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

func (srv *ResolverService) AddService(service ResolverQuery) error {
	_, loaded := srv.serviceByName.LoadOrStore(service.Domain, service)
	if loaded {
		return fmt.Errorf("service[%v] entry already installed", service)
	}
	return nil
}

func (srv *ResolverService) AddServices(services ...ResolverQuery) error {
	var errs []error
	for _, service := range services {
		err := srv.AddService(service)
		if err != nil {
			errs = append(errs, err)
		}

	}

	return errors_.NewAggregate(errs)
}
