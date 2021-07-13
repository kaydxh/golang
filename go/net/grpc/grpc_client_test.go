package grpc_test

import (
	"testing"
	"time"

	grpc_ "github.com/kaydxh/golang/go/net/grpc"
)

func TestGetGrpcClientConn(t *testing.T) {

	var (
		serverAddress     = "127.0.0.1:8001"
		connectionTimeout = 5 * time.Second
	)

	conn, err := grpc_.GetGrpcClientConn(serverAddress, connectionTimeout)
	if err != nil {
		t.Fatalf("failed to get local addrs, err: %v", err)
		return
	}
	defer conn.Close()
}
