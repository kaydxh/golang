package tutorial

import (
	"fmt"
	"testing"
	"time"
)

// A ServerOption sets options.
type ServerOption interface {
	apply(*Server)
}

// EmptyServerOption does not alter the configuration. It can be embedded
// in another structure to build custom options.
//
// This API is EXPERIMENTAL.
type EmptyServerOption struct{}

func (EmptyServerOption) apply(*Server) {}

// ServerOptionFunc wraps a function that modifies Client into an
// implementation of the ServerOption interface.
type ServerOptionFunc func(*Server)

func (f ServerOptionFunc) apply(do *Server) {
	f(do)
}

// sample code for option, default for nothing to change
func _ServerOptionWithDefault() ServerOption {
	return ServerOptionFunc(func(*Server) {
		// nothing to change
	})
}
func (o *Server) ApplyOptions(options ...ServerOption) *Server {
	for _, opt := range options {
		if opt == nil {
			continue
		}
		opt.apply(o)
	}
	return o
}

func WithProtocol(protocol string) ServerOption {
	return ServerOptionFunc(func(s *Server) {
		s.opts.Protocol = protocol
	})
}

func WithTimeout(timeout time.Duration) ServerOption {
	return ServerOptionFunc(func(s *Server) {
		s.opts.Timeout = timeout
	})
}

func WithMaxConns(maxConns int) ServerOption {
	return ServerOptionFunc(func(s *Server) {
		s.opts.MaxConns = maxConns
	})
}

type Server struct {
	Addr string
	Port int
	opts struct {
		Protocol string
		Timeout  time.Duration
		MaxConns int
	}
}

func NewServer(addr string, port int, opts ...ServerOption) *Server {
	s := &Server{
		Addr: addr,
		Port: port,
	}
	for _, opt := range opts {
		s.ApplyOptions(opt)
	}

	return s
}

func TestNewServer(t *testing.T) {
	testCases := []struct {
		Addr     string
		Port     int
		Protocol string
		Timeout  time.Duration
		MaxConns int
	}{
		{
			Addr:     "127.0.0.1",
			Port:     10000,
			Protocol: "tcp",
			Timeout:  time.Second,
		},
	}

	for i, testCase := range testCases {
		t.Run(fmt.Sprintf("%d", i), func(t *testing.T) {
			s := NewServer(testCase.Addr, testCase.Port, WithTimeout(testCase.Timeout), WithProtocol(testCase.Protocol))
			t.Logf("s: %v", s)
		})
	}
}
