package consul

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/hashicorp/consul/api"
	"github.com/kaydxh/golang/go/errors"
	"github.com/sirupsen/logrus"
	"go.uber.org/atomic"
)

type ServiceRegistryServer struct {
	ConsulAddress string
	ServiceId     string
	ServiceName   string
	Tag           []string
	Ip            string
	Port          int
	TTL           time.Duration
	CheckInterval time.Duration

	inShutdown atomic.Bool // true when when server is in shutdown

	mu     sync.Mutex
	cancel func()
}

func NewServiceRegistry(
	consulAddr string,
	serviceName string,
	serviceAddress string,
) (*ServiceRegistryServer, error) {
	logger := logrus.WithField("module", "service_registry").
		WithField("service_name", serviceName).
		WithField("addr", serviceAddress)

	host, port, err := net.SplitHostPort(serviceAddress)
	if err != nil {
		logger.WithError(err).Errorln("malformed service serviceAddress")
		return nil, fmt.Errorf("malformed service serviceAddress: %w", err)
	}

	nport, err := strconv.Atoi(port)
	if err != nil {
		logger.WithField("port", port).
			WithError(err).
			Errorln("malformed service port, must be a number")
		return nil, fmt.Errorf("malformed service port: %w", err)
	}
	serviceId := fmt.Sprintf("%v-%v-%v", serviceName, host, port)
	s := &ServiceRegistryServer{
		ConsulAddress: consulAddr,
		ServiceId:     serviceId,
		ServiceName:   serviceName,
		Tag:           []string{},
		Port:          nport,
		Ip:            host,
		TTL:           300 * time.Second,
		CheckInterval: 10 * time.Second,
	}
	return s, nil
}

// Run will initialize the backend. It must not block, but may run go routines in the background.
func (srv *ServiceRegistryServer) Run(ctx context.Context) error {
	logger := srv.logger().
		WithField("service_name", srv.ServiceName).
		WithField("service_id", srv.ServiceId)
	logger.Infoln("ConsulRegistry Run")
	if srv.inShutdown.Load() {
		logger.Infoln("ConsulRegistry Shutdown")
		return fmt.Errorf("server closed")
	}
	go func() {
		errors.HandleError(srv.Serve(ctx))
	}()
	return nil
}

func (srv *ServiceRegistryServer) Serve(ctx context.Context) error {
	logger := srv.logger().
		WithField("service_name", srv.ServiceName).
		WithField("service_id", srv.ServiceId)
	logger.Infoln("ConsulRegistry Serve")

	if srv.inShutdown.Load() {
		logger.Infoln("ConsulRegistry Shutdown")
		return fmt.Errorf("server closed")
	}

	defer srv.inShutdown.Store(true)
	ctx, cancel := context.WithCancel(ctx)
	srv.mu.Lock()
	srv.cancel = cancel
	srv.mu.Unlock()

	t := time.NewTicker(time.Second * 30)
	defer t.Stop()
	for {
		select {
		case <-t.C:
			logger.Infoln("registering service to consul")
			err := srv.Register()
			if err != nil {
				logger.WithError(err).Errorln("register service failed")
				continue
			}
			srv.logger().Info("register service by consul")

		case <-ctx.Done():
			logger.Infoln("unregistering service to consul")
			err := srv.UnRegister()
			if err != nil {
				logger.WithError(err).Errorln("unregister service failed")
				return err
			}
			srv.logger().Info("unregister service by consul")
			return nil
		}
	}
}

func (srv *ServiceRegistryServer) Shutdown() {
	srv.inShutdown.Store(true)
	srv.mu.Lock()
	defer srv.mu.Unlock()
	if srv.cancel != nil {
		srv.cancel()
	}
}

func (srv *ServiceRegistryServer) Register() error {
	config := api.DefaultConfig()
	config.Address = srv.ConsulAddress
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}
	agent := client.Agent()

	checkUrl := fmt.Sprintf("http://%v:%v/api/%v/v1/health", srv.Ip, srv.Port, srv.ServiceName)

	reg := &api.AgentServiceRegistration{
		ID:      srv.ServiceId,
		Name:    srv.ServiceName,
		Tags:    srv.Tag,
		Port:    srv.Port,
		Address: srv.Ip,
		Check: &api.AgentServiceCheck{
			Interval:                       srv.CheckInterval.String(),
			HTTP:                           checkUrl,
			DeregisterCriticalServiceAfter: srv.TTL.String(),
		},
	}

	return agent.ServiceRegister(reg)
}

func (srv *ServiceRegistryServer) UnRegister() error {
	config := api.DefaultConfig()
	config.Address = srv.ConsulAddress
	client, err := api.NewClient(config)
	if err != nil {
		return err
	}

	return client.Agent().ServiceDeregister(srv.ServiceId)
}

func (srv *ServiceRegistryServer) logger() logrus.FieldLogger {
	return logrus.
		WithField("module", "service_registry").
		WithField("consul", srv.ConsulAddress).
		WithField("service_name", srv.ServiceName).
		WithField("service_id", srv.ServiceId).
		WithField("ip", srv.Ip).WithField("port", srv.Port)
}
