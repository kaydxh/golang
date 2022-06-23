package instance_test

import (
	"context"
	"testing"

	instance_ "github.com/kaydxh/golang/pkg/pool/instance"
)

func TestNewPool(t *testing.T) {
	ctx := context.Background()
	pool := instance_.NewPool()
	err := pool.GlobalInit(ctx)
	if err != nil {
		t.Errorf("global init err: %v", err)
	}

}
