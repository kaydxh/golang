package date_test

import (
	"testing"
	"time"

	"context"

	http_ "github.com/kaydxh/golang/go/net/http"
	date_ "github.com/kaydxh/golang/go/net/http/example"
)

func TestNowPbJson(t *testing.T) {

	ctx := context.Background()
	client, err := http_.NewClient()
	if err != nil {
		t.Errorf("failed to new http client, err: %v", err)
	}
	factory, err := http_.NewFactory[date_.NowRequest, date_.NowResponse](http_.FactoryConfig{
		Addr:    "http://localhost:10001/Now",
		Timeout: 5 * time.Second,
		Client:  client,
	},
	)
	if err != nil {
		t.Errorf("failed to new factory, err: %v", err)
	}
	repository, err := factory.NewClient(ctx)
	if err != nil {
		t.Errorf("failed to new respository client, err: %v", err)
	}

	nowResponse, err := repository.PostPbJson(ctx, &date_.NowRequest{})
	if err != nil {
		t.Errorf("failed to call Now, err: %v", err)
	}

	t.Logf("resp: %v", nowResponse)

}

func TestNowPb(t *testing.T) {

	ctx := context.Background()
	client, err := http_.NewClient()
	if err != nil {
		t.Errorf("failed to new http client, err: %v", err)
	}
	factory, err := http_.NewFactory[date_.NowRequest, date_.NowResponse](http_.FactoryConfig{
		Addr:    "http://localhost:10001/Now",
		Timeout: 5 * time.Second,
		Client:  client,
	},
	)
	if err != nil {
		t.Errorf("failed to new factory, err: %v", err)
	}
	repository, err := factory.NewClient(ctx)
	if err != nil {
		t.Errorf("failed to new respository client, err: %v", err)
	}

	nowResponse, err := repository.PostPb(ctx, &date_.NowRequest{})
	if err != nil {
		t.Errorf("failed to call Now, err: %v", err)
	}

	t.Logf("resp: %v", nowResponse)

}
