package webserver

import (
	"context"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kaydxh/golang/pkg/consul"
	"github.com/searKing/golang/third_party/github.com/grpc-ecosystem/grpc-gateway/grpc"
	"github.com/searKing/sole/pkg/webserver/healthz"
)

type WebHandler interface {
	SetRoutes(ginRouter gin.IRouter, grpcRouter *grpc.Gateway)
}

type GenericWebServer struct {
	// Server Register. The backend is started after the server starts listening.
	ServiceRegistryBackend *consul.ServiceRegistryServer

	ginBackend  *gin.Engine
	grpcBackend *grpc.Gateway

	// PostStartHooks are each called after the server has started listening, in a separate go func for each
	// with no guarantee of ordering between them.  The map key is a name used for error reporting.
	// It may kill the process with a panic if it wishes to by returning an error.
	postStartHookLock sync.Mutex
	//	postStartHooks       map[string]postStartHookEntry
	postStartHooksCalled bool

	preShutdownHookLock sync.Mutex
	//	preShutdownHooks       map[string]preShutdownHookEntry
	preShutdownHooksCalled bool

	// healthz checks
	healthzLock            sync.Mutex
	healthzChecks          []healthz.HealthCheck
	healthzChecksInstalled bool
	// livez checks
	livezLock            sync.Mutex
	livezChecks          []healthz.HealthCheck
	livezChecksInstalled bool
	// readyz checks
	readyzLock            sync.Mutex
	readyzChecks          []healthz.HealthCheck
	readyzChecksInstalled bool
	livezGracePeriod      time.Duration

	// the readiness stop channel is used to signal that the apiserver has initiated a shutdown sequence, this
	// will cause readyz to return unhealthy.
	readinessStopCh chan struct{}

	// ShutdownDelayDuration allows to block shutdown for some time, e.g. until endpoints pointing to this API server
	// have converged on all node. During this time, the API server keeps serving, /healthz will return 200,
	// but /readyz will return failure.
	ShutdownDelayDuration time.Duration

	// The limit on the request body size that would be accepted and decoded in a write request.
	// 0 means no limit.
	maxRequestBodyBytes int64

	// WebServerID is the ID of this Web server
	WebServerID string
}

// preparedGenericWebServer is a private wrapper that enforces a call of PrepareRun() before Run can be invoked.
type preparedGenericWebServer struct {
	*GenericWebServer
}

func (s *GenericWebServer) PrepareRun() (preparedGenericWebServer, error) {
	return preparedGenericWebServer{s}, nil
}

func (s preparedGenericWebServer) NonBlockingRun(ctx context.Context) (context.Context, error) {

	ctx, cancel := context.WithCancel(ctx)
	go func() {
		defer cancel()
		err := s.grpcBackend.ListenAndServe()
		if err != nil {
			return
		}
	}()

	return ctx, nil
}

func (s preparedGenericWebServer) Run(ctx context.Context) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	ctx, err := s.NonBlockingRun(ctx)
	if err != nil {
		return err
	}

	<-ctx.Done()

	return nil
}

func (s *GenericWebServer) InstallWebHandlers(handlers ...WebHandler) {
	for _, h := range handlers {
		if h == nil {
			continue
		}
		h.SetRoutes(s.ginBackend, s.grpcBackend)
	}
}
