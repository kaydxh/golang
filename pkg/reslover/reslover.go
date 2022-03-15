package reslover

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
	"github.com/serialx/hashring"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

type ResloverQueryMap sync.Map

type ResloverOptions struct {
	ResloverType    Reslover_ResloverType
	LoadBalanceMode Reslover_LoadBalanceMode
}

type ResloverQuery struct {
	Domain   string
	nodes    []string
	hashring *hashring.HashRing
	Opts     ResloverOptions
}

//NewDefaultResloverQuery, dns reslover, bls consist hash
func NewDefaultResloverQuery(domain string) *ResloverQuery {
	rq := &ResloverQuery{
		Domain: domain,
	}
	rq.SetDefault()
	return rq
}

func defaultResloverOptions() ResloverOptions {
	return ResloverOptions{
		ResloverType:    Reslover_reslover_type_dns,
		LoadBalanceMode: Reslover_load_balance_mode_consist,
	}
}

func (r *ResloverQuery) SetDefault() {
	r.Opts = defaultResloverOptions()
}

type ResloverService struct {
	resolverInterval time.Duration
	serviceByName    ResloverQueryMap
	inShutdown       atomic.Bool
	mu               sync.Mutex
	cancel           func()
}

func NewResloverService(resolverInterval time.Duration, services ...ResloverQuery) *ResloverService {
	rs := &ResloverService{
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

func (srv *ResloverService) logger() logrus.FieldLogger {
	return logrus.WithField("module", "ResloverService")
}

func (srv *ResloverService) Run(ctx context.Context) error {
	logger := srv.logger()
	logger.Infoln("ResloverService Run")
	if srv.inShutdown.Load() {
		logger.Infoln("ResloverService Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(srv.Serve(ctx))
	}()
	return nil

}

func (srv *ResloverService) Serve(ctx context.Context) error {
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
		logger.Infof("querying services")
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

func (srv *ResloverService) Shutdown() {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.cancel != nil {
		srv.cancel()
	}
}

func (srv *ResloverService) QueryServices() (err error) {
	var (
		errs []error
	)

	logger := srv.logger()
	logger.Info("QueryServices.....")
	srv.serviceByName.Range(func(name string, service ResloverQuery) bool {
		if service.Opts.ResloverType == Reslover_reslover_type_dns {
			service.nodes, err = net_.LookupHostIPv4(service.Domain)
			if err != nil {
				errs = append(errs, fmt.Errorf("failed to query service: %v, err: %v", service.Domain, err))
				return true
			}

			if service.Opts.LoadBalanceMode == Reslover_load_balance_mode_consist {
				service.hashring = hashring.New(service.nodes)
			}
			srv.serviceByName.Store(name, service)
			return true
		}

		return true
	})

	return errors_.NewAggregate(errs)
}

func (srv *ResloverService) PickNode(name string, consistKey string) (node string, has bool) {
	service, has := srv.serviceByName.Load(name)
	if !has {
		return "", false
	}

	if service.Opts.LoadBalanceMode == Reslover_load_balance_mode_consist {
		if service.hashring == nil {
			srv.QueryServices()
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

func (srv *ResloverService) AddService(service ResloverQuery) error {
	_, loaded := srv.serviceByName.LoadOrStore(service.Domain, service)
	if loaded {
		return fmt.Errorf("service entry already installed")
	}
	return nil
}
